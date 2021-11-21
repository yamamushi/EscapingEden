package server

import (
	"bufio"
	"github.com/yamamushi/EscapingEden/ui"
	"log"
	"net"
	"sync"
	"time"
)

// Connection is a connection to a client in case we need to store any extra details later
type Connection struct {
	ID           string
	conn         net.Conn
	mutex        sync.Mutex
	outputBuffer []byte
	console      *ui.Console
	manager      *ConnectionManager
}

// NewConnection creates a new connection
func NewConnection(conn net.Conn, id string, manager *ConnectionManager) *Connection {
	connection := &Connection{
		conn:    conn,
		ID:      id,
		manager: manager,
	}
	go connection.Handle()
	return connection
}

func (c *Connection) Write(msg []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, err := c.conn.Write(msg)
	if err != nil {
		log.Println(err)
	}
}

// handleConnection handles a single connection
func (c *Connection) Handle() {
	/*
		We need to setup the session, and to do that we need to do some communication with the client.
	*/
	//RequestTerminalType(conn)
	w, h, err := RequestTerminalSize(c.conn)
	if err != nil {
		log.Println("Client unsupported, disconnecting")
		c.manager.HandleDisconnect(c)
		return
	}
	c.console = ui.NewConsole(w, h, c.manager.CMReceiveMessages)
	log.Println("Initializing Console")
	c.console.Init()
	log.Println("Console Initialized")
	time.Sleep(time.Second * 1)

	// Enter our client loop
	for {
		userInput, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			log.Println("Client closed connection")
			c.manager.HandleDisconnect(c)
			return
		}
		c.console.HandleInput(userInput)
		log.Println("Continuing client loop")

		if c.console.GetShutdown() {
			log.Println("Client requested shutdown")
			c.conn.Write([]byte("Goodbye!\n"))
			c.conn.Close()
			c.manager.HandleDisconnect(c)
			return
		}
		log.Println("Writing in client loop")

		c.Write(c.outputBuffer)
		log.Println("Finished Writing in client loop")

	}
}
