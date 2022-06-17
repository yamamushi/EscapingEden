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
				//log.Println("Console received login response")
				loginMessage := messages.WindowMessage{Type: messages.WM_LoginResponse, Data: consoleMessage.Data}
				c.LoginWindowMessages <- loginMessage
			case messages.Console_Message_Chat:
				//log.Println("Chat message received from manager")
				chatMessage := messages.ChatMessage{Type: messages.Chat_Message_Normal, Content: consoleMessage.Data.(string)}
				c.ChatMessageReceive <- chatMessage
				continue
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
