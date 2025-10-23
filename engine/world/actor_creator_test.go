package world

import (
	"testing"

	"rp-go/engine/data"
	"rp-go/engine/ecs"
)

func TestNewActorCreatorAndSpawn(t *testing.T) {
	db := data.ActorDatabase{
		Actors: []data.ActorTemplate{
			{
				Name:      "drone",
				Archetype: "enemy",
				Sprite: data.ActorSpriteTemplate{
					Image:  "../../assets/entities/drone.png",
					Width:  32,
					Height: 32,
				},
				Velocity: &data.ActorVelocity{VX: 1, VY: 2},
				AI: &data.ActorAITemplate{
					Speed: 3.0,
				},
			},
		},
	}

	ac := NewActorCreator(db)
	if ac == nil {
		t.Fatal("expected ActorCreator instance")
	}

	w := ecs.NewWorld()
	entity, err := ac.Spawn(w, "drone", ecs.Position{X: 10, Y: 20})
	if err != nil {
		t.Fatalf("spawn error: %v", err)
	}
	if entity == nil {
		t.Fatal("spawn returned nil entity")
	}
}

