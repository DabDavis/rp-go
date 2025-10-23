package ai

import (
	"math"

	"rp-go/engine/ecs"
)

func (s *System) behaviorPatrol(_ *ecs.World, _ *ecs.Entity, pos *ecs.Position, vel *ecs.Velocity, p map[string]any) bool {
	wpList, ok := p["waypoints"].([]any)
	if !ok || len(wpList) == 0 {
		return false
	}
	speed := getFloat(p, "speed", 2.0)
	wp := wpList[0].(map[string]any)
	tx := getFloat(wp, "x", pos.X)
	ty := getFloat(wp, "y", pos.Y)
	dx := tx - pos.X
	dy := ty - pos.Y
	dist := math.Hypot(dx, dy)
	if dist < 2 {
		return false
	}
	vel.VX = (dx / dist) * speed
	vel.VY = (dy / dist) * speed
	return true
}

