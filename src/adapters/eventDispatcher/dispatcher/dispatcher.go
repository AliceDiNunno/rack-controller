package dispatcher

type Dispatcher struct {
	listeners map[string][]func(interface{})
}

func (d *Dispatcher) Dispatch(event string, payload interface{}) {
	for _, listener := range d.listeners[event] {
		//TODO: track events
		go listener(payload)
	}
}

func (d *Dispatcher) RegisterForEvent(event string, callback func(interface{})) {
	if d.listeners[event] == nil {
		d.listeners[event] = make([]func(interface{}), 0)
	}

	d.listeners[event] = append(d.listeners[event], callback)
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		listeners: make(map[string][]func(interface{})),
	}
}
