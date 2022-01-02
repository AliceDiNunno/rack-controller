package templates

import (
	"fmt"
	"github.com/AliceDiNunno/rack-controller/src/adapters/gateway/kubernetes/utils"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/davecgh/go-spew/spew"
	v13 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	v16 "k8s.io/client-go/applyconfigurations/core/v1"
	v15 "k8s.io/client-go/applyconfigurations/meta/v1"
)

func ExecPortTemplate(namespace string, name string, portsToExpose []domain.Port) v16.ServiceApplyConfiguration {
	apiVersion := "v1"
	deploymentKind := "Service"
	appNamespace := namespace
	ports := utils.PortsFromDomain(portsToExpose)
	serviceName := fmt.Sprintf("%s", name)
	specType := v13.ServiceTypeNodePort

	var portList []v16.ServicePortApplyConfiguration

	for _, port := range ports {
		portList = append(portList, v16.ServicePortApplyConfiguration{
			Name: port.Name,
			Port: port.ContainerPort,
			TargetPort: &intstr.IntOrString{
				Type:   1,
				StrVal: *port.Name,
			},
			NodePort: port.HostPort,
		})
		spew.Dump(portList)
	}

	exposePortConfiguration := v16.ServiceApplyConfiguration{
		TypeMetaApplyConfiguration: v15.TypeMetaApplyConfiguration{
			Kind:       &deploymentKind,
			APIVersion: &apiVersion,
		},
		ObjectMetaApplyConfiguration: &v15.ObjectMetaApplyConfiguration{
			Name:      &serviceName,
			Namespace: &appNamespace,
		},
		Spec: &v16.ServiceSpecApplyConfiguration{
			Ports: portList,
			Selector: map[string]string{
				"app": name,
			},
			Type: &specType,
		},
	}

	return exposePortConfiguration
}
