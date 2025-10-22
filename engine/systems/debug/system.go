package debug

import (
	"fmt"
	"image/color"
	"strings"

	"golang.org/x/image/font/basicfont"
	"rp-go/engine/ecs"
	"rp-go/engine/platform"
)

type System struct {
	overlayText string
}

func (s *System) Update(*ecs.World) {}

func (s *System) Draw(w *ecs.World, screen *platform.Image) {
	DrawDebugGrid(w, screen)
	s.overlayText = buildOverlay(w)
}

func (s *System) DrawOverlay(w *ecs.World, screen *platform.Image) {
	if s.overlayText == "" {
		s.overlayText = buildOverlay(w)
	}
	platform.DrawText(screen, s.overlayText, basicfont.Face7x13, 10, 20, color.White)
}

func buildOverlay(w *ecs.World) string {
	var cam *ecs.Camera
	for _, e := range w.Entities {
		if c, ok := e.Get("Camera").(*ecs.Camera); ok {
			cam = c
			break
		}
	}

	var builder strings.Builder
	fmt.Fprintf(&builder, "FPS: %.0f\nEntities: %d", platform.ActualFPS(), len(w.Entities))

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
		fmt.Fprintf(&builder, "Cam: (%.1f, %.1f)\n", cam.X, cam.Y)
		fmt.Fprintf(&builder, "Rotation: %.1f deg\n", cam.Rotation*180/3.14159)
		fmt.Fprintf(&builder, "Scale: %.2f -> %.2f\n", cam.Scale, targetScale)
		fmt.Fprintf(&builder, "Bounds: %.2f - %.2f\n", minScale, maxScale)
		fmt.Fprintf(&builder, "Default Scale: %.2f", defaultScale)
	}

	return builder.String()
}
