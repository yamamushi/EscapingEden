package gamewindow

import (
	"github.com/yamamushi/EscapingEden/edenitems"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/terminals"
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"github.com/yamamushi/EscapingEden/ui/window"
	"sync"
)

// GameWindow is a window for users to login as a character, create a new one, manage their settings or log out.
type GameWindow struct {
	window.Window
	windowState GameWindowState

	// Initialize the window
	gwInitialized bool

	// Vars for navigation
	characterCreatorState CharacterCreatorState

	// Vars for command input
	commandMutex sync.Mutex

	// Character Info
	characterID    string
	characterMutex sync.Mutex

	log logging.LoggerType

	// Current Map inside the game window (upon a redraw we need to resize the map and redraw it too)
	visibleMap types.PointMap
	mapMutex   sync.Mutex

	FrameCounter       int
	FrameCounterMutext sync.Mutex

	StatusBarMessage string
	StatusBarMutex   sync.Mutex

	Menus      []MenuBoxType
	MenusMutex sync.Mutex
	CloseMenu  bool

	Inventory            []edenitems.Item
	InventoryMutex       sync.Mutex
	InventoryDisplayType edenitems.ItemType
}

// GameWindowState is an enum for storing game window state
type GameWindowState int

const (
	GW_NullState GameWindowState = iota
	GW_DefaultView
)

type CharacterCreatorState int

// NewGameWindow creates a new login window
func NewGameWindow(x, y, width, height, consoleWidth, consoleHeight int, input, output chan messages.WindowMessage,
	log logging.LoggerType, term terminals.TerminalType) *GameWindow {
	gw := &GameWindow{}
	gw.Log = log
	gw.Terminal = term
	gw.ID = config.WindowGameDisplay
	// if x or y are less than 1 set them to 1
	if x < 1 {
		x = 1
	}
	if y < 1 {
		y = 1
	}
	gw.X = x
	gw.Y = y

	// if w or h are less than 1 set them to 1
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	gw.Width = width
	gw.Height = height
	gw.ConsoleWidth = consoleWidth
	gw.ConsoleHeight = consoleHeight
	gw.Bordered = true
	gw.ConsoleReceive = input
	gw.ConsoleSend = output
	gw.windowState = GW_DefaultView
	gw.log = log
	//gw.log.Println(logging.LogInfo, "Character ID: ", gw.characterID)
	go gw.Listen()
	gw.SetupVisibleMap()
	return gw
}

// UpdateContents updates the contents of the login window
func (gw *GameWindow) UpdateContents() {
	switch gw.windowState {
	case GW_DefaultView:
		//gw.log.Println(logging.LogInfo, "Requesting Window View")
		gw.SendToConsole(messages.WindowMessage{Type: messages.WM_GameCommand, Data: messages.GameManagerMessage{Type: messages.GameManager_GetCharacterView, Data: messages.GameMessageData{CharacterID: gw.GetCharacterInfoField("id"), Data: messages.GameViewDimensions{Width: gw.Width, Height: gw.Height}}}})
		gw.PrintStringToMap(gw.X+1, gw.Y+1, "Game Window", gw.Terminal.Bold())
		gw.DrawStatusBar()
		gw.DrawMenus()
		if gw.CloseMenu {
			for _, menu := range gw.Menus {
				gw.RemoveMenuBox(menu)
			}
			gw.CloseMenu = false
		}

		// At center of window draw an @
		//gw.DrawToVisibleMap(gw.Width/2, (gw.Height/2)-1, "@", gw.CharacterInfo.FGColor.FG()+gw.CharacterInfo.BGColor.BG())
		gw.DrawMap()
		//xgw.RequestFlushFromConsole()
	}
}

func (gw *GameWindow) DrawStatusBar() {
	gw.StatusBarMutex.Lock()
	defer gw.StatusBarMutex.Unlock()
	messageLen := len(gw.StatusBarMessage)

	// Clear the status bar
	gw.ClearStatusBar()

	//gw.PrintStringToMap(gw.X+gw.Width-messageLen-2, gw.Y+gw.Height-4, gw.StatusBarMessage, gw.Terminal.Bold())
	gw.PrintStringToStatusBar(gw.Width-messageLen-2, 0, gw.StatusBarMessage, gw.Terminal.Bold())
}

