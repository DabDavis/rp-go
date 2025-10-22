package ecs

// Component defines data-only objects attached to entities.
type Component interface {
	Name() string
}

// Tag is a simple flag-style component (no data, just presence).
type Tag string

func (t Tag) Name() string { return string(t) }

