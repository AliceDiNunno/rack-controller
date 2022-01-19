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
	Slug        string
	Project     Project
	ProjectId   uuid.UUID
}

func (r environmentRepo) CreateEnvironment(environment *domain.Environment) *e.Error {
	if err := r.db.Create(environment).Error; err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (r environmentRepo) GetEnvironments(projectID uuid.UUID) ([]domain.Environment, *e.Error) {
	var environments []Environment

	if err := r.db.Where("project_id = ?", projectID).Find(&environments).Error; err != nil {
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

func (r environmentRepo) GetEnvironmentByName(projectID uuid.UUID, name string) (*domain.Environment, *e.Error) {
	var environment Environment

	if err := r.db.Where("project_id = ? AND display_name = ?", projectID, name).First(&environment).Error; err != nil {
		return nil, e.Wrap(err)
	}

	envToReturn := environmentToDomain(environment)

	return &envToReturn, nil
}

func (r environmentRepo) GetEnvironmentByID(projectID uuid.UUID, ID uuid.UUID) (*domain.Environment, *e.Error) {
	var environment Environment

	if err := r.db.Where("project_id = ? AND id = ?", projectID, ID).First(&environment).Error; err != nil {
		return nil, e.Wrap(err)
	}

	envToReturn := environmentToDomain(environment)

	return &envToReturn, nil
}

func environmentsToDomain(environment []Environment) []domain.Environment {
	environmentSlice := []domain.Environment{}

	for _, p := range environment {
		environmentSlice = append(environmentSlice, environmentToDomain(p))
	}

	return environmentSlice
}

func environmentToDomain(environment Environment) domain.Environment {
	return domain.Environment{
		ID:          environment.ID,
		DisplayName: environment.DisplayName,
		ProjectId:   environment.ProjectId,
		Slug:        environment.Slug,
	}
}

func NewEnvironmentRepo(db *gorm.DB) environmentRepo {
	return environmentRepo{
		db: db,
	}
}
