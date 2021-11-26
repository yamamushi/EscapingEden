package ui

import (
	"encoding/json"
	"log"
	"sync"
)

// Implements a Chat Window

type ChatWindow struct {
	Window
	History      []string
	HistoryIndex int
	cwMutex      sync.Mutex
}

func NewChatWindow(x, y, w, h int, input, output chan string) *ChatWindow {
	cw := new(ChatWindow)
	cw.ID = CHATBOX
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

	cw.ConsoleReceive = input
	cw.ConsoleSend = output

	cw.History = append(cw.History, "Hello World")
	cw.HistoryIndex = 0
	go cw.Listen()

	return cw
}

func (cw *ChatWindow) HandleInput(input Input) {
	cw.cwMutex.Lock()
	defer cw.cwMutex.Unlock()
	if cw.GetActive() {
		log.Println("ChatWindow Handling input")
	}

	if input.Data == "up" {
		log.Println("ChatWindow Handling input up")
	}
	if input.Data == "down" {
		log.Println("ChatWindow Handling input down")
	}
	if input.Data == "left" {
		log.Println("ChatWindow Handling input left")
	}
	if input.Data == "right" {
		log.Println("ChatWindow Handling input right")
	}

	// Send a console message to the ConsoleSend channel
	message := ConsoleMessage{Message: input.Data, Type: "chat"}
	output, err := json.Marshal(message)
	if err == nil {
		log.Println("Sending message on cw.ConsoleSend")
		cw.ConsoleSend <- string(output)
		log.Println("Message Sent")
	} else {
		log.Println(err.Error())
	}
}

// ConsoleMessage is called by console to manually write a console message to the history
func (cw *ChatWindow) ConsoleMessage(message string) {
	cw.cwMutex.Lock()
	defer cw.cwMutex.Unlock()
	cw.History = append(cw.History, message)
}

// Listens for any messages on cw.ReceiveMessages Chan and handles them
func (cw *ChatWindow) Listen() {
	for {
		select {
		case msg := <-cw.ConsoleReceive:
			log.Println("Chat window received message from console")

			message := ConsoleMessage{}
			err := json.Unmarshal([]byte(msg), &message)
			if err == nil {
				cw.cwMutex.Lock()
				cw.History = append(cw.History, message.Message)
				cw.cwMutex.Unlock()
			} else {
				log.Println(err.Error())
			}
		}
	}
}

func (cw *ChatWindow) UpdateContents() {
	cw.cwMutex.Lock()
	defer cw.cwMutex.Unlock()

	// only keep the newest 500 messages in cw.history
	if len(cw.History) > 500 {
		cw.History = cw.History[len(cw.History)-500:]
	}

	output := ""
	for i := 0; i < len(cw.History); i++ {
		output += cw.History[i] + "\n"
	}

	cw.SetContents(output)
}
