package ui

import (
	"encoding/json"
	"log"
	"sync"
)

type PopupBox struct {
	Window
	pbMutex sync.Mutex
}

type PopupBoxConfig struct {
	X       int
	Y       int
	Width   int
	Height  int
	Content string
}

func PopupConfig(x, y, width, height int, content string) *PopupBoxConfig {
	return &PopupBoxConfig{x, y, width, height, content}
}

func (c *PopupBoxConfig) String() string {
	output, _ := json.Marshal(c)
	return string(output)
}

func NewPopupBox(x, y, w, h, consoleWidth, consoleHeight int, input, output chan string) *PopupBox {
	pb := &PopupBox{}
	pb.ID = POPUPBOX
	// if x or y are less than 1 set them to 1
	if x < 1 {
		x = 1
	}
	if y < 1 {
		y = 1
	}
	pb.X = x
	pb.Y = y

	// if w or h are less than 1 set them to 1
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	pb.Width = w
	pb.Height = h
	pb.ConsoleWidth = consoleWidth
	pb.ConsoleHeight = consoleHeight
	pb.Bordered = true
	pb.ConsoleReceive = input
	pb.ConsoleSend = output
	pb.ScrollingSupported = true

	return pb
}

func (pb *PopupBox) HandleInput(input Input) {
	pb.pbMutex.Lock()
	defer pb.pbMutex.Unlock()

	if pb.GetActive() {
		log.Println("PopupBox Handling input")
	}

	switch input.Type {
	case InputUp:
		log.Println("PopupBox Up")
		pb.DecreaseContentPos()
		return
	case InputDown:
		log.Println("PopupBox Down")
		pb.IncreaseContentPos()
		return
	case InputRight:
		log.Println("PopupBox Handling input right - attempting to close popup")
		message := ConsoleMessage{Type: "popupbox", Message: "close"}
		pb.ConsoleSend <- message.String()
		log.Println("PopupBox sent close message to console")
	}

}

func (pb *PopupBox) UpdateContents() {
	pb.pbMutex.Lock()
	defer pb.pbMutex.Unlock()
}
