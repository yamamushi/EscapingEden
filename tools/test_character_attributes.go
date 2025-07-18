package main

import (
	"fmt"

	"github.com/yamamushi/EscapingEden/messages"
)

func main() {
	fmt.Println("=== Character Attributes Test ===")
	fmt.Println()

	// Test 1: Create a character with default attributes
	fmt.Println("1. Testing default attributes...")
	character := messages.CharacterInfo{
		Name: "TestCharacter",
	}
	character.Attributes.SetDefaultAttributes()

	fmt.Printf("   Character: %s\n", character.Name)
	fmt.Printf("   Strength: %d (modifier: %+d)\n",
		character.Attributes.Strength, character.Attributes.GetAttributeModifier("strength"))
	fmt.Printf("   Dexterity: %d (modifier: %+d)\n",
		character.Attributes.Dexterity, character.Attributes.GetAttributeModifier("dexterity"))
	fmt.Printf("   Constitution: %d (modifier: %+d)\n",
		character.Attributes.Constitution, character.Attributes.GetAttributeModifier("constitution"))
	fmt.Printf("   Intelligence: %d (modifier: %+d)\n",
		character.Attributes.Intelligence, character.Attributes.GetAttributeModifier("intelligence"))
	fmt.Printf("   Wisdom: %d (modifier: %+d)\n",
		character.Attributes.Wisdom, character.Attributes.GetAttributeModifier("wisdom"))
	fmt.Printf("   Charisma: %d (modifier: %+d)\n",
		character.Attributes.Charisma, character.Attributes.GetAttributeModifier("charisma"))
	fmt.Printf("   Total: %d\n", character.Attributes.GetTotal())
	fmt.Println()

	// Test 2: Test attribute modification
	fmt.Println("2. Testing attribute modification...")
	character.Attributes.SetAttribute("strength", 18)
	character.Attributes.SetAttribute("dex", 14)
	character.Attributes.SetAttribute("intelligence", 16)

	fmt.Printf("   Modified Strength: %d (modifier: %+d)\n",
		character.Attributes.GetAttribute("strength"), character.Attributes.GetAttributeModifier("str"))
	fmt.Printf("   Modified Dexterity: %d (modifier: %+d)\n",
		character.Attributes.GetAttribute("dex"), character.Attributes.GetAttributeModifier("dexterity"))
	fmt.Printf("   Modified Intelligence: %d (modifier: %+d)\n",
		character.Attributes.GetAttribute("intelligence"), character.Attributes.GetAttributeModifier("int"))
	fmt.Printf("   New Total: %d\n", character.Attributes.GetTotal())
	fmt.Println()

	// Test 3: Test empty attributes detection
	fmt.Println("3. Testing empty attributes detection...")
	emptyCharacter := messages.CharacterInfo{Name: "EmptyCharacter"}
	fmt.Printf("   Empty character attributes empty: %t\n", emptyCharacter.Attributes.IsEmpty())
	fmt.Printf("   Character with attributes empty: %t\n", character.Attributes.IsEmpty())
	fmt.Println()

	// Test 4: Test various attribute score modifiers
	fmt.Println("4. Testing attribute modifiers for various scores...")
	testScores := []int{1, 3, 8, 10, 12, 14, 16, 18, 20, 25, 30}

	for _, score := range testScores {
		testAttr := messages.Attributes{Strength: score}
		modifier := testAttr.GetAttributeModifier("strength")
		fmt.Printf("   Score %d → Modifier %+d\n", score, modifier)
	}

	fmt.Println()
	fmt.Println("✓ All attribute tests completed successfully!")
}
