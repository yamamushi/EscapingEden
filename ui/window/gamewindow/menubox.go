package gamewindow

import (
	"github.com/yamamushi/EscapingEden/edenutil"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/ui/types"
)

type MenuBoxType interface {
	HandleInput(*GameWindow, types.InputType, string)
	Draw(*GameWindow)
	Clear(*GameWindow)
	DrawTitle(*GameWindow)
	DrawBorder(*GameWindow)
	DrawMenuItems(*GameWindow)
	DrawPopupMenu(*GameWindow)
	PrintToMenu(*GameWindow, int, int, string, string)
	CloseMenus(*GameWindow)
}

type MenuBox struct {
	// The menu box's position
	X, Y int
	// The menu box's width and height
	Width, Height int
	// The menu box's title
	Title string
	// The menu box's options
	Options                  []MenuBoxOption
	ResponseCallback         interface{}
	PopupMenu                *MenuBox
	CallbackData             interface{} // Arbitrary data we can unpack later
	CallbackStatusBarMessage string
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

func (mb *MenuBox) CloseMenus(gw *GameWindow) {
	gw.CloseMenus = true
}

func (mb *MenuBox) HandleInput(gw *GameWindow, inputType types.InputType, input string) {
	gw.Log.Println(logging.LogInfo, "Menubox received input: ", input)
	// Handle input for the menu box
	// First check if the input is a keybind
	if inputType == types.InputEscape {
		mb.CloseMenus(gw)
		return
	}
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
	mb.DrawMenuItems(gw)
	mb.DrawBorder(gw)
	mb.DrawTitle(gw)
	mb.DrawPopupMenu(gw)
	gw.SetStatusBarMessage(mb.CallbackStatusBarMessage)
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
			gw.DrawToVisibleMap(mb.X+i, mb.Y, edenutil.UCHorizontalBorder, "\033[32m")
			// Draw the bottom of the box
			gw.DrawToVisibleMap(mb.X+i, mb.Y+mb.Height-1, edenutil.UCHorizontalBorder, "\033[32m")

		} else {
			gw.DrawToVisibleMap(mb.X+i, mb.Y, edenutil.UCHorizontalBorder, "")
			// Draw the bottom of the box
			gw.DrawToVisibleMap(mb.X+i, mb.Y+mb.Height-1, edenutil.UCHorizontalBorder, "")
		}

	}
	// Draw the left and right of the box
	for i := 0; i < mb.Height; i++ {
		// Draw the left of the box
		if gw.Active {
			gw.DrawToVisibleMap(mb.X, mb.Y+i, edenutil.UCVerticalBorder, "\033[32m")
			// Draw the right of the box
			gw.DrawToVisibleMap(mb.X+mb.Width-1, mb.Y+i, edenutil.UCVerticalBorder, "\033[32m")
		} else {
			gw.DrawToVisibleMap(mb.X, mb.Y+i, edenutil.UCVerticalBorder, "")
			// Draw the right of the box
			gw.DrawToVisibleMap(mb.X+mb.Width-1, mb.Y+i, edenutil.UCVerticalBorder, "")
		}
	}
	// Draw the corners of the box
	if gw.Active {
		gw.DrawToVisibleMap(mb.X, mb.Y, edenutil.UCTopLeftBorder, "\033[32m")
		gw.DrawToVisibleMap(mb.X+mb.Width-1, mb.Y, edenutil.UCTopRightBorder, "\033[32m")
		gw.DrawToVisibleMap(mb.X, mb.Y+mb.Height-1, edenutil.UCBottomLeftBorder, "\033[32m")
		gw.DrawToVisibleMap(mb.X+mb.Width-1, mb.Y+mb.Height-1, edenutil.UCBottomRightBorder, "\033[32m")
	} else {
		gw.DrawToVisibleMap(mb.X, mb.Y, edenutil.UCTopLeftBorder, "")
		gw.DrawToVisibleMap(mb.X+mb.Width-1, mb.Y, edenutil.UCTopRightBorder, "")
		gw.DrawToVisibleMap(mb.X, mb.Y+mb.Height-1, edenutil.UCBottomLeftBorder, "")
		gw.DrawToVisibleMap(mb.X+mb.Width-1, mb.Y+mb.Height-1, edenutil.UCBottomRightBorder, "")
	}
}

func (mb *MenuBox) DrawMenuItems(gw *GameWindow) {
	// Range through the menu items and draw them
	for i, option := range mb.Options {
		if option.SkipDraw {
			continue
		}
		// Draw the keybind
		mb.PrintToMenu(gw, 2, i+2, option.Keybind+")", "")
		// Draw the name
		mb.PrintToMenu(gw, 5, i+2, option.Name, "")
	}
}

func (mb *MenuBox) DrawPopupMenu(gw *GameWindow) {
	// Draw the popup menu
	if mb.PopupMenu != nil {
		mb.PopupMenu.Draw(gw)
	}
}

func (mb *MenuBox) PrintToMenu(gw *GameWindow, x, y int, input string, escapeCode string) {
	if x > mb.Width-1 || x < 0 {
		return
	}
	if y > mb.Height-1 || y < 0 {
		return
	}
	for i, character := range input {
		gw.DrawToVisibleMap(mb.X+x+i, mb.Y+y, string(character), escapeCode)
	}
}
