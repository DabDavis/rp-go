package world

import (
	"fmt"
	"rp-go/engine/data"
	"rp-go/engine/ecs"
	"rp-go/engine/gfx"
)

// nextID returns a unique, thread-safe actor ID per template.
func (c *ActorCreator) nextID(template string) string {
	if c == nil {
		return fmt.Sprintf("%s-0000", template)
	}
	if c.counters == nil {
		c.counters = make(map[string]int)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.counters[template]++
	return fmt.Sprintf("%s-%04d", template, c.counters[template])
}

// Templates returns all known actor template names.
func (c *ActorCreator) Templates() []string {
	if c == nil {
		return nil
	}
	names := make([]string, 0, len(c.templates))
	for k := range c.templates {
		names = append(names, k)
	}
	return names
}

// buildSprite safely constructs an ecs.Sprite from data.ActorSpriteTemplate.
func buildSprite(st data.ActorSpriteTemplate) *ecs.Sprite {
	img := gfx.LoadImage(st.Image)
	sprite := &ecs.Sprite{
		Image:          img,
		Width:          st.Width,
		Height:         st.Height,
		Rotation:       st.Rotation,
		FlipHorizontal: st.FlipHorizontal,
		PixelPerfect:   st.PixelPerfect,
	}
	if img != nil && (sprite.Width == 0 || sprite.Height == 0) {
		b := img.Bounds()
		sprite.Width, sprite.Height = b.Dx(), b.Dy()
	}
	return sprite
}

