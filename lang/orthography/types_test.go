package orthography

import (
	"encoding/json"
	"math/rand/v2"
	"testing"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
)

func TestGraphemeType_String(t *testing.T) {
	tests := []struct {
		name     string
		gType    GraphemeType
		expected string
	}{
		{"unknown", GraphemeTypeUnknown, "unknown"},
		{"letter", GraphemeTypeLetter, "letter"},
		{"syllable", GraphemeTypeSyllable, "syllable"},
		{"logogram", GraphemeTypeLogogram, "logogram"},
		{"invalid", GraphemeType(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gType.String(); got != tt.expected {
				t.Errorf("GraphemeType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestWritingStyle_String(t *testing.T) {
	tests := []struct {
		name     string
		wStyle   WritingStyle
		expected string
	}{
		{"unknown", WritingStyleUnknown, "unknown"},
		{"alphabetic", WritingStyleAlphabetic, "alphabetic"},
		{"syllabic", WritingStyleSyllabic, "syllabic"},
		{"logographic", WritingStyleLogographic, "logographic"},
		{"invalid", WritingStyle(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.wStyle.String(); got != tt.expected {
				t.Errorf("WritingStyle.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGrapheme_JSON(t *testing.T) {
	phoneme := &phoneme.Phoneme{
		Symbol: "p",
		Type:   phoneme.PhonemeTypeConsonant,
		Weight: 1.0,
	}

	grapheme := Grapheme{
		Symbol:     "p",
		Type:       GraphemeTypeLetter,
		PhonemeRef: phoneme,
		Weight:     1.0,
	}

	// Test marshaling
	data, err := json.Marshal(grapheme)
	if err != nil {
		t.Fatalf("Failed to marshal grapheme: %v", err)
	}

	// Test unmarshaling
	var unmarshaled Grapheme
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal grapheme: %v", err)
	}

	if unmarshaled.Symbol != grapheme.Symbol {
		t.Errorf("Symbol mismatch: got %s, want %s", unmarshaled.Symbol, grapheme.Symbol)
	}

	if unmarshaled.Type != grapheme.Type {
		t.Errorf("Type mismatch: got %v, want %v", unmarshaled.Type, grapheme.Type)
	}

	if unmarshaled.Weight != grapheme.Weight {
		t.Errorf("Weight mismatch: got %f, want %f", unmarshaled.Weight, grapheme.Weight)
	}
}

func TestGraphemeList_UnweightedChoice(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	graphemes := GraphemeList{
		{Symbol: "a", Type: GraphemeTypeLetter, Weight: 1.0},
		{Symbol: "b", Type: GraphemeTypeLetter, Weight: 2.0},
		{Symbol: "c", Type: GraphemeTypeLetter, Weight: 3.0},
	}

	// Test with populated list
	selected := graphemes.UnweightedChoice(rng)
	if selected == nil {
		t.Error("Expected non-nil grapheme, got nil")
	}

	// Test with empty list
	emptyList := GraphemeList{}
	selected = emptyList.UnweightedChoice(rng)
	if selected != nil {
		t.Errorf("Expected nil for empty list, got %v", selected)
	}
}

func TestGraphemeList_WeightedChoice(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	graphemes := GraphemeList{
		{Symbol: "a", Type: GraphemeTypeLetter, Weight: 1.0},
		{Symbol: "b", Type: GraphemeTypeLetter, Weight: 2.0},
		{Symbol: "c", Type: GraphemeTypeLetter, Weight: 3.0},
	}

	// Test with populated list
	selected := graphemes.WeightedChoice(rng)
	if selected == nil {
		t.Error("Expected non-nil grapheme, got nil")
	}

	// Test with empty list
	emptyList := GraphemeList{}
	selected = emptyList.WeightedChoice(rng)
	if selected != nil {
		t.Errorf("Expected nil for empty list, got %v", selected)
	}
}

func TestGraphemeList_WeightedChoice_ZeroWeight(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	// Test with zero weights (should panic)
	graphemes := GraphemeList{
		{Symbol: "a", Type: GraphemeTypeLetter, Weight: 0.0},
		{Symbol: "b", Type: GraphemeTypeLetter, Weight: 0.0},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for zero total weight, but none occurred")
		}
	}()

	graphemes.WeightedChoice(rng)
}

func TestOrthographyMapping_JSON(t *testing.T) {
	phoneme := &phoneme.Phoneme{
		Symbol: "p",
		Type:   phoneme.PhonemeTypeConsonant,
		Weight: 1.0,
	}

	mapping := OrthographyMapping{
		Phoneme: phoneme,
		Graphemes: []Grapheme{
			{Symbol: "p", Type: GraphemeTypeLetter, Weight: 1.0},
		},
		Primary: true,
	}

	// Test marshaling
	data, err := json.Marshal(mapping)
	if err != nil {
		t.Fatalf("Failed to marshal mapping: %v", err)
	}

	// Test unmarshaling
	var unmarshaled OrthographyMapping
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal mapping: %v", err)
	}

	if unmarshaled.Primary != mapping.Primary {
		t.Errorf("Primary mismatch: got %v, want %v", unmarshaled.Primary, mapping.Primary)
	}

	if len(unmarshaled.Graphemes) != len(mapping.Graphemes) {
		t.Errorf("Graphemes count mismatch: got %d, want %d",
			len(unmarshaled.Graphemes), len(mapping.Graphemes))
	}
}

func TestWritingSystem_JSON(t *testing.T) {
	system := WritingSystem{
		Style: WritingStyleAlphabetic,
		Graphemes: []Grapheme{
			{Symbol: "a", Type: GraphemeTypeLetter, Weight: 1.0},
			{Symbol: "b", Type: GraphemeTypeLetter, Weight: 1.0},
		},
		Mappings: []OrthographyMapping{},
		Name:     "Test System",
		Culture:  "Test Culture",
	}

	// Test marshaling
	data, err := json.Marshal(system)
	if err != nil {
		t.Fatalf("Failed to marshal writing system: %v", err)
	}

	// Test unmarshaling
	var unmarshaled WritingSystem
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal writing system: %v", err)
	}

	if unmarshaled.Name != system.Name {
		t.Errorf("Name mismatch: got %s, want %s", unmarshaled.Name, system.Name)
	}

	if unmarshaled.Culture != system.Culture {
		t.Errorf("Culture mismatch: got %s, want %s", unmarshaled.Culture, system.Culture)
	}

	if unmarshaled.Style != system.Style {
		t.Errorf("Style mismatch: got %v, want %v", unmarshaled.Style, system.Style)
	}
}
