package dashboard

// This is where we hit when we've finally logged in -> This should probably be refactored into something else but
// this was good enough to get by for now.
func (dw *DashboardWindow) drawMenu() {
	dw.dwMutex.Lock()
	defer dw.dwMutex.Unlock()
	if !(dw.dwInitialized) {
		dw.RequestFlushFromConsole()
		dw.dwInitialized = true
	}

	dw.PrintLn(dw.X+10, dw.Y+5, "Welcome "+dw.GetUserInfoField("username")+"!", "")
	dw.PrintLn(dw.X+10, dw.Y+7, "Please select an option from the menu below to continue", "")

	dw.PrintLn(dw.X+11, dw.Y+9, "(", "")
	dw.PrintLn(dw.X+12, dw.Y+9, "a", dw.Terminal.Bold())
	dw.PrintLn(dw.X+13, dw.Y+9, ") Login as last character used (character name here)", "")

	dw.PrintLn(dw.X+11, dw.Y+10, "(", "")
	dw.PrintLn(dw.X+12, dw.Y+10, "b", dw.Terminal.Bold())
	dw.PrintLn(dw.X+13, dw.Y+10, ") Manage characters", "")

	dw.PrintLn(dw.X+11, dw.Y+11, "(", "")
	dw.PrintLn(dw.X+12, dw.Y+11, "c", dw.Terminal.Bold())
	dw.PrintLn(dw.X+13, dw.Y+11, ") Manage account settings", "")

	dw.PrintLn(dw.X+11, dw.Y+12, "(", "")
	dw.PrintLn(dw.X+12, dw.Y+12, "d", dw.Terminal.Bold())
	dw.PrintLn(dw.X+13, dw.Y+12, ") Logout", "")

}
