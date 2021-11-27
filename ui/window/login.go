package window

import (
	"github.com/yamamushi/EscapingEden/ui/console"
	"log"
	"strings"
	"sync"
)

type LoginWindow struct {
	Window
	credentials       *LoginCreds
	lwMutex           sync.Mutex
	windowState       LoginWindowState
	loginState        LoginState
	registrationState RegistrationState

	optionSelected     int
	currentOptionCount int
}

type LoginCreds struct {
	Username string
	Hash     string
}

type LoginWindowState int

const (
	LoginWindowMenu LoginWindowState = iota
	LoginWindowLogin
	LoginWindowRegister
)

type LoginState int

const (
	LoginUsername LoginState = iota
	LoginPassword
	LoginSubmit
)

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

func NewLoginWindow(x, y, w, h, consoleWidth, consoleHeight int, input, output chan string) *LoginWindow {
	lw := &LoginWindow{}
	lw.ID = LOGINMENU
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
	lw.windowState = LoginWindowMenu
	lw.credentials = &LoginCreds{}
	return lw
}

func (lw *LoginWindow) HandleInput(input Input) {
	switch lw.windowState {
	case LoginWindowMenu:
		lw.handleMenuInput(input.Data)
	case LoginWindowLogin:
		lw.handleLoginInput(input.Data)
	case LoginWindowRegister:
		lw.handleRegistrationInput(input)
	}
}

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

func (lw *LoginWindow) drawMenu() {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	// First we are going to setup our default login screen
	lw.PrintLn(lw.X+10, lw.Y+5, "Welcome to Escaping Eden!", "")
	lw.PrintLn(lw.X+10, lw.Y+6, "Please select a menu option from below:", "")

	lw.PrintLn(lw.X+11, lw.Y+8, "(", "")
	lw.PrintLn(lw.X+12, lw.Y+8, "l", "\033[1m")
	lw.PrintLn(lw.X+13, lw.Y+8, ")ogin", "")

	lw.PrintLn(lw.X+11, lw.Y+9, "(", "")
	lw.PrintLn(lw.X+12, lw.Y+9, "r", "\033[1m")
	lw.PrintLn(lw.X+13, lw.Y+9, ")egister", "")

	lw.PrintLn(lw.X+11, lw.Y+10, "(", "")
	lw.PrintLn(lw.X+12, lw.Y+10, "q", "\033[1m")
	lw.PrintLn(lw.X+13, lw.Y+10, ")uit", "")
	return
}

func (lw *LoginWindow) handleMenuInput(input string) {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	if !lw.GetActive() {
		return
	}
	input = strings.ToLower(input)

	if input != "l" && input != "r" && input != "q" && input != "login" && input != "register" && input != "quit" {
		lw.Error("Invalid input received")
		return
	}

	input = input[:1]

	switch input {
	case "l":
		log.Println("Login selected")
		lw.windowState = LoginWindowLogin
		lw.loginState = LoginUsername
		lw.ResetWindowDrawings() // Whenever we switch to a different window state, we need to reset the drawings
		return
	case "r":
		log.Println("Register selected")
		lw.windowState = LoginWindowRegister
		lw.registrationState = RegistrationMain
		lw.ResetWindowDrawings()
		return
	case "q":
		log.Println("Quit selected")
		lw.Quit()
		return
	}
}

func (lw *LoginWindow) handleLoginInput(input string) {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	if !lw.GetActive() {
		return
	}

	input = strings.ToLower(input[:1])

	switch lw.loginState {
	case LoginUsername:
		lw.credentials.Username = input
		lw.loginState = LoginPassword
	case LoginPassword:
		lw.credentials.Hash = input
		lw.loginState = LoginSubmit
	case LoginSubmit:
		lw.ConsoleSend <- "login:" + lw.credentials.Username + ":" + lw.credentials.Hash
		lw.loginState = LoginUsername
	}
}

func (lw *LoginWindow) handleRegistrationInput(input Input) {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	if !lw.GetActive() {
		return
	}
	log.Println("Handling registration input")

	switch lw.registrationState {
	case RegistrationMain:
		switch input.Type {
		case InputCharacter:
			if input.Data == "r" {
				log.Println("Opening rules popup")
				lw.RequestPopupFromConsole(lw.ConsoleWidth/2-40, lw.ConsoleHeight/2-10, 100, 20, "This is a test of a really long string with a bunch of random content to see if the content buffer will scroll or not correctly")
			}
			return
		case InputLeft:
			log.Println("Left arrow pressed")
			lw.optionSelected = 1
			return
		case InputRight:
			log.Println("Right arrow pressed")
			lw.optionSelected = 2
			return
		case InputReturn:
			log.Println("Return pressed")
			if lw.optionSelected == 1 {
				lw.windowState = LoginWindowMenu
			}
			if lw.optionSelected == 2 {
				lw.registrationState = RegistrationUsername
			}
			lw.optionSelected = 0
			lw.ResetWindowDrawings() // Whenever we switch to a different window state, we need to reset the drawings
			return
		default:
			return
		}
		lw.registrationState = RegistrationUsername
	case RegistrationUsername:
		lw.registrationState = RegistrationPassword
	case RegistrationPassword:
		lw.registrationState = RegistrationPasswordConfirm
	case RegistrationPasswordConfirm:
		lw.registrationState = RegistrationEmail
	case RegistrationEmail:
		lw.registrationState = RegistrationDiscord
	case RegistrationDiscord:
		lw.registrationState = RegistrationSubmit
	case RegistrationSubmit:
		lw.ConsoleSend <- "register:" + lw.credentials.Username + ":" + lw.credentials.Hash
		lw.registrationState = RegistrationUsername
	}
}

func (lw *LoginWindow) drawLoginMenu() {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	lw.SetContents("Login")
}

func (lw *LoginWindow) drawRegistrationMenu() {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	lw.SetContents("Welcome to Escaping Eden! You are about to embark upon a journey into ")
	//lw.PrintLn(lw.Y+2, lw.X+2, "Welcome to Escaping Eden! You are about to embark upon a journey into ", "")

	// Bold the text for the back and continue buttons
	if lw.optionSelected == 1 {
		fg := console.RGBCode(0, 0, 0)
		bg := console.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+5, lw.Height+1, "<Back>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+5, lw.Height+1, "<Back>", "\033[1m")
	}
	if lw.optionSelected == 2 {
		fg := console.RGBCode(0, 0, 0)
		bg := console.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+lw.Width-15, lw.Height+1, "<Continue>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+lw.Width-15, lw.Height+1, "<Continue>", "\033[1m")
	}

}
