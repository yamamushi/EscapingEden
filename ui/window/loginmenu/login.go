package login

import "github.com/yamamushi/EscapingEden/ui/util"

type LoginState int

const (
	LoginNull LoginState = iota
	LoginUserInfo
	LoginPending
	LoginForgotPassword
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
	LoginForgotPasswordNull
)

type LoginForgotPasswordData struct {
	Username string
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
			// If we didn't get an error, we load our user info screen
			lw.windowState = LoginWindowUserDashboard
		}
		lw.loginResponseReceived = false
		lw.RequestFlushFromConsole()
	}

	lw.PrintLn(lw.X+lw.Width/2-5, lw.Y+lw.Height/2, "Login Pending...", lw.Terminal.Bold())
}

func (lw *LoginWindow) drawLoginMenuForgotPassword() {
	lw.loginForgotPasswordMutex.Lock()
	defer lw.loginForgotPasswordMutex.Unlock()

	if lw.loginForgotPasswordState == LoginForgotPasswordUsername {
		lw.PrintLn(lw.X+9, lw.Y+5, "Username:", lw.Terminal.Bold())
	} else {
		lw.PrintLn(lw.X+9, lw.Y+5, "Username:", "")
	}
	username := ""
	// We only want the last 12 characters of the username
	if len(lw.loginForgotPasswordData.Username) > 12 {
		username = lw.loginForgotPasswordData.Username[len(lw.loginForgotPasswordData.Username)-12:]
	} else {
		username = lw.loginForgotPasswordData.Username
	}
	lw.PrintLn(lw.X+16, lw.Y+5, username, "")

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
