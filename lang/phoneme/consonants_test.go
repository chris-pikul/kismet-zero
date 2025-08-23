package phoneme

import (
	"encoding/json"
	"testing"
)

func TestConsonantPlace(t *testing.T) {
	testCases := []struct {
		name     string
		place    ConsonantPlace
		expected string
	}{
		{"Unknown", ConsonantPlaceUnknown, "unknown"},
		{"Bilabial", ConsonantPlaceBilabial, "bilabial"},
		{"Labiodental", ConsonantPlaceLabiodental, "labiodental"},
		{"Dental", ConsonantPlaceDental, "dental"},
		{"Alveolar", ConsonantPlaceAlveolar, "alveolar"},
		{"Postalveolar", ConsonantPlacePostalveolar, "postalveolar"},
		{"Retroflex", ConsonantPlaceRetroflex, "retroflex"},
		{"Palatal", ConsonantPlacePalatal, "palatal"},
		{"Velar", ConsonantPlaceVelar, "velar"},
		{"Uvular", ConsonantPlaceUvular, "uvular"},
		{"Pharyngeal", ConsonantPlacePharyngeal, "pharyngeal"},
		{"Glottal", ConsonantPlaceGlottal, "glottal"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.place.String() != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, tc.place.String())
			}
		})
	}
}

func TestConsonantPlaceString_Invalid(t *testing.T) {
	invalidPlace := ConsonantPlace(255)
	if invalidPlace.String() != "unknown" {
		t.Errorf("Expected 'unknown' for invalid place, got %s", invalidPlace.String())
	}
}

func TestParseConsonantPlace(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected ConsonantPlace
	}{
		{"Valid bilabial", "bilabial", ConsonantPlaceBilabial},
		{"Valid velar", "velar", ConsonantPlaceVelar},
		{"Valid alveolar", "alveolar", ConsonantPlaceAlveolar},
		{"Invalid input", "invalid", ConsonantPlaceUnknown},
		{"Empty string", "", ConsonantPlaceUnknown},
		{"Case sensitive", "BILABIAL", ConsonantPlaceUnknown},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseConsonantPlace(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestConsonantPlaceJSON(t *testing.T) {
	testCases := []struct {
		name     string
		place    ConsonantPlace
		expected string
	}{
		{"Bilabial", ConsonantPlaceBilabial, `"bilabial"`},
		{"Velar", ConsonantPlaceVelar, `"velar"`},
		{"Unknown", ConsonantPlaceUnknown, `"unknown"`},
	}

	for _, tc := range testCases {
		t.Run(tc.name+"_Marshal", func(t *testing.T) {
			data, err := json.Marshal(tc.place)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if string(data) != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, string(data))
			}
		})

		t.Run(tc.name+"_Unmarshal", func(t *testing.T) {
			var place ConsonantPlace
			err := json.Unmarshal([]byte(tc.expected), &place)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if place != tc.place {
				t.Errorf("Expected %v, got %v", tc.place, place)
			}
		})
	}
}

func TestConsonantPlaceUnmarshalJSON_Invalid(t *testing.T) {
	testCases := []string{
		`"invalid"`,
		`"BILABIAL"`,
		`"notaplace"`,
	}

	for _, tc := range testCases {
		t.Run("Invalid_"+tc, func(t *testing.T) {
			var place ConsonantPlace
			err := json.Unmarshal([]byte(tc), &place)
			if err == nil {
				t.Errorf("Expected error for invalid input %s", tc)
			}
		})
	}
}

func TestConsonantVoicing(t *testing.T) {
	testCases := []struct {
		name     string
		voicing  ConsonantVoicing
		expected string
	}{
		{"Unknown", ConsonantVoicingUnknown, "unknown"},
		{"Voiced", ConsonantVoicingVoiced, "voiced"},
		{"Unvoiced", ConsonantVoicingUnvoiced, "unvoiced"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.voicing.String() != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, tc.voicing.String())
			}
		})
	}
}

func TestConsonantVoicingString_Invalid(t *testing.T) {
	invalidVoicing := ConsonantVoicing(255)
	if invalidVoicing.String() != "unknown" {
		t.Errorf("Expected 'unknown' for invalid voicing, got %s", invalidVoicing.String())
	}
}

