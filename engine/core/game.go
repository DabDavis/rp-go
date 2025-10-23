package core

import (
    "rp-go/engine/data"
    "rp-go/engine/ecs"
    "rp-go/engine/events"
    "rp-go/engine/platform"

    "rp-go/engine/scenes/space"
    "rp-go/engine/systems/actor"
    "rp-go/engine/systems/ai"
    "rp-go/engine/systems/background"
    "rp-go/engine/systems/camera"
    "rp-go/engine/systems/debug"
    "rp-go/engine/systems/devconsole"
    "rp-go/engine/systems/entitylist"
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

// NewGameWorld creates and initializes the full engine world with
// all core systems registered in deterministic update order.
func NewGameWorld() *GameWorld {
    // Load render configuration (viewport, scaling, etc.)
    cfg := data.LoadRenderConfig("engine/data/render_config.json")

    // Initialize ECS world and event bus
    w := ecs.NewWorld()
    w.EventBus = events.NewBus()

    // Scene manager handles scene transitions and deferred loading.
    sceneManager := &scene.Manager{}

    // Create system instances
    actorSystem := actor.NewSystem()
    aiSystem := ai.NewSystem()
    aiSystem.SetActorLookup(actorSystem.Registry())

    // Developer console + overlays share the actor registry
    consoleSystem := devconsole.NewSystem(actorSystem.Registry())
    entityListSystem := entitylist.NewSystem(actorSystem.Registry())

    /* -------------------------- Simulation Phase --------------------------- */
    // These systems update gameplay state but do not draw.
    simulationSystems := []ecs.System{
        sceneManager,           // Scene transitions, state management
        actorSystem,            // Actor registry and ownership
        &input.System{},        // Player + controller input
        aiSystem,               // AI controllers and decisions
        &movement.System{},     // Position/velocity updates
        camera.NewSystem(camera.Config{
            MinScale: cfg.Viewport.MinScale,
            MaxScale: cfg.Viewport.MaxScale,
            ZoomStep: cfg.Viewport.ZoomStep,
            ZoomLerp: cfg.Viewport.ZoomLerp,
        }),
    }

    /* --------------------------- Rendering Phase --------------------------- */
    // These systems draw world-space and overlay visuals.
    renderingSystems := []ecs.System{
        &background.System{},   // 🌌 Parallax background stars
        &render.System{},       // World-space sprite rendering
        entityListSystem,       // Entity overlay (actors + positions)
        &debug.System{},        // Diagnostic overlay (FPS, entities, etc.)
        consoleSystem,          // Developer console overlay (F12 toggle)
    }

    /* --------------------------- System Binding ---------------------------- */
    // Add systems in simulation phase first, then rendering phase.
    for _, system := range simulationSystems {
        w.AddSystem(system)
    }
    for _, system := range renderingSystems {
        w.AddSystem(system)
    }

    // Start in the default scene
    sceneManager.QueueScene(&space.Scene{})

    return &GameWorld{
        World:  w,
        Config: cfg,
    }
}

// Update advances the world simulation by one frame and flushes events.
func (g *GameWorld) Update() {
    g.World.Update()
    if bus, ok := g.World.EventBus.(*events.TypedBus); ok && bus != nil {
        bus.Flush()
    }
}

// Draw executes only world-space (camera-affected) rendering systems.
// Overlay systems (debug, HUD, console) are drawn later in main.go.
func (g *GameWorld) Draw(screen *platform.Image) {
    g.World.DrawWorld(screen)
}

