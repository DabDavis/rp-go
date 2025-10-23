package debug

import (
	"image/color"
	"rp-go/engine/ecs"
	"rp-go/engine/ui/window"
)

type StatsWindow struct {
	component *window.Component
	content   *StatsContent
	cfg       Config
}

func NewStatsWindow(cfg Config) *StatsWindow {
	return &StatsWindow{
		cfg:     cfg,
		content: &StatsContent{lineHeight: 16, baselineOffset: 12},
	}
}

func (s *StatsWindow) Ensure(world *ecs.World) {
	if s.component != nil {
		return
	}

	entity := world.NewEntity()
	comp := window.NewComponent("debug.stats", "Debug Stats", window.Bounds{
		X:      s.cfg.Margin,
		Y:      s.cfg.Margin,
		Width:  s.cfg.StatsWidth,
		Height: 120,
	}, s.content)
	comp.Layer = ecs.LayerDebug
	comp.Order = 10
	comp.Padding = 12
	comp.TitleBarHeight = 26
	comp.Background = color.RGBA{12, 16, 24, 220}
	comp.Border = color.RGBA{90, 130, 200, 200}
	comp.TitleBar = color.RGBA{30, 50, 90, 230}
	comp.TitleColor = color.RGBA{230, 240, 255, 255}
	comp.Movable = true
	comp.Closable = true

	entity.Add(comp)
	s.component = comp
}

func (s *StatsWindow) Update(world *ecs.World) {
	if s.component == nil {
		return
	}
	s.content.Refresh(world)
	lines := len(s.content.lines)
	padding := s.component.Padding
	title := s.component.TitleBarHeight
	total := title + padding*2 + lines*s.content.lineHeight
	if total < title+padding*2+s.content.lineHeight {
		total = title + padding*2 + s.content.lineHeight
	}
	s.component.Bounds.Height = total
}

func (s *StatsWindow) Hide(world *ecs.World) {
	if s.component != nil {
		s.component.Visible = false
	}
}

