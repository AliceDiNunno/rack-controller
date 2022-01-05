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
	Project     Project
	ProjectId   uuid.UUID
}

func NewEnvironmentRepo(db *gorm.DB) environmentRepo {
	return environmentRepo{
		db: db,
	}
}
