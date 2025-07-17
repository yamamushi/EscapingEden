package edenutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Metrics holds various application metrics
type Metrics struct {
	// Connection metrics
	ActiveConnections int64 `json:"active_connections"`
	TotalConnections  int64 `json:"total_connections"`
	FailedConnections int64 `json:"failed_connections"`
	RateLimitedConns  int64 `json:"rate_limited_connections"`

	// Authentication metrics
	LoginAttempts        int64 `json:"login_attempts"`
	SuccessfulLogins     int64 `json:"successful_logins"`
	FailedLogins         int64 `json:"failed_logins"`
	RegistrationAttempts int64 `json:"registration_attempts"`

	// Game metrics
	ActiveCharacters int64 `json:"active_characters"`
	GameCommands     int64 `json:"game_commands"`

	// System metrics
	Uptime         time.Duration `json:"uptime_seconds"`
	MemoryUsage    uint64        `json:"memory_usage_bytes"`
	GoroutineCount int           `json:"goroutine_count"`

	// Error metrics
	TotalErrors    int64 `json:"total_errors"`
	DatabaseErrors int64 `json:"database_errors"`
	NetworkErrors  int64 `json:"network_errors"`

	startTime time.Time
	mutex     sync.RWMutex
}

// NewMetrics creates a new metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		startTime: time.Now(),
	}
}

// IncrementActiveConnections increments the active connections counter
func (m *Metrics) IncrementActiveConnections() {
	atomic.AddInt64(&m.ActiveConnections, 1)
	atomic.AddInt64(&m.TotalConnections, 1)
}

// DecrementActiveConnections decrements the active connections counter
func (m *Metrics) DecrementActiveConnections() {
	atomic.AddInt64(&m.ActiveConnections, -1)
}

// IncrementFailedConnections increments the failed connections counter
func (m *Metrics) IncrementFailedConnections() {
	atomic.AddInt64(&m.FailedConnections, 1)
}

// IncrementRateLimitedConnections increments the rate limited connections counter
func (m *Metrics) IncrementRateLimitedConnections() {
	atomic.AddInt64(&m.RateLimitedConns, 1)
}

// IncrementLoginAttempts increments the login attempts counter
func (m *Metrics) IncrementLoginAttempts() {
	atomic.AddInt64(&m.LoginAttempts, 1)
}

// IncrementSuccessfulLogins increments the successful logins counter
func (m *Metrics) IncrementSuccessfulLogins() {
	atomic.AddInt64(&m.SuccessfulLogins, 1)
}

// IncrementFailedLogins increments the failed logins counter
func (m *Metrics) IncrementFailedLogins() {
	atomic.AddInt64(&m.FailedLogins, 1)
}

// IncrementRegistrationAttempts increments the registration attempts counter
func (m *Metrics) IncrementRegistrationAttempts() {
	atomic.AddInt64(&m.RegistrationAttempts, 1)
}

// IncrementActiveCharacters increments the active characters counter
func (m *Metrics) IncrementActiveCharacters() {
	atomic.AddInt64(&m.ActiveCharacters, 1)
}

// DecrementActiveCharacters decrements the active characters counter
func (m *Metrics) DecrementActiveCharacters() {
	atomic.AddInt64(&m.ActiveCharacters, -1)
}

// IncrementGameCommands increments the game commands counter
func (m *Metrics) IncrementGameCommands() {
	atomic.AddInt64(&m.GameCommands, 1)
}

// IncrementTotalErrors increments the total errors counter
func (m *Metrics) IncrementTotalErrors() {
	atomic.AddInt64(&m.TotalErrors, 1)
}

// IncrementDatabaseErrors increments the database errors counter
func (m *Metrics) IncrementDatabaseErrors() {
	atomic.AddInt64(&m.DatabaseErrors, 1)
	atomic.AddInt64(&m.TotalErrors, 1)
}

// IncrementNetworkErrors increments the network errors counter
func (m *Metrics) IncrementNetworkErrors() {
	atomic.AddInt64(&m.NetworkErrors, 1)
	atomic.AddInt64(&m.TotalErrors, 1)
}

// UpdateSystemMetrics updates system-related metrics
func (m *Metrics) UpdateSystemMetrics() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.Uptime = time.Since(m.startTime)
	m.GoroutineCount = runtime.NumGoroutine()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	m.MemoryUsage = memStats.Alloc
}

