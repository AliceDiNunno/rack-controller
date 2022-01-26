package request

import (
	domain "github.com/AliceDiNunno/rack-controller/src/core/domain/eventDomain"
	"github.com/google/uuid"
)

type ItemCreationRequest struct {
	ProjectKey uuid.UUID `binding:"required" json:"project_key"`

	Identification domain.LogIdentification `json:"identification"`
	Data           domain.LogData           `json:"data"`
}
