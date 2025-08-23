package orthography

import (
	"testing"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
)

func TestNewOrthographyMapper(t *testing.T) {
	mapper := NewOrthographyMapper()

	if mapper == nil {
		t.Fatal("Expected non-nil mapper")
	}

	if mapper.Count() != 0 {
		t.Errorf("Expected empty mapper, got count %d", mapper.Count())
	}
}

func TestOrthographyMapper_AddMapping(t *testing.T) {
	mapper := NewOrthographyMapper()

	// Create a simple phoneme using the common pool
	pool := phoneme.CommonPool
	if len(pool.Consonants) == 0 {
		t.Fatal("Common pool should have consonants")
	}

	phoneme := &pool.Consonants[0]

	mapping := &OrthographyMapping{
		Phoneme: phoneme,
		Graphemes: []Grapheme{
			{Symbol: "p", Type: GraphemeTypeLetter, Weight: 1.0},
		},
		Primary: true,
	}

	// Test successful addition
	err := mapper.AddMapping(mapping)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mapper.Count() != 1 {
		t.Errorf("Expected count 1, got %d", mapper.Count())
	}

	// Test nil mapping
	err = mapper.AddMapping(nil)
	if err == nil {
		t.Error("Expected error for nil mapping")
	}

	// Test nil phoneme
	invalidMapping := &OrthographyMapping{
		Phoneme: nil,
		Graphemes: []Grapheme{
			{Symbol: "p", Type: GraphemeTypeLetter, Weight: 1.0},
		},
		Primary: true,
	}

	err = mapper.AddMapping(invalidMapping)
	if err == nil {
		t.Error("Expected error for nil phoneme")
	}

	// Test empty phoneme symbol - skip this test for now to avoid complex phoneme creation
	// TODO: Implement proper test for empty phoneme symbol validation
}

func TestOrthographyMapper_GetMapping(t *testing.T) {
	mapper := NewOrthographyMapper()

	// Use phoneme from common pool
	pool := phoneme.CommonPool
	phoneme := &pool.Consonants[0]

	mapping := &OrthographyMapping{
		Phoneme: phoneme,
		Graphemes: []Grapheme{
			{Symbol: "p", Type: GraphemeTypeLetter, Weight: 1.0},
		},
		Primary: true,
	}

	// Add mapping
	err := mapper.AddMapping(mapping)
	if err != nil {
		t.Fatalf("Failed to add mapping: %v", err)
	}

	// Test retrieval
	retrieved := mapper.GetMapping(phoneme.Symbol)
	if retrieved == nil {
		t.Fatal("Expected non-nil mapping")
	}

	if retrieved.Phoneme.Symbol != phoneme.Symbol {
		t.Errorf("Expected phoneme symbol '%s', got %s", phoneme.Symbol, retrieved.Phoneme.Symbol)
	}

	// Test non-existent mapping
	nonExistent := mapper.GetMapping("q")
	if nonExistent != nil {
		t.Errorf("Expected nil for non-existent mapping, got %v", nonExistent)
	}
}

func TestOrthographyMapper_GetGraphemesForPhoneme(t *testing.T) {
	mapper := NewOrthographyMapper()

	// Use phoneme from common pool
	pool := phoneme.CommonPool
	phoneme := &pool.Consonants[0]

	mapping := &OrthographyMapping{
		Phoneme: phoneme,
		Graphemes: []Grapheme{
			{Symbol: "p", Type: GraphemeTypeLetter, Weight: 1.0},
			{Symbol: "pp", Type: GraphemeTypeLetter, Weight: 0.5},
		},
		Primary: true,
	}

	// Add mapping
	err := mapper.AddMapping(mapping)
	if err != nil {
		t.Fatalf("Failed to add mapping: %v", err)
	}

	// Test retrieval
	graphemes := mapper.GetGraphemesForPhoneme(phoneme.Symbol)
	if len(graphemes) != 2 {
		t.Errorf("Expected 2 graphemes, got %d", len(graphemes))
	}

	// Test non-existent phoneme
	nonExistentGraphemes := mapper.GetGraphemesForPhoneme("q")
	if len(nonExistentGraphemes) != 0 {
		t.Errorf("Expected empty slice for non-existent phoneme, got %d", len(nonExistentGraphemes))
	}
}

func TestOrthographyMapper_HasMapping(t *testing.T) {
	mapper := NewOrthographyMapper()

	// Use phoneme from common pool
	pool := phoneme.CommonPool
	phoneme := &pool.Consonants[0]

	mapping := &OrthographyMapping{
		Phoneme: phoneme,
		Graphemes: []Grapheme{
			{Symbol: "p", Type: GraphemeTypeLetter, Weight: 1.0},
		},
		Primary: true,
	}

	// Test before adding
	if mapper.HasMapping(phoneme.Symbol) {
		t.Error("Expected no mapping before addition")
	}

	// Add mapping
	err := mapper.AddMapping(mapping)
	if err != nil {
		t.Fatalf("Failed to add mapping: %v", err)
	}

	// Test after adding
	if !mapper.HasMapping(phoneme.Symbol) {
		t.Error("Expected mapping after addition")
	}

	// Test non-existent
	if mapper.HasMapping("q") {
		t.Error("Expected no mapping for non-existent phoneme")
	}
}

