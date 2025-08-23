package phoneme

import (
	"encoding/json"
	"testing"
)

func TestVowelHeight_Constants(t *testing.T) {
	// Test that constants have expected values
	if VowelHeightUnknown != 0 {
		t.Errorf("VowelHeightUnknown should be 0, got %d", VowelHeightUnknown)
	}
	if VowelHeightLow != 1 {
		t.Errorf("VowelHeightLow should be 1, got %d", VowelHeightLow)
	}
	if VowelHeightMid != 2 {
		t.Errorf("VowelHeightMid should be 2, got %d", VowelHeightMid)
	}
	if VowelHeightHigh != 3 {
		t.Errorf("VowelHeightHigh should be 3, got %d", VowelHeightHigh)
	}
}

func TestVowelHeight_String(t *testing.T) {
	tests := []struct {
		name     string
		height   VowelHeight
		expected string
	}{
		{"unknown", VowelHeightUnknown, "unknown"},
		{"low", VowelHeightLow, "low"},
		{"mid", VowelHeightMid, "mid"},
		{"high", VowelHeightHigh, "high"},
		{"out of bounds", VowelHeight(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.height.String()
			if result != tt.expected {
				t.Errorf("VowelHeight(%d).String() = %s, want %s", tt.height, result, tt.expected)
			}
		})
	}
}

func TestParseVowelHeight(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected VowelHeight
	}{
		{"low", "low", VowelHeightLow},
		{"mid", "mid", VowelHeightMid},
		{"high", "high", VowelHeightHigh},
		{"unknown", "unknown", VowelHeightUnknown},
		{"LOW", "LOW", VowelHeightUnknown}, // Case sensitive
		{"Low", "Low", VowelHeightUnknown}, // Case sensitive
		{"", "", VowelHeightUnknown},
		{"invalid", "invalid", VowelHeightUnknown},
		{"low ", "low ", VowelHeightUnknown}, // Extra space
		{" low", " low", VowelHeightUnknown}, // Leading space
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseVowelHeight(tt.input)
			if result != tt.expected {
				t.Errorf("ParseVowelHeight(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestVowelHeight_UnmarshalText(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expected    VowelHeight
		expectError bool
	}{
		{"low", []byte("low"), VowelHeightLow, false},
		{"mid", []byte("mid"), VowelHeightMid, false},
		{"high", []byte("high"), VowelHeightHigh, false},
		{"unknown", []byte("unknown"), VowelHeightUnknown, false},
		{"invalid", []byte("invalid"), VowelHeightUnknown, true},
		{"empty", []byte(""), VowelHeightUnknown, true},
		{"case sensitive", []byte("LOW"), VowelHeightUnknown, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var vh VowelHeight
			err := vh.UnmarshalText(tt.input)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if vh != tt.expected {
				t.Errorf("UnmarshalText(%q) = %v, want %v", tt.input, vh, tt.expected)
			}
		})
	}
}

func TestVowelHeight_JSON(t *testing.T) {
	// Test JSON marshaling
	vh := VowelHeightMid
	data, err := json.Marshal(vh)
	if err != nil {
		t.Fatalf("Failed to marshal VowelHeight: %v", err)
	}

	// Should marshal to string representation
	expected := `"mid"`
	if string(data) != expected {
		t.Errorf("JSON marshaling failed: got %s, want %s", string(data), expected)
	}

	// Test JSON unmarshaling
	var unmarshaled VowelHeight
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal VowelHeight: %v", err)
	}

	if unmarshaled != vh {
		t.Errorf("JSON unmarshaling failed: got %v, want %v", unmarshaled, vh)
	}
}

func TestVowelHeight_JSONUnmarshalError(t *testing.T) {
	// Test invalid JSON unmarshaling
	invalidJSON := `"invalid_height"`
	var vh VowelHeight
	err := json.Unmarshal([]byte(invalidJSON), &vh)

	if err == nil {
		t.Error("Expected error for invalid vowel height")
	}

	if vh != VowelHeightUnknown {
		t.Errorf("Expected VowelHeightUnknown for invalid input, got %v", vh)
	}
}

