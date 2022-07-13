package dashboard

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/ui/types"
	"unicode"
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

	case CharacterCreatorCharacterDetails:
		switch input.Type {
		case types.InputDown:
			if dw.charCreatorOptionSelected < 3 && dw.charColorOptionActive == false {
				dw.charCreatorOptionSelected += 1
				if dw.charCreatorOptionSelected == 3 {
					dw.charCreatorNavOptionSelected = 2
				} else {
					dw.charCreatorNavOptionSelected = 0
				}
			} else {
				if dw.charColorOption < 2 && dw.charColorOptionActive == true {
					dw.charColorOption += 1
				}
			}
			return
		case types.InputUp:
			if dw.charCreatorOptionSelected > 0 && dw.charColorOptionActive == false {
				dw.charCreatorOptionSelected -= 1
				dw.charCreatorNavOptionSelected = 0
			} else {
				if dw.charColorOption > 0 && dw.charColorOptionActive == true {
					dw.charColorOption -= 1
				}
			}
			return
		case types.InputRight:
			if dw.charCreatorOptionSelected == 2 {
				dw.charColorOptionActive = true
			} else {
				dw.charCreatorNavOptionSelected = 2 // submit
				dw.charCreatorOptionSelected = 3
			}
			return

		case types.InputLeft:
			if dw.charCreatorOptionSelected == 2 {
				dw.charColorOptionActive = false
			} else {
				dw.charCreatorNavOptionSelected = 1 // cancel
				dw.charCreatorOptionSelected = 3
			}
			return

		case types.InputCharacter:
			if dw.charCreatorOptionSelected == 1 {
				dw.handleCharacterCreatorUsernameInput(input.Data)
				dw.RequestFlushFromConsole()
				return
			}
		case types.InputBackspace:
			if dw.charCreatorOptionSelected == 1 {
				dw.handleCharacterCreatorBackspaceInput()
				dw.RequestFlushFromConsole()
				return
			}

		case types.InputReturn:
			if dw.charCreatorNavOptionSelected == 1 {
				// go back to the menu
				dw.windowState = DashboardMainMenu
				dw.characterCreatorState = CharacterCreatorDefaultNull
				dw.charCreatorNavOptionSelected = 0 // reset to default
				dw.charCreatorOptionSelected = 0    // reset to default
				dw.charColorOption = 0              // reset color option
				dw.charCreatorUsername = ""         // reset username
				dw.charCreatorUsernameError = ""    // reset username error
				if dw.GetUserInfoField("lastcharacter") == "" {
					dw.firstTimeLogin = true // reset first time login
				}
				dw.RequestFlushFromConsole()
				return
			}
			if dw.charCreatorNavOptionSelected == 2 {
				// Go to character details confirmation screen
				dw.Log.Println(logging.LogInfo, "Character submitted for confirmation")
			}
		}
	}
}

func (dw *DashboardWindow) handleCharacterCreatorUsernameInput(character string) {
	// If character is not a letter
	for _, r := range character {
		if !unicode.IsLetter(r) {
			dw.charCreatorUsernameError = "Name must be letters only"
			return
		}
	}
	dw.charCreatorUsernameError = ""

	if len(dw.charCreatorUsername) < 16 {
		dw.charCreatorUsername += character
	}
}

func (dw *DashboardWindow) handleCharacterCreatorBackspaceInput() {
	dw.charCreatorUsernameError = ""
	if len(dw.charCreatorUsername) > 0 {
		dw.charCreatorUsername = dw.charCreatorUsername[:len(dw.charCreatorUsername)-1]
	}
}
