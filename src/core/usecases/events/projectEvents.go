package events

import "github.com/AliceDiNunno/rack-controller/src/core/domain"

func (h EventHandler) ProjectCreatedEvent(data interface{}) {
	project, ok := data.(domain.Project)

	if !ok {
		return
	}

	h

	h.cluster.CreateNamespace()
}
