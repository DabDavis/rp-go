package actor

import (
	"sort"
	"strings"

	"rp-go/engine/ecs"
)

// Registry maintains indexed views over actors currently registered in the world.
type Registry struct {
	byID        map[string]*ecs.Entity
	byArchetype map[string][]*ecs.Entity
}

// NewRegistry constructs an empty actor registry ready for reuse.
func NewRegistry() *Registry {
	return &Registry{
		byID:        make(map[string]*ecs.Entity),
		byArchetype: make(map[string][]*ecs.Entity),
	}
}

// Reset clears all cached actor indices.
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
}

// Add indexes an actor entity for lookups by ID and archetype.
func (r *Registry) Add(actor *ecs.Actor, entity *ecs.Entity) {
	if r == nil || actor == nil || entity == nil || actor.ID == "" {
		return
	}
	r.byID[actor.ID] = entity
	if actor.Archetype != "" {
		r.byArchetype[actor.Archetype] = append(r.byArchetype[actor.Archetype], entity)
	}
}

// FindByID retrieves the entity associated with the supplied actor ID.
func (r *Registry) FindByID(id string) (*ecs.Entity, bool) {
	if r == nil {
		return nil, false
	}
	entity, ok := r.byID[id]
	return entity, ok
}

// FindByArchetype returns the entities for all actors matching the archetype.
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

// FindByTemplatePrefix returns entities whose actor IDs share the provided prefix.
func (r *Registry) FindByTemplatePrefix(prefix string) []*ecs.Entity {
	if r == nil || prefix == "" {
		return nil
	}
	matches := make([]*ecs.Entity, 0, len(r.byID))
	for id, entity := range r.byID {
		if strings.HasPrefix(id, prefix) {
			matches = append(matches, entity)
		}
	}
	sort.Slice(matches, func(i, j int) bool { return matches[i].ID < matches[j].ID })
	return matches
}

// System keeps an index of all Actor components for other systems to query.
type System struct {
	registry *Registry
}

// NewSystem constructs an Actor system with an empty registry.
func NewSystem() *System {
	return &System{registry: NewRegistry()}
}

// Registry exposes the actor registry so other systems can reuse it without
// iterating the entire entity list each frame.
func (s *System) Registry() *Registry {
	if s.registry == nil {
		s.registry = NewRegistry()
	}
	return s.registry
}

func (s *System) Update(w *ecs.World) {
	if s.registry == nil {
		s.registry = NewRegistry()
	}
	s.registry.Reset()

	primaryAssigned := false

	for _, e := range w.Entities {
		actor, ok := e.Get("Actor").(*ecs.Actor)
		if !ok || actor == nil {
			continue
		}

		s.registry.Add(actor, e)

		if controller, ok := e.Get("PlayerInput").(*ecs.PlayerInput); ok && controller != nil {
			if _, hasAI := e.Get("AIController").(*ecs.AIController); hasAI {
				controller.Enabled = false
				continue
			}
			if !primaryAssigned {
				controller.Enabled = true
				primaryAssigned = true
			} else {
				controller.Enabled = false
			}
		}
	}
}
