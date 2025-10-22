//go:build !headless

package platform_desktop

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

func DrawText(dst *Image, str string, face font.Face, x, y int, clr color.Color) {
	if dst == nil {
		return
	}
	text.Draw(dst.native, str, face, x, y, clr)
}
