package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"rp-go/engine/core"
	"rp-go/engine/ecs"
	"rp-go/engine/platform"
)

type Game struct {
	world     *core.GameWorld
	offscreen *platform.Image
}

func (g *Game) Update() error {
	g.world.Update()
	return nil
}

func (g *Game) Draw(screen *platform.Image) {
	cfg := g.world.Config
	w := g.world.World
	cam := getActiveCamera(w)

	// Create offscreen buffer for world rendering
	if g.offscreen == nil {
		g.offscreen = platform.NewImage(cfg.Viewport.Width, cfg.Viewport.Height)
	}

	/* ---------------------------------------------------------------------- */
	/*                              WORLD PASS                                */
	/* ---------------------------------------------------------------------- */
	g.offscreen.Clear()
	g.world.Draw(g.offscreen) // world.DrawWorld internally

	/* ---------------------------------------------------------------------- */
	/*                            COMPOSITE TO SCREEN                         */
	/* ---------------------------------------------------------------------- */
	op := platform.NewDrawImageOptions()
	op.SetFilter(platform.FilterNearest)

	if cam != nil {
		// Apply camera zoom
		op.Scale(cam.Scale, cam.Scale)

		// Center world on screen
		windowW := float64(cfg.Window.Width)
		windowH := float64(cfg.Window.Height)
		offW := float64(cfg.Viewport.Width)
		offH := float64(cfg.Viewport.Height)
		op.Translate(
			windowW/2-offW*cam.Scale/2,
			windowH/2-offH*cam.Scale/2,
		)
	}

	screen.DrawImage(g.offscreen, op)

	/* ---------------------------------------------------------------------- */
	/*                              OVERLAY PASS                              */
	/* ---------------------------------------------------------------------- */
	// Draw screen-space systems (HUD, debug) directly to window
	w.DrawOverlay(screen)
}

func (g *Game) Layout(outW, outH int) (int, int) {
	cfg := g.world.Config
	return cfg.Window.Width, cfg.Window.Height
}

func main() {
	gameWorld := core.NewGameWorld()
	cfg := gameWorld.Config
	game := &Game{world: gameWorld}

	headless := flag.Bool("headless", false, "run without opening a window")
	frames := flag.Int("frames", 120, "number of frames to run in headless mode")
	flag.Parse()

	// Allow environment variables to override flags
	if envHeadless := os.Getenv("RP_HEADLESS"); envHeadless != "" {
		if v, err := strconv.ParseBool(envHeadless); err == nil {
			*headless = v
		}
	}
	if envFrames := os.Getenv("RP_HEADLESS_FRAMES"); envFrames != "" {
		if v, err := strconv.Atoi(envFrames); err == nil {
			*frames = v
		}
	}

	// Headless simulation mode
	if *headless {
		if err := platform.RunHeadless(game, *frames, cfg.Viewport.Width, cfg.Viewport.Height); err != nil {
			log.Fatal(err)
		}
		log.Printf("Headless run complete (%d frames)\n", *frames)
		return
	}

	// Normal game mode
	platform.SetWindowSize(cfg.Window.Width, cfg.Window.Height)
	platform.SetWindowTitle("rp-go: ECS Camera Prototype")

	if err := platform.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// Utility: get first active camera in the world
func getActiveCamera(w *ecs.World) *ecs.Camera {
	for _, e := range w.Entities {
		if c, ok := e.Get("Camera").(*ecs.Camera); ok {
			return c
		}
	}
	return nil
}

