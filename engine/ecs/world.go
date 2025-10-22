package ecs

import "github.com/hajimehoshi/ebiten/v2"

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

func (w *World) Draw(screen *ebiten.Image) {
	for _, s := range w.Systems {
		s.Draw(w, screen)
	}
}

