package ecs

import "rp-go/engine/platform"

type World struct {
	nextID   EntityID
	Entities []*Entity
	Systems  []System
	EventBus any
}

func NewWorld() *World {
	return &World{}
}

func (w *World) NewEntity() *Entity {
	e := NewEntity(w.nextID)
	w.nextID++
	w.Entities = append(w.Entities, e)
	return e
}

func (w *World) AddSystem(s System) {
	w.Systems = append(w.Systems, s)
}

func (w *World) Update() {
	for _, s := range w.Systems {
		s.Update(w)
	}
}

// DrawWorld: draws only LayerWorld systems (camera-affected)
func (w *World) DrawWorld(screen *platform.Image) {
	for _, s := range w.Systems {
		if ls, ok := s.(LayeredSystem); ok {
			if ls.Layer() != LayerWorld {
				continue
			}
		}
		s.Draw(w, screen)
	}
}

// DrawOverlay: draws only LayerOverlay systems (screen-space)
func (w *World) DrawOverlay(screen *platform.Image) {
	for _, s := range w.Systems {
		if ls, ok := s.(LayeredSystem); ok {
			if ls.Layer() != LayerOverlay {
				continue
			}
		}
		s.Draw(w, screen)
	}
}

