package main

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
	server := NewServer("localhost", "8080")
	server.Start()
}
