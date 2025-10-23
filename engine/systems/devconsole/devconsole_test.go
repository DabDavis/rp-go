package devconsole

import (
	"testing"

	"rp-go/engine/ecs"
)

/*───────────────────────────────────────────────*
 | MOCK HELPERS                                  |
 *───────────────────────────────────────────────*/

var mockEntityCounter int

// mockEntity simulates a minimal ECS entity with an ID and components.
func mockEntity(id string, comps map[string]interface{}) *ecs.Entity {
	mockEntityCounter++
	e := ecs.NewEntity(ecs.EntityID(mockEntityCounter))
	e.AddNamed("ID", id)
	for k, v := range comps {
		e.AddNamed(k, v)
	}
	return e
}

// mockWorld returns a simple ECS world containing two mock entities.
func mockWorld() *ecs.World {
	w := ecs.NewWorld()

	e1 := mockEntity("hero", map[string]interface{}{
		"Actor":    &ecs.Actor{ID: "hero"},
		"Position": &ecs.Position{X: 10, Y: 5},
	})
	e2 := mockEntity("npc_guard", map[string]interface{}{
		"Actor":    &ecs.Actor{ID: "npc_guard"},
		"Position": &ecs.Position{X: -3, Y: 7},
	})

	w.EntitiesManager().Add(e1)
	w.EntitiesManager().Add(e2)
	return w
}

// mockRegistry fakes an actor registry.
type mockRegistry struct {
	entities []*ecs.Entity
}

func (m *mockRegistry) FindByID(id string) (*ecs.Entity, bool) {
	for _, e := range m.entities {
		if actor, ok := e.Get("Actor").(*ecs.Actor); ok && actor.ID == id {
			return e, true
		}
	}
	return nil, false
}

func (m *mockRegistry) All() []*ecs.Entity {
	return append([]*ecs.Entity{}, m.entities...)
}

/*───────────────────────────────────────────────*
 | TESTS                                         |
 *───────────────────────────────────────────────*/

// TestConsoleState_BasicLifecycle verifies open/close state toggling and logging.
func TestConsoleState_BasicLifecycle(t *testing.T) {
	reg := &mockRegistry{}
	console := NewConsoleState(nil)
	if console.Open {
		t.Fatal("expected console to start closed")
	}

	console.Open = true
	console.Log("test message")
	if len(console.LogMessages) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(console.LogMessages))
	}

	console.Open = false
	if consoleOpen.Load() {
		t.Errorf("expected global IsOpen() to be false when closed")
	}
}

// TestConsoleState_HistoryNavigation ensures history navigation works correctly.
func TestConsoleState_HistoryNavigation(t *testing.T) {
	console := NewConsoleState(nil)
	console.History = []string{"help", "spawn test", "list"}
	console.HistoryIdx = -1

	console.NavigateHistory(-1)
	if console.InputBuffer != "list" {
		t.Errorf("expected 'list', got %q", console.InputBuffer)
	}

	console.NavigateHistory(-1)
	if console.InputBuffer != "spawn test" {
		t.Errorf("expected 'spawn test', got %q", console.InputBuffer)
	}

	console.NavigateHistory(1)
	if console.InputBuffer != "list" {
		t.Errorf("expected 'list' after moving down, got %q", console.InputBuffer)
	}
}

// TestConsoleState_PushHistory ensures command history rolls correctly.
func TestConsoleState_PushHistory(t *testing.T) {
	console := NewConsoleState(nil)
	for i := 0; i < maxHistoryStored+2; i++ {
		console.PushHistory("cmd")
	}
	if len(console.History) != maxHistoryStored {
		t.Errorf("expected max %d history entries, got %d", maxHistoryStored, len(console.History))
	}
}

// TestCollectActors ensures collectActors returns sorted actor list.
func TestCollectActors(t *testing.T) {
	w := mockWorld()
	reg := &mockRegistry{entities: w.EntitiesManager().All()}
	console := NewConsoleState(nil)
	console.Registry = reg

	list := console.collectActors(w)
	if len(list) != 2 {
		t.Fatalf("expected 2 actors, got %d", len(list))
	}
	if list[0] > list[1] {
		t.Errorf("expected sorted output, got: %v", list)
	}
}

// TestFindActorByID ensures actor lookup by ID works.
func TestFindActorByID(t *testing.T) {
	w := mockWorld()
	reg := &mockRegistry{entities: w.EntitiesManager().All()}
	console := NewConsoleState(reg)

	e := console.findActorByID(w, "hero")
	if e == nil {
		t.Fatal("expected to find 'hero'")
	}

	if _, ok := e.Get("Actor").(*ecs.Actor); !ok {
		t.Error("expected entity to have Actor component")
	}

	if console.findActorByID(w, "missing") != nil {
		t.Error("expected nil for missing actor")
	}
}

