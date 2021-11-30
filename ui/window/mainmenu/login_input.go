package login

// handleLoginInput handles input for the login window
func (lw *LoginWindow) handleLoginInput(input string) {
	lw.lwMutex.Lock()
	defer lw.lwMutex.Unlock()

	if !lw.GetActive() {
		return
	}

	//input = strings.ToLower(input[:1])

	switch lw.loginState {
	case LoginUsername:
		lw.credentials.Username = input
		lw.loginState = LoginPassword
	case LoginPassword:
		lw.credentials.Hash = input
		lw.loginState = LoginSubmit
	case LoginSubmit:
		lw.ConsoleSend <- "login:" + lw.credentials.Username + ":" + lw.credentials.Hash
		lw.loginState = LoginUsername
	}
}
