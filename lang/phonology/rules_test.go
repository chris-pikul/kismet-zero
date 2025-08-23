package phonology

import (
	"testing"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
)

func TestMaxClusterSizeRule(t *testing.T) {
	rule := NewMaxClusterSizeRule(2)

	testCases := []struct {
		name    string
		seq     []phoneme.Phoneme
		isValid bool
	}{
		{
			name:    "Empty sequence",
			seq:     []phoneme.Phoneme{},
			isValid: true,
		},
		{
			name: "Single consonant",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
			},
			isValid: true,
		},
		{
			name: "Two consonants",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "t", Type: phoneme.PhonemeTypeConsonant},
			},
			isValid: true,
		},
		{
			name: "Three consonants - invalid",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "t", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "p", Type: phoneme.PhonemeTypeConsonant},
			},
			isValid: false,
		},
		{
			name: "Consonant cluster broken by vowel",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "t", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
				{Symbol: "p", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "s", Type: phoneme.PhonemeTypeConsonant},
			},
			isValid: true,
		},
		{
			name: "Multiple consonant clusters",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "t", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
				{Symbol: "p", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "s", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "i", Type: phoneme.PhonemeTypeVowel},
			},
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := rule.Validate(tc.seq)
			if result != tc.isValid {
				t.Errorf("Expected validity %v, got %v", tc.isValid, result)
			}
		})
	}
}

func TestVelarNasalOnsetRule(t *testing.T) {
	rule := NewVelarNasalOnsetRule()

	testCases := []struct {
		name    string
		seq     []phoneme.Phoneme
		isValid bool
	}{
		{
			name:    "Empty sequence",
			seq:     []phoneme.Phoneme{},
			isValid: true,
		},
		{
			name: "Single phoneme",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
			},
			isValid: true,
		},
		{
			name: "Consonant + Vowel",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
			},
			isValid: true,
		},
		{
			name: "Consonant + Consonant (non-nasal)",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{Manner: phoneme.ConsonantMannerPlosive}},
				{Symbol: "t", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{Manner: phoneme.ConsonantMannerPlosive}},
			},
			isValid: true,
		},
		{
			name: "Consonant + Nasal - invalid",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Manner: phoneme.ConsonantMannerPlosive, Place: phoneme.ConsonantPlaceVelar, Voicing: phoneme.ConsonantVoicingUnvoiced}},
				{Symbol: "n", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Manner: phoneme.ConsonantMannerNasal, Place: phoneme.ConsonantPlaceAlveolar, Voicing: phoneme.ConsonantVoicingVoiced}},
			},
			isValid: false,
		},
		{
			name: "Vowel + Consonant + Nasal",
			seq: []phoneme.Phoneme{
				{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{Manner: phoneme.ConsonantMannerPlosive}},
				{Symbol: "n", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{Manner: phoneme.ConsonantMannerNasal}},
			},
			isValid: true, // Not at onset
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := rule.Validate(tc.seq)
			if result != tc.isValid {
				t.Errorf("Expected validity %v, got %v", tc.isValid, result)
			}
		})
	}
}

func TestVowelRequirementRule(t *testing.T) {
	rule := NewVowelRequirementRule()

	testCases := []struct {
		name    string
		seq     []phoneme.Phoneme
		isValid bool
	}{
		{
			name:    "Empty sequence",
			seq:     []phoneme.Phoneme{},
			isValid: false,
		},
		{
			name: "Single vowel",
			seq: []phoneme.Phoneme{
				{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
			},
			isValid: true,
		},
		{
			name: "Single consonant",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
			},
			isValid: false,
		},
		{
			name: "Consonant + Vowel",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
			},
			isValid: true,
		},
		{
			name: "Vowel + Consonant",
			seq: []phoneme.Phoneme{
				{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
			},
			isValid: true,
		},
		{
			name: "Multiple consonants + Vowel",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "t", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
			},
			isValid: true,
		},
		{
			name: "Multiple consonants only",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "t", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "p", Type: phoneme.PhonemeTypeConsonant},
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := rule.Validate(tc.seq)
			if result != tc.isValid {
				t.Errorf("Expected validity %v, got %v", tc.isValid, result)
			}
		})
	}
}

