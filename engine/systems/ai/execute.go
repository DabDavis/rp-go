package ai

import "rp-go/engine/ecs"

func (s *System) executeAction(w *ecs.World, e *ecs.Entity, pos *ecs.Position, vel *ecs.Velocity, act AIActionInstance) bool {
	fn, ok := GlobalBehaviorCatalog.Get(act.Type)
	if !ok {
		return false
	}
	return fn(w, e, pos, vel, act.Params)
}

