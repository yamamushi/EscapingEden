package commands

import (
	"github.com/yamamushi/EscapingEden/edendb"
	"github.com/yamamushi/EscapingEden/edentypes"
	"github.com/yamamushi/EscapingEden/logging"
)

// DigCommandHandler handles the dig command
type DigCommandHandler struct {
	db  edendb.DatabaseType
	log logging.LoggerType
	// Function callbacks to avoid import cycle
	getCharacterFunc     func(string) (interface{}, error)
	getInventoryFunc     func(string) ([]edentypes.Item, error)
	handleDigRequestFunc func(string, string, int, int) error
}

// NewDigCommandHandler creates a new dig command handler
func NewDigCommandHandler(
	db edendb.DatabaseType,
	log logging.LoggerType,
	getCharacterFunc func(string) (interface{}, error),
	getInventoryFunc func(string) ([]edentypes.Item, error),
	handleDigRequestFunc func(string, string, int, int) error,
) *DigCommandHandler {
	return &DigCommandHandler{
		db:                   db,
		log:                  log,
		getCharacterFunc:     getCharacterFunc,
		getInventoryFunc:     getInventoryFunc,
		handleDigRequestFunc: handleDigRequestFunc,
	}
}

// GetCommandType returns the command type
func (h *DigCommandHandler) GetCommandType() CommandType {
	return CommandDig
}

// GetDescription returns the command description
func (h *DigCommandHandler) GetDescription() string {
	return "Dig in a direction using a tool"
}

// GetRequiredPermissions returns the required permissions
func (h *DigCommandHandler) GetRequiredPermissions() []Permission {
	return []Permission{PermissionDig}
}

// GetValidationRules returns the validation rules
func (h *DigCommandHandler) GetValidationRules() []ValidationRule {
	return []ValidationRule{
		{
			Parameter: "tool_id",
			Required:  true,
			Type:      "string",
			Validator: func(value interface{}) error {
				validator := NewIDValidator(true)
				return validator.Validate(value)
			},
		},
		{
			Parameter: "direction",
			Required:  true,
			Type:      "string",
			Validator: func(value interface{}) error {
				validator := &DirectionValidator{}
				return validator.Validate(value)
			},
		},
	}
}

// Validate validates the dig command
func (h *DigCommandHandler) Validate(cmd *PlayerCommand) error {
	// Validate parameters against rules
	if err := ValidateParameters(cmd.Parameters, h.GetValidationRules()); err != nil {
		return err
	}

	// Get parameters
	toolID := cmd.Parameters["tool_id"].(string)
	direction := cmd.Parameters["direction"].(string)

	// Get character
	_, err := h.getCharacterFunc(cmd.PlayerID)
	if err != nil {
		return NewCommandError(ErrorResourceNotFound, "character not found").WithDetails(err.Error())
	}

	// Check if player has the tool
	item := edentypes.Item{}
	err = h.db.One("Items", "ID", toolID, &item)
	if err != nil {
		return NewCommandError(ErrorToolNotFound, "tool not found").WithDetails(err.Error())
	}

	// Check if tool belongs to player
	hasItem := false
	items, err := h.getInventoryFunc(cmd.PlayerID)
	if err != nil {
		return NewCommandError(ErrorSystemError, "failed to get inventory").WithDetails(err.Error())
	}

	for _, invItem := range items {
		if invItem.ID == toolID {
			hasItem = true
			break
		}
	}

	if !hasItem {
		return NewCommandError(ErrorToolNotFound, "tool not in inventory")
	}

	// Check if tool can dig
	if item.Type != edentypes.ItemTool || !item.Attributes["digging"] {
		return NewCommandError(ErrorToolIncompatible, "this tool cannot be used for digging")
	}

	// Calculate target coordinates
	deltaX, deltaY := h.getDirectionDeltas(direction)
	cmd.Parameters["delta_x"] = deltaX
	cmd.Parameters["delta_y"] = deltaY

	return nil
}

// Execute executes the dig command
func (h *DigCommandHandler) Execute(cmd *PlayerCommand) (*CommandResult, error) {
	// Get parameters
	toolID := cmd.Parameters["tool_id"].(string)
	deltaX := cmd.Parameters["delta_x"].(int)
	deltaY := cmd.Parameters["delta_y"].(int)

	// Execute dig action
	err := h.handleDigRequestFunc(toolID, cmd.PlayerID, deltaX, deltaY)
	if err != nil {
		return nil, NewCommandError(ErrorTileBlocked, "cannot dig there").WithDetails(err.Error())
	}

	// Create result
	result := &CommandResult{
		Success: true,
		Message: "You dig successfully.",
		Data: map[string]interface{}{
			"tool_id":   toolID,
			"delta_x":   deltaX,
			"delta_y":   deltaY,
			"player_id": cmd.PlayerID,
		},
		SideEffects: []SideEffect{
			{
				Type:        "tile_changed",
				Description: "Tile was dug",
				Data: map[string]interface{}{
					"delta_x": deltaX,
					"delta_y": deltaY,
				},
			},
		},
	}

	return result, nil
}

// getDirectionDeltas converts a direction string to x,y deltas
func (h *DigCommandHandler) getDirectionDeltas(direction string) (int, int) {
	switch direction {
	// Cardinal directions
	case "n", "north", "k":
		return 0, -1
	case "s", "south", "j":
		return 0, 1
	case "e", "east", "l":
		return 1, 0
	case "w", "west", "h":
		return -1, 0
	// Diagonal directions
	case "ne", "northeast", "u":
		return 1, -1
	case "nw", "northwest", "y":
		return -1, -1
	case "se", "southeast":
		return 1, 1
	case "sw", "southwest", "b":
		return -1, 1
	default:
		return 0, 0
	}
}
