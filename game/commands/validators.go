package commands

import (
	"fmt"
	"reflect"
	"regexp"
)

// Common validators for command parameters

// StringValidator validates string parameters
type StringValidator struct {
	MinLength int
	MaxLength int
	Pattern   *regexp.Regexp
	Required  bool
}

func (v *StringValidator) Validate(value interface{}) error {
	if value == nil {
		if v.Required {
			return NewCommandError(ErrorMissingParameter, "parameter is required")
		}
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return NewCommandError(ErrorParameterType, "parameter must be a string")
	}

	if v.MinLength > 0 && len(str) < v.MinLength {
		return NewCommandError(ErrorInvalidParameter,
			fmt.Sprintf("string must be at least %d characters", v.MinLength))
	}

	if v.MaxLength > 0 && len(str) > v.MaxLength {
		return NewCommandError(ErrorInvalidParameter,
			fmt.Sprintf("string must be at most %d characters", v.MaxLength))
	}

	if v.Pattern != nil && !v.Pattern.MatchString(str) {
		return NewCommandError(ErrorInvalidParameter, "string format is invalid")
	}

	return nil
}

func (v *StringValidator) GetErrorMessage() string {
	return "invalid string parameter"
}

// IntValidator validates integer parameters
type IntValidator struct {
	Min      *int
	Max      *int
	Required bool
}

func (v *IntValidator) Validate(value interface{}) error {
	if value == nil {
		if v.Required {
			return NewCommandError(ErrorMissingParameter, "parameter is required")
		}
		return nil
	}

	// Handle different numeric types
	var intVal int
	switch val := value.(type) {
	case int:
		intVal = val
	case int32:
		intVal = int(val)
	case int64:
		intVal = int(val)
	case float64:
		intVal = int(val)
	default:
		return NewCommandError(ErrorParameterType, "parameter must be an integer")
	}

	if v.Min != nil && intVal < *v.Min {
		return NewCommandError(ErrorInvalidParameter,
			fmt.Sprintf("value must be at least %d", *v.Min))
	}

	if v.Max != nil && intVal > *v.Max {
		return NewCommandError(ErrorInvalidParameter,
			fmt.Sprintf("value must be at most %d", *v.Max))
	}

	return nil
}

func (v *IntValidator) GetErrorMessage() string {
	return "invalid integer parameter"
}

// DirectionValidator validates movement directions
type DirectionValidator struct{}

var validDirections = map[string]bool{
	"n": true, "north": true,
	"s": true, "south": true,
	"e": true, "east": true,
	"w": true, "west": true,
	"ne": true, "northeast": true,
	"nw": true, "northwest": true,
	"se": true, "southeast": true,
	"sw": true, "southwest": true,
	// Vi-style directions (h=west, j=south, k=north, l=east, y=nw, u=ne, b=sw)
	"h": true, "j": true, "k": true, "l": true,
	"y": true, "u": true, "b": true,
}

func (v *DirectionValidator) Validate(value interface{}) error {
	if value == nil {
		return NewCommandError(ErrorMissingParameter, "direction is required")
	}

	direction, ok := value.(string)
	if !ok {
		return NewCommandError(ErrorParameterType, "direction must be a string")
	}

	if !validDirections[direction] {
		return NewCommandError(ErrorInvalidDirection,
			"invalid direction (use n/s/e/w/ne/nw/se/sw or h/j/k/l/y/u/b)")
	}

	return nil
}

func (v *DirectionValidator) GetErrorMessage() string {
	return "invalid direction"
}

// IDValidator validates ID parameters (UUIDs, etc.)
type IDValidator struct {
	Pattern  *regexp.Regexp
	Required bool
}

func NewIDValidator(required bool) *IDValidator {
	// Basic ID pattern - alphanumeric with hyphens
	pattern := regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)
	return &IDValidator{
		Pattern:  pattern,
		Required: required,
	}
}

