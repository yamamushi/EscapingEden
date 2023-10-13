package game

import (
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/types"
)

func (gm *GameManager) GetCharacterView(charID string) messages.GameCharView {

	view := messages.GameCharView{}

	// Load the map

	plane := make([][]types.Point, len(gm.MapChunks[0].TileMap))
	// We want to get the first layer of the map, z = 0 (the ground layer)
	for i := 0; i < len(gm.MapChunks[0].TileMap); i++ {
		plane[i] = make([]types.Point, len(gm.MapChunks[0].TileMap[0]))
		for j := 0; j < len(gm.MapChunks[0].TileMap[0]); j++ {
			if gm.MapChunks[0].TileMap[i][j][0].Passable {
				plane[i][j].Character = "."
			} else {
				plane[i][j].Character = "#"
			}
		}
	}

	view.View = plane
	return view

}
