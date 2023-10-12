package ui

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/terminals"
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"github.com/yamamushi/EscapingEden/ui/window"
	"github.com/yamamushi/EscapingEden/ui/window/chat"
	"github.com/yamamushi/EscapingEden/ui/window/loginmenu"
	"github.com/yamamushi/EscapingEden/ui/window/toolbox"
	"sync"
)

const (
	MINWIDTH  = 155
	MINHEIGHT = 30
)

// Console is the main UI object.
type Console struct {
	ConnectionID string // The ID of the connection using this console
	Terminal     terminals.TerminalType

	Log logging.LoggerType

	Height int // The height of the console
	Width  int // The width of the console

	UserInfo      messages.UserInfo // The user info for the user using this console, generated after logging in, and updated as needed
	userInfoMutex sync.Mutex

	CharacterInfo      messages.CharacterInfo // The character info for the user using this console, generated after logging in, and updated as needed
	characterInfoMutex sync.Mutex

	Windows            []window.WindowType // The list of windows that are currently in the console
	ConsoleCommands    string
	LastSentOutput     string
	LastActiveWindow   window.WindowType
	mutex              sync.Mutex
	Shutdown           bool
	consoleInitialized bool
	initConsole        sync.Once

	// Channels for communicating with ConnectionManager
	WindowMessages  chan messages.WindowMessage
	SendMessages    chan messages.ConnectionManagerMessage
	ReceiveMessages chan messages.ConsoleMessage

	// Channels for communicating with windows
	LoginWindowMessages    chan messages.WindowMessage
	UserDashboardMessages  chan messages.WindowMessage
	GameWindowMessages     chan messages.WindowMessage
	ChatMessageReceive     chan messages.ChatMessage
	ChatWindowMessages     chan messages.WindowMessage
	ToolboxWindowMessages  chan messages.WindowMessage
	PopupBoxWindowMessages chan messages.WindowMessage

	// Our console character map
	// Memory is cheap with these
	// And we're going to use it
	PointMap         types.PointMap
	LastSentPointMap types.PointMap
	pmapMutex        sync.Mutex

	escapeBuffer           string
	escapeSequence         bool
	returnSequence         bool
	forceScreenRefresh     bool
	abortSend              bool
	abortSync              sync.Mutex
	resizeActive           bool
	userLoggedIn           bool
	userLoggedInMutex      sync.Mutex
	characterLoggedIn      bool // Whether the user has logged in with a character, chat windows will only be active for input if this is true
	characterLoggedInMutex sync.Mutex
	flushWindowList        []config.WindowID
}

// NewConsole creates a new console with no windows.
func NewConsole(height int, width int, connectionID string, outputChannel chan messages.ConnectionManagerMessage,
	log logging.LoggerType, term terminals.TerminalType) *Console {
	// Set up a new Chat Window and add it to the console at the bottom.
	receiveFromConnectionManager := make(chan messages.ConsoleMessage)
	windowMessages := make(chan messages.WindowMessage)
	loginMessages := make(chan messages.WindowMessage)
	gameWindowMessages := make(chan messages.WindowMessage)
	chatMessageReceive := make(chan messages.ChatMessage)
	chatMessageSend := make(chan messages.WindowMessage)
	toolboxMessages := make(chan messages.WindowMessage)
	popupBoxMessage := make(chan messages.WindowMessage)
	userDashboardMessages := make(chan messages.WindowMessage)

	return &Console{Height: height, Width: width, ConnectionID: connectionID, SendMessages: outputChannel,
		ReceiveMessages: receiveFromConnectionManager, WindowMessages: windowMessages, LoginWindowMessages: loginMessages,
		ChatMessageReceive: chatMessageReceive, ChatWindowMessages: chatMessageSend, UserDashboardMessages: userDashboardMessages, ToolboxWindowMessages: toolboxMessages,
		GameWindowMessages:     gameWindowMessages,
		PopupBoxWindowMessages: popupBoxMessage, Log: log, Terminal: term}
}

// Init is called once to initialize the console, it does things like create the default windows, and launches the
// Capture Message goroutines. It also sets up some default console commands that should be sent on first connect.
// These can be appended to later if we want to send a quick console command, but they are flushed after every write.
// Ie they are single use.
func (c *Console) Init() {

	c.consoleInitialized = false
	// First we set up our login window
	loginWindow := login.NewLoginWindow(0, 0, c.Width-50, c.Height-13, c.Width, c.Height,
		c.LoginWindowMessages, c.WindowMessages, c.Log, c.Terminal)
	loginWindow.Init()
	c.AddWindow(loginWindow)

	// Next we set up our chat window
	chatWindow := chat.NewChatWindow(0, c.Height-10, c.Width-50, 9, c.Width, c.Height,
		c.ChatMessageReceive, c.ChatWindowMessages, c.WindowMessages, c.Log, c.Terminal)
	chatWindow.Init()
	c.AddWindow(chatWindow)

	// Then we add our toolbox last
	toolboxWindow := toolbox.NewToolboxWindow(c.Width-48, 0, 48, c.Height-2, c.Width,
		c.Height, c.ToolboxWindowMessages, c.WindowMessages, c.Log, c.Terminal)
	toolboxWindow.Init()
	c.AddWindow(toolboxWindow)

	c.PointMap = types.NewPointMap(c.Width, c.Height)
	c.LastSentPointMap = types.NewPointMap(c.Width, c.Height)

	go c.CaptureWindowMessages()
	go c.CaptureManagerMessages()

	// Set our default active window to the login window, we will pass this to another window after we log in.
	c.SetActiveWindow(loginWindow)

	c.ConsoleCommands += c.Terminal.HideCursor() + c.Terminal.ClearTerminal()
}

// SetManagerSendChannel sets the channel for sending messages to the ConnectionManager.
func (c *Console) SetManagerSendChannel(ch chan messages.ConnectionManagerMessage) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.SendMessages = ch
}

// SetManagerReceiveChannel sets the channel for receiving messages from the ConnectionManager.
func (c *Console) SetManagerReceiveChannel(ch chan messages.ConsoleMessage) {
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
	//log.Println("Forcing redraw")
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
		case config.WindowUserDashboard:
			w.UpdateParams(0, 0, c.Width-50, c.Height-13, c.Width, c.Height)
		case config.WindowGameDisplay:
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
