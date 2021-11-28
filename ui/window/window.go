package window

import (
	"fmt"
	"github.com/yamamushi/EscapingEden/ui/console"
	"log"
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
	ClearMap(X, Y, height, width, startX, startY int)
	Draw(X, Y, height, width, startX, startY int)
	HandleInput(input Input)
	Init()

	HandleReceive(message console.ConsoleMessage)

	DrawBorder(X, Y, height, width int)
	DrawContents(X, Y, visibleHeight, visibleWidth, startX, startY int)
	UpdateContents()
	SetContents(string)
	PrintLn(X, Y int, text string, escaoeCode string)
	PointMapToString() string
	FlushLastSent()

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

	CheckScrollBufferNew() bool
	SetScrollBufferNew(bool)
	SupportsScrolling() bool
	GetContentStartPos() int

	Error(string)
	Quit()
}

type Window struct {
	ID int
	X  int // The X position of the Window
	Y  int // The Y position of the Window

	StartX int // When window content is rendered, it is a 2D array, so this is the starting X position of the content
	StartY int // When window content is rendered, it is a 2D array, so this is the starting Y position of the content

	Contents           string // The contents of the window
	ContentStartPos    int    // The starting position of the content
	ScrollingSupported bool
	ScrollBufferHasNew bool // Indicates that the scroll buffer has new content

	Width         int // The width of the Window
	Height        int // The height of the Window
	ConsoleWidth  int // The width of the console
	ConsoleHeight int // The height of the console

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

	DirectionInput InputType // The last direction input from the user

	mutex               sync.Mutex
	pmapMutex           sync.Mutex
	pointMap            console.PointMap
	lastSentContents    console.PointMap // The last contents sent to the client
	pointMapInitialized bool
}

// Draw returns a string of the Window's contents
func (w *Window) Draw(X int, Y int, visibleHeight, visibleWidth int, startX, startY int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.DrawContents(X, Y, visibleHeight, visibleWidth, startX, startY)
}

func (w *Window) Init() {
	w.pointMap = console.NewPointMap(w.ConsoleWidth, w.ConsoleHeight)
	w.lastSentContents = console.NewPointMap(w.ConsoleWidth, w.ConsoleHeight)
	w.pointMapInitialized = true
}

func (w *Window) HandleInput(input Input) {
	return
}

func (w *Window) Error(err string) {
	consoleMessage := &console.ConsoleMessage{Type: "error", Message: err}
	w.ConsoleSend <- consoleMessage.String()
}

func (w *Window) Quit() {
	consoleMessage := &console.ConsoleMessage{Type: "quit"}
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

func (w *Window) CheckScrollBufferNew() bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.ScrollBufferHasNew
}

func (w *Window) SetScrollBufferNew(new bool) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.ScrollBufferHasNew = new
}

func (w *Window) GetY() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.Y
}

func (w *Window) SupportsScrolling() bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.ScrollingSupported
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

func (w *Window) GetContentStartPos() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.ContentStartPos
}

func (w *Window) SetContentStartPos(pos int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.ContentStartPos = pos
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
	w.DirectionInput = InputUp
}

func (w *Window) DecreaseContentPos() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.DirectionInput = InputDown
}

// MoveCursorTopLeft moves the cursor to the top left of the Window and returns as a string
func (w *Window) MoveCursorTopLeft() string {
	// set cursor to top left of window taking into account border
	return fmt.Sprintf("\033[%d;%dH", w.Y+1, w.X+1)
}

func (w *Window) HandleReceive(message console.ConsoleMessage) {
	w.ConsoleReceive <- message.String()
}

