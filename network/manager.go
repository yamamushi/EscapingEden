package network

import (
	"github.com/yamamushi/EscapingEden/messages"
	"log"
	"sync"
)

// ConnectionManager synchronizes connection output globally
type ConnectionManager struct {
	// Mutex for locking
	mutex         sync.Mutex
	connectionMap *sync.Map

	startedNotify chan bool

	// Channel for receiving messages
	CMReceiveMessages chan messages.ConnectionManagerMessage
}

// NewConnectionManager creates a new ConnectionManager
func NewConnectionManager(connectionMap *sync.Map, startedNotify chan bool) *ConnectionManager {

	return &ConnectionManager{
		connectionMap: connectionMap,
		startedNotify: startedNotify,
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
func (cm *ConnectionManager) Run() {
	go cm.MessageParser() // Launch our goroutine that listens for incoming messages
}

// MessageParser performs a non-blocking check for messages on cm.CMReceiveMessages
func (cm *ConnectionManager) MessageParser() {
	log.Println("ConnectionManager is now listening for incoming messages")
	cm.startedNotify <- true

	cm.CMReceiveMessages = make(chan messages.ConnectionManagerMessage)
	for {
		select {
		case managerMessage := <-cm.CMReceiveMessages:
			log.Println("Message received from cm.CMReceiveMessages")

			switch managerMessage.Type {
			case messages.ConnectManager_Message_Chat:
				// For every connection, send the message to the console channel
				cm.connectionMap.Range(func(key, value interface{}) bool {
					if conn, ok := value.(*Connection); ok {
						outMessage := messages.ConsoleMessage{Message: managerMessage.SenderConsoleID + ": " + managerMessage.Message, Type: messages.Console_Message_Chat}
						log.Println("Chat message found, sending to conn.console.ReceiveMessages")
						conn.SendToConsole(outMessage)
					}
					return true
				})
			case messages.ConnectManager_Message_Broadcast:
				// For every connection, send the message to the console channel
				cm.connectionMap.Range(func(key, value interface{}) bool {
					if _, ok := value.(*Connection); ok {
						// json marshal message to string
						log.Println("Broadcast message found, sending to conn.console.ReceiveMessages")
						//conn.console.ReceiveMessages <- output
					}
					return true
				})
			case messages.ConnectManager_Message_Error:
				// For every connection, send the message to the console channel
				cm.connectionMap.Range(func(key, value interface{}) bool {
					if conn, ok := value.(*Connection); ok {
						if managerMessage.RecipientConsoleID == conn.ID {
							log.Println("Error message found, sending to conn.console.ReceiveMessages")
							consoleMessage := messages.ConsoleMessage{Type: messages.Console_Message_Error, Message: managerMessage.Message}
							conn.SendToConsole(consoleMessage)

						}
					}
					return true
				})
			case messages.ConnectManager_Message_Quit:
				cm.connectionMap.Range(func(key, value interface{}) bool {
					if conn, ok := value.(*Connection); ok {
						if managerMessage.RecipientConsoleID == conn.ID {
							log.Println("Quit message found, sending to conn.console.ReceiveMessages")
							consoleMessage := messages.ConsoleMessage{Type: messages.Console_Message_Quit}
							conn.SendToConsole(consoleMessage)
						}
					}
					return true
				})
			default:
				log.Println("Unknown message type received: ", managerMessage.Type, managerMessage.SenderConsoleID, managerMessage.RecipientConsoleID)
			}

		}
	}
}
