package phoneme

import (
	"encoding/json"
	"testing"
)

func TestCommonPool_Structure(t *testing.T) {
	// Test that CommonPool is not nil
	if CommonPool.Consonants == nil {
		t.Fatal("CommonPool.Consonants should not be nil")
	}

	if CommonPool.Vowels == nil {
		t.Fatal("CommonPool.Vowels should not be nil")
	}

	// Test that CommonPool has expected number of phonemes
	expectedConsonants := 17
	expectedVowels := 6

	if len(CommonPool.Consonants) != expectedConsonants {
		t.Errorf("Expected %d consonants in CommonPool, got %d", expectedConsonants, len(CommonPool.Consonants))
	}

	if len(CommonPool.Vowels) != expectedVowels {
		t.Errorf("Expected %d vowels in CommonPool, got %d", expectedVowels, len(CommonPool.Vowels))
	}
}

func TestCommonPool_Consonants(t *testing.T) {
	// Test that all consonants have valid types
	for i, consonant := range CommonPool.Consonants {
		if consonant.Type != PhonemeTypeConsonant {
			t.Errorf("Consonant %d (%s) has invalid type: %v", i, consonant.Symbol, consonant.Type)
		}

		if consonant.Consonant == nil {
			t.Errorf("Consonant %d (%s) missing ConsonantSpec", i, consonant.Symbol)
		}

		if consonant.Symbol == "" {
			t.Errorf("Consonant %d has empty symbol", i)
		}

		if consonant.Weight <= 0 {
			t.Errorf("Consonant %d (%s) has invalid weight: %f", i, consonant.Symbol, consonant.Weight)
		}
	}

	// Test specific consonant categories
	plosives := []string{"p", "b", "t", "d", "k", "g"}
	nasals := []string{"m", "n", "ŋ"}
	fricatives := []string{"f", "v"}
	sibilants := []string{"s", "z"}
	approximants := []string{"l", "r", "j", "w"}

	// Verify plosives
	for _, symbol := range plosives {
		found := false
		for _, consonant := range CommonPool.Consonants {
			if consonant.Symbol == symbol && consonant.Consonant.Manner == ConsonantMannerPlosive {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Plosive %s not found in CommonPool", symbol)
		}
	}

	// Verify nasals
	for _, symbol := range nasals {
		found := false
		for _, consonant := range CommonPool.Consonants {
			if consonant.Symbol == symbol && consonant.Consonant.Manner == ConsonantMannerNasal {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Nasal %s not found in CommonPool", symbol)
		}
	}

	// Verify fricatives
	for _, symbol := range fricatives {
		found := false
		for _, consonant := range CommonPool.Consonants {
			if consonant.Symbol == symbol && consonant.Consonant.Manner == ConsonantMannerFricative {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Fricative %s not found in CommonPool", symbol)
		}
	}

	// Verify sibilants
	for _, symbol := range sibilants {
		found := false
		for _, consonant := range CommonPool.Consonants {
			if consonant.Symbol == symbol && consonant.Consonant.Manner == ConsonantMannerSibilants {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Sibilant %s not found in CommonPool", symbol)
		}
	}

	// Verify approximants
	for _, symbol := range approximants {
		found := false
		for _, consonant := range CommonPool.Consonants {
			if consonant.Symbol == symbol && (consonant.Consonant.Manner == ConsonantMannerLateralApproximant ||
				consonant.Consonant.Manner == ConsonantMannerTrill ||
				consonant.Consonant.Manner == ConsonantMannerSemivowel) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Approximant %s not found in CommonPool", symbol)
		}
	}
}

func TestCommonPool_Vowels(t *testing.T) {
	// Test that all vowels have valid types
	for i, vowel := range CommonPool.Vowels {
		if vowel.Type != PhonemeTypeVowel {
			t.Errorf("Vowel %d (%s) has invalid type: %v", i, vowel.Symbol, vowel.Type)
		}

		if vowel.Vowel == nil {
			t.Errorf("Vowel %d (%s) missing VowelSpec", i, vowel.Symbol)
		}

		if vowel.Symbol == "" {
			t.Errorf("Vowel %d has empty symbol", i)
		}

		if vowel.Weight <= 0 {
			t.Errorf("Vowel %d (%s) has invalid weight: %f", i, vowel.Symbol, vowel.Weight)
		}
	}

	// Test specific vowel categories
	frontVowels := []string{"a", "e", "i"}
	backVowels := []string{"o", "u"}
	centralVowels := []string{"ə"}

	// Verify front vowels
	for _, symbol := range frontVowels {
		found := false
		for _, vowel := range CommonPool.Vowels {
			if vowel.Symbol == symbol && vowel.Vowel.Backness == VowelBacknessFront {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Front vowel %s not found in CommonPool", symbol)
		}
	}

	// Verify back vowels
	for _, symbol := range backVowels {
		found := false
		for _, vowel := range CommonPool.Vowels {
			if vowel.Symbol == symbol && vowel.Vowel.Backness == VowelBacknessBack {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Back vowel %s not found in CommonPool", symbol)
		}
	}

	// Verify central vowels
	for _, symbol := range centralVowels {
		found := false
		for _, vowel := range CommonPool.Vowels {
			if vowel.Symbol == symbol && vowel.Vowel.Backness == VowelBacknessCentral {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Central vowel %s not found in CommonPool", symbol)
		}
	}

	// Test vowel heights
	lowVowels := []string{"a"}
	midVowels := []string{"e", "o", "ə"}
	highVowels := []string{"i", "u"}

	// Verify low vowels
	for _, symbol := range lowVowels {
		found := false
		for _, vowel := range CommonPool.Vowels {
			if vowel.Symbol == symbol && vowel.Vowel.Height == VowelHeightLow {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Low vowel %s not found in CommonPool", symbol)
		}
	}

	// Verify mid vowels
	for _, symbol := range midVowels {
		found := false
		for _, vowel := range CommonPool.Vowels {
			if vowel.Symbol == symbol && vowel.Vowel.Height == VowelHeightMid {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Mid vowel %s not found in CommonPool", symbol)
		}
	}

	// Verify high vowels
	for _, symbol := range highVowels {
		found := false
		for _, vowel := range CommonPool.Vowels {
			if vowel.Symbol == symbol && vowel.Vowel.Height == VowelHeightHigh {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("High vowel %s not found in CommonPool", symbol)
		}
	}
}

func TestCommonPool_Weights(t *testing.T) {
	// Test that weights are reasonable
	for i, consonant := range CommonPool.Consonants {
		if consonant.Weight < 0.5 || consonant.Weight > 2.0 {
			t.Errorf("Consonant %d (%s) has unusual weight: %f", i, consonant.Symbol, consonant.Weight)
		}
	}

	for i, vowel := range CommonPool.Vowels {
		if vowel.Weight < 0.5 || vowel.Weight > 2.0 {
			t.Errorf("Vowel %d (%s) has unusual weight: %f", i, vowel.Symbol, vowel.Weight)
		}
	}

	// Test specific weight values
	// Most phonemes should have weight 1.0
	weight1Count := 0
	for _, consonant := range CommonPool.Consonants {
		if consonant.Weight == 1.0 {
			weight1Count++
		}
	}

	// Most consonants should have weight 1.0
	if weight1Count < len(CommonPool.Consonants)*3/4 {
		t.Errorf("Expected most consonants to have weight 1.0, but only %d of %d do",
			weight1Count, len(CommonPool.Consonants))
	}

	// Test that ŋ has weight 0.8 (less common)
	found := false
	for _, consonant := range CommonPool.Consonants {
		if consonant.Symbol == "ŋ" && consonant.Weight == 0.8 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected ŋ to have weight 0.8")
	}

	// Test that r has weight 0.8 (less common)
	found = false
	for _, consonant := range CommonPool.Consonants {
		if consonant.Symbol == "r" && consonant.Weight == 0.8 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected r to have weight 0.8")
	}

	// Test that ə has weight 0.7 (less common)
	found = false
	for _, vowel := range CommonPool.Vowels {
		if vowel.Symbol == "ə" && vowel.Weight == 0.7 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected ə to have weight 0.7")
	}
}

func TestCommonPool_Uniqueness(t *testing.T) {
	// Test that all symbols are unique
	consonantSymbols := make(map[string]bool)
	vowelSymbols := make(map[string]bool)

	for _, consonant := range CommonPool.Consonants {
		if consonantSymbols[consonant.Symbol] {
			t.Errorf("Duplicate consonant symbol: %s", consonant.Symbol)
		}
		consonantSymbols[consonant.Symbol] = true
	}

	for _, vowel := range CommonPool.Vowels {
		if vowelSymbols[vowel.Symbol] {
			t.Errorf("Duplicate vowel symbol: %s", vowel.Symbol)
		}
		vowelSymbols[vowel.Symbol] = true
	}

	// Test that no symbol appears in both consonants and vowels
	for symbol := range consonantSymbols {
		if vowelSymbols[symbol] {
			t.Errorf("Symbol %s appears in both consonants and vowels", symbol)
		}
	}
}

func TestCommonPool_JSON(t *testing.T) {
	// Test that CommonPool can be marshaled to JSON
	data, err := json.Marshal(CommonPool)
	if err != nil {
		t.Fatalf("Failed to marshal CommonPool: %v", err)
	}

	// Test that CommonPool can be unmarshaled from JSON
	var unmarshaled PhonemePool
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal CommonPool: %v", err)
	}

	// Verify the unmarshaled pool has the same structure
	if len(unmarshaled.Consonants) != len(CommonPool.Consonants) {
		t.Errorf("Unmarshaled consonants count mismatch: got %d, want %d",
			len(unmarshaled.Consonants), len(CommonPool.Consonants))
	}

	if len(unmarshaled.Vowels) != len(CommonPool.Vowels) {
		t.Errorf("Unmarshaled vowels count mismatch: got %d, want %d",
			len(unmarshaled.Vowels), len(CommonPool.Vowels))
	}
}

func TestCommonPool_ContentValidation(t *testing.T) {
	// Test that all phonemes have valid specifications
	for i, consonant := range CommonPool.Consonants {
		if consonant.Consonant.Manner == ConsonantMannerUnknown {
			t.Errorf("Consonant %d (%s) has unknown manner", i, consonant.Symbol)
		}
	}

	for i, vowel := range CommonPool.Vowels {
		if vowel.Vowel.Height == VowelHeightUnknown {
			t.Errorf("Vowel %d (%s) has unknown height", i, vowel.Symbol)
		}
		if vowel.Vowel.Backness == VowelBacknessUnknown {
			t.Errorf("Vowel %d (%s) has unknown backness", i, vowel.Symbol)
		}
	}

	// Test that rounded vowels are back vowels
	for i, vowel := range CommonPool.Vowels {
		if vowel.Vowel.Rounded && vowel.Vowel.Backness != VowelBacknessBack {
			t.Errorf("Vowel %d (%s) is rounded but not back", i, vowel.Symbol)
		}
	}

	// Test that front vowels are not rounded (except for special cases)
	for i, vowel := range CommonPool.Vowels {
		if vowel.Vowel.Backness == VowelBacknessFront && vowel.Vowel.Rounded {
			t.Errorf("Vowel %d (%s) is front but rounded", i, vowel.Symbol)
		}
	}
}
