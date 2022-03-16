package postgres

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	logDomain "github.com/AliceDiNunno/rack-controller/src/core/domain/eventDomain"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type eventsRepo struct {
	db *gorm.DB
}

type TracebackEntry struct {
	gorm.Model

	ID uuid.UUID

	Filename string
	Line     int
	Method   string

	Traceback   *Traceback
	TracebackID uuid.UUID
}

type Traceback struct {
	gorm.Model

	ID uuid.UUID

	Message    string
	Traceback  []TracebackEntry
	LogEntry   *LogEntry
	LogEntryID uuid.UUID
}

type LogEntry struct {
	gorm.Model

	//Object metadata
	ID        uuid.UUID `gorm:"uniqueIndex:idx_id"`
	CreatedAt time.Time
	UpdatedAt time.Time

	//Entry metadata
	Timestamp   time.Time
	GroupingID  string `gorm:"index:idx_group"`
	Fingerprint string

	//Identification (for reproduction)
	Platform    string
	Source      string
	ProjectID   uuid.UUID `gorm:"index:idx_project"`
	Hostname    string
	Environment string
	Level       string
	Version     string
	Trace       *Traceback
	TraceID     uuid.UUID
	NestedTrace []Traceback
	UserID      *uuid.UUID
	IPAddress   string
	StatusCode  int

	//Entry Data
	Message          string
	AdditionalFields string
}

func tracebackToDomain(traceback *Traceback) logDomain.Traceback {
	if traceback == nil {
		return logDomain.Traceback{}
	}

	var tracebackEntries []logDomain.TracebackEntry

	for _, entry := range traceback.Traceback {
		tracebackEntries = append(tracebackEntries, logDomain.TracebackEntry{
			Filename: entry.Filename,
			Line:     entry.Line,
			Method:   entry.Method,
		})
	}

	return logDomain.Traceback{
		Message:   traceback.Message,
		Traceback: tracebackEntries,
	}
}

func tracebackFromDomain(traceback logDomain.Traceback) Traceback {
	var tracebackEntries []TracebackEntry
	for _, entry := range traceback.Traceback {
		tracebackEntries = append(tracebackEntries, TracebackEntry{
			ID:       uuid.New(),
			Filename: entry.Filename,
			Line:     entry.Line,
			Method:   entry.Method,
		})
	}

	return Traceback{
		ID:        uuid.New(),
		Message:   traceback.Message,
		Traceback: tracebackEntries,
	}
}

func tracebacksToDomain(traceback []Traceback) []logDomain.Traceback {
	var tracebacks []logDomain.Traceback

	for _, tb := range traceback {
		tracebacks = append(tracebacks, tracebackToDomain(&tb))
	}

	return tracebacks
}

func tracebacksFromDomain(traceback []logDomain.Traceback) []Traceback {
	var tracebacks []Traceback

	for _, tb := range traceback {
		tracebacks = append(tracebacks, tracebackFromDomain(tb))
	}

	return tracebacks
}

func logEntryToDomain(entry *LogEntry) *logDomain.LogEntry {
	return &logDomain.LogEntry{
		ID:        entry.ID,
		CreatedAt: entry.CreatedAt,
		UpdatedAt: entry.UpdatedAt,
		ProjectID: entry.ProjectID,
		Identification: logDomain.LogIdentification{
			Client: logDomain.LogClientIdentification{
				UserID:    entry.UserID,
				IPAddress: entry.IPAddress,
			},
			Deployment: logDomain.LogDeploymentIdentification{
				Platform:    entry.Platform,
				Source:      entry.Source,
				Hostname:    entry.Hostname,
				Environment: entry.Environment,
				Version:     entry.Version,
			},
		},
		Data: logDomain.LogData{
			Timestamp:   entry.Timestamp,
			GroupingID:  entry.GroupingID,
			Fingerprint: entry.Fingerprint,
			Level:       entry.Level,
			Trace:       tracebackToDomain(entry.Trace),
			NestedTrace: tracebacksToDomain(entry.NestedTrace),
			Message:     entry.Message,
			StatusCode:  entry.StatusCode,
			//AdditionalFields: entry.AdditionalFields,
		},
	}
}

func entriesToDomain(entries []LogEntry) []logDomain.LogEntry {
	var logEntries []logDomain.LogEntry

	for _, entry := range entries {
		logEntries = append(logEntries, *logEntryToDomain(&entry))
	}

	return logEntries
}

