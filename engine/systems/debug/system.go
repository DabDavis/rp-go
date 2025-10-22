package debug

import (
	"fmt"
	"image/color"
	"strings"

	"golang.org/x/image/font/basicfont"
	"rp-go/engine/ecs"
	"rp-go/engine/platform"
)

type System struct{}

func (s *System) Update(*ecs.World) {}

func (s *System) Draw(w *ecs.World, screen *platform.Image) {
	var cam *ecs.Camera
	for _, e := range w.Entities {
		if c, ok := e.Get("Camera").(*ecs.Camera); ok {
			cam = c
			break
		}
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("FPS: %.0f\nEntities: %d", platform.ActualFPS(), len(w.Entities)))
	builder.WriteString(fmt.Sprintf("FPS: %.0f\nEntities: %d", platform.ActualFPS(), len(w.Entities)))
	builder.WriteString(fmt.Sprintf("FPS: %.0f\nEntities: %d", platform.ActualFPS(), len(w.Entities)))

	if cam != nil {
		targetScale := cam.TargetScale
		if targetScale <= 0 {
			targetScale = cam.Scale
		}
		minScale := cam.MinScale
		if minScale <= 0 {
			minScale = cam.Scale
		}
		maxScale := cam.MaxScale
		if maxScale <= 0 {
			maxScale = cam.Scale
		}
		defaultScale := cam.DefaultScale
		if defaultScale <= 0 {
			defaultScale = cam.Scale
		}

		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("Cam: (%.1f, %.1f)\n", cam.X, cam.Y))
		builder.WriteString(fmt.Sprintf("Rotation: %.1f°\n", cam.Rotation*180/3.14159))
		builder.WriteString(fmt.Sprintf("Scale: %.2f → %.2f\n", cam.Scale, targetScale))
		builder.WriteString(fmt.Sprintf("Bounds: %.2f – %.2f\n", minScale, maxScale))
		builder.WriteString(fmt.Sprintf("Default Scale: %.2f", defaultScale))
	}

	platform.DrawText(screen, builder.String(), basicfont.Face7x13, 10, 20, color.White)
	platform.DrawText(screen, builder.String(), basicfont.Face7x13, 10, 20, color.White)
	platform.DrawText(screen, builder.String(), basicfont.Face7x13, 10, 20, color.White)
}
