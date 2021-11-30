package ui

import (
	"encoding/json"
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"log"
	"strconv"
)

// CaptureWindowMessages is a goroutine that listens for messages from the windows and parses them to determine
// Where they should go, or if any action should be taken from them. Launching a popup box, or sending a message
// or even quitting the session. There are many types of messages the console may have to parse from a window.
func (c *Console) CaptureWindowMessages() {
	for {
		select {
		case message := <-c.WindowMessages:
			log.Println("Console received window message")
			consoleMessage := &types.ConsoleMessage{}
			err := json.Unmarshal([]byte(message), consoleMessage)
			if err != nil {
				log.Println("Error unmarshalling consoleMessage: ", err)
				continue
			}
			//log.Println("MessageType: ", consoleMessage.Type)
			switch consoleMessage.Type {
			case "console":
				//log.Println("Received console message from a window")
				switch consoleMessage.Message {
				case "popup":
					options := &config.WindowConfig{}
					err = json.Unmarshal([]byte(consoleMessage.Options), options)
					if err != nil {
						log.Println("Error unmarshalling popup box options: ", err)
						continue
					} else {
						//log.Println("Popup box options: ", options.String())
						c.OpenPopup(options)
						continue
					}
				case "help":
					options := &config.WindowConfig{}
					err = json.Unmarshal([]byte(consoleMessage.Options), options)
					if err != nil {
						log.Println("Error unmarshalling help box options: ", err)
						continue
					} else {
						//log.Println("Help box options: ", options.String())
						c.ToggleHelp(options)
					}
				case "refresh":
					//log.Println("Refresh message received")
					//if !c.IsPopupOpen() {
					c.AbortSend()
					c.ForceRedraw()
					continue
					//}
				case "flush":
					log.Println("Flush message received")
					windowID, err := strconv.Atoi(consoleMessage.Options)
					if err != nil {
						log.Println("Could not parse flush options", err)
					}
					c.flushWindowList = append(c.flushWindowList, config.WindowID(windowID))
					continue

				default:
					continue
				}

			case "popupbox":
				log.Println("Popup box message: ", consoleMessage.Message)
				c.HandlePopupMessage(consoleMessage)
				continue
			case "help":
				log.Println("Help message: ", consoleMessage.Message)
				c.HandleHelpMessage(consoleMessage)
				continue
			case "error":
				log.Println("Error message: ", consoleMessage.Message)
				consoleMessage.RecipientID = c.ConnectionID
				c.SendMessages <- consoleMessage.String()
				continue
			case "quit":
				log.Println("Sending Quit request to ConnectionManager")
				consoleMessage.RecipientID = c.ConnectionID
				c.SendMessages <- consoleMessage.String()
				continue
			case "chat":
				log.Println("Chat message: ", consoleMessage.Message)
				consoleMessage.SenderID = c.ConnectionID
				c.SendMessages <- consoleMessage.String()
				continue
			default:
				continue
			}
		}
	}
}
