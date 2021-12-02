package ui

// ResetTerminal resets the terminal using the escape sequence
func (c *Console) ResetTerminal() string {
	return "\033c"
}
