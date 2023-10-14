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
			//gm.Log.Println(logging.LogInfo, "Game Manager received message")
			switch managerMessage.Type {
			case messages.GameManager_GetCharacterPosition: // Non functional
				charID := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).CharacterID
				gm.Log.Println(logging.LogInfo, "Game Manager received position request for ", charID)
				gm.Log.Println(logging.LogInfo, "Game Manager received position request from ", managerMessage.SenderConsoleID)

				// Do Something to get the position based on the provided character ID
				response := messages.ConnectionManagerMessage{
					Type:               messages.ConnectManager_Message_GameCommandResponse,
					RecipientConsoleID: managerMessage.SenderConsoleID,
					Data: messages.GameMessage{Type: messages.GM_CharacterPosition, Data: messages.GameMessageData{
						CharacterID: charID,
						Data:        "80, 80",
					},
					},
				}
				gm.Log.Println(logging.LogInfo, "GameManager", "Sending position request response")
				gm.SendChannel <- response

			case messages.GameManager_NotifyLoggedInCharacter:
				charID := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).CharacterID
				if charID == "" {
					continue
				}
				gm.Log.Println(logging.LogInfo, "Game Manager received login notification for ", charID)
				err := gm.LoadCharacter(charID)
				if err != nil {
					gm.Log.Println(logging.LogError, "Game Manager failed to load character", err.Error())
					response := messages.ConnectionManagerMessage{
						Type:               messages.ConnectManager_Message_GameCommandResponse,
						RecipientConsoleID: managerMessage.SenderConsoleID,
						Data:               messages.GameMessage{Type: messages.GM_FailedLoadCharacter},
					}
					gm.SendChannel <- response
				}

			case messages.GameManager_NotifyLoggedOutCharacter:
				charID := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).CharacterID
				if charID == "" {
					continue
				}
				gm.Log.Println(logging.LogInfo, "Game Manager received logout notification for ", charID)
				gm.RemoveFromLiveCharacterList(charID)

			case messages.GameManager_NotifyDisconnect:
				connectionID := managerMessage.Data.(string)
				gm.Log.Println(logging.LogInfo, "Game Manager received disconnect notification for ", connectionID)

			case messages.GameManager_GetCharacterView:
				//gm.Log.Println(logging.LogInfo, "Game Manager received character view request ")

				charID := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).CharacterID
				if charID == "" {
					continue
				}
				width := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).Data.(messages.GameViewDimensions).Width
				height := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).Data.(messages.GameViewDimensions).Height
				view, err := gm.GetCharacterView(charID, width, height)
				if err != nil {
					response := messages.ConnectionManagerMessage{
						Type:               messages.ConnectManager_Message_GameCommandResponse,
						RecipientConsoleID: managerMessage.SenderConsoleID,
						Data:               messages.GameMessage{Type: messages.GM_FailedLoadView},
					}
					gm.SendChannel <- response
				} else {
					response := messages.ConnectionManagerMessage{
						Type:               messages.ConnectManager_Message_GameCommandResponse,
						RecipientConsoleID: managerMessage.SenderConsoleID,
						Data: messages.GameMessage{Type: messages.GM_CharacterView, Data: messages.GameMessageData{
							CharacterID: charID,
							Data:        view,
						},
						},
					}
					//gm.Log.Println(logging.LogInfo, "GameManager", "Sending view request response")
					gm.SendChannel <- response
				}
			}
		}
	}
}
