package input

import (
	"math"

	"rp-go/engine/ecs"
	"rp-go/engine/platform"
	"rp-go/engine/systems/devconsole"
)

type System struct{}

// Default movement speed (units per frame)
const moveSpeed = 3.0

func (s *System) Update(w *ecs.World) {
	for _, e := range w.Entities {
		controller, hasController := e.Get("PlayerInput").(*ecs.PlayerInput)
		if !hasController || controller == nil || !controller.Enabled {
			continue
		}

		v, ok := e.Get("Velocity").(*ecs.Velocity)
		if !ok {
			continue
		}

		sprite, hasSprite := e.Get("Sprite").(*ecs.Sprite)

		// Local velocity deltas
		vx, vy := 0.0, 0.0

		/* ---------------------------- Keyboard movement --------------------------- */
		if platform.IsKeyPressed(platform.KeyArrowLeft) || platform.IsKeyPressed(platform.KeyA) {
			vx -= 1
		}
		if platform.IsKeyPressed(platform.KeyArrowRight) || platform.IsKeyPressed(platform.KeyD) {
			vx += 1
		}
		if platform.IsKeyPressed(platform.KeyArrowUp) || platform.IsKeyPressed(platform.KeyW) {
			vy -= 1
		}
		if platform.IsKeyPressed(platform.KeyArrowDown) || platform.IsKeyPressed(platform.KeyS) {
			vy += 1
		}

		/* ----------------------------- Gamepad movement --------------------------- */
		for _, id := range platform.GamepadIDs() {
			if !platform.IsStandardGamepadLayoutAvailable(id) {
				continue
			}

			// Analog stick input
			padVX := platform.StandardGamepadAxisValue(id, platform.StandardGamepadAxisLeftStickHorizontal)
			padVY := platform.StandardGamepadAxisValue(id, platform.StandardGamepadAxisLeftStickVertical)

			// Deadzone
			if math.Abs(padVX) < 0.1 {
				padVX = 0
			}
			if math.Abs(padVY) < 0.1 {
				padVY = 0
			}

			// Virtual D-pad override
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
				vx = padVX
				vy = padVY
				break
			}
		}

		/* ----------------------------- Normalize motion --------------------------- */
		// Prevent faster diagonal movement (e.g. √2 speed)
		mag := math.Hypot(vx, vy)
		if mag > 1 {
			vx /= mag
			vy /= mag
		}

		v.VX = vx * moveSpeed
		v.VY = vy * moveSpeed

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
