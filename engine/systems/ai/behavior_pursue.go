package ai

import (
	"math"

	"rp-go/engine/ecs"
)

func (s *System) behaviorPursue(w *ecs.World, e *ecs.Entity, pos *ecs.Position, vel *ecs.Velocity, p map[string]any) bool {
	targetName, _ := p["target"].(string)
	speed := getFloat(p, "speed", 2.0)
	maxDist := getFloat(p, "engage_distance", 300)
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
	dx := tp.X - pos.X
	dy := tp.Y - pos.Y
	dist := math.Hypot(dx, dy)
	if dist > maxDist || dist < 1 {
		return false
	}
	vel.VX = (dx / dist) * speed
	vel.VY = (dy / dist) * speed
	return true
}

