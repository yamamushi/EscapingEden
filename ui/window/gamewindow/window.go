package gamewindow

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/terminals"
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"github.com/yamamushi/EscapingEden/ui/window"
	"sync"
)

// DashboardWindow is a window for users to login as a character, create a new one, manage their settings or log out.
type GameWindow struct {
	window.Window
	dwMutex     sync.Mutex
	windowState GameWindowState

	// Initialize the window
	gwInitialized bool

	// Vars for navigation
	characterCreatorState CharacterCreatorState
}

// GameWindowState is an enum for storing game window state
type GameWindowState int

const (
	GW_NullState GameWindowState = iota
	GW_DefaultView
)

type CharacterCreatorState int

// NewGameWindow creates a new login window
func NewGameWindow(x, y, width, height, consoleWidth, consoleHeight int, input, output chan messages.WindowMessage,
	log logging.LoggerType, term terminals.TerminalType) *GameWindow {
	gw := &GameWindow{}
	gw.Log = log
	gw.Terminal = term
	gw.ID = config.WindowGameDisplay
	// if x or y are less than 1 set them to 1
	if x < 1 {
		x = 1
	}
	if y < 1 {
		y = 1
	}
	gw.X = x
	gw.Y = y

	// if w or h are less than 1 set them to 1
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	gw.Width = width
	gw.Height = height
	gw.ConsoleWidth = consoleWidth
	gw.ConsoleHeight = consoleHeight
	gw.Bordered = true
	gw.ConsoleReceive = input
	gw.ConsoleSend = output
	gw.windowState = GW_DefaultView
	return gw
}

// HandleInput handles input for the login window
func (gw *GameWindow) HandleInput(input types.Input) {
	switch gw.windowState {
	case GW_DefaultView:
		// Handle input for the main view
	}
}

// UpdateContents updates the contents of the login window
func (gw *GameWindow) UpdateContents() {
	switch gw.windowState {
	case GW_DefaultView:
		gw.PrintLn(gw.X+2, gw.Y+2, "Game Window", gw.Terminal.Bold())

		// At center of window draw an @
		gw.PrintLn(gw.X+gw.Width/2, gw.Y+gw.Height/2, "@", gw.CharacterInfo.FGColor.FG()+gw.CharacterInfo.BGColor.BG())
	}
}
