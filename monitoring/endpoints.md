# Monitoring Endpoints Documentation

This document describes all available monitoring endpoints for the Escaping Eden game server.

## Base URL

All endpoints are relative to the game server base URL (default: `http://localhost:3000`)

## Health Check Endpoints

### GET /health
Basic health check endpoint.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "uptime": 3600,
  "version": "1.0.0"
}
```

### GET /health/detailed
Comprehensive health information.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "uptime": 3600,
  "version": "1.0.0",
  "components": {
    "database": {
      "status": "healthy",
      "connections": 5,
      "response_time_ms": 12
    },
    "game_manager": {
      "status": "healthy",
      "active_characters": 25,
      "loaded_chunks": 150
    },
    "network": {
      "status": "healthy",
      "active_connections": 30,
      "total_connections": 1250
    }
  }
}
```

## Metrics Endpoints

### GET /metrics
Prometheus-compatible metrics endpoint.

**Response Format:** Prometheus text format
```
# HELP eden_active_players Number of currently active players
# TYPE eden_active_players gauge
eden_active_players 25

# HELP eden_total_connections Total number of connections since startup
# TYPE eden_total_connections counter
eden_total_connections 1250

# HELP eden_command_duration_seconds Command execution duration
# TYPE eden_command_duration_seconds histogram
eden_command_duration_seconds_bucket{command="move",le="0.1"} 1000
eden_command_duration_seconds_bucket{command="move",le="0.5"} 1200
eden_command_duration_seconds_bucket{command="move",le="1.0"} 1250
```

### GET /metrics/json
JSON format metrics for easier consumption by web dashboards.

**Response:**
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "server": {
    "uptime_seconds": 3600,
    "cpu_usage_percent": 15.5,
    "memory_usage_bytes": 134217728,
    "memory_usage_percent": 12.8,
    "goroutines": 45,
    "gc_cycles": 12
  },
  "game": {
    "active_players": 25,
    "total_characters": 150,
    "loaded_chunks": 75,
    "cached_chunks": 100,
    "commands_per_second": 12.5,
    "average_response_time_ms": 25
  },
  "database": {
    "connections_active": 5,
    "connections_idle": 3,
    "queries_per_second": 8.2,
    "average_query_time_ms": 15,
    "cache_hit_ratio": 0.85
  },
  "network": {
    "active_connections": 30,
    "bytes_sent": 1048576,
    "bytes_received": 524288,
    "messages_per_second": 45.2
  }
}
```

## Statistics Endpoints

### GET /stats/players
Player-related statistics.

**Response:**
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "active_players": 25,
  "peak_players_today": 45,
  "total_registered": 1250,
  "players_by_level": {
    "1-10": 15,
    "11-20": 8,
    "21-30": 2
  },
  "geographic_distribution": {
    "US": 15,
    "EU": 8,
    "ASIA": 2
  }
}
```

### GET /stats/world
World and chunk statistics.

**Response:**
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "total_chunks": 1000,
  "loaded_chunks": 75,
  "cached_chunks": 100,
  "chunk_cache_hit_ratio": 0.92,
  "world_dimensions": {
    "x": 10,
    "y": 10,
    "z": 10
  },
  "chunk_size": 255,
  "popular_areas": [
    {
      "chunk_id": "0-0-0",
      "player_count": 8,
      "coordinates": {"x": 0, "y": 0, "z": 0}
    }
  ]
}
```

### GET /stats/commands
Command execution statistics.

**Response:**
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "total_commands": 50000,
  "commands_per_second": 12.5,
  "success_rate": 0.98,
  "by_type": {
    "move": {
      "count": 25000,
      "success_rate": 0.99,
      "avg_duration_ms": 15,
      "errors": {
        "blocked": 150,
        "invalid": 25
      }
    },
    "dig": {
      "count": 15000,
      "success_rate": 0.95,
      "avg_duration_ms": 45,
      "errors": {
        "no_tool": 500,
        "blocked": 250
      }
    },
    "build": {
      "count": 8000,
      "success_rate": 0.92,
      "avg_duration_ms": 85,
      "errors": {
        "no_materials": 400,
        "blocked": 200
      }
    }
  }
}
```

### GET /stats/performance
Performance metrics over time.

**Query Parameters:**
- `period`: Time period (1h, 6h, 24h, 7d) - default: 1h
- `resolution`: Data point resolution (1m, 5m, 15m, 1h) - default: 5m

**Response:**
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "period": "1h",
  "resolution": "5m",
  "data_points": [
    {
      "timestamp": "2024-01-15T09:30:00Z",
      "cpu_percent": 12.5,
      "memory_percent": 15.2,
      "active_players": 20,
      "commands_per_second": 8.5,
      "response_time_ms": 22
    },
    {
      "timestamp": "2024-01-15T09:35:00Z",
      "cpu_percent": 14.1,
      "memory_percent": 15.8,
      "active_players": 23,
      "commands_per_second": 11.2,
      "response_time_ms": 28
    }
  ]
}
```

## Real-time Endpoints

### WebSocket /ws/metrics
Real-time metrics stream via WebSocket.

**Connection:** `ws://localhost:3000/ws/metrics`

**Message Format:**
```json
{
  "type": "metrics_update",
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {
    "active_players": 25,
    "cpu_percent": 15.5,
    "memory_percent": 12.8,
    "commands_per_second": 12.5
  }
}
```

### Server-Sent Events /events/metrics
Alternative real-time stream using SSE.

**Connection:** `GET /events/metrics`
**Headers:** `Accept: text/event-stream`

**Event Format:**
```
event: metrics
data: {"active_players": 25, "cpu_percent": 15.5}

event: alert
data: {"level": "warning", "message": "High CPU usage detected"}
```

## Error Responses

All endpoints return consistent error responses:

```json
{
  "error": {
    "code": "ENDPOINT_NOT_FOUND",
    "message": "The requested endpoint does not exist",
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

**Common Error Codes:**
- `ENDPOINT_NOT_FOUND` - 404
- `INTERNAL_ERROR` - 500
- `RATE_LIMITED` - 429
- `UNAUTHORIZED` - 401

## Rate Limiting

- Health endpoints: 60 requests/minute
- Metrics endpoints: 30 requests/minute
- Statistics endpoints: 10 requests/minute
- Real-time connections: 5 concurrent per IP

## Authentication

Currently, all endpoints are public. Future versions may require API keys:

```
Authorization: Bearer your-api-key-here
```

## CORS Policy

The server accepts requests from:
- `localhost` (any port)
- Same origin as game server
- Configured allowed origins in server config