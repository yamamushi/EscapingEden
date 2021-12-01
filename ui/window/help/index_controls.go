package help

import (
	"github.com/yamamushi/EscapingEden/ui/types"
	"log"
)

func (hw *HelpWindow) handleIndexInput(input types.Input) {

	switch input.Type {
	case types.InputCharacter:
		switch input.Data {
		case "c":
			hw.HelpPage = hw.LastHelpPage
			hw.HandleStateChange()
			return
		case "H":
			hw.HelpPage = types.HelpPageMain
			hw.HandleStateChange()
			hw.scrollInitialized = false
			return
		}
	case types.InputUp:
		log.Println("Help Index handling up command")
		hw.IndexSelectionUp()
	case types.InputDown:
		log.Println("Help Index handling down command")
		hw.IndexSelectionDown()
	case types.InputReturn:
		log.Println("Help Index handling return command")
		hw.HelpPage = types.HelpPage(hw.indexSelection)
		hw.HandleStateChange()
	}

}
