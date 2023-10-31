package window

import (
	"github.com/yamamushi/EscapingEden/edenutil"
)

// DrawBorder draws the Window's border

// DrawBorder returns the border of a window using code page 437 characters as a string
func (w *Window) DrawBorder(winX int, winY int) {
	// Draw top border using code page 437 characters starting at position winX, winY

	// Move cursor to top left corner of window
	// Draw top left corner
	if w.Active {
		w.PrintChar(winX, winY, edenutil.UCTopLeftBorder, "\033[32m")

	} else {
		w.PrintChar(winX, winY, edenutil.UCTopLeftBorder, w.Terminal.Bold())
	}

	// Draw left border
	for i := 1; i < w.Height+1; i++ {
		// Inserts a vertical line
		if w.Active {
			w.PrintChar(winX, winY+i, edenutil.UCVerticalBorder, "\033[32m")
		} else {
			w.PrintChar(winX, winY+i, edenutil.UCVerticalBorder, w.Terminal.Bold())
		}
	}
	// Draw bottom left corner
	if w.Active {
		w.PrintChar(winX, winY+w.Height+1, edenutil.UCBottomLeftBorder, "\033[32m")
	} else {
		w.PrintChar(winX, winY+w.Height+1, edenutil.UCBottomLeftBorder, w.Terminal.Bold())
	}

	// Draw top border
	for i := 1; i < w.Width; i++ {
		// Inserts a horizontal line
		if w.Active {
			w.PrintCharColor(winX+i, winY, edenutil.UCHorizontalBorder, "\033[32m")
		} else {
			w.PrintCharColor(winX+i, winY, edenutil.UCHorizontalBorder, w.Terminal.Bold())
		}
	}

	// Draw top right corner
	if w.Active {
		w.PrintChar(winX+w.Width, winY, edenutil.UCTopRightBorder, "\033[32m")
	} else {
		w.PrintChar(winX+w.Width, winY, edenutil.UCTopRightBorder, w.Terminal.Bold())
	}

	// Draw right border
	for i := 1; i < w.Height+1; i++ {
		// Inserts a vertical line
		if w.Active {
			w.PrintChar(winX+w.Width, winY+i, edenutil.UCVerticalBorder, "\033[32m")
		} else {
			w.PrintChar(winX+w.Width, winY+i, edenutil.UCVerticalBorder, w.Terminal.Bold())
		}
	}

	// Draw bottom right corner
	if w.Active {
		w.PrintChar(winX+w.Width, winY+w.Height+1, edenutil.UCBottomRightBorder, "\033[32m")
	} else {
		w.PrintChar(winX+w.Width, winY+w.Height+1, edenutil.UCBottomRightBorder, w.Terminal.Bold())
	}

	// Draw bottom border
	for i := 1; i < w.Width; i++ {
		// Inserts a horizontal line
		if w.Active {
			w.PrintCharColor(winX+i, winY+w.Height+1, edenutil.UCHorizontalBorder, "\033[32m")
		} else {
			w.PrintCharColor(winX+i, winY+w.Height+1, edenutil.UCHorizontalBorder, w.Terminal.Bold())
		}
	}
}
