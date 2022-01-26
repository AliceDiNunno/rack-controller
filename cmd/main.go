package main

import (
	"github.com/AliceDiNunno/rack-controller/src/adapters/cluster/kubernetes"
	"github.com/AliceDiNunno/rack-controller/src/adapters/eventDispatcher/dispatcher"
	"github.com/AliceDiNunno/rack-controller/src/adapters/persistence/mongodb"
	"github.com/AliceDiNunno/rack-controller/src/adapters/persistence/postgres"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest"
	"github.com/AliceDiNunno/rack-controller/src/config"
	"github.com/AliceDiNunno/rack-controller/src/core/usecases"
	"github.com/AliceDiNunno/rack-controller/src/core/usecases/events"
	log "github.com/sirupsen/logrus"
	stdlog "log"
)

func main() {
	//Disabling timestamp for log output since we rely on an external tool
	stdlog.SetFlags(0)

	//Loading Configuration
	config.LoadEnv()
	globalConfiguration := config.LoadGlobalConfiguration()
	ginConfiguration := config.LoadGinConfiguration()
	dbConfig := config.LoadGormConfiguration()
	initialUserConfiguration := config.LoadInitialUserConfiguration()
	clusterConfig := config.LoadClusterConfig()
	mongoConfig := config.LoadMongodbConfiguration()

	//Loading the kubernetes client
	kubernetesInstance, err := kubernetes.LoadInstance(clusterConfig)

	if err != nil {
		log.Fatalln(err)
	}

	mongo := mongodb.StartMongodbDatabase(mongoConfig)
	var logCollection usecases.LogCollection

	logCollection = mongodb.NewLogCollectionRepo(mongo)

	//Loading the database
	db := postgres.StartGormDatabase(dbConfig)
	err = db.AutoMigrate(
		//Migrating user tables
		&postgres.User{}, &postgres.JwtSignature{}, &postgres.UserToken{},
		//Migrating kubernetes-related tables
		&postgres.Environment{}, &postgres.Project{}, &postgres.Config{}, &postgres.Service{})
	if err != nil {
		log.Fatalln(err)
	}

	userRepo := postgres.NewUserRepo(db)
	tokenRepo := postgres.NewUserTokenRepo(db)
	jwtSignatureRepo := postgres.NewJwtSignatureRepo(db)
	projectRepo := postgres.NewProjectRepo(db)
	environmentRepo := postgres.NewEnvironmentRepo(db)
	serviceRepo := postgres.NewServiceRepo(db)
	configRepo := postgres.NewConfigRepo(db)

	//Loading the event dispatcher
	var eventDispatcher = dispatcher.NewDispatcher()
	usecasesHandler := usecases.NewInteractor(userRepo, tokenRepo, jwtSignatureRepo,
		projectRepo, environmentRepo, serviceRepo, configRepo,
		logCollection,
		kubernetesInstance, eventDispatcher)

	if initialUserConfiguration != nil {
		err := usecasesHandler.CreateInitialUser(initialUserConfiguration)
		if err != nil {
			log.Warning(err.Err.Error())
		}
	}

	//Loading the rest api
	restServer := rest.NewServer(globalConfiguration, ginConfiguration)
	routesHandler := rest.NewRouter(usecasesHandler)

	eventHandler := events.NewEventHandler(usecasesHandler, eventDispatcher, kubernetesInstance)
	eventHandler.SetEvents()

	rest.SetRoutes(restServer, routesHandler)

	//Starting the server
	restServer.Start()
}
