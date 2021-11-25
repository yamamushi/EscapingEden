package ui

import (
	"log"
	"sync"
	"time"
)

type ToolboxWindow struct {
	Window
	twMutex sync.Mutex
}

func NewToolboxWindow(x, y, w, h int, input, output chan string) *ToolboxWindow {
	lw := &ToolboxWindow{}
	lw.ID = TOOLBOX
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

	return lw
}

func (tw *ToolboxWindow) HandleInput(input string) {
	tw.twMutex.Lock()
	defer tw.twMutex.Unlock()

	if tw.GetActive() {
		log.Println("Toolbox Handling input")
	}

	if len(input) > 0 {
		log.Println(input[len(input)-1])
	}
}

func (tw *ToolboxWindow) UpdateContents() {
	tw.twMutex.Lock()
	defer tw.twMutex.Unlock()

	// current time with second accuracy as a string
	time := time.Now().Format("15:04:05")

	tw.SetContents("Current server time: " + time)
	log.Println("Toolbox contents updated")
}
