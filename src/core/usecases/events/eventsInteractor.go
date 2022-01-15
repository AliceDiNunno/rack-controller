package events

import (
	"github.com/AliceDiNunno/rack-controller/src/adapters/eventDispatcher/dispatcher"
	"github.com/AliceDiNunno/rack-controller/src/core/usecases"
	"github.com/AliceDiNunno/rack-controller/src/core/usecases/kubernetes"
)

type EventHandler struct {
	ucHandler  usecases.Usecases
	dispatcher *dispatcher.Dispatcher
	cluster    kubernetes.Kubernetes
}

func NewEventHandler(ucHandler usecases.Usecases, dp *dispatcher.Dispatcher, cluster kubernetes.Kubernetes) EventHandler {
	return EventHandler{
		ucHandler:  ucHandler,
		dispatcher: dp,
		cluster:    cluster,
	}
}
