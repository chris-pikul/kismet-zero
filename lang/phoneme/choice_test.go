package phoneme

import (
	"math/rand/v2"
	"testing"
)

func TestSelectFromPool_EmptyPool(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	pool := NewPool()

	// Test with empty pool
	selected, err := SelectFromPool(pool, rng, 5, ChoiceOptions{})
	if err != nil {
		t.Errorf("Expected no error for empty pool, got %v", err)
	}

	if len(selected) != 0 {
		t.Errorf("Expected 0 phonemes from empty pool, got %d", len(selected))
	}
}

func TestSelectFromPool_ZeroCount(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	pool := NewPool()

	// Add some phonemes
	consonant := Phoneme{
		Symbol:    "p",
		Type:      PhonemeTypeConsonant,
		Weight:    1.0,
		Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive},
	}
	pool.AddConsonant(consonant)

	// Test with zero count
	selected, err := SelectFromPool(pool, rng, 0, ChoiceOptions{})
	if err != nil {
		t.Errorf("Expected no error for zero count, got %v", err)
	}

	if len(selected) != 0 {
		t.Errorf("Expected 0 phonemes for zero count, got %d", len(selected))
	}
}

func TestSelectFromPool_UnweightedSelection(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	pool := NewPool()

	// Add phonemes with different weights
	consonants := []Phoneme{
		{Symbol: "p", Type: PhonemeTypeConsonant, Weight: 1.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}},
		{Symbol: "b", Type: PhonemeTypeConsonant, Weight: 10.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}},
		{Symbol: "t", Type: PhonemeTypeConsonant, Weight: 100.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}},
	}

	for _, c := range consonants {
		pool.AddConsonant(c)
	}

	// Test unweighted selection
	options := ChoiceOptions{UseWeighted: false}
	selected, err := SelectFromPool(pool, rng, 2, options)
	if err != nil {
		t.Fatalf("Failed to select phonemes: %v", err)
	}

	if len(selected) != 2 {
		t.Errorf("Expected 2 phonemes, got %d", len(selected))
	}

	// Verify all selected phonemes are from the pool
	for _, phoneme := range selected {
		found := false
		for _, expected := range consonants {
			if phoneme.Symbol == expected.Symbol {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Selected phoneme %s not found in pool", phoneme.Symbol)
		}
	}
}

func TestSelectFromPool_WeightedSelection(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	pool := NewPool()

	// Add phonemes with very different weights
	consonants := []Phoneme{
		{Symbol: "rare", Type: PhonemeTypeConsonant, Weight: 1.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}},
		{Symbol: "common", Type: PhonemeTypeConsonant, Weight: 10.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}},
		{Symbol: "very_common", Type: PhonemeTypeConsonant, Weight: 100.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}},
	}

	for _, c := range consonants {
		pool.AddConsonant(c)
	}

	// Test weighted selection
	options := ChoiceOptions{UseWeighted: true}
	selected, err := SelectFromPool(pool, rng, 2, options)
	if err != nil {
		t.Fatalf("Failed to select phonemes: %v", err)
	}

	if len(selected) != 2 {
		t.Errorf("Expected 2 phonemes, got %d", len(selected))
	}

	// Verify all selected phonemes are from the pool
	for _, phoneme := range selected {
		found := false
		for _, expected := range consonants {
			if phoneme.Symbol == expected.Symbol {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Selected phoneme %s not found in pool", phoneme.Symbol)
		}
	}
}

func TestSelectFromPool_FilterByType(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	pool := NewPool()

	// Add both consonants and vowels
	consonant := Phoneme{
		Symbol:    "p",
		Type:      PhonemeTypeConsonant,
		Weight:    1.0,
		Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive},
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

	// Test filtering by consonant type
	consonantType := PhonemeTypeConsonant
	options := ChoiceOptions{
		UseWeighted:  false,
		FilterByType: &consonantType,
	}

	selected, err := SelectFromPool(pool, rng, 1, options)
	if err != nil {
		t.Fatalf("Failed to select consonants: %v", err)
	}

	if len(selected) != 1 {
		t.Errorf("Expected 1 consonant, got %d", len(selected))
	}

	if selected[0].Type != PhonemeTypeConsonant {
		t.Errorf("Expected consonant type, got %v", selected[0].Type)
	}

	// Test filtering by vowel type
	vowelType := PhonemeTypeVowel
	options = ChoiceOptions{
		UseWeighted:  false,
		FilterByType: &vowelType,
	}

	selected, err = SelectFromPool(pool, rng, 1, options)
	if err != nil {
		t.Fatalf("Failed to select vowels: %v", err)
	}

	if len(selected) != 1 {
		t.Errorf("Expected 1 vowel, got %d", len(selected))
	}

	if selected[0].Type != PhonemeTypeVowel {
		t.Errorf("Expected vowel type, got %v", selected[0].Type)
	}
}

