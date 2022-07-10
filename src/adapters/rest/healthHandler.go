package rest

import (
	"github.com/AliceDiNunno/rack-controller/src/adapters/ip"
	"github.com/gin-gonic/gin"
	"time"
)

type HealthData struct {
	Running bool
	Date    time.Time
	IP      string
}

func (rH RoutesHandler) GetHealthHandler(c *gin.Context) {
	ipCollector := ip.NewIPCollector()

	ipData, err := ipCollector.GetLocalIP()

	localIP := ""

	if err != nil {
		localIP = "Unable to get IP"
	} else {
		localIP = ipData.Query
	}

	rH.handleSuccess(c, HealthData{
		Running: true,
		Date:    time.Now(),
		IP:      localIP,
	})
}
