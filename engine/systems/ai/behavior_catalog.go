package ai

import (
	"fmt"
	"sync"

	"rp-go/engine/ecs"
)

/*───────────────────────────────────────────────*
 | BEHAVIOR FUNCTION TYPE                        |
 *───────────────────────────────────────────────*/

// BehaviorFunc defines a single AI action executor.
// Returns true if this action took control (e.g. movement was applied).
type BehaviorFunc func(world *ecs.World, entity *ecs.Entity, pos *ecs.Position, vel *ecs.Velocity, params map[string]any) bool

/*───────────────────────────────────────────────*
 | BEHAVIOR CATALOG                              |
 *───────────────────────────────────────────────*/

// BehaviorCatalog is a global registry of AI behavior handlers.
// It enables modular registration and hot-swapping of behaviors.
type BehaviorCatalog struct {
	mu        sync.RWMutex
	behaviors map[string]BehaviorFunc
}

// GlobalBehaviorCatalog is the shared instance used by all AI systems.
var GlobalBehaviorCatalog = NewBehaviorCatalog()

// NewBehaviorCatalog constructs an empty registry.
func NewBehaviorCatalog() *BehaviorCatalog {
	return &BehaviorCatalog{
		behaviors: make(map[string]BehaviorFunc),
	}
}

/*───────────────────────────────────────────────*
 | REGISTRATION                                  |
 *───────────────────────────────────────────────*/

// Register associates a behavior type name (e.g. "pursue") with a handler.
func (c *BehaviorCatalog) Register(name string, fn BehaviorFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if name == "" || fn == nil {
		return
	}
	c.behaviors[name] = fn
	fmt.Printf("[AI] Registered behavior: %s\n", name)
}

// Unregister removes a behavior from the catalog.
func (c *BehaviorCatalog) Unregister(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.behaviors, name)
	fmt.Printf("[AI] Unregistered behavior: %s\n", name)
}

// Get retrieves a behavior handler by name.
func (c *BehaviorCatalog) Get(name string) (BehaviorFunc, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	fn, ok := c.behaviors[name]
	return fn, ok
}

// List returns a snapshot of all registered behavior names.
func (c *BehaviorCatalog) List() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]string, 0, len(c.behaviors))
	for name := range c.behaviors {
		out = append(out, name)
	}
	return out
}

/*───────────────────────────────────────────────*
 | BOOTSTRAP DEFAULTS                            |
 *───────────────────────────────────────────────*/

// RegisterDefaultBehaviors installs built-in behaviors into the catalog.
// This should be called once during engine initialization.
func RegisterDefaultBehaviors(sys *System) {
	GlobalBehaviorCatalog.Register("pursue", sys.behaviorPursue)
	GlobalBehaviorCatalog.Register("patrol", sys.behaviorPatrol)
	GlobalBehaviorCatalog.Register("retreat", sys.behaviorRetreat)
	GlobalBehaviorCatalog.Register("follow", sys.behaviorFollow)
	GlobalBehaviorCatalog.Register("idle", func(*ecs.World, *ecs.Entity, *ecs.Position, *ecs.Velocity, map[string]any) bool {
		return false
	})
	fmt.Printf("[AI] Default behaviors registered (%d total)\n", len(GlobalBehaviorCatalog.List()))
}

