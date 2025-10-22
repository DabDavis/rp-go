package debug

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"rp-go/engine/ecs"
)

type System struct{}

func (s *System) Update(*ecs.World) {}

func (s *System) Draw(w *ecs.World, screen *ebiten.Image) {
	var cam *ecs.Camera
	for _, e := range w.Entities {
		if c, ok := e.Get("Camera").(*ecs.Camera); ok {
			cam = c
			break
		}
	}

	msg := fmt.Sprintf("FPS: %.0f\nEntities: %d", ebiten.ActualFPS(), len(w.Entities))
	if cam != nil {
		msg += fmt.Sprintf("\nCam(%.1f, %.1f) Scale %.2f Rot %.1f°",
			cam.X, cam.Y, cam.Scale, cam.Rotation*180/3.14159)
	}
	text.Draw(screen, msg, basicfont.Face7x13, 10, 20, color.White)
}

