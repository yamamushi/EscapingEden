package messages

import (
	"strings"
	"time"

	"github.com/yamamushi/EscapingEden/edentypes"
	"github.com/yamamushi/EscapingEden/ui/util"
)

// Equipment represents all equipped items on a character using item IDs for database references
type Equipment struct {
	// Hands
	LeftHandID  string `json:"left_hand_id,omitempty"`
	RightHandID string `json:"right_hand_id,omitempty"`

	// Head and Face
	HeadID string `json:"head_id,omitempty"`
	FaceID string `json:"face_id,omitempty"`
	NeckID string `json:"neck_id,omitempty"`

	// Upper Body
	ShouldersID string `json:"shoulders_id,omitempty"`
	BackID      string `json:"back_id,omitempty"`
	BodyID      string `json:"body_id,omitempty"`

	// Arms
	RightArmID string `json:"right_arm_id,omitempty"`
	LeftArmID  string `json:"left_arm_id,omitempty"`

	// Lower Body
	WaistID string `json:"waist_id,omitempty"`

	// Legs
	RightLegID string `json:"right_leg_id,omitempty"`
	LeftLegID  string `json:"left_leg_id,omitempty"`

	// Feet
	RightFootID string `json:"right_foot_id,omitempty"`
	LeftFootID  string `json:"left_foot_id,omitempty"`

	// Rings - Left Hand
	LeftIndexRingID  string `json:"left_index_ring_id,omitempty"`
	LeftMiddleRingID string `json:"left_middle_ring_id,omitempty"`
	LeftRingRingID   string `json:"left_ring_ring_id,omitempty"`
	LeftPinkyRingID  string `json:"left_pinky_ring_id,omitempty"`

	// Rings - Right Hand
	RightIndexRingID  string `json:"right_index_ring_id,omitempty"`
	RightMiddleRingID string `json:"right_middle_ring_id,omitempty"`
	RightRingRingID   string `json:"right_ring_ring_id,omitempty"`
	RightPinkyRingID  string `json:"right_pinky_ring_id,omitempty"`
}

// Attributes represents a character's core attributes
type Attributes struct {
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
}

type CharacterInfo struct {
	ID          string `storm:"index"`
	UserID      string
	Name        string         `storm:"unique"`
	FGColor     util.ColorCode // The escape code of the FG color of the character
	BGColor     util.ColorCode
	InventoryID string `storm:"unique"`
	Initialized bool
	Position    struct {
		MapChunkID string
		X          int
		Y          int
		Z          int
	}
	LastLoginTime time.Time
	FirstLogin    int // 0 = false, 1 = true
	Error         string

	// The character's inventory
	Inventory    []edentypes.Item
	CurrentMapID string

	// The character's equipped items
	Equipment Equipment

	// The character's attributes
	Attributes Attributes
}

type PlayerViewHistoryCoordinate struct {
	Gx, Gy, Gz, X, Y, Z int
}

func (c *CharacterInfo) GetID() string {
	return c.ID
}

func (c *CharacterInfo) GetColorFG() util.ColorCode {
	return c.FGColor
}

func (c *CharacterInfo) GetColorBG() util.ColorCode {
	return c.BGColor
}

func (c *CharacterInfo) GetName() string {
	return c.Name
}

func (c *CharacterInfo) GetInventoryID() string {
	return c.InventoryID
}

// Equipment helper methods

// GetEquippedItemID returns the ID of the item equipped in the specified slot
func (e *Equipment) GetEquippedItemID(slot string) string {
	switch slot {
	case "left_hand":
		return e.LeftHandID
	case "right_hand":
		return e.RightHandID
	case "head":
		return e.HeadID
	case "face":
		return e.FaceID
	case "neck":
		return e.NeckID
	case "shoulders":
		return e.ShouldersID
	case "back":
		return e.BackID
	case "body":
		return e.BodyID
	case "right_arm":
		return e.RightArmID
	case "left_arm":
		return e.LeftArmID
	case "waist":
		return e.WaistID
	case "right_leg":
		return e.RightLegID
	case "left_leg":
		return e.LeftLegID
	case "right_foot":
		return e.RightFootID
	case "left_foot":
		return e.LeftFootID
	case "left_index_ring":
		return e.LeftIndexRingID
	case "left_middle_ring":
		return e.LeftMiddleRingID
	case "left_ring_ring":
		return e.LeftRingRingID
	case "left_pinky_ring":
		return e.LeftPinkyRingID
	case "right_index_ring":
		return e.RightIndexRingID
	case "right_middle_ring":
		return e.RightMiddleRingID
	case "right_ring_ring":
		return e.RightRingRingID
	case "right_pinky_ring":
		return e.RightPinkyRingID
	default:
		return ""
	}
}

