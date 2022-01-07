package clusterDomain

import (
	"time"
)

type DeploymentDiscoveredEndpoint struct {
	Ip   string
	Port int32
	//A public endpoint will be fetched from ingress
	//A private endpoint will be fetched from node IPs and the service with the web port.
	IsPublic bool
	IsSecure bool
}

type Deployment struct {
	Id                string
	Name              string
	ImageName         string
	CreationDate      time.Time
	Generation        int64
	Replicas          int64
	UpdatedReplicas   int64
	AvailableReplicas int64
	ReadyReplicas     int64
	Container         Container
	Condition         DeploymentCondition

	Endpoint []DeploymentDiscoveredEndpoint
}

type DeploymentCondition struct {
	Available   bool
	Progressing bool
}

type DeploymentCreationRequest struct {
	TemplateId     int
	DeploymentName string
	ImageName      string
	Ports          []Port
	Environment    []Environment
	ConfigMaps     []string
	Secrets        []string
	Replicas       int
	Memory         int
	CPU            int
}
