package edenconfig

import (
	"os"
	"strconv"
	"strings"
)

// LoadFromEnvironment loads configuration values from environment variables
// Environment variables override config file values
func (c *Config) LoadFromEnvironment() {
	// Server configuration
	if host := os.Getenv("EDEN_SERVER_HOST"); host != "" {
		c.Server.Host = host
	}
	if port := os.Getenv("EDEN_SERVER_PORT"); port != "" {
		c.Server.Port = port
	}
	if timeout := os.Getenv("EDEN_SERVER_SHUTDOWN_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			c.Server.ShutdownTimeout = t
		}
	}

	// Logger configuration
	if logType := os.Getenv("EDEN_LOGGER_TYPE"); logType != "" {
		c.Logger.Type = logType
	}
	if logPath := os.Getenv("EDEN_LOGGER_PATH"); logPath != "" {
		c.Logger.Path = logPath
	}

	// Database configuration
	if dbType := os.Getenv("EDEN_DB_TYPE"); dbType != "" {
		c.DB.Type = dbType
	}
	if dbPath := os.Getenv("EDEN_DB_PATH"); dbPath != "" {
		c.DB.Path = dbPath
	}
	if dbHost := os.Getenv("EDEN_DB_HOST"); dbHost != "" {
		c.DB.Host = dbHost
	}
	if dbPort := os.Getenv("EDEN_DB_PORT"); dbPort != "" {
		c.DB.Port = dbPort
	}
	if dbUser := os.Getenv("EDEN_DB_USER"); dbUser != "" {
		c.DB.User = dbUser
	}
	if dbPassword := os.Getenv("EDEN_DB_PASSWORD"); dbPassword != "" {
		c.DB.Password = dbPassword
	}
	if dbName := os.Getenv("EDEN_DB_NAME"); dbName != "" {
		c.DB.DBName = dbName
	}

	// Discord configuration
	if botToken := os.Getenv("EDEN_DISCORD_BOT_TOKEN"); botToken != "" {
		c.Discord.BotToken = botToken
	}
	if guildID := os.Getenv("EDEN_DISCORD_GUILD_ID"); guildID != "" {
		c.Discord.GuildID = guildID
	}
	if adminIDs := os.Getenv("EDEN_DISCORD_ADMIN_IDS"); adminIDs != "" {
		c.Discord.AdminIDs = strings.Split(adminIDs, ",")
		// Trim whitespace from each ID
		for i, id := range c.Discord.AdminIDs {
			c.Discord.AdminIDs[i] = strings.TrimSpace(id)
		}
	}
	if regChannelID := os.Getenv("EDEN_DISCORD_REGISTRATION_CHANNEL_ID"); regChannelID != "" {
		c.Discord.RegistrationChannelID = regChannelID
	}
	if roleID := os.Getenv("EDEN_DISCORD_REGISTERED_ROLE_ID"); roleID != "" {
		c.Discord.RegisteredRoleID = roleID
	}
	if requireDiscord := os.Getenv("EDEN_DISCORD_REQUIRE_DISCORD"); requireDiscord != "" {
		if require, err := strconv.ParseBool(requireDiscord); err == nil {
			c.Discord.RequireDiscord = require
		}
	}

	// WorldGen configuration
	if dimensions := os.Getenv("EDEN_WORLDGEN_DIMENSIONS"); dimensions != "" {
		c.WorldGen.Dimensions = dimensions
	}
	if chunkSize := os.Getenv("EDEN_WORLDGEN_CHUNK_SIZE"); chunkSize != "" {
		if size, err := strconv.Atoi(chunkSize); err == nil {
			c.WorldGen.ChunkSize = size
		}
	}
}

// GetEnvironmentOverrides returns a map of environment variables that would override config
func GetEnvironmentOverrides() map[string]string {
	overrides := make(map[string]string)

	envVars := []string{
		"EDEN_SERVER_HOST",
		"EDEN_SERVER_PORT",
		"EDEN_SERVER_SHUTDOWN_TIMEOUT",
		"EDEN_LOGGER_TYPE",
		"EDEN_LOGGER_PATH",
		"EDEN_DB_TYPE",
		"EDEN_DB_PATH",
		"EDEN_DB_HOST",
		"EDEN_DB_PORT",
		"EDEN_DB_USER",
		"EDEN_DB_PASSWORD",
		"EDEN_DB_NAME",
		"EDEN_DISCORD_BOT_TOKEN",
		"EDEN_DISCORD_GUILD_ID",
		"EDEN_DISCORD_ADMIN_IDS",
		"EDEN_DISCORD_REGISTRATION_CHANNEL_ID",
		"EDEN_DISCORD_REGISTERED_ROLE_ID",
		"EDEN_DISCORD_REQUIRE_DISCORD",
		"EDEN_WORLDGEN_DIMENSIONS",
		"EDEN_WORLDGEN_CHUNK_SIZE",
	}

	for _, envVar := range envVars {
		if value := os.Getenv(envVar); value != "" {
			overrides[envVar] = value
		}
	}

	return overrides
}

// SetDefaults sets default values for configuration
func (c *Config) SetDefaults() {
	// Server defaults
	if c.Server.Host == "" {
		c.Server.Host = "localhost"
	}
	if c.Server.Port == "" {
		c.Server.Port = "8080"
	}
	if c.Server.ShutdownTimeout == 0 {
		c.Server.ShutdownTimeout = 10
	}

	// Logger defaults
	if c.Logger.Type == "" {
		c.Logger.Type = "console"
	}

	// Database defaults
	if c.DB.Type == "" {
		c.DB.Type = "bolt"
	}
	if c.DB.Path == "" && strings.ToLower(c.DB.Type) == "bolt" {
		c.DB.Path = "eden.db"
	}
	if c.DB.Port == "" && strings.ToLower(c.DB.Type) == "mongodb" {
		c.DB.Port = "27017"
	}

	// Discord defaults
	// Default to NOT requiring Discord registration (false)
	// This can be overridden in config file or environment variable
	if !c.Discord.RequireDiscord && os.Getenv("EDEN_DISCORD_REQUIRE_DISCORD") == "" {
		c.Discord.RequireDiscord = false
	}

	// WorldGen defaults
	if c.WorldGen.Dimensions == "" {
		c.WorldGen.Dimensions = "100,100,10"
	}
	if c.WorldGen.ChunkSize == 0 {
		c.WorldGen.ChunkSize = 64
	}
}
