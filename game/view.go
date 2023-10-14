package game

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/types"
	"math"
)

func (gm *GameManager) GetCharacterView(charID string, width, height int) (messages.GameCharView, error) {
	//gm.Log.Println(logging.LogInfo, "Getting character view dimensions", width, height)

	character, err := gm.GetCharacter(charID)
	if err != nil {
		gm.Log.Println(logging.LogInfo, "error getting character", err.Error())
		return messages.GameCharView{}, err
	}
	gm.activeCharactersMutex.Lock()
	defer gm.activeCharactersMutex.Unlock()

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

	view := messages.GameCharView{}
	// Load the character

	// Prepare the output view
	plane := make([][]types.Point, width)
	for i := range plane {
		plane[i] = make([]types.Point, height)
	}

	// A static set of coordinates for now
	posX := 75 //character.Position.X
	posY := 15 //character.Position.Y
	/*

		for i := range plane {
			for j := range plane[i] {
				distance := math.Sqrt(float64((i-posX)*(i-posX) + (j-posY)*(j-posY)))
				if distance <= float64(radius) {
					if gm.MapChunks[0].TileMap[i][j][0].Passable {
						plane[i][j].Character = "."
					} else {
						plane[i][j].Character = "#"
					}
				} else {
					plane[i][j].Character = " "
				}
			}
		}
	*/

	radius := 5

	offsetX := character.Position.X - width/2
	offsetY := character.Position.Y - height/2

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {

			// Calculate the position in string[x][y] from which to draw
			mapX := j + offsetX
			mapY := i + offsetY

			// Ensure you're within the bounds of string[x][y]
			if mapX >= 0 && mapX < len(gm.MapChunks[0].TileMap) && mapY >= 0 && mapY < len(gm.MapChunks[0].TileMap) {
				distance := math.Sqrt(float64((i-posX)*(i-posX) + (j-posY)*(j-posY)))
				if distance <= float64(radius) {
					if gm.MapChunks[0].TileMap[i][j][0].Passable {
						plane[i][j].Character = "."
					} else {
						plane[i][j].Character = "#"
					}
				} else {
					plane[i][j].Character = " "
				}
			} else {
				plane[i][j].Character = " "
			}
		}
	}

	// Set the character's position on the map
	plane[posX][posY].Character = "@" // This is the character
	plane[posX][posY].EscapeCode = character.FGColor.FG() + character.BGColor.BG()

	view.View = plane
	return view, nil

}
