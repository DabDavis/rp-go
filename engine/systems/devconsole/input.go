package devconsole

import (
    "rp-go/engine/ecs"
    "rp-go/engine/platform"
)

// UpdateInput reads keyboard input, handles toggling,
// and triggers commands.
func (s *ConsoleState) UpdateInput(w *ecs.World) {
    // Toggle console on F12
    if platform.IsKeyJustPressed(platform.KeyF12) {
        s.Open = !s.Open
        s.JustOpened = s.Open
        if !s.Open {
            s.HistoryIdx = -1
        }
    }

    consoleOpen.Store(s.Open)
    if !s.Open {
        return
    }

    s.CursorTick++

    if s.JustOpened {
        s.Log("Developer console opened. Type 'help' for commands.")
        s.JustOpened = false
    }

    // Capture text input
    for _, char := range platform.InputChars() {
        if char != '\r' && char != '\n' {
            s.InputBuffer += string(char)
        }
    }

    // Editing and exit keys
    if platform.IsKeyJustPressed(platform.KeyBackspace) && len(s.InputBuffer) > 0 {
        s.InputBuffer = s.InputBuffer[:len(s.InputBuffer)-1]
    }
    if platform.IsKeyJustPressed(platform.KeyEscape) {
        s.Open = false
        consoleOpen.Store(false)
        s.HistoryIdx = -1
        return
    }

    // Command history
    if platform.IsKeyJustPressed(platform.KeyArrowUp) {
        s.NavigateHistory(-1)
    } else if platform.IsKeyJustPressed(platform.KeyArrowDown) {
        s.NavigateHistory(1)
    }

    // Execute on Enter
    if platform.IsKeyJustPressed(platform.KeyEnter) {
        s.ExecuteCommand(w, s.InputBuffer)
    }
}

