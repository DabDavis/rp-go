//go:build headless

package platform

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"sync"

	"golang.org/x/image/font"
)

type Image struct {
	mu     sync.Mutex
	rgba   *image.RGBA
	width  int
	height int
}

func newImageFromNative(img *Image) *Image {
	return img
}

func NativeImage(img *Image) *Image {
	return img
}

func NewImage(width, height int) *Image {
	if width <= 0 {
		width = 1
	}
	if height <= 0 {
		height = 1
	}
	return &Image{rgba: image.NewRGBA(image.Rect(0, 0, width, height)), width: width, height: height}
}

func NewImageFromImage(src image.Image) *Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	return &Image{rgba: dst, width: bounds.Dx(), height: bounds.Dy()}
}

func (img *Image) Clear() {
	if img == nil {
		return
	}
	img.Fill(color.RGBA{})
}

func (img *Image) Fill(c color.Color) {
	if img == nil {
		return
	}
	img.mu.Lock()
	defer img.mu.Unlock()
	draw.Draw(img.rgba, img.rgba.Bounds(), &image.Uniform{C: c}, image.Point{}, draw.Src)
}

func (img *Image) Bounds() image.Rectangle {
	if img == nil {
		return image.Rect(0, 0, 0, 0)
	}
	return image.Rect(0, 0, img.width, img.height)
}

func (img *Image) DrawImage(src *Image, _ *DrawImageOptions) {
	if img == nil || src == nil {
		return
	}
	img.mu.Lock()
	defer img.mu.Unlock()
	src.mu.Lock()
	defer src.mu.Unlock()
	draw.Draw(img.rgba, img.rgba.Bounds(), src.rgba, image.Point{}, draw.Over)
}

type DrawImageOptions struct{}

func NewDrawImageOptions() *DrawImageOptions { return &DrawImageOptions{} }

func (op *DrawImageOptions) SetFilter(Filter)       {}
func (op *DrawImageOptions) Scale(x, y float64)     {}
func (op *DrawImageOptions) Rotate(theta float64)   {}
func (op *DrawImageOptions) Translate(x, y float64) {}

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

func IsKeyPressed(Key) bool     { return false }
func IsKeyJustPressed(Key) bool { return false }
func Wheel() (float64, float64) { return 0, 0 }
func ActualFPS() float64        { return 60 }
func SetWindowSize(int, int)    {}
func SetWindowTitle(string)     {}

type GamepadID int

type GamepadLayoutID int

type StandardGamepadButton int

const (
	StandardGamepadButtonLeft StandardGamepadButton = iota
	StandardGamepadButtonRight
	StandardGamepadButtonUp
	StandardGamepadButtonDown
)

type StandardGamepadAxis int

const (
	StandardGamepadAxisLeftStickHorizontal StandardGamepadAxis = iota
	StandardGamepadAxisLeftStickVertical
)

func GamepadIDs() []GamepadID { return nil }

func StandardGamepadLayoutID(GamepadID) (GamepadLayoutID, bool) { return 0, false }

func StandardGamepadAxisValue(GamepadID, StandardGamepadAxis) float64 { return 0 }

func IsStandardGamepadButtonPressed(GamepadID, StandardGamepadButton) bool { return false }

type Game interface {
	Update() error
	Draw(screen *Image)
	Layout(outsideWidth, outsideHeight int) (int, int)
}

func RunGame(Game) error {
	return errors.New("headless build does not support interactive RunGame")
}

func RunHeadless(game Game, frames int, width, height int) error {
	if frames <= 0 {
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

func DrawText(dst *Image, _ string, _ font.Face, _ int, _ int, _ color.Color) {}
