package orthography

import (
	"math/rand/v2"
	"testing"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
)

func TestNewOrthographyGenerator(t *testing.T) {
	generator := NewOrthographyGenerator()

	if generator == nil {
		t.Fatal("Expected non-nil generator")
	}

	// Check that default styles are registered
	alphabeticStyle, err := generator.GetStyle(WritingStyleAlphabetic)
	if err != nil {
		t.Errorf("Expected alphabetic style to be available: %v", err)
	}
	if alphabeticStyle == nil {
		t.Error("Expected non-nil alphabetic style")
	}

	syllabicStyle, err := generator.GetStyle(WritingStyleSyllabic)
	if err != nil {
		t.Errorf("Expected syllabic style to be available: %v", err)
	}
	if syllabicStyle == nil {
		t.Error("Expected non-nil syllabic style")
	}

	logographicStyle, err := generator.GetStyle(WritingStyleLogographic)
	if err != nil {
		t.Errorf("Expected logographic style to be available: %v", err)
	}
	if logographicStyle == nil {
		t.Error("Expected non-nil logographic style")
	}
}

func TestOrthographyGenerator_RegisterStyle(t *testing.T) {
	generator := NewOrthographyGenerator()

	// Create a custom style
	customStyle := &AlphabeticStyle{
		consonantChars:  []string{"x", "y", "z"},
		vowelChars:      []string{"q"},
		consonantWeight: 2.0,
		vowelWeight:     2.0,
	}

	// Register the custom style
	generator.RegisterStyle(customStyle)

	// Verify it's available
	retrievedStyle, err := generator.GetStyle(WritingStyleAlphabetic)
	if err != nil {
		t.Errorf("Expected custom style to be available: %v", err)
	}

	if retrievedStyle == nil {
		t.Error("Expected non-nil custom style")
	}
}

func TestOrthographyGenerator_GetStyle(t *testing.T) {
	generator := NewOrthographyGenerator()

	// Test existing style
	style, err := generator.GetStyle(WritingStyleAlphabetic)
	if err != nil {
		t.Errorf("Expected no error for existing style: %v", err)
	}
	if style == nil {
		t.Error("Expected non-nil style")
	}

	// Test non-existent style
	_, err = generator.GetStyle(WritingStyle(99))
	if err == nil {
		t.Error("Expected error for non-existent style")
	}
}

func TestOrthographyGenerator_GenerateWritingSystem(t *testing.T) {
	generator := NewOrthographyGenerator()
	pool := &phoneme.CommonPool
	rng := rand.New(rand.NewPCG(1, 2))

	// Test alphabetic style
	system, err := generator.GenerateWritingSystem(
		pool,
		WritingStyleAlphabetic,
		"Test Alphabetic",
		"Test Culture",
		rng,
	)

	if err != nil {
		t.Fatalf("Expected no error: %v", err)
	}

	if system == nil {
		t.Fatal("Expected non-nil writing system")
	}

	if system.Name != "Test Alphabetic" {
		t.Errorf("Expected name 'Test Alphabetic', got %s", system.Name)
	}

	if system.Culture != "Test Culture" {
		t.Errorf("Expected culture 'Test Culture', got %s", system.Culture)
	}

	if system.Style != WritingStyleAlphabetic {
		t.Errorf("Expected alphabetic style, got %v", system.Style)
	}

	if len(system.Graphemes) == 0 {
		t.Error("Expected non-empty graphemes")
	}

	if len(system.Mappings) == 0 {
		t.Error("Expected non-empty mappings")
	}

	// Test with nil pool
	_, err = generator.GenerateWritingSystem(
		nil,
		WritingStyleAlphabetic,
		"Test",
		"Test",
		rng,
	)

	if err == nil {
		t.Error("Expected error for nil pool")
	}
}

func TestOrthographyGenerator_GenerateRandomWritingSystem(t *testing.T) {
	generator := NewOrthographyGenerator()
	pool := &phoneme.CommonPool
	rng := rand.New(rand.NewPCG(1, 2))

	system, err := generator.GenerateRandomWritingSystem(
		pool,
		"Test Random",
		"Test Culture",
		rng,
	)

	if err != nil {
		t.Fatalf("Expected no error: %v", err)
	}

	if system == nil {
		t.Fatal("Expected non-nil writing system")
	}

	if system.Name != "Test Random" {
		t.Errorf("Expected name 'Test Random', got %s", system.Name)
	}

	if system.Culture != "Test Culture" {
		t.Errorf("Expected culture 'Test Culture', got %s", system.Culture)
	}

	// Verify it has a valid style
	if system.Style == WritingStyleUnknown {
		t.Error("Expected non-unknown writing style")
	}
}

