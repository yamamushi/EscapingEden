# Character Attributes System

This document describes the character attributes system in Escaping Eden.

## Overview

Characters now have six core attributes that define their capabilities:

- **Strength (STR)**: Physical power and melee combat ability
- **Dexterity (DEX)**: Agility, reflexes, and ranged combat ability  
- **Constitution (CON)**: Health, stamina, and physical resilience
- **Intelligence (INT)**: Reasoning ability, memory, and magical power
- **Wisdom (WIS)**: Awareness, insight, and willpower
- **Charisma (CHA)**: Force of personality and social skills

## Default Values

New characters start with all attributes set to **10** (average human level).

## Attribute Modifiers

The system uses D&D-style modifiers calculated as: `(attribute - 10) / 2`

| Score | Modifier | Score | Modifier |
|-------|----------|-------|----------|
| 1     | -5       | 16    | +3       |
| 3     | -4       | 18    | +4       |
| 8     | -1       | 20    | +5       |
| 10    | +0       | 25    | +7       |
| 12    | +1       | 30    | +10      |
| 14    | +2       |       |          |

## API Reference

### Attributes Struct

```go
type Attributes struct {
    Strength     int `json:"strength"`
    Dexterity    int `json:"dexterity"`
    Constitution int `json:"constitution"`
    Intelligence int `json:"intelligence"`
    Wisdom       int `json:"wisdom"`
    Charisma     int `json:"charisma"`
}
```

### Methods

#### GetAttribute(name string) int
Returns the value of a specific attribute. Accepts both full names and abbreviations:
- "strength" or "str"
- "dexterity" or "dex"  
- "constitution" or "con"
- "intelligence" or "int"
- "wisdom" or "wis"
- "charisma" or "cha"

#### SetAttribute(name string, value int) bool
Sets the value of a specific attribute. Returns `true` if successful, `false` if attribute name is invalid.

#### GetAttributeModifier(name string) int
Returns the D&D-style modifier for an attribute score.

#### GetTotal() int
Returns the sum of all six attributes.

#### SetDefaultAttributes()
Sets all attributes to 10 (average).

#### IsEmpty() bool
Returns `true` if all attributes are zero (uninitialized).

## Usage Examples

### Creating a Character with Default Attributes
```go
character := messages.CharacterInfo{Name: "Hero"}
character.Attributes.SetDefaultAttributes()
// All attributes are now 10
```

### Modifying Attributes
```go
character.Attributes.SetAttribute("strength", 18)
character.Attributes.SetAttribute("dex", 14)
character.Attributes.SetAttribute("intelligence", 16)
```

### Getting Attribute Values and Modifiers
```go
str := character.Attributes.GetAttribute("strength")     // 18
strMod := character.Attributes.GetAttributeModifier("str") // +4
total := character.Attributes.GetTotal()                 // Sum of all attributes
```

## Database Migration

### For Existing Characters

Use the provided utility tool to update existing character records:

```bash
# Update all characters with default attributes
go run tools/update_character_attributes.go

# Use custom config file
go run tools/update_character_attributes.go /path/to/server.conf
```

The tool will:
- Connect to your database
- Find all existing characters
- Add default attributes (10 in each) to characters that don't have them
- Skip characters that already have attributes
- Provide a summary of changes made

### For New Characters

New characters automatically receive default attributes during creation.

## Testing

Test the attributes system:

```bash
go run tools/test_character_attributes.go
```

This will test:
- Default attribute initialization
- Attribute modification
- Modifier calculations
- Empty attribute detection

## Integration

The attributes are automatically included in:
- Character creation process
- Character database records
- Character info messages
- Game state management

Characters created after this update will automatically have attributes. Existing characters can be updated using the migration tool.