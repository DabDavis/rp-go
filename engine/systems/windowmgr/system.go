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
	layers []ecs.DrawLayer
}

// NewSystem creates a window manager that draws HUD and Console overlays.
func NewSystem() *System {
	return &System{
		layers: []ecs.DrawLayer{
			ecs.LayerHUD,
			ecs.LayerConsole,
		},
	}
}

// Layer reports the nominal ECS layer this system participates in.
func (s *System) Layer() ecs.DrawLayer { return ecs.LayerHUD }

// Update is currently a no-op; windows are purely rendered components.
func (s *System) Update(*ecs.World) {}

// Draw renders all visible windows for the handled layers.
func (s *System) Draw(world *ecs.World, screen *platform.Image) {
	if world == nil || screen == nil {
		return
	}

	for _, win := range s.collectWindows(world) {
		s.drawWindow(world, screen, win)
	}
}

func (s *System) collectWindows(world *ecs.World) []*window.Component {
	if world == nil {
		return nil
	}

	var windows []*window.Component
	for _, e := range world.Entities {
		if e == nil {
			continue
		}
		comp, _ := e.Get("Window").(*window.Component)
		if comp == nil || !comp.Visible || !s.handlesLayer(comp.Layer) {
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

func (s *System) handlesLayer(layer ecs.DrawLayer) bool {
	for _, l := range s.layers {
		if layer == l {
			return true
		}
	}
	return false
}

func (s *System) drawWindow(world *ecs.World, screen *platform.Image, comp *window.Component) {
	if comp == nil {
		return
	}
	b := comp.Bounds
	if b.Width <= 0 || b.Height <= 0 {
		return
	}

	canvas := platform.NewImage(b.Width, b.Height)
	bg := colorOrDefault(comp.Background, color.RGBA{8, 12, 20, 200})
	canvas.FillRect(0, 0, b.Width, b.Height, bg)

	titleBarHeight := max(0, comp.TitleBarHeight)
	if titleBarHeight > 0 {
		header := colorOrDefault(comp.TitleBar, color.RGBA{20, 36, 80, 220})
		canvas.FillRect(0, 0, b.Width, titleBarHeight, header)
	}

	border := colorOrDefault(comp.Border, color.RGBA{180, 210, 255, 120})
	drawBorder(canvas, b.Width, b.Height, border)

	if titleBarHeight > 0 && comp.Title != "" {
		textColor := colorOrDefault(comp.TitleColor, color.White)
		textX := max(8, comp.Padding)
		baseline := max(12, titleBarHeight/2+6)
		platform.DrawText(canvas, comp.Title, basicfont.Face7x13, textX, baseline, textColor)
	}

	if comp.Content != nil {
		contentBounds := comp.ContentBounds()
		if contentBounds.Width > 0 && contentBounds.Height > 0 {
			comp.Content.Draw(world, canvas, contentBounds)
		}
	}

	op := platform.NewDrawImageOptions()
	op.Translate(float64(b.X), float64(b.Y))
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
	img.FillRect(0, 0, width, 1, border)
	img.FillRect(0, height-1, width, 1, border)
	img.FillRect(0, 0, 1, height, border)
	img.FillRect(width-1, 0, 1, height, border)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

