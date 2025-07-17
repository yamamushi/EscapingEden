package examples

import (
	"fmt"
	"time"

	"github.com/yamamushi/EscapingEden/logging"
	xterm_256color "github.com/yamamushi/EscapingEden/terminals/xterm-256color"
	"github.com/yamamushi/EscapingEden/ui"
	"github.com/yamamushi/EscapingEden/ui/renderer"
)

// OptimizedRenderingExample demonstrates the new efficient rendering system
func OptimizedRenderingExample() {
	// Create a test logger
	logger := &TestLogger{}

	// Create terminal interface
	terminal := xterm_256color.NewXterm256ColorTerminal()

	// Create optimized console
	console := ui.NewOptimizedConsole(120, 40, terminal, logger, "example-connection")
	console.Initialize()

	// Create game windows with optimal layout
	setupGameWindows(console)

	// Demonstrate efficient updates
	demonstrateEfficientUpdates(console)

	// Show performance statistics
	showPerformanceStats(console)

	// Demonstrate migration from old system
	demonstrateMigration(logger, terminal)
}

// setupGameWindows creates a typical game UI layout
func setupGameWindows(console *ui.OptimizedConsole) {
	// Main game view (left side, most of screen)
	gameWindow := console.CreateWindow("game", 0, 0, 80, 35)
	gameWindow.SetBorder(true, "Game World", packStyle(255, 0, 1)) // White text, bold
	gameWindow.SetTextStyle(packStyle(255, 0, 0))                  // Normal white text

	// Chat window (bottom left)
	chatWindow := console.CreateWindow("chat", 0, 35, 80, 5)
	chatWindow.SetBorder(true, "Chat", packStyle(255, 0, 1))
	chatWindow.SetTextStyle(packStyle(255, 0, 0))

	// Status window (right side)
	statusWindow := console.CreateWindow("status", 80, 0, 40, 20)
	statusWindow.SetBorder(true, "Status", packStyle(255, 0, 1))
	statusWindow.SetTextStyle(packStyle(255, 0, 0))

	// Inventory window (right side, bottom)
	inventoryWindow := console.CreateWindow("inventory", 80, 20, 40, 15)
	inventoryWindow.SetBorder(true, "Inventory", packStyle(255, 0, 1))
	inventoryWindow.SetTextStyle(packStyle(255, 0, 0))

	// Minimap window (right side, very bottom)
	minimapWindow := console.CreateWindow("minimap", 80, 35, 40, 5)
	minimapWindow.SetBorder(true, "Map", packStyle(255, 0, 1))
	minimapWindow.SetTextStyle(packStyle(255, 0, 0))

	// Set initial content
	setInitialContent(gameWindow, chatWindow, statusWindow, inventoryWindow, minimapWindow)
}

// setInitialContent populates windows with example content
func setInitialContent(game, chat, status, inventory, minimap *renderer.OptimizedWindow) {
	// Game world view
	gameContent := []string{
		"You are standing in a vast meadow.",
		"Tall grass sways gently in the breeze.",
		"To the north, you see a dark forest.",
		"To the east, a mountain range looms.",
		"To the south, a river flows peacefully.",
		"To the west, rolling hills extend to the horizon.",
		"",
		"A path leads north into the forest.",
		"",
		"> _",
	}
	game.SetContent(gameContent)

	// Chat history
	chat.AppendLine("Welcome to Escaping Eden!")
	chat.AppendLine("Type 'help' for commands.")
	chat.AppendLine("Player joined the game.")

	// Status information
	statusContent := []string{
		"Health: 100/100",
		"Mana:   50/50",
		"Level:  1",
		"XP:     0/100",
		"",
		"Strength:     10",
		"Dexterity:    12",
		"Intelligence: 14",
		"Constitution: 11",
		"",
		"Location:",
		"Peaceful Meadow",
		"",
		"Weather: Sunny",
		"Time: Midday",
	}
	status.SetContent(statusContent)

	// Inventory
	inventoryContent := []string{
		"Worn leather armor",
		"Rusty iron sword",
		"Small health potion x3",
		"Bread loaf",
		"Water flask",
		"",
		"Gold: 25 coins",
		"",
		"Weight: 15/50 lbs",
	}
	inventory.SetContent(inventoryContent)

	// Minimap
	minimapContent := []string{
		"  ^^^  ",
		" ^^F^^ ",
		"~~~P~~~",
		" ~~R~~ ",
	}
	minimap.SetContent(minimapContent)
}

