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

// SystemInspectorContent lists all registered ECS systems with their layers.
type SystemInspectorContent struct {
	lines          []string
	lineHeight     int
	baselineOffset int
}

// Refresh collects system info from the ECS world.
func (c *SystemInspectorContent) Refresh(world *ecs.World) {
	if world == nil {
		c.lines = []string{"No world context"}
		return
	}

	if len(world.Systems) == 0 {
		c.lines = []string{"No systems registered"}
		return
	}

	type sysEntry struct {
		Name  string
		Layer ecs.DrawLayer
	}

	var list []sysEntry
	for _, s := range world.Systems {
		name := ecs.SystemName(s)
		layer := ecs.LayerNone
		if layered, ok := s.(ecs.LayeredSystem); ok {
			layer = layered.Layer()
		}
		list = append(list, sysEntry{Name: name, Layer: layer})
	}

	sort.SliceStable(list, func(i, j int) bool {
		if list[i].Layer == list[j].Layer {
			return list[i].Name < list[j].Name
		}
		return list[i].Layer < list[j].Layer
	})

	lines := []string{
		fmt.Sprintf("Systems: %d", len(list)),
		"──────────────────────────────",
	}

	for _, e := range list {
		layerName := fmt.Sprintf("%d", e.Layer)
		switch e.Layer {
		case ecs.LayerBackground:
			layerName = "Background"
		case ecs.LayerWorld:
			layerName = "World"
		case ecs.LayerForeground:
			layerName = "Foreground"
		case ecs.LayerHUD:
			layerName = "HUD"
		case ecs.LayerEntityList:
			layerName = "EntityList"
		case ecs.LayerDebug:
			layerName = "Debug"
		case ecs.LayerConsole:
			layerName = "Console"
		case ecs.LayerNone:
			layerName = "-"
		}

		lines = append(lines, fmt.Sprintf("[%s] %s", layerName, e.Name))
	}

	c.lines = lines
}

// Draw renders the system list text to the window.
func (c *SystemInspectorContent) Draw(_ *ecs.World, canvas *platform.Image, bounds window.Bounds) {
	if canvas == nil || len(c.lines) == 0 {
		return
	}
	baseline := bounds.Y + c.baselineOffset
	textX := bounds.X
	if textX < 0 {
		textX = 0
	}
	for _, line := range c.lines {
		platform.DrawText(canvas, line, basicfont.Face7x13, textX, baseline, color.RGBA{210, 230, 255, 255})
		baseline += c.lineHeight
	}
}

