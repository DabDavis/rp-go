package entitylist

import (
	"fmt"
	"image/color"

	"golang.org/x/image/font/basicfont"
	"rp-go/engine/ecs"
	"rp-go/engine/platform"
	"rp-go/engine/systems/actor"
)

// System renders a lightweight overlay showing the currently registered
// actors and their coordinates.
type System struct {
	registry *actor.Registry
}

// NewSystem constructs an entity list overlay backed by the actor registry.
func NewSystem(registry *actor.Registry) *System {
	return &System{registry: registry}
}

func (s *System) Layer() ecs.DrawLayer { return ecs.LayerEntityList }

func (s *System) Update(*ecs.World) {}

func (s *System) Draw(w *ecs.World, screen *platform.Image) {
	if screen == nil {
		return
	}
	entries := s.collectEntries(w)
	if len(entries) == 0 {
		return
	}

	bounds := screen.Bounds()
	width := 280
	height := len(entries)*16 + 32
	if height < 120 {
		height = 120
	}

	overlay := platform.NewImage(width, height)
	overlay.FillRect(0, 0, width, height, color.RGBA{0, 0, 0, 180})

	y := 20
	for _, line := range entries {
		platform.DrawText(overlay, line, basicfont.Face7x13, 12, y, color.White)
		y += 16
	}

	op := platform.NewDrawImageOptions()
	op.Translate(float64(bounds.Dx()-width-16), 16)
	screen.DrawImage(overlay, op)
}

func (s *System) collectEntries(w *ecs.World) []string {
	var entities []*ecs.Entity
	if s.registry != nil {
		entities = s.registry.All()
	}
	if len(entities) == 0 {
		return nil
	}

	entries := make([]string, 0, len(entities))
	for _, entity := range entities {
		actorComp, _ := entity.Get("Actor").(*ecs.Actor)
		if actorComp == nil {
			continue
		}
		pos, _ := entity.Get("Position").(*ecs.Position)
		if pos != nil {
			entries = append(entries, fmt.Sprintf("%s  (%.1f, %.1f)", actorComp.ID, pos.X, pos.Y))
		} else {
			entries = append(entries, fmt.Sprintf("%s  (-, -)", actorComp.ID))
		}
	}
	return entries
}
