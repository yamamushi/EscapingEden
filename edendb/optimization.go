package edendb

import (
	"fmt"
	"sync"
	"time"
)

// ConnectionPool manages database connections for better performance
type ConnectionPool struct {
	connections chan DatabaseType
	factory     func() (DatabaseType, error)
	maxSize     int
	timeout     time.Duration
	mutex       sync.RWMutex
	closed      bool
}

// NewConnectionPool creates a new database connection pool
func NewConnectionPool(maxSize int, timeout time.Duration, factory func() (DatabaseType, error)) *ConnectionPool {
	return &ConnectionPool{
		connections: make(chan DatabaseType, maxSize),
		factory:     factory,
		maxSize:     maxSize,
		timeout:     timeout,
	}
}

// Get retrieves a connection from the pool
func (cp *ConnectionPool) Get() (DatabaseType, error) {
	cp.mutex.RLock()
	if cp.closed {
		cp.mutex.RUnlock()
		return nil, ErrPoolClosed
	}
	cp.mutex.RUnlock()

	select {
	case conn := <-cp.connections:
		return conn, nil
	case <-time.After(cp.timeout):
		// Try to create a new connection if pool is empty
		return cp.factory()
	}
}

// Put returns a connection to the pool
func (cp *ConnectionPool) Put(conn DatabaseType) {
	cp.mutex.RLock()
	if cp.closed {
		cp.mutex.RUnlock()
		return
	}
	cp.mutex.RUnlock()

	select {
	case cp.connections <- conn:
		// Connection returned to pool
	default:
		// Pool is full, connection will be discarded
	}
}

// Close closes the connection pool
func (cp *ConnectionPool) Close() {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	if cp.closed {
		return
	}

	cp.closed = true
	close(cp.connections)

	// Close all connections in the pool
	for conn := range cp.connections {
		// Note: DatabaseType interface doesn't have Close method
		// This would need to be added to the interface
		_ = conn
	}
}

// Size returns the current number of connections in the pool
func (cp *ConnectionPool) Size() int {
	return len(cp.connections)
}

// BatchOperations provides utilities for batching database operations
type BatchOperations struct {
	db        DatabaseType
	batchSize int
	timeout   time.Duration
}

// NewBatchOperations creates a new batch operations handler
func NewBatchOperations(db DatabaseType, batchSize int, timeout time.Duration) *BatchOperations {
	return &BatchOperations{
		db:        db,
		batchSize: batchSize,
		timeout:   timeout,
	}
}

// BatchInsert performs batch insert operations
func (bo *BatchOperations) BatchInsert(collection string, records []interface{}) error {
	if len(records) == 0 {
		return nil
	}

	// Process in batches
	for i := 0; i < len(records); i += bo.batchSize {
		end := i + bo.batchSize
		if end > len(records) {
			end = len(records)
		}

		batch := records[i:end]
		for _, record := range batch {
			if err := bo.db.AddRecord(collection, record); err != nil {
				return err
			}
		}

		// Add small delay between batches to prevent overwhelming the database
		if end < len(records) {
			time.Sleep(10 * time.Millisecond)
		}
	}

	return nil
}

// BatchUpdate performs batch update operations
func (bo *BatchOperations) BatchUpdate(collection string, records []interface{}) error {
	if len(records) == 0 {
		return nil
	}

	// Process in batches
	for i := 0; i < len(records); i += bo.batchSize {
		end := i + bo.batchSize
		if end > len(records) {
			end = len(records)
		}

		batch := records[i:end]
		for _, record := range batch {
			if err := bo.db.UpdateRecord(collection, record); err != nil {
				return err
			}
		}

		// Add small delay between batches
		if end < len(records) {
			time.Sleep(10 * time.Millisecond)
		}
	}

	return nil
}

// QueryCache provides caching for database queries
type QueryCache struct {
	cache   map[string]*CacheEntry
	mutex   sync.RWMutex
	maxSize int
	ttl     time.Duration
}

// CacheEntry represents a cached query result
type CacheEntry struct {
	Data      interface{}
	Timestamp time.Time
}

// NewQueryCache creates a new query cache
func NewQueryCache(maxSize int, ttl time.Duration) *QueryCache {
	cache := &QueryCache{
		cache:   make(map[string]*CacheEntry),
		maxSize: maxSize,
		ttl:     ttl,
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// Get retrieves a cached query result
func (qc *QueryCache) Get(key string) (interface{}, bool) {
	qc.mutex.RLock()
	defer qc.mutex.RUnlock()

	entry, exists := qc.cache[key]
	if !exists {
		return nil, false
	}

	// Check if entry has expired
	if time.Since(entry.Timestamp) > qc.ttl {
		return nil, false
	}

	return entry.Data, true
}

// Set stores a query result in the cache
func (qc *QueryCache) Set(key string, data interface{}) {
	qc.mutex.Lock()
	defer qc.mutex.Unlock()

	// Remove oldest entries if cache is full
	if len(qc.cache) >= qc.maxSize {
		qc.evictOldest()
	}

	qc.cache[key] = &CacheEntry{
		Data:      data,
		Timestamp: time.Now(),
	}
}

// evictOldest removes the oldest cache entry
func (qc *QueryCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range qc.cache {
		if oldestKey == "" || entry.Timestamp.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.Timestamp
		}
	}

	if oldestKey != "" {
		delete(qc.cache, oldestKey)
	}
}

// cleanup removes expired entries periodically
func (qc *QueryCache) cleanup() {
	ticker := time.NewTicker(qc.ttl / 2)
	defer ticker.Stop()

	for range ticker.C {
		qc.mutex.Lock()
		now := time.Now()

		for key, entry := range qc.cache {
			if now.Sub(entry.Timestamp) > qc.ttl {
				delete(qc.cache, key)
			}
		}

		qc.mutex.Unlock()
	}
}

// DatabaseMetrics tracks database performance metrics
type DatabaseMetrics struct {
	QueryCount    int64
	ErrorCount    int64
	TotalDuration time.Duration
	mutex         sync.RWMutex
}

// NewDatabaseMetrics creates a new metrics tracker
func NewDatabaseMetrics() *DatabaseMetrics {
	return &DatabaseMetrics{}
}

// RecordQuery records a database query execution
func (dm *DatabaseMetrics) RecordQuery(duration time.Duration, err error) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dm.QueryCount++
	dm.TotalDuration += duration

	if err != nil {
		dm.ErrorCount++
	}
}

// GetStats returns current database statistics
func (dm *DatabaseMetrics) GetStats() (queryCount int64, errorCount int64, avgDuration time.Duration) {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	queryCount = dm.QueryCount
	errorCount = dm.ErrorCount

	if dm.QueryCount > 0 {
		avgDuration = dm.TotalDuration / time.Duration(dm.QueryCount)
	}

	return
}

// Reset resets all metrics
func (dm *DatabaseMetrics) Reset() {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dm.QueryCount = 0
	dm.ErrorCount = 0
	dm.TotalDuration = 0
}

// Common errors
var (
	ErrPoolClosed = fmt.Errorf("connection pool is closed")
)
