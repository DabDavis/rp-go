package space

import (
	"fmt"

	"rp-go/engine/ecs"
	"rp-go/engine/gfx"
	"rp-go/engine/platform"
	"rp-go/engine/world"
)

/*───────────────────────────────────────────────*
 | SCENE STRUCTURE                               |
 *───────────────────────────────────────────────*/

type Scene struct {
	init     bool
	ctx      *world.WorldContext
}

/*───────────────────────────────────────────────*
 | CORE                                           |
 *───────────────────────────────────────────────*/

func (s *Scene) Name() string { return "space" }

func (s *Scene) Init(w *ecs.World) {
	if s.init {
		return
	}
	s.init = true
	fmt.Println("[SCENE] Init:", s.Name())

	// ---------------------------------------------------------------------
	// Initialize world context (Data + AI + Creator)
	// ---------------------------------------------------------------------
	s.ctx = world.InitWorld(w)
	if s.ctx == nil || s.ctx.Creator == nil {
		fmt.Println("[SCENE] Missing world context — cannot spawn entities")
		return
	}

	// ---------------------------------------------------------------------
	// Preload additional shared assets
	// ---------------------------------------------------------------------
	gfx.PreloadImages("assets/entities/ship.png", "assets/entities/planet.png")

	// ---------------------------------------------------------------------
	// Player Entity
	// ---------------------------------------------------------------------
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

	player.Add(&ecs.Sprite{
		Image:        gfx.LoadImage("assets/entities/ship.png"),
		Width:        64,
		Height:       64,
		PixelPerfect: true,
	})
	fmt.Printf("[SCENE] Player spawned (entity %d)\n", player.ID)

	// ---------------------------------------------------------------------
	// Camera Entity
	// ---------------------------------------------------------------------
	cam := w.NewEntity()
	camComp := &ecs.Camera{
		X:      100,
		Y:      100,
		Scale:  1.5,
		Target: player,
	}
	cam.Add(camComp)
	fmt.Printf("[SCENE] Camera created (entity %d)\n", cam.ID)

	// ---------------------------------------------------------------------
	// Planet Entity
	// ---------------------------------------------------------------------
	planet := w.NewEntity()
	planet.Add(&ecs.Position{X: 350, Y: 180})
	planet.Add(&ecs.Sprite{
		Image:        gfx.LoadImage("assets/entities/planet.png"),
		Width:        128,
		Height:       128,
		PixelPerfect: true,
	})
	fmt.Printf("[SCENE] Planet created (entity %d)\n", planet.ID)

	// ---------------------------------------------------------------------
	// Enemy Fleet (data-driven)
	// ---------------------------------------------------------------------
	templates := []string{
		"dark-elf-ship-scout",
		"dark-elf-ship-vanguard",
		"dark-elf-ship-raider",
		"dark-elf-ship-evader",
		"dark-elf-ship-commander",
	}

	for i, name := range templates {
		x := 260 + float64(i*72)
		y := 220.0
		enemy, err := s.ctx.Creator.Spawn(w, name, ecs.Position{X: x, Y: y})
		if err != nil {
			fmt.Printf("[SCENE] Spawn failed for %s: %v\n", name, err)
			continue
		}

		if meta, ok := enemy.Get("Actor").(*ecs.Actor); ok {
			fmt.Printf("[SCENE] Spawned %s (%s) entity %d\n", meta.ID, name, enemy.ID)
		}

		if s.ctx.AIComposer != nil {
			s.ctx.AIComposer.BindEntityAI(w, enemy)
		}
	}

	fmt.Printf("[SCENE] Ready: %s\n", s.Name())
}

/*───────────────────────────────────────────────*
 | FRAME EVENTS                                  |
 *───────────────────────────────────────────────*/

func (s *Scene) Update(w *ecs.World) {}

func (s *Scene) Draw(w *ecs.World, screen *platform.Image) {}

func (s *Scene) Unload(w *ecs.World) {
	fmt.Println("[SCENE] Unload:", s.Name())
	s.ctx = nil
}