// demonstrateEfficientUpdates shows how the system efficiently handles updates
func demonstrateEfficientUpdates(console *ui.OptimizedConsole) {
	fmt.Println("=== Demonstrating Efficient Updates ===")

	// Get initial render size
	initialOutput := console.Render()
	fmt.Printf("Initial render size: %d bytes\n", len(initialOutput))

	// Simulate small updates that would be expensive in the old system

	// 1. Update health in status window
	statusWindow := console.GetWindow("status")
	statusContent := []string{
		"Health: 95/100", // Only this line changed
		"Mana:   50/50",
		"Level:  1",
		"XP:     0/100",
		"",
		"Strength:     10",
		"Dexterity:    12",
		"Intelligence: 14",
		"Constitution: 11",
		"",
		"Location:",
		"Peaceful Meadow",
		"",
		"Weather: Sunny",
		"Time: Midday",
	}
	statusWindow.SetContent(statusContent)

	healthUpdateOutput := console.Render()
	fmt.Printf("Health update size: %d bytes (%.1f%% of initial)\n",
		len(healthUpdateOutput),
		float64(len(healthUpdateOutput))/float64(len(initialOutput))*100)

	// 2. Add chat message
	chatWindow := console.GetWindow("chat")
	chatWindow.AppendLine("You feel slightly wounded.")

	chatUpdateOutput := console.Render()
	fmt.Printf("Chat update size: %d bytes (%.1f%% of initial)\n",
		len(chatUpdateOutput),
		float64(len(chatUpdateOutput))/float64(len(initialOutput))*100)

	// 3. Update single character in game view (cursor movement)
	gameWindow := console.GetWindow("game")
	gameContent := []string{
		"You are standing in a vast meadow.",
		"Tall grass sways gently in the breeze.",
		"To the north, you see a dark forest.",
		"To the east, a mountain range looms.",
		"To the south, a river flows peacefully.",
		"To the west, rolling hills extend to the horizon.",
		"",
		"A path leads north into the forest.",
		"",
		"> n_", // Cursor moved, user typed 'n'
	}
	gameWindow.SetContent(gameContent)

	cursorUpdateOutput := console.Render()
	fmt.Printf("Cursor update size: %d bytes (%.1f%% of initial)\n",
		len(cursorUpdateOutput),
		float64(len(cursorUpdateOutput))/float64(len(initialOutput))*100)

	// 4. No changes - should produce no output
	noChangeOutput := console.Render()
	fmt.Printf("No changes size: %d bytes\n", len(noChangeOutput))
}

