package login

import (
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/util"
)

type LoginState int

const (
	LoginNull LoginState = iota
	LoginUserInfo
	LoginPending
	LoginForgotPassword
	LoginForgotPasswordPending
	LoginForgotPasswordSuccess
	LoginForgotPasswordFailed
)

// LoginUserInfoState is an enum for storing login state
type LoginUserInfoState int

const (
	LoginUserInfoUsername LoginUserInfoState = iota
	LoginUserInfoPassword
	LoginUserInfoForgotPassword
	LoginUserInfoNull
)

type LoginSubmitData struct {
	Username string
	Password string
	Error    string
}

type LoginForgotPasswordState int

const (
	LoginForgotPasswordUsername LoginForgotPasswordState = iota
	LoginForgotPasswordDiscord
	LoginForgotPasswordNull
)

type LoginForgotPasswordPendingState int

const (
	LoginForgotPendingCode = iota
	LoginForgotPendingNull
)

type LoginForgotPasswordSuccessState int

const (
	LoginForgotPasswordSuccessNull = iota
	LoginForgotPasswordSuccessEntry
	LoginForgotPasswordSuccessConfirm
)

type LoginForgotPasswordSuccessData struct {
	Password        string
	PasswordConfirm string
	Error           string
}

// drawLoginMenu draws the login window
func (lw *LoginWindow) drawLoginMenu() {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	//lw.SetContents("handleLogin")

	switch lw.loginState {
	case LoginUserInfo:
		lw.drawLoginMenuUserInfo()
		return
	case LoginPending:
		lw.drawLoginMenuPending()
		return
	case LoginForgotPassword:
		lw.drawLoginMenuForgotPassword()
		return
	case LoginForgotPasswordPending:
		lw.drawLoginMenuForgotPasswordPending()
		return
	case LoginForgotPasswordSuccess:
		lw.drawLoginMenuForgotPasswordSuccess()
		return
	case LoginForgotPasswordFailed:
		lw.drawLoginMenuForgotPasswordFailed()
		return
	}

}

func (lw *LoginWindow) drawLoginMenuUserInfo() {
	lw.loginSubmitMutex.Lock()
	defer lw.loginSubmitMutex.Unlock()

	errorFG := util.RGBCode(255, 255, 255)
	errorBG := util.RGBCode(255, 0, 0)

	if lw.loginMenuState == LoginUserInfoUsername {
		lw.PrintLn(lw.X+6, lw.Y+5, "Username:", lw.Terminal.Bold())
	} else {
		lw.PrintLn(lw.X+6, lw.Y+5, "Username:", "")
	}
	username := ""
	// We only want the last 12 characters of the username
	if len(lw.loginSubmitData.Username) > 12 {
		username = lw.loginSubmitData.Username[len(lw.loginSubmitData.Username)-12:]
	} else {
		username = lw.loginSubmitData.Username
	}
	lw.PrintLn(lw.X+16, lw.Y+5, username, "")

	if lw.loginMenuState == LoginUserInfoPassword {
		lw.PrintLn(lw.X+6, lw.Y+6, "Password:", lw.Terminal.Bold())
	} else {
		lw.PrintLn(lw.X+6, lw.Y+6, "Password:", "")
	}
	for i := 0; i < len(lw.loginSubmitData.Password) && i < 12; i++ {
		lw.PrintLn(lw.X+16+i, lw.Y+6, "*", "")
	}

	// Draw the back and submit buttons
	if lw.loginMenuState == LoginUserInfoForgotPassword {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+7, lw.Y+8, "<Forgot Password>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+7, lw.Y+8, "<Forgot Password>", lw.Terminal.Bold())
	}

	if lw.loginMenuMessage != "" {
		// Draw the message in green
		lw.PrintLn(lw.X+7, lw.Y+10, lw.loginMenuMessage, lw.Terminal.Bold()+util.RGBCode(0, 255, 0).FG())
	}

	//lw.loginSubmitData.Error = "This is a test error message"

	if lw.loginSubmitData.Error != "" {
		lw.PrintLnColor(lw.X+5, lw.Y+10, "Error logging in: "+lw.loginSubmitData.Error, errorFG.FG()+errorBG.BG())
	}

	// Draw the back and submit buttons
	if lw.loginNavOptionSelected == 1 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", lw.Terminal.Bold())
	}

	if lw.loginNavOptionSelected == 2 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+lw.Width-12, lw.Y+lw.Height, "<Submit>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+lw.Width-12, lw.Y+lw.Height, "<Submit>", lw.Terminal.Bold())
	}
}

