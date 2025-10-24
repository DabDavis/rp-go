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

// World owns all entities, systems, and draw layers.
type World struct {
	nextID        EntityID
	Entities      []*Entity
	Systems       []System
	EventBus      any
	entityManager *EntityManager

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
 | CONSTRUCTOR                                  |
 *───────────────────────────────────────────────*/

func NewWorld() *World {
	w := &World{
		drawBuckets:   make(map[DrawLayer][]drawEntry, 8),
		worldLayers:   []DrawLayer{LayerBackground, LayerWorld, LayerForeground},
		overlayLayers: []DrawLayer{LayerHUD, LayerEntityList, LayerDebug, LayerConsole},
		Systems:       make([]System, 0, 16),
		Entities:      make([]*Entity, 0, 256),
	}
	w.entityManager = newEntityManager(w)
	return w
}

/*───────────────────────────────────────────────*
 | ENTITY MANAGEMENT                            |
 *───────────────────────────────────────────────*/

func (w *World) NewEntity() *Entity {
	e := NewEntity(w.nextID)
	w.nextID++
	w.Entities = append(w.Entities, e)
	return e
}

func (w *World) RemoveEntity(target *Entity) {
	if w == nil || target == nil {
		return
	}
	for i, e := range w.Entities {
		if e == target {
			w.Entities = append(w.Entities[:i], w.Entities[i+1:]...)
			break
		}
	}
}

func (w *World) RemoveEntityByID(id EntityID) {
	if w == nil {
		return
	}
	for _, e := range w.Entities {
		if e != nil && e.ID == id {
			w.RemoveEntity(e)
			return
		}
	}
}

// EntitiesManager returns the entity manager for iteration utilities.
func (w *World) EntitiesManager() *EntityManager {
	if w == nil {
		return nil
	}
	if w.entityManager == nil {
		w.entityManager = newEntityManager(w)
	}
	return w.entityManager
}

/*───────────────────────────────────────────────*
 | SYSTEM MANAGEMENT                            |
 *───────────────────────────────────────────────*/

// AddSystem registers a new system into the ECS world.
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

	// Register drawable system if applicable
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

	w.Systems = append(w.Systems, s)
}

/*───────────────────────────────────────────────*
 | SYSTEM LOOKUP                                |
 *───────────────────────────────────────────────*/

// FindSystem returns the first system matching the given type.
// Example usage:
//   ai := w.FindSystem((*ai.System)(nil)).(*ai.System)
func (w *World) FindSystem(ptrType any) System {
	if w == nil || ptrType == nil {
		return nil
	}
	for _, entry := range w.systemEntries {
		if sameType(entry.system, ptrType) {
			return entry.system
		}
	}
	return nil
}

// sameType compares dynamic types via reflection without allocating.
func sameType(a, b any) bool {
	if a == nil || b == nil {
		return false
	}
	return fmt.Sprintf("%T", a) == fmt.Sprintf("%T", b)
}

/*───────────────────────────────────────────────*
 | UPDATE LOOP                                  |
 *───────────────────────────────────────────────*/

var EnableProfiling bool // Toggle profiling per system

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
 | DRAWING PIPELINE                             |
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
 | LAYER MANAGEMENT                             |
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
 | INTERNAL HELPERS                             |
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
	out := make([]DrawLayer, 0, len(layers))
	for _, layer := range layers {
		out = appendLayerIfMissing(out, layer)
	}
	return out
}

func isOverlayLayer(layer DrawLayer) bool {
	return layer >= LayerHUD
}

