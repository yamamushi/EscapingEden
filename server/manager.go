package server

import (
	"sync"
)

// Manager synchronizes connection output globally
type ConnectionManager struct {
	// Mutex for locking
	mutex         sync.Mutex
	connectionMap *sync.Map
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
	for {
		cm.connectionMap.Range(func(key, value interface{}) bool {
			// Send the message "Boo!" to the client
			if conn, ok := value.(*Connection); ok {
				// convert string to byte array
				if conn.console != nil {
					output := conn.console.Draw()
					if string(output) != "" {
						conn.Write(output)
					}
				}
			}
			// do something with key and/or value
			return true // return true to continue iterating
		})

		// Sleep a random amount of seconds between 1 and 5
		//time.Sleep(time.Duration(rand.Intn(5-1)+1) * time.Second)
	}
}
