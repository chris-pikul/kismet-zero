package morphology

import (
	"math/rand/v2"
	"testing"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
)

func TestMorphemeTypeString(t *testing.T) {
	tests := []struct {
		name         string
		morphemeType MorphemeType
		expected     string
	}{
		{"Unknown", MorphemeTypeUnknown, "unknown"},
		{"Root", MorphemeTypeRoot, "root"},
		{"Prefix", MorphemeTypePrefix, "prefix"},
		{"Suffix", MorphemeTypeSuffix, "suffix"},
		{"Infix", MorphemeTypeInfix, "infix"},
		{"Circumfix", MorphemeTypeCircumfix, "circumfix"},
		{"Invalid", MorphemeTypeCircumfix + 1, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.morphemeType.String()
			if result != tt.expected {
				t.Errorf("MorphemeType.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMorphemeFrequencyString(t *testing.T) {
	tests := []struct {
		name      string
		frequency MorphemeFrequency
		expected  string
	}{
		{"Unknown", MorphemeFrequencyUnknown, "unknown"},
		{"Rare", MorphemeFrequencyRare, "rare"},
		{"Uncommon", MorphemeFrequencyUncommon, "uncommon"},
		{"Common", MorphemeFrequencyCommon, "common"},
		{"VeryCommon", MorphemeFrequencyVeryCommon, "very_common"},
		{"Invalid", MorphemeFrequencyVeryCommon + 1, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.frequency.String()
			if result != tt.expected {
				t.Errorf("MorphemeFrequency.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestWordCategoryString(t *testing.T) {
	tests := []struct {
		name     string
		category WordCategory
		expected string
	}{
		{"Unknown", WordCategoryUnknown, "unknown"},
		{"Noun", WordCategoryNoun, "noun"},
		{"Verb", WordCategoryVerb, "verb"},
		{"Adjective", WordCategoryAdjective, "adjective"},
		{"Adverb", WordCategoryAdverb, "adverb"},
		{"Pronoun", WordCategoryPronoun, "pronoun"},
		{"Preposition", WordCategoryPreposition, "preposition"},
		{"Conjunction", WordCategoryConjunction, "conjunction"},
		{"Interjection", WordCategoryInterjection, "interjection"},
		{"Invalid", WordCategoryInterjection + 1, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.category.String()
			if result != tt.expected {
				t.Errorf("WordCategory.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestWordFrequencyString(t *testing.T) {
	tests := []struct {
		name      string
		frequency WordFrequency
		expected  string
	}{
		{"Unknown", WordFrequencyUnknown, "unknown"},
		{"Rare", WordFrequencyRare, "rare"},
		{"Uncommon", WordFrequencyUncommon, "uncommon"},
		{"Common", WordFrequencyCommon, "common"},
		{"VeryCommon", WordFrequencyVeryCommon, "very_common"},
		{"Invalid", WordFrequencyVeryCommon + 1, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.frequency.String()
			if result != tt.expected {
				t.Errorf("WordFrequency.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNameTypeString(t *testing.T) {
	tests := []struct {
		name     string
		nameType NameType
		expected string
	}{
		{"Unknown", NameTypeUnknown, "unknown"},
		{"Personal", NameTypePersonal, "personal"},
		{"Family", NameTypeFamily, "family"},
		{"Place", NameTypePlace, "place"},
		{"Title", NameTypeTitle, "title"},
		{"Deity", NameTypeDeity, "deity"},
		{"Artifact", NameTypeArtifact, "artifact"},
		{"Invalid", NameTypeArtifact + 1, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.nameType.String()
			if result != tt.expected {
				t.Errorf("NameType.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMorphemeListUnweightedChoice(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	// Test empty list
	emptyList := MorphemeList{}
	result := emptyList.UnweightedChoice(rng)
	if result != nil {
		t.Errorf("UnweightedChoice on empty list should return nil, got %v", result)
	}

	// Test single item
	singleItem := MorphemeList{
		{ID: "test1", Meaning: "test", Weight: 1.0},
	}
	result = singleItem.UnweightedChoice(rng)
	if result == nil || result.ID != "test1" {
		t.Errorf("UnweightedChoice on single item should return the item, got %v", result)
	}

	// Test multiple items
	multipleItems := MorphemeList{
		{ID: "test1", Meaning: "test1", Weight: 1.0},
		{ID: "test2", Meaning: "test2", Weight: 2.0},
		{ID: "test3", Meaning: "test3", Weight: 3.0},
	}

	// Run multiple times to ensure we get different results
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		result = multipleItems.UnweightedChoice(rng)
		if result != nil {
			seen[result.ID] = true
		}
	}

	// We should see all items at least once
	if len(seen) < 2 {
		t.Errorf("UnweightedChoice should return different items, only saw %d", len(seen))
	}
}

func TestMorphemeListWeightedChoice(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	// Test empty list
	emptyList := MorphemeList{}
	result := emptyList.WeightedChoice(rng)
	if result != nil {
		t.Errorf("WeightedChoice on empty list should return nil, got %v", result)
	}

	// Test single item
	singleItem := MorphemeList{
		{ID: "test1", Meaning: "test", Weight: 1.0},
	}
	result = singleItem.WeightedChoice(rng)
	if result == nil || result.ID != "test1" {
		t.Errorf("WeightedChoice on single item should return the item, got %v", result)
	}

	// Test multiple items with different weights
	multipleItems := MorphemeList{
		{ID: "test1", Meaning: "test1", Weight: 1.0},
		{ID: "test2", Meaning: "test2", Weight: 2.0},
		{ID: "test3", Meaning: "test3", Weight: 3.0},
	}

	// Run multiple times to ensure we get different results
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		result = multipleItems.WeightedChoice(rng)
		if result != nil {
			seen[result.ID] = true
		}
	}

	// We should see all items at least once
	if len(seen) < 2 {
		t.Errorf("WeightedChoice should return different items, only saw %d", len(seen))
	}
}

func TestMorphemeListWeightedChoicePanic(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	// Test list with zero weights
	zeroWeightList := MorphemeList{
		{ID: "test1", Meaning: "test1", Weight: 0.0},
		{ID: "test2", Meaning: "test2", Weight: 0.0},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("WeightedChoice with zero weights should panic")
		}
	}()

	zeroWeightList.WeightedChoice(rng)
}

func TestWordListUnweightedChoice(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	// Test empty list
	emptyList := WordList{}
	result := emptyList.UnweightedChoice(rng)
	if result != nil {
		t.Errorf("UnweightedChoice on empty list should return nil, got %v", result)
	}

	// Test single item
	singleItem := WordList{
		{ID: "test1", Meaning: "test", Weight: 1.0},
	}
	result = singleItem.UnweightedChoice(rng)
	if result == nil || result.ID != "test1" {
		t.Errorf("UnweightedChoice on single item should return the item, got %v", result)
	}

	// Test multiple items
	multipleItems := WordList{
		{ID: "test1", Meaning: "test1", Weight: 1.0},
		{ID: "test2", Meaning: "test2", Weight: 2.0},
		{ID: "test3", Meaning: "test3", Weight: 3.0},
	}

	// Run multiple times to ensure we get different results
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		result = multipleItems.UnweightedChoice(rng)
		if result != nil {
			seen[result.ID] = true
		}
	}

	// We should see all items at least once
	if len(seen) < 2 {
		t.Errorf("UnweightedChoice should return different items, only saw %d", len(seen))
	}
}

func TestWordListWeightedChoice(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	// Test empty list
	emptyList := WordList{}
	result := emptyList.WeightedChoice(rng)
	if result != nil {
		t.Errorf("WeightedChoice on empty list should return nil, got %v", result)
	}

	// Test single item
	singleItem := WordList{
		{ID: "test1", Meaning: "test", Weight: 1.0},
	}
	result = singleItem.WeightedChoice(rng)
	if result == nil || result.ID != "test1" {
		t.Errorf("WeightedChoice on single item should return the item, got %v", result)
	}

	// Test multiple items with different weights
	multipleItems := WordList{
		{ID: "test1", Meaning: "test1", Weight: 1.0},
		{ID: "test2", Meaning: "test2", Weight: 2.0},
		{ID: "test3", Meaning: "test3", Weight: 3.0},
	}

	// Run multiple times to ensure we get different results
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		result = multipleItems.WeightedChoice(rng)
		if result != nil {
			seen[result.ID] = true
		}
	}

	// We should see all items at least once
	if len(seen) < 2 {
		t.Errorf("WeightedChoice should return different items, only saw %d", len(seen))
	}
}

func TestWordListWeightedChoicePanic(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	// Test list with zero weights
	zeroWeightList := WordList{
		{ID: "test1", Meaning: "test1", Weight: 0.0},
		{ID: "test2", Meaning: "test2", Weight: 0.0},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("WeightedChoice with zero weights should panic")
		}
	}()

	zeroWeightList.WeightedChoice(rng)
}

func TestMorphemeStruct(t *testing.T) {
	phoneme1 := phoneme.Phoneme{Symbol: "t", Weight: 1.0}
	phoneme2 := phoneme.Phoneme{Symbol: "a", Weight: 1.0}

	morpheme := &Morpheme{
		ID:        "test_morpheme",
		Type:      MorphemeTypeRoot,
		Meaning:   "test",
		Phonemes:  []phoneme.Phoneme{phoneme1, phoneme2},
		Syllables: []string{"ta"},
		Weight:    1.0,
		Culture:   "test_culture",
		Frequency: MorphemeFrequencyCommon,
	}

	if morpheme.ID != "test_morpheme" {
		t.Errorf("Morpheme ID mismatch: got %s, want test_morpheme", morpheme.ID)
	}

	if morpheme.Type != MorphemeTypeRoot {
		t.Errorf("Morpheme Type mismatch: got %v, want %v", morpheme.Type, MorphemeTypeRoot)
	}

	if len(morpheme.Phonemes) != 2 {
		t.Errorf("Morpheme Phonemes count mismatch: got %d, want 2", len(morpheme.Phonemes))
	}

	if len(morpheme.Syllables) != 1 {
		t.Errorf("Morpheme Syllables count mismatch: got %d, want 1", len(morpheme.Syllables))
	}
}

func TestWordStruct(t *testing.T) {
	morpheme1 := &Morpheme{ID: "root", Meaning: "walk", Weight: 1.0}
	morpheme2 := &Morpheme{ID: "suffix", Meaning: "ing", Weight: 0.5}

	word := &Word{
		ID:        "test_word",
		Morphemes: []Morpheme{*morpheme1, *morpheme2},
		Meaning:   "walking",
		Phonemes:  []phoneme.Phoneme{},
		Written:   "walking",
		Category:  WordCategoryVerb,
		Culture:   "test_culture",
		Frequency: WordFrequencyCommon,
		Weight:    0.75,
	}

	if word.ID != "test_word" {
		t.Errorf("Word ID mismatch: got %s, want test_word", word.ID)
	}

	if len(word.Morphemes) != 2 {
		t.Errorf("Word Morphemes count mismatch: got %d, want 2", len(word.Morphemes))
	}

	if word.Category != WordCategoryVerb {
		t.Errorf("Word Category mismatch: got %v, want %v", word.Category, WordCategoryVerb)
	}

	if word.Weight != 0.75 {
		t.Errorf("Word Weight mismatch: got %f, want 0.75", word.Weight)
	}
}

func TestNameStruct(t *testing.T) {
	name := &Name{
		ID:        "test_name",
		Type:      NameTypePersonal,
		Value:     "John",
		Meaning:   "noble person (male)",
		Culture:   "test_culture",
		Gender:    "male",
		Formality: "formal",
	}

	if name.ID != "test_name" {
		t.Errorf("Name ID mismatch: got %s, want test_name", name.ID)
	}

	if name.Type != NameTypePersonal {
		t.Errorf("Name Type mismatch: got %v, want %v", name.Type, NameTypePersonal)
	}

	if name.Value != "John" {
		t.Errorf("Name Value mismatch: got %s, want John", name.Value)
	}

	if name.Gender != "male" {
		t.Errorf("Name Gender mismatch: got %s, want male", name.Gender)
	}
}
