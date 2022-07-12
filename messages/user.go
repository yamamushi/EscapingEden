package messages

import "time"

// This struct is used to store some general information about an account after a successful login.
type UserInfo struct {
	Username   string
	Character  string
	DiscordTag string
	LastLogin  time.Time
	LastLogout time.Time
}