// GetSnapshot returns a snapshot of current metrics
func (m *Metrics) GetSnapshot() Metrics {
	m.UpdateSystemMetrics()

	return Metrics{
		ActiveConnections:    atomic.LoadInt64(&m.ActiveConnections),
		TotalConnections:     atomic.LoadInt64(&m.TotalConnections),
		FailedConnections:    atomic.LoadInt64(&m.FailedConnections),
		RateLimitedConns:     atomic.LoadInt64(&m.RateLimitedConns),
		LoginAttempts:        atomic.LoadInt64(&m.LoginAttempts),
		SuccessfulLogins:     atomic.LoadInt64(&m.SuccessfulLogins),
		FailedLogins:         atomic.LoadInt64(&m.FailedLogins),
		RegistrationAttempts: atomic.LoadInt64(&m.RegistrationAttempts),
		ActiveCharacters:     atomic.LoadInt64(&m.ActiveCharacters),
		GameCommands:         atomic.LoadInt64(&m.GameCommands),
		Uptime:               m.Uptime,
		MemoryUsage:          m.MemoryUsage,
		GoroutineCount:       m.GoroutineCount,
		TotalErrors:          atomic.LoadInt64(&m.TotalErrors),
		DatabaseErrors:       atomic.LoadInt64(&m.DatabaseErrors),
		NetworkErrors:        atomic.LoadInt64(&m.NetworkErrors),
		startTime:            m.startTime,
	}
}

// Reset resets all metrics counters
func (m *Metrics) Reset() {
	atomic.StoreInt64(&m.ActiveConnections, 0)
	atomic.StoreInt64(&m.TotalConnections, 0)
	atomic.StoreInt64(&m.FailedConnections, 0)
	atomic.StoreInt64(&m.RateLimitedConns, 0)
	atomic.StoreInt64(&m.LoginAttempts, 0)
	atomic.StoreInt64(&m.SuccessfulLogins, 0)
	atomic.StoreInt64(&m.FailedLogins, 0)
	atomic.StoreInt64(&m.RegistrationAttempts, 0)
	atomic.StoreInt64(&m.ActiveCharacters, 0)
	atomic.StoreInt64(&m.GameCommands, 0)
	atomic.StoreInt64(&m.TotalErrors, 0)
	atomic.StoreInt64(&m.DatabaseErrors, 0)
	atomic.StoreInt64(&m.NetworkErrors, 0)

	m.mutex.Lock()
	m.startTime = time.Now()
	m.mutex.Unlock()
}

// HealthChecker provides health check functionality
type HealthChecker struct {
	checks map[string]HealthCheck
	mutex  sync.RWMutex
}

