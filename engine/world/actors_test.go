package world

import (
	"fmt"
	"testing"

	"rp-go/engine/data"
	"rp-go/engine/ecs"
)

func TestActorCreatorUniqueIDs(t *testing.T) {
	db := data.ActorDatabase{Actors: []data.ActorTemplate{
		{
			Name:      "dark-elf-ship-1",
			Archetype: "enemy",
		},
	}}

	creator := NewActorCreator(db)
	world := ecs.NewWorld()

	seen := make(map[string]struct{})
	for i := 0; i < 5; i++ {
		entity, err := creator.Spawn(world, "dark-elf-ship-1", ecs.Position{})
		if err != nil {
			t.Fatalf("spawn failed: %v", err)
		}
		actorComp, _ := entity.Get("Actor").(*ecs.Actor)
		if actorComp == nil {
			t.Fatalf("actor component missing on entity %d", entity.ID)
		}

		if _, dup := seen[actorComp.ID]; dup {
			t.Fatalf("duplicate actor id generated: %s", actorComp.ID)
		}
		seen[actorComp.ID] = struct{}{}

		expected := fmt.Sprintf("dark-elf-ship-1-%04d", i+1)
		if actorComp.ID != expected {
			t.Fatalf("expected id %q, got %q", expected, actorComp.ID)
		}
	}
}

func TestActorCreatorBuildsAIComponent(t *testing.T) {
	db := data.ActorDatabase{Actors: []data.ActorTemplate{
		{
			Name:      "scout",
			Archetype: "enemy",
			AI: &data.ActorAITemplate{
				Speed:  3.5,
				Pursue: &data.ActorAIPursue{Target: "player", EngageDistance: 200, Speed: 4},
				Patrol: &data.ActorAIPatrol{
					Variant: "loop",
					Speed:   2.5,
					Waypoints: []data.ActorAIWaypoint{
						{X: 10, Y: 20},
						{X: 30, Y: 40},
					},
				},
			},
		},
	}}

	creator := NewActorCreator(db)
	world := ecs.NewWorld()
	entity, err := creator.Spawn(world, "scout", ecs.Position{X: 5, Y: 5})
	if err != nil {
		t.Fatalf("spawn failed: %v", err)
	}

	ai, ok := entity.Get("AIController").(*ecs.AIController)
	if !ok || ai == nil {
		t.Fatalf("expected AIController component on spawned entity")
	}

	if ai.Pursue == nil || ai.Pursue.Target != "player" {
		t.Fatalf("expected pursue behavior to target player")
	}

	if ai.Patrol == nil || len(ai.Patrol.Waypoints) != 2 {
		t.Fatalf("expected patrol waypoints to be copied")
	}

	if !ai.PatrolState.Forward {
		t.Fatalf("expected patrol state to be initialized in forward direction")
	}

	if _, ok := entity.Get("Velocity").(*ecs.Velocity); !ok {
		t.Fatalf("expected velocity component to be present for AI-controlled actor")
	}
}
