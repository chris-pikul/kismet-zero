package evolution

import (
	"testing"
	"time"

	"github.com/chris-pikul/kismet-zero/lang/orthography"
)

func TestNewOrthographicEvolutionEngine(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewOrthographicEvolutionEngine(config)

	if engine == nil {
		t.Fatal("Orthographic evolution engine should not be nil")
	}

	if engine.config.Seed != 42 {
		t.Errorf("Expected seed 42, got %d", engine.config.Seed)
	}
}

func TestOrthographicEvolutionEngine_ApplyOrthographicChanges(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewOrthographicEvolutionEngine(config)

	// Create a sample writing system
	writingSystem := createTestWritingSystem()

	// Apply orthographic changes
	changes := engine.ApplyOrthographicChanges(writingSystem, "test_era")

	// Changes might be empty if random evolution doesn't trigger
	// This is acceptable behavior for the evolution system
	if changes == nil {
		t.Error("Changes should not be nil")
	}

	// Test that changes have proper structure when they occur
	for _, change := range changes {
		if change.ID == "" {
			t.Error("Change should have an ID")
		}
		if change.Description == "" {
			t.Error("Change should have a description")
		}
		if change.Timestamp.IsZero() {
			t.Error("Change should have a timestamp")
		}
		if change.Era != "test_era" {
			t.Errorf("Expected era 'test_era', got '%s'", change.Era)
		}
	}
}

func TestOrthographicEvolutionEngine_SimulateOrthographicBorrowing(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewOrthographicEvolutionEngine(config)

	// Create two sample writing systems
	targetWS := createTestWritingSystem()
	sourceWS := createTestWritingSystem()
	sourceWS.Name = "Source Writing System"

	// Test borrowing graphemes
	changes := engine.SimulateOrthographicBorrowing(targetWS, sourceWS, "graphemes", 0.8)

	// Changes might be empty due to probability
	if changes != nil {
		for _, change := range changes {
			if change.Type != OrthographicChangeTypeBorrowing {
				t.Errorf("Expected type 'borrowing', got '%s'", change.Type.String())
			}
			if change.BorrowedFrom != sourceWS.Name {
				t.Errorf("Expected borrowed from '%s', got '%s'", sourceWS.Name, change.BorrowedFrom)
			}
		}
	}

	// Test borrowing mappings
	changes = engine.SimulateOrthographicBorrowing(targetWS, sourceWS, "mappings", 0.9)
	if changes != nil {
		for _, change := range changes {
			if change.Type != OrthographicChangeTypeBorrowing {
				t.Errorf("Expected type 'borrowing', got '%s'", change.Type.String())
			}
		}
	}

	// Test borrowing style
	changes = engine.SimulateOrthographicBorrowing(targetWS, sourceWS, "style", 1.0)
	if changes != nil {
		for _, change := range changes {
			if change.Type != OrthographicChangeTypeBorrowing {
				t.Errorf("Expected type 'borrowing', got '%s'", change.Type.String())
			}
		}
	}
}

func TestOrthographicEvolutionEngine_AdaptToPhonologicalChanges(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewOrthographicEvolutionEngine(config)

	// Create a sample writing system
	writingSystem := createTestWritingSystem()

	// Create sample phonological changes
	phonologicalChanges := []LinguisticChange{
		{
			ID:          "test_sound_change",
			Type:        ChangeTypeSoundShift,
			Description: "Test sound change",
			Timestamp:   time.Now(),
		},
		{
			ID:          "test_morphological_change",
			Type:        ChangeTypeMorphological,
			Description: "Test morphological change",
			Timestamp:   time.Now(),
		},
	}

	// Adapt to phonological changes
	changes := engine.AdaptToPhonologicalChanges(writingSystem, phonologicalChanges, "test_era")

	// Should have one adaptation change for the sound shift
	if len(changes) != 1 {
		t.Errorf("Expected 1 adaptation change, got %d", len(changes))
	}

	if changes[0].Type != OrthographicChangeTypeAdaptation {
		t.Errorf("Expected type 'adaptation', got '%s'", changes[0].Type.String())
	}
}

