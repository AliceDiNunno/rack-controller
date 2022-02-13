package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
)

func (i interactor) GetServiceInstances(service *domain.Service, environments *domain.Environment) ([]clusterDomain.Pod, *e.Error) {
	if service == nil {
		return nil, e.Wrap(domain.ErrServiceNotFound)
	}

	if environments == nil {
		return nil, e.Wrap(domain.ErrEnvironmentNotFound)
	}

	return i.kubeClient.GetPodsOfADeployment(environments.Slug, service.Slug)
}

func (i interactor) GetInstanceByName(service *domain.Service, environment *domain.Environment, name string) (*clusterDomain.Pod, *e.Error) {
	if service == nil {
		return nil, e.Wrap(domain.ErrServiceNotFound)
	}

	if environment == nil {
		return nil, e.Wrap(domain.ErrEnvironmentNotFound)
	}

	return i.kubeClient.GetPod(environment.Slug, name)
}

func (i interactor) GetInstanceLogs(service *domain.Service, environment *domain.Environment, instance *clusterDomain.Pod) (string, *e.Error) {
	if service == nil {
		return "", e.Wrap(domain.ErrServiceNotFound)
	}

	if environment == nil {
		return "", e.Wrap(domain.ErrEnvironmentNotFound)
	}

	if instance == nil {
		return "", e.Wrap(domain.ErrInstanceNotFound)
	}

	return i.kubeClient.GetPodLogs(environment.Slug, instance.Name)
}
