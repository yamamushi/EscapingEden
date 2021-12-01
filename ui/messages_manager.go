package ui

import (
	"github.com/yamamushi/EscapingEden/messages"
	"log"
)

// CaptureManagerMessages is a goroutine that listens for messages from the ConnectionManager and parses them to determine
// Where they should go, or if any action should be taken from them.
func (c *Console) CaptureManagerMessages() {
	for {
		select {
		case consoleMessage := <-c.ReceiveMessages:
			log.Println("Console received message from manager")

			switch consoleMessage.Type {
			case messages.Console_Message_Chat:
				log.Println("Chat message received from manager")
				chatMessage := messages.ChatMessage{Type: messages.Chat_Message_Normal, Content: consoleMessage.Message}
				c.ChatMessageReceive <- chatMessage
				continue
			case messages.Console_Message_Error:
				log.Println("Error message received from manager")
				chatMessage := messages.ChatMessage{Type: messages.Chat_Message_Normal, Content: "Error: " + consoleMessage.Message}
				c.ChatMessageReceive <- chatMessage
				continue
			case messages.Console_Message_Broadcast:
				log.Println("Broadcast message received from manager")
				chatMessage := messages.ChatMessage{Type: messages.Chat_Message_Normal, Content: "Broadcast: " + consoleMessage.Message}
				c.ChatMessageReceive <- chatMessage
				continue
			case messages.Console_Message_Quit:
				log.Println("Quit message received from manager")
				c.SetShutdown(true)
				continue
			}
		}
	}
}
