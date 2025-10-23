package ai

import (
    "math"

    "rp-go/engine/ecs"
)

func (s *System) applyFollowBehavior(w *ecs.World, ai *ecs.AIController, pos *ecs.Position, vel *ecs.Velocity) bool {
    cfg := ai.Follow
    if cfg == nil || cfg.Target == "" {
        return false
    }

    targetPos, ok := s.findTargetPosition(w, cfg.Target)
    if !ok {
        return false
    }

    dx := targetPos.X + cfg.OffsetX - pos.X
    dy := targetPos.Y + cfg.OffsetY - pos.Y
    dist := math.Hypot(dx, dy)
    if dist == 0 {
        return true
    }

    minDist := math.Max(cfg.MinDistance, 1)
    if dist <= minDist {
        return true
    }

    speed := ai.SpeedFor(cfg.Speed)
    if max := cfg.MaxDistance; max > minDist && dist < max {
        factor := (dist - minDist) / (max - minDist)
        speed *= math.Max(factor, 0)
    }

    applyDirectionalVelocity(vel, dx, dy, dist, speed)
    return true
}

