package ui

import (
	"github.com/yamamushi/EscapingEden/ui/config"
	"github.com/yamamushi/EscapingEden/ui/types"
	"log"
)

// HandleInput accepts a string terminated by a newline and processes it.
func (c *Console) HandleInput(rawInput byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if !c.IsConsoleValidSize() {
		// If the console isn't valid we don't want to accept any input
		// However if we receive the letter q, we will exit the program
		if rawInput == 'q' {
			c.Shutdown = true
		}
		return
	}
	if rawInput == 0 {
		// Ignore these null bytes
		return
	}

	//log.Println("Console received input: ", int(rawInput))

	if rawInput == 8 {
		options := &config.WindowConfig{X: c.Width/2 - 40, Y: c.Height/2 - 10, Width: 100, Height: 20, Page: 0}
		go c.ToggleHelp(options)
		return
	}

	if rawInput == 18 {
		// ctrl-r to force a screen refresh
		for _, w := range c.Windows {
			w.ResetWindowDrawings()
		}
		c.ClearPointMap()
		c.FlushLastSent()
		c.forceScreenRefresh = true
		return
	}

	// Captures things like the arrow keys.
	if rawInput == '\033' {
		c.escapeBuffer = "\\033" // Just used for logging the escape sequence buffer, nothing else
		c.escapeSequence = true  // Lets us know that the next few bytes are going to be related to the escape sequence
		return
	}

	// If we have an active escape sequence, we continue parsing it.
	if c.escapeSequence {
		// If we have a [ symbol, we know we are starting a new escape sequence, there is no ]
		if rawInput == '[' {
			c.escapeBuffer += "["
			return
		} else if c.escapeBuffer == "\\033[" {
			// If our escape buffer has an escape sequence, we know we are still parsing it.
			c.escapeBuffer += string(rawInput)
			switch rawInput {
			case 'A':
				c.InputToActiveWindow(types.Input{Type: types.InputUp})
			case 'B':
				c.InputToActiveWindow(types.Input{Type: types.InputDown})
			case 'C':
				c.InputToActiveWindow(types.Input{Type: types.InputRight})
			case 'D':
				c.InputToActiveWindow(types.Input{Type: types.InputLeft})
			default:
				log.Println("Unknown escape sequence: ", c.escapeBuffer)
			}
			c.escapeBuffer = ""
			c.escapeSequence = false
			return
		}
		c.escapeBuffer = ""
		c.escapeSequence = false
	}

	// If we have a backspace, we remove the last character from the input buffer.
	if rawInput == '\b' || rawInput == '\x7f' {
		c.InputToActiveWindow(types.Input{Type: types.InputBackspace})
		return
	}
	if rawInput == '\r' {
		c.InputToActiveWindow(types.Input{Type: types.InputReturn})
		return
	}
	if rawInput == '\t' {
		if !c.IsPopupOpen() && !c.IsHelpOpen() {
			c.SetNextActiveWindow()
			for _, w := range c.Windows {
				w.ResetWindowDrawings()
				w.FlushLastSent()
			}
			c.ResetWindowDrawings()
			c.FlushLastSent()
			c.forceScreenRefresh = true
		}
		return
	}
	if rawInput == '\n' {
		c.InputToActiveWindow(types.Input{Type: types.InputNewline})
		return
	}

	c.InputToActiveWindow(types.Input{Type: types.InputCharacter, Data: string(rawInput)})

}