func logEntryFromDomain(entry *logDomain.LogEntry) *LogEntry {
	trace := tracebackFromDomain(entry.Data.Trace)

	return &LogEntry{
		ID:          entry.ID,
		CreatedAt:   entry.CreatedAt,
		UpdatedAt:   entry.UpdatedAt,
		ProjectID:   entry.ProjectID,
		Timestamp:   entry.Data.Timestamp,
		GroupingID:  entry.Data.GroupingID,
		Fingerprint: entry.Data.Fingerprint,
		Level:       entry.Data.Level,
		Trace:       &trace,
		//AdditionalFields: entry.Data.AdditionalFields,
		NestedTrace: tracebacksFromDomain(entry.Data.NestedTrace),
		StatusCode:  entry.Data.StatusCode,
		Message:     entry.Data.Message,
		Platform:    entry.Identification.Deployment.Platform,
		Source:      entry.Identification.Deployment.Source,
		Hostname:    entry.Identification.Deployment.Hostname,
		Environment: entry.Identification.Deployment.Environment,
		Version:     entry.Identification.Deployment.Version,
		UserID:      entry.Identification.Client.UserID,
		IPAddress:   entry.Identification.Client.IPAddress,
	}
}

////////

func (c eventsRepo) AddLog(entry *logDomain.LogEntry) *e.Error {
	entryFromDomain := logEntryFromDomain(entry)

	if err := c.db.Create(entryFromDomain).Error; err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (c eventsRepo) ProjectVersions(project *domain.Project) ([]logDomain.LogEntry, *e.Error) {
	var logEntry []logDomain.LogEntry

	err := c.db.Distinct("version").Where("project_id = ?", project.ID).Find(&logEntry).Error

	if err != nil {
		return nil, e.Wrap(err)
	}

	return logEntry, nil
}

func (c eventsRepo) ProjectEnvironments(project *domain.Project) ([]logDomain.LogEntry, *e.Error) {
	var logEntry []logDomain.LogEntry

	err := c.db.Distinct("environment").Where("project_id = ?", project.ID).Find(&logEntry).Error

	if err != nil {
		return nil, e.Wrap(err)
	}

	return logEntry, nil
}

func (c eventsRepo) ProjectServers(project *domain.Project) ([]logDomain.LogEntry, *e.Error) {
	var logEntry []logDomain.LogEntry

	err := c.db.Distinct("environment").Where("project_id = ?", project.ID).Find(&logEntry).Error

	if err != nil {
		return nil, e.Wrap(err)
	}

	return logEntry, nil
}

func (c eventsRepo) ProjectGroupingIds(project *domain.Project) ([]logDomain.LogEntry, *e.Error) {
	var logEntry []LogEntry

	err := c.db.Distinct("grouping_id").Where("project_id = ?", project.ID).Find(&logEntry).Error

	if err != nil {
		return nil, e.Wrap(err)
	}

	return entriesToDomain(logEntry), nil
}

func (c eventsRepo) FindLastEntryForGroup(project *domain.Project, groupingId string) (*logDomain.LogEntry, *e.Error) {
	var logEntry *LogEntry

	err := c.db.Where("grouping_id = ?", groupingId).Order("created_at desc").Last(&logEntry).Error

	if err != nil {
		return nil, e.Wrap(err)
	}

	return logEntryToDomain(logEntry), nil
}

func (c eventsRepo) FindGroupOccurrences(project *domain.Project, groupingId string) ([]logDomain.LogEntry, *e.Error) {
	var entries []logDomain.LogEntry

	err := c.db.Where("grouping_id = ?", groupingId).Order("created_at desc").Last(&entries).Error

	if err != nil {
		spew.Dump(err)
		return nil, e.Wrap(err)
	}

	return entries, nil
}

func (c eventsRepo) FindGroupOccurrence(project *domain.Project, groupingId string, occurenceId string) (*logDomain.LogEntry, *e.Error) {
	return nil, nil

	/*var entry logEntry

	queryOptions := options.FindOneOptions{}

	queryOptions.SetSort(bson.D{{"created_at", -1}})

	id, err := primitive.ObjectIDFromHex(occurenceId)

	if err != nil {
		return nil, e.Wrap(err)
	}

	search := bson.D{{"grouping_id", groupingId}, {"_id", id}}

	err = c.collection.FindOne(context.Background(), search, &queryOptions).Decode(&entry)

	if err != nil {
		return nil, e.Wrap(err)
	}

	return logEntryToDomain(&entry), nil*/
}

func (c eventsRepo) IsGroupExist(project *domain.Project, groupingId string) bool {
	/*search := bson.D{{"grouping_id", groupingId}}
	count, err := c.collection.CountDocuments(context.Background(), search)

	if err != nil {
		return false
	}

	return count > 0*/

	//TODO: implement

	return false
}

///////

func NewEventsRepo(db *gorm.DB) eventsRepo {
	return eventsRepo{
		db: db,
	}
}
