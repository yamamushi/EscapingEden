# Configuration Guide

This guide covers all configuration options available in Escaping Eden's server configuration file.

## 📁 Configuration File

The server uses a configuration file (default: `server.conf`) in TOML format. You can specify a different config file using the `-c` flag:

```bash
./EscapingEden -c custom-config.conf
```

## 🔧 Configuration Sections

### [Server] Section

Controls basic server settings and network configuration.

```toml
[server]
host = "localhost"              # Server bind address
port = "8080"                   # Server port
shutdown_timeout = 30           # Graceful shutdown timeout in seconds
max_connections = 1000          # Maximum concurrent connections
read_timeout = 30               # Connection read timeout in seconds
write_timeout = 30              # Connection write timeout in seconds
```

**Options:**
- **host**: IP address or hostname to bind to
  - `"localhost"` - Local access only
  - `"0.0.0.0"` - Accept connections from any IP
  - `"192.168.1.100"` - Specific IP address
- **port**: TCP port number (1-65535)
- **shutdown_timeout**: Time to wait for graceful shutdown
- **max_connections**: Maximum simultaneous client connections
- **read_timeout**: Timeout for reading client data
- **write_timeout**: Timeout for writing to clients

### [Database] Section

Database connection and storage settings.

```toml
[database]
type = "boltdb"                 # Database type
path = "./eden.db"              # Database file path
backup_interval = 3600          # Backup interval in seconds
max_connections = 10            # Maximum database connections
connection_timeout = 30         # Connection timeout in seconds
```

**Database Types:**
- **boltdb**: Embedded key-value database (default)
- **mongodb**: MongoDB database (future support)

**Options:**
- **path**: File path for BoltDB database
- **backup_interval**: Automatic backup frequency
- **max_connections**: Connection pool size
- **connection_timeout**: Database connection timeout

### [WorldGen] Section

World generation parameters and settings.

```toml
[worldgen]
dimensions = "50x50x3"          # World size (width x height x depth)
seed = 12345                    # Random seed for generation
algorithm = "perlin"            # Terrain generation algorithm
biome_distribution = "random"   # Biome placement method
structure_density = 0.1         # Density of generated structures
resource_abundance = 1.0        # Resource generation multiplier
```

**Dimension Format:**
- Format: `"WIDTHxHEIGHTxDEPTH"`
- Example: `"100x100x5"` = 100x100 chunks, 5 levels deep
- Minimum: `"2x2x1"`
- Maximum: Limited by available disk space and memory

**Generation Algorithms:**
- **perlin**: Perlin noise-based terrain (default)
- **simplex**: Simplex noise terrain
- **cellular**: Cellular automata caves
- **flat**: Flat world for testing

**Biome Distribution:**
- **random**: Random biome placement
- **climate**: Climate-based biome zones
- **elevation**: Elevation-based biomes

### [Logging] Section

Logging configuration and output settings.

```toml
[logging]
level = "info"                  # Log level
output = "file"                 # Log output destination
file_path = "./server.log"      # Log file path
max_file_size = 100             # Max log file size in MB
max_backups = 5                 # Number of backup log files
max_age = 30                    # Max age of log files in days
compress = true                 # Compress old log files
```

**Log Levels:**
- **debug**: Detailed debugging information
- **info**: General information (default)
- **warn**: Warning messages only
- **error**: Error messages only
- **fatal**: Fatal errors only

**Output Destinations:**
- **file**: Write to log file
- **console**: Write to console/stdout
- **both**: Write to both file and console

### [Security] Section

Security and authentication settings.

```toml
[security]
enable_rate_limiting = true     # Enable connection rate limiting
max_requests_per_minute = 60    # Rate limit threshold
enable_ip_whitelist = false     # Enable IP address whitelist
whitelist_file = "whitelist.txt" # IP whitelist file path
session_timeout = 3600          # Session timeout in seconds
password_min_length = 8         # Minimum password length
require_email_verification = true # Require email verification
```

**Security Options:**
- **enable_rate_limiting**: Prevent connection flooding
- **max_requests_per_minute**: Rate limit threshold per IP
- **enable_ip_whitelist**: Restrict access to specific IPs
- **session_timeout**: Automatic session expiration
- **password_min_length**: Minimum password requirements
- **require_email_verification**: Email verification for new accounts

### [Performance] Section

Performance tuning and optimization settings.

```toml
[performance]
chunk_cache_size = 100          # Number of chunks to cache in memory
preload_radius = 2              # Chunks to preload around players
gc_interval = 300               # Garbage collection interval in seconds
max_goroutines = 1000           # Maximum concurrent goroutines
buffer_size = 8192              # Network buffer size in bytes
compression_level = 6           # Data compression level (1-9)
```

**Performance Options:**
- **chunk_cache_size**: Memory cache for map chunks
- **preload_radius**: Chunks to load around active players
- **gc_interval**: Force garbage collection frequency
- **max_goroutines**: Limit concurrent operations
- **buffer_size**: Network I/O buffer size
- **compression_level**: Balance between speed and compression

