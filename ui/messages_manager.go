package ui

import (
	"github.com/yamamushi/EscapingEden/messages"
)

// CaptureManagerMessages is a goroutine that listens for messages from the ConnectionManager and parses them to determine
// Where they should go, or if any action should be taken from them.
func (c *Console) CaptureManagerMessages() {
	for {
		select {
		case consoleMessage := <-c.ReceiveMessages:
			//log.Println("Console received message from manager")

			switch consoleMessage.Type {
			case messages.Console_Message_RegistrationResponse:
				//log.Println("Console received registration response")
				loginMessage := messages.WindowMessage{Type: messages.WM_RegistrationResponse, Data: consoleMessage.Data}
				c.LoginWindowMessages <- loginMessage
			case messages.Console_Message_ResetPasswordValidateResponse:
				//log.Println("Console received reset password validation response")
				loginMessage := messages.WindowMessage{Type: messages.WM_PasswordResetValidateResponse, Data: consoleMessage.Data}
				c.LoginWindowMessages <- loginMessage
			case messages.Console_Message_ProcessPasswordValidateResponse:
				//log.Println("Console received process password update response")
				loginMessage := messages.WindowMessage{Type: messages.WM_PasswordResetProcessResponse, Data: consoleMessage.Data}
				c.LoginWindowMessages <- loginMessage
			case messages.Console_Message_LoginResponse:
				//log.Println("Console received login response, sending to login window")
				loginMessage := messages.WindowMessage{Type: messages.WM_LoginResponse, Data: consoleMessage.Data}
				c.LoginWindowMessages <- loginMessage
			case messages.Console_Message_ValidateCharNameResponse:
				//log.Println("Console received validate character name response")
				charCreatorMessage := messages.WindowMessage{
					Type: messages.WM_RequestCharNameValidationResponse,
					Data: consoleMessage.Data,
				}
				c.UserDashboardMessages <- charCreatorMessage
			//case messages.Console_Message_LoginUser:
			//log.Println("Console received login user request")
			//		userInfo := consoleMessage.Data.(messages.UserInfo)
			//		c.UpdateUserInfo(userInfo)
			//		c.ChatMessageReceive <- messages.ChatMessage{Type: messages.Chat_Message_System, Content: "You are now logged in as " + userInfo.Username}
			case messages.Console_Message_LogoutUser:
				//log.Println("Console received logout user request")
				if c.IsUserLoggedIn() {
					c.LogoutUser()
					c.ChatMessageReceive <- messages.ChatMessage{Type: messages.Chat_Message_System, Content: consoleMessage.Data.(string)}
				}
			case messages.Console_Message_Chat:
				//log.Println("Chat message received from manager")
				if c.IsCharacterLoggedIn() {
					chatMessage := messages.ChatMessage{Type: messages.Chat_Message_Normal, Content: consoleMessage.Data.(string)}
					c.ChatMessageReceive <- chatMessage
					continue
				}
			case messages.Console_Message_Error:
				//log.Println("Error message received from manager")
				chatMessage := messages.ChatMessage{Type: messages.Chat_Message_Normal, Content: "Error: " + consoleMessage.Data.(string)}
				c.ChatMessageReceive <- chatMessage
				continue
			case messages.Console_Message_Broadcast:
				//log.Println("Broadcast message received from manager")
				chatMessage := messages.ChatMessage{Type: messages.Chat_Message_Normal, Content: "Broadcast: " + consoleMessage.Data.(string)}
				c.ChatMessageReceive <- chatMessage
				continue
			case messages.Console_Message_Quit:
				//log.Println("Quit message received from manager")
				c.SetShutdown(true)
				continue
			}
		}
	}
}
