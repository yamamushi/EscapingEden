package ui

import (
	"github.com/yamamushi/EscapingEden/messages"
)

// IsUserLoggedIn returns whether or not the user is logged in
func (c *Console) IsUserLoggedIn() bool {
	c.userLoggedInMutex.Lock()
	defer c.userLoggedInMutex.Unlock()
	return c.userLoggedIn
}

// LoginUser logs in the user, sets the userLoggedIn flag to true
func (c *Console) LoginUser(userInfo messages.UserInfo) {
	c.userLoggedInMutex.Lock()
	defer c.userLoggedInMutex.Unlock()
	c.UpdateUserInfo(userInfo)
	c.userLoggedIn = true
}

// LogoutUser logs out the user, sets the userLoggedIn flag to false and sets the active window to the login window
func (c *Console) LogoutUser() {
	c.userLoggedInMutex.Lock()
	defer c.userLoggedInMutex.Unlock()
	c.userLoggedIn = false
	c.UpdateUserInfo(messages.UserInfo{})
	c.SetActiveWindowNoThread(c.GetLoginWindow())
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
