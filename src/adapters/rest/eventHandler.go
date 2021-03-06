package rest

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (rH RoutesHandler) fetchingGroupMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := rH.getAuthenticatedUser(c)
		if user == nil {
			return
		}

		project := rH.getProject(c)
		if project == nil {
			return
		}

		id := c.Param("grouping_id")

		if id == "" {
			rH.handleError(c, e.Wrap(ErrFormValidation))
			return
		}

		c.Set("grouping", id)
	}
}

func (rH RoutesHandler) getGrouping(c *gin.Context) string {
	grouping, exists := c.Get("grouping")

	if !exists {
		return ""
	}

	foundGrouping := grouping.(string)

	return foundGrouping
}

func (rH RoutesHandler) PushEventsHandler(c *gin.Context) {
	var creationRequest request.ItemCreationRequest

	id, stderr := uuid.Parse(c.Param("project_id"))

	if stderr != nil {
		rH.handleError(c, e.Wrap(ErrFormValidation))
		return
	}

	stderr = c.ShouldBind(&creationRequest)

	if stderr != nil {
		rH.handleError(c, e.Wrap(ErrFormValidation))
		return
	}

	err := rH.usecases.PushNewEvent(id, &creationRequest)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

func (rH RoutesHandler) SearchEventsInGroupingHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	project := rH.getProject(c)
	if project == nil {
		return
	}

	groupingId := c.Param("grouping_id")

	events, err := rH.usecases.FetchGroupingIdContent(project, groupingId)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, events)
}

func (rH RoutesHandler) GetEventsOccurencesHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	project := rH.getProject(c)
	if project == nil {
		return
	}

	groupingId := rH.getGrouping(c)
	if groupingId == "" {
		return
	}

	events, err := rH.usecases.FetchGroupingIdOccurrences(project, groupingId)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, events)
}

func (rH RoutesHandler) GetSpecificEventsHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	project := rH.getProject(c)
	if project == nil {
		return
	}

	groupingId := rH.getGrouping(c)
	if groupingId == "" {
		return
	}

	logId := c.Param("log_id")
	if logId == "" {
		rH.handleError(c, e.Wrap(ErrFormValidation))
	}

	events, err := rH.usecases.FetchGroupOccurrence(project, groupingId, logId)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, events)
}

func (rH RoutesHandler) GetServerHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	project := rH.getProject(c)
	if project == nil {
		return
	}

	servers, err := rH.usecases.FetchProjectServers(project)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, servers)
}

func (rH RoutesHandler) GetVersionHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	project := rH.getProject(c)
	if project == nil {
		return
	}

	versions, err := rH.usecases.FetchProjectVersions(project)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, versions)
}

func (rH RoutesHandler) GetEnvironmentHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	project := rH.getProject(c)
	if project == nil {
		return
	}

	environments, err := rH.usecases.FetchProjectEnvironments(project)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, environments)
}

func (rH RoutesHandler) GetItemsHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	project := rH.getProject(c)
	if project == nil {
		return
	}

	result, err := rH.usecases.GetProjectsEvent(user, project)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	rH.handleSuccess(c, result)
}