func (gw *GameWindow) ClearStatusBar() {
	for i := 0; i < gw.Width; i++ {
		for j := 0; j < 4; j++ {
			gw.PrintStringToStatusBar(i, j, " ", "")
		}
	}
}

func (gw *GameWindow) PrintStringToStatusBar(x, y int, input string, escapeCode string) {
	// For every character in the input string, starting at x, y, print the character to the visible map
	// If x is greater than the width of the visible map, return
	if x > gw.Width-1 || x < 0 {
		return
	}
	// If y is greater than the height of the visible map, return
	if y > gw.Height-4 || y < 0 {
		return
	}
	for i, character := range input {
		// Using gw.DrawToVisibleMap for each point
		gw.DrawToVisibleMap(x+i, y+gw.Y+gw.Height-4, string(character), escapeCode)
	}
}

// Listen listens for any messages on cw.ReceiveMessages Chan and handles them
func (gw *GameWindow) Listen() {
	for {
		select {
		case receivedMessage := <-gw.ConsoleReceive:
			message := receivedMessage.Data.(messages.GameMessage).Type
			switch message {
			case messages.GM_CharacterPosition:
				//gw.log.Println(logging.LogInfo, "Game Window received message from console ", receivedMessage.Data.(messages.GameMessage).Data.Data)
				continue
			case messages.GM_CharacterView:
				//gw.log.Println(logging.LogInfo, "Game Window received view from console")
				gw.drawView(receivedMessage.Data.(messages.GameMessage).Data.Data.(messages.GameCharView))
			case messages.GM_Inventory:
				//gw.log.Println(logging.LogInfo, "Game Window received inventory from console")
				inventory := receivedMessage.Data.(messages.GameMessage).Data.Data.([]edenitems.Item)
				// cast the data to []*edenitems.Item
				gw.UpdateInventory(inventory)
				gw.DisplayInventory()
			}
		}
	}
}

func (gw *GameWindow) PrintStringToMap(x int, y int, input string, escapeCode string) {
	// For every character in the input string, starting at x, y, print the character to the visible map
	// If x is greater than the width of the visible map, return
	if x > len(gw.visibleMap)-1 || x < 0 {
		return
	}
	// If y is greater than the height of the visible map, return
	if y > len(gw.visibleMap[x])-1 || y < 0 {
		return
	}
	for i, character := range input {
		// Using gw.DrawToVisibleMap for each point
		gw.DrawToVisibleMap(x+i, y, string(character), escapeCode)
	}
}

func (gw *GameWindow) DrawToVisibleMap(x int, y int, character string, escapeCode string) {
	gw.mapMutex.Lock()
	defer gw.mapMutex.Unlock()

	if x > len(gw.visibleMap)-1 || x < 0 {
		return
	}
	if y > len(gw.visibleMap[x])-1 || y < 0 {
		return
	}
	gw.visibleMap[x][y] = types.Point{X: x, Y: y, Character: character, EscapeCode: escapeCode}
}

func (gw *GameWindow) DrawMap() {
	gw.mapMutex.Lock()
	defer gw.mapMutex.Unlock()

	for i := 0; i < gw.Width; i++ {
		for j := 0; j < gw.Height; j++ {
			gw.PrintLn(i+1, j+2, gw.visibleMap[i][j].Character, gw.visibleMap[i][j].EscapeCode)
		}
	}
}

func (gw *GameWindow) SetupVisibleMap() {
	gw.mapMutex.Lock()
	defer gw.mapMutex.Unlock()

	// Make a [][]Point of the size of the window
	gw.visibleMap = types.NewPointMap(gw.Width, gw.Height)

	// Fill with # for now
	for i := 0; i < gw.Width; i++ {
		for j := 0; j < gw.Height; j++ {
			if j < gw.Height-3 {
				gw.visibleMap[i][j] = types.Point{X: i, Y: j, Character: "#", EscapeCode: ""}
			} else {
				gw.visibleMap[i][j] = types.Point{X: i, Y: j, Character: " ", EscapeCode: ""}
			}
		}
	}

}

// UpdateParams is used when handling resize events to update the various window parameters in a safe state
func (gw *GameWindow) PostUpdateParams() {
	gw.SetupVisibleMap()
}
