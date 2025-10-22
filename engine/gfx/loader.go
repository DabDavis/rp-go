package gfx

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"sync"

	"rp-go/engine/platform"
)

// cachedImage wraps a lazily-loaded image resource with one-time initialization.
type cachedImage struct {
	once sync.Once
	img  *platform.Image
	err  error
}

// imageCache maps image file paths to their cached image objects.
var imageCache sync.Map // map[string]*cachedImage

// LoadImage returns an Ebiten-compatible image, caching the decoded result.
// Repeated calls with the same path reuse the same GPU resource.
func LoadImage(path string) *platform.Image {
	entryAny, _ := imageCache.LoadOrStore(path, &cachedImage{})
	entry := entryAny.(*cachedImage)

	entry.once.Do(func() {
		entry.img, entry.err = decodeImage(path)
		if entry.err != nil {
			fmt.Printf("[GFX] Failed to load image: %s (%v)\n", path, entry.err)
		}
	})

	if entry.err != nil {
		return nil
	}
	return entry.img
}

// PreloadImages loads a list of images concurrently, populating the cache
// ahead of time so gameplay can access them instantly later.
func PreloadImages(paths ...string) {
	var wg sync.WaitGroup
	wg.Add(len(paths))
	for _, path := range paths {
		path := path
		go func() {
			defer wg.Done()
			if LoadImage(path) == nil {
				fmt.Printf("[GFX] Preload failed for %s\n", path)
			}
		}()
	}
	wg.Wait()
}

// decodeImage decodes a PNG (or other supported formats) from disk and wraps
// it in a platform.Image for rendering.
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

