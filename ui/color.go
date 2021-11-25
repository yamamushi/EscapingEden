package ui

import "strconv"

// These are functions used for appending color codes to strings
// to make them print in different colors.

type ColorCode struct {
	R, G, B int
}

// SetRGB sets the xterm256 color to the given RGB values on the provided string, and appends a white on black code to the end
// of the provided string.
func SetRGB(cc *ColorCode, s string) string {
	return "\033[38;2;" + strconv.Itoa(int(cc.R)) + ";" + strconv.Itoa(int(cc.G)) + ";" + strconv.Itoa(int(cc.B)) + "m" + s + ResetColor()
}

// ResetColor outputs a black on white code string
func ResetColor() string {
	return "\033[0m"
}
