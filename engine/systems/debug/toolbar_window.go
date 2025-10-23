package debug

import (
	"image/color"

	"rp-go/engine/ecs"
	"rp-go/engine/events"
	"rp-go/engine/ui/button"
	"rp-go/engine/ui/layout"
	"rp-go/engine/ui/window"
)

/*───────────────────────────────────────────────*
 | TOOLBAR WINDOW                                |
 *───────────────────────────────────────────────*/

// ToolbarWindow provides top-bar buttons to toggle key debug windows.
type ToolbarWindow struct {
	cfg Config
	bus *events.TypedBus

	component *window.Component
	buttons   []*button.Button
}

/*───────────────────────────────────────────────*
 | CONSTRUCTOR                                   |
 *───────────────────────────────────────────────*/

func NewToolbarWindow(cfg Config, bus *events.TypedBus) *ToolbarWindow {
	return &ToolbarWindow{
		cfg: cfg,
		bus: bus,
	}
}

/*───────────────────────────────────────────────*
 | WINDOW CREATION                               |
 *───────────────────────────────────────────────*/

func (t *ToolbarWindow) Ensure(world *ecs.World) {
	if t.component != nil {
		return
	}

	entity := world.NewEntity()

	content := layout.NewHorizontal(6, 4)

	btns := []*button.Button{
		button.New("Stats", func() {
			events.Queue(t.bus, events.DebugToggleWindowEvent{ID: "debug.stats"})
		}),
		button.New("Entities", func() {
			events.Queue(t.bus, events.DebugToggleWindowEvent{ID: "debug.entities"})
		}),
		button.New("Systems", func() {
			events.Queue(t.bus, events.DebugToggleWindowEvent{ID: "debug.systems"})
		}),
		button.New("Composer", func() { // ✅ new button for AIComposer window
			events.Queue(t.bus, events.DebugToggleWindowEvent{ID: "debug.aicomposer"})
		}),
		button.New("Hide All", func() {
			events.Queue(t.bus, events.DebugToggleEvent{Enabled: false})
		}),
	}

	for _, b := range btns {
		content.Add(b)
	}

	comp := window.NewComponent("debug.toolbar", "Debug Toolbar", window.Bounds{
		X:      20,
		Y:      20,
		Width:  480,
		Height: 42,
	}, content)

	comp.Layer = ecs.LayerHUD
	comp.TitleBarHeight = 0
	comp.Padding = 8
	comp.Movable = false
	comp.Closable = false
	comp.Resizable = false
	comp.Background = color.RGBA{30, 30, 45, 220}
	comp.Border = color.RGBA{90, 100, 120, 160}

	entity.Add(comp)

	t.component = comp
	t.buttons = btns
}

/*───────────────────────────────────────────────*
 | WINDOW BEHAVIOR                               |
 *───────────────────────────────────────────────*/

func (t *ToolbarWindow) Hide(world *ecs.World) {
	if t.component == nil {
		return
	}
	t.component.Visible = false
}

func (t *ToolbarWindow) Show(world *ecs.World) {
	if t.component == nil {
		t.Ensure(world)
	}
	t.component.Visible = true
}

func (t *ToolbarWindow) Toggle(world *ecs.World) {
	if t.component == nil {
		t.Ensure(world)
		return
	}
	t.component.Visible = !t.component.Visible
}

