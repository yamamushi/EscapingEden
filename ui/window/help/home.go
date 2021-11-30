package help

// DrawIndex is a wrapper around the functions necessary to draw the current index page
func (hw *HelpWindow) DrawHome() {
	//hw.SetContents(strconv.Itoa(hw.indexPage))
	hw.DrawHomeInfo()
	hw.PrintControls()
}

func (hw *HelpWindow) DrawHomeInfo() {
	// Top Field
	windowTitle := "Escaping Eden Help"
	pageInfo := "Home"
	hw.PrintLn(hw.X+1, hw.Y+1, windowTitle, "\033[1m")
	hw.PrintLn(hw.X+hw.Width-len(pageInfo)-1, hw.Y+1, pageInfo, "")
}