### [Discord] Section

Discord bot integration settings.

```toml
[discord]
enabled = true                  # Enable Discord integration
bot_token = "your_bot_token"    # Discord bot token
guild_id = "your_guild_id"      # Discord server ID
admin_role = "Admin"            # Admin role name
moderator_role = "Moderator"    # Moderator role name
log_channel = "server-logs"     # Log channel name
```

**Discord Integration:**
- **enabled**: Enable/disable Discord features
- **bot_token**: Discord application bot token
- **guild_id**: Discord server (guild) ID
- **admin_role**: Role name for admin permissions
- **moderator_role**: Role name for moderator permissions
- **log_channel**: Channel for server log messages

## 📝 Example Configuration Files

### Development Configuration
```toml
[server]
host = "localhost"
port = "8080"
shutdown_timeout = 5

[database]
type = "boltdb"
path = "./dev-eden.db"

[worldgen]
dimensions = "10x10x3"
seed = 12345

[logging]
level = "debug"
output = "console"

[performance]
chunk_cache_size = 50
preload_radius = 1
```

### Production Configuration
```toml
[server]
host = "0.0.0.0"
port = "8080"
shutdown_timeout = 30
max_connections = 500

[database]
type = "boltdb"
path = "./production-eden.db"
backup_interval = 1800

[worldgen]
dimensions = "100x100x5"
seed = 98765

[logging]
level = "info"
output = "file"
file_path = "./logs/server.log"
max_file_size = 100
max_backups = 10

[security]
enable_rate_limiting = true
max_requests_per_minute = 30
session_timeout = 7200

[performance]
chunk_cache_size = 200
preload_radius = 3
compression_level = 6
```

### Testing Configuration
```toml
[server]
host = "localhost"
port = "0"  # Random port for testing
shutdown_timeout = 1

[database]
type = "boltdb"
path = ":memory:"  # In-memory database

[worldgen]
dimensions = "5x5x2"
seed = 1

[logging]
level = "debug"
output = "console"

[performance]
chunk_cache_size = 10
preload_radius = 1
```

## 🔧 Configuration Validation

The server validates configuration on startup and reports errors:

### Common Validation Errors

#### **Invalid Dimensions**
```
Error: invalid world dimensions format: "50x50" (expected "WxHxD")
```
**Fix**: Use format like `"50x50x3"`

#### **Invalid Port**
```
Error: invalid port number: 99999 (must be 1-65535)
```
**Fix**: Use valid port number (1-65535)

#### **Invalid File Path**
```
Error: cannot access database path: /invalid/path/eden.db
```
**Fix**: Ensure directory exists and is writable

#### **Invalid Log Level**
```
Error: invalid log level: "verbose" (must be debug, info, warn, error, fatal)
```
**Fix**: Use valid log level names

## 🛠️ Configuration Management

### Environment Variables
Override config values with environment variables:

```bash
export EDEN_SERVER_HOST="0.0.0.0"
export EDEN_SERVER_PORT="9090"
export EDEN_DATABASE_PATH="./custom.db"
./EscapingEden
```

### Configuration Profiles
Use different configs for different environments:

```bash
# Development
./EscapingEden -c configs/development.conf

# Staging
./EscapingEden -c configs/staging.conf

# Production
./EscapingEden -c configs/production.conf
```

### Dynamic Configuration
Some settings can be changed without restart:
- Log level (via admin commands)
- Rate limiting settings
- Cache sizes
- Performance parameters

## 📊 Performance Tuning

### Memory Usage
- **chunk_cache_size**: Higher = more memory, faster access
- **preload_radius**: Higher = more memory, smoother gameplay
- **buffer_size**: Larger = more memory per connection

### CPU Usage
- **compression_level**: Higher = more CPU, less bandwidth
- **gc_interval**: Lower = more CPU, less memory spikes
- **max_goroutines**: Balance concurrency vs overhead

### Disk Usage
- **backup_interval**: More frequent = more disk space
- **log rotation**: Configure to prevent disk filling
- **world dimensions**: Larger = more disk space

### Network Performance
- **buffer_size**: Optimize for connection speed
- **compression_level**: Balance CPU vs bandwidth
- **rate_limiting**: Prevent network congestion

## 🔐 Security Best Practices

### Production Security
1. **Change default ports** from 8080
2. **Enable rate limiting** to prevent abuse
3. **Use IP whitelisting** for admin access
4. **Set strong password requirements**
5. **Enable session timeouts**
6. **Regular security updates**

### Network Security
1. **Use reverse proxy** (nginx, Apache)
2. **Enable HTTPS** with SSL certificates
3. **Configure firewall** rules
4. **Monitor access logs**
5. **Use VPN** for remote admin access

### Data Security
1. **Regular database backups**
2. **Encrypt sensitive data**
3. **Secure file permissions**
4. **Monitor disk usage**
5. **Implement backup rotation**

This configuration system provides flexible control over all aspects of server operation while maintaining security and performance.