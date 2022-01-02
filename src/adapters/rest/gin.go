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

func NewServer(config config.GinConfig) GinServer {
	server := GinServer{
		Config: config,
		Router: gin.New(),
	}

	return server
}

func (server GinServer) Start() {
	if server.Config.Tls {
		httpsServer := http.Server{
			Addr:    fmt.Sprintf("%s:%d", server.Config.Host, server.Config.Port),
			Handler: server.Router,
		}

		err := httpsServer.ListenAndServeTLS("./certificates/server.cert", "./certificates/server.key")
		if err != nil {
			log.Fatalln("Could not start server", err)
		}
		log.Printf("Started server at %s on port %d\n", server.Config.Host, server.Config.Port)
	} else {
		if err := server.Router.Run(fmt.Sprintf("%s:%d", server.Config.Host, server.Config.Port)); err != nil {
			println("Couldn't start router")
		}
	}
}
