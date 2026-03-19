package events

type Registry struct {
	events            map[string][]*Event
	builtinEvents     map[string][]*Event
	skipBuiltinEvents bool
}

var DefaultRegistry = NewRegistry()

func NewRegistry() *Registry {
	return &Registry{
		events:            make(map[string][]*Event),
		builtinEvents:     make(map[string][]*Event),
		skipBuiltinEvents: false,
	}
}

func (r *Registry) RegisterBuiltin(event *Event) {
	if r.builtinEvents[event.Name] == nil {
		r.builtinEvents[event.Name] = make([]*Event, 0)
	}
	r.builtinEvents[event.Name] = append(r.builtinEvents[event.Name], event)

	if !r.skipBuiltinEvents {
		if r.events[event.Name] == nil {
			r.events[event.Name] = make([]*Event, 0)
		}
		r.events[event.Name] = append(r.events[event.Name], event)
	}
}

func (r *Registry) RegisterExplicit(event *Event) {
	if r.events[event.Name] == nil {
		r.events[event.Name] = make([]*Event, 0)
	}
	r.events[event.Name] = append(r.events[event.Name], event)
}

func (r *Registry) SetSkipBuiltinEvents(skip bool) {
	r.skipBuiltinEvents = skip
}

func (r *Registry) Clear() {
	r.events = make(map[string][]*Event)
	r.builtinEvents = make(map[string][]*Event)
}

func (r *Registry) GetEventsByName(name string) []*Event {
	return r.events[name]
}

func (r *Registry) GetAllEvents() map[string][]*Event {
	return r.events
}

func (r *Registry) GetBuiltinEvents() map[string][]*Event {
	return r.builtinEvents
}

func GetAllEvents() map[string][]*Event {
	return DefaultRegistry.GetAllEvents()
}
