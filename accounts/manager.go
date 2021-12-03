package accounts

import (
	"github.com/google/uuid"
	"github.com/yamamushi/EscapingEden/edendb"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

type AccountManager struct {
	ReceiveChannel chan messages.AccountManagerMessage    // We only receive player manager messages
	SendChannel    chan messages.ConnectionManagerMessage // We only send connection manager messages

	Log logging.LoggerType

	DB edendb.DatabaseType
}

func NewAccountManager(receiveChannel chan messages.AccountManagerMessage, sendChannel chan messages.ConnectionManagerMessage, db edendb.DatabaseType, log logging.LoggerType) *AccountManager {
	return &AccountManager{ReceiveChannel: receiveChannel, SendChannel: sendChannel, DB: db, Log: log}
}

// Init initializes the database for the account manager if needed
func (am *AccountManager) Init() error {

	id, _ := uuid.NewUUID()
	err := am.DB.AddRecord("Characters", &messages.Account{ID: id.String(), Username: "Test", HashedPassword: "Test"})
	if err != nil {
		if err.Error() == "already exists" {
			am.Log.Println(logging.LogInfo, "AccountManager", "Init", "Account already exists")
		}
	}
	_ = am.DB.UpdateRecord("Characters", &messages.Account{ID: id.String(), Username: "Test", HashedPassword: "Test"})

	//	search := messages.Account{Username: "Test"}
	result := messages.Account{}
	err = am.DB.One("Characters", "Username", "fdf", &result)
	if err != nil {
		am.Log.Println(logging.LogError, "Error finding account:", err)
		return err
	}

	am.Log.Println(logging.LogInfo, result.ID, result.Username, result.HashedPassword)

	return nil
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
