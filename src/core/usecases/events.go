package usecases

import (
	"github.com/AliceDiNunno/go-logger/src/core/domain"
	"github.com/AliceDiNunno/go-logger/src/core/domain/request"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (i interactor) FetchProjectVersions(project *domain.Project) ([]string, error) {
	versions, err := i.logCollection.ProjectVersions(project)

	if err != nil {
		return nil, domain.ErrUnknownDBError
	}

	return versions, nil
}

func (i interactor) FetchProjectServers(project *domain.Project) ([]string, error) {
	return i.logCollection.ProjectServers(project)
}

func (i interactor) PushNewLogEntry(id uuid.UUID, request *request.ItemCreationRequest) error {
	project, error := i.projectRepo.FindByIdAndKey(id, request.ProjectKey)

	if error != nil || project == nil {
		return domain.ErrProjectNotFound
	}

	logEntry := &domain.LogEntry{
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

func (i interactor) FetchGroupingIdContent(project *domain.Project, groupingId string) (*domain.LogEntry, error) {
	if !i.logCollection.IsGroupExist(project, groupingId) {
		return nil, domain.ErrGroupNotFound
	}

	return i.logCollection.FindLastEntryForGroup(project, groupingId)
}

func (i interactor) FetchGroupingIdOccurrences(project *domain.Project, groupingId string) ([]string, error) {
	if !i.logCollection.IsGroupExist(project, groupingId) {
		return nil, domain.ErrGroupNotFound
	}

	return i.logCollection.FindGroupOccurrences(project, groupingId)
}

func (i interactor) FetchGroupOccurrence(project *domain.Project, groupingId string, occurrence string) (*domain.LogEntry, error) {
	if !i.logCollection.IsGroupExist(project, groupingId) {
		return nil, domain.ErrGroupNotFound
	}

	return i.logCollection.FindGroupOccurrence(project, groupingId, occurrence)
}

func (i interactor) FetchProjectEnvironments(project *domain.Project) ([]string, error) {
	environments, err := i.logCollection.ProjectEnvironments(project)

	if err != nil {
		return nil, domain.ErrUnknownDBError
	}

	return environments, nil
}

func (i interactor) GetProjectsContent(user *domain.User, project *domain.Project) ([]string, error) {
	if user == nil {
		return nil, domain.ErrFailedToGetUser
	}
	return i.logCollection.ProjectGroupingIds(project)
}
