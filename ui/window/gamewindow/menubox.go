package gamewindow

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/edenutil"
)

type MenuBox struct {
	// The menu box's position
	X, Y int
	// The menu box's width and height
	Width, Height int
	// The menu box's title
	Title string
	// The menu box's options
	Options          []MenuBoxOption
	ResponseCallback interface{}
	PopupMenu        *MenuBox
	CallbackData     interface{} // Arbitrary data we can unpack later
}

type MenuBoxOption struct {
	Name       string
	SkipDraw   bool
	Data       interface{}
	Keybind    string
	ControlKey int
	Callback   interface{}
	Order      int
}

func (mb *MenuBox) HandleInput(gw *GameWindow, input string) {
	gw.Log.Println(logging.LogInfo, "Menubox received input: ", input)
	// Handle input for the menu box
	// First check if the input is a keybind
	if mb.ResponseCallback != nil {
		gw.Log.Println(logging.LogInfo, "Menubox called response callback ")

		switch mb.ResponseCallback.(type) {
		case func(*MenuBox, string):
			mb.ResponseCallback.(func(*MenuBox, string))(mb, input)
		default:
			gw.Log.Println(logging.LogInfo, "Menubox called response callback and failed")
		}
		return
	}
	for _, option := range mb.Options {
		if input == option.Keybind {
			gw.Log.Println(logging.LogInfo, "Menubox received input for ", option.Name)
			// If it is, call the callback
			switch option.Callback.(type) {
			case func(box *MenuBox):
				option.Callback.(func(*MenuBox))(mb)
			case func():
				option.Callback.(func())()
			case func(string):
				option.Callback.(func(*MenuBox, string))(mb, input)
			}
		} else if int(input[0]) == option.ControlKey {
			// If it is, call the callback
			switch option.Callback.(type) {
			case func(*MenuBox):
				option.Callback.(func(*MenuBox))(mb)
			case func(*MenuBox, string):
				option.Callback.(func(*MenuBox, string))(mb, input)
			}
		}
	}
}

func (mb *MenuBox) Draw(gw *GameWindow) {
	mb.Clear(gw)
	mb.DrawBorder(gw)
	mb.DrawTitle(gw)
	mb.DrawMenuItems(gw)
	mb.DrawPopupMenu(gw)
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
			gw.DrawToVisibleMap(mb.X+i, mb.Y, UCHorizontalBorder, SHGreen)
			// Draw the bottom of the box
			gw.DrawToVisibleMap(mb.X+i, mb.Y+mb.Height-1, UCHorizontalBorder, SHGreen)

		} else {
			gw.DrawToVisibleMap(mb.X+i, mb.Y, UCHorizontalBorder, "")
			// Draw the bottom of the box
			gw.DrawToVisibleMap(mb.X+i, mb.Y+mb.Height-1, UCHorizontalBorder, "")
		}

	}
	// Draw the left and right of the box
	for i := 0; i < mb.Height; i++ {
		// Draw the left of the box
		if gw.Active {
			gw.DrawToVisibleMap(mb.X, mb.Y+i, UCVerticalBorder, SHGreen)
			// Draw the right of the box
			gw.DrawToVisibleMap(mb.X+mb.Width-1, mb.Y+i, UCVerticalBorder, SHGreen)
		} else {
			gw.DrawToVisibleMap(mb.X, mb.Y+i, UCVerticalBorder, "")
			// Draw the right of the box
			gw.DrawToVisibleMap(mb.X+mb.Width-1, mb.Y+i, UCVerticalBorder, "")
		}
	}
	// Draw the corners of the box
	if gw.Active {
		gw.DrawToVisibleMap(mb.X, mb.Y, UCTopLeftBorder, SHGreen)
		gw.DrawToVisibleMap(mb.X+mb.Width-1, mb.Y, UCTopRightBorder, SHGreen)
		gw.DrawToVisibleMap(mb.X, mb.Y+mb.Height-1, UCBottomLeftBorder, SHGreen)
		gw.DrawToVisibleMap(mb.X+mb.Width-1, mb.Y+mb.Height-1, UCBottomRightBorder, SHGreen)
	} else {
		gw.DrawToVisibleMap(mb.X, mb.Y, UCTopLeftBorder, "")
		gw.DrawToVisibleMap(mb.X+mb.Width-1, mb.Y, UCTopRightBorder, "")
		gw.DrawToVisibleMap(mb.X, mb.Y+mb.Height-1, UCBottomLeftBorder, "")
		gw.DrawToVisibleMap(mb.X+mb.Width-1, mb.Y+mb.Height-1, UCBottomRightBorder, "")
	}
}

func (mb *MenuBox) DrawMenuItems(gw *GameWindow) {
	// Range through the menu items and draw them
	for i, option := range mb.Options {
		if option.SkipDraw {
			continue
		}
		// Draw the keybind
		gw.PrintStringToMap(mb.X+2, mb.Y+i+2, option.Keybind+")", "")
		// Draw the name
		gw.PrintStringToMap(mb.X+5, mb.Y+i+2, option.Name, "")
	}
}

func (mb *MenuBox) DrawPopupMenu(gw *GameWindow) {
	// Draw the popup menu
	if mb.PopupMenu != nil {
		mb.PopupMenu.Draw(gw)
	}
}
