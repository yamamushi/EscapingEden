package network

import (
	"github.com/google/uuid"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
	"net"
	"sync"
)

/*
Server manages connections
*/

// Server manages new connections
type Server struct {
	Host string
	Port string

	ConnectMap        *sync.Map
	ConnectionManager *ConnectionManager
	Log               logging.LoggerType
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
func (s *Server) Start(startedNotify chan bool, cmReceiveMessage chan messages.ConnectionManagerMessage, amReceiveMessages chan messages.AccountManagerMessage) error {
	l, err := net.Listen("tcp", s.Host+":"+s.Port)
	if err != nil {
		return err
	}

	// Using sync.Map to not deal with concurrency slice/map issues
	s.ConnectionManager = NewConnectionManager(s.ConnectMap, cmReceiveMessage, amReceiveMessages, s.Log)
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
		id := uuid.New().String()
		s.Log.Println(logging.LogInfo, "New connection from", conn.RemoteAddr(), "with id", id, "accepted")
		s.ConnectionManager.AddConnection(NewConnection(conn, id, s.ConnectionManager, s.Log))
	}
}
