package clusterDomain

var (
	QuantityKilobytes = int64(1024)
	QuantityMegabytes = QuantityKilobytes * 1024
	QuantityGigabytes = QuantityMegabytes * 1024
)

type PersistentVolume struct {
	Name        string
	StorageSize int64
	MountPath   string
}

type PersistentVolumeClaim struct {
	Name        string
	StorageSize int64
}

type VolumeDeployment struct {
	Name      string
	ClaimName string

	DeploymentName string

	MountPath   string
	StorageSize int64
	SubPath     string
}

//
