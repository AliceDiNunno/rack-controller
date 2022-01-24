package rest

import (
	"errors"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"runtime"
	"strings"
)

var (
	ErrFormValidation = errors.New("failed to validate form")
	ErrUrlValidation  = errors.New("failed to validate url")
	ErrNotFound       = errors.New("endpoint not found")
	ErrUnauthorized   = errors.New("unauthorized")

	ErrAuthorizationHeaderMissing = errors.New("authorization header missing")
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
)

func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}

func getFunctionName(depth int) string {
	function := getFrame(depth).Function
	functionSplitted := strings.Split(function, "/")
	functionName := functionSplitted[len(functionSplitted)-1:][0]

	specifiedFunctionActionSplitted := strings.Split(functionName, ".")
	specifiedFunctionName := specifiedFunctionActionSplitted[2]

	return specifiedFunctionName
}

func codeForError(err error) int {
	switch err {
	case ErrFormValidation, ErrUrlValidation:
		return http.StatusBadRequest
	case ErrNotFound:
		return http.StatusNotFound
	case ErrUnauthorized:
		return http.StatusUnauthorized
	}
	return http.StatusInternalServerError
}

func (rH RoutesHandler) handleError(c *gin.Context, err *e.Error) {
	var depth = 2
	errName := getFunctionName(depth) + ": " + err.Err.Error()
	code := codeForError(err.Err)

	fields := log.Fields{
		"code": code,
		"ip":   c.ClientIP(),
		"path": c.Request.RequestURI,
		"err":  &err,
	}

	log.WithFields(fields).Error(errName)
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
