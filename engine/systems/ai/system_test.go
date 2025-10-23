package ai

import (
	"testing"

	"rp-go/engine/ecs"
	actor "rp-go/engine/systems/actor"
)

func TestAISystemFollowBehavior(t *testing.T) {
	w := ecs.NewWorld()

	target := w.NewEntity()
	target.Add(&ecs.Actor{ID: "player", Archetype: "player"})
	target.Add(&ecs.Position{X: 10, Y: 0})

	follower := w.NewEntity()
	follower.Add(&ecs.Position{X: 100, Y: 0})
	follower.Add(&ecs.Velocity{})
	follower.Add(&ecs.AIController{
		Speed: 2.5,
		Follow: &ecs.AIFollowBehavior{
			Target:      "player",
			MinDistance: 4,
			Speed:       2,
		},
	})

	system := NewSystem()
	system.Update(w)

	vel, _ := follower.Get("Velocity").(*ecs.Velocity)
	if vel == nil {
		t.Fatalf("expected velocity component")
	}
	if vel.VX >= 0 {
		t.Fatalf("expected follower to move toward the player (negative X velocity), got %f", vel.VX)
	}
}

func TestAISystemPursueEngageDistance(t *testing.T) {
	w := ecs.NewWorld()

	target := w.NewEntity()
	target.Add(&ecs.Actor{ID: "player"})
	target.Add(&ecs.Position{X: 200, Y: 0})

	chaser := w.NewEntity()
	chaser.Add(&ecs.Position{X: 0, Y: 0})
	chaser.Add(&ecs.Velocity{})
	chaser.Add(&ecs.AIController{
		Speed: 3,
		Pursue: &ecs.AIPursueBehavior{
			Target:         "player",
			EngageDistance: 100,
			Speed:          4,
		},
	})

	system := NewSystem()
	system.Update(w)

	vel, _ := chaser.Get("Velocity").(*ecs.Velocity)
	if vel == nil {
		t.Fatalf("expected velocity component on chaser")
	}
	if vel.VX != 0 || vel.VY != 0 {
		t.Fatalf("expected no pursuit outside engage distance, got vx=%f vy=%f", vel.VX, vel.VY)
	}

	targetPos, _ := target.Get("Position").(*ecs.Position)
	targetPos.X = 60

	system.Update(w)
	if vel.VX <= 0 {
		t.Fatalf("expected chaser to move toward the player after entering range, got %f", vel.VX)
	}
}

func TestAISystemPatrolPingPong(t *testing.T) {
	w := ecs.NewWorld()

	entity := w.NewEntity()
	entity.Add(&ecs.Position{X: 260, Y: 220})
	entity.Add(&ecs.Velocity{})

	ai := &ecs.AIController{
		Speed: 1.5,
		Patrol: &ecs.AIPathBehavior{
			Variant: "pingpong",
			Speed:   1.5,
			Waypoints: []ecs.AIWaypoint{
				{X: 260, Y: 220},
				{X: 300, Y: 220},
			},
		},
	}
	ai.PatrolState.Reset()
	entity.Add(ai)

	system := NewSystem()
	system.Update(w)

	if ai.PatrolState.Index != 1 {
		t.Fatalf("expected patrol to advance to second waypoint, got index %d", ai.PatrolState.Index)
	}
	vel, _ := entity.Get("Velocity").(*ecs.Velocity)
	if vel == nil || vel.VX <= 0 {
		t.Fatalf("expected positive X velocity toward second waypoint, got %+v", vel)
	}

	pos, _ := entity.Get("Position").(*ecs.Position)
	pos.X = 300
	pos.Y = 220

	system.Update(w)
	if ai.PatrolState.Index != 0 {
		t.Fatalf("expected patrol to bounce back to first waypoint, got index %d", ai.PatrolState.Index)
	}
	if ai.PatrolState.Forward {
		t.Fatalf("expected patrol direction to reverse after reaching end")
	}
}

func TestAISystemRetreatStopsAfterSafeDistance(t *testing.T) {
	w := ecs.NewWorld()

	threat := w.NewEntity()
	threat.Add(&ecs.Actor{ID: "player"})
	threat.Add(&ecs.Position{X: 30, Y: 0})

	runner := w.NewEntity()
	runner.Add(&ecs.Position{X: 0, Y: 0})
	runner.Add(&ecs.Velocity{})
	runner.Add(&ecs.AIController{
		Speed: 2.8,
		Retreat: &ecs.AIRetreatBehavior{
			Target:          "player",
			TriggerDistance: 50,
			SafeDistance:    120,
			Speed:           3.1,
		},
	})

	system := NewSystem()
	system.Update(w)

	vel, _ := runner.Get("Velocity").(*ecs.Velocity)
	if vel == nil || vel.VX >= 0 {
		t.Fatalf("expected runner to flee away from player, got %+v", vel)
	}

	threatPos, _ := threat.Get("Position").(*ecs.Position)
	threatPos.X = 200

	system.Update(w)
	if vel.VX != 0 || vel.VY != 0 {
		t.Fatalf("expected runner to stop after reaching safety, got %+v", vel)
	}
}

func TestAISystemTravelOnceCompletes(t *testing.T) {
	w := ecs.NewWorld()

	traveler := w.NewEntity()
	traveler.Add(&ecs.Position{X: 200, Y: 120})
	traveler.Add(&ecs.Velocity{})

	ai := &ecs.AIController{
		Speed: 2.5,
		Travel: &ecs.AIPathBehavior{
			Variant: "once",
			Speed:   3,
			Waypoints: []ecs.AIWaypoint{
				{X: 120, Y: 120},
				{X: 200, Y: 120},
			},
		},
	}
	ai.TravelState.Reset()
	ai.TravelState.Index = 1
	traveler.Add(ai)

	system := NewSystem()
	system.Update(w)

	if !ai.TravelState.Completed {
		t.Fatalf("expected travel behavior to mark path as completed")
	}

	vel, _ := traveler.Get("Velocity").(*ecs.Velocity)
	if vel == nil || vel.VX != 0 || vel.VY != 0 {
		t.Fatalf("expected traveler to stop after completing path, got %+v", vel)
	}
}

func TestAISystemUsesActorLookup(t *testing.T) {
	w := ecs.NewWorld()

	actorSystem := actor.NewSystem()

	target := w.NewEntity()
	target.Add(&ecs.Actor{ID: "player", Archetype: "player"})
	target.Add(&ecs.Position{X: 50, Y: 0})

	pursuer := w.NewEntity()
	pursuer.Add(&ecs.Position{X: 0, Y: 0})
	pursuer.Add(&ecs.Velocity{})
	pursuer.Add(&ecs.AIController{
		Pursue: &ecs.AIPursueBehavior{
			Target: "player",
			Speed:  2,
		},
	})

	actorSystem.Update(w)

	system := NewSystem()
	system.SetActorLookup(actorSystem.Registry())
	system.Update(w)

	vel, _ := pursuer.Get("Velocity").(*ecs.Velocity)
	if vel == nil || vel.VX <= 0 {
		t.Fatalf("expected pursuer to chase player using registry lookup, got %+v", vel)
	}
}