func TestOrthographyGenerator_ValidateWritingSystem(t *testing.T) {
	generator := NewOrthographyGenerator()

	// Test valid system
	validSystem := &WritingSystem{
		Style: WritingStyleAlphabetic,
		Graphemes: []Grapheme{
			{Symbol: "a", Type: GraphemeTypeLetter, Weight: 1.0},
		},
		Mappings: []OrthographyMapping{
			{
				Phoneme: &phoneme.CommonPool.Consonants[0],
				Graphemes: []Grapheme{
					{Symbol: "a", Type: GraphemeTypeLetter, Weight: 1.0},
				},
				Primary: true,
			},
		},
		Name:    "Valid System",
		Culture: "Test Culture",
	}

	err := generator.ValidateWritingSystem(validSystem)
	if err != nil {
		t.Errorf("Expected no error for valid system: %v", err)
	}

	// Test nil system
	err = generator.ValidateWritingSystem(nil)
	if err == nil {
		t.Error("Expected error for nil system")
	}

	// Test system without name
	noNameSystem := &WritingSystem{
		Style: WritingStyleAlphabetic,
		Graphemes: []Grapheme{
			{Symbol: "a", Type: GraphemeTypeLetter, Weight: 1.0},
		},
		Mappings: []OrthographyMapping{},
		Name:     "",
		Culture:  "Test Culture",
	}

	err = generator.ValidateWritingSystem(noNameSystem)
	if err == nil {
		t.Error("Expected error for system without name")
	}

	// Test system without graphemes
	noGraphemesSystem := &WritingSystem{
		Style:     WritingStyleAlphabetic,
		Graphemes: []Grapheme{},
		Mappings:  []OrthographyMapping{},
		Name:      "No Graphemes",
		Culture:   "Test Culture",
	}

	err = generator.ValidateWritingSystem(noGraphemesSystem)
	if err == nil {
		t.Error("Expected error for system without graphemes")
	}

	// Test system without mappings
	noMappingsSystem := &WritingSystem{
		Style: WritingStyleAlphabetic,
		Graphemes: []Grapheme{
			{Symbol: "a", Type: GraphemeTypeLetter, Weight: 1.0},
		},
		Mappings: []OrthographyMapping{},
		Name:     "No Mappings",
		Culture:  "Test Culture",
	}

	err = generator.ValidateWritingSystem(noMappingsSystem)
	if err == nil {
		t.Error("Expected error for system without mappings")
	}
}

func TestOrthographyGenerator_GetWritingSystemStats(t *testing.T) {
	generator := NewOrthographyGenerator()

	// Test with valid system
	system := &WritingSystem{
		Style: WritingStyleAlphabetic,
		Graphemes: []Grapheme{
			{Symbol: "a", Type: GraphemeTypeLetter, Weight: 1.0},
			{Symbol: "b", Type: GraphemeTypeLetter, Weight: 1.0},
		},
		Mappings: []OrthographyMapping{
			{
				Phoneme: &phoneme.CommonPool.Consonants[0],
				Graphemes: []Grapheme{
					{Symbol: "a", Type: GraphemeTypeLetter, Weight: 1.0},
				},
				Primary: true,
			},
		},
		Name:    "Test System",
		Culture: "Test Culture",
	}

	stats := generator.GetWritingSystemStats(system)

	if stats == nil {
		t.Fatal("Expected non-nil stats")
	}

	if stats["name"] != "Test System" {
		t.Errorf("Expected name 'Test System', got %v", stats["name"])
	}

	if stats["style"] != "alphabetic" {
		t.Errorf("Expected style 'alphabetic', got %v", stats["style"])
	}

	if stats["culture"] != "Test Culture" {
		t.Errorf("Expected culture 'Test Culture', got %v", stats["culture"])
	}

	if stats["graphemeCount"] != 2 {
		t.Errorf("Expected 2 graphemes, got %v", stats["graphemeCount"])
	}

	if stats["mappingCount"] != 1 {
		t.Errorf("Expected 1 mapping, got %v", stats["mappingCount"])
	}

	// Test with nil system
	stats = generator.GetWritingSystemStats(nil)
	if stats != nil {
		t.Error("Expected nil stats for nil system")
	}
}

func TestOrthographyGenerator_Integration(t *testing.T) {
	generator := NewOrthographyGenerator()
	pool := &phoneme.CommonPool
	rng := rand.New(rand.NewPCG(1, 2))

	// Test all writing system styles
	styles := []WritingStyle{
		WritingStyleAlphabetic,
		WritingStyleSyllabic,
		WritingStyleLogographic,
	}

	for _, style := range styles {
		t.Run(style.String(), func(t *testing.T) {
			system, err := generator.GenerateWritingSystem(
				pool,
				style,
				"Test "+style.String(),
				"Test Culture",
				rng,
			)

			if err != nil {
				t.Fatalf("Failed to generate %s system: %v", style.String(), err)
			}

			if system == nil {
				t.Fatal("Expected non-nil writing system")
			}

			// Validate the generated system
			err = generator.ValidateWritingSystem(system)
			if err != nil {
				t.Errorf("Generated system failed validation: %v", err)
			}

			// Check stats
			stats := generator.GetWritingSystemStats(system)
			if stats == nil {
				t.Error("Expected non-nil stats")
			}
		})
	}
}
