package render

import (
	"image/color"
	"golang.org/x/image/font/basicfont"

	"rp-go/engine/ecs"
	"rp-go/engine/platform"
	"rp-go/engine/systems/windowmgr"
	"rp-go/engine/ui/window"
)

// WindowRenderer draws UI windows registered in the global window registry.
type WindowRenderer struct {
	layer ecs.DrawLayer
}

func NewWindowRenderer(layer ecs.DrawLayer) *WindowRenderer {
	return &WindowRenderer{layer: layer}
}

func (r *WindowRenderer) Layer() ecs.DrawLayer { return r.layer }
func (r *WindowRenderer) Update(*ecs.World)    {}

// Draw renders all visible windows for the configured layer.
func (r *WindowRenderer) Draw(world *ecs.World, screen *platform.Image) {
	if world == nil || screen == nil {
		return
	}

	registry := windowmgr.SharedRegistry()
	for _, comp := range registry.Windows(r.layer) {
		drawWindow(world, screen, comp)
	}
}

func drawWindow(world *ecs.World, screen *platform.Image, comp *window.Component) {
	if comp == nil || !comp.Visible {
		return
	}

	b := comp.Bounds
	if b.Width <= 0 || b.Height <= 0 {
		return
	}

	canvas := platform.NewImage(b.Width, b.Height)
	background := colorOrDefault(comp.Background, color.RGBA{8, 12, 20, 200})
	canvas.FillRect(0, 0, b.Width, b.Height, background)

	titleBarHeight := comp.TitleBarHeight
	if titleBarHeight > 0 {
		header := colorOrDefault(comp.TitleBar, color.RGBA{20, 36, 80, 220})
		canvas.FillRect(0, 0, b.Width, titleBarHeight, header)
	}

	borderColor := colorOrDefault(comp.Border, color.RGBA{180, 210, 255, 120})
	drawBorder(canvas, b.Width, b.Height, borderColor)

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

	if !comp.Minimized && comp.Content != nil {
		contentBounds := comp.ContentBounds()
		if contentBounds.Width > 0 && contentBounds.Height >= 0 {
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
	if img == nil {
		return
	}
	img.FillRect(0, 0, width, 1, border)
	img.FillRect(0, height-1, width, 1, border)
	img.FillRect(0, 0, 1, height, border)
	img.FillRect(width-1, 0, 1, height, border)
}

