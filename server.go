package main

import (
	"github.com/google/uuid"
	"log"
	"net"
	"sync"
)

/*
Server manages connections
*/

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
func (s *Server) Start() {
	l, err := net.Listen("tcp", s.Host+":"+s.Port)
	if err != nil {
		return
	}

	defer l.Close()
	// Using sync.Map to not deal with concurrency slice/map issues
	s.ConnectionManager = NewConnectionManager(s.ConnectMap)
	go s.ConnectionManager.Run()

	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		id := uuid.New().String()
		log.Println("Storing Connection")
		s.ConnectionManager.AddConnection(NewConnection(conn, id, s.ConnectionManager))
	}
}
