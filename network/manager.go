package network

import (
	"github.com/google/uuid"
	"github.com/yamamushi/EscapingEden/edendb"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"net"
	"sync"
	"time"
)

// ConnectionManager synchronizes connection output globally
type ConnectionManager struct {
	Log logging.LoggerType
	// Mutex for locking
	mutex         sync.Mutex
	connectionMap *sync.Map

	DB edendb.DatabaseType

	// Channel for receiving messages
	CMReceiveMessages chan messages.ConnectionManagerMessage

	// Our AccountManager
	//AccountManager *accounts.AccountManager
	// Account manager outbound channel
	AMSendMessages chan messages.AccountManagerMessage

	// Our CharacterManager outbound Channel
	CMSendMessages chan messages.CharacterManagerMessage

	// Our EdenBot Manager
	EBSendMessages chan messages.EdenbotMessage

	GMSendMessages chan messages.GameManagerMessage
}

// NewConnectionManager creates a new ConnectionManager
func NewConnectionManager(connectionMap *sync.Map, receiveMessages chan messages.ConnectionManagerMessage,
	accountManagerMessages chan messages.AccountManagerMessage, characterManagerReceiveMessages chan messages.CharacterManagerMessage, ebSendMessages chan messages.EdenbotMessage, gmSendMessages chan messages.GameManagerMessage, db edendb.DatabaseType, log logging.LoggerType) *ConnectionManager {
	return &ConnectionManager{
		connectionMap:     connectionMap,
		CMReceiveMessages: receiveMessages,
		AMSendMessages:    accountManagerMessages,
		CMSendMessages:    characterManagerReceiveMessages,
		EBSendMessages:    ebSendMessages,
		GMSendMessages:    gmSendMessages,
		Log:               log,
		DB:                db,
	}
}

// AddConnection adds a connection to the manager
func (cm *ConnectionManager) AddConnection(connection *Connection) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	if cm.checkIPBadLogins(connection) {
		cm.connectionMap.Store(connection.ID, connection)
	} else {
		connection.Write([]byte("This IP is temporarily banned. Please try again later."))
		connection.Close()
	}
}

// HandleDisconnect handles disconnect events
func (cm *ConnectionManager) HandleDisconnect(connection *Connection) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.connectionMap.Delete(connection.ID)
	cm.GMSendMessages <- messages.GameManagerMessage{Type: messages.GameManager_NotifyDisconnect, Data: connection.ID}
}

// Run launches the message parser handler
func (cm *ConnectionManager) Run(startedNotify chan bool) {
	go cm.MessageParser(startedNotify) // Launch our goroutine that listens for incoming messages
}

