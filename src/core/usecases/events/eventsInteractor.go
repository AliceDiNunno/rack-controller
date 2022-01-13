package events

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/eventDispatcher/dispatcher"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	"github.com/AliceDiNunno/rack-controller/src/core/usecases"
)

type Kubernetes interface {
	GetNodes() ([]clusterDomain.Node, *e.Error)
	GetNode(string) (*clusterDomain.Node, *e.Error)

	ListDeployments(namespace string) ([]clusterDomain.Deployment, *e.Error)
	GetDeployment(namespace string, deploymentName string) (*clusterDomain.Deployment, *e.Error)
	GetPodsOfADeployment(namespace string, deploymentName string) ([]clusterDomain.Pod, *e.Error)
	GetEnvironmentOfADeployment(namespace string, deploymentName string) ([]clusterDomain.Environment, *e.Error)
	GetPortsOfADeployment(namespace string, deploymentName string) ([]clusterDomain.Port, *e.Error)
	GetConfigMapsOfADeployment(namespace string, name string) ([]string, *e.Error)
	GetSecretsOfADeployment(namespace string, name string) ([]string, *e.Error)
	DeleteDeployment(namespace string, name string) *e.Error
	RestartDeployment(namespace string, name string) *e.Error
	CreateDeployment(namespace string, request clusterDomain.DeploymentCreationRequest) *e.Error
	ExposePorts(namespace string, name string, ports []clusterDomain.Port) *e.Error

	GetPods(namespace string) ([]clusterDomain.Pod, *e.Error)
	GetPod(namespace string, podName string) (*clusterDomain.Pod, *e.Error)
	GetPodLogs(namespace string, podName string) (string, *e.Error)
	DeletePod(namespace string, podName string) *e.Error

	CreateNamespace(namespace string) *e.Error
	GetNamespaces() ([]string, *e.Error)

	GetConfigMapList(namespace string) ([]string, *e.Error)
	GetConfigMap(namespace string, name string) (clusterDomain.ConfigMap, *e.Error)
	CreateConfigMap(namespace string, request clusterDomain.ConfigMapCreationRequest) *e.Error
	UpdateConfigMap(namespace string, name string, request clusterDomain.ConfigMapUpdateRequest) *e.Error
	DeleteConfigMap(namespace string, name string) *e.Error

	GetSecretsList(namespace string) ([]string, *e.Error)
	GetSecret(namespace string, name string) (clusterDomain.Secret, *e.Error)
	CreateSecret(namespace string, request clusterDomain.SecretCreationRequest) *e.Error
	UpdateSecret(namespace string, name string, request clusterDomain.SecretUpdateRequest) *e.Error
	DeleteSecret(namespace string, name string) *e.Error
}

type EventHandler struct {
	ucHandler  usecases.Usecases
	dispatcher *dispatcher.Dispatcher
	cluster    Kubernetes
}

func NewEventHandler(ucHandler usecases.Usecases, dp *dispatcher.Dispatcher, cluster Kubernetes) EventHandler {
	return EventHandler{
		ucHandler:  ucHandler,
		dispatcher: dp,
		cluster:    cluster,
	}
}
