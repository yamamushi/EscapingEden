package main

import (
	"fmt"

	"github.com/yamamushi/EscapingEden/edentypes"
	"github.com/yamamushi/EscapingEden/messages"
)

func main() {
	fmt.Println("=== Equipment System Test ===")
	fmt.Println()

	// Create a test character with equipment
	character := messages.CharacterInfo{
		Name: "TestCharacter",
	}

	// Create some test items with IDs
	sword := &edentypes.Item{
		ID:         "sword_001",
		Name:       "Iron Sword",
		Equippable: true,
		Symbol:     "†",
	}

	helmet := &edentypes.Item{
		ID:         "helmet_001",
		Name:       "Iron Helmet",
		Equippable: true,
		Symbol:     "⌐",
	}

	ring := &edentypes.Item{
		ID:         "ring_001",
		Name:       "Gold Ring",
		Equippable: true,
		Symbol:     "○",
	}

	// Test equipment operations
	fmt.Println("1. Testing equipment slots...")
	slots := character.Equipment.GetAllEquipmentSlots()
	fmt.Printf("   Total equipment slots: %d\n", len(slots))
	fmt.Println("   Available slots:")
	for i, slot := range slots {
		fmt.Printf("     %2d. %s\n", i+1, slot)
	}
	fmt.Println()

	fmt.Println("2. Testing item equipping by ID...")

	// Equip items by ID
	success := character.Equipment.SetEquippedItemID("right_hand", sword.ID)
	fmt.Printf("   Equip sword to right hand: %t\n", success)

	success = character.Equipment.SetEquippedItemID("head", helmet.ID)
	fmt.Printf("   Equip helmet to head: %t\n", success)

	success = character.Equipment.SetEquippedItemID("left_ring_ring", ring.ID)
	fmt.Printf("   Equip ring to left ring finger: %t\n", success)

	// Test pinky ring
	pinkyRing := &edentypes.Item{
		ID:         "pinky_ring_001",
		Name:       "Silver Pinky Ring",
		Equippable: true,
		Symbol:     "◯",
	}

	success = character.Equipment.SetEquippedItemID("right_pinky_ring", pinkyRing.ID)
	fmt.Printf("   Equip pinky ring to right pinky: %t\n", success)

	// Test invalid slot
	success = character.Equipment.SetEquippedItemID("invalid_slot", sword.ID)
	fmt.Printf("   Equip to invalid slot: %t\n", success)
	fmt.Println()

	fmt.Println("3. Testing equipment ID retrieval...")

	// Check equipped item IDs
	equippedSwordID := character.Equipment.GetEquippedItemID("right_hand")
	fmt.Printf("   Right hand ID: %s\n", equippedSwordID)

	equippedHelmetID := character.Equipment.GetEquippedItemID("head")
	fmt.Printf("   Head ID: %s\n", equippedHelmetID)

	equippedRingID := character.Equipment.GetEquippedItemID("left_ring_ring")
	fmt.Printf("   Left ring finger ID: %s\n", equippedRingID)

	// Check empty slot
	emptySlotID := character.Equipment.GetEquippedItemID("left_hand")
	fmt.Printf("   Left hand ID (should be empty): '%s'\n", emptySlotID)
	fmt.Println()

	fmt.Println("4. Testing equipment status...")
	fmt.Printf("   Total equipped items: %d\n", character.Equipment.GetEquippedItemCount())
	fmt.Printf("   Right hand empty: %t\n", character.Equipment.IsSlotEmpty("right_hand"))
	fmt.Printf("   Left hand empty: %t\n", character.Equipment.IsSlotEmpty("left_hand"))
	fmt.Println()

	fmt.Println("5. Equipment summary (by ID):")
	for _, slot := range slots {
		itemID := character.Equipment.GetEquippedItemID(slot)
		if itemID != "" {
			fmt.Printf("   %-18s: %s\n", slot, itemID)
		}
	}

	fmt.Println()
	fmt.Println("=== Equipment System Test Complete ===")
}
