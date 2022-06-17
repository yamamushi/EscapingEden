package edenbot

import (
	"github.com/yamamushi/EscapingEden/edenutil"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"log"
	"regexp"
	"strconv"
	"time"
)

func (eb *EdenBot) ForgotPassword(data messages.AccountForgotPasswordData) {
	eb.Log.Println(logging.LogInfo, "Edenbot Manager received a forgot password request")
	eb.Log.Println(logging.LogInfo, "Edenbot processing forgot password request for: ", data.DiscordTag)

	// Try to find the user in the database
	foundAccount, edenbotErr := eb.GetAccountByDiscordTag(data.DiscordTag)
	if edenbotErr != messages.Edenbot_Error_Null {
		eb.Log.Println(logging.LogWarn, "Attempted password reset for account not found: ", data.DiscordTag)
		return
	}

	if foundAccount != nil {
		if foundAccount.Username != data.Username {
			eb.Log.Println(logging.LogWarn, "Found discord account but invalid username: ", data.DiscordTag, " expected: ", foundAccount.Username)
			return
		}
	}

	// Generate a random token for the user to use to reset their password, store it in their account and send it to them
	if foundAccount.PasswordResetCode != "" {
		eb.Log.Println(logging.LogWarn, "Found existing validation code for: ", data.DiscordTag, ", skipping generation but sending notification")
		err := eb.SendPasswordResetToken(foundAccount)
		if err != nil {
			eb.Log.Println(logging.LogError, "Error sending password reset token: ", err)
		}
		return
	}
	eb.Log.Println(logging.LogInfo, "Generating password reset token for: ", data.DiscordTag)
	token, err := eb.GeneratePasswordResetToken(foundAccount)
	if err != nil {
		eb.Log.Println(logging.LogError, "Error generating password reset token for: ", data.DiscordTag, ": ", err)
		return
	}

	// We remove all non-alphanumeric characters from the token so that it can be entered easier
	// This doesn't make it easier to easier to guess the token
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	processedToken := reg.ReplaceAllString(token, "")
	foundAccount.PasswordResetCode = processedToken
	eb.Log.Println(logging.LogInfo, "Saving password reset token for: ", data.DiscordTag)
	edenbotErr = eb.SaveAccount(*foundAccount)
	if err != nil {
		eb.Log.Println(logging.LogError, "Error saving password reset token for: ", data.DiscordTag, ": ", edenbotErr)
		return
	}

}

func (eb *EdenBot) GeneratePasswordResetToken(account *messages.Account) (token string, err error) {
	id, err := strconv.Atoi(account.DiscordID)
	if err != nil {
		return "", err
	}
	edenutil.SeedRandom(time.Now().UnixNano() * int64(id))
	token, err = edenutil.GenerateRandomString(16)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (eb *EdenBot) SendPasswordResetToken(account *messages.Account) (err error) {
	// Send a message to the user with the token
	output := "You have requested a password reset.\n"
	output += "Please use the following token to reset your password:\n"
	output += "```" + account.PasswordResetCode + "```"
	_, err = eb.SendPrivateMessage(account.DiscordID, output)
	return err
}
