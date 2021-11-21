package ui

import (
	"sync"
)

// Implements a Chat Window

type ChatWindow struct {
	Window
	History      []string
	HistoryIndex int
	cwMutex      sync.Mutex
}

func NewChatWindow(x, y, w, h int) *ChatWindow {
	cw := new(ChatWindow)
	// if x or y are less than 1 set them to 1
	if x < 1 {
		x = 1
	}
	if y < 1 {
		y = 1
	}
	cw.X = x
	cw.Y = y

	// if w or h are less than 1 set them to 1
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	cw.Width = w
	cw.Height = h
	cw.Bordered = true

	cw.History = append(cw.History, "Hello World")
	cw.HistoryIndex = 0

	return cw
}

func (cw *ChatWindow) HandleInput(input string) {
	cw.cwMutex.Lock()
	defer cw.cwMutex.Unlock()

	cw.History = append(cw.History, input)
	cw.HistoryIndex = len(cw.History) - 1
}

func (cw *ChatWindow) UpdateContent() {
	cw.cwMutex.Lock()
	defer cw.cwMutex.Unlock()

	cw.SetContent(cw.History[cw.HistoryIndex])
}
