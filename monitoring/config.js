// Escaping Eden Monitoring Dashboard Configuration

const config = {
    // Game server configuration
    gameServerUrl: 'http://localhost:3000',
    monitoringServerUrl: 'http://localhost:8080',  // Monitoring endpoints are on port 8080
    pprofServerUrl: 'http://localhost:6060',       // pprof endpoints are on port 6060
    
    // Refresh intervals (in milliseconds)
    refreshInterval: 5000,        // Main dashboard refresh (5 seconds)
    fastRefreshInterval: 2000,    // Fast metrics (2 seconds)
    slowRefreshInterval: 30000,   // Slow metrics (30 seconds)
    
    // Chart configuration
    charts: {
        maxDataPoints: 50,        // Maximum points to show on time-series charts
        animationDuration: 750,   // Chart animation duration in ms
        responsive: true,
        maintainAspectRatio: false
    },
    
    // API endpoints - Real endpoints available on the server
    endpoints: {
        // Health endpoints (port 8080)
        health: '/health',
        healthLive: '/health/live',
        healthReady: '/health/ready',
        metrics: '/metrics',
        
        // Chunk monitoring endpoints (port 8080)
        chunksStats: '/chunks/stats',
        chunksHealth: '/chunks/health',
        chunksCache: '/chunks/cache',
        chunksPreload: '/chunks/preload',
        chunksRegistry: '/chunks/registry',
        
        // pprof profiling endpoints (port 6060)
        pprofIndex: ':6060/debug/pprof/',
        pprofProfile: ':6060/debug/pprof/profile',
        pprofHeap: ':6060/debug/pprof/heap',
        pprofGoroutine: ':6060/debug/pprof/goroutine',
        pprofAllocs: ':6060/debug/pprof/allocs',
        pprofBlock: ':6060/debug/pprof/block',
        pprofMutex: ':6060/debug/pprof/mutex',
        
        // Legacy endpoints (may not exist yet)
        statsPlayers: '/stats/players',
        statsWorld: '/stats/world',
        statsCommands: '/stats/commands',
        statsPerformance: '/stats/performance',
        websocket: '/ws/metrics',
        events: '/events/metrics'
    },
    
    // Alert thresholds
    alerts: {
        cpu: {
            warning: 70,    // CPU usage warning threshold (%)
            critical: 90    // CPU usage critical threshold (%)
        },
        memory: {
            warning: 80,    // Memory usage warning threshold (%)
            critical: 95    // Memory usage critical threshold (%)
        },
        responseTime: {
            warning: 1000,  // Response time warning threshold (ms)
            critical: 5000  // Response time critical threshold (ms)
        },
        errorRate: {
            warning: 0.05,  // Error rate warning threshold (5%)
            critical: 0.15  // Error rate critical threshold (15%)
        },
        players: {
            max: 1000       // Maximum expected players
        }
    },
    
    // Chart colors
    colors: {
        primary: '#0d6efd',
        success: '#198754',
        info: '#0dcaf0',
        warning: '#ffc107',
        danger: '#dc3545',
        secondary: '#6c757d',
        light: '#f8f9fa',
        dark: '#212529',
        
        // Chart color palette
        palette: [
            '#0d6efd', '#198754', '#ffc107', '#dc3545',
            '#0dcaf0', '#6f42c1', '#fd7e14', '#20c997',
            '#e83e8c', '#6c757d'
        ],
        
        // Gradient colors for area charts
        gradients: {
            blue: ['rgba(13, 110, 253, 0.2)', 'rgba(13, 110, 253, 0.05)'],
            green: ['rgba(25, 135, 84, 0.2)', 'rgba(25, 135, 84, 0.05)'],
            yellow: ['rgba(255, 193, 7, 0.2)', 'rgba(255, 193, 7, 0.05)'],
            red: ['rgba(220, 53, 69, 0.2)', 'rgba(220, 53, 69, 0.05)']
        }
    },
    
    // Feature flags
    features: {
        realTimeUpdates: true,      // Enable WebSocket/SSE updates
        worldHeatmap: true,         // Enable world activity heatmap
        alertNotifications: true,   // Enable browser notifications
        soundAlerts: false,         // Enable sound alerts
        exportData: true,           // Enable data export functionality
        darkMode: 'auto'            // 'auto', 'light', 'dark'
    },
    
    // Data retention
    dataRetention: {
        realTime: 300,      // Keep 5 minutes of real-time data
        shortTerm: 3600,    // Keep 1 hour of short-term data
        longTerm: 86400     // Keep 24 hours of long-term data
    },
    
    // Network configuration
    network: {
        timeout: 10000,             // Request timeout (ms)
        retryAttempts: 3,           // Number of retry attempts
        retryDelay: 1000,           // Delay between retries (ms)
        maxConcurrentRequests: 5    // Maximum concurrent API requests
    },
    
    // UI preferences
    ui: {
        compactMode: false,         // Use compact layout
        showTooltips: true,         // Show chart tooltips
        animateUpdates: true,       // Animate value changes
        autoHideAlerts: 30000,      // Auto-hide alerts after 30 seconds
        maxAlertsShown: 10          // Maximum alerts to show in sidebar
    }
};

// Environment-specific overrides
if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
    // Development environment
    config.gameServerUrl = 'http://localhost:3000';
    config.refreshInterval = 2000; // Faster refresh for development
} else if (window.location.hostname.includes('staging')) {
    // Staging environment
    config.gameServerUrl = 'https://staging-api.escapingeden.com';
} else {
    // Production environment
    config.gameServerUrl = 'https://api.escapingeden.com';
    config.refreshInterval = 10000; // Slower refresh for production
}

// Load user preferences from localStorage
function loadUserPreferences() {
    const saved = localStorage.getItem('eden-dashboard-config');
    if (saved) {
        try {
            const userConfig = JSON.parse(saved);
            Object.assign(config, userConfig);
        } catch (e) {
            console.warn('Failed to load user preferences:', e);
        }
    }
}

// Save user preferences to localStorage
function saveUserPreferences() {
    try {
        const userConfig = {
            gameServerUrl: config.gameServerUrl,
            refreshInterval: config.refreshInterval,
            features: config.features,
            ui: config.ui
        };
        localStorage.setItem('eden-dashboard-config', JSON.stringify(userConfig));
    } catch (e) {
        console.warn('Failed to save user preferences:', e);
    }
}

// Initialize configuration
loadUserPreferences();

// Export configuration
window.dashboardConfig = config;