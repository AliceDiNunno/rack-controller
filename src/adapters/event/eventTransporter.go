package event

import (
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
	print("PUSH NEW LOG ENTRY")

	return i.usecases.PushNewEvent(id, &request.ItemCreationRequest{
		ProjectKey: rqst.ProjectKey,
		Identification: request.LogIdentification{
			Client: request.LogClientIdentification{
				UserID:    nil, //TODO
				IPAddress: rqst.Identification.Client.IPAddress,
			},
			Deployment: request.LogDeploymentIdentification{
				Platform:    rqst.Identification.Deployment.Platform,
				Source:      rqst.Identification.Deployment.Source,
				Hostname:    rqst.Identification.Deployment.Hostname,
				Environment: rqst.Identification.Deployment.Environment,
				Version:     rqst.Identification.Deployment.Version,
			},
		},
		Data: request.LogData{
			Timestamp:        rqst.Data.Timestamp,
			GroupingID:       rqst.Data.GroupingID,
			Fingerprint:      rqst.Data.Fingerprint,
			Level:            rqst.Data.Level,
			Trace:            nil, //TODO
			NestedTrace:      nil, //TODO
			Message:          rqst.Data.Message,
			Module:           rqst.Data.Module,
			StatusCode:       rqst.Data.StatusCode,
			AdditionalFields: nil, //TODO
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