func (w *Window) RequestPopupFromConsole(x, y, width, height int, content string) {
	log.Println("Requesting popup from console")
	config := PopupConfig(x, y, width, height, content)
	log.Println(config.String())
	request := console.ConsoleMessage{Type: "console", Message: "popup", Options: config.String()}
	log.Println(request.String())
	w.ConsoleSend <- request.String()
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
func (w *Window) DrawBorder(winX int, winY int, visibleLength, visibleHeight int) {
	// Draw top border using code page 437 characters starting at position winX, winY

	// Move cursor to top left corner of window
	// Draw top left corner
	if w.Active {
		w.PrintChar(winX, winY, "\u250c", "\033[32m")

	} else {
		w.PrintChar(winX, winY, "\u250c", "")
	}

	// Draw left border
	for i := 1; i < visibleHeight+1; i++ {
		// Inserts a vertical line
		if w.Active {
			w.PrintChar(winX, winY+i, "\u2502", "\033[32m")
		} else {
			w.PrintChar(winX, winY+i, "\u2502", "")
		}
	}
	// Draw bottom left corner
	if w.Active {
		w.PrintChar(winX, winY+visibleHeight+1, "\u2514", "\033[32m")
	} else {
		w.PrintChar(winX, winY+visibleHeight+1, "\u2514", "")
	}

	// Draw top border
	for i := 1; i < visibleLength; i++ {
		// Inserts a horizontal line
		if w.Active {
			w.PrintChar(winX+i, winY, "\u2500", "\033[32m")
		} else {
			w.PrintChar(winX+i, winY, "\u2500", "")
		}
	}

	// Draw top right corner
	if w.Active {
		w.PrintChar(winX+visibleLength, winY, "\u2510", "\033[32m")
	} else {
		w.PrintChar(winX+visibleLength, winY, "\u2510", "")
	}

	// Draw right border
	for i := 1; i < visibleHeight+1; i++ {
		// Inserts a vertical line
		if w.Active {
			w.PrintChar(winX+visibleLength, winY+i, "\u2502", "\033[32m")
		} else {
			w.PrintChar(winX+visibleLength, winY+i, "\u2502", "")
		}
	}

	// Draw bottom right corner
	if w.Active {
		w.PrintChar(winX+visibleLength, winY+visibleHeight+1, "\u2518", "\033[32m")
	} else {
		w.PrintChar(winX+visibleLength, winY+visibleHeight+1, "\u2518", "")
	}

	// Draw bottom border
	for i := 1; i < visibleLength; i++ {
		// Inserts a horizontal line
		if w.Active {
			w.PrintChar(winX+i, winY+visibleHeight+1, "\u2500", "\033[32m")
		} else {
			w.PrintChar(winX+i, winY+visibleHeight+1, "\u2500", "")
		}
	}
}

func (w *Window) ContentToLines(winX int, winY int, visibleLength int) ([]string, int) {
	// maxLength is the maximum length of the window subtracting the border
	maxLength := visibleLength

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
	return lines, len(lines)
}

// Parse contents reads a string one character at a time, placing it within the bounds of the window and returns the string
func (w *Window) DrawContents(winX int, winY int, visibleLength, visibleHeight int, startX, startY int) {
	// maxHeight is the maximum height of the window subtracting the border
	maxHeight := visibleHeight
	lines, _ := w.ContentToLines(winX, winY, visibleLength)
	currentLine := winY + 1

	if len(lines) > maxHeight {
		if w.ScrollingSupported {
			if w.DirectionInput == InputUp {
				w.ContentStartPos++
				w.DirectionInput = 0
			} else if w.DirectionInput == InputDown {
				w.ContentStartPos--
				w.DirectionInput = 0
			}
		}

		if w.ContentStartPos == 0 {
			w.ScrollBufferHasNew = false
		}

		contentStartPos := 0
		if len(lines)-maxHeight+w.ContentStartPos-1 < 0 {
			contentStartPos = 0 - len(lines) + maxHeight + 1
			w.ContentStartPos = 0 - len(lines) + maxHeight + 1
		} else if len(lines)-maxHeight+w.ContentStartPos-1 > len(lines)-maxHeight-1 {
			contentStartPos = 0
			w.ContentStartPos = 0
		} else {
			contentStartPos = w.ContentStartPos
		}

		for i := len(lines) - maxHeight + contentStartPos - 1; i < len(lines); i++ {
			if currentLine > maxHeight+winY+1 {
				break
			}
			// Print current line
			w.PrintLn(winX+1, currentLine, lines[i], "")
			// Fill the rest of the line with spaces
			for j := len(lines[i]) + 1; j < visibleLength; j++ {
				if w.GetCharAt(winX+j, currentLine) == "" {
					w.PrintChar(winX+j, currentLine, " ", "")
				}
			}
			// increment currentLine
			currentLine++
		}

		// Draw our arrows last
		if len(lines)-maxHeight+contentStartPos-1 > 0 {
			// draw an up arrow in grey
			w.PrintChar(winX+visibleLength, winY+1, "\u2191", "\033[37m")
		}
		if len(lines)-maxHeight+contentStartPos-1 < len(lines)-maxHeight-1 {
			if w.ScrollBufferHasNew {
				// Draw down arrow in red if there is new content
				w.PrintChar(winX+visibleLength, winY+visibleHeight+1, "\u2193", "\033[31m")
			} else {
				// Draw down arrow in grey if there is no new content
				w.PrintChar(winX+visibleLength, winY+visibleHeight+1, "\u2193", "\033[37m")
			}
		}
	} else {
		if w.ScrollingSupported {
			// If the length of content doesn't exceed our visible height, we don't need to scroll
			// And we can discard the DirectionInput
			if w.DirectionInput == InputUp {
				w.ContentStartPos = 0
				w.DirectionInput = 0
			} else if w.DirectionInput == InputDown {
				w.ContentStartPos = 0
				w.DirectionInput = 0
			}
		}

		for i := 0; i < len(lines); i++ {
			w.PrintLn(winX+1, currentLine, lines[i], "")

			if i == len(lines)-1 {
				for lineNumber := i; lineNumber < maxHeight; lineNumber++ {
					// increment currentLine
					currentLine++
					for j := 0; j < visibleLength; j++ {
						if w.GetCharAt(winX+j, currentLine) == "" {
							w.PrintLn(winX+j, currentLine, " ", "")
						}
					}
				}
			}
			// increment currentLine
			currentLine++
		}
	}
}

func (w *Window) PrintLn(X int, Y int, text string, escapeCode string) {
	w.pmapMutex.Lock()
	defer w.pmapMutex.Unlock()
	if X > len(w.pointMap)-1 {
		return
	}
	if Y > len(w.pointMap[X])-1 {
		return
	}

	for i, character := range text {
		//log.Println("inserting character:", string(character))
		// For the point at X, Y+1, set the character to the character at the current index of the text string
		w.pointMap[X+i][Y] = console.Point{X: X + i, Y: Y, Character: string(character), EscapeCode: escapeCode}
		//log.Println("pointMap:", X, Y+i, w.pointMap[X][Y+i].Character)
	}
}

func (w *Window) PrintChar(X int, Y int, text string, escapeCode string) {
	w.pmapMutex.Lock()
	defer w.pmapMutex.Unlock()
	if X > len(w.pointMap)-1 {
		return
	}
	if Y > len(w.pointMap[X])-1 {
		return
	}
	w.pointMap[X][Y] = console.Point{X: X, Y: Y, Character: text, EscapeCode: escapeCode}
}

func (w *Window) GetCharAt(X, Y int) string {
	w.pmapMutex.Lock()
	defer w.pmapMutex.Unlock()
	if X > len(w.pointMap)-1 {
		return ""
	}
	if Y > len(w.pointMap[X])-1 {
		return ""
	}
	return w.pointMap[X][Y].Print()
}

func (w *Window) PointMapToString() string {
	w.pmapMutex.Lock()
	defer w.pmapMutex.Unlock()

	// iterate through entire w.pointMap and print out the character at each point
	output := ""
	lastSentChar := ""
	lastSentEscape := ""
	lastY := 0
	lastX := 0
	bufferCount := 0

	for y := 0; y < len(w.pointMap[0]); y++ {
		for x := 0; x < len(w.pointMap); x++ {
			if w.pointMap[x][y].Character != "" || w.pointMap[x][y].EscapeCode != "" {

				if w.lastSentContents[x][y].Print() != w.pointMap[x][y].Print() {
					pointMapChar := w.pointMap[x][y].Character
					pointMapEscape := w.pointMap[x][y].EscapeCode
					// If this character is the last one sent, then we increase the buffer count
					// and repeat
					if pointMapChar == lastSentChar && pointMapEscape == lastSentEscape &&
						y == lastY && (x)-lastX == 1 {
						log.Println(pointMapChar)
						bufferCount++
						lastX = x
					} else {
						// If we reached a new character, and the buffer count is greater than 0
						// We need to print the repeated last character bufferCount times
						if bufferCount > 0 {
							repeatCode := lastSentEscape + "\033[" + strconv.Itoa(bufferCount) + "b" + "\033[0m"
							output += repeatCode
							// Finally Reset the buffer count
							bufferCount = 0
						} else {
							// If the buffer count was already 0, we update the last sent character
							// And reset the buffer count for verbosity
							lastSentChar = pointMapChar
							lastSentEscape = pointMapEscape
							lastY = w.pointMap[x][y].Y
							lastX = w.pointMap[x][y].X
							bufferCount = 0
						}
						// Now that we have dealt with the buffer count, we can print the new character
						output += w.pointMap[x][y].Print()

					}

					// Finally, no matter what we do with the character, we still append it to
					// The last sent contents, as printing it will still take up column spaces
					w.lastSentContents[x][y] = w.pointMap[x][y]
				}
			}
		}
	}
	return output
}

func (w *Window) ClearMap(winX int, winY int, visibleLength, visibleHeight int, startX, startY int) {
	// First clear the window before we redraw it
	for i := winX; i < visibleLength+2; i++ {
		for j := winY; j < visibleHeight+2; j++ {
			w.PrintChar(i, j, " ", "")
		}
	}
}

func (w *Window) FlushLastSent() {
	w.pmapMutex.Lock()
	defer w.pmapMutex.Unlock()
	w.lastSentContents = console.NewPointMap(w.ConsoleWidth, w.ConsoleHeight)
}

func (w *Window) ResetWindowDrawings() {
	w.SetContents("")
	w.ClearMap(w.X, w.Y, w.Width, w.Height, 0, 0)
	w.FlushLastSent()
}
