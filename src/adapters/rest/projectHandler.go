package rest

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (rH RoutesHandler) getProjectMiddleware(context *gin.Context) {
	user := rH.getAuthenticatedUser(context)
	if user == nil {
		return
	}

	id, stderr := uuid.Parse(context.Param("project_id"))

	if stderr != nil {
		rH.handleError(context, e.Wrap(ErrUrlValidation))
		return
	}

	project, err := rH.usecases.GetProjectByID(user, id)

	if err != nil {
		rH.handleError(context, err.Append(domain.ErrProjectNotFound))
		return
	}

	context.Set("project", project)
}

func (rH RoutesHandler) getProject(c *gin.Context) *domain.Project {
	auth, exists := c.Get("project")

	if !exists {
		return nil
	}

	project := auth.(*domain.Project)

	return project
}

func (rH RoutesHandler) getProjectsHandler(context *gin.Context) {
	user := rH.getAuthenticatedUser(context)

	if user == nil {
		rH.handleError(context, e.Wrap(ErrUnauthorized))
		return
	}

	projects, err := rH.usecases.GetUserProjects(user)

	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(200, success(projects))
}

func (rH RoutesHandler) createProjectHandler(context *gin.Context) {
	user := rH.getAuthenticatedUser(context)

	if user == nil {
		rH.handleError(context, e.Wrap(ErrUnauthorized))
		return
	}

	var creationRequest request.CreateProjectRequest

	if err := context.ShouldBindJSON(&creationRequest); err != nil {
		rH.handleError(context, e.Wrap(ErrFormValidation))
		return
	}

	project, err := rH.usecases.CreateProject(user, creationRequest)

	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(201, success(project))
}

func (rH RoutesHandler) getProjectHandler(context *gin.Context) {
	project := rH.getProject(context)

	if project == nil {
		return
	}

	context.JSON(200, success(project))
}

func (rH RoutesHandler) deleteProjectHandler(context *gin.Context) {
	project := rH.getProject(context)
	if project == nil {
		return
	}

	err := rH.usecases.DeleteProject(project)
	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(200, success(nil))
}

func (rH RoutesHandler) getProjectConfigHandler(context *gin.Context) {
	project := rH.getProject(context)

	if project == nil {
		return
	}

	config, err := rH.usecases.GetProjectConfig(project)

	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(200, success(config))
}

func (rH RoutesHandler) updateProjectConfigHandler(context *gin.Context) {
	project := rH.getProject(context)

	if project == nil {
		return
	}

	var configRequest request.UpdateRequest

	if err := context.ShouldBindJSON(&configRequest); err != nil {
		rH.handleError(context, e.Wrap(ErrFormValidation))
		return
	}

	env := clusterDomain.EnvironmentListFromMap(configRequest.Data)
	err := rH.usecases.UpdateProjectConfig(project, env)

	if err != nil {
		rH.handleError(context, err)
		return
	}

	context.JSON(200, success(env))
}
