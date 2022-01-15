package rest

import (
	"fmt"
	"github.com/AliceDiNunno/rack-controller/src/config"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type GinServer struct {
	Config config.GinConfig
	Router *gin.Engine
}

func NewServer(globalConfig config.GlobalConfig, config config.GinConfig) GinServer {
	if globalConfig.CurrentEnvironment == "prod" || globalConfig.CurrentEnvironment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	server := GinServer{
		Config: config,
		Router: gin.New(),
	}

	return server
}

func (server GinServer) Start() {
	if server.Config.TlsEnabled {
		httpsServer := http.Server{
			Addr:    fmt.Sprintf("%s:%d", server.Config.ListenAddress, server.Config.Port),
			Handler: server.Router,
		}

		err := httpsServer.ListenAndServeTLS("./certificates/server.cert", "./certificates/server.key")
		if err != nil {
			log.Fatalln("Could not start server", err)
		}
		log.Printf("Started server at %s on port %d\n", server.Config.ListenAddress, server.Config.Port)
	} else {
		if err := server.Router.Run(fmt.Sprintf("%s:%d", server.Config.ListenAddress, server.Config.Port)); err != nil {
			println("Couldn't start router")
		}
	}
}
