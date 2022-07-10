package postgres

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type addonRepo struct {
	db *gorm.DB
}

type Addon struct {
	gorm.Model

	ID          uuid.UUID
	DisplayName string
	Slug        string
	Type        int
	Service     Service
	ServiceID   uuid.UUID
}

func (s addonRepo) GetAddons(service *domain.Service) ([]domain.Addon, *e.Error) {
	var addons []Addon

	if err := s.db.Where("service_id = ?", service.ID).Find(&addons).Error; err != nil {
		return nil, e.Wrap(err)
	}

	return addonsToDomain(addons), nil
}

func (s addonRepo) CreateAddon(d *domain.Addon) (*domain.Addon, *e.Error) {
	return s.CreateOrUpdateAddon(d)
}

func (s addonRepo) CreateOrUpdateAddon(addon *domain.Addon) (*domain.Addon, *e.Error) {
	addonToSave := addonFromDomain(*addon)

	if err := s.db.Where("service_id = ? AND display_name = ?", addon.ServiceID, addon.DisplayName).First(&addonToSave).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, e.Wrap(err)
		}

		if err := s.db.Create(&addonToSave).Error; err != nil {
			return nil, e.Wrap(err)
		}
	} else {
		if err := s.db.Save(&addonToSave).Error; err != nil {
			return nil, e.Wrap(err)
		}
	}

	return addon, nil
}

func (s addonRepo) UpdateAddon(addon *domain.Addon) *e.Error {
	if err := s.db.Save(addon).Error; err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (s addonRepo) GetAddonById(service *domain.Service, ID uuid.UUID) (*domain.Addon, *e.Error) {
	var addon Addon

	if err := s.db.Where("service_id = ? AND id = ?", service.ID, ID).First(&addon).Error; err != nil {
		return nil, e.Wrap(err)
	}

	addonToReturn := addonToDomain(addon)

	return &addonToReturn, nil
}

func (s addonRepo) DeleteAddon(addon *domain.Addon) *e.Error {
	addonToDelete := addonFromDomain(*addon)

	if err := s.db.Delete(&addonToDelete).Error; err != nil {
		return e.Wrap(err)
	}

	return nil
}

func addonsToDomain(addons []Addon) []domain.Addon {
	addonsSlice := []domain.Addon{}

	for _, p := range addons {
		addonsSlice = append(addonsSlice, addonToDomain(p))
	}

	return addonsSlice
}

func addonFromDomain(addon domain.Addon) Addon {
	return Addon{
		ID:          addon.ID,
		DisplayName: addon.DisplayName,
		Type:        addon.Type,
		Slug:        addon.Slug,
		ServiceID:   addon.ServiceID,
	}
}

func addonToDomain(addon Addon) domain.Addon {
	return domain.Addon{
		ID:          addon.ID,
		DisplayName: addon.DisplayName,
		Type:        addon.Type,
		Slug:        addon.Slug,
		ServiceID:   addon.ServiceID,
	}
}

func NewAddonRepo(db *gorm.DB) addonRepo {
	return addonRepo{
		db: db,
	}
}
