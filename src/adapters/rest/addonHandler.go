package rest

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	request "github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (rH RoutesHandler) getServiceAddonMiddleware(c *gin.Context) {
	service := rH.getService(c)

	if service == nil {
		return
	}

	addonID, stderr := uuid.Parse(c.Param("addon_id"))

	if stderr != nil {
		rH.handleError(c, e.Wrap(ErrUrlValidation))
		return
	}

	addon, err := rH.usecases.GetAddonById(service, addonID)

	if err != nil {
		rH.handleError(c, e.Wrap(domain.ErrAddonNotFound))
		return
	}

	c.Set("addon", addon)
}

func (rH RoutesHandler) getAddon(c *gin.Context) *domain.Addon {
	auth, exists := c.Get("addon")

	if !exists {
		return nil
	}

	addon := auth.(*domain.Addon)

	return addon
}

func (rH RoutesHandler) getServiceAddonsHandler(c *gin.Context) {
	service := rH.getService(c)

	if service == nil {
		return
	}

	addons, err := rH.usecases.GetAddons(service)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, addons)
}

func (rH RoutesHandler) createServiceAddonHandler(c *gin.Context) {
	service := rH.getService(c)

	if service == nil {
		return
	}

	var addonRequest request.AddonCreationRequest

	if stderr := c.ShouldBindJSON(&addonRequest); stderr != nil {
		rH.handleError(c, e.Wrap(stderr).Append(ErrFormValidation))
		return
	}

	addon, err := rH.usecases.CreateAddon(service, &addonRequest)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, addon)
}

func (rH RoutesHandler) getSelectedServiceAddonHandler(c *gin.Context) {

}

func (rH RoutesHandler) deleteSelectedServiceAddonHandler(c *gin.Context) {

}
