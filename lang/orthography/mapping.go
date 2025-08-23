package orthography

import (
	"fmt"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
)

// OrthographyMapper provides utilities for managing grapheme-phoneme mappings
// and converting between spoken and written forms.
type OrthographyMapper struct {
	mappings map[string]*OrthographyMapping // key: phoneme symbol
}

// NewOrthographyMapper creates a new empty orthography mapper.
func NewOrthographyMapper() *OrthographyMapper {
	return &OrthographyMapper{
		mappings: make(map[string]*OrthographyMapping),
	}
}

// AddMapping adds a new mapping from a phoneme to graphemes.
// If a mapping already exists for the phoneme, it will be replaced.
func (om *OrthographyMapper) AddMapping(mapping *OrthographyMapping) error {
	if mapping == nil || mapping.Phoneme == nil {
		return fmt.Errorf("mapping and phoneme cannot be nil")
	}

	if mapping.Phoneme.Symbol == "" {
		return fmt.Errorf("phoneme symbol cannot be empty")
	}

	om.mappings[mapping.Phoneme.Symbol] = mapping
	return nil
}

// GetMapping retrieves the mapping for a specific phoneme by its symbol.
// Returns nil if no mapping exists.
func (om *OrthographyMapper) GetMapping(phonemeSymbol string) *OrthographyMapping {
	return om.mappings[phonemeSymbol]
}

// GetGraphemesForPhoneme retrieves all graphemes that can represent a given phoneme.
// Returns an empty slice if no mapping exists.
func (om *OrthographyMapper) GetGraphemesForPhoneme(phonemeSymbol string) []Grapheme {
	mapping := om.GetMapping(phonemeSymbol)
	if mapping == nil {
		return []Grapheme{}
	}
	return mapping.Graphemes
}

// GetPrimaryGrapheme retrieves the primary grapheme for a given phoneme.
// Returns nil if no mapping exists or no primary grapheme is set.
func (om *OrthographyMapper) GetPrimaryGrapheme(phonemeSymbol string) *Grapheme {
	mapping := om.GetMapping(phonemeSymbol)
	if mapping == nil {
		return nil
	}

	for _, grapheme := range mapping.Graphemes {
		if mapping.Primary {
			return &grapheme
		}
	}

	// If no primary is explicitly set, return the first grapheme
	if len(mapping.Graphemes) > 0 {
		return &mapping.Graphemes[0]
	}

	return nil
}

// HasMapping checks if a mapping exists for the given phoneme symbol.
func (om *OrthographyMapper) HasMapping(phonemeSymbol string) bool {
	_, exists := om.mappings[phonemeSymbol]
	return exists
}

// RemoveMapping removes the mapping for a specific phoneme.
// Returns true if a mapping was removed, false if none existed.
func (om *OrthographyMapper) RemoveMapping(phonemeSymbol string) bool {
	if _, exists := om.mappings[phonemeSymbol]; exists {
		delete(om.mappings, phonemeSymbol)
		return true
	}
	return false
}

// GetAllMappings returns all current mappings.
func (om *OrthographyMapper) GetAllMappings() []*OrthographyMapping {
	mappings := make([]*OrthographyMapping, 0, len(om.mappings))
	for _, mapping := range om.mappings {
		mappings = append(mappings, mapping)
	}
	return mappings
}

// Count returns the total number of mappings.
func (om *OrthographyMapper) Count() int {
	return len(om.mappings)
}

// IsComplete checks if the mapper has mappings for all phonemes in the given phoneme pool.
func (om *OrthographyMapper) IsComplete(pool *phoneme.PhonemePool) bool {
	if pool == nil {
		return false
	}

	// Check consonants
	for _, consonant := range pool.Consonants {
		if !om.HasMapping(consonant.Symbol) {
			return false
		}
	}

	// Check vowels
	for _, vowel := range pool.Vowels {
		if !om.HasMapping(vowel.Symbol) {
			return false
		}
	}

	return true
}

// GetCoveragePercentage returns the percentage of phonemes in the pool that have mappings.
func (om *OrthographyMapper) GetCoveragePercentage(pool *phoneme.PhonemePool) float64 {
	if pool == nil || pool.Size() == 0 {
		return 0.0
	}

	mappedCount := 0
	totalCount := pool.Size()

	// Count mapped consonants
	for _, consonant := range pool.Consonants {
		if om.HasMapping(consonant.Symbol) {
			mappedCount++
		}
	}

	// Count mapped vowels
	for _, vowel := range pool.Vowels {
		if om.HasMapping(vowel.Symbol) {
			mappedCount++
		}
	}

	return float64(mappedCount) / float64(totalCount) * 100.0
}
