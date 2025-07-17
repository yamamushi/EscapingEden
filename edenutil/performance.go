package edenutil

import (
	"runtime"
	"sync"
	"time"
)

// ObjectPool provides a generic object pool for reducing allocations
type ObjectPool[T any] struct {
	pool sync.Pool
	new  func() T
}

// NewObjectPool creates a new object pool
func NewObjectPool[T any](newFunc func() T) *ObjectPool[T] {
	return &ObjectPool[T]{
		pool: sync.Pool{
			New: func() interface{} {
				return newFunc()
			},
		},
		new: newFunc,
	}
}

// Get retrieves an object from the pool
func (op *ObjectPool[T]) Get() T {
	return op.pool.Get().(T)
}

// Put returns an object to the pool
func (op *ObjectPool[T]) Put(obj T) {
	op.pool.Put(obj)
}

// MemoryStats provides memory usage statistics
type MemoryStats struct {
	Alloc         uint64      // bytes allocated and not yet freed
	TotalAlloc    uint64      // bytes allocated (even if freed)
	Sys           uint64      // bytes obtained from system
	Lookups       uint64      // number of pointer lookups
	Mallocs       uint64      // number of mallocs
	Frees         uint64      // number of frees
	HeapAlloc     uint64      // bytes allocated and not yet freed (same as Alloc)
	HeapSys       uint64      // bytes obtained from system
	HeapIdle      uint64      // bytes in idle spans
	HeapInuse     uint64      // bytes in non-idle span
	HeapReleased  uint64      // bytes released to the OS
	HeapObjects   uint64      // total number of allocated objects
	StackInuse    uint64      // bytes used by stack spans
	StackSys      uint64      // bytes obtained from system for stack
	MSpanInuse    uint64      // bytes used by mspan structures
	MSpanSys      uint64      // bytes obtained from system for mspan
	MCacheInuse   uint64      // bytes used by mcache structures
	MCacheSys     uint64      // bytes obtained from system for mcache
	GCSys         uint64      // bytes used for garbage collection system metadata
	OtherSys      uint64      // bytes used for other system allocations
	NextGC        uint64      // next collection will happen when HeapAlloc ≥ this amount
	LastGC        uint64      // end time of last collection (nanoseconds since 1970)
	PauseTotalNs  uint64      // cumulative nanoseconds in GC stop-the-world pauses
	PauseNs       [256]uint64 // circular buffer of recent GC pause durations
	PauseEnd      [256]uint64 // circular buffer of recent GC pause end times
	NumGC         uint32      // number of completed GC cycles
	NumForcedGC   uint32      // number of GC cycles that were forced by the application
	GCCPUFraction float64     // fraction of CPU time used by GC
}

// GetMemoryStats returns current memory statistics
func GetMemoryStats() MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return MemoryStats{
		Alloc:         m.Alloc,
		TotalAlloc:    m.TotalAlloc,
		Sys:           m.Sys,
		Lookups:       m.Lookups,
		Mallocs:       m.Mallocs,
		Frees:         m.Frees,
		HeapAlloc:     m.HeapAlloc,
		HeapSys:       m.HeapSys,
		HeapIdle:      m.HeapIdle,
		HeapInuse:     m.HeapInuse,
		HeapReleased:  m.HeapReleased,
		HeapObjects:   m.HeapObjects,
		StackInuse:    m.StackInuse,
		StackSys:      m.StackSys,
		MSpanInuse:    m.MSpanInuse,
		MSpanSys:      m.MSpanSys,
		MCacheInuse:   m.MCacheInuse,
		MCacheSys:     m.MCacheSys,
		GCSys:         m.GCSys,
		OtherSys:      m.OtherSys,
		NextGC:        m.NextGC,
		LastGC:        m.LastGC,
		PauseTotalNs:  m.PauseTotalNs,
		PauseNs:       m.PauseNs,
		PauseEnd:      m.PauseEnd,
		NumGC:         m.NumGC,
		NumForcedGC:   m.NumForcedGC,
		GCCPUFraction: m.GCCPUFraction,
	}
}

