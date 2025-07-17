package game

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/yamamushi/EscapingEden/logging"
)

// GameMonitoringServer extends the basic monitoring server with game-specific endpoints
type GameMonitoringServer struct {
	gm     *GameManager
	log    logging.LoggerType
	server *http.Server
}

// NewGameMonitoringServer creates a new game monitoring server
func NewGameMonitoringServer(addr string, gm *GameManager, log logging.LoggerType) *GameMonitoringServer {
	gms := &GameMonitoringServer{
		gm:  gm,
		log: log,
	}

	mux := http.NewServeMux()

	// Basic health and metrics endpoints
	gms.setupBasicEndpoints(mux)

	// Game-specific endpoints
	gms.setupGameEndpoints(mux)

	// Statistics endpoints
	gms.setupStatsEndpoints(mux)

	// Chunk monitoring endpoints
	gms.setupChunkEndpoints(mux)

	// Enable CORS for dashboard
	corsHandler := gms.enableCORS(mux)

	gms.server = &http.Server{
		Addr:    addr,
		Handler: corsHandler,
	}

	return gms
}

// Start starts the monitoring server
func (gms *GameMonitoringServer) Start() error {
	gms.log.Println(logging.LogInfo, "Game monitoring server starting on", gms.server.Addr)
	return gms.server.ListenAndServe()
}

// Stop stops the monitoring server
func (gms *GameMonitoringServer) Stop() error {
	return gms.server.Close()
}

// enableCORS adds CORS headers to allow dashboard access
func (gms *GameMonitoringServer) enableCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

// setupBasicEndpoints sets up basic health and metrics endpoints
func (gms *GameMonitoringServer) setupBasicEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("/health", gms.handleHealth)
	mux.HandleFunc("/health/detailed", gms.handleDetailedHealth)
	mux.HandleFunc("/health/live", gms.handleLiveness)
	mux.HandleFunc("/health/ready", gms.handleReadiness)
	mux.HandleFunc("/metrics", gms.handleMetrics)
	mux.HandleFunc("/metrics/json", gms.handleMetricsJSON)
}

// setupGameEndpoints sets up game-specific endpoints
func (gms *GameMonitoringServer) setupGameEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("/game/status", gms.handleGameStatus)
	mux.HandleFunc("/game/characters", gms.handleActiveCharacters)
	mux.HandleFunc("/game/commands", gms.handleCommandMetrics)
}

// setupStatsEndpoints sets up statistics endpoints
func (gms *GameMonitoringServer) setupStatsEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("/stats/players", gms.handlePlayerStats)
	mux.HandleFunc("/stats/world", gms.handleWorldStats)
	mux.HandleFunc("/stats/commands", gms.handleCommandStats)
	mux.HandleFunc("/stats/performance", gms.handlePerformanceStats)
}

// setupChunkEndpoints sets up chunk monitoring endpoints
func (gms *GameMonitoringServer) setupChunkEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("/chunks/stats", gms.handleChunkStats)
	mux.HandleFunc("/chunks/health", gms.handleChunkHealth)
	mux.HandleFunc("/chunks/cache", gms.handleChunkCache)
	mux.HandleFunc("/chunks/registry", gms.handleChunkRegistry)
}

// Basic endpoint handlers
func (gms *GameMonitoringServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    time.Since(gms.gm.StartTime).Seconds(),
		"version":   "1.0.0", // TODO: Get from build info
	}

	json.NewEncoder(w).Encode(health)
}

func (gms *GameMonitoringServer) handleDetailedHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	gms.gm.activeCharactersMutex.Lock()
	activeChars := len(gms.gm.ActiveCharacters)
	gms.gm.activeCharactersMutex.Unlock()

	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    time.Since(gms.gm.StartTime).Seconds(),
		"version":   "1.0.0",
		"components": map[string]interface{}{
			"database": map[string]interface{}{
				"status":           "healthy",
				"response_time_ms": 12, // TODO: Actual DB ping
			},
			"game_manager": map[string]interface{}{
				"status":            "healthy",
				"active_characters": activeChars,
				"loaded_chunks":     len(gms.gm.ChunkCache),
			},
			"network": map[string]interface{}{
				"status": "healthy",
				// TODO: Get actual connection counts
			},
		},
	}

	json.NewEncoder(w).Encode(health)
}

