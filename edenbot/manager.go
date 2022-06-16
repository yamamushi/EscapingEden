package edenbot

import (
	"github.com/yamamushi/EscapingEden/edenconfig"
	"github.com/yamamushi/EscapingEden/edendb"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

type EdenBot struct {
	DB     edendb.DatabaseType
	Log    logging.LoggerType
	Config *edenconfig.Config

	InputChannel  chan messages.EdenbotMessage
	OutputChannel chan messages.SystemManagerMessage
}

func NewEdenBot(input chan messages.EdenbotMessage, output chan messages.SystemManagerMessage, db edendb.DatabaseType,
	logger logging.LoggerType, conf *edenconfig.Config) *EdenBot {
	return &EdenBot{
		DB:            db,
		Log:           logger,
		Config:        conf,
		InputChannel:  input,
		OutputChannel: output,
	}
}

func (eb *EdenBot) Init() error {
	return nil //TODO
}

func (eb *EdenBot) Run(startNotify chan bool) error {
	// First we start the discord bot using the credentials from eb.config
	// Then we start the manager service

	go eb.HandleInput(startNotify)
	return nil
}
