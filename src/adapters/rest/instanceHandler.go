package rest

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	"github.com/gin-gonic/gin"
)

func (rH RoutesHandler) getServiceInstanceMiddleware(c *gin.Context) {
	service := rH.getService(c)

	if service == nil {
		return
	}

	environment := rH.getEnvironment(c)

	if environment == nil {
		return
	}

	instanceName := c.Param("instance_name")

	instance, err := rH.usecases.GetInstanceByName(service, environment, instanceName)

	if err != nil {
		rH.handleError(c, e.Wrap(domain.ErrInstanceNotFound))
		return
	}

	c.Set("instance", instance)
}

func (rH RoutesHandler) getInstance(c *gin.Context) *clusterDomain.Pod {
	auth, exists := c.Get("instance")

	if !exists {
		return nil
	}

	instance := auth.(*clusterDomain.Pod)

	return instance
}

func (rH RoutesHandler) getServiceInstancesHandler(c *gin.Context) {
	service := rH.getService(c)

	if service == nil {
		return
	}

	environments := rH.getEnvironment(c)

	if environments == nil {
		return
	}

	instances, err := rH.usecases.GetServiceInstances(service, environments)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, instances)
}

func (rH RoutesHandler) getServiceInstanceHandler(c *gin.Context) {
	service := rH.getService(c)

	if service == nil {
		return
	}

	environment := rH.getEnvironment(c)

	if environment == nil {
		return
	}

	instance := rH.getInstance(c)

	if instance == nil {
		return
	}

	rH.handleSuccess(c, instance)
}

func (rH RoutesHandler) deleteServiceInstanceHandler(c *gin.Context) {
	service := rH.getService(c)

	if service == nil {
		return
	}

	environment := rH.getEnvironment(c)

	if environment == nil {
		return
	}

	instance := rH.getInstance(c)

	if instance == nil {
		return
	}

	err := rH.usecases.DeleteInstance(service, environment, instance)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, nil)
}

func (rH RoutesHandler) getServiceInstanceEventsHandler(c *gin.Context) {

}

func (rH RoutesHandler) getServiceInstanceLogsHandler(c *gin.Context) {
	service := rH.getService(c)

	if service == nil {
		return
	}

	environment := rH.getEnvironment(c)

	if environment == nil {
		return
	}

	instance := rH.getInstance(c)

	if instance == nil {
		return
	}

	logs, err := rH.usecases.GetInstanceLogs(service, environment, instance)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, logs)
}
