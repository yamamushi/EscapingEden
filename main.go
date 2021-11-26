package main

import (
	"github.com/yamamushi/EscapingEden/server"
	"log"
)

/*
Escaping Eden is a simple text adventure mud ;)
*/

func main() {
	log.Println("Starting server...")
	server := server.NewServer("localhost", "8080")
	server.Start()
	log.Println("Server exited")
}
