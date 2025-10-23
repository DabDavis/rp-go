package windowmgr

import (
	"sync"

	"rp-go/engine/ecs"
	"rp-go/engine/ui/window"
)

// Registry stores the latest window components by draw layer.
type Registry struct {
	mu      sync.RWMutex
	windows map[ecs.DrawLayer][]*window.Component
}

// Add registers a window for the specified layer.
func (r *Registry) Add(layer ecs.DrawLayer, comp *window.Component) {
	if comp == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.windows == nil {
		r.windows = make(map[ecs.DrawLayer][]*window.Component)
	}
	r.windows[layer] = append(r.windows[layer], comp)
}

// Windows returns all windows for the given layer.
func (r *Registry) Windows(layer ecs.DrawLayer) []*window.Component {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]*window.Component(nil), r.windows[layer]...)
}

// Clear removes all windows.
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.windows = make(map[ecs.DrawLayer][]*window.Component)
}

/* -------------------------------------------------------------------------- */
/*                             Shared Registry API                            */
/* -------------------------------------------------------------------------- */

var (
	sharedRegistry *Registry
	once           sync.Once
)

// SharedRegistry returns the global singleton registry instance.
func SharedRegistry() *Registry {
	once.Do(func() {
		sharedRegistry = &Registry{
			windows: make(map[ecs.DrawLayer][]*window.Component),
		}
	})
	return sharedRegistry
}

