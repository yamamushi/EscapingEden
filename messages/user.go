package messages

import (
	"time"
)

// This struct is used to store some general information about an account after a successful login.
type UserInfo struct {
	Username   string
	Character  string
	DiscordTag string
	LastLogin  time.Time
	LastLogout time.Time
}

func (u *UserInfo) GetUsername() string {
	return u.Username
}

func (u *UserInfo) GetCharacter() string {
	return u.Character
}

func (u *UserInfo) GetDiscordTag() string {
	return u.DiscordTag
}

func (u *UserInfo) GetLastLogin() time.Time {
	return u.LastLogin
}

func (u *UserInfo) GetLastLogout() time.Time {
	return u.LastLogout
}
