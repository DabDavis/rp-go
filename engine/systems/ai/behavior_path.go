package ai

import (
    "math"
    "strings"

    "rp-go/engine/ecs"
)

func (s *System) applyPathBehavior(state *ecs.AIPathState, cfg *ecs.AIPathBehavior, pos *ecs.Position, vel *ecs.Velocity, speed float64) bool {
    if cfg == nil || len(cfg.Waypoints) == 0 || state == nil {
        return false
    }

    total := len(cfg.Waypoints)
    if total == 0 {
        return false
    }

    if state.Index < 0 || state.Index >= total {
        state.Index = (state.Index%total + total) % total
    }

    variant := strings.ToLower(cfg.Variant)
    if variant == "" {
        variant = "loop"
    }

    if variant == "once" && state.Completed {
        return true
    }

    target := cfg.Waypoints[state.Index]
    dx, dy := target.X-pos.X, target.Y-pos.Y
    dist := math.Hypot(dx, dy)

    if dist <= waypointTolerance {
        s.advancePathState(state, variant, total)
        if variant == "once" && state.Completed {
            return true
        }
        target = cfg.Waypoints[state.Index]
        dx, dy = target.X-pos.X, target.Y-pos.Y
        dist = math.Hypot(dx, dy)
    }

    applyDirectionalVelocity(vel, dx, dy, dist, math.Max(speed, ecs.DefaultAISpeed))
    return true
}

func (s *System) advancePathState(state *ecs.AIPathState, variant string, total int) {
    if total == 0 {
        return
    }

    switch variant {
    case "pingpong":
        if total == 1 {
            state.Completed = true
            return
        }
        if !state.Forward {
            if state.Index <= 0 {
                state.Forward = true
                state.Index = 1
            } else {
                state.Index--
            }
        } else if state.Index >= total-1 {
            state.Forward = false
            state.Index = total - 2
            if state.Index < 0 {
                state.Index = 0
            }
        } else {
            state.Index++
        }

    case "once":
        if state.Index >= total-1 {
            state.Completed = true
        } else {
            state.Index++
        }

    case "random":
        if total > 1 {
            next := state.Index
            for next == state.Index {
                next = s.rng.Intn(total)
            }
            state.Index = next
        } else {
            state.Completed = true
        }

    default: // loop
        state.Index = (state.Index + 1) % total
    }
}

