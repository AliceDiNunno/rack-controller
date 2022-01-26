package domain

import "github.com/google/uuid"

type Project struct {
	ID          uuid.UUID
	DisplayName string
	Slug        string

	EventKey uuid.UUID

	UserID uuid.UUID
}

func (p *Project) Initialize() {
	p.ID = uuid.New()
}
