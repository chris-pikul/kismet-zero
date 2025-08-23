package phoneme

import (
	"encoding/json"
	"testing"
)

func TestConsonantManner_Constants(t *testing.T) {
	// Test that constants have expected values
	if ConsonantMannerUnknown != 0 {
		t.Errorf("ConsonantMannerUnknown should be 0, got %d", ConsonantMannerUnknown)
	}
	if ConsonantMannerPlosive != 1 {
		t.Errorf("ConsonantMannerPlosive should be 1, got %d", ConsonantMannerPlosive)
	}
	if ConsonantMannerNasal != 2 {
		t.Errorf("ConsonantMannerNasal should be 2, got %d", ConsonantMannerNasal)
	}
	if ConsonantMannerFricative != 3 {
		t.Errorf("ConsonantMannerFricative should be 3, got %d", ConsonantMannerFricative)
	}
	if ConsonantMannerSibilants != 4 {
		t.Errorf("ConsonantMannerSibilants should be 4, got %d", ConsonantMannerSibilants)
	}
	if ConsonantMannerLateralFricative != 5 {
		t.Errorf("ConsonantMannerLateralFricative should be 5, got %d", ConsonantMannerLateralFricative)
	}
	if ConsonantMannerAffricate != 6 {
		t.Errorf("ConsonantMannerAffricate should be 6, got %d", ConsonantMannerAffricate)
	}
	if ConsonantMannerVibrant != 7 {
		t.Errorf("ConsonantMannerVibrant should be 7, got %d", ConsonantMannerVibrant)
	}
	if ConsonantMannerFlap != 8 {
		t.Errorf("ConsonantMannerFlap should be 8, got %d", ConsonantMannerFlap)
	}
	if ConsonantMannerTrill != 9 {
		t.Errorf("ConsonantMannerTrill should be 9, got %d", ConsonantMannerTrill)
	}
	if ConsonantMannerApproximant != 10 {
		t.Errorf("ConsonantMannerApproximant should be 10, got %d", ConsonantMannerApproximant)
	}
	if ConsonantMannerSemivowel != 11 {
		t.Errorf("ConsonantMannerSemivowel should be 11, got %d", ConsonantMannerSemivowel)
	}
	if ConsonantMannerDiphthong != 12 {
		t.Errorf("ConsonantMannerDiphthong should be 12, got %d", ConsonantMannerDiphthong)
	}
	if ConsonantMannerLateralApproximant != 13 {
		t.Errorf("ConsonantMannerLateralApproximant should be 13, got %d", ConsonantMannerLateralApproximant)
	}
	if ConsonantMannerEjective != 14 {
		t.Errorf("ConsonantMannerEjective should be 14, got %d", ConsonantMannerEjective)
	}
	if ConsonantMannerImplosive != 15 {
		t.Errorf("ConsonantMannerImplosive should be 15, got %d", ConsonantMannerImplosive)
	}
	if ConsonantMannerClick != 16 {
		t.Errorf("ConsonantMannerClick should be 16, got %d", ConsonantMannerClick)
	}
	if ConsonantMannerPercussive != 17 {
		t.Errorf("ConsonantMannerPercussive should be 17, got %d", ConsonantMannerPercussive)
	}
}

func TestConsonantManner_String(t *testing.T) {
	tests := []struct {
		name     string
		manner   ConsonantManner
		expected string
	}{
		{"unknown", ConsonantMannerUnknown, "unknown"},
		{"plosive", ConsonantMannerPlosive, "plosive"},
		{"nasal", ConsonantMannerNasal, "nasal"},
		{"fricative", ConsonantMannerFricative, "fricative"},
		{"sibilants", ConsonantMannerSibilants, "sibilants"},
		{"lateral-fricative", ConsonantMannerLateralFricative, "lateral-fricative"},
		{"affricate", ConsonantMannerAffricate, "affricate"},
		{"vibrant", ConsonantMannerVibrant, "vibrant"},
		{"flap", ConsonantMannerFlap, "flap"},
		{"trill", ConsonantMannerTrill, "trill"},
		{"approximant", ConsonantMannerApproximant, "approximant"},
		{"semivowel", ConsonantMannerSemivowel, "semivowel"},
		{"diphthong", ConsonantMannerDiphthong, "diphthong"},
		{"lateral-approximant", ConsonantMannerLateralApproximant, "lateral-approximant"},
		{"ejective", ConsonantMannerEjective, "ejective"},
		{"implosive", ConsonantMannerImplosive, "implosive"},
		{"click", ConsonantMannerClick, "click"},
		{"percussive", ConsonantMannerPercussive, "percussive"},
		{"out of bounds", ConsonantManner(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.manner.String()
			if result != tt.expected {
				t.Errorf("ConsonantManner(%d).String() = %s, want %s", tt.manner, result, tt.expected)
			}
		})
	}
}

