package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
)

func (i interactor) GetNodes() ([]clusterDomain.Node, *e.Error) {
	return i.kubeClient.GetNodes()
}
