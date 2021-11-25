package ui

import (
	"log"
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
	RegistrationUsername RegistrationState = iota
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
	lw.ManagerSend = output
	lw.windowState = LoginWindowMenu

	lw.credentials = &LoginCreds{}
	return lw
}

func (lw *LoginWindow) HandleInput(input string) {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	if lw.GetActive() {
		log.Println("LoginWindow Handling input")
	}

	if input != "cherry" {
		lw.Error("Invalid input")
		return
	}
}

func (lw *LoginWindow) UpdateContents() {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()
	// First we are going to setup our default login screen
	lw.Contents = ""

	switch lw.windowState {
	case LoginWindowMenu:
		lw.Contents = lw.DrawMenu()
	case LoginWindowLogin:
		lw.Contents = "Login"
	case LoginWindowRegister:
		lw.Contents = "Register"
	}

	// parse the current state as a string
	switch lw.loginState {
	case LoginUsername:
		//lw.Contents += lw.CenterText(SetRGB(&ColorCode{0xAA, 0x00, 0x00}, "Username:"), 10)
	case LoginPassword:
		//lw.Contents += lw.CenterText(SetRGB(&ColorCode{0xAA, 0x00, 0x00}, "Password:"), 5)
	case LoginSubmit:
		//lw.Contents += lw.CenterText(SetRGB(&ColorCode{0xAA, 0x00, 0x00}, "Submitting..."), 5)
	}
}

func (lw *LoginWindow) DrawMenu() string {

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

	return output

}
