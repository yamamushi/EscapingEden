# Role-Based Permission System

This document describes the comprehensive role-based permission system implemented in Escaping Eden.

## Overview

The role system provides hierarchical access control with four distinct roles:
- **User** (Level 0) - Default role for all players
- **Moderator** (Level 1) - Basic moderation capabilities
- **Admin** (Level 2) - Full game administration
- **SuperAdmin** (Level 3) - System administration and role management

## Role Hierarchy

The system uses inheritance-based permissions where higher roles automatically inherit all permissions from lower roles:

```
SuperAdmin (3)
    ├── All Admin permissions
    ├── Role management
    └── System administration
    
Admin (2)
    ├── All Moderator permissions
    ├── Player teleportation
    ├── Position management
    └── Emergency commands
    
Moderator (1)
    ├── All User permissions
    ├── Chat moderation
    ├── Basic admin commands
    └── Player assistance
    
User (0)
    ├── Play game
    ├── Chat
    └── Basic commands
```

## Automatic SuperAdmin Assignment

The system automatically assigns the first registered user as SuperAdmin:

1. **First User**: When the first account is created, they automatically become SuperAdmin
2. **Existing Systems**: If no SuperAdmin exists, the oldest account (by registration/login time) is promoted
3. **Fail-Safe**: Ensures there's always at least one SuperAdmin in the system

## Role Management Commands

### For SuperAdmins Only

#### `/role list`
Lists all users and their current roles.
```
/role list
```
**Output:**
```
User Roles:
- alice: SuperAdmin
- bob: Admin
- charlie: Moderator
- dave: User
```

#### `/role set <username> <role>`
Sets a user's role. SuperAdmins can set any role.
```
/role set bob admin
/role set charlie moderator
/role set dave user
```

#### `/role stats`
Shows statistics about role distribution.
```
/role stats
```
**Output:**
```
Role Statistics:
- SuperAdmin: 1
- Admin: 2
- Moderator: 3
- User: 15
Total Users: 21
```

### For All Staff

#### `/whoami`
Shows your current role and permissions.
```
/whoami
```
**Output:**
```
Character: YourCharacter
Role: Admin
Permissions: Play game, Chat, Basic commands, Moderate chat, Kick players, Temporary bans, Teleport players, Reset player positions, Access emergency commands, View player locations, Server management
```

## Permission-Based Commands

### Moderator Commands (Level 1+)
- `/spawns` - List available safe spawn points
- `/listspawns` - Same as spawns
- `/whoami` - Show role and permissions
- `/adminhelp` - Show available commands

### Admin Commands (Level 2+)
- `/teleport <player> <location>` - Teleport player to location
- `/tp <player> <location>` - Short form of teleport
- `/tphere <player>` - Teleport player to your location
- `/summon <player>` - Same as tphere
- `/emergency <player>` - Teleport to emergency safe area
- `/rescue <player>` - Same as emergency
- `/reset <player>` - Reset player position to default spawn
- `/resetpos <player>` - Same as reset
- `/locate <player>` - Get player's current location
- `/where <player>` - Same as locate

### SuperAdmin Commands (Level 3+)
- `/role list` - List all users and roles
- `/role set <user> <role>` - Set user role
- `/role stats` - Show role statistics

## Permission Validation

### Hierarchy Enforcement
- **SuperAdmins** can promote/demote anyone to any role
- **Admins** can only set Moderator and User roles (cannot create other Admins)
- **Moderators** and **Users** cannot change roles

### Command Access Control
Each command checks permissions before execution:
```go
if !userRole.HasPermission(RequiredRole) {
    // Access denied message
    return
}
```

### Error Messages
Clear feedback when permissions are insufficient:
```
Access denied. You need Admin role or higher to use teleport commands. Your current role: Moderator
```

## Database Integration

### Account Structure
The role is stored in the Account record:
```go
type Account struct {
    // ... other fields
    Role int // 0=User, 1=Moderator, 2=Admin, 3=SuperAdmin
}
```

### Automatic Migration
- Existing accounts without role field default to User (0)
- First user detection and promotion happens on server startup
- Role changes are immediately persisted to database

## Technical Implementation

