package rest

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (rH RoutesHandler) getProjectEnvironmentsMiddleware(context *gin.Context) {
	project := rH.getProject(context)

	user := rH.getAuthenticatedUser(context)
	if user == nil {
		return
	}

	id, stderr := uuid.Parse(context.Param("environment_id"))

	if stderr != nil {
		rH.handleError(context, e.Wrap(ErrFormValidation))
		return
	}

	environment, err := rH.usecases.GetEnvironmentByID(project, id)

	if err != nil {
		rH.handleError(context, err.Append(domain.ErrEnvironmentNotFound))
		return
	}

	context.Set("environment", environment)
}

func (rH RoutesHandler) getEnvironment(c *gin.Context) *domain.Environment {
	auth, exists := c.Get("environment")

	if !exists {
		return nil
	}

	environment := auth.(*domain.Environment)

	return environment
}

func (rH RoutesHandler) getEnvironmentHandler(context *gin.Context) {

}

func (rH RoutesHandler) getEnvironmentsHandler(context *gin.Context) {
	project := rH.getProject(context)
	if project == nil {
		return
	}

	environments, err := rH.usecases.GetEnvironments(project)
	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(200, success(environments))
}

func (rH RoutesHandler) createEnvironmentHandler(context *gin.Context) {
	project := rH.getProject(context)
	if project == nil {
		return
	}

	var environment request.EnvironmentCreationRequest
	stderr := context.BindJSON(&environment)
	if stderr != nil {
		rH.handleError(context, e.Wrap(ErrFormValidation))
		return
	}

	err := rH.usecases.CreateEnvironment(project, &environment)
	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(201, success(environment))
}

func (rH RoutesHandler) deleteEnvironmentHandler(context *gin.Context) {

}

func (rH RoutesHandler) getEnvironmentConfigHandler(context *gin.Context) {
	environment := rH.getEnvironment(context)

	if environment == nil {
		return
	}

	config, err := rH.usecases.GetEnvironmentConfig(environment)

	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(200, success(config))
}

func (rH RoutesHandler) updateEnvironmentConfigHandler(context *gin.Context) {
	environment := rH.getEnvironment(context)

	if environment == nil {
		return
	}

	var configRequest request.UpdateConfigRequest

	if err := context.ShouldBindJSON(&configRequest); err != nil {
		rH.handleError(context, e.Wrap(ErrFormValidation))
		return
	}

	err := rH.usecases.UpdateEnvironmentConfig(environment, configRequest)

	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(200, success(environment))
}
