package login

func (lw *LoginWindow) drawUserDashboard() {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	lw.PrintLn(lw.X+1, lw.Y+1, "Welcome "+lw.loginResponse.Account.Username+"!", "")

}
