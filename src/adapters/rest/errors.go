package rest

import (
	"errors"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var (
	ErrFormValidation = errors.New("failed to validate form")
	ErrNotFound       = errors.New("endpoint not found")

	ErrAuthorizationHeaderMissing = errors.New("authorization header missing")
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
)

func codeForError(err error) int {
	switch err {
	case ErrFormValidation:
		return http.StatusBadRequest
	case ErrNotFound:
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func (rH RoutesHandler) handleError(c *gin.Context, err *e.Error) {
	code := codeForError(err.Err)

	fields := log.Fields{
		"code": code,
		"ip":   c.ClientIP(),
		"path": c.Request.RequestURI,
	}

	log.WithFields(fields).Error(err.Err.Error())
	hostname, _ := os.Hostname()
	c.AbortWithStatusJSON(code, Status{
		Success: false,
		Message: err.Err.Error(),
		Data:    nil,
		Host:    hostname,
	})
}

func (rH RoutesHandler) endpointNotFound(c *gin.Context) {
	rH.handleError(c, e.Wrap(ErrNotFound))
}
