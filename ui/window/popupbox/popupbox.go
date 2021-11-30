package popupbox

import (
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/window"
	"sync"
)

// PopupBox is a window that can be used to display a message to the user.
type PopupBox struct {
	window.Window
	pbMutex sync.Mutex
}

// NewPopupBox creates a new PopupBox.
func NewPopupBox(x, y, w, h, consoleWidth, consoleHeight int, input, output chan string) *PopupBox {
	pb := &PopupBox{}
	pb.ID = config.WindowPopupBox
	// if x or y are less than 1 set them to 1
	if x < 1 {
		x = 1
	}
	if y < 1 {
		y = 1
	}
	pb.X = x
	pb.Y = y

	// if w or h are less than 1 set them to 1
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	pb.Width = w
	pb.Height = h
	pb.ConsoleWidth = consoleWidth
	pb.ConsoleHeight = consoleHeight
	pb.Bordered = true
	pb.ConsoleReceive = input
	pb.ConsoleSend = output
	pb.ScrollingSupported = true

	return pb
}

// UpdateContents updates the contents of the PopupBox.
func (pb *PopupBox) UpdateContents() {
	pb.pbMutex.Lock()
	defer pb.pbMutex.Unlock()

}
