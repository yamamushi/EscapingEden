package accounts

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

func (am *AccountManager) Start(started chan bool) error {
	go am.HandleMessages(started)
	return nil
}

func (am *AccountManager) HandleMessages(started chan bool) {
	am.Log.Println(logging.LogInfo, "Account Manager now handling messages")
	started <- true
	for {
		select {
		case managerMessage := <-am.ReceiveChannel:
			//am.Log.Println(logging.LogInfo, "Account Manager received message")
			switch managerMessage.Type {
			case messages.AccountManager_Message_Register:
				req := managerMessage.Data.(messages.AccountRegistrationRequest)
				//am.Log.Println(logging.LogInfo, "Account Manager received registration request")
				registrationResponse := am.CreateAccount(req.Username, req.Password, req.DiscordID)
				response := messages.ConnectionManagerMessage{
					Type:               messages.ConnectManager_Message_RegisterResponse,
					RecipientConsoleID: managerMessage.SenderConsoleID,
					Data:               registrationResponse,
				}
				am.Log.Println(logging.LogInfo, "Sending registration response")
				am.SendChannel <- response

			case messages.AccountManager_Message_Login:
				req := managerMessage.Data.(messages.AccountLoginRequest)
				//am.Log.Println(logging.LogInfo, "Account Manager received login request")

				loginResponse := am.handleLogin(req.Username, req.Password, managerMessage.SenderConsoleID)
				response := messages.ConnectionManagerMessage{
					Type:               messages.ConnectManager_Message_LoginResponse,
					RecipientConsoleID: managerMessage.SenderConsoleID,
					Data:               loginResponse,
				}
				//am.Log.Println(logging.LogInfo, "Sending login response")
				am.SendChannel <- response
			case messages.AccountManager_Message_ResetPasswordValidate:
				req := managerMessage.Data.(messages.AccountProcessForgotPasswordData)
				//am.Log.Println(logging.LogInfo, "Account Manager received reset password validation request")
				//am.Log.Println(logging.LogInfo, "Requested for: ", req.Username)
				validated := am.handleValidatePasswordReset(req)

				response := messages.ConnectionManagerMessage{
					Type:               messages.ConnectManager_Message_ValidatePasswordResetResponse,
					RecipientConsoleID: managerMessage.SenderConsoleID,
					Data:               validated,
				}
				//am.Log.Println(logging.LogInfo, "Sending reset password validation response")
				am.SendChannel <- response

			case messages.AccountManager_Message_ResetPasswordProcess:
				req := managerMessage.Data.(messages.AccountProcessForgotPasswordData)
				//am.Log.Println(logging.LogInfo, "Account Manager received reset password process request")

				status := am.handleChangePassword(req)
				response := messages.ConnectionManagerMessage{
					Type:               messages.ConnectManager_Message_ProcessPasswordResetResponse,
					RecipientConsoleID: managerMessage.SenderConsoleID,
					Data:               status,
				}
				//am.Log.Println(logging.LogInfo, "Sending process reset password response")
				am.SendChannel <- response

			case messages.AccountManager_Message_UpdateCharacterHistory:
				//am.Log.Println(logging.LogInfo, "Account Manager received update character history request")
				req := managerMessage.Data.(messages.CharacterInfo)
				response := messages.ConnectionManagerMessage{
					Type:               messages.ConnectManager_Message_UpdateAccountHistoryResponse,
					RecipientConsoleID: managerMessage.SenderConsoleID,
					Data: messages.CharManagerUpdateHistoryResponse{
						Data:              am.UpdateLoginCharacterHistory(req),
						RespondingManager: "account",
					},
				}
				am.SendChannel <- response
			}
		}
	}
}
