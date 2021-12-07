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
				req := managerMessage.RegistrationRequest
				am.Log.Println(logging.LogInfo, "Account Manager received registration request")
				registrationResponse := messages.AccountRegistrationResponse{Error: am.CreateAccount(req.Username, req.Password, req.Email)}
				response := messages.ConnectionManagerMessage{
					Type:               messages.ConnectManager_Message_RegisterResponse,
					RecipientConsoleID: managerMessage.SenderSessionID,
					Data:               registrationResponse,
				}
				am.Log.Println(logging.LogInfo, "Sending registration response")
				am.SendChannel <- response
			}
		}
	}
}