func TestSelectFromPool_FilterByManner(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	pool := NewPool()

	// Add consonants with different manners
	consonants := []Phoneme{
		{Symbol: "p", Type: PhonemeTypeConsonant, Weight: 1.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}},
		{Symbol: "m", Type: PhonemeTypeConsonant, Weight: 1.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerNasal}},
		{Symbol: "f", Type: PhonemeTypeConsonant, Weight: 1.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerFricative}},
	}

	for _, c := range consonants {
		pool.AddConsonant(c)
	}

	// Test filtering by plosive manner
	plosiveManner := ConsonantMannerPlosive
	consonantType := PhonemeTypeConsonant
	options := ChoiceOptions{
		UseWeighted:    false,
		FilterByType:   &consonantType,
		FilterByManner: &plosiveManner,
	}

	selected, err := SelectFromPool(pool, rng, 1, options)
	if err != nil {
		t.Fatalf("Failed to select plosives: %v", err)
	}

	if len(selected) != 1 {
		t.Errorf("Expected 1 plosive, got %d", len(selected))
	}

	if selected[0].Consonant.Manner != ConsonantMannerPlosive {
		t.Errorf("Expected plosive manner, got %v", selected[0].Consonant.Manner)
	}
}

func TestSelectFromPool_FilterByHeight(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	pool := NewPool()

	// Add vowels with different heights
	vowels := []Phoneme{
		{Symbol: "a", Type: PhonemeTypeVowel, Weight: 1.0, Vowel: &VowelSpec{Height: VowelHeightLow, Backness: VowelBacknessFront, Rounded: false}},
		{Symbol: "e", Type: PhonemeTypeVowel, Weight: 1.0, Vowel: &VowelSpec{Height: VowelHeightMid, Backness: VowelBacknessFront, Rounded: false}},
		{Symbol: "i", Type: PhonemeTypeVowel, Weight: 1.0, Vowel: &VowelSpec{Height: VowelHeightHigh, Backness: VowelBacknessFront, Rounded: false}},
	}

	for _, v := range vowels {
		pool.AddVowel(v)
	}

	// Test filtering by high height
	highHeight := VowelHeightHigh
	vowelType := PhonemeTypeVowel
	options := ChoiceOptions{
		UseWeighted:    false,
		FilterByType:   &vowelType,
		FilterByHeight: &highHeight,
	}

	selected, err := SelectFromPool(pool, rng, 1, options)
	if err != nil {
		t.Fatalf("Failed to select high vowels: %v", err)
	}

	if len(selected) != 1 {
		t.Errorf("Expected 1 high vowel, got %d", len(selected))
	}

	if selected[0].Vowel.Height != VowelHeightHigh {
		t.Errorf("Expected high height, got %v", selected[0].Vowel.Height)
	}
}

func TestSelectFromPool_NoDuplicates(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	pool := NewPool()

	// Add a few phonemes
	consonants := []Phoneme{
		{Symbol: "p", Type: PhonemeTypeConsonant, Weight: 1.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}},
		{Symbol: "b", Type: PhonemeTypeConsonant, Weight: 1.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}},
		{Symbol: "t", Type: PhonemeTypeConsonant, Weight: 1.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}},
	}

	for _, c := range consonants {
		pool.AddConsonant(c)
	}

	// Try to select more phonemes than available
	options := ChoiceOptions{UseWeighted: false}
	selected, err := SelectFromPool(pool, rng, 5, options)
	if err != nil {
		t.Fatalf("Failed to select phonemes: %v", err)
	}

	// Should only get the available phonemes
	if len(selected) != 3 {
		t.Errorf("Expected 3 phonemes (all available), got %d", len(selected))
	}

	// Verify no duplicates
	symbols := make(map[string]bool)
	for _, phoneme := range selected {
		if symbols[phoneme.Symbol] {
			t.Errorf("Duplicate phoneme selected: %s", phoneme.Symbol)
		}
		symbols[phoneme.Symbol] = true
	}
}

func TestSelectRandomConsonants(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	pool := NewPool()

	// Add some consonants
	consonants := []Phoneme{
		{Symbol: "p", Type: PhonemeTypeConsonant, Weight: 1.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}},
		{Symbol: "b", Type: PhonemeTypeConsonant, Weight: 1.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}},
		{Symbol: "m", Type: PhonemeTypeConsonant, Weight: 1.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerNasal}},
	}

	for _, c := range consonants {
		pool.AddConsonant(c)
	}

	// Test without manner filter
	selected, err := SelectRandomConsonants(pool, rng, 2, nil)
	if err != nil {
		t.Fatalf("Failed to select random consonants: %v", err)
	}

	if len(selected) != 2 {
		t.Errorf("Expected 2 consonants, got %d", len(selected))
	}

	// Verify all are consonants
	for _, phoneme := range selected {
		if phoneme.Type != PhonemeTypeConsonant {
			t.Errorf("Expected consonant type, got %v", phoneme.Type)
		}
	}

	// Test with manner filter
	plosiveManner := ConsonantMannerPlosive
	selected, err = SelectRandomConsonants(pool, rng, 1, &plosiveManner)
	if err != nil {
		t.Fatalf("Failed to select plosive consonants: %v", err)
	}

	if len(selected) != 1 {
		t.Errorf("Expected 1 plosive consonant, got %d", len(selected))
	}

	if selected[0].Consonant.Manner != ConsonantMannerPlosive {
		t.Errorf("Expected plosive manner, got %v", selected[0].Consonant.Manner)
	}
}

