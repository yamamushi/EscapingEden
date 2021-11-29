package network

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
