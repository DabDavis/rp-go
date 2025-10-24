package scene

import (
	"fmt"

	"rp-go/engine/ecs"
	"rp-go/engine/events"
	"rp-go/engine/platform"
)

/*───────────────────────────────────────────────*
 | SCENE MANAGER                                 |
 *───────────────────────────────────────────────*/

type Manager struct {
	current ecs.Scene
	next    ecs.Scene
	inited  bool
}

/*───────────────────────────────────────────────*
 | LIFECYCLE MANAGEMENT                          |
 *───────────────────────────────────────────────*/

func (m *Manager) Update(w *ecs.World) {
	if !m.inited {
		m.inited = true
		if bus, ok := w.EventBus.(*events.TypedBus); ok {
			events.Subscribe(bus, func(e events.SceneChangeEvent) {
				fmt.Printf("[SCENE] Switch %s → %s\n", m.name(m.current), e.Target)
				if scn, ok := e.Scene.(ecs.Scene); ok {
					m.QueueScene(scn)
				}
			})
		}
	}

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
	m.next = scene
}

/*───────────────────────────────────────────────*
 | UTILITIES                                     |
 *───────────────────────────────────────────────*/

func (m *Manager) name(s ecs.Scene) string {
	if s == nil {
		return "(none)"
	}
	return s.Name()
}

