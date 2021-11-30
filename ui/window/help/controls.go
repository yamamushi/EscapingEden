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
			hw.ForceConsoleRefresh()
		} else if hw.HelpPage != types.HelpPageIndex-1 {
			hw.HelpPage += 1
			hw.ForceConsoleRefresh()
		}
		return
	case types.InputLeft:
		log.Println("Help Window Handling input left")
		if hw.HelpPage == types.HelpPageIndex {
			if hw.indexPage > 0 {
				hw.indexPage -= 1
				hw.ForceConsoleRefresh()
			} else {
				if hw.HelpPage == types.HelpPageMain {
					hw.HelpPage -= 1
					hw.ForceConsoleRefresh()
				}
			}
		} else {
			if hw.HelpPage != types.HelpPageMain {
				hw.HelpPage -= 1
				hw.ForceConsoleRefresh()
			}
		}
		return
	case types.InputCharacter:
		switch input.Data {
		case "c":
			if hw.HelpPage == types.HelpPageIndex {
				hw.HelpPage = hw.LastHelpPage
			} else {
				message := types.ConsoleMessage{Type: "help", Message: "close"}
				hw.scrollInitialized = false
				hw.ConsoleSend <- message.String()
				log.Println("Help sent close message to console")
			}

		case "h":
			hw.HelpPage = types.HelpPageMain
			hw.ResetWindowDrawings()
			hw.scrollInitialized = false
			hw.ForceConsoleRefresh()
			return
		case "i":
			// If we're not on the index then load it, otherwise toss it
			if hw.HelpPage != types.HelpPageIndex {
				hw.LastHelpPage = hw.HelpPage
				hw.HelpPage = types.HelpPageIndex
				hw.ResetWindowDrawings()
				hw.scrollInitialized = false
				hw.ForceConsoleRefresh()
			}
			return
		case "n":
			if hw.HelpPage != types.HelpPageIndex {
				hw.HelpPage += 1
				hw.ForceConsoleRefresh()
			}
			return
		case "p":
			if hw.HelpPage != types.HelpPageMain && hw.HelpPage != types.HelpPageIndex {
				hw.HelpPage -= 1
				hw.ForceConsoleRefresh()
			}
			return
		default:
			return
		}
	default:
		log.Println("Unhandled Input event in Help Window")
	}
}
