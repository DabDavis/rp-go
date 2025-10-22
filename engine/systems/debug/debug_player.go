package debug

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"rp-go/engine/ecs"
)

func DrawPlayerInfo(w *ecs.World, screen *ebiten.Image) {
	var playerPos *ecs.Position
	for _, e := range w.Entities {
		if e.Has("CameraTarget") {
			if pos, ok := e.Get("Position").(*ecs.Position); ok {
				playerPos = pos
				break
			}
		}
	}
	if playerPos == nil {
		return
	}

	msg := fmt.Sprintf(
		"Player: (%.1f, %.1f)\n",
		playerPos.X, playerPos.Y,
	)
	text.Draw(screen, msg, basicfont.Face7x13, 10, 60, color.RGBA{200, 255, 200, 255})
}

