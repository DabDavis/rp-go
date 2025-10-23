package windowmgr

import (
	"image/color"
	"sort"

	"golang.org/x/image/font/basicfont"

	"rp-go/engine/ecs"
	"rp-go/engine/platform"
	"rp-go/engine/ui/window"
)

// System draws reusable UI windows registered as ECS components.
type System struct {
	layer ecs.DrawLayer
}

// NewSystem creates a window manager for the HUD overlay layer.
func NewSystem() *System {
	return &System{layer: ecs.LayerHUD}
}

// Layer reports which draw layer the system renders in.
func (s *System) Layer() ecs.DrawLayer { return s.layer }

// Update is currently a no-op; windows are purely rendered components.
func (s *System) Update(*ecs.World) {}

// Draw renders all visible windows matching the manager's layer.
func (s *System) Draw(world *ecs.World, screen *platform.Image) {
	if world == nil || screen == nil {
		return
	}

	windows := s.collectWindows(world)
	for _, win := range windows {
		s.drawWindow(world, screen, win)
	}
}

func (s *System) collectWindows(world *ecs.World) []*window.Component {
	if world == nil {
		return nil
	}
	windows := make([]*window.Component, 0, len(world.Entities))
	for _, entity := range world.Entities {
		if entity == nil {
			continue
		}
		comp, _ := entity.Get("Window").(*window.Component)
		if comp == nil || !comp.Visible {
			continue
		}
		if comp.Layer != s.layer {
			continue
		}
		windows = append(windows, comp)
	}

	sort.SliceStable(windows, func(i, j int) bool {
		if windows[i].Order == windows[j].Order {
			return windows[i].ID < windows[j].ID
		}
		return windows[i].Order < windows[j].Order
	})

	return windows
}

func (s *System) drawWindow(world *ecs.World, screen *platform.Image, comp *window.Component) {
	if comp == nil {
		return
	}
	bounds := comp.Bounds
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	canvas := platform.NewImage(bounds.Width, bounds.Height)
	background := colorOrDefault(comp.Background, color.RGBA{8, 12, 20, 200})
	canvas.FillRect(0, 0, bounds.Width, bounds.Height, background)

	titleBarHeight := comp.TitleBarHeight
	if titleBarHeight < 0 {
		titleBarHeight = 0
	}
	if titleBarHeight > 0 {
		header := colorOrDefault(comp.TitleBar, color.RGBA{20, 36, 80, 220})
		canvas.FillRect(0, 0, bounds.Width, titleBarHeight, header)
	}

	borderColor := colorOrDefault(comp.Border, color.RGBA{180, 210, 255, 120})
	drawBorder(canvas, bounds.Width, bounds.Height, borderColor)

	if titleBarHeight > 0 && comp.Title != "" {
		textColor := colorOrDefault(comp.TitleColor, color.White)
		textX := comp.Padding
		if textX <= 0 {
			textX = 8
		}
		baseline := titleBarHeight/2 + 6
		if baseline < 12 {
			baseline = 12
		}
		platform.DrawText(canvas, comp.Title, basicfont.Face7x13, textX, baseline, textColor)
	}

	if comp.Content != nil {
		contentBounds := comp.ContentBounds()
		if contentBounds.Width > 0 && contentBounds.Height >= 0 {
			comp.Content.Draw(world, canvas, contentBounds)
		}
	}

	op := platform.NewDrawImageOptions()
	op.Translate(float64(bounds.X), float64(bounds.Y))
	screen.DrawImage(canvas, op)
}

func colorOrDefault(c color.Color, fallback color.Color) color.Color {
	if c == nil {
		return fallback
	}
	return c
}

func drawBorder(img *platform.Image, width, height int, border color.Color) {
	if img == nil || width <= 0 || height <= 0 {
		return
	}
	// Draw a simple 1px border around the window.
	img.FillRect(0, 0, width, 1, border)
	img.FillRect(0, height-1, width, 1, border)
	img.FillRect(0, 0, 1, height, border)
	img.FillRect(width-1, 0, 1, height, border)
}
