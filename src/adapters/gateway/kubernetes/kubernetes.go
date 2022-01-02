package kubernetes

import (
	"github.com/AliceDiNunno/rack-controller/src/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type kubernetesInstance struct {
	Client *kubernetes.Clientset
}

func LoadInstance(config config.ClusterConfig) (*kubernetesInstance, error) {
	var configLocation *rest.Config

	if config.KubeConfig != "IN_CLUSTER" {
		fileconfig, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
		if err != nil {
			return nil, err
		}
		configLocation = fileconfig
	} else {
		clusterConfig, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		configLocation = clusterConfig
	}

	clientset, err := kubernetes.NewForConfig(configLocation)
	if err != nil {
		return nil, err
	}

	return &kubernetesInstance{
		Client: clientset,
	}, nil
}
