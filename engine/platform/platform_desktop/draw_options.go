//go:build !headless

package platform_desktop

import "github.com/hajimehoshi/ebiten/v2"

type DrawImageOptions struct {
	native *ebiten.DrawImageOptions
}

func NewDrawImageOptions() *DrawImageOptions {
	return &DrawImageOptions{native: &ebiten.DrawImageOptions{}}
}

type Filter int

const (
	FilterNearest Filter = iota
)

func (op *DrawImageOptions) SetFilter(f Filter) {
	if op == nil {
		return
	}
	if f == FilterNearest {
		op.native.Filter = ebiten.FilterNearest
	}
}

func (op *DrawImageOptions) Scale(x, y float64)     { op.native.GeoM.Scale(x, y) }
func (op *DrawImageOptions) Rotate(theta float64)   { op.native.GeoM.Rotate(theta) }
func (op *DrawImageOptions) Translate(x, y float64) { op.native.GeoM.Translate(x, y) }
