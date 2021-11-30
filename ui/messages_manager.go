package ui

import (
	"encoding/json"
	"github.com/yamamushi/EscapingEden/ui/types"
	"log"
)

// CaptureManagerMessages is a goroutine that listens for messages from the ConnectionManager and parses them to determine
// Where they should go, or if any action should be taken from them.
func (c *Console) CaptureManagerMessages() {
	for {
		select {
		case message := <-c.ReceiveMessages:
			log.Println("Console received message from manager")
			consoleMessage := &types.ConsoleMessage{}
			err := json.Unmarshal([]byte(message), consoleMessage)
			if err != nil {
				log.Println("Error unmarshalling consoleMessage: ", err)
				continue
			}

			switch consoleMessage.Type {
			case "chat":
				log.Println("Chat message received from manager")
				c.ChatMessages <- consoleMessage.String()
				continue
			case "error":
				log.Println("Error message received from manager")
				consoleMessage.Message = "Error: " + consoleMessage.Message
				c.ChatMessages <- consoleMessage.String()
				continue
			case "broadcast":
				log.Println("Broadcast message received from manager")
				consoleMessage.Message = "Server Message: " + consoleMessage.Message
				c.ChatMessages <- consoleMessage.String()
				continue
			case "quit":
				log.Println("Quit message received from manager")
				c.SetShutdown(true)
				continue
			}
		}
	}
}
