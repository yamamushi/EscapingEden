package ui

import (
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

func NewLoginWindow(x, y, w, h int, input, output chan string) *LoginWindow {
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
		lw.handleRegistrationInput(input.Data)
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
	lw.SetContents("")

	artConvert := NewArtConvert()
	artWork, err := artConvert.OpenAt("./assets/ascii/menuIsland.txt", lw, 10, 10)
	if err != nil {
		log.Println(err)
	}
	output := artWork

	output += lw.PrintAt(10, 10, "Welcome to "+BoldText("Escaping Eden"))
	output += lw.PrintAt(12, 10, "Please select a menu option from below")

	output += lw.PrintAt(14, 10, "("+BoldText("l")+")login")
	output += lw.PrintAt(15, 10, "("+BoldText("r")+")register")
	output += lw.PrintAt(16, 10, "("+BoldText("q")+")quit")

	output += ResetStyle()
	lw.SetContents(output)
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
	case "r":
		log.Println("Register selected")
		lw.windowState = LoginWindowRegister
		lw.registrationState = RegistrationMain
	case "q":
		log.Println("Quit selected")
		lw.Quit()
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

func (lw *LoginWindow) handleRegistrationInput(input string) {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	if !lw.GetActive() {
		return
	}

	switch lw.registrationState {
	case RegistrationMain:
		lw.registrationState = RegistrationUsername
	case RegistrationUsername:
		lw.credentials.Username = input
		lw.registrationState = RegistrationPassword
	case RegistrationPassword:
		lw.credentials.Hash = input
		lw.registrationState = RegistrationPasswordConfirm
	case RegistrationPasswordConfirm:
		lw.credentials.Hash = input
		lw.registrationState = RegistrationEmail
	case RegistrationEmail:
		lw.credentials.Hash = input
		lw.registrationState = RegistrationDiscord
	case RegistrationDiscord:
		lw.credentials.Hash = input
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
}
