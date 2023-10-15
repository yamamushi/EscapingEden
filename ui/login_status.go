package ui

import (
	"github.com/yamamushi/EscapingEden/messages"
)

// IsUserLoggedIn returns whether or not the user is logged in
func (c *Console) IsUserLoggedIn() bool {
	c.userLoggedInMutex.Lock()
	defer c.userLoggedInMutex.Unlock()
	return c.userLoggedIn
}

// LoginUser logs in the user, sets the userLoggedIn flag to true
func (c *Console) LoginUser(userInfo messages.UserInfo) {
	c.userLoggedInMutex.Lock()
	defer c.userLoggedInMutex.Unlock()
	// First we update the user info
	c.UpdateUserInfo(userInfo)
	c.userLoggedIn = true

	// Now we close the login window and open the user Dashboard window as active
	//c.Log.Println(logging.LogInfo, "Removing Login Window")
	c.RemoveWindow(c.GetLoginWindow().GetID())
	//c.Log.Println(logging.LogInfo, "Creating Dashboard Window")
	c.SetActiveWindowNoThread(c.GetUserDashboard())
	c.UpdateWindowsUserInfo()
	//c.Log.Println(logging.LogInfo, "Dashboard Active")
}

// LogoutUser logs out the user, sets the userLoggedIn flag to false and sets the active window to the login window
func (c *Console) LogoutUser() {
	c.userLoggedInMutex.Lock()
	defer c.userLoggedInMutex.Unlock()
	c.userLoggedIn = false
	c.UpdateUserInfo(messages.UserInfo{})
	c.LogoutCharacter()
	c.RemoveWindow(c.GetUserDashboard().GetID())
	c.SetActiveWindowNoThread(c.GetLoginWindow())
}

func (c *Console) IsCharacterLoggedIn() bool {
	c.characterLoggedInMutex.Lock()
	defer c.characterLoggedInMutex.Unlock()
	return c.characterLoggedIn
}

// LoginCharacter logs in the character, sets the characterLoggedIn flag to true
func (c *Console) LoginCharacter(charInfo messages.CharacterInfo) {
	c.characterLoggedInMutex.Lock()
	defer c.characterLoggedInMutex.Unlock()
	c.RemoveWindow(c.GetUserDashboard().GetID())
	c.SetActiveWindowNoThread(c.GetGameWindow())
	c.UpdateCharacterInfo(charInfo)
	c.currentCharID = charInfo.ID
	c.SendMessages <- messages.ConnectionManagerMessage{Data: messages.GameManagerMessage{Data: messages.GameMessageData{CharacterID: charInfo.ID}, Type: messages.GameManager_NotifyLoggedInCharacter}, Type: messages.ConnectManager_Message_GameCommand, SenderConsoleID: c.ConnectionID}
	if !c.characterLoggedIn {
		chatMessage := messages.ChatMessage{}
		if int(charInfo.FirstLogin) == 1 {
			chatMessage = messages.ChatMessage{Type: messages.Chat_Message_System, Content: "Welcome " + c.GetCharacterName() + "!"}
		} else {
			chatMessage = messages.ChatMessage{Type: messages.Chat_Message_System, Content: "Welcome back " + c.GetCharacterName() + "!"}
		}
		c.ChatMessageReceive <- chatMessage
	}
	c.characterLoggedIn = true
	c.ClearPointMap()
	c.FlushLastSent()
	//c.forceScreenRefresh = true
}

// LogoutCharacter logs out the character, sets the characterLoggedIn flag to false
func (c *Console) LogoutCharacter() {
	c.characterLoggedInMutex.Lock()
	defer c.characterLoggedInMutex.Unlock()
	c.characterLoggedIn = false
	c.SendMessages <- messages.ConnectionManagerMessage{Data: messages.GameManagerMessage{Data: messages.GameMessageData{CharacterID: c.currentCharID}, Type: messages.GameManager_NotifyLoggedInCharacter}, Type: messages.ConnectManager_Message_GameCommand, SenderConsoleID: c.ConnectionID}
	c.UpdateCharacterInfo(messages.CharacterInfo{})
	c.RemoveWindow(c.GetGameWindow().GetID())
	c.SetActiveWindowNoThread(c.GetUserDashboard())
}