func TestParseConsonantManner(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ConsonantManner
	}{
		{"plosive", "plosive", ConsonantMannerPlosive},
		{"nasal", "nasal", ConsonantMannerNasal},
		{"fricative", "fricative", ConsonantMannerFricative},
		{"sibilants", "sibilants", ConsonantMannerSibilants},
		{"lateral-fricative", "lateral-fricative", ConsonantMannerLateralFricative},
		{"affricate", "affricate", ConsonantMannerAffricate},
		{"vibrant", "vibrant", ConsonantMannerVibrant},
		{"flap", "flap", ConsonantMannerFlap},
		{"trill", "trill", ConsonantMannerTrill},
		{"approximant", "approximant", ConsonantMannerApproximant},
		{"semivowel", "semivowel", ConsonantMannerSemivowel},
		{"diphthong", "diphthong", ConsonantMannerDiphthong},
		{"lateral-approximant", "lateral-approximant", ConsonantMannerLateralApproximant},
		{"ejective", "ejective", ConsonantMannerEjective},
		{"implosive", "implosive", ConsonantMannerImplosive},
		{"click", "click", ConsonantMannerClick},
		{"percussive", "percussive", ConsonantMannerPercussive},
		{"unknown", "unknown", ConsonantMannerUnknown},
		{"PLOSIVE", "PLOSIVE", ConsonantMannerUnknown}, // Case sensitive
		{"Plosive", "Plosive", ConsonantMannerUnknown}, // Case sensitive
		{"", "", ConsonantMannerUnknown},
		{"invalid", "invalid", ConsonantMannerUnknown},
		{"plosive ", "plosive ", ConsonantMannerUnknown}, // Extra space
		{" plosive", " plosive", ConsonantMannerUnknown}, // Leading space
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseConsonantManner(tt.input)
			if result != tt.expected {
				t.Errorf("ParseConsonantManner(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConsonantManner_UnmarshalText(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expected    ConsonantManner
		expectError bool
	}{
		{"plosive", []byte("plosive"), ConsonantMannerPlosive, false},
		{"nasal", []byte("nasal"), ConsonantMannerNasal, false},
		{"fricative", []byte("fricative"), ConsonantMannerFricative, false},
		{"unknown", []byte("unknown"), ConsonantMannerUnknown, false},
		{"invalid", []byte("invalid"), ConsonantMannerUnknown, true},
		{"empty", []byte(""), ConsonantMannerUnknown, true},
		{"case sensitive", []byte("PLOSIVE"), ConsonantMannerUnknown, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cm ConsonantManner
			err := cm.UnmarshalText(tt.input)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if cm != tt.expected {
				t.Errorf("UnmarshalText(%q) = %v, want %v", tt.input, cm, tt.expected)
			}
		})
	}
}

func TestConsonantManner_JSON(t *testing.T) {
	// Test JSON marshaling
	cm := ConsonantMannerPlosive
	data, err := json.Marshal(cm)
	if err != nil {
		t.Fatalf("Failed to marshal ConsonantManner: %v", err)
	}

	// Should marshal to string representation
	expected := `"plosive"`
	if string(data) != expected {
		t.Errorf("JSON marshaling failed: got %s, want %s", string(data), expected)
	}

	// Test JSON unmarshaling
	var unmarshaled ConsonantManner
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ConsonantManner: %v", err)
	}

	if unmarshaled != cm {
		t.Errorf("JSON unmarshaling failed: got %v, want %v", unmarshaled, cm)
	}
}

func TestConsonantManner_JSONUnmarshalError(t *testing.T) {
	// Test invalid JSON unmarshaling
	invalidJSON := `"invalid_manner"`
	var cm ConsonantManner
	err := json.Unmarshal([]byte(invalidJSON), &cm)

	if err == nil {
		t.Error("Expected error for invalid consonant manner")
	}

	if cm != ConsonantMannerUnknown {
		t.Errorf("Expected ConsonantMannerUnknown for invalid input, got %v", cm)
	}
}

func TestConsonantSpec_JSON(t *testing.T) {
	spec := ConsonantSpec{
		Manner: ConsonantMannerNasal,
	}

	// Test JSON marshaling
	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("Failed to marshal ConsonantSpec: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled ConsonantSpec
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ConsonantSpec: %v", err)
	}

	if unmarshaled.Manner != spec.Manner {
		t.Errorf("Manner mismatch: got %v, want %v", unmarshaled.Manner, spec.Manner)
	}
}

func TestConsonantManner_EdgeCases(t *testing.T) {
	// Test boundary conditions
	cm := ConsonantManner(0)
	if cm.String() != "unknown" {
		t.Errorf("ConsonantManner(0).String() = %s, want 'unknown'", cm.String())
	}

	cm = ConsonantManner(17)
	if cm.String() != "percussive" {
		t.Errorf("ConsonantManner(17).String() = %s, want 'percussive'", cm.String())
	}

	// Test out of bounds
	cm = ConsonantManner(255)
	if cm.String() != "unknown" {
		t.Errorf("ConsonantManner(255).String() = %s, want 'unknown'", cm.String())
	}
}

func TestConsonantManner_Consistency(t *testing.T) {
	// Test that parsing and string conversion are consistent
	manners := []ConsonantManner{
		ConsonantMannerUnknown, ConsonantMannerPlosive, ConsonantMannerNasal,
		ConsonantMannerFricative, ConsonantMannerSibilants, ConsonantMannerLateralFricative,
		ConsonantMannerAffricate, ConsonantMannerVibrant, ConsonantMannerFlap,
		ConsonantMannerTrill, ConsonantMannerApproximant, ConsonantMannerSemivowel,
		ConsonantMannerDiphthong, ConsonantMannerLateralApproximant, ConsonantMannerEjective,
		ConsonantMannerImplosive, ConsonantMannerClick, ConsonantMannerPercussive,
	}

	for _, cm := range manners {
		str := cm.String()
		parsed := ParseConsonantManner(str)
		if parsed != cm {
			t.Errorf("Inconsistency: %v.String() = %q, ParseConsonantManner(%q) = %v",
				cm, str, str, parsed)
		}
	}
}
