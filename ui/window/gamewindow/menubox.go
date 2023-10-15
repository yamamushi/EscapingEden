package gamewindow

type MenuBox struct {
	// The menu box's position
	X, Y int
	// The menu box's width and height
	Width, Height int
	// The menu box's title
	Title string
	// The menu box's options
	Options []struct {
		Data interface{}
	}
}

func (mb *MenuBox) Draw(gw *GameWindow) {
	mb.Clear(gw)
	mb.DrawBorder(gw)
	mb.DrawTitle(gw)
}

func (mb *MenuBox) Clear(gw *GameWindow) {
	for i := 0; i < mb.Width; i++ {
		for j := 0; j < mb.Height; j++ {
			gw.DrawToVisibleMap(mb.X+i, mb.Y+j, " ", "")
		}
	}
}

func (mb *MenuBox) DrawTitle(gw *GameWindow) {
	// Draw the title centered on the top border
	// First calculate the x position of the title
	titleX := mb.X + (mb.Width/2 - len(mb.Title)/2)
	// Then draw the title
	gw.PrintStringToMap(titleX, mb.Y, mb.Title, "")
}

func (mb *MenuBox) DrawBorder(gw *GameWindow) {
	// First draw the box
	// Draw the top and bottom of the box
	for i := 0; i < mb.Width; i++ {
		// Draw the top of the box
		if gw.Active {
			gw.DrawToVisibleMap(mb.X+i, mb.Y, "\u2500", "\033[32m")
			// Draw the bottom of the box
			gw.DrawToVisibleMap(mb.X+i, mb.Y+mb.Height-1, "\u2500", "\033[32m")

		} else {
			gw.DrawToVisibleMap(mb.X+i, mb.Y, "\u2500", "")
			// Draw the bottom of the box
			gw.DrawToVisibleMap(mb.X+i, mb.Y+mb.Height-1, "\u2500", "")
		}

	}
	// Draw the left and right of the box
	for i := 0; i < mb.Height; i++ {
		// Draw the left of the box
		if gw.Active {
			gw.DrawToVisibleMap(mb.X, mb.Y+i, "\u2502", "\033[32m")
			// Draw the right of the box
			gw.DrawToVisibleMap(mb.X+mb.Width-1, mb.Y+i, "\u2502", "\033[32m")
		} else {
			gw.DrawToVisibleMap(mb.X, mb.Y+i, "\u2502", "")
			// Draw the right of the box
			gw.DrawToVisibleMap(mb.X+mb.Width-1, mb.Y+i, "\u2502", "")
		}
	}
	// Draw the corners of the box
	if gw.Active {
		gw.DrawToVisibleMap(mb.X, mb.Y, "\u250C", "\033[32m")
		gw.DrawToVisibleMap(mb.X+mb.Width-1, mb.Y, "\u2510", "\033[32m")
		gw.DrawToVisibleMap(mb.X, mb.Y+mb.Height-1, "\u2514", "\033[32m")
		gw.DrawToVisibleMap(mb.X+mb.Width-1, mb.Y+mb.Height-1, "\u2518", "\033[32m")
	} else {
		gw.DrawToVisibleMap(mb.X, mb.Y, "\u250C", "")
		gw.DrawToVisibleMap(mb.X+mb.Width-1, mb.Y, "\u2510", "")
		gw.DrawToVisibleMap(mb.X, mb.Y+mb.Height-1, "\u2514", "")
		gw.DrawToVisibleMap(mb.X+mb.Width-1, mb.Y+mb.Height-1, "\u2518", "")
	}
}
