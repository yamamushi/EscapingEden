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
	fmt.Println("Syncing character inventories with JSON definitions...")

	updatedCount := 0

	// Update character inventories
	var characters []messages.CharacterInfo
	err = db.All(&characters)
	if err != nil {
		log.Printf("Failed to get characters: %v", err)
		return
	}

	for i, character := range characters {
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
				// Update item properties from JSON while preserving player-specific data
				oldEquippable := character.Inventory[j].Equippable

				// Update properties from JSON
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

				// Preserve player-specific fields like ID, Hotkey, current Durability
				// (Don't overwrite durability if item has been used)

				if oldEquippable != def.Equippable {
					fmt.Printf("  Updated %s equippable: %t -> %t\n", item.Name, oldEquippable, def.Equippable)
				}

				characterUpdated = true
				updatedCount++
			} else {
				fmt.Printf("  Warning: No JSON definition found for item '%s'\n", item.Name)
			}
		}

		// Save the character if any items were updated
		if characterUpdated {
			characters[i] = character
			err = db.Update(&characters[i])
			if err != nil {
				log.Printf("Failed to save character %s: %v", character.Name, err)
			} else {
				fmt.Printf("Updated character: %s (%d items)\n", character.Name, len(character.Inventory))
			}
		}
	}

	// Also update standalone items in the database
	var items []edentypes.Item
	err = db.All(&items)
	if err == nil {
		for _, item := range items {
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
				// Update properties from JSON while preserving instance-specific data
				item.Symbol = def.Symbol
				item.FGColor = def.FGColor
				item.BGColor = def.BGColor
				item.Category = def.Category
				item.Tags = def.Tags
				item.Equippable = def.Equippable
				item.MaxStack = def.MaxStack
				item.Weight = def.Weight
				item.Description = def.Description
				item.Type = edentypes.ItemType(def.Type)
				item.Stackable = def.Stackable

				err = db.Update(&item)
				if err != nil {
					log.Printf("Failed to update standalone item %s: %v", item.Name, err)
				}
			}
		}
	}

	fmt.Printf("Item sync completed! Updated %d items across all characters.\n", updatedCount)
}
