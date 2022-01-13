package events

import "github.com/AliceDiNunno/rack-controller/src/core/domain"

func (h EventHandler) SetEvents() {
	h.dispatcher.RegisterForEvent(domain.EventProjectCreated, h.ProjectCreatedEvent)
}
