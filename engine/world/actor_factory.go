package world

import (
	"rp-go/engine/ecs"
)

/*───────────────────────────────────────────────*
 | FACTORY INTERFACE                             |
 *───────────────────────────────────────────────*/

// ActorFactory defines the minimal interface required for any factory
// that can spawn actors or entities from templates.
type ActorFactory interface {
	Spawn(w *ecs.World, template string, pos ecs.Position) (*ecs.Entity, error)
	Templates() []string
}

/*───────────────────────────────────────────────*
 | FACTORY UTILITIES                             |
 *───────────────────────────────────────────────*/

// TrySpawn attempts to create an entity and logs errors internally.
// Returns nil if the spawn failed.
func TrySpawn(factory ActorFactory, w *ecs.World, template string, pos ecs.Position) *ecs.Entity {
	if factory == nil || w == nil {
		return nil
	}
	entity, err := factory.Spawn(w, template, pos)
	if err != nil {
		println("[WORLD] Spawn failed:", err.Error())
		return nil
	}
	return entity
}

// HasTemplate checks if a factory supports a given template name.
func HasTemplate(factory ActorFactory, name string) bool {
	if factory == nil {
		return false
	}
	for _, n := range factory.Templates() {
		if n == name {
			return true
		}
	}
	return false
}