// BatchProcessor processes items in batches for better performance
type BatchProcessor[T any] struct {
	batchSize int
	timeout   time.Duration
	processor func([]T) error
	buffer    []T
	mutex     sync.Mutex
	timer     *time.Timer
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor[T any](batchSize int, timeout time.Duration, processor func([]T) error) *BatchProcessor[T] {
	bp := &BatchProcessor[T]{
		batchSize: batchSize,
		timeout:   timeout,
		processor: processor,
		buffer:    make([]T, 0, batchSize),
	}

	bp.timer = time.AfterFunc(timeout, bp.flush)
	bp.timer.Stop()

	return bp
}

// Add adds an item to the batch
func (bp *BatchProcessor[T]) Add(item T) error {
	bp.mutex.Lock()
	defer bp.mutex.Unlock()

	bp.buffer = append(bp.buffer, item)

	// Start timer if this is the first item
	if len(bp.buffer) == 1 {
		bp.timer.Reset(bp.timeout)
	}

	// Process if batch is full
	if len(bp.buffer) >= bp.batchSize {
		return bp.processBatch()
	}

	return nil
}

// Flush processes any remaining items in the buffer
func (bp *BatchProcessor[T]) Flush() error {
	bp.mutex.Lock()
	defer bp.mutex.Unlock()

	return bp.processBatch()
}

// flush is called by the timer
func (bp *BatchProcessor[T]) flush() {
	bp.mutex.Lock()
	defer bp.mutex.Unlock()

	bp.processBatch()
}

// processBatch processes the current batch (must be called with mutex held)
func (bp *BatchProcessor[T]) processBatch() error {
	if len(bp.buffer) == 0 {
		return nil
	}

	bp.timer.Stop()

	// Process the batch
	err := bp.processor(bp.buffer)

	// Clear the buffer
	bp.buffer = bp.buffer[:0]

	return err
}

// LRUCache implements a Least Recently Used cache
type LRUCache[K comparable, V any] struct {
	capacity int
	cache    map[K]*lruNode[K, V]
	head     *lruNode[K, V]
	tail     *lruNode[K, V]
	mutex    sync.RWMutex
}

type lruNode[K comparable, V any] struct {
	key   K
	value V
	prev  *lruNode[K, V]
	next  *lruNode[K, V]
}

// NewLRUCache creates a new LRU cache
func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	cache := &LRUCache[K, V]{
		capacity: capacity,
		cache:    make(map[K]*lruNode[K, V]),
	}

	// Create dummy head and tail nodes
	cache.head = &lruNode[K, V]{}
	cache.tail = &lruNode[K, V]{}
	cache.head.next = cache.tail
	cache.tail.prev = cache.head

	return cache
}

// Get retrieves a value from the cache
func (lru *LRUCache[K, V]) Get(key K) (V, bool) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	if node, exists := lru.cache[key]; exists {
		// Move to front
		lru.moveToFront(node)
		return node.value, true
	}

	var zero V
	return zero, false
}

// Put adds or updates a value in the cache
func (lru *LRUCache[K, V]) Put(key K, value V) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	if node, exists := lru.cache[key]; exists {
		// Update existing node
		node.value = value
		lru.moveToFront(node)
		return
	}

	// Create new node
	node := &lruNode[K, V]{key: key, value: value}
	lru.cache[key] = node
	lru.addToFront(node)

	// Remove least recently used if over capacity
	if len(lru.cache) > lru.capacity {
		tail := lru.removeTail()
		delete(lru.cache, tail.key)
	}
}

// moveToFront moves a node to the front of the list
func (lru *LRUCache[K, V]) moveToFront(node *lruNode[K, V]) {
	lru.removeNode(node)
	lru.addToFront(node)
}

// addToFront adds a node to the front of the list
func (lru *LRUCache[K, V]) addToFront(node *lruNode[K, V]) {
	node.prev = lru.head
	node.next = lru.head.next
	lru.head.next.prev = node
	lru.head.next = node
}

// removeNode removes a node from the list
func (lru *LRUCache[K, V]) removeNode(node *lruNode[K, V]) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

// removeTail removes and returns the tail node
func (lru *LRUCache[K, V]) removeTail() *lruNode[K, V] {
	tail := lru.tail.prev
	lru.removeNode(tail)
	return tail
}
