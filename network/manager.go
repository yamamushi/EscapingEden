package network

import (
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"sync"
)

// ConnectionManager synchronizes connection output globally
type ConnectionManager struct {
	Log logging.LoggerType
	// Mutex for locking
	mutex         sync.Mutex
	connectionMap *sync.Map

	// Channel for receiving messages
	CMReceiveMessages chan messages.ConnectionManagerMessage

	// Our AccountManager
	//AccountManager *accounts.AccountManager
	// Account manager outbound channel
	AMSendMessages chan messages.AccountManagerMessage

	// Our CharacterManager outbound Channel
	CharacterManagerMessages chan messages.CharacterManagerMessage
}

// NewConnectionManager creates a new ConnectionManager
func NewConnectionManager(connectionMap *sync.Map, receiveMessages chan messages.ConnectionManagerMessage,
	accountManagerMessages chan messages.AccountManagerMessage, characterManagerReceiveMessages chan messages.CharacterManagerMessage, log logging.LoggerType) *ConnectionManager {
	return &ConnectionManager{
		connectionMap:            connectionMap,
		CMReceiveMessages:        receiveMessages,
		AMSendMessages:           accountManagerMessages,
		CharacterManagerMessages: characterManagerReceiveMessages,
		Log:                      log,
	}
}

// AddConnection adds a connection to the manager
func (cm *ConnectionManager) AddConnection(connection *Connection) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.connectionMap.Store(connection.ID, connection)
}

// HandleDisconnect handles disconnect events
func (cm *ConnectionManager) HandleDisconnect(connection *Connection) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.connectionMap.Delete(connection.ID)
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
						outMessage := messages.ConsoleMessage{Data: managerMessage.SenderConsoleID + ": " + managerMessage.Data.(string), Type: messages.Console_Message_Chat}
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
					cm.AMSendMessages <- messages.AccountManagerMessage{Type: messages.AccountManager_Message_Register, Data: registrationRequest, SenderSessionID: managerMessage.SenderConsoleID}
				}()

			case messages.ConnectManager_Message_Login:
				go func() {
					//cm.Log.Println(logging.LogInfo, "Sending login request to AccountManager")
					loginRequest := managerMessage.Data.(messages.AccountLoginRequest)
					cm.AMSendMessages <- messages.AccountManagerMessage{Type: messages.AccountManager_Message_Login, Data: loginRequest, SenderSessionID: managerMessage.SenderConsoleID}
				}()

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

			default:
				cm.Log.Println(logging.LogError, "Unknown message type received: ", managerMessage.Type, managerMessage.SenderConsoleID, managerMessage.RecipientConsoleID)
			}

		}
	}
}
