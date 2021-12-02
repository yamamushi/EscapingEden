package util

import (
	"github.com/yamamushi/EscapingEden/ui/types"
	"strconv"
)

// These are functions used for appending color codes to strings
// to make them print in different colors.

// ColorCode stores RGB values for a color.
type ColorCode struct {
	R, G, B int
}

// FG returns a string that sets the foreground color to the given color.
func (c *ColorCode) FG() string {
	return "\033[38;2;" + strconv.Itoa(c.R) + ";" + strconv.Itoa(c.G) + ";" + strconv.Itoa(c.B) + "m"
}

// BG returns a string that sets the background color to the given color.
func (c *ColorCode) BG() string {
	return "\033[48;2;" + strconv.Itoa(c.R) + ";" + strconv.Itoa(c.G) + ";" + strconv.Itoa(c.B) + "m"
}

// SetRGB sets the xterm256 color to the given RGB values on the provided string, and appends a white on black code to the end
// of the provided string.
func SetRGB(cc *ColorCode, s string) []*types.Point {
	var output []*types.Point
	for _, character := range s {
		output = append(output, &types.Point{0, 0, "\033[38;2;" + strconv.Itoa(int(cc.R)) + ";" + strconv.Itoa(int(cc.G)) + ";" + strconv.Itoa(int(cc.B)) + "m", string(character)})
	}
	return output
}

// RGBCode returns a ColorCode with the given RGB values.
func RGBCode(r, g, b int) *ColorCode {
	return &ColorCode{r, g, b}
}
