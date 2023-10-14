package edenbot

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"strings"
)

func (eb *EdenBot) IsUserInDiscordServer(username string) (*discordgo.User, error) {

	/*usertag := strings.Split(username, "#")
	if len(usertag) != 2 {
		eb.Log.Println(logging.LogInfo, "Invalid username: ", username)
		return nil, errors.New("invalid discord username: " + username)
	}*/

	members, err := eb.dg.GuildMembersSearch(eb.Config.Discord.GuildID, username, 100)
	if err != nil {
		eb.Log.Println(logging.LogError, "Error searching for user: ", err)
		return nil, errors.New("error searching for user")
	}
	if len(members) == 0 {
		eb.Log.Println(logging.LogError, "User not found in server")
		return nil, errors.New("user not found in server")
	}

	for _, member := range members {
		eb.Log.Println(logging.LogInfo, "Checking user: ", member.User.Discriminator)
		if member.User.Discriminator == username {
			eb.Log.Println(logging.LogInfo, "User found in server")
			return member.User, nil
		}
		if member.User.Username == username {
			eb.Log.Println(logging.LogInfo, "User found in server")
			return member.User, nil
		}
	}

	eb.Log.Println(logging.LogInfo, "User not found in server", username)
	return nil, errors.New("user not found in server")
}

func (eb *EdenBot) ValidateUser(username string, userID string) error {

	outputMessage := "Hello, " + username + "!\n\n"
	outputMessage += "Thank you for registering an account on Escaping Eden! You must now verify your account before you " +
		"can log in.\n\n"
	outputMessage += "Please respond to this message with your verification code using the !verify command.\n\n"
	outputMessage += "Example: ```!verify 123456```\n\n"

	_, err := eb.SendPrivateMessage(userID, outputMessage)
	if err != nil {
		return err
	}
	return nil
}

func (eb *EdenBot) HandleVerify(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, "!verify") {
		eb.Log.Println(logging.LogInfo, "Verifying user: ", m.Author.Username)
		account := messages.Account{}
		err := eb.DB.One("Accounts", "DiscordID", m.Author.ID, &account)
		if err != nil {
			eb.Log.Println(logging.LogError, "Error finding account: ", err)
			_, _ = eb.dg.ChannelMessageSend(m.ChannelID, "Error finding account: "+err.Error())
			return
		}

		if account.ValidationStatus == 1 {
			_, _ = eb.dg.ChannelMessageSend(m.ChannelID, "Your account is already validated!")
			return
		}

		if account.ValidationCode == m.Content[8:] {
			account.ValidationStatus = 1
			dbError := eb.DB.UpdateRecord("Accounts", &account)
			if dbError != nil {
				eb.Log.Println(logging.LogError, "Error updating account: ", dbError)
				_, _ = eb.dg.ChannelMessageSend(m.ChannelID, "Error updating account in database: "+dbError.Error())
				return
			}
		} else {
			_, _ = eb.dg.ChannelMessageSend(m.ChannelID, "Invalid validation code!")
			return
		}
		_, _ = eb.dg.ChannelMessageSend(m.ChannelID, "Account validation successful, you may now enter Eden.")
		err = eb.dg.GuildMemberRoleAdd(eb.Config.Discord.GuildID, m.Author.ID, eb.Config.Discord.RegisteredRoleID)
		if err != nil {
			eb.Log.Println(logging.LogError, "Error adding registered role to user: ", err, m.Author.ID)
		}
		eb.NotifyRegistrationChannel(m.Author.ID)
	}
}

func (eb *EdenBot) NotifyRegistrationChannel(userID string) {
	user, err := eb.dg.User(userID)
	if err != nil {
		eb.Log.Println(logging.LogError, "NotifyRegistrationChannel error finding user: ", err)
		return
	}
	_, err = eb.dg.ChannelMessageSend(eb.Config.Discord.RegistrationChannelID, user.Mention()+" has been successfully registered to enter Eden!")
	if err != nil {
		eb.Log.Println(logging.LogError, "NotifyRegistrationChannel error sending message: ", err)
	}
}
