package edenconfig

import (
	"github.com/BurntSushi/toml"
)

// Config is the main configuration struct
type Config struct {
	Logger   LoggerConfig   `toml:"logging"`
	Server   ServerConfig   `toml:"server"`
	DB       DatabaseConfig `toml:"database"`
	Discord  DiscordConfig  `toml:"discord"`
	WorldGen WorldGenConfig `toml:"worldgen"`
}

// LoggerConfig holds logging configuration
type LoggerConfig struct {
	Type       string `toml:"type"`
	Path       string `toml:"path"`
	DebugChunk bool   `toml:"debug_chunk"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host            string `toml:"host"`
	Port            string `toml:"port"`
	ShutdownTimeout int    `toml:"shutdown_timeout"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Type     string `toml:"type"`
	Path     string `toml:"path"`
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"dbname"`
}

// DiscordConfig holds Discord bot configuration
type DiscordConfig struct {
	BotToken              string   `toml:"bot_token"`
	GuildID               string   `toml:"guild_id"`
	AdminIDs              []string `toml:"admin_ids"`
	RegistrationChannelID string   `toml:"registration_channel_id"`
	RegisteredRoleID      string   `toml:"registered_role_id"`
	RequireDiscord        bool     `toml:"require_discord"`
}

// WorldGenConfig holds world generation configuration
type WorldGenConfig struct {
	Dimensions string `toml:"dimensions"`
	ChunkSize  int    `toml:"chunk_size"`
}

// ReadConfig function
func ReadConfig(path string) (config Config, err error) {
	var conf Config

	// Decode the TOML file
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return conf, err
	}

	// Set defaults for any missing values
	conf.SetDefaults()

	// Load environment variable overrides
	conf.LoadFromEnvironment()

	// Validate the configuration
	if err := conf.Validate(); err != nil {
		return conf, err
	}

	return conf, nil
}
