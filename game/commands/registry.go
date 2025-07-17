package commands

import (
	"fmt"
	"sync"
	"time"

	"github.com/yamamushi/EscapingEden/logging"
)

// CommandRegistry manages command handlers and execution
type CommandRegistry struct {
	handlers    map[CommandType]CommandHandler
	middleware  []CommandMiddleware
	permissions PermissionChecker
	logger      CommandLogger
	metrics     CommandMetrics
	gameState   GameStateProviderInterface
	log         logging.LoggerType
	mutex       sync.RWMutex
}

// NewCommandRegistry creates a new command registry
func NewCommandRegistry(log logging.LoggerType) *CommandRegistry {
	return &CommandRegistry{
		handlers:   make(map[CommandType]CommandHandler),
		middleware: make([]CommandMiddleware, 0),
		log:        log,
	}
}

// RegisterHandler registers a command handler
func (cr *CommandRegistry) RegisterHandler(handler CommandHandler) error {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	cmdType := handler.GetCommandType()
	if _, exists := cr.handlers[cmdType]; exists {
		return fmt.Errorf("handler for command type %s already registered", cmdType)
	}

	cr.handlers[cmdType] = handler
	cr.log.Println(logging.LogInfo, fmt.Sprintf("Registered command handler: %s", cmdType))
	return nil
}

// UnregisterHandler removes a command handler
func (cr *CommandRegistry) UnregisterHandler(cmdType CommandType) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	delete(cr.handlers, cmdType)
	cr.log.Println(logging.LogInfo, fmt.Sprintf("Unregistered command handler: %s", cmdType))
}

// AddMiddleware adds middleware to the command processing pipeline
func (cr *CommandRegistry) AddMiddleware(middleware CommandMiddleware) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	cr.middleware = append(cr.middleware, middleware)
	cr.log.Println(logging.LogInfo, "Added command middleware")
}

// SetPermissionChecker sets the permission checker
func (cr *CommandRegistry) SetPermissionChecker(checker PermissionChecker) {
	cr.permissions = checker
}

// SetLogger sets the command logger
func (cr *CommandRegistry) SetLogger(logger CommandLogger) {
	cr.logger = logger
}

// SetMetrics sets the metrics collector
func (cr *CommandRegistry) SetMetrics(metrics CommandMetrics) {
	cr.metrics = metrics
}

// SetGameStateProvider sets the game state provider
func (cr *CommandRegistry) SetGameStateProvider(provider GameStateProviderInterface) {
	cr.gameState = provider
}

// ProcessCommand processes a command through the full pipeline
func (cr *CommandRegistry) ProcessCommand(cmd *PlayerCommand) (*CommandResult, error) {
	startTime := time.Now()

	// Set timestamp if not already set
	if cmd.Timestamp.IsZero() {
		cmd.Timestamp = startTime
	}

	// Get handler
	cr.mutex.RLock()
	handler, exists := cr.handlers[cmd.Type]
	cr.mutex.RUnlock()

	if !exists {
		err := NewCommandError(ErrorInvalidParameter, fmt.Sprintf("unknown command type: %s", cmd.Type))
		cr.logError(cmd, err)
		return &CommandResult{
			Success: false,
			Error:   err,
		}, err
	}

	// Check permissions
	if cr.permissions != nil {
		requiredPerms := handler.GetRequiredPermissions()
		for _, perm := range requiredPerms {
			if !cr.permissions.HasPermission(cmd.PlayerID, perm) {
				err := NewCommandError(ErrorInsufficientPermission,
					fmt.Sprintf("insufficient permission: %s", perm))
				cr.logError(cmd, err)
				return &CommandResult{
					Success: false,
					Error:   err,
				}, err
			}
		}
	}

	// Process through middleware chain
	result, err := cr.processWithMiddleware(cmd, handler)

	// Record metrics
	duration := time.Since(startTime)
	if cr.metrics != nil {
		cr.metrics.RecordExecution(cmd.Type, duration, result != nil && result.Success)
		if err != nil {
			if cmdErr, ok := err.(*CommandError); ok {
				cr.metrics.RecordError(cmd.Type, cmdErr.Code)
			}
		}
	}

	// Log command execution
	if cr.logger != nil {
		cr.logger.LogCommand(cmd, result, duration)
		if err != nil {
			cr.logger.LogError(cmd, err)
		}
	}

	return result, err
}

