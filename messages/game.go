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
	GM_LoadingMessage

	// Message to the game manager

	// Response types from the game manager
	GM_CharacterPosition
	GM_CharacterView
	GM_Inventory
	GM_Dig
	GM_BuildWall
	GM_SystemMessage
	GM_AdminTeleportResponse
	GM_AdminSpawnListResponse
	GM_AdminLocationResponse
	GM_ActionQueueUpdate
	GM_ActionStarted
	GM_ActionCompleted
	GM_ActionFailed
)

type GameMessageCommand int

const (
	GMC_Null GameMessageCommand = iota
)

type GameMessage struct {
	Type    GameMessageType
	Message string // For simple text messages like loading messages
	Data    GameMessageData
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

type GameCharMine struct {
	DeltaX int
	DeltaY int
	ItemID string // Expected resource to mine
	ToolID string // Mining tool to use
}

type GameAdminTeleport struct {
	AdminID        string
	TargetPlayerID string
	Location       string
}

type GameAdminListSpawns struct {
	AdminID string
}

type GameAdminGetLocation struct {
	AdminID        string
	TargetPlayerID string
}

// Action queue related message types
type GameActionInfo struct {
	Type          string  `json:"type"`
	Description   string  `json:"description"`
	Progress      float64 `json:"progress"` // 0.0 to 1.0
	TicksLeft     int     `json:"ticksLeft"`
	TotalTicks    int     `json:"totalTicks"`
	Interruptible bool    `json:"interruptible"`
}

type GameActionQueueUpdate struct {
	CharacterID   string            `json:"characterId"`
	CurrentAction *GameActionInfo   `json:"currentAction"`
	QueuedActions []*GameActionInfo `json:"queuedActions"`
	MaxQueue      int               `json:"maxQueue"`
}

type GameActionStarted struct {
	CharacterID string `json:"characterId"`
	ActionType  string `json:"actionType"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
}

type GameActionCompleted struct {
	CharacterID string      `json:"characterId"`
	ActionType  string      `json:"actionType"`
	Success     bool        `json:"success"`
	Result      interface{} `json:"result"`
}

type GameActionFailed struct {
	CharacterID string `json:"characterId"`
	ActionType  string `json:"actionType"`
	Reason      string `json:"reason"`
}

// GameQueueAction represents a request to queue an action
type GameQueueAction struct {
	CharacterID string      `json:"characterId"`
	ActionType  string      `json:"actionType"`
	Data        interface{} `json:"data"`
}
