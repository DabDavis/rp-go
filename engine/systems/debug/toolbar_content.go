package debug

import (
	"image/color"

	"golang.org/x/image/font/basicfont"

	"rp-go/engine/ecs"
	"rp-go/engine/events"
	"rp-go/engine/platform"
	"rp-go/engine/ui/window"
)

// ToolbarContent provides clickable buttons for debug toggles.
type ToolbarContent struct {
	buttons []ToolbarButton
}

// ToolbarButton defines a simple button rect + label + action.
type ToolbarButton struct {
	Label string
	OnClick func()
	Bounds window.Bounds
}

// Init initializes the toolbar buttons and binds them to events.
func (c *ToolbarContent) Init(bus *events.TypedBus) {
	c.buttons = []ToolbarButton{
		{
			Label: "Stats",
			OnClick: func() {
				events.Queue(bus, events.DebugToggleEvent{Enabled: true})
			},
		},
		{
			Label: "Entities",
			OnClick: func() {
				events.Queue(bus, events.WindowRestoredEvent{ID: "debug.entities"})
			},
		},
		{
			Label: "Systems",
			OnClick: func() {
				events.Queue(bus, events.WindowRestoredEvent{ID: "debug.systems"})
			},
		},
	}
}

// Draw renders the toolbar buttons.
func (c *ToolbarContent) Draw(_ *ecs.World, canvas *platform.Image, bounds window.Bounds) {
	if canvas == nil {
		return
	}

	const btnWidth, btnHeight = 60, 20
	x := bounds.X + 8
	y := bounds.Y + 8
	for i := range c.buttons {
		b := &c.buttons[i]
		b.Bounds = window.Bounds{X: x, Y: y, Width: btnWidth, Height: btnHeight}

		drawButton(canvas, b.Bounds, b.Label)
		x += btnWidth + 8
	}
}

// drawButton draws a minimal button with a label.
func drawButton(canvas *platform.Image, bounds window.Bounds, label string) {
	btnColor := color.RGBA{30, 45, 90, 220}
	borderColor := color.RGBA{180, 210, 255, 180}
	textColor := color.White

	// Background
	canvas.FillRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, btnColor)
	// Border
	canvas.FillRect(bounds.X, bounds.Y, bounds.Width, 1, borderColor)
	canvas.FillRect(bounds.X, bounds.Y+bounds.Height-1, bounds.Width, 1, borderColor)
	canvas.FillRect(bounds.X, bounds.Y, 1, bounds.Height, borderColor)
	canvas.FillRect(bounds.X+bounds.Width-1, bounds.Y, 1, bounds.Height, borderColor)

	// ✅ FIX: Use a guaranteed non-nil font
	font := basicfont.Face7x13
	textX := bounds.X + 8
	textY := bounds.Y + bounds.Height/2 + 5

	platform.DrawText(canvas, label, font, textX, textY, textColor)
}

