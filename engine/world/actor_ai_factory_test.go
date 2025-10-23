package world

import (
	"reflect"
	"testing"

	"rp-go/engine/data"
	"rp-go/engine/ecs"
)

func TestBuildAIController_EmptyTemplateReturnsNil(t *testing.T) {
	if ctrl := BuildAIController(nil); ctrl != nil {
		t.Fatalf("expected nil for nil template, got %#v", ctrl)
	}

	if ctrl := BuildAIController(&data.ActorAITemplate{}); ctrl != nil {
		t.Fatalf("expected nil for empty AI template, got %#v", ctrl)
	}
}

func TestBuildAIController_FullTemplate(t *testing.T) {
	tpl := &data.ActorAITemplate{
		Speed: 2.5,
		Follow: &data.ActorAIFollow{
			Target:      "player",
			OffsetX:     10,
			OffsetY:     -5,
			MinDistance: 5,
			MaxDistance: 20,
			Speed:       3,
		},
		Pursue: &data.ActorAIPursue{
			Target:         "enemy",
			EngageDistance: 100,
			Speed:          4,
		},
		Patrol: &data.ActorAIPatrol{
			Variant: "loop",
			Speed:   2,
			Waypoints: []data.ActorAIWaypoint{
				{X: 10, Y: 10}, {X: 20, Y: 20},
			},
		},
		Retreat: &data.ActorAIRetreat{
			Target:          "player",
			TriggerDistance: 50,
			SafeDistance:    100,
			Speed:           3,
		},
		Travel: &data.ActorAITravel{
			Variant: "once",
			Speed:   5,
			Waypoints: []data.ActorAIWaypoint{
				{X: 0, Y: 0}, {X: 100, Y: 100},
			},
		},
	}

	ctrl := BuildAIController(tpl)
	if ctrl == nil {
		t.Fatalf("expected non-nil AI controller")
	}
	if ctrl.Speed != tpl.Speed {
		t.Errorf("expected speed %.2f, got %.2f", tpl.Speed, ctrl.Speed)
	}
	if ctrl.Follow == nil || ctrl.Pursue == nil || ctrl.Patrol == nil || ctrl.Retreat == nil || ctrl.Travel == nil {
		t.Errorf("expected all behaviors populated, got %+v", ctrl)
	}
	if len(ctrl.Patrol.Waypoints) != 2 {
		t.Errorf("expected 2 patrol waypoints, got %d", len(ctrl.Patrol.Waypoints))
	}
	if ctrl.PatrolState.Index != 0 {
		t.Errorf("expected initial patrol index 0, got %d", ctrl.PatrolState.Index)
	}
}

func TestConvertAIWaypoints(t *testing.T) {
	in := []data.ActorAIWaypoint{{X: 1, Y: 2}, {X: 3, Y: 4}}
	out := convertAIWaypoints(in)
	if len(out) != len(in) {
		t.Fatalf("expected %d waypoints, got %d", len(in), len(out))
	}
	for i := range in {
		if !reflect.DeepEqual(out[i], ecs.AIWaypoint{X: in[i].X, Y: in[i].Y}) {
			t.Errorf("mismatch at %d: %+v != %+v", i, in[i], out[i])
		}
	}
}