func (lw *LoginWindow) drawLoginMenuPending() {
	lw.loginStatusMutex.Lock()
	defer lw.loginStatusMutex.Unlock()

	if lw.loginResponseReceived {
		if lw.loginSubmitData.Error != "" {
			// If we got an error, we go back to the user info screen, display that error and wait
			lw.loginState = LoginUserInfo
		} else {
			// If we didn't get an error, we load our user info screen, and we also notify the console the user has logged in
			lw.windowState = LoginWindowUserDashboard
			// Create a console message with type Console_Message_LoginUser, we don't pack any data with this message (yet, TBD)
			msg := messages.WindowMessage{Type: messages.WM_ConsoleCommand, Command: messages.WMC_SetLoggedIn, TargetID: lw.GetID()}
			// Send the message to the console so that we can enable the full dashboard
			lw.SendToConsole(msg)
		}
		lw.loginResponseReceived = false
		lw.RequestFlushFromConsole()
	}

	lw.PrintLn(lw.X+lw.Width/2-5, lw.Y+lw.Height/2, "Login Pending...", lw.Terminal.Bold())
}

func (lw *LoginWindow) drawLoginMenuForgotPassword() {
	lw.loginForgotPasswordMutex.Lock()
	defer lw.loginForgotPasswordMutex.Unlock()

	lw.PrintLn(lw.X+6, lw.Y+2, "Please enter your username and discord account to begin the password reset process.", "")

	if lw.loginForgotPasswordState == LoginForgotPasswordUsername {
		lw.PrintLn(lw.X+9, lw.Y+5, "Username:", lw.Terminal.Bold())
	} else {
		lw.PrintLn(lw.X+9, lw.Y+5, "Username:", "")
	}
	username := ""
	// We only want to print the last 12 characters of the username
	if len(lw.loginForgotPasswordData.Username) > 12 {
		username = lw.loginForgotPasswordData.Username[len(lw.loginForgotPasswordData.Username)-12:]
	} else {
		username = lw.loginForgotPasswordData.Username
	}
	lw.PrintLn(lw.X+19, lw.Y+5, username, "")

	if lw.loginForgotPasswordState == LoginForgotPasswordDiscord {
		lw.PrintLn(lw.X+5, lw.Y+6, "Discord User:", lw.Terminal.Bold())
	} else {
		lw.PrintLn(lw.X+5, lw.Y+6, "Discord User:", "")
	}
	discordUser := ""
	// We only want the last 12 characters of the username
	if len(lw.loginForgotPasswordData.DiscordTag) > 12 {
		discordUser = lw.loginForgotPasswordData.DiscordTag[len(lw.loginForgotPasswordData.DiscordTag)-12:]
	} else {
		discordUser = lw.loginForgotPasswordData.DiscordTag
	}
	lw.PrintLn(lw.X+19, lw.Y+6, discordUser, "")

	// Draw the back and submit buttons
	if lw.loginForgotPasswordOptionSelected == 1 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", lw.Terminal.Bold())
	}

	if lw.loginForgotPasswordOptionSelected == 2 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+lw.Width-12, lw.Y+lw.Height, "<Submit>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+lw.Width-12, lw.Y+lw.Height, "<Submit>", lw.Terminal.Bold())
	}
}

