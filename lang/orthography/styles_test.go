package orthography

import (
	"math/rand/v2"
	"testing"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
)

func TestNewAlphabeticStyle(t *testing.T) {
	style := NewAlphabeticStyle()

	if style == nil {
		t.Fatal("Expected non-nil alphabetic style")
	}

	if style.GetStyle() != WritingStyleAlphabetic {
		t.Errorf("Expected alphabetic style, got %v", style.GetStyle())
	}

	if len(style.consonantChars) == 0 {
		t.Error("Expected consonant characters")
	}

	if len(style.vowelChars) == 0 {
		t.Error("Expected vowel characters")
	}
}

func TestAlphabeticStyle_GenerateGraphemes(t *testing.T) {
	style := NewAlphabeticStyle()
	pool := &phoneme.CommonPool
	rng := rand.New(rand.NewPCG(1, 2))

	graphemes := style.GenerateGraphemes(pool, rng)

	if len(graphemes) == 0 {
		t.Error("Expected non-empty graphemes")
	}

	// Check that we have graphemes for both consonants and vowels
	consonantCount := 0
	vowelCount := 0

	for _, grapheme := range graphemes {
		if grapheme.Type == GraphemeTypeLetter {
			if grapheme.PhonemeRef != nil {
				if grapheme.PhonemeRef.Type == phoneme.PhonemeTypeConsonant {
					consonantCount++
				} else if grapheme.PhonemeRef.Type == phoneme.PhonemeTypeVowel {
					vowelCount++
				}
			}
		}
	}

	if consonantCount == 0 {
		t.Error("Expected consonant graphemes")
	}

	if vowelCount == 0 {
		t.Error("Expected vowel graphemes")
	}
}

func TestAlphabeticStyle_ValidateMapping(t *testing.T) {
	style := NewAlphabeticStyle()

	// Test valid mapping
	validMapping := &OrthographyMapping{
		Phoneme: &phoneme.CommonPool.Consonants[0],
		Graphemes: []Grapheme{
			{Symbol: "p", Type: GraphemeTypeLetter, Weight: 1.0},
		},
		Primary: true,
	}

	if !style.ValidateMapping(validMapping) {
		t.Error("Expected valid mapping to pass validation")
	}

	// Test invalid mapping (multi-character grapheme)
	invalidMapping := &OrthographyMapping{
		Phoneme: &phoneme.CommonPool.Consonants[0],
		Graphemes: []Grapheme{
			{Symbol: "pp", Type: GraphemeTypeLetter, Weight: 1.0},
		},
		Primary: true,
	}

	if style.ValidateMapping(invalidMapping) {
		t.Error("Expected invalid mapping to fail validation")
	}

	// Test nil mapping
	if style.ValidateMapping(nil) {
		t.Error("Expected nil mapping to fail validation")
	}

	// Test nil phoneme
	nilPhonemeMapping := &OrthographyMapping{
		Phoneme: nil,
		Graphemes: []Grapheme{
			{Symbol: "p", Type: GraphemeTypeLetter, Weight: 1.0},
		},
		Primary: true,
	}

	if style.ValidateMapping(nilPhonemeMapping) {
		t.Error("Expected nil phoneme mapping to fail validation")
	}
}

func TestNewSyllabicStyle(t *testing.T) {
	style := NewSyllabicStyle()

	if style == nil {
		t.Fatal("Expected non-nil syllabic style")
	}

	if style.GetStyle() != WritingStyleSyllabic {
		t.Errorf("Expected syllabic style, got %v", style.GetStyle())
	}

	if len(style.syllablePatterns) == 0 {
		t.Error("Expected syllable patterns")
	}
}

func TestSyllabicStyle_GenerateGraphemes(t *testing.T) {
	style := NewSyllabicStyle()
	pool := &phoneme.CommonPool
	rng := rand.New(rand.NewPCG(1, 2))

	graphemes := style.GenerateGraphemes(pool, rng)

	if len(graphemes) == 0 {
		t.Error("Expected non-empty graphemes")
	}

	// Check that all graphemes are syllables
	for _, grapheme := range graphemes {
		if grapheme.Type != GraphemeTypeSyllable {
			t.Errorf("Expected syllable type, got %v", grapheme.Type)
		}

		if len(grapheme.Symbol) < 2 {
			t.Errorf("Expected multi-character syllable, got %s", grapheme.Symbol)
		}
	}
}

