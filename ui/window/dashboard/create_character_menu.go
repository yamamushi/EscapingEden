package dashboard

import "github.com/yamamushi/EscapingEden/ui/util"

func (dw *DashboardWindow) drawCreateCharacterMenu() {
	switch dw.characterCreatorState {
	case CharacterCreatorFirstTimeLoginWelcome:
		dw.drawFirstTimeLogin()
	case CharacterCreatorCharacterDetails:
		dw.drawCharacterCreator()
	case CharacterCreatorConfirmCharacter:
		dw.drawConfirmCharacter()
	}
}

func (dw *DashboardWindow) drawFirstTimeLogin() {
	//dw.Log.Println(logging.LogInfo, "Drawing first time login window")
	welcomeHeader := "Welcome traveller, to the world of Eden."
	dw.PrintLn(dw.X+(dw.Width/2)-(len(welcomeHeader)/2), dw.Y+2, welcomeHeader, "")

	welcomeMessageP1 := "This is the first time you are logging in to the game, so a character will need to be created for you."
	welcomeMessageP2 := "As this is still in early development, character creation is limited in scope. "
	welcomeMessageP3 := "You will be asked to choose a character color, as well as a username."
	welcomeMessageP4 := "There will be frequent wipes of this server, so please be patient."
	welcomeMessageP5 := "You may press Enter to continue."
	continueButton := "<Continue>"

	dw.PrintLn(dw.X+2, dw.Y+6, welcomeMessageP1, "")
	dw.PrintLn(dw.X+2, dw.Y+8, welcomeMessageP2, "")
	dw.PrintLn(dw.X+2, dw.Y+10, welcomeMessageP3, "")
	dw.PrintLn(dw.X+2, dw.Y+12, welcomeMessageP4, "")

	dw.PrintLn(dw.X+(dw.Width/2)-(len(welcomeMessageP5)/2), dw.Y+dw.Height-4, welcomeMessageP5, "")

	// Draw continue button centered
	dw.PrintLn(dw.X+(dw.Width/2)-(len(continueButton)/2), dw.Y+dw.Height-2, continueButton, dw.Terminal.Bold())

}

func (dw *DashboardWindow) drawCharacterCreator() {
	welcomeHeader := "Character Creation"
	dw.PrintLn(dw.X+(dw.Width/2)-(len(welcomeHeader)/2), dw.Y+2, welcomeHeader, "")

	usernameField := "Name: "
	if dw.charCreatorOptionSelected == 1 {
		dw.PrintLn(dw.X+2, dw.Y+4, usernameField, dw.Terminal.Bold())
	} else {
		dw.PrintLn(dw.X+2, dw.Y+4, usernameField, "")
	}
	dw.PrintLn(dw.X+2+len(usernameField), dw.Y+4, dw.charCreatorName, "")

	if dw.charCreatorUsernameError != "" {
		colorRed := util.ColorCode{255, 0, 0}
		dw.PrintLn(dw.X+dw.Width-2-len(dw.charCreatorUsernameError), dw.Y+4, dw.charCreatorUsernameError, colorRed.FG())
	}

	colorField := "Color -> (Press right arrow to select, left to deselect) "
	if dw.charCreatorOptionSelected == 2 {
		if dw.charColorOptionActive {
			greenColor := util.ColorCode{0, 255, 0}
			dw.PrintLn(dw.X+2, dw.Y+6, colorField, dw.Terminal.Bold()+greenColor.FG())
		} else {
			dw.PrintLn(dw.X+2, dw.Y+6, colorField, dw.Terminal.Bold())
		}
	} else {
		dw.PrintLn(dw.X+2, dw.Y+6, colorField, "")
	}

	greenOption := "( ) Green"
	redOption := "( ) Red"
	blueOption := "( ) Blue"

	colorRed := util.ColorCode{255, 0, 0}
	dw.PrintLn(dw.X+4, dw.Y+8, redOption, colorRed.FG())
	colorBlue := util.ColorCode{0, 0, 255}
	dw.PrintLn(dw.X+4, dw.Y+9, blueOption, colorBlue.FG())
	colorGreen := util.ColorCode{0, 255, 0}
	dw.PrintLn(dw.X+4, dw.Y+10, greenOption, colorGreen.FG())

	dw.PrintLn(dw.X+5, dw.Y+8+dw.charColorOption, "*", dw.Terminal.Bold())

	footerMessage := "When you are ready, select submit to create your character."
	dw.PrintLn(dw.X+(dw.Width/2)-(len(footerMessage)/2), dw.Y+dw.Height-4, footerMessage, "")

	submitOption := "<Submit>"
	cancelOption := "<Cancel>"
	if dw.charCreatorNavOptionSelected == 1 {
		dw.PrintLn(dw.X+2, dw.Y+dw.Height-1, cancelOption, dw.Terminal.Bold())
	} else {
		dw.PrintLn(dw.X+2, dw.Y+dw.Height-1, cancelOption, "")
	}

	if dw.charCreatorNavOptionSelected == 2 {
		dw.PrintLn(dw.X+dw.Width-2-len(submitOption), dw.Y+dw.Height-1, submitOption, dw.Terminal.Bold())
	} else {
		dw.PrintLn(dw.X+dw.Width-2-len(submitOption), dw.Y+dw.Height-1, submitOption, "")
	}

}

func (dw *DashboardWindow) drawConfirmCharacter() {
	dw.PrintLn(dw.X+(dw.Width/2)-(len("Confirm Character")/2), dw.Y+2, "Confirm Character", "")

	colorRed := util.ColorCode{255, 0, 0}
	colorBlue := util.ColorCode{0, 0, 255}
	colorGreen := util.ColorCode{0, 255, 0}
	color := util.ColorCode{}

	if dw.charColorOption == 0 {
		color = colorRed
	} else if dw.charColorOption == 1 {
		color = colorBlue
	} else if dw.charColorOption == 2 {
		color = colorGreen
	}

	dw.PrintLn(dw.X+2, dw.Y+4, "Name: "+dw.charCreatorName, "")

	colorText := "Color:"
	dw.PrintLn(dw.X+2, dw.Y+6, "Color:", "")
	dw.PrintLn(dw.X+3+len(colorText), dw.Y+6, "@", color.FG())

	footerMessage := "If everything looks correct, select confirm to create your character."
	dw.PrintLn(dw.X+(dw.Width/2)-(len(footerMessage)/2), dw.Y+dw.Height-4, footerMessage, "")

	submitOption := "<Confirm>"
	cancelOption := "<Cancel>"
	if dw.charCreatorConfirmNavOptionSelected == 1 {
		dw.PrintLn(dw.X+2, dw.Y+dw.Height-1, cancelOption, dw.Terminal.Bold())
	} else {
		dw.PrintLn(dw.X+2, dw.Y+dw.Height-1, cancelOption, "")
	}

	if dw.charCreatorConfirmNavOptionSelected == 2 {
		dw.PrintLn(dw.X+dw.Width-2-len(submitOption), dw.Y+dw.Height-1, submitOption, dw.Terminal.Bold())
	} else {
		dw.PrintLn(dw.X+dw.Width-2-len(submitOption), dw.Y+dw.Height-1, submitOption, "")
	}
}
