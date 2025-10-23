package devconsole

import (
    "fmt"
    "sort"
    "strconv"
    "strings"

    "rp-go/engine/ecs"
    "rp-go/engine/data"
    "rp-go/engine/world"
)

func (s *ConsoleState) ExecuteCommand(w *ecs.World, input string) {
    command := strings.TrimSpace(input)
    if command == "" {
        return
    }
    s.PushHistory(command)

    fields := strings.Fields(command)

    switch strings.ToLower(fields[0]) {
    case "help":
        s.Log("Commands: help, spawn <template> [x y], remove <actorID>, move <actorID> <x y>, list")
    case "spawn":
        s.HandleSpawn(w, fields)
    case "remove", "rm":
        s.HandleRemove(w, fields)
    case "move", "teleport":
        s.HandleMove(w, fields)
    case "list":
        s.HandleList(w)
    default:
        s.Log(fmt.Sprintf("Unknown command: %s", fields[0]))
    }
}

func (s *ConsoleState) HandleSpawn(w *ecs.World, fields []string) {
    if len(fields) < 2 {
        s.Log("Usage: spawn <template> [x y]")
        return
    }

    if s.Creator == nil {
        db := data.LoadActorDatabase("engine/data/actors.json")
        s.Creator = world.NewActorCreator(db)
    }

    x, y := 0.0, 0.0
    if len(fields) >= 3 {
        if parsed, err := strconv.ParseFloat(fields[2], 64); err == nil {
            x = parsed
        } else {
            s.Log(fmt.Sprintf("Invalid X coordinate: %q", fields[2]))
            return
        }
    }
    if len(fields) >= 4 {
        if parsed, err := strconv.ParseFloat(fields[3], 64); err == nil {
            y = parsed
        } else {
            s.Log(fmt.Sprintf("Invalid Y coordinate: %q", fields[3]))
            return
        }
    }

    entity, err := s.Creator.Spawn(w, fields[1], ecs.Position{X: x, Y: y})
    if err != nil {
        s.Log(err.Error())
        if templates := s.listTemplates(); len(templates) > 0 {
            s.Log("Templates: " + strings.Join(templates, ", "))
        }
        return
    }

    if actorComp, ok := entity.Get("Actor").(*ecs.Actor); ok && actorComp != nil {
        s.Log(fmt.Sprintf("Spawned %s at (%.1f, %.1f)", actorComp.ID, x, y))
    } else {
        s.Log(fmt.Sprintf("Spawned entity %d", entity.ID))
    }
}

func (s *ConsoleState) HandleRemove(w *ecs.World, fields []string) {
    if len(fields) < 2 {
        s.Log("Usage: remove <actorID>")
        return
    }
    target := s.findActorByID(w, fields[1])
    if target == nil {
        s.Log(fmt.Sprintf("Actor %q not found", fields[1]))
        return
    }
    w.RemoveEntity(target)
    s.Log(fmt.Sprintf("Removed actor %s", fields[1]))
}

func (s *ConsoleState) HandleMove(w *ecs.World, fields []string) {
    if len(fields) < 4 {
        s.Log("Usage: move <actorID> <x> <y>")
        return
    }

    x, errX := strconv.ParseFloat(fields[2], 64)
    y, errY := strconv.ParseFloat(fields[3], 64)
    if errX != nil || errY != nil {
        s.Log("Invalid coordinates. Expected numbers for x and y.")
        return
    }

    target := s.findActorByID(w, fields[1])
    if target == nil {
        s.Log(fmt.Sprintf("Actor %q not found", fields[1]))
        return
    }

    pos, ok := target.Get("Position").(*ecs.Position)
    if !ok {
        pos = &ecs.Position{}
        target.Add(pos)
    }
    pos.X, pos.Y = x, y

    s.Log(fmt.Sprintf("Moved %s to (%.1f, %.1f)", fields[1], x, y))
}

func (s *ConsoleState) HandleList(w *ecs.World) {
    entities := s.collectActors(w)
    if len(entities) == 0 {
        s.Log("No actors registered.")
        return
    }
    for _, entry := range entities {
        s.Log(entry)
    }
}

func (s *ConsoleState) listTemplates() []string {
    if s.Creator == nil {
        return nil
    }
    names := s.Creator.Templates()
    sort.Strings(names)
    return names
}

