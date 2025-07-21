package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// TestActionRegistry_LoadFromFile tests loading configuration from a valid JSON file
func TestActionRegistry_LoadFromFile(t *testing.T) {
	mockLogger := &MockLogger{}
	registry := NewActionRegistry(mockLogger)

	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test_actions.json")

	validConfig := map[string]interface{}{
		"tickRate":     200,
		"maxQueueSize": 3,
		"actions": map[string]interface{}{
			"move": map[string]interface{}{
				"tickCost":      5,
				"description":   "Moving one tile",
				"canQueue":      true,
				"interruptible": true,
				"category":      "movement",
			},
			"build": map[string]interface{}{
				"tickCost":      15,
				"description":   "Building something",
				"canQueue":      true,
				"interruptible": false,
				"category":      "construction",
			},
		},
	}

	// Write config to file
	data, err := json.MarshalIndent(validConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	err = ioutil.WriteFile(configPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Test loading
	err = registry.LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("Failed to load valid configuration: %v", err)
	}

	// Verify loaded values
	if registry.GetTickRate() != 200 {
		t.Errorf("Expected tick rate 200, got %d", registry.GetTickRate())
	}

	if registry.GetMaxQueueSize() != 3 {
		t.Errorf("Expected max queue size 3, got %d", registry.GetMaxQueueSize())
	}

	// Test action definitions
	moveAction, err := registry.GetActionDefinition("move")
	if err != nil {
		t.Fatalf("Failed to get move action: %v", err)
	}

	if moveAction.TickCost != 5 {
		t.Errorf("Expected move tick cost 5, got %d", moveAction.TickCost)
	}

	if !moveAction.Interruptible {
		t.Error("Expected move action to be interruptible")
	}

	buildAction, err := registry.GetActionDefinition("build")
	if err != nil {
		t.Fatalf("Failed to get build action: %v", err)
	}

	if buildAction.Interruptible {
		t.Error("Expected build action to not be interruptible")
	}
}

// TestActionRegistry_LoadFromFile_InvalidJSON tests loading from malformed JSON
func TestActionRegistry_LoadFromFile_InvalidJSON(t *testing.T) {
	mockLogger := &MockLogger{}
	registry := NewActionRegistry(mockLogger)

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "invalid.json")

	// Write invalid JSON
	invalidJSON := `{
		"tickRate": 200,
		"maxQueueSize": 3,
		"actions": {
			"move": {
				"tickCost": 5,
				"description": "Moving"
				// Missing comma - invalid JSON
			}
		}
	}`

	err := ioutil.WriteFile(configPath, []byte(invalidJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid config file: %v", err)
	}

	// Test loading should fail
	err = registry.LoadFromFile(configPath)
	if err == nil {
		t.Error("Expected error when loading invalid JSON")
	}
}

// TestActionRegistry_LoadFromFile_ValidationErrors tests configuration validation
func TestActionRegistry_LoadFromFile_ValidationErrors(t *testing.T) {
	mockLogger := &MockLogger{}
	registry := NewActionRegistry(mockLogger)

	tempDir := t.TempDir()

	testCases := []struct {
		name   string
		config map[string]interface{}
	}{
		{
			name: "negative tick rate",
			config: map[string]interface{}{
				"tickRate":     -100,
				"maxQueueSize": 3,
				"actions": map[string]interface{}{
					"move": map[string]interface{}{
						"tickCost":    5,
						"description": "Moving",
						"canQueue":    true,
					},
				},
			},
		},
		{
			name: "zero max queue size",
			config: map[string]interface{}{
				"tickRate":     200,
				"maxQueueSize": 0,
				"actions": map[string]interface{}{
					"move": map[string]interface{}{
						"tickCost":    5,
						"description": "Moving",
						"canQueue":    true,
					},
				},
			},
		},
		{
			name: "no actions defined",
			config: map[string]interface{}{
				"tickRate":     200,
				"maxQueueSize": 3,
				"actions":      map[string]interface{}{},
			},
		},
		{
			name: "negative tick cost",
			config: map[string]interface{}{
				"tickRate":     200,
				"maxQueueSize": 3,
				"actions": map[string]interface{}{
					"move": map[string]interface{}{
						"tickCost":    -5,
						"description": "Moving",
						"canQueue":    true,
					},
				},
			},
		},
		{
			name: "empty description",
			config: map[string]interface{}{
				"tickRate":     200,
				"maxQueueSize": 3,
				"actions": map[string]interface{}{
					"move": map[string]interface{}{
						"tickCost":    5,
						"description": "",
						"canQueue":    true,
					},
				},
			},
		},
		{
			name: "invalid category",
			config: map[string]interface{}{
				"tickRate":     200,
				"maxQueueSize": 3,
				"actions": map[string]interface{}{
					"move": map[string]interface{}{
						"tickCost":    5,
						"description": "Moving",
						"canQueue":    true,
						"category":    "invalid_category",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			configPath := filepath.Join(tempDir, fmt.Sprintf("%s.json", tc.name))

			data, err := json.MarshalIndent(tc.config, "", "  ")
			if err != nil {
				t.Fatalf("Failed to marshal test config: %v", err)
			}

			err = ioutil.WriteFile(configPath, data, 0644)
			if err != nil {
				t.Fatalf("Failed to write test config file: %v", err)
			}

			err = registry.LoadFromFile(configPath)
			if err == nil {
				t.Errorf("Expected validation error for %s", tc.name)
			}
		})
	}
}

// TestActionRegistry_GetActionDefinition tests retrieving action definitions
func TestActionRegistry_GetActionDefinition(t *testing.T) {
	mockLogger := &MockLogger{}
	registry := NewActionRegistry(mockLogger)

	// Manually set up some actions for testing
	registry.Actions = map[string]*ActionDefinition{
		"move": {
			TickCost:      5,
			Description:   "Moving one tile",
			CanQueue:      true,
			Interruptible: true,
			Category:      "movement",
		},
		"build": {
			TickCost:      15,
			Description:   "Building something",
			CanQueue:      true,
			Interruptible: false,
			Category:      "construction",
		},
	}

	// Test getting existing action
	moveAction, err := registry.GetActionDefinition("move")
	if err != nil {
		t.Fatalf("Failed to get move action: %v", err)
	}

	if moveAction.TickCost != 5 {
		t.Errorf("Expected tick cost 5, got %d", moveAction.TickCost)
	}

	// Test getting non-existent action
	_, err = registry.GetActionDefinition("nonexistent")
	if err == nil {
		t.Error("Expected error when getting non-existent action")
	}

	// Test that returned definition is a copy (modification doesn't affect original)
	moveAction.TickCost = 999
	originalMove, _ := registry.GetActionDefinition("move")
	if originalMove.TickCost != 5 {
		t.Error("Action definition should be returned as a copy")
	}
}

// TestActionRegistry_ValidateAction tests action validation
func TestActionRegistry_ValidateAction(t *testing.T) {
	mockLogger := &MockLogger{}
	registry := NewActionRegistry(mockLogger)

	registry.Actions = map[string]*ActionDefinition{
		"queueable": {
			TickCost:    5,
			Description: "Can be queued",
			CanQueue:    true,
		},
		"not_queueable": {
			TickCost:    5,
			Description: "Cannot be queued",
			CanQueue:    false,
		},
	}

	// Test valid queueable action
	err := registry.ValidateAction("queueable")
	if err != nil {
		t.Errorf("Expected no error for queueable action, got: %v", err)
	}

	// Test non-queueable action
	err = registry.ValidateAction("not_queueable")
	if err == nil {
		t.Error("Expected error for non-queueable action")
	}

	// Test non-existent action
	err = registry.ValidateAction("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent action")
	}
}

// TestActionRegistry_GetActionsByCategory tests filtering actions by category
func TestActionRegistry_GetActionsByCategory(t *testing.T) {
	mockLogger := &MockLogger{}
	registry := NewActionRegistry(mockLogger)

	registry.Actions = map[string]*ActionDefinition{
		"move": {
			TickCost:    5,
			Description: "Moving",
			Category:    "movement",
		},
		"run": {
			TickCost:    3,
			Description: "Running",
			Category:    "movement",
		},
		"build": {
			TickCost:    15,
			Description: "Building",
			Category:    "construction",
		},
		"uncategorized": {
			TickCost:    10,
			Description: "No category",
			Category:    "",
		},
	}

	// Test getting movement actions
	movementActions := registry.GetActionsByCategory("movement")
	if len(movementActions) != 2 {
		t.Errorf("Expected 2 movement actions, got %d", len(movementActions))
	}

	if _, exists := movementActions["move"]; !exists {
		t.Error("Expected 'move' action in movement category")
	}

	if _, exists := movementActions["run"]; !exists {
		t.Error("Expected 'run' action in movement category")
	}

	// Test getting construction actions
	constructionActions := registry.GetActionsByCategory("construction")
	if len(constructionActions) != 1 {
		t.Errorf("Expected 1 construction action, got %d", len(constructionActions))
	}

	// Test getting non-existent category
	nonExistentActions := registry.GetActionsByCategory("nonexistent")
	if len(nonExistentActions) != 0 {
		t.Errorf("Expected 0 actions for non-existent category, got %d", len(nonExistentActions))
	}

	// Test getting empty category
	emptyActions := registry.GetActionsByCategory("")
	if len(emptyActions) != 1 {
		t.Errorf("Expected 1 action with empty category, got %d", len(emptyActions))
	}
}

// TestActionRegistry_CreateDefaultConfiguration tests creating default config
func TestActionRegistry_CreateDefaultConfiguration(t *testing.T) {
	mockLogger := &MockLogger{}
	registry := NewActionRegistry(mockLogger)

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "default_config.json")

	// Test creating default configuration
	err := registry.CreateDefaultConfiguration(configPath)
	if err != nil {
		t.Fatalf("Failed to create default configuration: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Default configuration file was not created")
	}

	// Test loading the created configuration
	newRegistry := NewActionRegistry(mockLogger)
	err = newRegistry.LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("Failed to load created default configuration: %v", err)
	}

	// Verify default values
	if newRegistry.GetTickRate() != 200 {
		t.Errorf("Expected default tick rate 200, got %d", newRegistry.GetTickRate())
	}

	if newRegistry.GetMaxQueueSize() != 3 {
		t.Errorf("Expected default max queue size 3, got %d", newRegistry.GetMaxQueueSize())
	}

	// Verify default actions exist
	expectedActions := []string{"move", "mine", "build_wall", "dig"}
	for _, actionType := range expectedActions {
		_, err := newRegistry.GetActionDefinition(actionType)
		if err != nil {
			t.Errorf("Expected default action '%s' not found: %v", actionType, err)
		}
	}
}