func TestOrthographicEvolutionEngine_CalculateOrthographicComplexity(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewOrthographicEvolutionEngine(config)

	// Test with nil writing system
	complexity := engine.CalculateOrthographicComplexity(nil)
	if complexity != 0.0 {
		t.Errorf("Expected complexity 0.0 for nil writing system, got %f", complexity)
	}

	// Test with empty writing system
	emptyWS := &orthography.WritingSystem{
		Style:     orthography.WritingStyleAlphabetic,
		Graphemes: []orthography.Grapheme{},
		Mappings:  []orthography.OrthographyMapping{},
	}
	complexity = engine.CalculateOrthographicComplexity(emptyWS)
	if complexity != 0.1 {
		t.Errorf("Expected complexity 0.1 for alphabetic system, got %f", complexity)
	}

	// Test with complex writing system
	complexWS := createTestWritingSystem()
	complexity = engine.CalculateOrthographicComplexity(complexWS)
	if complexity <= 0.0 {
		t.Errorf("Expected positive complexity for complex writing system, got %f", complexity)
	}
}

func TestOrthographicEvolutionEngine_GenerateOrthographicReform(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewOrthographicEvolutionEngine(config)

	// Create a sample writing system
	writingSystem := createTestWritingSystem()

	// Test spelling reform
	reform := engine.GenerateOrthographicReform(writingSystem, "spelling", "test_era")
	if reform.Type != OrthographicChangeTypeScriptReform {
		t.Errorf("Expected type 'script_reform', got '%s'", reform.Type.String())
	}
	if reform.ScriptReformType != "spelling" {
		t.Errorf("Expected reform type 'spelling', got '%s'", reform.ScriptReformType)
	}
	if reform.ComplexityChange >= 0 {
		t.Errorf("Expected negative complexity change for reform, got %f", reform.ComplexityChange)
	}

	// Test character reform
	reform = engine.GenerateOrthographicReform(writingSystem, "character", "test_era")
	if reform.ScriptReformType != "character" {
		t.Errorf("Expected reform type 'character', got '%s'", reform.ScriptReformType)
	}

	// Test system reform
	reform = engine.GenerateOrthographicReform(writingSystem, "system", "test_era")
	if reform.ScriptReformType != "system" {
		t.Errorf("Expected reform type 'system', got '%s'", reform.ScriptReformType)
	}

	// Test default reform
	reform = engine.GenerateOrthographicReform(writingSystem, "unknown", "test_era")
	if reform.ScriptReformType != "unknown" {
		t.Errorf("Expected reform type 'unknown', got '%s'", reform.ScriptReformType)
	}
}

func TestOrthographicChangeType_String(t *testing.T) {
	// Test all change types
	testCases := []struct {
		changeType OrthographicChangeType
		expected   string
	}{
		{OrthographicChangeTypeUnknown, "unknown"},
		{OrthographicChangeTypeScriptReform, "script_reform"},
		{OrthographicChangeTypeBorrowing, "borrowing"},
		{OrthographicChangeTypeSimplification, "simplification"},
		{OrthographicChangeTypeStandardization, "standardization"},
		{OrthographicChangeTypeInnovation, "innovation"},
		{OrthographicChangeTypeAdaptation, "adaptation"},
	}

	for _, tc := range testCases {
		result := tc.changeType.String()
		if result != tc.expected {
			t.Errorf("Expected '%s' for %v, got '%s'", tc.expected, tc.changeType, result)
		}
	}

	// Test out of range value
	outOfRange := OrthographicChangeType(99)
	if outOfRange.String() != "unknown" {
		t.Errorf("Expected 'unknown' for out of range value, got '%s'", outOfRange.String())
	}
}

func TestOrthographicChange_Structure(t *testing.T) {
	// Test that OrthographicChange struct has all required fields
	change := OrthographicChange{
		ID:                "test_id",
		Type:              OrthographicChangeTypeSimplification,
		Description:       "Test description",
		Details:           "Test details",
		AffectedGraphemes: []string{"test_grapheme"},
		AffectedMappings:  []string{"test_mapping"},
		AffectedRules:     []string{"test_rule"},
		ComplexityChange:  -0.2,
		Timestamp:         time.Now(),
		Era:               "test_era",
		Trigger:           "test_trigger",
		Intensity:         0.5,
		ScriptReformType:  "test_reform",
		BorrowedFrom:      "test_source",
	}

	if change.ID == "" {
		t.Error("OrthographicChange should have ID field")
	}
	if change.Type == OrthographicChangeTypeUnknown {
		t.Error("OrthographicChange should have valid Type field")
	}
	if change.Description == "" {
		t.Error("OrthographicChange should have Description field")
	}
	if change.Timestamp.IsZero() {
		t.Error("OrthographicChange should have Timestamp field")
	}
}

