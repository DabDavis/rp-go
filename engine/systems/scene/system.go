package scene

import (
	"fmt"

	"rp-go/engine/ecs"
	"rp-go/engine/events"
	"rp-go/engine/platform"
)

type Manager struct {
	current ecs.Scene
	next    ecs.Scene
	init    bool
}

func (m *Manager) Update(w *ecs.World) {
	// Subscribe once to SceneChangeEvent
	if !m.init {
		m.init = true
		if bus, ok := w.EventBus.(*events.TypedBus); ok && bus != nil {
			events.Subscribe(bus, func(e events.SceneChangeEvent) {
				fmt.Printf("[SCENE] Switching from %s → %s\n", m.currentName(), e.Target)
				if scn, ok := e.Scene.(ecs.Scene); ok {
					m.QueueScene(scn)
				}
			})
		}
	}

	// Scene transition
	if m.next != nil {
		if m.current != nil {
			m.current.Unload(w)
		}
		m.current = m.next
		m.next = nil
		m.current.Init(w)
		fmt.Printf("[SCENE] Active: %s\n", m.current.Name())
	}

	if m.current != nil {
		m.current.Update(w)
	}
}

func (m *Manager) Draw(w *ecs.World, screen *platform.Image) {
	if m.current != nil {
		m.current.Draw(w, screen)
	}
}

func (m *Manager) QueueScene(scene ecs.Scene) {
	if scene == nil {
		return
	}

	// Always funnel through m.next so the lifecycle consistently
	// triggers Init/Unload inside Update. The previous implementation
	// assigned the very first scene straight to m.current, which meant
	// Init was never called and no entities (camera, player, etc.) were
	// spawned, leading to the render system spamming "No camera found".
	m.next = scene
}

func (m *Manager) currentName() string {
	if m.current != nil {
		return m.current.Name()
	}
	return "(none)"
}
