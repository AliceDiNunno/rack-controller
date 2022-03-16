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
	Slug        string
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

func (s serviceRepo) GetServiceByName(projectID uuid.UUID, name string) (*domain.Service, *e.Error) {
	var service Service

	if err := s.db.Where("project_id = ? AND display_name = ?", projectID, name).First(&service).Error; err != nil {
		return nil, e.Wrap(err)
	}

	serviceToReturn := serviceToDomain(service)

	return &serviceToReturn, nil
}

func (s serviceRepo) CreateOrUpdateService(service *domain.Service) *e.Error {
	serviceToSave := serviceFromDomain(*service)

	if err := s.db.Where("project_id = ? AND display_name = ?", service.ProjectID, service.DisplayName).First(&serviceToSave).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return e.Wrap(err)
		}

		if err := s.db.Create(&serviceToSave).Error; err != nil {
			return e.Wrap(err)
		}
	} else {
		if err := s.db.Save(&serviceToSave).Error; err != nil {
			return e.Wrap(err)
		}
	}

	return nil
}

func (s serviceRepo) UpdateService(service *domain.Service) *e.Error {
	if err := s.db.Save(service).Error; err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (s serviceRepo) GetServiceById(projectID uuid.UUID, ID uuid.UUID) (*domain.Service, *e.Error) {
	var service Service

	if err := s.db.Where("project_id = ? AND id = ?", projectID, ID).First(&service).Error; err != nil {
		return nil, e.Wrap(err)
	}

	serviceToReturn := serviceToDomain(service)

	return &serviceToReturn, nil
}

func (s serviceRepo) DeleteService(service *domain.Service) *e.Error {
	serviceToDelete := serviceFromDomain(*service)

	if err := s.db.Delete(&serviceToDelete).Error; err != nil {
		return e.Wrap(err)
	}

	return nil
}

func servicesToDomain(services []Service) []domain.Service {
	servicesSlice := []domain.Service{}

	for _, p := range services {
		servicesSlice = append(servicesSlice, serviceToDomain(p))
	}

	return servicesSlice
}

func serviceFromDomain(service domain.Service) Service {
	return Service{
		ID:          service.ID,
		DisplayName: service.DisplayName,
		ImageName:   service.ImageName,
		Slug:        service.Slug,
		ProjectID:   service.ProjectID,
	}
}

func serviceToDomain(service Service) domain.Service {
	return domain.Service{
		ID:          service.ID,
		DisplayName: service.DisplayName,
		ImageName:   service.ImageName,
		Slug:        service.Slug,
		ProjectID:   service.ProjectID,
	}
}

func NewServiceRepo(db *gorm.DB) serviceRepo {
	return serviceRepo{
		db: db,
	}
}
