package space

import (
	"fmt"

	"rp-go/engine/ecs"
	"rp-go/engine/events"
	"rp-go/engine/gfx"
	"rp-go/engine/platform"
	"rp-go/engine/systems/aicomposer"
	dataSys "rp-go/engine/systems/data"
	"rp-go/engine/world"
)

type Scene struct {
	initialized   bool
	actorCreator  *world.ActorCreator
	dataSystem    *dataSys.System
	composer      *aicomposer.System
}

// Name returns the scene identifier.
func (s *Scene) Name() string { return "space" }

/*───────────────────────────────────────────────*
 | INITIALIZATION                                |
 *───────────────────────────────────────────────*/

func (s *Scene) Init(w *ecs.World) {
	if s.initialized {
		return
	}
	s.initialized = true
	fmt.Println("[SCENE] Initializing: Space")

	// -------------------------------------------------------------------------
	// Resolve core systems (data + AIComposer)
	// -------------------------------------------------------------------------
	if s.dataSystem == nil {
		if sys := w.FindSystem((*dataSys.System)(nil)); sys != nil {
			s.dataSystem = sys.(*dataSys.System)
		}
	}
	if s.composer == nil {
		if sys := w.FindSystem((*aicomposer.System)(nil)); sys != nil {
			s.composer = sys.(*aicomposer.System)
		}
	}

	// -------------------------------------------------------------------------
	// Prepare ActorCreator using the centralized actor DB
	// -------------------------------------------------------------------------
	db := s.dataSystem.Actors
	if len(db.Actors) == 0 {
		db = s.dataSystem.Actors // ensures fallback even if uninitialized
	}
	s.actorCreator = world.NewActorCreator(db)
	s.actorCreator.PreloadImages()

	// -------------------------------------------------------------------------
	// Warm Texture Cache
	// -------------------------------------------------------------------------
	gfx.PreloadImages(
		"assets/entities/ship.png",
		"assets/entities/planet.png",
	)

	// -------------------------------------------------------------------------
	// === Player Ship ===
	// -------------------------------------------------------------------------
	player := w.NewEntity()
	player.Add(&ecs.Actor{
		ID:         "player",
		Archetype:  "ship",
		Persistent: true,
	})
	player.Add(&ecs.Position{X: 100, Y: 100})
	player.Add(&ecs.Velocity{})
	player.Add(&ecs.PlayerInput{Enabled: true})
	player.Add(&ecs.CameraTarget{})

	playerImg := gfx.LoadImage("assets/entities/ship.png")
	player.Add(&ecs.Sprite{
		Image:        playerImg,
		Width:        64,
		Height:       64,
		PixelPerfect: true,
	})

	// -------------------------------------------------------------------------
	// === Camera Entity ===
	// -------------------------------------------------------------------------
	cam := w.NewEntity()
	camComp := &ecs.Camera{
		X:      100,
		Y:      100,
		Scale:  1.5,
		Target: player,
	}
	cam.Add(camComp)
	fmt.Printf("[SCENE] Camera entity created (ID %d) -> %+v\n", cam.ID, *camComp)

	// -------------------------------------------------------------------------
	// === Planet Entity ===
	// -------------------------------------------------------------------------
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

	// -------------------------------------------------------------------------
	// === Dark Elf Patrol ===
	// -------------------------------------------------------------------------
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
				fmt.Printf("[SCENE] Spawned %s (%s) [entity %d]\n", meta.ID, template, enemy.ID)
			}
		}

		// Notify the AIComposer to attach behavior
		if s.composer != nil {
			s.composer.BindEntityAI(w, enemy)
		}
	}

	// -------------------------------------------------------------------------
	// Scene ready
	// -------------------------------------------------------------------------
	fmt.Println("[SCENE] Active:", s.Name())
}

/*───────────────────────────────────────────────*
 | FRAME EVENTS                                  |
 *───────────────────────────────────────────────*/

func (s *Scene) Update(w *ecs.World) {}

func (s *Scene) Draw(w *ecs.World, screen *platform.Image) {}

func (s *Scene) Unload(w *ecs.World) {
	fmt.Println("[SCENE] Unloading:", s.Name())
}

