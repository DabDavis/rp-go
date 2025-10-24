package world

import (
	"rp-go/engine/ecs"
)

/*───────────────────────────────────────────────*
 | ACTOR REGISTRY ADAPTER                        |
 *───────────────────────────────────────────────*/

// ActorRegistryAdapter is a helper for systems or tools that need
// to query ECS actors without depending on systems/actor directly.
type ActorRegistryAdapter struct {
	World *ecs.World
}

// FindByID looks up an entity with a matching Actor.ID.
func (r *ActorRegistryAdapter) FindByID(id string) (*ecs.Entity, bool) {
	if r == nil || r.World == nil {
		return nil, false
	}
	manager := r.World.EntitiesManager()
	found := (*ecs.Entity)(nil)
	manager.ForEach(func(e *ecs.Entity) {
		if found != nil {
			return
		}
		act, _ := e.Get("Actor").(*ecs.Actor)
		if act != nil && act.ID == id {
			found = e
		}
	})
	return found, found != nil
}

// All returns all entities with an Actor component.
func (r *ActorRegistryAdapter) All() []*ecs.Entity {
	if r == nil || r.World == nil {
		return nil
	}
	out := []*ecs.Entity{}
	manager := r.World.EntitiesManager()
	manager.ForEach(func(e *ecs.Entity) {
		if e.Has("Actor") {
			out = append(out, e)
		}
	})
	return out
}

