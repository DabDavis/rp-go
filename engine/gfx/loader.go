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
	img  *platform.Image
	img  *platform.Image
	err  error
}

var imageCache sync.Map // map[string]*cachedImage

// LoadImage returns an Ebiten image, caching decoded results so repeated
// requests (even across goroutines) reuse the same GPU resource.
func LoadImage(path string) *platform.Image {

func LoadImage(path string) *platform.Image {
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

// PreloadImages eagerly loads a list of image paths using a worker-per-path
// fan-out. It reuses the LoadImage cache so subsequent calls are instantaneous.
func PreloadImages(paths ...string) {
	var wg sync.WaitGroup
	wg.Add(len(paths))
	for _, path := range paths {
		path := path
		go func() {
			defer wg.Done()
			LoadImage(path)
		}()
	}
	wg.Wait()
}

func decodeImage(path string) (*platform.Image, error) {
func decodeImage(path string) (*platform.Image, error) {
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
	return platform.NewImageFromImage(img), nil
	return platform.NewImageFromImage(img), nil
}
