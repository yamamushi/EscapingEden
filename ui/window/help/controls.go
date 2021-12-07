package help

import (
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/types"
)

// HandleInput handles input events for HelpWindow.
func (hw *HelpWindow) HandleInput(input types.Input) {
	hw.hwMutex.Lock()
	defer hw.hwMutex.Unlock()

	if hw.GetActive() {
		//log.Println("Help Window Handling input")
	} else {
		//log.Println("Help window received an input event that it should not have")
		return
	}

	if hw.HelpPage == types.HelpPageIndex {
		hw.handleIndexInput(input)
		return
	}

	switch input.Type {

	case types.InputUp:
		//log.Println("Help Window Handling input up")
		hw.DecreaseContentPos()
		hw.ResetWindowDrawings()
		return
	case types.InputDown:
		//log.Println("Help Window Handling input down")
		hw.IncreaseContentPos()
		hw.ResetWindowDrawings()
		return
	case types.InputRight:
		//log.Println("Help Window Handling input right")
		if hw.HelpPage == types.HelpPageIndex {
			//hw.indexPage += 1  - Commented out right now because our index generation is very simple
			//hw.ForceConsoleRefresh()
		} else if hw.HelpPage != types.HelpPageIndex-1 {
			hw.HelpPage += 1
			hw.HandleStateChange()
		}
		return
	case types.InputLeft:
		//log.Println("Help Window Handling input left")
		if hw.HelpPage == types.HelpPageIndex {
			if hw.indexPage > 0 {
				//hw.indexPage -= 1 - Commented out right now because our index generation is very simple
				//hw.ForceConsoleRefresh()
			} else {
				if hw.HelpPage == types.HelpPageMain {
					hw.HelpPage -= 1
					hw.HandleStateChange()
				}
			}
		} else {
			if hw.HelpPage != types.HelpPageMain {
				hw.HelpPage -= 1
				hw.HandleStateChange()
			}
		}
		return
	case types.InputCharacter:
		switch input.Data {
		case "c":
			if hw.HelpPage == types.HelpPageIndex {
				hw.HelpPage = hw.LastHelpPage
			} else {
				message := messages.WindowMessage{Type: messages.WM_ConsoleCommand, Command: messages.WMC_ToggleHelp}
				hw.scrollInitialized = false
				hw.ConsoleSend <- message
				//log.Println("Help sent toggle window message to console")
			}

		case "h":
			hw.HelpPage = types.HelpPageMain
			hw.scrollInitialized = false
			hw.HandleStateChange()
			return
		case "i":
			// If we're not on the index then load it, otherwise toss it
			if hw.HelpPage != types.HelpPageIndex {
				hw.LastHelpPage = hw.HelpPage
				hw.HelpPage = types.HelpPageIndex
				hw.scrollInitialized = false
				hw.HandleStateChange()
			}
			return
		case "n":
			if hw.HelpPage != types.HelpPageIndex {
				hw.HelpPage += 1
				hw.HandleStateChange()
			}
			return
		case "p":
			if hw.HelpPage != types.HelpPageMain && hw.HelpPage != types.HelpPageIndex {
				hw.HelpPage -= 1
				hw.HandleStateChange()
			}
			return
		default:
			return
		}
	default:
		return
	}
}
