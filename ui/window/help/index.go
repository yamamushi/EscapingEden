package help

import (
	"github.com/yamamushi/EscapingEden/ui/types"
	"strconv"
	"strings"
)

type IndexPageNames []string

// DrawIndex is a wrapper around the functions necessary to draw the current index page
func (hw *HelpWindow) DrawIndex() {
	//hw.SetContents(strconv.Itoa(hw.indexPage))
	hw.DrawIndexInfo()
	hw.PrintIndexControls()
}

// PrintIndexControls prints the controls for the help index
func (hw *HelpWindow) PrintIndexControls() {
	// Bottom Field
	separator := " | "
	homeCommand := "[H]ome"
	homeDistance := 1

	selectCommand := "[Enter] Select Entry"
	selectDistance := len(separator) + len(homeCommand) + homeDistance

	closeCommand := "[ ]lose index"
	closeDistance := len(separator) + len(selectCommand) + selectDistance

	var shift int
	commandList := homeCommand + separator + selectCommand + separator + closeCommand

	//hw.PrintLn(hw.X+(hw.Width/2)-(len(commandList)/2)-1, hw.Y+hw.Height, commandList, "")
	shift = (hw.Width / 2) - (len(commandList) / 2) - 1

	hw.PrintLn(hw.X+shift, hw.Y+hw.Height, commandList, "")
	hw.PrintChar(hw.X+shift+homeDistance, hw.Y+hw.Height, "H", hw.Terminal.Bold())
	hw.PrintLn(hw.X+shift+selectDistance, hw.Y+hw.Height, "Enter", hw.Terminal.Bold())
	hw.PrintChar(hw.X+shift+closeDistance, hw.Y+hw.Height, "c", hw.Terminal.Bold())
}

// DrawIndexInfo draws the index page scraping from types.HelpPage const types
func (hw *HelpWindow) DrawIndexInfo() {
	// Top Field
	windowTitle := "Escaping Eden Help"
	var pageInfo string
	if hw.indexPage > 0 {
		pageInfo = "Index  (Page " + strconv.Itoa(hw.indexPage) + ")"
	} else {
		pageInfo = "Index Main"
	}
	hw.PrintLn(hw.X+1, hw.Y+1, windowTitle, hw.Terminal.Bold())
	hw.PrintLn(hw.X+hw.Width-len(pageInfo)-1, hw.Y+1, pageInfo, "")

	hw.PrintLn(hw.X+3, hw.Y+4, "Welcome to the help index, you can use your arrow keys or h/j/k/l keys to navigate", "")
	hw.PrintChar(hw.X+61, hw.Y+4, "h", hw.Terminal.Bold())
	hw.PrintChar(hw.X+63, hw.Y+4, "j", hw.Terminal.Bold())
	hw.PrintChar(hw.X+65, hw.Y+4, "k", hw.Terminal.Bold())
	hw.PrintChar(hw.X+67, hw.Y+4, "l", hw.Terminal.Bold())

	hw.generateIndexCommands()

}

// Right now we don't have that many help pages so the index generation is going to be very simple
// As we add more pages, I'll update this to do more interesting things :)

func (hw *HelpWindow) generateIndexCommands() {
	hw.indexMutex.Lock()
	defer hw.indexMutex.Unlock()

	var pageNames []string

	for i := 0; i < int(types.HelpPageIndex); i++ {
		helpPageName := types.HelpPage(i).String()
		pageNames = append(pageNames, helpPageName)
	}

	for i := 0; i < len(pageNames); i++ {
		hw.PrintLn(hw.X+3, hw.Y+7+i, "[ ] "+strings.Title(pageNames[i]), "")
	}

	hw.PrintChar(hw.X+4, hw.Y+7+hw.indexSelection, "*", hw.Terminal.Bold())
}
