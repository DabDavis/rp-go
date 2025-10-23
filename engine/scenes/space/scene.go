package space

import (
	"fmt"

	"rp-go/engine/data"
	"rp-go/engine/ecs"
	"rp-go/engine/gfx"
	"rp-go/engine/platform"
	"rp-go/engine/world"
)

type Scene struct {
	initialized  bool
	actorCreator *world.ActorCreator
}

func (s *Scene) Name() string { return "space" }

func (s *Scene) Init(w *ecs.World) {
	if s.initialized {
		return
	}
	s.initialized = true

	fmt.Println("[SCENE] Initializing: Space")

	actorDB := data.LoadActorDatabase("engine/data/actors.json")
	s.actorCreator = world.NewActorCreator(actorDB)
	s.actorCreator.PreloadImages()

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
	ship.Add(&ecs.PlayerInput{Enabled: true})
	ship.Add(&ecs.CameraTarget{})

	shipImg := gfx.LoadImage("assets/entities/ship.png")
	ship.Add(&ecs.Sprite{
		Image:        shipImg,
		Width:        64,
		Height:       64,
		PixelPerfect: true,
	})

	// === Camera ===
	cam := w.NewEntity()
	camComp := &ecs.Camera{
		X:      100,
		Y:      100,
		Scale:  1.5,
		Target: ship,
	}
	cam.Add(camComp)
	fmt.Printf("[SCENE] Camera entity created (ID %d) -> %+v\n", cam.ID, *camComp)

	// === Planet ===
	planet := w.NewEntity()
	planet.Add(&ecs.Position{X: 350, Y: 180})
	planetImg := gfx.LoadImage("assets/entities/planet.png")
	planet.Add(&ecs.Sprite{
		Image:        planetImg,
		Width:        128,
		Height:       128,
		PixelPerfect: true,
	})
	fmt.Printf("[SCENE] Planet entity created (ID %d)\n", planet.ID)

	// === Dark Elf Patrol ===
	templates := []string{
		"dark-elf-ship-scout",
		"dark-elf-ship-vanguard",
		"dark-elf-ship-raider",
		"dark-elf-ship-evader",
		"dark-elf-ship-commander",
	}

	for i, template := range templates {
		spawnX := 260 + float64(i*72)
		spawnY := 220.0
		enemy, err := s.actorCreator.Spawn(w, template, ecs.Position{X: spawnX, Y: spawnY})
		if err != nil {
			fmt.Printf("[SCENE] Failed to spawn patrol ship: %v\n", err)
			continue
		}
		if actor := enemy.Get("Actor"); actor != nil {
			if meta, ok := actor.(*ecs.Actor); ok {
				fmt.Printf("[SCENE] Spawned %s (%s) with AI (entity %d)\n", meta.ID, template, enemy.ID)
			}
		}
	}
}

func (s *Scene) Update(w *ecs.World) {}

func (s *Scene) Draw(w *ecs.World, screen *platform.Image) {
}

func (s *Scene) Unload(w *ecs.World) {
	fmt.Println("[SCENE] Unloading: Space")
}
