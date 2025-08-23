package phoneme

import (
	"encoding/json"
	"testing"
)

func TestNewPool(t *testing.T) {
	pool := NewPool()

	if pool == nil {
		t.Fatal("NewPool() returned nil")
	}

	if len(pool.Consonants) != 0 {
		t.Errorf("Expected empty consonants list, got %d", len(pool.Consonants))
	}

	if len(pool.Vowels) != 0 {
		t.Errorf("Expected empty vowels list, got %d", len(pool.Vowels))
	}
}

func TestPhonemePool_AddConsonant(t *testing.T) {
	pool := NewPool()

	// Test adding valid consonant
	consonant := Phoneme{
		Symbol: "p",
		Type:   PhonemeTypeConsonant,
		Weight: 1.0,
		Consonant: &ConsonantSpec{
			Manner: ConsonantMannerPlosive,
		},
	}

	pool.AddConsonant(consonant)

	if len(pool.Consonants) != 1 {
		t.Errorf("Expected 1 consonant, got %d", len(pool.Consonants))
	}

	if pool.Consonants[0].Symbol != "p" {
		t.Errorf("Expected consonant symbol 'p', got %s", pool.Consonants[0].Symbol)
	}

	// Test adding vowel as consonant (should be ignored)
	vowel := Phoneme{
		Symbol: "a",
		Type:   PhonemeTypeVowel,
		Weight: 1.0,
		Vowel: &VowelSpec{
			Height:   VowelHeightLow,
			Backness: VowelBacknessFront,
			Rounded:  false,
		},
	}

	pool.AddConsonant(vowel)

	if len(pool.Consonants) != 1 {
		t.Errorf("Expected still 1 consonant after adding vowel, got %d", len(pool.Consonants))
	}
}

func TestPhonemePool_AddVowel(t *testing.T) {
	pool := NewPool()

	// Test adding valid vowel
	vowel := Phoneme{
		Symbol: "a",
		Type:   PhonemeTypeVowel,
		Weight: 1.0,
		Vowel: &VowelSpec{
			Height:   VowelHeightLow,
			Backness: VowelBacknessFront,
			Rounded:  false,
		},
	}

	pool.AddVowel(vowel)

	if len(pool.Vowels) != 1 {
		t.Errorf("Expected 1 vowel, got %d", len(pool.Vowels))
	}

	if pool.Vowels[0].Symbol != "a" {
		t.Errorf("Expected vowel symbol 'a', got %s", pool.Vowels[0].Symbol)
	}

	// Test adding consonant as vowel (should be ignored)
	consonant := Phoneme{
		Symbol: "p",
		Type:   PhonemeTypeConsonant,
		Weight: 1.0,
		Consonant: &ConsonantSpec{
			Manner: ConsonantMannerPlosive,
		},
	}

	pool.AddVowel(consonant)

	if len(pool.Vowels) != 1 {
		t.Errorf("Expected still 1 vowel after adding consonant, got %d", len(pool.Vowels))
	}
}

func TestPhonemePool_DuplicatePrevention(t *testing.T) {
	pool := NewPool()

	consonant1 := Phoneme{
		Symbol: "p",
		Type:   PhonemeTypeConsonant,
		Weight: 1.0,
		Consonant: &ConsonantSpec{
			Manner: ConsonantMannerPlosive,
		},
	}

	consonant2 := Phoneme{
		Symbol: "p", // Same symbol
		Type:   PhonemeTypeConsonant,
		Weight: 2.0, // Different weight
		Consonant: &ConsonantSpec{
			Manner: ConsonantMannerFricative, // Different manner
		},
	}

	pool.AddConsonant(consonant1)
	pool.AddConsonant(consonant2)

	if len(pool.Consonants) != 1 {
		t.Errorf("Expected 1 consonant after adding duplicate, got %d", len(pool.Consonants))
	}

	// Should keep the first one
	if pool.Consonants[0].Weight != 1.0 {
		t.Errorf("Expected weight 1.0, got %f", pool.Consonants[0].Weight)
	}
	if pool.Consonants[0].Consonant.Manner != ConsonantMannerPlosive {
		t.Errorf("Expected manner plosive, got %v", pool.Consonants[0].Consonant.Manner)
	}
}

