package world

import (
	"fmt"
	"sync"

	"rp-go/engine/data"
	"rp-go/engine/ecs"
	"rp-go/engine/gfx"
)

// ActorCreator spawns ECS entities based on JSON-defined templates.
type ActorCreator struct {
	templates map[string]data.ActorTemplate
	counters  map[string]int
	mu        sync.Mutex
}

// NewActorCreator constructs an ActorCreator from a loaded template database.
func NewActorCreator(db data.ActorDatabase) *ActorCreator {
	templates := make(map[string]data.ActorTemplate, len(db.Actors))
	for _, tpl := range db.Actors {
		templates[tpl.Name] = tpl
	}
	return &ActorCreator{
		templates: templates,
		counters:  make(map[string]int, len(templates)),
	}
}

// PreloadImages warms the graphics cache for every actor template sprite.
func (c *ActorCreator) PreloadImages() {
	if c == nil {
		return
	}
	paths := make([]string, 0, len(c.templates))
	for _, tpl := range c.templates {
		if tpl.Sprite.Image != "" {
			paths = append(paths, tpl.Sprite.Image)
		}
	}
	if len(paths) > 0 {
		gfx.PreloadImages(paths...)
	}
}

// Spawn instantiates a new actor entity using the specified template.
func (c *ActorCreator) Spawn(w *ecs.World, template string, position ecs.Position) (*ecs.Entity, error) {
	if c == nil {
		return nil, fmt.Errorf("actor creator is nil")
	}
	if w == nil {
		return nil, fmt.Errorf("world is nil")
	}

	tpl, ok := c.templates[template]
	if !ok {
		return nil, fmt.Errorf("actor template %q not found", template)
	}

	entity := w.NewEntity()

	id := c.nextID(template)
	entity.Add(&ecs.Actor{
		ID:         id,
		Archetype:  tpl.Archetype,
		Persistent: tpl.Persistent,
	})

	entity.Add(&ecs.Position{X: position.X, Y: position.Y})

	if tpl.Velocity != nil {
		entity.Add(&ecs.Velocity{VX: tpl.Velocity.VX, VY: tpl.Velocity.VY})
	}

	if tpl.Sprite.Image != "" {
		spriteImage := gfx.LoadImage(tpl.Sprite.Image)
		sprite := &ecs.Sprite{
			Image:          spriteImage,
			Width:          tpl.Sprite.Width,
			Height:         tpl.Sprite.Height,
			Rotation:       tpl.Sprite.Rotation,
			FlipHorizontal: tpl.Sprite.FlipHorizontal,
			PixelPerfect:   tpl.Sprite.PixelPerfect,
		}

		if sprite.Width == 0 || sprite.Height == 0 {
			if spriteImage != nil {
				bounds := spriteImage.Bounds()
				sprite.Width = bounds.Dx()
				sprite.Height = bounds.Dy()
			}
		}

		entity.Add(sprite)
	}

	return entity, nil
}

// nextID increments and returns a unique actor ID for the given template name.
func (c *ActorCreator) nextID(template string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counters[template]++
	return fmt.Sprintf("%s-%04d", template, c.counters[template])
}

// Templates returns the list of registered template names.
func (c *ActorCreator) Templates() []string {
	if c == nil {
		return nil
	}
	names := make([]string, 0, len(c.templates))
	for name := range c.templates {
		names = append(names, name)
	}
	return names
}
