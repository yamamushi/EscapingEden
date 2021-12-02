package help

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/ui/util"
)

// DrawHome is a wrapper around the functions necessary to draw the home screen of the Help Window.
func (hw *HelpWindow) DrawHome() {
	//hw.SetContents(strconv.Itoa(hw.indexPage))
	hw.DrawHomeInfo()
	hw.PrintControls()
}

// DrawHomeInfo draws the content for the top of field of the Help Window on the home screen, also draws the main home screen.
func (hw *HelpWindow) DrawHomeInfo() {
	// Top Field
	windowTitle := "Escaping Eden Help"
	pageInfo := "Home"
	hw.PrintLn(hw.X+1, hw.Y+1, windowTitle, "\033[1m")
	hw.PrintLn(hw.X+hw.Width-len(pageInfo)-1, hw.Y+1, pageInfo, "")

	helpTitle := "Welcome to the Escaping Eden Help Interface!"
	hw.PrintLn(hw.X+(hw.Width/2)-(len(helpTitle)/2)+1, hw.Y+3, helpTitle, "\033[1m")

	helpIntro := "Inside you'll find articles covering a variety of topics about Escaping Eden"
	hw.PrintLn(hw.X+(hw.Width/2)-(len(helpIntro)/2), hw.Y+5, helpIntro, "")

	artFile, err := util.OpenASCIIArtFile("assets/ascii/spire.txt")
	if err != nil {
		hw.Log.Println(logging.LogError, err.Error())
	} else {
		for y, line := range artFile.Data {
			for x, point := range line {
				printX := x + hw.X + (hw.Width / 2) - (artFile.Width / 2)
				printY := y + hw.Y + 7 // - artFile.Height
				if printX < hw.X+hw.Width+1 && printY < hw.Y+hw.Height+3 {
					hw.PrintChar(printX, printY, point.Character, point.EscapeCode)
				}
			}
		}
	}

}
