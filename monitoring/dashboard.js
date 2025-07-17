// Escaping Eden Monitoring Dashboard JavaScript

class MonitoringDashboard {
    constructor() {
        this.config = window.dashboardConfig;
        this.charts = {};
        this.data = {
            players: [],
            system: [],
            commands: [],
            database: [],
            alerts: []
        };
        this.isConnected = false;
        this.lastUpdate = null;
        this.refreshTimer = null;
        this.websocket = null;
        
        this.init();
    }
    
    init() {
        this.setupCharts();
        this.setupEventListeners();
        this.startDataRefresh();
        this.setupWebSocket();
        this.loadConfiguration();
    }
    
    setupCharts() {
        // Players Over Time Chart
        this.charts.players = new Chart(document.getElementById('playersChart'), {
            type: 'line',
            data: {
                labels: [],
                datasets: [{
                    label: 'Active Players',
                    data: [],
                    borderColor: this.config.colors.primary,
                    backgroundColor: this.createGradient('playersChart', this.config.colors.gradients.blue),
                    fill: true,
                    tension: 0.4
                }]
            },
            options: this.getChartOptions('Players')
        });
        
        // System Resources Chart
        this.charts.system = new Chart(document.getElementById('systemChart'), {
            type: 'line',
            data: {
                labels: [],
                datasets: [
                    {
                        label: 'CPU Usage (%)',
                        data: [],
                        borderColor: this.config.colors.danger,
                        backgroundColor: 'transparent',
                        yAxisID: 'y'
                    },
                    {
                        label: 'Memory Usage (%)',
                        data: [],
                        borderColor: this.config.colors.warning,
                        backgroundColor: 'transparent',
                        yAxisID: 'y'
                    }
                ]
            },
            options: {
                ...this.getChartOptions('System Resources'),
                scales: {
                    y: {
                        type: 'linear',
                        display: true,
                        position: 'left',
                        max: 100,
                        ticks: {
                            callback: function(value) {
                                return value + '%';
                            }
                        }
                    }
                }
            }
        });
        
        // Commands Chart
        this.charts.commands = new Chart(document.getElementById('commandsChart'), {
            type: 'bar',
            data: {
                labels: [],
                datasets: [
                    {
                        label: 'Successful',
                        data: [],
                        backgroundColor: this.config.colors.success
                    },
                    {
                        label: 'Failed',
                        data: [],
                        backgroundColor: this.config.colors.danger
                    }
                ]
            },
            options: {
                ...this.getChartOptions('Commands per Minute'),
                scales: {
                    x: {
                        stacked: true
                    },
                    y: {
                        stacked: true
                    }
                }
            }
        });
        
        // Database Performance Chart
        this.charts.database = new Chart(document.getElementById('databaseChart'), {
            type: 'line',
            data: {
                labels: [],
                datasets: [
                    {
                        label: 'Query Time (ms)',
                        data: [],
                        borderColor: this.config.colors.info,
                        backgroundColor: 'transparent',
                        yAxisID: 'y'
                    },
                    {
                        label: 'Queries/sec',
                        data: [],
                        borderColor: this.config.colors.success,
                        backgroundColor: 'transparent',
                        yAxisID: 'y1'
                    }
                ]
            },
            options: {
                ...this.getChartOptions('Database Performance'),
                scales: {
                    y: {
                        type: 'linear',
                        display: true,
                        position: 'left'
                    },
                    y1: {
                        type: 'linear',
                        display: true,
                        position: 'right',
                        grid: {
                            drawOnChartArea: false
                        }
                    }
                }
            }
        });
        
        // World Activity Heatmap
        this.setupWorldHeatmap();
    }
    
