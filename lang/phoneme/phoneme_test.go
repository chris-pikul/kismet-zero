package phoneme

import (
	"encoding/json"
	"math/rand/v2"
	"testing"
)

func TestPhoneme_JSON(t *testing.T) {
	consonant := Phoneme{
		Symbol: "p",
		Type:   PhonemeTypeConsonant,
		Weight: 1.5,
		Consonant: &ConsonantSpec{
			Manner: ConsonantMannerPlosive,
		},
	}

	vowel := Phoneme{
		Symbol: "a",
		Type:   PhonemeTypeVowel,
		Weight: 2.0,
		Vowel: &VowelSpec{
			Height:    VowelHeightLow,
			Backness:  VowelBacknessFront,
			Rounded:   false,
			Diphthong: false,
		},
	}

	// Test consonant JSON marshaling
	consonantData, err := json.Marshal(consonant)
	if err != nil {
		t.Fatalf("Failed to marshal consonant: %v", err)
	}

	var unmarshaledConsonant Phoneme
	err = json.Unmarshal(consonantData, &unmarshaledConsonant)
	if err != nil {
		t.Fatalf("Failed to unmarshal consonant: %v", err)
	}

	if unmarshaledConsonant.Symbol != consonant.Symbol {
		t.Errorf("Symbol mismatch: got %s, want %s", unmarshaledConsonant.Symbol, consonant.Symbol)
	}
	if unmarshaledConsonant.Type != consonant.Type {
		t.Errorf("Type mismatch: got %v, want %v", unmarshaledConsonant.Type, consonant.Type)
	}
	if unmarshaledConsonant.Weight != consonant.Weight {
		t.Errorf("Weight mismatch: got %f, want %f", unmarshaledConsonant.Weight, consonant.Weight)
	}
	if unmarshaledConsonant.Consonant == nil {
		t.Error("Consonant spec should not be nil")
	}
	if unmarshaledConsonant.Consonant.Manner != consonant.Consonant.Manner {
		t.Errorf("Consonant manner mismatch: got %v, want %v", unmarshaledConsonant.Consonant.Manner, consonant.Consonant.Manner)
	}

	// Test vowel JSON marshaling
	vowelData, err := json.Marshal(vowel)
	if err != nil {
		t.Fatalf("Failed to marshal vowel: %v", err)
	}

	var unmarshaledVowel Phoneme
	err = json.Unmarshal(vowelData, &unmarshaledVowel)
	if err != nil {
		t.Fatalf("Failed to unmarshal vowel: %v", err)
	}

	if unmarshaledVowel.Symbol != vowel.Symbol {
		t.Errorf("Symbol mismatch: got %s, want %s", unmarshaledVowel.Symbol, vowel.Symbol)
	}
	if unmarshaledVowel.Type != vowel.Type {
		t.Errorf("Type mismatch: got %v, want %v", unmarshaledVowel.Type, vowel.Type)
	}
	if unmarshaledVowel.Weight != vowel.Weight {
		t.Errorf("Weight mismatch: got %f, want %f", unmarshaledVowel.Weight, vowel.Weight)
	}
	if unmarshaledVowel.Vowel == nil {
		t.Error("Vowel spec should not be nil")
	}
	if unmarshaledVowel.Vowel.Height != vowel.Vowel.Height {
		t.Errorf("Vowel height mismatch: got %v, want %v", unmarshaledVowel.Vowel.Height, vowel.Vowel.Height)
	}
}