### Role Manager
Central component handling all role operations:
```go
type RoleManager struct {
    gm  *GameManager
    log logging.LoggerType
}
```

### Key Methods
- `GetCharacterRole(characterID)` - Get role by character
- `GetUserRole(accountID)` - Get role by account
- `SetUserRole(accountID, role)` - Set user role
- `HasPermission(characterID, requiredRole)` - Check permissions
- `EnsureFirstUserIsSuperAdmin()` - Auto-promote first user

### Permission Checking
```go
func (r Role) HasPermission(requiredRole Role) bool {
    return r >= requiredRole
}
```

## Security Features

### Audit Logging
All role changes and admin commands are logged:
```
[INFO] Changed role for account abc123 (alice) from User to Admin
[INFO] Admin alice teleported player bob to emergency safe area
```

### Validation
- Role hierarchy is enforced (Admins can't promote to Admin+)
- Invalid role names are rejected
- Non-existent users are handled gracefully

### Fail-Safe Mechanisms
- Always ensures at least one SuperAdmin exists
- Graceful degradation when role system fails
- Default to User role for any errors

## Usage Examples

### Promoting a User to Moderator
```
SuperAdmin: /role set newmod moderator
System: Successfully set newmod's role to Moderator
```

### Admin Helping Stuck Player
```
Admin: /locate stuckplayer
System: Player stuckplayer is at local (50,60,0) global (1050,1060,0) in chunk corrupted-chunk
Admin: /emergency stuckplayer
System: Successfully teleported player stuckplayer to emergency safe area
```

### Checking Your Permissions
```
User: /whoami
System: Character: MyCharacter
        Role: User
        Permissions: Play game, Chat, Basic commands

Moderator: /adminhelp
System: Admin Commands (Your role: Moderator):
        
        Moderator Commands:
        /spawns - List available safe spawn points
        ...
```

### Permission Denied Example
```
User: /teleport someone somewhere
System: Access denied. You need Admin role or higher to use teleport commands. Your current role: User
```

## Best Practices

### For SuperAdmins
1. **Promote Carefully**: Only promote trusted users to Admin roles
2. **Regular Audits**: Use `/role stats` to monitor role distribution
3. **Backup Admins**: Ensure multiple SuperAdmins exist for redundancy

### For Admins
1. **Use Appropriate Commands**: Use `/emergency` for stuck players, `/reset` for corrupted data
2. **Document Actions**: Communicate with players about teleports
3. **Escalate When Needed**: Contact SuperAdmins for role-related issues

### For Moderators
1. **Know Your Limits**: Understand which commands you can and cannot use
2. **Help Players**: Use `/spawns` to guide players to safe locations
3. **Report Issues**: Escalate serious problems to Admins

## Migration and Compatibility

### Existing Accounts
- All existing accounts default to User role
- First registered user is automatically promoted to SuperAdmin
- No data loss or compatibility issues

### Database Schema
- Role field added to Account structure
- Backward compatible with existing databases
- Automatic migration on first startup

### Command Compatibility
- All existing admin commands now require appropriate roles
- New role management commands added
- Existing functionality preserved

## Troubleshooting

### Common Issues

#### "No SuperAdmin Found"
- System automatically promotes first user
- Check logs for promotion messages
- Manually set role in database if needed

#### "Permission Denied for Valid Admin"
- Verify character is linked to correct account
- Check role assignment with `/whoami`
- Ensure database role field is set correctly

#### "Role Commands Not Working"
- Verify you have SuperAdmin role
- Check command syntax: `/role set username rolename`
- Ensure target user exists

### Debug Commands
```bash
# Check role in database (example)
/role list  # Shows all users and roles
/whoami     # Shows your current role
/role stats # Shows role distribution
```

## Future Enhancements

### Planned Features
1. **Time-Limited Roles**: Temporary promotions with expiration
2. **Custom Permissions**: Fine-grained permission control
3. **Role Templates**: Predefined role configurations
4. **Audit Dashboard**: Web interface for role management
5. **Integration**: Discord role synchronization

### Extensibility
The role system is designed to be easily extended:
- Add new roles by extending the Role enum
- Implement custom permission checks
- Add role-specific features and commands
- Integrate with external authentication systems

This role system provides a solid foundation for managing user permissions while maintaining security and ease of use.