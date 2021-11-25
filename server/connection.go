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

	if c.console.GetShutdown() {
		log.Println("Client requested shutdown")
		c.conn.Write([]byte("Goodbye!\n"))
		c.conn.Close()
		c.manager.HandleDisconnect(c)
		return
	}

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
	c.console = ui.NewConsole(w, h, c.ID, c.manager.CMReceiveMessages)
	log.Println("Initializing Console")
	c.console.Init()
	log.Println("Console Initialized")
	time.Sleep(time.Second * 1)

	// Enter our client loop
	for {

		var buff []byte
		reader := bufio.NewReader(c.conn)
		for {
			readByte, err := reader.ReadByte()
			if err != nil {
				log.Println("Client closed connection")
				c.manager.HandleDisconnect(c)
				return
			}
			if readByte == '\n' {
				if len(buff) == 0 {
					buff = append(buff, readByte)
					continue
				}
				break
			}
			if readByte == '\b' {
				if len(buff) == 0 {
					continue
				}
				c.console.SetBackspaceReceived()
			}
			buff = append(buff, readByte)
		}

		//userInput, err := bufio.NewReader(c.conn).ReadString('\n')
		userInput := string(buff)
		log.Println("User Input Buffer received: " + userInput)
		if err != nil {
			log.Println("Client closed connection")
			c.manager.HandleDisconnect(c)
			return
		}
		c.console.HandleInput(userInput)
		c.Write(c.outputBuffer)
		if c.console.GetShutdown() {
			log.Println("Client requested shutdown")
			c.conn.Write([]byte("Goodbye!\n"))
			c.conn.Close()
			c.manager.HandleDisconnect(c)
			return
		}
		log.Println("Finished Writing in client loop")

	}
}
