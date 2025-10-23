package input

import (
	"math"

	"rp-go/engine/platform"
)

// Manager centralizes polling game input devices and exposes normalized
// movement vectors for systems to consume.
type Manager struct {
	moveX float64
	moveY float64
}

var defaultManager = &Manager{}

// ManagerInstance returns the shared input manager used by the engine.
func ManagerInstance() *Manager {
	return defaultManager
}

// Poll refreshes the cached input state for the current frame.
func (m *Manager) Poll() {
	if m == nil {
		return
	}

	vx, vy := 0.0, 0.0

	// Keyboard
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

	// Gamepad
	for _, id := range platform.GamepadIDs() {
		if !platform.IsStandardGamepadLayoutAvailable(id) {
			continue
		}

		padVX := platform.StandardGamepadAxisValue(id, platform.StandardGamepadAxisLeftStickHorizontal)
		padVY := platform.StandardGamepadAxisValue(id, platform.StandardGamepadAxisLeftStickVertical)

		if math.Abs(padVX) < 0.1 {
			padVX = 0
		}
		if math.Abs(padVY) < 0.1 {
			padVY = 0
		}

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

	mag := math.Hypot(vx, vy)
	if mag > 1 {
		vx /= mag
		vy /= mag
	}

	m.moveX = vx
	m.moveY = vy
}

// Movement returns the latest normalized movement vector sampled during Poll.
func (m *Manager) Movement() (float64, float64) {
	if m == nil {
		return 0, 0
	}
	return m.moveX, m.moveY
}