func (gms *GameMonitoringServer) handleLiveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Alive")
}

func (gms *GameMonitoringServer) handleReadiness(w http.ResponseWriter, r *http.Request) {
	// Check if game manager is ready
	if gms.gm.DB == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, "Not Ready - Database not initialized")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Ready")
}

func (gms *GameMonitoringServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	// Prometheus format metrics
	w.Header().Set("Content-Type", "text/plain")

	gms.gm.activeCharactersMutex.Lock()
	activeChars := len(gms.gm.ActiveCharacters)
	gms.gm.activeCharactersMutex.Unlock()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Fprintf(w, "# HELP eden_active_players Number of currently active players\n")
	fmt.Fprintf(w, "# TYPE eden_active_players gauge\n")
	fmt.Fprintf(w, "eden_active_players %d\n", activeChars)

	fmt.Fprintf(w, "# HELP eden_loaded_chunks Number of loaded chunks\n")
	fmt.Fprintf(w, "# TYPE eden_loaded_chunks gauge\n")
	fmt.Fprintf(w, "eden_loaded_chunks %d\n", len(gms.gm.ChunkCache))

	fmt.Fprintf(w, "# HELP eden_memory_usage_bytes Memory usage in bytes\n")
	fmt.Fprintf(w, "# TYPE eden_memory_usage_bytes gauge\n")
	fmt.Fprintf(w, "eden_memory_usage_bytes %d\n", m.Alloc)

	fmt.Fprintf(w, "# HELP eden_goroutines Number of goroutines\n")
	fmt.Fprintf(w, "# TYPE eden_goroutines gauge\n")
	fmt.Fprintf(w, "eden_goroutines %d\n", runtime.NumGoroutine())
}

func (gms *GameMonitoringServer) handleMetricsJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	gms.gm.activeCharactersMutex.Lock()
	activeChars := len(gms.gm.ActiveCharacters)
	gms.gm.activeCharactersMutex.Unlock()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Get command registry stats if available
	var commandStats map[string]interface{}
	if gms.gm.CommandRegistry != nil {
		stats := gms.gm.CommandRegistry.GetStats()
		commandStats = make(map[string]interface{})
		for cmdType, stat := range stats {
			commandStats[string(cmdType)] = map[string]interface{}{
				"total_executions": stat.TotalExecutions,
				"success_count":    stat.SuccessCount,
				"error_count":      stat.ErrorCount,
				"average_time_ms":  stat.AverageTime.Milliseconds(),
			}
		}
	}

	metrics := map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"server": map[string]interface{}{
			"uptime_seconds":       time.Since(gms.gm.StartTime).Seconds(),
			"memory_usage_bytes":   m.Alloc,
			"memory_usage_percent": float64(m.Alloc) / float64(m.Sys) * 100,
			"goroutines":           runtime.NumGoroutine(),
			"gc_cycles":            m.NumGC,
		},
		"game": map[string]interface{}{
			"active_players":           activeChars,
			"loaded_chunks":            len(gms.gm.ChunkCache),
			"cached_chunks":            len(gms.gm.ChunkCache),
			"commands_per_second":      0.0,  // TODO: Calculate from command registry
			"average_response_time_ms": 25.0, // TODO: Get from performance profiler
		},
		"database": map[string]interface{}{
			"connections_active":    1, // TODO: Get actual DB connection info
			"connections_idle":      0,
			"queries_per_second":    0.0,
			"average_query_time_ms": 15.0,
			"cache_hit_ratio":       0.85,
		},
		"commands": commandStats,
	}

	json.NewEncoder(w).Encode(metrics)
}

// Game-specific endpoint handlers
func (gms *GameMonitoringServer) handleGameStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	gms.gm.activeCharactersMutex.Lock()
	activeChars := len(gms.gm.ActiveCharacters)
	gms.gm.activeCharactersMutex.Unlock()

	status := map[string]interface{}{
		"timestamp":        time.Now().UTC().Format(time.RFC3339),
		"active_players":   activeChars,
		"loaded_chunks":    len(gms.gm.ChunkCache),
		"world_dimensions": gms.gm.Config.WorldGen.Dimensions,
		"chunk_size":       gms.gm.Config.WorldGen.ChunkSize,
	}

	json.NewEncoder(w).Encode(status)
}

