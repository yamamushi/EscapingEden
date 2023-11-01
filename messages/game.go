package messages

import (
	"github.com/yamamushi/EscapingEden/ui/types"
)

/*
Window messages are sent to/from the console from/to windows.
*/

type GameMessageType int

const (
	GM_Null GameMessageType = iota
	// These will send messages to the connection manager ultimately
	GM_Error
	GM_FailedLoadCharacter
	GM_FailedLoadView
	GM_QuitConsole
	GM_FailedLoadInventory
	GM_FailedDig
	GM_FailedBuildWall

	// Message to the game manager

	// Response types from the game manager
	GM_CharacterPosition
	GM_CharacterView
	GM_Inventory
	GM_Dig
	GM_BuildWall
)

type GameMessageCommand int

const (
	GMC_Null WindowMessageCommand = iota
)

type GameMessage struct {
	Type GameMessageType
	Data GameMessageData
}

type GameMessageData struct {
	CharacterID string
	Data        interface{}
}

type GameViewDimensions struct {
	Width  int
	Height int
}

type GameCharPosition struct {
	X int
	Y int
}

type GameCharView struct {
	View [][]types.Point
}

type GameCharMove struct {
	DeltaX int
	DeltaY int
}

type GameCharDig struct {
	DeltaX int
	DeltaY int
	ItemID string
}

type GameCharBuildWall struct {
	DeltaX int
	DeltaY int
	ItemID string
	ToolID string // Unused for now, but will be used for tools that are required to build
}
