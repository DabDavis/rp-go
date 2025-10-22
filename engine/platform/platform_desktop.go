//go:build !headless

package platform

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Image struct {
	native *ebiten.Image
}

func newImageFromNative(img *ebiten.Image) *Image {
	if img == nil {
		return nil
	}
	return &Image{native: img}
}

func NativeImage(img *Image) *ebiten.Image {
	if img == nil {
		return nil
	}
	return img.native
}

func NewImage(width, height int) *Image {
	return &Image{native: ebiten.NewImage(width, height)}
}

func NewImageFromImage(src image.Image) *Image {
	return &Image{native: ebiten.NewImageFromImage(src)}
}

func (img *Image) Clear() {
	img.native.Clear()
}

func (img *Image) Fill(c color.Color) {
	img.native.Fill(c)
}

func (img *Image) Bounds() image.Rectangle {
	return img.native.Bounds()
}

func (img *Image) DrawImage(src *Image, op *DrawImageOptions) {
	if img == nil || src == nil {
		return
	}
	var nativeOp *ebiten.DrawImageOptions
	if op != nil {
		nativeOp = op.native
	}
	img.native.DrawImage(src.native, nativeOp)
}

type DrawImageOptions struct {
	native *ebiten.DrawImageOptions
}

func NewDrawImageOptions() *DrawImageOptions {
	return &DrawImageOptions{native: &ebiten.DrawImageOptions{}}
}

func (op *DrawImageOptions) SetFilter(f Filter) {
	if op == nil {
		return
	}
	switch f {
	case FilterNearest:
		op.native.Filter = ebiten.FilterNearest
	}
}

func (op *DrawImageOptions) Scale(x, y float64) {
	if op == nil {
		return
	}
	op.native.GeoM.Scale(x, y)
}

func (op *DrawImageOptions) Rotate(theta float64) {
	if op == nil {
		return
	}
	op.native.GeoM.Rotate(theta)
}

func (op *DrawImageOptions) Translate(x, y float64) {
	if op == nil {
		return
	}
	op.native.GeoM.Translate(x, y)
}

type Filter int

const (
	FilterNearest Filter = iota
)

type Key = ebiten.Key

const (
	KeyArrowLeft  Key = ebiten.KeyArrowLeft
	KeyArrowRight Key = ebiten.KeyArrowRight
	KeyArrowUp    Key = ebiten.KeyArrowUp
	KeyArrowDown  Key = ebiten.KeyArrowDown
	KeyA          Key = ebiten.KeyA
	KeyD          Key = ebiten.KeyD
	KeyW          Key = ebiten.KeyW
	KeyS          Key = ebiten.KeyS
	KeyQ          Key = ebiten.KeyQ
	KeyE          Key = ebiten.KeyE
	KeyMinus      Key = ebiten.KeyMinus
	KeyEqual      Key = ebiten.KeyEqual
	Key0          Key = ebiten.Key0
	KeyKP0        Key = ebiten.KeyKP0
	KeyKPAdd      Key = ebiten.KeyKPAdd
	KeyKPSubtract Key = ebiten.KeyKPSubtract
)

func IsKeyPressed(k Key) bool {
	return ebiten.IsKeyPressed(k)
}

func IsKeyJustPressed(k Key) bool {
	return inpututil.IsKeyJustPressed(k)
}

type GamepadID = ebiten.GamepadID

type GamepadLayoutID = ebiten.StandardGamepadLayoutID

type StandardGamepadButton = ebiten.StandardGamepadButton

const (
	StandardGamepadButtonLeft  StandardGamepadButton = ebiten.StandardGamepadButtonLeft
	StandardGamepadButtonRight StandardGamepadButton = ebiten.StandardGamepadButtonRight
	StandardGamepadButtonUp    StandardGamepadButton = ebiten.StandardGamepadButtonUp
	StandardGamepadButtonDown  StandardGamepadButton = ebiten.StandardGamepadButtonDown
)

type StandardGamepadAxis = ebiten.StandardGamepadAxis

const (
	StandardGamepadAxisLeftStickHorizontal StandardGamepadAxis = ebiten.StandardGamepadAxisLeftStickHorizontal
	StandardGamepadAxisLeftStickVertical   StandardGamepadAxis = ebiten.StandardGamepadAxisLeftStickVertical
)

func GamepadIDs() []GamepadID {
	return ebiten.GamepadIDs()
}

func StandardGamepadLayoutID(id GamepadID) (GamepadLayoutID, bool) {
	return ebiten.StandardGamepadLayoutID(id)
}

func StandardGamepadAxisValue(id GamepadID, axis StandardGamepadAxis) float64 {
	return ebiten.StandardGamepadAxisValue(id, axis)
}

func IsStandardGamepadButtonPressed(id GamepadID, button StandardGamepadButton) bool {
	return ebiten.IsStandardGamepadButtonPressed(id, button)
}

func Wheel() (float64, float64) {
	return ebiten.Wheel()
}

func ActualFPS() float64 {
	return ebiten.ActualFPS()
}

func SetWindowSize(w, h int) {
	ebiten.SetWindowSize(w, h)
}

func SetWindowTitle(title string) {
	ebiten.SetWindowTitle(title)
}

type Game interface {
	Update() error
	Draw(screen *Image)
	Layout(outsideWidth, outsideHeight int) (int, int)
}

type gameAdapter struct {
	game Game
}

func (g *gameAdapter) Update() error {
	return g.game.Update()
}

func (g *gameAdapter) Draw(screen *ebiten.Image) {
	g.game.Draw(newImageFromNative(screen))
}

func (g *gameAdapter) Layout(outW, outH int) (int, int) {
	return g.game.Layout(outW, outH)
}

func RunGame(game Game) error {
	return ebiten.RunGame(&gameAdapter{game: game})
}

func RunHeadless(game Game, frames int, width, height int) error {
	if frames <= 0 {
		return nil
	}
	offscreen := NewImage(width, height)
	for i := 0; i < frames; i++ {
		if err := game.Update(); err != nil {
			return err
		}
		offscreen.Clear()
		game.Draw(offscreen)
	}
	return nil
}

func DrawText(dst *Image, str string, face font.Face, x, y int, clr color.Color) {
	if dst == nil {
		return
	}
	text.Draw(dst.native, str, face, x, y, clr)
}
