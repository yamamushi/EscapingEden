package game

import "time"

type World struct {
	// World is a struct for storing world data, because it has the potential to grow quite large
	// It instead contains references to other pieces of data
	// This is a work in progress
	ID           string
	RegionIDs    []string
	CreationTime time.Time
}

type MapChunk struct {
	ID string
	// The chunk's position in the world as a string [x,y,z]
	GlobalPosition struct {
		X int
		Y int
		Z int
	}
	// The chunk's data
	TileMap [][][]Tile
	//Mutex   sync.Mutex
}

type Tile struct {
	TileType string
}

type TileInfo struct {
	TileType     string
	Passable     bool
	BlocksVision bool
}
