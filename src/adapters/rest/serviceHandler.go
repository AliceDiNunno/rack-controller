package rest

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (rH RoutesHandler) getServiceMiddleware(context *gin.Context) {
	project := rH.getProject(context)

	user := rH.getAuthenticatedUser(context)
	if user == nil {
		return
	}

	id, stderr := uuid.Parse(context.Param("service_id"))

	if stderr != nil {
		rH.handleError(context, e.Wrap(ErrFormValidation))
		return
	}

	service, err := rH.usecases.GetServiceById(project, id)

	if err != nil {
		rH.handleError(context, err.Append(domain.ErrEnvironmentNotFound))
		return
	}

	context.Set("service", service)
}

func (rH RoutesHandler) getService(c *gin.Context) *domain.Service {
	auth, exists := c.Get("service")

	if !exists {
		return nil
	}

	service := auth.(*domain.Service)

	return service
}

func (rH RoutesHandler) getServiceOfEnvironmentHandler(context *gin.Context) {
	service := rH.getService(context)

	if service == nil {
		return
	}

	environment := rH.getEnvironment(context)

	if environment == nil {
		return
	}

	serviceOfEnvironment, err := rH.usecases.GetServiceOfEnvironment(service, environment)

	if err != nil {
		rH.handleError(context, err.Append(domain.ErrEnvironmentNotFound))
		return
	}

	context.JSON(200, success(serviceOfEnvironment))
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
	service := rH.getService(context)
	if service == nil {
		return
	}

	err := rH.usecases.DeleteService(service)
	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(200, success(nil))
}

func (rH RoutesHandler) getServiceHandler(context *gin.Context) {

}

func (rH RoutesHandler) updateServiceHandler(context *gin.Context) {

}

func (rH RoutesHandler) getServiceConfigHandler(context *gin.Context) {
	service := rH.getService(context)

	if service == nil {
		return
	}

	config, err := rH.usecases.GetServiceConfig(service)

	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(200, success(config))
}

func (rH RoutesHandler) restartServiceHandler(context *gin.Context) {
	print("Endpoint hit")
	service := rH.getService(context)

	if service == nil {
		return
	}

	err := rH.usecases.RestartService(service)

	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(200, success(service))
}

func (rH RoutesHandler) updateServiceConfigHandler(context *gin.Context) {
	service := rH.getService(context)

	if service == nil {
		return
	}

	var configRequest request.UpdateConfigRequest

	if err := context.ShouldBindJSON(&configRequest); err != nil {
		rH.handleError(context, e.Wrap(ErrFormValidation))
		return
	}

	err := rH.usecases.UpdateServiceConfig(service, configRequest)

	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(200, success(service))
}
