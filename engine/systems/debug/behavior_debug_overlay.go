package debug

import (
	"fmt"
	"image/color"
	"rp-go/engine/ecs"
	"rp-go/engine/platform"
	"rp-go/engine/systems/ai"
	"rp-go/engine/ui/window"

	"golang.org/x/image/font/basicfont"
)

// BehaviorDebugOverlay displays current AI actions for the selected entity.
type BehaviorDebugOverlay struct {
	component *window.Component
	selected  *ecs.Entity
	cfg       Config
}

// NewBehaviorDebugOverlay constructs the overlay window.
func NewBehaviorDebugOverlay(cfg Config) *BehaviorDebugOverlay {
	return &BehaviorDebugOverlay{cfg: cfg}
}

// Ensure creates the overlay window if needed.
func (o *BehaviorDebugOverlay) Ensure(world *ecs.World) {
	if o.component != nil {
		return
	}
	o.component = window.NewComponent("debug.ai.behaviors", "AI Behaviors", window.Bounds{
		X:      o.cfg.Margin,
		Y:      o.cfg.Margin + 200,
		Width:  260,
		Height: 140,
	}, o)
	o.component.Layer = ecs.LayerDebug
	o.component.Movable = true
	o.component.Closable = true
	o.component.Visible = true
	o.component.Background = color.RGBA{20, 20, 30, 220}
	o.component.Border = color.RGBA{120, 150, 255, 200}
}

// Draw lists the actions of the selected entity’s AIController.
func (o *BehaviorDebugOverlay) Draw(world *ecs.World, screen *platform.Image, bounds window.Bounds) {
	if o.selected == nil || screen == nil {
		return
	}

	ctrl, _ := o.selected.Get("AIController").(*ai.AIController)
	if ctrl == nil {
		return
	}

	lines := []string{"Active Behaviors:"}
	for _, act := range ctrl.Actions {
		lines = append(lines, fmt.Sprintf("[%d] %s", act.Priority, act.Name))
	}

	y := bounds.Y + 16
	for _, line := range lines {
		platform.DrawText(screen, line, basicfont.Face7x13, bounds.X+8, y, color.RGBA{200, 220, 255, 255})
		y += 14
	}
}

// SetSelected sets the entity whose behaviors are displayed.
func (o *BehaviorDebugOverlay) SetSelected(e *ecs.Entity) {
	o.selected = e
}

