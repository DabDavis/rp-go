package devconsole

import (
    "image/color"
    "golang.org/x/image/font/basicfont"
    "rp-go/engine/platform"
    "rp-go/engine/ecs"
)

const (
    maxLogEntries    = 12
    maxHistoryStored = 32
)

// Render draws the developer console overlay.
func (s *ConsoleState) Render(w *ecs.World, screen *platform.Image) {
    if !s.Open || screen == nil {
        return
    }

    bounds := screen.Bounds()
    width, height := bounds.Dx(), bounds.Dy()
    if width == 0 || height == 0 {
        return
    }

    consoleHeight := height / 3
    if consoleHeight < 180 {
        consoleHeight = 180
    }

    overlay := platform.NewImage(width, consoleHeight)
    overlay.FillRect(0, 0, width, consoleHeight, color.RGBA{0, 0, 0, 200})

    y := 24
    visible := s.LogMessages
    if len(visible) > maxLogEntries {
        visible = visible[len(visible)-maxLogEntries:]
    }
    for _, line := range visible {
        platform.DrawText(overlay, line, basicfont.Face7x13, 12, y, color.White)
        y += 16
    }

    prompt := "> " + s.InputBuffer
    if (s.CursorTick/20)%2 == 0 {
        prompt += "_"
    }
    platform.DrawText(overlay, prompt, basicfont.Face7x13, 12, consoleHeight-16, color.White)

    op := platform.NewDrawImageOptions()
    op.Translate(0, float64(height-consoleHeight))
    screen.DrawImage(overlay, op)
}

