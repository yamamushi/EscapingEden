package game

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/types"
)

func (gm *GameManager) GetCharacterView(charID string, width, height int) (messages.GameCharView, error) {
	//gm.Log.Println(logging.LogInfo, "Getting character view dimensions", width, height)

	character, err := gm.GetCharacter(charID)
	if err != nil {
		//gm.Log.Println(logging.LogInfo, "error getting character", err.Error())
		return messages.GameCharView{}, err
	}
	gm.activeCharactersMutex.Lock()
	posX := character.Position.X
	posY := character.Position.Y
	charSymbol := "@"
	charEscapeCode := character.FGColor.FG() + character.BGColor.BG()

	// If character is not initialized, we need to load it into the map
	if !character.Initialized {
		gm.Log.Println(logging.LogInfo, "character not initialized, loading into map")
		character.Position = struct {
			X int
			Y int
			Z int
		}{X: len(gm.MapChunks[0].TileMap) / 2, Y: len(gm.MapChunks[0].TileMap[0]) / 2, Z: 0}
		character.Initialized = true
	}
	defer gm.activeCharactersMutex.Unlock()

	view := messages.GameCharView{}

	// Prepare the output view
	plane := make([][]types.Point, width)
	for i := range plane {
		plane[i] = make([]types.Point, height)
	}

	// plane[i][j] is the window drawing we're sending to the client
	// We want to center the view around the player's position, posX and posY
	// So we need to offset the current MapChunk starting draw position by the player's position so that
	// When they are drawn they are always at the center and the map is drawn around them

	// log width and height
	//gm.Log.Println(logging.LogInfo, "width", width, "height", height)

	offsetX := width / 2
	offsetY := height / 2
	radius := 4

	// Now we loop through the plane, do our checks for each point and draw
	// Prepare vars
	tilemapXLen := len(gm.MapChunks[0].TileMap)
	tilemapYLen := len(gm.MapChunks[0].TileMap[0])
	// Loop through the screen
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			// Calculate which part of the map we're drawing based on player's position and screen position
			mapX := posX - offsetX + j
			mapY := posY - offsetY + i

			// Check if the map coordinate is within bounds
			if mapX >= 0 && mapX < tilemapXLen &&
				mapY >= 0 && mapY < tilemapYLen {
				// If we're at the player position, draw the player
				if mapX == posX && mapY == posY {
					//gm.Log.Println(logging.LogInfo, "drawing player at", mapX, mapY, j, i)
					plane[j][i].Character = charSymbol
					plane[j][i].EscapeCode = charEscapeCode
				} else {
					distanceSquared := float64((mapX-posX)*(mapX-posX) + (mapY-posY)*(mapY-posY))
					if distanceSquared <= float64(radius*radius) {
						playercheck := gm.GetCharacterAt(mapX, mapY)
						if playercheck != nil && playercheck.ID != charID {
							plane[j][i].Character = "@"
							plane[j][i].EscapeCode = playercheck.FGColor.FG() + playercheck.BGColor.BG()
						} else {
							if gm.MapChunks[0].TileMap[mapX][mapY][0].Passable {
								plane[j][i].Character = "."
							} else {
								plane[j][i].Character = "#"
							}
						}
					}
				}
			} else {
				plane[j][i].Character = "#" // Draw out of bounds as walls for simplicity
			}
		}
	}
	/*

	   |-----------------|
	   |                 |
	   |                 |
	   |   posX, posY    |
	   |                 |
	   |-----------------|

	*/

	// Set the character's position on the map
	//plane[posX][posY].Character = "@" // This is the character
	//plane[posX][posY].EscapeCode = character.FGColor.FG() + character.BGColor.BG()

	view.View = plane
	return view, nil

}
