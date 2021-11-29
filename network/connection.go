package network

import (
	"bufio"
	"github.com/yamamushi/EscapingEden/ui"
	"log"
	"net"
	"sync"
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

func (c *Connection) Write(msg []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, err := c.conn.Write(msg)
	if err != nil {
		//log.Println(err)
		return err
	}
	return nil
}

// handleConnection handles a single connection
func (c *Connection) Handle() {
	/*
		We need to setup the session, and to do that we need to do some communication with the client.
	*/
	log.Println("Setting up new Connection for:", c.ID)

	log.Println("Requesting terminal type for:", c.ID)
	termType, err := RequestTerminalType(c.conn)
	if err != nil {
		log.Println("Error requesting terminal type:", c.ID, " Message: ", err, "Closing connection")
		c.manager.HandleDisconnect(c)
		return
	}
	if termType == "xterm-256color" {
		log.Println("Unsupported terminal type:", c.ID, " TermType:", termType, " Closing connection")
		c.manager.HandleDisconnect(c)
		return
	}

	log.Println("Requesting terminal size for:", c.ID)
	width, height, err := RequestTerminalSize(c.conn)
	if err != nil {
		log.Println("Terminal Size Request Error:", c.ID, " Message: ", err, " Closing connection")
		c.manager.HandleDisconnect(c)
		return
	}

	log.Println("Enabling LineMode on client:", c.ID)
	err = EnableLineMode(c.conn)
	if err != nil {
		log.Println("LineMode Enable Error:", c.ID, " Message: ", err, " Closing connection")
		c.manager.HandleDisconnect(c)
		return
	}

	log.Println("Disabling Echo on client:", c.ID)
	err = DisableEcho(c.conn) // Flush our IAC codes, we don't care about the responses (for now)
	if err != nil {
		log.Println("Disable Echo Error:", c.ID, " Message: ", err, " Closing connection")
		c.manager.HandleDisconnect(c)
		return
	}

	log.Println("Flushing client terminal and hiding cursor:", c.ID)
	_, err = c.conn.Write([]byte("\033[2J"))
	if err != nil {
		log.Println("Flush Error:", c.ID, " Message: ", err, " Closing connection")
		c.manager.HandleDisconnect(c)
		return
	}
	_, err = c.conn.Write([]byte("\033[?25l"))
	if err != nil {
		log.Println("Hide Cursor Error:", c.ID, " Message: ", err, " Closing connection")
		c.manager.HandleDisconnect(c)
		return
	}

	c.console = ui.NewConsole(width, height, c.ID, c.manager.CMReceiveMessages)
	log.Println("Initializing Console for:", c.ID)
	c.console.Init()
	log.Println("Console Initialized for:", c.ID)

	log.Println("Launching Write Handler for:", c.ID)
	go c.WriteHandler()
	log.Println("Launching Read Handler for:", c.ID)
	go c.ReadHandler()
	log.Println("Connection successfully created for:", c.ID)
}

func (c *Connection) ReadHandler() {
	// Enter our client loop
	reader := bufio.NewReader(c.conn)
	for {
		readByte, err := reader.ReadByte()
		if err != nil {
			log.Println("Client ", c.ID, " closed connection")
			c.manager.HandleDisconnect(c)
			return
		}

		if readByte == IAC {
			//log.Println("IAC received")
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

func (c *Connection) WriteHandler() {
	for {
		if c.console.GetShutdown() {
			log.Println("Client requested shutdown")
			c.conn.Write([]byte("\033[2J"))
			c.conn.Write([]byte("\033[;H" + "See you back soon! Goodbye :)\r\n"))
			c.conn.Close()
			c.manager.HandleDisconnect(c)
			return
		}

		output := c.console.Draw()
		if len(output) > 0 {
			err := c.Write(output)
			if err != nil {
				log.Println("Client disconnected")
				c.conn.Close()
				c.manager.HandleDisconnect(c)
				return
			}
		}
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
			//log.Println("IAC subnegotiation received")
			c.iacSubnegotationActive = true
			return
		}
		if c.iacSubnegotationActive {
			if readByte == 31 && !c.iacResizeActive {
				//log.Println("IAC Window resize received")
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
					//log.Println("IAC Window resize received: ", c.iacWindowResizeX, c.iacWindowResizeY)
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
				//log.Println("IAC subnegotiation complete")

				return
			}
		}

		log.Println("!!! If you're reading this, IAC parsing failed and "+
			"you're about to see garbage get sent to your windows. Unhandled IAC command: ", int(readByte))
		c.iacActive = false
	}
}
