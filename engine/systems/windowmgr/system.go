package windowmgr

import (
	"rp-go/engine/ecs"
	"rp-go/engine/ui/window"
)

// System scans the ECS world for window components and maintains the shared
// registry consumed by rendering systems. It performs no rendering directly.
type System struct {
	registry *Registry
}

var globalRegistry = NewRegistry()

// NewSystem constructs a window manager bound to the global registry.
func NewSystem() *System {
	return &System{registry: globalRegistry}
}

// SharedRegistry exposes the global window registry for read-only usage.
func SharedRegistry() *Registry {
	return globalRegistry
}

// Update rebuilds the registry from the current set of entities each frame.
func (s *System) Update(world *ecs.World) {
	if world == nil {
		return
	}
	if s.registry == nil {
		s.registry = globalRegistry
	}

	manager := world.EntitiesManager()
	if manager == nil {
		return
	}

	s.registry.Reset()
	manager.ForEach(func(e *ecs.Entity) {
		comp, _ := e.Get("Window").(*window.Component)
		if comp == nil || !comp.Visible {
			return
		}
		s.registry.Add(comp.Layer, comp)
	})
	s.registry.Finalize()
}
