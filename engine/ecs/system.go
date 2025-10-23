package ecs

import (
	"reflect"
	"rp-go/engine/platform"
)

/*───────────────────────────────────────────────*
 | SYSTEM INTERFACES                             |
 *───────────────────────────────────────────────*/

// System defines any update-only system (simulation, physics, AI, etc.)
type System interface {
	Update(world *World)
}

// DrawableSystem extends System with draw capability.
type DrawableSystem interface {
	System
	Draw(world *World, screen *platform.Image)
}

// PrioritizedSystem allows systems to specify execution priority.
// Lower numbers execute first.
type PrioritizedSystem interface {
	System
	Priority() int
}

// LayeredSystem marks systems that render to a specific draw layer.
type LayeredSystem interface {
	DrawableSystem
	Layer() DrawLayer
}

/*───────────────────────────────────────────────*
 | DRAW LAYERS                                   |
 *───────────────────────────────────────────────*/

type DrawLayer int

const (
	LayerBackground DrawLayer = iota // Parallax, distant scenery
	LayerWorld                       // Primary entities, world space
	LayerForeground                  // Effects rendered after entities
	LayerHUD                         // UI overlays in screen space
	LayerEntityList                  // Debug entity lists
	LayerDebug                       // Diagnostics, bounding boxes
	LayerConsole                     // Developer console, drawn last
)

/*───────────────────────────────────────────────*
 | UTILS                                         |
 *───────────────────────────────────────────────*/

// SystemName returns a readable name for debug or profiling logs.
func SystemName(s System) string {
	if s == nil {
		return "<nil>"
	}
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

