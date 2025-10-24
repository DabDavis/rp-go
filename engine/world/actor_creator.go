package world

import (
	"fmt"
	"sync"

	"rp-go/engine/data"
	"rp-go/engine/ecs"
	"rp-go/engine/gfx"
)

/*───────────────────────────────────────────────*
 | ACTOR CREATOR                                 |
 *───────────────────────────────────────────────*/

// ActorCreator spawns ECS entities from JSON-defined templates (actors.json).
// It assigns unique IDs and attaches core components (Position, Sprite, etc.)
// but does NOT attach AI — that's handled later by the aicomposer system.
type ActorCreator struct {
	templates map[string]data.ActorTemplate
	counters  map[string]int
	mu        sync.Mutex
}

// NewActorCreator constructs a new creator from a loaded ActorDatabase.
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
 | ENTITY CREATION                               |
 *───────────────────────────────────────────────*/

// Spawn instantiates an ECS entity by template name and position.
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

	// --- Actor Metadata ---
	actor := &ecs.Actor{
		ID:         c.nextID(template),
		Archetype:  tpl.Archetype,
		Persistent: tpl.Persistent,
	}
	e.Add(actor)

	// --- Transform Components ---
	e.Add(&ecs.Position{X: pos.X, Y: pos.Y})
	if tpl.Velocity != nil {
		e.Add(&ecs.Velocity{VX: tpl.Velocity.VX, VY: tpl.Velocity.VY})
	}

	// --- Sprite Component ---
	if tpl.Sprite.Image != "" {
		e.Add(buildSprite(tpl.Sprite))
	}

	// --- AI References (handled by AIComposer) ---
	if len(tpl.AIRefs) > 0 {
		actor.AIRefs = append([]string{}, tpl.AIRefs...)
	}

	return e, nil
}

/*───────────────────────────────────────────────*
 | IMAGE PRELOADING                              |
 *───────────────────────────────────────────────*/

// PreloadImages preloads all sprite textures used by templates.
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
 | UTILITIES                                     |
 *───────────────────────────────────────────────*/

// Templates lists all available template names (used by DevConsole).
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

// nextID generates a unique, per-template instance ID.
func (c *ActorCreator) nextID(template string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counters[template]++
	return fmt.Sprintf("%s-%03d", template, c.counters[template])
}

// buildSprite constructs an ECS sprite from a data template.
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

