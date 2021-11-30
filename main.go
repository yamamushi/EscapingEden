package main

/*
Escaping Eden is a simple text adventure mud ;)
*/

import (
	"flag"
	"fmt"
	"github.com/yamamushi/EscapingEden/edenconfig"
	"github.com/yamamushi/EscapingEden/network"
	"log"
	"os"
	"os/signal"
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
		log.Println("Config file is missing: ", ConfPath)
		flag.Usage()
		os.Exit(1)
	}
}

// main is the entry point for Escaping Eden
func main() {
	log.Println("Preparing to launch Escaping Eden v" + EscapingEdenVersion)

	log.Println("Reading config file at ", ConfPath)
	conf, err := edenconfig.ReadConfig(ConfPath)
	if err != nil {
		log.Println("Error reading config: ", err)
		os.Exit(1)
	}

	server := network.NewServer(conf.Server.Host, conf.Server.Port)
	log.Println("Starting Escaping Eden server...")
	startNotify := make(chan bool)
	err = server.Start(startNotify)
	if err != nil {
		log.Println("Error starting server: ", err)
		os.Exit(1)
	}
	ticker := time.NewTicker(1 * time.Second)
	select {
	case <-startNotify:
		log.Println("Escaping Eden is now running.  Press CTRL-C to exit.")
		break
	case <-ticker.C:
		fmt.Print(".")
	}

	// Wait here until CTRL-C or other term signal is received.

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-osSignal
	log.Println("Caught interrupt signal, shutting down...")
	log.Println("Server exited cleanly.")
}
