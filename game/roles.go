package game

import (
	"fmt"
	"strings"
	"time"

	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

// Role represents a user role in the system
type Role int

const (
	RoleUser Role = iota
	RoleModerator
	RoleAdmin
	RoleSuperAdmin
)

// String returns the string representation of a role
func (r Role) String() string {
	switch r {
	case RoleUser:
		return "User"
	case RoleModerator:
		return "Moderator"
	case RoleAdmin:
		return "Admin"
	case RoleSuperAdmin:
		return "SuperAdmin"
	default:
		return "Unknown"
	}
}

// RoleFromString converts a string to a Role
func RoleFromString(s string) Role {
	switch strings.ToLower(s) {
	case "user":
		return RoleUser
	case "moderator", "mod":
		return RoleModerator
	case "admin":
		return RoleAdmin
	case "superadmin", "super":
		return RoleSuperAdmin
	default:
		return RoleUser
	}
}

// HasPermission checks if a role has permission to perform an action
func (r Role) HasPermission(requiredRole Role) bool {
	return r >= requiredRole
}

// GetPermissions returns a list of permissions for a role
func (r Role) GetPermissions() []string {
	permissions := []string{}

	switch r {
	case RoleSuperAdmin:
		permissions = append(permissions, "All Admin permissions")
		permissions = append(permissions, "Manage admin roles")
		permissions = append(permissions, "System administration")
		fallthrough
	case RoleAdmin:
		permissions = append(permissions, "All Moderator permissions")
		permissions = append(permissions, "Teleport players")
		permissions = append(permissions, "Reset player positions")
		permissions = append(permissions, "Access emergency commands")
		permissions = append(permissions, "View player locations")
		permissions = append(permissions, "Server management")
		fallthrough
	case RoleModerator:
		permissions = append(permissions, "All User permissions")
		permissions = append(permissions, "Moderate chat")
		permissions = append(permissions, "Kick players")
		permissions = append(permissions, "Temporary bans")
		fallthrough
	case RoleUser:
		permissions = append(permissions, "Play game")
		permissions = append(permissions, "Chat")
		permissions = append(permissions, "Basic commands")
	}

	return permissions
}

// RoleManager handles role-related operations
type RoleManager struct {
	gm  *GameManager
	log logging.LoggerType
}

// NewRoleManager creates a new role manager
func NewRoleManager(gm *GameManager) *RoleManager {
	return &RoleManager{
		gm:  gm,
		log: gm.Log,
	}
}

// GetUserRole gets the role of a user by their account ID
func (rm *RoleManager) GetUserRole(accountID string) Role {
	account := messages.Account{}
	err := rm.gm.DB.One("Accounts", "ID", accountID, &account)
	if err != nil {
		rm.log.Println(logging.LogWarn, fmt.Sprintf("Failed to get account %s: %v", accountID, err))
		return RoleUser
	}

	return Role(account.Role)
}

// GetUserRoleByUsername gets the role of a user by their username
func (rm *RoleManager) GetUserRoleByUsername(username string) Role {
	account := messages.Account{}
	err := rm.gm.DB.One("Accounts", "Username", username, &account)
	if err != nil {
		rm.log.Println(logging.LogWarn, fmt.Sprintf("Failed to get account by username %s: %v", username, err))
		return RoleUser
	}

	return Role(account.Role)
}

// GetCharacterRole gets the role of a character by their character ID
func (rm *RoleManager) GetCharacterRole(characterID string) Role {
	character, err := rm.gm.GetCharacter(characterID)
	if err != nil {
		rm.log.Println(logging.LogWarn, fmt.Sprintf("Failed to get character %s: %v", characterID, err))
		return RoleUser
	}

	return rm.GetUserRole(character.UserID)
}

// SetUserRole sets the role of a user
func (rm *RoleManager) SetUserRole(accountID string, role Role) error {
	account := messages.Account{}
	err := rm.gm.DB.One("Accounts", "ID", accountID, &account)
	if err != nil {
		return fmt.Errorf("failed to get account %s: %w", accountID, err)
	}

	oldRole := Role(account.Role)
	account.Role = int(role)

	err = rm.gm.DB.UpdateRecord("Accounts", &account)
	if err != nil {
		return fmt.Errorf("failed to update account role: %w", err)
	}

	rm.log.Println(logging.LogInfo, fmt.Sprintf("Changed role for account %s (%s) from %s to %s",
		accountID, account.Username, oldRole.String(), role.String()))

	return nil
}

// SetUserRoleByUsername sets the role of a user by username
func (rm *RoleManager) SetUserRoleByUsername(username string, role Role) error {
	account := messages.Account{}
	err := rm.gm.DB.One("Accounts", "Username", username, &account)
	if err != nil {
		return fmt.Errorf("failed to get account by username %s: %w", username, err)
	}

	return rm.SetUserRole(account.ID, role)
}

// HasPermission checks if a character has permission for a specific role level
func (rm *RoleManager) HasPermission(characterID string, requiredRole Role) bool {
	userRole := rm.GetCharacterRole(characterID)
	return userRole.HasPermission(requiredRole)
}

// HasPermissionByUsername checks if a user has permission for a specific role level
func (rm *RoleManager) HasPermissionByUsername(username string, requiredRole Role) bool {
	userRole := rm.GetUserRoleByUsername(username)
	return userRole.HasPermission(requiredRole)
}

// EnsureFirstUserIsSuperAdmin ensures the first registered user becomes a super admin
func (rm *RoleManager) EnsureFirstUserIsSuperAdmin() error {
	// Count total accounts
	var accounts []messages.Account
	err := rm.gm.DB.All("Accounts", &accounts)
	if err != nil {
		return fmt.Errorf("failed to count accounts: %w", err)
	}

	if len(accounts) == 0 {
		rm.log.Println(logging.LogInfo, "No accounts found, first user will become SuperAdmin")
		return nil
	}

	// Check if we already have a super admin
	hasSuperAdmin := false
	for _, account := range accounts {
		if Role(account.Role) == RoleSuperAdmin {
			hasSuperAdmin = true
			break
		}
	}

	if !hasSuperAdmin {
		// Find the oldest account (first registered)
		var oldestAccount *messages.Account
		var oldestTime time.Time

		for i, account := range accounts {
			// Use account creation time or last login time as proxy for registration time
			accountTime := account.LastLoginTime
			if accountTime.IsZero() {
				// If no login time, this might be a very old account, prioritize it
				accountTime = time.Unix(0, 0)
			}

			if oldestAccount == nil || accountTime.Before(oldestTime) {
				oldestAccount = &accounts[i]
				oldestTime = accountTime
			}
		}

		if oldestAccount != nil {
			err = rm.SetUserRole(oldestAccount.ID, RoleSuperAdmin)
			if err != nil {
				return fmt.Errorf("failed to set first user as super admin: %w", err)
			}

			rm.log.Println(logging.LogInfo, fmt.Sprintf("Set first user %s as SuperAdmin", oldestAccount.Username))
		}
	}

	return nil
}

// ListUsersByRole returns all users with a specific role
func (rm *RoleManager) ListUsersByRole(role Role) ([]messages.Account, error) {
	var allAccounts []messages.Account
	err := rm.gm.DB.All("Accounts", &allAccounts)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	var filteredAccounts []messages.Account
	for _, account := range allAccounts {
		if Role(account.Role) == role {
			filteredAccounts = append(filteredAccounts, account)
		}
	}

	return filteredAccounts, nil
}

// GetRoleStats returns statistics about role distribution
func (rm *RoleManager) GetRoleStats() map[Role]int {
	stats := make(map[Role]int)

	var accounts []messages.Account
	err := rm.gm.DB.All("Accounts", &accounts)
	if err != nil {
		rm.log.Println(logging.LogWarn, fmt.Sprintf("Failed to get accounts for stats: %v", err))
		return stats
	}

	for _, account := range accounts {
		role := Role(account.Role)
		stats[role]++
	}

	return stats
}

// ValidateRoleHierarchy ensures role changes follow hierarchy rules
func (rm *RoleManager) ValidateRoleHierarchy(adminRole Role, targetRole Role) error {
	// Super admins can set any role
	if adminRole == RoleSuperAdmin {
		return nil
	}

	// Admins can set moderator and user roles, but not admin or super admin
	if adminRole == RoleAdmin {
		if targetRole == RoleAdmin || targetRole == RoleSuperAdmin {
			return fmt.Errorf("admins cannot promote users to admin or super admin roles")
		}
		return nil
	}

	// Moderators and users cannot change roles
	return fmt.Errorf("insufficient permissions to change roles")
}
