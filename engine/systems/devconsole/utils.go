package devconsole

import (
    "fmt"
    "sort"
    "rp-go/engine/ecs"
)

func (s *ConsoleState) Log(message string) {
    if message == "" {
        return
    }
    s.LogMessages = append(s.LogMessages, message)
    if len(s.LogMessages) > maxLogEntries {
        s.LogMessages = s.LogMessages[len(s.LogMessages)-maxLogEntries:]
    }
}

func (s *ConsoleState) findActorByID(w *ecs.World, id string) *ecs.Entity {
    if s.Registry != nil {
        if e, ok := s.Registry.FindByID(id); ok {
            return e
        }
    }
    for _, e := range w.Entities {
        if actorComp, ok := e.Get("Actor").(*ecs.Actor); ok && actorComp != nil && actorComp.ID == id {
            return e
        }
    }
    return nil
}

func (s *ConsoleState) collectActors(w *ecs.World) []string {
    var entities []*ecs.Entity
    if s.Registry != nil {
        entities = s.Registry.All()
    }
    if len(entities) == 0 {
        for _, e := range w.Entities {
            if actorComp, ok := e.Get("Actor").(*ecs.Actor); ok && actorComp != nil {
                entities = append(entities, e)
            }
        }
    }

    if len(entities) == 0 {
        return nil
    }

    descriptions := make([]string, 0, len(entities))
    for _, e := range entities {
        actorComp, _ := e.Get("Actor").(*ecs.Actor)
        pos, _ := e.Get("Position").(*ecs.Position)
        if actorComp == nil {
            continue
        }
        if pos != nil {
            descriptions = append(descriptions, fmt.Sprintf("%s (%.1f, %.1f)", actorComp.ID, pos.X, pos.Y))
        } else {
            descriptions = append(descriptions, fmt.Sprintf("%s (no position)", actorComp.ID))
        }
    }

    sort.Strings(descriptions)
    if len(descriptions) > maxLogEntries {
        return descriptions[:maxLogEntries]
    }
    return descriptions
}

