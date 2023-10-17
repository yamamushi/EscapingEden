package ui

import (
	"github.com/yamamushi/EscapingEden/messages"
	"time"
)

func (c *Console) UpdateUserInfo(userInfo messages.UserInfo) {
	c.userInfoMutex.Lock()
	defer c.userInfoMutex.Unlock()
	c.UserInfo = userInfo
	//c.Log.Println(logging.LogInfo, "Updated user info")
	c.UpdateWindowsUserInfo()
	//c.Log.Println(logging.LogInfo, "Updated windows user info")
}

func (c *Console) UpdateWindowsUserInfo() {
	for _, w := range c.Windows {
		w.SetUserInfo(c.UserInfo)
	}
}

func (c *Console) GetUserName() string {
	c.userInfoMutex.Lock()
	defer c.userInfoMutex.Unlock()
	return c.UserInfo.Username
}

func (c *Console) GetDiscordTag() string {
	c.userInfoMutex.Lock()
	defer c.userInfoMutex.Unlock()
	return c.UserInfo.DiscordTag
}

func (c *Console) GetLastLogin() time.Time {
	c.userInfoMutex.Lock()
	defer c.userInfoMutex.Unlock()
	return c.UserInfo.LastLogin
}

func (c *Console) GetLastLogout() time.Time {
	c.userInfoMutex.Lock()
	defer c.userInfoMutex.Unlock()
	return c.UserInfo.LastLogout
}
