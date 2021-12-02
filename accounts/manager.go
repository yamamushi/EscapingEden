package accounts

import (
	"github.com/yamamushi/EscapingEden/db"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

type AccountManager struct {
	ReceiveChannel chan messages.AccountManagerMessage    // We only receive player manager messages
	SendChannel    chan messages.ConnectionManagerMessage // We only send connection manager messages

	Log logging.LoggerType

	DB db.DatabaseType
}

func NewAccountManager(receiveChannel chan messages.AccountManagerMessage, sendChannel chan messages.ConnectionManagerMessage, db db.DatabaseType, log logging.LoggerType) *AccountManager {
	return &AccountManager{ReceiveChannel: receiveChannel, SendChannel: sendChannel, DB: db, Log: log}
}

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
				registrationResponse := messages.AccountRegistrationResponse{Success: true, Message: "Registration Successful"}
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
