package windowmgr

import (
	"sort"

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

// Registry exposes the global window registry for read-only usage.
func Registry() *Registry {
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
		s.registry.Add(comp)
	})
	s.registry.Finalize()
}

// Registry maintains ordered windows grouped by draw layer.
type Registry struct {
	layers     map[ecs.DrawLayer][]*window.Component
	layerOrder []ecs.DrawLayer
}

// NewRegistry creates an empty registry instance.
func NewRegistry() *Registry {
	return &Registry{
		layers:     make(map[ecs.DrawLayer][]*window.Component),
		layerOrder: make([]ecs.DrawLayer, 0, 8),
	}
}

// Reset clears all tracked windows while preserving allocated slices.
func (r *Registry) Reset() {
	if r == nil {
		return
	}
	for key := range r.layers {
		slice := r.layers[key]
		for i := range slice {
			slice[i] = nil
		}
		r.layers[key] = r.layers[key][:0]
	}
	r.layerOrder = r.layerOrder[:0]
}

// Add registers a window component for later rendering.
func (r *Registry) Add(comp *window.Component) {
	if r == nil || comp == nil {
		return
	}
	layer := comp.Layer
	if _, ok := r.layers[layer]; !ok {
		r.layers[layer] = make([]*window.Component, 0, 4)
		r.layerOrder = append(r.layerOrder, layer)
	}
	r.layers[layer] = append(r.layers[layer], comp)
}

// Finalize sorts windows within each layer by order/id for stable rendering.
func (r *Registry) Finalize() {
	if r == nil {
		return
	}
	for layer, windows := range r.layers {
		sort.SliceStable(windows, func(i, j int) bool {
			if windows[i].Order == windows[j].Order {
				return windows[i].ID < windows[j].ID
			}
			return windows[i].Order < windows[j].Order
		})
		r.layers[layer] = windows
	}
}

// Layers returns the list of layers that currently contain windows.
func (r *Registry) Layers() []ecs.DrawLayer {
	if r == nil {
		return nil
	}
	result := make([]ecs.DrawLayer, len(r.layerOrder))
	copy(result, r.layerOrder)
	return result
}

// Windows returns an ordered copy of the windows for the requested layer.
func (r *Registry) Windows(layer ecs.DrawLayer) []*window.Component {
	if r == nil {
		return nil
	}
	src := r.layers[layer]
	if len(src) == 0 {
		return nil
	}
	out := make([]*window.Component, len(src))
	copy(out, src)
	return out
}
