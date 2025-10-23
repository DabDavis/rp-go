package hud

import (
	"fmt"
	"image/color"

	"golang.org/x/image/font/basicfont"

	"rp-go/engine/ecs"
	"rp-go/engine/platform"
)

// System renders a lightweight heads-up display with player stats and controls.
type System struct{}

// Layer ensures HUD draws before other overlay diagnostics but after world content.
func (s *System) Layer() ecs.DrawLayer { return ecs.LayerHUD }

func (s *System) Update(*ecs.World) {}

func (s *System) Draw(w *ecs.World, screen *platform.Image) {
	if screen == nil {
		return
	}

	var position *ecs.Position
	var velocity *ecs.Velocity

	for _, entity := range w.Entities {
		actor, _ := entity.Get("Actor").(*ecs.Actor)
		if actor == nil || actor.ID != "player" {
			continue
		}
		position, _ = entity.Get("Position").(*ecs.Position)
		velocity, _ = entity.Get("Velocity").(*ecs.Velocity)
		break
	}

	title := "Pilot HUD"
	lines := []string{
		"Controls:",
		"  Move: WASD / Arrow Keys",
		"  Zoom: Mouse Wheel or +/-",
		"  Reset Zoom: 0",
		"  Toggle Console: F12",
	}

	if position != nil {
		lines = append(lines, fmt.Sprintf("Player Position: (%.0f, %.0f)", position.X, position.Y))
	}
	if velocity != nil {
		lines = append(lines, fmt.Sprintf("Velocity: (%.1f, %.1f)", velocity.VX, velocity.VY))
	}

	lineHeight := 16
	padding := 12
	titleBaseline := padding + 10

	width := 260
	height := padding*2 + lineHeight + len(lines)*lineHeight

	overlay := platform.NewImage(width, height)
	overlay.FillRect(0, 0, width, height, color.RGBA{0, 0, 0, 170})

	platform.DrawText(overlay, title, basicfont.Face7x13, padding, titleBaseline, color.RGBA{200, 220, 255, 255})

	y := titleBaseline + lineHeight
	for _, line := range lines {
		platform.DrawText(overlay, line, basicfont.Face7x13, padding, y, color.White)
		y += lineHeight
	}

	op := platform.NewDrawImageOptions()
	op.Translate(16, 16)
	screen.DrawImage(overlay, op)
}
