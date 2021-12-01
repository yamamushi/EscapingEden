package login

// drawLoginMenu draws the login window
func (lw *LoginWindow) drawLoginMenu() {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	lw.SetContents("Login")
}
