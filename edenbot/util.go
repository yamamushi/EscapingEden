package edenbot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/yamamushi/EscapingEden/messages"
)

func (eb *EdenBot) GetAccountByDiscordTag(discordTag string) (*messages.Account, messages.EdenbotErrorType) {
	result := messages.Account{}
	err := eb.DB.One("Accounts", "DiscordTag", discordTag, &result)
	if err != nil {
		if err.Error() == "not found" {
			return nil, messages.Edenbot_Error_Null // no error, no account
		}
		return nil, messages.Edenbot_Error_DB
	}
	return &result, messages.Edenbot_Error_Null
}

func (eb *EdenBot) GetAccountByDiscordID(discordID string) (*messages.Account, messages.EdenbotErrorType) {
	result := messages.Account{}
	err := eb.DB.One("Accounts", "DiscordID", discordID, &result)
	if err != nil {
		if err.Error() == "not found" {
			return nil, messages.Edenbot_Error_Null // no error, no account
		}
		return nil, messages.Edenbot_Error_DB
	}
	return &result, messages.Edenbot_Error_Null
}

func (eb *EdenBot) GetAccountByUsername(username string) (*messages.Account, messages.EdenbotErrorType) {
	result := messages.Account{}
	err := eb.DB.One("Accounts", "Username", username, &result)
	if err != nil {
		if err.Error() == "not found" {
			return nil, messages.Edenbot_Error_Null // no error, no account
		}
		return nil, messages.Edenbot_Error_DB
	}
	return &result, messages.Edenbot_Error_Null
}

func (eb *EdenBot) SaveAccount(account messages.Account) messages.EdenbotErrorType {
	err := eb.DB.UpdateRecord("Accounts", &account)
	if err != nil {
		return messages.Edenbot_Error_DB
	}
	return messages.Edenbot_Error_Null
}

func (eb *EdenBot) SendPrivateMessage(discordID string, output string) (message *discordgo.Message, err error) {
	userprivatechannel, err := eb.dg.UserChannelCreate(discordID)
	if err != nil {
		return nil, err
	}

	message, err = eb.dg.ChannelMessageSend(userprivatechannel.ID, output)
	if err != nil {
		return nil, err
	}
	return message, nil
}
