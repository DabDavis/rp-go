package world

import (
	"fmt"
	"rp-go/engine/data"
	"rp-go/engine/ecs"
	"rp-go/engine/gfx"
	"sync"
)

/*───────────────────────────────────────────────*
 | ACTOR CREATOR                                 |
 *───────────────────────────────────────────────*/

// ActorCreator spawns ECS entities from JSON-defined templates (actors.json).
// It automatically assigns unique IDs, sets up base components,
// and leaves AI binding to the AIComposer system.
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

/*───────────────────────────────────────────────*
 | PRELOAD / CACHING                             |
 *───────────────────────────────────────────────*/

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

/*───────────────────────────────────────────────*
 | ENTITY SPAWNING                               |
 *───────────────────────────────────────────────*/

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

	// --- Actor metadata
	actor := &ecs.Actor{
		ID:         c.nextID(template),
		Archetype:  tpl.Archetype,
		Persistent: tpl.Persistent,
	}
	e.Add(actor)

	// --- Transform
	e.Add(&ecs.Position{X: pos.X, Y: pos.Y})
	if tpl.Velocity != nil {
		e.Add(&ecs.Velocity{VX: tpl.Velocity.VX, VY: tpl.Velocity.VY})
	}

	// --- Sprite
	if tpl.Sprite.Image != "" {
		e.Add(buildSprite(tpl.Sprite))
	}

	// --- AI references (deferred to AIComposer)
	if len(tpl.AIRefs) > 0 {
		// Store references directly in the Actor for AIComposer to process.
		actor.AIRefs = append([]string{}, tpl.AIRefs...)
	}

	return e, nil
}

/*───────────────────────────────────────────────*
 | HELPERS                                       |
 *───────────────────────────────────────────────*/

// nextID generates a unique per-template instance ID.
func (c *ActorCreator) nextID(template string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counters[template]++
	return fmt.Sprintf("%s-%03d", template, c.counters[template])
}

// buildSprite converts a JSON sprite definition into an ECS Sprite component.
func buildSprite(st data.ActorSpriteTemplate) *ecs.Sprite {
	img := gfx.LoadImage(st.Image)
	return &ecs.Sprite{
		Image:          img,
		Width:          st.Width,
		Height:         st.Height,
		Rotation:       st.Rotation,
		FlipHorizontal: st.FlipHorizontal,
		PixelPerfect:   st.PixelPerfect,
	}
}

