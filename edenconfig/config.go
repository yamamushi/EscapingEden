package edenconfig

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

// Config is the main configuration struct
type Config struct {
	Logger  loggingConfig  `toml:"logging"`
	Server  serverConfig   `toml:"server"`
	DB      dbConfig       `toml:"database"`
	Discord discordConfig  `toml:"discord"`
	World   worldgenConfig `toml:"worldgen"`
}

type loggingConfig struct {
	Type string `toml:"type"`
	Path string `toml:"path"`
}

// serverConfig struct
type serverConfig struct {
	Host            string `toml:"host"`
	Port            string `toml:"port"`
	ShutdownTimeout int    `toml:"shutdown_timeout"`
}

// dbConfig struct
type dbConfig struct {
	Type string `toml:"type"`
	Path string `toml:"path"`
	Host string `toml:"host"`
	Port string `toml:"port"`
	User string `toml:"user"`
	Pass string `toml:"pass"`
	Name string `toml:"name"`
}

// discordConfig struct
type discordConfig struct {
	Token                 string   `toml:"bot_token"`
	GuildID               string   `toml:"guild_id"`
	AdminIDs              []string `toml:"admin_ids"`
	RegistrationChannelID string   `toml:"registration_channel_id"`
	RegisteredRoleID      string   `toml:"registered_role_id"`
}

type worldgenConfig struct {
	Dimensions string `toml:"dimensions"`
}

// ReadConfig function
func ReadConfig(path string) (config Config, err error) {
	var conf Config
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		fmt.Println(err)
		return conf, err
	}
	return conf, nil
}
