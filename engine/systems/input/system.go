package input

import (
	"math"

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
		sprite, hasSprite := e.Get("Sprite").(*ecs.Sprite)
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

		if hasSprite {
			switch {
			case v.VX < 0:
				sprite.Rotation = math.Pi
			case v.VX > 0:
				sprite.Rotation = 0
			case v.VY < 0:
				sprite.Rotation = -math.Pi / 2
			case v.VY > 0:
				sprite.Rotation = math.Pi / 2
			}
		}
	}
}

func (s *System) Draw(*ecs.World, *platform.Image) {}
