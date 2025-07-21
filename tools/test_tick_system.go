package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// ActionDefinition represents a single action configuration
type ActionDefinition struct {
	TickCost      int    `json:"tickCost"`
	Description   string `json:"description"`
	CanQueue      bool   `json:"canQueue"`
	Interruptible bool   `json:"interruptible"`
	Category      string `json:"category"`
}

// ActionRegistry represents the complete actions.json configuration
type ActionRegistry struct {
	TickRate     int                          `json:"tickRate"`
	MaxQueueSize int                          `json:"maxQueueSize"`
	Actions      map[string]*ActionDefinition `json:"actions"`
}

// TestResults holds the results of various tests
type TestResults struct {
	JSONValidation   bool
	ConfigValidation bool
	TimingAccuracy   bool
	PerformanceTest  bool
	Errors           []string
	Warnings         []string
	TimingResults    TimingTestResults
}

// TimingTestResults holds timing test specific results
type TimingTestResults struct {
	AverageTickDuration time.Duration
	MaxTickDuration     time.Duration
	MinTickDuration     time.Duration
	TickVariance        time.Duration
	TicksProcessed      int
}

// Colors for output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

func main() {
	fmt.Println(ColorCyan + "=================================" + ColorReset)
	fmt.Println(ColorCyan + "  Game Tick System Test Tool" + ColorReset)
	fmt.Println(ColorCyan + "=================================" + ColorReset)
	fmt.Println()

	// Find the actions.json file
	configPath := findConfigFile()
	if configPath == "" {
		log.Fatal(ColorRed + "Error: Could not find actions.json configuration file" + ColorReset)
	}

	fmt.Printf(ColorBlue+"Using configuration file: %s\n"+ColorReset, configPath)
	fmt.Println()

	// Run all tests
	results := runAllTests(configPath)

	// Display results
	displayResults(results)

	// Exit with appropriate code
	if !results.JSONValidation || !results.ConfigValidation {
		os.Exit(1)
	}

	if len(results.Errors) > 0 {
		os.Exit(1)
	}

	fmt.Println(ColorGreen + "All tests passed! Tick system configuration is valid." + ColorReset)
}

// findConfigFile looks for actions.json in common locations
func findConfigFile() string {
	possiblePaths := []string{
		"config/actions.json",
		"../config/actions.json",
		"actions.json",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// runAllTests executes all validation and performance tests
func runAllTests(configPath string) TestResults {
	results := TestResults{
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
	}

	fmt.Println(ColorYellow + "Running tests..." + ColorReset)
	fmt.Println()

	// Test 1: JSON Validation
	fmt.Print("1. JSON Syntax Validation... ")
	registry, jsonValid := testJSONValidation(configPath)
	results.JSONValidation = jsonValid
	if jsonValid {
		fmt.Println(ColorGreen + "PASS" + ColorReset)
	} else {
		fmt.Println(ColorRed + "FAIL" + ColorReset)
		results.Errors = append(results.Errors, "JSON validation failed")
		return results // Can't continue without valid JSON
	}

	// Test 2: Configuration Validation
	fmt.Print("2. Configuration Validation... ")
	configValid := testConfigValidation(registry, &results)
	results.ConfigValidation = configValid
	if configValid {
		fmt.Println(ColorGreen + "PASS" + ColorReset)
	} else {
		fmt.Println(ColorRed + "FAIL" + ColorReset)
	}

	// Test 3: Timing Accuracy Test
	fmt.Print("3. Tick Timing Accuracy... ")
	timingResults := testTimingAccuracy(registry.TickRate)
	results.TimingResults = timingResults
	results.TimingAccuracy = timingResults.TickVariance < time.Duration(registry.TickRate/10)*time.Millisecond
	if results.TimingAccuracy {
		fmt.Println(ColorGreen + "PASS" + ColorReset)
	} else {
		fmt.Println(ColorYellow + "WARN" + ColorReset)
		results.Warnings = append(results.Warnings, "Tick timing variance is high")
	}

	// Test 4: Performance Test
	fmt.Print("4. Performance Test... ")
	perfResult := testPerformance(registry)
	results.PerformanceTest = perfResult
	if perfResult {
		fmt.Println(ColorGreen + "PASS" + ColorReset)
	} else {
		fmt.Println(ColorYellow + "WARN" + ColorReset)
		results.Warnings = append(results.Warnings, "Performance test showed potential issues")
	}

	return results
}

// testJSONValidation validates the JSON syntax and structure
func testJSONValidation(configPath string) (*ActionRegistry, bool) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf(ColorRed+"Error reading file: %v\n"+ColorReset, err)
		return nil, false
	}

	var registry ActionRegistry
	err = json.Unmarshal(data, &registry)
	if err != nil {
		fmt.Printf(ColorRed+"JSON parsing error: %v\n"+ColorReset, err)
		return nil, false
	}

	return &registry, true
}

