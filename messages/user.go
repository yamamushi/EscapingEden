package messages

import (
	"time"
)

// This struct is used to store some general information about an account after a successful login.
type UserInfo struct {
	ID                string
	Username          string
	CharacterName     string
	LastCharacterID   string
	LastCharacterName string
	DiscordTag        string
	LastLogin         time.Time
	LastLogout        time.Time
}

func (u *UserInfo) GetID() string {
	return u.ID
}

func (u *UserInfo) GetUsername() string {
	return u.Username
}

func (u *UserInfo) GetCharacter() string {
	return u.CharacterName
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

func (u *UserInfo) GetLastCharacterID() string {
	return u.LastCharacterID
}

func (u *UserInfo) GetLastCharacterName() string {
	return u.LastCharacterName
}
