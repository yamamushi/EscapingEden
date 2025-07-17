# Monitoring Endpoints Guide

This guide covers all HTTP monitoring endpoints available in Escaping Eden for server health monitoring, performance tracking, and debugging.

## 🌐 Overview

Escaping Eden provides several HTTP endpoints for monitoring server health, performance metrics, and chunk loading statistics. These endpoints are essential for:

- **Server Health Monitoring** - Check if the server is running properly
- **Performance Tracking** - Monitor chunk loading and cache performance  
- **Debugging** - Investigate issues with world generation and loading
- **Profiling** - Analyze memory usage and performance bottlenecks

## 🔧 Enabling Monitoring

### pprof Profiling (Built-in)
The server includes Go's built-in pprof profiling endpoints. To enable them, uncomment this line in `main.go`:

```go
//go http.ListenAndServe("localhost:6060", nil)
```

Change it to:
```go
go http.ListenAndServe("localhost:6060", nil)
```

**Default URL**: `http://localhost:6060`

### Custom Monitoring Endpoints
The server also provides custom monitoring endpoints through the chunk monitoring system and metrics server.

## 📊 Available Endpoints

### 1. **pprof Profiling Endpoints** (Port 6060)

#### **Main pprof Dashboard**
```
http://localhost:6060/debug/pprof/
```
**Description**: Interactive profiling dashboard with links to all profiling tools.

#### **CPU Profiling**
```
http://localhost:6060/debug/pprof/profile
```
**Description**: 30-second CPU profile for performance analysis.
**Usage**: Download and analyze with `go tool pprof`

#### **Memory Profiling**
```
http://localhost:6060/debug/pprof/heap
```
**Description**: Current memory heap profile.
**Usage**: Analyze memory usage and potential leaks.

#### **Goroutine Analysis**
```
http://localhost:6060/debug/pprof/goroutine
```
**Description**: Current goroutine stack traces.
**Usage**: Debug concurrency issues and goroutine leaks.

#### **Memory Allocations**
```
http://localhost:6060/debug/pprof/allocs
```
**Description**: Memory allocation profile.
**Usage**: Track memory allocation patterns.

#### **Blocking Profile**
```
http://localhost:6060/debug/pprof/block
```
**Description**: Goroutine blocking profile.
**Usage**: Identify synchronization bottlenecks.

#### **Mutex Contention**
```
http://localhost:6060/debug/pprof/mutex
```
**Description**: Mutex contention profile.
**Usage**: Find lock contention issues.

### 2. **Chunk Monitoring Endpoints**

#### **Chunk Statistics**
```
http://localhost:8080/chunks/stats
```
**Method**: GET  
**Response**: JSON  
**Description**: Detailed chunk loading and cache statistics.

**Example Response**:
```json
{
  "total_chunks": 2500,
  "loaded_chunks": 1247,
  "cache_hits": 8934,
  "cache_misses": 156,
  "cache_size": 100,
  "load_time_avg_ms": 45.2,
  "generation_time_avg_ms": 234.7
}
```

#### **Chunk Health Check**
```
http://localhost:8080/chunks/health
```
**Method**: GET  
**Response**: JSON  
**Description**: Overall chunk system health status.

**Example Response**:
```json
{
  "status": "healthy",
  "world_validated": true,
  "chunks_loaded": 1247,
  "chunks_total": 2500,
  "health_percentage": 49.88
}
```

**HTTP Status Codes**:
- `200 OK` - System is healthy
- `202 Accepted` - System is incomplete but functional
- `503 Service Unavailable` - System is not ready
- `500 Internal Server Error` - System error

#### **Cache Statistics**
```
http://localhost:8080/chunks/cache
```
**Method**: GET  
**Response**: JSON  
**Description**: Detailed chunk cache performance metrics.

**Example Response**:
```json
{
  "cache_size": 100,
  "cache_used": 87,
  "cache_hit_rate": 98.3,
  "cache_miss_rate": 1.7,
  "evictions": 23,
  "memory_usage_mb": 145.7
}
```

#### **Chunk Preload Request**
```
http://localhost:8080/chunks/preload
```
**Method**: POST  
**Content-Type**: application/json  
**Description**: Request preloading of specific chunks.

**Request Body**:
```json
{
  "chunks": [
    {"x": 0, "y": 0, "z": 0},
    {"x": 1, "y": 0, "z": 0},
    {"x": 0, "y": 1, "z": 0}
  ]
}
```

**Response**:
```json
{
  "status": "success",
  "preloaded": 3,
  "failed": 0,
  "message": "Chunks preloaded successfully"
}
```

#### **Registry Information**
```
http://localhost:8080/chunks/registry
```
**Method**: GET  
**Response**: JSON  
**Description**: Chunk registry metadata and statistics.

**Example Response**:
```json
{
  "total_expected": 2500,
  "registered_chunks": 1247,
  "valid_chunks": 1200,
  "registry_version": 1,
  "last_updated": "2025-07-16T19:00:00Z"
}
```

