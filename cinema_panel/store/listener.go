package store

type ListenerRegistry struct {
	listeners map[interface{}]func()
}

func NewListenerRegistry() *ListenerRegistry {
	return &ListenerRegistry{
		listeners: make(map[interface{}]func()),
	}
}

func (r *ListenerRegistry) Add(key interface{}, listener func()) {
	if key == nil {
		key = new(int)
	}
	if _, ok := r.listeners[key]; ok {
		panic("duplicate listener key")
	}
	r.listeners[key] = listener
}

func (r *ListenerRegistry) Remove(key interface{}) {
	delete(r.listeners, key)
}

func (r *ListenerRegistry) Fire() {
	for _, l := range r.listeners {
		l()
	}
}
