package util

import (
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

// RGBCode returns a ColorCode with the given RGB values.
func RGBCode(r, g, b int) *ColorCode {
	return &ColorCode{r, g, b}
}
