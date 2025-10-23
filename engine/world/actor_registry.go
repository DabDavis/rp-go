package world

import "rp-go/engine/systems/actor"

// GlobalActorRegistry provides an optional shared registry instance for world-level queries.
// This is useful for scripting, triggers, or saving/loading state.
var GlobalActorRegistry *actor.Registry

// SetGlobalRegistry assigns a shared actor registry from the active ActorSystem.
func SetGlobalRegistry(r *actor.Registry) {
	GlobalActorRegistry = r
}

// Actors returns a sorted snapshot of all active actor entities.
func Actors() []*actor.Registry {
	if GlobalActorRegistry == nil {
		return nil
	}
	return []*actor.Registry{GlobalActorRegistry}
}

