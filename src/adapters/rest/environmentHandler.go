package rest

import (
	"github.com/gin-gonic/gin"
)

func (rH RoutesHandler) getProjectEnvironmentsMiddleware(context *gin.Context) {
	/*user := rH.getAuthenticatedUser(context)
	if user == nil {
		return
	}

	id, stderr := uuid.Parse(context.Param("project_id"))

	if stderr != nil {
		rH.handleError(context, e.Wrap(ErrFormValidation))
		return
	}

	project, err := rH.usecases.GetProjectByID(user, id)

	if err != nil {
		rH.handleError(context, err.Append(domain.ErrProjectNotFound))
		return
	}

	context.Set("project", project)*/
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

}

func (rH RoutesHandler) deleteEnvironmentHandler(context *gin.Context) {

}