func (lw *LoginWindow) drawLoginMenuForgotPasswordPending() {
	lw.loginForgotPasswordPendingMutex.Lock()
	defer lw.loginForgotPasswordPendingMutex.Unlock()

	lw.PrintLn(lw.X+5, lw.Y+3, "If an account exists with the details you provided, a validation code will be provided", util.RGBCode(255, 0, 0).FG())
	lw.PrintLn(lw.X+5, lw.Y+4, "to you via discord by Eden Bot, please enter it below to continue.", util.RGBCode(255, 0, 0).FG())

	if lw.loginForgotPasswordPendingState == LoginForgotPendingCode {
		lw.PrintLn(lw.X+5, lw.Y+6, "Validation Code:", lw.Terminal.Bold())
	} else {
		lw.PrintLn(lw.X+5, lw.Y+6, "Validation Code:", "")
	}
	validationCode := ""
	// We only want the last 12 characters of the username
	if len(lw.loginProcessForgotPasswordPendingData.Code) > 12 {
		validationCode = lw.loginProcessForgotPasswordPendingData.Code[len(lw.loginProcessForgotPasswordPendingData.Code)-12:]
	} else {
		validationCode = lw.loginProcessForgotPasswordPendingData.Code
	}
	lw.PrintLn(lw.X+22, lw.Y+6, validationCode, "")

	// Draw the back and submit buttons
	if lw.loginForgotPasswordPendingOptionSelected == 1 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", lw.Terminal.Bold())
	}

	if lw.loginForgotPasswordPendingOptionSelected == 2 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+lw.Width-12, lw.Y+lw.Height, "<Submit>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+lw.Width-12, lw.Y+lw.Height, "<Submit>", lw.Terminal.Bold())
	}
}

func (lw *LoginWindow) drawLoginMenuForgotPasswordSuccess() {
	lw.loginForgotPasswordSuccessMutex.Lock()
	defer lw.loginForgotPasswordSuccessMutex.Unlock()

	lw.PrintLn(lw.X+5, lw.Y+3, "Password reset validated, please enter a new password below to continue.", "")

	if lw.loginForgotPasswordSuccessState == LoginForgotPasswordSuccessEntry {
		lw.PrintLn(lw.X+9, lw.Y+6, "New Password:", lw.Terminal.Bold())
	} else {
		lw.PrintLn(lw.X+9, lw.Y+6, "New Password:", "")
	}
	password := ""
	// We only want the last 12 characters of the password
	if len(lw.loginForgotPasswordNewPasswordData.Password) > 12 {
		password = lw.loginForgotPasswordNewPasswordData.Password[len(lw.loginForgotPasswordNewPasswordData.Password)-12:]
	} else {
		password = lw.loginForgotPasswordNewPasswordData.Password
	}
	lw.PrintLn(lw.X+23, lw.Y+6, password, "")

	if lw.loginForgotPasswordSuccessState == LoginForgotPasswordSuccessConfirm {
		lw.PrintLn(lw.X+5, lw.Y+7, "Confirm Password:", lw.Terminal.Bold())
	} else {
		lw.PrintLn(lw.X+5, lw.Y+7, "Confirm Password:", "")
	}
	confirm := ""
	// We only want the last 12 characters of the password
	if len(lw.loginForgotPasswordNewPasswordData.PasswordConfirm) > 12 {
		confirm = lw.loginForgotPasswordNewPasswordData.PasswordConfirm[len(lw.loginForgotPasswordNewPasswordData.PasswordConfirm)-12:]
	} else {
		confirm = lw.loginForgotPasswordNewPasswordData.PasswordConfirm
	}
	lw.PrintLn(lw.X+23, lw.Y+7, confirm, "")

	// Draw the error message if there is one
	if lw.loginForgotPasswordNewPasswordData.Error != "" {
		lw.PrintLn(lw.X+5, lw.Y+10, lw.loginForgotPasswordNewPasswordData.Error, util.RGBCode(255, 0, 0).FG())
	}

	if lw.loginForgotPasswordSuccessOptionSelected == 2 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+lw.Width-12, lw.Y+lw.Height, "<Submit>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+lw.Width-12, lw.Y+lw.Height, "<Submit>", lw.Terminal.Bold())
	}
}

func (lw *LoginWindow) drawLoginMenuForgotPasswordFailed() {

	lw.PrintLn(lw.X+5, lw.Y+3, "The provided authorization code was invalid, please try again.", "")

	// Draw the back button
	if lw.loginForgotPasswordFailedOptionSelected == 1 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", lw.Terminal.Bold())
	}
}
