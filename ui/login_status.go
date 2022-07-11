package ui

// IsUserLoggedIn returns whether or not the user is logged in
func (c *Console) IsUserLoggedIn() bool {
	return c.userLoggedIn
}

// LoginUser logs in the user, sets the userLoggedIn flag to true
func (c *Console) LoginUser() {
	// No mutex because the caller has locked the mutex
	c.userLoggedIn = true
}

// LogoutUser logs out the user, sets the userLoggedIn flag to false
func (c *Console) LogoutUser() {
	c.userLoggedIn = false
}
