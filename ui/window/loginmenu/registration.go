package login

import "github.com/yamamushi/EscapingEden/ui/util"

// RegistrationState is an enum for storing registration state
type RegistrationState int

const (
	RegistrationMain RegistrationState = iota
	RegistrationUserInfo
	RegistrationSuccess
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
	case RegistrationSuccess:
		lw.drawRegistrationSuccess()
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
	lw.PrintChar(lw.X+43, lw.Y+5, "r", "\033[1m")
	lw.PrintLn(lw.X+45, lw.Y+8, "ctrl-r", "\033[1m")
	lw.PrintLn(lw.X+20, lw.Y+9, "ctrl-h", "\033[1m")
	lw.PrintChar(lw.X+20, lw.Y+7, "b", "\033[1m")
	lw.PrintChar(lw.X+2, lw.Y+12, "d", "\033[1m")
	lw.SetContents(content)

	// We eventually want to embed all of this in an easier to use way
	lw.PrintLn(lw.X+1, lw.Y+lw.Height-2, "When you are ready, and have agreed to the [r]ules, please select <Continue> below.", "")
	lw.PrintLn(lw.X+67, lw.Y+lw.Height-2, "<Continue>", "\033[1m")
	lw.PrintChar(lw.X+45, lw.Y+lw.Height-2, "r", "\033[1m")
	// Bold the text for the back and continue buttons
	if lw.registrationNavOptionSelected == 1 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", "\033[1m")
	}
	if lw.registrationNavOptionSelected == 2 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+lw.Width-15, lw.Y+lw.Height, "<Continue>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+lw.Width-15, lw.Y+lw.Height, "<Continue>", "\033[1m")
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
		lw.PrintLn(lw.X+12, lw.Y+7, "Username:", "\033[1m")
	} else {
		lw.PrintLn(lw.X+12, lw.Y+7, "Username:", "")
	}
	lw.PrintLn(lw.X+22, lw.Y+7, lw.registrationSubmitData.Username, "")
	lw.PrintLn(lw.X+36, lw.Y+7, lw.registrationErrorData.UsernameError(), errorFG.FG()+errorBG.BG())

	if lw.registrationUserInfoOptionSelected == UserInfoPassword {
		lw.PrintLn(lw.X+12, lw.Y+8, "Password:", "\033[1m")
	} else {
		lw.PrintLn(lw.X+12, lw.Y+8, "Password:", "")
	}
	lw.PrintLn(lw.X+22, lw.Y+8, lw.registrationSubmitData.Password, "")
	lw.PrintLn(lw.X+36, lw.Y+8, lw.registrationErrorData.PasswordError(), errorFG.FG()+errorBG.BG())

	if lw.registrationUserInfoOptionSelected == UserInfoPasswordConfirm {
		lw.PrintLn(lw.X+4, lw.Y+9, "Confirm Password:", "\033[1m")
	} else {
		lw.PrintLn(lw.X+4, lw.Y+9, "Confirm Password:", "")
	}
	lw.PrintLn(lw.X+22, lw.Y+9, lw.registrationSubmitData.PasswordConfirm, "")
	lw.PrintLn(lw.X+36, lw.Y+9, lw.registrationErrorData.PasswordConfirmError(), errorFG.FG()+errorBG.BG())

	if lw.registrationUserInfoOptionSelected == UserInfoEmail {
		lw.PrintLn(lw.X+15, lw.Y+10, "Email:", "\033[1m")
	} else {
		lw.PrintLn(lw.X+15, lw.Y+10, "Email:", "")
	}
	lw.PrintLn(lw.X+22, lw.Y+10, lw.registrationSubmitData.Email, "")
	lw.PrintLn(lw.X+36, lw.Y+10, lw.registrationErrorData.EmailError(), errorFG.FG()+errorBG.BG())

	lw.PrintLn(lw.X+20, lw.Y+14, "Do you agree to the rules?     (Space to toggle)", "")
	lw.PrintLn(lw.X+20, lw.Y+13, lw.registrationErrorData.RulesError(), errorFG.FG()+errorBG.BG())

	if lw.registrationUserInfoOptionSelected == UserInfoAgreeRules {
		lw.PrintLn(lw.X+47, lw.Y+14, "[ ]", "\033[1m")
	} else {
		lw.PrintLn(lw.X+47, lw.Y+14, "[ ]", "")
	}
	if lw.registrationAgreeRules {
		lw.PrintChar(lw.X+48, lw.Y+14, "\u2666", "\033[1m")
	}

	if lw.registrationNavOptionSelected == 1 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", "\033[1m")
	}

	if lw.registrationNavOptionSelected == 2 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+lw.Width-12, lw.Y+lw.Height, "<Submit>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+lw.Width-12, lw.Y+lw.Height, "<Submit>", "\033[1m")
	}
}

func (lw *LoginWindow) drawRegistrationSuccess() {
	lw.registrationStatusMutex.Lock()
	defer lw.registrationStatusMutex.Unlock()
	if lw.registrationResponse.Success {
		lw.PrintLn(lw.X+2, lw.Y+4, "Registration successful!", "")
	} else {
		lw.PrintLn(lw.X+2, lw.Y+4, "Something went wrong!: "+lw.registrationResponse.Message, "")
	}
}
