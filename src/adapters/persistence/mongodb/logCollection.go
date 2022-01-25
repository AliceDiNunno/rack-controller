package mongodb

import (
	"context"
	"fmt"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	logDomain "github.com/AliceDiNunno/rack-controller/src/core/domain/eventDomain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type logEntry struct {
	//Object metadata
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`

	//Entry metadata
	Timestamp   time.Time `bson:"timestamp"`
	GroupingID  string    `bson:"grouping_id"`
	Fingerprint string    `bson:"fingerprint"`

	//Identification (for reproduction)
	Platform    string                `bson:"platform"`
	Source      string                `bson:"source"`
	ProjectID   uuid.UUID             `bson:"project"`
	Hostname    string                `bson:"hostname"`
	Environment string                `bson:"environment"`
	Level       string                `bson:"level"`
	Version     string                `bson:"version"`
	Trace       logDomain.Traceback   `bson:"trace"`
	NestedTrace []logDomain.Traceback `bson:"nested_trace"`
	UserID      *uuid.UUID            `bson:"user_id"`
	IPAddress   string                `bson:"ip_address"`
	StatusCode  int                   `bson:"status_code"`

	//Entry Data
	Message          string                 `bson:"message"`
	AdditionalFields map[string]interface{} `bson:"additional"`
}

type SearchLogsFilter struct {
	ProjectID      *uuid.UUID
	UserID         *uuid.UUID
	Fingerprint    string
	GroupingID     string
	ServerHostname string
	Environment    string
	Level          string
	IPAddress      string
}

type logCollection struct {
	db         *mongo.Client
	collection *mongo.Collection
}

func logEntryToDomain(entry *logEntry) *logDomain.LogEntry {
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
			Timestamp:        entry.Timestamp,
			GroupingID:       entry.GroupingID,
			Fingerprint:      entry.Fingerprint,
			Level:            entry.Level,
			Trace:            entry.Trace,
			NestedTrace:      entry.NestedTrace,
			Message:          entry.Message,
			StatusCode:       entry.StatusCode,
			AdditionalFields: entry.AdditionalFields,
		},
	}
}

func logEntryFromDomain(entry *logDomain.LogEntry) *logEntry {
	return &logEntry{
		ID:               entry.ID,
		CreatedAt:        entry.CreatedAt,
		UpdatedAt:        entry.UpdatedAt,
		ProjectID:        entry.ProjectID,
		Timestamp:        entry.Data.Timestamp,
		GroupingID:       entry.Data.GroupingID,
		Fingerprint:      entry.Data.Fingerprint,
		Level:            entry.Data.Level,
		Trace:            entry.Data.Trace,
		AdditionalFields: entry.Data.AdditionalFields,
		NestedTrace:      entry.Data.NestedTrace,
		StatusCode:       entry.Data.StatusCode,
		Message:          entry.Data.Message,
		Platform:         entry.Identification.Deployment.Platform,
		Source:           entry.Identification.Deployment.Source,
		Hostname:         entry.Identification.Deployment.Hostname,
		Environment:      entry.Identification.Deployment.Environment,
		Version:          entry.Identification.Deployment.Version,
		UserID:           entry.Identification.Client.UserID,
		IPAddress:        entry.Identification.Client.IPAddress,
	}
}

func (c logCollection) AddLog(entry *logDomain.LogEntry) *e.Error {
	entryFromDomain := logEntryFromDomain(entry)

	_, err := c.collection.InsertOne(context.Background(), entryFromDomain)

	return e.Wrap(err)
}

func interfaceArrayToStringArray(input []interface{}) []string {
	output := make([]string, len(input))
	for i, v := range input {
		output[i] = v.(string)
	}
	return output
}

func (c logCollection) ProjectVersions(project *domain.Project) ([]string, *e.Error) {
	result, err := c.collection.Distinct(context.Background(), "version", bson.M{"project": project.ID})

	if err != nil {
		return nil, e.Wrap(err)
	}

	return interfaceArrayToStringArray(result), nil
}

func (c logCollection) ProjectEnvironments(project *domain.Project) ([]string, *e.Error) {
	result, err := c.collection.Distinct(context.Background(), "environment", bson.M{"project": project.ID})

	if err != nil {
		return nil, e.Wrap(err)
	}

	return interfaceArrayToStringArray(result), nil
}

func (c logCollection) ProjectServers(project *domain.Project) ([]string, *e.Error) {
	result, err := c.collection.Distinct(context.Background(), "hostname", bson.M{"project": project.ID})

	if err != nil {
		return nil, e.Wrap(err)
	}

	return interfaceArrayToStringArray(result), nil
}

func (c logCollection) ProjectGroupingIds(project *domain.Project) ([]string, *e.Error) {
	result, err := c.collection.Distinct(context.Background(), "grouping_id", bson.M{"project": project.ID})

	if err != nil {
		return nil, e.Wrap(err)
	}

	return interfaceArrayToStringArray(result), nil
}

func (c logCollection) FindLastEntryForGroup(project *domain.Project, groupingId string) (*logDomain.LogEntry, *e.Error) {
	var entry logEntry

	queryOptions := options.FindOneOptions{}

	queryOptions.SetSort(bson.D{{"created_at", -1}})

	err := c.collection.FindOne(context.Background(), bson.D{{"grouping_id", groupingId}}, &queryOptions).Decode(&entry)

	if err != nil {
		return nil, e.Wrap(err)
	}

	return logEntryToDomain(&entry), nil
}

func (c logCollection) FindGroupOccurrences(project *domain.Project, groupingId string) ([]string, *e.Error) {
	var entries []string

	queryOptions := options.FindOptions{}

	queryOptions.SetSort(bson.D{{"created_at", -1}})

	cur, err := c.collection.Find(context.Background(), bson.D{{"grouping_id", groupingId}}, &queryOptions)

	if err != nil {
		return nil, e.Wrap(err)
	}

	for cur.Next(context.Background()) {
		var elem logEntry
		err := cur.Decode(&elem)
		if err == nil {
			entries = append(entries, fmt.Sprintf("%q", elem.ID.Hex()))
		}
	}

	return entries, nil
}
func (c logCollection) FindGroupOccurrence(project *domain.Project, groupingId string, occurenceId string) (*logDomain.LogEntry, *e.Error) {
	var entry logEntry

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

	return logEntryToDomain(&entry), nil
}

func (c logCollection) IsGroupExist(project *domain.Project, groupingId string) bool {
	search := bson.D{{"grouping_id", groupingId}}
	count, err := c.collection.CountDocuments(context.Background(), search)

	if err != nil {
		return false
	}

	return count > 0
}

func NewLogCollectionRepo(db *mongo.Client) logCollection {
	collection := db.Database("logger").Collection("logs")

	return logCollection{
		db:         db,
		collection: collection,
	}
}
