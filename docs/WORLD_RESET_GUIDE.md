# World Reset Guide - `--reset-world` Command

This guide explains how to use the `--reset-world` command line argument to safely delete and regenerate the game world.

## Overview

The `--reset-world` flag provides a safe way to completely reset the game world, deleting all existing map chunks and regenerating them based on the current configuration file. This is useful for:

- **Development**: Testing world generation changes
- **Configuration Updates**: Applying new world dimension settings
- **Corruption Recovery**: Starting fresh when world data is corrupted
- **Server Maintenance**: Periodic world resets for gameplay reasons

## Usage

### Basic Command
```bash
./EscapingEden --reset-world
```

### With Custom Config
```bash
./EscapingEden --reset-world -c custom-server.conf
```

## Safety Features

### 1. User Confirmation Required
The system will **always** prompt for confirmation before proceeding:

```
⚠️  WORLD RESET WARNING ⚠️
This will permanently delete the entire world and regenerate it from scratch.
All existing map chunks, player positions, and world data will be lost.
World dimensions from config: 50x50x3
This action cannot be undone!

Are you sure you want to continue? (Y/n): 
```

**Important**: You must type exactly `Y` (uppercase) to proceed. Any other input will cancel the operation.

### 2. Comprehensive Logging
All world reset operations are logged with timestamps:
- Reset initiation
- File deletion progress
- Directory recreation
- Completion status

### 3. Graceful Error Handling
The system handles various error conditions:
- Missing world directory (safe to proceed)
- Permission issues (clear error messages)
- Disk space problems (early detection)
- Configuration file issues (validation before reset)

## What Gets Deleted

### Files and Directories Removed:
- **All map chunk files** (`*.map` files in `./assets/world/`)
- **Chunk registry** (`chunk_registry.json`)
- **World metadata** (any additional world-related files)
- **Temporary files** (cache files, partial generations)

### What's Preserved:
- **Player accounts** (stored in database)
- **Character data** (names, stats, inventories)
- **Server configuration** (unchanged)
- **Logs** (reset operation is logged)
- **Other assets** (non-world related files)

## Step-by-Step Process

### 1. Pre-Reset Validation
```
Reading config file at: server.conf
Validating world generation settings...
```

### 2. User Confirmation
```
⚠️  WORLD RESET WARNING ⚠️
[Warning details displayed]
Are you sure you want to continue? (Y/n): Y
```

### 3. World Deletion
```
🗑️  Deleting existing world...
🗂️  Found 1247 files to delete...
📁 World directory recreated and ready for new world generation
✅ World deletion complete!
```

### 4. Server Startup
```
🌍 Server will now start and generate a fresh world based on your config.
Starting async world generation...
Processed 0/2500 chunks (Generated: 0, Skipped: 0)
```

## Example Usage Scenarios

### Scenario 1: Development Testing
```bash
# Developer wants to test new world generation settings
./EscapingEden --reset-world -c test-config.conf
```

**Output:**
```
⚠️  WORLD RESET WARNING ⚠️
World dimensions from config: 100x100x5
Are you sure you want to continue? (Y/n): Y

🗑️  Deleting existing world...
✅ World deletion complete!
🌍 Server will now start and generate a fresh world based on your config.
```

### Scenario 2: Production Reset
```bash
# Server admin performing scheduled world reset
./EscapingEden --reset-world
```

**Output:**
```
⚠️  WORLD RESET WARNING ⚠️
World dimensions from config: 50x50x3
Are you sure you want to continue? (Y/n): Y

🗑️  Deleting existing world...
🗂️  Found 7500 files to delete...
✅ World deletion complete!
```

### Scenario 3: Cancelled Reset
```bash
./EscapingEden --reset-world
```

**Output:**
```
⚠️  WORLD RESET WARNING ⚠️
Are you sure you want to continue? (Y/n): n
World reset cancelled.
```

## Configuration Impact

### World Dimensions
The new world will be generated using current config settings:
```ini
[WorldGen]
Dimensions = "50x50x3"  # Width x Height x Depth
```

### Generation Settings
All world generation parameters from the config file will be applied:
- Terrain algorithms
- Biome distribution  
- Structure placement
- Resource generation

## Player Impact

### Character Positions
- **Existing characters** will be automatically moved to safe spawn points
- **Position reset system** ensures no players are stuck
- **Fallback mechanisms** handle invalid positions gracefully

### Inventory and Stats
- **Character inventories** are preserved (stored in database)
- **Character stats** remain unchanged
- **Account information** is unaffected

### Login Experience
After world reset, players will experience:
1. **Automatic position correction** to safe spawn points
2. **Notification messages** explaining the reset
3. **Normal gameplay** in the fresh world

## Troubleshooting

### Common Issues

#### Permission Denied
```
Error: failed to delete world directory: permission denied
```
**Solution**: Ensure the server has write permissions to the `./assets/world/` directory.

#### Config File Missing
```
Error: config file missing: server.conf
```
**Solution**: Specify the correct config file path with `-c` flag.

#### Disk Space Issues
```
Error: failed to recreate world directory: no space left on device
```
**Solution**: Free up disk space before attempting world reset.

### Recovery Procedures

#### Partial Reset Failure
If the reset process fails partway through:
1. **Check logs** for specific error messages
2. **Manually delete** remaining files in `./assets/world/`
3. **Restart server** to trigger fresh world generation

#### Config Validation Errors
If world generation fails after reset:
1. **Verify config syntax** in `server.conf`
2. **Check world dimensions** are valid
3. **Review generation parameters** for conflicts

## Best Practices

### Before Reset
1. **Backup important data** (if needed)
2. **Notify players** of scheduled reset
3. **Test config changes** in development first
4. **Ensure adequate disk space** for new world

### During Reset
1. **Monitor the process** for errors
2. **Don't interrupt** the deletion process
3. **Wait for completion** before starting server

### After Reset
1. **Verify world generation** completes successfully
2. **Test player login** and positioning
3. **Monitor server performance** during initial generation
4. **Check logs** for any issues

## Advanced Usage

### Scripted Resets
For automated environments:
```bash
#!/bin/bash
echo "Y" | ./EscapingEden --reset-world -c production.conf
```

**Warning**: Use scripted resets carefully in production environments.

### Development Workflow
```bash
# Development cycle with world resets
./EscapingEden --reset-world -c dev.conf
# Test world generation
# Modify config
./EscapingEden --reset-world -c dev.conf
# Repeat
```

## Security Considerations

### Access Control
- **Limit access** to the `--reset-world` flag
- **Use proper file permissions** on server executables
- **Audit reset operations** through log monitoring

### Data Protection
- **Understand data loss** implications
- **Backup strategies** for critical data
- **Recovery procedures** for accidental resets

## Monitoring and Logging

### Log Messages
World reset operations generate detailed logs:
```
[INFO] World reset initiated by --reset-world flag
[INFO] Deleting world directory with 1247 files
[INFO] World directory deleted and recreated
[INFO] World reset completed successfully
```

### Performance Monitoring
After reset, monitor:
- **World generation progress**
- **Server memory usage**
- **Disk I/O during generation**
- **Player connection stability**

This world reset functionality provides a safe, reliable way to refresh the game world while preserving player data and maintaining server stability.