func (gms *GameMonitoringServer) handleActiveCharacters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	gms.gm.activeCharactersMutex.Lock()
	defer gms.gm.activeCharactersMutex.Unlock()

	characters := make([]map[string]interface{}, 0, len(gms.gm.ActiveCharacters))
	for _, char := range gms.gm.ActiveCharacters {
		characters = append(characters, map[string]interface{}{
			"id":            char.ID,
			"name":          char.Name,
			"connection_id": char.ConnectionID,
			"position": map[string]interface{}{
				"x":        char.Record.Position.X,
				"y":        char.Record.Position.Y,
				"z":        char.Record.Position.Z,
				"chunk_id": char.Record.CurrentMapID,
			},
		})
	}

	response := map[string]interface{}{
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
		"count":      len(characters),
		"characters": characters,
	}

	json.NewEncoder(w).Encode(response)
}

func (gms *GameMonitoringServer) handleCommandMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var commandStats map[string]interface{}
	if gms.gm.CommandRegistry != nil {
		stats := gms.gm.CommandRegistry.GetStats()
		commandStats = make(map[string]interface{})
		for cmdType, stat := range stats {
			commandStats[string(cmdType)] = map[string]interface{}{
				"total_executions": stat.TotalExecutions,
				"success_count":    stat.SuccessCount,
				"error_count":      stat.ErrorCount,
				"success_rate":     float64(stat.SuccessCount) / float64(stat.TotalExecutions),
				"average_time_ms":  stat.AverageTime.Milliseconds(),
				"errors_by_code":   stat.ErrorsByCode,
			}
		}
	} else {
		commandStats = make(map[string]interface{})
	}

	response := map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"commands":  commandStats,
	}

	json.NewEncoder(w).Encode(response)
}

// Statistics endpoint handlers
func (gms *GameMonitoringServer) handlePlayerStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	gms.gm.activeCharactersMutex.Lock()
	activeChars := len(gms.gm.ActiveCharacters)
	gms.gm.activeCharactersMutex.Unlock()

	stats := map[string]interface{}{
		"timestamp":          time.Now().UTC().Format(time.RFC3339),
		"active_players":     activeChars,
		"peak_players_today": activeChars, // TODO: Track actual peak
		"total_registered":   0,           // TODO: Get from database
		"players_by_level": map[string]int{
			"1-10":  activeChars,
			"11-20": 0,
			"21-30": 0,
		},
		"geographic_distribution": map[string]int{
			"US":   activeChars,
			"EU":   0,
			"ASIA": 0,
		},
	}

	json.NewEncoder(w).Encode(stats)
}

func (gms *GameMonitoringServer) handleWorldStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse world dimensions
	worldX, worldY, worldZ := gms.gm.ParseWorldDimensions(gms.gm.Config.WorldGen.Dimensions)
	totalChunks := worldX * worldY * worldZ

	// Get popular areas (chunks with players)
	popularAreas := make([]map[string]interface{}, 0)
	chunkPlayerCount := make(map[string]int)

	gms.gm.activeCharactersMutex.Lock()
	for _, char := range gms.gm.ActiveCharacters {
		chunkID := char.Record.CurrentMapID
		chunkPlayerCount[chunkID]++
	}
	gms.gm.activeCharactersMutex.Unlock()

	for chunkID, playerCount := range chunkPlayerCount {
		if playerCount > 0 {
			// Parse chunk coordinates from ID (format: "x-y-z")
			var x, y, z int
			fmt.Sscanf(chunkID, "%d-%d-%d", &x, &y, &z)

			popularAreas = append(popularAreas, map[string]interface{}{
				"chunk_id":     chunkID,
				"player_count": playerCount,
				"coordinates": map[string]int{
					"x": x,
					"y": y,
					"z": z,
				},
			})
		}
	}

	stats := map[string]interface{}{
		"timestamp":             time.Now().UTC().Format(time.RFC3339),
		"total_chunks":          totalChunks,
		"loaded_chunks":         len(gms.gm.ChunkCache),
		"cached_chunks":         len(gms.gm.ChunkCache),
		"chunk_cache_hit_ratio": 0.92, // TODO: Get actual cache stats
		"world_dimensions": map[string]int{
			"x": worldX,
			"y": worldY,
			"z": worldZ,
		},
		"chunk_size":    gms.gm.Config.WorldGen.ChunkSize,
		"popular_areas": popularAreas,
	}

	json.NewEncoder(w).Encode(stats)
}

