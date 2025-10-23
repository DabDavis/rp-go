package debug

import (
	"fmt"
	"image/color"
	"sort"

	"golang.org/x/image/font/basicfont"

	"rp-go/engine/ecs"
	"rp-go/engine/platform"
	"rp-go/engine/ui/window"
)

// EntitiesContent lists entities with position and sprite info.
type EntitiesContent struct {
	lines           []string
	lineHeight      int
	baselineOffset  int
	lastLoggedFrame int
}

// Refresh rebuilds the list of visible entities.
func (c *EntitiesContent) Refresh(world *ecs.World, frame int, maxLines int) {
	if world == nil {
		c.lines = []string{"No world context"}
		return
	}

	lines := make([]string, 0, maxLines)
	manager := world.EntitiesManager()
	if manager != nil {
		collected := make([]string, 0, manager.Count())
		manager.ForEach(func(e *ecs.Entity) {
			pos, _ := e.Get("Position").(*ecs.Position)
			sprite, _ := e.Get("Sprite").(*ecs.Sprite)
			if pos == nil {
				return
			}

			spriteSize := "?"
			if sprite != nil && sprite.Image != nil {
				w, h := sprite.NativeSize()
				scale := sprite.PixelScale()
				// Fix: w and h are float64 — use %.0f for cleaner formatting
				spriteSize = fmt.Sprintf("%.0fx%.0f x%.2f", w, h, scale)
			}

			collected = append(collected,
				fmt.Sprintf("%d @ (%.0f, %.0f) %s", e.ID, pos.X, pos.Y, spriteSize))
		})

		sort.Strings(collected)
		if len(collected) > maxLines {
			collected = collected[:maxLines]
		}
		lines = append(lines, collected...)

		// Log once every few frames
		if frame-c.lastLoggedFrame >= 15 {
			for _, line := range collected {
				fmt.Println("[ENTITY DEBUG]", line)
			}
			c.lastLoggedFrame = frame
		}
	}

	if len(lines) == 0 {
		lines = append(lines, "No entities with Position component")
	}

	c.lines = lines
}

// Draw renders the entity diagnostics text.
func (c *EntitiesContent) Draw(_ *ecs.World, canvas *platform.Image, bounds window.Bounds) {
	if canvas == nil || len(c.lines) == 0 {
		return
	}

	baseline := bounds.Y + c.baselineOffset
	textX := bounds.X
	if textX < 0 {
		textX = 0
	}

	for _, line := range c.lines {
		platform.DrawText(canvas, line, basicfont.Face7x13, textX, baseline,
			color.RGBA{200, 220, 255, 255})
		baseline += c.lineHeight
	}
}

