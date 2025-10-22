package planet

import (
	"rp-go/engine/ecs"
	"rp-go/engine/gfx"
	"rp-go/engine/platform"
)

type Scene struct{ initialized bool }

func (s *Scene) Name() string { return "planet" }

func (s *Scene) Init(w *ecs.World) {
	if s.initialized {
		return
	}
	s.initialized = true

	var player *ecs.Entity
	for _, e := range w.Entities {
		if a, ok := e.Get("Actor").(*ecs.Actor); ok && a.ID == "player" {
			player = e
			break
		}
	}

	if player == nil {
		player = w.NewEntity()
		player.Add(&ecs.Actor{
			ID:         "player",
			Archetype:  "character",
			Persistent: true,
		})
		player.Add(&ecs.Position{X: 100, Y: 100})
		player.Add(&ecs.Velocity{})
		player.Add(&ecs.CameraTarget{})
	}

	playerImg := gfx.LoadImage("assets/entities/player.png")
	player.Add(&ecs.Sprite{Image: playerImg, Width: 32, Height: 32})

	cam := w.NewEntity()
	cam.Add(&ecs.Camera{X: 0, Y: 0, Scale: 2.0, Target: player})
}

func (s *Scene) Update(w *ecs.World)                       {}
func (s *Scene) Draw(w *ecs.World, screen *platform.Image) {}
func (s *Scene) Unload(w *ecs.World)                       {}
