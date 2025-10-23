package world

import "rp-go/engine/data"
import "rp-go/engine/ecs"

func buildAIController(tpl *data.ActorAITemplate) *ecs.AIController {
	if tpl == nil {
		return nil
	}

	ctrl := &ecs.AIController{Speed: tpl.Speed}
	var active bool

	if f := tpl.Follow; f != nil {
		ctrl.Follow = &ecs.AIFollowBehavior{
			Target:      f.Target,
			OffsetX:     f.OffsetX,
			OffsetY:     f.OffsetY,
			MinDistance: f.MinDistance,
			MaxDistance: f.MaxDistance,
			Speed:       f.Speed,
		}
		active = true
	}
	if p := tpl.Pursue; p != nil {
		ctrl.Pursue = &ecs.AIPursueBehavior{
			Target:         p.Target,
			EngageDistance: p.EngageDistance,
			Speed:          p.Speed,
		}
		active = true
	}
	if pa := tpl.Patrol; pa != nil {
		ctrl.Patrol = &ecs.AIPathBehavior{
			Variant:   pa.Variant,
			Waypoints: convertAIWaypoints(pa.Waypoints),
			Speed:     pa.Speed,
		}
		ctrl.PatrolState.Reset()
		active = true
	}
	if r := tpl.Retreat; r != nil {
		ctrl.Retreat = &ecs.AIRetreatBehavior{
			Target:          r.Target,
			TriggerDistance: r.TriggerDistance,
			SafeDistance:    r.SafeDistance,
			Speed:           r.Speed,
		}
		active = true
	}
	if t := tpl.Travel; t != nil {
		ctrl.Travel = &ecs.AIPathBehavior{
			Variant:   t.Variant,
			Waypoints: convertAIWaypoints(t.Waypoints),
			Speed:     t.Speed,
		}
		ctrl.TravelState.Reset()
		active = true
	}

	if !active && ctrl.Speed <= 0 {
		return nil
	}
	return ctrl
}

func convertAIWaypoints(src []data.ActorAIWaypoint) []ecs.AIWaypoint {
	if len(src) == 0 {
		return nil
	}
	out := make([]ecs.AIWaypoint, len(src))
	for i, pt := range src {
		out[i] = ecs.AIWaypoint{X: pt.X, Y: pt.Y}
	}
	return out
}
// BuildAIController is the public entry point for constructing an ECS AIController
// from a data.ActorAITemplate. It wraps the internal buildAIController.
func BuildAIController(tpl *data.ActorAITemplate) *ecs.AIController {
	return buildAIController(tpl)
}

