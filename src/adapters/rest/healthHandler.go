package rest

import (
	"github.com/gin-gonic/gin"
	"time"
)

type HealthData struct {
	Running bool
	Date    time.Time
}

func (rH RoutesHandler) GetHealthHandler(c *gin.Context) {
	c.JSON(200, success(HealthData{
		Running: true,
		Date:    time.Now(),
	}))
}
