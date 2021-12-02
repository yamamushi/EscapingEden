package main

/*
These functions are used to initialize various components, to keep main clean :)
*/
import (
	"errors"
	"fmt"
	"github.com/yamamushi/EscapingEden/accounts"
	"github.com/yamamushi/EscapingEden/db"
	"github.com/yamamushi/EscapingEden/db/boltdb"
	"github.com/yamamushi/EscapingEden/edenconfig"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/logging/logconsole"
	"github.com/yamamushi/EscapingEden/logging/logfile"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/network"
	"os"
	"time"
)

// InitAll initializes all components
func InitAll(conf edenconfig.Config, notifyDone chan bool) {
	// Setup logging
	log, err := InitLogger(conf)
	if err != nil {
		fmt.Println("Error initializing logger: ", err)
		os.Exit(1)
	}

	// Setup database
	dbConn, err := InitDB(conf, log)
	if err != nil {
		log.Println(logging.LogFatal, "Error initializing database: ", err)
	}

	// Setup channels for account manager and connection manager
	accountManagerReceiver := make(chan messages.AccountManagerMessage)
	connectionManagerReceive := make(chan messages.ConnectionManagerMessage)

	// Initialize account manager
	_, err = InitAccountManager(accountManagerReceiver, connectionManagerReceive, dbConn, log)
	if err != nil {
		// Fatal errors will os.Exit(1)
		log.Println(logging.LogFatal, "Error initializing account manager: ", err)
	}

	// Initialize the server, and by proxy, the connection manager
	_, err = InitServer(conf, accountManagerReceiver, connectionManagerReceive, log)
	if err != nil {
		log.Println(logging.LogFatal, "Error initializing server: ", err)
	}

	// Non-blocking send to notifyDone
	select {
	case notifyDone <- true:
		//
	default:
		//
	}
}

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
func InitDB(conf edenconfig.Config, log logging.LoggerType) (db.DatabaseType, error) {
	log.Println(logging.LogInfo, "Initializing Database connnection")
	if conf.DB.Type == "bolt" {
		dbConn, err := boltdb.NewBoltDB(conf.DB.Path)
		if err != nil {
			return nil, err
		}
		return dbConn, nil
	}
	return nil, errors.New("Invalid Database Type found - " + conf.DB.Type)
}

// InitAccountManager initializes the account manager
func InitAccountManager(receiver chan messages.AccountManagerMessage, sender chan messages.ConnectionManagerMessage,
	dbConn db.DatabaseType, log logging.LoggerType) (*accounts.AccountManager, error) {
	log.Println(logging.LogInfo, "Starting Account Manager...")

	startNotify := make(chan bool)

	accountManager := accounts.NewAccountManager(receiver, sender, dbConn, log)
	err := accountManager.Start(startNotify)
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

// InitServer initializes the server
func InitServer(conf edenconfig.Config, accountManagerReceive chan messages.AccountManagerMessage,
	connectionManagerReceive chan messages.ConnectionManagerMessage, log logging.LoggerType) (*network.Server, error) {
	log.Println(logging.LogInfo, "Starting Server...")

	startNotify := make(chan bool)

	server := network.NewServer(conf.Server.Host, conf.Server.Port, log)
	err := server.Start(startNotify, connectionManagerReceive, accountManagerReceive)
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
