package window

import (
	"encoding/json"
	"fmt"
	"github.com/yamamushi/EscapingEden/ui/types"
	"log"
	"strconv"
	"sync"
)

type ID int

// These are used as ID's for tracking drawing and other behavior.
// They are unique, and only one window of any type can be open at any time.
const (
	DEBUGBOX ID = iota
	HELPBOX
	CHATBOX
	LOGINMENU
	TOOLBOX
	POPUPBOX
)

type Config struct {
	X       int
	Y       int
	Width   int
	Height  int
	Content string
	Page    types.HelpPage
}

func NewWindowConfig(x, y, width, height int, content string) *Config {
	return &Config{X: x, Y: y, Width: width, Height: height, Content: content}
}

func (c *Config) String() string {
	output, _ := json.Marshal(c)
	return string(output)
}

type WindowType interface {
	ClearMap(X, Y, height, width, startX, startY int)
	Draw(X, Y int)
	HandleInput(input types.Input)
	Init()

	HandleReceive(message types.ConsoleMessage)

	DrawBorder(X, Y int)
	DrawContents(X, Y int)
	UpdateContents()
	SetContents(string)
	PrintLn(X, Y int, text string, escaoeCode string)
	PointMapToString() string
	FlushLastSent()
	ResetWindowDrawings()

	GetID() ID
	GetX() int
	GetY() int
	UpdateParams(x, y, height, width, consoleWidth, consoleHeight int)
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
	GetConfig() *Config

	CheckScrollBufferNew() bool
	SetScrollBufferNew(bool)
	SupportsScrolling() bool
	GetContentStartPos() int

	LockMutex()
	UnlockMutex()

	Error(string)
	Quit()
}

type Window struct {
	ID ID
	X  int // The X position of the Window
	Y  int // The Y position of the Window

	StartX int // When window content is rendered, it is a 2D array, so this is the starting X position of the content
	StartY int // When window content is rendered, it is a 2D array, so this is the starting Y position of the content

	Contents            string // The contents of the window
	ContentStartPos     int    // The starting position of the content
	ScrollingSupported  bool
	ScrollBufferHasNew  bool // Indicates that the scroll buffer has new content
	ScrollBufferLimit   int  // The maximum number of lines that can be stored in the scroll buffer
	ScrollBufferCharMod int  // The character limiter of the scroll buffer

	// For Window Paging, only HELPBOX uses this
	Page int

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

	DirectionInput types.InputType // The last direction input from the user

	mutex               sync.Mutex
	pmapMutex           sync.Mutex
	pointMap            types.PointMap
	lastSentContents    types.PointMap // The last contents sent to the client
	pointMapInitialized bool
}

// Draw returns a string of the Window's contents
func (w *Window) Draw(X int, Y int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.DrawContents(X, Y)
}

func (w *Window) Init() {
	w.pointMap = types.NewPointMap(w.ConsoleWidth, w.ConsoleHeight)
	w.lastSentContents = types.NewPointMap(w.ConsoleWidth, w.ConsoleHeight)
	w.pointMapInitialized = true
}

func (w *Window) HandleInput(input types.Input) {
	return
}

func (w *Window) Error(err string) {
	consoleMessage := &types.ConsoleMessage{Type: "error", Message: err}
	w.ConsoleSend <- consoleMessage.String()
}

func (w *Window) Quit() {
	consoleMessage := &types.ConsoleMessage{Type: "quit"}
	w.ConsoleSend <- consoleMessage.String()
}

// These functions implement the default WindowType interface for Window
func (w *Window) GetID() ID {
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

func (w *Window) LockMutex() {
	w.mutex.Lock()
}

func (w *Window) UnlockMutex() {
	w.mutex.Unlock()
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
	w.DirectionInput = types.InputUp
}

func (w *Window) DecreaseContentPos() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.DirectionInput = types.InputDown
}

// MoveCursorTopLeft moves the cursor to the top left of the Window and returns as a string
func (w *Window) MoveCursorTopLeft() string {
	// set cursor to top left of window taking into account border
	return fmt.Sprintf("\033[%d;%dH", w.Y+1, w.X+1)
}

func (w *Window) HandleReceive(message types.ConsoleMessage) {
	w.ConsoleReceive <- message.String()
}

func (w *Window) GetConfig() *Config {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return NewWindowConfig(w.X, w.Y, w.Width, w.Height, w.Contents)
}

