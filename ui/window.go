package ui

import (
	"fmt"
	"strconv"
	"sync"
)

const (
	// Some default window ID's that are used by the console
	DEBUGBOX   = 0
	CONSOLE    = 1
	CHATBOX    = 2
	INVENTORY  = 3
	MINIMAP    = 4
	PLAYERINFO = 5
	PLAYERLIST = 6
	STATUS     = 7
	TARGET     = 8
	TARGETINFO = 9
	TARGETLIST = 10
	LOGINMENU  = 11
	WORLDMAP   = 12
	TOOLBOX    = 13
	POPUPBOX   = 14
)

type WindowType interface {
	Draw(X, Y, height, width, startX, startY int) string
	HandleInput(input Input)

	DrawBorder(X, Y, height, width int) string
	UpdateContents()
	SetContents(string)
	PrintAt(X, Y int, text string) string

	GetID() int
	GetX() int
	GetY() int
	GetStartX() int
	GetStartY() int
	GetWidth() int
	GetHeight() int
	GetContents() string
	GetActive() bool
	SetActive(bool)
	GetHidden() bool
	GetBordered() bool
	GetFG() int
	GetBG() int
	GetBorderFG() int
	GetBorderBG() int

	Error(string)
	Quit()
}

type Window struct {
	ID int
	X  int // The X position of the Window
	Y  int // The Y position of the Window

	StartX int // When window content is rendered, it is a 2D array, so this is the starting X position of the content
	StartY int // When window content is rendered, it is a 2D array, so this is the starting Y position of the content

	Contents         string // The contents of the window
	ContentStartPos  int    // The starting position of the content
	LastSentContents string // The last contents sent to the client

	Width    int  // The width of the Window
	Height   int  // The height of the Window
	Active   bool // Is the Window active?
	Hidden   bool // Is the Window hidden?
	Bordered bool // Is the Window bordered?

	FG int // The foreground color of the Window Text
	BG int // The background color of the Window Text

	BorderFG int // The foreground color of the Window Border
	BorderBG int // The background color of the Window Border

	// Channels for communicating with ConnectionManager
	ConsoleSend    chan string // Send messages to the Console
	ConsoleReceive chan string // Receive messages from the Console

	mutex sync.Mutex
}

// Draw returns a string of the Window's contents
func (w *Window) Draw(X int, Y int, height, width int, startX, startY int) string {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	output := w.MoveCursorTopLeft()

	output += w.ParseContents(X, Y, height, width, startX, startY)
	output += w.DrawBorder(X, Y, height, width)

	return output
}

func (w *Window) HandleInput(input Input) {
	return
}

func (w *Window) Error(err string) {
	consoleMessage := &ConsoleMessage{Type: "error", Message: err}
	w.ConsoleSend <- consoleMessage.String()
}

func (w *Window) Quit() {
	consoleMessage := &ConsoleMessage{Type: "quit"}
	w.ConsoleSend <- consoleMessage.String()
}

// These functions implement the default WindowType interface for Window
func (w *Window) GetID() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.ID
}

// Move takes X and Y coordinates and moves the window to those coordinates
func (w *Window) Move(X int, Y int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.X = X
	w.Y = Y
}

// updateContent updates the contents of the window
func (w *Window) UpdateContents() {
	return
}

func (w *Window) GetX() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.X
}

func (w *Window) GetY() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.Y
}

func (w *Window) GetStartX() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.StartX
}

func (w *Window) GetStartY() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.StartY
}

func (w *Window) GetWidth() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.Width
}

func (w *Window) GetHeight() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.Height
}

func (w *Window) GetContents() string {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.Contents
}

// SetContent sets the contents of the window
func (w *Window) SetContents(content string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.Contents = content
}

// GetActive returns whether or not the Window is active
func (w *Window) GetActive() bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.Active
}

// SetActive sets the active state of the Window
func (w *Window) SetActive(active bool) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.Active = active
}

// GetHidden returns whether or not the Window is hidden
func (w *Window) GetHidden() bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.Hidden
}

// GetBordered returns whether or not the Window is bordered
func (w *Window) GetBordered() bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.Bordered
}

// GetFG returns the foreground color of the Window
func (w *Window) GetFG() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.FG
}

// GetBG returns the background color of the Window
func (w *Window) GetBG() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.BG
}

// GetBorderFG returns the foreground color of the Window Border
func (w *Window) GetBorderFG() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.BorderFG
}

// GetBorderBG returns the background color of the Window Border
func (w *Window) GetBorderBG() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.BorderBG
}

func (w *Window) IncreaseContentPos() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.ContentStartPos++
}

func (w *Window) DecreaseContentPos() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.ContentStartPos--
}

