package messages

/*
Window messages are sent to/from the console from/to windows.
*/

type GameMessageType int

const (
	GM_Null WindowMessageType = iota
	// These will send messages to the connection manager ultimately
	GM_Error
	GM_QuitConsole

	// This sends a message to the console, needed to notify console of various status changes
	GM_ConsoleCommand

	// These messages are sent to the connection manager for processing upstream
	// We can't do anything with them here, so we just pass them along
	GM_RequestCharacterPosition
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
