package help

// Implements a window very similar to a popupbox, but with more controls, and
// Options to open to a specific help page.

import (
	"github.com/yamamushi/EscapingEden/ui/types"
	"github.com/yamamushi/EscapingEden/ui/window"
	"sync"
)

// While a normal popupbox only has controls to close the window (return), the help screen
// It is able to navigate the help pages, which are defined here as consts.
type HelpWindow struct {
	window.Window

	HelpPage types.HelpPage

	// Threading stuff if we need it
	twMutex sync.Mutex
}

func NewHelpWindow(x, y, w, h, consoleWidth, consoleHeight int, page types.HelpPage, input, output chan string) *HelpWindow {
	lw := &HelpWindow{}
	lw.ID = window.HELPBOX
	// if x or y are less than 1 set them to 1
	if x < 1 {
		x = 1
	}
	if y < 1 {
		y = 1
	}
	lw.X = x
	lw.Y = y

	// if w or h are less than 1 set them to 1
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	lw.Width = w
	lw.Height = h
	lw.ConsoleWidth = consoleWidth
	lw.ConsoleHeight = consoleHeight
	lw.Bordered = true
	lw.ConsoleReceive = input
	lw.ConsoleSend = output

	if page < 0 {
		page = 0
	}
	lw.HelpPage = page

	return lw
}
