package ai

import "rp-go/engine/ecs"

// executeAction dispatches an AI action to its registered behavior handler.
// Returns true if the action executed successfully this frame.
func (s *System) executeAction(
    w *ecs.World,
    e *ecs.Entity,
    pos *ecs.Position,
    vel *ecs.Velocity,
    act ecs.AIActionInstance,
) bool {
    fn, ok := GlobalBehaviorCatalog.Get(act.Type)
    if !ok {
        return false
    }
    return fn(w, e, pos, vel, act.Params)
}

