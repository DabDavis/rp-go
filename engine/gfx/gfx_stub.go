//go:build test
// +build test

package gfx

import (
	"image"
	"log"
)

// Image is a lightweight placeholder for tests.
type Image struct{}

// LoadImage fakes a texture load for unit tests (no filesystem access).
func LoadImage(path string) *Image {
	log.Printf("[gfx_stub] Pretending to load image: %s", path)
	return &Image{}
}

// PreloadImages fakes sprite preloading for tests.
func PreloadImages(paths ...string) {
	for _, p := range paths {
		log.Printf("[gfx_stub] Pretending to preload: %s", p)
	}
}

// Bounds returns a fixed-size dummy rectangle.
func (i *Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, 64, 64)
}

