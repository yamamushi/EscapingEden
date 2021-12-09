package login

func (lw *LoginWindow) LoginSubmit() {
	lw.loginSubmitMutex.Lock()
	defer lw.loginSubmitMutex.Unlock()
	lw.loginSubmitData.Error = "" // clear error
	// check if username is empty
	if lw.loginSubmitData.Email == "" {
		lw.loginSubmitData.Error = "Email cannot be empty"
		return
	}
	// check if password is empty
	if lw.loginSubmitData.Password == "" {
		lw.loginSubmitData.Error = "Password cannot be empty"
		return
	}

	go lw.HandleReceiveChannel() // We're going to start listening for responses now

	return
}
