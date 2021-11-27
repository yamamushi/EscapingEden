package ui

import (
	"encoding/json"
	"log"
	"sync"
)

// Implements a Chat Window

type ChatWindow struct {
	Window
	History       []string
	HistoryIndex  int
	cwMutex       sync.Mutex
	cwInputBuffer string
}

func NewChatWindow(x, y, w, h, consoleWidth, consoleHeight int, input, output chan string) *ChatWindow {
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
	cw.ConsoleWidth = consoleWidth
	cw.ConsoleHeight = consoleHeight
	cw.Bordered = true

	cw.ConsoleReceive = input
	cw.ConsoleSend = output
	cw.ScrollingSupported = true

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

	switch input.Type {
	case InputUp:
		log.Println("ChatWindow Up")
		cw.DecreaseContentPos()
		return
	case InputDown:
		log.Println("ChatWindow Down")
		cw.IncreaseContentPos()
		return
	case InputLeft:
		log.Println("ChatWindow Left")
		cw.RequestPopupFromConsole(cw.ConsoleWidth/2-40, cw.ConsoleHeight/2-10, 150, 20, "This is a test of a really long string with a bunch of random content to see if the content buffer will scroll or not correctly")
		return
	case InputRight:
		log.Println("ChatWindow Right")
		return
	case InputNewline:
		log.Println("ChatWindow Newline")
		return
	case InputBackspace:
		log.Println("ChatWindow Backspace")
		// remove one character from the input buffer
		cw.cwInputBuffer = cw.cwInputBuffer[:len(cw.cwInputBuffer)-1]
		return
	case InputReturn:
		log.Println("ChatWindow Return")
		if cw.cwInputBuffer != "" {
			// Send a console message to the ConsoleSend channel
			message := ConsoleMessage{Message: cw.cwInputBuffer, Type: "chat"}
			output, err := json.Marshal(message)
			if err == nil {
				log.Println("Sending message on cw.ConsoleSend")
				cw.ConsoleSend <- string(output)
				log.Println("Message Sent")
			} else {
				log.Println(err.Error())
			}
			cw.cwInputBuffer = ""
		} else {
			log.Println("cw.cwInputBuffer is empty during InputReturn chat request")
		}
		return
	}

	log.Println("Chatwindow Receive: ", input.Data)
	cw.cwInputBuffer += input.Data

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
				// We know if our starting position is less than 0, and we append a new message, then there is
				// Content in the scroll buffer that has not been displayed yet.
				if cw.GetContentStartPos() < 0 {
					cw.DecreaseContentPos()
					cw.SetScrollBufferNew(true)
				} else {
					cw.SetScrollBufferNew(false)
				}

				log.Println("content start pos: ", cw.GetContentStartPos())
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
