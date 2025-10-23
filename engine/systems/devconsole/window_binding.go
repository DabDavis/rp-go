package devconsole

import (
	"rp-go/engine/ecs"
	"rp-go/engine/ui/window"
)

func (s *ConsoleState) ensureWindow(world *ecs.World, cfg Config) {
	if s == nil || world == nil {
		return
	}
	if s.windowComponent != nil && s.windowEntity != nil && s.windowEntity.Has("Window") {
		return
	}

	content := newConsoleWindowContent(s)
	bounds := window.Bounds{X: cfg.Margin, Y: cfg.Margin, Width: cfg.MinWidth, Height: cfg.MinHeight}
	component := window.NewComponent("console.dev", "Developer Console", bounds, content)
	component.Layer = ecs.LayerConsole
	component.Visible = false
	component.Order = 200
	component.Padding = 12
	component.TitleBarHeight = 28

	entity := world.NewEntity()
	entity.Add(component)

	s.windowEntity = entity
	s.windowComponent = component
	s.windowContent = content
}

func (s *ConsoleState) syncWindowVisibility() {
	if s == nil || s.windowComponent == nil {
		return
	}
	s.windowComponent.Visible = s.Open
}

func (s *ConsoleState) applyLayout(cfg Config) {
	if s == nil || s.windowComponent == nil {
		return
	}
	margin := cfg.Margin
	if margin < 0 {
		margin = 0
	}
	availableWidth := cfg.ViewportWidth - margin*2
	if availableWidth <= 0 {
		availableWidth = cfg.ViewportWidth
	}
	width := availableWidth
	if width < cfg.MinWidth {
		width = cfg.MinWidth
	}
	if cfg.ViewportWidth > 0 && width > cfg.ViewportWidth {
		width = cfg.ViewportWidth
	}
	if width < 1 {
		width = 1
	}

	height := int(float64(cfg.ViewportHeight) * cfg.HeightRatio)
	if height < cfg.MinHeight {
		height = cfg.MinHeight
	}
	if cfg.ViewportHeight > 0 && height > cfg.ViewportHeight {
		height = cfg.ViewportHeight
	}
	if height < 1 {
		height = 1
	}

	y := cfg.ViewportHeight - height - margin
	if y < margin {
		y = margin
	}
	s.windowComponent.Bounds = window.Bounds{
		X:      margin,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

func (s *ConsoleState) promptString() string {
	if s == nil {
		return "> "
	}
	prompt := "> " + s.InputBuffer
	if (s.CursorTick/20)%2 == 0 {
		prompt += "_"
	}
	return prompt
}
