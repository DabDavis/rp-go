package main

import (
	"log"
	"rp-go/engine/core"
	"rp-go/engine/ecs"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	world     *core.GameWorld
	offscreen *ebiten.Image
}

func (g *Game) Update() error {
	g.world.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	cfg := g.world.Config
	w := g.world.World
	cam := getActiveCamera(w)

	if g.offscreen == nil {
		g.offscreen = ebiten.NewImage(cfg.Viewport.Width, cfg.Viewport.Height)
	}

	// ✅ Draw world into offscreen buffer (1:1 internal pixels)
	g.offscreen.Clear()
	w.Draw(g.offscreen)

	if cam == nil {
		screen.DrawImage(g.offscreen, nil)
		return
	}

	// ✅ Composite offscreen to window, applying zoom & rotation
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterNearest

	// Apply camera scale and rotation
	op.GeoM.Scale(cam.Scale, cam.Scale)
	op.GeoM.Rotate(cam.Rotation) // rotation placeholder (0 by default)

	// Center on screen
	windowW := float64(cfg.Window.Width)
	windowH := float64(cfg.Window.Height)
	offW := float64(cfg.Viewport.Width)
	offH := float64(cfg.Viewport.Height)
	op.GeoM.Translate(
		windowW/2-offW*cam.Scale/2,
		windowH/2-offH*cam.Scale/2,
	)

	screen.DrawImage(g.offscreen, op)
}

func (g *Game) Layout(outW, outH int) (int, int) {
	cfg := g.world.Config
	return cfg.Window.Width, cfg.Window.Height
}

func main() {
	gameWorld := core.NewGameWorld()
	cfg := gameWorld.Config

	game := &Game{world: gameWorld}

	ebiten.SetWindowSize(cfg.Window.Width, cfg.Window.Height)
	ebiten.SetWindowTitle("rp-go: ECS Camera Prototype")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// Utility: get first camera entity
func getActiveCamera(w *ecs.World) *ecs.Camera {
	for _, e := range w.Entities {
		if c, ok := e.Get("Camera").(*ecs.Camera); ok {
			return c
		}
	}
	return nil
}

