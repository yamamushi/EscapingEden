package main

/*
These functions are used to initialize various components, to keep main clean :)
*/
import (
	"errors"
	"github.com/yamamushi/EscapingEden/accounts"
	"github.com/yamamushi/EscapingEden/character"
	"github.com/yamamushi/EscapingEden/edenbot"
	"github.com/yamamushi/EscapingEden/edenconfig"
	"github.com/yamamushi/EscapingEden/edendb"
	"github.com/yamamushi/EscapingEden/edendb/bolt"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/logging/logconsole"
	"github.com/yamamushi/EscapingEden/logging/logfile"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/network"
	"strings"
	"time"
)

// InitLogger initializes the logger
func InitLogger(conf edenconfig.Config) (logging.LoggerType, error) {
	switch conf.Logger.Type {
	case "console":
		return logconsole.NewConsoleLogger(), nil
	case "file":
		return logfile.NewFileLogger(conf.Logger.Path)
	default:
		return nil, errors.New("Unknown logger type: " + conf.Logger.Type)
	}
}

// InitDB initializes the database
func InitDB(conf edenconfig.Config, log logging.LoggerType) (edendb.DatabaseType, error) {
	log.Println(logging.LogInfo, "Initializing Database connnection")
	if strings.ToLower(conf.DB.Type) == "bolt" {
		dbConn, err := bolt.NewBoltDB(conf.DB.Path)
		if err != nil {
			return nil, err
		}
		log.Println(logging.LogInfo, "Database connection initialized.")
		return dbConn, nil
	}
	return nil, errors.New("Invalid Database Type found - " + conf.DB.Type)
}

func InitEdenbot(input chan messages.EdenbotMessage, output chan messages.SystemManagerMessage, dbConn edendb.DatabaseType,
	log logging.LoggerType, conf *edenconfig.Config) (*edenbot.EdenBot, error) {
	log.Println(logging.LogInfo, "Starting Edenbot...")

	startNotify := make(chan bool)

	edenBot := edenbot.NewEdenBot(input, output, dbConn, log, conf)
	err := edenBot.Init()
	if err != nil {
		return nil, err
	}
	err = edenBot.Run(startNotify)
	if err != nil {
		return nil, err
	}

	ticker := time.NewTicker(1 * time.Second)
	select {
	case <-startNotify:
		log.Println(logging.LogInfo, "Eden Bot started.")
		break
	case <-ticker.C:
		//fmt.Print(".")
		// no-op
	}

	return edenBot, nil
}

// InitAccountManager initializes the account manager
func InitAccountManager(receiver chan messages.AccountManagerMessage, sender chan messages.ConnectionManagerMessage,
	dbConn edendb.DatabaseType, log logging.LoggerType, edenbot edenbot.EdenBot) (*accounts.AccountManager, error) {
	log.Println(logging.LogInfo, "Starting Account Manager...")

	startNotify := make(chan bool)

	accountManager := accounts.NewAccountManager(receiver, sender, dbConn, log, edenbot)
	// We need to init first because we need to ensure the database has been initialized
	// For the Account Manager to work.
	err := accountManager.Init()
	if err != nil {
		return nil, err
	}
	err = accountManager.Start(startNotify)
	if err != nil {
		return nil, err
	}

	ticker := time.NewTicker(1 * time.Second)

	select {
	case <-startNotify:
		log.Println(logging.LogInfo, "Account Manager started.")
		break
	case <-ticker.C:
		//fmt.Print(".")
		// no-op
	}
	return accountManager, nil
}

// InitCharacterManager initializes the character manager
func InitCharacterManager(input chan messages.CharacterManagerMessage, output chan messages.ConnectionManagerMessage, dbConn edendb.DatabaseType, conf *edenconfig.Config, log logging.LoggerType) (*character.CharacterManager, error) {
	log.Println(logging.LogInfo, "Starting Character Manager...")
	characterManager := character.NewCharacterManager(input, output, dbConn, log, conf)

	startNotify := make(chan bool)
	err := characterManager.Run(startNotify)
	if err != nil {
		return nil, err
	}

	ticker := time.NewTicker(1 * time.Second)
	select {
	case <-startNotify:
		log.Println(logging.LogInfo, "Character Manager started.")
		break
	case <-ticker.C:
		//fmt.Print(".")
		// no-op
	}
	return characterManager, nil
}

// InitServer initializes the server
func InitServer(conf edenconfig.Config,
	accountManagerReceive chan messages.AccountManagerMessage, characterManagerReceiver chan messages.CharacterManagerMessage,
	connectionManagerReceive chan messages.ConnectionManagerMessage,
	db edendb.DatabaseType, log logging.LoggerType) (*network.Server, error) {
	log.Println(logging.LogInfo, "Starting Server...")

	startNotify := make(chan bool)

	server := network.NewServer(conf.Server.Host, conf.Server.Port, log)
	err := server.Start(startNotify, connectionManagerReceive, accountManagerReceive, characterManagerReceiver, db)
	if err != nil {
		return nil, err
	}

	ticker := time.NewTicker(1 * time.Second)

	select {
	case <-startNotify:
		break
	case <-ticker.C:
		//fmt.Print(".")
		// no-op
	}
	return server, nil
}
