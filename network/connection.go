package network

import (
	"bufio"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"github.com/yamamushi/EscapingEden/terminals"
	xterm_256color "github.com/yamamushi/EscapingEden/terminals/xterm-256color"
	"github.com/yamamushi/EscapingEden/ui"
	"net"
	"strings"
	"sync"
)

// Connection is a connection to a client in case we need to store any extra details later
type Connection struct {
	ID      string
	conn    net.Conn
	mutex   sync.Mutex
	Console *ui.Console
	manager *ConnectionManager

	Log logging.LoggerType

	// These are for working with IAC commands coming into the handler
	iacBuffer               string
	iacActive               bool
	iacSubnegotationActive  bool
	iacResizeActive         bool
	cleanupAfterResize      bool
	cleanupStage            int
	cleanupAfterResizeMutex sync.Mutex
	iacParamIndex           int
	iacWindowResizeX        int
	iacWindowResizeY        int
}

// NewConnection creates a new connection
func NewConnection(conn net.Conn, id string, manager *ConnectionManager, log logging.LoggerType) *Connection {
	connection := &Connection{
		conn:    conn,
		ID:      id,
		manager: manager,
		Log:     log,
	}
	go connection.Handle()
	return connection
}

// Write writes a byte to the connection
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

func (c *Connection) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.conn.Close()
}

// Handle handles a single connection
func (c *Connection) Handle() {
	/*
		We need to setup the session, and to do that we need to do some communication with the client.
	*/
	//c.Log.Println(logging.LogInfo, "Setting up new Connection for:", c.ID)

	//log.Println("Requesting terminal type for:", c.ID)
	termType, err := RequestTerminalType(c.conn)
	if err != nil {
		c.Log.Println(logging.LogWarn, "Error requesting terminal type:", c.ID, " Message: ", err, "Closing connection")
		c.manager.HandleDisconnect(c)
		return
	}
	termType = strings.ToLower(termType)
	var termTypeID terminals.TermTypeID
	if termType == "xterm-256color" {
		termTypeID = terminals.TermTypeXTerm256Color
	} else {
		c.Log.Println(logging.LogWarn, "Unsupported terminal type:", c.ID, " Terminal:", termType, " Closing connection")
		c.manager.HandleDisconnect(c)
		// This is one of the FEW places we use \033c
		c.Write([]byte("\033cUnsupported terminal type, sorry only xterm-256color is supported at the " +
			"moment and you're using " + termType + " >_>\r\n"))
		c.Close()
		return
	}
	// Set up the terminal type
	var terminal terminals.TerminalType
	switch termTypeID {
	case terminals.TermTypeXTerm256Color:
		terminal = xterm_256color.New()
	}

	//log.Println("Requesting terminal size for:", c.ID)
	width, height, err := RequestTerminalSize(c.conn)
	if err != nil {
		c.Log.Println(logging.LogWarn, "Terminal Size Request Error:", c.ID, " Message: ", err, " Closing connection")
		c.manager.HandleDisconnect(c)
		return
	}

	//log.Println("Enabling LineMode on client:", c.ID)
	err = EnableLineMode(c.conn)
	if err != nil {
		c.Log.Println(logging.LogWarn, "LineMode Enable Error:", c.ID, " Message: ", err, " Closing connection")
		c.manager.HandleDisconnect(c)
		return
	}

	//log.Println("Disabling Echo on client:", c.ID)
	err = DisableEcho(c.conn) // Flush our IAC codes, we don't care about the responses (for now)
	if err != nil {
		c.Log.Println(logging.LogWarn, "Disable Echo Error:", c.ID, " Message: ", err, " Closing connection")
		c.manager.HandleDisconnect(c)
		return
	}

	//log.Println("Flushing client terminal and hiding cursor:", c.ID)
	_, err = c.conn.Write([]byte(terminal.ClearTerminal()))
	if err != nil {
		c.Log.Println(logging.LogWarn, "Flush Error:", c.ID, " Message: ", err, " Closing connection")
		c.manager.HandleDisconnect(c)
		return
	}
	_, err = c.conn.Write([]byte(terminal.HideCursor()))
	if err != nil {
		c.Log.Println(logging.LogWarn, "Hide Cursor Error:", c.ID, " Message: ", err, " Closing connection")
		c.manager.HandleDisconnect(c)
		return
	}

	c.Console = ui.NewConsole(width, height, c.ID, c.manager.CMReceiveMessages, c.Log, terminal)
	//log.Println("Initializing terminal type for:", c.ID)
	//log.Println("Initializing Console for:", c.ID)
	c.Console.Init()
	//log.Println("Console Initialized for:", c.ID)

	//log.Println("Launching Write Handler for:", c.ID)
	go c.WriteHandler()
	//log.Println("Launching Read Handler for:", c.ID)
	go c.ReadHandler()
	c.Log.Println(logging.LogInfo, "Connection successfully created for:", c.ID)
}

