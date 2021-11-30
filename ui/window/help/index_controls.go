package help

import "github.com/yamamushi/EscapingEden/ui/types"

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
	}

}
