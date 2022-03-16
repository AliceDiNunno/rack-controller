package rest

import (
	"github.com/gin-gonic/gin"
)

func (rH RoutesHandler) GetDomainsHandler(c *gin.Context) {
	domains, err := rH.usecases.GetDomainNames()

	if err != nil {
		rH.handleError(c, err)
		return
	}

	c.JSON(200, success(domains))
}
