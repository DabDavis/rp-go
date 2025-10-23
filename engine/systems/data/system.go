package data

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"rp-go/engine/data"
	"rp-go/engine/ecs"
	"rp-go/engine/events"
)

/*───────────────────────────────────────────────*
 | DATA SYSTEM                                   |
 *───────────────────────────────────────────────*/

// System manages engine-wide configuration, JSON databases, and live reloads.
type System struct {
	Config     data.RenderConfig   // Render config
	Actors     data.ActorDatabase  // Actor definitions
	AICatalog  data.AIActionCatalog // AI behavior definitions

	reloadMgr  *HotReloadManager
	subscriber *DataSubscriber
	lastUpdate time.Time
	mu         sync.RWMutex
}

// NewSystem initializes the DataSystem and registers all engine data files.
func NewSystem() *System {
	s := &System{
		reloadMgr:  NewHotReloadManager(),
		subscriber: NewDataSubscriber(),
	}
	s.RegisterDataFile("render_config", "engine/data/render_config.json")
	s.RegisterDataFile("actor_db", "engine/data/actors.json")
	s.RegisterDataFile("ai_catalog", "engine/data/ai.json")
	return s
}

/*───────────────────────────────────────────────*
 | LIFECYCLE                                     |
 *───────────────────────────────────────────────*/

func (s *System) Update(world *ecs.World) {
	if s.reloadMgr == nil {
		s.reloadMgr = NewHotReloadManager()
	}

	now := time.Now()
	if now.Sub(s.lastUpdate) < 500*time.Millisecond {
		return
	}
	s.lastUpdate = now

	changed := s.reloadMgr.CheckChanges()
	if len(changed) == 0 {
		s.ensureLoaded()
		return
	}

	for _, path := range changed {
		s.reloadFile(world, path)
	}
}

// ReloadAll forces full reload of all known data and emits global event.
func (s *System) ReloadAll(world *ecs.World) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Config = data.LoadRenderConfig("engine/data/render_config.json")
	s.Actors = data.LoadActorDatabase("engine/data/actors.json")
	s.AICatalog = data.LoadAICatalog("engine/data/ai.json")

	fmt.Println("[DATA] Reloaded all configuration, actors, and AI catalog")

	if bus, ok := world.EventBus.(*events.TypedBus); ok {
		evt := events.DataReloaded{Path: "engine/data", Type: "all"}
		events.Queue(bus, evt)
		go s.subscriber.Notify(evt)
	}
}

/*───────────────────────────────────────────────*
 | RELOAD LOGIC                                  |
 *───────────────────────────────────────────────*/

func (s *System) reloadFile(world *ecs.World, path string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	base := filepath.Base(path)
	var evt events.DataReloaded

	switch base {
	case "render_config.json":
		s.Config = data.LoadRenderConfig(path)
		fmt.Println("[DATA] Reloaded render_config")
		evt = events.DataReloaded{Path: path, Type: "render_config"}

	case "actors.json":
		s.Actors = data.LoadActorDatabase(path)
		fmt.Println("[DATA] Reloaded actor_db")
		evt = events.DataReloaded{Path: path, Type: "actor_db"}

	case "ai.json":
		s.AICatalog = data.LoadAICatalog(path)
		fmt.Println("[DATA] Reloaded ai_catalog")
		evt = events.DataReloaded{Path: path, Type: "ai_catalog"}

	default:
		fmt.Printf("[DATA] Reloaded generic file: %s\n", path)
		evt = events.DataReloaded{Path: path, Type: "generic"}
	}

	if world != nil && world.EventBus != nil {
		if bus, ok := world.EventBus.(*events.TypedBus); ok {
			events.Queue(bus, evt)
		}
	}
	go s.subscriber.Notify(evt)
}

// ensureLoaded guarantees all data sets are loaded once at startup.
func (s *System) ensureLoaded() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Config.Window.Width == 0 {
		s.Config = data.LoadRenderConfig("engine/data/render_config.json")
	}
	if len(s.Actors.Actors) == 0 {
		s.Actors = data.LoadActorDatabase("engine/data/actors.json")
	}
	if len(s.AICatalog.Actions) == 0 {
		s.AICatalog = data.LoadAICatalog("engine/data/ai.json")
	}
}

/*───────────────────────────────────────────────*
 | MANAGEMENT                                    |
 *───────────────────────────────────────────────*/

func (s *System) RegisterDataFile(kind string, path string) {
	s.reloadMgr.Watch(path)
	fmt.Printf("[DATA] Registered %s for hot reload: %s\n", kind, path)
}

func (s *System) Subscriber() *DataSubscriber { return s.subscriber }

/*───────────────────────────────────────────────*
 | HOT RELOAD MANAGER                            |
 *───────────────────────────────────────────────*/

type HotReloadManager struct {
	mu    sync.RWMutex
	files map[string]time.Time
}

func NewHotReloadManager() *HotReloadManager {
	return &HotReloadManager{files: make(map[string]time.Time)}
}

func (h *HotReloadManager) Watch(path string) {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Printf("[HOTRELOAD] Missing %s (will monitor when available)\n", path)
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.files[path] = info.ModTime()
	fmt.Printf("[HOTRELOAD] Watching %s\n", path)
}

func (h *HotReloadManager) CheckChanges() []string {
	h.mu.Lock()
	defer h.mu.Unlock()

	var changed []string
	for path, lastMod := range h.files {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		if info.ModTime().After(lastMod) {
			h.files[path] = info.ModTime()
			changed = append(changed, path)
		}
	}
	return changed
}

