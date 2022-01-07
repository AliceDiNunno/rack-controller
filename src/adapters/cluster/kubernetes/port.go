package kubernetes

import (
	"context"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/adapters/cluster/kubernetes/templates"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k8s kubernetesInstance) ExposePorts(namespace string, name string, ports []clusterDomain.Port) *e.Error {
	data := templates.ExecPortTemplate(namespace, name, ports)

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
