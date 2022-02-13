package domain

import "github.com/google/uuid"

type Service struct {
	ID uuid.UUID

	DisplayName string
	ImageName   string
	Slug        string

	ProjectID uuid.UUID
}

type ServiceDetail struct {
	Service

	RequestedInstances int
	RunningInstances   int
}

func (s *Service) Initialize() {
	s.ID = uuid.New()
}
