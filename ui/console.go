package ui

import (
	"encoding/json"
	"github.com/yamamushi/EscapingEden/ui/console"
	"github.com/yamamushi/EscapingEden/ui/window"
	"log"
	"strconv"
	"sync"
)

type Console struct {
	ConnectionID string // The ID of the connection using this console
	Height       int    // The height of the console
	Width        int    // The width of the console

	Windows            []window.WindowType // The list of windows that are currently in the console
	ConsoleCommands    string
	LastSentOutput     string
	LastActiveWindow   window.WindowType
	mutex              sync.Mutex
	Shutdown           bool
	consoleInitialized bool
	initConsole        sync.Once

	// Channels for communicating with ConnectionManager
	WindowMessages  chan string
	SendMessages    chan string
	ReceiveMessages chan string

	// Channels for communicating with windows
	LoginMessages    chan string
	ChatMessages     chan string
	ToolboxMessages  chan string
	PopupBoxMessages chan string

	escapeBuffer       string
	escapeSequence     bool
	returnSequence     bool
	forceScreenRefresh bool
}

// NewConsole creates a new console with no windows.
func NewConsole(height int, width int, connectionID string, outputChannel chan string) *Console {
	// Setup a new Chat Window and add it to the console at the bottom.
	receiver := make(chan string)
	windowMessages := make(chan string)
	loginMessages := make(chan string)
	chatMessages := make(chan string)
	toolboxMessages := make(chan string)
	popupBoxMessage := make(chan string)

	return &Console{Height: height, Width: width, ConnectionID: connectionID, SendMessages: outputChannel,
		ReceiveMessages: receiver, WindowMessages: windowMessages, LoginMessages: loginMessages,
		ChatMessages: chatMessages, ToolboxMessages: toolboxMessages, PopupBoxMessages: popupBoxMessage}
}

func (c *Console) Init() {

	c.consoleInitialized = false
	// First we setup our login window
	loginWindow := window.NewLoginWindow(0, 0, c.Width-50, c.Height-13, c.Width, c.Height, c.LoginMessages, c.WindowMessages)
	loginWindow.Init()
	c.AddWindow(loginWindow)

	// Next we setup our chat window
	chatWindow := window.NewChatWindow(0, c.Height-10, c.Width-50, c.Height, c.Width, c.Height, c.ChatMessages, c.WindowMessages)
	chatWindow.Init()
	c.AddWindow(chatWindow)

	// Then we add our toolbox last
	toolboxWindow := window.NewToolboxWindow(c.Width-48, 0, 50, c.Height, c.Width, c.Height, c.ToolboxMessages, c.WindowMessages)
	toolboxWindow.Init()
	c.AddWindow(toolboxWindow)

	go c.CaptureWindowMessages()
	go c.CaptureManagerMessages()

	c.SetActiveWindow(chatWindow) // Set our default active window to the login window, we will pass this to another
	// window after we log in.

	c.ConsoleCommands += c.HideCursor() + c.ResetTerminal() // + c.DrawPrompt() + c.MoveCursorToTopLeft()
}

func (c *Console) CaptureWindowMessages() {
	for {
		select {
		case message := <-c.WindowMessages:
			log.Println("Client received window message")
			consoleMessage := &console.ConsoleMessage{}
			err := json.Unmarshal([]byte(message), consoleMessage)
			if err != nil {
				log.Println("Error unmarshalling consoleMessage: ", err)
				continue
			}
			log.Println("MessageType: ", consoleMessage.Type)
			switch consoleMessage.Type {
			case "console":
				switch consoleMessage.Message {
				case "popup":
					options := &window.PopupBoxConfig{}
					err := json.Unmarshal([]byte(consoleMessage.Options), options)
					if err != nil {
						log.Println("Error unmarshalling popup box options: ", err)
						continue
					} else {
						log.Println("Popup box options: ", options.String())
						c.OpenPopup(options)
					}
				}
			case "popupbox":
				log.Println("Popup box message: ", consoleMessage.Message)
				c.HandlePopupMessage(consoleMessage)
				continue
			case "error":
				log.Println("Error message: ", consoleMessage.Message)
				consoleMessage.RecipientID = c.ConnectionID
				c.SendMessages <- consoleMessage.String()
				continue
			case "quit":
				log.Println("Sending Quit request to ConnectionManager")
				consoleMessage.RecipientID = c.ConnectionID
				c.SendMessages <- consoleMessage.String()
				continue
			case "chat":
				log.Println("Chat message: ", consoleMessage.Message)
				consoleMessage.SenderID = c.ConnectionID
				c.SendMessages <- consoleMessage.String()
				continue
			}
		}
	}
}

