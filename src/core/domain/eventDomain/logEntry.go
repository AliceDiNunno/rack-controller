package domain

import (
	"github.com/google/uuid"
	"time"
)

type LogEntry struct {
	//Object metadata
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time

	//Identification (for reproduction)
	ProjectID      uuid.UUID
	Identification LogIdentification
	Data           LogData
}

type LogData struct {
	Timestamp        time.Time
	GroupingID       string
	Fingerprint      string `binding:"required,omitempty"`
	Level            string `binding:"required,omitempty"`
	Trace            Traceback
	NestedTrace      []Traceback
	Message          string `binding:"required,omitempty"`
	StatusCode       int
	AdditionalFields map[string]interface{}
}

type LogClientIdentification struct {
	UserID    *uuid.UUID
	IPAddress string
}

type LogDeploymentIdentification struct {
	Platform    string
	Source      string `binding:"required,omitempty"` //Source is either server or client
	Hostname    string `binding:"required,omitempty"` //Hostname can be the name of the server or the client device
	Environment string `binding:"required,omitempty"`
	Version     string `binding:"required,omitempty"`
}

type LogIdentification struct {
	Client     LogClientIdentification
	Deployment LogDeploymentIdentification
}
