package space

import (
	"fmt"
	"image/color"

	"rp-go/engine/ecs"
	"rp-go/engine/gfx"
	"rp-go/engine/platform"
)

type Scene struct{ initialized bool }

func (s *Scene) Name() string { return "space" }

func (s *Scene) Init(w *ecs.World) {
	if s.initialized {
		return
	}
	s.initialized = true

	fmt.Println("[SCENE] Initializing: Space")

	// Warm the texture cache so the first frame doesn't block on disk IO.
	gfx.PreloadImages(
		"assets/entities/ship.png",
		"assets/entities/planet.png",
	)

	// === Player Ship ===
	ship := w.NewEntity()
	ship.Add(&ecs.Actor{
		ID:         "player",
		Archetype:  "ship",
		Persistent: true,
	})
	ship.Add(&ecs.Position{X: 100, Y: 100})
	ship.Add(&ecs.Velocity{})
	ship.Add(&ecs.CameraTarget{})

	shipImg := gfx.LoadImage("assets/entities/ship.png")
	ship.Add(&ecs.Sprite{Image: shipImg, Width: 64, Height: 64})

	// === Camera ===
	cam := w.NewEntity()
	camComp := &ecs.Camera{X: 100, Y: 100, Scale: 1.5, Target: ship}
	cam.Add(camComp)
	fmt.Printf("[SCENE] Camera entity created (ID %d) -> %+v\n", cam.ID, *camComp)

	// === Planet ===
	planet := w.NewEntity()
	planet.Add(&ecs.Position{X: 350, Y: 180})
	planetImg := gfx.LoadImage("assets/entities/planet.png")
	planet.Add(&ecs.Sprite{Image: planetImg, Width: 128, Height: 128})
	fmt.Printf("[SCENE] Planet entity created (ID %d)\n", planet.ID)
}

func (s *Scene) Update(w *ecs.World) {}

func (s *Scene) Draw(w *ecs.World, screen *platform.Image) {
	// Proper background fill
	screen.Fill(color.RGBA{0, 0, 32, 255})
}

func (s *Scene) Unload(w *ecs.World) {
	fmt.Println("[SCENE] Unloading: Space")
}