// processWithMiddleware processes command through middleware chain
func (cr *CommandRegistry) processWithMiddleware(cmd *PlayerCommand, handler CommandHandler) (*CommandResult, error) {
	// Create execution chain
	execute := func(cmd *PlayerCommand) (*CommandResult, error) {
		return cr.executeCommand(cmd, handler)
	}

	// Apply middleware in reverse order
	for i := len(cr.middleware) - 1; i >= 0; i-- {
		middleware := cr.middleware[i]
		nextExecute := execute
		execute = func(cmd *PlayerCommand) (*CommandResult, error) {
			return middleware.Process(cmd, nextExecute)
		}
	}

	return execute(cmd)
}

// executeCommand executes the actual command
func (cr *CommandRegistry) executeCommand(cmd *PlayerCommand, handler CommandHandler) (*CommandResult, error) {
	// Validate command
	if err := handler.Validate(cmd); err != nil {
		if cmdErr, ok := err.(*CommandError); ok {
			return &CommandResult{
				Success: false,
				Error:   cmdErr,
			}, err
		}

		// Convert generic error to CommandError
		cmdErr := NewCommandError(ErrorInvalidParameter, err.Error())
		return &CommandResult{
			Success: false,
			Error:   cmdErr,
		}, cmdErr
	}

	// Execute command
	result, err := handler.Execute(cmd)
	if err != nil {
		if cmdErr, ok := err.(*CommandError); ok {
			result = &CommandResult{
				Success: false,
				Error:   cmdErr,
			}
		} else {
			// Convert generic error to CommandError
			cmdErr := NewCommandError(ErrorSystemError, err.Error())
			result = &CommandResult{
				Success: false,
				Error:   cmdErr,
			}
		}
	}

	// Update game state if successful
	if result.Success && cr.gameState != nil {
		if err := cr.gameState.UpdateGameState(result); err != nil {
			cr.log.Println(logging.LogWarn, fmt.Sprintf("Failed to update game state: %v", err))
		}
	}

	return result, err
}

// GetRegisteredCommands returns all registered command types
func (cr *CommandRegistry) GetRegisteredCommands() []CommandType {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	commands := make([]CommandType, 0, len(cr.handlers))
	for cmdType := range cr.handlers {
		commands = append(commands, cmdType)
	}
	return commands
}

// GetHandler returns the handler for a specific command type
func (cr *CommandRegistry) GetHandler(cmdType CommandType) (CommandHandler, bool) {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	handler, exists := cr.handlers[cmdType]
	return handler, exists
}

// GetCommandInfo returns information about a command
func (cr *CommandRegistry) GetCommandInfo(cmdType CommandType) (*CommandInfo, error) {
	handler, exists := cr.GetHandler(cmdType)
	if !exists {
		return nil, fmt.Errorf("command type %s not found", cmdType)
	}

	return &CommandInfo{
		Type:        cmdType,
		Description: handler.GetDescription(),
		Permissions: handler.GetRequiredPermissions(),
		Rules:       handler.GetValidationRules(),
	}, nil
}

// GetStats returns command execution statistics
func (cr *CommandRegistry) GetStats() map[CommandType]CommandStats {
	if cr.metrics == nil {
		return make(map[CommandType]CommandStats)
	}
	return cr.metrics.GetStats()
}

// CommandInfo contains information about a command
type CommandInfo struct {
	Type        CommandType      `json:"type"`
	Description string           `json:"description"`
	Permissions []Permission     `json:"permissions"`
	Rules       []ValidationRule `json:"rules"`
}

// logError logs command errors
func (cr *CommandRegistry) logError(cmd *PlayerCommand, err error) {
	cr.log.Println(logging.LogError, fmt.Sprintf("Command error - Type: %s, Player: %s, Error: %v",
		cmd.Type, cmd.PlayerID, err))
}
