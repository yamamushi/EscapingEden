package server

import (
	"encoding/json"
	"log"
	"sync"
)

// Manager synchronizes connection output globally
type ConnectionManager struct {
	// Mutex for locking
	mutex         sync.Mutex
	connectionMap *sync.Map

	// Channel for receiving messages
	CMReceiveMessages chan string
}

func NewConnectionManager(connectionMap *sync.Map) *ConnectionManager {

	return &ConnectionManager{
		connectionMap: connectionMap,
	}
}

// Adds a connection to the manager
func (cm *ConnectionManager) AddConnection(connection *Connection) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.connectionMap.Store(connection.ID, connection)
}

// Handle Disconnects
func (cm *ConnectionManager) HandleDisconnect(connection *Connection) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.connectionMap.Delete(connection.ID)
}

func (cm *ConnectionManager) Run() {
	go cm.MessageParser() // Launch our goroutine that listens for incoming messages
}

// Non blocking check for messages
func (cm *ConnectionManager) MessageParser() {
	log.Println("ConnectionManager Listening for incoming messages")
	cm.CMReceiveMessages = make(chan string)
	for {
		select {
		case message := <-cm.CMReceiveMessages:
			log.Println("Message received from cm.CMReceiveMessages")

			managerMessage := &ManagerMessage{}
			err := json.Unmarshal([]byte(message), managerMessage)
			if err != nil {
				continue
			}
			log.Println("Message received from manager: ", managerMessage.Message)

			switch managerMessage.Type {
			case "chat":
				// For every connection, send the message to the console channel
				cm.connectionMap.Range(func(key, value interface{}) bool {
					if conn, ok := value.(*Connection); ok {
						log.Println(conn.ID)
						log.Println(managerMessage.SenderID)
						log.Println(managerMessage.Message)
						outMessage := &ManagerMessage{Message: managerMessage.SenderID + ": " + managerMessage.Message, Type: "chat"}
						// json marshal message to string
						output, err := json.Marshal(outMessage)
						if err == nil {
							log.Println("Chat message found, sending to conn.console.ReceiveMessages")
							conn.console.ReceiveMessages <- string(output)
						}
					}
					return true
				})
			case "broadcast":
				// For every connection, send the message to the console channel
				cm.connectionMap.Range(func(key, value interface{}) bool {
					if conn, ok := value.(*Connection); ok {
						// json marshal message to string
						output, err := json.Marshal(managerMessage)
						if err == nil {
							log.Println("Broadcast message found, sending to conn.console.ReceiveMessages")
							conn.console.ReceiveMessages <- string(output)
						}
					}
					return true
				})
			case "error":
				// For every connection, send the message to the console channel
				cm.connectionMap.Range(func(key, value interface{}) bool {
					if conn, ok := value.(*Connection); ok {
						if managerMessage.RecipientID == conn.ID {
							// json marshal message to string
							output, err := json.Marshal(managerMessage)
							if err == nil {
								log.Println("Error message found, sending to conn.console.ReceiveMessages")
								conn.console.ReceiveMessages <- string(output)
							}
						}
					}
					return true
				})
			case "quit":
				cm.connectionMap.Range(func(key, value interface{}) bool {
					if conn, ok := value.(*Connection); ok {
						if managerMessage.RecipientID == conn.ID {
							// json marshal message to string
							output, err := json.Marshal(managerMessage)
							if err == nil {
								log.Println("Quit message found, sending to conn.console.ReceiveMessages")
								conn.console.ReceiveMessages <- string(output)
							}
						}
					}
					return true
				})
			default:
				log.Println("Unknown message type received: ", managerMessage.Type, managerMessage.SenderID, managerMessage.RecipientID)
			}

		}
	}
}
