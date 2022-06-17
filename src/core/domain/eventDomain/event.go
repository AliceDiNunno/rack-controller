package domain

import (
	"github.com/google/uuid"
	"time"
)

type EventOccurrence struct {
	//Object metadata
	ID        uuid.UUID
	Timestamp time.Time
	ProjectID uuid.UUID

	EventID uuid.UUID

	Platform    string
	Source      string //Source is either server or client
	Hostname    string //Hostname can be the name of the server or the client device
	Environment string
	Version     string

	UserID    *uuid.UUID
	IPAddress string
}

type Event struct {
	//Object metadata
	ID        uuid.UUID
	Timestamp time.Time
	ProjectID uuid.UUID

	GroupingID  string
	Fingerprint string

	Level      string
	Module     string
	StatusCode int

	Message string

	Occurrences int64
}
