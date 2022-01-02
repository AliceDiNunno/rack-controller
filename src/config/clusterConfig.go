package config

type ClusterConfig struct {
	KubeConfig string
}

func LoadClusterConfig() ClusterConfig {
	kubeConfig, err := GetEnvString("CLUSTER_KUBECONFIG_PATH")

	if err != nil {
		kubeConfig = ""
	}

	return ClusterConfig{
		KubeConfig: kubeConfig,
	}
}