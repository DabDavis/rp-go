package debug

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"rp-go/engine/ecs"
)

func DrawEntityDiagnostics(w *ecs.World, screen *ebiten.Image, frame int) {
	y := 100

	for _, e := range w.Entities {
		pos, ok1 := e.Get("Position").(*ecs.Position)
		sprite, ok2 := e.Get("Sprite").(*ecs.Sprite)
		if !ok1 || !ok2 || sprite.Image == nil {
			continue
		}

		imgW := float64(sprite.Image.Bounds().Dx())
		imgH := float64(sprite.Image.Bounds().Dy())
		playerScale := float64(sprite.Width) / imgW
		totalScale := playerScale // (cam.Scale handled in render system)

		entityInfo := fmt.Sprintf(
			"Entity %d | World(%.1f, %.1f) | Sprite %.0fx%.0f | Scale: %.2f",
			e.ID, pos.X, pos.Y, imgW, imgH, totalScale,
		)

		text.Draw(screen, entityInfo, basicfont.Face7x13, 10, y, color.RGBA{180, 220, 255, 255})
		y += 14

		if frame%15 == 0 {
			fmt.Println("[ENTITY DEBUG]", entityInfo)
		}
	}
}

