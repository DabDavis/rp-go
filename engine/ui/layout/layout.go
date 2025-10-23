package layout

import "rp-go/engine/ui/button"

// Horizontal represents a row layout with consistent padding.
type Horizontal struct {
	PaddingX int
	PaddingY int
	Buttons  []*button.Button
}

// NewHorizontal constructs a horizontal layout container.
func NewHorizontal(paddingX, paddingY int) *Horizontal {
	return &Horizontal{
		PaddingX: paddingX,
		PaddingY: paddingY,
		Buttons:  []*button.Button{},
	}
}

// Add inserts a button to the layout.
func (h *Horizontal) Add(b *button.Button) {
	h.Buttons = append(h.Buttons, b)
}

// Elements returns all buttons.
func (h *Horizontal) Elements() []*button.Button {
	return h.Buttons
}