func (c *Console) CaptureManagerMessages() {
	for {
		select {
		case message := <-c.ReceiveMessages:
			log.Println("Console received message from manager")
			consoleMessage := &console.ConsoleMessage{}
			err := json.Unmarshal([]byte(message), consoleMessage)
			if err != nil {
				log.Println("Error unmarshalling consoleMessage: ", err)
				continue
			}

			switch consoleMessage.Type {
			case "chat":
				log.Println("Chat message received from manager")
				c.ChatMessages <- consoleMessage.String()
			case "error":
				log.Println("Error message received from manager")
				consoleMessage.Message = "Error: " + consoleMessage.Message
				c.ChatMessages <- consoleMessage.String()
			case "broadcast":
				log.Println("Broadcast message received from manager")
				consoleMessage.Message = "Server Message: " + consoleMessage.Message
				c.ChatMessages <- consoleMessage.String()
			case "quit":
				log.Println("Quit message received from manager")
				c.SetShutdown(true)
			}
		}
	}
}

// SetManagerSendChannel sets the channel for sending messages to the ConnectionManager.
func (c *Console) SetManagerSendChannel(ch chan string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.SendMessages = ch
}

// SetManagerReceiveChannel sets the channel for receiving messages from the ConnectionManager.
func (c *Console) SetManagerReceiveChannel(ch chan string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.ReceiveMessages = ch
}

// AddWindow adds a window to the console if it is not already in the console by ID.
func (c *Console) AddWindow(w window.WindowType) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, window := range c.Windows {
		if window.GetID() == w.GetID() {
			log.Println("duplicate window: ", window.GetID())
			return
		}
	}
	c.Windows = append(c.Windows, w)
}

// RemoveWindow removes a window from the console by ID.
func (c *Console) RemoveWindow(id int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i, window := range c.Windows {
		if window.GetID() == id {
			c.Windows = append(c.Windows[:i], c.Windows[i+1:]...)
			return
		}
	}
}

// Draw returns the console as a byte array.
func (c *Console) Draw() []byte {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var s string
	//s = s + c.ConsoleCommands
	//c.ConsoleCommands = ""

	if !c.consoleInitialized {
		c.consoleInitialized = true
		return []byte(s)
	}

	if c.forceScreenRefresh {
		c.forceScreenRefresh = false
		s = s + c.ClearTerminal()
		return []byte(s)
	}

	for _, window := range c.Windows {
		if !window.GetHidden() {
			windowDraw := c.DrawWindow(window)
			if windowDraw != "" {
				s = s + windowDraw
			}
		}
	}

	//return []byte(s)

	// If the last output was not the same as the current output, we send it to the client and update the last output.
	if c.LastSentOutput != s && s != "" {
		log.Println("Sending new output to client, length:", len(s))
		c.LastSentOutput = s
		return []byte(s)
	} else {
		return []byte("")
	}
}

