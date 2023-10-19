package popupbox

import (
	"github.com/yamamushi/EscapingEden/ui/util"
	"github.com/yamamushi/EscapingEden/edenutil"
)

// DrawBorder returns the border of a window using code page 437 characters as a string
func (pb *PopupBox) DrawBorder(winX int, winY int) {
	// Draw top border using code page 437 characters starting at position winX, winY

	// Move cursor to top left corner of window
	// Draw top left corner
	if pb.Active {
		pb.PrintChar(winX, winY, UCTopLeftBorder, SHGreen)

	} else {
		pb.PrintChar(winX, winY, UCTopLeftBorder, "")
	}

	// Draw left border
	for i := 1; i < pb.Height+1; i++ {
		// Inserts a vertical line
		if pb.Active {
			pb.PrintChar(winX, winY+i, UCVerticalBorder, SHGreen)
		} else {
			pb.PrintChar(winX, winY+i, UCVerticalBorder, "")
		}
	}
	// Draw bottom left corner
	if pb.Active {
		pb.PrintChar(winX, winY+pb.Height+1, UCBottomLeftBorder, SHGreen)
	} else {
		pb.PrintChar(winX, winY+pb.Height+1, UCBottomLeftBorder, "")
	}

	// Draw top border
	for i := 1; i < pb.Width; i++ {
		// Inserts a horizontal line
		if pb.Active {
			pb.PrintChar(winX+i, winY, UCHorizontalBorder, SHGreen)
		} else {
			pb.PrintChar(winX+i, winY, UCHorizontalBorder, "")
		}
	}

	// Draw top right corner
	if pb.Active {
		pb.PrintChar(winX+pb.Width, winY, UCTopRightBorder, SHGreen)
	} else {
		pb.PrintChar(winX+pb.Width, winY, UCTopRightBorder, "")
	}

	// Draw right border
	for i := 1; i < pb.Height+1; i++ {
		// Inserts a vertical line
		if pb.Active {
			pb.PrintChar(winX+pb.Width, winY+i, UCVerticalBorder, SHGreen)
		} else {
			pb.PrintChar(winX+pb.Width, winY+i, UCVerticalBorder, "")
		}
	}

	// Draw bottom right corner
	if pb.Active {
		pb.PrintChar(winX+pb.Width, winY+pb.Height+1, UCBottomRightBorder, SHGreen)
	} else {
		pb.PrintChar(winX+pb.Width, winY+pb.Height+1, UCBottomRightBorder, "")
	}

	// Draw bottom border
	for i := 1; i < pb.Width; i++ {
		// Inserts a horizontal line
		if pb.Active {
			pb.PrintChar(winX+i, winY+pb.Height+1, UCHorizontalBorder, SHGreen)
		} else {
			pb.PrintChar(winX+i, winY+pb.Height+1, UCHorizontalBorder, "")
		}
	}

	// Prints the world "Close" at the bottom center
	if pb.Active {
		var colorCode string
		// set colorCode background to green
		colorCode = util.RGBCode(0, 0, 0).FG()
		colorCode += util.RGBCode(0, 255, 0).BG()

		pb.PrintLn(winX+pb.Width/2-3, winY+pb.Height+1, "Close", colorCode)
	} else {
		//log.Println("Writing Close to window")
		pb.PrintLn(winX+pb.Width/2-3, winY+pb.Height+1, "Close", "")
	}
}