// showPerformanceStats displays comprehensive performance information
func showPerformanceStats(console *ui.OptimizedConsole) {
	fmt.Println("\n=== Performance Statistics ===")

	stats := console.GetStats()

	fmt.Printf("Console: %dx%d (%d total cells)\n",
		stats.Width, stats.Height, stats.Width*stats.Height)
	fmt.Printf("Total renders: %d\n", stats.RenderingStats.TotalRenders)
	fmt.Printf("Total bytes sent: %d\n", stats.RenderingStats.BytesSent)
	fmt.Printf("Average render time: %v\n", stats.RenderingStats.AverageRenderTime)
	fmt.Printf("Last render size: %d bytes\n", stats.RenderingStats.LastRenderSize)
	fmt.Printf("Efficiency ratio: %.3f\n", stats.RenderingStats.EfficiencyRatio)

	fmt.Printf("\nWindow Statistics:\n")
	fmt.Printf("  Total windows: %d\n", stats.WindowStats.TotalWindows)
	fmt.Printf("  Visible windows: %d\n", stats.WindowStats.VisibleWindows)
	fmt.Printf("  Windows with changes: %d\n", stats.WindowStats.WindowsWithChanges)

	fmt.Printf("\nRenderer Statistics:\n")
	fmt.Printf("  Dirty regions: %d\n", stats.WindowStats.RendererStats.DirtyRegions)
	fmt.Printf("  Dirty cells: %d\n", stats.WindowStats.RendererStats.DirtyCells)
	fmt.Printf("  Total cells: %d\n", stats.WindowStats.RendererStats.TotalCells)

	// Get efficiency report
	efficiency := console.GetEfficiencyReport()
	fmt.Printf("\nEfficiency Report:\n")
	fmt.Printf("  Average bytes per render: %.1f\n", efficiency.AverageBytesPerRender)
	fmt.Printf("  Max bytes per full screen: %.1f\n", efficiency.MaxBytesPerFullScreen)
	fmt.Printf("  Data reduction: %.1f%%\n", efficiency.DataReduction*100)
	fmt.Printf("  Render frequency: %.1f Hz\n", efficiency.RenderFrequency)
}

// demonstrateMigration shows how to migrate from old to new system
func demonstrateMigration(logger logging.LoggerType, terminal *xterm_256color.Xterm256ColorTerminal) {
	fmt.Println("\n=== Migration Demonstration ===")

	// Create old console (simulated)
	oldConsole := ui.NewConsole(120, 40, terminal, logger)
	oldConsole.Initialize()

	// Create migrator
	migrator := ui.NewConsoleMigrator(logger)

	// Migrate to new system
	newConsole := migrator.MigrateConsole(oldConsole)

	// Create adapter for backward compatibility
	adapter := ui.NewConsoleAdapter(newConsole)

	// Use adapter with old-style API calls
	adapter.PrintToChat("Migration completed successfully!")
	adapter.UpdateStatus([]string{"Status: Migrated", "Performance: Optimized"})

	// Render and show size
	output := adapter.Draw()
	fmt.Printf("Migrated console render size: %d bytes\n", len(output))

	// Show performance comparison
	benchmark := ui.NewPerformanceBenchmark(logger)
	results := benchmark.BenchmarkComparison(120, 40, terminal)

	fmt.Printf("\nPerformance Comparison (120x40 terminal):\n")
	fmt.Printf("Full screen update: %d -> %d bytes (%.1f%% reduction)\n",
		results.FullScreenUpdate.OldSystemBytes,
		results.FullScreenUpdate.NewSystemBytes,
		results.FullScreenUpdate.DataReduction*100)

	fmt.Printf("Single line update: %d -> %d bytes (%.1f%% reduction)\n",
		results.SingleLineUpdate.OldSystemBytes,
		results.SingleLineUpdate.NewSystemBytes,
		results.SingleLineUpdate.DataReduction*100)

	fmt.Printf("Chat message append: %d -> %d bytes (%.1f%% reduction)\n",
		results.ChatMessageAppend.OldSystemBytes,
		results.ChatMessageAppend.NewSystemBytes,
		results.ChatMessageAppend.DataReduction*100)
}

// packStyle creates a packed style value from color components
func packStyle(fg, bg, attrs uint8) uint32 {
	return uint32(fg) | (uint32(bg) << 8) | (uint32(attrs) << 16)
}

// TestLogger implements a simple logger for examples
type TestLogger struct{}

func (tl *TestLogger) GetTypeID() logging.LoggerTypeID {
	return logging.LoggerTypeID_Console
}

func (tl *TestLogger) Println(level logging.LogLevel, message string, v ...interface{}) {
	fmt.Printf("[%s] %s", level.String(), message)
	if len(v) > 0 {
		fmt.Printf(" %v", v)
	}
	fmt.Println()
}

