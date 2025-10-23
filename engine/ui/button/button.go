package button

import (
	"image/color"

	"rp-go/engine/platform"
)

// Button represents a clickable text button.
type Button struct {
	Label    string
	OnClick  func()
	Width    int
	Height   int
	Hovered  bool
	Active   bool
	Disabled bool
}

// New creates a new button with label and callback.
func New(label string, onClick func()) *Button {
	return &Button{
		Label:   label,
		OnClick: onClick,
		Width:   80,
		Height:  22,
	}
}

// Draw renders the button background and label.
func (b *Button) Draw(screen *platform.Image, x, y int) {
	bg := color.RGBA{50, 60, 80, 200}
	if b.Hovered {
		bg = color.RGBA{70, 90, 120, 220}
	}
	if b.Active {
		bg = color.RGBA{100, 120, 180, 255}
	}
	if b.Disabled {
		bg = color.RGBA{80, 80, 80, 180}
	}

	screen.FillRect(x, y, b.Width, b.Height, bg)
	platform.DrawText(screen, b.Label, platform.DefaultFont(), x+6, y+16, color.White)
}

// Click invokes the OnClick handler if enabled.
func (b *Button) Click() {
	if !b.Disabled && b.OnClick != nil {
		b.OnClick()
	}
}

