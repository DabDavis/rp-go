package movement

import (
	"math"

	"rp-go/engine/ecs"
	"rp-go/engine/events"
	"rp-go/engine/platform"
)

// System updates entity positions based on velocity and
// rotates sprites to face their direction of travel.
type System struct{}

func (s *System) Update(w *ecs.World) {
	for _, e := range w.Entities {
		pos, hasPos := e.Get("Position").(*ecs.Position)
		vel, hasVel := e.Get("Velocity").(*ecs.Velocity)
		if !hasPos || !hasVel {
			continue
		}

		// Skip stationary entities.
		if vel.VX == 0 && vel.VY == 0 {
			continue
		}

		// Move entity.
		pos.X += vel.VX
		pos.Y += vel.VY

		// Rotate sprite toward movement direction.
		if spr, ok := e.Get("Sprite").(*ecs.Sprite); ok {
			// Offset by +90° (π/2 radians) because sprite art faces upward by default.
			spr.Rotation = math.Atan2(vel.VY, vel.VX) + math.Pi/2
		}

		// Publish movement event safely.
		if bus, ok := w.EventBus.(*events.TypedBus); ok && bus != nil {
			events.Queue(bus, events.EntityMovedEvent{
				EntityID: int(e.ID),
				X:        pos.X,
				Y:        pos.Y,
			})
		}
	}
}

func (s *System) Draw(*ecs.World, *platform.Image) {}

