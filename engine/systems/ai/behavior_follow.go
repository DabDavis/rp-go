package ai

import (
	"math"

	"rp-go/engine/ecs"
)

func (s *System) behaviorFollow(w *ecs.World, e *ecs.Entity, pos *ecs.Position, vel *ecs.Velocity, p map[string]any) bool {
	targetName, _ := p["target"].(string)
	speed := getFloat(p, "speed", 2.2)
	offsetX := getFloat(p, "offset_x", 0)
	offsetY := getFloat(p, "offset_y", 0)
	minDist := getFloat(p, "min_distance", 32)

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
	tx := tp.X + offsetX
	ty := tp.Y + offsetY
	dx := tx - pos.X
	dy := ty - pos.Y
	dist := math.Hypot(dx, dy)
	if dist < minDist {
		return false
	}
	vel.VX = (dx / dist) * speed
	vel.VY = (dy / dist) * speed
	return true
}

