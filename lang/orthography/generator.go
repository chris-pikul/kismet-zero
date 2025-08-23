package orthography

import (
	"fmt"
	"math/rand/v2"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
)

// OrthographyGenerator provides high-level functionality for creating and managing
// writing systems for different cultures and languages.
type OrthographyGenerator struct {
	styles map[WritingStyle]WritingSystemStyle
}

// NewOrthographyGenerator creates a new orthography generator with all available
// writing system styles.
func NewOrthographyGenerator() *OrthographyGenerator {
	generator := &OrthographyGenerator{
		styles: make(map[WritingStyle]WritingSystemStyle),
	}

	// Register default writing system styles
	generator.RegisterStyle(NewAlphabeticStyle())
	generator.RegisterStyle(NewSyllabicStyle())
	generator.RegisterStyle(NewLogographicStyle())

	return generator
}

// RegisterStyle adds a new writing system style to the generator.
func (og *OrthographyGenerator) RegisterStyle(style WritingSystemStyle) {
	og.styles[style.GetStyle()] = style
}

// GetStyle retrieves a writing system style by type.
func (og *OrthographyGenerator) GetStyle(style WritingStyle) (WritingSystemStyle, error) {
	if styleImpl, exists := og.styles[style]; exists {
		return styleImpl, nil
	}
	return nil, fmt.Errorf("writing system style '%s' not found", style)
}

// GenerateWritingSystem creates a complete writing system for a given phoneme pool
// and writing style.
func (og *OrthographyGenerator) GenerateWritingSystem(
	pool *phoneme.PhonemePool,
	style WritingStyle,
	name string,
	culture string,
	rng *rand.Rand,
) (*WritingSystem, error) {
	if pool == nil {
		return nil, fmt.Errorf("phoneme pool cannot be nil")
	}

	styleImpl, err := og.GetStyle(style)
	if err != nil {
		return nil, fmt.Errorf("failed to get writing system style: %w", err)
	}

	// Generate graphemes for the writing system
	graphemes := styleImpl.GenerateGraphemes(pool, rng)
	if len(graphemes) == 0 {
		return nil, fmt.Errorf("failed to generate graphemes for writing system")
	}

	// Create mappings from phonemes to graphemes
	mappings := og.createMappings(pool, graphemes, styleImpl)

	writingSystem := &WritingSystem{
		Style:     style,
		Graphemes: graphemes,
		Mappings:  mappings,
		Name:      name,
		Culture:   culture,
	}

	return writingSystem, nil
}

// createMappings creates orthography mappings between phonemes and graphemes.
func (og *OrthographyGenerator) createMappings(
	pool *phoneme.PhonemePool,
	graphemes []Grapheme,
	style WritingSystemStyle,
) []OrthographyMapping {
	var mappings []OrthographyMapping

	// Create mappings for consonants
	for i, consonant := range pool.Consonants {
		if i < len(graphemes) {
			mapping := OrthographyMapping{
				Phoneme:   &consonant,
				Graphemes: []Grapheme{graphemes[i]},
				Primary:   true,
			}

			// Validate the mapping for the writing system style
			if style.ValidateMapping(&mapping) {
				mappings = append(mappings, mapping)
			}
		}
	}

	// Create mappings for vowels
	vowelStart := len(pool.Consonants)
	for i, vowel := range pool.Vowels {
		graphemeIndex := vowelStart + i
		if graphemeIndex < len(graphemes) {
			mapping := OrthographyMapping{
				Phoneme:   &vowel,
				Graphemes: []Grapheme{graphemes[graphemeIndex]},
				Primary:   true,
			}

			// Validate the mapping for the writing system style
			if style.ValidateMapping(&mapping) {
				mappings = append(mappings, mapping)
			}
		}
	}

	return mappings
}

// GenerateRandomWritingSystem creates a writing system with a randomly selected style.
func (og *OrthographyGenerator) GenerateRandomWritingSystem(
	pool *phoneme.PhonemePool,
	name string,
	culture string,
	rng *rand.Rand,
) (*WritingSystem, error) {
	// Get available styles
	availableStyles := make([]WritingStyle, 0, len(og.styles))
	for style := range og.styles {
		availableStyles = append(availableStyles, style)
	}

	if len(availableStyles) == 0 {
		return nil, fmt.Errorf("no writing system styles available")
	}

	// Randomly select a style
	selectedStyle := availableStyles[rng.IntN(len(availableStyles))]

	return og.GenerateWritingSystem(pool, selectedStyle, name, culture, rng)
}

// ValidateWritingSystem checks if a writing system is valid and complete.
func (og *OrthographyGenerator) ValidateWritingSystem(system *WritingSystem) error {
	if system == nil {
		return fmt.Errorf("writing system cannot be nil")
	}

	if system.Name == "" {
		return fmt.Errorf("writing system must have a name")
	}

	if len(system.Graphemes) == 0 {
		return fmt.Errorf("writing system must have at least one grapheme")
	}

	if len(system.Mappings) == 0 {
		return fmt.Errorf("writing system must have at least one mapping")
	}

	// Check if all graphemes have valid types
	for _, grapheme := range system.Graphemes {
		if grapheme.Type == GraphemeTypeUnknown {
			return fmt.Errorf("grapheme '%s' has unknown type", grapheme.Symbol)
		}
	}

	// Check if all mappings are valid
	styleImpl, err := og.GetStyle(system.Style)
	if err != nil {
		return fmt.Errorf("invalid writing system style: %w", err)
	}

	for _, mapping := range system.Mappings {
		if !styleImpl.ValidateMapping(&mapping) {
			return fmt.Errorf("mapping for phoneme '%s' is invalid for style '%s'",
				mapping.Phoneme.Symbol, system.Style.String())
		}
	}

	return nil
}

// GetWritingSystemStats returns statistics about a writing system.
func (og *OrthographyGenerator) GetWritingSystemStats(system *WritingSystem) map[string]interface{} {
	if system == nil {
		return nil
	}

	stats := map[string]interface{}{
		"name":          system.Name,
		"style":         system.Style.String(),
		"culture":       system.Culture,
		"graphemeCount": len(system.Graphemes),
		"mappingCount":  len(system.Mappings),
	}

	// Count graphemes by type
	typeCounts := make(map[string]int)
	for _, grapheme := range system.Graphemes {
		typeStr := grapheme.Type.String()
		typeCounts[typeStr]++
	}
	stats["graphemeTypeCounts"] = typeCounts

	return stats
}
