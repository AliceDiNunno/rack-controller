package request

import (
	"github.com/AliceDiNunno/go-logger/src/core/domain"
	"github.com/google/uuid"
)

type ItemCreationRequest struct {
	ProjectKey uuid.UUID `binding:"required" json:"project_key"`

	Identification domain.LogIdentification `json:"identification"`
	Data           domain.LogData           `json:"data"`
}
