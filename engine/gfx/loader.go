package gfx

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"sync"

	"rp-go/engine/platform"
)

type cachedImage struct {
	once sync.Once
	img  *platform.Image
	err  error
}

var imageCache sync.Map // map[string]*cachedImage

// LoadImage returns an image, caching the decoded result so repeated calls reuse
// the same underlying resource.
func LoadImage(path string) *platform.Image {
	entryAny, _ := imageCache.LoadOrStore(path, &cachedImage{})
	entry := entryAny.(*cachedImage)

	entry.once.Do(func() {
		entry.img, entry.err = decodeImage(path)
		if entry.err != nil {
			fmt.Printf("[GFX] failed to load image %s: %v\n", path, entry.err)
		}
	})

	if entry.err != nil {
		return nil
	}
	return entry.img
}

// PreloadImages eagerly loads the provided paths so the cache is primed before
// gameplay needs them. The work is performed serially to keep the loader easy to
// reason about in all build modes.
func PreloadImages(paths ...string) {
	for _, path := range paths {
		if LoadImage(path) == nil {
			fmt.Printf("[GFX] preload failed for %s\n", path)
		}
	}
}

func decodeImage(path string) (*platform.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return platform.NewImageFromImage(img), nil
}
