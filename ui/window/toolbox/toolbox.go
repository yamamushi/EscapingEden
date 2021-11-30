package toolbox

import (
	"github.com/yamamushi/EscapingEden/edenutil"
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"github.com/yamamushi/EscapingEden/ui/window"
	"log"
	"sync"
	"time"
)

// ToolboxWindow is a window that contains a toolbox for misc use
type ToolboxWindow struct {
	window.Window
	twMutex sync.Mutex
}

// NewToolboxWindow creates a new toolbox window
func NewToolboxWindow(x, y, w, h, consoleWidth, consoleHeight int, input, output chan string) *ToolboxWindow {
	lw := &ToolboxWindow{}
	lw.ID = config.WindowToolBox
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

	return lw
}

// HandleInput handles input for the toolbox window
func (tw *ToolboxWindow) HandleInput(input types.Input) {
	tw.twMutex.Lock()
	defer tw.twMutex.Unlock()

	if tw.GetActive() {
		log.Println("Toolbox Handling input")
	}

	if len(input.Data) > 0 {
		log.Println(input.Data)
	}
}

// UpdateContents updates the contents of the toolbox window
func (tw *ToolboxWindow) UpdateContents() {
	tw.twMutex.Lock()
	defer tw.twMutex.Unlock()

	// current time with second accuracy as a string
	serverTime := time.Now().Format("15:04:05")
	edenTime := edenutil.EdenTime.CurrentTimeString(edenutil.EdenTime{})
	tw.SetContents("Current Server Time: " + serverTime + "\n" + "  Current Eden Time: " + edenTime)
}
