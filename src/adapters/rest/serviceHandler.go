package rest

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (rH RoutesHandler) getServiceMiddleware(c *gin.Context) {
	project := rH.getProject(c)

	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	id, stderr := uuid.Parse(c.Param("service_id"))

	if stderr != nil {
		rH.handleError(c, e.Wrap(ErrFormValidation))
		return
	}

	service, err := rH.usecases.GetServiceById(project, id)

	if err != nil {
		rH.handleError(c, err.Append(domain.ErrEnvironmentNotFound))
		return
	}

	c.Set("service", service)
}

func (rH RoutesHandler) getService(c *gin.Context) *domain.Service {
	auth, exists := c.Get("service")

	if !exists {
		return nil
	}

	service := auth.(*domain.Service)

	return service
}

func (rH RoutesHandler) getServiceOfEnvironmentHandler(c *gin.Context) {
	service := rH.getService(c)

	if service == nil {
		return
	}

	environment := rH.getEnvironment(c)

	if environment == nil {
		return
	}

	serviceOfEnvironment, err := rH.usecases.GetServiceOfEnvironment(service, environment)

	if err != nil {
		rH.handleError(c, err.Append(domain.ErrEnvironmentNotFound))
		return
	}

	rH.handleSuccess(c, serviceOfEnvironment)
}

func (rH RoutesHandler) getServicesHandler(c *gin.Context) {
	project := rH.getProject(c)
	if project == nil {
		return
	}

	services, err := rH.usecases.GetServices(project)
	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, services)
}

func (rH RoutesHandler) createServiceHandler(c *gin.Context) {
	project := rH.getProject(c)
	if project == nil {
		return
	}

	var service request.ServiceCreationRequest
	stderr := c.BindJSON(&service)
	if stderr != nil {
		rH.handleError(c, e.Wrap(ErrFormValidation))
		return
	}

	err := rH.usecases.CreateService(project, &service)
	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, service)
}

func (rH RoutesHandler) deleteServiceHandler(c *gin.Context) {
	service := rH.getService(c)
	if service == nil {
		return
	}

	err := rH.usecases.DeleteService(service)
	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, nil)
}

func (rH RoutesHandler) getServiceHandler(c *gin.Context) {

}

func (rH RoutesHandler) updateServiceHandler(c *gin.Context) {

}

func (rH RoutesHandler) getServiceConfigHandler(c *gin.Context) {
	service := rH.getService(c)

	if service == nil {
		return
	}

	config, err := rH.usecases.GetServiceConfig(service)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, config)
}

func (rH RoutesHandler) restartServiceHandler(c *gin.Context) {
	print("Endpoint hit")
	service := rH.getService(c)

	if service == nil {
		return
	}

	err := rH.usecases.RestartService(service)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, service)
}

func (rH RoutesHandler) updateServiceConfigHandler(c *gin.Context) {
	service := rH.getService(c)

	if service == nil {
		return
	}

	var configRequest request.UpdateRequest

	if err := c.ShouldBindJSON(&configRequest); err != nil {
		rH.handleError(c, e.Wrap(ErrFormValidation))
		return
	}

	err := rH.usecases.UpdateServiceConfig(service, clusterDomain.EnvironmentListFromMap(configRequest.Data))

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, service)
}
