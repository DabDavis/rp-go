package core

import (
	"rp-go/engine/data"
	"rp-go/engine/ecs"
	"rp-go/engine/events"
	"rp-go/engine/platform"

	"rp-go/engine/scenes/space"
	"rp-go/engine/systems/camera"
	"rp-go/engine/systems/debug"
	"rp-go/engine/systems/input"
	"rp-go/engine/systems/movement"
	"rp-go/engine/systems/background"
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

	// Scene manager FIRST — it creates entities (ship, camera, planet)
	sm := &scene.Manager{}
	w.AddSystem(sm)

	// Core systems in logical update order
	w.AddSystem(&background.System{}) // 🌌 Draws parallax stars
	w.AddSystem(&input.System{})
	w.AddSystem(&movement.System{})
	w.AddSystem(camera.NewSystem(camera.Config{
		MinScale: cfg.Viewport.MinScale,
		MaxScale: cfg.Viewport.MaxScale,
		ZoomStep: cfg.Viewport.ZoomStep,
		ZoomLerp: cfg.Viewport.ZoomLerp,
	}))
	w.AddSystem(&render.System{}) // Draws world-space entities
	w.AddSystem(&debug.System{})  // Overlay (UI/debug info)

	// Start in the space scene
	sm.QueueScene(&space.Scene{})

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

