package kubernetes

import (
	"context"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
