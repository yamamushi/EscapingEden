package ui

import (
	"log"
	"strconv"
	"strings"
	"sync"
)

type Console struct {
	Height int // The height of the console
	Width  int // The width of the console

	Windows         []WindowType // The list of windows that are currently in the console
	ConsoleCommands string
	LastSentOutput  string
	mutex           sync.Mutex
	Shutdown        bool

	// Channels for communicating with ConnectionManager
	SendMessages    chan string
	ReceiveMessages chan string
}

// NewConsole creates a new console with no windows.
func NewConsole(height int, width int, outputChannel chan string) *Console {
	// Setup a new Chat Window and add it to the console at the bottom.
	receiver := make(chan string)

	return &Console{Height: height, Width: width, SendMessages: outputChannel, ReceiveMessages: receiver}
}

func (c *Console) Init() {

	// First we setup our login window
	loginWindow := NewLoginWindow(0, 0, c.Width-50, c.Height-12, c.ReceiveMessages, c.SendMessages)
	c.AddWindow(loginWindow)
	c.SetActiveWindow(loginWindow) // Set our default active window to the login window, we will pass this to another
	// window after we log in.

	// Next we setup our chat window
	chatWindow := NewChatWindow(0, c.Height-10, c.Width-50, 10, c.ReceiveMessages, c.SendMessages)
	c.AddWindow(chatWindow)

	// Then we add our toolbox last
	toolboxWindow := NewToolboxWindow(c.Width-48, 0, 50, c.Height, c.ReceiveMessages, c.SendMessages)
	c.AddWindow(toolboxWindow)

	c.ConsoleCommands += c.HardClear() + c.MoveCursorToTopLeft()
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
func (c *Console) AddWindow(w WindowType) {
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
	s = s + c.ConsoleCommands
	//c.ConsoleCommands = ""

	for _, window := range c.Windows {
		if !window.GetHidden() {

			window.UpdateContents()
			s = s + c.DrawWindow(window)

			if window == c.Windows[len(c.Windows)-1] {
				s = s + c.DrawPrompt()
			}
		}
	}

	// If the last output was not the same as the current output, we send it to the client and update the last output.
	if c.LastSentOutput != s && s != "" {
		c.LastSentOutput = s
		return []byte(s)
	} else {
		return []byte("")
	}
}

// HandleInput accepts a string terminated by a newline and processes it.
func (c *Console) HandleInput(input string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// If input contains a newline remove it
	input = strings.TrimRight(input, "\n")
	input = strings.TrimRight(input, "\r")
	input = strings.TrimRight(input, "\r\n")
	/*
		// If length of input is greater than 80
		if len(input) > 80 {
			c.GetChatWindow().ConsoleMessage("Input too long (Max 80 characters). Please try again.")
			return
		}
	*/

	log.Println("Input recieved: " + input)
	if input == "quit" {
		c.SetShutdown(true)
		return
	}

	for _, window := range c.Windows {
		if window.GetActive() {
			window.HandleInput(input)
			log.Println("Input Handled on window: ", window.GetID())
			return
		}
	}
}

// GetChatWindow returns the chat window.
func (c *Console) GetChatWindow() *ChatWindow {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, window := range c.Windows {
		if window.GetID() == 0 {
			return window.(*ChatWindow)
		}
	}
	return nil
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
func (c *Console) GetWindowAttrs(window WindowType) (X int, Y int, visibleLength int, visibleHeight int) {
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
func (c *Console) DrawWindow(window WindowType) (content string) {
	// Get Window Attrs
	winX, winY, visibleLength, visibleHeight := c.GetWindowAttrs(window)

	// Draw content of window
	content += window.Draw(winX, winY, visibleLength, visibleHeight, 0, 0)

	return content
}

// SetActiveWindow sets the active window and sets all other windows to inactive
func (c *Console) SetActiveWindow(window WindowType) {
	for _, w := range c.Windows {
		if w.GetID() == window.GetID() {
			log.Println("Active window set to: ", w.GetID())
			w.SetActive(true)
		} else {
			w.SetActive(false)
		}
	}
}
