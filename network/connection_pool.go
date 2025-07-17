package network

import (
	"sync"
	"time"
)

// ConnectionPool manages a pool of connections for better resource management
type ConnectionPool struct {
	connections map[string]*Connection
	mutex       sync.RWMutex
	maxIdle     int
	maxActive   int
	idleTimeout time.Duration
}

// NewConnectionPool creates a new connection pool
func NewConnectionPool(maxIdle, maxActive int, idleTimeout time.Duration) *ConnectionPool {
	pool := &ConnectionPool{
		connections: make(map[string]*Connection),
		maxIdle:     maxIdle,
		maxActive:   maxActive,
		idleTimeout: idleTimeout,
	}

	// Start cleanup goroutine
	go pool.cleanup()

	return pool
}

// Add adds a connection to the pool
func (cp *ConnectionPool) Add(id string, conn *Connection) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	// Don't exceed max active connections
	if len(cp.connections) >= cp.maxActive {
		return
	}

	cp.connections[id] = conn
}

// Get retrieves a connection from the pool
func (cp *ConnectionPool) Get(id string) (*Connection, bool) {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()

	conn, exists := cp.connections[id]
	return conn, exists
}

// Remove removes a connection from the pool
func (cp *ConnectionPool) Remove(id string) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	delete(cp.connections, id)
}

// Count returns the number of connections in the pool
func (cp *ConnectionPool) Count() int {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()

	return len(cp.connections)
}

// Range iterates over all connections in the pool
func (cp *ConnectionPool) Range(fn func(string, *Connection) bool) {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()

	for id, conn := range cp.connections {
		if !fn(id, conn) {
			break
		}
	}
}

// cleanup removes idle connections periodically
func (cp *ConnectionPool) cleanup() {
	ticker := time.NewTicker(cp.idleTimeout)
	defer ticker.Stop()

	for range ticker.C {
		cp.mutex.Lock()
		now := time.Now()

		for id, conn := range cp.connections {
			// Remove connections that have been idle too long
			if now.Sub(conn.lastActivity) > cp.idleTimeout {
				delete(cp.connections, id)
			}
		}

		cp.mutex.Unlock()
	}
}
