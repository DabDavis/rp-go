//go:build !headless

package platform_desktop

import (
	"image/color"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

// DrawText safely draws text to an Ebiten image.
// Falls back to a built-in font if 'face' is nil.
func DrawText(dst *Image, str string, face font.Face, x, y int, clr color.Color) {
	if dst == nil || str == "" {
		return
	}

	// --- Use fallback font if face is nil
	if face == nil {
		face = basicfont.Face7x13
		logOnce("[WARN] DrawText called with nil font; using basicfont.Face7x13 fallback")
	}

	// --- Defensive color check (nil color causes panic in some Ebiten builds)
	if clr == nil {
		clr = color.White
	}

	text.Draw(dst.native, str, face, x, y, clr)
}

/*───────────────────────────────────────────────*
 | INTERNAL UTILITIES                            |
 *───────────────────────────────────────────────*/

var (
	onceWarnings sync.Map
)

// logOnce prints a message only once (per session) to avoid console spam.
func logOnce(msg string) {
	if _, exists := onceWarnings.LoadOrStore(msg, true); !exists {
		log.Println(msg)
	}
}