func TestPhonemePool_GetConsonantBySymbol(t *testing.T) {
	pool := NewPool()

	consonant := Phoneme{
		Symbol: "p",
		Type:   PhonemeTypeConsonant,
		Weight: 1.0,
		Consonant: &ConsonantSpec{
			Manner: ConsonantMannerPlosive,
		},
	}

	pool.AddConsonant(consonant)

	// Test successful retrieval
	retrieved := pool.GetConsonantBySymbol("p")
	if retrieved == nil {
		t.Fatal("Expected non-nil consonant")
	}

	if retrieved.Symbol != "p" {
		t.Errorf("Expected symbol 'p', got %s", retrieved.Symbol)
	}

	// Test retrieval of non-existent consonant
	retrieved = pool.GetConsonantBySymbol("q")
	if retrieved != nil {
		t.Error("Expected nil for non-existent consonant")
	}
}

func TestPhonemePool_GetVowelBySymbol(t *testing.T) {
	pool := NewPool()

	vowel := Phoneme{
		Symbol: "a",
		Type:   PhonemeTypeVowel,
		Weight: 1.0,
		Vowel: &VowelSpec{
			Height:   VowelHeightLow,
			Backness: VowelBacknessFront,
			Rounded:  false,
		},
	}

	pool.AddVowel(vowel)

	// Test successful retrieval
	retrieved := pool.GetVowelBySymbol("a")
	if retrieved == nil {
		t.Fatal("Expected non-nil vowel")
	}

	if retrieved.Symbol != "a" {
		t.Errorf("Expected symbol 'a', got %s", retrieved.Symbol)
	}

	// Test retrieval of non-existent vowel
	retrieved = pool.GetVowelBySymbol("e")
	if retrieved != nil {
		t.Error("Expected nil for non-existent vowel")
	}
}

func TestPhonemePool_Size(t *testing.T) {
	pool := NewPool()

	// Test empty pool
	if pool.Size() != 0 {
		t.Errorf("Expected size 0 for empty pool, got %d", pool.Size())
	}

	// Test with consonants only
	consonant := Phoneme{
		Symbol: "p",
		Type:   PhonemeTypeConsonant,
		Weight: 1.0,
		Consonant: &ConsonantSpec{
			Manner: ConsonantMannerPlosive,
		},
	}
	pool.AddConsonant(consonant)

	if pool.Size() != 1 {
		t.Errorf("Expected size 1 after adding consonant, got %d", pool.Size())
	}

	// Test with vowels only
	vowel := Phoneme{
		Symbol: "a",
		Type:   PhonemeTypeVowel,
		Weight: 1.0,
		Vowel: &VowelSpec{
			Height:   VowelHeightLow,
			Backness: VowelBacknessFront,
			Rounded:  false,
		},
	}
	pool.AddVowel(vowel)

	if pool.Size() != 2 {
		t.Errorf("Expected size 2 after adding vowel, got %d", pool.Size())
	}
}

