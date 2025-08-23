package lang

import (
	"testing"
)

func TestLanguageID_String(t *testing.T) {
	tests := []struct {
		name     string
		id       LanguageID
		expected string
	}{
		{
			name: "full ID",
			id: LanguageID{
				Family:   "elv",
				Branch:   "wood",
				Language: "sindarin",
				Dialect:  "noldorin",
				Script:   "tengwar",
				Region:   "middle-earth",
			},
			expected: "elv-wood-sindarin-noldorin-tengwar-middle-earth",
		},
		{
			name: "minimal ID",
			id: LanguageID{
				Family: "hum",
			},
			expected: "hum",
		},
		{
			name: "partial ID",
			id: LanguageID{
				Family:   "orc",
				Branch:   "black",
				Language: "speech",
			},
			expected: "orc-black-speech",
		},
		{
			name:     "empty ID",
			id:       LanguageID{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.id.String()
			if result != tt.expected {
				t.Errorf("LanguageID.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseLanguageID(t *testing.T) {
	tests := []struct {
		name        string
		idStr       string
		expected    LanguageID
		expectError bool
	}{
		{
			name:  "full ID",
			idStr: "elv-wood-sindarin-noldorin-tengwar-middleearth",
			expected: LanguageID{
				Family:   "elv",
				Branch:   "wood",
				Language: "sindarin",
				Dialect:  "noldorin",
				Script:   "tengwar",
				Region:   "middleearth",
			},
			expectError: false,
		},
		{
			name:  "minimal ID",
			idStr: "hum",
			expected: LanguageID{
				Family: "hum",
			},
			expectError: false,
		},
		{
			name:  "partial ID",
			idStr: "orc-black-speech",
			expected: LanguageID{
				Family:   "orc",
				Branch:   "black",
				Language: "speech",
			},
			expectError: false,
		},
		{
			name:        "empty string",
			idStr:       "",
			expected:    LanguageID{},
			expectError: true,
		},
		{
			name:        "too many parts",
			idStr:       "a-b-c-d-e-f-g",
			expected:    LanguageID{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseLanguageID(tt.idStr)

			if tt.expectError {
				if err == nil {
					t.Errorf("ParseLanguageID() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseLanguageID() unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("ParseLanguageID() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLanguageType_String(t *testing.T) {
	tests := []struct {
		name     string
		langType LanguageType
		expected string
	}{
		{"unknown", LanguageTypeUnknown, "unknown"},
		{"natural", LanguageTypeNatural, "natural"},
		{"constructed", LanguageTypeConstructed, "constructed"},
		{"divine", LanguageTypeDivine, "divine"},
		{"ancient", LanguageTypeAncient, "ancient"},
		{"modern", LanguageTypeModern, "modern"},
		{"invalid", LanguageType(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.langType.String()
			if result != tt.expected {
				t.Errorf("LanguageType.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLanguageComplexity_String(t *testing.T) {
	tests := []struct {
		name       string
		complexity LanguageComplexity
		expected   string
	}{
		{"unknown", LanguageComplexityUnknown, "unknown"},
		{"simple", LanguageComplexitySimple, "simple"},
		{"moderate", LanguageComplexityModerate, "moderate"},
		{"complex", LanguageComplexityComplex, "complex"},
		{"very_complex", LanguageComplexityVeryComplex, "very_complex"},
		{"invalid", LanguageComplexity(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.complexity.String()
			if result != tt.expected {
				t.Errorf("LanguageComplexity.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewLanguage(t *testing.T) {
	id := LanguageID{Family: "elv", Language: "sindarin"}
	seed := int64(12345)

	lang := NewLanguage(id, "Sindarin", LanguageTypeNatural, seed)

	if lang.ID != id {
		t.Errorf("NewLanguage() ID = %v, want %v", lang.ID, id)
	}

	if lang.Name != "Sindarin" {
		t.Errorf("NewLanguage() Name = %v, want %v", lang.Name, "Sindarin")
	}

	if lang.Type != LanguageTypeNatural {
		t.Errorf("NewLanguage() Type = %v, want %v", lang.Type, LanguageTypeNatural)
	}

	if lang.Complexity != LanguageComplexityModerate {
		t.Errorf("NewLanguage() Complexity = %v, want %v", lang.Complexity, LanguageComplexityModerate)
	}

	if lang.Seed != seed {
		t.Errorf("NewLanguage() Seed = %v, want %v", lang.Seed, seed)
	}

	if lang.RNG == nil {
		t.Error("NewLanguage() RNG is nil")
	}

	if lang.CreatedAt.IsZero() {
		t.Error("NewLanguage() CreatedAt is zero")
	}
}

func TestLanguage_Setters(t *testing.T) {
	lang := NewLanguage(LanguageID{Family: "test"}, "Test", LanguageTypeNatural, 12345)

	// Test SetComplexity
	lang.SetComplexity(LanguageComplexityComplex)
	if lang.Complexity != LanguageComplexityComplex {
		t.Errorf("SetComplexity() failed, got %v, want %v", lang.Complexity, LanguageComplexityComplex)
	}

	// Test SetCulture
	lang.SetCulture("test-culture")
	if lang.Culture != "test-culture" {
		t.Errorf("SetCulture() failed, got %v, want %v", lang.Culture, "test-culture")
	}

	// Test SetDescription
	description := "A test language for testing purposes"
	lang.SetDescription(description)
	if lang.Description != description {
		t.Errorf("SetDescription() failed, got %v, want %v", lang.Description, description)
	}
}

func TestLanguage_IsComplete(t *testing.T) {
	lang := NewLanguage(LanguageID{Family: "test"}, "Test", LanguageTypeNatural, 12345)

	// Initially incomplete
	if lang.IsComplete() {
		t.Error("IsComplete() should return false for new language")
	}

	// Check missing components
	missing := lang.GetMissingComponents()
	expectedMissing := []string{"phonology", "orthography", "morphology", "grammar"}

	if len(missing) != len(expectedMissing) {
		t.Errorf("GetMissingComponents() returned %d components, want %d", len(missing), len(expectedMissing))
	}

	// Check that all expected components are in the missing list
	for _, expected := range expectedMissing {
		found := false
		for _, actual := range missing {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetMissingComponents() missing expected component: %s", expected)
		}
	}
}

func TestLanguage_Clone(t *testing.T) {
	original := NewLanguage(LanguageID{Family: "elv"}, "Original", LanguageTypeNatural, 12345)
	original.SetComplexity(LanguageComplexityComplex)
	original.SetCulture("elvish")
	original.SetDescription("Original language")

	newID := LanguageID{Family: "elv", Branch: "clone"}
	newSeed := int64(54321)

	clone := original.Clone(newID, newSeed)

	// Check that clone has new ID and seed
	if clone.ID != newID {
		t.Errorf("Clone() ID = %v, want %v", clone.ID, newID)
	}

	if clone.Seed != newSeed {
		t.Errorf("Clone() Seed = %v, want %v", clone.Seed, newSeed)
	}

	// Check that clone has parent reference
	if clone.ParentID == nil || *clone.ParentID != original.ID {
		t.Errorf("Clone() ParentID = %v, want %v", clone.ParentID, original.ID)
	}

	// Check that clone has updated name
	expectedName := "Original (clone)"
	if clone.Name != expectedName {
		t.Errorf("Clone() Name = %v, want %v", clone.Name, expectedName)
	}

	// Check that clone has copied properties
	if clone.Complexity != original.Complexity {
		t.Errorf("Clone() Complexity = %v, want %v", clone.Complexity, original.Complexity)
	}

	if clone.Culture != original.Culture {
		t.Errorf("Clone() Culture = %v, want %v", clone.Culture, original.Culture)
	}

	if clone.Description != original.Description {
		t.Errorf("Clone() Description = %v, want %v", clone.Description, original.Description)
	}

	// Check that clone has new creation time
	if clone.CreatedAt.Equal(original.CreatedAt) {
		t.Error("Clone() CreatedAt should be different from original")
	}
}

func TestLanguage_String(t *testing.T) {
	id := LanguageID{Family: "elv", Language: "sindarin"}
	lang := NewLanguage(id, "Sindarin", LanguageTypeNatural, 12345)

	expected := "Language(elv-sindarin: Sindarin)"
	result := lang.String()

	if result != expected {
		t.Errorf("Language.String() = %v, want %v", result, expected)
	}
}
