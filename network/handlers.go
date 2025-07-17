package network

import (
	"net"
	"time"

	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

// MessageHandlers contains all message handling logic for the connection manager
type MessageHandlers struct {
	cm *ConnectionManager
}

// NewMessageHandlers creates a new message handlers instance
func NewMessageHandlers(cm *ConnectionManager) *MessageHandlers {
	return &MessageHandlers{cm: cm}
}

// HandleBroadcast handles broadcast messages to all connections
func (mh *MessageHandlers) HandleBroadcast(data interface{}) {
	mh.cm.connectionMap.Range(func(key, value interface{}) bool {
		if conn, ok := value.(*Connection); ok {
			consoleMessage := messages.ConsoleMessage{
				Type: messages.Console_Message_Broadcast,
				Data: data,
			}
			conn.SendToConsole(consoleMessage)
		}
		return true
	})
}

// HandleServerShutdown handles server shutdown by disconnecting all clients
func (mh *MessageHandlers) HandleServerShutdown() {
	mh.cm.Log.Println(logging.LogInfo, "Server shutdown message received on Connection Manager, disconnecting all clients.")
	mh.cm.connectionMap.Range(func(key, value interface{}) bool {
		if conn, ok := value.(*Connection); ok {
			mh.cm.HandleDisconnect(conn)
			conn.Write([]byte("\033c\r\nThe Escaping Eden server has shut down for maintenance.\r\n" +
				"For the latest status updates, be sure to check the Escaping Eden Discord at: https://discord.gg/uMxZnjJGGu\r\n\r\n"))
			conn.Close()
		}
		return true
	})
}

// HandleError handles error messages to specific connections
func (mh *MessageHandlers) HandleError(recipientID string, data interface{}) {
	mh.cm.connectionMap.Range(func(key, value interface{}) bool {
		if conn, ok := value.(*Connection); ok {
			if recipientID == conn.ID {
				consoleMessage := messages.ConsoleMessage{
					Type: messages.Console_Message_Error,
					Data: data,
				}
				conn.SendToConsole(consoleMessage)
			}
		}
		return true
	})
}

// HandleQuit handles quit messages to specific connections
func (mh *MessageHandlers) HandleQuit(recipientID string) {
	mh.cm.connectionMap.Range(func(key, value interface{}) bool {
		if conn, ok := value.(*Connection); ok {
			if recipientID == conn.ID {
				consoleMessage := messages.ConsoleMessage{Type: messages.Console_Message_Quit}
				conn.SendToConsole(consoleMessage)
			}
		}
		return true
	})
}

// HandleRegistration forwards registration requests to account manager
func (mh *MessageHandlers) HandleRegistration(senderID string, data interface{}) {
	go func() {
		registrationRequest := data.(messages.AccountRegistrationRequest)
		mh.cm.AMSendMessages <- messages.AccountManagerMessage{
			Type:            messages.AccountManager_Message_Register,
			Data:            registrationRequest,
			SenderConsoleID: senderID,
		}
	}()
}

// HandleLogin forwards login requests to account manager
func (mh *MessageHandlers) HandleLogin(senderID string, data interface{}) {
	go func() {
		loginRequest := data.(messages.AccountLoginRequest)
		mh.cm.AMSendMessages <- messages.AccountManagerMessage{
			Type:            messages.AccountManager_Message_Login,
			Data:            loginRequest,
			SenderConsoleID: senderID,
		}
	}()
}

// HandleBadLogins manages bad login attempts and potential IP blocking
func (mh *MessageHandlers) HandleBadLogins(conn *Connection) bool {
	ipAddress, _, _ := net.SplitHostPort(conn.conn.RemoteAddr().String())

	// Get current bad login count for this IP
	badLogins := mh.cm.getBadLoginCount(ipAddress)
	badLogins++

	// Update bad login count
	mh.cm.setBadLoginCount(ipAddress, badLogins)

	// If too many bad logins, disconnect
	if badLogins >= 5 {
		mh.cm.Log.Println(logging.LogWarn, "Too many bad login attempts from IP:", ipAddress)
		conn.Write([]byte("\r\nToo many failed login attempts. Connection terminated.\r\n"))
		return true // disconnect
	}

	return false // don't disconnect
}

// CheckIPBadLogins checks if an IP has too many bad login attempts
func (mh *MessageHandlers) CheckIPBadLogins(conn *Connection) bool {
	ipAddress, _, _ := net.SplitHostPort(conn.conn.RemoteAddr().String())
	badLogins := mh.cm.getBadLoginCount(ipAddress)

	// Allow connection if less than threshold
	return badLogins < 5
}

// BadLoginTracker manages bad login attempts per IP
type BadLoginTracker struct {
	attempts map[string]int
	lastSeen map[string]time.Time
}

// NewBadLoginTracker creates a new bad login tracker
func NewBadLoginTracker() *BadLoginTracker {
	return &BadLoginTracker{
		attempts: make(map[string]int),
		lastSeen: make(map[string]time.Time),
	}
}

// Increment increases the bad login count for an IP
func (blt *BadLoginTracker) Increment(ip string) {
	blt.attempts[ip]++
	blt.lastSeen[ip] = time.Now()
}

// Get returns the bad login count for an IP
func (blt *BadLoginTracker) Get(ip string) int {
	return blt.attempts[ip]
}

// Cleanup removes old entries to prevent memory leaks
func (blt *BadLoginTracker) Cleanup() {
	now := time.Now()
	for ip, lastSeen := range blt.lastSeen {
		// Remove entries older than 1 hour
		if now.Sub(lastSeen) > time.Hour {
			delete(blt.attempts, ip)
			delete(blt.lastSeen, ip)
		}
	}
}
