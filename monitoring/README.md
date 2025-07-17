# Escaping Eden Monitoring Dashboard

A comprehensive real-time monitoring dashboard for the Escaping Eden game server.

## Features

- **Real-time Metrics**: Live graphs updating every 5-30 seconds
- **Server Health**: CPU, memory, connection status
- **Game Statistics**: Active players, world chunks, command metrics
- **Performance Monitoring**: Response times, error rates, throughput
- **Database Metrics**: Query performance, connection pool status
- **Network Analytics**: Connection patterns, bandwidth usage

## Quick Start

1. **Start the monitoring server:**
   ```bash
   cd monitoring
   python3 -m http.server 8080
   ```

2. **Open the dashboard:**
   ```
   http://localhost:8080
   ```

3. **Configure endpoints:**
   Edit `config.js` to point to your game server endpoints.

## Dashboard Components

### System Health
- CPU usage and load average
- Memory consumption and garbage collection
- Disk I/O and storage usage
- Network interface statistics

### Game Metrics
- Active player count over time
- World chunk loading/unloading
- Command execution rates and success/failure
- Character movement patterns

### Performance Analytics
- Request/response latencies
- Error rate trends
- Database query performance
- Cache hit/miss ratios

### Real-time Alerts
- Server overload warnings
- Database connection issues
- High error rate notifications
- Performance degradation alerts

## Configuration

The dashboard reads configuration from `config.js`:

```javascript
const config = {
    gameServerUrl: 'http://localhost:3000',
    refreshInterval: 5000, // 5 seconds
    endpoints: {
        health: '/health',
        metrics: '/metrics',
        stats: '/stats'
    }
};
```

## Endpoints Documentation

See `endpoints.md` for detailed API documentation.

## Browser Compatibility

- Chrome 80+
- Firefox 75+
- Safari 13+
- Edge 80+

## Dependencies

- Chart.js for graphs
- Socket.io for real-time updates (optional)
- Bootstrap for responsive UI