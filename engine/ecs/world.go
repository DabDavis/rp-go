package ecs

import (
	"sort"

	"rp-go/engine/platform"
)

type World struct {
	nextID   EntityID
	Entities []*Entity
	Systems  []System
	EventBus any

	systemEntries []systemEntry
	drawBuckets   map[DrawLayer][]drawEntry
	worldLayers   []DrawLayer
	overlayLayers []DrawLayer
	nextOrder     int
}

type systemEntry struct {
	system   System
	priority int
	order    int
}

type drawEntry struct {
	system   DrawableSystem
	priority int
	order    int
}

func NewWorld() *World {
	return &World{
		drawBuckets:   make(map[DrawLayer][]drawEntry),
		worldLayers:   []DrawLayer{LayerBackground, LayerWorld, LayerForeground},
		overlayLayers: []DrawLayer{LayerHUD, LayerDebug},
	}
}

func (w *World) NewEntity() *Entity {
	e := NewEntity(w.nextID)
	w.nextID++
	w.Entities = append(w.Entities, e)
	return e
}

func (w *World) AddSystem(s System) {
	w.Systems = append(w.Systems, s)

	entry := systemEntry{
		system:   s,
		priority: systemPriority(s),
		order:    w.nextOrder,
	}
	w.nextOrder++

	w.systemEntries = append(w.systemEntries, entry)
	stableSortSystems(w.systemEntries)

	if drawable, ok := s.(DrawableSystem); ok {
		layer := resolveLayer(drawable)
		w.ensureLayerRegistered(layer)
		draw := drawEntry{
			system:   drawable,
			priority: entry.priority,
			order:    entry.order,
		}
		bucket := append(w.drawBuckets[layer], draw)
		stableSortDraw(bucket)
		w.drawBuckets[layer] = bucket
	}
}

func (w *World) Update() {
	for _, entry := range w.systemEntries {
		entry.system.Update(w)
	}
}

// DrawWorld renders all world-space layers (background → foreground).
func (w *World) DrawWorld(screen *platform.Image) {
	w.drawLayerGroup(screen, w.worldLayers)
}

// DrawOverlay renders all overlay layers (HUD → debug).
func (w *World) DrawOverlay(screen *platform.Image) {
	w.drawLayerGroup(screen, w.overlayLayers)
}

// DrawLayer renders every system registered for the supplied layer.
func (w *World) DrawLayer(screen *platform.Image, layer DrawLayer) {
	w.drawLayerGroup(screen, []DrawLayer{layer})
}

// DrawLayers renders each layer in sequence.
func (w *World) DrawLayers(screen *platform.Image, layers ...DrawLayer) {
	w.drawLayerGroup(screen, layers)
}

// SetWorldLayers overrides the draw order for world-space layers.
func (w *World) SetWorldLayers(layers ...DrawLayer) {
	w.worldLayers = uniqueLayers(layers)
}

// SetOverlayLayers overrides the draw order for overlay layers.
func (w *World) SetOverlayLayers(layers ...DrawLayer) {
	w.overlayLayers = uniqueLayers(layers)
}

func (w *World) drawLayerGroup(screen *platform.Image, layers []DrawLayer) {
	for _, layer := range layers {
		entries := w.drawBuckets[layer]
		for _, entry := range entries {
			entry.system.Draw(w, screen)
		}
	}
}

func (w *World) ensureLayerRegistered(layer DrawLayer) {
	if _, exists := w.drawBuckets[layer]; !exists {
		w.drawBuckets[layer] = nil

		if isOverlayLayer(layer) {
			w.overlayLayers = appendLayerIfMissing(w.overlayLayers, layer)
		} else {
			w.worldLayers = appendLayerIfMissing(w.worldLayers, layer)
		}
	}
}

func systemPriority(s System) int {
	if ps, ok := s.(PrioritizedSystem); ok {
		return ps.Priority()
	}
	return 0
}

func resolveLayer(s System) DrawLayer {
	if ls, ok := s.(LayeredSystem); ok {
		return ls.Layer()
	}
	return LayerWorld
}

func stableSortSystems(entries []systemEntry) {
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].priority == entries[j].priority {
			return entries[i].order < entries[j].order
		}
		return entries[i].priority < entries[j].priority
	})
}

func stableSortDraw(entries []drawEntry) {
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].priority == entries[j].priority {
			return entries[i].order < entries[j].order
		}
		return entries[i].priority < entries[j].priority
	})
}

func appendLayerIfMissing(layers []DrawLayer, layer DrawLayer) []DrawLayer {
	for _, existing := range layers {
		if existing == layer {
			return layers
		}
	}
	return append(layers, layer)
}

func uniqueLayers(layers []DrawLayer) []DrawLayer {
	if len(layers) == 0 {
		return nil
	}
	out := make([]DrawLayer, 0, len(layers))
	for _, layer := range layers {
		out = appendLayerIfMissing(out, layer)
	}
	return out
}

func isOverlayLayer(layer DrawLayer) bool {
	return layer >= LayerHUD
}
