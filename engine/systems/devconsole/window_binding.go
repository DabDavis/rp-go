package devconsole

import (
	"rp-go/engine/ecs"
	"rp-go/engine/platform"
	"rp-go/engine/ui/window"
)

const (
	consoleMargin = 16
)

func (s *ConsoleState) ensureWindow(world *ecs.World) {
	if s == nil || world == nil {
		return
	}
	if s.windowComponent != nil && s.windowEntity != nil && s.windowEntity.Has("Window") {
		return
	}

	content := newConsoleWindowContent(s)
	bounds := window.Bounds{X: consoleMargin, Y: consoleMargin, Width: 320, Height: 240}
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

func (s *ConsoleState) layoutWindow(screen *platform.Image) {
	if s == nil || s.windowComponent == nil || screen == nil {
		return
	}
	bounds := screen.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width <= 0 || height <= 0 {
		return
	}

	windowWidth := width - consoleMargin*2
	if windowWidth < 360 {
		windowWidth = 360
	}
	windowHeight := height / 3
	if windowHeight < 180 {
		windowHeight = 180
	}
	windowY := height - windowHeight - consoleMargin
	if windowY < consoleMargin {
		windowY = consoleMargin
	}

	s.windowComponent.Bounds = window.Bounds{
		X:      consoleMargin,
		Y:      windowY,
		Width:  windowWidth,
		Height: windowHeight,
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