// MessageParser performs a non-blocking check for messages on cm.CMReceiveMessages
func (cm *ConnectionManager) MessageParser(startedNotify chan bool) {
	cm.Log.Println(logging.LogInfo, "Connection Manager is now listening for incoming messages")
	startedNotify <- true
	for {
		select {
		case managerMessage := <-cm.CMReceiveMessages:
			//cm.Log.Println(logging.LogInfo, "Message received from cm.CMReceiveMessages")
			switch managerMessage.Type {
			case messages.ConnectManager_Message_Chat:
				// For every connection, send the message to the Console channel
				cm.connectionMap.Range(func(key, value interface{}) bool {
					if conn, ok := value.(*Connection); ok {
						//cm.Log.Println(logging.LogInfo, "Chat message found, sending to conn.Console.ReceiveMessages")
						//outMessage := messages.ConsoleMessage{Data: managerMessage.SenderConsoleID + ": " + managerMessage.Data.(string), Type: messages.Console_Message_Chat}
						outMessage := messages.ConsoleMessage{Data: managerMessage.Data.(string), Type: messages.Console_Message_Chat}
						conn.SendToConsole(outMessage)
					}
					return true
				})
			case messages.ConnectManager_Message_Broadcast:
				// For every connection, send the message to the Console channel
				cm.Log.Println(logging.LogInfo, "Broadcast message received on Connection Manager, sending to all connected clients: ", managerMessage.Data)
				cm.connectionMap.Range(func(key, value interface{}) bool {
					if conn, ok := value.(*Connection); ok {
						consoleMessage := messages.ConsoleMessage{Type: messages.Console_Message_Broadcast, Data: managerMessage.Data}
						conn.SendToConsole(consoleMessage)
					}
					return true
				})
			case messages.ConnectManager_Message_ServerShutdown:
				cm.Log.Println(logging.LogInfo, "Server shutdown message received on Connection Manager, disconnecting all clients.")
				cm.connectionMap.Range(func(key, value interface{}) bool {
					if conn, ok := value.(*Connection); ok {
						cm.HandleDisconnect(conn)
						conn.Write([]byte("\033c\r\nThe Escaping Eden server has shut down for maintenance.\r\n" +
							"For the latest status updates, be sure to check the Escaping Eden Discord at: https://discord.gg/uMxZnjJGGu\r\n\r\n"))
						conn.Close()
					}
					return true
				})
			case messages.ConnectManager_Message_Error:
				// For every connection, send the message to the Console channel
				cm.connectionMap.Range(func(key, value interface{}) bool {
					if conn, ok := value.(*Connection); ok {
						if managerMessage.RecipientConsoleID == conn.ID {
							//cm.Log.Println(logging.LogInfo, "Error message found, sending to conn.Console.ReceiveMessages")
							consoleMessage := messages.ConsoleMessage{Type: messages.Console_Message_Error, Data: managerMessage.Data}
							conn.SendToConsole(consoleMessage)

						}
					}
					return true
				})
			case messages.ConnectManager_Message_Quit:
				cm.connectionMap.Range(func(key, value interface{}) bool {
					if conn, ok := value.(*Connection); ok {
						if managerMessage.RecipientConsoleID == conn.ID {
							//cm.Log.Println(logging.LogInfo, "Quit message found, sending to conn.Console.ReceiveMessages")
							consoleMessage := messages.ConsoleMessage{Type: messages.Console_Message_Quit}
							conn.SendToConsole(consoleMessage)
						}
					}
					return true
				})
			case messages.ConnectManager_Message_Register:
				go func() {
					//cm.Log.Println(logging.LogInfo, "Sending registration request to AccountManager")
					registrationRequest := managerMessage.Data.(messages.AccountRegistrationRequest)
					cm.AMSendMessages <- messages.AccountManagerMessage{Type: messages.AccountManager_Message_Register, Data: registrationRequest, SenderConsoleID: managerMessage.SenderConsoleID}
				}()

			case messages.ConnectManager_Message_AccountLogin:
				go func() {
					//cm.Log.Println(logging.LogInfo, "Sending login request to AccountManager")
					loginRequest := managerMessage.Data.(messages.AccountLoginRequest)
					cm.AMSendMessages <- messages.AccountManagerMessage{Type: messages.AccountManager_Message_Login, Data: loginRequest, SenderConsoleID: managerMessage.SenderConsoleID}
				}()

			case messages.ConnectManager_Message_CharacterLoggedInNotify:
				go func() {
					amMessage := messages.AccountManagerMessage{
						Type:            messages.AccountManager_Message_UpdateCharacterHistory,
						Data:            managerMessage.Data,
						SenderConsoleID: managerMessage.SenderConsoleID,
					}
					cm.AMSendMessages <- amMessage

					cmMessage := messages.CharacterManagerMessage{
						Type:            messages.CharManager_UpdateLoginHistory,
						Data:            managerMessage.Data,
						SenderConsoleID: managerMessage.SenderConsoleID,
					}
					cm.CMSendMessages <- cmMessage
				}()

			case messages.ConnectManager_Message_RequestCharacterByID:
				cm.Log.Println(logging.LogInfo, "Requesting character by ID: ", managerMessage.Data)
				go func() {
					cmMessage := messages.CharacterManagerMessage{
						Type:            messages.CharManager_RequestCharacterByID,
						Data:            managerMessage.Data,
						SenderConsoleID: managerMessage.SenderConsoleID,
					}
					cm.CMSendMessages <- cmMessage
				}()

			case messages.ConnectManager_Message_CharacterRequestResponse:
				go func() {
					cm.Log.Println(logging.LogInfo, "Sending character to client that requested it")
					//log.Println(managerMessage.RecipientConsoleID)
					cm.connectionMap.Range(func(key, value interface{}) bool {
						if conn, ok := value.(*Connection); ok {
							if managerMessage.RecipientConsoleID == conn.ID {
								cm.Log.Println(logging.LogInfo, "Sending character to console: ", conn.ID)
								consoleMessage := messages.ConsoleMessage{
									Type: messages.Console_Message_CharacterRequestResponse,
									Data: managerMessage.Data,
								}
								conn.SendToConsole(consoleMessage)
							}
						}
						return true
					})
				}()

			case messages.ConnectManager_Message_RequestPasswordReset:
				// Sending request to EdenBot to reset password for a user
				// If this doesn't work, we don't tell the user
				go func() {
					cm.Log.Println(logging.LogInfo, "Sending password reset request to EdenBot")
					cm.EBSendMessages <- messages.EdenbotMessage{
						Type:       messages.Edenbot_Message_ForgotPassword,
						Data:       managerMessage.Data.(messages.AccountForgotPasswordData),
						SourceType: "console",
						SourceID:   managerMessage.SenderConsoleID,
					}
				}()

			case messages.ConnectManager_Message_CharNameValidation:
				go func() {
					cm.Log.Println(logging.LogInfo, "Sending character name validation request to CharacterManager")
					cm.CMSendMessages <- messages.CharacterManagerMessage{
						Type:            messages.CharManager_CheckName,
						Data:            managerMessage.Data,
						SenderConsoleID: managerMessage.SenderConsoleID,
					}
				}()

			case messages.ConnectManager_Message_CharacterCreation:
				go func() {
					cm.Log.Println(logging.LogInfo, "Sending character creation request to CharacterManager")
					cm.CMSendMessages <- messages.CharacterManagerMessage{
						Type:            messages.CharManager_CreateCharacter,
						Data:            managerMessage.Data,
						SenderConsoleID: managerMessage.SenderConsoleID,
					}
				}()

			case messages.ConnectManager_Message_GameCommand:
				go func() {
					cm.Log.Println(logging.LogInfo, "Sending game command to game manager from", managerMessage.SenderConsoleID)
					cm.GMSendMessages <- messages.GameManagerMessage{
						Type:            managerMessage.Data.(messages.GameManagerMessage).Type,
						Data:            managerMessage.Data,
						SenderConsoleID: managerMessage.SenderConsoleID,
					}
				}()

			case messages.ConnectManager_Message_GameCommandResponse:
				go func() {
					cm.Log.Println(logging.LogInfo, "Sending game manager response to client console", managerMessage.SenderConsoleID)
					cm.connectionMap.Range(func(key, value interface{}) bool {
						if conn, ok := value.(*Connection); ok {
							if managerMessage.RecipientConsoleID == conn.ID {
								//cm.Log.Println(logging.LogInfo, "Quit message found, sending to conn.Console.ReceiveMessages")
								consoleMessage := messages.ConsoleMessage{
									Type: messages.Console_Message_GameCommandResponse,
									Data: managerMessage.Data,
								}
								conn.SendToConsole(consoleMessage)
							}
						}
						return true
					})
				}()

			case messages.ConnectManager_Message_CharacterCreationResponse:
				go func() {
					cm.Log.Println(logging.LogInfo, "Sending character name validation response to Console")
					cm.connectionMap.Range(func(key, value interface{}) bool {
						if conn, ok := value.(*Connection); ok {
							if managerMessage.RecipientConsoleID == conn.ID {
								//cm.Log.Println(logging.LogInfo, "Quit message found, sending to conn.Console.ReceiveMessages")
								consoleMessage := messages.ConsoleMessage{
									Type: messages.Console_Message_CharacterCreationResponse,
									Data: managerMessage.Data,
								}
								conn.SendToConsole(consoleMessage)
							}
						}
						return true
					})
				}()

			case messages.ConnectManager_Message_CharNameValidationResponse:
				go func() {
					cm.Log.Println(logging.LogInfo, "Sending character name validation response to Console")
					cm.connectionMap.Range(func(key, value interface{}) bool {
						if conn, ok := value.(*Connection); ok {
							if managerMessage.RecipientConsoleID == conn.ID {
								//cm.Log.Println(logging.LogInfo, "Quit message found, sending to conn.Console.ReceiveMessages")
								consoleMessage := messages.ConsoleMessage{
									Type: messages.Console_Message_ValidateCharNameResponse,
									Data: managerMessage.Data,
								}
								conn.SendToConsole(consoleMessage)
							}
						}
						return true
					})
				}()

			case messages.ConnectManager_Message_ValidatePasswordReset:
				// Send the password reset validation request to the AccountManager
				go func() {
					//cm.Log.Println(logging.LogInfo, "Sending password reset validation request to AccountManager")
					cm.AMSendMessages <- messages.AccountManagerMessage{Type: messages.AccountManager_Message_ResetPasswordValidate, Data: managerMessage.Data.(messages.AccountProcessForgotPasswordData), SenderConsoleID: managerMessage.SenderConsoleID}
				}()

			case messages.ConnectManager_Message_ProcessPasswordReset:
				// Send the new password data to the AccountManager
				go func() {
					cm.Log.Println(logging.LogInfo, "Sending new password request to AccountManager")
					cm.AMSendMessages <- messages.AccountManagerMessage{Type: messages.AccountManager_Message_ResetPasswordProcess, Data: managerMessage.Data.(messages.AccountProcessForgotPasswordData), SenderConsoleID: managerMessage.SenderConsoleID}
				}()

			case messages.ConnectManager_Message_ValidatePasswordResetResponse:
				// Send the password reset validation response to the Console
				cm.connectionMap.Range(func(key, value interface{}) bool {
					if conn, ok := value.(*Connection); ok {
						if managerMessage.RecipientConsoleID == conn.ID {
							//cm.Log.Println(logging.LogInfo, "Sending registration response to Console that requested registration")
							consoleMessage := messages.ConsoleMessage{Type: messages.Console_Message_ResetPasswordValidateResponse, Data: managerMessage.Data}
							conn.SendToConsole(consoleMessage)
						}
					}
					return true
				})

			case messages.ConnectManager_Message_ProcessPasswordResetResponse:
				// Send the password reset validation response to the Console
				cm.connectionMap.Range(func(key, value interface{}) bool {
					if conn, ok := value.(*Connection); ok {
						if managerMessage.RecipientConsoleID == conn.ID {
							//cm.Log.Println(logging.LogInfo, "Sending registration response to Console that requested registration")
							consoleMessage := messages.ConsoleMessage{Type: messages.Console_Message_ProcessPasswordValidateResponse, Data: managerMessage.Data}
							conn.SendToConsole(consoleMessage)
						}
					}
					return true
				})

			case messages.ConnectManager_Message_RegisterResponse:
				cm.connectionMap.Range(func(key, value interface{}) bool {
					if conn, ok := value.(*Connection); ok {
						if managerMessage.RecipientConsoleID == conn.ID {
							//cm.Log.Println(logging.LogInfo, "Sending registration response to Console that requested registration")
							consoleMessage := messages.ConsoleMessage{Type: messages.Console_Message_RegistrationResponse, Data: managerMessage.Data}
							conn.SendToConsole(consoleMessage)
						}
					}
					return true
				})

			case messages.ConnectManager_Message_LoginResponse:
				cm.connectionMap.Range(func(key, value interface{}) bool {
					if conn, ok := value.(*Connection); ok {
						if managerMessage.RecipientConsoleID == conn.ID {
							//cm.Log.Println(logging.LogInfo, "Sending login response to Console that requested login")
							consoleMessage := messages.ConsoleMessage{Type: messages.Console_Message_LoginResponse, Data: managerMessage.Data}
							conn.SendToConsole(consoleMessage)
						}
					}
					return true
				})

			case messages.ConnectManager_Message_ForceLogout:
				cm.connectionMap.Range(func(key, value interface{}) bool {
					if conn, ok := value.(*Connection); ok {
						if managerMessage.RecipientConsoleID == conn.ID {
							cm.Log.Println(logging.LogInfo, "Sending logout to existing connection")
							consoleMessage := messages.ConsoleMessage{Type: messages.Console_Message_LogoutUser, Data: managerMessage.Data}
							conn.SendToConsole(consoleMessage)
						}
					}
					return true
				})

			case messages.ConnectManager_Message_BadLoginAttempt:
				cm.connectionMap.Range(func(key, value interface{}) bool {
					if conn, ok := value.(*Connection); ok {
						// Note the distinction here, the SenderConsoleID is the ID of the console that sent the message,
						if managerMessage.SenderConsoleID == conn.ID {
							if cm.handleBadLogins(conn) {
								cm.HandleDisconnect(conn)
								conn.Write([]byte("\033c\r\nToo many bad login attempts, your IP has been temporarily blocked.\r\n"))
								conn.Close()
							}
						}
					}
					return true
				})

			case messages.ConnectManager_Message_UpdateCharacterHistoryResponse:

				cm.connectionMap.Range(func(key, value interface{}) bool {
					if conn, ok := value.(*Connection); ok {
						if managerMessage.RecipientConsoleID == conn.ID {
							cm.Log.Println(logging.LogInfo, "Sending character history update response to console")
							//cm.Log.Println(logging.LogInfo, "Sending login response to Console that requested login")
							consoleMessage := messages.ConsoleMessage{Type: messages.Console_Message_CharacterHistoryCharacterUpdateResponse, Data: managerMessage.Data}
							conn.SendToConsole(consoleMessage)
						}
					}
					return true
				})

			case messages.ConnectManager_Message_UpdateAccountHistoryResponse:
				cm.connectionMap.Range(func(key, value interface{}) bool {
					if conn, ok := value.(*Connection); ok {
						if managerMessage.RecipientConsoleID == conn.ID {
							cm.Log.Println(logging.LogInfo, "Sending account history update response to console")
							//cm.Log.Println(logging.LogInfo, "Sending login response to Console that requested login")
							consoleMessage := messages.ConsoleMessage{Type: messages.Console_Message_CharacterHistoryAccountUpdateResponse, Data: managerMessage.Data}
							conn.SendToConsole(consoleMessage)
						}
					}
					return true
				})

			default:
				cm.Log.Println(logging.LogError, "Unknown message type received: ", managerMessage.Type, managerMessage.SenderConsoleID, managerMessage.RecipientConsoleID)
			}

		}
	}
}

