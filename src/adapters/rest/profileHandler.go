package rest

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/gin-gonic/gin"
)

//TODO: separate user and profile
func (rH RoutesHandler) getProfileHandler(context *gin.Context) {
	user := rH.getAuthenticatedUser(context)

	if user == nil {
		rH.handleError(context, e.Wrap(ErrUnauthorized))
		return
	}

	context.JSON(200, success(user))
}

func (rH RoutesHandler) getRolesHandler(context *gin.Context) {
	//TODO: implement
}

func (rH RoutesHandler) getPermissionsHandler(context *gin.Context) {
	//TODO: implement
}
