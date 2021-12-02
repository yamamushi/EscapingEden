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
	"syscall"
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
	fmt.Println("Preparing to launch Escaping Eden v" + EscapingEdenVersion)
	fmt.Println("Reading config file at:", ConfPath)
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

	// Wait here until CTRL-C or other term signal is received.
	log.Println(logging.LogInfo, "Escaping Eden is now running. Press <ctrl-c> to exit.")

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-osSignal
	
	log.Println(logging.LogInfo, "Caught interrupt signal, shutting down...")
	log.Println(logging.LogInfo, "Server exited cleanly.")
}
