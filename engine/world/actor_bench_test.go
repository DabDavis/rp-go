package world

import (
	"fmt"
	"testing"

	"rp-go/engine/data"
	"rp-go/engine/ecs"
)

func mockBenchDB() data.ActorDatabase {
	const N = 5
	actors := make([]data.ActorTemplate, 0, N)
	for i := 0; i < N; i++ {
		name := fmt.Sprintf("enemy_type_%d", i)
		actors = append(actors, data.ActorTemplate{
			Name:       name,
			Archetype:  "enemy",
			Persistent: true,
			Sprite: data.ActorSpriteTemplate{
				Image:  "../../assets/entities/enemy.png",
				Width:  32,
				Height: 32,
			},
			Velocity: &data.ActorVelocity{VX: 1, VY: 1},
			AI: &data.ActorAITemplate{
				Speed: 2.0,
				Patrol: &data.ActorAIPatrol{
					Variant: "loop",
					Waypoints: []data.ActorAIWaypoint{
						{X: 0, Y: 0}, {X: 10, Y: 10}, {X: 20, Y: 20},
					},
					Speed: 1.5,
				},
			},
		})
	}
	return data.ActorDatabase{Actors: actors}
}

func BenchmarkActorCreator_Spawn(b *testing.B) {
	db := mockBenchDB()
	creator := NewActorCreator(db)
	world := ecs.NewWorld()
	pos := ecs.Position{X: 100, Y: 200}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := creator.Spawn(world, "enemy_type_0", pos)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBuildAIController(b *testing.B) {
	tpl := &data.ActorAITemplate{
		Speed: 2.5,
		Follow: &data.ActorAIFollow{
			Target: "player", Speed: 1.2,
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if ctrl := BuildAIController(tpl); ctrl == nil {
			b.Fatal("nil AIController")
		}
	}
}

func BenchmarkNextID(b *testing.B) {
	c := &ActorCreator{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = c.nextID("enemy")
		}
	})
}

