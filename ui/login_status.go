package ui

// IsUserLoggedIn returns whether or not the user is logged in
func (c *Console) IsUserLoggedIn() bool {
	c.userLoggedInMutex.Lock()
	defer c.userLoggedInMutex.Unlock()
	return c.userLoggedIn
}

// LoginUser logs in the user, sets the userLoggedIn flag to true
func (c *Console) LoginUser() {
	c.userLoggedInMutex.Lock()
	defer c.userLoggedInMutex.Unlock()
	c.userLoggedIn = true
}

// LogoutUser logs out the user, sets the userLoggedIn flag to false
func (c *Console) LogoutUser() {
	c.userLoggedInMutex.Lock()
	defer c.userLoggedInMutex.Unlock()
	c.userLoggedIn = false
}

func (c *Console) IsCharacterLoggedIn() bool {
	c.characterLoggedInMutex.Lock()
	defer c.characterLoggedInMutex.Unlock()
	return c.characterLoggedIn
}

// LoginCharacter
func (c *Console) LoginCharacter() {
	c.characterLoggedInMutex.Lock()
	defer c.characterLoggedInMutex.Unlock()
	c.characterLoggedIn = true
}

func (c *Console) LogoutCharacter() {
	c.characterLoggedInMutex.Lock()
	defer c.characterLoggedInMutex.Unlock()
	c.characterLoggedIn = false
}
