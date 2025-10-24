package world

import (
	"fmt"

	"rp-go/engine/ecs"
	"rp-go/engine/events"
	"rp-go/engine/platform"

	"rp-go/engine/systems/ai"
	"rp-go/engine/systems/aicomposer"
	"rp-go/engine/systems/data"
	"rp-go/engine/systems/devconsole"
	"rp-go/engine/systems/render"
	"rp-go/engine/systems/scene"
	"rp-go/engine/systems/ui"
)

/*───────────────────────────────────────────────*
 | WORLD INITIALIZATION                          |
 *───────────────────────────────────────────────*/

// InitWorld constructs and wires up all ECS systems and the event bus.
func InitWorld() *ecs.World {
	fmt.Println("[WORLD] Initializing ECS world")

	w := ecs.NewWorld()

	// Central event bus for cross-system communication
	bus := events.NewTypedBus()
	w.EventBus = bus
	fmt.Println("[WORLD] Event bus initialized")

	/*───────────────────────────────────────────────*
	 | CORE SYSTEMS                                 |
	 *───────────────────────────────────────────────*/

	dataSys := data.NewSystem()
	w.AddSystem(dataSys)

	aiSys := ai.NewSystem(dataSys.AICatalog)
	w.AddSystem(aiSys)

	composerSys := aicomposer.NewSystem(dataSys, aiSys)
	w.AddSystem(composerSys)

	renderSys := render.NewSystem()
	w.AddSystem(renderSys)

	uiSys := ui.NewSystem()
	w.AddSystem(uiSys)

	sceneSys := &scene.Manager{}
	w.AddSystem(sceneSys)

	consoleSys := devconsole.NewSystem()
	w.AddSystem(consoleSys)

	fmt.Println("[WORLD] Core systems registered")

	/*───────────────────────────────────────────────*
	 | DATA-DEPENDENT LINKING                        |
	 *───────────────────────────────────────────────*/

	// Subscribe systems to reload events
	events.Subscribe(bus, func(e events.DataReloaded) {
		if dataSys != nil {
			dataSys.OnDataReload(e)
		}
		if aiSys != nil {
			aiSys.OnDataReload(e, dataSys.AICatalog)
		}
		if composerSys != nil {
			composerSys.OnDataReload(e)
		}
	})

	fmt.Println("[WORLD] Data reload subscriptions active")

	/*───────────────────────────────────────────────*
	 | FINALIZATION                                  |
	 *───────────────────────────────────────────────*/

	fmt.Println("[WORLD] Initialization complete")
	return w
}

/*───────────────────────────────────────────────*
 | ENTRY-POINT HELPERS                            |
 *───────────────────────────────────────────────*/

// RunWorld starts the main update/draw loop.
func RunWorld(w *ecs.World, screen *platform.Image) {
	if w == nil {
		fmt.Println("[WORLD] Cannot run: world is nil")
		return
	}

	fmt.Println("[WORLD] Running main loop")

	for {
		w.Update()
		w.DrawWorld(screen)
		w.DrawOverlay(screen)
		// Platform-specific frame sync (if available)
		if platform.SyncFrame != nil {
			platform.SyncFrame()
		}
	}
}

