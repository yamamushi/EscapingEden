package network

import (
	"github.com/google/uuid"
	"github.com/yamamushi/EscapingEden/edenutil"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"net"
	"strings"
	"sync"
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
}

// NewServer creates a new server
func NewServer(host string, port string, log logging.LoggerType) *Server {
	return &Server{
		Host:       host,
		Port:       port,
		ConnectMap: &sync.Map{},
		Log:        log,
	}
}

// Start starts the server
func (s *Server) Start(startedNotify chan bool, cmReceiveMessage chan messages.ConnectionManagerMessage, amReceiveMessages chan messages.AccountManagerMessage, characterManagerReceiveMessages chan messages.CharacterManagerMessage) error {
	l, err := net.Listen("tcp", s.Host+":"+s.Port)
	if err != nil {
		return err
	}
	s.ConnectionManagerSend = cmReceiveMessage
	// Using sync.Map to not deal with concurrency slice/map issues
	s.ConnectionManager = NewConnectionManager(s.ConnectMap, cmReceiveMessage, amReceiveMessages, characterManagerReceiveMessages, s.Log)
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

		addressSlice := strings.Split(conn.RemoteAddr().String(), ":")
		if edenutil.CheckBlacklist(addressSlice[0], edenutil.BlackListIPs) {
			s.Log.Println(logging.LogWarn, "Connection from blacklisted IP: ", conn.RemoteAddr().String())
			_, _ = conn.Write([]byte("\r\nConnections from this IP are not allowed."))
			_ = conn.Close()
			continue
		}

		id := uuid.New().String()
		s.Log.Println(logging.LogInfo, "New connection from", conn.RemoteAddr(), "with id", id, "accepted")
		s.ConnectionManager.AddConnection(NewConnection(conn, id, s.ConnectionManager, s.Log))
	}
}
