package messages

type EdenbotMessageType int

const (
	Edenbot_Message_Null EdenbotMessageType = iota
	Edenbot_Message_Error
	Edenbot_Message_NewRegistration
	Edenbot_Message_PlayerLoggedIn
	Edenbot_Message_PlayerLoggedOut
	Edenbot_Message_CharacterCreated
	Edenbot_Message_ForgotPassword
	Edenbot_Message_Shutdown
)

// EdenbotMessage is a message sent meant for parsing by Edenbot.
type EdenbotMessage struct {
	Type       EdenbotMessageType `json:"type"`
	SourceType string             `json:"source_type"`
	SourceID   string             `json:"source_id"`
	Data       interface{}        `json:"data"`
}

type EdenbotErrorType int

const (
	Edenbot_Error_Null EdenbotErrorType = iota
	Edenbot_Error_DB
)
