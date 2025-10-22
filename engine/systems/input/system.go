package input

import (
	"math"

	"rp-go/engine/ecs"
	"rp-go/engine/platform"
)

type System struct{}

// movement speed per frame
const moveSpeed = 3.0

func (s *System) Update(w *ecs.World) {
	for _, e := range w.Entities {
		v, ok := e.Get("Velocity").(*ecs.Velocity)
		if !ok {
			continue
		}

		sprite, hasSprite := e.Get("Sprite").(*ecs.Sprite)

		// Local velocity deltas
		vx, vy := 0.0, 0.0

		/* ---------------------------- Keyboard movement --------------------------- */
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

		/* ----------------------------- Gamepad movement --------------------------- */
		for _, id := range platform.GamepadIDs() {
			if !platform.IsStandardGamepadLayoutAvailable(id) {
				continue
			}

			// Analog axes
			padVX := platform.StandardGamepadAxisValue(id, platform.StandardGamepadAxisLeftStickHorizontal)
			padVY := platform.StandardGamepadAxisValue(id, platform.StandardGamepadAxisLeftStickVertical)

			// Deadzone
			if math.Abs(padVX) < 0.1 {
				padVX = 0
			}
			if math.Abs(padVY) < 0.1 {
				padVY = 0
			}

			// Virtual digital overrides
			if platform.IsGamepadLeft(id) {
				padVX = -1
			} else if platform.IsGamepadRight(id) {
				padVX = 1
			}
			if platform.IsGamepadUp(id) {
				padVY = -1
			} else if platform.IsGamepadDown(id) {
				padVY = 1
			}

			if padVX != 0 || padVY != 0 {
				vx = moveSpeed * padVX
				vy = moveSpeed * padVY
				break // take first active gamepad
			}
		}

		/* ----------------------------- Apply velocity ----------------------------- */
		v.VX, v.VY = vx, vy

		/* ---------------------------- Sprite orientation --------------------------- */
		if hasSprite {
			if v.VX != 0 || v.VY != 0 {
				sprite.Rotation = math.Atan2(v.VY, v.VX)
			}
			if v.VX < 0 {
				sprite.FlipHorizontal = true
			} else if v.VX > 0 {
				sprite.FlipHorizontal = false
			}
		}
	}
}

func (s *System) Draw(*ecs.World, *platform.Image) {}

