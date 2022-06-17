package postgres

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	eventDomain "github.com/AliceDiNunno/rack-controller/src/core/domain/eventDomain"
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

	Message   string
	Traceback []TracebackEntry
	Event     *Event
	EventID   uuid.UUID
}

type EventOccurrence struct {
	gorm.Model

	ID        uuid.UUID
	Timestamp time.Time
	Event     Event
	EventID   uuid.UUID

	Platform    string
	Source      string //Source is either server or client
	Hostname    string //Hostname can be the name of the server or the client device
	Environment string
	Version     string

	UserID    *uuid.UUID
	IPAddress string
}

type Event struct {
	gorm.Model

	//Object metadata
	ID        uuid.UUID
	Timestamp time.Time
	Project   Project
	ProjectID uuid.UUID

	//Event metadata
	GroupingID  string
	Fingerprint string

	//Identification (for reproduction)
	Trace   *Traceback
	TraceID uuid.UUID

	Occurrences []EventOccurrence

	NestedTrace []Traceback

	Level      string
	Module     string
	StatusCode int

	//Entry Data
	Message string
}

func tracebackToDomain(traceback *Traceback) eventDomain.Traceback {
	if traceback == nil {
		return eventDomain.Traceback{}
	}

	var tracebackEntries []eventDomain.TracebackEntry

	for _, entry := range traceback.Traceback {
		tracebackEntries = append(tracebackEntries, eventDomain.TracebackEntry{
			Filename: entry.Filename,
			Line:     entry.Line,
			Method:   entry.Method,
		})
	}

	return eventDomain.Traceback{
		Message:   traceback.Message,
		Traceback: tracebackEntries,
	}
}

func tracebackFromDomain(traceback eventDomain.Traceback) Traceback {
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

func tracebacksToDomain(traceback []Traceback) []eventDomain.Traceback {
	var tracebacks []eventDomain.Traceback

	for _, tb := range traceback {
		tracebacks = append(tracebacks, tracebackToDomain(&tb))
	}

	return tracebacks
}

func tracebacksFromDomain(traceback []eventDomain.Traceback) []Traceback {
	var tracebacks []Traceback

	for _, tb := range traceback {
		tracebacks = append(tracebacks, tracebackFromDomain(tb))
	}

	return tracebacks
}

func entriesToDomain(entries []Event) []eventDomain.Event {
	var eventEntries []eventDomain.Event

	for _, entry := range entries {
		eventEntries = append(eventEntries, *eventEntryToDomain(&entry, 0))
	}

	return eventEntries
}

func eventEntryFromDomain(entry *eventDomain.Event) *Event {
	return &Event{
		ID:          entry.ID,
		Timestamp:   entry.Timestamp,
		ProjectID:   entry.ProjectID,
		GroupingID:  entry.GroupingID,
		Fingerprint: entry.Fingerprint,
		Message:     entry.Message,
		Level:       entry.Level,
		StatusCode:  entry.StatusCode,
		Module:      entry.Module,
	}
}

func eventEntryToDomain(entry *Event, occurrences int64) *eventDomain.Event {
	return &eventDomain.Event{
		ID:          entry.ID,
		Timestamp:   entry.Timestamp,
		ProjectID:   entry.ProjectID,
		GroupingID:  entry.GroupingID,
		Fingerprint: entry.Fingerprint,
		Message:     entry.Message,
		Level:       entry.Level,
		StatusCode:  entry.StatusCode,
		Module:      entry.Module,
		Occurrences: occurrences,
	}
}

func eventOccurrenceFromDomain(entry *eventDomain.EventOccurrence) *EventOccurrence {
	return &EventOccurrence{
		ID:          entry.ID,
		Timestamp:   entry.Timestamp,
		EventID:     entry.EventID,
		Platform:    entry.Platform,
		Source:      entry.Source,
		Hostname:    entry.Hostname,
		Environment: entry.Environment,
		Version:     entry.Version,
		UserID:      entry.UserID,
		IPAddress:   entry.IPAddress,
	}
}

////////

func (c eventsRepo) AddEvent(entry *eventDomain.Event) *e.Error {
	entryFromDomain := eventEntryFromDomain(entry)

	if err := c.db.Create(entryFromDomain).Error; err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (c eventsRepo) AddOccurrence(occurrence *eventDomain.EventOccurrence) *e.Error {
	entryFromDomain := eventOccurrenceFromDomain(occurrence)

	if err := c.db.Create(entryFromDomain).Error; err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (c eventsRepo) ProjectVersions(project *domain.Project) ([]eventDomain.Event, *e.Error) {
	var eventEntry []eventDomain.Event

	err := c.db.Distinct("version").Where("project_id = ?", project.ID).Find(&eventEntry).Error

	if err != nil {
		return nil, e.Wrap(err)
	}

	return eventEntry, nil
}

func (c eventsRepo) ProjectEnvironments(project *domain.Project) ([]eventDomain.Event, *e.Error) {
	var eventEntry []eventDomain.Event

	err := c.db.Distinct("environment").Where("project_id = ?", project.ID).Find(&eventEntry).Error

	if err != nil {
		return nil, e.Wrap(err)
	}

	return eventEntry, nil
}

func (c eventsRepo) ProjectServers(project *domain.Project) ([]eventDomain.Event, *e.Error) {
	var eventEntry []eventDomain.Event

	err := c.db.Distinct("environment").Where("project_id = ?", project.ID).Find(&eventEntry).Error

	if err != nil {
		return nil, e.Wrap(err)
	}

	return eventEntry, nil
}

func (c eventsRepo) ProjectGroupingIds(project *domain.Project) ([]eventDomain.Event, *e.Error) {
	var eventEntry []Event

	err := c.db.Distinct("grouping_id").Where("project_id = ?", project.ID).Find(&eventEntry).Error

	if err != nil {
		return nil, e.Wrap(err)
	}

	return entriesToDomain(eventEntry), nil
}

func (c eventsRepo) FindLastEntryForGroup(project *domain.Project, groupingId string) (*eventDomain.Event, *e.Error) {
	var eventEntry *Event
	var count int64

	data := c.db.Where("grouping_id = ?", groupingId).Order("created_at desc").Last(&eventEntry)

	data.Count(&count)

	if err := data.Error; err != nil {
		spew.Dump(err)
		return nil, e.Wrap(err)
	}

	return eventEntryToDomain(eventEntry, count), nil
}

func (c eventsRepo) FindGroupOccurrences(project *domain.Project, groupingId string) ([]eventDomain.Event, *e.Error) {
	var entries []eventDomain.Event

	data := c.db.Where("grouping_id = ?", groupingId).Order("created_at desc").Last(&entries)

	if err := data.Error; err != nil {
		spew.Dump(err)
		return nil, e.Wrap(err)
	}

	return entries, nil
}

func (c eventsRepo) FindGroupOccurrence(project *domain.Project, groupingId string, occurenceId string) (*eventDomain.Event, *e.Error) {
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

	return eventEntryToDomain(&entry), nil*/
}

func (c eventsRepo) IsGroupExist(project *domain.Project, groupingId string) bool {
	entry, err := c.FindLastEntryForGroup(project, groupingId)

	return entry != nil || err == nil
}

///////

func NewEventsRepo(db *gorm.DB) eventsRepo {
	return eventsRepo{
		db: db,
	}
}
