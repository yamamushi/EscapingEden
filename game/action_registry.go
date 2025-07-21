package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/yamamushi/EscapingEden/logging"
)

// ActionDefinition represents a single action configuration
type ActionDefinition struct {
	TickCost      int    `json:"tickCost"`
	Description   string `json:"description"`
	CanQueue      bool   `json:"canQueue"`
	Interruptible bool   `json:"interruptible"`
	Category      string `json:"category"`
}

// ActionRegistry manages the configuration of all game actions
type ActionRegistry struct {
	TickRate     int                          `json:"tickRate"`
	MaxQueueSize int                          `json:"maxQueueSize"`
	Actions      map[string]*ActionDefinition `json:"actions"`
	filePath     string
	mutex        sync.RWMutex
	log          logging.LoggerType
}

// NewActionRegistry creates a new ActionRegistry
func NewActionRegistry(log logging.LoggerType) *ActionRegistry {
	return &ActionRegistry{
		Actions: make(map[string]*ActionDefinition),
		mutex:   sync.RWMutex{},
		log:     log,
	}
}

// LoadFromFile loads action configuration from a JSON file
func (ar *ActionRegistry) LoadFromFile(filePath string) error {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()

	ar.filePath = filePath

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("configuration file not found: %s", filePath)
	}

	// Read file contents
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read configuration file: %v", err)
	}

	// Parse JSON
	var tempRegistry ActionRegistry
	err = json.Unmarshal(data, &tempRegistry)
	if err != nil {
		return fmt.Errorf("failed to parse JSON configuration: %v", err)
	}

	// Validate configuration
	err = ar.validateConfiguration(&tempRegistry)
	if err != nil {
		return fmt.Errorf("configuration validation failed: %v", err)
	}

	// Apply configuration
	ar.TickRate = tempRegistry.TickRate
	ar.MaxQueueSize = tempRegistry.MaxQueueSize
	ar.Actions = tempRegistry.Actions

	ar.log.Println(logging.LogInfo, fmt.Sprintf("Loaded action configuration from %s: %d actions, tick rate %dms",
		filePath, len(ar.Actions), ar.TickRate))

	return nil
}

// GetActionDefinition returns the definition for a specific action type
func (ar *ActionRegistry) GetActionDefinition(actionType string) (*ActionDefinition, error) {
	ar.mutex.RLock()
	defer ar.mutex.RUnlock()

	definition, exists := ar.Actions[actionType]
	if !exists {
		return nil, fmt.Errorf("action type '%s' not found", actionType)
	}

	// Return a copy to prevent external modification
	return &ActionDefinition{
		TickCost:      definition.TickCost,
		Description:   definition.Description,
		CanQueue:      definition.CanQueue,
		Interruptible: definition.Interruptible,
		Category:      definition.Category,
	}, nil
}

// GetAllActionTypes returns a list of all available action types
func (ar *ActionRegistry) GetAllActionTypes() []string {
	ar.mutex.RLock()
	defer ar.mutex.RUnlock()

	types := make([]string, 0, len(ar.Actions))
	for actionType := range ar.Actions {
		types = append(types, actionType)
	}

	return types
}

// GetTickRate returns the configured tick rate
func (ar *ActionRegistry) GetTickRate() int {
	ar.mutex.RLock()
	defer ar.mutex.RUnlock()
	return ar.TickRate
}

// GetMaxQueueSize returns the configured maximum queue size
func (ar *ActionRegistry) GetMaxQueueSize() int {
	ar.mutex.RLock()
	defer ar.mutex.RUnlock()
	return ar.MaxQueueSize
}

// ValidateAction checks if an action type is valid and can be queued
func (ar *ActionRegistry) ValidateAction(actionType string) error {
	ar.mutex.RLock()
	defer ar.mutex.RUnlock()

	definition, exists := ar.Actions[actionType]
	if !exists {
		return fmt.Errorf("unknown action type: %s", actionType)
	}

	if !definition.CanQueue {
		return fmt.Errorf("action type '%s' cannot be queued", actionType)
	}

	return nil
}

// ReloadFromFile reloads the configuration from the file
func (ar *ActionRegistry) ReloadFromFile() error {
	if ar.filePath == "" {
		return errors.New("no file path set for reload")
	}

	ar.log.Println(logging.LogInfo, "Reloading action configuration...")

	err := ar.LoadFromFile(ar.filePath)
	if err != nil {
		ar.log.Println(logging.LogError, fmt.Sprintf("Failed to reload configuration: %v", err))
		return err
	}

	ar.log.Println(logging.LogInfo, "Action configuration reloaded successfully")
	return nil
}

