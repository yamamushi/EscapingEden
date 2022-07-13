package chat

import (
	"github.com/yamamushi/EscapingEden/edenutil"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/terminals"
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"github.com/yamamushi/EscapingEden/ui/util"
	"github.com/yamamushi/EscapingEden/ui/window"
	"sync"
	"time"
)

const CHAR_BUFFER_SIZE = 128

// Implements a Chat Window

// ChatWindow is a chat window
type ChatWindow struct {
	window.Window
	History       []string
	HistoryIndex  int
	cwMutex       sync.Mutex
	cwInputBuffer string
	ChatReceiver  chan messages.ChatMessage
}

// NewChatWindow creates a new chat window
func NewChatWindow(x, y, w, h, consoleWidth, consoleHeight int, chatInput chan messages.ChatMessage,
	windowInput, output chan messages.WindowMessage, log logging.LoggerType, term terminals.TerminalType) *ChatWindow {
	cw := new(ChatWindow)
	cw.Log = log
	cw.Terminal = term
	cw.ID = config.WindowChatBox
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

	cw.ChatReceiver = chatInput
	cw.ConsoleReceive = windowInput
	cw.ConsoleSend = output
	cw.ScrollingSupported = true

	// Initializing a default chat message

	cw.History = append(cw.History, "Welcome traveller!")
	cw.History = append(cw.History, "The current server time is: "+time.Now().Format("2006-01-02 15:04:05"))
	cw.History = append(cw.History, "The current time in Freeport is: "+edenutil.EdenTime.TimeStamp(edenutil.EdenTime{}))
	//cw.History = append(cw.History, "There are currently __ players online. ")
	cw.History = append(cw.History, "")
	cw.History = append(cw.History, "There are no current active world events.")
	cw.History = append(cw.History, "")

	cw.HistoryIndex = 0
	go cw.Listen()

	return cw
}

// HandleInput handles for the chat window
func (cw *ChatWindow) HandleInput(input types.Input) {
	cw.cwMutex.Lock()
	defer cw.cwMutex.Unlock()
	if !cw.GetActive() {
		// return
	}

	switch input.Type {
	case types.InputUp:
		//log.Println("ChatWindow Up")
		cw.DecreaseContentPos()
		cw.ResetWindowDrawings()
		cw.RequestFlushFromConsole()
		return
	case types.InputDown:
		//log.Println("ChatWindow Down")
		cw.IncreaseContentPos()
		cw.ResetWindowDrawings()
		cw.RequestFlushFromConsole()
		return
	case types.InputLeft:
		//log.Println("ChatWindow Left")
		cw.RequestPopupFromConsole(cw.ConsoleWidth/2-40, cw.ConsoleHeight/2-10, 50, 20, "This is a test of a really long string with a bunch of random content to see if the content buffer will scroll or not correctly")
		return
	case types.InputRight:
		//log.Println("ChatWindow Right")
		return
	case types.InputNewline:
		//log.Println("ChatWindow Newline")
		return
	case types.InputBackspace:
		//log.Println("ChatWindow Backspace")
		// remove one character from the input buffer if there is one
		if len(cw.cwInputBuffer) > 0 {
			cw.cwInputBuffer = cw.cwInputBuffer[:len(cw.cwInputBuffer)-1]
		}
		return
	case types.InputReturn:
		//log.Println("ChatWindow Return")
		if cw.cwInputBuffer != "" {
			// Send a console message to the ConsoleSend channel
			consoleMessage := messages.WindowMessage{Data: cw.cwInputBuffer, Type: messages.WM_ParseChat}
			cw.ConsoleSend <- consoleMessage
			cw.cwInputBuffer = ""
		}
		return
	case types.InputCharacter:
		cw.HandleCharacterInput(input.Data)
	}
	//log.Println("Chatwindow Receive: ", input.Data)
}

func (cw *ChatWindow) DrawInputLine() {

	colorCode := new(util.ColorCode)
	if len(cw.cwInputBuffer) >= CHAR_BUFFER_SIZE {
		colorCode = util.RGBCode(255, 0, 0)
	} else {
		colorCode = util.RGBCode(255, 255, 255)
	}
	// Clear the input line first with the correct color
	for i := 0; i < cw.Width; i++ {
		cw.PrintLn(cw.X+i, cw.Y+cw.Height, " ", colorCode.FG())
	}
	// Draw a > at the bottom of the chat window
	if len("> "+cw.cwInputBuffer)+2 > cw.Width {
		// Only show the last part of the input buffer that fits in the window
		// If the input buffer is 256 characters, we draw it in red to indicate that no more characters can be added
		cw.PrintLn(cw.X+1, cw.Y+cw.Height, "> "+cw.cwInputBuffer[len(cw.cwInputBuffer)-cw.Width+4:], colorCode.FG())
	} else {
		cw.PrintLn(cw.X+1, cw.Y+cw.Height, "> "+cw.cwInputBuffer, colorCode.FG())
	}
}

func (cw *ChatWindow) HandleCharacterInput(character string) {
	// We set a 256 character limit on the input buffer
	if len(cw.cwInputBuffer) < CHAR_BUFFER_SIZE {
		cw.cwInputBuffer += character
	}
}

func (cw *ChatWindow) SendChatMessage() {

}

// ConsoleMessage is called by console to manually write a console message to the history
func (cw *ChatWindow) ConsoleMessage(message string) {
	cw.cwMutex.Lock()
	defer cw.cwMutex.Unlock()
	cw.History = append(cw.History, message)
}

// Listen listens for any messages on cw.ReceiveMessages Chan and handles them
func (cw *ChatWindow) Listen() {
	for {
		select {
		case chatMessage := <-cw.ChatReceiver:
			//log.Println("Chat Window received chat message from console")
			cw.cwMutex.Lock()
			cw.History = append(cw.History, chatMessage.Content)
			// We know if our starting position is less than 0, and we append a new message, then there is
			// Content in the scroll buffer that has not been displayed yet.
			if cw.GetContentStartPos() < 0 {
				cw.DecreaseContentPos()
				cw.SetScrollBufferNew(true)
			} else {
				cw.SetScrollBufferNew(false)
				cw.ResetWindowDrawings()
				cw.RequestFlushFromConsole()
			}
			//log.Println("content start pos: ", cw.GetContentStartPos())
			cw.cwMutex.Unlock()
		}
	}
}

// UpdateContents updates the contents of the chat window
func (cw *ChatWindow) UpdateContents() {
	cw.cwMutex.Lock()
	defer cw.cwMutex.Unlock()
	// First clear the window
	cw.ResetWindowDrawings()

	// only keep the newest 500 messages in cw.history
	if len(cw.History) > 500 {
		cw.History = cw.History[len(cw.History)-500:]
	}

	output := ""
	for i := 0; i < len(cw.History); i++ {
		output += cw.History[i] + "\n"
	}

	cw.SetContents(output)
	cw.DrawInputLine()

}
