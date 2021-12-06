package login

import (
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/ui/util"
)

// RegistrationState is an enum for storing registration state
type RegistrationState int

const (
	RegistrationMain RegistrationState = iota
	RegistrationUserInfo
	RegistrationPending
	RegistrationSuccess
	RegistrationFailure
)

type RegistrationUserInfoState int

const (
	UserInfoUsername RegistrationUserInfoState = iota
	UserInfoPassword
	UserInfoPasswordConfirm
	UserInfoEmail
	UserInfoAgreeRules
	UserInfoNULL
)

// drawRegistrationMenu draws the registration window
func (lw *LoginWindow) drawRegistrationMenu() {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	switch lw.registrationState {
	case RegistrationMain:
		lw.drawRegistrationWelcome()
	case RegistrationUserInfo:
		lw.drawRegistrationUserInfo()
	case RegistrationPending:
		lw.drawRegistrationStatus()
	case RegistrationSuccess:
		lw.drawRegistrationSuccess()
	case RegistrationFailure:
		lw.drawRegistrationFailure()
	}

}

func (lw *LoginWindow) drawRegistrationWelcome() {
	content, err := util.OpenFileAsText("assets/text/welcome.txt")
	if err != nil {
		lw.Error("Error opening rules file:" + err.Error())
		return
	}
	// This isn't pretty but it works
	// Perhaps in the future we can have embedded text file reading
	// Technically we could pull this off with the art reader too
	// But that's a bit more overkill for this
	lw.PrintChar(lw.X+43, lw.Y+5, "r", lw.Terminal.Bold())
	lw.PrintLn(lw.X+45, lw.Y+8, "ctrl-r", lw.Terminal.Bold())
	lw.PrintLn(lw.X+20, lw.Y+9, "ctrl-h", lw.Terminal.Bold())
	lw.PrintChar(lw.X+20, lw.Y+7, "b", lw.Terminal.Bold())
	lw.PrintChar(lw.X+2, lw.Y+12, "d", lw.Terminal.Bold())
	lw.SetContents(content)

	// We eventually want to embed all of this in an easier to use way
	lw.PrintLn(lw.X+1, lw.Y+lw.Height-2, "When you are ready, and have agreed to the [r]ules, please select <Continue> below.", "")
	lw.PrintLn(lw.X+67, lw.Y+lw.Height-2, "<Continue>", lw.Terminal.Bold())
	lw.PrintChar(lw.X+45, lw.Y+lw.Height-2, "r", lw.Terminal.Bold())
	// Bold the text for the back and continue buttons
	if lw.registrationNavOptionSelected == 1 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", lw.Terminal.Bold())
	}
	if lw.registrationNavOptionSelected == 2 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+lw.Width-15, lw.Y+lw.Height, "<Continue>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+lw.Width-15, lw.Y+lw.Height, "<Continue>", lw.Terminal.Bold())
	}
}