    setupWorldHeatmap() {
        const canvas = document.getElementById('worldChart');
        const ctx = canvas.getContext('2d');
        
        // Simple heatmap implementation
        this.worldHeatmap = {
            canvas: canvas,
            ctx: ctx,
            data: [],
            render: (data) => {
                const width = canvas.width;
                const height = canvas.height;
                
                // Clear canvas
                ctx.clearRect(0, 0, width, height);
                
                // Draw grid and activity data
                if (data && data.length > 0) {
                    const cellSize = Math.min(width / 10, height / 10);
                    
                    data.forEach(chunk => {
                        const x = chunk.coordinates.x * cellSize;
                        const y = chunk.coordinates.y * cellSize;
                        const intensity = Math.min(chunk.player_count / 10, 1);
                        
                        ctx.fillStyle = `rgba(13, 110, 253, ${intensity})`;
                        ctx.fillRect(x, y, cellSize, cellSize);
                        
                        // Draw player count
                        if (chunk.player_count > 0) {
                            ctx.fillStyle = '#fff';
                            ctx.font = '12px Arial';
                            ctx.textAlign = 'center';
                            ctx.fillText(chunk.player_count, x + cellSize/2, y + cellSize/2);
                        }
                    });
                }
            }
        };
    }
    
    getChartOptions(title) {
        return {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                title: {
                    display: false
                },
                legend: {
                    display: true,
                    position: 'top'
                }
            },
            scales: {
                x: {
                    display: true,
                    grid: {
                        display: false
                    }
                },
                y: {
                    display: true,
                    beginAtZero: true
                }
            },
            animation: {
                duration: this.config.charts.animationDuration
            }
        };
    }
    
    createGradient(canvasId, colors) {
        const canvas = document.getElementById(canvasId);
        const ctx = canvas.getContext('2d');
        const gradient = ctx.createLinearGradient(0, 0, 0, canvas.height);
        gradient.addColorStop(0, colors[0]);
        gradient.addColorStop(1, colors[1]);
        return gradient;
    }
    
    setupEventListeners() {
        // Configuration modal
        document.getElementById('save-config').addEventListener('click', () => {
            this.saveConfiguration();
        });
        
        // Auto-refresh toggle
        document.getElementById('auto-refresh').addEventListener('change', (e) => {
            if (e.target.checked) {
                this.startDataRefresh();
            } else {
                this.stopDataRefresh();
            }
        });
        
        // Window visibility change
        document.addEventListener('visibilitychange', () => {
            if (document.hidden) {
                this.stopDataRefresh();
            } else {
                this.startDataRefresh();
            }
        });
    }
    
    setupWebSocket() {
        if (!this.config.features.realTimeUpdates) return;
        
        try {
            const wsUrl = this.config.gameServerUrl.replace('http', 'ws') + this.config.endpoints.websocket;
            this.websocket = new WebSocket(wsUrl);
            
            this.websocket.onopen = () => {
                console.log('WebSocket connected');
                this.updateConnectionStatus(true);
            };
            
            this.websocket.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data);
                    this.handleRealtimeUpdate(data);
                } catch (e) {
                    console.error('Failed to parse WebSocket message:', e);
                }
            };
            
            this.websocket.onclose = () => {
                console.log('WebSocket disconnected');
                this.updateConnectionStatus(false);
                
                // Attempt to reconnect after 5 seconds
                setTimeout(() => {
                    this.setupWebSocket();
                }, 5000);
            };
            
            this.websocket.onerror = (error) => {
                console.error('WebSocket error:', error);
                this.updateConnectionStatus(false);
            };
        } catch (e) {
            console.warn('WebSocket not supported or failed to connect:', e);
        }
    }
    
    handleRealtimeUpdate(data) {
        if (data.type === 'metrics_update') {
            this.updateMetrics(data.data);
        } else if (data.type === 'alert') {
            this.addAlert(data.data);
        }
    }
    
    startDataRefresh() {
        this.stopDataRefresh();
        
        // Initial load
        this.refreshData();
        
        // Set up periodic refresh
        this.refreshTimer = setInterval(() => {
            this.refreshData();
        }, this.config.refreshInterval);
    }
    
    stopDataRefresh() {
        if (this.refreshTimer) {
            clearInterval(this.refreshTimer);
            this.refreshTimer = null;
        }
    }
    
    async refreshData() {
        try {
            this.updateConnectionStatus(true);
            
            // Fetch real server endpoints in parallel
            const promises = [];
            const results = {};
            
            // Try to fetch from real endpoints first
            try {
                promises.push(
                    this.fetchHealth().then(data => results.health = data).catch(e => console.warn('Health endpoint failed:', e)),
                    this.fetchChunkStats().then(data => results.chunkStats = data).catch(e => console.warn('Chunk stats endpoint failed:', e)),
                    this.fetchChunkHealth().then(data => results.chunkHealth = data).catch(e => console.warn('Chunk health endpoint failed:', e)),
                    this.fetchChunkCache().then(data => results.chunkCache = data).catch(e => console.warn('Chunk cache endpoint failed:', e)),
                    this.fetchChunkRegistry().then(data => results.chunkRegistry = data).catch(e => console.warn('Chunk registry endpoint failed:', e))
                );
                
                // Also try legacy endpoints (may not exist yet)
                promises.push(
                    this.fetchMetrics().then(data => results.metrics = data).catch(e => console.warn('Metrics endpoint failed:', e)),
                    this.fetchPlayerStats().then(data => results.players = data).catch(e => console.warn('Player stats endpoint failed:', e)),
                    this.fetchWorldStats().then(data => results.world = data).catch(e => console.warn('World stats endpoint failed:', e)),
                    this.fetchCommandStats().then(data => results.commands = data).catch(e => console.warn('Command stats endpoint failed:', e))
                );
                
                await Promise.all(promises);
                
                // Update UI with available data
                if (results.health) this.updateHealthStatus(results.health);
                if (results.chunkStats) this.updateChunkStats(results.chunkStats);
                if (results.chunkHealth) this.updateChunkHealth(results.chunkHealth);
                if (results.chunkCache) this.updateChunkCache(results.chunkCache);
                if (results.chunkRegistry) this.updateChunkRegistry(results.chunkRegistry);
                
                // Legacy endpoints
                if (results.metrics) this.updateMetrics(results.metrics);
                if (results.players) this.updatePlayerStats(results.players);
                if (results.world) this.updateWorldStats(results.world);
                if (results.commands) this.updateCommandStats(results.commands);
                
                this.lastUpdate = new Date();
                this.updateLastUpdateTime();
                
            } catch (error) {
                console.error('Failed to refresh data:', error);
                this.updateConnectionStatus(false);
                this.addAlert({
                    level: 'danger',
                    message: 'Failed to fetch data from server: ' + error.message,
                    timestamp: new Date()
                });
            }
            
        } catch (error) {
            console.error('Critical error in refreshData:', error);
            this.updateConnectionStatus(false);
        }
    }
    
    async fetchMetrics() {
        const response = await fetch(this.config.gameServerUrl + this.config.endpoints.metrics);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return await response.json();
    }
    
    async fetchPlayerStats() {
        const response = await fetch(this.config.gameServerUrl + this.config.endpoints.statsPlayers);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return await response.json();
    }
    
    async fetchWorldStats() {
        const response = await fetch(this.config.gameServerUrl + this.config.endpoints.statsWorld);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return await response.json();
    }
    
    async fetchCommandStats() {
        const response = await fetch(this.config.gameServerUrl + this.config.endpoints.statsCommands);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return await response.json();
    }
    
    // New methods for real server endpoints
    async fetchHealth() {
        const response = await fetch(this.config.monitoringServerUrl + this.config.endpoints.health);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return await response.json();
    }
    
    async fetchHealthLive() {
        const response = await fetch(this.config.monitoringServerUrl + this.config.endpoints.healthLive);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return await response.json();
    }
    
    async fetchHealthReady() {
        const response = await fetch(this.config.monitoringServerUrl + this.config.endpoints.healthReady);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return await response.json();
    }
    
    async fetchChunkStats() {
        const response = await fetch(this.config.monitoringServerUrl + this.config.endpoints.chunksStats);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return await response.json();
    }
    
    async fetchChunkHealth() {
        const response = await fetch(this.config.monitoringServerUrl + this.config.endpoints.chunksHealth);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return await response.json();
    }
    
    async fetchChunkCache() {
        const response = await fetch(this.config.monitoringServerUrl + this.config.endpoints.chunksCache);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return await response.json();
    }
    
    async fetchChunkRegistry() {
        const response = await fetch(this.config.monitoringServerUrl + this.config.endpoints.chunksRegistry);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return await response.json();
    }
    
    async preloadChunks(chunks) {
        const response = await fetch(this.config.monitoringServerUrl + this.config.endpoints.chunksPreload, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ chunks })
        });
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return await response.json();
    }
    
    updateMetrics(data) {
        // Update status cards
        document.getElementById('active-players').textContent = data.game?.active_players || '--';
        document.getElementById('server-uptime').textContent = this.formatUptime(data.server?.uptime_seconds);
        document.getElementById('commands-per-sec').textContent = (data.game?.commands_per_second || 0).toFixed(1);
        document.getElementById('response-time').textContent = (data.game?.average_response_time_ms || 0) + 'ms';
        
        // Update charts with new data
        this.addDataPoint(this.charts.players, new Date(), data.game?.active_players || 0);
        this.addDataPoint(this.charts.system, new Date(), [
            data.server?.cpu_usage_percent || 0,
            data.server?.memory_usage_percent || 0
        ]);
        
        // Check for alerts
        this.checkAlerts(data);
    }
    
    updatePlayerStats(data) {
        // Update player-related UI elements
        if (data.peak_players_today) {
            // Update peak players display if needed
        }
    }
    
    updateWorldStats(data) {
        // Update world heatmap
        if (this.config.features.worldHeatmap && data.popular_areas) {
            this.worldHeatmap.render(data.popular_areas);
        }
    }
    
    updateCommandStats(data) {
        if (!data.by_type) return;
        
        const tbody = document.querySelector('#commands-table tbody');
        tbody.innerHTML = '';
        
        Object.entries(data.by_type).forEach(([command, stats]) => {
            const row = tbody.insertRow();
            row.innerHTML = `
                <td><strong>${command}</strong></td>
                <td>${stats.count.toLocaleString()}</td>
                <td>
                    <span class="badge ${this.getSuccessRateClass(stats.success_rate)}">
                        ${(stats.success_rate * 100).toFixed(1)}%
                    </span>
                </td>
                <td>${stats.avg_duration_ms}ms</td>
                <td>${Object.values(stats.errors || {}).reduce((a, b) => a + b, 0)}</td>
                <td>
                    <div class="progress" style="height: 20px;">
                        <div class="progress-bar" style="width: ${Math.min(stats.count / 1000 * 100, 100)}%"></div>
                    </div>
                </td>
            `;
        });
        
        // Update commands chart
        const labels = Object.keys(data.by_type);
        const successful = labels.map(cmd => data.by_type[cmd].count * data.by_type[cmd].success_rate);
        const failed = labels.map(cmd => data.by_type[cmd].count * (1 - data.by_type[cmd].success_rate));
        
        this.charts.commands.data.labels = labels;
        this.charts.commands.data.datasets[0].data = successful;
        this.charts.commands.data.datasets[1].data = failed;
        this.charts.commands.update();
    }
    
    // New update methods for real server endpoints
    updateHealthStatus(data) {
        // Update overall health status
        const healthElement = document.getElementById('server-health');
        if (healthElement) {
            const status = data.status || 'unknown';
            const statusClass = status === 'healthy' ? 'success' : 
                               status === 'degraded' ? 'warning' : 'danger';
            
            healthElement.innerHTML = `
                <span class="badge bg-${statusClass}">
                    <i class="fas fa-${status === 'healthy' ? 'check' : 'exclamation-triangle'} me-1"></i>
                    ${status.toUpperCase()}
                </span>
            `;
        }
        
        // Update health details if available
        if (data.components) {
            this.updateHealthComponents(data.components);
        }
    }
    
    updateHealthComponents(components) {
        const container = document.getElementById('health-components');
        if (!container) return;
        
        container.innerHTML = Object.entries(components).map(([name, component]) => {
            const statusClass = component.status === 'healthy' ? 'success' : 
                               component.status === 'degraded' ? 'warning' : 'danger';
            
            return `
                <div class="health-component mb-2">
                    <div class="d-flex justify-content-between align-items-center">
                        <span class="fw-bold">${name}</span>
                        <span class="badge bg-${statusClass}">${component.status}</span>
                    </div>
                    ${component.message ? `<small class="text-muted">${component.message}</small>` : ''}
                </div>
            `;
        }).join('');
    }
    
    updateChunkStats(data) {
        // Update chunk statistics cards
        document.getElementById('total-chunks').textContent = (data.total_chunks || 0).toLocaleString();
        document.getElementById('loaded-chunks').textContent = (data.loaded_chunks || 0).toLocaleString();
        document.getElementById('cache-hits').textContent = (data.cache_hits || 0).toLocaleString();
        document.getElementById('cache-misses').textContent = (data.cache_misses || 0).toLocaleString();
        
        // Calculate and display cache hit rate
        const hitRate = data.cache_hits && data.cache_misses ? 
            (data.cache_hits / (data.cache_hits + data.cache_misses) * 100) : 0;
        document.getElementById('cache-hit-rate').textContent = hitRate.toFixed(1) + '%';
        
        // Update load times
        document.getElementById('avg-load-time').textContent = (data.load_time_avg_ms || 0).toFixed(1) + 'ms';
        document.getElementById('avg-gen-time').textContent = (data.generation_time_avg_ms || 0).toFixed(1) + 'ms';
        
        // Add chunk performance data to charts
        this.addDataPoint(this.charts.database, new Date(), [
            data.load_time_avg_ms || 0,
            (data.cache_hits || 0) / 60 // Approximate chunks per second
        ]);
    }
    
    updateChunkHealth(data) {
        const healthElement = document.getElementById('chunk-health');
        if (healthElement) {
            const status = data.status || 'unknown';
            const statusClass = status === 'healthy' ? 'success' : 
                               status === 'incomplete' ? 'warning' : 'danger';
            
            healthElement.innerHTML = `
                <span class="badge bg-${statusClass}">
                    <i class="fas fa-${status === 'healthy' ? 'check' : 'exclamation-triangle'} me-1"></i>
                    ${status.toUpperCase()}
                </span>
            `;
        }
        
        // Update chunk health percentage
        const healthPercent = document.getElementById('chunk-health-percent');
        if (healthPercent && data.health_percentage !== undefined) {
            healthPercent.textContent = data.health_percentage.toFixed(1) + '%';
        }
        
        // Update world validation status
        const worldValidated = document.getElementById('world-validated');
        if (worldValidated) {
            const validated = data.world_validated;
            worldValidated.innerHTML = `
                <span class="badge bg-${validated ? 'success' : 'warning'}">
                    <i class="fas fa-${validated ? 'check' : 'clock'} me-1"></i>
                    ${validated ? 'VALIDATED' : 'PENDING'}
                </span>
            `;
        }
    }
    
    updateChunkCache(data) {
        // Update cache statistics
        document.getElementById('cache-size').textContent = (data.cache_size || 0).toLocaleString();
        document.getElementById('cache-used').textContent = (data.cache_used || 0).toLocaleString();
        document.getElementById('cache-evictions').textContent = (data.evictions || 0).toLocaleString();
        
        // Update memory usage
        const memoryUsage = document.getElementById('cache-memory-usage');
        if (memoryUsage && data.memory_usage_mb !== undefined) {
            memoryUsage.textContent = data.memory_usage_mb.toFixed(1) + ' MB';
        }
        
        // Update cache utilization progress bar
        const utilizationBar = document.getElementById('cache-utilization-bar');
        if (utilizationBar && data.cache_size && data.cache_used) {
            const utilization = (data.cache_used / data.cache_size) * 100;
            utilizationBar.style.width = utilization + '%';
            utilizationBar.textContent = utilization.toFixed(1) + '%';
            
            // Change color based on utilization
            utilizationBar.className = `progress-bar ${
                utilization > 90 ? 'bg-danger' : 
                utilization > 70 ? 'bg-warning' : 'bg-success'
            }`;
        }
    }
    
    updateChunkRegistry(data) {
        // Update registry statistics
        document.getElementById('registry-total').textContent = (data.total_expected || 0).toLocaleString();
        document.getElementById('registry-registered').textContent = (data.registered_chunks || 0).toLocaleString();
        document.getElementById('registry-valid').textContent = (data.valid_chunks || 0).toLocaleString();
        
        // Update registry version and last updated
        document.getElementById('registry-version').textContent = data.registry_version || '--';
        
        const lastUpdated = document.getElementById('registry-last-updated');
        if (lastUpdated && data.last_updated) {
            const updateTime = new Date(data.last_updated);
            lastUpdated.textContent = updateTime.toLocaleString();
        }
        
        // Calculate and display registry completion percentage
        const completionPercent = document.getElementById('registry-completion');
        if (completionPercent && data.total_expected && data.registered_chunks) {
            const completion = (data.registered_chunks / data.total_expected) * 100;
            completionPercent.textContent = completion.toFixed(1) + '%';
            
            // Update progress bar
            const progressBar = document.getElementById('registry-progress-bar');
            if (progressBar) {
                progressBar.style.width = completion + '%';
                progressBar.className = `progress-bar ${
                    completion === 100 ? 'bg-success' : 
                    completion > 80 ? 'bg-info' : 'bg-warning'
                }`;
            }
        }
    }
    
    addDataPoint(chart, timestamp, value) {
        const label = timestamp.toLocaleTimeString();
        
        chart.data.labels.push(label);
        
        if (Array.isArray(value)) {
            value.forEach((val, index) => {
                chart.data.datasets[index].data.push(val);
            });
        } else {
            chart.data.datasets[0].data.push(value);
        }
        
        // Keep only the last N data points
        const maxPoints = this.config.charts.maxDataPoints;
        if (chart.data.labels.length > maxPoints) {
            chart.data.labels.shift();
            chart.data.datasets.forEach(dataset => {
                dataset.data.shift();
            });
        }
        
        chart.update('none'); // Update without animation for real-time feel
    }
    
    checkAlerts(data) {
        const alerts = [];
        
        // CPU usage alert
        if (data.server?.cpu_usage_percent > this.config.alerts.cpu.critical) {
            alerts.push({
                level: 'danger',
                message: `Critical CPU usage: ${data.server.cpu_usage_percent.toFixed(1)}%`,
                timestamp: new Date()
            });
        } else if (data.server?.cpu_usage_percent > this.config.alerts.cpu.warning) {
            alerts.push({
                level: 'warning',
                message: `High CPU usage: ${data.server.cpu_usage_percent.toFixed(1)}%`,
                timestamp: new Date()
            });
        }
        
        // Memory usage alert
        if (data.server?.memory_usage_percent > this.config.alerts.memory.critical) {
            alerts.push({
                level: 'danger',
                message: `Critical memory usage: ${data.server.memory_usage_percent.toFixed(1)}%`,
                timestamp: new Date()
            });
        } else if (data.server?.memory_usage_percent > this.config.alerts.memory.warning) {
            alerts.push({
                level: 'warning',
                message: `High memory usage: ${data.server.memory_usage_percent.toFixed(1)}%`,
                timestamp: new Date()
            });
        }
        
        // Response time alert
        if (data.game?.average_response_time_ms > this.config.alerts.responseTime.critical) {
            alerts.push({
                level: 'danger',
                message: `Critical response time: ${data.game.average_response_time_ms}ms`,
                timestamp: new Date()
            });
        } else if (data.game?.average_response_time_ms > this.config.alerts.responseTime.warning) {
            alerts.push({
                level: 'warning',
                message: `High response time: ${data.game.average_response_time_ms}ms`,
                timestamp: new Date()
            });
        }
        
        alerts.forEach(alert => this.addAlert(alert));
    }
    
    addAlert(alert) {
        this.data.alerts.unshift(alert);
        
        // Keep only recent alerts
        if (this.data.alerts.length > this.config.ui.maxAlertsShown) {
            this.data.alerts = this.data.alerts.slice(0, this.config.ui.maxAlertsShown);
        }
        
        this.updateAlertsDisplay();
        
        // Show browser notification if enabled
        if (this.config.features.alertNotifications && 'Notification' in window) {
            if (Notification.permission === 'granted') {
                new Notification('Escaping Eden Alert', {
                    body: alert.message,
                    icon: '/favicon.ico'
                });
            }
        }
    }
    
    updateAlertsDisplay() {
        const container = document.getElementById('alerts-container');
        
        if (this.data.alerts.length === 0) {
            container.innerHTML = `
                <div class="text-muted text-center">
                    <i class="fas fa-check-circle fa-2x mb-2"></i>
                    <p>No alerts at this time</p>
                </div>
            `;
            return;
        }
        
        container.innerHTML = this.data.alerts.map(alert => `
            <div class="alert-item alert-${alert.level}">
                <div class="d-flex justify-content-between align-items-start">
                    <div>
                        <i class="fas fa-${this.getAlertIcon(alert.level)} me-2"></i>
                        ${alert.message}
                    </div>
                    <small class="alert-timestamp">${alert.timestamp.toLocaleTimeString()}</small>
                </div>
            </div>
        `).join('');
    }
    
    getAlertIcon(level) {
        switch (level) {
            case 'danger': return 'exclamation-triangle';
            case 'warning': return 'exclamation-circle';
            case 'info': return 'info-circle';
            default: return 'bell';
        }
    }
    
    getSuccessRateClass(rate) {
        if (rate >= 0.95) return 'success-rate-high';
        if (rate >= 0.8) return 'success-rate-medium';
        return 'success-rate-low';
    }
    
    formatUptime(seconds) {
        if (!seconds) return '--';
        
        const days = Math.floor(seconds / 86400);
        const hours = Math.floor((seconds % 86400) / 3600);
        const minutes = Math.floor((seconds % 3600) / 60);
        
        if (days > 0) return `${days}d ${hours}h`;
        if (hours > 0) return `${hours}h ${minutes}m`;
        return `${minutes}m`;
    }
    
    updateConnectionStatus(connected) {
        this.isConnected = connected;
        const statusElement = document.getElementById('connection-status');
        
        if (connected) {
            statusElement.className = 'badge bg-success me-3';
            statusElement.innerHTML = '<i class="fas fa-circle me-1"></i>Connected';
        } else {
            statusElement.className = 'badge bg-danger me-3';
            statusElement.innerHTML = '<i class="fas fa-circle me-1"></i>Disconnected';
        }
    }
    
    updateLastUpdateTime() {
        if (this.lastUpdate) {
            document.getElementById('last-update').textContent = 
                `Last update: ${this.lastUpdate.toLocaleTimeString()}`;
        }
    }
    
    loadConfiguration() {
        document.getElementById('server-url').value = this.config.gameServerUrl;
        document.getElementById('refresh-interval').value = this.config.refreshInterval / 1000;
        document.getElementById('auto-refresh').checked = true;
    }
    
    saveConfiguration() {
        this.config.gameServerUrl = document.getElementById('server-url').value;
        this.config.refreshInterval = parseInt(document.getElementById('refresh-interval').value) * 1000;
        
        // Save to localStorage
        saveUserPreferences();
        
        // Restart data refresh with new settings
        this.startDataRefresh();
        
        // Close modal
        const modal = bootstrap.Modal.getInstance(document.getElementById('configModal'));
        modal.hide();
        
        // Show success message
        this.addAlert({
            level: 'info',
            message: 'Configuration saved successfully',
            timestamp: new Date()
        });
    }
}

// Initialize dashboard when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.dashboard = new MonitoringDashboard();
    
    // Request notification permission
    if ('Notification' in window && Notification.permission === 'default') {
        Notification.requestPermission();
    }
});