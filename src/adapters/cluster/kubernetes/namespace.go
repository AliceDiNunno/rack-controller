package kubernetes

import (
	"context"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/davecgh/go-spew/spew"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k8s kubernetesInstance) CreateNamespace(namespace string) *e.Error {
	_, err := k8s.Client.CoreV1().Namespaces().Create(context.Background(), &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: namespace,
		},
	}, v1.CreateOptions{})
	if err != nil {
		spew.Dump(err)
		return e.Wrap(err)
	}
	return nil
}

func (k8s kubernetesInstance) GetNamespaces() ([]string, *e.Error) {
	namespaces, err := k8s.Client.CoreV1().Namespaces().List(context.Background(), v1.ListOptions{})

	if err != nil {
		return nil, e.Wrap(err)
	}

	var mapList []string

	for _, namespaceEntry := range namespaces.Items {
		mapList = append(mapList, namespaceEntry.Name)
	}

	return mapList, e.Wrap(err)
}
