package window

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/terminals"
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"sync"
)

type WindowType interface {
	ClearMap(X, Y, height, width, startX, startY int)
	Draw(X, Y int)
	HandleInput(input types.Input)
	Init()

	DrawBorder(X, Y int)
	DrawContents(X, Y int)
	UpdateContents()
	SetContents(string)
	PrintLn(X, Y int, text string, escapeCode string)
	GetPointMap() types.PointMap
	FlushLastSent()
	ResetWindowDrawings()

	GetID() config.WindowID
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
	GetConfig() *config.WindowConfig

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
	ID config.WindowID

	Log      logging.LoggerType
	Terminal terminals.TerminalType

	X int // The X position of the Window
	Y int // The Y position of the Window

	StartX int // When window content is rendered, it is a 2D array, so this is the starting X position of the content
	StartY int // When window content is rendered, it is a 2D array, so this is the starting Y position of the content

	Contents            string // The contents of the window
	ContentStartPos     int    // The starting position of the content
	ScrollingSupported  bool
	ScrollBufferHasNew  bool // Indicates that the scroll buffer has new content
	ScrollBufferLimit   int  // The maximum number of lines that can be stored in the scroll buffer
	ScrollBufferCharMod int  // The character limiter of the scroll buffer

	// For Window Paging, only WindowHelpBox uses this
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
	ConsoleSend    chan messages.WindowMessage // Send messages to the Console
	ConsoleReceive chan messages.WindowMessage // Receive messages from the Console

	DirectionInput types.InputType // The last direction input from the user

	mutex               sync.Mutex
	pmapMutex           sync.Mutex
	pointMap            types.PointMap
	lastSentContents    types.PointMap // The last contents sent to the client
	pointMapInitialized bool
}

func (w *Window) Init() {
	w.pointMap = types.NewPointMap(w.ConsoleWidth, w.ConsoleHeight)
	w.lastSentContents = types.NewPointMap(w.ConsoleWidth, w.ConsoleHeight)
	w.pointMapInitialized = true
}

func (w *Window) HandleInput(input types.Input) {
	return
}

// These functions implement the default WindowType interface for Window

// GetID returns the ID of the Window (types.WindowType)
func (w *Window) GetID() config.WindowID {
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

// UpdateContents updates the contents of the window
func (w *Window) UpdateContents() {
	return
}

// GetX returns the X position of the window
func (w *Window) GetX() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.X
}

// GetY returns the Y position of the window
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

func (w *Window) GetConfig() *config.WindowConfig {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return config.NewWindowConfig(w.X, w.Y, w.Width, w.Height, w.Contents)
}
