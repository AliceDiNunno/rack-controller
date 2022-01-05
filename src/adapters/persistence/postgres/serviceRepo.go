package postgres

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	gorm.Model

	ID          uuid.UUID
	DisplayName string
	ImageName   string
}
