package messages

type ConnectManagerMessageType int

const (
	ConnectManager_Message_Null ConnectManagerMessageType = iota
	ConnectManager_Message_Error
	ConnectManager_Message_Broadcast
	ConnectManager_Message_ServerShutdown
	ConnectManager_Message_Quit
	ConnectManager_Message_Chat
	ConnectManager_Message_Register
	ConnectManager_Message_RegisterResponse
	ConnectManager_Message_Login
	ConnectManager_Message_LoginResponse
)

// ConnectionManagerMessage is a message sent meant for parsing by the ConnectionManager
// It is identical to ConsoleMessage - This is so that either side can parse messages
// without having to know the other side.
type ConnectionManagerMessage struct {
	Type               ConnectManagerMessageType `json:"type"`
	SenderConsoleID    string                    `json:"sender_id"`
	RecipientConsoleID string                    `json:"recipient_id"`
	Data               interface{}               `json:"data"`
}

// MessageType is the type of message
type MessageType string

// MessageTypes are the types of messages

const (
	MessageTypeChat MessageType = "chat"
)