// testConfigValidation validates the configuration values
func testConfigValidation(registry *ActionRegistry, results *TestResults) bool {
	valid := true

	// Validate tick rate
	if registry.TickRate <= 0 {
		results.Errors = append(results.Errors, "tickRate must be positive")
		valid = false
	} else if registry.TickRate < 50 {
		results.Warnings = append(results.Warnings, "tickRate below 50ms may cause performance issues")
	} else if registry.TickRate > 1000 {
		results.Warnings = append(results.Warnings, "tickRate above 1000ms may feel unresponsive")
	}

	// Validate max queue size
	if registry.MaxQueueSize <= 0 {
		results.Errors = append(results.Errors, "maxQueueSize must be positive")
		valid = false
	} else if registry.MaxQueueSize > 10 {
		results.Warnings = append(results.Warnings, "maxQueueSize above 10 may cause memory issues")
	}

	// Validate actions
	if len(registry.Actions) == 0 {
		results.Errors = append(results.Errors, "no actions defined")
		valid = false
	}

	for actionName, action := range registry.Actions {
		// Validate tick cost
		if action.TickCost <= 0 {
			results.Errors = append(results.Errors, fmt.Sprintf("action '%s': tickCost must be positive", actionName))
			valid = false
		}

		// Validate description
		if strings.TrimSpace(action.Description) == "" {
			results.Warnings = append(results.Warnings, fmt.Sprintf("action '%s': description is empty", actionName))
		}

		// Validate category
		validCategories := []string{"movement", "resource", "construction", "terrain", "combat", "magic", "social", "crafting"}
		categoryValid := false
		for _, validCat := range validCategories {
			if action.Category == validCat {
				categoryValid = true
				break
			}
		}
		if !categoryValid && action.Category != "" {
			results.Warnings = append(results.Warnings, fmt.Sprintf("action '%s': unknown category '%s'", actionName, action.Category))
		}

		// Check for reasonable tick costs
		if action.TickCost > 100 {
			results.Warnings = append(results.Warnings, fmt.Sprintf("action '%s': very high tick cost (%d), may feel slow", actionName, action.TickCost))
		}
	}

	// Check for required actions
	requiredActions := []string{"move", "dig", "mine", "build_wall"}
	for _, required := range requiredActions {
		if _, exists := registry.Actions[required]; !exists {
			results.Warnings = append(results.Warnings, fmt.Sprintf("recommended action '%s' not found", required))
		}
	}

	return valid
}

// testTimingAccuracy tests the accuracy of tick timing
func testTimingAccuracy(tickRate int) TimingTestResults {
	const testDuration = 2 * time.Second
	numTicks := int(testDuration / (time.Duration(tickRate) * time.Millisecond))

	results := TimingTestResults{
		MinTickDuration: time.Hour, // Start with a very high value
	}

	ticker := time.NewTicker(time.Duration(tickRate) * time.Millisecond)
	defer ticker.Stop()

	startTime := time.Now()
	lastTick := startTime
	var durations []time.Duration
	var mutex sync.Mutex

	// Collect timing data
	go func() {
		for i := 0; i < numTicks; i++ {
			<-ticker.C
			now := time.Now()
			duration := now.Sub(lastTick)

			mutex.Lock()
			durations = append(durations, duration)
			if duration > results.MaxTickDuration {
				results.MaxTickDuration = duration
			}
			if duration < results.MinTickDuration {
				results.MinTickDuration = duration
			}
			mutex.Unlock()

			lastTick = now
		}
	}()

	// Wait for test completion
	time.Sleep(testDuration + 100*time.Millisecond)

	// Calculate results
	mutex.Lock()
	results.TicksProcessed = len(durations)
	if results.TicksProcessed > 0 {
		var total time.Duration
		for _, d := range durations {
			total += d
		}
		results.AverageTickDuration = total / time.Duration(results.TicksProcessed)

		// Calculate variance
		var varianceSum time.Duration
		for _, d := range durations {
			diff := d - results.AverageTickDuration
			if diff < 0 {
				diff = -diff
			}
			varianceSum += diff
		}
		results.TickVariance = varianceSum / time.Duration(results.TicksProcessed)
	}
	mutex.Unlock()

	return results
}