### 3. **General Health Endpoints**

#### **Overall Health**
```
http://localhost:8080/health
```
**Method**: GET  
**Response**: JSON  
**Description**: Overall server health status.

#### **Liveness Probe**
```
http://localhost:8080/health/live
```
**Method**: GET  
**Response**: JSON  
**Description**: Simple liveness check for container orchestration.

#### **Readiness Probe**
```
http://localhost:8080/health/ready
```
**Method**: GET  
**Response**: JSON  
**Description**: Readiness check indicating if server can accept traffic.

#### **Metrics Endpoint**
```
http://localhost:8080/metrics
```
**Method**: GET  
**Response**: JSON  
**Description**: Comprehensive server metrics.

## 🛠️ Using the Endpoints

### Browser Access
Simply open your web browser and navigate to any endpoint:

```
http://localhost:6060/debug/pprof/
http://localhost:8080/chunks/stats
http://localhost:8080/health
```

### Command Line with curl
```bash
# Get chunk statistics
curl http://localhost:8080/chunks/stats | jq

# Check server health
curl http://localhost:8080/health

# Get CPU profile (saves to file)
curl http://localhost:6060/debug/pprof/profile > cpu.prof

# Preload specific chunks
curl -X POST http://localhost:8080/chunks/preload \
  -H "Content-Type: application/json" \
  -d '{"chunks": [{"x": 0, "y": 0, "z": 0}]}'
```

### Monitoring Scripts
```bash
#!/bin/bash
# Simple health check script
HEALTH=$(curl -s http://localhost:8080/health | jq -r '.status')
if [ "$HEALTH" != "healthy" ]; then
    echo "Server unhealthy: $HEALTH"
    exit 1
fi
echo "Server is healthy"
```

## 📈 Monitoring Best Practices

### Regular Health Checks
Set up automated health checks:
```bash
# Cron job for health monitoring (every 5 minutes)
*/5 * * * * curl -f http://localhost:8080/health > /dev/null || echo "Server down" | mail admin@example.com
```

### Performance Monitoring
Monitor key metrics regularly:
- **Chunk cache hit rate** - Should be >95%
- **Memory usage** - Watch for memory leaks
- **Load times** - Monitor chunk loading performance
- **Goroutine count** - Check for goroutine leaks

### Alerting Thresholds
Set up alerts for:
- **Health status** != "healthy"
- **Cache hit rate** < 90%
- **Memory usage** > 80% of available
- **Load time** > 1000ms average

## 🔍 Troubleshooting

### Common Issues

#### **Endpoints Not Accessible**
**Problem**: Cannot access monitoring endpoints  
**Solutions**:
1. Check if pprof is enabled in `main.go`
2. Verify server is running on expected ports
3. Check firewall settings
4. Ensure server started successfully

#### **Empty or Error Responses**
**Problem**: Endpoints return errors or empty data  
**Solutions**:
1. Check server logs for errors
2. Verify chunk system is initialized
3. Ensure world generation completed
4. Check database connectivity

#### **High Memory Usage**
**Problem**: Memory usage continuously increasing  
**Investigation**:
1. Check `/debug/pprof/heap` for memory leaks
2. Monitor `/chunks/cache` for cache issues
3. Review `/debug/pprof/goroutine` for goroutine leaks
4. Analyze allocation patterns with `/debug/pprof/allocs`

### Performance Analysis

#### **CPU Profiling**
```bash
# Capture 30-second CPU profile
curl http://localhost:6060/debug/pprof/profile > cpu.prof

# Analyze with go tool
go tool pprof cpu.prof
```

#### **Memory Analysis**
```bash
# Get heap profile
curl http://localhost:6060/debug/pprof/heap > heap.prof

# Analyze memory usage
go tool pprof heap.prof
```

#### **Goroutine Analysis**
```bash
# Get goroutine dump
curl http://localhost:6060/debug/pprof/goroutine > goroutines.prof

# Analyze goroutines
go tool pprof goroutines.prof
```

## 🔐 Security Considerations

### Access Control
- **Restrict access** to monitoring endpoints in production
- **Use reverse proxy** with authentication for external access
- **Monitor access logs** for unauthorized attempts
- **Consider VPN** for remote monitoring access

### Sensitive Information
- **Memory dumps** may contain sensitive data
- **Profiles** can reveal application internals
- **Limit profile retention** time
- **Secure profile storage** if saved

## 📊 Integration Examples

### Prometheus Integration
```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'escaping-eden'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 30s
```

### Grafana Dashboard
Create dashboards monitoring:
- Server health status
- Chunk loading performance
- Memory and CPU usage
- Cache hit rates
- Player connection metrics

### Docker Health Checks
```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1
```

This monitoring system provides comprehensive visibility into server performance and health, enabling proactive maintenance and quick issue resolution.