// SetEquippedItemID equips an item by ID in the specified slot
func (e *Equipment) SetEquippedItemID(slot string, itemID string) bool {
	switch slot {
	case "left_hand":
		e.LeftHandID = itemID
	case "right_hand":
		e.RightHandID = itemID
	case "head":
		e.HeadID = itemID
	case "face":
		e.FaceID = itemID
	case "neck":
		e.NeckID = itemID
	case "shoulders":
		e.ShouldersID = itemID
	case "back":
		e.BackID = itemID
	case "body":
		e.BodyID = itemID
	case "right_arm":
		e.RightArmID = itemID
	case "left_arm":
		e.LeftArmID = itemID
	case "waist":
		e.WaistID = itemID
	case "right_leg":
		e.RightLegID = itemID
	case "left_leg":
		e.LeftLegID = itemID
	case "right_foot":
		e.RightFootID = itemID
	case "left_foot":
		e.LeftFootID = itemID
	case "left_index_ring":
		e.LeftIndexRingID = itemID
	case "left_middle_ring":
		e.LeftMiddleRingID = itemID
	case "left_ring_ring":
		e.LeftRingRingID = itemID
	case "left_pinky_ring":
		e.LeftPinkyRingID = itemID
	case "right_index_ring":
		e.RightIndexRingID = itemID
	case "right_middle_ring":
		e.RightMiddleRingID = itemID
	case "right_ring_ring":
		e.RightRingRingID = itemID
	case "right_pinky_ring":
		e.RightPinkyRingID = itemID
	default:
		return false
	}
	return true
}

// GetAllEquipmentSlots returns a list of all equipment slot names
func (e *Equipment) GetAllEquipmentSlots() []string {
	return []string{
		"left_hand", "right_hand",
		"head", "face", "neck",
		"shoulders", "back", "body",
		"right_arm", "left_arm",
		"waist",
		"right_leg", "left_leg",
		"right_foot", "left_foot",
		"left_index_ring", "left_middle_ring", "left_ring_ring", "left_pinky_ring",
		"right_index_ring", "right_middle_ring", "right_ring_ring", "right_pinky_ring",
	}
}

// IsSlotEmpty checks if an equipment slot is empty
func (e *Equipment) IsSlotEmpty(slot string) bool {
	return e.GetEquippedItemID(slot) == ""
}

// GetEquippedItemCount returns the total number of equipped items
func (e *Equipment) GetEquippedItemCount() int {
	count := 0
	for _, slot := range e.GetAllEquipmentSlots() {
		if !e.IsSlotEmpty(slot) {
			count++
		}
	}
	return count
}

// Attributes helper methods

// GetAttribute returns the value of a specific attribute by name
func (a *Attributes) GetAttribute(name string) int {
	switch strings.ToLower(name) {
	case "strength", "str":
		return a.Strength
	case "dexterity", "dex":
		return a.Dexterity
	case "constitution", "con":
		return a.Constitution
	case "intelligence", "int":
		return a.Intelligence
	case "wisdom", "wis":
		return a.Wisdom
	case "charisma", "cha":
		return a.Charisma
	default:
		return 0
	}
}

// SetAttribute sets the value of a specific attribute by name
func (a *Attributes) SetAttribute(name string, value int) bool {
	switch strings.ToLower(name) {
	case "strength", "str":
		a.Strength = value
	case "dexterity", "dex":
		a.Dexterity = value
	case "constitution", "con":
		a.Constitution = value
	case "intelligence", "int":
		a.Intelligence = value
	case "wisdom", "wis":
		a.Wisdom = value
	case "charisma", "cha":
		a.Charisma = value
	default:
		return false
	}
	return true
}

// GetAttributeModifier returns the D&D-style modifier for an attribute (-5 to +10 for scores 1-30)
func (a *Attributes) GetAttributeModifier(name string) int {
	score := a.GetAttribute(name)
	return (score - 10) / 2
}

// GetTotal returns the sum of all attributes
func (a *Attributes) GetTotal() int {
	return a.Strength + a.Dexterity + a.Constitution + a.Intelligence + a.Wisdom + a.Charisma
}

// SetDefaultAttributes sets default starting attributes (typically 10 for average)
func (a *Attributes) SetDefaultAttributes() {
	a.Strength = 10
	a.Dexterity = 10
	a.Constitution = 10
	a.Intelligence = 10
	a.Wisdom = 10
	a.Charisma = 10
}

// IsEmpty returns true if all attributes are zero (indicating uninitialized)
func (a *Attributes) IsEmpty() bool {
	return a.Strength == 0 && a.Dexterity == 0 && a.Constitution == 0 &&
		a.Intelligence == 0 && a.Wisdom == 0 && a.Charisma == 0
}
