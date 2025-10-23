package debug

import (
	"image/color"
	"rp-go/engine/ecs"
	"rp-go/engine/ui/window"
)

type EntitiesWindow struct {
	component *window.Component
	content   *EntitiesContent
	cfg       Config
	frame     int
}

func NewEntitiesWindow(cfg Config) *EntitiesWindow {
	return &EntitiesWindow{
		cfg:     cfg,
		content: &EntitiesContent{lineHeight: 16, baselineOffset: 12},
	}
}

func (w *EntitiesWindow) Ensure(world *ecs.World) {
	if w.component != nil {
		return
	}

	entity := world.NewEntity()
	comp := window.NewComponent("debug.entities", "Entity Diagnostics", window.Bounds{
		X:      w.cfg.ViewportWidth - w.cfg.EntitiesWidth - w.cfg.Margin,
		Y:      w.cfg.Margin,
		Width:  w.cfg.EntitiesWidth,
		Height: 120,
	}, w.content)
	comp.Layer = ecs.LayerDebug
	comp.Order = 20
	comp.Padding = 12
	comp.TitleBarHeight = 26
	comp.Background = color.RGBA{12, 16, 24, 220}
	comp.Border = color.RGBA{120, 160, 235, 200}
	comp.TitleBar = color.RGBA{36, 58, 110, 230}
	comp.TitleColor = color.RGBA{235, 245, 255, 255}
	comp.Movable = true
	comp.Closable = true

	entity.Add(comp)
	w.component = comp
}

func (w *EntitiesWindow) Update(world *ecs.World) {
	if w.component == nil {
		return
	}
	w.frame++
	w.content.Refresh(world, w.frame, w.cfg.MaxEntities)
	lines := len(w.content.lines)
	padding := w.component.Padding
	title := w.component.TitleBarHeight
	total := title + padding*2 + lines*w.content.lineHeight
	if total < title+padding*2+w.content.lineHeight {
		total = title + padding*2 + w.content.lineHeight
	}
	w.component.Bounds.Height = total
}

func (w *EntitiesWindow) Hide(world *ecs.World) {
	if w.component != nil {
		w.component.Visible = false
	}
}

