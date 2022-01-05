package postgres

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type serviceRepo struct {
	db *gorm.DB
}

type Service struct {
	gorm.Model

	ID          uuid.UUID
	DisplayName string
	ImageName   string
	Project     Project
	ProjectId   uuid.UUID
}

func NewServiceRepo(db *gorm.DB) serviceRepo {
	return serviceRepo{
		db: db,
	}
}
