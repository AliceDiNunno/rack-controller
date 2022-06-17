package domain

import "github.com/google/uuid"

const AddonTypePostgres = 0

type Addon struct {
	ID uuid.UUID

	DisplayName string
	Type        int
	Slug        string

	ServiceID uuid.UUID
}

func (a *Addon) Initialize() {
	a.ID = uuid.New()
}
