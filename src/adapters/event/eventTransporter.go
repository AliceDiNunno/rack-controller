package event

import (
	"errors"
	glc "github.com/AliceDiNunno/go-logger-client"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	request "github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	domain "github.com/AliceDiNunno/rack-controller/src/core/domain/eventDomain"
	"github.com/AliceDiNunno/rack-controller/src/core/usecases"
	"github.com/google/uuid"
)

type InternalEventTransporter struct {
	usecases usecases.Usecases
}

func (i InternalEventTransporter) PushNewLogEntry(id uuid.UUID, rqst *glc.ItemCreationRequest) *e.Error {
	projectKey, err := uuid.Parse(rqst.ProjectKey)

	if err != nil {
		return e.Wrap(errors.New("Failed to parse project key"))
	}

	return i.usecases.PushNewLogEntry(id, &request.ItemCreationRequest{
		ProjectKey: projectKey,

		Identification: domain.LogIdentification{
			Client: domain.LogClientIdentification{
				UserID:    rqst.Identification.Client.UserID,
				IPAddress: rqst.Identification.Client.IPAddress,
			},
			Deployment: domain.LogDeploymentIdentification{
				Platform:    rqst.Identification.Deployment.Platform,
				Source:      rqst.Identification.Deployment.Source,
				Hostname:    rqst.Identification.Deployment.Hostname,
				Environment: rqst.Identification.Deployment.Environment,
				Version:     rqst.Identification.Deployment.Version,
			},
		},
		Data: domain.LogData{
			Timestamp:        rqst.Data.Timestamp,
			GroupingID:       rqst.Data.GroupingID,
			Fingerprint:      rqst.Data.Fingerprint,
			Level:            rqst.Data.Level,
			Trace:            tracebackFromClient(rqst.Data.Trace),
			NestedTrace:      tracebacksFromClient(rqst.Data.NestedTrace),
			Message:          rqst.Data.Message,
			StatusCode:       rqst.Data.StatusCode,
			AdditionalFields: rqst.Data.AdditionalFields,
		},
	})
}

func tracebackFromClient(traceback *glc.Traceback) domain.Traceback {
	var frames []domain.TracebackEntry

	for _, frame := range traceback.Traceback {
		frames = append(frames, domain.TracebackEntry{
			Filename: frame.Filename,
			Method:   frame.Method,
			Line:     frame.Line,
		})
	}

	return domain.Traceback{
		Message:   traceback.Message,
		Traceback: frames,
	}
}

func tracebacksFromClient(tracebacks []*glc.Traceback) []domain.Traceback {
	var result []domain.Traceback

	for _, traceback := range tracebacks {
		result = append(result, tracebackFromClient(traceback))
	}

	return result
}

func NewEventTransporter(usecases usecases.Usecases) *InternalEventTransporter {
	return &InternalEventTransporter{
		usecases: usecases,
	}
}
