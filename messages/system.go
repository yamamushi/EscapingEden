package messages

// This file stores system message types, these go to the system manager
// Which in turn can parse them and send them to the appropriate system

type SystemManagerMessagType int

const (
	SystemManager_Message_Null SystemManagerMessagType = iota
	SystemManager_Message_Error
)

type SystemManagerMessage struct {
	Type SystemManagerMessagType `json:"type"`
}
