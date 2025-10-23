package windowmgr

// SharedRegistry exposes the global registry created by system.go.
// This lightweight bridge file exists for clarity and modular imports.
func SharedRegistry() *Registry {
	return globalRegistry
}