type BadLoginRecord struct {
	ID        string `storm:"index"`
	IPAddress string `storm:"unique"`
	Timestamp time.Time
	Attempts  int
}

func (cm *ConnectionManager) handleBadLogins(conn *Connection) (disconnect bool) {
	// First we get the IP address of the connection
	ipAddress, _, _ := net.SplitHostPort(conn.conn.RemoteAddr().String())

	// Then we check if we have a record for this IP address
	var badLoginRecord BadLoginRecord

	err := cm.DB.One("bad_logins", "IPAddress", ipAddress, &badLoginRecord)
	if err != nil {
		// If we don't have a record, we create one
		id := uuid.New().String()
		badLoginRecord = BadLoginRecord{ID: id, IPAddress: ipAddress, Timestamp: time.Now(), Attempts: 1}
		err = cm.DB.AddRecord("bad_logins", &badLoginRecord)
		if err != nil {
			cm.Log.Println(logging.LogError, "Error adding bad login record: ", err)
		}
		return false
	} else {
		// If we do have a record, we increment the number of attempts
		badLoginRecord.Attempts++
		badLoginRecord.Timestamp = time.Now()
		err = cm.DB.UpdateRecord("bad_logins", &badLoginRecord)
		if err != nil {
			cm.Log.Println(logging.LogError, "Error updating bad login record: ", err)
		}

		// If we have more than 3 bad login attempts recorded, we disconnect the connection
		if badLoginRecord.Attempts > 3 {
			return true
		}
	}
	return false
}

func (cm *ConnectionManager) checkIPBadLogins(conn *Connection) (allowed bool) {
	// First we get the IP address of the connection
	ipAddress, _, _ := net.SplitHostPort(conn.conn.RemoteAddr().String())

	// Then we check if we have a record for this IP address
	var badLoginRecord BadLoginRecord

	err := cm.DB.One("bad_logins", "IPAddress", ipAddress, &badLoginRecord)
	if err != nil {
		// If we don't have a record, we allow the connection
		return true
	} else {
		// If we do have a record, we check the number of bad logins
		if badLoginRecord.Attempts > 3 {
			// If less than 30 minutes have passed since the last bad login, we don't allow the connection
			if time.Since(badLoginRecord.Timestamp) < time.Minute*30 {
				return false
			} else {
				// If more than 30 minutes have passed, we remove the bad login record and allow the connection
				err = cm.DB.RemoveRecord("bad_logins", &badLoginRecord)
				if err != nil {
					cm.Log.Println(logging.LogError, "Error updating bad login record: ", err)
				}
				return true
			}
		}
	}

	return true
}
