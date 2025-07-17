package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/yamamushi/EscapingEden/edenconfig"
	"github.com/yamamushi/EscapingEden/logging"
)

// handleWorldReset handles the --reset-world command line argument
func handleWorldReset(conf edenconfig.Config, log logging.LoggerType) error {
	fmt.Println("\n⚠️  WORLD RESET WARNING ⚠️")
	fmt.Println("This will permanently delete the entire world and regenerate it from scratch.")
	fmt.Println("All existing map chunks, player positions, and world data will be lost.")
	fmt.Printf("World dimensions from config: %s\n", conf.WorldGen.Dimensions)
	fmt.Println("\nThis action cannot be undone!")

	// Prompt for confirmation
	fmt.Print("\nAre you sure you want to continue? (Y/n): ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read user input: %w", err)
	}

	// Clean up the response
	response = strings.TrimSpace(strings.ToUpper(response))

	// Check if user confirmed with Y
	if response != "Y" {
		fmt.Println("World reset cancelled.")
		os.Exit(0)
	}

	fmt.Println("\n🗑️  Deleting existing world...")
	log.Println(logging.LogInfo, "World reset initiated by --reset-world flag")

	// Delete the world directory
	worldDir := "./assets/world"
	err = deleteWorldDirectory(worldDir, log)
	if err != nil {
		return fmt.Errorf("failed to delete world directory: %w", err)
	}

	fmt.Println("✅ World deletion complete!")
	fmt.Println("🌍 Server will now start and generate a fresh world based on your config.")
	log.Println(logging.LogInfo, "World reset completed successfully")

	return nil
}

// deleteWorldDirectory safely deletes the world directory and all its contents
func deleteWorldDirectory(worldDir string, log logging.LoggerType) error {
	// Check if world directory exists
	if _, err := os.Stat(worldDir); os.IsNotExist(err) {
		log.Println(logging.LogInfo, "World directory doesn't exist, nothing to delete")
		fmt.Println("ℹ️  World directory doesn't exist, nothing to delete")
		return nil
	}

	// Get directory info for logging
	entries, err := os.ReadDir(worldDir)
	if err != nil {
		return fmt.Errorf("failed to read world directory: %w", err)
	}

	fileCount := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			fileCount++
		}
	}

	log.Println(logging.LogInfo, fmt.Sprintf("Deleting world directory with %d files", fileCount))
	fmt.Printf("🗂️  Found %d files to delete...\n", fileCount)

	// Remove the entire directory
	err = os.RemoveAll(worldDir)
	if err != nil {
		return fmt.Errorf("failed to remove world directory: %w", err)
	}

	// Recreate the empty directory
	err = os.MkdirAll(worldDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to recreate world directory: %w", err)
	}

	log.Println(logging.LogInfo, "World directory deleted and recreated")
	fmt.Println("📁 World directory recreated and ready for new world generation")

	return nil
}