// TestActionRegistry_ReloadFromFile tests configuration reloading
func TestActionRegistry_ReloadFromFile(t *testing.T) {
	mockLogger := &MockLogger{}
	registry := NewActionRegistry(mockLogger)

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "reload_test.json")

	// Create initial configuration
	initialConfig := map[string]interface{}{
		"tickRate":     200,
		"maxQueueSize": 3,
		"actions": map[string]interface{}{
			"move": map[string]interface{}{
				"tickCost":    5,
				"description": "Moving",
				"canQueue":    true,
			},
		},
	}

	data, _ := json.MarshalIndent(initialConfig, "", "  ")
	ioutil.WriteFile(configPath, data, 0644)

	// Load initial configuration
	err := registry.LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("Failed to load initial configuration: %v", err)
	}

	// Verify initial values
	if registry.GetTickRate() != 200 {
		t.Errorf("Expected initial tick rate 200, got %d", registry.GetTickRate())
	}

	// Modify configuration file
	modifiedConfig := map[string]interface{}{
		"tickRate":     300,
		"maxQueueSize": 5,
		"actions": map[string]interface{}{
			"move": map[string]interface{}{
				"tickCost":    7,
				"description": "Moving faster",
				"canQueue":    true,
			},
			"build": map[string]interface{}{
				"tickCost":    20,
				"description": "Building",
				"canQueue":    true,
			},
		},
	}

	data, _ = json.MarshalIndent(modifiedConfig, "", "  ")
	ioutil.WriteFile(configPath, data, 0644)

	// Reload configuration
	err = registry.ReloadFromFile()
	if err != nil {
		t.Fatalf("Failed to reload configuration: %v", err)
	}

	// Verify updated values
	if registry.GetTickRate() != 300 {
		t.Errorf("Expected reloaded tick rate 300, got %d", registry.GetTickRate())
	}

	if registry.GetMaxQueueSize() != 5 {
		t.Errorf("Expected reloaded max queue size 5, got %d", registry.GetMaxQueueSize())
	}

	// Verify new action exists
	_, err = registry.GetActionDefinition("build")
	if err != nil {
		t.Errorf("Expected new 'build' action after reload: %v", err)
	}

	// Test reload without file path set
	emptyRegistry := NewActionRegistry(mockLogger)
	err = emptyRegistry.ReloadFromFile()
	if err == nil {
		t.Error("Expected error when reloading without file path")
	}
}

