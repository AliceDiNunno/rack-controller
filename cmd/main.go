package main

import (
	"fmt"
	"github.com/AliceDiNunno/rack-controller/src/adapters/cluster/kubernetes"
	"github.com/AliceDiNunno/rack-controller/src/adapters/eventDispatcher/dispatcher"
	"github.com/AliceDiNunno/rack-controller/src/adapters/persistence/postgres"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest"
	"github.com/AliceDiNunno/rack-controller/src/config"
	"github.com/AliceDiNunno/rack-controller/src/core/usecases"
	"github.com/AliceDiNunno/rack-controller/src/core/usecases/events"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	stdlog "log"
)

func main() {
	stdlog.SetFlags(0)

	config.LoadEnv()

	ginConfiguration := config.LoadGinConfiguration()
	dbConfig := config.LoadGormConfiguration()
	initialUserConfiguration := config.LoadInitialUserConfiguration()
	clusterConfig := config.LoadClusterConfig()

	var userRepo usecases.UserRepository
	var tokenRepo usecases.UserTokenRepository
	var jwtSignatureRepo usecases.JwtSignatureRepository
	var projectRepo usecases.ProjectRepository
	var environmentRepo usecases.EnvironmentRepository
	var serviceRepo usecases.ServiceRepository

	kubernetesInstance, err := kubernetes.LoadInstance(clusterConfig)

	if err != nil {
		log.Fatalln(err)
	}

	var db *gorm.DB
	if dbConfig.Engine == "POSTGRES" {
		db = postgres.StartGormDatabase(dbConfig)
		err := db.AutoMigrate(
			//Migrating user tables
			&postgres.User{}, &postgres.JwtSignature{}, &postgres.UserToken{},
			//Migrating events tables
			&postgres.Project{}, &postgres.Environment{}, &postgres.Service{})
		if err != nil {
			log.Fatalln(err)
		}

		userRepo = postgres.NewUserRepo(db)
		tokenRepo = postgres.NewUserTokenRepo(db)
		jwtSignatureRepo = postgres.NewJwtSignatureRepo(db)
		projectRepo = postgres.NewProjectRepo(db)
		environmentRepo = postgres.NewEnvironmentRepo(db)
		serviceRepo = postgres.NewServiceRepo(db)
	} else {
		log.Fatalln(fmt.Sprintf("Database engine \"%s\" not supported", dbConfig.Engine))
	}

	var eventDispatcher = dispatcher.NewDispatcher()
	usecasesHandler := usecases.NewInteractor(userRepo, tokenRepo, jwtSignatureRepo,
		projectRepo, environmentRepo, serviceRepo,
		kubernetesInstance, eventDispatcher)

	if initialUserConfiguration != nil {
		err := usecasesHandler.CreateInitialUser(initialUserConfiguration)
		if err != nil {
			log.Warning(err.Err.Error())
		}
	}

	restServer := rest.NewServer(ginConfiguration)
	routesHandler := rest.NewRouter(usecasesHandler)

	eventHandler := events.NewEventHandler(usecasesHandler, eventDispatcher)
	eventHandler.SetEvents()

	rest.SetRoutes(restServer, routesHandler)

	restServer.Start()
}
