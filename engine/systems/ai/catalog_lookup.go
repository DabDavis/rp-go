package ai

import "rp-go/engine/data"

// AIActionCatalogLookup provides fast lookup by name.
type AIActionCatalogLookup struct {
	byName map[string]data.AIActionTemplate
}

func NewCatalogLookup(cat data.AIActionCatalog) *AIActionCatalogLookup {
	m := make(map[string]data.AIActionTemplate)
	for _, a := range cat.Actions {
		m[a.Name] = a
	}
	return &AIActionCatalogLookup{byName: m}
}

func (l *AIActionCatalogLookup) Get(name string) (data.AIActionTemplate, bool) {
	if l == nil {
		return data.AIActionTemplate{}, false
	}
	v, ok := l.byName[name]
	return v, ok
}

