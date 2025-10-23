package windowmgr

import (
	"sort"

	"rp-go/engine/ecs"
	"rp-go/engine/events"
	"rp-go/engine/platform"
	"rp-go/engine/ui/window"
)

// System tracks and updates all window components in the ECS world.
// It maintains a shared global registry for the renderer to read.
type System struct {
	registry *Registry
}

var globalRegistry = NewRegistry()

// NewSystem constructs the manager bound to the global registry.
func NewSystem() *System {
	return &System{registry: globalRegistry}
}

// SharedRegistry exposes the global registry for rendering systems.
func SharedRegistry() *Registry {
	return globalRegistry
}

// Layer reports which ECS layer this system belongs to (HUD layer).
func (s *System) Layer() ecs.DrawLayer { return ecs.LayerHUD }

// Update rebuilds the registry and handles user interactions.
func (s *System) Update(world *ecs.World) {
	if world == nil {
		return
	}

	if s.registry == nil {
		s.registry = globalRegistry
	}

	s.registry.Reset()

	for _, e := range world.Entities {
		comp, _ := e.Get("Window").(*window.Component)
		if comp == nil || !comp.Visible {
			continue
		}
		s.registry.Add(comp)
	}

	s.registry.Finalize()

	// Handle user interaction (mouse clicks, drags)
	bus, _ := world.EventBus.(*events.TypedBus)
	UpdateWindowInteractions(s.registry.All(), bus)
}

// Draw is intentionally empty; rendering is handled in render.WindowRenderer.
func (s *System) Draw(*ecs.World, *platform.Image) {}

// Registry maintains ordered windows grouped by draw layer.
type Registry struct {
	layers     map[ecs.DrawLayer][]*window.Component
	layerOrder []ecs.DrawLayer
}

// NewRegistry creates an empty registry.
func NewRegistry() *Registry {
	return &Registry{
		layers:     make(map[ecs.DrawLayer][]*window.Component),
		layerOrder: make([]ecs.DrawLayer, 0, 8),
	}
}

func (r *Registry) Reset() {
	for k := range r.layers {
		r.layers[k] = r.layers[k][:0]
	}
	r.layerOrder = r.layerOrder[:0]
}

func (r *Registry) Add(comp *window.Component) {
	if comp == nil {
		return
	}
	layer := comp.Layer
	if _, ok := r.layers[layer]; !ok {
		r.layers[layer] = make([]*window.Component, 0, 4)
		r.layerOrder = append(r.layerOrder, layer)
	}
	r.layers[layer] = append(r.layers[layer], comp)
}

func (r *Registry) Finalize() {
	for layer, list := range r.layers {
		sort.SliceStable(list, func(i, j int) bool {
			if list[i].Order == list[j].Order {
				return list[i].ID < list[j].ID
			}
			return list[i].Order < list[j].Order
		})
		r.layers[layer] = list
	}
}

func (r *Registry) Windows(layer ecs.DrawLayer) []*window.Component {
	if r == nil {
		return nil
	}
	return append([]*window.Component(nil), r.layers[layer]...)
}

// All returns all windows across all layers (for input handling).
func (r *Registry) All() []*window.Component {
	if r == nil {
		return nil
	}
	var all []*window.Component
	for _, layer := range r.layerOrder {
		all = append(all, r.layers[layer]...)
	}
	return all
}

