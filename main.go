package main

/*
Escaping Eden is a simple text adventure mud ;)
*/

import (
	"flag"
	"fmt"
	"github.com/yamamushi/EscapingEden/edenconfig"
	"github.com/yamamushi/EscapingEden/logging"
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
		fmt.Println("Config file is missing: ", ConfPath)
		flag.Usage()
		os.Exit(1)
	}
}

// main is the entry point for Escaping Eden
func main() {
	timestamp := time.Now().Format("01/02/2006 15:04:05")

	fmt.Println("Preparing to launch Escaping Eden v"+EscapingEdenVersion, "at", timestamp)
	fmt.Println("Reading config file at:", ConfPath)
	conf, err := edenconfig.ReadConfig(ConfPath)
	if err != nil {
		fmt.Println("Error reading config: ", err)
		os.Exit(1)
	}

	notifyDone := make(chan bool)
	if conf.Logger.Type != "console" {
		go InitAll(conf, notifyDone)
		ticker := time.NewTicker(100 * time.Millisecond)
		done := false
		for {
			select {
			case <-notifyDone:
				done = true
				break
			case <-ticker.C:
				fmt.Print(".")
			}
			if done {
				break
			}
		}
	} else {
		fmt.Println("[Info] Logging in console mode\n")
		InitAll(conf, notifyDone)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("\nEscaping Eden is now running. Press <ctrl-c> to exit.\n")

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-osSignal

	log.Println(logging.LogInfo, "Caught interrupt signal, shutting down...")
	log.Println(logging.LogInfo, "Server exited cleanly.")
}
