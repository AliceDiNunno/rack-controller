package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
)

func (i interactor) GetServices(project *domain.Project) ([]domain.Service, *e.Error) {
	services, err := i.serviceRepository.GetServices(project.ID)

	if err != nil {
		return nil, err
	}

	return services, nil
}

func (i interactor) GetServiceById(project *domain.Project, id uuid.UUID) (*domain.Service, *e.Error) {
	service, err := i.serviceRepository.GetServiceById(project.ID, id)

	if err != nil {
		return nil, err
	}

	return service, nil
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

	var service *domain.Service

	service, err := i.serviceRepository.GetServiceByName(project.ID, r.Name)

	if err != nil || service == nil {
		service = &domain.Service{
			DisplayName: r.Name,
			ProjectID:   project.ID,
			Slug:        slugify(r.Name),
		}

		service.Initialize()
	}

	environments, err := i.environmentRepository.GetEnvironments(project.ID)

	if err != nil {
		return e.Wrap(domain.ErrUnableToCreateProject)
	}

	spew.Dump(service)

	newDeployment := clusterDomain.DeploymentCreationRequest{
		DeploymentName: service.Slug,
		ImageName:      r.ImageName,
		Ports:          nil,
		Environment:    nil,
		ConfigMaps:     nil,
		Secrets:        nil,
		Replicas:       4,
		Memory:         0,
		CPU:            0,
	}

	for _, env := range environments {
		config := []clusterDomain.Environment{}

		projectConfig := i.ConfigForProject(project)
		for k, v := range projectConfig {
			config = append(config, clusterDomain.Environment{
				Name:  k,
				Value: v,
			})
		}

		environmentConfig := i.ConfigForEnvironment(&env)
		for k, v := range environmentConfig {
			config = append(config, clusterDomain.Environment{
				Name:  k,
				Value: v,
			})
		}

		serviceConfig := i.ConfigForService(service)
		for k, v := range serviceConfig {
			config = append(config, clusterDomain.Environment{
				Name:  k,
				Value: v,
			})
		}

		newDeployment.Environment = config

		err := i.kubeClient.CreateDeployment(env.Slug, newDeployment)
		spew.Dump(err)
	}

	return i.serviceRepository.CreateOrUpdateService(service)
}

func (i interactor) GetServiceConfig(service *domain.Service) ([]clusterDomain.Environment, *e.Error) {
	if service == nil {
		return nil, e.Wrap(domain.ErrProjectNotFound)
	}

	config, err := i.configRepository.GetConfigByObjectID(service.ID)

	if err != nil {
		return nil, err.Append(domain.UnableToGetConfig)
	}

	return config, nil
}

func (i interactor) UpdateServiceConfig(service *domain.Service, envVariables []clusterDomain.Environment) *e.Error {
	if service == nil {
		return e.Wrap(domain.ErrProjectNotFound)
	}

	err := i.configRepository.SetConfig(service.ID, envVariables)

	if err != nil {
		return err.Append(domain.UnableToUpdateConfig)
	}

	return nil
}

func (i interactor) ConfigForService(service *domain.Service) map[string]string {
	if service == nil {
		return nil
	}

	config := map[string]string{
		"SERVICE_NAME": service.DisplayName,
		"SERVICE_SLUG": service.Slug,
		"API_PREFIX":   service.Slug,
		"DB_NAME":      service.Slug,
	}

	userConfig, err := i.configRepository.GetConfigByObjectID(service.ID)

	if err != nil {
		return config
	}

	for _, env := range userConfig {
		config[env.Name] = env.Value
	}

	return config
}
