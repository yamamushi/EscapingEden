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
	ID      string
	conn    net.Conn
	mutex   sync.Mutex
	console *ui.Console
	manager *ConnectionManager

	// These are for working with IAC commands coming into the handler
	iacBuffer              string
	iacActive              bool
	iacSubnegotationActive bool
	iacResizeActive        bool
	iacParamIndex          int
	iacWindowResizeX       int
	iacWindowResizeY       int
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
	err = EnableLineMode(c.conn)
	err = DisableEcho(c.conn) // Flush our IAC codes, we don't care about the responses (for now)
	if err != nil {
		log.Println("Client unsupported, disconnecting: ", err)
		c.manager.HandleDisconnect(c)
		return
	}
	c.conn.Write([]byte("\033[2J"))
	c.conn.Write([]byte("\033[?25l"))

	c.console = ui.NewConsole(w, h, c.ID, c.manager.CMReceiveMessages)
	log.Println("Initializing Console")
	c.console.Init()
	log.Println("Console Initialized")
	time.Sleep(time.Second * 1)

	reader := bufio.NewReader(c.conn)

	// Enter our client loop
	for {
		readByte, err := reader.ReadByte()
		if err != nil {
			log.Println("Client closed connection")
			c.manager.HandleDisconnect(c)
			return
		}

		if readByte == IAC {
			log.Println("IAC received")
			c.iacActive = true
			continue
		}
		if c.iacActive {
			c.HandleIAC(readByte)
			continue
		}

		c.console.HandleInput(readByte)
	}
}

func (c *Connection) HandleIAC(readByte byte) {
	log.Println("Byte: ", readByte)

	if c.iacActive {
		if readByte == 0 {
			// We don't care about null bytes, just ignore them
			// And a window resize shouldn't be sending a null byte like this
			// We should probably handle this better later
			// But for now, just ignore it
			return
		}
		if readByte == 250 {
			log.Println("IAC subnegotiation received")
			c.iacSubnegotationActive = true
			return
		}
		if c.iacSubnegotationActive {
			if readByte == 31 && !c.iacResizeActive {
				log.Println("IAC Window resize received")
				c.iacResizeActive = true
				return
			}
			if c.iacResizeActive {
				if c.iacParamIndex == 0 {
					c.iacParamIndex = 1
					c.iacWindowResizeX = int(readByte)
					return
				}
				if c.iacParamIndex == 1 {
					c.iacParamIndex = 0
					c.iacWindowResizeY = int(readByte)
					c.iacResizeActive = false
					log.Println("IAC Window resize received: ", c.iacWindowResizeX, c.iacWindowResizeY)
					return
				}
			}
			if readByte == 240 {
				// Only when we receive the final SE bit are we done with the subnegotiation
				c.iacSubnegotationActive = false
				c.iacActive = false
				c.iacResizeActive = false
				c.iacParamIndex = 0
				//if c.iacWindowResizeX
				if c.iacWindowResizeX > 0 && c.iacWindowResizeY > 0 {
					c.console.HandleResize(c.iacWindowResizeX, c.iacWindowResizeY)
				}
				c.iacWindowResizeX = 0
				c.iacWindowResizeY = 0
				log.Println("IAC subnegotiation complete")

				return
			}
		}

		log.Println("Unhandled IAC command received: ", int(readByte))
		log.Println("If you're reading this, you're about to see garbage get sent to your windows.")
		c.iacActive = false
	}
}
