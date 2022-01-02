package rest

import (
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/gin-gonic/gin"
)

//if fetchingUserMiddleware is a root middleware there is no reason than domain.user is empty at this point
//We will assume fetchingUserMiddleware is a root middleware
func (rH RoutesHandler) getAuthenticatedUser(c *gin.Context) *domain.User {
	auth, exists := c.Get("authenticatedUser")

	if !exists {
		return nil
	}

	authenticatedUser := auth.(domain.User)

	return &authenticatedUser
}