func TestConsonantHarmonyRule(t *testing.T) {
	rule := NewConsonantHarmonyRule()

	testCases := []struct {
		name    string
		seq     []phoneme.Phoneme
		isValid bool
	}{
		{
			name:    "Empty sequence",
			seq:     []phoneme.Phoneme{},
			isValid: true,
		},
		{
			name: "Single phoneme",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
			},
			isValid: true,
		},
		{
			name: "Consonant + Vowel",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
			},
			isValid: true,
		},
		{
			name: "Vowel + Consonant",
			seq: []phoneme.Phoneme{
				{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
			},
			isValid: true,
		},
		{
			name: "Valid bilabial + labiodental",
			seq: []phoneme.Phoneme{
				{Symbol: "p", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Place: phoneme.ConsonantPlaceBilabial}},
				{Symbol: "f", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Place: phoneme.ConsonantPlaceLabiodental}},
			},
			isValid: true,
		},
		{
			name: "Valid alveolar + postalveolar",
			seq: []phoneme.Phoneme{
				{Symbol: "t", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Place: phoneme.ConsonantPlaceAlveolar}},
				{Symbol: "ʃ", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Place: phoneme.ConsonantPlacePostalveolar}},
			},
			isValid: true,
		},
		{
			name: "Invalid bilabial + velar",
			seq: []phoneme.Phoneme{
				{Symbol: "p", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Place: phoneme.ConsonantPlaceBilabial}},
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Place: phoneme.ConsonantPlaceVelar}},
			},
			isValid: false,
		},
		{
			name: "Invalid dental + glottal",
			seq: []phoneme.Phoneme{
				{Symbol: "θ", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Place: phoneme.ConsonantPlaceDental}},
				{Symbol: "h", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Place: phoneme.ConsonantPlaceGlottal}},
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := rule.Validate(tc.seq)
			if result != tc.isValid {
				t.Errorf("Expected validity %v, got %v", tc.isValid, result)
			}
		})
	}
}

func TestSyllableLengthRule(t *testing.T) {
	testCases := []struct {
		name      string
		minLength int
		maxLength int
		seq       []phoneme.Phoneme
		isValid   bool
	}{
		{
			name:      "Empty sequence, min 1",
			minLength: 1,
			maxLength: 5,
			seq:       []phoneme.Phoneme{},
			isValid:   false,
		},
		{
			name:      "Empty sequence, min 0",
			minLength: 0,
			maxLength: 5,
			seq:       []phoneme.Phoneme{},
			isValid:   true,
		},
		{
			name:      "Single phoneme, min 1 max 5",
			minLength: 1,
			maxLength: 5,
			seq: []phoneme.Phoneme{
				{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
			},
			isValid: true,
		},
		{
			name:      "Single phoneme, min 2 max 5",
			minLength: 2,
			maxLength: 5,
			seq: []phoneme.Phoneme{
				{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
			},
			isValid: false,
		},
		{
			name:      "Three phonemes, min 1 max 5",
			minLength: 1,
			maxLength: 5,
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
				{Symbol: "t", Type: phoneme.PhonemeTypeConsonant},
			},
			isValid: true,
		},
		{
			name:      "Six phonemes, min 1 max 5",
			minLength: 1,
			maxLength: 5,
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "t", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
				{Symbol: "p", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "s", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "i", Type: phoneme.PhonemeTypeVowel},
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rule := NewSyllableLengthRule(tc.minLength, tc.maxLength)
			result := rule.Validate(tc.seq)
			if result != tc.isValid {
				t.Errorf("Expected validity %v, got %v", tc.isValid, result)
			}
		})
	}
}

