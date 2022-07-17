package clusterDomain

type Service struct {
	Name string

	DeploymentName string

	PortName   string
	Protocol   string
	Port       int32
	TargetPort int
	NodePort   int32
}