// TestActionRegistry_ConcurrentAccess tests thread safety
func TestActionRegistry_ConcurrentAccess(t *testing.T) {
	mockLogger := &MockLogger{}
	registry := NewActionRegistry(mockLogger)

	// Set up initial configuration
	registry.TickRate = 200
	registry.MaxQueueSize = 3
	registry.Actions = map[string]*ActionDefinition{
		"move": {
			TickCost:    5,
			Description: "Moving",
			CanQueue:    true,
		},
	}

	var wg sync.WaitGroup
	numGoroutines := 50

	// Concurrently read configuration
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Read operations
			_ = registry.GetTickRate()
			_ = registry.GetMaxQueueSize()
			_, _ = registry.GetActionDefinition("move")
			_ = registry.GetAllActionTypes()
			_ = registry.GetActionsByCategory("movement")
			_ = registry.ValidateAction("move")
		}(i)
	}

	wg.Wait()

	// Verify data integrity after concurrent access
	if registry.GetTickRate() != 200 {
		t.Error("Data corruption detected after concurrent access")
	}
}

// TestActionRegistry_GetConfigurationSummary tests configuration summary
func TestActionRegistry_GetConfigurationSummary(t *testing.T) {
	mockLogger := &MockLogger{}
	registry := NewActionRegistry(mockLogger)

	registry.TickRate = 200
	registry.MaxQueueSize = 3
	registry.filePath = "/test/path/actions.json"
	registry.Actions = map[string]*ActionDefinition{
		"move":          {Category: "movement"},
		"run":           {Category: "movement"},
		"build":         {Category: "construction"},
		"uncategorized": {Category: ""},
	}

	summary := registry.GetConfigurationSummary()

	if summary["tickRate"] != 200 {
		t.Errorf("Expected tick rate 200 in summary, got %v", summary["tickRate"])
	}

	if summary["actionCount"] != 4 {
		t.Errorf("Expected action count 4 in summary, got %v", summary["actionCount"])
	}

	categories, ok := summary["categories"].(map[string]int)
	if !ok {
		t.Fatal("Expected categories to be map[string]int")
	}

	if categories["movement"] != 2 {
		t.Errorf("Expected 2 movement actions in summary, got %d", categories["movement"])
	}

	if categories["construction"] != 1 {
		t.Errorf("Expected 1 construction action in summary, got %d", categories["construction"])
	}

	if categories["uncategorized"] != 1 {
		t.Errorf("Expected 1 uncategorized action in summary, got %d", categories["uncategorized"])
	}
}

// BenchmarkActionRegistry_GetActionDefinition benchmarks action definition retrieval
func BenchmarkActionRegistry_GetActionDefinition(b *testing.B) {
	mockLogger := &MockLogger{}
	registry := NewActionRegistry(mockLogger)

	// Set up many actions for benchmarking
	registry.Actions = make(map[string]*ActionDefinition)
	for i := 0; i < 1000; i++ {
		actionName := fmt.Sprintf("action_%d", i)
		registry.Actions[actionName] = &ActionDefinition{
			TickCost:    i + 1,
			Description: fmt.Sprintf("Action %d", i),
			CanQueue:    true,
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		actionName := fmt.Sprintf("action_%d", i%1000)
		_, _ = registry.GetActionDefinition(actionName)
	}
}

// BenchmarkActionRegistry_ValidateAction benchmarks action validation
func BenchmarkActionRegistry_ValidateAction(b *testing.B) {
	mockLogger := &MockLogger{}
	registry := NewActionRegistry(mockLogger)

	registry.Actions = map[string]*ActionDefinition{
		"move": {
			TickCost:    5,
			Description: "Moving",
			CanQueue:    true,
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = registry.ValidateAction("move")
	}
}
