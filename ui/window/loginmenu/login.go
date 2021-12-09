package login

import "github.com/yamamushi/EscapingEden/ui/util"

type LoginState int

const (
	LoginNull LoginState = iota
	LoginUserInfo
	LoginPending
	LoginFailure
	LoginSuccess
)

// LoginMenuState is an enum for storing login state
type LoginMenuState int

const (
	LoginMenuEmail LoginMenuState = iota
	LoginMenuPassword
	LoginMenuNull
)

type LoginSubmitData struct {
	Email    string
	Password string
	Error    string
}

// drawLoginMenu draws the login window
func (lw *LoginWindow) drawLoginMenu() {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	//lw.SetContents("Login")

	switch lw.loginState {
	case LoginUserInfo:
		lw.drawLoginMenuUserInfo()
		return
	case LoginPending:
		lw.drawLoginMenuPending()
		return
	case LoginFailure:
		lw.drawLoginMenuFailure()
		return
	}

}

func (lw *LoginWindow) drawLoginMenuUserInfo() {
	lw.loginSubmitMutex.Lock()
	defer lw.loginSubmitMutex.Unlock()

	errorFG := util.RGBCode(255, 255, 255)
	errorBG := util.RGBCode(255, 0, 0)

	lw.PrintLn(lw.X+9, lw.Y+5, "Email:", "")
	email := ""
	// We only want the last 12 characters of the email
	if len(lw.loginSubmitData.Email) > 12 {
		email = lw.loginSubmitData.Email[len(lw.loginSubmitData.Email)-12:]
	} else {
		email = lw.loginSubmitData.Email
	}
	lw.PrintLn(lw.X+16, lw.Y+5, email, "")

	lw.PrintLn(lw.X+6, lw.Y+6, "Password:", "")
	for i := 0; i < len(lw.loginSubmitData.Password) && i < 12; i++ {
		lw.PrintLn(lw.X+16+i, lw.Y+6, "*", "")
	}

	//lw.loginSubmitData.Error = "This is a test error message"

	if lw.loginSubmitData.Error != "" {
		lw.PrintLnColor(lw.X+5, lw.Y+8, "Error logging in: "+lw.loginSubmitData.Error, errorFG.FG()+errorBG.BG())
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

func (lw *LoginWindow) drawLoginMenuFailure() {

}
