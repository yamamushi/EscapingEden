package main

import (
	server2 "github.com/yamamushi/EscapingEden/server"
	"log"
)

/*
Escaping Eden is a simple text adventure mud ;)
*/

func main() {
	/*
		dungeon := createDemo()
		// display the description for the player's room
		fmt.Print(dungeon.Player.Room.Describe())
		console := NewConsole(80, 25)
		console.Clear()
		for {
			console.Prompt()
			input := console.Read()
			err := ParseCommand(input, dungeon)
			if err != nil {
				fmt.Println(err)
			}
		}
	*/
	log.Println("Starting server...")
	server := server2.NewServer("localhost", "8080")
	server.Start()
	log.Println("Server exited")
}