// testPerformance runs a basic performance test
func testPerformance(registry *ActionRegistry) bool {
	// Simulate processing multiple action queues
	const numPlayers = 100
	const actionsPerPlayer = 3

	startTime := time.Now()

	// Simulate action processing
	var wg sync.WaitGroup
	for i := 0; i < numPlayers; i++ {
		wg.Add(1)
		go func(playerID int) {
			defer wg.Done()

			// Simulate processing actions for this player
			for j := 0; j < actionsPerPlayer; j++ {
				// Simulate some work (validation, state checking, etc.)
				time.Sleep(time.Microsecond * 10)
			}
		}(i)
	}

	wg.Wait()
	processingTime := time.Since(startTime)

	// Performance should complete well within one tick
	maxAllowedTime := time.Duration(registry.TickRate/4) * time.Millisecond
	return processingTime < maxAllowedTime
}

// displayResults shows the test results in a formatted way
func displayResults(results TestResults) {
	fmt.Println()
	fmt.Println(ColorCyan + "=================================" + ColorReset)
	fmt.Println(ColorCyan + "           Test Results" + ColorReset)
	fmt.Println(ColorCyan + "=================================" + ColorReset)
	fmt.Println()

	// Show test status
	fmt.Printf("JSON Validation:      %s\n", getStatusString(results.JSONValidation))
	fmt.Printf("Config Validation:    %s\n", getStatusString(results.ConfigValidation))
	fmt.Printf("Timing Accuracy:      %s\n", getStatusString(results.TimingAccuracy))
	fmt.Printf("Performance Test:     %s\n", getStatusString(results.PerformanceTest))
	fmt.Println()

	// Show timing details
	if results.TimingAccuracy || len(results.Warnings) > 0 {
		fmt.Println(ColorBlue + "Timing Test Details:" + ColorReset)
		fmt.Printf("  Ticks Processed:    %d\n", results.TimingResults.TicksProcessed)
		fmt.Printf("  Average Duration:   %v\n", results.TimingResults.AverageTickDuration)
		fmt.Printf("  Min Duration:       %v\n", results.TimingResults.MinTickDuration)
		fmt.Printf("  Max Duration:       %v\n", results.TimingResults.MaxTickDuration)
		fmt.Printf("  Variance:           %v\n", results.TimingResults.TickVariance)
		fmt.Println()
	}

	// Show errors
	if len(results.Errors) > 0 {
		fmt.Println(ColorRed + "Errors:" + ColorReset)
		for _, err := range results.Errors {
			fmt.Printf("  • %s\n", err)
		}
		fmt.Println()
	}

	// Show warnings
	if len(results.Warnings) > 0 {
		fmt.Println(ColorYellow + "Warnings:" + ColorReset)
		for _, warning := range results.Warnings {
			fmt.Printf("  • %s\n", warning)
		}
		fmt.Println()
	}

	// Show recommendations
	if len(results.Warnings) > 0 || !results.TimingAccuracy {
		fmt.Println(ColorPurple + "Recommendations:" + ColorReset)
		if !results.TimingAccuracy {
			fmt.Println("  • Consider adjusting system load or tick rate for better timing accuracy")
		}
		if len(results.Warnings) > 0 {
			fmt.Println("  • Review warnings above and adjust configuration as needed")
		}
		fmt.Println("  • Run this test on the target server hardware for accurate results")
		fmt.Println("  • Monitor tick performance in production using server logs")
		fmt.Println()
	}
}

// getStatusString returns a colored status string
func getStatusString(passed bool) string {
	if passed {
		return ColorGreen + "PASS" + ColorReset
	}
	return ColorRed + "FAIL" + ColorReset
}
