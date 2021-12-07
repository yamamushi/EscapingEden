package character

import (
	"github.com/yamamushi/EscapingEden/edenconfig"
	"github.com/yamamushi/EscapingEden/edendb"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

type CharacterManager struct {
	DB     edendb.DatabaseType
	Log    logging.LoggerType
	Config *edenconfig.Config

	InputChannel  chan messages.CharacterManagerMessage
	OutputChannel chan messages.ConnectionManagerMessage
}

func NewCharacterManager(input chan messages.CharacterManagerMessage, output chan messages.ConnectionManagerMessage, db edendb.DatabaseType,
	logger logging.LoggerType, conf *edenconfig.Config) *CharacterManager {
	return &CharacterManager{
		DB:            db,
		Log:           logger,
		Config:        conf,
		InputChannel:  input,
		OutputChannel: output,
	}
}

func (cm *CharacterManager) Run(startNotify chan bool) error {
	go cm.HandleInput(startNotify)
	return nil
}
