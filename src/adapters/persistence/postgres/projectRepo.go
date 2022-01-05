package postgres

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type projectRepo struct {
	db *gorm.DB
}

type Project struct {
	gorm.Model
	ClusterModel

	ID           uuid.UUID
	DisplayName  string
	Environments []Environment
	Services     []Service
}

func NewProjectRepo(db *gorm.DB) projectRepo {
	return projectRepo{
		db: db,
	}
}
