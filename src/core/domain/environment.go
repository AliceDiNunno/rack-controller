package domain

import "github.com/google/uuid"

type Environment struct {
	ID          uuid.UUID
	DisplayName string

	ProjectId uuid.UUID
}

func (e *Environment) Initialize() {
	e.ID = uuid.New()
}
