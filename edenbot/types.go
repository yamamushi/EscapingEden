package edenbot

// DiscordUser represents a Discord user
type DiscordUser struct {
	ID       string
	Username string
	Tag      string
}

// EdenBotInterface defines the interface for EdenBot operations
type EdenBotInterface interface {
	IsUserInDiscordServer(discordTag string) (*DiscordUser, error)
	ValidateUser(username, discordID string) error
}
