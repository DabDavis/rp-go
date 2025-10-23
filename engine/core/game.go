package core

import (
	"rp-go/engine/data"
	"rp-go/engine/ecs"
	"rp-go/engine/events"
	"rp-go/engine/platform"

	"rp-go/engine/scenes/space"

	"rp-go/engine/systems/actor"
	"rp-go/engine/systems/ai"
	"rp-go/engine/systems/aicomposer"
	"rp-go/engine/systems/background"
	"rp-go/engine/systems/camera"
	dataSys "rp-go/engine/systems/data" // renamed to avoid collision
	"rp-go/engine/systems/debug"
	"rp-go/engine/systems/devconsole"
	"rp-go/engine/systems/entitylist"
	"rp-go/engine/systems/hud"
	"rp-go/engine/systems/input"
	"rp-go/engine/systems/movement"
	"rp-go/engine/systems/render"
	"rp-go/engine/systems/scene"
	"rp-go/engine/systems/windowmgr"
)

/*───────────────────────────────────────────────*
 | GAME WORLD                                    |
 *───────────────────────────────────────────────*/

type GameWorld struct {
	World  *ecs.World
	Config data.RenderConfig
}

/*───────────────────────────────────────────────*
 | INITIALIZATION                                |
 *───────────────────────────────────────────────*/

func NewGameWorld() *GameWorld {
	// -------------------------------------------------------------------------
	// ECS + EventBus Setup
	// -------------------------------------------------------------------------
	w := ecs.NewWorld()
	w.EventBus = events.NewBus()

	// -------------------------------------------------------------------------
	// Data System (config, actor db, ai.json hot reload)
	// -------------------------------------------------------------------------
	dataSystem := dataSys.NewSystem()
	w.AddSystem(dataSystem)

	cfg := dataSystem.Config
	if cfg.Window.Width == 0 {
		cfg = data.LoadRenderConfig("engine/data/render_config.json")
	}

	// -------------------------------------------------------------------------
	// Core Managers and Systems
	// -------------------------------------------------------------------------
	sceneManager := &scene.Manager{}
	actorSystem := actor.NewSystem()

	// --- AI Layer ------------------------------------------------------------
	aiSystem := ai.NewSystem(dataSystem.AICatalog)
	aiSystem.SetActorLookup(actorSystem.Registry())

	composerSystem := aicomposer.NewSystem(dataSystem, aiSystem)

	// --- Developer + Debug Layer --------------------------------------------
	consoleSystem := devconsole.NewSystem(actorSystem.Registry(), devconsole.Config{
		Margin:         16,
		ViewportWidth:  cfg.Viewport.Width,
		ViewportHeight: cfg.Viewport.Height,
	})

	debugSystem := debug.NewSystem(debug.Config{
		Margin:         16,
		ViewportWidth:  cfg.Viewport.Width,
		ViewportHeight: cfg.Viewport.Height,
	})
	debugSystem.AttachComposer(composerSystem) // ✅ hook AIComposer debug window

	entityListSystem := entitylist.NewSystem(actorSystem.Registry())

	// -------------------------------------------------------------------------
	// Simulation Phase — world state and logic
	// -------------------------------------------------------------------------
	simulationSystems := []ecs.System{
		dataSystem,         // config hot-reload + JSON catalogs
		sceneManager,       // scene transitions, loading/unloading
		actorSystem,        // actor registration
		composerSystem,     // auto-binds AIControllers from refs
		&input.System{},    // player + input control
		aiSystem,           // AI decision-making & movement
		&movement.System{}, // position/velocity propagation
		camera.NewSystem(camera.Config{
			MinScale: cfg.Viewport.MinScale,
			MaxScale: cfg.Viewport.MaxScale,
			ZoomStep: cfg.Viewport.ZoomStep,
			ZoomLerp: cfg.Viewport.ZoomLerp,
		}),
	}

	// -------------------------------------------------------------------------
	// Rendering Phase — visuals & overlays
	// -------------------------------------------------------------------------
	hudSystem := hud.NewSystem()
	windowSystem := windowmgr.NewSystem()

	windowRenderers := []ecs.System{
		render.NewWindowRenderer(ecs.LayerHUD),
		render.NewWindowRenderer(ecs.LayerDebug),
		render.NewWindowRenderer(ecs.LayerConsole),
	}

	renderingSystems := []ecs.System{
		&background.System{}, // parallax stars
		&render.System{},     // world-space drawables
		hudSystem,            // reusable HUD content
		windowSystem,         // modular window overlays
	}
	renderingSystems = append(renderingSystems, windowRenderers...)
	renderingSystems = append(renderingSystems,
		entityListSystem, // entity info overlay
		debugSystem,      // debug metrics, including composer window
		consoleSystem,    // developer console overlay
	)

	// -------------------------------------------------------------------------
	// System Registration
	// -------------------------------------------------------------------------
	for _, sys := range simulationSystems {
		w.AddSystem(sys)
	}
	for _, sys := range renderingSystems {
		w.AddSystem(sys)
	}

	// -------------------------------------------------------------------------
	// Initial Scene Setup
	// -------------------------------------------------------------------------
	sceneManager.QueueScene(&space.Scene{})

	// -------------------------------------------------------------------------
	// Return Assembled World
	// -------------------------------------------------------------------------
	return &GameWorld{
		World:  w,
		Config: cfg,
	}
}

/*───────────────────────────────────────────────*
 | MAIN LOOP                                     |
 *───────────────────────────────────────────────*/

// Update advances the world simulation by one frame and flushes events.
func (g *GameWorld) Update() {
	g.World.Update()
	if bus, ok := g.World.EventBus.(*events.TypedBus); ok && bus != nil {
		bus.Flush()
	}
}

// Draw executes all world-space renderers.
// Overlays (HUD, console, debug) are drawn later in main.go.
func (g *GameWorld) Draw(screen *platform.Image) {
	g.World.DrawWorld(screen)
}

