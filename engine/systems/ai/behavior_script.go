package ai

import (
	"encoding/json"
	"time"

	"rp-go/engine/ecs"
)

// ScriptStep represents one step in a scripted AI sequence.
type ScriptStep struct {
	Action string             `json:"action"`
	Params map[string]any     `json:"params"`
	Delay  int                `json:"delay_ms,omitempty"` // optional pause before next step
}

// ScriptState stores runtime progress of a script.
type ScriptState struct {
	Current int
	NextAt  time.Time
}

// behaviorScript executes a JSON-defined action sequence.
// Example in ai.json:
// {
//   "name": "patrol_then_retreat",
//   "type": "script",
//   "params": {
//     "steps": [
//       {"action": "patrol", "params": {"speed": 2.2}, "delay_ms": 0},
//       {"action": "retreat", "params": {"target": "player"}, "delay_ms": 500}
//     ]
//   }
// }
func (s *System) behaviorScript(w *ecs.World, e *ecs.Entity, pos *ecs.Position, vel *ecs.Velocity, params map[string]any) bool {
	rawSteps, ok := params["steps"]
	if !ok {
		return false
	}

	var steps []ScriptStep
	switch v := rawSteps.(type) {
	case []any:
		b, _ := json.Marshal(v)
		_ = json.Unmarshal(b, &steps)
	case []ScriptStep:
		steps = v
	default:
		return false
	}

	state, _ := e.Get("AIScriptState").(*ScriptState)
	if state == nil {
		state = &ScriptState{}
		e.AddNamed("AIScriptState", state)
	}

	if time.Now().Before(state.NextAt) {
		return true // waiting
	}
	if state.Current >= len(steps) {
		state.Current = 0
	}

	step := steps[state.Current]
	act := AIActionInstance{
		Name:     step.Action,
		Type:     step.Action,
		Priority: 0,
		Params:   step.Params,
	}
	executed := s.executeAction(w, e, pos, vel, act)
	if executed {
		if step.Delay > 0 {
			state.NextAt = time.Now().Add(time.Duration(step.Delay) * time.Millisecond)
		}
		state.Current++
	}
	return executed
}

