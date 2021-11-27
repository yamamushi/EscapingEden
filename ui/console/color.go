package console

import (
	"strconv"
)

// These are functions used for appending color codes to strings
// to make them print in different colors.

type ColorCode struct {
	R, G, B int
}

func (c *ColorCode) FG() string {
	return "\033[38;2;" + strconv.Itoa(c.R) + ";" + strconv.Itoa(c.G) + ";" + strconv.Itoa(c.B) + "m"
}

func (c *ColorCode) BG() string {
	return "\033[48;2;" + strconv.Itoa(c.R) + ";" + strconv.Itoa(c.G) + ";" + strconv.Itoa(c.B) + "m"
}

// SetRGB sets the xterm256 color to the given RGB values on the provided string, and appends a white on black code to the end
// of the provided string.
func SetRGB(cc *ColorCode, s string) []*Point {
	var output []*Point
	for _, character := range s {
		output = append(output, &Point{0, 0, "\033[38;2;" + strconv.Itoa(int(cc.R)) + ";" + strconv.Itoa(int(cc.G)) + ";" + strconv.Itoa(int(cc.B)) + "m", string(character)})
	}
	return output
}

func RGBCode(r, g, b int) *ColorCode {
	return &ColorCode{r, g, b}
}

// ResetStyle outputs a black on white code string
func ResetStyle() string {
	return "\033[0m"
}

func BoldText(s string) []*Point {
	var output []*Point
	for _, character := range s {
		output = append(output, &Point{0, 0, "\033[1m", string(character)})
	}
	return output
}
