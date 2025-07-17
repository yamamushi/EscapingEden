package edenutil

import (
	"testing"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		username string
		expected bool
	}{
		{"validuser", true},
		{"valid_user", true},
		{"user123", true},
		{"User123", true},
		{"ab", false}, // too short
		{"thisusernameiswaytoolongtobevalid", false}, // too long
		{"user@name", false},                         // invalid character
		{"user name", false},                         // space
		{"user-name", false},                         // hyphen
		{"", false},                                  // empty
	}

	for _, test := range tests {
		result := ValidateUsername(test.username)
		if result != test.expected {
			t.Errorf("ValidateUsername(%q) = %v, expected %v", test.username, result, test.expected)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password string
		expected bool
	}{
		{"Password123!", true}, // has upper, lower, number, special
		{"password123!", true}, // has lower, number, special (3 of 4)
		{"PASSWORD123!", true}, // has upper, number, special (3 of 4)
		{"Password!", true},    // has upper, lower, special (3 of 4)
		{"password", false},    // only lowercase
		{"PASSWORD", false},    // only uppercase
		{"12345678", false},    // only numbers
		{"!@#$%^&*", false},    // only special
		{"Pass1!", false},      // too short
		{"password123", false}, // only 2 types (lower, number)
		{"", false},            // empty
	}

	for _, test := range tests {
		result := ValidatePassword(test.password)
		if result != test.expected {
			t.Errorf("ValidatePassword(%q) = %v, expected %v", test.password, result, test.expected)
		}
	}
}

func TestValidateDiscordID(t *testing.T) {
	tests := []struct {
		discordID string
		expected  bool
	}{
		{"123456789012345678", true},    // valid 18-digit ID
		{"12345678901234567", true},     // valid 17-digit ID
		{"1234567890123456789", true},   // valid 19-digit ID
		{"123456789012345", false},      // too short
		{"12345678901234567890", false}, // too long
		{"12345678901234567a", false},   // contains letter
		{"", false},                     // empty
		{"abc", false},                  // non-numeric
	}

	for _, test := range tests {
		result := ValidateDiscordID(test.discordID)
		if result != test.expected {
			t.Errorf("ValidateDiscordID(%q) = %v, expected %v", test.discordID, result, test.expected)
		}
	}
}

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal text", "normal text"},
		{"text\x00with\x01null", "textwithnull"},             // removes null and control chars
		{"text\nwith\twhitespace", "text\nwith\twhitespace"}, // keeps newlines and tabs
		{"text\rwith\rcarriage", "text\rwith\rcarriage"},     // keeps carriage returns
		{"", ""},
	}

	for _, test := range tests {
		result := SanitizeInput(test.input)
		if result != test.expected {
			t.Errorf("SanitizeInput(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"this is a long string", 10, "this is a "},
		{"exact", 5, "exact"},
		{"", 5, ""},
		{"test", 0, ""},
	}

	for _, test := range tests {
		result := TruncateString(test.input, test.maxLen)
		if result != test.expected {
			t.Errorf("TruncateString(%q, %d) = %q, expected %q", test.input, test.maxLen, result, test.expected)
		}
	}
}
