package popupbox

import (
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"github.com/yamamushi/EscapingEden/ui/util"
	"github.com/yamamushi/EscapingEden/ui/window"
	"log"
	"sync"
)

type PopupBox struct {
	window.Window
	pbMutex sync.Mutex
}

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

func (pb *PopupBox) HandleInput(input types.Input) {
	pb.pbMutex.Lock()
	defer pb.pbMutex.Unlock()

	if pb.GetActive() {
		log.Println("PopupBox Handling input")
	}

	switch input.Type {
	case types.InputUp:
		log.Println("PopupBox Up")
		pb.DecreaseContentPos()
		return
	case types.InputDown:
		log.Println("PopupBox Down")
		pb.IncreaseContentPos()
		return
	case types.InputReturn:
		log.Println("PopupBox Handling input return - attempting to close popup")
		message := types.ConsoleMessage{Type: "popupbox", Message: "close"}
		pb.ConsoleSend <- message.String()
		log.Println("PopupBox sent close message to console")
	}

}

// DrawBorder returns the border of a window using code page 437 characters as a string
func (pb *PopupBox) DrawBorder(winX int, winY int) {
	// Draw top border using code page 437 characters starting at position winX, winY

	// Move cursor to top left corner of window
	// Draw top left corner
	if pb.Active {
		pb.PrintChar(winX, winY, "\u250c", "\033[32m")

	} else {
		pb.PrintChar(winX, winY, "\u250c", "")
	}

	// Draw left border
	for i := 1; i < pb.Height+1; i++ {
		// Inserts a vertical line
		if pb.Active {
			pb.PrintChar(winX, winY+i, "\u2502", "\033[32m")
		} else {
			pb.PrintChar(winX, winY+i, "\u2502", "")
		}
	}
	// Draw bottom left corner
	if pb.Active {
		pb.PrintChar(winX, winY+pb.Height+1, "\u2514", "\033[32m")
	} else {
		pb.PrintChar(winX, winY+pb.Height+1, "\u2514", "")
	}

	// Draw top border
	for i := 1; i < pb.Width; i++ {
		// Inserts a horizontal line
		if pb.Active {
			pb.PrintChar(winX+i, winY, "\u2500", "\033[32m")
		} else {
			pb.PrintChar(winX+i, winY, "\u2500", "")
		}
	}

	// Draw top right corner
	if pb.Active {
		pb.PrintChar(winX+pb.Width, winY, "\u2510", "\033[32m")
	} else {
		pb.PrintChar(winX+pb.Width, winY, "\u2510", "")
	}

	// Draw right border
	for i := 1; i < pb.Height+1; i++ {
		// Inserts a vertical line
		if pb.Active {
			pb.PrintChar(winX+pb.Width, winY+i, "\u2502", "\033[32m")
		} else {
			pb.PrintChar(winX+pb.Width, winY+i, "\u2502", "")
		}
	}

	// Draw bottom right corner
	if pb.Active {
		pb.PrintChar(winX+pb.Width, winY+pb.Height+1, "\u2518", "\033[32m")
	} else {
		pb.PrintChar(winX+pb.Width, winY+pb.Height+1, "\u2518", "")
	}

	// Draw bottom border
	for i := 1; i < pb.Width; i++ {
		// Inserts a horizontal line
		if pb.Active {
			pb.PrintChar(winX+i, winY+pb.Height+1, "\u2500", "\033[32m")
		} else {
			pb.PrintChar(winX+i, winY+pb.Height+1, "\u2500", "")
		}
	}

	// Prints the world "Close" at the bottom center
	if pb.Active {
		var colorCode string
		// set colorCode background to green
		colorCode = util.RGBCode(0, 0, 0).FG()
		colorCode += util.RGBCode(0, 255, 0).BG()

		pb.PrintLn(winX+pb.Width/2-3, winY+pb.Height+1, "Close", colorCode)
	} else {
		log.Println("Writing Close to window")
		pb.PrintLn(winX+pb.Width/2-3, winY+pb.Height+1, "Close", "")
	}
}

func (pb *PopupBox) UpdateContents() {
	pb.pbMutex.Lock()
	defer pb.pbMutex.Unlock()

}
