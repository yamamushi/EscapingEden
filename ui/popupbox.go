package ui

import (
	"log"
	"sync"
	"time"
)

type PopupBox struct {
	Window
	pbMutex sync.Mutex
}

func NewPopupBox(x, y, w, h int, input, output chan string) *PopupBox {
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
	pb.Bordered = true
	pb.ConsoleReceive = input
	pb.ConsoleSend = output

	return pb
}

func (pb *PopupBox) HandleInput(input string) {
	pb.pbMutex.Lock()
	defer pb.pbMutex.Unlock()

	if pb.GetActive() {
		log.Println("PopupBox Handling input")
	}

	if len(input) > 0 {
		log.Println(input[len(input)-1])
	}
}

func (pb *PopupBox) UpdateContents() {
	pb.pbMutex.Lock()
	defer pb.pbMutex.Unlock()

	// current time with second accuracy as a string
	time := time.Now().Format("15:04:05")

	pb.SetContents("Current server time: " + time)
}
