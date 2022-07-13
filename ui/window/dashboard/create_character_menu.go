package dashboard

func (dw *DashboardWindow) drawCreateCharacterMenu() {
	if dw.firstTimeLogin {
		dw.drawFirstTimeLogin()
	} else {
		dw.drawCharacterCreator()
	}
}

func (dw *DashboardWindow) drawFirstTimeLogin() {
	//dw.Log.Println(logging.LogInfo, "Drawing first time login window")
	welcomeHeader := "Welcome traveller, to the world of Eden."
	dw.PrintLn(dw.X+(dw.Width/2)-(len(welcomeHeader)/2), dw.Y+2, welcomeHeader, "")
}

func (dw *DashboardWindow) drawCharacterCreator() {
	//dw.Log.Println(logging.LogInfo, "Drawing create character window")
	dw.PrintLn(dw.X+(dw.Width/2), dw.Y+2, "Character Creator", "")
}
