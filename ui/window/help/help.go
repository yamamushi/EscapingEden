package help

// Implements a window very similar to a popupbox, but with more controls, and
// Options to open to a specific help page.

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/terminals"
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"github.com/yamamushi/EscapingEden/ui/util"
	"github.com/yamamushi/EscapingEden/ui/window"
	"strconv"
	"strings"
	"sync"
)

// While a normal popupbox only has controls to close the window (return), the help screen
// It is able to navigate the help pages, which are defined here as consts.

// HelpWindow is a window that displays help pages.
type HelpWindow struct {
	window.Window

	HelpPage     types.HelpPage
	LastHelpPage types.HelpPage

	// Threading stuff if we need it
	hwMutex           sync.Mutex
	scrollInitialized bool

	indexMutex              sync.Mutex
	indexPage               int
	indexSelection          int
	IndexPageNames          IndexPageNames
	currentIndexRowSelected int
}

// NewHelpWindow creates a new help window.
func NewHelpWindow(x, y, w, h, consoleWidth, consoleHeight int, page types.HelpPage,
	input, output chan messages.WindowMessage, log logging.LoggerType, term terminals.TerminalType) *HelpWindow {
	hw := &HelpWindow{}
	hw.Log = log
	hw.Terminal = term
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

// UpdateContents updates the contents of the window.
// This method
func (hw *HelpWindow) UpdateContents() {
	hw.hwMutex.Lock()
	defer hw.hwMutex.Unlock()
	hw.LoadPage(hw.HelpPage)
}

// PrintPageInfo prints the title page info to the top field of the Help window
func (hw *HelpWindow) PrintPageInfo(page types.HelpPage) {
	// Top Field
	windowTitle := "Escaping Eden Help"
	pageInfo := strings.Title(page.String()) + " (Page " + strconv.Itoa(int(page)) + ")"
	hw.PrintLn(hw.X+1, hw.Y+1, windowTitle, hw.Terminal.Bold())
	hw.PrintLn(hw.X+hw.Width-len(pageInfo)-1, hw.Y+1, pageInfo, "")
}

// LoadPage takes a types.HelpPage and loads the corresponding help page into the window.
func (hw *HelpWindow) LoadPage(page types.HelpPage) {
	if page == types.HelpPageIndex {
		//log.Println("Help Window Loading index page")
		hw.DrawIndex()
		return
	}
	if page == types.HelpPageMain {
		//log.Println("Help Window Loading main page")
		hw.DrawHome()
		return
	}

	// If we're not on the Index or the Main page, we need to load some content :D
	// Prints our top field content
	hw.PrintPageInfo(page)

	// Prints our bottom field content
	hw.PrintControls()

	// We load the text file for the help page
	content, err := util.OpenFileAsText("assets/text/help/" + page.String() + ".txt")
	if err != nil {
		//hw.Error("Error opening file:" + err.Error())

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

func (hw *HelpWindow) HandleStateChange() {
	hw.ResetWindowDrawings()
	hw.ForceConsoleRefresh()
}

func (hw *HelpWindow) NumberOfHelpPages() int {
	return int(types.HelpPageIndex - 1) // The index is always the last page, we don't count it
}

func (hw *HelpWindow) IndexSelectionUp() {
	hw.indexMutex.Lock()
	defer hw.indexMutex.Unlock()
	if hw.indexSelection > 0 {
		hw.indexSelection--
	}
}

func (hw *HelpWindow) IndexSelectionDown() {
	hw.indexMutex.Lock()
	defer hw.indexMutex.Unlock()
	if hw.indexSelection < hw.NumberOfHelpPages() {
		hw.indexSelection++
	}
}
