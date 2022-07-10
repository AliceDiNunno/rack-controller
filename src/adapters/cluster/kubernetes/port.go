package kubernetes

import (
	"context"
	"fmt"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/cluster/kubernetes/utils"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	"github.com/davecgh/go-spew/spew"
	v12 "k8s.io/api/core/v1"
	v13 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	v16 "k8s.io/client-go/applyconfigurations/core/v1"
	v15 "k8s.io/client-go/applyconfigurations/meta/v1"
)

func ExecPortTemplate(namespace string, name string, portsToExpose []clusterDomain.Port) v16.ServiceApplyConfiguration {
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

func (k8s kubernetesInstance) ExposePorts(namespace string, name string, ports []clusterDomain.Port) *e.Error {
	data := ExecPortTemplate(namespace, name, ports)

	_, err := k8s.Client.CoreV1().Services(namespace).Apply(context.Background(), &data, v1.ApplyOptions{FieldManager: "rack-controller"})

	if err != nil {
		return e.Wrap(err)
	}
	return nil
}

func (k8s kubernetesInstance) getExposedPorts(namespace string, name string) (*v12.Service, *e.Error) {
	data, err := k8s.Client.CoreV1().Services(namespace).Get(context.Background(), name, v1.GetOptions{})

	if err != nil {
		return nil, e.Wrap(err)
	}

	return data, nil
}
