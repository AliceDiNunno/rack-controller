package event

type dispatcher struct {
	listeners map[string][]func(interface{})
}

func (d *dispatcher) Dispatch(event string, payload interface{}) {
	for _, listener := range d.listeners[event] {
		listener(payload)
	}
}

func (d *dispatcher) RegisterForEvent(event string, callback func(interface{})) {
	if d.listeners[event] == nil {
		d.listeners[event] = make([]func(interface{}), 0)
	}

	d.listeners[event] = append(d.listeners[event], callback)
}

func NewDispatcher() *dispatcher {
	return &dispatcher{
		listeners: make(map[string][]func(interface{})),
	}
}
