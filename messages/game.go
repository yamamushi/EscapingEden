package messages

/*
Window messages are sent to/from the console from/to windows.
*/

type GameMessageType int

const (
	GM_Null GameMessageType = iota
	// These will send messages to the connection manager ultimately
	GM_Error
	GM_QuitConsole

	// Message to the game manager

	// Response types from the game manager
	GM_CharacterPosition
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

type GameCharPosition struct {
	X int
	Y int
}