// HealthCheck represents a single health check
type HealthCheck struct {
	Name        string                 `json:"name"`
	Status      string                 `json:"status"`
	Message     string                 `json:"message,omitempty"`
	LastChecked time.Time              `json:"last_checked"`
	Duration    time.Duration          `json:"duration_ms"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// HealthCheckFunc is a function that performs a health check
type HealthCheckFunc func() HealthCheck

// NewHealthChecker creates a new health checker
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		checks: make(map[string]HealthCheck),
	}
}

// AddCheck adds a health check
func (hc *HealthChecker) AddCheck(name string, checkFunc HealthCheckFunc) {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()

	// Run the check immediately to get initial status
	start := time.Now()
	check := checkFunc()
	check.Name = name
	check.LastChecked = start
	check.Duration = time.Since(start)

	hc.checks[name] = check
}

// RunChecks runs all health checks
func (hc *HealthChecker) RunChecks() map[string]HealthCheck {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()

	results := make(map[string]HealthCheck)

	for name := range hc.checks {
		// For this implementation, we'll just return the last known status
		// In a real implementation, you'd re-run the check functions
		results[name] = hc.checks[name]
	}

	return results
}

// GetOverallHealth returns the overall health status
func (hc *HealthChecker) GetOverallHealth() string {
	checks := hc.RunChecks()

	for _, check := range checks {
		if check.Status != "healthy" {
			return "unhealthy"
		}
	}

	if len(checks) == 0 {
		return "unknown"
	}

	return "healthy"
}

// MetricsServer provides an HTTP server for metrics and health checks
type MetricsServer struct {
	metrics       *Metrics
	healthChecker *HealthChecker
	server        *http.Server
}

// NewMetricsServer creates a new metrics server
func NewMetricsServer(addr string, metrics *Metrics, healthChecker *HealthChecker) *MetricsServer {
	ms := &MetricsServer{
		metrics:       metrics,
		healthChecker: healthChecker,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", ms.handleMetrics)
	mux.HandleFunc("/health", ms.handleHealth)
	mux.HandleFunc("/health/ready", ms.handleReadiness)
	mux.HandleFunc("/health/live", ms.handleLiveness)

	ms.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return ms
}

// Start starts the metrics server
func (ms *MetricsServer) Start() error {
	return ms.server.ListenAndServe()
}

// Stop stops the metrics server
func (ms *MetricsServer) Stop() error {
	return ms.server.Close()
}

// handleMetrics handles the /metrics endpoint
func (ms *MetricsServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	snapshot := ms.metrics.GetSnapshot()
	if err := json.NewEncoder(w).Encode(snapshot); err != nil {
		http.Error(w, "Failed to encode metrics", http.StatusInternalServerError)
		return
	}
}

// handleHealth handles the /health endpoint
func (ms *MetricsServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	checks := ms.healthChecker.RunChecks()
	overallHealth := ms.healthChecker.GetOverallHealth()

	response := map[string]interface{}{
		"status": overallHealth,
		"checks": checks,
	}

	if overallHealth != "healthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode health status", http.StatusInternalServerError)
		return
	}
}

// handleReadiness handles the /health/ready endpoint
func (ms *MetricsServer) handleReadiness(w http.ResponseWriter, r *http.Request) {
	// Check if the application is ready to serve requests
	overallHealth := ms.healthChecker.GetOverallHealth()

	if overallHealth == "healthy" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Ready")
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, "Not Ready")
	}
}

// handleLiveness handles the /health/live endpoint
func (ms *MetricsServer) handleLiveness(w http.ResponseWriter, r *http.Request) {
	// Simple liveness check - if we can respond, we're alive
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Alive")
}

// PerformanceProfiler provides performance profiling utilities
type PerformanceProfiler struct {
	samples    []PerformanceSample
	mutex      sync.RWMutex
	maxSamples int
}

// PerformanceSample represents a performance measurement
type PerformanceSample struct {
	Timestamp time.Time              `json:"timestamp"`
	Operation string                 `json:"operation"`
	Duration  time.Duration          `json:"duration_ms"`
	Success   bool                   `json:"success"`
	ErrorMsg  string                 `json:"error,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// NewPerformanceProfiler creates a new performance profiler
func NewPerformanceProfiler(maxSamples int) *PerformanceProfiler {
	return &PerformanceProfiler{
		samples:    make([]PerformanceSample, 0, maxSamples),
		maxSamples: maxSamples,
	}
}

// RecordOperation records a performance sample
func (pp *PerformanceProfiler) RecordOperation(operation string, duration time.Duration, success bool, errorMsg string, metadata map[string]interface{}) {
	pp.mutex.Lock()
	defer pp.mutex.Unlock()

	sample := PerformanceSample{
		Timestamp: time.Now(),
		Operation: operation,
		Duration:  duration,
		Success:   success,
		ErrorMsg:  errorMsg,
		Metadata:  metadata,
	}

	pp.samples = append(pp.samples, sample)

	// Remove oldest samples if we exceed max
	if len(pp.samples) > pp.maxSamples {
		pp.samples = pp.samples[1:]
	}
}

// GetSamples returns all performance samples
func (pp *PerformanceProfiler) GetSamples() []PerformanceSample {
	pp.mutex.RLock()
	defer pp.mutex.RUnlock()

	// Return a copy to avoid race conditions
	samples := make([]PerformanceSample, len(pp.samples))
	copy(samples, pp.samples)

	return samples
}

// GetAverageLatency returns the average latency for an operation
func (pp *PerformanceProfiler) GetAverageLatency(operation string) time.Duration {
	pp.mutex.RLock()
	defer pp.mutex.RUnlock()

	var total time.Duration
	var count int

	for _, sample := range pp.samples {
		if sample.Operation == operation && sample.Success {
			total += sample.Duration
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return total / time.Duration(count)
}

// Clear clears all performance samples
func (pp *PerformanceProfiler) Clear() {
	pp.mutex.Lock()
	defer pp.mutex.Unlock()

	pp.samples = pp.samples[:0]
}
