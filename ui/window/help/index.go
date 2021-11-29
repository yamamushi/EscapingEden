package help

import "strconv"

func (hw *HelpWindow) DrawIndexInfo() {
	// Top Field
	windowTitle := "Escaping Eden Help"
	pageInfo := "Index  (Page " + strconv.Itoa(hw.indexPage+1) + ")"
	hw.PrintLn(hw.X+1, hw.Y+1, windowTitle, "\033[1m")
	hw.PrintLn(hw.X+hw.Width-len(pageInfo)-1, hw.Y+1, pageInfo, "")
}

func (hw *HelpWindow) PrintIndexControls() {
	// Bottom Field
	separator := " | "
	homeCommand := "[h]ome"
	homeDistance := 1

	var prevCommand, nextCommand string
	var prevDistance, nextDistance int

	if hw.indexPage > 0 {
		prevCommand = "[ / ]revious"
		prevDistance = len(homeCommand+separator) + homeDistance
		nextCommand = "[ / ]ext"
		nextDistance = len(prevCommand+separator) + prevDistance
	} else {
		nextCommand = "[ / ]ext"
		nextDistance = len(homeCommand+separator) + homeDistance
	}

	closeCommand := "[ ]lose index"
	closeDistance := len(nextCommand+separator) + nextDistance

	var commandList string
	var shift int
	if hw.indexPage > 0 {
		commandList = homeCommand + separator + prevCommand + separator + nextCommand +
			separator + closeCommand

	} else {
		commandList = homeCommand + separator + nextCommand +
			separator + closeCommand
	}

	//hw.PrintLn(hw.X+(hw.Width/2)-(len(commandList)/2)-1, hw.Y+hw.Height, commandList, "")
	shift = (hw.Width / 2) - (len(commandList) / 2) - 1

	hw.PrintLn(hw.X+shift, hw.Y+hw.Height, commandList, "")
	hw.PrintChar(hw.X+shift+homeDistance, hw.Y+hw.Height, "h", "\033[1m")

	if hw.indexPage > 0 {
		hw.PrintChar(hw.X+shift+prevDistance, hw.Y+hw.Height, "\u25C4", "\033[1m")
		hw.PrintChar(hw.X+shift+prevDistance+2, hw.Y+hw.Height, "p", "\033[1m")
	}

	hw.PrintChar(hw.X+shift+nextDistance, hw.Y+hw.Height, "\u25BA", "\033[1m")
	hw.PrintChar(hw.X+shift+nextDistance+2, hw.Y+hw.Height, "n", "\033[1m")
	hw.PrintChar(hw.X+shift+closeDistance, hw.Y+hw.Height, "c", "\033[1m")
}

func (hw *HelpWindow) DrawIndex() {
	hw.SetContents(strconv.Itoa(hw.indexPage))
	hw.DrawIndexInfo()
	hw.PrintIndexControls()
}

//Need to finish index generation