// GameUpdateExample shows how to efficiently update game state
func GameUpdateExample(console *ui.OptimizedConsole) {
	fmt.Println("\n=== Game Update Example ===")

	gameWindow := console.GetWindow("game")
	chatWindow := console.GetWindow("chat")
	statusWindow := console.GetWindow("status")

	// Simulate player movement
	for i := 0; i < 5; i++ {
		// Update game view
		gameContent := []string{
			fmt.Sprintf("You are in room %d.", i+1),
			"The walls are made of ancient stone.",
			"Torches flicker on the walls.",
			"",
			fmt.Sprintf("Exits: %s", getExits(i)),
			"",
			"> _",
		}
		gameWindow.SetContent(gameContent)

		// Add movement message to chat
		chatWindow.AppendLine(fmt.Sprintf("You move to room %d.", i+1))

		// Update status (simulate health/mana changes)
		statusContent := []string{
			fmt.Sprintf("Health: %d/100", 100-i*2),
			fmt.Sprintf("Mana:   %d/50", 50-i),
			"Level:  1",
			"XP:     0/100",
			"",
			"Location:",
			fmt.Sprintf("Room %d", i+1),
		}
		statusWindow.SetContent(statusContent)

		// Render and show efficiency
		output := console.Render()
		fmt.Printf("Room %d update: %d bytes\n", i+1, len(output))

		// Small delay to simulate real gameplay
		time.Sleep(100 * time.Millisecond)
	}
}

// getExits returns example exits for a room
func getExits(roomNum int) string {
	exits := []string{"north", "south", "east", "west"}
	return exits[roomNum%len(exits)]
}

// CombatExample demonstrates efficient combat updates
func CombatExample(console *ui.OptimizedConsole) {
	fmt.Println("\n=== Combat Example ===")

	gameWindow := console.GetWindow("game")
	chatWindow := console.GetWindow("chat")
	statusWindow := console.GetWindow("status")

	playerHealth := 100
	enemyHealth := 80

	for round := 1; round <= 5 && playerHealth > 0 && enemyHealth > 0; round++ {
		// Player attacks
		damage := 15 + round*2
		enemyHealth -= damage
		if enemyHealth < 0 {
			enemyHealth = 0
		}

		chatWindow.AppendLine(fmt.Sprintf("You hit the orc for %d damage!", damage))

		// Enemy attacks back
		if enemyHealth > 0 {
			enemyDamage := 10 + round
			playerHealth -= enemyDamage
			if playerHealth < 0 {
				playerHealth = 0
			}
			chatWindow.AppendLine(fmt.Sprintf("The orc hits you for %d damage!", enemyDamage))
		}

		// Update game view with combat status
		gameContent := []string{
			"=== COMBAT ===",
			"",
			fmt.Sprintf("You:     %d/100 HP", playerHealth),
			fmt.Sprintf("Orc:     %d/80 HP", enemyHealth),
			"",
			"Commands: [a]ttack, [d]efend, [r]un",
			"",
			"> _",
		}
		gameWindow.SetContent(gameContent)

		// Update status
		statusContent := []string{
			fmt.Sprintf("Health: %d/100", playerHealth),
			"Mana:   50/50",
			"Level:  1",
			"",
			"In Combat!",
			fmt.Sprintf("Round: %d", round),
		}
		statusWindow.SetContent(statusContent)

		// Render
		output := console.Render()
		fmt.Printf("Combat round %d: %d bytes\n", round, len(output))

		time.Sleep(200 * time.Millisecond)
	}

	// Combat end
	if playerHealth <= 0 {
		chatWindow.AppendLine("You have been defeated!")
	} else {
		chatWindow.AppendLine("You are victorious!")
	}

	finalOutput := console.Render()
	fmt.Printf("Combat end: %d bytes\n", len(finalOutput))
}
