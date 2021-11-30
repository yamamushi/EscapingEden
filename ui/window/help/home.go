package help

// DrawHome is a wrapper around the functions necessary to draw the home screen of the Help Window.
func (hw *HelpWindow) DrawHome() {
	//hw.SetContents(strconv.Itoa(hw.indexPage))
	hw.DrawHomeInfo()
	hw.PrintControls()
}

// DrawHomeInfo draws the content for the top of field of the Help Window on the home screen.
func (hw *HelpWindow) DrawHomeInfo() {
	// Top Field
	windowTitle := "Escaping Eden Help"
	pageInfo := "Home"
	hw.PrintLn(hw.X+1, hw.Y+1, windowTitle, "\033[1m")
	hw.PrintLn(hw.X+hw.Width-len(pageInfo)-1, hw.Y+1, pageInfo, "")
}