func TestVoicingAssimilationRule(t *testing.T) {
	testCases := []struct {
		name            string
		requireMatching bool
		seq             []phoneme.Phoneme
		isValid         bool
	}{
		{
			name:            "Empty sequence",
			requireMatching: true,
			seq:             []phoneme.Phoneme{},
			isValid:         true,
		},
		{
			name:            "Single consonant",
			requireMatching: true,
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
			},
			isValid: true,
		},
		{
			name:            "Matching voicing - both unvoiced",
			requireMatching: true,
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Voicing: phoneme.ConsonantVoicingUnvoiced}},
				{Symbol: "t", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Voicing: phoneme.ConsonantVoicingUnvoiced}},
			},
			isValid: true,
		},
		{
			name:            "Matching voicing - both voiced",
			requireMatching: true,
			seq: []phoneme.Phoneme{
				{Symbol: "b", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Voicing: phoneme.ConsonantVoicingVoiced}},
				{Symbol: "d", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Voicing: phoneme.ConsonantVoicingVoiced}},
			},
			isValid: true,
		},
		{
			name:            "Mismatched voicing - invalid when required",
			requireMatching: true,
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Voicing: phoneme.ConsonantVoicingUnvoiced}},
				{Symbol: "b", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Voicing: phoneme.ConsonantVoicingVoiced}},
			},
			isValid: false,
		},
		{
			name:            "Mismatched voicing - valid when not required",
			requireMatching: false,
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Voicing: phoneme.ConsonantVoicingUnvoiced}},
				{Symbol: "b", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Voicing: phoneme.ConsonantVoicingVoiced}},
			},
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rule := NewVoicingAssimilationRule(tc.requireMatching)
			result := rule.Validate(tc.seq)
			if result != tc.isValid {
				t.Errorf("Expected validity %v, got %v", tc.isValid, result)
			}
		})
	}
}

func TestMannerRestrictionsRule(t *testing.T) {
	rule := NewMannerRestrictionsRule()

	testCases := []struct {
		name    string
		seq     []phoneme.Phoneme
		isValid bool
	}{
		{
			name:    "Empty sequence",
			seq:     []phoneme.Phoneme{},
			isValid: true,
		},
		{
			name: "Single consonant",
			seq: []phoneme.Phoneme{
				{Symbol: "s", Type: phoneme.PhonemeTypeConsonant},
			},
			isValid: true,
		},
		{
			name: "Valid plosive + fricative",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Manner: phoneme.ConsonantMannerPlosive}},
				{Symbol: "s", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Manner: phoneme.ConsonantMannerFricative}},
			},
			isValid: true,
		},
		{
			name: "Invalid fricative + fricative",
			seq: []phoneme.Phoneme{
				{Symbol: "f", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Manner: phoneme.ConsonantMannerFricative}},
				{Symbol: "s", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Manner: phoneme.ConsonantMannerFricative}},
			},
			isValid: false,
		},
		{
			name: "Invalid fricative + sibilant",
			seq: []phoneme.Phoneme{
				{Symbol: "f", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Manner: phoneme.ConsonantMannerFricative}},
				{Symbol: "s", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Manner: phoneme.ConsonantMannerSibilants}},
			},
			isValid: false,
		},
		{
			name: "Invalid nasal + nasal",
			seq: []phoneme.Phoneme{
				{Symbol: "n", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Manner: phoneme.ConsonantMannerNasal}},
				{Symbol: "m", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{
					Manner: phoneme.ConsonantMannerNasal}},
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := rule.Validate(tc.seq)
			if result != tc.isValid {
				t.Errorf("Expected validity %v, got %v", tc.isValid, result)
			}
		})
	}
}