func TestSyllabicStyle_ValidateMapping(t *testing.T) {
	style := NewSyllabicStyle()

	// Test valid mapping
	validMapping := &OrthographyMapping{
		Phoneme: &phoneme.CommonPool.Consonants[0],
		Graphemes: []Grapheme{
			{Symbol: "ba", Type: GraphemeTypeSyllable, Weight: 1.0},
		},
		Primary: true,
	}

	if !style.ValidateMapping(validMapping) {
		t.Error("Expected valid mapping to pass validation")
	}

	// Test invalid mapping (single character)
	invalidMapping := &OrthographyMapping{
		Phoneme: &phoneme.CommonPool.Consonants[0],
		Graphemes: []Grapheme{
			{Symbol: "b", Type: GraphemeTypeSyllable, Weight: 1.0},
		},
		Primary: true,
	}

	if style.ValidateMapping(invalidMapping) {
		t.Error("Expected invalid mapping to fail validation")
	}
}

func TestNewLogographicStyle(t *testing.T) {
	style := NewLogographicStyle()

	if style == nil {
		t.Fatal("Expected non-nil logographic style")
	}

	if style.GetStyle() != WritingStyleLogographic {
		t.Errorf("Expected logographic style, got %v", style.GetStyle())
	}

	if len(style.logogramSymbols) == 0 {
		t.Error("Expected logogram symbols")
	}
}

func TestLogographicStyle_GenerateGraphemes(t *testing.T) {
	style := NewLogographicStyle()
	pool := &phoneme.CommonPool
	rng := rand.New(rand.NewPCG(1, 2))

	graphemes := style.GenerateGraphemes(pool, rng)

	if len(graphemes) == 0 {
		t.Error("Expected non-empty graphemes")
	}

	// Check that all graphemes are logograms
	for _, grapheme := range graphemes {
		if grapheme.Type != GraphemeTypeLogogram {
			t.Errorf("Expected logogram type, got %v", grapheme.Type)
		}
	}
}

func TestLogographicStyle_ValidateMapping(t *testing.T) {
	style := NewLogographicStyle()

	// Test valid mapping (logographic systems are permissive)
	validMapping := &OrthographyMapping{
		Phoneme: &phoneme.CommonPool.Consonants[0],
		Graphemes: []Grapheme{
			{Symbol: "日", Type: GraphemeTypeLogogram, Weight: 1.0},
		},
		Primary: true,
	}

	if !style.ValidateMapping(validMapping) {
		t.Error("Expected valid mapping to pass validation")
	}

	// Test nil mapping
	if style.ValidateMapping(nil) {
		t.Error("Expected nil mapping to fail validation")
	}

	// Test nil phoneme
	nilPhonemeMapping := &OrthographyMapping{
		Phoneme: nil,
		Graphemes: []Grapheme{
			{Symbol: "日", Type: GraphemeTypeLogogram, Weight: 1.0},
		},
		Primary: true,
	}

	if style.ValidateMapping(nilPhonemeMapping) {
		t.Error("Expected nil phoneme mapping to fail validation")
	}
}

func TestWritingSystemStyle_Interface(t *testing.T) {
	// Test that all styles implement the interface correctly
	styles := []WritingSystemStyle{
		NewAlphabeticStyle(),
		NewSyllabicStyle(),
		NewLogographicStyle(),
	}

	for _, style := range styles {
		if style == nil {
			t.Error("Expected non-nil style")
			continue
		}

		// Test GetStyle method
		styleType := style.GetStyle()
		if styleType == WritingStyleUnknown {
			t.Error("Expected non-unknown style type")
		}

		// Test that we can call GenerateGraphemes
		pool := &phoneme.CommonPool
		rng := rand.New(rand.NewPCG(1, 2))
		graphemes := style.GenerateGraphemes(pool, rng)

		if len(graphemes) == 0 {
			t.Error("Expected non-empty graphemes from style")
		}
	}
}
