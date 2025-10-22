//go:build headless

package platform

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"math"

	"golang.org/x/image/font"
)

type Image struct {
	rgba *image.RGBA
}

type DrawImageOptions struct {
	translateX float64
	translateY float64
	scaleX     float64
	scaleY     float64
	rotation   float64
	filter     Filter
}

type Filter int

const (
	FilterNearest Filter = iota
)

type Key int

const (
	KeyArrowLeft Key = iota
	KeyArrowRight
	KeyArrowUp
	KeyArrowDown
	KeyA
	KeyD
	KeyW
	KeyS
	KeyQ
	KeyE
	KeyMinus
	KeyEqual
	Key0
	KeyKP0
	KeyKPAdd
	KeyKPSubtract
)

type GamepadID int

type StandardGamepadButton int

type StandardGamepadAxis int

type Game interface {
	Update() error
	Draw(screen *Image)
	Layout(outsideWidth, outsideHeight int) (int, int)
}

func NewImage(width, height int) *Image {
	if width <= 0 {
		width = 1
	}
	if height <= 0 {
		height = 1
	}
	return &Image{rgba: image.NewRGBA(image.Rect(0, 0, width, height))}
}

func NewImageFromImage(src image.Image) *Image {
	if src == nil {
		return nil
	}
	bounds := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(dst, dst.Bounds(), src, bounds.Min, draw.Src)
	return &Image{rgba: dst}
}

func NewDrawImageOptions() *DrawImageOptions {
	return &DrawImageOptions{scaleX: 1, scaleY: 1}
}

func (op *DrawImageOptions) SetFilter(f Filter) { op.filter = f }

func (op *DrawImageOptions) Scale(x, y float64) {
	if op == nil {
		return
	}
	op.scaleX *= x
	op.scaleY *= y
}

func (op *DrawImageOptions) Rotate(theta float64) {
	if op == nil {
		return
	}
	op.rotation += theta
}

func (op *DrawImageOptions) Translate(x, y float64) {
	if op == nil {
		return
	}
	op.translateX += x
	op.translateY += y
}

func (img *Image) Clear() {
	if img == nil || img.rgba == nil {
		return
	}
	for i := range img.rgba.Pix {
		img.rgba.Pix[i] = 0
	}
}

func (img *Image) Fill(c color.Color) {
	if img == nil || img.rgba == nil {
		return
	}
	draw.Draw(img.rgba, img.rgba.Bounds(), &image.Uniform{C: c}, image.Point{}, draw.Src)
}

func (img *Image) Bounds() image.Rectangle {
	if img == nil || img.rgba == nil {
		return image.Rect(0, 0, 0, 0)
	}
	return img.rgba.Bounds()
}

func (img *Image) DrawImage(src *Image, op *DrawImageOptions) {
	if img == nil || img.rgba == nil || src == nil || src.rgba == nil {
		return
	}
	dx := 0
	dy := 0
	if op != nil {
		dx = int(math.Round(op.translateX))
		dy = int(math.Round(op.translateY))
	}
	dstRect := src.rgba.Bounds().Add(image.Pt(dx, dy))
	draw.Draw(img.rgba, dstRect, src.rgba, src.rgba.Bounds().Min, draw.Over)
}

func (img *Image) FillRect(x, y, w, h int, c color.Color) {
	if img == nil || img.rgba == nil || w <= 0 || h <= 0 {
		return
	}
	rect := image.Rect(x, y, x+w, y+h)
	draw.Draw(img.rgba, rect, &image.Uniform{C: c}, image.Point{}, draw.Src)
}

func IsKeyPressed(Key) bool     { return false }
func IsKeyJustPressed(Key) bool { return false }

func GamepadIDs() []GamepadID { return nil }

func IsStandardGamepadLayoutAvailable(GamepadID) bool { return false }

const (
	StandardGamepadAxisLeftStickHorizontal StandardGamepadAxis = iota
	StandardGamepadAxisLeftStickVertical
)

func StandardGamepadAxisValue(GamepadID, StandardGamepadAxis) float64 { return 0 }

func IsStandardGamepadButtonPressed(GamepadID, StandardGamepadButton) bool { return false }
func IsGamepadLeft(GamepadID) bool                                         { return false }
func IsGamepadRight(GamepadID) bool                                        { return false }
func IsGamepadUp(GamepadID) bool                                           { return false }
func IsGamepadDown(GamepadID) bool                                         { return false }

func Wheel() (float64, float64) { return 0, 0 }

func ActualFPS() float64 { return 60 }

func SetWindowSize(int, int) {}
func SetWindowTitle(string)  {}

func RunGame(Game) error { return errors.New("headless build does not support interactive RunGame") }

func RunHeadless(game Game, frames, width, height int) error {
	if game == nil || frames <= 0 {
		return nil
	}
	screen := NewImage(width, height)
	for i := 0; i < frames; i++ {
		if err := game.Update(); err != nil {
			return err
		}
		screen.Clear()
		game.Draw(screen)
	}
	return nil
}

func DrawText(*Image, string, font.Face, int, int, color.Color) {}
