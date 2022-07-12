package login

// This is where we hit when we've finally logged in -> This should probably be refactored into something else but
// this was good enough to get by for now.
func (lw *LoginWindow) drawUserDashboard() {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	lw.PrintLn(lw.X+10, lw.Y+5, "Welcome "+lw.GetUserInfoField("username")+"!", "")
	lw.PrintLn(lw.X+10, lw.Y+7, "Please select an option from the menu below to continue", "")

	lw.PrintLn(lw.X+11, lw.Y+9, "(", "")
	lw.PrintLn(lw.X+12, lw.Y+9, "a", lw.Terminal.Bold())
	lw.PrintLn(lw.X+13, lw.Y+9, ") Login as last character used (character name here)", "")

	lw.PrintLn(lw.X+11, lw.Y+10, "(", "")
	lw.PrintLn(lw.X+12, lw.Y+10, "b", lw.Terminal.Bold())
	lw.PrintLn(lw.X+13, lw.Y+10, ") Manage characters", "")

	lw.PrintLn(lw.X+11, lw.Y+11, "(", "")
	lw.PrintLn(lw.X+12, lw.Y+11, "c", lw.Terminal.Bold())
	lw.PrintLn(lw.X+13, lw.Y+11, ") Manage account settings", "")

	lw.PrintLn(lw.X+11, lw.Y+12, "(", "")
	lw.PrintLn(lw.X+12, lw.Y+12, "d", lw.Terminal.Bold())
	lw.PrintLn(lw.X+13, lw.Y+12, ") Logout", "")

}
