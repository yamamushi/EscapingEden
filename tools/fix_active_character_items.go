package main

import (
	"fmt"
	"log"

	"github.com/asdine/storm/v3"
	"github.com/yamamushi/EscapingEden/edentypes"
	"github.com/yamamushi/EscapingEden/messages"
)

func main() {
	// Open the database
	db, err := storm.Open("eden.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Load item registry from JSON files
	itemRegistry := edentypes.NewItemRegistry()
	err = itemRegistry.LoadItemsFromDirectory("assets/items")
	if err != nil {
		log.Fatal("Failed to load item definitions:", err)
	}

	fmt.Printf("Loaded %d item definitions from JSON\n", len(itemRegistry.Items))
	fmt.Println("Looking for characters to fix...")

	// Check for characters
	var characters []messages.CharacterInfo
	err = db.All(&characters)
	if err != nil {
		log.Printf("Failed to get characters: %v", err)
	}

	fmt.Printf("Found %d characters in database\n", len(characters))

	if len(characters) == 0 {
		fmt.Println("No characters found in database.")
		fmt.Println("This means your character data isn't being saved yet.")
		fmt.Println("")
		fmt.Println("SOLUTION:")
		fmt.Println("1. Log in to the game with your character")
		fmt.Println("2. Make sure to save/logout properly so the character gets saved to database")
		fmt.Println("3. Then run this tool again")
		fmt.Println("")
		fmt.Println("Alternatively, we can force-create a character entry...")
		return
	}

	// Update character inventories
	updatedCount := 0
	for i, character := range characters {
		fmt.Printf("Processing character: %s\n", character.Name)
		characterUpdated := false

		// Update each item in the character's inventory
		for j, item := range character.Inventory {
			// Find the JSON definition for this item by name
			var def edentypes.ItemDefinition
			var found bool
			for _, itemDef := range itemRegistry.Items {
				if itemDef.Name == item.Name {
					def = itemDef
					found = true
					break
				}
			}

			if found {
				oldEquippable := character.Inventory[j].Equippable

				// Update properties from JSON while preserving player-specific data
				character.Inventory[j].Symbol = def.Symbol
				character.Inventory[j].FGColor = def.FGColor
				character.Inventory[j].BGColor = def.BGColor
				character.Inventory[j].Category = def.Category
				character.Inventory[j].Tags = def.Tags
				character.Inventory[j].Equippable = def.Equippable
				character.Inventory[j].MaxStack = def.MaxStack
				character.Inventory[j].Weight = def.Weight
				character.Inventory[j].Description = def.Description
				character.Inventory[j].Type = edentypes.ItemType(def.Type)
				character.Inventory[j].Stackable = def.Stackable

				if oldEquippable != def.Equippable {
					fmt.Printf("  ✓ Fixed %s: equippable %t -> %t\n", item.Name, oldEquippable, def.Equippable)
				}

				characterUpdated = true
				updatedCount++
			} else {
				fmt.Printf("  ⚠ Warning: No JSON definition found for item '%s'\n", item.Name)
			}
		}

		// Save the character if any items were updated
		if characterUpdated {
			characters[i] = character
			err = db.Update(&characters[i])
			if err != nil {
				log.Printf("Failed to save character %s: %v", character.Name, err)
			} else {
				fmt.Printf("✓ Updated character: %s (%d items fixed)\n", character.Name, len(character.Inventory))
			}
		}
	}

	fmt.Printf("\n🎉 Item fix completed! Updated %d items across all characters.\n", updatedCount)
	fmt.Println("Your pickaxe should now show 'Equippable: yes' in the game!")
}
