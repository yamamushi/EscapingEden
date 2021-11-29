package edenconfig

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type Config struct {
	Server serverConfig `toml:"server"`
	DB     dbConfig     `toml:"database"`
}

// serverConfig struct
type serverConfig struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

// dbConfig struct
type dbConfig struct {
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
