package main

/*
Escaping Eden is a simple text adventure mud ;)
*/

import (
	"context"
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/yamamushi/EscapingEden/edenconfig"
	"github.com/yamamushi/EscapingEden/edenutil"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

const EscapingEdenVersion = "0.0.1"

// Variables used for command line parameters

var (
	ConfPath   string
	ResetWorld bool
)

// init is called before main()
func init() {
	// Read our command line options
	flag.StringVar(&ConfPath, "c", "server.conf", "Path to Config File")
	flag.BoolVar(&ResetWorld, "reset-world", false, "Delete and regenerate the world (requires confirmation)")
	flag.Parse()
}

// main is the entry point for Escaping Eden
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// run contains the main application logic with proper error handling
func run() error {
	// Validate config file exists
	if _, err := os.Stat(ConfPath); err != nil {
		return fmt.Errorf("config file missing: %s", ConfPath)
	}

	timestamp := time.Now().Format("01/02/2006 15:04:05")
	//go http.ListenAndServe("localhost:6060", nil)

	fmt.Println("Preparing to launch Escaping Eden v"+EscapingEdenVersion, "at", timestamp)
	fmt.Println("Reading config file at:", ConfPath+"\n")
	conf, err := edenconfig.ReadConfig(ConfPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	// Setup logging
	log, err := InitLogger(conf)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Handle world reset if requested
	if ResetWorld {
		err := handleWorldReset(conf, log)
		if err != nil {
			return fmt.Errorf("failed to reset world: %w", err)
		}
	}

	// Setup context for graceful shutdown (will be used in future improvements)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup context manager for coordinated shutdown
	ctxManager := edenutil.NewContextManager()
	defer func() {
		ctxManager.Cancel()
		ctxManager.Wait()
	}()

	// Setup database
	dbConn, err := InitDB(conf, log)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Setup Edenbot
	edenbotInput := make(chan messages.EdenbotMessage)
	edenbotOutput := make(chan messages.SystemManagerMessage) // TODO: this needs a whole new manager created
	edenBot, err := InitEdenbot(edenbotInput, edenbotOutput, dbConn, log, &conf)
	if err != nil {
		return fmt.Errorf("failed to initialize edenbot: %w", err)
	}

	// Setup channels for account manager and connection manager
	accountManagerReceiver := make(chan messages.AccountManagerMessage)
	connectionManagerReceive := make(chan messages.ConnectionManagerMessage)

	// Initialize account manager
	_, err = InitAccountManager(accountManagerReceiver, connectionManagerReceive, dbConn, log, edenBot, &conf)
	if err != nil {
		return fmt.Errorf("failed to initialize account manager: %w", err)
	}

	// Initialize the character manager
	characterManagerReceiver := make(chan messages.CharacterManagerMessage)
	_, err = InitCharacterManager(characterManagerReceiver, connectionManagerReceive, dbConn, &conf, log)
	if err != nil {
		return fmt.Errorf("failed to initialize character manager: %w", err)
	}

	gameManagerReceiver := make(chan messages.GameManagerMessage)
	gameManager, err := InitGameManager(gameManagerReceiver, connectionManagerReceive, dbConn, log, &conf)
	if err != nil {
		return fmt.Errorf("failed to initialize game manager: %w", err)
	}

	// Initialize the server, and by proxy, the connection manager
	server, err := InitServer(conf, accountManagerReceiver, characterManagerReceiver, connectionManagerReceive, edenbotInput, gameManagerReceiver, dbConn, log)
	if err != nil {
		return fmt.Errorf("failed to initialize server: %w", err)
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
	fmt.Println("Closing all client connections...")
	managerMessage = messages.ConnectionManagerMessage{Type: messages.ConnectManager_Message_ServerShutdown}
	server.ConnectionManagerSend <- managerMessage

	// Wait a shorter time for graceful disconnects (reduced to prevent hanging)
	gracefulDisconnectTime := time.Second * time.Duration(conf.Server.ShutdownTimeout/2)
	if gracefulDisconnectTime < 3*time.Second {
		gracefulDisconnectTime = 3 * time.Second
	}
	fmt.Printf("Waiting %.0f seconds for graceful disconnects...\n", gracefulDisconnectTime.Seconds())
	time.Sleep(gracefulDisconnectTime)

	// Run cleanup with built-in timeout protection
	fmt.Println("Running final cleanup...")
	gameManager.Cleanup()

	log.Println(logging.LogInfo, "Server exited cleanly.")
	if log.GetTypeID() != logging.LoggerTypeID_Console {
		fmt.Println("Server exited cleanly.")
	}

	return nil
}
