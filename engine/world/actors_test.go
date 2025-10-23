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

func TestActorCreatorAddsAIController(t *testing.T) {
	db := data.ActorDatabase{Actors: []data.ActorTemplate{
		{
			Name:      "ai-test",
			Archetype: "enemy",
			AI: &data.ActorAITemplate{
				Pursue: &data.ActorAIPursue{Target: "player", Speed: 3.5},
			},
		},
	}}

	creator := NewActorCreator(db)
	world := ecs.NewWorld()

	entity, err := creator.Spawn(world, "ai-test", ecs.Position{})
	if err != nil {
		t.Fatalf("spawn failed: %v", err)
	}

	controller, ok := entity.Get("AIController").(*ecs.AIController)
	if !ok || controller == nil {
		t.Fatalf("expected AIController component on spawned entity")
	}
	if controller.Pursue == nil || controller.Pursue.Target != "player" {
		t.Fatalf("expected pursue behavior to be configured on controller, got %+v", controller.Pursue)
	}
}
