package rest

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/gin-gonic/gin"
)

func (rH RoutesHandler) getServiceMiddleware(context *gin.Context) {

}

func (rH RoutesHandler) getServiceOfEnvironmentHandler(context *gin.Context) {

}

func (rH RoutesHandler) getServicesHandler(context *gin.Context) {
	project := rH.getProject(context)
	if project == nil {
		return
	}

	services, err := rH.usecases.GetServices(project)
	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(200, success(services))
}

func (rH RoutesHandler) createServiceHandler(context *gin.Context) {
	project := rH.getProject(context)
	if project == nil {
		return
	}

	var service request.ServiceCreationRequest
	stderr := context.BindJSON(&service)
	if stderr != nil {
		rH.handleError(context, e.Wrap(ErrFormValidation))
		return
	}

	err := rH.usecases.CreateService(project, &service)
	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(201, success(service))
}

func (rH RoutesHandler) deleteServiceHandler(context *gin.Context) {

}

func (rH RoutesHandler) getServiceHandler(context *gin.Context) {

}

func (rH RoutesHandler) updateServiceHandler(context *gin.Context) {

}

func (rH RoutesHandler) getServiceConfigHandler(context *gin.Context) {

}

func (rH RoutesHandler) updateServiceConfigHandler(context *gin.Context) {

}
