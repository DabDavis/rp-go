//go:build !headless

package platform_desktop

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
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

func (img *Image) Clear() { img.native.Clear() }
func (img *Image) Fill(c color.Color) { img.native.Fill(c) }
func (img *Image) Bounds() image.Rectangle { return img.native.Bounds() }

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

// Helper rectangle fill (for starfield, etc.)
func (img *Image) FillRect(x, y, w, h int, c color.Color) {
	if img == nil {
		return
	}
	rect := image.Rect(x, y, x+w, y+h)
	sub := img.native.SubImage(rect).(*ebiten.Image)
	sub.Fill(c)
}

