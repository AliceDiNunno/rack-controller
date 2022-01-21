package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/rest/request"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
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

	service := domain.Service{
		DisplayName: r.Name,
		ProjectID:   project.ID,
		Slug:        slugify(r.Name),
	}

	service.Initialize()

	environments, err := i.environmentRepository.GetEnvironments(project.ID)

	if err != nil {
		return e.Wrap(domain.ErrUnableToCreateProject)
	}

	newDeployment := clusterDomain.DeploymentCreationRequest{
		DeploymentName: service.Slug,
		ImageName:      r.ImageName,
		Ports:          nil,
		Environment:    nil,
		ConfigMaps:     nil,
		Secrets:        nil,
		Replicas:       3,
		Memory:         0,
		CPU:            0,
	}

	for _, env := range environments {

		envVars := []clusterDomain.Environment{}

		projectEnvVars := i.EnvVariablesForProject(project)
		for k, v := range projectEnvVars {
			envVars = append(envVars, clusterDomain.Environment{
				Name:  k,
				Value: v,
			})
		}

		environmentEnvVars := i.EnvVariablesForEnvironment(&env)
		for k, v := range environmentEnvVars {
			envVars = append(envVars, clusterDomain.Environment{
				Name:  k,
				Value: v,
			})
		}

		newDeployment.Environment = envVars

		err := i.kubeClient.CreateDeployment(env.Slug, newDeployment)
		spew.Dump(err)
	}

	env, err := i.serviceRepository.GetServiceByName(project.ID, r.Name)

	if err == nil && env != nil {
		return i.serviceRepository.UpdateService(&service)
	}
	return i.serviceRepository.CreateService(&service)
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
