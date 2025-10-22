package input

import (
	"rp-go/engine/ecs"
	"rp-go/engine/platform"
)

type System struct{}

func (s *System) Update(w *ecs.World) {
	for _, e := range w.Entities {
		v, ok := e.Get("Velocity").(*ecs.Velocity)
		if !ok {
			continue
		}
		v.VX, v.VY = 0, 0
		if platform.IsKeyPressed(platform.KeyArrowLeft) || platform.IsKeyPressed(platform.KeyA) {
			v.VX = -2
		}
		if platform.IsKeyPressed(platform.KeyArrowRight) || platform.IsKeyPressed(platform.KeyD) {
			v.VX = 2
		}
		if platform.IsKeyPressed(platform.KeyArrowUp) || platform.IsKeyPressed(platform.KeyW) {
			v.VY = -2
		}
		if platform.IsKeyPressed(platform.KeyArrowDown) || platform.IsKeyPressed(platform.KeyS) {
			v.VY = 2
		}
	}
}

func (s *System) Draw(*ecs.World, *platform.Image) {}
