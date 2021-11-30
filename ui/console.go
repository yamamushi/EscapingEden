package ui

import (
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"github.com/yamamushi/EscapingEden/ui/window"
	"github.com/yamamushi/EscapingEden/ui/window/chat"
	"github.com/yamamushi/EscapingEden/ui/window/mainmenu"
	"github.com/yamamushi/EscapingEden/ui/window/toolbox"
	"log"
	"sync"
)

const (
	MINWIDTH  = 155
	MINHEIGHT = 30
)

// Console is the main UI object.
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

	// Our console character map
	// Memory is cheap with these
	// And we're going to use it
	PointMap         types.PointMap
	LastSentPointMap types.PointMap
	pmapMutex        sync.Mutex

	escapeBuffer       string
	escapeSequence     bool
	returnSequence     bool
	forceScreenRefresh bool
	abortSend          bool
	abortSync          sync.Mutex
	resizeActive       bool
	userLoggedIn       bool
	flushWindowList    []config.WindowID
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

	c.PointMap = types.NewPointMap(c.Width, c.Height)
	c.LastSentPointMap = types.NewPointMap(c.Width, c.Height)

	go c.CaptureWindowMessages()
	go c.CaptureManagerMessages()

	c.SetActiveWindow(loginWindow) // Set our default active window to the login window, we will pass this to another
	// window after we log in.

	c.ConsoleCommands += c.HideCursor() + c.ResetTerminal() // + c.DrawPrompt() + c.MoveCursorToTopLeft()
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

// DrawPrompt returns the prompt for the console
func (c *Console) DrawPrompt() string {
	output := c.MoveCursorToBottomLeft()
	return output + "> "
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

// ForceRedraw forces a redraw of the console.
func (c *Console) ForceRedraw() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	log.Println("Forcing redraw")
	for _, w := range c.Windows {
		w.FlushLastSent()
	}
	c.ClearPointMap()
	c.FlushLastSent()
	c.forceScreenRefresh = true
}

// AbortSend tells the console to abort sending the last content it has in its buffer.
// We do this when we flush the console to avoid sending duplicate content.
func (c *Console) AbortSend() {
	c.abortSync.Lock()
	defer c.abortSync.Unlock()
	c.abortSend = true
}

// HandleResize handles the resize event for all windows
func (c *Console) HandleResize(newWidth, newHeight int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.pmapMutex.Lock()
	defer c.pmapMutex.Unlock()

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

	for _, w := range c.Windows {
		w.ResetWindowDrawings()
	}
	c.PointMap = types.NewPointMap(c.Width, c.Height)
	c.LastSentPointMap = types.NewPointMap(c.Width, c.Height)

	c.forceScreenRefresh = true
	c.resizeActive = true
	c.abortSend = true
}

// IsConsoleValidSize returns whether or not the console meets the global size requirements
func (c *Console) IsConsoleValidSize() bool {
	return c.Width > MINWIDTH && c.Height > MINHEIGHT
}

// IsUserLoggedIn returns whether or not the user is logged in
func (c *Console) IsUserLoggedIn() bool {
	c.mutex.Lock()
	c.mutex.Unlock()
	return c.userLoggedIn
}