func TestVowelBackness_Constants(t *testing.T) {
	// Test that constants have expected values
	if VowelBacknessUnknown != 0 {
		t.Errorf("VowelBacknessUnknown should be 0, got %d", VowelBacknessUnknown)
	}
	if VowelBacknessFront != 1 {
		t.Errorf("VowelBacknessFront should be 1, got %d", VowelBacknessFront)
	}
	if VowelBacknessCentral != 2 {
		t.Errorf("VowelBacknessCentral should be 2, got %d", VowelBacknessCentral)
	}
	if VowelBacknessBack != 3 {
		t.Errorf("VowelBacknessBack should be 3, got %d", VowelBacknessBack)
	}
}

func TestVowelBackness_String(t *testing.T) {
	tests := []struct {
		name     string
		backness VowelBackness
		expected string
	}{
		{"unknown", VowelBacknessUnknown, "unknown"},
		{"front", VowelBacknessFront, "front"},
		{"central", VowelBacknessCentral, "central"},
		{"back", VowelBacknessBack, "back"},
		{"out of bounds", VowelBackness(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.backness.String()
			if result != tt.expected {
				t.Errorf("VowelBackness(%d).String() = %s, want %s", tt.backness, result, tt.expected)
			}
		})
	}
}

func TestParseVowelBackness(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected VowelBackness
	}{
		{"front", "front", VowelBacknessFront},
		{"central", "central", VowelBacknessCentral},
		{"back", "back", VowelBacknessBack},
		{"unknown", "unknown", VowelBacknessUnknown},
		{"FRONT", "FRONT", VowelBacknessUnknown}, // Case sensitive
		{"Front", "Front", VowelBacknessUnknown}, // Case sensitive
		{"", "", VowelBacknessUnknown},
		{"invalid", "invalid", VowelBacknessUnknown},
		{"front ", "front ", VowelBacknessUnknown}, // Extra space
		{" front", " front", VowelBacknessUnknown}, // Leading space
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseVowelBackness(tt.input)
			if result != tt.expected {
				t.Errorf("ParseVowelBackness(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestVowelBackness_UnmarshalText(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expected    VowelBackness
		expectError bool
	}{
		{"front", []byte("front"), VowelBacknessFront, false},
		{"central", []byte("central"), VowelBacknessCentral, false},
		{"back", []byte("back"), VowelBacknessBack, false},
		{"unknown", []byte("unknown"), VowelBacknessUnknown, false},
		{"invalid", []byte("invalid"), VowelBacknessUnknown, true},
		{"empty", []byte(""), VowelBacknessUnknown, true},
		{"case sensitive", []byte("FRONT"), VowelBacknessUnknown, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var vb VowelBackness
			err := vb.UnmarshalText(tt.input)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if vb != tt.expected {
				t.Errorf("UnmarshalText(%q) = %v, want %v", tt.input, vb, tt.expected)
			}
		})
	}
}

func TestVowelBackness_JSON(t *testing.T) {
	// Test JSON marshaling
	vb := VowelBacknessBack
	data, err := json.Marshal(vb)
	if err != nil {
		t.Fatalf("Failed to marshal VowelBackness: %v", err)
	}

	// Should marshal to string representation
	expected := `"back"`
	if string(data) != expected {
		t.Errorf("JSON marshaling failed: got %s, want %s", string(data), expected)
	}

	// Test JSON unmarshaling
	var unmarshaled VowelBackness
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal VowelBackness: %v", err)
	}

	if unmarshaled != vb {
		t.Errorf("JSON unmarshaling failed: got %v, want %v", unmarshaled, vb)
	}
}

func TestVowelBackness_JSONUnmarshalError(t *testing.T) {
	// Test invalid JSON unmarshaling
	invalidJSON := `"invalid_backness"`
	var vb VowelBackness
	err := json.Unmarshal([]byte(invalidJSON), &vb)

	if err == nil {
		t.Error("Expected error for invalid vowel backness")
	}

	if vb != VowelBacknessUnknown {
		t.Errorf("Expected VowelBacknessUnknown for invalid input, got %v", vb)
	}
}

func TestVowelSpec_JSON(t *testing.T) {
	spec := VowelSpec{
		Height:    VowelHeightHigh,
		Backness:  VowelBacknessFront,
		Rounded:   true,
		Diphthong: false,
	}

	// Test JSON marshaling
	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("Failed to marshal VowelSpec: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled VowelSpec
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal VowelSpec: %v", err)
	}

	if unmarshaled.Height != spec.Height {
		t.Errorf("Height mismatch: got %v, want %v", unmarshaled.Height, spec.Height)
	}
	if unmarshaled.Backness != spec.Backness {
		t.Errorf("Backness mismatch: got %v, want %v", unmarshaled.Backness, spec.Backness)
	}
	if unmarshaled.Rounded != spec.Rounded {
		t.Errorf("Rounded mismatch: got %t, want %t", unmarshaled.Rounded, spec.Rounded)
	}
	if unmarshaled.Diphthong != spec.Diphthong {
		t.Errorf("Diphthong mismatch: got %t, want %t", unmarshaled.Diphthong, spec.Diphthong)
	}
}

