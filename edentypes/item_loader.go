package edentypes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// ItemDefinition represents the JSON structure for item definitions
type ItemDefinition struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Weight      float64  `json:"weight"`
	Type        int      `json:"type"` // ItemType as int for JSON
	Stackable   bool     `json:"stackable"`
	Equippable  bool     `json:"equippable"`
	Symbol      string   `json:"symbol"`
	FGColor     int      `json:"fg_color"` // Foreground color (integer)
	BGColor     int      `json:"bg_color"` // Background color (integer)
	Value       int      `json:"value"`
	Durability  int      `json:"durability"`
	MaxStack    int      `json:"max_stack"`
	Rarity      string   `json:"rarity"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
}

// ItemRegistry holds all loaded item definitions
type ItemRegistry struct {
	Items map[string]ItemDefinition
}

// NewItemRegistry creates a new item registry
func NewItemRegistry() *ItemRegistry {
	return &ItemRegistry{
		Items: make(map[string]ItemDefinition),
	}
}

// LoadItemsFromDirectory loads all item JSON files from the specified directory
func (ir *ItemRegistry) LoadItemsFromDirectory(directory string) error {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("failed to read items directory %s: %w", directory, err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".json") {
			filePath := filepath.Join(directory, file.Name())
			if err := ir.LoadItemFromFile(filePath); err != nil {
				return fmt.Errorf("failed to load item from %s: %w", filePath, err)
			}
		}
	}

	return nil
}

// LoadItemFromFile loads a single item definition from a JSON file
func (ir *ItemRegistry) LoadItemFromFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var itemDef ItemDefinition
	if err := json.Unmarshal(data, &itemDef); err != nil {
		return fmt.Errorf("failed to parse JSON from %s: %w", filePath, err)
	}

	ir.Items[itemDef.ID] = itemDef
	return nil
}

// GetItemDefinition returns the item definition for the given ID
func (ir *ItemRegistry) GetItemDefinition(id string) (ItemDefinition, bool) {
	item, exists := ir.Items[id]
	return item, exists
}

// CreateItemFromDefinition creates an Item instance from an ItemDefinition
func (ir *ItemRegistry) CreateItemFromDefinition(id string, hotkey string) (*Item, error) {
	def, exists := ir.GetItemDefinition(id)
	if !exists {
		return nil, fmt.Errorf("item definition not found for ID: %s", id)
	}

	// Convert attributes (for now, just copy from the old system)
	attributes := make(map[string]bool)
	for _, tag := range def.Tags {
		attributes[tag] = true
	}

	item := &Item{
		ID:          def.ID,
		Name:        def.Name,
		Description: def.Description,
		Weight:      def.Weight,
		Type:        ItemType(def.Type),
		Stackable:   def.Stackable,
		Equippable:  def.Equippable,
		Hotkey:      hotkey,
		Attributes:  attributes,
		Symbol:      def.Symbol,
		FGColor:     def.FGColor,
		BGColor:     def.BGColor,
		Value:       def.Value,
		Durability:  def.Durability,
		MaxStack:    def.MaxStack,
		Rarity:      def.Rarity,
		Category:    def.Category,
		Tags:        def.Tags,
	}

	return item, nil
}

// GetAllItemIDs returns a slice of all loaded item IDs
func (ir *ItemRegistry) GetAllItemIDs() []string {
	ids := make([]string, 0, len(ir.Items))
	for id := range ir.Items {
		ids = append(ids, id)
	}
	return ids
}

// GetItemsByCategory returns all items in a specific category
func (ir *ItemRegistry) GetItemsByCategory(category string) []ItemDefinition {
	var items []ItemDefinition
	for _, item := range ir.Items {
		if item.Category == category {
			items = append(items, item)
		}
	}
	return items
}

// GetItemsByRarity returns all items of a specific rarity
func (ir *ItemRegistry) GetItemsByRarity(rarity string) []ItemDefinition {
	var items []ItemDefinition
	for _, item := range ir.Items {
		if item.Rarity == rarity {
			items = append(items, item)
		}
	}
	return items
}
