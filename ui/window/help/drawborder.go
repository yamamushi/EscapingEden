package help

// DrawBorder returns the border of a window using code page 437 characters as a string
func (hw *HelpWindow) DrawBorder(winX int, winY int) {
	// Draw top border using code page 437 characters starting at position winX, winY

	// Move cursor to top left corner of window
	// Draw top left corner
	if hw.Active {
		hw.PrintChar(winX, winY, "\u250c", "\033[32m")

	} else {
		hw.PrintChar(winX, winY, "\u250c", "")
	}

	// Draw left border
	for i := 1; i < hw.Height+1; i++ {
		// Inserts a vertical line
		if hw.Active {
			hw.PrintChar(winX, winY+i, "\u2502", "\033[32m")
		} else {
			hw.PrintChar(winX, winY+i, "\u2502", "")
		}
	}

	// Draw bottom left corner
	if hw.Active {
		hw.PrintChar(winX, winY+hw.Height+1, "\u2514", "\033[32m")
	} else {
		hw.PrintChar(winX, winY+hw.Height+1, "\u2514", "")
	}

	// Draw top border
	for i := 1; i < hw.Width; i++ {
		// Inserts a horizontal line
		if hw.Active {
			hw.PrintChar(winX+i, winY, "\u2500", "\033[32m")
		} else {
			hw.PrintChar(winX+i, winY, "\u2500", "")
		}
	}

	// Draw top right corner
	if hw.Active {
		hw.PrintChar(winX+hw.Width, winY, "\u2510", "\033[32m")
	} else {
		hw.PrintChar(winX+hw.Width, winY, "\u2510", "")
	}

	// Draw right border
	for i := 1; i < hw.Height+1; i++ {
		// Inserts a vertical line
		if hw.Active {
			hw.PrintChar(winX+hw.Width, winY+i, "\u2502", "\033[32m")
		} else {
			hw.PrintChar(winX+hw.Width, winY+i, "\u2502", "")
		}
	}

	// Draw bottom right corner
	if hw.Active {
		hw.PrintChar(winX+hw.Width, winY+hw.Height+1, "\u2518", "\033[32m")
	} else {
		hw.PrintChar(winX+hw.Width, winY+hw.Height+1, "\u2518", "")
	}

	// Draw bottom border
	for i := 1; i < hw.Width; i++ {
		// Inserts a horizontal line
		if hw.Active {
			hw.PrintChar(winX+i, winY+hw.Height+1, "\u2500", "\033[32m")
		} else {
			hw.PrintChar(winX+i, winY+hw.Height+1, "\u2500", "")
		}
	}

	// Print the borders for the top and bottom fields
	if hw.Active {
		// Left Side
		hw.PrintChar(winX, winY+2, "\u251C", "\033[32m")
		hw.PrintChar(winX, winY+hw.Height-hw.ScrollBufferLimit+1, "\u251C", "\033[32m")

		// Right Side
		hw.PrintChar(winX+hw.Width, winY+2, "\u2524", "\033[32m")
		hw.PrintChar(winX+hw.Width, winY+hw.Height-hw.ScrollBufferLimit+1, "\u2524", "\033[32m")

		// Top Border
		for i := 1; i < hw.Width; i++ {
			// Inserts a horizontal line
			hw.PrintChar(winX+i, winY+2, "\u2500", "\033[32m")
		}

		// Draw bottom border
		for i := 1; i < hw.Width; i++ {
			// Inserts a horizontal line
			hw.PrintChar(winX+i, winY+hw.Height-hw.ScrollBufferLimit+1, "\u2500", "\033[32m")
		}

	} else {
		// Left Side
		hw.PrintChar(winX, winY+2, "\u251C", "")
		hw.PrintChar(winX, winY+hw.Height-hw.ScrollBufferLimit+1, "\u251C", "")

		// Right Side
		hw.PrintChar(winX+hw.Width, winY+2, "\u2524", "")
		hw.PrintChar(winX+hw.Width, winY+hw.Height-hw.ScrollBufferLimit+1, "\u2524", "")

		// Top Border
		for i := 1; i < hw.Width; i++ {
			// Inserts a horizontal line
			hw.PrintChar(winX+i, winY+2, "\u2500", "")
		}

		// Draw bottom border
		for i := 1; i < hw.Width; i++ {
			// Inserts a horizontal line
			hw.PrintChar(winX+i, winY+hw.Height-hw.ScrollBufferLimit+1, "\u2500", "")
		}
	}
}
