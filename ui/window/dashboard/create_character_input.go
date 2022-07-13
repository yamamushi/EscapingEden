package dashboard

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/ui/types"
)

// handleMenuInput handles input for the login window
func (dw *DashboardWindow) handleCreateCharacterMenuInput(input types.Input) {
	dw.dwMutex.Lock()
	defer dw.dwMutex.Unlock()

	if !dw.GetActive() {
		return
	}

	switch dw.characterCreatorState {
	case CharacterCreatorDefaultNull:
		return // do nothing
	case CharacterCreatorFirstTimeLoginWelcome:
		switch input.Type {
		case types.InputReturn:
			// Continue to character details input screen
			dw.firstTimeLogin = false
			dw.RequestFlushFromConsole()
		}
	}

	switch input.Type {
	case types.InputCharacter:
		switch input.Data {
		case "a":
			dw.Log.Println(logging.LogInfo, "Letter A pressed on create character menu")
		}
	}
}
