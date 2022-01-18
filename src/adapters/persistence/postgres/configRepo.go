package postgres

import (
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type configRepo struct {
	db *gorm.DB
}

type Config struct {
	gorm.Model
	ClusterModel

	ID        uuid.UUID
	Name      string
	Value     string
	Project   Project
	ProjectId uuid.UUID
}

func environmentsToDomain(project []Config) []domain.Environment {
	environmentSlice := []domain.Environment{}

	for _, p := range project {
		environmentSlice = append(environmentSlice, environmentToDomain(p))
	}

	return environmentSlice
}

func environmentToDomain(environment Config) domain.Environment {
	return domain.Environment{
		ID:          environment.ID,
		DisplayName: environment.DisplayName,
		ProjectId:   environment.ProjectId,
		Slug:        environment.Slug,
	}
}

func NewEnvironmentRepo(db *gorm.DB) configRepo {
	return configRepo{
		db: db,
	}
}
