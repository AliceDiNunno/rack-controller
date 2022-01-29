package event

import (
	"errors"
	glc "github.com/AliceDiNunno/go-logger-client"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	request "github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	domain "github.com/AliceDiNunno/rack-controller/src/core/domain/eventDomain"
	"github.com/AliceDiNunno/rack-controller/src/core/usecases"
	"github.com/google/uuid"
	"time"
)

type InternalEventTransporter struct {
	usecases usecases.Usecases
}

func (i InternalEventTransporter) PushNewLogEntry(id uuid.UUID, rqst *glc.ItemCreationRequest) *e.Error {
	projectKey, err := uuid.Parse(rqst.ProjectKey)

	println("PUSHING LOG ENTRY")

	if err != nil {
		return e.Wrap(errors.New("Failed to parse project key"))
	}

	return i.usecases.PushNewLogEntry(id, &request.ItemCreationRequest{
		ProjectKey: projectKey,

		Identification: domain.LogIdentification{
			Client: domain.LogClientIdentification{
				UserID:    nil,
				IPAddress: "",
			},
			Deployment: domain.LogDeploymentIdentification{
				Environment: "",
				Version:     "",
			},
		},

		Data: domain.LogData{
			Message:   rqst.Data.Message,
			Timestamp: time.Now(),
		},
	})
}

func NewEventTransporter(usecases usecases.Usecases) *InternalEventTransporter {
	receiver := InternalEventTransporter{
		usecases: usecases,
	}

	return &receiver
}
