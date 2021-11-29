package help

import (
	"github.com/yamamushi/EscapingEden/ui/types"
	"log"
)

func (hw *HelpWindow) HandleInput(input types.Input) {
	hw.hwMutex.Lock()
	defer hw.hwMutex.Unlock()

	if hw.GetActive() {
		log.Println("Help Window Handling input")
	} else {
		log.Println("Help window received an input event that it should not have")
		return
	}

	switch input.Type {

	case types.InputUp:
		log.Println("Help Window Handling input up")
		hw.DecreaseContentPos()
		hw.ResetWindowDrawings()
		return
	case types.InputDown:
		log.Println("Help Window Handling input down")
		hw.IncreaseContentPos()
		hw.ResetWindowDrawings()
		return
	case types.InputRight:
		log.Println("Help Window Handling input right")
		if hw.HelpPage == types.HelpPageIndex {
			hw.indexPage += 1
		}
		hw.ForceConsoleRefresh()
		return
	case types.InputLeft:
		log.Println("Help Window Handling input left")
		if hw.HelpPage == types.HelpPageIndex {
			if hw.indexPage > 0 {
				hw.indexPage -= 1
			}
		}
		hw.ForceConsoleRefresh()
		return
	case types.InputCharacter:
		switch input.Data {
		case "c":
			message := types.ConsoleMessage{Type: "help", Message: "close"}
			hw.scrollInitialized = false
			hw.ConsoleSend <- message.String()
			log.Println("Help sent close message to console")
		case "h":
			hw.HelpPage = types.HelpPageMain
			hw.ResetWindowDrawings()
			hw.scrollInitialized = false
			hw.ForceConsoleRefresh()
			return
		case "i":
			hw.HelpPage = types.HelpPageIndex
			hw.ResetWindowDrawings()
			hw.scrollInitialized = false
			hw.ForceConsoleRefresh()
			return
		}

	default:
		log.Println("Unhandled Input event in Help Window")
	}
}

func (hw *HelpWindow) PrintControls() {
	// Bottom Field
	separator := " | "
	homeCommand := "[h]ome"
	homeDistance := 1

	var prevCommand, nextCommand string
	var prevDistance, nextDistance int
	if hw.HelpPage == types.HelpPageMain {
		nextCommand = "[ ]ext"
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
	closeCommand := "[ ]lose"
	closeDistance := len(indexCommand+separator) + indexDistance

	var commandList string
	var shift int
	if hw.HelpPage == types.HelpPageMain {
		commandList = homeCommand + separator + nextCommand + separator + scrollUpCommand +
			separator + scrollDownCommand + separator + indexCommand + separator + closeCommand
	} else {
		log.Println("Printing controls for non home page: ", prevDistance)
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
