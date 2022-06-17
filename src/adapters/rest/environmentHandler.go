package rest

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (rH RoutesHandler) getProjectEnvironmentsMiddleware(c *gin.Context) {
	project := rH.getProject(c)

	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	id, stderr := uuid.Parse(c.Param("environment_id"))

	if stderr != nil {
		rH.handleError(c, e.Wrap(ErrFormValidation))
		return
	}

	environment, err := rH.usecases.GetEnvironmentByID(project, id)

	if err != nil {
		rH.handleError(c, err.Append(domain.ErrEnvironmentNotFound))
		return
	}

	c.Set("environment", environment)
}

func (rH RoutesHandler) getEnvironment(c *gin.Context) *domain.Environment {
	auth, exists := c.Get("environment")

	if !exists {
		return nil
	}

	environment := auth.(*domain.Environment)

	return environment
}

func (rH RoutesHandler) getEnvironmentHandler(c *gin.Context) {

}

func (rH RoutesHandler) getEnvironmentsHandler(c *gin.Context) {
	project := rH.getProject(c)
	if project == nil {
		return
	}

	environments, err := rH.usecases.GetEnvironments(project)
	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, environments)
}

func (rH RoutesHandler) createEnvironmentHandler(c *gin.Context) {
	project := rH.getProject(c)
	if project == nil {
		return
	}

	var environment request.EnvironmentCreationRequest
	stderr := c.BindJSON(&environment)
	if stderr != nil {
		rH.handleError(c, e.Wrap(ErrFormValidation))
		return
	}

	err := rH.usecases.CreateEnvironment(project, &environment)
	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, environment)
}

func (rH RoutesHandler) deleteEnvironmentHandler(c *gin.Context) {
	print("DELENV: start")
	environment := rH.getEnvironment(c)
	if environment == nil {
		return
	}

	print("DELENV: got environment")
	err := rH.usecases.DeleteEnvironment(environment)
	if err != nil {
		print("DELENV: error " + err.Err.Error())
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, nil)
}

func (rH RoutesHandler) getEnvironmentConfigHandler(c *gin.Context) {
	environment := rH.getEnvironment(c)

	if environment == nil {
		return
	}

	config, err := rH.usecases.GetEnvironmentConfig(environment)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, config)
}

func (rH RoutesHandler) updateEnvironmentConfigHandler(c *gin.Context) {
	environment := rH.getEnvironment(c)

	if environment == nil {
		return
	}

	var configRequest request.UpdateConfigData

	if err := c.ShouldBindJSON(&configRequest); err != nil {
		rH.handleError(c, e.Wrap(ErrFormValidation))
		return
	}
	/*
		err := rH.usecases.UpdateEnvironmentConfig(environment, configRequest)

		if err != nil {
			rH.handleError(c, err)
			return
		}*/

	rH.handleSuccess(c, environment)
}
