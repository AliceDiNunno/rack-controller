package domain

import "github.com/google/uuid"

type Service struct {
	ID uuid.UUID

	DisplayName string
	ImageName   string
	Slug        string

	ProjectID uuid.UUID
}

func (s *Service) Initialize() {
	s.ID = uuid.New()
}
