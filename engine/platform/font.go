package platform

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

// DefaultFont returns a globally available fallback font
// for all UI and debug text rendering.
func DefaultFont() font.Face {
	return basicfont.Face7x13
}

