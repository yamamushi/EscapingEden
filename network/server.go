package network

import (
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/yamamushi/EscapingEden/edendb"
	"github.com/yamamushi/EscapingEden/edenutil"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

/*
Server manages connections
*/

// Server manages new connections
type Server struct {
	Host string
	Port string

	ConnectMap            *sync.Map
	ConnectionManager     *ConnectionManager
	ConnectionManagerSend chan messages.ConnectionManagerMessage
	Log                   logging.LoggerType
	RateLimiter           *edenutil.ConnectionRateLimiter
}

// NewServer creates a new server
func NewServer(host string, port string, log logging.LoggerType) *Server {
	// Allow 10 connections per minute per IP
	rateLimiter := edenutil.NewConnectionRateLimiter(10, time.Minute)

	return &Server{
		Host:        host,
		Port:        port,
		ConnectMap:  &sync.Map{},
		Log:         log,
		RateLimiter: rateLimiter,
	}
}

// Start starts the server
func (s *Server) Start(startedNotify chan bool, cmReceiveMessage chan messages.ConnectionManagerMessage, amReceiveMessages chan messages.AccountManagerMessage, characterManagerReceiveMessages chan messages.CharacterManagerMessage,
	ebSendMessages chan messages.EdenbotMessage, gmSendMessages chan messages.GameManagerMessage, db edendb.DatabaseType) error {
	l, err := net.Listen("tcp", s.Host+":"+s.Port)
	if err != nil {
		return err
	}
	s.ConnectionManagerSend = cmReceiveMessage
	// Using sync.Map to not deal with concurrency slice/map issues
	s.ConnectionManager = NewConnectionManager(s.ConnectMap, cmReceiveMessage, amReceiveMessages, characterManagerReceiveMessages, ebSendMessages, gmSendMessages, db, s.Log)
	go s.ConnectionManager.Run(startedNotify)
	go s.Listen(l)
	return nil
}

// Listen listens for new connections and adds them to the connection manager
func (s *Server) Listen(l net.Listener) {
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			s.Log.Println(logging.LogWarn, "Error Accepting Connection: ", err)
			continue
		}

		ipAddress, _, _ := net.SplitHostPort(conn.RemoteAddr().String())

		// Check blacklist
		if edenutil.CheckBlacklist(ipAddress, edenutil.BlackListIPs) {
			s.Log.Println(logging.LogWarn, "Connection from blacklisted IP: ", conn.RemoteAddr().String())
			_, _ = conn.Write([]byte("\r\nConnections from this IP are not allowed."))
			_ = conn.Close()
			continue
		}

		// Check rate limit
		if !s.RateLimiter.Allow(ipAddress) {
			s.Log.Println(logging.LogWarn, "Rate limit exceeded for IP: ", ipAddress)
			_, _ = conn.Write([]byte("\r\nToo many connections. Please try again later."))
			_ = conn.Close()
			continue
		}

		id := uuid.New().String()
		ipaddress, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
		s.Log.Println(logging.LogInfo, "New connection accepted from: "+ipaddress+" with id: "+id)
		s.ConnectionManager.AddConnection(NewConnection(conn, id, s.ConnectionManager, s.Log))
	}
}
