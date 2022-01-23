package rest

import "github.com/gin-gonic/gin"

func SetRoutes(server GinServer, routesHandler RoutesHandler) {
	r := server.Router

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

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
	profileEndpoint := authenticatedEndpoint.Group("/me")
	profileEndpoint.GET("", routesHandler.getProfileHandler)
	profileEndpoint.GET("/roles", routesHandler.getRolesHandler)
	profileEndpoint.GET("/permissions", routesHandler.getPermissionsHandler)

	projectsEndpoint := authenticatedEndpoint.Group("/projects")
	projectsEndpoint.GET("", routesHandler.getProjectsHandler)
	projectsEndpoint.POST("", routesHandler.createProjectHandler)

	selectedProjectEndpoint := projectsEndpoint.Group("/:project_id", routesHandler.getProjectMiddleware)
	selectedProjectEndpoint.GET("", routesHandler.getProjectHandler)
	selectedProjectEndpoint.PUT("", routesHandler.updateProjectHandler)
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

	serviceSelectedEnvironmentEndpoint := selectedServiceEndpoint.Group("/environments/:environment_id", routesHandler.getProjectEnvironmentsMiddleware)
	serviceSelectedEnvironmentEndpoint.GET("", routesHandler.getServiceOfEnvironmentHandler)

	serviceInstanceEndpoint := serviceSelectedEnvironmentEndpoint.Group("/instances")
	serviceInstanceEndpoint.GET("", routesHandler.getServiceInstancesHandler)

	selectedServiceInstanceEndpoint := serviceInstanceEndpoint.Group("/:instance_id", routesHandler.getServiceInstanceMiddleware)
	selectedServiceInstanceEndpoint.GET("", routesHandler.getServiceInstanceHandler)
	selectedServiceInstanceEndpoint.DELETE("", routesHandler.deleteServiceInstanceHandler)
	selectedServiceInstanceEndpoint.GET("/logs", routesHandler.getServiceInstanceLogsHandler)
	selectedServiceInstanceEndpoint.GET("/events", routesHandler.getServiceInstanceEventsHandler)
}
