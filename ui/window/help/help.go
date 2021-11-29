package help

// Implements a window very similar to a popupbox, but with more controls, and
// Options to open to a specific help page.

import (
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"github.com/yamamushi/EscapingEden/ui/util"
	"github.com/yamamushi/EscapingEden/ui/window"
	"log"
	"strconv"
	"strings"
	"sync"
)

// While a normal popupbox only has controls to close the window (return), the help screen
// It is able to navigate the help pages, which are defined here as consts.
type HelpWindow struct {
	window.Window

	HelpPage     types.HelpPage
	LastHelpPage types.HelpPage

	// Threading stuff if we need it
	hwMutex           sync.Mutex
	scrollInitialized bool
	indexPage         int
}

func NewHelpWindow(x, y, w, h, consoleWidth, consoleHeight int, page types.HelpPage, input, output chan string) *HelpWindow {
	hw := &HelpWindow{}
	hw.ID = config.WindowHelpBox
	// if x or y are less than 1 set them to 1
	if x < 1 {
		x = 1
	}
	if y < 1 {
		y = 1
	}
	hw.X = x
	hw.Y = y

	// if w or h are less than 1 set them to 1
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	hw.Width = w
	hw.Height = h
	hw.ConsoleWidth = consoleWidth
	hw.ConsoleHeight = consoleHeight
	hw.Bordered = true
	hw.ConsoleReceive = input
	hw.ConsoleSend = output

	if page < 0 {
		page = 0
	}
	hw.HelpPage = page
	hw.StartY = 2
	hw.ScrollBufferLimit = 2
	hw.ScrollingSupported = true
	hw.ScrollBufferCharMod = 1

	return hw
}

func (hw *HelpWindow) UpdateContents() {
	hw.hwMutex.Lock()
	defer hw.hwMutex.Unlock()

	switch hw.HelpPage {
	case types.HelpPageMain:
		hw.LoadPage(types.HelpPageMain)
	case types.HelpPageRules:
		hw.LoadPage(types.HelpPageRules)
	case types.HelpPageDeath:
		hw.LoadPage(types.HelpPageDeath)
	case types.HelpPageAbout:
		hw.LoadPage(types.HelpPageAbout)
	case types.HelpPageControls:
		hw.LoadPage(types.HelpPageControls)
	case types.HelpPageCredits:
		hw.LoadPage(types.HelpPageCredits)
	case types.HelpPageIndex:
		hw.LoadPage(types.HelpPageIndex)
	default:
		hw.LoadPage(types.HelpPageMain)
	}
}

func (hw *HelpWindow) PrintPageInfo(page types.HelpPage) {
	// Top Field
	windowTitle := "Escaping Eden Help"
	pageInfo := strings.Title(page.String()) + " (Page " + strconv.Itoa(int(page)) + ")"
	hw.PrintLn(hw.X+1, hw.Y+1, windowTitle, "\033[1m")
	hw.PrintLn(hw.X+hw.Width-len(pageInfo)-1, hw.Y+1, pageInfo, "")
}

func (hw *HelpWindow) LoadPage(page types.HelpPage) {
	if page == types.HelpPageIndex {
		log.Println("Help Window Loading page index")
		hw.DrawIndex()
		return
	}
	// Prints our top field content
	hw.PrintPageInfo(page)

	// Prints our bottom field content
	hw.PrintControls()

	// We load the text file for the help page
	content, err := util.OpenFileAsText("assets/text/" + page.String() + ".txt")
	if err != nil {
		hw.Error("Error opening file:" + err.Error())
		return
	}
	// Now we set the content accordingly
	hw.SetContents(content)

	// If our scroll has been initialized already, we don't modify it because that would mess up
	// Any scrolling that was input.
	if !hw.scrollInitialized {
		// We need the visiblelLength of the window to determine the line count
		visibleLength := hw.Width - 1
		_, lineCount := hw.ContentToLines(hw.X, hw.Y, visibleLength)

		// We set the visible height with our scroll buffer limit accordingly
		visibleHeight := hw.Height - 1 - hw.ScrollBufferLimit - hw.StartY

		if lineCount > visibleHeight {
			hw.SetContentStartPos(-lineCount + visibleHeight)
			hw.scrollInitialized = true
		} else {
			hw.SetContentStartPos(0)
			hw.scrollInitialized = true
		}
	}
	return
}
