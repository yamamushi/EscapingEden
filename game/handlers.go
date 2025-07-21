package game

import (
	"fmt"
	"github.com/yamamushi/EscapingEden/edenutil"
	"github.com/yamamushi/EscapingEden/game/commands"
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
					continue
				}

				// Ensure character is in a safe location when logging in
				err = gm.EnsureSafeLogin(charID)
				if err != nil {
					gm.Log.Println(logging.LogWarn, "Failed to ensure safe login for character:", err)
					// Even if safety check fails, let them login - they'll be in emergency chunk
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
				// Legacy movement handling - now redirected to action queue system
				charID := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).CharacterID
				if charID == "" {
					continue
				}

				moveData := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).Data.(messages.GameCharMove)
				deltax := moveData.DeltaX
				deltay := moveData.DeltaY

				// Convert to direction string for action queue
				direction := gm.deltaToDirection(deltax, deltay)

				// Queue movement action instead of immediate execution
				actionData := map[string]interface{}{
					"direction": direction,
					"deltaX":    deltax,
					"deltaY":    deltay,
				}

				queueAction := messages.GameQueueAction{
					CharacterID: charID,
					ActionType:  "move",
					Data:        actionData,
				}

				// Queue the action through the ActionManager
				if gm.ActionManager != nil {
					err := gm.ActionManager.QueueAction(queueAction.CharacterID, queueAction.ActionType, queueAction.Data)
					if err != nil {
						// Send error response
						response := messages.ConnectionManagerMessage{
							Type:               messages.ConnectManager_Message_GameCommandResponse,
							RecipientConsoleID: managerMessage.SenderConsoleID,
							Data: messages.GameMessage{
								Type: messages.GM_SystemMessage,
								Data: messages.GameMessageData{
									CharacterID: queueAction.CharacterID,
									Data:        fmt.Sprintf("Failed to queue movement: %v", err),
								},
							},
						}
						gm.SendChannel <- response
					}
				} else {
					// Fallback to old system if ActionManager is not available
					gm.MovePlayer(charID, deltax, deltay)
				}

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
				// Convert digging to action queue system
				charID := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).CharacterID
				if charID == "" {
					continue
				}

				digData := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).Data.(messages.GameCharDig)
				deltaX := digData.DeltaX
				deltaY := digData.DeltaY
				itemID := digData.ItemID

				// Queue digging action instead of immediate execution
				actionData := map[string]interface{}{
					"itemID": itemID,
					"deltaX": deltaX,
					"deltaY": deltaY,
				}

				queueAction := messages.GameQueueAction{
					CharacterID: charID,
					ActionType:  "dig",
					Data:        actionData,
				}

				// Queue the action through the ActionManager
				if gm.ActionManager != nil {
					err := gm.ActionManager.QueueAction(queueAction.CharacterID, queueAction.ActionType, queueAction.Data)
					if err != nil {
						// Send error response
						response := messages.ConnectionManagerMessage{
							Type:               messages.ConnectManager_Message_GameCommandResponse,
							RecipientConsoleID: managerMessage.SenderConsoleID,
							Data: messages.GameMessage{
								Type: messages.GM_SystemMessage,
								Data: messages.GameMessageData{
									CharacterID: queueAction.CharacterID,
									Data:        fmt.Sprintf("Failed to queue digging: %v", err),
								},
							},
						}
						gm.SendChannel <- response
					}
				} else {
					// Fallback to old system if ActionManager is not available
					err := gm.HandleDigRequest(itemID, charID, deltaX, deltaY)
					if err != nil {
						response := messages.ConnectionManagerMessage{
							Type:               messages.ConnectManager_Message_GameCommandResponse,
							RecipientConsoleID: managerMessage.SenderConsoleID,
							Data:               messages.GameMessage{Type: messages.GM_FailedDig},
						}
						gm.SendChannel <- response
					}
				}

			case messages.GameManager_BuildWallCommand:
				// Convert building to action queue system
				charID := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).CharacterID
				if charID == "" {
					continue
				}

				buildData := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).Data.(messages.GameCharBuildWall)
				deltaX := buildData.DeltaX
				deltaY := buildData.DeltaY
				itemID := buildData.ItemID
				toolID := buildData.ToolID

				// Queue building action instead of immediate execution
				actionData := map[string]interface{}{
					"itemID": itemID,
					"toolID": toolID,
					"deltaX": deltaX,
					"deltaY": deltaY,
				}

				queueAction := messages.GameQueueAction{
					CharacterID: charID,
					ActionType:  "build_wall",
					Data:        actionData,
				}

				// Queue the action through the ActionManager
				if gm.ActionManager != nil {
					err := gm.ActionManager.QueueAction(queueAction.CharacterID, queueAction.ActionType, queueAction.Data)
					if err != nil {
						// Send error response
						response := messages.ConnectionManagerMessage{
							Type:               messages.ConnectManager_Message_GameCommandResponse,
							RecipientConsoleID: managerMessage.SenderConsoleID,
							Data: messages.GameMessage{
								Type: messages.GM_SystemMessage,
								Data: messages.GameMessageData{
									CharacterID: queueAction.CharacterID,
									Data:        fmt.Sprintf("Failed to queue building: %v", err),
								},
							},
						}
						gm.SendChannel <- response
					}
				} else {
					// Fallback to old system if ActionManager is not available
					err := gm.HandleBuildWallRequest(itemID, toolID, charID, deltaX, deltaY)
					if err != nil {
						response := messages.ConnectionManagerMessage{
							Type:               messages.ConnectManager_Message_GameCommandResponse,
							RecipientConsoleID: managerMessage.SenderConsoleID,
							Data:               messages.GameMessage{Type: messages.GM_FailedBuildWall},
						}
						gm.SendChannel <- response
					}
				}

			case messages.GameManager_MineCommand:
				// Convert mining to action queue system
				charID := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).CharacterID
				if charID == "" {
					continue
				}

				mineData := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).Data.(messages.GameCharMine)
				deltaX := mineData.DeltaX
				deltaY := mineData.DeltaY
				itemID := mineData.ItemID
				toolID := mineData.ToolID

				// Queue mining action instead of immediate execution
				actionData := map[string]interface{}{
					"itemID": itemID,
					"toolID": toolID,
					"deltaX": deltaX,
					"deltaY": deltaY,
				}

				queueAction := messages.GameQueueAction{
					CharacterID: charID,
					ActionType:  "mine",
					Data:        actionData,
				}

				// Queue the action through the ActionManager
				if gm.ActionManager != nil {
					err := gm.ActionManager.QueueAction(queueAction.CharacterID, queueAction.ActionType, queueAction.Data)
					if err != nil {
						// Send error response
						response := messages.ConnectionManagerMessage{
							Type:               messages.ConnectManager_Message_GameCommandResponse,
							RecipientConsoleID: managerMessage.SenderConsoleID,
							Data: messages.GameMessage{
								Type: messages.GM_SystemMessage,
								Data: messages.GameMessageData{
									CharacterID: queueAction.CharacterID,
									Data:        fmt.Sprintf("Failed to queue mining: %v", err),
								},
							},
						}
						gm.SendChannel <- response
					}
				} else {
					// Fallback to old system if ActionManager is not available
					err := gm.HandleMineRequest(itemID, toolID, charID, deltaX, deltaY)
					if err != nil {
						response := messages.ConnectionManagerMessage{
							Type:               messages.ConnectManager_Message_GameCommandResponse,
							RecipientConsoleID: managerMessage.SenderConsoleID,
							Data: messages.GameMessage{
								Type: messages.GM_SystemMessage,
								Data: messages.GameMessageData{
									CharacterID: charID,
									Data:        fmt.Sprintf("Mining failed: %v", err),
								},
							},
						}
						gm.SendChannel <- response
					}
				}

			case messages.GameManager_AdminCommand:
				gm.Log.Println(logging.LogInfo, "Game Manager received admin command")
				senderID := managerMessage.SenderConsoleID
				commandString := managerMessage.Data.(messages.GameMessageData).Data.(string)

				// Parse and handle the admin command
				handled := gm.ParseAdminCommand(senderID, commandString)
				if !handled {
					// Send error message for unrecognized command
					response := messages.ConnectionManagerMessage{
						Type:               messages.ConnectManager_Message_GameCommandResponse,
						RecipientConsoleID: senderID,
						Data: messages.GameMessage{
							Type: messages.GM_SystemMessage,
							Data: messages.GameMessageData{
								CharacterID: senderID,
								Data:        "Unknown command. Type /adminhelp for available commands.",
							},
						},
					}
					gm.SendChannel <- response
				}

			case messages.GameManager_QueueAction:
				//gm.Log.Println(logging.LogInfo, "Game Manager received queue action request")
				queueAction := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).Data.(messages.GameQueueAction)

				// Queue the action through the ActionManager
				if gm.ActionManager != nil {
					err := gm.ActionManager.QueueAction(queueAction.CharacterID, queueAction.ActionType, queueAction.Data)
					if err != nil {
						// Send error response
						response := messages.ConnectionManagerMessage{
							Type:               messages.ConnectManager_Message_GameCommandResponse,
							RecipientConsoleID: managerMessage.SenderConsoleID,
							Data: messages.GameMessage{
								Type: messages.GM_SystemMessage,
								Data: messages.GameMessageData{
									CharacterID: queueAction.CharacterID,
									Data:        fmt.Sprintf("Failed to queue action: %v", err),
								},
							},
						}
						gm.SendChannel <- response
					}
				}

			case messages.GameManager_InterruptAction:
				//gm.Log.Println(logging.LogInfo, "Game Manager received interrupt action request")
				charID := managerMessage.Data.(messages.GameManagerMessage).Data.(messages.GameMessageData).CharacterID
				if charID == "" {
					continue
				}

				// Interrupt the current action through the ActionManager
				if gm.ActionManager != nil {
					err := gm.ActionManager.InterruptCurrentAction(charID)
					if err != nil {
						// Send error response
						response := messages.ConnectionManagerMessage{
							Type:               messages.ConnectManager_Message_GameCommandResponse,
							RecipientConsoleID: managerMessage.SenderConsoleID,
							Data: messages.GameMessage{
								Type: messages.GM_SystemMessage,
								Data: messages.GameMessageData{
									CharacterID: charID,
									Data:        fmt.Sprintf("Failed to interrupt action: %v", err),
								},
							},
						}
						gm.SendChannel <- response
					} else {
						// Send success response
						response := messages.ConnectionManagerMessage{
							Type:               messages.ConnectManager_Message_GameCommandResponse,
							RecipientConsoleID: managerMessage.SenderConsoleID,
							Data: messages.GameMessage{
								Type: messages.GM_SystemMessage,
								Data: messages.GameMessageData{
									CharacterID: charID,
									Data:        "Action interrupted.",
								},
							},
						}
						gm.SendChannel <- response
					}
				}

			case messages.GameManager_NewCommand:
				//gm.Log.Println(logging.LogInfo, "Game Manager received new command")
				cmd := managerMessage.Data.(*commands.PlayerCommand)

				// Process the command through the new system
				result, err := gm.ProcessCommand(cmd)

				// Send response back to the client
				var responseType messages.GameMessageType
				var responseData interface{}

				if err != nil || (result != nil && !result.Success) {
					// Command failed
					responseType = messages.GM_SystemMessage
					if result != nil && result.Error != nil {
						responseData = result.Error.Error()
					} else if err != nil {
						responseData = err.Error()
					} else {
						responseData = "Command failed"
					}
				} else if result != nil {
					// Command succeeded
					responseType = messages.GM_SystemMessage
					responseData = result.Message
				}

				// Send response
				response := messages.ConnectionManagerMessage{
					Type:               messages.ConnectManager_Message_GameCommandResponse,
					RecipientConsoleID: managerMessage.SenderConsoleID,
					Data: messages.GameMessage{
						Type: responseType,
						Data: messages.GameMessageData{
							CharacterID: cmd.PlayerID,
							Data:        responseData,
						},
					},
				}
				gm.SendChannel <- response

			}
		}
	}
}

// deltaToDirection converts movement deltas to direction strings
func (gm *GameManager) deltaToDirection(deltaX, deltaY int) string {
	switch {
	case deltaX == 0 && deltaY == -1:
		return "north"
	case deltaX == 0 && deltaY == 1:
		return "south"
	case deltaX == 1 && deltaY == 0:
		return "east"
	case deltaX == -1 && deltaY == 0:
		return "west"
	case deltaX == 1 && deltaY == -1:
		return "northeast"
	case deltaX == -1 && deltaY == -1:
		return "northwest"
	case deltaX == 1 && deltaY == 1:
		return "southeast"
	case deltaX == -1 && deltaY == 1:
		return "southwest"
	default:
		return "unknown"
	}
}
