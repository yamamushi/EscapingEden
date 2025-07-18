package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yamamushi/EscapingEden/edenconfig"
	"github.com/yamamushi/EscapingEden/edendb"
	"github.com/yamamushi/EscapingEden/edendb/bolt"
	"github.com/yamamushi/EscapingEden/messages"
)

func main() {
	fmt.Println("=== Character Attributes Update Tool ===")
	fmt.Println("This tool will update existing character records to include default attributes.")
	fmt.Println()

	// Check if config file exists
	configPath := "server.conf"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("Config file not found: %s", configPath)
	}

	// Load configuration
	fmt.Printf("Loading configuration from: %s\n", configPath)
	conf, err := edenconfig.ReadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	// Initialize database
	fmt.Println("Connecting to database...")
	var db edendb.DatabaseType
	if strings.ToLower(conf.DB.Type) == "bolt" {
		db, err = bolt.NewBoltDB(conf.DB.Path)
		if err != nil {
			log.Fatalf("Failed to create database connection: %v", err)
		}
	} else {
		log.Fatalf("Unsupported database type: %s", conf.DB.Type)
	}

	// Get all characters
	fmt.Println("Retrieving all character records...")
	var characters []messages.CharacterInfo
	err = db.All("Characters", &characters)
	if err != nil {
		log.Fatalf("Failed to retrieve characters: %v", err)
	}

	fmt.Printf("Found %d character records\n", len(characters))

	if len(characters) == 0 {
		fmt.Println("No characters found. Nothing to update.")
		return
	}

	// Process each character
	updatedCount := 0
	skippedCount := 0

	for i, character := range characters {
		fmt.Printf("\nProcessing character %d/%d: %s (ID: %s)\n",
			i+1, len(characters), character.Name, character.ID)

		// Check if attributes are already set (not empty)
		if !character.Attributes.IsEmpty() {
			fmt.Printf("  ✓ Character already has attributes (Total: %d), skipping\n",
				character.Attributes.GetTotal())
			skippedCount++
			continue
		}

		// Set default attributes
		fmt.Println("  → Setting default attributes...")
		character.Attributes.SetDefaultAttributes()

		// Update the character in the database
		err = db.UpdateRecord("Characters", &character)
		if err != nil {
			fmt.Printf("  ✗ Failed to update character %s: %v\n", character.Name, err)
			continue
		}

		fmt.Printf("  ✓ Updated character with default attributes:\n")
		fmt.Printf("    - Strength: %d\n", character.Attributes.Strength)
		fmt.Printf("    - Dexterity: %d\n", character.Attributes.Dexterity)
		fmt.Printf("    - Constitution: %d\n", character.Attributes.Constitution)
		fmt.Printf("    - Intelligence: %d\n", character.Attributes.Intelligence)
		fmt.Printf("    - Wisdom: %d\n", character.Attributes.Wisdom)
		fmt.Printf("    - Charisma: %d\n", character.Attributes.Charisma)
		fmt.Printf("    - Total: %d\n", character.Attributes.GetTotal())

		updatedCount++
	}

	// Summary
	fmt.Println("\n=== Update Summary ===")
	fmt.Printf("Total characters processed: %d\n", len(characters))
	fmt.Printf("Characters updated: %d\n", updatedCount)
	fmt.Printf("Characters skipped (already had attributes): %d\n", skippedCount)

	if updatedCount > 0 {
		fmt.Println("\n✓ Character attributes update completed successfully!")
	} else {
		fmt.Println("\n→ No characters needed updating.")
	}
}
