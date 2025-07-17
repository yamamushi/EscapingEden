package edenconfig

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("config validation error in field '%s': %s", ve.Field, ve.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "no validation errors"
	}

	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Error())
	}

	return fmt.Sprintf("configuration validation failed: %s", strings.Join(messages, "; "))
}

// Validate validates the entire configuration
func (c *Config) Validate() error {
	var errors ValidationErrors

	// Validate server configuration
	if err := c.validateServer(); err != nil {
		if ve, ok := err.(ValidationErrors); ok {
			errors = append(errors, ve...)
		} else {
			errors = append(errors, ValidationError{Field: "server", Message: err.Error()})
		}
	}

	// Validate logger configuration
	if err := c.validateLogger(); err != nil {
		if ve, ok := err.(ValidationErrors); ok {
			errors = append(errors, ve...)
		} else {
			errors = append(errors, ValidationError{Field: "logger", Message: err.Error()})
		}
	}

	// Validate database configuration
	if err := c.validateDatabase(); err != nil {
		if ve, ok := err.(ValidationErrors); ok {
			errors = append(errors, ve...)
		} else {
			errors = append(errors, ValidationError{Field: "database", Message: err.Error()})
		}
	}

	// Validate Discord configuration
	if err := c.validateDiscord(); err != nil {
		if ve, ok := err.(ValidationErrors); ok {
			errors = append(errors, ve...)
		} else {
			errors = append(errors, ValidationError{Field: "discord", Message: err.Error()})
		}
	}

	// Validate WorldGen configuration
	if err := c.validateWorldGen(); err != nil {
		if ve, ok := err.(ValidationErrors); ok {
			errors = append(errors, ve...)
		} else {
			errors = append(errors, ValidationError{Field: "worldgen", Message: err.Error()})
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// validateServer validates server configuration
func (c *Config) validateServer() error {
	var errors ValidationErrors

	// Validate host
	if c.Server.Host == "" {
		errors = append(errors, ValidationError{Field: "server.host", Message: "host cannot be empty"})
	} else {
		// Check if host is a valid IP or hostname
		if ip := net.ParseIP(c.Server.Host); ip == nil {
			// Not an IP, check if it's a valid hostname
			if c.Server.Host != "localhost" && !isValidHostname(c.Server.Host) {
				errors = append(errors, ValidationError{Field: "server.host", Message: "invalid hostname or IP address"})
			}
		}
	}

	// Validate port
	if c.Server.Port == "" {
		errors = append(errors, ValidationError{Field: "server.port", Message: "port cannot be empty"})
	} else {
		if port, err := strconv.Atoi(c.Server.Port); err != nil {
			errors = append(errors, ValidationError{Field: "server.port", Message: "port must be a number"})
		} else if port < 1 || port > 65535 {
			errors = append(errors, ValidationError{Field: "server.port", Message: "port must be between 1 and 65535"})
		}
	}

	// Validate shutdown timeout
	if c.Server.ShutdownTimeout < 0 {
		errors = append(errors, ValidationError{Field: "server.shutdown_timeout", Message: "shutdown timeout cannot be negative"})
	} else if c.Server.ShutdownTimeout > 300 {
		errors = append(errors, ValidationError{Field: "server.shutdown_timeout", Message: "shutdown timeout should not exceed 300 seconds"})
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// validateLogger validates logger configuration
func (c *Config) validateLogger() error {
	var errors ValidationErrors

	// Validate logger type
	validTypes := []string{"console", "file", "disabled"}
	if !contains(validTypes, c.Logger.Type) {
		errors = append(errors, ValidationError{
			Field:   "logger.type",
			Message: fmt.Sprintf("invalid logger type '%s', must be one of: %s", c.Logger.Type, strings.Join(validTypes, ", ")),
		})
	}

	// If file logger, validate path
	if c.Logger.Type == "file" && c.Logger.Path == "" {
		errors = append(errors, ValidationError{Field: "logger.path", Message: "path is required for file logger"})
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// validateDatabase validates database configuration
func (c *Config) validateDatabase() error {
	var errors ValidationErrors

	// Validate database type
	validTypes := []string{"bolt", "mongodb"}
	if !contains(validTypes, strings.ToLower(c.DB.Type)) {
		errors = append(errors, ValidationError{
			Field:   "database.type",
			Message: fmt.Sprintf("invalid database type '%s', must be one of: %s", c.DB.Type, strings.Join(validTypes, ", ")),
		})
	}

	// Validate path for bolt
	if strings.ToLower(c.DB.Type) == "bolt" && c.DB.Path == "" {
		errors = append(errors, ValidationError{Field: "database.path", Message: "path is required for bolt database"})
	}

	// Validate MongoDB configuration
	if strings.ToLower(c.DB.Type) == "mongodb" {
		if c.DB.Host == "" {
			errors = append(errors, ValidationError{Field: "database.host", Message: "host is required for mongodb"})
		}
		if c.DB.Port == "" {
			errors = append(errors, ValidationError{Field: "database.port", Message: "port is required for mongodb"})
		} else {
			if port, err := strconv.Atoi(c.DB.Port); err != nil {
				errors = append(errors, ValidationError{Field: "database.port", Message: "port must be a number"})
			} else if port < 1 || port > 65535 {
				errors = append(errors, ValidationError{Field: "database.port", Message: "port must be between 1 and 65535"})
			}
		}
		if c.DB.DBName == "" {
			errors = append(errors, ValidationError{Field: "database.dbname", Message: "database name is required for mongodb"})
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// validateDiscord validates Discord configuration
func (c *Config) validateDiscord() error {
	var errors ValidationErrors

	// Bot token validation (basic check)
	if c.Discord.BotToken == "" {
		errors = append(errors, ValidationError{Field: "discord.bot_token", Message: "bot token cannot be empty"})
	} else if len(c.Discord.BotToken) < 50 {
		errors = append(errors, ValidationError{Field: "discord.bot_token", Message: "bot token appears to be invalid (too short)"})
	}

	// Guild ID validation
	if c.Discord.GuildID == "" {
		errors = append(errors, ValidationError{Field: "discord.guild_id", Message: "guild ID cannot be empty"})
	} else if !isValidDiscordID(c.Discord.GuildID) {
		errors = append(errors, ValidationError{Field: "discord.guild_id", Message: "invalid guild ID format"})
	}

	// Admin IDs validation
	if len(c.Discord.AdminIDs) == 0 {
		errors = append(errors, ValidationError{Field: "discord.admin_ids", Message: "at least one admin ID is required"})
	} else {
		for i, adminID := range c.Discord.AdminIDs {
			if !isValidDiscordID(adminID) {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("discord.admin_ids[%d]", i),
					Message: "invalid admin ID format",
				})
			}
		}
	}

	// Channel ID validation
	if c.Discord.RegistrationChannelID != "" && !isValidDiscordID(c.Discord.RegistrationChannelID) {
		errors = append(errors, ValidationError{Field: "discord.registration_channel_id", Message: "invalid channel ID format"})
	}

	// Role ID validation
	if c.Discord.RegisteredRoleID != "" && !isValidDiscordID(c.Discord.RegisteredRoleID) {
		errors = append(errors, ValidationError{Field: "discord.registered_role_id", Message: "invalid role ID format"})
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// validateWorldGen validates world generation configuration
func (c *Config) validateWorldGen() error {
	var errors ValidationErrors

	// Validate dimensions
	if c.WorldGen.Dimensions == "" {
		errors = append(errors, ValidationError{Field: "worldgen.dimensions", Message: "dimensions cannot be empty"})
	} else {
		parts := strings.Split(c.WorldGen.Dimensions, ",")
		if len(parts) != 3 {
			errors = append(errors, ValidationError{Field: "worldgen.dimensions", Message: "dimensions must be in format 'x,y,z'"})
		} else {
			for i, part := range parts {
				if dim, err := strconv.Atoi(strings.TrimSpace(part)); err != nil {
					errors = append(errors, ValidationError{
						Field:   fmt.Sprintf("worldgen.dimensions[%d]", i),
						Message: "dimension must be a number",
					})
				} else if dim < 1 {
					errors = append(errors, ValidationError{
						Field:   fmt.Sprintf("worldgen.dimensions[%d]", i),
						Message: "dimension must be positive",
					})
				} else if dim > 1000 {
					errors = append(errors, ValidationError{
						Field:   fmt.Sprintf("worldgen.dimensions[%d]", i),
						Message: "dimension should not exceed 1000 for performance reasons",
					})
				}
			}
		}
	}

	// Validate chunk size
	if c.WorldGen.ChunkSize < 16 {
		errors = append(errors, ValidationError{Field: "worldgen.chunk_size", Message: "chunk size must be at least 16"})
	} else if c.WorldGen.ChunkSize > 512 {
		errors = append(errors, ValidationError{Field: "worldgen.chunk_size", Message: "chunk size should not exceed 512 for performance reasons"})
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// Helper functions

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// isValidHostname checks if a string is a valid hostname
func isValidHostname(hostname string) bool {
	if len(hostname) == 0 || len(hostname) > 253 {
		return false
	}

	// Simple hostname validation
	for _, char := range hostname {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') || char == '-' || char == '.') {
			return false
		}
	}

	return true
}

// isValidDiscordID checks if a string is a valid Discord ID
func isValidDiscordID(id string) bool {
	if len(id) < 17 || len(id) > 19 {
		return false
	}

	for _, char := range id {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}
