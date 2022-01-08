package postgres

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
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

func (r environmentRepo) GetEnvironments(projectID uuid.UUID) ([]domain.Environment, *e.Error) {
	var environments []Environment

	if err := r.db.Find(&environments).Error; err != nil {
		return nil, e.Wrap(err)
	}

	return environmentsToDomain(environments), nil

}

func (r environmentRepo) GetEnvironmentByProjectId(project uuid.UUID) ([]domain.Environment, *e.Error) {
	var environments []Environment

	if err := r.db.Where("project_id = ?", project).Find(&environments).Error; err != nil {
		return nil, e.Wrap(err)
	}

	return environmentsToDomain(environments), nil
}

func environmentsToDomain(project []Environment) []domain.Environment {
	environmentSlice := []domain.Environment{}

	for _, p := range project {
		environmentSlice = append(environmentSlice, environmentToDomain(p))
	}

	return environmentSlice
}

func environmentToDomain(project Environment) domain.Environment {
	return domain.Environment{
		ID:          project.ID,
		DisplayName: project.DisplayName,
		ProjectId:   project.ProjectId,
	}
}

func NewEnvironmentRepo(db *gorm.DB) environmentRepo {
	return environmentRepo{
		db: db,
	}
}
