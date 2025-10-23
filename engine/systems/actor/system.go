package actor

import "rp-go/engine/ecs"

// System keeps an index of all Actor components for other systems to query.
type System struct {
	registry *Registry
}

// NewSystem constructs an Actor system with an empty registry.
func NewSystem() *System {
	return &System{registry: NewRegistry()}
}

// Registry exposes the actor registry so other systems can reuse it without
// iterating the entire entity list each frame.
func (s *System) Registry() *Registry {
	if s.registry == nil {
		s.registry = NewRegistry()
	}
	return s.registry
}

func (s *System) Update(w *ecs.World) {
	if s.registry == nil {
		s.registry = NewRegistry()
	}
	s.registry.Reset()

	primaryAssigned := false

	for _, e := range w.Entities {
		actor, ok := e.Get("Actor").(*ecs.Actor)
		if !ok || actor == nil {
			continue
		}

		s.registry.Add(actor, e)

		if controller, ok := e.Get("PlayerInput").(*ecs.PlayerInput); ok && controller != nil {
			if _, hasAI := e.Get("AIController").(*ecs.AIController); hasAI {
				controller.Enabled = false
				continue
			}
			if !primaryAssigned {
				controller.Enabled = true
				primaryAssigned = true
			} else {
				controller.Enabled = false
			}
		}
	}
}
