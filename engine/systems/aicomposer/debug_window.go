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

/*───────────────────────────────────────────────*
| DEBUG WINDOW STRUCTURE                        |
*───────────────────────────────────────────────*/

// DebugWindow displays a real-time list of AI-composed entities
// and their currently bound AI actions.
type DebugWindow struct {
	component *window.Component
	content   *ComposerDebugContent
	visible   bool
}

/*───────────────────────────────────────────────*
| CONSTRUCTION / VISIBILITY                     |
*───────────────────────────────────────────────*/

// NewDebugWindow constructs a new AI Composer debug overlay.
func NewDebugWindow() *DebugWindow {
	return &DebugWindow{
		content: &ComposerDebugContent{
			lineHeight:     16,
			baselineOffset: 12,
		},
		visible: true,
	}
}

// Ensure adds the window to the ECS world if it hasn’t been added yet.
func (w *DebugWindow) Ensure(world *ecs.World) {
	if !w.visible || w.component != nil || world == nil {
		return
	}

	entity := world.NewEntity()
	comp := window.NewComponent(
		"debug.aicomposer",
		"AI Composer",
		window.Bounds{X: 32, Y: 320, Width: 380, Height: 160},
		w.content,
	)

	comp.Layer = ecs.LayerDebug
	comp.Order = 50
	comp.Movable = true
	comp.Closable = true
	comp.TitleBarHeight = 24
	comp.Padding = 8
	comp.Background = color.RGBA{18, 20, 30, 230}
	comp.Border = color.RGBA{80, 120, 255, 180}
	comp.TitleBar = color.RGBA{30, 50, 100, 220}
	comp.TitleColor = color.RGBA{235, 245, 255, 255}

	entity.Add(comp)
	w.component = comp
}

// Update refreshes the window’s content from the composer system.
func (w *DebugWindow) Update(world *ecs.World, composer *System) {
	if !w.visible || w.component == nil {
		return
	}
	w.content.Refresh(world, composer)
	w.component.Bounds.Height = w.content.estimateHeight()
}

// Hide toggles visibility off.
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

/*───────────────────────────────────────────────*
| REFRESH LOGIC                                 |
*───────────────────────────────────────────────*/

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

	lines := []string{
		"AI Composer Active Entities:",
		"--------------------------------",
	}

	// Collect and sort processed IDs
	var ids []ecs.EntityID
	for id := range composer.processed {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	for _, id := range ids {
		entity := world.GetEntity(id)
		if entity == nil {
			continue
		}

		act, _ := entity.Get("Actor").(*ecs.Actor)
		ctrl, _ := entity.Get("AIController").(*ecs.AIController)
		if act == nil || ctrl == nil {
			continue
		}

		lines = append(lines, fmt.Sprintf("[%3d] %-18s (%d actions)", entity.ID, act.ID, len(ctrl.Actions)))
		for _, a := range ctrl.Actions {
			lines = append(lines, fmt.Sprintf("   • %s [%s]", a.Name, a.Type))
		}
	}

	c.lines = lines
}

/*───────────────────────────────────────────────*
| DRAW FUNCTION                                 |
*───────────────────────────────────────────────*/

func (c *ComposerDebugContent) Draw(_ *ecs.World, dst *platform.Image, bounds window.Bounds) {
	if dst == nil || len(c.lines) == 0 {
		return
	}

	y := bounds.Y + c.baselineOffset
	for _, line := range c.lines {
		platform.DrawText(dst, line, basicfont.Face7x13, bounds.X+4, y, color.RGBA{220, 235, 255, 255})
		y += c.lineHeight
	}
}

/*───────────────────────────────────────────────*
| SIZE HELPER                                   |
*───────────────────────────────────────────────*/

func (c *ComposerDebugContent) estimateHeight() int {
	h := len(c.lines)*c.lineHeight + 32
	if h < 64 {
		return 64
	}
	return h
}