func TestParseConsonantVoicing(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected ConsonantVoicing
	}{
		{"Valid voiced", "voiced", ConsonantVoicingVoiced},
		{"Valid unvoiced", "unvoiced", ConsonantVoicingUnvoiced},
		{"Valid unknown", "unknown", ConsonantVoicingUnknown},
		{"Invalid input", "invalid", ConsonantVoicingUnknown},
		{"Empty string", "", ConsonantVoicingUnknown},
		{"Case sensitive", "VOICED", ConsonantVoicingUnknown},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseConsonantVoicing(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestConsonantVoicingJSON(t *testing.T) {
	testCases := []struct {
		name     string
		voicing  ConsonantVoicing
		expected string
	}{
		{"Voiced", ConsonantVoicingVoiced, `"voiced"`},
		{"Unvoiced", ConsonantVoicingUnvoiced, `"unvoiced"`},
		{"Unknown", ConsonantVoicingUnknown, `"unknown"`},
	}

	for _, tc := range testCases {
		t.Run(tc.name+"_Marshal", func(t *testing.T) {
			data, err := json.Marshal(tc.voicing)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if string(data) != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, string(data))
			}
		})

		t.Run(tc.name+"_Unmarshal", func(t *testing.T) {
			var voicing ConsonantVoicing
			err := json.Unmarshal([]byte(tc.expected), &voicing)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if voicing != tc.voicing {
				t.Errorf("Expected %v, got %v", tc.voicing, voicing)
			}
		})
	}
}

func TestConsonantVoicingUnmarshalJSON_Invalid(t *testing.T) {
	testCases := []string{
		`"invalid"`,
		`"VOICED"`,
		`"notavoicing"`,
	}

	for _, tc := range testCases {
		t.Run("Invalid_"+tc, func(t *testing.T) {
			var voicing ConsonantVoicing
			err := json.Unmarshal([]byte(tc), &voicing)
			if err == nil {
				t.Errorf("Expected error for invalid input %s", tc)
			}
		})
	}
}

func TestConsonantSpecComplete(t *testing.T) {
	// Test the updated ConsonantSpec with all fields
	spec := ConsonantSpec{
		Manner:  ConsonantMannerPlosive,
		Place:   ConsonantPlaceVelar,
		Voicing: ConsonantVoicingUnvoiced,
	}

	// Test JSON marshaling and unmarshaling
	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("Unexpected error marshaling: %v", err)
	}

	var unmarshaled ConsonantSpec
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Unexpected error unmarshaling: %v", err)
	}

	if unmarshaled.Manner != spec.Manner {
		t.Errorf("Expected manner %v, got %v", spec.Manner, unmarshaled.Manner)
	}
	if unmarshaled.Place != spec.Place {
		t.Errorf("Expected place %v, got %v", spec.Place, unmarshaled.Place)
	}
	if unmarshaled.Voicing != spec.Voicing {
		t.Errorf("Expected voicing %v, got %v", spec.Voicing, unmarshaled.Voicing)
	}
}

func TestConsonantSpecExamples(t *testing.T) {
	// Test some realistic consonant specifications
	testCases := []struct {
		name   string
		spec   ConsonantSpec
		symbol string // IPA symbol for reference
	}{
		{
			name: "k - voiceless velar plosive",
			spec: ConsonantSpec{
				Manner:  ConsonantMannerPlosive,
				Place:   ConsonantPlaceVelar,
				Voicing: ConsonantVoicingUnvoiced,
			},
			symbol: "k",
		},
		{
			name: "b - voiced bilabial plosive",
			spec: ConsonantSpec{
				Manner:  ConsonantMannerPlosive,
				Place:   ConsonantPlaceBilabial,
				Voicing: ConsonantVoicingVoiced,
			},
			symbol: "b",
		},
		{
			name: "s - voiceless alveolar fricative",
			spec: ConsonantSpec{
				Manner:  ConsonantMannerFricative,
				Place:   ConsonantPlaceAlveolar,
				Voicing: ConsonantVoicingUnvoiced,
			},
			symbol: "s",
		},
		{
			name: "n - voiced alveolar nasal",
			spec: ConsonantSpec{
				Manner:  ConsonantMannerNasal,
				Place:   ConsonantPlaceAlveolar,
				Voicing: ConsonantVoicingVoiced,
			},
			symbol: "n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Verify the spec is well-formed
			if tc.spec.Manner == ConsonantMannerUnknown {
				t.Error("Manner should not be unknown")
			}
			if tc.spec.Place == ConsonantPlaceUnknown {
				t.Error("Place should not be unknown")
			}
			if tc.spec.Voicing == ConsonantVoicingUnknown {
				t.Error("Voicing should not be unknown")
			}

			// Test JSON round-trip
			data, err := json.Marshal(tc.spec)
			if err != nil {
				t.Fatalf("Unexpected error marshaling: %v", err)
			}

			var unmarshaled ConsonantSpec
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Fatalf("Unexpected error unmarshaling: %v", err)
			}

			if unmarshaled != tc.spec {
				t.Errorf("JSON round-trip failed: expected %+v, got %+v", tc.spec, unmarshaled)
			}
		})
	}
}