// HandleInput accepts a string terminated by a newline and processes it.
func (c *Console) HandleInput(rawInput byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if rawInput == 0 {
		// Ignore these null bytes
		return
	}

	if rawInput == 18 {
		// ctrl-r to force a screen refresh
		for _, w := range c.Windows {
			w.FlushLastSent()
		}
		c.forceScreenRefresh = true
		return
	}

	// Captures things like the arrow keys.
	if rawInput == '\033' {
		c.escapeBuffer = "\\033" // Just used for logging the escape sequence buffer, nothing else
		c.escapeSequence = true  // Lets us know that the next few bytes are going to be related to the escape sequence
		return
	}

	// If we have an active escape sequence, we continue parsing it.
	if c.escapeSequence {
		// If we have a [, we know we are starting a new escape sequence.
		if rawInput == '[' {
			c.escapeBuffer += "["
			return
		} else if c.escapeBuffer == "\\033[" {
			// If our escape buffer has an escape sequence, we know we are still parsing it.
			c.escapeBuffer += string(rawInput)
			switch rawInput {
			case 'A':
				c.InputToActiveWindow(window.Input{Type: window.InputUp})
			case 'B':
				c.InputToActiveWindow(window.Input{Type: window.InputDown})
			case 'C':
				c.InputToActiveWindow(window.Input{Type: window.InputRight})
			case 'D':
				c.InputToActiveWindow(window.Input{Type: window.InputLeft})
			default:
				log.Println("Unknown escape sequence: ", c.escapeBuffer)
			}
			c.escapeBuffer = ""
			c.escapeSequence = false
			return
		}
		c.escapeBuffer = ""
		c.escapeSequence = false
	}

	// If we have a backspace, we remove the last character from the input buffer.
	if rawInput == '\b' || rawInput == '\x7f' {
		c.InputToActiveWindow(window.Input{Type: window.InputBackspace})
		return
	}
	if rawInput == '\r' {
		c.InputToActiveWindow(window.Input{Type: window.InputReturn})
		return
	}
	if rawInput == '\t' {
		c.InputToActiveWindow(window.Input{Type: window.InputTab})
		return
	}
	if rawInput == '\n' {
		c.InputToActiveWindow(window.Input{Type: window.InputNewline})
		return
	}

	c.InputToActiveWindow(window.Input{Type: window.InputCharacter, Data: string(rawInput)})
}

func (c *Console) InputToActiveWindow(input window.Input) {
	for _, window := range c.Windows {
		if window.GetActive() {
			window.HandleInput(input)
			log.Println("Input Handled on window: ", window.GetID())
			return
		}
	}
}

// GetShutdown returns the shutdown status of the console.
func (c *Console) GetShutdown() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.Shutdown
}

// SetShutdown sets the shutdown status of the console.
func (c *Console) SetShutdown(status bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Shutdown = status
}

// DisableCursor disables the cursor.

// Moves the cursor to the top left corner of the console
func (c *Console) MoveCursorToTopLeft() string {
	return "\033[1;1H"
}

// Moves the cursor to the bottom left corner of the console
func (c *Console) MoveCursorToBottomLeft() string {
	return "\033[" + strconv.Itoa(c.Height) + ";0H"
}

// DrawPrompt returns the prompt for the console
func (c *Console) DrawPrompt() string {
	output := c.MoveCursorToBottomLeft()
	return output + "> "
}

// ScrollLock locks the scroll
func (c *Console) ScrollLock() string {
	return "\033[?1049h"
}

// ScrollUnlock unlocks the scroll
func (c *Console) ScrollUnlock() string {
	return "\033[?1049l"
}

// ClearTerminal
func (c *Console) ClearTerminal() string {
	return "\033[2J\n"
}

// ClearNotPrompt will clear each line of the console except the prompt
func (c *Console) ClearNotPrompt() string {
	var s string
	// save cursor position
	//s = c.SaveCursor()
	for i := 0; i < c.Height; i++ {
		// Move cursor to line i
		s = s + "\033[" + strconv.Itoa(i+1) + ";0H"
		// Clear line
		s = s + "\033[2K"
	}
	return s
}

func (c *Console) SaveCursor() string {
	return "\033[s"
}

func (c *Console) RestoreCursor() string {
	return "\033[u"
}

func (c *Console) HideCursor() string {
	return "\033[?25l"
}

func (c *Console) ShowCursor() string {
	return "\033[?25h"
}

// HardClear Terminal clears each line individually for height of console
func (c *Console) HardClear() string {
	var s string
	for i := 0; i < c.Height; i++ {
		// Move cursor to line i
		s = s + "\033[" + strconv.Itoa(i+1) + ";0H"
		// Clear line
		s = s + "\033[2K"
	}
	return s
}

