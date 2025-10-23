package debug

import (
	"image/color"
	"rp-go/engine/ecs"
	"rp-go/engine/ui/window"
)

// SystemWindow displays all registered ECS systems and their layers.
type SystemWindow struct {
	component *window.Component
	content   *SystemInspectorContent
	cfg       Config
}

// NewSystemWindow creates the inspector window.
func NewSystemWindow(cfg Config) *SystemWindow {
	return &SystemWindow{
		cfg:     cfg,
		content: &SystemInspectorContent{lineHeight: 16, baselineOffset: 12},
	}
}

// Ensure creates the window entity once and adds it to the ECS world.
func (w *SystemWindow) Ensure(world *ecs.World) {
	if w.component != nil {
		return
	}

	entity := world.NewEntity()
	comp := window.NewComponent("debug.systems", "System Inspector", window.Bounds{
		X:      w.cfg.Margin,
		Y:      w.cfg.ViewportHeight - 200 - w.cfg.Margin,
		Width:  360,
		Height: 180,
	}, w.content)
	comp.Layer = ecs.LayerDebug
	comp.Order = 30
	comp.Padding = 12
	comp.TitleBarHeight = 26
	comp.Background = color.RGBA{10, 15, 25, 220}
	comp.Border = color.RGBA{140, 160, 255, 200}
	comp.TitleBar = color.RGBA{40, 60, 120, 230}
	comp.TitleColor = color.RGBA{230, 240, 255, 255}
	comp.Movable = true
	comp.Closable = true

	entity.Add(comp)
	w.component = comp
}

// Update refreshes content every frame.
func (w *SystemWindow) Update(world *ecs.World) {
	if w.component == nil {
		return
	}
	w.content.Refresh(world)
	lines := len(w.content.lines)
	padding := w.component.Padding
	title := w.component.TitleBarHeight
	total := title + padding*2 + lines*w.content.lineHeight
	if total < title+padding*2+w.content.lineHeight {
		total = title + padding*2 + w.content.lineHeight
	}
	w.component.Bounds.Height = total
}

// Hide toggles visibility.
func (w *SystemWindow) Hide(world *ecs.World) {
	if w.component != nil {
		w.component.Visible = false
	}
}

