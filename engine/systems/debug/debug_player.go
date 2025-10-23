package debug

import (
	"fmt"
	"image/color"

	"golang.org/x/image/font/basicfont"
	"rp-go/engine/ecs"
	"rp-go/engine/platform"
)

func DrawPlayerInfo(w *ecs.World, screen *platform.Image) {
	var playerPos *ecs.Position
	if manager := w.EntitiesManager(); manager != nil {
		manager.ForEach(func(e *ecs.Entity) {
			if playerPos != nil {
				return
			}
			if e.Has("CameraTarget") {
				if pos, ok := e.Get("Position").(*ecs.Position); ok {
					playerPos = pos
				}
			}
		})
	}
	if playerPos == nil {
		return
	}

	msg := fmt.Sprintf(
		"Player: (%.1f, %.1f)\n",
		playerPos.X, playerPos.Y,
	)
	platform.DrawText(screen, msg, basicfont.Face7x13, 10, 60, color.RGBA{200, 255, 200, 255})
}
