package main

/*
Escaping Eden is a simple text adventure mud ;)
*/

import (
	"flag"
	"fmt"
	"github.com/yamamushi/EscapingEden/edenconfig"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const EscapingEdenVersion = "0.0.1"

// Variables used for command line parameters

var (
	ConfPath string
)

// init is called before main()
func init() {
	// Read our command line options
	flag.StringVar(&ConfPath, "c", "server.conf", "Path to Config File")
	flag.Parse()

	_, err := os.Stat(ConfPath)
	if err != nil {
		fmt.Println("Config file is missing: ", ConfPath)
		flag.Usage()
		os.Exit(1)
	}
}

// main is the entry point for Escaping Eden
func main() {
	timestamp := time.Now().Format("01/02/2006 15:04:05")

	fmt.Println("Preparing to launch Escaping Eden v"+EscapingEdenVersion, "at", timestamp)
	fmt.Println("Reading config file at:", ConfPath+"\n")
	conf, err := edenconfig.ReadConfig(ConfPath)
	if err != nil {
		fmt.Println("Error reading config: ", err)
		os.Exit(1)
	}

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

	// Setup Edenbot
	edenbotInput := make(chan messages.EdenbotMessage)
	edenbotOutput := make(chan messages.SystemManagerMessage) // TODO: this needs a whole new manager created
	edenBot, err := InitEdenbot(edenbotInput, edenbotOutput, dbConn, log, &conf)
	if err != nil {
		log.Println(logging.LogFatal, "Error initializing edenbot: ", err)
	}

	// Setup channels for account manager and connection manager
	accountManagerReceiver := make(chan messages.AccountManagerMessage)
	connectionManagerReceive := make(chan messages.ConnectionManagerMessage)

	// Initialize account manager
	_, err = InitAccountManager(accountManagerReceiver, connectionManagerReceive, dbConn, log, *edenBot)
	if err != nil {
		// Fatal errors will os.Exit(1)
		log.Println(logging.LogFatal, "Error initializing account manager: ", err)
	}

	// Initialize the character manager
	characterManagerReceiver := make(chan messages.CharacterManagerMessage)
	_, err = InitCharacterManager(characterManagerReceiver, connectionManagerReceive, dbConn, &conf, log)

	// Initialize the server, and by proxy, the connection manager
	server, err := InitServer(conf, accountManagerReceiver, characterManagerReceiver, connectionManagerReceive, dbConn, log)
	if err != nil {
		log.Println(logging.LogFatal, "Error initializing server: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("\nEscaping Eden is now running. Press <ctrl-c> to exit.\n")

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-osSignal

	log.Println(logging.LogInfo, "Caught interrupt signal, shutting down...")
	if log.GetTypeID() != logging.LoggerTypeID_Console {
		fmt.Println("Caught interrupt signal, shutting down...")
	}
	// Issue a shutdown request to edenbot
	edenbotInput <- messages.EdenbotMessage{Type: messages.Edenbot_Message_Shutdown}

	// We need to notify our connections we're shutting down :D
	managerMessage := messages.ConnectionManagerMessage{
		Type: messages.ConnectManager_Message_Broadcast,
		Data: "Server shutting down in " + strconv.Itoa(conf.Server.ShutdownTimeout) + " seconds...",
	}
	server.ConnectionManagerSend <- managerMessage

	// We sleep for the configured ShutdownTimeout
	time.Sleep(time.Second * time.Duration(conf.Server.ShutdownTimeout))

	// Now we tell the connection manager we're shutting down and to close all connections
	managerMessage = messages.ConnectionManagerMessage{Type: messages.ConnectManager_Message_ServerShutdown}
	server.ConnectionManagerSend <- managerMessage

	// We sleep for the configured ShutdownTimeout
	time.Sleep(time.Second * time.Duration(conf.Server.ShutdownTimeout))

	log.Println(logging.LogInfo, "Server exited cleanly.")
	if log.GetTypeID() != logging.LoggerTypeID_Console {
		fmt.Println("Server exited cleanly.")
	}
}