// ResetTerminal
func (c *Console) ResetTerminal() string {
	return "\033c"
}

// GetWindowAttrs Takes window as an argument and returns the x,y position and visible height and length of the window
func (c *Console) GetWindowAttrs(window window.WindowType) (X int, Y int, visibleLength int, visibleHeight int) {
	if (window.GetWidth() + window.GetX()) > c.Width-2 {
		visibleLength = c.Width - window.GetX() - 2
	} else {
		visibleLength = window.GetWidth() - window.GetX()
	}

	if (window.GetHeight() + window.GetY()) > c.Height-2 {
		visibleHeight = c.Height - window.GetY() - 2
	} else {
		visibleHeight = window.GetHeight() - window.GetY()
	}
	return window.GetX(), window.GetY(), visibleLength, visibleHeight
}

// DrawWindow takes a WindowType as an argument and draws the content of the window within the window border
func (c *Console) DrawWindow(window window.WindowType) (content string) {
	// Get Window Attrs
	winX, winY, visibleLength, visibleHeight := c.GetWindowAttrs(window)

	// First we want to clear the window for any new content coming in
	window.ClearMap(winX, winY, visibleLength, visibleHeight, 0, 0)

	// If this is the active window we're going to force redraw it every request
	//if window.GetActive() {
	//	window.FlushLastSent()
	//}

	// Now we want to check for any new content updates
	window.UpdateContents()

	// Draw the contents of the window
	window.Draw(winX, winY, visibleLength, visibleHeight, 0, 0)

	// Now we want to draw the window border
	window.DrawBorder(winX, winY, visibleLength+1, visibleHeight+1)

	// Now we get the window's content as a string from it's PointMap
	content = window.PointMapToString()
	return content
}

// SetActiveWindow sets the active window and sets all other windows to inactive
func (c *Console) SetActiveWindow(window window.WindowType) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i, w := range c.Windows {
		if w.GetID() == window.GetID() {
			log.Println("Active window set to: ", w.GetID())
			w.SetActive(true)
			c.Windows = append(c.Windows[:i], c.Windows[i+1:]...)
			c.Windows = append(c.Windows, window)
		} else {
			w.SetActive(false)
		}
	}

}

func (c *Console) GetActiveWindow() window.WindowType {
	for _, w := range c.Windows {
		if w.GetActive() {
			return w
		}
	}
	return nil
}

func (c *Console) OpenPopup(options *window.PopupBoxConfig) {
	//popupBox := NewPopupBox(c.Width/2-40, c.Height/2-10, 80, 20, c.PopupBoxMessages, c.WindowMessages)
	log.Println(options)
	popupBox := window.NewPopupBox(options.X, options.Y, options.Width, options.Height, c.Width, c.Height, c.PopupBoxMessages, c.WindowMessages)
	popupBox.Init()
	popupBox.SetContents(options.Content)
	c.LastActiveWindow = c.GetActiveWindow() // Save the last active window
	c.AddWindow(popupBox)                    // Add the popup to the list of windows
	c.SetActiveWindow(popupBox)              // Set the popup as the active window
	//popupBox.FlushLastSent()
	popupBox.ClearMap(popupBox.X, popupBox.Y, options.Width, options.Height, 0, 0)
}

func (c *Console) ClosePopup() {
	// Loop through windows and remove the popup
	for _, w := range c.Windows {
		if w.GetID() == window.POPUPBOX {
			c.RemoveWindow(w.GetID())
			c.SetActiveWindow(c.LastActiveWindow)
			break
		}
	}
	// When we close the popup, our window is all garbage so we're going to force a redraw on everything
	c.ForceRedraw()
}

func (c *Console) HandlePopupMessage(message *console.ConsoleMessage) {
	switch message.Message {
	case "close":
		c.ClosePopup()
	}

}

func (c *Console) ForceRedraw() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, w := range c.Windows {
		w.FlushLastSent()
	}
	c.forceScreenRefresh = true
}
