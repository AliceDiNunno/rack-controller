package domain

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LogEntry struct {
	//Object metadata
	ID        primitive.ObjectID `json:"id"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`

	//Identification (for reproduction)
	ProjectID      uuid.UUID         `json:"project_id"`
	Identification LogIdentification `json:"identification"`
	Data           LogData           `json:"data"`
}

type LogData struct {
	Timestamp        time.Time              `json:"timestamp"`
	GroupingID       string                 `json:"grouping_id"`
	Fingerprint      string                 `binding:"required,omitempty" json:"fingerprint"`
	Level            string                 `binding:"required,omitempty" json:"level"`
	Trace            Traceback              `json:"trace"`
	NestedTrace      []Traceback            `json:"nested_trace"`
	Message          string                 `binding:"required,omitempty" json:"message"`
	StatusCode       int                    `json:"status_code"`
	AdditionalFields map[string]interface{} `json:"additional_fields"`
}

type LogClientIdentification struct {
	UserID    *uuid.UUID `json:"user_id"`
	IPAddress string     `json:"ip_address"`
}

type LogDeploymentIdentification struct {
	Platform    string `json:"platform"`
	Source      string `binding:"required,omitempty" json:"source"`   //Source is either server or client
	Hostname    string `binding:"required,omitempty" json:"hostname"` //Hostname can be the name of the server or the client device
	Environment string `binding:"required,omitempty" json:"environment"`
	Version     string `binding:"required,omitempty" json:"version"`
}

type LogIdentification struct {
	Client     LogClientIdentification     `json:"client"`
	Deployment LogDeploymentIdentification `json:"deployment"`
}
