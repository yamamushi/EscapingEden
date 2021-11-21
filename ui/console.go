package ui

import (
	"log"
	"strconv"
	"sync"
)

type Console struct {
	Height int // The height of the console
	Width  int // The width of the console

	Windows         []WindowType // The list of windows that are currently in the console
	ConsoleCommands string
	LastSentOutput  string
	mutex           sync.Mutex
}

// NewConsole creates a new console with no windows.
func NewConsole(height int, width int) *Console {
	return &Console{Height: height, Width: width}
}

// AddWindow adds a window to the console if it is not already in the console by ID.
func (c *Console) AddWindow(w WindowType) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, window := range c.Windows {
		if window.GetID() == w.GetID() {
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
			window.UpdateContent()
			s = s + c.DrawWindow(window)

			if window == c.Windows[len(c.Windows)-1] {
				s = s + c.MoveCursorToBottomLeft()
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

	log.Println("Input recieved: " + input)
	for _, window := range c.Windows {
		window.HandleInput(input)
	}
}

func (c *Console) Init() {
	// Setup a new Chat Window and add it to the console
	chatWindow := NewChatWindow(5, 10, c.Width/2, c.Height/2)
	c.AddWindow(chatWindow)
	c.ConsoleCommands += c.HardClear() + c.MoveCursorToTopLeft()
}

// Moves the cursor to the top left corner of the console
func (c *Console) MoveCursorToTopLeft() string {
	return "\033[1;1H"
}

// Moves the cursor to the bottom left corner of the console
func (c *Console) MoveCursorToBottomLeft() string {
	return "\033[" + strconv.Itoa(c.Height) + ";0H"
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
	content += window.Draw(winX, winY, visibleLength, visibleHeight)

	return content
}
