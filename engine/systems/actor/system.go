package actor

import (
	"fmt"

	"rp-go/engine/ecs"
)

type System struct {
	actors map[string]*ecs.Entity
}

func (s *System) Update(w *ecs.World) {
	if s.actors == nil {
		s.actors = make(map[string]*ecs.Entity)
	}

	for _, e := range w.Entities {
		if a, ok := e.Get("Actor").(*ecs.Actor); ok {
			if _, exists := s.actors[a.ID]; !exists {
				s.actors[a.ID] = e
				fmt.Printf("[ACTOR] Registered %s (%s)\n", a.ID, a.Archetype)
			}
		}
	}
}
