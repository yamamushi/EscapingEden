package login

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/ui/types"
)

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
	case LoginForgotPasswordPending:
		lw.handleForgotPasswordPendingInput(input)
	case LoginUserInfo:
		lw.handleLoginUserInfoInput(input)
	case LoginForgotPasswordSuccess:
		lw.handleForgotPasswordSuccessInput(input)
	case LoginForgotPasswordFailed:
		lw.handleForgotPasswordFailedInput(input)
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
			lw.loginSubmitData = LoginSubmitData{} // reset our login data so if we come back this screen is clean

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
		case LoginForgotPasswordDiscord:
			lw.loginForgotPasswordState = LoginForgotPasswordUsername
		case LoginForgotPasswordNull:
			lw.loginForgotPasswordState = LoginForgotPasswordDiscord
		}
	case types.InputDown:
		lw.loginForgotPasswordOptionSelected = 0
		switch lw.loginForgotPasswordState {
		case LoginForgotPasswordNull:
			lw.loginForgotPasswordState = LoginForgotPasswordUsername
		case LoginForgotPasswordUsername:
			lw.loginForgotPasswordState = LoginForgotPasswordDiscord
		case LoginForgotPasswordDiscord:
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
			lw.loginForgotPasswordState = LoginForgotPasswordDiscord
		case LoginForgotPasswordDiscord:
			lw.loginForgotPasswordState = LoginForgotPasswordNull
			lw.loginForgotPasswordOptionSelected = 2
		case LoginForgotPasswordNull:
			// Go Back to the login screen
			if lw.loginForgotPasswordOptionSelected == 1 {
				lw.loginMenuState = LoginUserInfoUsername
				lw.loginState = LoginUserInfo
			}
			// Submit the forgot password request
			if lw.loginForgotPasswordOptionSelected == 2 {
				lw.loginForgotPasswordOptionSelected = 0
				lw.forgotPasswordSubmit()
				lw.loginState = LoginForgotPasswordPending
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
	case LoginForgotPasswordDiscord:
		if lw.loginForgotPasswordData.DiscordUser != "" {
			lw.loginForgotPasswordData.DiscordUser = lw.loginForgotPasswordData.DiscordUser[:len(lw.loginForgotPasswordData.DiscordUser)-1]
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
	case LoginForgotPasswordDiscord:
		if len(lw.loginForgotPasswordData.Username) < 128 {
			lw.loginForgotPasswordData.DiscordUser += input
		}
	}
}

func (lw *LoginWindow) handleForgotPasswordPendingInput(input types.Input) {
	switch input.Type {
	case types.InputCharacter:
		lw.handleForgotPasswordPendingCharInput(input.Data)
		return
	case types.InputBackspace:
		lw.handleForgotPasswordPendingBackspace()
		return
	case types.InputUp:
		lw.loginForgotPasswordOptionSelected = 0
		switch lw.loginForgotPasswordPendingState {
		case LoginForgotPendingCode:
			lw.loginForgotPasswordPendingState = LoginForgotPendingCode
		case LoginForgotPendingNull:
			lw.loginForgotPasswordPendingState = LoginForgotPendingCode
		}
	case types.InputDown:
		lw.loginForgotPasswordPendingOptionSelected = 0
		switch lw.loginForgotPasswordPendingState {
		case LoginForgotPendingNull:
			lw.loginForgotPasswordPendingState = LoginForgotPendingCode
		case LoginForgotPendingCode:
			lw.loginForgotPasswordPendingState = LoginForgotPendingNull
			lw.loginForgotPasswordPendingOptionSelected = 2 // submit
		}
	case types.InputLeft:
		//log.Println("Left arrow pressed")
		lw.loginForgotPasswordPendingState = LoginForgotPendingNull
		lw.loginForgotPasswordPendingOptionSelected = 1
		return
	case types.InputRight:
		//log.Println("Right arrow pressed")
		lw.loginForgotPasswordPendingState = LoginForgotPendingNull
		lw.loginForgotPasswordPendingOptionSelected = 2
		return
	case types.InputReturn:
		switch lw.loginForgotPasswordPendingState {
		case LoginForgotPendingCode:
			lw.loginForgotPasswordPendingState = LoginForgotPendingNull
			lw.loginForgotPasswordPendingOptionSelected = 2
		case LoginForgotPendingNull:
			lw.loginForgotPasswordPendingState = LoginForgotPendingNull
			// Go Back to the login screen
			if lw.loginForgotPasswordPendingOptionSelected == 1 {
				// lw.windowState = LoginWindowMenu
				lw.loginState = LoginUserInfo
			}
			// Submit the forgot password request
			if lw.loginForgotPasswordPendingOptionSelected == 2 {
				lw.loginForgotPasswordPendingOptionSelected = 0
				lw.forgotPasswordValidate()
				//lw.loginForgotPasswordState = LoginForgotPasswordPending
			}
			lw.loginForgotPasswordPendingOptionSelected = 0
			// Whenever we switch to a different window state, we need to reset the console
			lw.RequestFlushFromConsole()
		}
		return
	}
}

func (lw *LoginWindow) handleForgotPasswordPendingBackspace() {
	lw.loginForgotPasswordPendingMutex.Lock()
	defer lw.loginForgotPasswordPendingMutex.Unlock()

	switch lw.loginForgotPasswordPendingState {
	case LoginForgotPendingCode:
		if lw.loginForgotPasswordPendingData.Code != "" {
			lw.loginForgotPasswordPendingData.Code = lw.loginForgotPasswordPendingData.Code[:len(lw.loginForgotPasswordPendingData.Code)-1]
		}
	}
}

func (lw *LoginWindow) handleForgotPasswordPendingCharInput(input string) {
	lw.loginForgotPasswordPendingMutex.Lock()
	defer lw.loginForgotPasswordPendingMutex.Unlock()

	switch lw.loginForgotPasswordPendingState {
	case LoginForgotPendingCode:
		if len(lw.loginForgotPasswordPendingData.Code) < 128 {
			lw.loginForgotPasswordPendingData.Code += input
		}
	}
}

func (lw *LoginWindow) handleForgotPasswordSuccessInput(input types.Input) {
	switch input.Type {
	case types.InputCharacter:
		lw.handleForgotPasswordSuccessCharInput(input.Data)
		return
	case types.InputBackspace:
		lw.handleForgotPasswordSuccessBackspace()
		return
	case types.InputDown:
		lw.loginForgotPasswordSuccessOptionSelected = 0
		switch lw.loginForgotPasswordSuccessState {
		case LoginForgotPasswordSuccessNull:
			lw.loginForgotPasswordSuccessState = LoginForgotPasswordSuccessEntry
		case LoginForgotPasswordSuccessEntry:
			lw.loginForgotPasswordSuccessState = LoginForgotPasswordSuccessConfirm
		case LoginForgotPasswordSuccessConfirm:
			lw.loginForgotPasswordSuccessState = LoginForgotPasswordSuccessNull
			lw.loginForgotPasswordSuccessOptionSelected = 2
		}
	case types.InputUp:
		lw.loginForgotPasswordSuccessOptionSelected = 0
		switch lw.loginForgotPasswordSuccessState {
		case LoginForgotPasswordSuccessEntry:
			lw.loginForgotPasswordSuccessState = LoginForgotPasswordSuccessNull
			lw.loginForgotPasswordSuccessOptionSelected = 2
		case LoginForgotPasswordSuccessConfirm:
			lw.loginForgotPasswordSuccessState = LoginForgotPasswordSuccessEntry
		case LoginForgotPasswordSuccessNull:
			lw.loginForgotPasswordSuccessState = LoginForgotPasswordSuccessConfirm
		}
	case types.InputLeft:
		lw.loginForgotPasswordSuccessState = LoginForgotPasswordSuccessNull
		lw.loginForgotPasswordSuccessOptionSelected = 0
		return
	case types.InputRight:
		lw.loginForgotPasswordSuccessState = LoginForgotPasswordSuccessNull
		lw.loginForgotPasswordSuccessOptionSelected = 2
		return
	case types.InputReturn:
		switch lw.loginForgotPasswordSuccessState {
		case LoginForgotPasswordSuccessEntry:
			lw.loginForgotPasswordSuccessState = LoginForgotPasswordSuccessConfirm
		case LoginForgotPasswordSuccessConfirm:
			lw.loginForgotPasswordSuccessState = LoginForgotPasswordSuccessNull
			lw.loginForgotPasswordSuccessOptionSelected = 2
		case LoginForgotPasswordSuccessNull:
			// Submit the forgot password request
			if lw.loginForgotPasswordSuccessOptionSelected == 2 {
				lw.loginForgotPasswordSuccessOptionSelected = LoginForgotPendingNull
				lw.Log.Println(logging.LogInfo, "Password Reset Submitted")
				//lw.loginForgotPasswordState = LoginForgotPasswordPending
			}
			lw.loginForgotPasswordSuccessOptionSelected = 0
			// Whenever we switch to a different window state, we need to reset the console
			lw.RequestFlushFromConsole()
		}
		return
	}
}

func (lw *LoginWindow) handleForgotPasswordSuccessBackspace() {
	lw.loginForgotPasswordSuccessMutex.Lock()
	defer lw.loginForgotPasswordSuccessMutex.Unlock()

	switch lw.loginForgotPasswordSuccessState {
	case LoginForgotPasswordSuccessEntry:
		if lw.loginForgotPasswordSuccessData.Password != "" {
			lw.loginForgotPasswordSuccessData.Password = lw.loginForgotPasswordSuccessData.Password[:len(lw.loginForgotPasswordSuccessData.Password)-1]
		}
	case LoginForgotPasswordSuccessConfirm:
		if lw.loginForgotPasswordSuccessData.PasswordConfirm != "" {
			lw.loginForgotPasswordSuccessData.PasswordConfirm = lw.loginForgotPasswordSuccessData.PasswordConfirm[:len(lw.loginForgotPasswordSuccessData.PasswordConfirm)-1]
		}
	}
}

func (lw *LoginWindow) handleForgotPasswordSuccessCharInput(input string) {
	lw.loginForgotPasswordSuccessMutex.Lock()
	defer lw.loginForgotPasswordSuccessMutex.Unlock()

	switch lw.loginForgotPasswordSuccessState {
	case LoginForgotPasswordSuccessEntry:
		if len(lw.loginForgotPasswordSuccessData.Password) < 128 {
			lw.loginForgotPasswordSuccessData.Password += input
		}
	case LoginForgotPasswordSuccessConfirm:
		if len(lw.loginForgotPasswordSuccessData.PasswordConfirm) < 128 {
			lw.loginForgotPasswordSuccessData.PasswordConfirm += input
		}
	}
}

func (lw *LoginWindow) handleForgotPasswordFailedInput(input types.Input) {
	switch input.Type {
	case types.InputDown:
		lw.loginForgotPasswordFailedOptionSelected = 1
	case types.InputUp:
		lw.loginForgotPasswordFailedOptionSelected = 0
	case types.InputLeft:
		lw.loginForgotPasswordFailedOptionSelected = 0
	case types.InputRight:
		lw.loginForgotPasswordFailedOptionSelected = 1
	case types.InputReturn:
		switch lw.loginForgotPasswordFailedOptionSelected {
		case 0:
			// do nothing
		case 1:
			// Go back to the code entry screen
			lw.loginForgotPasswordPendingData.Code = ""                 // Reset the code so we can enter it fresh
			lw.loginForgotPasswordPendingState = LoginForgotPendingCode // reset to the entry line
			lw.loginState = LoginForgotPasswordPending
		}

	}
}
