package actor

import "rp-go/engine/ecs"

// System maintains and updates the global Actor registry each frame.
// It ensures that only one PlayerInput is active at a time and that AI
// and PlayerInput components do not conflict on the same entity.
type System struct {
	registry *Registry
}

// NewSystem constructs a fresh Actor system with an empty registry.
func NewSystem() *System {
	return &System{registry: NewRegistry()}
}

// Registry exposes the actor registry for other systems to query.
func (s *System) Registry() *Registry {
	if s.registry == nil {
		s.registry = NewRegistry()
	}
	return s.registry
}

// Update rebuilds the actor registry from all entities currently in the world.
// It ensures that only one PlayerInput is enabled, disabling others and
// automatically disabling input for entities that also have AI controllers.
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

		// Index this actor for global lookups
		s.registry.Add(actor, e)

		// Manage PlayerInput exclusivity
		if controller, ok := e.Get("PlayerInput").(*ecs.PlayerInput); ok && controller != nil {
			// Disable input if AI is present
			if _, hasAI := e.Get("AIController").(*ecs.AIController); hasAI {
				controller.Enabled = false
				continue
			}

			// Allow only the first PlayerInput to be active
			if !primaryAssigned {
				controller.Enabled = true
				primaryAssigned = true
			} else {
				controller.Enabled = false
			}
		}
	}
}

