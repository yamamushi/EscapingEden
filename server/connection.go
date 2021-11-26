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
	log.Println("Enabling LineMode on client")
	EnableLineMode(c.conn)
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
			// if byte is return character
			if readByte == '\r' {
				if len(buff) == 0 {
					buff = append(buff, readByte)
					continue
				}
				break
			}
			if readByte == '\n' {
				if len(buff) == 0 {
					buff = append(buff, readByte)
					continue
				}
				break
			}
			// if byte is a backspace sequence
			if readByte == '\b' {
				if len(buff) > 0 {
					buff = buff[:len(buff)-1]
				}
				c.console.SetBackspaceReceived(1)
				continue
			}
			// if control key is pressed log it
			if readByte == '\x1b' {
				log.Println("Control key pressed")
			}
			if readByte == '\x7f' {
				c.console.SetBackspaceReceived(3)
				continue
			}
			// if ascll code is a control character, log it
			if readByte < 32 {
				log.Println("Control character received:", readByte)
			}

			// If up arrow pressed, move cursor up and erase
			if readByte == '\033' {
				log.Println("Read escape sequence")
				// Read the next byte
				readByte, err = reader.ReadByte()
				if err != nil {
					log.Println("Client closed connection")
					c.manager.HandleDisconnect(c)
					return
				}
				log.Println("Read byte after escape: ", string(readByte))
				if readByte == '[' {
					// Read the next byte
					readByte, err = reader.ReadByte()
					if err != nil {
						log.Println("Client closed connection")
						c.manager.HandleDisconnect(c)
						return
					}
					if readByte == 'A' {
						log.Println("Up arrow pressed")
						c.console.HandleMovement("up")
					}
					if readByte == 'B' {
						log.Println("Down arrow pressed")
						c.console.HandleMovement("down")
					}
					if readByte == 'C' {
						log.Println("Right arrow pressed")
						c.console.HandleMovement("right")
					}
					if readByte == 'D' {
						log.Println("Left arrow pressed")
						c.console.HandleMovement("left")
					}
				}
				c.console.SetBackspaceReceived(3)
			}
			buff = append(buff, readByte)
			log.Println("Read byte: ", string(readByte))
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