func TestPhonemeList_UnweightedChoice(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	list := PhonemeList{
		{Symbol: "a", Type: PhonemeTypeVowel, Weight: 1.0},
		{Symbol: "b", Type: PhonemeTypeConsonant, Weight: 2.0},
		{Symbol: "c", Type: PhonemeTypeConsonant, Weight: 3.0},
	}

	// Test with valid list
	choice := list.UnweightedChoice(rng)
	if choice == nil {
		t.Fatal("Expected non-nil choice from non-empty list")
	}

	// Verify choice is from the list
	found := false
	for _, phoneme := range list {
		if phoneme.Symbol == choice.Symbol {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Choice %s not found in original list", choice.Symbol)
	}

	// Test with empty list
	emptyList := PhonemeList{}
	choice = emptyList.UnweightedChoice(rng)
	if choice != nil {
		t.Error("Expected nil choice from empty list")
	}
}

func TestPhonemeList_WeightedChoice(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	// Test with valid weighted list
	list := PhonemeList{
		{Symbol: "a", Type: PhonemeTypeVowel, Weight: 1.0},
		{Symbol: "b", Type: PhonemeTypeConsonant, Weight: 2.0},
		{Symbol: "c", Type: PhonemeTypeConsonant, Weight: 3.0},
	}

	choice := list.WeightedChoice(rng)
	if choice == nil {
		t.Fatal("Expected non-nil choice from non-empty list")
	}

	// Verify choice is from the list
	found := false
	for _, phoneme := range list {
		if phoneme.Symbol == choice.Symbol {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Choice %s not found in original list", choice.Symbol)
	}

	// Test with empty list
	emptyList := PhonemeList{}
	choice = emptyList.WeightedChoice(rng)
	if choice != nil {
		t.Error("Expected nil choice from empty list")
	}

	// Test with zero weights
	zeroWeightList := PhonemeList{
		{Symbol: "a", Type: PhonemeTypeVowel, Weight: 0.0},
		{Symbol: "b", Type: PhonemeTypeConsonant, Weight: 0.0},
	}

	// This should panic due to zero total weight
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for zero weight list")
		}
	}()
	zeroWeightList.WeightedChoice(rng)
}

func TestPhonemeList_WeightedChoiceDistribution(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	// Create a list with very different weights to test distribution
	list := PhonemeList{
		{Symbol: "rare", Type: PhonemeTypeVowel, Weight: 1.0},
		{Symbol: "common", Type: PhonemeTypeConsonant, Weight: 10.0},
		{Symbol: "very_common", Type: PhonemeTypeConsonant, Weight: 100.0},
	}

	// Run many selections to test distribution
	selected := make(map[string]int)
	totalSelections := 10000

	for i := 0; i < totalSelections; i++ {
		choice := list.WeightedChoice(rng)
		if choice != nil {
			selected[choice.Symbol]++
		}
	}

	// Verify that higher weights result in more selections
	if selected["rare"] >= selected["common"] {
		t.Errorf("Rare phoneme selected %d times, common selected %d times - weights not working",
			selected["rare"], selected["common"])
	}
	if selected["common"] >= selected["very_common"] {
		t.Errorf("Common phoneme selected %d times, very common selected %d times - weights not working",
			selected["common"], selected["very_common"])
	}

	// Verify all phonemes were selected at least once
	for _, phoneme := range list {
		if selected[phoneme.Symbol] == 0 {
			t.Errorf("Phoneme %s was never selected", phoneme.Symbol)
		}
	}
}

func TestPhonemeList_EdgeCases(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	// Test single phoneme list
	singleList := PhonemeList{
		{Symbol: "x", Type: PhonemeTypeConsonant, Weight: 1.0},
	}

	choice := singleList.UnweightedChoice(rng)
	if choice == nil || choice.Symbol != "x" {
		t.Error("Single phoneme list should always return the same phoneme")
	}

	choice = singleList.WeightedChoice(rng)
	if choice == nil || choice.Symbol != "x" {
		t.Error("Single phoneme list should always return the same phoneme")
	}

	// Test with very large weights
	largeWeightList := PhonemeList{
		{Symbol: "a", Type: PhonemeTypeVowel, Weight: 1e6},
		{Symbol: "b", Type: PhonemeTypeConsonant, Weight: 2e6},
	}

	choice = largeWeightList.WeightedChoice(rng)
	if choice == nil {
		t.Error("Large weight list should work correctly")
	}
}