func TestDefaultRules(t *testing.T) {
	rules := DefaultRules()

	if len(rules) != 6 {
		t.Errorf("Expected 6 default rules, got %d", len(rules))
	}

	// Test that each rule type is present
	ruleTypes := make(map[string]bool)
	for _, rule := range rules {
		switch rule.(type) {
		case *MaxClusterSizeRule:
			ruleTypes["MaxClusterSize"] = true
		case *VelarNasalOnsetRule:
			ruleTypes["VelarNasalOnset"] = true
		case *VowelRequirementRule:
			ruleTypes["VowelRequirement"] = true
		case *SyllableLengthRule:
			ruleTypes["SyllableLength"] = true
		case *VoicingAssimilationRule:
			ruleTypes["VoicingAssimilation"] = true
		case *MannerRestrictionsRule:
			ruleTypes["MannerRestrictions"] = true
		default:
			t.Errorf("Unexpected rule type: %T", rule)
		}
	}

	expectedTypes := []string{"MaxClusterSize", "VelarNasalOnset", "VowelRequirement", "SyllableLength", "VoicingAssimilation", "MannerRestrictions"}
	for _, expectedType := range expectedTypes {
		if !ruleTypes[expectedType] {
			t.Errorf("Expected rule type %s not found", expectedType)
		}
	}
}

func TestStrictDefaultRules(t *testing.T) {
	rules := StrictDefaultRules()

	if len(rules) != 7 {
		t.Errorf("Expected 7 strict default rules, got %d", len(rules))
	}

	// Test that each rule type is present
	ruleTypes := make(map[string]bool)
	for _, rule := range rules {
		switch rule.(type) {
		case *MaxClusterSizeRule:
			ruleTypes["MaxClusterSize"] = true
		case *VelarNasalOnsetRule:
			ruleTypes["VelarNasalOnset"] = true
		case *VowelRequirementRule:
			ruleTypes["VowelRequirement"] = true
		case *SyllableLengthRule:
			ruleTypes["SyllableLength"] = true
		case *ConsonantHarmonyRule:
			ruleTypes["ConsonantHarmony"] = true
		case *VoicingAssimilationRule:
			ruleTypes["VoicingAssimilation"] = true
		case *MannerRestrictionsRule:
			ruleTypes["MannerRestrictions"] = true
		default:
			t.Errorf("Unexpected rule type: %T", rule)
		}
	}

	expectedTypes := []string{"MaxClusterSize", "VelarNasalOnset", "VowelRequirement", "SyllableLength", "ConsonantHarmony", "VoicingAssimilation", "MannerRestrictions"}
	for _, expectedType := range expectedTypes {
		if !ruleTypes[expectedType] {
			t.Errorf("Expected rule type %s not found", expectedType)
		}
	}
}

func TestRuleCombinations(t *testing.T) {
	// Test that multiple rules work together
	pool := phoneme.NewPool()
	phonology := NewPhonology(pool)

	// Add multiple rules
	phonology.AddRule(NewMaxClusterSizeRule(2))
	phonology.AddRule(NewVowelRequirementRule())
	phonology.AddRule(NewSyllableLengthRule(1, 4))

	testCases := []struct {
		name    string
		seq     []phoneme.Phoneme
		isValid bool
	}{
		{
			name: "Valid CV syllable",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
			},
			isValid: true,
		},
		{
			name: "Invalid: no vowel",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "t", Type: phoneme.PhonemeTypeConsonant},
			},
			isValid: false,
		},
		{
			name: "Invalid: too many consonants",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "t", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "p", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
			},
			isValid: false,
		},
		{
			name: "Invalid: too long",
			seq: []phoneme.Phoneme{
				{Symbol: "k", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "t", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "a", Type: phoneme.PhonemeTypeVowel},
				{Symbol: "p", Type: phoneme.PhonemeTypeConsonant},
				{Symbol: "s", Type: phoneme.PhonemeTypeConsonant},
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := phonology.ApplyRules(tc.seq)
			if result != tc.isValid {
				t.Errorf("Expected validity %v, got %v", tc.isValid, result)
			}
		})
	}
}
