package accounts

import (
	"github.com/yamamushi/EscapingEden/messages"
	"log"
)

type AccountManager struct {
	ReceiveChannel chan messages.AccountManagerMessage    // We only receive player manager messages
	SendChannel    chan messages.ConnectionManagerMessage // We only send connection manager messages

}

func NewAccountManager(receiveChannel chan messages.AccountManagerMessage, sendChannel chan messages.ConnectionManagerMessage) *AccountManager {
	return &AccountManager{ReceiveChannel: receiveChannel, SendChannel: sendChannel}
}

func (am *AccountManager) Start(started chan bool) error {
	go am.HandleMessages(started)
	return nil
}

func (am *AccountManager) HandleMessages(started chan bool) {
	log.Println("Account Manager now handling messages")

	started <- true
	for {
		select {
		case managerMessage := <-am.ReceiveChannel:
			switch managerMessage.Type {
			case messages.AccountManager_Message_Register:
				registrationResponse := messages.AccountRegistrationResponse{Success: true, Message: "Registration Successful"}
				response := messages.ConnectionManagerMessage{
					Type:               messages.ConnectManager_Message_RegisterResponse,
					RecipientConsoleID: managerMessage.SenderSessionID,
					Data:               registrationResponse,
				}
				am.SendChannel <- response
			}
		}
	}
}
