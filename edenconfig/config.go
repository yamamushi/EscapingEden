package edenconfig

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

// Config is the main configuration struct
type Config struct {
	Logger loggingConfig `toml:"logging"`
	Server serverConfig  `toml:"server"`
	DB     dbConfig      `toml:"database"`
}

type loggingConfig struct {
	Type string `toml:"type"`
	Path string `toml:"path"`
}

// serverConfig struct
type serverConfig struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
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

// ReadConfig function
func ReadConfig(path string) (config Config, err error) {
	var conf Config
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		fmt.Println(err)
		return conf, err
	}
	return conf, nil
}
