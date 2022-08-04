package kubernetes

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
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
	AddSecretsToDeployment(namespace string, deploymentSlug string, secretSlug string) *e.Error
	DeleteDeployment(namespace string, name string) *e.Error
	RestartDeployment(namespace string, name string) *e.Error
	CreateDeployment(namespace string, request clusterDomain.DeploymentCreationRequest) *e.Error
	ExposePorts(namespace string, name string, ports []clusterDomain.Port) *e.Error

	GetPods(namespace string) ([]clusterDomain.Pod, *e.Error)
	GetPod(namespace string, podName string) (*clusterDomain.Pod, *e.Error)
	GetPodLogs(namespace string, podName string) (string, *e.Error)
	GetPodsOfANode(node string) ([]clusterDomain.Pod, *e.Error)
	DeletePod(namespace string, podName string) *e.Error

	CreateNamespace(namespace string) *e.Error
	GetNamespaces() ([]string, *e.Error)
	NamespaceExists(namespace string) bool

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
	DeleteNamespace(namespace string) *e.Error
	CreateService(namespace string, request clusterDomain.Service) *e.Error

	CreatePersistentVolume(namespace string, request clusterDomain.PersistentVolume) *e.Error
	CreatePersistentVolumeClaim(namespace string, request clusterDomain.PersistentVolumeClaim) *e.Error
	AddVolumeToDeployment(namespace string, deploymentSlug string, volume clusterDomain.VolumeDeployment) *e.Error
}
