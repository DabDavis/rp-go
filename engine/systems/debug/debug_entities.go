package debug

import (
	"fmt"
	"image/color"

	"golang.org/x/image/font/basicfont"
	"rp-go/engine/ecs"
	"rp-go/engine/platform"
)

// DrawEntityDiagnostics overlays basic per-entity info.
func DrawEntityDiagnostics(w *ecs.World, screen *platform.Image, frame int) {
	y := 100

	for _, e := range w.Entities {
		pos, ok1 := e.Get("Position").(*ecs.Position)
		sprite, ok2 := e.Get("Sprite").(*ecs.Sprite)
		if !ok1 || !ok2 || sprite.Texture == nil {
			continue
		}

		imgW := float64(sprite.Image.Bounds().Dx())
		imgH := float64(sprite.Image.Bounds().Dy())
		spriteScale := float64(sprite.Width) / imgW

		entityInfo := fmt.Sprintf(
			"Entity %d | World(%.1f, %.1f) | Sprite %.0fx%.0f | Scale: %.2f",
			e.ID, pos.X, pos.Y, imgW, imgH, spriteScale,
		)

		platform.DrawText(screen, entityInfo, basicfont.Face7x13, 10, y, color.RGBA{180, 220, 255, 255})
		y += 14

		if frame%15 == 0 {
			fmt.Println("[ENTITY DEBUG]", entityInfo)
		}
	}
}
