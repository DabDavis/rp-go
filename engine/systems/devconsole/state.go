package devconsole

import (
    "sync/atomic"
    "rp-go/engine/systems/actor"
)

// Global atomic flag for console visibility.
var consoleOpen atomic.Bool

// IsOpen reports whether the console overlay is currently active.
func IsOpen() bool { return consoleOpen.Load() }

// ConsoleState contains all runtime data for the dev console.
type ConsoleState struct {
    Registry    *actor.Registry
    Creator     ActorSpawner
    Open        bool
    JustOpened  bool
    CursorTick  int
    InputBuffer string
    History     []string
    HistoryIdx  int
    LogMessages []string
}

// ActorRegistry is a minimal interface for actor lookups.
type ActorRegistry interface {
    FindByID(id string) (*actor.Entity, bool)
    All() []*actor.Entity
}

// ActorSpawner abstracts spawning from templates.
type ActorSpawner interface {
    Spawn(w *actor.World, template string, pos actor.Position) (*actor.Entity, error)
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

