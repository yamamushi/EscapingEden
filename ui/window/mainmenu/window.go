package login

import (
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"github.com/yamamushi/EscapingEden/ui/window"
	"sync"
)

// LoginWindow is a window for logins
type LoginWindow struct {
	window.Window
	credentials       *LoginCreds
	lwMutex           sync.Mutex
	windowState       LoginWindowState
	loginState        LoginState
	registrationState RegistrationState

	optionSelected     int
	currentOptionCount int
}

// LoginCreds is a struct for storing login credentials
type LoginCreds struct {
	Username string
	Hash     string
}

// LoginWindowState is an enum for storing login window state
type LoginWindowState int

const (
	LoginWindowMenu LoginWindowState = iota
	LoginWindowLogin
	LoginWindowRegister
)

// LoginState is an enum for storing login state
type LoginState int

const (
	LoginUsername LoginState = iota
	LoginPassword
	LoginSubmit
)

// RegistrationState is an enum for storing registration state
type RegistrationState int

const (
	RegistrationMain RegistrationState = iota
	RegistrationUsername
	RegistrationPassword
	RegistrationPasswordConfirm
	RegistrationEmail
	RegistrationDiscord
	RegistrationSubmit
)

// NewLoginWindow creates a new login window
func NewLoginWindow(x, y, width, height, consoleWidth, consoleHeight int, input, output chan string) *LoginWindow {
	lw := &LoginWindow{}
	lw.ID = config.WindowLoginMenu
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
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	lw.Width = width
	lw.Height = height
	lw.ConsoleWidth = consoleWidth
	lw.ConsoleHeight = consoleHeight
	lw.Bordered = true
	lw.ConsoleReceive = input
	lw.ConsoleSend = output
	lw.windowState = LoginWindowMenu
	lw.credentials = &LoginCreds{}
	return lw
}

// HandleInput handles input for the login window
func (lw *LoginWindow) HandleInput(input types.Input) {
	switch lw.windowState {
	case LoginWindowMenu:
		lw.handleMenuInput(input.Data)
	case LoginWindowLogin:
		lw.handleLoginInput(input.Data)
	case LoginWindowRegister:
		lw.handleRegistrationInput(input)
	}
}

// UpdateContents updates the contents of the login window
func (lw *LoginWindow) UpdateContents() {
	switch lw.windowState {
	case LoginWindowMenu:
		lw.drawMenu()
	case LoginWindowLogin:
		lw.drawLoginMenu()
	case LoginWindowRegister:
		lw.drawRegistrationMenu()
	}
}
