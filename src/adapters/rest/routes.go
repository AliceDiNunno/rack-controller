package rest

import (
	"github.com/gin-gonic/gin"
)

/*
func (rH RoutesHandler) debugRoute(c *gin.Context) {
	print("ROUTE")

	print(c.FullPath())
}*/

func SetRoutes(server GinServer, routesHandler RoutesHandler) {
	r := server.Router

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	//r.Use(routesHandler.debugRoute)

	r.NoRoute(routesHandler.endpointNotFound)

	r.GET("/health", routesHandler.GetHealthHandler)

	authenticationEndpoint := r.Group("/authentication")

	//This is not oauth, but we need to be able to authenticate with a token
	//this is a private project, so we don't need to worry about security too much
	//fetching a token with a username and password
	authenticationEndpoint.POST("/token", routesHandler.createAuthTokenHandler)
	authenticationEndpoint.POST("/jwt", routesHandler.createJwtTokenHandler)
	authenticationEndpoint.DELETE("/jwt", routesHandler.deleteJwtTokenHandler)

	authenticatedEndpoint := r.Group("/", routesHandler.verifyAuthenticationMiddleware)

	nodeEndpoint := authenticatedEndpoint.Group("/nodes")
	nodeEndpoint.GET("", routesHandler.getNodesHandler)
	nodeEndpoint.GET("/:node_id", routesHandler.getSpecificNodeHandler)
	nodeEndpoint.GET("/:node_id/instances", routesHandler.getSpecificNodeInstancesHandler)

	profileEndpoint := authenticatedEndpoint.Group("/me")
	profileEndpoint.GET("", routesHandler.getProfileHandler)
	profileEndpoint.GET("/roles", routesHandler.getRolesHandler)
	profileEndpoint.GET("/permissions", routesHandler.getPermissionsHandler)

	projectsEndpoint := authenticatedEndpoint.Group("/projects")
	projectsEndpoint.GET("", routesHandler.getProjectsHandler)
	projectsEndpoint.POST("", routesHandler.createProjectHandler)
	projectsEndpoint.POST("/events", routesHandler.PushEventsHandler) //Push a log

	selectedProjectEndpoint := projectsEndpoint.Group("/:project_id", routesHandler.getProjectMiddleware)
	selectedProjectEndpoint.GET("", routesHandler.getProjectHandler)
	selectedProjectEndpoint.DELETE("", routesHandler.deleteProjectHandler)
	selectedProjectEndpoint.GET("/config", routesHandler.getProjectConfigHandler)
	selectedProjectEndpoint.POST("/config", routesHandler.updateProjectConfigHandler)

	//This endpoint is used to get the environment list and create an environment
	environmentEndpoint := selectedProjectEndpoint.Group("/environments")
	environmentEndpoint.GET("", routesHandler.getEnvironmentsHandler)
	environmentEndpoint.POST("", routesHandler.createEnvironmentHandler)

	selectedEnvironmentEndpoint := environmentEndpoint.Group("/:environment_id", routesHandler.getProjectEnvironmentsMiddleware)
	selectedEnvironmentEndpoint.GET("", routesHandler.getEnvironmentHandler)
	selectedEnvironmentEndpoint.DELETE("", routesHandler.deleteEnvironmentHandler)
	selectedEnvironmentEndpoint.GET("/config", routesHandler.getEnvironmentConfigHandler)
	selectedEnvironmentEndpoint.POST("/config", routesHandler.updateEnvironmentConfigHandler)

	serviceEndpoint := selectedProjectEndpoint.Group("/services")
	serviceEndpoint.GET("", routesHandler.getServicesHandler)
	serviceEndpoint.POST("", routesHandler.createServiceHandler)

	selectedServiceEndpoint := serviceEndpoint.Group("/:service_id", routesHandler.getServiceMiddleware)
	selectedServiceEndpoint.GET("", routesHandler.getServiceHandler)
	selectedServiceEndpoint.DELETE("", routesHandler.deleteServiceHandler)
	selectedServiceEndpoint.GET("/config", routesHandler.getServiceConfigHandler)
	selectedServiceEndpoint.POST("/config", routesHandler.updateServiceConfigHandler)

	selectedServiceEndpoint.GET("/add-ons", routesHandler.getServiceAddonsHandler)
	selectedServiceEndpoint.POST("/add-ons", routesHandler.createServiceAddonHandler)

	selectedAddonEndpoint := selectedServiceEndpoint.Group("/add-ons/:addon_id", routesHandler.getServiceAddonMiddleware)
	selectedAddonEndpoint.GET("", routesHandler.getSelectedServiceAddonHandler)
	selectedAddonEndpoint.DELETE("", routesHandler.deleteSelectedServiceAddonHandler)

	serviceSelectedEnvironmentEndpoint := selectedServiceEndpoint.Group("/environments/:environment_id", routesHandler.getProjectEnvironmentsMiddleware)
	serviceSelectedEnvironmentEndpoint.GET("", routesHandler.getServiceOfEnvironmentHandler)
	serviceSelectedEnvironmentEndpoint.POST("/restart", routesHandler.restartServiceHandler)

	serviceInstanceEndpoint := serviceSelectedEnvironmentEndpoint.Group("/instances")
	serviceInstanceEndpoint.GET("", routesHandler.getServiceInstancesHandler)

	selectedServiceInstanceEndpoint := serviceInstanceEndpoint.Group("/:instance_name", routesHandler.getServiceInstanceMiddleware)
	selectedServiceInstanceEndpoint.GET("", routesHandler.getServiceInstanceHandler)
	selectedServiceInstanceEndpoint.DELETE("", routesHandler.deleteServiceInstanceHandler)
	selectedServiceInstanceEndpoint.GET("/logs", routesHandler.getServiceInstanceLogsHandler)

	//Domain and ingress endpoints
	domainEndpoint := authenticatedEndpoint.Group("/domains")
	domainEndpoint.GET("", routesHandler.GetDomainsHandler)

	//Events routes
	serviceEvents := selectedProjectEndpoint.Group("/events")
	environmentGroup := serviceEvents.Group("/environment")
	environmentGroup.GET("", routesHandler.GetEnvironmentHandler) //Getting all declared environments for a project

	versionGroup := serviceEvents.Group("/version")
	versionGroup.GET("", routesHandler.GetVersionHandler) //Getting all declared version for a project

	serverGroup := serviceEvents.Group("/server")
	serverGroup.GET("", routesHandler.GetServerHandler) //Getting all declared servers for a project

	itemsGroup := serviceEvents.Group("/items")
	itemsGroup.GET("", routesHandler.GetItemsHandler) //Search all grouping ids

	logsGroup := itemsGroup.Group("/:grouping_id", routesHandler.fetchingGroupMiddleware())
	logsGroup.GET("/occurrences", routesHandler.GetEventsOccurencesHandler)       //Getting a specific log id
	logsGroup.GET("", routesHandler.SearchEventsInGroupingHandler)                //Search all logs (corresponding to a grouping ID)
	logsGroup.GET("/occurrences/:log_id", routesHandler.GetSpecificEventsHandler) //Getting a specific log id
}
