package core

import (
	"rp-go/engine/data"
	"rp-go/engine/ecs"
	"rp-go/engine/events"
	"rp-go/engine/platform"

	"rp-go/engine/scenes/space"
	"rp-go/engine/systems/ai"
	"rp-go/engine/systems/background"
	"rp-go/engine/systems/camera"
	"rp-go/engine/systems/debug"
	"rp-go/engine/systems/input"
	"rp-go/engine/systems/movement"
	"rp-go/engine/systems/render"
	"rp-go/engine/systems/scene"
)

// GameWorld bundles the ECS world and runtime configuration.
type GameWorld struct {
	World  *ecs.World
	Config data.RenderConfig
}

// NewGameWorld creates and initializes a full engine world.
func NewGameWorld() *GameWorld {
	cfg := data.LoadRenderConfig("engine/data/render_config.json")
	w := ecs.NewWorld()

	// Wire up the typed event bus so systems can coordinate without direct dependencies.
	w.EventBus = events.NewBus()

	// Systems are registered in two passes: simulation first, then rendering.
	// This keeps the update loop free of draw calls so the world can advance
	// headlessly inside tests or server-side simulations.
	sceneManager := &scene.Manager{}

	simulationSystems := []ecs.System{
		sceneManager,
		&input.System{},
		ai.NewSystem(),
		&movement.System{},
		camera.NewSystem(camera.Config{
			MinScale: cfg.Viewport.MinScale,
			MaxScale: cfg.Viewport.MaxScale,
			ZoomStep: cfg.Viewport.ZoomStep,
			ZoomLerp: cfg.Viewport.ZoomLerp,
		}),
	}

	renderingSystems := []ecs.System{
		&background.System{}, // 🌌 Parallax stars
		&render.System{},     // World-space sprites
		&debug.System{},      // Overlay diagnostics
	}

	for _, sys := range append(simulationSystems, renderingSystems...) {
		w.AddSystem(sys)
	}

	// Start in the space scene
	sceneManager.QueueScene(&space.Scene{})

	return &GameWorld{World: w, Config: cfg}
}

// Update runs one ECS update tick and flushes the event bus.
func (g *GameWorld) Update() {
	g.World.Update()
	if bus, ok := g.World.EventBus.(*events.TypedBus); ok && bus != nil {
		bus.Flush()
	}
}

// Draw executes only world-space (camera-affected) rendering systems.
// Overlay systems (debug, HUD) are drawn in main.go after compositing.
func (g *GameWorld) Draw(screen *platform.Image) {
	g.World.DrawWorld(screen)
}
