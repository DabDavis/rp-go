package hud

import (
	"fmt"
	"image/color"

	"golang.org/x/image/font/basicfont"

	"rp-go/engine/ecs"
	"rp-go/engine/platform"
	"rp-go/engine/ui/window"
)

// System maintains the pilot HUD window and populates its content.
type System struct {
	windowEntity    *ecs.Entity
	windowComponent *window.Component
	content         *pilotHUDContent
}

// NewSystem creates a HUD system backed by the shared window manager.
func NewSystem() *System {
	return &System{
		content: &pilotHUDContent{
			lineHeight:     16,
			baselineOffset: 12,
		},
	}
}

// Update ensures the HUD window exists and refreshes its dynamic content.
func (s *System) Update(world *ecs.World) {
	if world == nil {
		return
	}

	if s.content == nil {
		s.content = &pilotHUDContent{
			lineHeight:     16,
			baselineOffset: 12,
		}
	}

	if s.windowComponent == nil || s.windowEntity == nil || !s.windowEntity.Has("Window") {
		s.attachWindow(world)
	}

	if s.windowComponent == nil {
		return
	}

	s.content.Refresh(world)

	padding := s.windowComponent.Padding
	if padding < 0 {
		padding = 0
	}
	titleBar := s.windowComponent.TitleBarHeight
	if titleBar < 0 {
		titleBar = 0
	}
	contentHeight := len(s.content.lines) * s.content.lineHeight
	minimum := titleBar + padding*2
	totalHeight := minimum + contentHeight
	if totalHeight < minimum {
		totalHeight = minimum
	}
	s.windowComponent.Bounds.Height = totalHeight
}

func (s *System) attachWindow(world *ecs.World) {
	if world == nil {
		return
	}
	entity := world.NewEntity()
	bounds := window.Bounds{X: 16, Y: 16, Width: 320, Height: 0}
	component := window.NewComponent("hud.pilot", "Pilot HUD", bounds, s.content)
	component.Order = 10
	component.Padding = 12
	component.TitleBarHeight = 24
	component.Background = color.RGBA{8, 12, 20, 200}
	component.Border = color.RGBA{120, 160, 255, 180}
	component.TitleBar = color.RGBA{25, 40, 92, 220}
	component.TitleColor = color.RGBA{235, 244, 255, 255}
	entity.Add(component)

	s.windowEntity = entity
	s.windowComponent = component
}

type pilotHUDContent struct {
	lines          []string
	lineHeight     int
	baselineOffset int
}

func (c *pilotHUDContent) Refresh(world *ecs.World) {
	lines := []string{
		"Controls:",
		"  Move: WASD / Arrow Keys",
		"  Zoom: Mouse Wheel or +/-",
		"  Reset Zoom: 0",
		"  Toggle Console: F12",
	}

	var position *ecs.Position
	var velocity *ecs.Velocity

	for _, entity := range world.Entities {
		if entity == nil {
			continue
		}
		actor, _ := entity.Get("Actor").(*ecs.Actor)
		if actor == nil || actor.ID != "player" {
			continue
		}
		position, _ = entity.Get("Position").(*ecs.Position)
		velocity, _ = entity.Get("Velocity").(*ecs.Velocity)
		break
	}

	if position != nil {
		lines = append(lines, fmt.Sprintf("Position: (%.0f, %.0f)", position.X, position.Y))
	}
	if velocity != nil {
		lines = append(lines, fmt.Sprintf("Velocity: (%.1f, %.1f)", velocity.VX, velocity.VY))
	}

	c.lines = lines
}

func (c *pilotHUDContent) Draw(world *ecs.World, canvas *platform.Image, bounds window.Bounds) {
	if canvas == nil {
		return
	}
	if len(c.lines) == 0 {
		return
	}

	baseline := bounds.Y + c.baselineOffset
	lineHeight := c.lineHeight
	textX := bounds.X
	if textX < 0 {
		textX = 0
	}

	for _, line := range c.lines {
		platform.DrawText(canvas, line, basicfont.Face7x13, textX, baseline, color.White)
		baseline += lineHeight
	}
}
