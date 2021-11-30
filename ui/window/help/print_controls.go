package help

import "github.com/yamamushi/EscapingEden/ui/types"

// PrintControls prints the default controls for the Help window.
func (hw *HelpWindow) PrintControls() {
	// Bottom Field
	separator := " | "
	homeCommand := "[h]ome"
	homeDistance := 1

	var prevCommand, nextCommand string
	var prevDistance, nextDistance int
	if hw.HelpPage == types.HelpPageMain {
		nextCommand = "[ / ]ext"
		nextDistance = len(homeCommand+separator) + homeDistance
	} else {
		prevCommand = "[ / ]revious"
		prevDistance = len(homeCommand+separator) + homeDistance
		nextCommand = "[ / ]ext"
		nextDistance = len(prevCommand+separator) + prevDistance
	}

	scrollUpCommand := "[ ] Scroll Up"
	scrollUpDistance := len(nextCommand+separator) + nextDistance
	scrollDownCommand := "[ ] Scroll Down"
	scrollDownDistance := len(scrollUpCommand+separator) + scrollUpDistance
	indexCommand := "[ ]ndex"
	indexDistance := len(scrollDownCommand+separator) + scrollDownDistance
	closeCommand := "[ ]lose help"
	closeDistance := len(indexCommand+separator) + indexDistance

	var commandList string
	var shift int
	if hw.HelpPage == types.HelpPageMain {
		commandList = homeCommand + separator + nextCommand + separator + scrollUpCommand +
			separator + scrollDownCommand + separator + indexCommand + separator + closeCommand
	} else {
		//log.Println("Printing controls for non home page: ", prevDistance)
		commandList = homeCommand + separator + prevCommand + separator + nextCommand + separator + scrollUpCommand +
			separator + scrollDownCommand + separator + indexCommand + separator + closeCommand
	}
	shift = (hw.Width / 2) - (len(commandList) / 2) - 1

	hw.PrintLn(hw.X+shift, hw.Y+hw.Height, commandList, "")
	hw.PrintChar(hw.X+shift+homeDistance, hw.Y+hw.Height, "h", "\033[1m")

	if hw.HelpPage != types.HelpPageMain {
		hw.PrintChar(hw.X+shift+prevDistance, hw.Y+hw.Height, "\u25C4", "\033[1m")
		hw.PrintChar(hw.X+shift+prevDistance+2, hw.Y+hw.Height, "p", "\033[1m")
	}

	hw.PrintChar(hw.X+shift+nextDistance, hw.Y+hw.Height, "\u25BA", "\033[1m")
	hw.PrintChar(hw.X+shift+nextDistance+2, hw.Y+hw.Height, "n", "\033[1m")
	hw.PrintChar(hw.X+shift+scrollUpDistance, hw.Y+hw.Height, "\u25B2", "\033[1m")
	hw.PrintChar(hw.X+shift+scrollDownDistance, hw.Y+hw.Height, "\u25BC", "\033[1m")
	hw.PrintChar(hw.X+shift+indexDistance, hw.Y+hw.Height, "i", "\033[1m")
	hw.PrintChar(hw.X+shift+closeDistance, hw.Y+hw.Height, "c", "\033[1m")
}
