package dashboard

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/ui/types"
)

// handleMenuInput handles input for the login window
func (dw *DashboardWindow) handleCreateCharacterInput(input types.Input) {
	dw.dwMutex.Lock()
	defer dw.dwMutex.Unlock()

	if !dw.GetActive() {
		return
	}

	switch input.Type {
	case types.InputCharacter:
		switch input.Data {
		case "a":
			dw.Log.Println(logging.LogInfo, "Letter A pressed on create character menu")
		}
	}
}
