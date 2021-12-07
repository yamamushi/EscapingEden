package types

import (
	"github.com/yamamushi/EscapingEden/terminals"
)

// Point is a 2D point.
type Point struct {
	X, Y       int
	EscapeCode string
	Character  string
	NoReset    bool
}

// PointMap is a map of points.
type PointMap [][]Point

// Print prints a point as an escaped string.
func (p *Point) Print(term terminals.TerminalType) string {
	var output string
	// Move cursor to X Y
	output += term.MoveCursor(p.X, p.Y)
	if p.EscapeCode != "" {
		output += p.EscapeCode
		output += p.Character
		if !p.NoReset {
			output += term.Reset()
		}
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
