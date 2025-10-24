package world

import (
	"fmt"
	"rp-go/engine/gfx"
)

// PreloadAllImages iterates over all registered creators to preload textures.
// You can call this once during startup to warm GPU caches.
func PreloadAllImages(creators ...*ActorCreator) {
	total := 0
	for _, c := range creators {
		if c == nil {
			continue
		}
		paths := make([]string, 0)
		for _, tpl := range c.templates {
			if tpl.Sprite.Image != "" {
				paths = append(paths, tpl.Sprite.Image)
			}
		}
		if len(paths) > 0 {
			gfx.PreloadImages(paths...)
			total += len(paths)
		}
	}
	fmt.Printf("[WORLD] Preloaded %d textures for all creators\n", total)
}

