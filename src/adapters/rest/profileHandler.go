package rest

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/gin-gonic/gin"
)

//TODO: separate user and profile
func (rH RoutesHandler) getProfileHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)

	if user == nil {
		rH.handleError(c, e.Wrap(ErrUnauthorized))
		return
	}

	rH.handleSuccess(c, user)
}

func (rH RoutesHandler) getRolesHandler(c *gin.Context) {
	//TODO: implement
}

func (rH RoutesHandler) getPermissionsHandler(c *gin.Context) {
	//TODO: implement
}
