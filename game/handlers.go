package game

import (
	"github.com/yamamushi/EscapingEden/edenutil"
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
				//gm.Log.Println(logging.LogInfo, "Game Manager received login notification for ", charID)
				err := gm.LoadCharacter(charID, managerMessage.SenderConsoleID)
				if err != nil {
					gm.Log.Println(logging.LogError, "Game Manager failed to load character", err.Error())
					response := messages.ConnectionManagerMessage{
						Type:               messages.ConnectManager_Message_GameCommandResponse,
						RecipientConsoleID: managerMessage.SenderConsoleID,
						Data:               messages.GameMessage{Type: messages.GM_FailedLoadCharacter},
					}
					gm.SendChannel <- response
				}
				response := messages.ConnectionManagerMessage{
					Type: messages.ConnectManager_Message_Broadcast,
					Data: edenutil.EdenTime{}.CurrentTimeString() + " - " + gm.GetCharacterName(charID) + " entered the world.",
				}
				gm.SendChannel <- response

			case messages.GameManager_NotifyLoggedOutCharacter:
				charID := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).CharacterID
				if charID == "" {
					continue
				}
				gm.Log.Println(logging.LogInfo, "Game Manager received logout notification for ", charID)
				gm.RemoveFromLiveCharacterList(charID)

			case messages.GameManager_NotifyDisconnect:
				connectionID := managerMessage.Data.(string)
				//gm.Log.Println(logging.LogInfo, "Game Manager received disconnect notification for ", connectionID)
				gm.RemoveFromLiveCharacterList(connectionID)

			case messages.GameManager_RequestInventory:
				charID := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).CharacterID
				if charID == "" {
					continue
				}
				//gm.Log.Println(logging.LogInfo, "Game Manager received inventory request for ", charID)
				inventory, err := gm.GetCharacterInventory(charID)
				if err != nil {
					response := messages.ConnectionManagerMessage{
						Type:               messages.ConnectManager_Message_GameCommandResponse,
						RecipientConsoleID: managerMessage.SenderConsoleID,
						Data:               messages.GameMessage{Type: messages.GM_FailedLoadInventory},
					}
					gm.SendChannel <- response
				} else {
					response := messages.ConnectionManagerMessage{
						Type:               messages.ConnectManager_Message_GameCommandResponse,
						RecipientConsoleID: managerMessage.SenderConsoleID,
						Data: messages.GameMessage{Type: messages.GM_Inventory, Data: messages.GameMessageData{
							CharacterID: charID,
							Data:        inventory,
						},
						},
					}
					//gm.Log.Println(logging.LogInfo, "GameManager", "Sending inventory request response")
					gm.SendChannel <- response
				}

			case messages.GameManager_MoveCharacter:
				//gm.Log.Println(logging.LogInfo, "Game Manager received move request")
				charID := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).CharacterID
				if charID == "" {
					continue
				}
				//gm.Log.Println(logging.LogInfo, "Game Manager received move request for ", charID)
				deltax := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).Data.(messages.GameCharMove).DeltaX
				deltay := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).Data.(messages.GameCharMove).DeltaY
				gm.MovePlayer(charID, deltax, deltay)

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

			case messages.GameManager_DigCommand:
				//gm.Log.Println(logging.LogInfo, "Game Manager received dig request")
				charID := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).CharacterID
				if charID == "" {
					continue
				}
				deltaX := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).Data.(messages.GameCharDig).DeltaX
				deltaY := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).Data.(messages.GameCharDig).DeltaY
				itemID := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).Data.(messages.GameCharDig).ItemID
				err := gm.HandleDigRequest(itemID, charID, deltaX, deltaY)
				if err != nil {
					//gm.Log.Println(logging.LogError, "Game Manager failed to dig", err.Error())
					response := messages.ConnectionManagerMessage{
						Type:               messages.ConnectManager_Message_GameCommandResponse,
						RecipientConsoleID: managerMessage.SenderConsoleID,
						Data:               messages.GameMessage{Type: messages.GM_FailedDig},
					}
					gm.SendChannel <- response
				}

			case messages.GameManager_BuildWallCommand:
				//gm.Log.Println(logging.LogInfo, "Game Manager received dig request")
				charID := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).CharacterID
				if charID == "" {
					continue
				}
				deltaX := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).Data.(messages.GameCharBuildWall).DeltaX
				deltaY := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).Data.(messages.GameCharBuildWall).DeltaY
				itemID := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).Data.(messages.GameCharBuildWall).ItemID
				toolID := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).Data.(messages.GameCharBuildWall).ToolID
				err := gm.HandleBuildWallRequest(itemID, toolID, charID, deltaX, deltaY)
				if err != nil {
					//gm.Log.Println(logging.LogError, "Game Manager failed to dig", err.Error())
					response := messages.ConnectionManagerMessage{
						Type:               messages.ConnectManager_Message_GameCommandResponse,
						RecipientConsoleID: managerMessage.SenderConsoleID,
						Data:               messages.GameMessage{Type: messages.GM_FailedDig},
					}
					gm.SendChannel <- response
				}

			}
		}
	}
}
