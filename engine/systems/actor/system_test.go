package actor

import (
	"testing"

	"rp-go/engine/ecs"
)

func TestRegistryIndexesActors(t *testing.T) {
	w := ecs.NewWorld()

	player := w.NewEntity()
	player.Add(&ecs.Actor{ID: "player", Archetype: "player"})

	enemyA := w.NewEntity()
	enemyA.Add(&ecs.Actor{ID: "dark-elf-ship-scout-0001", Archetype: "enemy"})

	enemyB := w.NewEntity()
	enemyB.Add(&ecs.Actor{ID: "dark-elf-ship-raider-0002", Archetype: "enemy"})

	system := NewSystem()
	system.Update(w)

	registry := system.Registry()

	foundPlayer, ok := registry.FindByID("player")
	if !ok || foundPlayer != player {
		t.Fatalf("expected to locate player entity, got %v", foundPlayer)
	}

	enemies := registry.FindByArchetype("enemy")
	if len(enemies) != 2 {
		t.Fatalf("expected two enemy actors, got %d", len(enemies))
	}
	if enemies[0].ID > enemies[1].ID {
		t.Fatalf("expected enemies slice to be sorted by entity ID")
	}

	templated := registry.FindByTemplatePrefix("dark-elf-ship-raider")
	if len(templated) != 1 || templated[0] != enemyB {
		t.Fatalf("expected to find raider template match, got %+v", templated)
	}
}

func TestActorSystemEnforcesSinglePlayerInput(t *testing.T) {
	w := ecs.NewWorld()

	pilotOne := w.NewEntity()
	pilotOne.Add(&ecs.Actor{ID: "pilot-one", Archetype: "ship"})
	pilotOne.Add(&ecs.PlayerInput{Enabled: true})

	pilotTwo := w.NewEntity()
	pilotTwo.Add(&ecs.Actor{ID: "pilot-two", Archetype: "ship"})
	pilotTwo.Add(&ecs.PlayerInput{Enabled: true})

	autopilot := w.NewEntity()
	autopilot.Add(&ecs.Actor{ID: "auto", Archetype: "ship"})
	autopilot.Add(&ecs.PlayerInput{Enabled: true})
	autopilot.Add(&ecs.AIController{})

	system := NewSystem()
	system.Update(w)

	if ctrl, _ := pilotOne.Get("PlayerInput").(*ecs.PlayerInput); ctrl == nil || !ctrl.Enabled {
		t.Fatalf("expected first player input to remain enabled")
	}

	if ctrl, _ := pilotTwo.Get("PlayerInput").(*ecs.PlayerInput); ctrl == nil || ctrl.Enabled {
		t.Fatalf("expected second player input to be disabled")
	}

	if ctrl, _ := autopilot.Get("PlayerInput").(*ecs.PlayerInput); ctrl == nil || ctrl.Enabled {
		t.Fatalf("expected entities with AI to have player input disabled")
	}
}
