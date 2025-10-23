package ai

import (
	"math"
	"math/rand"
	"strings"
	"time"

	"rp-go/engine/ecs"
	"rp-go/engine/platform"
)

const waypointTolerance = 4.0

// ActorLookup exposes indexed queries over actor entities.
type ActorLookup interface {
	FindByID(id string) (*ecs.Entity, bool)
	FindByArchetype(archetype string) []*ecs.Entity
	FindByTemplatePrefix(prefix string) []*ecs.Entity
}

// System drives AI-controlled actors each frame.
type System struct {
	rng    *rand.Rand
	lookup ActorLookup
}

// NewSystem constructs an AI system with its own random source.
func NewSystem() *System {
	return &System{rng: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (s *System) ensureRNG() {
	if s.rng == nil {
		s.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
}

// SetActorLookup shares the actor registry so the AI can resolve targets
// without iterating over every entity in the world each frame.
func (s *System) SetActorLookup(lookup ActorLookup) {
	s.lookup = lookup
}

func (s *System) Update(w *ecs.World) {
	s.ensureRNG()
	for _, e := range w.Entities {
		ai, ok := e.Get("AIController").(*ecs.AIController)
		if !ok || ai == nil {
			continue
		}

		pos, hasPos := e.Get("Position").(*ecs.Position)
		if !hasPos {
			continue
		}

		vel, hasVel := e.Get("Velocity").(*ecs.Velocity)
		if !hasVel {
			vel = &ecs.Velocity{}
			e.Add(vel)
		}

		vel.VX, vel.VY = 0, 0

		if ai.Retreat != nil && s.applyRetreat(w, ai, pos, vel) {
			continue
		}
		if ai.Travel != nil && s.applyPathBehavior(&ai.TravelState, ai.Travel, pos, vel, ai.SpeedFor(ai.Travel.Speed)) {
			continue
		}
		if ai.Patrol != nil && s.applyPathBehavior(&ai.PatrolState, ai.Patrol, pos, vel, ai.SpeedFor(ai.Patrol.Speed)) {
			continue
		}
		if ai.Pursue != nil && s.applyPursue(w, ai, pos, vel) {
			continue
		}
		if ai.Follow != nil && s.applyFollow(w, ai, pos, vel) {
			continue
		}
	}
}

func (s *System) Draw(*ecs.World, *platform.Image) {}

func (s *System) applyFollow(w *ecs.World, ai *ecs.AIController, pos *ecs.Position, vel *ecs.Velocity) bool {
	cfg := ai.Follow
	if cfg == nil || cfg.Target == "" {
		return false
	}

	targetPos, ok := s.findTargetPosition(w, cfg.Target)
	if !ok {
		return false
	}

	desiredX := targetPos.X + cfg.OffsetX
	desiredY := targetPos.Y + cfg.OffsetY
	dx := desiredX - pos.X
	dy := desiredY - pos.Y
	dist := math.Hypot(dx, dy)

	minDist := cfg.MinDistance
	if minDist <= 0 {
		minDist = 1
	}

	if dist <= minDist {
		vel.VX, vel.VY = 0, 0
		return true
	}

	maxDist := cfg.MaxDistance
	if maxDist > 0 && maxDist < minDist {
		maxDist = minDist
	}

	speed := ai.SpeedFor(cfg.Speed)
	if maxDist > minDist && dist < maxDist {
		factor := (dist - minDist) / (maxDist - minDist)
		if factor < 0 {
			factor = 0
		}
		speed *= factor
	}
	if dist < speed {
		speed = dist
	}
	if dist == 0 {
		vel.VX, vel.VY = 0, 0
		return true
	}

	if speed <= 0 {
		vel.VX, vel.VY = 0, 0
		return true
	}

	vel.VX = (dx / dist) * speed
	vel.VY = (dy / dist) * speed
	return true
}

func (s *System) applyPursue(w *ecs.World, ai *ecs.AIController, pos *ecs.Position, vel *ecs.Velocity) bool {
	cfg := ai.Pursue
	if cfg == nil || cfg.Target == "" {
		return false
	}

	targetPos, ok := s.findTargetPosition(w, cfg.Target)
	if !ok {
		return false
	}

	dx := targetPos.X - pos.X
	dy := targetPos.Y - pos.Y
	dist := math.Hypot(dx, dy)

	if cfg.EngageDistance > 0 && dist > cfg.EngageDistance {
		return false
	}
	if dist == 0 {
		vel.VX, vel.VY = 0, 0
		return true
	}

	speed := ai.SpeedFor(cfg.Speed)
	if dist < speed {
		speed = dist
	}
	vel.VX = (dx / dist) * speed
	vel.VY = (dy / dist) * speed
	return true
}

func (s *System) applyRetreat(w *ecs.World, ai *ecs.AIController, pos *ecs.Position, vel *ecs.Velocity) bool {
	cfg := ai.Retreat
	if cfg == nil || cfg.Target == "" {
		return false
	}

	targetPos, ok := s.findTargetPosition(w, cfg.Target)
	if !ok {
		return false
	}

	dx := pos.X - targetPos.X
	dy := pos.Y - targetPos.Y
	dist := math.Hypot(dx, dy)

	trigger := cfg.TriggerDistance
	if trigger <= 0 {
		trigger = 150
	}

	if dist > trigger {
		if cfg.SafeDistance > 0 && dist < cfg.SafeDistance {
			// continue retreating until reaching safety
		} else {
			return false
		}
	}

	safe := cfg.SafeDistance
	if safe <= trigger {
		safe = trigger * 1.25
	}

	if dist >= safe {
		vel.VX, vel.VY = 0, 0
		return true
	}

	if dist == 0 {
		vel.VX = ai.SpeedFor(cfg.Speed)
		vel.VY = 0
		return true
	}

	speed := ai.SpeedFor(cfg.Speed)
	if dist < speed {
		speed = dist
	}
	vel.VX = (dx / dist) * speed
	vel.VY = (dy / dist) * speed
	return true
}

func (s *System) applyPathBehavior(state *ecs.AIPathState, cfg *ecs.AIPathBehavior, pos *ecs.Position, vel *ecs.Velocity, speed float64) bool {
	if cfg == nil || len(cfg.Waypoints) == 0 || state == nil {
		return false
	}

	total := len(cfg.Waypoints)
	if total == 0 {
		return false
	}

	if state.Index < 0 || state.Index >= total {
		if total == 0 {
			return false
		}
		state.Index = state.Index % total
		if state.Index < 0 {
			state.Index += total
		}
	}

	variant := strings.ToLower(cfg.Variant)
	if variant == "" {
		variant = "loop"
	}

	if variant == "once" && state.Completed {
		vel.VX, vel.VY = 0, 0
		return true
	}

	target := cfg.Waypoints[state.Index]
	dx := target.X - pos.X
	dy := target.Y - pos.Y
	dist := math.Hypot(dx, dy)

	if dist <= waypointTolerance {
		s.advancePath(state, variant, total)
		if variant == "once" && state.Completed {
			vel.VX, vel.VY = 0, 0
			return true
		}
		target = cfg.Waypoints[state.Index]
		dx = target.X - pos.X
		dy = target.Y - pos.Y
		dist = math.Hypot(dx, dy)
	}

	if dist == 0 {
		vel.VX, vel.VY = 0, 0
		return true
	}

	moveSpeed := speed
	if moveSpeed <= 0 {
		moveSpeed = ecs.DefaultAISpeed
	}
	if dist < moveSpeed {
		moveSpeed = dist
	}

	vel.VX = (dx / dist) * moveSpeed
	vel.VY = (dy / dist) * moveSpeed
	return true
}

func (s *System) advancePath(state *ecs.AIPathState, variant string, total int) {
	if total == 0 {
		return
	}
	switch variant {
	case "pingpong":
		if total == 1 {
			state.Completed = true
			return
		}
		if !state.Forward {
			if state.Index <= 0 {
				state.Forward = true
				state.Index = 1
			} else {
				state.Index--
			}
			return
		}
		if state.Index >= total-1 {
			state.Forward = false
			state.Index = total - 2
			if state.Index < 0 {
				state.Index = 0
			}
			return
		}
		state.Index++
	case "once":
		if state.Index >= total-1 {
			state.Completed = true
			return
		}
		state.Index++
	case "random":
		if total == 1 {
			state.Completed = true
			return
		}
		next := state.Index
		for total > 1 && next == state.Index {
			next = s.rng.Intn(total)
		}
		state.Index = next
	default: // loop
		state.Index = (state.Index + 1) % total
	}
}

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
		archetype := strings.TrimPrefix(target, "archetype:")
		candidates = s.lookup.FindByArchetype(archetype)
	case strings.HasPrefix(target, "template:"):
		prefix := strings.TrimPrefix(target, "template:")
		candidates = s.lookup.FindByTemplatePrefix(prefix)
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

	selector := func(*ecs.Actor) bool { return false }
	switch {
	case strings.HasPrefix(target, "archetype:"):
		archetype := strings.TrimPrefix(target, "archetype:")
		selector = func(a *ecs.Actor) bool { return a.Archetype == archetype }
	case strings.HasPrefix(target, "template:"):
		template := strings.TrimPrefix(target, "template:")
		selector = func(a *ecs.Actor) bool {
			return strings.HasPrefix(a.ID, template)
		}
	default:
		selector = func(a *ecs.Actor) bool { return a.ID == target }
	}

	for _, e := range w.Entities {
		actor, ok := e.Get("Actor").(*ecs.Actor)
		if !ok || actor == nil {
			continue
		}
		if !selector(actor) {
			continue
		}
		pos, ok := e.Get("Position").(*ecs.Position)
		if ok {
			return pos, true
		}
	}
	return nil, false
}
