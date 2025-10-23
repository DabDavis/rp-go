package ecs

// EntityManager centralizes iteration and lookup logic for entities within a
// world. It provides helpers for common traversal patterns so systems no
// longer duplicate iteration boilerplate.
type EntityManager struct {
	world *World
}

func newEntityManager(world *World) *EntityManager {
	return &EntityManager{world: world}
}

// ForEach executes fn for every non-nil entity currently alive in the world.
func (m *EntityManager) ForEach(fn func(*Entity)) {
	if m == nil || m.world == nil || fn == nil {
		return
	}
	for _, entity := range m.world.Entities {
		if entity != nil {
			fn(entity)
		}
	}
}

// ForEachComponent executes fn for each entity that has the named component.
func (m *EntityManager) ForEachComponent(name string, fn func(*Entity, Component)) {
	if m == nil || m.world == nil || fn == nil {
		return
	}
	for _, entity := range m.world.Entities {
		if entity == nil {
			continue
		}
		if component := entity.Get(name); component != nil {
			fn(entity, component)
		}
	}
}

// FirstComponent returns the first entity that owns the named component along
// with the component instance. If none are found it returns nil, nil.
func (m *EntityManager) FirstComponent(name string) (*Entity, Component) {
	if m == nil || m.world == nil {
		return nil, nil
	}
	for _, entity := range m.world.Entities {
		if entity == nil {
			continue
		}
		if component := entity.Get(name); component != nil {
			return entity, component
		}
	}
	return nil, nil
}

// Count returns the total number of live entities.
func (m *EntityManager) Count() int {
	if m == nil || m.world == nil {
		return 0
	}
	return len(m.world.Entities)
}
