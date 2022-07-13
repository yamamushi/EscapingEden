package dashboard

// This is where we hit when we've finally logged in -> This should probably be refactored into something else but
// this was good enough to get by for now.
func (dw *DashboardWindow) drawMenu() {
	dw.dwMutex.Lock()
	defer dw.dwMutex.Unlock()

	lastCharacter := dw.GetUserInfoField("lastcharactername")

	if !(dw.dwInitialized) {
		dw.RequestFlushFromConsole()
		if lastCharacter == "" {
			dw.firstTimeLogin = true
		} else {
			dw.lastCharacterID = dw.GetUserInfoField("lastcharacterid")
			dw.lastCharacterName = dw.GetUserInfoField("lastcharactername")
		}
		dw.dwInitialized = true
	}

	dw.PrintLn(dw.X+10, dw.Y+5, "Welcome "+dw.GetUserInfoField("username")+"!", "")
	dw.PrintLn(dw.X+10, dw.Y+7, "Please select an option from the menu below to continue", "")

	dw.PrintLn(dw.X+11, dw.Y+9, "(", "")
	dw.PrintLn(dw.X+12, dw.Y+9, "a", dw.Terminal.Bold())

	offset := 0
	if dw.firstTimeLogin {
		dw.PrintLn(dw.X+13, dw.Y+9, ") Create a new character ", "")
	} else {
		dw.PrintLn(dw.X+13, dw.Y+9, ") Login as last character used ("+dw.lastCharacterName+")", "")
		dw.PrintLn(dw.X+11, dw.Y+10, "(", "")
		dw.PrintLn(dw.X+12, dw.Y+10, "b", dw.Terminal.Bold())
		dw.PrintLn(dw.X+13, dw.Y+10, ") Manage characters", "")
		offset = 1
	}

	dw.PrintLn(dw.X+11, dw.Y+10+offset, "(", "")
	dw.PrintLn(dw.X+12, dw.Y+10+offset, "c", dw.Terminal.Bold())
	dw.PrintLn(dw.X+13, dw.Y+10+offset, ") Manage account settings", "")

	dw.PrintLn(dw.X+11, dw.Y+11+offset, "(", "")
	dw.PrintLn(dw.X+12, dw.Y+11+offset, "d", dw.Terminal.Bold())
	dw.PrintLn(dw.X+13, dw.Y+11+offset, ") Logout", "")

}
