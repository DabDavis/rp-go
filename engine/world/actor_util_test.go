package world

import (
	"testing"

	"rp-go/engine/data"
)

func TestNextID(t *testing.T) {
	ac := &ActorCreator{counters: make(map[string]int)}
	id1 := ac.nextID("drone")
	id2 := ac.nextID("drone")

	if id1 == id2 {
		t.Fatalf("expected unique IDs, got %q and %q", id1, id2)
	}
}

func TestTemplateNames(t *testing.T) {
	ac := &ActorCreator{
		templates: map[string]data.ActorTemplate{
			"enemy": {
				Name: "enemy",
				Sprite: data.ActorSpriteTemplate{
					Image: "enemy.png",
				},
			},
		},
	}
	names := ac.Templates()
	if len(names) != 1 || names[0] != "enemy" {
		t.Fatalf("unexpected template names: %v", names)
	}
}

