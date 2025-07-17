package accounts

import (
	"fmt"
	"testing"

	"github.com/yamamushi/EscapingEden/edenbot"
	"github.com/yamamushi/EscapingEden/messages"
	testutils "github.com/yamamushi/EscapingEden/testing"
)

func TestCreateAccount_ValidInput(t *testing.T) {
	helper := testutils.NewTestHelper(t)

	// Create mock edenbot
	mockBot := &MockEdenBot{}

	// Create account manager
	am := &AccountManager{
		DB:  helper.GetDB(),
		Log: helper.GetLogger(),
		EB:  mockBot,
	}

	// Test valid account creation
	response := am.CreateAccount("testuser", "TestPass123!", "testuser#1234")

	helper.AssertEqual(messages.AMError_Null, response.Error)
}

func TestCreateAccount_InvalidUsername(t *testing.T) {
	helper := testutils.NewTestHelper(t)

	mockBot := &MockEdenBot{}

	am := &AccountManager{
		DB:  helper.GetDB(),
		Log: helper.GetLogger(),
		EB:  mockBot,
	}

	// Test invalid username (too short)
	response := am.CreateAccount("ab", "TestPass123!", "testuser#1234")
	helper.AssertEqual(messages.AMError_InvalidUsername, response.Error)

	// Test invalid username (too long)
	response = am.CreateAccount("thisusernameiswaytoolongtobevalid", "TestPass123!", "testuser#1234")
	helper.AssertEqual(messages.AMError_InvalidUsername, response.Error)

	// Test invalid username (invalid characters)
	response = am.CreateAccount("test@user", "TestPass123!", "testuser#1234")
	helper.AssertEqual(messages.AMError_InvalidUsername, response.Error)
}

func TestCreateAccount_InvalidPassword(t *testing.T) {
	helper := testutils.NewTestHelper(t)

	mockBot := &MockEdenBot{}

	am := &AccountManager{
		DB:  helper.GetDB(),
		Log: helper.GetLogger(),
		EB:  mockBot,
	}

	// Test invalid password (too short)
	response := am.CreateAccount("testuser", "short", "testuser#1234")
	helper.AssertEqual(messages.AMError_InvalidPassword, response.Error)

	// Test invalid password (too simple)
	response = am.CreateAccount("testuser", "password", "testuser#1234")
	helper.AssertEqual(messages.AMError_InvalidPassword, response.Error)
}

func TestCreateAccount_InvalidDiscordID(t *testing.T) {
	helper := testutils.NewTestHelper(t)

	mockBot := &MockEdenBot{}

	am := &AccountManager{
		DB:  helper.GetDB(),
		Log: helper.GetLogger(),
		EB:  mockBot,
	}

	// Test invalid Discord ID
	response := am.CreateAccount("testuser", "TestPass123!", "invalid_discord_id")
	helper.AssertEqual(messages.AMError_InvalidDiscordID, response.Error)
}

func TestHashPassword(t *testing.T) {
	helper := testutils.NewTestHelper(t)

	am := &AccountManager{
		DB:  helper.GetDB(),
		Log: helper.GetLogger(),
	}

	password := "TestPassword123!"
	hash, err := am.HashPassword(password)

	helper.AssertNoError(err)

	// Hash should not be empty
	if hash == "" {
		t.Fatal("Hash should not be empty")
	}

	// Hash should not equal original password
	if hash == password {
		t.Fatal("Hash should not equal original password")
	}
}

func TestComparePasswords(t *testing.T) {
	helper := testutils.NewTestHelper(t)

	am := &AccountManager{
		DB:  helper.GetDB(),
		Log: helper.GetLogger(),
	}

	password := "TestPassword123!"
	hash, err := am.HashPassword(password)
	helper.AssertNoError(err)

	// Correct password should match
	if !am.ComparePasswords(hash, password) {
		t.Fatal("Correct password should match hash")
	}

	// Incorrect password should not match
	if am.ComparePasswords(hash, "WrongPassword") {
		t.Fatal("Incorrect password should not match hash")
	}
}

// MockEdenBot implements a mock EdenBot for testing
type MockEdenBot struct {
	shouldFailUserCheck  bool
	shouldFailValidation bool
}

// Ensure MockEdenBot implements EdenBotInterface
var _ edenbot.EdenBotInterface = (*MockEdenBot)(nil)

func (meb *MockEdenBot) IsUserInDiscordServer(discordTag string) (*edenbot.DiscordUser, error) {
	if meb.shouldFailUserCheck {
		return nil, fmt.Errorf("user not in server")
	}

	return &edenbot.DiscordUser{
		ID:       "123456789",
		Username: "testuser",
		Tag:      discordTag,
	}, nil
}

func (meb *MockEdenBot) ValidateUser(username, discordID string) error {
	if meb.shouldFailValidation {
		return fmt.Errorf("validation failed")
	}

	return nil
}