func TestSelectRandomVowels(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	pool := NewPool()

	// Add some vowels
	vowels := []Phoneme{
		{Symbol: "a", Type: PhonemeTypeVowel, Weight: 1.0, Vowel: &VowelSpec{Height: VowelHeightLow, Backness: VowelBacknessFront, Rounded: false}},
		{Symbol: "e", Type: PhonemeTypeVowel, Weight: 1.0, Vowel: &VowelSpec{Height: VowelHeightMid, Backness: VowelBacknessFront, Rounded: false}},
		{Symbol: "i", Type: PhonemeTypeVowel, Weight: 1.0, Vowel: &VowelSpec{Height: VowelHeightHigh, Backness: VowelBacknessFront, Rounded: false}},
	}

	for _, v := range vowels {
		pool.AddVowel(v)
	}

	// Test without height filter
	selected, err := SelectRandomVowels(pool, rng, 2, nil)
	if err != nil {
		t.Fatalf("Failed to select random vowels: %v", err)
	}

	if len(selected) != 2 {
		t.Errorf("Expected 2 vowels, got %d", len(selected))
	}

	// Verify all are vowels
	for _, phoneme := range selected {
		if phoneme.Type != PhonemeTypeVowel {
			t.Errorf("Expected vowel type, got %v", phoneme.Type)
		}
	}

	// Test with height filter
	highHeight := VowelHeightHigh
	selected, err = SelectRandomVowels(pool, rng, 1, &highHeight)
	if err != nil {
		t.Fatalf("Failed to select high vowels: %v", err)
	}

	if len(selected) != 1 {
		t.Errorf("Expected 1 high vowel, got %d", len(selected))
	}

	if selected[0].Vowel.Height != VowelHeightHigh {
		t.Errorf("Expected high height, got %v", selected[0].Vowel.Height)
	}
}

func TestSelectWeightedConsonants(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	pool := NewPool()

	// Add consonants with different weights
	consonants := []Phoneme{
		{Symbol: "rare", Type: PhonemeTypeConsonant, Weight: 1.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}},
		{Symbol: "common", Type: PhonemeTypeConsonant, Weight: 10.0, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}},
	}

	for _, c := range consonants {
		pool.AddConsonant(c)
	}

	// Test weighted selection
	selected, err := SelectWeightedConsonants(pool, rng, 1, nil)
	if err != nil {
		t.Fatalf("Failed to select weighted consonants: %v", err)
	}

	if len(selected) != 1 {
		t.Errorf("Expected 1 consonant, got %d", len(selected))
	}

	// Verify it's a consonant
	if selected[0].Type != PhonemeTypeConsonant {
		t.Errorf("Expected consonant type, got %v", selected[0].Type)
	}
}

func TestSelectWeightedVowels(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	pool := NewPool()

	// Add vowels with different weights
	vowels := []Phoneme{
		{Symbol: "rare", Type: PhonemeTypeVowel, Weight: 1.0, Vowel: &VowelSpec{Height: VowelHeightLow, Backness: VowelBacknessFront, Rounded: false}},
		{Symbol: "common", Type: PhonemeTypeVowel, Weight: 10.0, Vowel: &VowelSpec{Height: VowelHeightMid, Backness: VowelBacknessFront, Rounded: false}},
	}

	for _, v := range vowels {
		pool.AddVowel(v)
	}

	// Test weighted selection
	selected, err := SelectWeightedVowels(pool, rng, 1, nil)
	if err != nil {
		t.Fatalf("Failed to select weighted vowels: %v", err)
	}

	if len(selected) != 1 {
		t.Errorf("Expected 1 vowel, got %d", len(selected))
	}

	// Verify it's a vowel
	if selected[0].Type != PhonemeTypeVowel {
		t.Errorf("Expected vowel type, got %v", selected[0].Type)
	}
}

func TestSelectFromPool_EdgeCases(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	pool := NewPool()

	// Test with negative count
	selected, err := SelectFromPool(pool, rng, -1, ChoiceOptions{})
	if err != nil {
		t.Errorf("Expected no error for negative count, got %v", err)
	}

	if len(selected) != 0 {
		t.Errorf("Expected 0 phonemes for negative count, got %d", len(selected))
	}

	// Test with nil pool
	selected, err = SelectFromPool(nil, rng, 1, ChoiceOptions{})
	if err != nil {
		t.Errorf("Expected no error for nil pool, got %v", err)
	}

	if len(selected) != 0 {
		t.Errorf("Expected 0 phonemes for nil pool, got %d", len(selected))
	}
}
