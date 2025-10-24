package planet

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"rp-go/engine/ecs"
	"rp-go/engine/events"
	"rp-go/engine/gfx"
	"rp-go/engine/platform"
	"rp-go/engine/world"
)

/*───────────────────────────────────────────────*
 | PLANET SCENE                                  |
 *───────────────────────────────────────────────*/

type Scene struct {
	init       bool
	ctx        *world.WorldContext
	landTimer  float64
	player     *ecs.Entity
}

/*───────────────────────────────────────────────*
 | CORE                                           |
 *───────────────────────────────────────────────*/

func (s *Scene) Name() string { return "planet" }

func (s *Scene) Init(w *ecs.World) {
	if s.init {
		return
	}
	s.init = true
	fmt.Println("[SCENE] Init:", s.Name())

	// ------------------------------------------------------------
	// World context (data, AI, assets)
	// ------------------------------------------------------------
	s.ctx = world.InitWorld(w)
	if s.ctx == nil || s.ctx.Creator == nil {
		fmt.Println("[PLANET] Missing world context")
		return
	}

	gfx.PreloadImages(
		"assets/entities/ship.png",
		"assets/entities/building.png",
		"assets/entities/lander.png",
	)

	// ------------------------------------------------------------
	// Player Ship (landing)
	// ------------------------------------------------------------
	s.player = w.NewEntity()
	s.player.Add(&ecs.Actor{
		ID:         "player",
		Archetype:  "ship",
		Persistent: true,
	})
	s.player.Add(&ecs.Position{X: 320, Y: -200}) // start offscreen
	s.player.Add(&ecs.Velocity{})
	s.player.Add(&ecs.Sprite{
		Image:        gfx.LoadImage("assets/entities/lander.png"),
		Width:        64,
		Height:       64,
		PixelPerfect: true,
	})

	// ------------------------------------------------------------
	// Simple landing area (background)
	// ------------------------------------------------------------
	ground := w.NewEntity()
	ground.Add(&ecs.Position{X: 0, Y: 300})
	ground.Add(&ecs.Sprite{
		Image:        gfx.LoadImage("assets/entities/building.png"),
		Width:        128,
		Height:       128,
		PixelPerfect: true,
	})
	fmt.Printf("[PLANET] Ground object created (entity %d)\n", ground.ID)

	// ------------------------------------------------------------
	// Camera follows player
	// ------------------------------------------------------------
	cam := w.NewEntity()
	camComp := &ecs.Camera{
		X:      320,
		Y:      240,
		Scale:  1.8,
		Target: s.player,
	}
	cam.Add(camComp)

	fmt.Println("[PLANET] Landing sequence starting")
}

/*───────────────────────────────────────────────*
 | UPDATE                                         |
 *───────────────────────────────────────────────*/

func (s *Scene) Update(w *ecs.World) {
	if s.player == nil {
		return
	}

	// Smooth landing animation
	pos, _ := s.player.Get("Position").(*ecs.Position)
	if pos.Y < 240 {
		s.landTimer += 0.03
		pos.Y = float64(240 - int(200*math.Cos(s.landTimer)))
		if pos.Y >= 240 {
			pos.Y = 240
			fmt.Println("[PLANET] Landing complete.")
		}
	}

	// Player can press ENTER to return to space
	if inpututil.IsKeyJustPressed(platform.KeyEnter) {
		if bus, ok := w.EventBus.(*events.TypedBus); ok && bus != nil {
			events.Publish(bus, events.SceneChangeEvent{
				Target: "space",
				Scene:  &space.Scene{},
			})
		}
	}
}

/*───────────────────────────────────────────────*
 | DRAW                                           |
 *───────────────────────────────────────────────*/

func (s *Scene) Draw(w *ecs.World, screen *platform.Image) {
	screen.Fill(color.RGBA{R: 20, G: 10, B: 30, A: 255})

	// Optional: text overlay
	platform.DrawText(screen, "Planet Surface", platform.DefaultFont(), 20, 32, color.White)
	platform.DrawText(screen, "Press ENTER to return to space", platform.DefaultFont(), 20, 56, color.RGBA{200, 200, 220, 255})
}

/*───────────────────────────────────────────────*
 | UNLOAD                                         |
 *───────────────────────────────────────────────*/

func (s *Scene) Unload(w *ecs.World) {
	fmt.Println("[SCENE] Unload:", s.Name())
	s.ctx = nil
}

