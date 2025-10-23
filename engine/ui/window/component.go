package window

import (
	"image/color"

	"rp-go/engine/ecs"
	"rp-go/engine/platform"
)

/* -------------------------------------------------------------------------- */
/*                                   Bounds                                   */
/* -------------------------------------------------------------------------- */

// Bounds describes a rectangular region in screen space.
type Bounds struct {
	X      int
	Y      int
	Width  int
	Height int
}

// Inset shrinks the bounds by the given padding on all sides.
func (b Bounds) Inset(padding int) Bounds {
	if padding <= 0 {
		return b
	}
	nx := b.X + padding
	ny := b.Y + padding
	nw := b.Width - padding*2
	nh := b.Height - padding*2
	if nw < 0 {
		nw = 0
	}
	if nh < 0 {
		nh = 0
	}
	return Bounds{X: nx, Y: ny, Width: nw, Height: nh}
}

// Contains checks if a point (x, y) lies inside the bounds.
func (b Bounds) Contains(x, y int) bool {
	return x >= b.X && x < b.X+b.Width && y >= b.Y && y < b.Y+b.Height
}

/* -------------------------------------------------------------------------- */
/*                                   Content                                  */
/* -------------------------------------------------------------------------- */

// Content represents drawable window content.
type Content interface {
	Draw(world *ecs.World, canvas *platform.Image, bounds Bounds)
}

// RendererFunc adapts a plain function to Content.
type RendererFunc func(world *ecs.World, canvas *platform.Image, bounds Bounds)

// Draw renders the function based content.
func (fn RendererFunc) Draw(world *ecs.World, canvas *platform.Image, bounds Bounds) {
	if fn == nil {
		return
	}
	fn(world, canvas, bounds)
}

/* -------------------------------------------------------------------------- */
/*                                  Component                                 */
/* -------------------------------------------------------------------------- */

// Component is an ECS component describing a reusable UI window.
type Component struct {
	ID      string
	Title   string
	Bounds  Bounds
	Visible bool
	Order   int
	Layer   ecs.DrawLayer

	Padding        int
	TitleBarHeight int

	Background color.Color
	Border     color.Color
	TitleBar   color.Color
	TitleColor color.Color

	Content Content

	// --- Interactive features ---
	Movable   bool // allows dragging via title bar
	Locked    bool // disables user input entirely
	Minimized bool // window is collapsed (content hidden)
	Closable  bool // shows and enables a close button
	Closed    bool // true if user clicked close
}

/* -------------------------------------------------------------------------- */
/*                               Constructors                                 */
/* -------------------------------------------------------------------------- */

// NewComponent creates a window component with sensible defaults.
func NewComponent(id, title string, bounds Bounds, content Content) *Component {
	return &Component{
		ID:             id,
		Title:          title,
		Bounds:         bounds,
		Visible:        true,
		Order:          0,
		Layer:          ecs.LayerHUD,
		Padding:        8,
		TitleBarHeight: 24,
		Background:     color.RGBA{0, 0, 0, 180},
		Border:         color.RGBA{255, 255, 255, 60},
		TitleBar:       color.RGBA{30, 45, 90, 220},
		TitleColor:     color.RGBA{220, 235, 255, 255},
		Content:        content,
		Movable:        true,  // default: draggable
		Closable:       true,  // default: shows a close button
		Locked:         false, // default: interactive
		Minimized:      false,
		Closed:         false,
	}
}

// Name implements ecs.Component.
func (c *Component) Name() string { return "Window" }

/* -------------------------------------------------------------------------- */
/*                               Helper Methods                               */
/* -------------------------------------------------------------------------- */

// ContentBounds returns the drawable region for the window content.
func (c *Component) ContentBounds() Bounds {
	if c == nil {
		return Bounds{}
	}
	if c.Minimized {
		// When minimized, content area is 0 height.
		return Bounds{X: 0, Y: 0, Width: 0, Height: 0}
	}

	padding := max(0, c.Padding)
	titleBar := max(0, c.TitleBarHeight)

	width := c.Bounds.Width - padding*2
	height := c.Bounds.Height - titleBar - padding*2
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	return Bounds{
		X:      padding,
		Y:      titleBar + padding,
		Width:  width,
		Height: height,
	}
}

// ToggleMinimize flips the minimized state.
func (c *Component) ToggleMinimize() { c.Minimized = !c.Minimized }

// Close marks the window as closed and invisible.
func (c *Component) Close() {
	c.Closed = true
	c.Visible = false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

