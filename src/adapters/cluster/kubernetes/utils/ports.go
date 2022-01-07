package utils

import (
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	v13 "k8s.io/api/core/v1"
	v16 "k8s.io/client-go/applyconfigurations/core/v1"
)

func PortFromDomain(port clusterDomain.Port) v16.ContainerPortApplyConfiguration {
	return v16.ContainerPortApplyConfiguration{
		Name:          &port.Name,
		HostPort:      &port.ExposedPort,
		ContainerPort: &port.ServicePort,
		Protocol:      (*v13.Protocol)(&port.NetworkProtocol),
		HostIP:        nil,
	}
}

func PortsFromDomain(domainPorts []clusterDomain.Port) []v16.ContainerPortApplyConfiguration {
	var ports []v16.ContainerPortApplyConfiguration

	for _, port := range domainPorts {
		ports = append(ports, PortFromDomain(port))
	}

	return ports
}
