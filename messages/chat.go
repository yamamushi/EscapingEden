package messages

type ChatMessageType int

const (
	Chat_Message_Null ChatMessageType = iota
	Chat_Message_Normal
	Chat_Message_Whisper
	Chat_Message_Broadcast
)

type ChatMessageChannel int

const (
	Chat_Message_Channel_Null ChatMessageChannel = iota
	Chat_Message_Channel_System
	Chat_Message_Channel_World
	Chat_Message_Channel_Local
	Chat_Message_Channel_Private
)

type ChatMessage struct {
	Type    ChatMessageType
	Content string
	Channel ChatMessageChannel
}
