package help

import (
	"github.com/yamamushi/EscapingEden/edenutil"
)

// DrawBorder returns the border of a window using code page 437 characters as a string
func (hw *HelpWindow) DrawBorder(winX int, winY int) {
	// Draw top border using code page 437 characters starting at position winX, winY

	// Move cursor to top left corner of window
	// Draw top left corner
	if hw.Active {
		hw.PrintChar(winX, winY, UCTopLeftBorder, SHGreen)

	} else {
		hw.PrintChar(winX, winY, UCTopLeftBorder, "")
	}

	// Draw left border
	for i := 1; i < hw.Height+1; i++ {
		// Inserts a vertical line
		if hw.Active {
			hw.PrintChar(winX, winY+i, UCVerticalBorder, SHGreen)
		} else {
			hw.PrintChar(winX, winY+i, UCVerticalBorder, "")
		}
	}

	// Draw bottom left corner
	if hw.Active {
		hw.PrintChar(winX, winY+hw.Height+1, UCBottomLeftBorder, SHGreen)
	} else {
		hw.PrintChar(winX, winY+hw.Height+1, UCBottomLeftBorder, "")
	}

	// Draw top border
	for i := 1; i < hw.Width; i++ {
		// Inserts a horizontal line
		if hw.Active {
			hw.PrintCharColor(winX+i, winY, UCHorizontalBorder, SHGreen)
		} else {
			hw.PrintChar(winX+i, winY, UCHorizontalBorder, "")
		}
	}

	// Draw top right corner
	if hw.Active {
		hw.PrintChar(winX+hw.Width, winY, UCTopRightBorder, SHGreen)
	} else {
		hw.PrintChar(winX+hw.Width, winY, UCTopRightBorder, "")
	}

	// Draw right border
	for i := 1; i < hw.Height+1; i++ {
		// Inserts a vertical line
		if hw.Active {
			hw.PrintChar(winX+hw.Width, winY+i, UCVerticalBorder, SHGreen)
		} else {
			hw.PrintChar(winX+hw.Width, winY+i, UCVerticalBorder, "")
		}
	}

	// Draw bottom right corner
	if hw.Active {
		hw.PrintChar(winX+hw.Width, winY+hw.Height+1, UCBottomRightBorder, SHGreen)
	} else {
		hw.PrintChar(winX+hw.Width, winY+hw.Height+1, UCBottomRightBorder, "")
	}

	// Draw bottom border
	for i := 1; i < hw.Width; i++ {
		// Inserts a horizontal line
		if hw.Active {
			hw.PrintCharColor(winX+i, winY+hw.Height+1, UCHorizontalBorder, SHGreen)
		} else {
			hw.PrintChar(winX+i, winY+hw.Height+1, UCHorizontalBorder, "")
		}
	}

	// Print the borders for the top and bottom fields
	if hw.Active {
		// Left Side
		hw.PrintChar(winX, winY+2, UCVerticalGameBorder, SHGreen)
		hw.PrintChar(winX, winY+hw.Height-hw.ScrollBufferLimit+1, UCVerticalGameBorder, SHGreen)

		// Right Side
		hw.PrintChar(winX+hw.Width, winY+2, "\u2524", SHGreen)
		hw.PrintChar(winX+hw.Width, winY+hw.Height-hw.ScrollBufferLimit+1, "\u2524", SHGreen)

		// Top Border
		for i := 1; i < hw.Width; i++ {
			// Inserts a horizontal line
			hw.PrintCharColor(winX+i, winY+2, UCHorizontalBorder, SHGreen)
		}

		// Draw bottom border
		for i := 1; i < hw.Width; i++ {
			// Inserts a horizontal line
			hw.PrintCharColor(winX+i, winY+hw.Height-hw.ScrollBufferLimit+1, UCHorizontalBorder, SHGreen)
		}

	} else {
		// Left Side
		hw.PrintChar(winX, winY+2, UCVerticalGameBorder, "")
		hw.PrintChar(winX, winY+hw.Height-hw.ScrollBufferLimit+1, UCVerticalGameBorder, "")

		// Right Side
		hw.PrintChar(winX+hw.Width, winY+2, "\u2524", "")
		hw.PrintChar(winX+hw.Width, winY+hw.Height-hw.ScrollBufferLimit+1, "\u2524", "")

		// Top Border
		for i := 1; i < hw.Width; i++ {
			// Inserts a horizontal line
			hw.PrintCharColor(winX+i, winY+2, UCHorizontalBorder, "")
		}

		// Draw bottom border
		for i := 1; i < hw.Width; i++ {
			// Inserts a horizontal line
			hw.PrintCharColor(winX+i, winY+hw.Height-hw.ScrollBufferLimit+1, UCHorizontalBorder, "")
		}
	}
}
