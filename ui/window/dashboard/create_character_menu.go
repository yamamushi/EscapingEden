package dashboard

import "github.com/yamamushi/EscapingEden/ui/util"

func (dw *DashboardWindow) drawCreateCharacterMenu() {
	if dw.firstTimeLogin {
		dw.characterCreatorState = CharacterCreatorFirstTimeLoginWelcome
		dw.drawFirstTimeLogin()
	} else {
		dw.characterCreatorState = CharacterCreatorCharacterDetails
		dw.drawCharacterCreator()
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
	dw.PrintLn(dw.X+2, dw.Y+4, usernameField, "")

	colorField := "Color: "
	dw.PrintLn(dw.X+2, dw.Y+6, colorField, "")
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

	footerMessage := "When you are ready, press Enter to create your character."
	dw.PrintLn(dw.X+(dw.Width/2)-(len(footerMessage)/2), dw.Y+dw.Height-4, footerMessage, "")

	submitOption := "<Submit>"
	cancelOption := "<Cancel>"
	dw.PrintLn(dw.X+2, dw.Y+dw.Height-1, cancelOption, "")
	dw.PrintLn(dw.X+dw.Width-2-len(submitOption), dw.Y+dw.Height-1, submitOption, "")

}
