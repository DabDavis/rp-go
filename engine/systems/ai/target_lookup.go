package ai

import (
	"strings"

	"rp-go/engine/ecs"
)

func (s *System) findTargetPosition(w *ecs.World, target string) (*ecs.Position, bool) {
	if target == "" {
		return nil, false
	}

	if pos, ok := s.lookupTargetPosition(target); ok {
		return pos, true
	}
	return fallbackFindTargetPosition(w, target)
}

func (s *System) lookupTargetPosition(target string) (*ecs.Position, bool) {
	if s.lookup == nil {
		return nil, false
	}

	var candidates []*ecs.Entity
	switch {
	case strings.HasPrefix(target, "archetype:"):
		candidates = s.lookup.FindByArchetype(strings.TrimPrefix(target, "archetype:"))
	case strings.HasPrefix(target, "template:"):
		candidates = s.lookup.FindByTemplatePrefix(strings.TrimPrefix(target, "template:"))
	default:
		if entity, ok := s.lookup.FindByID(target); ok && entity != nil {
			candidates = []*ecs.Entity{entity}
		}
	}

	for _, entity := range candidates {
		if pos, ok := entity.Get("Position").(*ecs.Position); ok {
			return pos, true
		}
	}
	return nil, false
}

func fallbackFindTargetPosition(w *ecs.World, target string) (*ecs.Position, bool) {
	if w == nil || target == "" {
		return nil, false
	}

	var selector func(*ecs.Actor) bool
	switch {
	case strings.HasPrefix(target, "archetype:"):
		val := strings.TrimPrefix(target, "archetype:")
		selector = func(a *ecs.Actor) bool { return a.Archetype == val }
	case strings.HasPrefix(target, "template:"):
		val := strings.TrimPrefix(target, "template:")
		selector = func(a *ecs.Actor) bool { return strings.HasPrefix(a.ID, val) }
	default:
		selector = func(a *ecs.Actor) bool { return a.ID == target }
	}

	manager := w.EntitiesManager()
	if manager == nil {
		return nil, false
	}

	foundPos := (*ecs.Position)(nil)
	manager.ForEach(func(e *ecs.Entity) {
		if foundPos != nil {
			return
		}
		actor, ok := e.Get("Actor").(*ecs.Actor)
		if !ok || actor == nil || !selector(actor) {
			return
		}
		if pos, ok := e.Get("Position").(*ecs.Position); ok {
			foundPos = pos
		}
	})
	if foundPos != nil {
		return foundPos, true
	}
	return nil, false
}
