package ai

import (
    "math"

    "rp-go/engine/ecs"
)

func (s *System) applyPursueBehavior(w *ecs.World, ai *ecs.AIController, pos *ecs.Position, vel *ecs.Velocity) bool {
    cfg := ai.Pursue
    if cfg == nil || cfg.Target == "" {
        return false
    }

    targetPos, ok := s.findTargetPosition(w, cfg.Target)
    if !ok {
        return false
    }

    dx, dy := targetPos.X-pos.X, targetPos.Y-pos.Y
    dist := math.Hypot(dx, dy)
    if dist == 0 {
        return true
    }

    if cfg.EngageDistance > 0 && dist > cfg.EngageDistance {
        return false
    }

    applyDirectionalVelocity(vel, dx, dy, dist, ai.SpeedFor(cfg.Speed))
    return true
}

