package ecs

import (
	"fmt"
	"sort"
	"time"

	"rp-go/engine/platform"
)

/*───────────────────────────────────────────────*
 | WORLD STRUCTURE                               |
 *───────────────────────────────────────────────*/

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

/*───────────────────────────────────────────────*
 | WORLD LIFECYCLE                               |
 *───────────────────────────────────────────────*/

func NewWorld() *World {
	return &World{
		drawBuckets:   make(map[DrawLayer][]drawEntry, 8),
		worldLayers:   []DrawLayer{LayerBackground, LayerWorld, LayerForeground},
		overlayLayers: []DrawLayer{LayerHUD, LayerEntityList, LayerDebug, LayerConsole},
		Systems:       make([]System, 0, 16),
		Entities:      make([]*Entity, 0, 256),
	}
}

func (w *World) NewEntity() *Entity {
	e := NewEntity(w.nextID)
	w.nextID++
	w.Entities = append(w.Entities, e)
	return e
}

/*───────────────────────────────────────────────*
 | SYSTEM MANAGEMENT                             |
 *───────────────────────────────────────────────*/

func (w *World) AddSystem(s System) {
	if s == nil {
		return
	}

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

/*───────────────────────────────────────────────*
 | UPDATE LOOP                                   |
 *───────────────────────────────────────────────*/

var EnableProfiling bool // Toggle for profiling per-system timings

func (w *World) Update() {
	for _, entry := range w.systemEntries {
		if EnableProfiling {
			start := time.Now()
			entry.system.Update(w)
			elapsed := time.Since(start)
			if elapsed > 2*time.Millisecond {
				fmt.Printf("[ECS] %s took %v\n", SystemName(entry.system), elapsed)
			}
		} else {
			entry.system.Update(w)
		}
	}
}

/*───────────────────────────────────────────────*
 | ENTITY MANAGEMENT                             |
 *───────────────────────────────────────────────*/

func (w *World) RemoveEntity(target *Entity) {
	if w == nil || target == nil {
		return
	}
	for i, entity := range w.Entities {
		if entity == target {
			w.Entities = append(w.Entities[:i], w.Entities[i+1:]...)
			break
		}
	}
}

func (w *World) RemoveEntityByID(id EntityID) {
	if w == nil {
		return
	}
	for _, entity := range w.Entities {
		if entity != nil && entity.ID == id {
			w.RemoveEntity(entity)
			return
		}
	}
}

/*───────────────────────────────────────────────*
 | DRAWING PIPELINE                              |
 *───────────────────────────────────────────────*/

func (w *World) DrawWorld(screen *platform.Image) {
	w.drawLayerGroup(screen, w.worldLayers)
}

func (w *World) DrawOverlay(screen *platform.Image) {
	w.drawLayerGroup(screen, w.overlayLayers)
}

func (w *World) DrawLayer(screen *platform.Image, layer DrawLayer) {
	w.drawLayerGroup(screen, []DrawLayer{layer})
}

func (w *World) DrawLayers(screen *platform.Image, layers ...DrawLayer) {
	w.drawLayerGroup(screen, layers)
}

func (w *World) drawLayerGroup(screen *platform.Image, layers []DrawLayer) {
	if w == nil || screen == nil {
		return
	}
	for _, layer := range layers {
		if entries, ok := w.drawBuckets[layer]; ok {
			for _, entry := range entries {
				if entry.system != nil {
					entry.system.Draw(w, screen)
				}
			}
		}
	}
}

/*───────────────────────────────────────────────*
 | LAYER MANAGEMENT                              |
 *───────────────────────────────────────────────*/

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

func (w *World) SetWorldLayers(layers ...DrawLayer) {
	w.worldLayers = uniqueLayers(layers)
}

func (w *World) SetOverlayLayers(layers ...DrawLayer) {
	w.overlayLayers = uniqueLayers(layers)
}

/*───────────────────────────────────────────────*
 | INTERNAL HELPERS                              |
 *───────────────────────────────────────────────*/

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