func (lw *LoginWindow) drawRegistrationUserInfo() {
	lw.registrationSubmitMutex.Lock()
	defer lw.registrationSubmitMutex.Unlock()
	// If we're touching lw.registrationErrorData, we need to lock it

	lw.PrintLn(lw.X+4, lw.Y+4, "Please enter your user registration information below.", "")
	lw.PrintLn(lw.X+4, lw.Y+5, "(You can use your arrow keys, or enter, to navigate between fields)", "")

	errorFG := util.RGBCode(255, 255, 255)
	errorBG := util.RGBCode(255, 0, 0)

	if lw.registrationUserInfoOptionSelected == UserInfoUsername {
		lw.PrintLn(lw.X+12, lw.Y+7, "Username:", lw.Terminal.Bold())
	} else {
		lw.PrintLn(lw.X+12, lw.Y+7, "Username:", "")
	}

	username := ""
	// We only want the last 12 characters of the username
	if len(lw.registrationSubmitData.Username) >= 12 {
		username = lw.registrationSubmitData.Username[len(lw.registrationSubmitData.Username)-12:]
	} else {
		username = lw.registrationSubmitData.Username
	}

	lw.PrintLn(lw.X+22, lw.Y+7, username, "")
	lw.PrintLnColor(lw.X+41, lw.Y+7, lw.registrationErrorData.UsernameError(), errorFG.FG()+errorBG.BG())

	if lw.registrationUserInfoOptionSelected == UserInfoPassword {
		lw.PrintLn(lw.X+12, lw.Y+8, "Password:", lw.Terminal.Bold())
	} else {
		lw.PrintLn(lw.X+12, lw.Y+8, "Password:", "")
	}

	password := ""
	// We only want the last 12 characters of the password
	if len(lw.registrationSubmitData.Password) >= 12 {
		password = lw.registrationSubmitData.Password[len(lw.registrationSubmitData.Password)-12:]
	} else {
		password = lw.registrationSubmitData.Password
	}
	for i := 0; i < len(password); i++ {
		lw.PrintChar(lw.X+22+i, lw.Y+8, "*", lw.Terminal.Bold())
	}
	lw.PrintLnColor(lw.X+41, lw.Y+8, lw.registrationErrorData.PasswordError(), errorFG.FG()+errorBG.BG())

	if lw.registrationUserInfoOptionSelected == UserInfoPasswordConfirm {
		lw.PrintLn(lw.X+4, lw.Y+9, "Confirm Password:", lw.Terminal.Bold())
	} else {
		lw.PrintLn(lw.X+4, lw.Y+9, "Confirm Password:", "")
	}

	passwordConfirm := ""
	// We only want the last 12 characters of the password
	if len(lw.registrationSubmitData.PasswordConfirm) > 12 {
		passwordConfirm = lw.registrationSubmitData.PasswordConfirm[len(lw.registrationSubmitData.Password)-12:]
	} else {
		passwordConfirm = lw.registrationSubmitData.PasswordConfirm
	}
	for i := 0; i < len(passwordConfirm) && i < 12; i++ {
		lw.PrintLn(lw.X+22+i, lw.Y+9, "*", "")
	}
	lw.PrintLnColor(lw.X+41, lw.Y+9, lw.registrationErrorData.PasswordConfirmError(), errorFG.FG()+errorBG.BG())

	if lw.registrationUserInfoOptionSelected == UserInfoEmail {
		lw.PrintLn(lw.X+15, lw.Y+10, "Email:", lw.Terminal.Bold())
	} else {
		lw.PrintLn(lw.X+15, lw.Y+10, "Email:", "")
	}

	email := ""
	// We only want the last 12 characters of the email
	if len(lw.registrationSubmitData.Email) > 12 {
		email = lw.registrationSubmitData.Email[len(lw.registrationSubmitData.Email)-12:]
	} else {
		email = lw.registrationSubmitData.Email
	}
	lw.PrintLn(lw.X+22, lw.Y+10, email, "")
	lw.PrintLnColor(lw.X+41, lw.Y+10, lw.registrationErrorData.EmailError(), errorFG.FG()+errorBG.BG())

	lw.PrintLn(lw.X+20, lw.Y+14, "Do you agree to the rules?     (Space to toggle)", "")
	lw.PrintLnColor(lw.X+20, lw.Y+13, lw.registrationErrorData.RulesError(), errorFG.FG()+errorBG.BG())

	if lw.registrationUserInfoOptionSelected == UserInfoAgreeRules {
		lw.PrintLn(lw.X+47, lw.Y+14, "[ ]", lw.Terminal.Bold())
	} else {
		lw.PrintLn(lw.X+47, lw.Y+14, "[ ]", "")
	}
	if lw.registrationAgreeRules {
		lw.PrintChar(lw.X+48, lw.Y+14, "\u2666", lw.Terminal.Bold())
	}

	if lw.registrationNavOptionSelected == 1 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", lw.Terminal.Bold())
	}

	if lw.registrationNavOptionSelected == 2 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+lw.Width-12, lw.Y+lw.Height, "<Submit>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+lw.Width-12, lw.Y+lw.Height, "<Submit>", lw.Terminal.Bold())
	}
}

func (lw *LoginWindow) drawRegistrationStatus() {
	lw.registrationStatusMutex.Lock()
	defer lw.registrationStatusMutex.Unlock()

	if lw.registrationResponseReceived {
		if lw.registrationResponse.Error == messages.AMError_Null {
			lw.registrationState = RegistrationSuccess
		} else {
			lw.registrationState = RegistrationFailure
		}
		lw.registrationResponseReceived = false
		lw.RequestFlushFromConsole()
	} else {
		lw.PrintLn(lw.X+2, lw.Y+4, "Please wait while your account registration is being processed...", "")
	}
}

func (lw *LoginWindow) drawRegistrationFailure() {
	lw.registrationErrorMutex.Lock()
	defer lw.registrationErrorMutex.Unlock()
	// We're locking this because we want to parse the error into lw.registrationErrorData

	// Print error in red
	lw.PrintLnColor(lw.X+2, lw.Y+2, "Something went wrong!: "+lw.registrationResponse.Error.Error(), "\033[31m")

	if lw.registrationErrorData.errorRequest != "" {
		lw.PrintLnColor(lw.X+2, lw.Y+4, "Please report this issue as something more serious may be wrong", "\033[31m")
	}

	lw.PrintLn(lw.X+2, lw.Y+lw.Height-2, "Please select <Back> to be taken back to the registration details screen.", "")
	lw.PrintLn(lw.X+16, lw.Y+lw.Height-2, "<Back>", lw.Terminal.Bold())

	if lw.registrationFailureOptionSelected == 1 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", lw.Terminal.Bold())
	}
}

func (lw *LoginWindow) drawRegistrationSuccess() {
	lw.PrintLnColor(lw.X+2, lw.Y+2, "Registration successful!", "\033[32m")

	lw.PrintLn(lw.X+2, lw.Y+6, "We look forward to seeing you soon in Eden!", "")

	lw.PrintLn(lw.X+2, lw.Y+lw.Height-2, "Please select <Continue> to be taken back to the login screen.", "")
	lw.PrintLn(lw.X+16, lw.Y+lw.Height-2, "<Continue>", lw.Terminal.Bold())

	if lw.registrationSuccessOptionSelected == 2 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+lw.Width-12, lw.Y+lw.Height, "<Continue>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+lw.Width-12, lw.Y+lw.Height, "<Continue>", lw.Terminal.Bold())
	}
}
