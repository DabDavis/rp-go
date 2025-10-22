package gfx

import (
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func LoadImage(path string) *ebiten.Image {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("[GFX] Failed to open image: %s (%v)\n", path, err)
		return nil
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("[GFX] Failed to decode image: %s (%v)\n", path, err)
		return nil
	}

	return ebiten.NewImageFromImage(img)
}

