package aicomposer

import (
	"fmt"
	"image/color"
	"sort"

	"golang.org/x/image/font/basicfont"

	"rp-go/engine/ecs"
	"rp-go/engine/platform"
	"rp-go/engine/ui/window"
)

// DebugWindow displays a real-time list of AIComposed entities and their behaviors.
type DebugWindow struct {
	component *window.Component
	content   *ComposerDebugContent
	visible   bool
}

// NewDebugWindow constructs a new AIComposer debug overlay.
func NewDebugWindow() *DebugWindow {
	return &DebugWindow{
		content: &ComposerDebugContent{
			lineHeight:     16,
			baselineOffset: 12,
		},
		visible: true,
	}
}

// Ensure adds the window to the ECS world if missing.
func (w *DebugWindow) Ensure(world *ecs.World) {
	if !w.visible || w.component != nil {
		return
	}

	entity := world.NewEntity()
	comp := window.NewComponent("debug.aicomposer", "AI Composer", window.Bounds{
		X:      32,
		Y:      320,
		Width:  360,
		Height: 160,
	}, w.content)

	comp.Layer = ecs.LayerDebug
	comp.Order = 50
	comp.Movable = true
	comp.Closable = true
	comp.TitleBarHeight = 26
	comp.Padding = 10
	comp.Background = color.RGBA{18, 20, 30, 220}
	comp.Border = color.RGBA{100, 140, 255, 180}
	comp.TitleBar = color.RGBA{40, 60, 110, 230}
	comp.TitleColor = color.RGBA{235, 245, 255, 255}

	entity.Add(comp)
	w.component = comp
}

// Update refreshes the content text from ECS world data.
func (w *DebugWindow) Update(world *ecs.World, composer *System) {
	if !w.visible || w.component == nil {
		return
	}
	w.content.Refresh(world, composer)
	w.component.Bounds.Height = w.content.estimateHeight()
}

// Hide toggles visibility.
func (w *DebugWindow) Hide() {
	w.visible = false
	if w.component != nil {
		w.component.Visible = false
	}
}

/*───────────────────────────────────────────────*
 | CONTENT RENDERER                              |
 *───────────────────────────────────────────────*/

type ComposerDebugContent struct {
	lines          []string
	lineHeight     int
	baselineOffset int
}

// Refresh rebuilds the entity/action list from the composer state.
func (c *ComposerDebugContent) Refresh(world *ecs.World, composer *System) {
	if composer == nil {
		c.lines = []string{"(no composer system)"}
		return
	}
	if world == nil {
		c.lines = []string{"(no world context)"}
		return
	}

	composer.mu.RLock()
	defer composer.mu.RUnlock()

	if len(composer.processed) == 0 {
		c.lines = []string{"No AI-composed entities"}
		return
	}

	lines := make([]string, 0, len(composer.processed)+4)
	lines = append(lines, "AI Composer Active Entities:")
	lines = append(lines, "--------------------------------")

	// Sort IDs for deterministic order
	var ids []ecs.EntityID
	for id := range composer.processed {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	for _, id := range ids {
		e := world.GetEntity(id)
		if e == nil {
			continue
		}

		act, _ := e.Get("Actor").(*ecs.Actor)
		ctrl, _ := e.Get("AIController").(*AIController)
		if act == nil || ctrl == nil {
			continue
		}

		line := fmt.Sprintf("[%3d] %-20s  (%d actions)", id, act.ID, len(ctrl.Actions))
		lines = append(lines, line)

		for _, a := range ctrl.Actions {
			lines = append(lines, fmt.Sprintf("   • %s [%s]", a.Name, a.Type))
		}
	}

	c.lines = lines
}

// Draw renders the text into the debug window.
func (c *ComposerDebugContent) Draw(_ *ecs.World, dst *platform.Image, bounds window.Bounds) {
	if dst == nil || len(c.lines) == 0 {
		return
	}

	baseline := bounds.Y + c.baselineOffset
	for _, line := range c.lines {
		platform.DrawText(dst, line, basicfont.Face7x13, bounds.X+4, baseline, color.RGBA{220, 235, 255, 255})
		baseline += c.lineHeight
	}
}

// estimateHeight computes the needed height for all lines.
func (c *ComposerDebugContent) estimateHeight() int {
	min := 48
	h := len(c.lines)*c.lineHeight + 32
	if h < min {
		return min
	}
	return h
}

