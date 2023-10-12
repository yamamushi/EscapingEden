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
