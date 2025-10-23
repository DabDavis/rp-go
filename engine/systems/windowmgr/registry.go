package windowmgr

import (
	"sort"
	"sync"

	"rp-go/engine/ecs"
	"rp-go/engine/ui/window"
)

// Registry stores window components grouped by draw layer for consumption by
// renderers. It does not perform any rendering itself.
type Registry struct {
	mu         sync.RWMutex
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

// Reset clears all registered windows while preserving the allocated slices so
// they can be reused on subsequent frames.
func (r *Registry) Reset() {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	for layer, slice := range r.layers {
		for i := range slice {
			slice[i] = nil
		}
		r.layers[layer] = r.layers[layer][:0]
	}
	r.layerOrder = r.layerOrder[:0]
}

// Add registers a window component for a given draw layer.
func (r *Registry) Add(layer ecs.DrawLayer, comp *window.Component) {
	if r == nil || comp == nil {
		return
	}
	r.add(layer, comp)
}

// add inserts the component for the provided layer. Callers must ensure comp is
// non-nil; the public Add helper enforces this.
func (r *Registry) add(layer ecs.DrawLayer, comp *window.Component) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.layers[layer]; !ok {
		r.layers[layer] = make([]*window.Component, 0, 4)
		r.layerOrder = append(r.layerOrder, layer)
	}
	r.layers[layer] = append(r.layers[layer], comp)
}

// Finalize sorts window entries inside each layer to guarantee stable output
// for renderers consuming the registry.
func (r *Registry) Finalize() {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()

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

// Layers returns the list of draw layers currently holding visible windows.
func (r *Registry) Layers() []ecs.DrawLayer {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]ecs.DrawLayer, len(r.layerOrder))
	copy(result, r.layerOrder)
	return result
}

// Windows returns a copy of the ordered windows for the provided layer.
func (r *Registry) Windows(layer ecs.DrawLayer) []*window.Component {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	src := r.layers[layer]
	if len(src) == 0 {
		return nil
	}
	out := make([]*window.Component, len(src))
	copy(out, src)
	return out
}
