package ui

import "encoding/json"

type ConsoleMessage struct {
	Type        string `json:"type"`
	Message     string `json:"message"`
	Options     string `json:"options"`
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`
}

func (c *ConsoleMessage) String() string {
	out, _ := json.Marshal(c)
	return string(out)
}
