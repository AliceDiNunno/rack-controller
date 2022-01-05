package main

import (
	"fmt"
	"github.com/AliceDiNunno/rack-controller/src/adapters/cluster/kubernetes"
	"github.com/AliceDiNunno/rack-controller/src/adapters/event"
	"github.com/AliceDiNunno/rack-controller/src/adapters/persistence/postgres"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest"
	"github.com/AliceDiNunno/rack-controller/src/config"
	usecases "github.com/AliceDiNunno/rack-controller/src/core/usecases"
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

	var userRepo usecases.UserRepo
	var tokenRepo usecases.UserTokenRepo
	var jwtSignatureRepo usecases.JwtSignatureRepo
	var projectRepo usecases.ProjectRepo
	var environmentRepo usecases.EnvironmentRepo
	var serviceRepo usecases.ServiceRepo

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
			//Migrating cluster tables
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

	usecasesHandler := usecases.NewInteractor(userRepo, tokenRepo, jwtSignatureRepo,
		projectRepo, environmentRepo, serviceRepo,
		kubernetesInstance)

	if initialUserConfiguration != nil {
		err := usecasesHandler.CreateInitialUser(initialUserConfiguration)
		if err != nil {
			log.Warning(err.Err.Error())
		}
	}

	eventManager := event.NewDispatcher()

	_ = eventManager

	restServer := rest.NewServer(ginConfiguration)
	routesHandler := rest.NewRouter(usecasesHandler)

	rest.SetRoutes(restServer, routesHandler)

	restServer.Start()
}
