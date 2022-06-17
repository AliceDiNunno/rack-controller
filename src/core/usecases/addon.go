package usecases

import (
	"fmt"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/google/uuid"

	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
)

func (i interactor) GetAddons(service *domain.Service) ([]domain.Addon, *e.Error) {
	addons, err := i.addonRepository.GetAddons(service)
	if err != nil {
		return nil, err
	}
	return addons, nil
}

func (i interactor) GetAddonById(service *domain.Service, id uuid.UUID) (*domain.Addon, *e.Error) {
	addon, err := i.addonRepository.GetAddonById(service, id)
	if err != nil {
		return nil, err.Append(domain.ErrAddonNotFound)
	}
	return addon, nil
}

func (i interactor) createPostgresAddon(service *domain.Service, r *request.AddonCreationRequest) (*domain.Addon, *e.Error) {
	addon := domain.Addon{
		DisplayName: fmt.Sprint("%s-postgres", service.DisplayName),
		Type:        domain.AddonTypePostgres,
		ServiceID:   service.ID,
	}

	addon.Slug = slugify(addon.DisplayName)
	addon.Initialize()

	//return i.addonRepository.CreateAddon(&addon)
}

func (i interactor) CreateAddon(service *domain.Service, r *request.AddonCreationRequest) (*domain.Addon, *e.Error) {
	foundAddons, err := i.addonRepository.GetAddons(service)

	if err != nil {
		return nil, err.Append(domain.ErrUnableToCreateAddon)
	}

	if len(foundAddons) > 0 && r.Type == domain.AddonTypePostgres {
		return nil, e.Wrap(domain.ErrAServiceCanOnlyHaveOnePostgresAddon)
	}

	if r.Type == domain.AddonTypePostgres {
		return i.createPostgresAddon(service, r)
	}

	return nil, err.Append(domain.ErrUnknownAddonType)
}
