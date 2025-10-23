package ai

import (
	"math"

	"rp-go/engine/ecs"
)

func (s *System) behaviorRetreat(w *ecs.World, e *ecs.Entity, pos *ecs.Position, vel *ecs.Velocity, p map[string]any) bool {
	targetName, _ := p["target"].(string)
	trigger := getFloat(p, "trigger_distance", 200)
	safe := getFloat(p, "safe_distance", 320)
	speed := getFloat(p, "speed", 2.4)
	if targetName == "" {
		return false
	}

	var target *ecs.Entity
	w.EntitiesManager().ForEach(func(ent *ecs.Entity) {
		if target != nil {
			return
		}
		act, _ := ent.Get("Actor").(*ecs.Actor)
		if act != nil && act.ID == targetName {
			target = ent
		}
	})
	if target == nil {
		return false
	}
	tp, _ := target.Get("Position").(*ecs.Position)
	if tp == nil {
		return false
	}
	dx := pos.X - tp.X
	dy := pos.Y - tp.Y
	dist := math.Hypot(dx, dy)
	if dist > safe {
		return false
	}
	if dist < trigger {
		vel.VX = (dx / dist) * speed
		vel.VY = (dy / dist) * speed
		return true
	}
	return false
}

