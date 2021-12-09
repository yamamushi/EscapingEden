package login

import "github.com/yamamushi/EscapingEden/ui/types"

// handleLoginInput handles input for the login window
func (lw *LoginWindow) handleLoginInput(input types.Input) {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	if !lw.GetActive() {
		return
	}

	//input = strings.ToLower(input[:1])

	switch input.Type {
	case types.InputCharacter:
		lw.handleLoginCharInput(input.Data)
		return
	case types.InputBackspace:
		lw.handleLoginBackspace()
		return
	case types.InputUp:
		switch lw.loginMenuState {
		case LoginMenuNull:
			lw.loginMenuState = LoginMenuPassword
		case LoginMenuPassword:
			lw.loginMenuState = LoginMenuEmail
		case LoginMenuEmail:
			lw.loginMenuState = LoginMenuNull
		}
	case types.InputDown:
		switch lw.loginMenuState {
		case LoginMenuNull:
			lw.loginMenuState = LoginMenuEmail
		case LoginMenuPassword:
			lw.loginMenuState = LoginMenuNull
			lw.loginNavOptionSelected = 2 // submit
		case LoginMenuEmail:
			lw.loginMenuState = LoginMenuPassword
		}
	case types.InputLeft:
		//log.Println("Left arrow pressed")
		lw.loginMenuState = LoginMenuNull
		lw.loginNavOptionSelected = 1
		return
	case types.InputRight:
		//log.Println("Right arrow pressed")
		lw.loginMenuState = LoginMenuNull
		lw.loginNavOptionSelected = 2
		return
	case types.InputReturn:
		switch lw.loginMenuState {
		case LoginMenuEmail:
			lw.loginMenuState = LoginMenuPassword
		case LoginMenuPassword:
			lw.loginMenuState = LoginMenuNull
			lw.loginNavOptionSelected = 2
		case LoginMenuNull:
			// Go Back to the main menu
			if lw.loginNavOptionSelected == 1 {
				lw.windowState = LoginWindowMenu
			}
			// Submit the login request
			if lw.loginNavOptionSelected == 2 {
				lw.LoginSubmit()
				lw.loginState = LoginPending
			}
			lw.loginNavOptionSelected = 0
			// Whenever we switch to a different window state, we need to reset the console
			lw.RequestFlushFromConsole()
		}
		return
	default:
		return
	}
}

func (lw *LoginWindow) handleLoginBackspace() {
	lw.loginSubmitMutex.Lock()
	defer lw.loginSubmitMutex.Unlock()

	switch lw.loginMenuState {
	case LoginMenuEmail:
		if lw.loginSubmitData.Email != "" {
			lw.loginSubmitData.Email = lw.loginSubmitData.Email[:len(lw.loginSubmitData.Email)-1]
		}
	case LoginMenuPassword:
		if lw.loginSubmitData.Password != "" {
			lw.loginSubmitData.Password = lw.loginSubmitData.Password[:len(lw.loginSubmitData.Password)-1]
		}
	}
}

func (lw *LoginWindow) handleLoginCharInput(input string) {
	lw.loginSubmitMutex.Lock()
	defer lw.loginSubmitMutex.Unlock()

	switch lw.loginMenuState {
	case LoginMenuEmail:
		if len(lw.loginSubmitData.Email) < 128 {
			lw.loginSubmitData.Email += input
		}
		//lw.loginMenuState = LoginMenuPassword
	case LoginMenuPassword:
		if len(lw.loginSubmitData.Password) < 32 {
			lw.loginSubmitData.Password += input
		}
	}
}
