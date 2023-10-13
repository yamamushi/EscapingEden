package ui

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"log"
)

// CaptureWindowMessages is a goroutine that listens for messages from the windows and parses them to determine
// Where they should go, or if any action should be taken from them. Launching a popup box, or sending a message
// or even quitting the session. There are many types of messages the console may have to parse from a window.
func (c *Console) CaptureWindowMessages() {
	for {
		select {
		case windowMessage := <-c.WindowMessages:
			//log.Println("Console received window message")
			//log.Println("MessageType: ", consoleMessage.Type)
			switch windowMessage.Type {
			case messages.WM_ConsoleCommand:
				//log.Println("Received console message from a window")
				switch windowMessage.Command {
				case messages.WMC_NewPopup:
					c.OpenPopup(&windowMessage.PopupOptions)
					continue
					// These following messages are sent into their respective windows
				case messages.WMC_ClosePopup:
					c.HandlePopupMessage(windowMessage)
					continue
				case messages.WMC_ToggleHelp:
					c.ToggleHelp(&windowMessage.HelpOptions)
					continue
				case messages.WMC_RefreshConsole:
					c.AbortSend()
					c.ForceRedraw()
					continue
				//case messages.WMC_UpdateUserInfoForAllWindows:
				//	c.userInfoMutex.Lock()
				//	defer c.userInfoMutex.Unlock()
				//	c.UpdateWindowsUserInfo()
				//	continue
				case messages.WMC_FlushConsoleBuffer:
					//log.Println("Flush message received")
					c.flushWindowList = append(c.flushWindowList, windowMessage.TargetID)
					continue
				case messages.WMC_SetAccountLoggedIn:
					//log.Println("Console received login user for " + c.ConnectionID)
					c.LoginUser(windowMessage.Data.(messages.UserInfo))
					chatMessage := messages.ChatMessage{Type: messages.Chat_Message_System, Content: "You have logged in as " + c.GetUserName() + "."}
					c.ChatMessageReceive <- chatMessage
					continue
				case messages.WMC_SetAccountLoggedOut:
					//log.Println("Console received logout user for " + c.ConnectionID)
					c.LogoutUser() // This also logs out a character, no need to force both.
					chatMessage := messages.ChatMessage{Type: messages.Chat_Message_System, Content: "You have logged out."}
					c.ChatMessageReceive <- chatMessage
					continue
				case messages.WMC_RequestCharacterByID:
					//log.Println("Console received request for character by ID")
					// Request Character by ID from character manager
					requestByID := messages.ConnectionManagerMessage{
						Type:            messages.ConnectManager_Message_RequestCharacterByID,
						Data:            windowMessage.Data,
						SenderConsoleID: c.ConnectionID,
					}
					c.SendMessages <- requestByID

				case messages.WMC_RequestCharacterHistoryUpdate:
					// Tell account manager AND character manager that a character is attempting to log in
					loggedInMessage := messages.ConnectionManagerMessage{
						Type:            messages.ConnectManager_Message_CharacterLoggedInNotify,
						Data:            windowMessage.Data,
						SenderConsoleID: c.ConnectionID,
					}
					c.SendMessages <- loggedInMessage

				case messages.WMC_SetCharacterLoggedIn:
					//log.Println("Console received login character for " + c.ConnectionID)
					charInfo := windowMessage.Data.(messages.CharacterInfo)
					c.Log.Println(logging.LogInfo, "Setting character info for ", charInfo.Name, " to ", charInfo.ID)
					c.LoginCharacter(charInfo)
					chatMessage := messages.ChatMessage{}
					if int(charInfo.FirstLogin) == 1 {
						chatMessage = messages.ChatMessage{Type: messages.Chat_Message_System, Content: "Welcome " + c.GetCharacterName() + "!"}
					} else {
						chatMessage = messages.ChatMessage{Type: messages.Chat_Message_System, Content: "Welcome back " + c.GetCharacterName() + "!"}
					}
					c.ChatMessageReceive <- chatMessage

				case messages.WMC_SetCharacterLoggedOut:
					log.Println("Console received logout character for " + c.ConnectionID)
					c.LogoutCharacter()
					chatMessage := messages.ChatMessage{Type: messages.Chat_Message_System, Content: "You have logged out."}
					c.ChatMessageReceive <- chatMessage
					continue
				default:
					continue
				}

			case messages.WM_ParseChat:
				//log.Println("Sending Chat message: ", windowMessage.Data.(string))
				managerMessage := messages.ConnectionManagerMessage{
					Type:            messages.ConnectManager_Message_Chat,
					Data:            windowMessage.Data,
					SenderConsoleID: c.ConnectionID,
				}
				c.SendMessages <- managerMessage
				continue

			case messages.WM_GameCommand:
				//c.Log.Println(logging.LogInfo, "Sending Game Command: ", windowMessage.Data.(messages.GameManagerMessage).Type)
				managerMessage := messages.ConnectionManagerMessage{
					Type:            messages.ConnectManager_Message_GameCommand,
					Data:            windowMessage.Data,
					SenderConsoleID: c.ConnectionID,
				}
				c.SendMessages <- managerMessage
				continue

			case messages.WM_RequestRegistration:
				//log.Println("Sending registration request to connection manager")
				managerMessage := messages.ConnectionManagerMessage{
					Type:            messages.ConnectManager_Message_Register,
					Data:            windowMessage.Data,
					SenderConsoleID: c.ConnectionID,
				}
				c.SendMessages <- managerMessage
				// These messages require serializing to send to ConnectionManager

			case messages.WM_RequestLogin:
				//log.Println("Sending login request to connection manager")
				managerMessage := messages.ConnectionManagerMessage{
					Type:            messages.ConnectManager_Message_AccountLogin,
					Data:            windowMessage.Data,
					SenderConsoleID: c.ConnectionID,
				}
				c.SendMessages <- managerMessage
			// These messages require serializing to send to ConnectionManager

			case messages.WM_RequestCharNameValidation:
				//log.Println("Sending character name validation request to connection manager")
				managerMessage := messages.ConnectionManagerMessage{
					Type:            messages.ConnectManager_Message_CharNameValidation,
					Data:            windowMessage.Data,
					SenderConsoleID: c.ConnectionID,
				}
				c.SendMessages <- managerMessage

			case messages.WM_RequestCharacterCreation:
				//log.Println("Sending character creation request to connection manager")
				managerMessage := messages.ConnectionManagerMessage{
					Type:            messages.ConnectManager_Message_CharacterCreation,
					Data:            windowMessage.Data,
					SenderConsoleID: c.ConnectionID,
				}
				c.SendMessages <- managerMessage

			case messages.WM_RequestForgotPassword:
				managerMessage := messages.ConnectionManagerMessage{
					Type:            messages.ConnectManager_Message_RequestPasswordReset,
					Data:            windowMessage.Data,
					SenderConsoleID: c.ConnectionID,
				}
				c.SendMessages <- managerMessage

			case messages.WM_ValidateForgotPassword:
				managerMessage := messages.ConnectionManagerMessage{
					Type:            messages.ConnectManager_Message_ValidatePasswordReset,
					Data:            windowMessage.Data,
					SenderConsoleID: c.ConnectionID,
				}
				c.SendMessages <- managerMessage

			case messages.WM_ProcessForgotPassword:
				managerMessage := messages.ConnectionManagerMessage{
					Type:            messages.ConnectManager_Message_ProcessPasswordReset,
					Data:            windowMessage.Data,
					SenderConsoleID: c.ConnectionID,
				}
				c.SendMessages <- managerMessage

			case messages.WM_BadLoginAttempt:
				managerMessage := messages.ConnectionManagerMessage{
					Type:            messages.ConnectManager_Message_BadLoginAttempt,
					SenderConsoleID: c.ConnectionID,
				}
				c.SendMessages <- managerMessage

			case messages.WM_Error:
				//log.Println("Sending Error message to Connection Manager: ", windowMessage.Data.(string))
				managerMessage := messages.ConnectionManagerMessage{
					Type:            messages.ConnectManager_Message_Error,
					Data:            windowMessage.Data,
					SenderConsoleID: c.ConnectionID,
				}
				c.SendMessages <- managerMessage
				continue
			case messages.WM_QuitConsole:
				//log.Println("Sending Quit request to ConnectionManager")
				managerMessage := messages.ConnectionManagerMessage{
					Type:               messages.ConnectManager_Message_Quit,
					RecipientConsoleID: c.ConnectionID,
				}
				c.SendMessages <- managerMessage
				continue

			default:
				continue
			}
		}
	}
}