func TestVowelSpec_EdgeCases(t *testing.T) {
	// Test with all boolean combinations
	specs := []VowelSpec{
		{Height: VowelHeightLow, Backness: VowelBacknessFront, Rounded: false, Diphthong: false},
		{Height: VowelHeightMid, Backness: VowelBacknessCentral, Rounded: true, Diphthong: true},
		{Height: VowelHeightHigh, Backness: VowelBacknessBack, Rounded: false, Diphthong: true},
		{Height: VowelHeightMid, Backness: VowelBacknessFront, Rounded: true, Diphthong: false},
	}

	for i, spec := range specs {
		data, err := json.Marshal(spec)
		if err != nil {
			t.Fatalf("Failed to marshal VowelSpec %d: %v", i, err)
		}

		var unmarshaled VowelSpec
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal VowelSpec %d: %v", i, err)
		}

		if unmarshaled.Height != spec.Height {
			t.Errorf("Spec %d: Height mismatch: got %v, want %v", i, unmarshaled.Height, spec.Height)
		}
		if unmarshaled.Backness != spec.Backness {
			t.Errorf("Spec %d: Backness mismatch: got %v, want %v", i, unmarshaled.Backness, spec.Backness)
		}
		if unmarshaled.Rounded != spec.Rounded {
			t.Errorf("Spec %d: Rounded mismatch: got %t, want %t", i, unmarshaled.Rounded, spec.Rounded)
		}
		if unmarshaled.Diphthong != spec.Diphthong {
			t.Errorf("Spec %d: Diphthong mismatch: got %t, want %t", i, unmarshaled.Diphthong, spec.Diphthong)
		}
	}
}

func TestVowelHeight_EdgeCases(t *testing.T) {
	// Test boundary conditions
	vh := VowelHeight(0)
	if vh.String() != "unknown" {
		t.Errorf("VowelHeight(0).String() = %s, want 'unknown'", vh.String())
	}

	vh = VowelHeight(3)
	if vh.String() != "high" {
		t.Errorf("VowelHeight(3).String() = %s, want 'high'", vh.String())
	}

	// Test out of bounds
	vh = VowelHeight(255)
	if vh.String() != "unknown" {
		t.Errorf("VowelHeight(255).String() = %s, want 'unknown'", vh.String())
	}
}

func TestVowelBackness_EdgeCases(t *testing.T) {
	// Test boundary conditions
	vb := VowelBackness(0)
	if vb.String() != "unknown" {
		t.Errorf("VowelBackness(0).String() = %s, want 'unknown'", vb.String())
	}

	vb = VowelBackness(3)
	if vb.String() != "back" {
		t.Errorf("VowelBackness(3).String() = %s, want 'back'", vb.String())
	}

	// Test out of bounds
	vb = VowelBackness(255)
	if vb.String() != "unknown" {
		t.Errorf("VowelBackness(255).String() = %s, want 'unknown'", vb.String())
	}
}

func TestVowelEnums_Consistency(t *testing.T) {
	// Test that parsing and string conversion are consistent for VowelHeight
	heights := []VowelHeight{VowelHeightUnknown, VowelHeightLow, VowelHeightMid, VowelHeightHigh}

	for _, vh := range heights {
		str := vh.String()
		parsed := ParseVowelHeight(str)
		if parsed != vh {
			t.Errorf("VowelHeight inconsistency: %v.String() = %q, ParseVowelHeight(%q) = %v",
				vh, str, str, parsed)
		}
	}

	// Test that parsing and string conversion are consistent for VowelBackness
	backnesses := []VowelBackness{VowelBacknessUnknown, VowelBacknessFront, VowelBacknessCentral, VowelBacknessBack}

	for _, vb := range backnesses {
		str := vb.String()
		parsed := ParseVowelBackness(str)
		if parsed != vb {
			t.Errorf("VowelBackness inconsistency: %v.String() = %q, ParseVowelBackness(%q) = %v",
				vb, str, str, parsed)
		}
	}
}