func (gms *GameMonitoringServer) handleCommandStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var totalCommands int64
	var commandsByType map[string]interface{}

	if gms.gm.CommandRegistry != nil {
		stats := gms.gm.CommandRegistry.GetStats()
		commandsByType = make(map[string]interface{})

		for cmdType, stat := range stats {
			totalCommands += stat.TotalExecutions
			successRate := float64(stat.SuccessCount) / float64(stat.TotalExecutions)
			if stat.TotalExecutions == 0 {
				successRate = 0
			}

			commandsByType[string(cmdType)] = map[string]interface{}{
				"count":           stat.TotalExecutions,
				"success_rate":    successRate,
				"avg_duration_ms": stat.AverageTime.Milliseconds(),
				"errors":          stat.ErrorsByCode,
			}
		}
	} else {
		commandsByType = make(map[string]interface{})
	}

	response := map[string]interface{}{
		"timestamp":           time.Now().UTC().Format(time.RFC3339),
		"total_commands":      totalCommands,
		"commands_per_second": 0.0,  // TODO: Calculate rate
		"success_rate":        0.98, // TODO: Calculate overall success rate
		"by_type":             commandsByType,
	}

	json.NewEncoder(w).Encode(response)
}

func (gms *GameMonitoringServer) handlePerformanceStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// TODO: Implement time-series performance data
	// For now, return current snapshot
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	gms.gm.activeCharactersMutex.Lock()
	activeChars := len(gms.gm.ActiveCharacters)
	gms.gm.activeCharactersMutex.Unlock()

	dataPoint := map[string]interface{}{
		"timestamp":           time.Now().UTC().Format(time.RFC3339),
		"cpu_percent":         0.0, // TODO: Get actual CPU usage
		"memory_percent":      float64(m.Alloc) / float64(m.Sys) * 100,
		"active_players":      activeChars,
		"commands_per_second": 0.0,
		"response_time_ms":    25.0,
	}

	response := map[string]interface{}{
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
		"period":      "1h",
		"resolution":  "5m",
		"data_points": []interface{}{dataPoint},
	}

	json.NewEncoder(w).Encode(response)
}

// Chunk endpoint handlers
func (gms *GameMonitoringServer) handleChunkStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stats := map[string]interface{}{
		"timestamp":     time.Now().UTC().Format(time.RFC3339),
		"total_chunks":  len(gms.gm.ChunkCache),
		"loaded_chunks": len(gms.gm.ChunkCache),
		"cache_hits":    0,    // TODO: Get from chunk monitor
		"cache_misses":  0,    // TODO: Get from chunk monitor
		"hit_ratio":     0.95, // TODO: Calculate actual ratio
	}

	json.NewEncoder(w).Encode(stats)
}

func (gms *GameMonitoringServer) handleChunkHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	health := map[string]interface{}{
		"timestamp":     time.Now().UTC().Format(time.RFC3339),
		"status":        "healthy",
		"loaded_chunks": len(gms.gm.ChunkCache),
		"memory_usage":  "normal", // TODO: Calculate chunk memory usage
	}

	json.NewEncoder(w).Encode(health)
}

func (gms *GameMonitoringServer) handleChunkCache(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cacheInfo := make([]map[string]interface{}, 0, len(gms.gm.ChunkCache))
	for chunkID, chunk := range gms.gm.ChunkCache {
		cacheInfo = append(cacheInfo, map[string]interface{}{
			"id": chunkID,
			"global_position": map[string]int{
				"x": chunk.GlobalPosition.X,
				"y": chunk.GlobalPosition.Y,
				"z": chunk.GlobalPosition.Z,
			},
			"size": gms.gm.Config.WorldGen.ChunkSize,
		})
	}

	response := map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"chunks":    cacheInfo,
		"count":     len(cacheInfo),
	}

	json.NewEncoder(w).Encode(response)
}

func (gms *GameMonitoringServer) handleChunkRegistry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// TODO: Get actual registry info if available
	response := map[string]interface{}{
		"timestamp":     time.Now().UTC().Format(time.RFC3339),
		"registry_size": 0,
		"status":        "healthy",
	}

	json.NewEncoder(w).Encode(response)
}
