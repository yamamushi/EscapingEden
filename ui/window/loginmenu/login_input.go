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
	switch lw.loginState {
	case LoginForgotPassword:
		lw.handleForgotPasswordInput(input)
	case LoginUserInfo:
		lw.handleLoginUserInfoInput(input)
	}

}

func (lw *LoginWindow) handleLoginUserInfoInput(input types.Input) {
	switch input.Type {
	case types.InputCharacter:
		lw.handleLoginCharInput(input.Data)
		return
	case types.InputBackspace:
		lw.handleLoginBackspace()
		return
	case types.InputUp:
		lw.loginNavOptionSelected = 0
		switch lw.loginMenuState {
		case LoginUserInfoNull:
			lw.loginMenuState = LoginUserInfoForgotPassword
		case LoginUserInfoForgotPassword:
			lw.loginMenuState = LoginUserInfoPassword
		case LoginUserInfoPassword:
			lw.loginMenuState = LoginUserInfoUsername
		case LoginUserInfoUsername:
			lw.loginMenuState = LoginUserInfoNull
		}
	case types.InputDown:
		lw.loginNavOptionSelected = 0
		switch lw.loginMenuState {
		case LoginUserInfoNull:
			lw.loginMenuState = LoginUserInfoUsername
		case LoginUserInfoPassword:
			lw.loginMenuState = LoginUserInfoForgotPassword
		case LoginUserInfoForgotPassword:
			lw.loginMenuState = LoginUserInfoNull
			lw.loginNavOptionSelected = 2 // submit
		case LoginUserInfoUsername:
			lw.loginMenuState = LoginUserInfoPassword
		}
	case types.InputLeft:
		//log.Println("Left arrow pressed")
		lw.loginMenuState = LoginUserInfoNull
		lw.loginNavOptionSelected = 1
		return
	case types.InputRight:
		//log.Println("Right arrow pressed")
		lw.loginMenuState = LoginUserInfoNull
		lw.loginNavOptionSelected = 2
		return
	case types.InputReturn:
		switch lw.loginMenuState {
		case LoginUserInfoUsername:
			lw.loginMenuState = LoginUserInfoPassword
		case LoginUserInfoPassword:
			lw.loginMenuState = LoginUserInfoNull
			lw.loginNavOptionSelected = 2
		case LoginUserInfoForgotPassword:
			// Reset our login data
			lw.loginMenuState = LoginUserInfoNull
			lw.loginNavOptionSelected = 0
			lw.loginSubmitData = LoginSubmitData{}

			// Send us to the forgot password state
			lw.loginState = LoginForgotPassword
			lw.loginForgotPasswordState = LoginForgotPasswordUsername
			lw.RequestFlushFromConsole()
		case LoginUserInfoNull:
			// Go Back to the main menu
			if lw.loginNavOptionSelected == 1 {
				lw.windowState = LoginWindowMenu
			}
			// Submit the login request
			if lw.loginNavOptionSelected == 2 {
				lw.loginSubmit()
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
	case LoginUserInfoUsername:
		if lw.loginSubmitData.Username != "" {
			lw.loginSubmitData.Username = lw.loginSubmitData.Username[:len(lw.loginSubmitData.Username)-1]
		}
	case LoginUserInfoPassword:
		if lw.loginSubmitData.Password != "" {
			lw.loginSubmitData.Password = lw.loginSubmitData.Password[:len(lw.loginSubmitData.Password)-1]
		}
	}
}

func (lw *LoginWindow) handleLoginCharInput(input string) {
	lw.loginSubmitMutex.Lock()
	defer lw.loginSubmitMutex.Unlock()

	switch lw.loginMenuState {
	case LoginUserInfoUsername:
		if len(lw.loginSubmitData.Username) < 128 {
			lw.loginSubmitData.Username += input
		}
		//lw.loginMenuState = LoginUserInfoPassword
	case LoginUserInfoPassword:
		if len(lw.loginSubmitData.Password) < 32 {
			lw.loginSubmitData.Password += input
		}
	}
}

func (lw *LoginWindow) handleForgotPasswordInput(input types.Input) {

	switch input.Type {
	case types.InputCharacter:
		lw.handleForgotPasswordCharInput(input.Data)
		return
	case types.InputBackspace:
		lw.handleForgotPasswordBackspace()
		return
	case types.InputUp:
		lw.loginForgotPasswordOptionSelected = 0
		switch lw.loginForgotPasswordState {
		case LoginForgotPasswordUsername:
			lw.loginForgotPasswordState = LoginForgotPasswordUsername
		case LoginForgotPasswordNull:
			lw.loginForgotPasswordState = LoginForgotPasswordUsername
		}
	case types.InputDown:
		lw.loginForgotPasswordOptionSelected = 0
		switch lw.loginForgotPasswordState {
		case LoginForgotPasswordNull:
			lw.loginForgotPasswordState = LoginForgotPasswordUsername
		case LoginForgotPasswordUsername:
			lw.loginForgotPasswordState = LoginForgotPasswordNull
			lw.loginForgotPasswordOptionSelected = 2 // submit
		}
	case types.InputLeft:
		//log.Println("Left arrow pressed")
		lw.loginForgotPasswordState = LoginForgotPasswordNull
		lw.loginForgotPasswordOptionSelected = 1
		return
	case types.InputRight:
		//log.Println("Right arrow pressed")
		lw.loginForgotPasswordState = LoginForgotPasswordNull
		lw.loginForgotPasswordOptionSelected = 2
		return
	case types.InputReturn:
		switch lw.loginForgotPasswordState {
		case LoginForgotPasswordUsername:
			lw.loginForgotPasswordState = LoginForgotPasswordNull
			lw.loginForgotPasswordOptionSelected = 2
		case LoginForgotPasswordNull:
			// Go Back to the login screen
			if lw.loginForgotPasswordOptionSelected == 1 {
				lw.loginState = LoginUserInfo
			}
			// Submit the forgot password request
			if lw.loginForgotPasswordOptionSelected == 2 {
				lw.loginForgotPasswordOptionSelected = 0
				lw.forgotPasswordSubmit()
				//lw.loginForgotPasswordState = LoginForgotPasswordPending
			}
			lw.loginForgotPasswordOptionSelected = 0
			// Whenever we switch to a different window state, we need to reset the console
			lw.RequestFlushFromConsole()
		}
		return
	}
}

func (lw *LoginWindow) handleForgotPasswordBackspace() {
	lw.loginForgotPasswordMutex.Lock()
	defer lw.loginForgotPasswordMutex.Unlock()

	switch lw.loginForgotPasswordState {
	case LoginForgotPasswordUsername:
		if lw.loginForgotPasswordData.Username != "" {
			lw.loginForgotPasswordData.Username = lw.loginForgotPasswordData.Username[:len(lw.loginForgotPasswordData.Username)-1]
		}
	}
}

func (lw *LoginWindow) handleForgotPasswordCharInput(input string) {
	lw.loginForgotPasswordMutex.Lock()
	defer lw.loginForgotPasswordMutex.Unlock()

	switch lw.loginForgotPasswordState {
	case LoginForgotPasswordUsername:
		if len(lw.loginForgotPasswordData.Username) < 128 {
			lw.loginForgotPasswordData.Username += input
		}
		//lw.loginMenuState = LoginUserInfoPassword
	}
}