func TestOrthographyMapper_RemoveMapping(t *testing.T) {
	mapper := NewOrthographyMapper()

	// Use phoneme from common pool
	pool := phoneme.CommonPool
	phoneme := &pool.Consonants[0]

	mapping := &OrthographyMapping{
		Phoneme: phoneme,
		Graphemes: []Grapheme{
			{Symbol: "p", Type: GraphemeTypeLetter, Weight: 1.0},
		},
		Primary: true,
	}

	// Add mapping
	err := mapper.AddMapping(mapping)
	if err != nil {
		t.Fatalf("Failed to add mapping: %v", err)
	}

	// Test removal of existing mapping
	removed := mapper.RemoveMapping(phoneme.Symbol)
	if !removed {
		t.Error("Expected true for successful removal")
	}

	if mapper.Count() != 0 {
		t.Errorf("Expected count 0 after removal, got %d", mapper.Count())
	}

	// Test removal of non-existent mapping
	removed = mapper.RemoveMapping("q")
	if removed {
		t.Error("Expected false for non-existent mapping removal")
	}
}

func TestOrthographyMapper_GetAllMappings(t *testing.T) {
	mapper := NewOrthographyMapper()

	// Use phonemes from common pool
	pool := phoneme.CommonPool
	if len(pool.Consonants) < 2 {
		t.Fatal("Common pool should have at least 2 consonants")
	}

	phoneme1 := &pool.Consonants[0]
	phoneme2 := &pool.Consonants[1]

	mapping1 := &OrthographyMapping{
		Phoneme: phoneme1,
		Graphemes: []Grapheme{
			{Symbol: "p", Type: GraphemeTypeLetter, Weight: 1.0},
		},
		Primary: true,
	}

	mapping2 := &OrthographyMapping{
		Phoneme: phoneme2,
		Graphemes: []Grapheme{
			{Symbol: "b", Type: GraphemeTypeLetter, Weight: 1.0},
		},
		Primary: true,
	}

	// Add mappings
	err := mapper.AddMapping(mapping1)
	if err != nil {
		t.Fatalf("Failed to add mapping1: %v", err)
	}

	err = mapper.AddMapping(mapping2)
	if err != nil {
		t.Fatalf("Failed to add mapping2: %v", err)
	}

	// Test retrieval
	allMappings := mapper.GetAllMappings()
	if len(allMappings) != 2 {
		t.Errorf("Expected 2 mappings, got %d", len(allMappings))
	}
}

func TestOrthographyMapper_IsComplete(t *testing.T) {
	mapper := NewOrthographyMapper()

	// Create a small phoneme pool for testing
	pool := phoneme.NewPool()

	// Add some phonemes from common pool
	if len(phoneme.CommonPool.Consonants) > 0 {
		pool.AddConsonant(phoneme.CommonPool.Consonants[0])
	}
	if len(phoneme.CommonPool.Vowels) > 0 {
		pool.AddVowel(phoneme.CommonPool.Vowels[0])
	}

	// Test incomplete mapping
	if mapper.IsComplete(pool) {
		t.Error("Expected incomplete mapping before adding mappings")
	}

	// Add mappings
	if len(pool.Consonants) > 0 {
		mapping1 := &OrthographyMapping{
			Phoneme: &pool.Consonants[0],
			Graphemes: []Grapheme{
				{Symbol: "p", Type: GraphemeTypeLetter, Weight: 1.0},
			},
			Primary: true,
		}

		err := mapper.AddMapping(mapping1)
		if err != nil {
			t.Fatalf("Failed to add mapping1: %v", err)
		}
	}

	if len(pool.Vowels) > 0 {
		mapping2 := &OrthographyMapping{
			Phoneme: &pool.Vowels[0],
			Graphemes: []Grapheme{
				{Symbol: "a", Type: GraphemeTypeLetter, Weight: 1.0},
			},
			Primary: true,
		}

		err := mapper.AddMapping(mapping2)
		if err != nil {
			t.Fatalf("Failed to add mapping2: %v", err)
		}
	}

	// Test complete mapping
	if !mapper.IsComplete(pool) {
		t.Error("Expected complete mapping after adding all mappings")
	}

	// Test with nil pool
	if mapper.IsComplete(nil) {
		t.Error("Expected false for nil pool")
	}
}
