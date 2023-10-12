package game

import (
	"github.com/google/uuid"
	"github.com/yamamushi/EscapingEden/edendb"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

type GameManager struct {
	ReceiveChannel chan messages.GameManagerMessage       // We only receive player manager messages
	SendChannel    chan messages.ConnectionManagerMessage // We only send connection manager messages

	Log logging.LoggerType
	DB  edendb.DatabaseType
}

func NewGameManager(receiveChannel chan messages.GameManagerMessage, sendChannel chan messages.ConnectionManagerMessage, db edendb.DatabaseType, log logging.LoggerType) *GameManager {
	return &GameManager{ReceiveChannel: receiveChannel, SendChannel: sendChannel, DB: db, Log: log}
}

// Init initializes the database for the game manager if needed
func (gm *GameManager) Init() error {

	id, _ := uuid.NewUUID()
	err := gm.DB.AddRecord("Game", &World{ID: id.String()})
	if err != nil {
		if err.Error() == "already exists" {
			gm.Log.Println(logging.LogInfo, "GameManager", "Init", messages.GMError_DBError.Error())
		}
	}
	_ = gm.DB.UpdateRecord("Game", &World{ID: id.String()})

	//	search := messages.Account{Username: "Test"}
	result := World{}
	err = gm.DB.One("Game", "ID", id.String(), &result)
	if err != nil {
		gm.Log.Println(logging.LogError, messages.GMError_DBError.Error(), err)
		return err
	}

	return nil
}
