package network

import (
	"github.com/google/uuid"
	"log"
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
}

// NewServer creates a new server
func NewServer(host string, port string) *Server {
	return &Server{
		Host:       host,
		Port:       port,
		ConnectMap: &sync.Map{},
	}
}

// Start starts the server
func (s *Server) Start(startedNotify chan bool) error {
	l, err := net.Listen("tcp", s.Host+":"+s.Port)
	if err != nil {
		return err
	}

	// Using sync.Map to not deal with concurrency slice/map issues
	s.ConnectionManager = NewConnectionManager(s.ConnectMap, startedNotify)
	go s.ConnectionManager.Run()
	go s.Listen(l)
	return nil
}

// Listen listens for new connections and adds them to the connection manager
func (s *Server) Listen(l net.Listener) {
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Fatal Error Accepting Connection: ", err)
			return
		}
		id := uuid.New().String()
		log.Println("Storing New Connection: ", id)
		s.ConnectionManager.AddConnection(NewConnection(conn, id, s.ConnectionManager))
	}
}
