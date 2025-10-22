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

	minScale := cam.MinScale
	if minScale <= 0 {
		minScale = cam.Scale
	}
	maxScale := cam.MaxScale
	if maxScale <= 0 {
		maxScale = cam.Scale
	}
	targetScale := cam.TargetScale
	if targetScale <= 0 {
		targetScale = cam.Scale
	}

	msg := fmt.Sprintf(
		"Camera: (%.1f, %.1f)\nScale: %.2f → %.2f\nBounds: %.2f–%.2f\nViewport: %.0fx%.0f\n",
		cam.X, cam.Y, cam.Scale, targetScale, minScale, maxScale, sw, sh,
	)

	text.Draw(screen, msg, basicfont.Face7x13, 10, 20, color.White)
}