// ReadHandler is launched as a goroutine that handles reading bytes from the connection
func (c *Connection) ReadHandler() {
	// Enter our client loop
	reader := bufio.NewReader(c.conn)
	for {
		readByte, err := reader.ReadByte()
		if err != nil {
			c.Log.Println(logging.LogInfo, "Client "+c.ID+" closed connection")
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

		c.Console.HandleInput(readByte)
	}
}

// WriteHandler is launched as a goroutine that handles writing bytes to the connection
// It is constantly pulling messages from the Console and writing them to the connection
func (c *Connection) WriteHandler() {
	for {
		if c.Console.GetShutdown() {
			c.Log.Println(logging.LogInfo, "Client "+c.ID+" requested shutdown")
			c.Write([]byte("\033c"))
			c.Write([]byte("\r\nFor the latest updates, be sure to check the Escaping Eden Discord at: https://discord.gg/uMxZnjJGGu\r\n" +
				"\r\nSee you back soon! Goodbye :)\r\n\r\n"))
			c.Close()
			c.manager.HandleDisconnect(c)
			return
		}

		output := c.Console.Draw()
		if len(output) > 0 {
			err := c.Write(output)
			if err != nil {
				c.Log.Println(logging.LogInfo, "Client "+c.ID+" disconnected")
				c.Close()
				c.manager.HandleDisconnect(c)
				return
			}
		}

		if c.GetResizeCleanup() {
			//log.Println("Cleanup requested")
			//log.Println(c.cleanupStage)
			switch c.cleanupStage {
			case 0:
				//c.Console.ForceRedraw()
				c.cleanupStage++
			case 1:
				//c.Console.ForceRedraw()
				c.cleanupStage++
			case 2:
				c.Console.ResetWindowDrawings()
				c.cleanupStage++
			case 3:
				c.Console.ClearPointMap()
				c.cleanupStage++
			case 4:
				c.Console.FlushLastSent()
				c.cleanupStage++
			case 5:
				c.ResizeCleanupComplete()
				c.cleanupStage = 0
			}

		}
	}
}

// HandleIAC handles IAC codes
func (c *Connection) HandleIAC(readByte byte) {
	//log.Println("Byte: ", readByte)

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
					c.Console.HandleResize(c.iacWindowResizeX, c.iacWindowResizeY)
					c.NotifyCleanupResize()
				}
				c.iacWindowResizeX = 0
				c.iacWindowResizeY = 0
				//log.Println("IAC subnegotiation complete")

				return
			}
		}

		c.Log.Println(logging.LogWarn, "!!! If you're reading this, IAC parsing failed and "+
			"you're about to see garbage get sent to your windows. Unhandled IAC command: ", int(readByte))
		c.iacActive = false
	}
}

func (c *Connection) NotifyCleanupResize() {
	c.cleanupAfterResizeMutex.Lock()
	defer c.cleanupAfterResizeMutex.Unlock()
	c.cleanupAfterResize = true
}

func (c *Connection) GetResizeCleanup() bool {
	c.cleanupAfterResizeMutex.Lock()
	defer c.cleanupAfterResizeMutex.Unlock()
	return c.cleanupAfterResize
}

func (c *Connection) ResizeCleanupComplete() {
	c.cleanupAfterResizeMutex.Lock()
	defer c.cleanupAfterResizeMutex.Unlock()
	c.cleanupAfterResize = false
}

func (c *Connection) SendToConsole(message messages.ConsoleMessage) {
	go func() {
		c.Console.ReceiveMessages <- message
	}()
}