// TestOrthographicEvolutionEngine_ApplyScriptReforms tests the script reform application.
func TestOrthographicEvolutionEngine_ApplyScriptReforms(t *testing.T) {
	config := EvolutionConfig{
		Seed:                    42,
		SoundShiftProbability:   0.3,
		MorphologyChangeRate:    0.2,
		OrthographyChangeRate:   0.1,
		ContactInfluenceRate:    0.4,
		CulturalInfluenceWeight: 0.3,
		DialectFormationRate:    0.2,
	}

	engine := NewOrthographicEvolutionEngine(config)
	writingSystem := createTestWritingSystem()

	// Test script reforms for different eras
	testCases := []struct {
		era string
	}{
		{"early_evolution"},
		{"middle_evolution"},
		{"late_evolution"},
	}

	for _, tc := range testCases {
		t.Run(tc.era, func(t *testing.T) {
			changes := engine.applyScriptReforms(writingSystem, tc.era)

			// The function should return changes (even if empty)
			if changes == nil {
				t.Logf("No script reform changes for %s (this is expected behavior)", tc.era)
			} else {
				t.Logf("Applied script reform changes for %s: %v", tc.era, changes)

				// Verify the changes have the expected properties
				for _, change := range changes {
					if change.Type != OrthographicChangeTypeScriptReform {
						t.Errorf("Expected change type OrthographicChangeTypeScriptReform, got %s", change.Type.String())
					}
					if change.Era != tc.era {
						t.Errorf("Expected era %s, got %s", tc.era, change.Era)
					}
					if change.Trigger != "orthographic_reform" {
						t.Errorf("Expected trigger 'orthographic_reform', got %s", change.Trigger)
					}
					if change.ScriptReformType == "" {
						t.Error("Expected ScriptReformType to be set")
					}
					if change.ComplexityChange >= 0 {
						t.Errorf("Expected negative complexity change for script reform, got %f", change.ComplexityChange)
					}
				}
			}
		})
	}
}

// TestOrthographicEvolutionEngine_ApplyStandardizationChanges tests the standardization changes.
func TestOrthographicEvolutionEngine_ApplyStandardizationChanges(t *testing.T) {
	config := EvolutionConfig{
		Seed:                    42,
		SoundShiftProbability:   0.3,
		MorphologyChangeRate:    0.2,
		OrthographyChangeRate:   0.1,
		ContactInfluenceRate:    0.4,
		CulturalInfluenceWeight: 0.3,
		DialectFormationRate:    0.2,
	}

	engine := NewOrthographicEvolutionEngine(config)
	writingSystem := createTestWritingSystem()

	// Test standardization changes for different eras
	testCases := []struct {
		era string
	}{
		{"early_evolution"},
		{"middle_evolution"},
		{"late_evolution"},
	}

	for _, tc := range testCases {
		t.Run(tc.era, func(t *testing.T) {
			changes := engine.applyStandardizationChanges(writingSystem, tc.era)

			// The function should return changes (even if empty)
			if changes == nil {
				t.Logf("No standardization changes for %s (this is expected behavior)", tc.era)
			} else {
				t.Logf("Applied standardization changes for %s: %v", tc.era, changes)

				// Verify the changes have the expected properties
				for _, change := range changes {
					if change.Type != OrthographicChangeTypeStandardization {
						t.Errorf("Expected change type OrthographicChangeTypeStandardization, got %s", change.Type.String())
					}
					if change.Era != tc.era {
						t.Errorf("Expected era %s, got %s", tc.era, change.Era)
					}
					if change.Trigger != "standardization" {
						t.Errorf("Expected trigger 'standardization', got %s", change.Trigger)
					}
				}
			}
		})
	}
}

// Helper function to create a test writing system
func createTestWritingSystem() *orthography.WritingSystem {
	// Create some test graphemes
	graphemes := []orthography.Grapheme{
		{Symbol: "a", Type: orthography.GraphemeTypeLetter, Weight: 1.0},
		{Symbol: "b", Type: orthography.GraphemeTypeLetter, Weight: 1.0},
		{Symbol: "c", Type: orthography.GraphemeTypeLetter, Weight: 1.0},
	}

	// Create some test mappings without depending on phoneme pool
	mappings := []orthography.OrthographyMapping{
		{
			Phoneme:   nil, // We don't need actual phonemes for this test
			Graphemes: []orthography.Grapheme{graphemes[0]},
			Primary:   true,
		},
		{
			Phoneme:   nil,                                                // We don't need actual phonemes for this test
			Graphemes: []orthography.Grapheme{graphemes[1], graphemes[2]}, // Irregular mapping
			Primary:   true,
		},
	}

	return &orthography.WritingSystem{
		Style:     orthography.WritingStyleAlphabetic,
		Graphemes: graphemes,
		Mappings:  mappings,
		Name:      "Test Writing System",
		Culture:   "test_culture",
	}
}
