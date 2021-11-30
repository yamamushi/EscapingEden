package types

import "strconv"

// Point is a 2D point.
type Point struct {
	X, Y       int
	EscapeCode string
	Character  string
}

// PointMap is a map of points.
type PointMap [][]Point

// Print prints a point as an escaped string.
func (p *Point) Print() string {
	var output string
	// Move cursor to X Y
	output += "\033[" + strconv.Itoa(p.Y) + ";" + strconv.Itoa(p.X) + "H"
	if p.EscapeCode != "" {
		//output += "\033[0m"
		output += p.EscapeCode
		output += p.Character
		output += "\033[0m"
	} else {
		output += p.Character
	}

	return output
}

// NewPointMap creates a new PointMap.
func NewPointMap(width, height int) PointMap {
	pm := make([][]Point, width+1)
	for i := range pm {
		pm[i] = make([]Point, height+1)
		for j := range pm[i] {
			pm[i][j] = Point{
				X:          i,
				Y:          j,
				EscapeCode: "",
				Character:  "",
			}
		}
	}
	return pm
}
