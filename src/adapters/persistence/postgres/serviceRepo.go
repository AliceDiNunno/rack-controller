package postgres

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
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
	ProjectID   uuid.UUID
}

func (s serviceRepo) GetServices(projectID uuid.UUID) ([]domain.Service, *e.Error) {
	var services []Service

	if err := s.db.Where("project_id = ?", projectID).Find(&services).Error; err != nil {
		return nil, e.Wrap(err)
	}

	return servicesToDomain(services), nil
}

func servicesToDomain(services []Service) []domain.Service {
	servicesSlice := []domain.Service{}

	for _, p := range services {
		servicesSlice = append(servicesSlice, serviceToDomain(p))
	}

	return servicesSlice
}

func serviceToDomain(service Service) domain.Service {
	return domain.Service{
		ID:          service.ID,
		DisplayName: service.DisplayName,
		ProjectID:   service.ProjectID,
	}
}

func NewServiceRepo(db *gorm.DB) serviceRepo {
	return serviceRepo{
		db: db,
	}
}
