package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/davecgh/go-spew/spew"
)

func (i interactor) GetServices(project *domain.Project) ([]domain.Service, *e.Error) {
	services, err := i.serviceRepository.GetServices(project.ID)

	if err != nil {
		return nil, err
	}

	return services, nil
}

func (i interactor) CreateService(project *domain.Project, r *request.ServiceCreationRequest) *e.Error {
	if project == nil {
		return e.Wrap(domain.ErrProjectNotFound)
	}

	if r == nil {
		return e.Wrap(domain.ErrInvalidRequest)
	}

	if r.Name == "" {
		return e.Wrap(domain.ErrServiceNameIsEmpty)
	}

	env, err := i.serviceRepository.GetServiceByName(project.ID, r.Name)

	spew.Dump(env, err)

	if err == nil && env != nil {
		return e.Wrap(domain.ErrServiceAlreadyExistsWithThisName)
	}

	service := domain.Service{
		DisplayName: r.Name,
		ProjectID:   project.ID,
	}

	service.Initialize()

	i.dispatcher.Dispatch("service.created", service)

	return i.serviceRepository.CreateService(&service)
}
