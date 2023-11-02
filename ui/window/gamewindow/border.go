package gamewindow

// We implement our own border drawing to account for status bars and other things

// DrawBorder returns the border of a window using code page 437 characters as a string
func (w *GameWindow) DrawBorder(winX int, winY int) {
	// Draw top border using code page 437 characters starting at position winX, winY
	statusBarHeight := 3

	// Move cursor to top left corner of window
	// Draw top left corner
	if w.Active {
		w.PrintChar(winX, winY, "\u250c", "\033[32m")

	} else {
		w.PrintChar(winX, winY, "\u250c", w.Terminal.Bold()+w.Terminal.Reset())
	}

	// Draw left border
	for i := 1; i < w.Height+1; i++ {
		// Inserts a vertical line
		if i == w.Height-statusBarHeight {
			if w.Active {
				w.PrintChar(winX, winY+i, "\u251C", "\033[32m")
			} else {
				w.PrintChar(winX, winY+i, "\u251C", w.Terminal.Bold()+w.Terminal.Reset())
			}
		} else {
			if w.Active {
				w.PrintChar(winX, winY+i, "\u2502", "\033[32m")
			} else {
				w.PrintChar(winX, winY+i, "\u2502", w.Terminal.Bold()+w.Terminal.Reset())
			}
		}

	}
	// Draw bottom left corner
	if w.Active {
		w.PrintChar(winX, winY+w.Height+1, "\u2514", "\033[32m")
	} else {
		w.PrintChar(winX, winY+w.Height+1, "\u2514", w.Terminal.Bold()+w.Terminal.Reset())
	}

	// Draw top border
	for i := 1; i < w.Width; i++ {
		// Inserts a horizontal line
		if w.Active {
			w.PrintCharColor(winX+i, winY, "\u2500", "\033[32m")
		} else {
			w.PrintCharColor(winX+i, winY, "\u2500", w.Terminal.Bold()+w.Terminal.Reset())
		}
	}

	// Draw top right corner
	if w.Active {
		w.PrintChar(winX+w.Width, winY, "\u2510", "\033[32m")
	} else {
		w.PrintChar(winX+w.Width, winY, "\u2510", w.Terminal.Bold()+w.Terminal.Reset())
	}

	// Draw right border
	for i := 1; i < w.Height+1; i++ {
		// Inserts a vertical line
		if i == w.Height-statusBarHeight {
			if w.Active {
				w.PrintChar(winX+w.Width, winY+i, "\u2524", "\033[32m")
			} else {
				w.PrintChar(winX+w.Width, winY+i, "\u2524", w.Terminal.Bold()+w.Terminal.Reset())
			}
		} else {
			if w.Active {
				w.PrintChar(winX+w.Width, winY+i, "\u2502", "\033[32m")
			} else {
				w.PrintChar(winX+w.Width, winY+i, "\u2502", w.Terminal.Bold()+w.Terminal.Reset())
			}
		}
	}

	// Draw bottom right corner
	if w.Active {
		w.PrintChar(winX+w.Width, winY+w.Height+1, "\u2518", "\033[32m")
	} else {
		w.PrintChar(winX+w.Width, winY+w.Height+1, "\u2518", w.Terminal.Bold()+w.Terminal.Reset())
	}

	// Draw bottom border
	for i := 1; i < w.Width; i++ {
		// Inserts a horizontal line
		if w.Active {
			w.PrintCharColor(winX+i, winY+w.Height+1, "\u2500", "\033[32m")
		} else {
			w.PrintCharColor(winX+i, winY+w.Height+1, "\u2500", w.Terminal.Bold()+w.Terminal.Reset())
		}
	}

	// Draw bottom statusbar border
	for i := 1; i < w.Width; i++ {
		// Inserts a horizontal line
		if w.Active {
			w.PrintCharColor(winX+i, winY+w.Height-statusBarHeight, "\u2500", "\033[32m")
		} else {
			w.PrintCharColor(winX+i, winY+w.Height-statusBarHeight, "\u2500", w.Terminal.Bold()+w.Terminal.Reset())
		}
	}
}
