package accounts

import (
	"github.com/google/uuid"
	"github.com/yamamushi/EscapingEden/edenbot"
	"github.com/yamamushi/EscapingEden/edendb"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

type AccountManager struct {
	ReceiveChannel chan messages.AccountManagerMessage    // We only receive player manager messages
	SendChannel    chan messages.ConnectionManagerMessage // We only send connection manager messages

	Log logging.LoggerType
	EB  edenbot.EdenBot

	DB edendb.DatabaseType
}

func NewAccountManager(receiveChannel chan messages.AccountManagerMessage, sendChannel chan messages.ConnectionManagerMessage, db edendb.DatabaseType, log logging.LoggerType, edenbot edenbot.EdenBot) *AccountManager {
	return &AccountManager{ReceiveChannel: receiveChannel, SendChannel: sendChannel, DB: db, Log: log, EB: edenbot}
}

// Init initializes the database for the account manager if needed
func (am *AccountManager) Init() error {

	id, _ := uuid.NewUUID()
	err := am.DB.AddRecord("Accounts", &messages.Account{ID: id.String(), Username: "Test", HashedPassword: "Test"})
	if err != nil {
		if err.Error() == "already exists" {
			am.Log.Println(logging.LogInfo, "AccountManager", "Init", messages.AMError_AccountAlreadyExists.Error())
		}
	}
	_ = am.DB.UpdateRecord("Characters", &messages.Account{ID: id.String(), Username: "Test", HashedPassword: "Test"})

	//	search := messages.Account{Username: "Test"}
	result := messages.Account{}
	err = am.DB.One("Accounts", "Username", "Test", &result)
	if err != nil {
		am.Log.Println(logging.LogError, messages.AMError_AccountDoesNotExist.Error(), err)
		return err
	}

	//am.Log.Println(logging.LogInfo, result.ID, result.Username, result.HashedPassword)

	return nil
}
