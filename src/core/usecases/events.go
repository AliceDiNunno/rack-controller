package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	eventDomain "github.com/AliceDiNunno/rack-controller/src/core/domain/eventDomain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/userDomain"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
)

func (i interactor) FetchProjectVersions(project *domain.Project) ([]eventDomain.Event, *e.Error) {
	versions, err := i.eventCollection.ProjectVersions(project)

	if err != nil {
		return nil, e.Wrap(eventDomain.ErrUnableToFindEvents)
	}

	return versions, nil
}

func (i interactor) FetchProjectServers(project *domain.Project) ([]eventDomain.Event, *e.Error) {
	return i.eventCollection.ProjectServers(project)
}

func (i interactor) PushNewEvent(id uuid.UUID, request *request.ItemCreationRequest) *e.Error {
	print("PUSH NEW EVENT")
	spew.Dump(request)
	projectKey, err := uuid.Parse(request.ProjectKey)

	if err != nil {
		return e.Wrap(eventDomain.ErrInvalidProjectKey)
	}

	project, error := i.projectRepository.GetProjectByIDAndKey(id, projectKey)

	if error != nil || project == nil {
		return e.Wrap(domain.ErrProjectNotFound)
	}

	eventEntry := &eventDomain.Event{
		ID:          uuid.New(),
		ProjectID:   id,
		Timestamp:   request.Data.Timestamp,
		Level:       request.Data.Level,
		StatusCode:  request.Data.StatusCode,
		Module:      request.Data.Module,
		Message:     request.Data.Message,
		Fingerprint: request.Data.Fingerprint,
	}

	if eventEntry.GroupingID == "" {
		eventEntry.GroupingID = eventEntry.Fingerprint
	}

	if !i.eventCollection.IsGroupExist(project, eventEntry.GroupingID) {
		i.eventCollection.AddEvent(eventEntry)
	}

	eventOccurrence := &eventDomain.EventOccurrence{
		ID:          uuid.New(),
		Timestamp:   request.Data.Timestamp,
		EventID:     eventEntry.ID,
		Platform:    request.Identification.Deployment.Platform,
		Source:      request.Identification.Deployment.Source,
		Hostname:    request.Identification.Deployment.Hostname,
		Environment: request.Identification.Deployment.Environment,
		Version:     request.Identification.Deployment.Version,
		UserID:      nil,
		IPAddress:   request.Identification.Client.IPAddress,
	}

	return i.eventCollection.AddOccurrence(eventOccurrence)
}

func (i interactor) FetchGroupingIdContent(project *domain.Project, groupingId string) (*eventDomain.Event, *e.Error) {
	if !i.eventCollection.IsGroupExist(project, groupingId) {
		return nil, e.Wrap(eventDomain.ErrGroupNotFound)
	}

	return i.eventCollection.FindLastEntryForGroup(project, groupingId)
}

func (i interactor) FetchGroupingIdOccurrences(project *domain.Project, groupingId string) ([]eventDomain.Event, *e.Error) {
	if !i.eventCollection.IsGroupExist(project, groupingId) {
		return nil, e.Wrap(eventDomain.ErrGroupNotFound)
	}

	return i.eventCollection.FindGroupOccurrences(project, groupingId)
}

func (i interactor) FetchGroupOccurrence(project *domain.Project, groupingId string, occurrence string) (*eventDomain.Event, *e.Error) {
	if !i.eventCollection.IsGroupExist(project, groupingId) {
		return nil, e.Wrap(eventDomain.ErrGroupNotFound)
	}

	return i.eventCollection.FindGroupOccurrence(project, groupingId, occurrence)
}

func (i interactor) FetchProjectEnvironments(project *domain.Project) ([]eventDomain.Event, *e.Error) {
	environments, err := i.eventCollection.ProjectEnvironments(project)

	if err != nil {
		return nil, e.Wrap(eventDomain.ErrUnableToFindEvents)
	}

	return environments, nil
}

//TODO: should be groupings not events
func (i interactor) GetProjectsEvent(user *userDomain.User, project *domain.Project) ([]eventDomain.Event, *e.Error) {
	if user == nil {
		return nil, e.Wrap(domain.ErrUserIsNil)
	}

	groupingIds, err := i.eventCollection.ProjectGroupingIds(project)

	if err != nil {
		spew.Dump(err)
		return nil, e.Wrap(eventDomain.ErrUnableToFindEvents)
	}

	groupings := []eventDomain.Event{}

	for _, groupingId := range groupingIds {
		grouping, err := i.eventCollection.FindLastEntryForGroup(project, groupingId.GroupingID)

		if err != nil {
			spew.Dump(err)
			return nil, e.Wrap(eventDomain.ErrUnableToFindEvents)
		}

		groupings = append(groupings, *grouping)
	}

	return groupings, nil
}
