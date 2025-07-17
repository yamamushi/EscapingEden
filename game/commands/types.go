package commands

import (
	"fmt"
	"time"
)

// CommandType represents the type of command being executed
type CommandType string

const (
	CommandDig      CommandType = "dig"
	CommandMove     CommandType = "move"
	CommandBuild    CommandType = "build"
	CommandTeleport CommandType = "teleport"
	CommandMine     CommandType = "mine"
	CommandCraft    CommandType = "craft"
	CommandAttack   CommandType = "attack"
	CommandUse      CommandType = "use"
	CommandDrop     CommandType = "drop"
	CommandPickup   CommandType = "pickup"
)

// PlayerCommand represents a unified command structure
type PlayerCommand struct {
	Type       CommandType            `json:"type"`
	PlayerID   string                 `json:"player_id"`
	Parameters map[string]interface{} `json:"parameters"`
	Context    *CommandContext        `json:"context,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
}

// CommandContext provides additional context for command execution
type CommandContext struct {
	SessionID  string                 `json:"session_id,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
	ClientInfo map[string]interface{} `json:"client_info,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// CommandResult represents the result of command execution
type CommandResult struct {
	Success     bool                   `json:"success"`
	Message     string                 `json:"message"`
	Data        map[string]interface{} `json:"data,omitempty"`
	SideEffects []SideEffect           `json:"side_effects,omitempty"`
	Error       *CommandError          `json:"error,omitempty"`
}

// SideEffect represents additional effects of command execution
type SideEffect struct {
	Type        string                 `json:"type"`
	Target      string                 `json:"target,omitempty"`
	Description string                 `json:"description"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// CommandHandler interface for command handlers
type CommandHandler interface {
	GetCommandType() CommandType
	GetDescription() string
	GetRequiredPermissions() []Permission
	GetValidationRules() []ValidationRule
	Validate(cmd *PlayerCommand) error
	Execute(cmd *PlayerCommand) (*CommandResult, error)
}

// CommandMiddleware interface for command middleware
type CommandMiddleware interface {
	Process(cmd *PlayerCommand, next func(*PlayerCommand) (*CommandResult, error)) (*CommandResult, error)
}

// PermissionChecker interface for permission checking
type PermissionChecker interface {
	HasPermission(playerID string, permission Permission) bool
	GetPlayerPermissions(playerID string) []Permission
}

// CommandLogger interface for command logging
type CommandLogger interface {
	LogCommand(cmd *PlayerCommand, result *CommandResult, duration time.Duration)
	LogError(cmd *PlayerCommand, err error)
}

// CommandMetrics interface for command metrics
type CommandMetrics interface {
	RecordExecution(cmdType CommandType, duration time.Duration, success bool)
	RecordError(cmdType CommandType, errorCode ErrorCode)
	GetStats() map[CommandType]CommandStats
}

// CommandStats represents command execution statistics
type CommandStats struct {
	TotalExecutions int64               `json:"total_executions"`
	SuccessCount    int64               `json:"success_count"`
	ErrorCount      int64               `json:"error_count"`
	AverageTime     time.Duration       `json:"average_time"`
	ErrorsByCode    map[ErrorCode]int64 `json:"errors_by_code"`
}

// GameStateProviderInterface defines the interface for game state access
type GameStateProviderInterface interface {
	GetCharacter(playerID string) (interface{}, error)
	GetMapChunk(chunkID string) (interface{}, error)
	GetTile(x, y, z int, chunkID string) (interface{}, error)
	UpdateGameState(result *CommandResult) error
}

// CommandError represents detailed error information
type CommandError struct {
	Code    ErrorCode              `json:"code"`
	Message string                 `json:"message"`
	Context map[string]interface{} `json:"context,omitempty"`
	Details string                 `json:"details,omitempty"`
}

// ErrorCode represents specific error types
type ErrorCode string

const (
	// Parameter errors
	ErrorInvalidParameter ErrorCode = "invalid_parameter"
	ErrorMissingParameter ErrorCode = "missing_parameter"
	ErrorParameterType    ErrorCode = "parameter_type_error"

	// Permission errors
	ErrorInsufficientPermission ErrorCode = "insufficient_permission"
	ErrorUnauthorized           ErrorCode = "unauthorized"

	// Resource errors
	ErrorResourceNotFound ErrorCode = "resource_not_found"
	ErrorResourceBusy     ErrorCode = "resource_busy"
	ErrorResourceLocked   ErrorCode = "resource_locked"

	// Game state errors
	ErrorInvalidState     ErrorCode = "invalid_state"
	ErrorInvalidPosition  ErrorCode = "invalid_position"
	ErrorInvalidDirection ErrorCode = "invalid_direction"
	ErrorInvalidTarget    ErrorCode = "invalid_target"

	// Item/Tool errors
	ErrorToolNotFound     ErrorCode = "tool_not_found"
	ErrorToolIncompatible ErrorCode = "tool_incompatible"
	ErrorItemNotFound     ErrorCode = "item_not_found"
	ErrorInventoryFull    ErrorCode = "inventory_full"

	// Map/Tile errors
	ErrorTileBlocked   ErrorCode = "tile_blocked"
	ErrorTileInvalid   ErrorCode = "tile_invalid"
	ErrorMapNotLoaded  ErrorCode = "map_not_loaded"
	ErrorChunkNotFound ErrorCode = "chunk_not_found"

	// System errors
	ErrorSystemError   ErrorCode = "system_error"
	ErrorDatabaseError ErrorCode = "database_error"
	ErrorNetworkError  ErrorCode = "network_error"
	ErrorTimeout       ErrorCode = "timeout"
)

// Error implements the error interface
func (e *CommandError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewCommandError creates a new command error
func NewCommandError(code ErrorCode, message string) *CommandError {
	return &CommandError{
		Code:    code,
		Message: message,
		Context: make(map[string]interface{}),
	}
}

// WithContext adds context to a command error
func (e *CommandError) WithContext(key string, value interface{}) *CommandError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithDetails adds detailed information to a command error
func (e *CommandError) WithDetails(details string) *CommandError {
	e.Details = details
	return e
}

// Permission represents a required permission for command execution
type Permission string

const (
	PermissionMove      Permission = "move"
	PermissionDig       Permission = "dig"
	PermissionBuild     Permission = "build"
	PermissionCraft     Permission = "craft"
	PermissionAttack    Permission = "attack"
	PermissionTeleport  Permission = "teleport"
	PermissionAdmin     Permission = "admin"
	PermissionModerator Permission = "moderator"
)

// ValidationRule represents a validation rule for command parameters
type ValidationRule struct {
	Parameter string
	Required  bool
	Type      string
	Validator func(interface{}) error
}
