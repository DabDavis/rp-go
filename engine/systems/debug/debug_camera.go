package debug

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"rp-go/engine/ecs"
)

func DrawCameraInfo(w *ecs.World, screen *ebiten.Image) {
	var cam *ecs.Camera
	for _, e := range w.Entities {
		if c, ok := e.Get("Camera").(*ecs.Camera); ok {
			cam = c
			break
		}
	}

	if cam == nil {
		return
	}

	sw := float64(screen.Bounds().Dx())
	sh := float64(screen.Bounds().Dy())
	
	msg := fmt.Sprintf(
	"Camera: (%.1f, %.1f)\nScale: %.2f (%.1f–%.1f)\nViewport: %.0fx%.0f\n",
	cam.X, cam.Y, cam.Scale, 0.5, 3.0, sw, sh,
	)

	text.Draw(screen, msg, basicfont.Face7x13, 10, 20, color.White)
}

