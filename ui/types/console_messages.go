package types

import "encoding/json"

// ConsoleMessage is a message that is meant for parsing by the Console
// They can be sent by the ConnectionManager or by windows into the Console
type ConsoleMessage struct {
	Type        string `json:"type"`
	Message     string `json:"message"`
	Options     string `json:"options"`
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`
	WindowID    int    `json:"window_id"`
}

// String returns a JSON string representation of the ConsoleMessage
func (c *ConsoleMessage) String() string {
	out, _ := json.Marshal(c)
	return string(out)
}
