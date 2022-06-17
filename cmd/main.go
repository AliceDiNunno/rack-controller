package main

import (
	"crypto/tls"
	glc "github.com/AliceDiNunno/go-logger-client"
	"github.com/AliceDiNunno/rack-controller/src/adapters/cluster/kubernetes"
	eventAdapter "github.com/AliceDiNunno/rack-controller/src/adapters/event"
	"github.com/AliceDiNunno/rack-controller/src/adapters/eventDispatcher/dispatcher"
	"github.com/AliceDiNunno/rack-controller/src/adapters/ip"
	"github.com/AliceDiNunno/rack-controller/src/adapters/ovh"
	"github.com/AliceDiNunno/rack-controller/src/adapters/persistence/postgres"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest"
	"github.com/AliceDiNunno/rack-controller/src/config"
	"github.com/AliceDiNunno/rack-controller/src/core/usecases"
	log "github.com/sirupsen/logrus"
	stdlog "log"
	"net/http"
)

func main() {
	//Disabling timestamp for log output since we rely on an external tool
	stdlog.SetFlags(0)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	//Loading Configuration
	config.LoadEnv()
	globalConfiguration := config.LoadGlobalConfiguration()
	ginConfiguration := config.LoadGinConfiguration()
	dbConfig := config.LoadGormConfiguration()
	initialUserConfiguration := config.LoadInitialUserConfiguration()
	clusterConfig := config.LoadClusterConfig()
	eventConfig := config.LoadEventConfiguration()
	ovhConfig := config.LoadOvhConfiguration()

	//Loading the kubernetes client
	kubernetesInstance, err := kubernetes.LoadInstance(clusterConfig)

	//Loading the IP Collector
	ipCollector := ip.NewIPCollector()

	//Loading the OVH Client
	ovhClient := ovh.NewOVHClient(ovhConfig)

	if err != nil {
		log.Fatalln(err)
	}

	var eventCollection usecases.EventRepository

	//Loading the database
	db := postgres.StartGormDatabase(dbConfig)
	err = db.AutoMigrate(
		//Migrating event tables
		&postgres.Event{}, &postgres.EventOccurrence{},
		&postgres.Traceback{},
		&postgres.TracebackEntry{},
		//Migrating user tables
		&postgres.User{}, &postgres.JwtSignature{}, &postgres.UserToken{},
		//Migrating kubernetes-related tables
		&postgres.Environment{}, &postgres.Project{}, &postgres.Addon{}, &postgres.Config{}, &postgres.Service{})
	if err != nil {
		log.Fatalln(err)
	}

	userRepo := postgres.NewUserRepo(db)
	tokenRepo := postgres.NewUserTokenRepo(db)
	jwtSignatureRepo := postgres.NewJwtSignatureRepo(db)
	projectRepo := postgres.NewProjectRepo(db)
	environmentRepo := postgres.NewEnvironmentRepo(db)
	serviceRepo := postgres.NewServiceRepo(db)
	eventCollection = postgres.NewEventsRepo(db)
	configRepo := postgres.NewConfigRepo(db)
	addonRepo := postgres.NewAddonRepo(db)

	//Loading the event dispatcher
	var eventDispatcher = dispatcher.NewDispatcher()
	usecasesHandler := usecases.NewInteractor(userRepo, tokenRepo, jwtSignatureRepo,
		projectRepo, environmentRepo, serviceRepo, configRepo, addonRepo,
		eventCollection,
		kubernetesInstance, eventDispatcher, ipCollector, ovhClient)

	internalEventTransporter := eventAdapter.NewEventTransporter(usecasesHandler)
	receiver := glc.NewInternalTransporter(internalEventTransporter, eventConfig)
	glc.SetupHook(eventConfig, receiver)

	if initialUserConfiguration != nil {
		err := usecasesHandler.CreateInitialUser(initialUserConfiguration)
		if err != nil {
			log.Println(err.Err.Error())
		}
	}

	//Loading the rest api
	restServer := rest.NewServer(globalConfiguration, ginConfiguration)
	routesHandler := rest.NewRouter(usecasesHandler)

	rest.SetRoutes(restServer, routesHandler)

	//Starting the server
	restServer.Start()
}
