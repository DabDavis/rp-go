package ai

// AIActionInstance is a runtime version of a behavior template.
type AIActionInstance struct {
	Name     string
	Type     string
	Priority int
	Params   map[string]any
}

// AIController holds active actions per entity.
type AIController struct {
	Actions []AIActionInstance
	Active  bool
}

