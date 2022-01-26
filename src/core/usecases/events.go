package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	eventDomain "github.com/AliceDiNunno/rack-controller/src/core/domain/eventDomain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/userDomain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (i interactor) FetchProjectVersions(project *domain.Project) ([]string, *e.Error) {
	versions, err := i.logCollection.ProjectVersions(project)

	if err != nil {
		return nil, e.Wrap(eventDomain.ErrUnableToFindEvents)
	}

	return versions, nil
}

func (i interactor) FetchProjectServers(project *domain.Project) ([]string, *e.Error) {
	return i.logCollection.ProjectServers(project)
}

func (i interactor) PushNewLogEntry(id uuid.UUID, request *request.ItemCreationRequest) *e.Error {
	project, error := i.projectRepository.GetProjectByIDAndKey(id, request.ProjectKey)

	if error != nil || project == nil {
		return e.Wrap(domain.ErrProjectNotFound)
	}

	logEntry := &eventDomain.LogEntry{
		ID:             primitive.NewObjectID(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ProjectID:      id,
		Identification: request.Identification,
		Data:           request.Data,
	}

	if logEntry.Data.GroupingID == "" {
		logEntry.Data.GroupingID = logEntry.Data.Fingerprint
	}

	return i.logCollection.AddLog(logEntry)
}

func (i interactor) FetchGroupingIdContent(project *domain.Project, groupingId string) (*eventDomain.LogEntry, *e.Error) {
	if !i.logCollection.IsGroupExist(project, groupingId) {
		return nil, e.Wrap(eventDomain.ErrGroupNotFound)
	}

	return i.logCollection.FindLastEntryForGroup(project, groupingId)
}

func (i interactor) FetchGroupingIdOccurrences(project *domain.Project, groupingId string) ([]string, *e.Error) {
	if !i.logCollection.IsGroupExist(project, groupingId) {
		return nil, e.Wrap(eventDomain.ErrGroupNotFound)
	}

	return i.logCollection.FindGroupOccurrences(project, groupingId)
}

func (i interactor) FetchGroupOccurrence(project *domain.Project, groupingId string, occurrence string) (*eventDomain.LogEntry, *e.Error) {
	if !i.logCollection.IsGroupExist(project, groupingId) {
		return nil, e.Wrap(eventDomain.ErrGroupNotFound)
	}

	return i.logCollection.FindGroupOccurrence(project, groupingId, occurrence)
}

func (i interactor) FetchProjectEnvironments(project *domain.Project) ([]string, *e.Error) {
	environments, err := i.logCollection.ProjectEnvironments(project)

	if err != nil {
		return nil, e.Wrap(eventDomain.ErrUnableToFindEvents)
	}

	return environments, nil
}

//TODO: should be groupings not events
func (i interactor) GetProjectsEvent(user *userDomain.User, project *domain.Project) ([]string, *e.Error) {
	if user == nil {
		return nil, e.Wrap(domain.ErrUserIsNil)
	}

	return i.logCollection.ProjectGroupingIds(project)
}
