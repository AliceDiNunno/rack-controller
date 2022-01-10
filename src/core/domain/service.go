package domain

import "github.com/google/uuid"

type Service struct {
	ID uuid.UUID

	DisplayName string
	ImageName   string

	ProjectID uuid.UUID
}
