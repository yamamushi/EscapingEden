package game

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

func (gm *GameManager) Start(started chan bool) error {
	go gm.HandleMessages(started)
	return nil
}

func (gm *GameManager) HandleMessages(started chan bool) {
	gm.Log.Println(logging.LogInfo, "Game Manager now handling messages")
	started <- true
	for {
		select {
		case managerMessage := <-gm.ReceiveChannel:
			gm.Log.Println(logging.LogInfo, "Game Manager received message")
			switch managerMessage.Type {
			case messages.GameManager_GetCharacterPosition:
				charID := managerMessage.Data.(messages.GameMessage).Data.CharacterID
				gm.Log.Println(logging.LogInfo, "Game Manager received position request for ", charID)
				// Do Something to get the position based on the provided character ID
				response := messages.ConnectionManagerMessage{
					Type:               messages.ConnectManager_Message_GameCommandResponse,
					RecipientConsoleID: managerMessage.SenderConsoleID,
					Data: messages.GameMessage{Data: messages.GameMessageData{
						CharacterID: charID,
						Data:        "Test response",
					},
					},
				}
				gm.Log.Println(logging.LogInfo, "GameManager", "Sending position request response")
				gm.SendChannel <- response

			}
		}
	}
}
