package ui

import "strconv"

// SaveCursor saves the cursor position using the escape sequence
func (c *Console) SaveCursor() string {
	return "\033[s"
}

// RestoreCursor restores the cursor position using the escape sequence
func (c *Console) RestoreCursor() string {
	return "\033[u"
}

// HideCursor sets the cursor position using the escape sequence
func (c *Console) HideCursor() string {
	return "\033[?25l"
}

// ShowCursor sets the cursor position using the escape sequence
func (c *Console) ShowCursor() string {
	return "\033[?25h"
}

// MoveCursorToTopLeft Moves the cursor to the top left corner of the console
func (c *Console) MoveCursorToTopLeft() string {
	return "\033[1;1H"
}

// MoveCursorToBottomLeft Moves the cursor to the bottom left corner of the console
func (c *Console) MoveCursorToBottomLeft() string {
	return "\033[" + strconv.Itoa(c.Height) + ";0H"
}
