package edenutil

import (
	"regexp"
	"strings"
	"unicode"
)

// ValidateUsername checks if a username is valid
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}

	// Only allow alphanumeric characters and underscores
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	return matched
}

// ValidatePassword checks if a password meets security requirements
func ValidatePassword(password string) bool {
	if len(password) < 8 || len(password) > 128 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// Require at least 3 of the 4 character types
	count := 0
	if hasUpper {
		count++
	}
	if hasLower {
		count++
	}
	if hasNumber {
		count++
	}
	if hasSpecial {
		count++
	}

	return count >= 3
}

// SanitizeInput removes potentially dangerous characters from input
func SanitizeInput(input string) string {
	// Remove null bytes and control characters except newlines and tabs
	var result strings.Builder
	for _, r := range input {
		if r == 0 || (r < 32 && r != 9 && r != 10 && r != 13) {
			continue
		}
		result.WriteRune(r)
	}
	return result.String()
}

// ValidateDiscordID checks if a Discord ID is valid format
func ValidateDiscordID(discordID string) bool {
	if len(discordID) < 17 || len(discordID) > 19 {
		return false
	}

	// Discord IDs are numeric
	matched, _ := regexp.MatchString("^[0-9]+$", discordID)
	return matched
}

// TruncateString safely truncates a string to a maximum length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