// MoveCursorTopLeft moves the cursor to the top left of the Window and returns as a string
func (w *Window) MoveCursorTopLeft() string {
	// set cursor to top left of window taking into account border
	return fmt.Sprintf("\033[%d;%dH", w.Y+1, w.X+1)
}

func (w *Window) PrintAt(X int, Y int, text string) string {
	// set cursor to X and Y
	return fmt.Sprintf("\033[%d;%dH%s", X+w.GetX(), Y+w.GetY(), text)
}

// CenterText takes a text string and outputs it at the center of the window
func (w *Window) CenterText(text string, line int) string {
	// get the length of the text
	length := len(text)

	// get the center of the window
	center := (w.Width / 2) - (length / 2)

	// return the text centered in the window
	return fmt.Sprintf("\033[%d;%dH%s", w.Y+line, w.X+center+1, text)
}

// DrawBorder draws the Window's border

// DrawBorder returns the border of a window using code page 437 characters as a string
func (w *Window) DrawBorder(winX int, winY int, visibleLength, visibleHeight int) string {
	// Draw top border using code page 437 characters starting at position winX, winY

	// Move cursor to top left corner of window
	border := "\033[" + strconv.Itoa(winY) + ";" + strconv.Itoa(winX) + "H"
	// Draw top left corner
	border += "\u250c"
	for i := 0; i < visibleLength; i++ {
		// Inserts a horizontal line
		border += "\u2500"
	}
	// insert the top right corner
	border += "\u2510"

	// For visibleHeight draw a left and right border at each line
	for i := 0; i < visibleHeight; i++ {
		// Move cursor to left side of window
		border += "\033[" + strconv.Itoa(winY+i+1) + ";" + strconv.Itoa(winX) + "H"
		// Draw left border
		border += "\u2502"
		// Move cursor to right side of window
		border += "\033[" + strconv.Itoa(winY+i+1) + ";" + strconv.Itoa(winX+visibleLength+1) + "H"
		// Draw right border
		border += "\u2502"
	}

	// Move cursor to bottom left corner of window
	border += "\033[" + strconv.Itoa(winY+visibleHeight+1) + ";" + strconv.Itoa(winX) + "H"
	// Draw bottom left corner
	border += "\u2514"
	for i := 0; i < visibleLength; i++ {
		// Inserts a horizontal line
		border += "\u2500"
	}
	// insert the bottom right corner
	border += "\u2518"

	return border
}

// Parse contents reads a string one character at a time, placing it within the bounds of the window and returns the string
func (w *Window) ParseContents(winX int, winY int, visibleLength, visibleHeight int, startX, startY int) string {
	// Parse contents of window into a string

	// maxLength is the maximum length of the window subtracting the border
	maxLength := visibleLength
	// maxHeight is the maximum height of the window subtracting the border
	maxHeight := visibleHeight

	currentColumn := winX + 1

	var lines []string
	parsed := ""

	// For every character in the contents
	for i := 0; i < len(w.Contents); i++ {
		if currentColumn > maxLength+winX {
			// reset currentColumn
			currentColumn = winX + 1
			// append the current line to the lines slice
			lines = append(lines, parsed)
			// reset parsed
			parsed = ""

		}
		// if the character is a newline
		if w.Contents[i] == '\n' {
			// reset currentColumn
			currentColumn = winX + 1
			// append the current line to the lines slice
			lines = append(lines, parsed)
			// reset parsed
			parsed = ""
			currentColumn++
			continue
		}

		// append the character to the parsed string
		parsed += string(w.Contents[i])
		// increment currentColumn
		currentColumn++

		// If this is the last character in the contents
		if i == len(w.Contents)-1 {
			// append the current line to the lines slice
			lines = append(lines, parsed)
		}
	}

	// Move cursor to top left corner of window accounting for the border
	// and the visible length and height of the window
	output := "\033[" + strconv.Itoa(winY+1) + ";" + strconv.Itoa(winX+1) + "H"
	currentLine := winY + 1

	// append the last maxHeight lines to the output string
	if len(lines) >= maxHeight {
		for i := len(lines) - maxHeight + w.ContentStartPos; i < len(lines); i++ {
			output += lines[i]
			// increment currentLine
			currentLine++
			// Move cursor down one line
			output += "\033[" + strconv.Itoa(currentLine) + ";" + strconv.Itoa(winX+1) + "H"

		}
	} else {
		for i := 0; i < len(lines); i++ {
			output += lines[i]
			// increment currentLine
			currentLine++
			// Move cursor down one line
			output += "\033[" + strconv.Itoa(currentLine) + ";" + strconv.Itoa(winX+1) + "H"
		}
	}

	return output
}
