package game

func (gm *GameManager) DrawRect(x1, y1, x2, y2 int) {

	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	// Draw a rectangle with the given coordinates to gm.MapChunks[0].TileMap
	// We want to get the first layer of the map, z = 0 (the ground layer)
	for x := 0; x < len(gm.MapChunks[0].TileMap); x++ {
		for y := 0; y < len(gm.MapChunks[0].TileMap[0]); y++ {
			// If x or y is the rectangle border, draw a wall.
			if (x == x1 || x == x2) && (y >= y1 && y <= y2) { // Left or right edges
				gm.MapChunks[0].TileMap[x][y][0] = Tile{TileType: "wall"}
			} else if (y == y1 || y == y2) && (x >= x1 && x <= x2) { // Top or bottom edges
				gm.MapChunks[0].TileMap[x][y][0] = Tile{TileType: "wall"}
			} /* else if x > x1 && x < x2 && y > y1 && y < y2 { // Inside the rectangle
				fmt.Print("^")
			}*/
		}
	}
}
