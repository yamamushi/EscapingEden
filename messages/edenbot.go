package messages

type EdenbotMessageType int

const (
	Edenbot_Message_Null EdenbotMessageType = iota
	Edenbot_Message_Error
	Edenbot_Message_NewRegistration
	Edenbot_Message_PlayerLoggedIn
	Edenbot_Message_PlayerLoggedOut
	Edenbot_Message_CharacterCreated
	Edenbot_Message_Shutdown
)

// EdenbotMessage is a message sent meant for parsing by Edenbot.
type EdenbotMessage struct {
	Type       EdenbotMessageType `json:"type"`
	SourceType string             `json:"source_type"`
	Data       interface{}        `json:"data"`
}
