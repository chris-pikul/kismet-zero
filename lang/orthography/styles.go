package orthography

import (
	"math/rand/v2"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
)

// WritingSystemStyle defines the interface for different writing system implementations.
type WritingSystemStyle interface {
	// GenerateGraphemes creates graphemes appropriate for this writing system style.
	GenerateGraphemes(pool *phoneme.PhonemePool, rng *rand.Rand) []Grapheme

	// GetStyle returns the WritingStyle enum value for this implementation.
	GetStyle() WritingStyle

	// ValidateMapping checks if a mapping is valid for this writing system style.
	ValidateMapping(mapping *OrthographyMapping) bool
}

// AlphabeticStyle implements an alphabetic writing system where each grapheme
// typically represents a single phoneme.
type AlphabeticStyle struct {
	// Character sets for different phoneme types
	consonantChars []string
	vowelChars     []string
	// Weight adjustments for character selection
	consonantWeight float32
	vowelWeight     float32
}

// NewAlphabeticStyle creates a new alphabetic writing system style with default
// Latin-like character sets.
func NewAlphabeticStyle() *AlphabeticStyle {
	return &AlphabeticStyle{
		consonantChars: []string{
			"b", "c", "d", "f", "g", "h", "j", "k", "l", "m",
			"n", "p", "q", "r", "s", "t", "v", "w", "x", "z",
		},
		vowelChars:      []string{"a", "e", "i", "o", "u", "y"},
		consonantWeight: 1.0,
		vowelWeight:     1.0,
	}
}

// GenerateGraphemes creates alphabetic graphemes for the given phoneme pool.
func (as *AlphabeticStyle) GenerateGraphemes(pool *phoneme.PhonemePool, rng *rand.Rand) []Grapheme {
	var graphemes []Grapheme

	// Generate consonant graphemes
	for i, consonant := range pool.Consonants {
		if i < len(as.consonantChars) {
			grapheme := Grapheme{
				Symbol:     as.consonantChars[i],
				Type:       GraphemeTypeLetter,
				PhonemeRef: &consonant,
				Weight:     as.consonantWeight,
			}
			graphemes = append(graphemes, grapheme)
		}
	}

	// Generate vowel graphemes
	for i, vowel := range pool.Vowels {
		if i < len(as.vowelChars) {
			grapheme := Grapheme{
				Symbol:     as.vowelChars[i],
				Type:       GraphemeTypeLetter,
				PhonemeRef: &vowel,
				Weight:     as.vowelWeight,
			}
			graphemes = append(graphemes, grapheme)
		}
	}

	return graphemes
}

// GetStyle returns the alphabetic writing style.
func (as *AlphabeticStyle) GetStyle() WritingStyle {
	return WritingStyleAlphabetic
}

// ValidateMapping checks if a mapping is valid for alphabetic systems.
func (as *AlphabeticStyle) ValidateMapping(mapping *OrthographyMapping) bool {
	if mapping == nil || mapping.Phoneme == nil {
		return false
	}

	// In alphabetic systems, each grapheme should typically map to one phoneme
	// and graphemes should be single characters
	for _, grapheme := range mapping.Graphemes {
		if len(grapheme.Symbol) != 1 {
			return false
		}
	}

	return true
}

// SyllabicStyle implements a syllabic writing system where graphemes represent
// syllables rather than individual phonemes.
type SyllabicStyle struct {
	// Syllable patterns to generate graphemes for
	syllablePatterns []string
	// Weight for syllable grapheme selection
	syllableWeight float32
}

// NewSyllabicStyle creates a new syllabic writing system style with common
// syllable patterns.
func NewSyllabicStyle() *SyllabicStyle {
	return &SyllabicStyle{
		syllablePatterns: []string{
			"ba", "be", "bi", "bo", "bu",
			"ka", "ke", "ki", "ko", "ku",
			"ma", "me", "mi", "mo", "mu",
			"na", "ne", "ni", "no", "nu",
			"pa", "pe", "pi", "po", "pu",
			"sa", "se", "si", "so", "su",
			"ta", "te", "ti", "to", "tu",
		},
		syllableWeight: 1.0,
	}
}

// GenerateGraphemes creates syllabic graphemes for the given phoneme pool.
func (ss *SyllabicStyle) GenerateGraphemes(pool *phoneme.PhonemePool, rng *rand.Rand) []Grapheme {
	var graphemes []Grapheme

	// For syllabic systems, we create graphemes that represent common syllable patterns
	// rather than individual phonemes
	for i, pattern := range ss.syllablePatterns {
		if i < len(ss.syllablePatterns) {
			grapheme := Grapheme{
				Symbol: pattern,
				Type:   GraphemeTypeSyllable,
				Weight: ss.syllableWeight,
			}
			graphemes = append(graphemes, grapheme)
		}
	}

	return graphemes
}

// GetStyle returns the syllabic writing style.
func (ss *SyllabicStyle) GetStyle() WritingStyle {
	return WritingStyleSyllabic
}

// ValidateMapping checks if a mapping is valid for syllabic systems.
func (ss *SyllabicStyle) ValidateMapping(mapping *OrthographyMapping) bool {
	if mapping == nil || mapping.Phoneme == nil {
		return false
	}

	// In syllabic systems, graphemes should represent syllables (multiple characters)
	for _, grapheme := range mapping.Graphemes {
		if len(grapheme.Symbol) < 2 {
			return false
		}
	}

	return true
}

// LogographicStyle implements a logographic writing system where graphemes
// represent words or morphemes.
type LogographicStyle struct {
	// Logogram symbols to use
	logogramSymbols []string
	// Weight for logogram selection
	logogramWeight float32
}

// NewLogographicStyle creates a new logographic writing system style with
// basic logogram symbols.
func NewLogographicStyle() *LogographicStyle {
	return &LogographicStyle{
		logogramSymbols: []string{
			"日", "月", "山", "水", "火", "土", "人", "口", "手", "足",
			"木", "金", "女", "男", "子", "大", "小", "上", "下", "中",
		},
		logogramWeight: 1.0,
	}
}

// GenerateGraphemes creates logographic graphemes for the given phoneme pool.
func (ls *LogographicStyle) GenerateGraphemes(pool *phoneme.PhonemePool, rng *rand.Rand) []Grapheme {
	var graphemes []Grapheme

	// For logographic systems, we create graphemes that represent concepts
	// rather than individual phonemes or syllables
	for i, symbol := range ls.logogramSymbols {
		if i < len(ls.logogramSymbols) {
			grapheme := Grapheme{
				Symbol: symbol,
				Type:   GraphemeTypeLogogram,
				Weight: ls.logogramWeight,
			}
			graphemes = append(graphemes, grapheme)
		}
	}

	return graphemes
}

// GetStyle returns the logographic writing style.
func (ls *LogographicStyle) GetStyle() WritingStyle {
	return WritingStyleLogographic
}

// ValidateMapping checks if a mapping is valid for logographic systems.
func (ls *LogographicStyle) ValidateMapping(mapping *OrthographyMapping) bool {
	if mapping == nil || mapping.Phoneme == nil {
		return false
	}

	// In logographic systems, graphemes should represent concepts
	// and may not have direct phoneme mappings
	return true
}
