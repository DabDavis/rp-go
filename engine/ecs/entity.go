package ecs

type EntityID int

type Entity struct {
	ID         EntityID
	Components map[string]Component
}

func NewEntity(id EntityID) *Entity {
	return &Entity{ID: id, Components: make(map[string]Component)}
}

func (e *Entity) Add(c Component) {
	e.Components[c.Name()] = c
}

func (e *Entity) Get(name string) Component {
	return e.Components[name]
}

func (e *Entity) Has(name string) bool {
	_, ok := e.Components[name]
	return ok
}

func (e *Entity) Remove(name string) {
	delete(e.Components, name)
}

