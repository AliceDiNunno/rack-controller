package rest

import (
	"github.com/AliceDiNunno/rack-controller/src/core/usecases"
	"github.com/gin-gonic/gin"
)

type RoutesHandler struct {
	usecases usecases.Usecases
}

func (rH RoutesHandler) getProjectsHandler(context *gin.Context) {

}

func (rH RoutesHandler) createProjectHandler(context *gin.Context) {

}

func (rH RoutesHandler) getProjectHandler(context *gin.Context) {

}

func (rH RoutesHandler) updateProjectHandler(context *gin.Context) {

}

func (rH RoutesHandler) deleteProjectHandler(context *gin.Context) {

}

func (rH RoutesHandler) getEnvironmentsHandler(context *gin.Context) {

}

func (rH RoutesHandler) createEnvironmentHandler(context *gin.Context) {

}

func (rH RoutesHandler) getServicesHandler(context *gin.Context) {

}

func (rH RoutesHandler) createServiceHandler(context *gin.Context) {

}

func (rH RoutesHandler) deleteServiceHandler(context *gin.Context) {

}

func (rH RoutesHandler) getServiceHandler(context *gin.Context) {

}

func (rH RoutesHandler) updateServiceHandler(context *gin.Context) {

}

func (rH RoutesHandler) deleteEnvironmentHandler(context *gin.Context) {

}

func (rH RoutesHandler) getServiceInstancesHandler(context *gin.Context) {

}

func (rH RoutesHandler) getServiceInstanceHandler(context *gin.Context) {

}

func NewRouter(ucHandler usecases.Usecases) RoutesHandler {
	return RoutesHandler{usecases: ucHandler}
}
