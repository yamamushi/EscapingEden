package messages

import "encoding/json"

type ConsoleMessageType int

const (
	Console_Message_Null ConsoleMessageType = iota
	Console_Message_Error
	Console_Message_Warning
	Console_Message_Log
	Console_Message_Quit
	Console_Message_Broadcast
	Console_Message_Chat // This needs to be refactored into the game engine channels, the connection manager does not need to know about chat.
	Console_Message_RegistrationRequest
	Console_Message_RegistrationResponse
)

// ConsoleMessage is a message that is meant for parsing by the Console
// They can be sent by the ConnectionManager or by windows into the Console
type ConsoleMessage struct {
	Type ConsoleMessageType
	Data interface{}
}

// String returns a JSON string representation of the ConsoleMessage
func (c *ConsoleMessage) String() string {
	out, _ := json.Marshal(c)
	return string(out)
}
