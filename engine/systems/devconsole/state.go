package devconsole

import (
	"sync/atomic"

	"rp-go/engine/ecs"
	"rp-go/engine/systems/actor"
)

// Global atomic flag for console visibility.
var consoleOpen atomic.Bool

// IsOpen reports whether the console overlay is currently active.
func IsOpen() bool { return consoleOpen.Load() }

// ConsoleState contains all runtime data for the dev console.
type ConsoleState struct {
	Registry    *actor.Registry   // reference to ECS actor registry
	Creator     ActorSpawner      // interface for spawning
	Open        bool
	JustOpened  bool
	CursorTick  int
	InputBuffer string
	History     []string
	HistoryIdx  int
	LogMessages []string
}

// ActorRegistry defines a minimal interface to query actors.
type ActorRegistry interface {
	FindByID(id string) (*ecs.Entity, bool)
	All() []*ecs.Entity
}

// ActorSpawner abstracts spawning from templates.
type ActorSpawner interface {
	Spawn(w *ecs.World, template string, pos ecs.Position) (*ecs.Entity, error)
	Templates() []string
}

// Factory for initializing a fresh console state.
func NewConsoleState(reg *actor.Registry) *ConsoleState {
	consoleOpen.Store(false)
	return &ConsoleState{
		Registry:   reg,
		HistoryIdx: -1,
	}
}

