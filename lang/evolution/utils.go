package evolution

import (
	"fmt"
	"time"
)

// Utils provides common utility functions for the evolution package.
// This consolidates duplicate utility patterns found throughout the package.

// ID generation utilities to eliminate duplicate ID generation patterns

// GenerateChangeID generates a unique ID for linguistic changes.
// This consolidates the pattern: fmt.Sprintf("change_type_%d", time.Now().Unix())
func GenerateChangeID(prefix, identifier string) string {
	return fmt.Sprintf("%s_%s_%d", prefix, identifier, time.Now().Unix())
}

// GenerateContactID generates a unique ID for contact events.
// This consolidates the pattern: fmt.Sprintf("contact_%s_%d", sourceLang.ID.String(), time.Now().Unix())
func GenerateContactID(prefix, identifier string) string {
	return fmt.Sprintf("%s_%s_%d", prefix, identifier, time.Now().Unix())
}

// GenerateDialectID generates a unique ID for dialect events.
// This consolidates the pattern: fmt.Sprintf("dialect_formation_%s_%d", dialectID, time.Now().Unix())
func GenerateDialectID(prefix, identifier string) string {
	return fmt.Sprintf("%s_%s_%d", prefix, identifier, time.Now().Unix())
}

// GenerateEvolutionID generates a unique ID for evolution events.
// This consolidates the pattern: fmt.Sprintf("evolution_%s_%d", language.ID.String(), time.Now().Unix())
func GenerateEvolutionID(prefix, identifier string) string {
	return fmt.Sprintf("%s_%s_%d", prefix, identifier, time.Now().Unix())
}

// GenerateSoundChangeID generates a unique ID for sound changes.
// This consolidates the pattern: fmt.Sprintf("sound_change_%s_%d", rule.ID, time.Now().Unix())
func GenerateSoundChangeID(prefix, identifier string) string {
	return fmt.Sprintf("%s_%s_%d", prefix, identifier, time.Now().Unix())
}

// GenerateTimestampID generates a unique ID with current timestamp.
// This is a generic function for any ID generation that needs a timestamp.
func GenerateTimestampID(prefix, identifier string) string {
	return fmt.Sprintf("%s_%s_%d", prefix, identifier, time.Now().Unix())
}

// GenerateNanoID generates a unique ID with nanosecond precision.
// This is for cases where higher precision is needed.
func GenerateNanoID(prefix, identifier string) string {
	return fmt.Sprintf("%s_%s_%d", prefix, identifier, time.Now().UnixNano())
}

// Enum utilities to eliminate duplicate String() method implementations

// EnumString provides a generic way to convert enum values to strings.
// This consolidates the pattern repeated 7 times across different enum types.
func EnumString(value int, enumArray []string, defaultVal string) string {
	if value < 0 || value >= len(enumArray) {
		return defaultVal
	}
	return enumArray[value]
}

// EnumStringByte provides a generic way to convert byte-based enum values to strings.
// This is specifically for byte-based enums like ContactType, ChangeType, etc.
func EnumStringByte(value byte, enumArray []string, defaultVal string) string {
	if int(value) < 0 || int(value) >= len(enumArray) {
		return defaultVal
	}
	return enumArray[value]
}

// Change creation factory functions to eliminate duplicate change creation patterns

// NewLinguisticChange creates a new LinguisticChange with common fields.
// This consolidates the pattern repeated 10+ times across different engines.
func NewLinguisticChange(changeType ChangeType, direction ChangeDirection,
	description, era, trigger string, intensity float32) *LinguisticChange {
	return &LinguisticChange{
		ID:          GenerateChangeID("change", fmt.Sprintf("%d", time.Now().Unix())),
		Type:        changeType,
		Direction:   direction,
		Description: description,
		Details:     "Change applied through evolution system",
		Timestamp:   time.Now(),
		Era:         era,
		Trigger:     trigger,
		Intensity:   intensity,
	}
}

// NewSoundChange creates a new LinguisticChange specifically for sound changes.
func NewSoundChange(description, era string, intensity float32) *LinguisticChange {
	return NewLinguisticChange(
		ChangeTypeSoundShift,
		ChangeDirectionModifying,
		description,
		era,
		"natural_evolution",
		intensity,
	)
}

// NewMorphologicalChange creates a new LinguisticChange specifically for morphological changes.
func NewMorphologicalChange(description, era string, intensity float32) *LinguisticChange {
	return NewLinguisticChange(
		ChangeTypeMorphological,
		ChangeDirectionModifying,
		description,
		era,
		"natural_evolution",
		intensity,
	)
}

// NewContactChange creates a new LinguisticChange specifically for contact-induced changes.
func NewContactChange(changeType ChangeType, description, era string, intensity float32, cultureInfluence string) *LinguisticChange {
	change := NewLinguisticChange(
		changeType,
		ChangeDirectionAdditive,
		description,
		era,
		"contact_evolution",
		intensity,
	)
	change.CultureInfluence = cultureInfluence
	return change
}

// NewAdaptationChange creates a new LinguisticChange specifically for adaptation changes.
func NewAdaptationChange(changeType ChangeType, description, era string, intensity float32) *LinguisticChange {
	return NewLinguisticChange(
		changeType,
		ChangeDirectionModifying,
		description,
		era,
		"adaptation",
		intensity,
	)
}
