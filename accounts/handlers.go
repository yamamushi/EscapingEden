package accounts

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

func (am *AccountManager) Start(started chan bool) error {
	go am.HandleMessages(started)
	return nil
}

func (am *AccountManager) HandleMessages(started chan bool) {
	am.Log.Println(logging.LogInfo, "Account Manager now handling messages")
	started <- true
	for {
		select {
		case managerMessage := <-am.ReceiveChannel:
			am.Log.Println(logging.LogInfo, "Account Manager received message")
			switch managerMessage.Type {
			case messages.AccountManager_Message_Register:
				req := managerMessage.Data.(messages.AccountRegistrationRequest)
				am.Log.Println(logging.LogInfo, "Account Manager received registration request")
				registrationResponse := am.CreateAccount(req.Username, req.Password, req.DiscordID)
				response := messages.ConnectionManagerMessage{
					Type:               messages.ConnectManager_Message_RegisterResponse,
					RecipientConsoleID: managerMessage.SenderSessionID,
					Data:               registrationResponse,
				}
				am.Log.Println(logging.LogInfo, "Sending registration response")
				am.SendChannel <- response

			case messages.AccountManager_Message_Login:
				req := managerMessage.Data.(messages.AccountLoginRequest)
				am.Log.Println(logging.LogInfo, "Account Manager received login request")

				loginResponse := am.handleLogin(req.Username, req.Password)
				response := messages.ConnectionManagerMessage{
					Type:               messages.ConnectManager_Message_LoginResponse,
					RecipientConsoleID: managerMessage.SenderSessionID,
					Data:               loginResponse,
				}
				am.Log.Println(logging.LogInfo, "Sending login response")
				am.SendChannel <- response
			}
		}
	}
}