func (w *Window) RequestPopupFromConsole(x, y, width, height int, content string) {
	log.Println("Requesting popup from console")
	config := NewWindowConfig(x, y, width, height, content)
	log.Println(config.String())
	request := types.ConsoleMessage{Type: "console", Message: "popup", Options: config.String()}
	log.Println(request.String())
	w.ConsoleSend <- request.String()
}

func (w *Window) RequestHelpFromConsole(page types.HelpPage) {
	log.Println("Requesting help from console")
	config := NewWindowConfig(w.ConsoleWidth/2-40, w.ConsoleHeight/2-10, 100, 20, "")
	log.Println(page)
	config.Page = page
	log.Println(config.String())
	request := types.ConsoleMessage{Type: "console", Message: "help", Options: config.String()}
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
func (w *Window) DrawBorder(winX int, winY int) {
	// Draw top border using code page 437 characters starting at position winX, winY

	// Move cursor to top left corner of window
	// Draw top left corner
	if w.Active {
		w.PrintChar(winX, winY, "\u250c", "\033[32m")

	} else {
		w.PrintChar(winX, winY, "\u250c", "")
	}

	// Draw left border
	for i := 1; i < w.Height+1; i++ {
		// Inserts a vertical line
		if w.Active {
			w.PrintChar(winX, winY+i, "\u2502", "\033[32m")
		} else {
			w.PrintChar(winX, winY+i, "\u2502", "")
		}
	}
	// Draw bottom left corner
	if w.Active {
		w.PrintChar(winX, winY+w.Height+1, "\u2514", "\033[32m")
	} else {
		w.PrintChar(winX, winY+w.Height+1, "\u2514", "")
	}

	// Draw top border
	for i := 1; i < w.Width; i++ {
		// Inserts a horizontal line
		if w.Active {
			w.PrintChar(winX+i, winY, "\u2500", "\033[32m")
		} else {
			w.PrintChar(winX+i, winY, "\u2500", "")
		}
	}

	// Draw top right corner
	if w.Active {
		w.PrintChar(winX+w.Width, winY, "\u2510", "\033[32m")
	} else {
		w.PrintChar(winX+w.Width, winY, "\u2510", "")
	}

	// Draw right border
	for i := 1; i < w.Height+1; i++ {
		// Inserts a vertical line
		if w.Active {
			w.PrintChar(winX+w.Width, winY+i, "\u2502", "\033[32m")
		} else {
			w.PrintChar(winX+w.Width, winY+i, "\u2502", "")
		}
	}

	// Draw bottom right corner
	if w.Active {
		w.PrintChar(winX+w.Width, winY+w.Height+1, "\u2518", "\033[32m")
	} else {
		w.PrintChar(winX+w.Width, winY+w.Height+1, "\u2518", "")
	}

	// Draw bottom border
	for i := 1; i < w.Width; i++ {
		// Inserts a horizontal line
		if w.Active {
			w.PrintChar(winX+i, winY+w.Height+1, "\u2500", "\033[32m")
		} else {
			w.PrintChar(winX+i, winY+w.Height+1, "\u2500", "")
		}
	}
}

func (w *Window) ContentToLines(winX int, winY int, visibleLength int) ([]string, int) {
	// Split the content into lines
	visibleLength = w.Width - 1 - w.ScrollBufferCharMod
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
func (w *Window) DrawContents(winX int, winY int) {
	visibleLength := w.Width - 1
	visibleHeight := w.Height - 1 - w.ScrollBufferLimit
	// maxHeight is the maximum height of the window subtracting the border
	maxHeight := visibleHeight
	lines, _ := w.ContentToLines(winX, winY, visibleLength)
	currentLine := winY + 1 + w.StartY // We use this for scrolling windows that want a little buffer at top

	if len(lines) > maxHeight {
		if w.ScrollingSupported {
			if w.DirectionInput == types.InputUp {
				w.ContentStartPos++
				w.DirectionInput = 0
			} else if w.DirectionInput == types.InputDown {
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
			for j, char := range lines[i] {
				if j+winX+1 > w.X+w.Width-1 {
					break
				}
				if w.GetCharAt(j+winX+1, currentLine) == " " || w.GetCharAt(j+winX+1, currentLine) == "" {
					w.PrintChar(j+winX+1, currentLine, string(char), "")
				}
			}

			//w.PrintLn(winX+1, currentLine, lines[i], "")
			// increment currentLine
			currentLine++
		}

		// Draw our arrows last
		if len(lines)-maxHeight+contentStartPos-1 > 0 {
			// draw an up arrow in grey
			w.PrintChar(winX+visibleLength, winY+1+w.StartY, "\u2191", "\033[37m")
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
			if w.DirectionInput == types.InputUp {
				w.ContentStartPos = 0
				w.DirectionInput = 0
			} else if w.DirectionInput == types.InputDown {
				w.ContentStartPos = 0
				w.DirectionInput = 0
			}
		}

		for i := 0; i < len(lines); i++ {

			for j, char := range lines[i] {
				if j+winX+1 > w.X+w.Width-1 {
					break
				}
				if w.GetCharAt(j+winX+1, currentLine) == " " || w.GetCharAt(j+winX+1, currentLine) == "" {
					w.PrintChar(j+winX+1, currentLine, string(char), "")
				}
				//w.PrintChar(i+winX+1, currentLine, string(char), "")
			} // increment currentLine
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
		w.pointMap[X+i][Y] = types.Point{X: X + i, Y: Y, Character: string(character), EscapeCode: escapeCode}
		//log.Println("pointMap:", X, Y+i, w.pointMap[X][Y+i].Character)
	}
}

func (w *Window) PrintChar(X int, Y int, text string, escapeCode string) {
	w.pmapMutex.Lock()
	defer w.pmapMutex.Unlock()
	if X > len(w.pointMap)-1 || X < 0 {
		return
	}
	if Y > len(w.pointMap[X])-1 || Y < 0 {
		return
	}
	w.pointMap[X][Y] = types.Point{X: X, Y: Y, Character: text, EscapeCode: escapeCode}
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
	return w.pointMap[X][Y].Character
}

func (w *Window) GetEscapeCodeAt(X, Y int) string {
	w.pmapMutex.Lock()
	defer w.pmapMutex.Unlock()
	if X > len(w.pointMap)-1 {
		return ""
	}
	if Y > len(w.pointMap[X])-1 {
		return ""
	}
	return w.pointMap[X][Y].EscapeCode
}

func (w *Window) IsPointBlank(X, Y int) bool {
	if w.GetEscapeCodeAt(X, Y) == "" && w.GetCharAt(X, Y) == "" {
		return true
	} else {
		return false
	}
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
						//log.Println(pointMapChar)
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
	for i := w.Y + 1; i < w.Y+w.Height+1; i++ {
		for j := w.X + 1; j < w.X+w.Width; j++ {
			//if w.GetCharAt(i, j) != " " { // && w.GetEscapeCodeAt(i, j) != "" {
			//log.Println("Blank point found: ", i, j)
			w.PrintChar(j, i, " ", "")
			//}
		}
	}
}

func (w *Window) FlushLastSent() {
	w.pmapMutex.Lock()
	defer w.pmapMutex.Unlock()
	w.lastSentContents = types.NewPointMap(w.ConsoleWidth, w.ConsoleHeight)
}

func (w *Window) ResetWindowDrawings() {
	w.FlushLastSent()
	w.SetContents("")
	w.ClearMap(w.X, w.Y, w.Width, w.Height, 0, 0)
}

func (w *Window) ForceConsoleRefresh() {
	w.ResetWindowDrawings()
	message := &types.ConsoleMessage{Type: "console", Message: "refresh", WindowID: int(w.GetID())}
	w.ConsoleSend <- message.String()
}

func (w *Window) UpdateParams(x, y, width, height, consoleWidth, consoleHeight int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.pmapMutex.Lock()
	defer w.pmapMutex.Unlock()

	// This function can probably also be used later for window moving

	// if x or y are less than 1 set them to 1
	if x < 1 {
		x = 1
	}
	if y < 1 {
		y = 1
	}
	w.X = x
	w.Y = y

	// if w or h are less than 1 set them to 1
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	w.Width = width
	w.Height = height
	w.ConsoleWidth = consoleWidth
	w.ConsoleHeight = consoleHeight

	w.pointMap = types.NewPointMap(w.ConsoleWidth, w.ConsoleHeight)
	w.lastSentContents = types.NewPointMap(w.ConsoleWidth, w.ConsoleHeight)
	w.pointMapInitialized = true

	for i := w.Y + 1; i < w.Y+w.Height+1; i++ {
		for j := w.X + 1; j < w.X+w.Width; j++ {
			if j > len(w.pointMap)-1 {
				return
			}
			if i > len(w.pointMap[i])-1 {
				return
			}
			w.pointMap[j][i] = types.Point{X: j, Y: i, Character: " ", EscapeCode: ""}
		}
	}
	w.lastSentContents = types.NewPointMap(w.ConsoleWidth, w.ConsoleHeight)
}
