# Admin Teleport Commands

This document describes the admin teleport commands available in Escaping Eden for helping players who are stuck or experiencing map loading issues.

## Overview

The admin teleport system provides administrators with powerful tools to:
- Teleport players to safe locations when they're stuck
- Move players to emergency safe areas when maps fail to load
- Locate players and get their current position
- List available safe spawn points

## Available Commands

### Basic Teleport Commands

#### `/teleport <player> <location>` or `/tp <player> <location>`
Teleports a player to a specified location.

**Examples:**
```
/teleport john emergency
/tp alice origin
/teleport bob 100,200,0
```

#### `/tphere <player>` or `/summon <player>`
Teleports a player to your current location.

**Examples:**
```
/tphere john
/summon alice
```

#### `/emergency <player>` or `/rescue <player>`
Immediately teleports a player to the emergency safe area. This is useful when a player is completely stuck or experiencing severe map issues.

**Examples:**
```
/emergency john
/rescue alice
```

#### `/reset <player>` or `/resetpos <player>`
Resets a player's position to the default spawn location. This is particularly useful for fixing corrupted character data or when a player has invalid map information.

**Examples:**
```
/reset john
/resetpos alice
```

### Information Commands

#### `/locate <player>` or `/where <player>`
Gets the current location of a player, showing both local and global coordinates.

**Examples:**
```
/locate john
/where alice
```

#### `/spawns` or `/listspawns`
Lists all available safe spawn points that can be used as teleport destinations.

**Example:**
```
/spawns
```

#### `/adminhelp` or `/ahelp`
Shows help information for all admin commands.

**Example:**
```
/adminhelp
```

## Location Formats

When using teleport commands, you can specify locations in several formats:

### Named Locations
- `emergency` - Emergency safe chunk (always available)
- `here` - Your current location (when using /tphere)
- `origin` - Origin spawn point
- `origin center` - Center of origin chunk
- `secondary spawn` - Secondary spawn location

### Coordinate Format
- `x,y,z` - Specific global coordinates (e.g., `100,200,0`)

### Special Keywords
- `emergency` - Teleports to the emergency safe chunk
- `here` - Uses admin's current position (for /tphere command)

## Safety Features

### Automatic Safety Checks
The system includes several safety features:

1. **Position Validation**: Ensures teleport destinations are safe and walkable
2. **Map Validation**: Verifies that destination maps exist and are accessible
3. **Emergency Fallback**: If a safe spawn fails, automatically uses emergency chunk
4. **Player Notifications**: Players are informed when they're teleported and why

### Safe Spawn Points
The system maintains a prioritized list of safe spawn points:

1. **Origin Center** (Priority 1) - Center of the origin chunk
2. **Origin Safe Zone** (Priority 2) - Safe area in origin chunk  
3. **Secondary Spawn** (Priority 3) - Alternative spawn location

### Emergency Safe Chunk
When all else fails, the system provides an emergency safe chunk that:
- Is always available in memory
- Contains a simple safe room with walkable floors
- Has basic features to prevent boredom
- Serves as a last resort for stuck players

## Integration with Map Fallback System

The admin teleport commands work seamlessly with the automatic map fallback system:

### Automatic Triggers
- When a player's map fails to load
- When a player ends up at invalid coordinates
- When a player is on an unwalkable tile
- During character login safety checks

### Manual Intervention
Administrators can manually trigger teleports when:
- Players report being stuck
- Map corruption is detected
- Emergency situations arise
- Testing and debugging

## Usage Examples

### Common Scenarios

#### Player Stuck in Invalid Map
```
Admin: /locate john
System: Player john is at local (50,60,0) global (1050,1060,0) in chunk corrupted-chunk-123
Admin: /emergency john
System: Successfully teleported player john to emergency safe area
```

#### Moving Player to Safe Spawn
```
Admin: /spawns
System: Available safe spawn points:
        1. Origin Center (Priority 1)
        2. Origin Safe Zone (Priority 2)
        3. Secondary Spawn (Priority 3)
Admin: /teleport alice origin
System: Successfully teleported player alice to Origin Center
```

#### Bringing Player to Admin Location
```
Admin: /tphere bob
System: Successfully teleported player bob to your location
```

#### Getting Player Location
```
Admin: /where charlie
System: Player charlie is at local (127,127,0) global (127,127,0) in chunk 0-0-0
```

## Error Handling

The system provides clear error messages for common issues:

- **Player not found**: "Failed to teleport player john: target player not found"
- **Invalid location**: "Invalid location format. Use: emergency, here, spawn name, or x,y,z coordinates"
- **Unsafe destination**: "Destination coordinates are not safe"
- **Map loading failure**: "Failed to get map chunk at 100,200,0: chunk not found"

## Security Considerations

### Admin Verification
Currently, the system logs all admin commands but does not implement strict permission checking. In a production environment, you should:

1. Implement proper admin role verification
2. Add audit logging for all teleport actions
3. Consider rate limiting for admin commands
4. Add confirmation prompts for emergency teleports

### Player Privacy
- All teleport actions are logged
- Players are notified when they're teleported
- The reason for teleportation is communicated to the player

## Technical Implementation

### Command Flow
1. Player types command in chat (starting with `/`)
2. Network manager detects admin command and routes to game manager
3. Game manager parses command using `ParseAdminCommand()`
4. Appropriate handler function executes the command
5. Response is sent back to the admin
6. Player is notified if they were affected

### Integration Points
- **Chat System**: Commands are entered through normal chat
- **Game Manager**: Processes and executes commands
- **Fallback System**: Provides safe locations and emergency handling
- **Character Management**: Updates player positions safely
- **Logging System**: Records all admin actions

## Future Enhancements

Potential improvements to the admin teleport system:

1. **Batch Operations**: Teleport multiple players at once
2. **Teleport History**: Track recent teleports for each player
3. **Custom Safe Zones**: Allow admins to define new safe spawn points
4. **Teleport Requests**: Let players request teleports from admins
5. **Map Repair Tools**: Commands to fix corrupted map chunks
6. **Player Grouping**: Teleport entire groups or parties together

## Troubleshooting

### Common Issues

#### Command Not Recognized
- Ensure command starts with `/`
- Check spelling of command name
- Use `/adminhelp` to see available commands

#### Player Not Found
- Verify player name spelling
- Ensure player is currently online
- Check if player ID is correct

#### Teleport Fails
- Check if destination location exists
- Verify coordinates are valid
- Try using emergency teleport as fallback

#### No Response from System
- Check server logs for errors
- Verify game manager is running
- Ensure fallback system is initialized

For additional support or to report issues with the admin teleport system, check the server logs or contact the development team.