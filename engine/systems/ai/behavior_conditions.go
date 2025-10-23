package ai

import (
	"math"

	"rp-go/engine/ecs"
)

// ConditionSet defines a set of preconditions before an action triggers.
type ConditionSet struct {
	Target     string   `json:"target,omitempty"`
	Within     float64  `json:"within,omitempty"`  // distance must be less than
	Beyond     float64  `json:"beyond,omitempty"`  // distance must be greater than
	HealthLess float64  `json:"health_lt,omitempty"`
	HealthMore float64  `json:"health_gt,omitempty"`
}

// checkConditions returns true if all conditions are satisfied.
func (s *System) checkConditions(w *ecs.World, e *ecs.Entity, cond map[string]any) bool {
	if len(cond) == 0 {
		return true
	}

	// Convert to struct for clarity
	c := ConditionSet{}
	if v, ok := cond["target"].(string); ok {
		c.Target = v
	}
	c.Within = getFloat(cond, "within", 0)
	c.Beyond = getFloat(cond, "beyond", 0)
	c.HealthLess = getFloat(cond, "health_lt", 0)
	c.HealthMore = getFloat(cond, "health_gt", 0)

	// Check health
	if hp, ok := e.Get("Health").(*ecs.Health); ok && hp != nil {
		val := hp.Value / hp.Max
		if c.HealthLess > 0 && val >= c.HealthLess {
			return false
		}
		if c.HealthMore > 0 && val <= c.HealthMore {
			return false
		}
	}

	// Check distance to target
	if c.Target != "" && (c.Within > 0 || c.Beyond > 0) {
		var target *ecs.Entity
		w.EntitiesManager().ForEach(func(ent *ecs.Entity) {
			if target != nil {
				return
			}
			act, _ := ent.Get("Actor").(*ecs.Actor)
			if act != nil && act.ID == c.Target {
				target = ent
			}
		})
		if target != nil {
			tp, _ := target.Get("Position").(*ecs.Position)
			p1, _ := e.Get("Position").(*ecs.Position)
			if tp != nil && p1 != nil {
				d := math.Hypot(tp.X-p1.X, tp.Y-p1.Y)
				if c.Within > 0 && d > c.Within {
					return false
				}
				if c.Beyond > 0 && d < c.Beyond {
					return false
				}
			}
		}
	}

	return true
}

