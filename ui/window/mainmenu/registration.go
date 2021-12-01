package login

import "github.com/yamamushi/EscapingEden/ui/util"

// drawRegistrationMenu draws the registration window
func (lw *LoginWindow) drawRegistrationMenu() {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	switch lw.registrationState {
	case RegistrationMain:
		lw.drawRegistrationWelcome()
	case RegistrationUserInfo:
		lw.drawRegistrationUserInfo()
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
	if lw.optionSelected == 1 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", "\033[1m")
	}
	if lw.optionSelected == 2 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+lw.Width-15, lw.Y+lw.Height, "<Continue>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+lw.Width-15, lw.Y+lw.Height, "<Continue>", "\033[1m")
	}
}

func (lw *LoginWindow) drawRegistrationUserInfo() {

	lw.PrintLn(lw.X+4, lw.Y+2, "Please enter your user registration information below.", "")
	lw.PrintLn(lw.X+4, lw.Y+3, "(You can use your arrow keys to navigate between fields)", "")

	lw.PrintLn(lw.X+12, lw.Y+7, "Username:", "\033[1m")
	fg := util.RGBCode(0, 255, 0)
	//bg := util.RGBCode(0, 255, 0)
	lw.PrintLn(lw.X+21, lw.Y+7, "         ", fg.FG()+"\033[4m")

	lw.PrintLn(lw.X+12, lw.Y+8, "Password:", "\033[1m")
	lw.PrintLn(lw.X+21, lw.Y+8, "         ", "\033[4m")

	lw.PrintLn(lw.X+4, lw.Y+9, "Confirm Password:", "\033[1m")
	lw.PrintLn(lw.X+21, lw.Y+9, "         ", "\033[4m")

	lw.PrintLn(lw.X+15, lw.Y+10, "Email:", "\033[1m")
	lw.PrintLn(lw.X+21, lw.Y+10, "         ", "\033[4m")

	lw.PrintLn(lw.X+13, lw.Y+11, "Discord:", "\033[1m")
	lw.PrintLn(lw.X+21, lw.Y+11, "         ", "\033[4m")

	lw.PrintLn(lw.X+13, lw.Y+13, "(Discord usernames accepted in the form of Username#0001)", "")

	if lw.optionSelected == 1 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+5, lw.Y+lw.Height, "<Back>", "\033[1m")
	}

	if lw.optionSelected == 2 {
		fg := util.RGBCode(0, 0, 0)
		bg := util.RGBCode(255, 255, 255)
		lw.PrintLn(lw.X+lw.Width-12, lw.Y+lw.Height, "<Submit>", fg.FG()+bg.BG())
	} else {
		lw.PrintLn(lw.X+lw.Width-12, lw.Y+lw.Height, "<Submit>", "\033[1m")
	}
}
