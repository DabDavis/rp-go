package movement

import (
	"rp-go/engine/ecs"
	"rp-go/engine/events"
	"rp-go/engine/platform"
)

type System struct{}

func (s *System) Update(w *ecs.World) {
	for _, e := range w.Entities {
		pos, ok1 := e.Get("Position").(*ecs.Position)
		vel, ok2 := e.Get("Velocity").(*ecs.Velocity)
		if !ok1 || !ok2 {
			continue
		}

		// 2D movement
		if vel.VX == 0 && vel.VY == 0 {
			continue
		}

		pos.X += vel.VX
		pos.Y += vel.VY

		// Publish movement event safely
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
