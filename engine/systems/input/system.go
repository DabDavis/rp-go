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
			vx -= moveSpeed
		}
		if platform.IsKeyPressed(platform.KeyArrowRight) || platform.IsKeyPressed(platform.KeyD) {
			vx += moveSpeed
		}
		if platform.IsKeyPressed(platform.KeyArrowUp) || platform.IsKeyPressed(platform.KeyW) {
			vy -= moveSpeed
		}
		if platform.IsKeyPressed(platform.KeyArrowDown) || platform.IsKeyPressed(platform.KeyS) {
			vy += moveSpeed
		}

		padVX, padVY := 0.0, 0.0
		for _, id := range platform.GamepadIDs() {
			if _, ok := platform.StandardGamepadLayoutID(id); !ok {
				continue
			}

			padVX = platform.StandardGamepadAxisValue(id, platform.StandardGamepadAxisLeftStickHorizontal)
			padVY = platform.StandardGamepadAxisValue(id, platform.StandardGamepadAxisLeftStickVertical)

			if math.Abs(padVX) < 0.1 {
				padVX = 0
			}
			if math.Abs(padVY) < 0.1 {
				padVY = 0
			}

			if platform.IsStandardGamepadButtonPressed(id, platform.StandardGamepadButtonLeft) {
				padVX = -1
			} else if platform.IsStandardGamepadButtonPressed(id, platform.StandardGamepadButtonRight) {
				padVX = 1
			}
			if platform.IsStandardGamepadButtonPressed(id, platform.StandardGamepadButtonUp) {
				padVY = -1
			} else if platform.IsStandardGamepadButtonPressed(id, platform.StandardGamepadButtonDown) {
				padVY = 1
			}

			if padVX != 0 || padVY != 0 {
				vx = moveSpeed * padVX
				vy = moveSpeed * padVY
				break
			}
		}

		v.VX, v.VY = vx, vy

		if hasSprite && (v.VX != 0 || v.VY != 0) {
			sprite.Rotation = math.Atan2(v.VY, v.VX)
		}

		if hasSprite {
			if v.VX < 0 {
				sprite.FlipHorizontal = true
			} else if v.VX > 0 {
				sprite.FlipHorizontal = false
			}
		}
	}
}

func (s *System) Draw(*ecs.World, *platform.Image) {}

