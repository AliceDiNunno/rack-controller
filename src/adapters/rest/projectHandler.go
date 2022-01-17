package rest

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
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

func (rH RoutesHandler) updateProjectHandler(context *gin.Context) {

}

func (rH RoutesHandler) deleteProjectHandler(context *gin.Context) {

}

func (rH RoutesHandler) getProjectEnvHandler(context *gin.Context) {

}

func (rH RoutesHandler) updateProjectEnvHandler(context *gin.Context) {

}
