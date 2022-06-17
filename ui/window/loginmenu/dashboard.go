package login

// This is where we hit when we've finally logged in -> This should probably be refactored into something else but
// this was good enough to get by for now.
func (lw *LoginWindow) drawUserDashboard() {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	lw.PrintLn(lw.X+1, lw.Y+1, "Welcome "+lw.loginResponse.Account.Username+"!", "")

}
