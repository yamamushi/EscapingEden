package console

import "strconv"

type Point struct {
	X, Y       int
	EscapeCode string
	Character  string
}

type PointMap [][]Point

func (p *Point) Print() string {
	var output string
	// Move cursor to X Y
	output += "\033[" + strconv.Itoa(p.Y) + ";" + strconv.Itoa(p.X) + "H"
	output += "\033[0m"
	output += p.EscapeCode
	output += p.Character
	output += "\033[0m"
	return output
}

func PointMapToString(pm PointMap) string {
	var s string
	for _, row := range pm {
		for _, p := range row {
			s += p.Print()
		}
	}
	return s
}

func NewPointMap(width, height int) PointMap {
	pm := make([][]Point, width+1)
	for i := range pm {
		pm[i] = make([]Point, height+1)
		for j := range pm[i] {
			pm[i][j] = Point{
				X:          i,
				Y:          j,
				EscapeCode: "",
				Character:  " ",
			}
		}
	}
	return pm
}
