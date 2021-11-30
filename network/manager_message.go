package network

// ManagerMessage is a message sent meant for parsing by the ConnectionManager
// It is identical to ConsoleMessage - This is so that either side can parse messages
// without having to know the other side.
type ManagerMessage struct {
	Type        string `json:"type"`
	Message     string `json:"message"`
	Options     string `json:"options"`
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`
	WindowID    int    `json:"window_id"`
}

// MessageType is the type of message
type MessageType string

// MessageTypes are the types of messages

const (
	MessageTypeChat MessageType = "chat"
)
