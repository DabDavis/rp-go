package ai

import (
    "math/rand"
    "time"

    "rp-go/engine/ecs"
)

func (s *System) ensureRNG() {
    if s.rng == nil {
        s.rng = newRNG()
    }
}

func newRNG() *rand.Rand {
    return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func applyDirectionalVelocity(v *ecs.Velocity, dx, dy, dist, speed float64) {
    if dist == 0 || speed <= 0 {
        v.VX, v.VY = 0, 0
        return
    }
    if dist < speed {
        speed = dist
    }
    v.VX = (dx / dist) * speed
    v.VY = (dy / dist) * speed
}

