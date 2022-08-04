package rest

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (rH RoutesHandler) getProjectMiddleware(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	id, stderr := uuid.Parse(c.Param("project_id"))

	if stderr != nil {
		rH.handleError(c, e.Wrap(ErrUrlValidation))
		return
	}

	project, err := rH.usecases.GetProjectByID(user, id)

	if err != nil {
		rH.handleError(c, err.Append(domain.ErrProjectNotFound))
		return
	}

	c.Set("project", project)
}

func (rH RoutesHandler) getProject(c *gin.Context) *domain.Project {
	auth, exists := c.Get("project")

	if !exists {
		return nil
	}

	project := auth.(*domain.Project)

	return project
}

func (rH RoutesHandler) getProjectsHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)

	if user == nil {
		rH.handleError(c, e.Wrap(ErrUnauthorized))
		return
	}

	projects, err := rH.usecases.GetUserProjects(user)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, projects)
}

func (rH RoutesHandler) createProjectHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)

	if user == nil {
		rH.handleError(c, e.Wrap(ErrUnauthorized))
		return
	}

	var creationRequest request.CreateProjectRequest

	if err := c.ShouldBindJSON(&creationRequest); err != nil {
		rH.handleError(c, e.Wrap(ErrFormValidation))
		return
	}

	project, err := rH.usecases.CreateProject(user, creationRequest)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, project)
}

func (rH RoutesHandler) getProjectHandler(c *gin.Context) {
	project := rH.getProject(c)

	if project == nil {
		return
	}

	rH.handleSuccess(c, project)
}

func (rH RoutesHandler) deleteProjectHandler(c *gin.Context) {
	project := rH.getProject(c)
	if project == nil {
		return
	}

	err := rH.usecases.DeleteProject(project)
	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, nil)
}

func (rH RoutesHandler) getProjectConfigHandler(c *gin.Context) {
	project := rH.getProject(c)

	if project == nil {
		return
	}

	config, err := rH.usecases.GetProjectConfig(project)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, config)
}

func (rH RoutesHandler) updateProjectConfigHandler(c *gin.Context) {
	project := rH.getProject(c)

	if project == nil {
		return
	}

	var configRequest request.UpdateConfigData

	if err := c.ShouldBindJSON(&configRequest); err != nil {
		rH.handleError(c, e.Wrap(ErrFormValidation))
		return
	}

	env := clusterDomain.EnvironmentListFromMap(configRequest)
	err := rH.usecases.UpdateProjectConfig(project, env)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, env)
}
