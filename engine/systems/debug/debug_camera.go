package debug

import (
	"fmt"
	"image/color"

	"golang.org/x/image/font/basicfont"
	"rp-go/engine/ecs"
	"rp-go/engine/platform"
)

func DrawCameraInfo(w *ecs.World, screen *platform.Image) {
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

	bounds := screen.Bounds()
	sw := float64(bounds.Dx())
	sh := float64(bounds.Dy())

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
		"Camera: (%.1f, %.1f)\nScale: %.2f -> %.2f\nBounds: %.2f-%.2f\nViewport: %.0fx%.0f\n",
		cam.X, cam.Y, cam.Scale, targetScale, minScale, maxScale, sw, sh,
	)

	platform.DrawText(screen, msg, basicfont.Face7x13, 10, 20, color.White)
}
