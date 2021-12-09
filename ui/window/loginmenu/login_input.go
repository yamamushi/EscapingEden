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
			lw.loginMenuState = LoginUserInfoEmail
		case LoginUserInfoEmail:
			lw.loginMenuState = LoginUserInfoNull
		}
	case types.InputDown:
		lw.loginNavOptionSelected = 0
		switch lw.loginMenuState {
		case LoginUserInfoNull:
			lw.loginMenuState = LoginUserInfoEmail
		case LoginUserInfoPassword:
			lw.loginMenuState = LoginUserInfoForgotPassword
		case LoginUserInfoForgotPassword:
			lw.loginMenuState = LoginUserInfoNull
			lw.loginNavOptionSelected = 2 // submit
		case LoginUserInfoEmail:
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
		case LoginUserInfoEmail:
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
			lw.loginForgotPasswordState = LoginForgotPasswordEmail
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
	case LoginUserInfoEmail:
		if lw.loginSubmitData.Email != "" {
			lw.loginSubmitData.Email = lw.loginSubmitData.Email[:len(lw.loginSubmitData.Email)-1]
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
	case LoginUserInfoEmail:
		if len(lw.loginSubmitData.Email) < 128 {
			lw.loginSubmitData.Email += input
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
		case LoginForgotPasswordEmail:
			lw.loginForgotPasswordState = LoginForgotPasswordEmail
		case LoginForgotPasswordNull:
			lw.loginForgotPasswordState = LoginForgotPasswordEmail
		}
	case types.InputDown:
		lw.loginForgotPasswordOptionSelected = 0
		switch lw.loginForgotPasswordState {
		case LoginForgotPasswordNull:
			lw.loginForgotPasswordState = LoginForgotPasswordEmail
		case LoginForgotPasswordEmail:
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
		case LoginForgotPasswordEmail:
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
	case LoginForgotPasswordEmail:
		if lw.loginForgotPasswordData.Email != "" {
			lw.loginForgotPasswordData.Email = lw.loginForgotPasswordData.Email[:len(lw.loginForgotPasswordData.Email)-1]
		}
	}
}

func (lw *LoginWindow) handleForgotPasswordCharInput(input string) {
	lw.loginForgotPasswordMutex.Lock()
	defer lw.loginForgotPasswordMutex.Unlock()

	switch lw.loginForgotPasswordState {
	case LoginForgotPasswordEmail:
		if len(lw.loginForgotPasswordData.Email) < 128 {
			lw.loginForgotPasswordData.Email += input
		}
		//lw.loginMenuState = LoginUserInfoPassword
	}
}
