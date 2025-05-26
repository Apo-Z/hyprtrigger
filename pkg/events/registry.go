package events

type Registry struct {
	events               map[string][]*Event
	builtinEvents        map[string][]*Event
	skipBuiltinEvents    bool
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
	r.rebuildEventsList()
}

func (r *Registry) Clear() {
	r.events = make(map[string][]*Event)
	r.builtinEvents = make(map[string][]*Event)
}

func (r *Registry) rebuildEventsList() {
	explicitEvents := make(map[string][]*Event)
	for eventName, eventList := range r.events {
		for _, event := range eventList {
			isBuiltin := false
			for _, builtinEvent := range r.builtinEvents[eventName] {
				if event == builtinEvent {
					isBuiltin = true
					break
				}
			}

			if !isBuiltin {
				if explicitEvents[eventName] == nil {
					explicitEvents[eventName] = make([]*Event, 0)
				}
				explicitEvents[eventName] = append(explicitEvents[eventName], event)
			}
		}
	}

	r.events = make(map[string][]*Event)

	if !r.skipBuiltinEvents {
		for eventName, eventList := range r.builtinEvents {
			r.events[eventName] = append(r.events[eventName], eventList...)
		}
	}

	for eventName, eventList := range explicitEvents {
		r.events[eventName] = append(r.events[eventName], eventList...)
	}
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

// Global functions for backward compatibility
func RegisterEvent(event *Event) {
	DefaultRegistry.RegisterBuiltin(event)
}

func RegisterEventExplicit(event *Event) {
	DefaultRegistry.RegisterExplicit(event)
}

func GetEventsByName(name string) []*Event {
	return DefaultRegistry.GetEventsByName(name)
}

func GetAllEvents() map[string][]*Event {
	return DefaultRegistry.GetAllEvents()
}

func SetSkipBuiltinEvents(skip bool) {
	DefaultRegistry.SetSkipBuiltinEvents(skip)
}

// ClearRegistry clears all events for reload
func ClearRegistry() {
	DefaultRegistry.Clear()
}
