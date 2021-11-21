package ui

type ConsoleMessage struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	Options  string `json:"options"`
	SenderID string `json:"sender_id"`
}
