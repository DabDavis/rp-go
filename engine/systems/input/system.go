package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"rp-go/engine/ecs"
)

type System struct{}

func (s *System) Update(w *ecs.World) {
	for _, e := range w.Entities {
		v, ok := e.Get("Velocity").(*ecs.Velocity)
		if !ok {
			continue
		}
		v.VX, v.VY = 0, 0
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			v.VX = -2
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			v.VX = 2
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
			v.VY = -2
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			v.VY = 2
		}
	}
}

func (s *System) Draw(*ecs.World, *ebiten.Image) {}

