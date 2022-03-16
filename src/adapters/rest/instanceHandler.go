package rest

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	"github.com/gin-gonic/gin"
)

func (rH RoutesHandler) getServiceInstanceMiddleware(context *gin.Context) {
	service := rH.getService(context)

	if service == nil {
		return
	}

	environment := rH.getEnvironment(context)

	if environment == nil {
		return
	}

	instanceName := context.Param("instance_name")

	instance, err := rH.usecases.GetInstanceByName(service, environment, instanceName)

	if err != nil {
		rH.handleError(context, e.Wrap(domain.ErrInstanceNotFound))
		return
	}

	context.Set("instance", instance)
}

func (rH RoutesHandler) getInstance(c *gin.Context) *clusterDomain.Pod {
	auth, exists := c.Get("instance")

	if !exists {
		return nil
	}

	instance := auth.(*clusterDomain.Pod)

	return instance
}

func (rH RoutesHandler) getServiceInstancesHandler(context *gin.Context) {
	service := rH.getService(context)

	if service == nil {
		return
	}

	environments := rH.getEnvironment(context)

	if environments == nil {
		return
	}

	instances, err := rH.usecases.GetServiceInstances(service, environments)

	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(200, success(instances))
}

func (rH RoutesHandler) getServiceInstanceHandler(context *gin.Context) {
	service := rH.getService(context)

	if service == nil {
		return
	}

	environment := rH.getEnvironment(context)

	if environment == nil {
		return
	}

	instance := rH.getInstance(context)

	if instance == nil {
		return
	}

	context.JSON(200, success(instance))
}

func (rH RoutesHandler) deleteServiceInstanceHandler(context *gin.Context) {
	service := rH.getService(context)

	if service == nil {
		return
	}

	environment := rH.getEnvironment(context)

	if environment == nil {
		return
	}

	instance := rH.getInstance(context)

	if instance == nil {
		return
	}

	err := rH.usecases.DeleteInstance(service, environment, instance)

	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(200, success(nil))
}

func (rH RoutesHandler) getServiceInstanceEventsHandler(context *gin.Context) {

}

func (rH RoutesHandler) getServiceInstanceLogsHandler(context *gin.Context) {
	service := rH.getService(context)

	if service == nil {
		return
	}

	environment := rH.getEnvironment(context)

	if environment == nil {
		return
	}

	instance := rH.getInstance(context)

	if instance == nil {
		return
	}

	logs, err := rH.usecases.GetInstanceLogs(service, environment, instance)

	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(200, success(logs))
}
