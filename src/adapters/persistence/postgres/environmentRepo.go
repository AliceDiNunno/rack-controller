package postgres

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type environmentRepo struct {
	db *gorm.DB
}

type Environment struct {
	gorm.Model
	ClusterModel

	ID          uuid.UUID
	DisplayName string
}
