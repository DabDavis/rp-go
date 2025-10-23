package core

import (
	"fmt"
	"testing"

	"rp-go/engine/ecs"
)

func TestNewGameWorldRegistersSystemPasses(t *testing.T) {
	world := NewGameWorld().World

	simulationTypes := map[string]struct{}{
		"*scene.Manager":   {},
		"*input.System":    {},
		"*ai.System":       {},
		"*movement.System": {},
		"*camera.System":   {},
	}
	renderingTypes := map[string]struct{}{
		"*background.System": {},
		"*render.System":     {},
		"*debug.System":      {},
	}

	var seenRendering bool
	for idx, sys := range world.Systems {
		typeName := fmt.Sprintf("%T", sys)
		if _, ok := renderingTypes[typeName]; ok {
			seenRendering = true
			continue
		}

		if _, ok := simulationTypes[typeName]; ok {
			if seenRendering {
				t.Fatalf("simulation system %s registered after rendering system at index %d", typeName, idx)
			}
			continue
		}

		if _, ok := sys.(ecs.DrawableSystem); ok {
			t.Fatalf("unexpected drawable system registered in game world: %s", typeName)
		}
	}

	for typeName := range simulationTypes {
		if !containsSystem(world.Systems, typeName) {
			t.Errorf("missing simulation system %s", typeName)
		}
	}
	for typeName := range renderingTypes {
		if !containsSystem(world.Systems, typeName) {
			t.Errorf("missing rendering system %s", typeName)
		}
	}
}

func containsSystem(systems []ecs.System, typeName string) bool {
	for _, sys := range systems {
		if fmt.Sprintf("%T", sys) == typeName {
			return true
		}
	}
	return false
}
