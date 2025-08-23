package phoneme

import (
	"encoding/json"
	"testing"
)

func TestPhonemeType_Constants(t *testing.T) {
	// Test that constants have expected values
	if PhonemeTypeUnknown != 0 {
		t.Errorf("PhonemeTypeUnknown should be 0, got %d", PhonemeTypeUnknown)
	}
	if PhonemeTypeConsonant != 1 {
		t.Errorf("PhonemeTypeConsonant should be 1, got %d", PhonemeTypeConsonant)
	}
	if PhonemeTypeVowel != 2 {
		t.Errorf("PhonemeTypeVowel should be 2, got %d", PhonemeTypeVowel)
	}
}

func TestPhonemeType_String(t *testing.T) {
	tests := []struct {
		name        string
		phonemeType PhonemeType
		expected    string
	}{
		{"unknown", PhonemeTypeUnknown, "unknown"},
		{"consonant", PhonemeTypeConsonant, "consonant"},
		{"vowel", PhonemeTypeVowel, "vowel"},
		{"out of bounds", PhonemeType(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.phonemeType.String()
			if result != tt.expected {
				t.Errorf("PhonemeType(%d).String() = %s, want %s", tt.phonemeType, result, tt.expected)
			}
		})
	}
}

func TestParsePhonemeType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected PhonemeType
	}{
		{"consonant", "consonant", PhonemeTypeConsonant},
		{"vowel", "vowel", PhonemeTypeVowel},
		{"unknown", "unknown", PhonemeTypeUnknown},
		{"CONSONANT", "CONSONANT", PhonemeTypeUnknown}, // Case sensitive
		{"Vowel", "Vowel", PhonemeTypeUnknown},         // Case sensitive
		{"", "", PhonemeTypeUnknown},
		{"invalid", "invalid", PhonemeTypeUnknown},
		{"consonant ", "consonant ", PhonemeTypeUnknown}, // Extra space
		{" consonant", " consonant", PhonemeTypeUnknown}, // Leading space
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParsePhonemeType(tt.input)
			if result != tt.expected {
				t.Errorf("ParsePhonemeType(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPhonemeType_UnmarshalText(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expected    PhonemeType
		expectError bool
	}{
		{"consonant", []byte("consonant"), PhonemeTypeConsonant, false},
		{"vowel", []byte("vowel"), PhonemeTypeVowel, false},
		{"unknown", []byte("unknown"), PhonemeTypeUnknown, false},
		{"invalid", []byte("invalid"), PhonemeTypeUnknown, true},
		{"empty", []byte(""), PhonemeTypeUnknown, true},
		{"case sensitive", []byte("CONSONANT"), PhonemeTypeUnknown, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pt PhonemeType
			err := pt.UnmarshalText(tt.input)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if pt != tt.expected {
				t.Errorf("UnmarshalText(%q) = %v, want %v", tt.input, pt, tt.expected)
			}
		})
	}
}

func TestPhonemeType_JSON(t *testing.T) {
	// Test JSON marshaling
	pt := PhonemeTypeConsonant
	data, err := json.Marshal(pt)
	if err != nil {
		t.Fatalf("Failed to marshal PhonemeType: %v", err)
	}

	// Should marshal to string representation
	expected := `"consonant"`
	if string(data) != expected {
		t.Errorf("JSON marshaling failed: got %s, want %s", string(data), expected)
	}

	// Test JSON unmarshaling
	var unmarshaled PhonemeType
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal PhonemeType: %v", err)
	}

	if unmarshaled != pt {
		t.Errorf("JSON unmarshaling failed: got %v, want %v", unmarshaled, pt)
	}
}

func TestPhonemeType_JSONUnmarshalError(t *testing.T) {
	// Test invalid JSON unmarshaling
	invalidJSON := `"invalid_type"`
	var pt PhonemeType
	err := json.Unmarshal([]byte(invalidJSON), &pt)

	if err == nil {
		t.Error("Expected error for invalid phoneme type")
	}

	if pt != PhonemeTypeUnknown {
		t.Errorf("Expected PhonemeTypeUnknown for invalid input, got %v", pt)
	}
}

func TestPhonemeType_EdgeCases(t *testing.T) {
	// Test boundary conditions
	pt := PhonemeType(0)
	if pt.String() != "unknown" {
		t.Errorf("PhonemeType(0).String() = %s, want 'unknown'", pt.String())
	}

	pt = PhonemeType(1)
	if pt.String() != "consonant" {
		t.Errorf("PhonemeType(1).String() = %s, want 'consonant'", pt.String())
	}

	pt = PhonemeType(2)
	if pt.String() != "vowel" {
		t.Errorf("PhonemeType(2).String() = %s, want 'vowel'", pt.String())
	}

	// Test out of bounds
	pt = PhonemeType(255)
	if pt.String() != "unknown" {
		t.Errorf("PhonemeType(255).String() = %s, want 'unknown'", pt.String())
	}
}

func TestPhonemeType_Consistency(t *testing.T) {
	// Test that parsing and string conversion are consistent
	types := []PhonemeType{PhonemeTypeUnknown, PhonemeTypeConsonant, PhonemeTypeVowel}

	for _, pt := range types {
		str := pt.String()
		parsed := ParsePhonemeType(str)
		if parsed != pt {
			t.Errorf("Inconsistency: %v.String() = %q, ParsePhonemeType(%q) = %v",
				pt, str, str, parsed)
		}
	}
}
