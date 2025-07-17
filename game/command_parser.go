package game

import (
	"fmt"
	"strings"

	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

// ParseAdminCommand parses admin commands from chat input
// Returns true if the command was handled, false if it's not an admin command
func (gm *GameManager) ParseAdminCommand(senderID, input string) bool {
	// Check if it's a command (starts with /)
	if !strings.HasPrefix(input, "/") {
		return false
	}

	// Remove the leading slash and split into parts
	command := strings.TrimPrefix(input, "/")
	parts := strings.Fields(command)

	if len(parts) == 0 {
		return false
	}

	commandName := strings.ToLower(parts[0])

	// Get user role for permission checking
	userRole := gm.RoleManager.GetCharacterRole(senderID)

	// Handle different admin commands with permission checks
	switch commandName {
	case "teleport", "tp":
		if !userRole.HasPermission(RoleAdmin) {
			gm.sendPermissionDeniedMessage(senderID, "teleport commands", RoleAdmin)
			return true
		}
		return gm.handleTeleportCommand(senderID, parts[1:])

	case "spawns", "listspawns":
		if !userRole.HasPermission(RoleModerator) {
			gm.sendPermissionDeniedMessage(senderID, "spawn list", RoleModerator)
			return true
		}
		return gm.handleListSpawnsCommand(senderID)

	case "locate", "where", "position":
		if !userRole.HasPermission(RoleAdmin) {
			gm.sendPermissionDeniedMessage(senderID, "player location commands", RoleAdmin)
			return true
		}
		return gm.handleLocateCommand(senderID, parts[1:])

	case "tphere", "summon":
		if !userRole.HasPermission(RoleAdmin) {
			gm.sendPermissionDeniedMessage(senderID, "teleport here commands", RoleAdmin)
			return true
		}
		return gm.handleTeleportHereCommand(senderID, parts[1:])

	case "emergency", "rescue":
		if !userRole.HasPermission(RoleAdmin) {
			gm.sendPermissionDeniedMessage(senderID, "emergency teleport commands", RoleAdmin)
			return true
		}
		return gm.handleEmergencyTeleportCommand(senderID, parts[1:])

	case "reset", "resetpos", "resetposition":
		if !userRole.HasPermission(RoleAdmin) {
			gm.sendPermissionDeniedMessage(senderID, "reset position commands", RoleAdmin)
			return true
		}
		return gm.handleResetPositionCommand(senderID, parts[1:])

	case "role", "roles":
		if !userRole.HasPermission(RoleSuperAdmin) {
			gm.sendPermissionDeniedMessage(senderID, "role management commands", RoleSuperAdmin)
			return true
		}
		return gm.handleRoleCommand(senderID, parts[1:])

	case "whoami":
		return gm.handleWhoAmICommand(senderID)

	case "adminhelp", "ahelp":
		if !userRole.HasPermission(RoleModerator) {
			gm.sendPermissionDeniedMessage(senderID, "admin help", RoleModerator)
			return true
		}
		return gm.handleAdminHelpCommand(senderID)

	case "repairitems", "repair":
		if !userRole.HasPermission(RoleAdmin) {
			gm.sendPermissionDeniedMessage(senderID, "repair items", RoleAdmin)
			return true
		}
		return gm.handleRepairItemsCommand(senderID)

	default:
		return false // Not an admin command
	}
}

// handleTeleportCommand handles /teleport <player> <location>
func (gm *GameManager) handleTeleportCommand(adminID string, args []string) bool {
	if len(args) < 2 {
		gm.sendAdminMessage(adminID, "Usage: /teleport <player> <location>")
		return true
	}

	targetPlayer := args[0]
	location := strings.Join(args[1:], " ")

	// Execute the teleport command directly
	err := gm.AdminTeleportCommand(adminID, targetPlayer, location)

	// Send response back to admin
	var responseMessage string
	if err != nil {
		responseMessage = fmt.Sprintf("Failed to teleport player %s: %v", targetPlayer, err)
	} else {
		responseMessage = fmt.Sprintf("Successfully teleported player %s to %s", targetPlayer, location)
	}

	gm.sendAdminMessage(adminID, responseMessage)
	return true
}

// handleListSpawnsCommand handles /spawns
func (gm *GameManager) handleListSpawnsCommand(adminID string) bool {
	spawns := gm.ListSafeSpawns()
	spawnList := "Available safe spawn points:\n"
	for i, spawn := range spawns {
		spawnList += fmt.Sprintf("%d. %s\n", i+1, spawn)
	}
	spawnList += "\nYou can also use:\n- 'emergency' for emergency chunk\n- 'here' to teleport to your location\n- 'x,y,z' for specific coordinates"

	gm.sendAdminMessage(adminID, spawnList)
	return true
}

// handleLocateCommand handles /locate <player>
func (gm *GameManager) handleLocateCommand(adminID string, args []string) bool {
	if len(args) < 1 {
		gm.sendAdminMessage(adminID, "Usage: /locate <player>")
		return true
	}

	targetPlayer := args[0]

	locationInfo, err := gm.AdminGetPlayerLocation(targetPlayer)

	var responseMessage string
	if err != nil {
		responseMessage = fmt.Sprintf("Failed to get location for player %s: %v", targetPlayer, err)
	} else {
		responseMessage = locationInfo
	}

	gm.sendAdminMessage(adminID, responseMessage)
	return true
}

// handleTeleportHereCommand handles /tphere <player>
func (gm *GameManager) handleTeleportHereCommand(adminID string, args []string) bool {
	if len(args) < 1 {
		gm.sendAdminMessage(adminID, "Usage: /tphere <player>")
		return true
	}

	targetPlayer := args[0]

	// Execute the teleport command directly with "here" as location
	err := gm.AdminTeleportCommand(adminID, targetPlayer, "here")

	// Send response back to admin
	var responseMessage string
	if err != nil {
		responseMessage = fmt.Sprintf("Failed to teleport player %s to your location: %v", targetPlayer, err)
	} else {
		responseMessage = fmt.Sprintf("Successfully teleported player %s to your location", targetPlayer)
	}

	gm.sendAdminMessage(adminID, responseMessage)
	return true
}

// handleEmergencyTeleportCommand handles /emergency <player>
func (gm *GameManager) handleEmergencyTeleportCommand(adminID string, args []string) bool {
	if len(args) < 1 {
		gm.sendAdminMessage(adminID, "Usage: /emergency <player>")
		return true
	}

	targetPlayer := args[0]

	// Execute the teleport command directly with "emergency" as location
	err := gm.AdminTeleportCommand(adminID, targetPlayer, "emergency")

	// Send response back to admin
	var responseMessage string
	if err != nil {
		responseMessage = fmt.Sprintf("Failed to teleport player %s to emergency area: %v", targetPlayer, err)
	} else {
		responseMessage = fmt.Sprintf("Successfully teleported player %s to emergency safe area", targetPlayer)
	}

	gm.sendAdminMessage(adminID, responseMessage)
	return true
}

// handleResetPositionCommand handles /reset <player>
func (gm *GameManager) handleResetPositionCommand(adminID string, args []string) bool {
	if len(args) < 1 {
		gm.sendAdminMessage(adminID, "Usage: /reset <player>")
		return true
	}

	targetPlayer := args[0]

	// Execute the reset command
	err := gm.ResetCharacterToDefaultSpawn(targetPlayer)

	// Send response back to admin
	var responseMessage string
	if err != nil {
		responseMessage = fmt.Sprintf("Failed to reset player %s position: %v", targetPlayer, err)
	} else {
		responseMessage = fmt.Sprintf("Successfully reset player %s to default spawn position", targetPlayer)
	}

	gm.sendAdminMessage(adminID, responseMessage)
	return true
}

// handleAdminHelpCommand handles /adminhelp
func (gm *GameManager) handleAdminHelpCommand(adminID string) bool {
	userRole := gm.RoleManager.GetCharacterRole(adminID)

	helpText := fmt.Sprintf("Admin Commands (Your role: %s):\n\n", userRole.String())

	// Commands available to all staff (Moderator+)
	if userRole.HasPermission(RoleModerator) {
		helpText += "Moderator Commands:\n"
		helpText += "/spawns                       - List available safe spawn points\n"
		helpText += "/listspawns                   - Same as spawns\n"
		helpText += "/whoami                       - Show your role and permissions\n"
		helpText += "/adminhelp                    - Show this help message\n"
		helpText += "/ahelp                        - Short form of adminhelp\n\n"
	}

	// Commands available to Admins+
	if userRole.HasPermission(RoleAdmin) {
		helpText += "Admin Commands:\n"
		helpText += "/teleport <player> <location> - Teleport player to location\n"
		helpText += "/tp <player> <location>       - Short form of teleport\n"
		helpText += "/tphere <player>              - Teleport player to your location\n"
		helpText += "/summon <player>              - Same as tphere\n"
		helpText += "/emergency <player>           - Teleport player to emergency safe area\n"
		helpText += "/rescue <player>              - Same as emergency\n"
		helpText += "/reset <player>               - Reset player to default spawn position\n"
		helpText += "/resetpos <player>            - Same as reset\n"
		helpText += "/locate <player>              - Get player's current location\n"
		helpText += "/where <player>               - Same as locate\n"
		helpText += "/repairitems                  - Repair item database with colors and symbols\n"
		helpText += "/repair                       - Same as repairitems\n\n"
	}

	// Commands available to SuperAdmins only
	if userRole.HasPermission(RoleSuperAdmin) {
		helpText += "SuperAdmin Commands:\n"
		helpText += "/role list                    - List all users and their roles\n"
		helpText += "/role set <user> <role>       - Set user role (user/moderator/admin/superadmin)\n"
		helpText += "/role stats                   - Show role statistics\n\n"
	}

	if userRole.HasPermission(RoleAdmin) {
		helpText += "Location formats:\n"
		helpText += "- \"emergency\" - Emergency safe chunk\n"
		helpText += "- \"here\" - Your current location\n"
		helpText += "- \"origin\" - Origin spawn point\n"
		helpText += "- \"x,y,z\" - Specific coordinates (e.g., \"100,200,0\")\n"
		helpText += "- Spawn point names (use /spawns to see available ones)\n\n"
	}

	helpText += fmt.Sprintf("Your permissions: %s", strings.Join(userRole.GetPermissions(), ", "))

	gm.sendAdminMessage(adminID, helpText)
	return true
}

// sendPermissionDeniedMessage sends a permission denied message to a user
func (gm *GameManager) sendPermissionDeniedMessage(userID, action string, requiredRole Role) {
	message := fmt.Sprintf("Access denied. You need %s role or higher to use %s. Your current role: %s",
		requiredRole.String(), action, gm.RoleManager.GetCharacterRole(userID).String())
	gm.sendAdminMessage(userID, message)
}

// handleWhoAmICommand handles /whoami
func (gm *GameManager) handleWhoAmICommand(userID string) bool {
	userRole := gm.RoleManager.GetCharacterRole(userID)
	character, err := gm.GetCharacter(userID)
	if err != nil {
		gm.sendAdminMessage(userID, "Failed to get character information")
		return true
	}

	message := fmt.Sprintf("Character: %s\nRole: %s\nPermissions: %s",
		character.Name, userRole.String(), strings.Join(userRole.GetPermissions(), ", "))
	gm.sendAdminMessage(userID, message)
	return true
}

// handleRoleCommand handles /role commands
func (gm *GameManager) handleRoleCommand(adminID string, args []string) bool {
	if len(args) == 0 {
		gm.sendAdminMessage(adminID, "Usage: /role <list|set|stats>\n/role list - List all users and roles\n/role set <username> <role> - Set user role\n/role stats - Show role statistics")
		return true
	}

	subCommand := strings.ToLower(args[0])

	switch subCommand {
	case "list":
		return gm.handleRoleListCommand(adminID)
	case "set":
		if len(args) < 3 {
			gm.sendAdminMessage(adminID, "Usage: /role set <username> <role>\nRoles: user, moderator, admin, superadmin")
			return true
		}
		return gm.handleRoleSetCommand(adminID, args[1], args[2])
	case "stats":
		return gm.handleRoleStatsCommand(adminID)
	default:
		gm.sendAdminMessage(adminID, "Unknown role command. Use: list, set, or stats")
		return true
	}
}

// handleRoleListCommand handles /role list
func (gm *GameManager) handleRoleListCommand(adminID string) bool {
	var accounts []messages.Account
	err := gm.DB.All("Accounts", &accounts)
	if err != nil {
		gm.sendAdminMessage(adminID, fmt.Sprintf("Failed to get accounts: %v", err))
		return true
	}

	if len(accounts) == 0 {
		gm.sendAdminMessage(adminID, "No accounts found")
		return true
	}

	message := "User Roles:\n"
	for _, account := range accounts {
		role := Role(account.Role)
		message += fmt.Sprintf("- %s: %s\n", account.Username, role.String())
	}

	gm.sendAdminMessage(adminID, message)
	return true
}

// handleRoleSetCommand handles /role set <username> <role>
func (gm *GameManager) handleRoleSetCommand(adminID, username, roleStr string) bool {
	// Get admin's role for validation
	adminRole := gm.RoleManager.GetCharacterRole(adminID)
	targetRole := RoleFromString(roleStr)

	// Validate role hierarchy
	err := gm.RoleManager.ValidateRoleHierarchy(adminRole, targetRole)
	if err != nil {
		gm.sendAdminMessage(adminID, fmt.Sprintf("Permission denied: %v", err))
		return true
	}

	// Set the role
	err = gm.RoleManager.SetUserRoleByUsername(username, targetRole)
	if err != nil {
		gm.sendAdminMessage(adminID, fmt.Sprintf("Failed to set role: %v", err))
		return true
	}

	gm.sendAdminMessage(adminID, fmt.Sprintf("Successfully set %s's role to %s", username, targetRole.String()))
	return true
}

// handleRoleStatsCommand handles /role stats
func (gm *GameManager) handleRoleStatsCommand(adminID string) bool {
	stats := gm.RoleManager.GetRoleStats()

	message := "Role Statistics:\n"
	total := 0
	for role, count := range stats {
		message += fmt.Sprintf("- %s: %d\n", role.String(), count)
		total += count
	}
	message += fmt.Sprintf("Total Users: %d", total)

	gm.sendAdminMessage(adminID, message)
	return true
}

// handleRepairItemsCommand handles /repairitems command
func (gm *GameManager) handleRepairItemsCommand(adminID string) bool {
	gm.sendAdminMessage(adminID, "Starting item database repair... This may take a moment.")

	// Run the repair function
	go func() {
		gm.RepairItemDatabase()
		gm.sendAdminMessage(adminID, "Item database repair completed. Check server logs for details.")
	}()

	return true
}

// sendAdminMessage sends a message directly to an admin
func (gm *GameManager) sendAdminMessage(adminID, message string) {
	response := messages.ConnectionManagerMessage{
		Type:               messages.ConnectManager_Message_GameCommandResponse,
		RecipientConsoleID: adminID,
		Data: messages.GameMessage{
			Type: messages.GM_SystemMessage,
			Data: messages.GameMessageData{
				CharacterID: adminID,
				Data:        message,
			},
		},
	}

	gm.SendChannel <- response
	gm.Log.Println(logging.LogInfo, fmt.Sprintf("Sent admin message to %s: %s", adminID, message))
}

// IsAdminCommand checks if a string is an admin command
func IsAdminCommand(input string) bool {
	if !strings.HasPrefix(input, "/") {
		return false
	}

	command := strings.TrimPrefix(input, "/")
	parts := strings.Fields(command)

	if len(parts) == 0 {
		return false
	}

	commandName := strings.ToLower(parts[0])
	adminCommands := []string{
		"teleport", "tp", "tphere", "summon", "emergency", "rescue",
		"locate", "where", "position", "spawns", "listspawns",
		"adminhelp", "ahelp", "reset", "resetpos", "resetposition",
		"role", "roles", "whoami", "repairitems", "repair",
	}

	for _, cmd := range adminCommands {
		if commandName == cmd {
			return true
		}
	}

	return false
}
