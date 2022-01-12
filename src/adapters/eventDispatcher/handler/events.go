package handler

import (
	"github.com/AliceDiNunno/rack-controller/src/adapters/eventDispatcher/dispatcher"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/AliceDiNunno/rack-controller/src/core/usecases"
)

type EventHandler struct {
	ucHandler  usecases.Usecases
	dispatcher *dispatcher.Dispatcher
}

func (h EventHandler) HandleProjectCreated(data interface{}) {
}

func (h EventHandler) SetEvents() {
	h.dispatcher.RegisterForEvent(domain.EventProjectCreated, h.HandleProjectCreated)
}

func NewEventHandler(ucHandler usecases.Usecases, dp *dispatcher.Dispatcher) EventHandler {
	return EventHandler{ucHandler, dp}
}
