package rest

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	event "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func (rH RoutesHandler) handleError(c *gin.Context, err *e.Error) {
	var depth = 2
	errName := getFunctionName(depth) + ": " + err.Err.Error()
	code := codeForError(err.Err)

	fields := event.Fields{
		"code":   code,
		"ip":     c.ClientIP(),
		"path":   c.Request.RequestURI,
		"module": c.FullPath(),
		"err":    &err,
	}

	spew.Dump(fields)

	event.WithFields(fields).Error(errName)
	hostname, _ := os.Hostname()
	c.AbortWithStatusJSON(code, Status{
		Success: false,
		Message: err.Err.Error(),
		Data:    nil,
		Host:    hostname,
	})
}

func (rH RoutesHandler) handleSuccess(c *gin.Context, data interface{}) {
	hostname, _ := os.Hostname()

	println(c.Request.RequestURI, spew.Sdump(data))

	c.JSON(http.StatusOK, Status{
		Success: true,
		Data:    data,
		Host:    hostname,
	})
}
