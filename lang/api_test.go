package lang

import (
	"testing"
)

// TestCreateRandomLanguage tests the basic language creation functionality.
func TestCreateRandomLanguage(t *testing.T) {
	tests := []struct {
		name    string
		family  string
		culture string
		seed    int64
		wantErr bool
	}{
		{
			name:    "elvish language",
			family:  "elvish",
			culture: "high_elf",
			seed:    42,
			wantErr: false,
		},
		{
			name:    "dwarvish language",
			family:  "dwarvish",
			culture: "mountain_dwarf",
			seed:    43,
			wantErr: false,
		},
		{
			name:    "orcish language",
			family:  "orcish",
			culture: "black_orc",
			seed:    44,
			wantErr: false,
		},
		{
			name:    "human language",
			family:  "human",
			culture: "northern_human",
			seed:    45,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lang, err := CreateRandomLanguage(tt.family, tt.culture, tt.seed)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateRandomLanguage() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("CreateRandomLanguage() unexpected error: %v", err)
				return
			}

			if lang == nil {
				t.Errorf("CreateRandomLanguage() returned nil language")
				return
			}

			// Check basic properties
			if lang.ID.Family != tt.family {
				t.Errorf("CreateRandomLanguage() family = %v, want %v", lang.ID.Family, tt.family)
			}

			if lang.Culture != tt.culture {
				t.Errorf("CreateRandomLanguage() culture = %v, want %v", lang.Culture, tt.culture)
			}

			if lang.Seed != tt.seed {
				t.Errorf("CreateRandomLanguage() seed = %v, want %v", lang.Seed, tt.seed)
			}

			// Check that language is complete
			if !lang.IsComplete() {
				t.Errorf("CreateRandomLanguage() language is not complete, missing: %v", lang.GetMissingComponents())
			}
		})
	}
}

// TestCulturalLanguageConfig tests the cultural configuration generation.
func TestCulturalLanguageConfig(t *testing.T) {
	tests := []struct {
		culture        string
		wantSize       int
		wantComplexity LanguageComplexity
	}{
		{
			culture:        "elvish",
			wantSize:       35,
			wantComplexity: LanguageComplexityComplex,
		},
		{
			culture:        "dwarvish",
			wantSize:       28,
			wantComplexity: LanguageComplexityModerate,
		},
		{
			culture:        "orcish",
			wantSize:       20,
			wantComplexity: LanguageComplexitySimple,
		},
		{
			culture:        "human",
			wantSize:       25,
			wantComplexity: LanguageComplexityModerate,
		},
		{
			culture:        "ancient",
			wantSize:       40,
			wantComplexity: LanguageComplexityVeryComplex,
		},
	}

	for _, tt := range tests {
		t.Run(tt.culture, func(t *testing.T) {
			config := CulturalLanguageConfig(tt.culture)

			if config.PhonemeInventorySize != tt.wantSize {
				t.Errorf("CulturalLanguageConfig() phoneme size = %v, want %v", config.PhonemeInventorySize, tt.wantSize)
			}

			if config.ComplexityTarget != tt.wantComplexity {
				t.Errorf("CulturalLanguageConfig() complexity = %v, want %v", config.ComplexityTarget, tt.wantComplexity)
			}
		})
	}
}

// TestGenerateNames tests the name generation functionality.
func TestGenerateNames(t *testing.T) {
	// Create a test language
	lang, err := CreateRandomLanguage("test", "elvish", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	tests := []struct {
		nameType        string
		count           int
		culturalContext string
		seed            int64
		wantErr         bool
	}{
		{
			nameType:        "person",
			count:           5,
			culturalContext: "noble",
			seed:            42,
			wantErr:         false,
		},
		{
			nameType:        "place",
			count:           3,
			culturalContext: "forest",
			seed:            43,
			wantErr:         false,
		},
		{
			nameType:        "deity",
			count:           2,
			culturalContext: "nature",
			seed:            44,
			wantErr:         false,
		},
		{
			nameType:        "artifact",
			count:           4,
			culturalContext: "magical",
			seed:            45,
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nameType, func(t *testing.T) {
			names, err := GenerateNames(lang, tt.nameType, tt.count, tt.culturalContext, tt.seed)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GenerateNames() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("GenerateNames() unexpected error: %v", err)
				return
			}

			if len(names) != tt.count {
				t.Errorf("GenerateNames() returned %d names, want %d", len(names), tt.count)
			}

			// Check that names are not empty
			for i, name := range names {
				if name == "" {
					t.Errorf("GenerateNames() name[%d] is empty", i)
				}
			}
		})
	}
}

// TestGenerateText tests the text generation functionality.
func TestGenerateText(t *testing.T) {
	// Create a test language
	lang, err := CreateRandomLanguage("test", "elvish", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	tests := []struct {
		textType        string
		length          string
		culturalContext string
		seed            int64
		wantErr         bool
	}{
		{
			textType:        "lore",
			length:          "medium",
			culturalContext: "ancient",
			seed:            42,
			wantErr:         false,
		},
		{
			textType:        "poetry",
			length:          "short",
			culturalContext: "nature",
			seed:            43,
			wantErr:         false,
		},
		{
			textType:        "dialogue",
			length:          "short",
			culturalContext: "formal",
			seed:            44,
			wantErr:         false,
		},
		{
			textType:        "inscription",
			length:          "short",
			culturalContext: "temple",
			seed:            45,
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.textType, func(t *testing.T) {
			text, err := GenerateText(lang, tt.textType, tt.length, tt.culturalContext, tt.seed)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GenerateText() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("GenerateText() unexpected error: %v", err)
				return
			}

			if text == "" {
				t.Errorf("GenerateText() returned empty text")
			}
		})
	}
}

// TestValidateLanguage tests the language validation functionality.
func TestValidateLanguage(t *testing.T) {
	// Test with a complete language
	completeLang, err := CreateRandomLanguage("test", "elvish", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	result := ValidateLanguage(completeLang)
	if !result.IsValid {
		t.Errorf("ValidateLanguage() should be valid for complete language, got: %v", result)
	}

	// Test with an incomplete language
	incompleteLang := NewLanguage(LanguageID{Family: "test"}, "test", LanguageTypeConstructed, 42)
	result = ValidateLanguage(incompleteLang)
	if result.IsValid {
		t.Errorf("ValidateLanguage() should not be valid for incomplete language")
	}

	if len(result.Issues) == 0 {
		t.Errorf("ValidateLanguage() should report issues for incomplete language")
	}
}

// TestGetLanguageInfo tests the language info extraction functionality.
func TestGetLanguageInfo(t *testing.T) {
	lang, err := CreateRandomLanguage("test", "elvish", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	info := GetLanguageInfo(lang)

	if info.BasicStats.PhonemeCount == 0 {
		t.Errorf("GetLanguageInfo() phoneme count should not be 0")
	}

	if info.CulturalFeatures.Culture != lang.Culture {
		t.Errorf("GetLanguageInfo() culture = %v, want %v", info.CulturalFeatures.Culture, lang.Culture)
	}

	if info.EvolutionHistory.CreatedAt != lang.CreatedAt {
		t.Errorf("GetLanguageInfo() created at = %v, want %v", info.EvolutionHistory.CreatedAt, lang.CreatedAt)
	}
}
