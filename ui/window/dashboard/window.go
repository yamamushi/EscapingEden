package dashboard

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
type DashboardWindow struct {
	window.Window
	dwMutex     sync.Mutex
	windowState DashboardState

	// Initialize the window
	dwInitialized bool

	// The current login states
	firstTimeLogin    bool
	lastCharacterName string
	lastCharacterID   string

	// Vars for navigation
	characterCreatorState CharacterCreatorState

	// Vars for CharacterCreator
	charColorOption                     int // 0 = red, 1 = green, 2 = blue
	charCreatorOptionSelected           int // 0 = none, 1 = name, 2 = color
	charColorOptionActive               bool
	charCreatorNavOptionSelected        int // 0 = none, 1 = cancel, 2 = submit
	charCreatorConfirmNavOptionSelected int // 0 = none, 1 = cancel, 2 = submit
	charCreatorName                     string
	charCreatorUsernameError            string
}

// LoginWindowState is an enum for storing login window state
type DashboardState int

const (
	DashboardMainMenu DashboardState = iota
	DashboardLogin
	DashboardCreateCharacter
	DashboardManageCharacters
	DashboardManageSettings
	DashboardLogout
)

type CharacterCreatorState int

const (
	CharacterCreatorDefaultNull = iota
	CharacterCreatorFirstTimeLoginWelcome
	CharacterCreatorCharacterDetails
	CharacterCreatorConfirmCharacter
	CharacterCreatorPending
)

// NewDashboardWindow creates a new login window
func NewDashboardWindow(x, y, width, height, consoleWidth, consoleHeight int, input, output chan messages.WindowMessage,
	log logging.LoggerType, term terminals.TerminalType) *DashboardWindow {
	dw := &DashboardWindow{}
	dw.Log = log
	dw.Terminal = term
	dw.ID = config.WindowUserDashboard
	// if x or y are less than 1 set them to 1
	if x < 1 {
		x = 1
	}
	if y < 1 {
		y = 1
	}
	dw.X = x
	dw.Y = y

	// if w or h are less than 1 set them to 1
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	dw.Width = width
	dw.Height = height
	dw.ConsoleWidth = consoleWidth
	dw.ConsoleHeight = consoleHeight
	dw.Bordered = true
	dw.ConsoleReceive = input
	dw.ConsoleSend = output
	dw.windowState = DashboardMainMenu
	return dw
}

// HandleInput handles input for the login window
func (dw *DashboardWindow) HandleInput(input types.Input) {
	switch dw.windowState {
	case DashboardMainMenu:
		dw.handleMenuInput(input)
	case DashboardCreateCharacter:
		dw.handleCreateCharacterMenuInput(input)
	case DashboardLogin:
		//dw.handleLoginInput(input)
	case DashboardManageCharacters:
		//dw.handleRegistrationInput(input)
	case DashboardManageSettings:
		//dw.handleSettingsInput(input)
	case DashboardLogout:
		//dw.handleLogoutInput(input)
	}
}

// UpdateContents updates the contents of the login window
func (dw *DashboardWindow) UpdateContents() {
	switch dw.windowState {
	case DashboardMainMenu:
		dw.drawMenu()
	case DashboardCreateCharacter:
		dw.drawCreateCharacterMenu()
	}
}