func TestPhonemePool_JSON(t *testing.T) {
	pool := NewPool()

	// Add some phonemes
	consonant := Phoneme{
		Symbol: "p",
		Type:   PhonemeTypeConsonant,
		Weight: 1.0,
		Consonant: &ConsonantSpec{
			Manner: ConsonantMannerPlosive,
		},
	}

	vowel := Phoneme{
		Symbol: "a",
		Type:   PhonemeTypeVowel,
		Weight: 1.0,
		Vowel: &VowelSpec{
			Height:   VowelHeightLow,
			Backness: VowelBacknessFront,
			Rounded:  false,
		},
	}

	pool.AddConsonant(consonant)
	pool.AddVowel(vowel)

	// Test JSON marshaling
	data, err := json.Marshal(pool)
	if err != nil {
		t.Fatalf("Failed to marshal pool: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled PhonemePool
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal pool: %v", err)
	}

	if len(unmarshaled.Consonants) != len(pool.Consonants) {
		t.Errorf("Consonants count mismatch: got %d, want %d",
			len(unmarshaled.Consonants), len(pool.Consonants))
	}

	if len(unmarshaled.Vowels) != len(pool.Vowels) {
		t.Errorf("Vowels count mismatch: got %d, want %d",
			len(unmarshaled.Vowels), len(pool.Vowels))
	}

	// Verify consonant details
	if unmarshaled.Consonants[0].Symbol != pool.Consonants[0].Symbol {
		t.Errorf("Consonant symbol mismatch: got %s, want %s",
			unmarshaled.Consonants[0].Symbol, pool.Consonants[0].Symbol)
	}

	// Verify vowel details
	if unmarshaled.Vowels[0].Symbol != pool.Vowels[0].Symbol {
		t.Errorf("Vowel symbol mismatch: got %s, want %s",
			unmarshaled.Vowels[0].Symbol, pool.Vowels[0].Symbol)
	}
}

func TestPhonemePool_EdgeCases(t *testing.T) {
	pool := NewPool()

	// Test adding phoneme with empty symbol
	emptySymbolPhoneme := Phoneme{
		Symbol: "",
		Type:   PhonemeTypeConsonant,
		Weight: 1.0,
		Consonant: &ConsonantSpec{
			Manner: ConsonantMannerPlosive,
		},
	}

	pool.AddConsonant(emptySymbolPhoneme)

	// Should still be added (no validation on empty symbols)
	if len(pool.Consonants) != 1 {
		t.Errorf("Expected 1 consonant after adding empty symbol, got %d", len(pool.Consonants))
	}

	// Test adding phoneme with zero weight
	zeroWeightPhoneme := Phoneme{
		Symbol: "z",
		Type:   PhonemeTypeConsonant,
		Weight: 0.0,
		Consonant: &ConsonantSpec{
			Manner: ConsonantMannerFricative,
		},
	}

	pool.AddConsonant(zeroWeightPhoneme)

	if len(pool.Consonants) != 2 {
		t.Errorf("Expected 2 consonants after adding zero weight, got %d", len(pool.Consonants))
	}

	// Test adding phoneme with negative weight
	negativeWeightPhoneme := Phoneme{
		Symbol: "n",
		Type:   PhonemeTypeConsonant,
		Weight: -1.0,
		Consonant: &ConsonantSpec{
			Manner: ConsonantMannerNasal,
		},
	}

	pool.AddConsonant(negativeWeightPhoneme)

	if len(pool.Consonants) != 3 {
		t.Errorf("Expected 3 consonants after adding negative weight, got %d", len(pool.Consonants))
	}
}

func TestPhonemePool_ComplexOperations(t *testing.T) {
	pool := NewPool()

	// Add multiple phonemes
	consonants := []Phoneme{
		{Symbol: "p", Type: PhonemeTypeConsonant, Weight: 1.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}},
		{Symbol: "b", Type: PhonemeTypeConsonant, Weight: 1.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}},
		{Symbol: "t", Type: PhonemeTypeConsonant, Weight: 1.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}},
		{Symbol: "d", Type: PhonemeTypeConsonant, Weight: 1.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}},
	}

	vowels := []Phoneme{
		{Symbol: "a", Type: PhonemeTypeVowel, Weight: 1.0, Vowel: &VowelSpec{Height: VowelHeightLow, Backness: VowelBacknessFront, Rounded: false}},
		{Symbol: "e", Type: PhonemeTypeVowel, Weight: 1.0, Vowel: &VowelSpec{Height: VowelHeightMid, Backness: VowelBacknessFront, Rounded: false}},
		{Symbol: "i", Type: PhonemeTypeVowel, Weight: 1.0, Vowel: &VowelSpec{Height: VowelHeightHigh, Backness: VowelBacknessFront, Rounded: false}},
	}

	for _, c := range consonants {
		pool.AddConsonant(c)
	}

	for _, v := range vowels {
		pool.AddVowel(v)
	}

	// Verify final state
	if pool.Size() != 7 {
		t.Errorf("Expected final size 7, got %d", pool.Size())
	}

	if len(pool.Consonants) != 4 {
		t.Errorf("Expected 4 consonants, got %d", len(pool.Consonants))
	}

	if len(pool.Vowels) != 3 {
		t.Errorf("Expected 3 vowels, got %d", len(pool.Vowels))
	}

	// Test retrieval of all phonemes
	for _, expected := range consonants {
		retrieved := pool.GetConsonantBySymbol(expected.Symbol)
		if retrieved == nil {
			t.Errorf("Failed to retrieve consonant %s", expected.Symbol)
		}
	}

	for _, expected := range vowels {
		retrieved := pool.GetVowelBySymbol(expected.Symbol)
		if retrieved == nil {
			t.Errorf("Failed to retrieve vowel %s", expected.Symbol)
		}
	}
}
