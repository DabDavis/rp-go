package input

import (
	"math"

	"rp-go/engine/ecs"
	inputmgr "rp-go/engine/input"
	"rp-go/engine/systems/devconsole"
)

// System processes player input from keyboard and gamepad,
// updating Velocity and Sprite components accordingly.
type System struct{}

// Default movement speed in world units per frame.
const moveSpeed = 3.0

// Update polls input devices and applies movement to entities
// with PlayerInput and Velocity components. When the developer
// console is open, all player input is ignored.
func (s *System) Update(w *ecs.World) {
	if w == nil {
		return
	}

	// Skip input when the developer console is open.
	if devconsole.IsOpen() {
		return
	}

	manager := w.EntitiesManager()
	if manager == nil {
		return
	}

	inputManager := inputmgr.ManagerInstance()
	inputManager.Poll()
	moveX, moveY := inputManager.Movement()

	manager.ForEach(func(e *ecs.Entity) {
		controller, hasController := e.Get("PlayerInput").(*ecs.PlayerInput)
		if !hasController || controller == nil || !controller.Enabled {
			return
		}

		v, ok := e.Get("Velocity").(*ecs.Velocity)
		if !ok {
			return
		}

		sprite, hasSprite := e.Get("Sprite").(*ecs.Sprite)

		v.VX = moveX * moveSpeed
		v.VY = moveY * moveSpeed

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
	})
}
