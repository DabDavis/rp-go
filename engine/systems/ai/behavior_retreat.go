package ai

import (
    "math"

    "rp-go/engine/ecs"
)

func (s *System) applyRetreatBehavior(w *ecs.World, ai *ecs.AIController, pos *ecs.Position, vel *ecs.Velocity) bool {
    cfg := ai.Retreat
    if cfg == nil || cfg.Target == "" {
        return false
    }

    targetPos, ok := s.findTargetPosition(w, cfg.Target)
    if !ok {
        return false
    }

    dx, dy := pos.X-targetPos.X, pos.Y-targetPos.Y
    dist := math.Hypot(dx, dy)

    trigger := math.Max(cfg.TriggerDistance, 150)
    safe := cfg.SafeDistance
    if safe <= trigger {
        safe = trigger * 1.25
    }

    if dist > trigger && (cfg.SafeDistance <= 0 || dist >= safe) {
        return false
    }
    if dist >= safe {
        vel.VX, vel.VY = 0, 0
        return true
    }

    applyDirectionalVelocity(vel, dx, dy, dist, ai.SpeedFor(cfg.Speed))
    return true
}

