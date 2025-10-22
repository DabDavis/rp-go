package core

import (
	"testing"

	"rp-go/engine/ecs"
	"rp-go/engine/platform"
)

func TestHeadlessWorldSpawnsCamera(t *testing.T) {
	world := NewGameWorld()
	screen := platform.NewImage(world.Config.Viewport.Width, world.Config.Viewport.Height)

	for i := 0; i < 5; i++ {
		world.Update()
		screen.Clear()
		world.Draw(screen)
	}

	foundCamera := false
	for _, e := range world.World.Entities {
		if _, ok := e.Get("Camera").(*ecs.Camera); ok {
			foundCamera = true
			break
		}
	}

	if !foundCamera {
		t.Fatalf("expected at least one camera entity after headless steps")
	}
}
