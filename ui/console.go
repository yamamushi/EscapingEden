package ui

import (
	"encoding/json"
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"github.com/yamamushi/EscapingEden/ui/window"
	"github.com/yamamushi/EscapingEden/ui/window/chat"
	"github.com/yamamushi/EscapingEden/ui/window/help"
	"github.com/yamamushi/EscapingEden/ui/window/login"
	"github.com/yamamushi/EscapingEden/ui/window/popupbox"
	"github.com/yamamushi/EscapingEden/ui/window/toolbox"
	"log"
	"strconv"
	"sync"
)

const (
	MINWIDTH  = 155
	MINHEIGHT = 30
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
	abortSend          bool
	abortSync          sync.Mutex
	resizeActive       bool
	userLoggedIn       bool
}

// NewConsole creates a new console with no windows.
func NewConsole(height int, width int, connectionID string, outputChannel chan string) *Console {
	// Set up a new Chat Window and add it to the console at the bottom.
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

// Init is called once to initialize the console, it does things like create the default windows, and launches the
// Capture Message goroutines. It also sets up some default console commands that should be sent on first connect.
// These can be appended to later if we want to send a quick console command, but they are flushed after every write.
// Ie they are single use.
func (c *Console) Init() {

	c.consoleInitialized = false
	// First we set up our login window
	loginWindow := login.NewLoginWindow(0, 0, c.Width-50, c.Height-13, c.Width, c.Height, c.LoginMessages, c.WindowMessages)
	loginWindow.Init()
	c.AddWindow(loginWindow)

	// Next we set up our chat window
	chatWindow := chat.NewChatWindow(0, c.Height-10, c.Width-50, 9, c.Width, c.Height, c.ChatMessages, c.WindowMessages)
	chatWindow.Init()
	c.AddWindow(chatWindow)

	// Then we add our toolbox last
	toolboxWindow := toolbox.NewToolboxWindow(c.Width-48, 0, 48, c.Height-2, c.Width, c.Height, c.ToolboxMessages, c.WindowMessages)
	toolboxWindow.Init()
	c.AddWindow(toolboxWindow)

	go c.CaptureWindowMessages()
	go c.CaptureManagerMessages()

	c.SetActiveWindow(loginWindow) // Set our default active window to the login window, we will pass this to another
	// window after we log in.

	c.ConsoleCommands += c.HideCursor() + c.ResetTerminal() // + c.DrawPrompt() + c.MoveCursorToTopLeft()
}

// CaptureWindowMessages is a goroutine that listens for messages from the windows and parses them to determine
// Where they should go, or if any action should be taken from them. Launching a popup box, or sending a message
// or even quitting the session. There are many types of messages the console may have to parse from a window.
func (c *Console) CaptureWindowMessages() {
	for {
		select {
		case message := <-c.WindowMessages:
			log.Println("Client received window message")
			consoleMessage := &types.ConsoleMessage{}
			err := json.Unmarshal([]byte(message), consoleMessage)
			if err != nil {
				log.Println("Error unmarshalling consoleMessage: ", err)
				continue
			}
			log.Println("MessageType: ", consoleMessage.Type)
			switch consoleMessage.Type {
			case "console":
				log.Println("Received console message")
				switch consoleMessage.Message {
				case "popup":
					options := &config.WindowConfig{}
					err = json.Unmarshal([]byte(consoleMessage.Options), options)
					if err != nil {
						log.Println("Error unmarshalling popup box options: ", err)
						continue
					} else {
						log.Println("Popup box options: ", options.String())
						c.OpenPopup(options)
					}
				case "help":
					options := &config.WindowConfig{}
					err = json.Unmarshal([]byte(consoleMessage.Options), options)
					if err != nil {
						log.Println("Error unmarshalling help box options: ", err)
						continue
					} else {
						log.Println("Help box options: ", options.String())
						c.ToggleHelp(options)
					}
				case "refresh":
					log.Println("Refresh message received")
					//if !c.IsPopupOpen() {
					c.AbortSend()
					c.ForceRedraw()
					//}
				}

			case "popupbox":
				log.Println("Popup box message: ", consoleMessage.Message)
				c.HandlePopupMessage(consoleMessage)
				continue
			case "help":
				log.Println("Help message: ", consoleMessage.Message)
				c.HandleHelpMessage(consoleMessage)
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

// CaptureManagerMessages is a goroutine that listens for messages from the ConnectionManager and parses them to determine
// Where they should go, or if any action should be taken from them.
func (c *Console) CaptureManagerMessages() {
	for {
		select {
		case message := <-c.ReceiveMessages:
			log.Println("Console received message from manager")
			consoleMessage := &types.ConsoleMessage{}
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

	for _, target := range c.Windows {
		if target.GetID() == w.GetID() {
			log.Println("duplicate window: ", target.GetID())
			return
		}
	}
	c.Windows = append(c.Windows, w)
}

// RemoveWindow removes a window from the console by ID.
func (c *Console) RemoveWindow(id config.WindowID) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i, target := range c.Windows {
		if target.GetID() == id {
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
	if !c.IsConsoleValidSize() {
		s = s + "\033[2J"
		s = s + "Invalid console size, Escaping Eden requires a terminal size of" + strconv.Itoa(MINWIDTH) + "x" + strconv.Itoa(MINHEIGHT) + "or greater.\r\n"
		s = s + "Please resize your terminal, or press q to disconnect.\n"
		s = s + "If your terminal is empty after resizing, you can press ctrl-r to force a screen refresh.\n"

		if c.LastSentOutput != s {
			c.LastSentOutput = s
			return []byte(s)
		} else {
			return []byte("")
		}
	}

	if !c.consoleInitialized {
		c.consoleInitialized = true
		return []byte(s)
	}

	if c.forceScreenRefresh {
		c.forceScreenRefresh = false
		s = s + c.ClearTerminal()
		return []byte(s)
	} else if c.resizeActive {
		log.Println("Handling resize in buffer")
		for _, w := range c.Windows {
			w.FlushLastSent()
		}
		c.resizeActive = false
	}

	//log.Println("Drawing console")
	for _, target := range c.Windows {
		if !target.GetHidden() {
			//log.Println("Drawing window: ", window.GetID())
			windowDraw := c.DrawWindow(target)
			if windowDraw != "" {
				s = s + windowDraw
			}
		}
	}

	//return []byte(s)
	// We do a last minute check for aborting sending messages
	// This is useful when a screen has asked for a screen refresh
	// And we don't want to send data that is going to get overwritten immediately
	c.abortSync.Lock()
	defer c.abortSync.Unlock()
	// If the last output was not the same as the current output, we send it to the client and update the last output.
	if c.LastSentOutput != s && s != "" && !c.abortSend {
		log.Println("Sending new output to client, length:", len(s))
		c.LastSentOutput = s
		return []byte(s)
	} else {
		c.abortSend = false
		return []byte("")
	}
}

// HandleInput accepts a string terminated by a newline and processes it.
func (c *Console) HandleInput(rawInput byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if !c.IsConsoleValidSize() {
		// If the console isn't valid we don't want to accept any input
		// However if we receive the letter q, we will exit the program
		if rawInput == 'q' {
			c.Shutdown = true
		}
		return
	}
	if rawInput == 0 {
		// Ignore these null bytes
		return
	}

	//log.Println("Console received input: ", int(rawInput))

	if rawInput == 8 {
		options := &config.WindowConfig{X: c.Width/2 - 40, Y: c.Height/2 - 10, Width: 100, Height: 20, Page: 0}
		go c.ToggleHelp(options)
		return
	}

	if rawInput == 18 {
		// ctrl-r to force a screen refresh
		for _, w := range c.Windows {
			w.ResetWindowDrawings()
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
		// If we have a [ symbol, we know we are starting a new escape sequence, there is no ]
		if rawInput == '[' {
			c.escapeBuffer += "["
			return
		} else if c.escapeBuffer == "\\033[" {
			// If our escape buffer has an escape sequence, we know we are still parsing it.
			c.escapeBuffer += string(rawInput)
			switch rawInput {
			case 'A':
				c.InputToActiveWindow(types.Input{Type: types.InputUp})
			case 'B':
				c.InputToActiveWindow(types.Input{Type: types.InputDown})
			case 'C':
				c.InputToActiveWindow(types.Input{Type: types.InputRight})
			case 'D':
				c.InputToActiveWindow(types.Input{Type: types.InputLeft})
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
		c.InputToActiveWindow(types.Input{Type: types.InputBackspace})
		return
	}
	if rawInput == '\r' {
		c.InputToActiveWindow(types.Input{Type: types.InputReturn})
		return
	}
	if rawInput == '\t' {
		if !c.IsPopupOpen() {
			c.SetNextActiveWindow()
			for _, w := range c.Windows {
				w.ResetWindowDrawings()
				//w.FlushLastSent()
			}
			c.forceScreenRefresh = true
		}
		return
	}
	if rawInput == '\n' {
		c.InputToActiveWindow(types.Input{Type: types.InputNewline})
		return
	}

	c.InputToActiveWindow(types.Input{Type: types.InputCharacter, Data: string(rawInput)})

}

// InputToActiveWindow sends an input to the active window.
func (c *Console) InputToActiveWindow(input types.Input) {
	for _, target := range c.Windows {
		if target.GetActive() {
			target.HandleInput(input)
			log.Println("Input Handled on window: ", target.GetID())
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

// MoveCursorToTopLeft Moves the cursor to the top left corner of the console
func (c *Console) MoveCursorToTopLeft() string {
	return "\033[1;1H"
}

// MoveCursorToBottomLeft Moves the cursor to the bottom left corner of the console
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

// ClearTerminal clears the terminal using the escape sequence
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

// SaveCursor saves the cursor position using the escape sequence
func (c *Console) SaveCursor() string {
	return "\033[s"
}

// RestoreCursor restores the cursor position using the escape sequence
func (c *Console) RestoreCursor() string {
	return "\033[u"
}

// HideCursor sets the cursor position using the escape sequence
func (c *Console) HideCursor() string {
	return "\033[?25l"
}

// ShowCursor sets the cursor position using the escape sequence
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

// ResetTerminal resets the terminal using the escape sequence
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
	window.Draw(winX, winY)

	// Now we want to draw the window border
	window.DrawBorder(winX, winY)

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
			//log.Println("Active window set to: ", w.GetID())
			w.SetActive(true)
			// We do this to move the window to the end of the slice
			// Since the last one will always be drawn last, ensuring it will be on top
			// of all other drawn windows
			c.Windows = append(c.Windows[:i], c.Windows[i+1:]...)
			c.Windows = append(c.Windows, window)
			//log.Println("Active Window: ", c.Windows[len(c.Windows)-1].GetID())
		} else {
			w.SetActive(false)
		}
	}
}

// GetActiveWindow returns the current active window
func (c *Console) GetActiveWindow() window.WindowType {
	for _, w := range c.Windows {
		if w.GetActive() {
			return w
		}
	}
	return nil
}

// OpenPopup opens a new popup window using the options
func (c *Console) OpenPopup(options *config.WindowConfig) {
	//log.Println(options)
	popupBox := popupbox.NewPopupBox(options.X, options.Y, options.Width, options.Height, c.Width, c.Height, c.PopupBoxMessages, c.WindowMessages)
	popupBox.Init()
	popupBox.SetContents(options.Content)
	c.LastActiveWindow = c.GetActiveWindow() // Save the last active window
	c.AddWindow(popupBox)                    // Add the popup to the list of windows
	c.SetActiveWindow(popupBox)              // Set the popup as the active window
	//popupBox.FlushLastSent()
	popupBox.FlushLastSent()
}

func (c *Console) ClosePopup() {
	// Loop through windows and remove the popup
	for _, w := range c.Windows {
		if w.GetID() == config.WindowPopupBox {
			c.RemoveWindow(w.GetID())
			c.SetActiveWindow(c.LastActiveWindow)
			break
		}
	}
	// When we close the popup, our window is all garbage so, we're going to force a re-draw on everything
	c.ForceRedraw()
}

func (c *Console) HandlePopupMessage(message *types.ConsoleMessage) {
	switch message.Message {
	case "close":
		c.ClosePopup()
	}
}

func (c *Console) IsPopupOpen() bool {
	for _, w := range c.Windows {
		if w.GetID() == config.WindowPopupBox {
			return true
		}
	}
	return false
}

// OpenPopup opens a new popup window using the options
func (c *Console) ToggleHelp(options *config.WindowConfig) {
	if !c.IsHelpOpen() {
		helpWindow := help.NewHelpWindow(options.X, options.Y, options.Width, options.Height, c.Width, c.Height, options.Page, c.PopupBoxMessages, c.WindowMessages)
		helpWindow.Init()
		helpWindow.SetContents(options.Content)
		c.LastActiveWindow = c.GetActiveWindow() // Save the last active window
		c.AddWindow(helpWindow)                  // Add the popup to the list of windows
		c.SetActiveWindow(helpWindow)            // Set the popup as the active window
		//popupBox.FlushLastSent()
		helpWindow.FlushLastSent()
	} else {
		c.CloseHelp()
	}
}

func (c *Console) CloseHelp() {
	// Loop through windows and remove the popup
	for _, w := range c.Windows {
		if w.GetID() == config.WindowHelpBox {
			c.RemoveWindow(w.GetID())
			c.SetActiveWindow(c.LastActiveWindow)
			break
		}
	}
	// When we close the popup, our window is all garbage so, we're going to force a re-draw on everything
	c.ForceRedraw()
}

func (c *Console) HandleHelpMessage(message *types.ConsoleMessage) {
	switch message.Message {
	case "close":
		c.CloseHelp()
	}
}

func (c *Console) IsHelpOpen() bool {
	for _, w := range c.Windows {
		if w.GetID() == config.WindowHelpBox {
			return true
		}
	}
	return false
}

func (c *Console) ForceRedraw() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	log.Println("Forcing redraw")
	for _, w := range c.Windows {
		w.FlushLastSent()
	}
	c.forceScreenRefresh = true
}

func (c *Console) ForceRedrawOn(windowType config.WindowID) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	log.Println("Forcing redraw on: ", windowType)
	for _, w := range c.Windows {
		if w.GetID() == windowType {
			log.Println("Flushing: ", w.GetID())
			w.FlushLastSent()
		}
	}
	//c.forceScreenRefresh = true
}

func (c *Console) AbortSend() {
	c.abortSync.Lock()
	defer c.abortSync.Unlock()
	c.abortSend = true
}

func (c *Console) SetNextActiveWindow() {
	// Set the active window to the first one in the list, because we know the last one is
	// Always the active one
	c.SetActiveWindowNoThread(c.Windows[0])
}

func (c *Console) SetPrevActiveWindow() {
	// Set the active window to the second to last one in the last, because the last one is
	// Always the active one

	// If the index is less than 2, then we only have one window in the list
	// In which case we don't want to do anything
	if len(c.Windows) < 2 {
		return
	}
	c.SetActiveWindowNoThread(c.Windows[len(c.Windows)-2])
}

// SetActiveWindowNoThread sets the active window and sets all other windows to inactive without locking
func (c *Console) SetActiveWindowNoThread(window window.WindowType) {
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

// HandleResize handles the resize event for all windows
func (c *Console) HandleResize(newWidth, newHeight int) {
	c.mutex.Lock()
	c.mutex.Unlock()

	c.Width = newWidth
	c.Height = newHeight
	for _, w := range c.Windows {
		switch w.GetID() {
		case config.WindowChatBox:
			w.UpdateParams(0, c.Height-10, c.Width-50, 9, c.Width, c.Height)
		case config.WindowLoginMenu:
			w.UpdateParams(0, 0, c.Width-50, c.Height-13, c.Width, c.Height)
		case config.WindowToolBox:
			w.UpdateParams(c.Width-48, 0, 48, c.Height-2, c.Width, c.Height)
		}

	}
	c.forceScreenRefresh = true
	c.resizeActive = true
	c.abortSend = true
}

// IsConsoleValidSize returns whether or not the console meets the global size requirements
func (c *Console) IsConsoleValidSize() bool {
	return c.Width > MINWIDTH && c.Height > MINHEIGHT
}

func (c *Console) IsUserLoggedIn() bool {
	c.mutex.Lock()
	c.mutex.Unlock()
	return c.userLoggedIn
}