func (v *IDValidator) Validate(value interface{}) error {
	if value == nil {
		if v.Required {
			return NewCommandError(ErrorMissingParameter, "ID is required")
		}
		return nil
	}

	id, ok := value.(string)
	if !ok {
		return NewCommandError(ErrorParameterType, "ID must be a string")
	}

	if id == "" {
		if v.Required {
			return NewCommandError(ErrorMissingParameter, "ID cannot be empty")
		}
		return nil
	}

	if v.Pattern != nil && !v.Pattern.MatchString(id) {
		return NewCommandError(ErrorInvalidParameter, "ID format is invalid")
	}

	return nil
}

func (v *IDValidator) GetErrorMessage() string {
	return "invalid ID parameter"
}

// ValidateParameters validates command parameters against rules
func ValidateParameters(params map[string]interface{}, rules []ValidationRule) error {
	for _, rule := range rules {
		value, exists := params[rule.Parameter]

		if !exists || value == nil {
			if rule.Required {
				return NewCommandError(ErrorMissingParameter,
					fmt.Sprintf("required parameter '%s' is missing", rule.Parameter))
			}
			continue
		}

		// Type validation
		if rule.Type != "" {
			if err := validateParameterType(value, rule.Type); err != nil {
				return err.WithContext("parameter", rule.Parameter)
			}
		}

		// Custom validation
		if rule.Validator != nil {
			if err := rule.Validator(value); err != nil {
				if cmdErr, ok := err.(*CommandError); ok {
					return cmdErr.WithContext("parameter", rule.Parameter)
				}
				return NewCommandError(ErrorInvalidParameter, err.Error()).
					WithContext("parameter", rule.Parameter)
			}
		}
	}

	return nil
}

// validateParameterType validates parameter type
func validateParameterType(value interface{}, expectedType string) *CommandError {
	actualType := reflect.TypeOf(value).String()

	switch expectedType {
	case "string":
		if _, ok := value.(string); !ok {
			return NewCommandError(ErrorParameterType,
				fmt.Sprintf("expected string, got %s", actualType))
		}
	case "int":
		switch value.(type) {
		case int, int32, int64, float64:
			// Accept various numeric types
		default:
			return NewCommandError(ErrorParameterType,
				fmt.Sprintf("expected integer, got %s", actualType))
		}
	case "bool":
		if _, ok := value.(bool); !ok {
			return NewCommandError(ErrorParameterType,
				fmt.Sprintf("expected boolean, got %s", actualType))
		}
	case "float":
		switch value.(type) {
		case float32, float64, int, int32, int64:
			// Accept numeric types
		default:
			return NewCommandError(ErrorParameterType,
				fmt.Sprintf("expected number, got %s", actualType))
		}
	}

	return nil
}

// Helper functions for creating common validators

// RequiredString creates a required string validator
func RequiredString(minLen, maxLen int) ValidationRule {
	return ValidationRule{
		Required: true,
		Type:     "string",
		Validator: func(value interface{}) error {
			validator := &StringValidator{
				MinLength: minLen,
				MaxLength: maxLen,
				Required:  true,
			}
			return validator.Validate(value)
		},
	}
}

// OptionalString creates an optional string validator
func OptionalString(minLen, maxLen int) ValidationRule {
	return ValidationRule{
		Required: false,
		Type:     "string",
		Validator: func(value interface{}) error {
			validator := &StringValidator{
				MinLength: minLen,
				MaxLength: maxLen,
				Required:  false,
			}
			return validator.Validate(value)
		},
	}
}

// RequiredDirection creates a required direction validator
func RequiredDirection() ValidationRule {
	return ValidationRule{
		Required: true,
		Type:     "string",
		Validator: func(value interface{}) error {
			validator := &DirectionValidator{}
			return validator.Validate(value)
		},
	}
}

// RequiredID creates a required ID validator
func RequiredID() ValidationRule {
	return ValidationRule{
		Required: true,
		Type:     "string",
		Validator: func(value interface{}) error {
			validator := NewIDValidator(true)
			return validator.Validate(value)
		},
	}
}
