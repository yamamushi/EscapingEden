package ui

import (
	"log"
	"sync"
)

type LoginWindow struct {
	Window
	credentials *LoginCreds
	lwMutex     sync.Mutex
}

type LoginCreds struct {
	Username string
	Hash     string
}

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
	lw.ManagerReceive = input
	lw.ManagerSend = output

	lw.credentials = &LoginCreds{}
	return lw
}

func (lw *LoginWindow) HandleInput(input string) {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()
	if lw.GetActive() {
		log.Println("LoginWindow Handling input")
	}
}

func (lw *LoginWindow) UpdateContent() {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()
	// First we are going to setup our default login screen

	// This is an ascii art of the login screen
	lw.Contents = lw.CenterText(SetRGB(&ColorCode{0xAA, 0x00, 0x00}, "Hello"), 5)
}
