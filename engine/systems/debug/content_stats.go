package debug

import (
	"fmt"
	"image/color"

	"golang.org/x/image/font/basicfont"

	"rp-go/engine/ecs"
	"rp-go/engine/platform"
	"rp-go/engine/ui/window"
)

// StatsContent renders frame rate, entity counts, and camera/player info.
type StatsContent struct {
	lines          []string
	lineHeight     int
	baselineOffset int
}

// Refresh gathers live debug stats from the ECS world.
func (c *StatsContent) Refresh(world *ecs.World) {
	if world == nil {
		c.lines = []string{"No world context"}
		return
	}

	lines := make([]string, 0, 8)
	fps := platform.ActualFPS()
	lines = append(lines, fmt.Sprintf("FPS: %.0f", fps))

	manager := world.EntitiesManager()
	if manager != nil {
		lines = append(lines, fmt.Sprintf("Entities: %d", manager.Count()))
	}

	// Camera info
	var cam *ecs.Camera
	if manager != nil {
		if _, comp := manager.FirstComponent("Camera"); comp != nil {
			cam, _ = comp.(*ecs.Camera)
		}
	}
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
		lines = append(lines,
			fmt.Sprintf("Camera: (%.1f, %.1f)", cam.X, cam.Y),
			fmt.Sprintf("Scale: %.2f → %.2f", cam.Scale, targetScale),
			fmt.Sprintf("Bounds: %.2f – %.2f", minScale, maxScale),
		)
	}

	// Player info
	var playerPos *ecs.Position
	if manager != nil {
		manager.ForEach(func(entity *ecs.Entity) {
			if playerPos != nil {
				return
			}
			if entity.Has("CameraTarget") {
				if pos, ok := entity.Get("Position").(*ecs.Position); ok {
					playerPos = pos
				}
			}
		})
	}
	if playerPos != nil {
		lines = append(lines, fmt.Sprintf("Player: (%.1f, %.1f)", playerPos.X, playerPos.Y))
	}

	if len(lines) == 0 {
		lines = append(lines, "No debug data available")
	}

	c.lines = lines
}

// Draw renders the debug text into the window’s content area.
func (c *StatsContent) Draw(_ *ecs.World, canvas *platform.Image, bounds window.Bounds) {
	if canvas == nil || len(c.lines) == 0 {
		return
	}
	baseline := bounds.Y + c.baselineOffset
	textX := bounds.X
	if textX < 0 {
		textX = 0
	}
	for _, line := range c.lines {
		platform.DrawText(canvas, line, basicfont.Face7x13, textX, baseline, color.White)
		baseline += c.lineHeight
	}
}

