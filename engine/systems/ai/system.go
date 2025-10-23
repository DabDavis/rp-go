package ai

import (
    "math/rand"
    "time"

    "rp-go/engine/ecs"
    "rp-go/engine/platform"
)

// ActorLookup allows indexed access to actors from the AI system.
type ActorLookup interface {
    FindByID(id string) (*ecs.Entity, bool)
    FindByArchetype(archetype string) []*ecs.Entity
    FindByTemplatePrefix(prefix string) []*ecs.Entity
}

// System updates all entities with AIController components.
type System struct {
    rng    *rand.Rand
    lookup ActorLookup
}

// NewSystem constructs a new AI system with its own RNG.
func NewSystem() *System {
    return &System{rng: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

// SetActorLookup injects a registry for target resolution.
func (s *System) SetActorLookup(lookup ActorLookup) {
    s.lookup = lookup
}

// Update runs per-frame AI logic for all AI-controlled entities.
func (s *System) Update(w *ecs.World) {
    s.ensureRNG()

    for _, e := range w.Entities {
        ai, ok := e.Get("AIController").(*ecs.AIController)
        if !ok || ai == nil {
            continue
        }

        pos, ok := e.Get("Position").(*ecs.Position)
        if !ok {
            continue
        }

        vel, ok := e.Get("Velocity").(*ecs.Velocity)
        if !ok {
            vel = &ecs.Velocity{}
            e.Add(vel)
        }

        vel.VX, vel.VY = 0, 0

        // Behavior priority: Retreat > Travel > Patrol > Pursue > Follow
        if ai.Retreat != nil && s.applyRetreatBehavior(w, ai, pos, vel) {
            continue
        }
        if ai.Travel != nil && s.applyPathBehavior(&ai.TravelState, ai.Travel, pos, vel, ai.SpeedFor(ai.Travel.Speed)) {
            continue
        }
        if ai.Patrol != nil && s.applyPathBehavior(&ai.PatrolState, ai.Patrol, pos, vel, ai.SpeedFor(ai.Patrol.Speed)) {
            continue
        }
        if ai.Pursue != nil && s.applyPursueBehavior(w, ai, pos, vel) {
            continue
        }
        if ai.Follow != nil && s.applyFollowBehavior(w, ai, pos, vel) {
            continue
        }
    }
}

func (s *System) Draw(*ecs.World, *platform.Image) {}

