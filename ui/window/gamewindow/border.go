package gamewindow

import (
	"github.com/yamamushi/EscapingEden/edenutil"
)

// We implement our own border drawing to account for status bars and other things

// DrawBorder returns the border of a window using code page 437 characters as a string
func (w *GameWindow) DrawBorder(winX int, winY int) {
	// Draw top border using code page 437 characters starting at position winX, winY
	statusBarHeight := 3

	// Move cursor to top left corner of window
	// Draw top left corner
	if w.Active {
		w.PrintChar(winX, winY, UCTopLeftBorder, SHGreen)

	} else {
		w.PrintChar(winX, winY, UCTopLeftBorder, w.Terminal.Bold())
	}

	// Draw left border
	for i := 1; i < w.Height+1; i++ {
		// Inserts a vertical line
		if i == w.Height-statusBarHeight {
			if w.Active {
				w.PrintChar(winX, winY+i, UCVerticalGameBorder, SHGreen)
			} else {
				w.PrintChar(winX, winY+i, UCVerticalGameBorder, w.Terminal.Bold())
			}
		} else {
			if w.Active {
				w.PrintChar(winX, winY+i, UCVerticalBorder, SHGreen)
			} else {
				w.PrintChar(winX, winY+i, UCVerticalBorder, w.Terminal.Bold())
			}
		}

	}
	// Draw bottom left corner
	if w.Active {
		w.PrintChar(winX, winY+w.Height+1, UCBottomLeftBorder, SHGreen)
	} else {
		w.PrintChar(winX, winY+w.Height+1, UCBottomLeftBorder, w.Terminal.Bold())
	}

	// Draw top border
	for i := 1; i < w.Width; i++ {
		// Inserts a horizontal line
		if w.Active {
			w.PrintCharColor(winX+i, winY, UCHorizontalBorder, SHGreen)
		} else {
			w.PrintCharColor(winX+i, winY, UCHorizontalBorder, w.Terminal.Bold())
		}
	}

	// Draw top right corner
	if w.Active {
		w.PrintChar(winX+w.Width, winY, UCTopRightBorder, SHGreen)
	} else {
		w.PrintChar(winX+w.Width, winY, UCTopRightBorder, w.Terminal.Bold())
	}

	// Draw right border
	for i := 1; i < w.Height+1; i++ {
		// Inserts a vertical line
		if i == w.Height-statusBarHeight {
			if w.Active {
				w.PrintChar(winX+w.Width, winY+i, "\u2524", SHGreen)
			} else {
				w.PrintChar(winX+w.Width, winY+i, "\u2524", w.Terminal.Bold())
			}
		} else {
			if w.Active {
				w.PrintChar(winX+w.Width, winY+i, UCVerticalBorder, SHGreen)
			} else {
				w.PrintChar(winX+w.Width, winY+i, UCVerticalBorder, w.Terminal.Bold())
			}
		}
	}

	// Draw bottom right corner
	if w.Active {
		w.PrintChar(winX+w.Width, winY+w.Height+1, UCBottomRightBorder, SHGreen)
	} else {
		w.PrintChar(winX+w.Width, winY+w.Height+1, UCBottomRightBorder, w.Terminal.Bold())
	}

	// Draw bottom border
	for i := 1; i < w.Width; i++ {
		// Inserts a horizontal line
		if w.Active {
			w.PrintCharColor(winX+i, winY+w.Height+1, UCHorizontalBorder, SHGreen)
		} else {
			w.PrintCharColor(winX+i, winY+w.Height+1, UCHorizontalBorder, w.Terminal.Bold())
		}
	}

	// Draw bottom statusbar border
	for i := 1; i < w.Width; i++ {
		// Inserts a horizontal line
		if w.Active {
			w.PrintCharColor(winX+i, winY+w.Height-statusBarHeight, UCHorizontalBorder, SHGreen)
		} else {
			w.PrintCharColor(winX+i, winY+w.Height-statusBarHeight, UCHorizontalBorder, w.Terminal.Bold())
		}
	}
}
