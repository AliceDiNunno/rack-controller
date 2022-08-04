package kubernetes

import (
	"context"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (k8s kubernetesInstance) GetServicesList(namespace string) ([]string, *e.Error) {
	services, err := k8s.Client.CoreV1().Services(namespace).List(context.Background(), v1.ListOptions{})

	if err != nil {
		return nil, e.Wrap(err).Append(ErrUnableToGetRessource)
	}

	var serviceNames []string
	for _, service := range services.Items {
		serviceNames = append(serviceNames, service.Name)
	}

	return serviceNames, nil
}

func (k8s kubernetesInstance) GetService(namespace string, name string) (clusterDomain.Service, *e.Error) {
	service, err := k8s.Client.CoreV1().Services(namespace).Get(context.Background(), name, v1.GetOptions{})

	if err != nil {
		return clusterDomain.Service{}, e.Wrap(err).Append(clusterDomain.ErrServiceNotFound)
	}

	return clusterDomain.Service{
		Name: service.Name,

		DeploymentName: service.Spec.Selector["app"],

		Port:       service.Spec.Ports[0].Port,
		Protocol:   string(service.Spec.Ports[0].Protocol),
		TargetPort: service.Spec.Ports[0].TargetPort.IntValue(),
		NodePort:   service.Spec.Ports[0].NodePort,
	}, nil
}

func (k8s kubernetesInstance) CreateService(namespace string, request clusterDomain.Service) *e.Error {
	service := v12.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      request.Name,
			Namespace: namespace,
		},
		Spec: v12.ServiceSpec{
			Type: v12.ServiceTypeNodePort,
			Selector: map[string]string{
				"app": request.DeploymentName,
			},
			Ports: []v12.ServicePort{
				v12.ServicePort{
					Name:     request.PortName,
					Port:     request.Port,
					Protocol: v12.Protocol(request.Protocol),
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(request.TargetPort),
					},
					NodePort: request.NodePort,
				},
			},
		},
	}

	_, err := k8s.Client.CoreV1().Services(namespace).Create(context.Background(), &service, v1.CreateOptions{})

	if err != nil {
		return e.Wrap(err).Append(ErrUnableToCreateService)
	}

	return nil
}
