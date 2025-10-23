package world

import (
	"fmt"
	"rp-go/engine/data"
	"rp-go/engine/ecs"
	"rp-go/engine/gfx"
	"sync"
)

// ActorCreator spawns ECS entities from JSON-defined templates.
type ActorCreator struct {
	templates map[string]data.ActorTemplate
	counters  map[string]int
	mu        sync.Mutex
}

// NewActorCreator constructs a creator from a loaded ActorDatabase.
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

// PreloadImages warms the graphics cache for all known sprite paths.
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

// Spawn instantiates an actor entity by template name and position.
func (c *ActorCreator) Spawn(w *ecs.World, template string, pos ecs.Position) (*ecs.Entity, error) {
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

	e := w.NewEntity()
	e.Add(&ecs.Actor{
		ID:         c.nextID(template),
		Archetype:  tpl.Archetype,
		Persistent: tpl.Persistent,
	})
	e.Add(&ecs.Position{X: pos.X, Y: pos.Y})

	if tpl.Velocity != nil {
		e.Add(&ecs.Velocity{VX: tpl.Velocity.VX, VY: tpl.Velocity.VY})
	}

	if tpl.Sprite.Image != "" {
		sprite := buildSprite(tpl.Sprite)
		e.Add(sprite)
	}

	if ai := buildAIController(tpl.AI); ai != nil {
		e.Add(ai)
	}

	return e, nil
}

