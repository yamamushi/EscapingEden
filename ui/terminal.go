package ui

import "strconv"

// ScrollLock locks the scroll
func (c *Console) ScrollLock() string {
	return "\033[?1049h"
}

// ScrollUnlock unlocks the scroll
func (c *Console) ScrollUnlock() string {
	return "\033[?1049l"
}

// ClearTerminal clears the terminal using the escape sequence
func (c *Console) ClearTerminal() string {
	return "\033[2J\n"
}

// ClearNotPrompt will clear each line of the console except the prompt
func (c *Console) ClearNotPrompt() string {
	var s string
	// save cursor position
	//s = c.SaveCursor()
	for i := 0; i < c.Height; i++ {
		// Move cursor to line i
		s = s + "\033[" + strconv.Itoa(i+1) + ";0H"
		// Clear line
		s = s + "\033[2K"
	}
	return s
}

// HardClear Terminal clears each line individually for height of console
func (c *Console) HardClear() string {
	var s string
	for i := 0; i < c.Height; i++ {
		// Move cursor to line i
		s = s + "\033[" + strconv.Itoa(i+1) + ";0H"
		// Clear line
		s = s + "\033[2K"
	}
	return s
}

// ResetTerminal resets the terminal using the escape sequence
func (c *Console) ResetTerminal() string {
	return "\033c"
}