// validateConfiguration validates the loaded configuration
func (ar *ActionRegistry) validateConfiguration(config *ActionRegistry) error {
	// Validate tick rate
	if config.TickRate <= 0 {
		return errors.New("tickRate must be positive")
	}

	if config.TickRate < 10 {
		return errors.New("tickRate too low (minimum 10ms)")
	}

	if config.TickRate > 5000 {
		return errors.New("tickRate too high (maximum 5000ms)")
	}

	// Validate max queue size
	if config.MaxQueueSize <= 0 {
		return errors.New("maxQueueSize must be positive")
	}

	if config.MaxQueueSize > 20 {
		return errors.New("maxQueueSize too high (maximum 20)")
	}

	// Validate actions
	if len(config.Actions) == 0 {
		return errors.New("no actions defined")
	}

	validCategories := map[string]bool{
		"movement":     true,
		"resource":     true,
		"construction": true,
		"terrain":      true,
		"combat":       true,
		"magic":        true,
		"social":       true,
		"crafting":     true,
	}

	for actionType, action := range config.Actions {
		// Validate action type name
		if actionType == "" {
			return errors.New("action type cannot be empty")
		}

		// Validate tick cost
		if action.TickCost <= 0 {
			return fmt.Errorf("action '%s': tickCost must be positive", actionType)
		}

		if action.TickCost > 1000 {
			return fmt.Errorf("action '%s': tickCost too high (maximum 1000)", actionType)
		}

		// Validate description
		if action.Description == "" {
			return fmt.Errorf("action '%s': description cannot be empty", actionType)
		}

		// Validate category (if provided)
		if action.Category != "" && !validCategories[action.Category] {
			return fmt.Errorf("action '%s': invalid category '%s'", actionType, action.Category)
		}
	}

	return nil
}

// GetActionsByCategory returns all actions in a specific category
func (ar *ActionRegistry) GetActionsByCategory(category string) map[string]*ActionDefinition {
	ar.mutex.RLock()
	defer ar.mutex.RUnlock()

	result := make(map[string]*ActionDefinition)
	for actionType, definition := range ar.Actions {
		if definition.Category == category {
			// Return a copy to prevent external modification
			result[actionType] = &ActionDefinition{
				TickCost:      definition.TickCost,
				Description:   definition.Description,
				CanQueue:      definition.CanQueue,
				Interruptible: definition.Interruptible,
				Category:      definition.Category,
			}
		}
	}

	return result
}

// GetConfigurationSummary returns a summary of the current configuration
func (ar *ActionRegistry) GetConfigurationSummary() map[string]interface{} {
	ar.mutex.RLock()
	defer ar.mutex.RUnlock()

	summary := map[string]interface{}{
		"tickRate":     ar.TickRate,
		"maxQueueSize": ar.MaxQueueSize,
		"actionCount":  len(ar.Actions),
		"filePath":     ar.filePath,
	}

	// Count actions by category
	categories := make(map[string]int)
	for _, action := range ar.Actions {
		if action.Category != "" {
			categories[action.Category]++
		} else {
			categories["uncategorized"]++
		}
	}
	summary["categories"] = categories

	return summary
}

// CreateDefaultConfiguration creates a default configuration file
func (ar *ActionRegistry) CreateDefaultConfiguration(filePath string) error {
	defaultConfig := &ActionRegistry{
		TickRate:     200,
		MaxQueueSize: 3,
		Actions: map[string]*ActionDefinition{
			"move": {
				TickCost:      5,
				Description:   "Moving one tile",
				CanQueue:      true,
				Interruptible: true,
				Category:      "movement",
			},
			"mine": {
				TickCost:      20,
				Description:   "Mining a tile",
				CanQueue:      true,
				Interruptible: false,
				Category:      "resource",
			},
			"build_wall": {
				TickCost:      15,
				Description:   "Building a wall",
				CanQueue:      true,
				Interruptible: false,
				Category:      "construction",
			},
			"dig": {
				TickCost:      10,
				Description:   "Digging a tile",
				CanQueue:      true,
				Interruptible: true,
				Category:      "terrain",
			},
		},
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal default configuration: %v", err)
	}

	// Write to file
	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write default configuration: %v", err)
	}

	ar.log.Println(logging.LogInfo, fmt.Sprintf("Created default configuration file: %s", filePath))
	return nil
}

// WatchForChanges starts watching the configuration file for changes (basic implementation)
func (ar *ActionRegistry) WatchForChanges() error {
	if ar.filePath == "" {
		return errors.New("no file path set for watching")
	}

	// Get initial file info
	initialInfo, err := os.Stat(ar.filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	// Start watching in a goroutine
	go func() {
		ticker := time.NewTicker(5 * time.Second) // Check every 5 seconds
		defer ticker.Stop()

		lastModTime := initialInfo.ModTime()

		for range ticker.C {
			currentInfo, err := os.Stat(ar.filePath)
			if err != nil {
				ar.log.Println(logging.LogWarn, fmt.Sprintf("Failed to check file modification time: %v", err))
				continue
			}

			if currentInfo.ModTime().After(lastModTime) {
				ar.log.Println(logging.LogInfo, "Configuration file changed, reloading...")
				err := ar.ReloadFromFile()
				if err != nil {
					ar.log.Println(logging.LogError, fmt.Sprintf("Failed to reload after file change: %v", err))
				}
				lastModTime = currentInfo.ModTime()
			}
		}
	}()

	ar.log.Println(logging.LogInfo, fmt.Sprintf("Started watching configuration file: %s", ar.filePath))
	return nil
}
