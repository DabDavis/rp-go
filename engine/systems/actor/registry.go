package actor

import (
	"sort"
	"strings"

	"rp-go/engine/ecs"
)

// Registry maintains indexed and sorted views over all Actor entities.
// It allows other systems (AI, debug overlays, UI, etc.) to efficiently
// query actors by ID, archetype, or template prefix without scanning
// the entire ECS world.
type Registry struct {
	byID        map[string]*ecs.Entity
	byArchetype map[string][]*ecs.Entity
	all         []*ecs.Entity
}

// NewRegistry constructs an empty registry ready for use.
func NewRegistry() *Registry {
	return &Registry{
		byID:        make(map[string]*ecs.Entity),
		byArchetype: make(map[string][]*ecs.Entity),
		all:         make([]*ecs.Entity, 0),
	}
}

// Reset clears all cached indices and slices for reuse.
func (r *Registry) Reset() {
	if r == nil {
		return
	}
	for id := range r.byID {
		delete(r.byID, id)
	}
	for archetype := range r.byArchetype {
		r.byArchetype[archetype] = r.byArchetype[archetype][:0]
	}
	r.all = r.all[:0]
}

// Add indexes an actor entity by its ID and archetype, caching it globally.
func (r *Registry) Add(actor *ecs.Actor, entity *ecs.Entity) {
	if r == nil || actor == nil || entity == nil || actor.ID == "" {
		return
	}
	r.byID[actor.ID] = entity
	if actor.Archetype != "" {
		r.byArchetype[actor.Archetype] = append(r.byArchetype[actor.Archetype], entity)
	}
	r.all = append(r.all, entity)
}

// FindByID retrieves the entity associated with a given actor ID.
func (r *Registry) FindByID(id string) (*ecs.Entity, bool) {
	if r == nil {
		return nil, false
	}
	e, ok := r.byID[id]
	return e, ok
}

// FindByArchetype returns all entities matching a given archetype, sorted by ID.
func (r *Registry) FindByArchetype(archetype string) []*ecs.Entity {
	if r == nil {
		return nil
	}
	slice := r.byArchetype[archetype]
	if len(slice) == 0 {
		return nil
	}
	out := make([]*ecs.Entity, len(slice))
	copy(out, slice)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// FindByTemplatePrefix returns all entities whose actor IDs share a prefix.
// This is useful for debugging or referencing templated spawns (e.g. "drone-", "npc-").
func (r *Registry) FindByTemplatePrefix(prefix string) []*ecs.Entity {
	if r == nil || prefix == "" {
		return nil
	}
	matches := make([]*ecs.Entity, 0, len(r.byID))
	for id, e := range r.byID {
		if strings.HasPrefix(id, prefix) {
			matches = append(matches, e)
		}
	}
	sort.Slice(matches, func(i, j int) bool { return matches[i].ID < matches[j].ID })
	return matches
}

// Entities returns all registered entities as a sorted, independent slice.
func (r *Registry) Entities() []*ecs.Entity {
	if r == nil || len(r.all) == 0 {
		return nil
	}

	out := make([]*ecs.Entity, len(r.all))
	copy(out, r.all)
	sort.Slice(out, func(i, j int) bool {
		ai, _ := out[i].Get("Actor").(*ecs.Actor)
		aj, _ := out[j].Get("Actor").(*ecs.Actor)
		if ai == nil || aj == nil {
			return out[i].ID < out[j].ID
		}
		if ai.ID == aj.ID {
			return out[i].ID < out[j].ID
		}
		return ai.ID < aj.ID
	})
	return out
}

// All returns every registered actor entity sorted by ID.
// Deprecated: use Entities() instead.
func (r *Registry) All() []*ecs.Entity {
	return r.Entities()
}

