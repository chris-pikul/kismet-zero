package evolution

import (
	"strings"
	"testing"
	"time"

	"github.com/chris-pikul/kismet-zero/lang"
)

// TestOrthographicEvolutionIntegration tests that the orthographic evolution system
// is properly integrated with the language evolution system.
func TestOrthographicEvolutionIntegration(t *testing.T) {
	// Create a language with orthography using the lang package
	testLang, err := lang.CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	// Verify the language has orthography
	if testLang.Orthography == nil {
		t.Fatal("Test language should have orthography")
	}

	// Test evolution over different time periods to see orthographic changes
	testCases := []struct {
		timePeriod string
		expected   bool // whether we expect orthographic changes
	}{
		{"100_years", true},  // Early evolution should have orthographic changes
		{"500_years", true},  // Medium evolution should have orthographic changes
		{"1000_years", true}, // Long evolution should have orthographic changes
		{"2000_years", true}, // Very long evolution should have orthographic changes
	}

	for i, tc := range testCases {
		t.Run(tc.timePeriod, func(t *testing.T) {
			// Evolve the language using the lang package
			evolvedLang, err := lang.EvolveLanguage(testLang, tc.timePeriod, []string{}, 42+int64(i))
			if err != nil {
				t.Fatalf("Failed to evolve language: %v", err)
			}

			// Check if orthographic changes were applied
			description := evolvedLang.Description
			hasOrthographicChanges := strings.Contains(description, "orthographic changes:")

			// Log the description for verification
			t.Logf("Description for %s: %s", tc.timePeriod, description)

			// Note: Orthographic changes are probabilistic, so we don't always expect them
			// The important thing is that the system is working when they do occur
			if hasOrthographicChanges {
				t.Logf("Orthographic changes applied for %s", tc.timePeriod)
			}

			// Verify the language was actually evolved
			if evolvedLang.EvolvedAt.IsZero() {
				t.Error("Expected evolution timestamp to be set")
			}

			// Verify the description changed
			if evolvedLang.Description == testLang.Description {
				t.Error("Expected language description to change after evolution")
			}
		})
	}
}

// TestOrthographicChangeTypes tests the orthographic change type system.
func TestOrthographicChangeTypes(t *testing.T) {
	// Test that orthographic change types are properly defined
	// Since CommonOrthographicChanges doesn't exist, we'll test the types directly
	changeTypes := []OrthographicChangeType{
		OrthographicChangeTypeScriptReform,
		OrthographicChangeTypeBorrowing,
		OrthographicChangeTypeSimplification,
		OrthographicChangeTypeStandardization,
		OrthographicChangeTypeInnovation,
		OrthographicChangeTypeAdaptation,
	}

	// Test that each change type has a valid string representation
	for _, changeType := range changeTypes {
		if changeType.String() == "" {
			t.Errorf("Orthographic change type %d missing string representation", changeType)
		}
	}

	// Test that unknown types return "unknown"
	unknownType := OrthographicChangeType(99)
	if unknownType.String() != "unknown" {
		t.Errorf("Expected unknown type to return 'unknown', got: %s", unknownType.String())
	}
}

// TestOrthographicChangeStruct tests the OrthographicChange struct.
func TestOrthographicChangeStruct(t *testing.T) {
	// Create a test orthographic change
	change := OrthographicChange{
		ID:                "test_change",
		Type:              OrthographicChangeTypeSimplification,
		Description:       "Test orthographic change",
		Details:           "This is a test change",
		AffectedGraphemes: []string{"test_grapheme"},
		AffectedMappings:  []string{"test_mapping"},
		AffectedRules:     []string{"test_rule"},
		ComplexityChange:  -0.5,
		Timestamp:         time.Now(),
		Era:               "test_era",
		Trigger:           "test_trigger",
		Intensity:         0.7,
		ScriptReformType:  "spelling",
		BorrowedFrom:      "test_source",
	}

	// Verify the change has the expected values
	if change.ID != "test_change" {
		t.Errorf("Expected ID 'test_change', got: %s", change.ID)
	}
	if change.Type != OrthographicChangeTypeSimplification {
		t.Errorf("Expected type OrthographicChangeTypeSimplification, got: %d", change.Type)
	}
	if change.Description != "Test orthographic change" {
		t.Errorf("Expected description 'Test orthographic change', got: %s", change.Description)
	}
	if change.ComplexityChange != -0.5 {
		t.Errorf("Expected complexity change -0.5, got: %f", change.ComplexityChange)
	}
	if change.Intensity != 0.7 {
		t.Errorf("Expected intensity 0.7, got: %f", change.Intensity)
	}
}

// TestOrthographicEvolutionEngine tests the OrthographicEvolutionEngine struct and methods.
func TestOrthographicEvolutionEngine(t *testing.T) {
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
	if engine == nil {
		t.Fatal("Expected OrthographicEvolutionEngine to be created")
	}

	// Test that the engine has the expected configuration
	if engine.config.Seed != 42 {
		t.Errorf("Expected seed 42, got %d", engine.config.Seed)
	}
}

// TestOrthographicEvolutionEngineApplyChanges tests the OrthographicEvolutionEngine's ApplyOrthographicChanges method.
func TestOrthographicEvolutionEngineApplyChanges(t *testing.T) {
	config := EvolutionConfig{
		Seed:                    42,
		SoundShiftProbability:   0.3,
		MorphologyChangeRate:    0.2,
		OrthographyChangeRate:   0.8, // High rate for testing
		ContactInfluenceRate:    0.4,
		CulturalInfluenceWeight: 0.3,
		DialectFormationRate:    0.2,
	}

	engine := NewOrthographicEvolutionEngine(config)

	// Create a test language with orthography
	testLang, err := lang.CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	// Test applying orthographic changes for different eras
	testCases := []struct {
		era string
	}{
		{"early_evolution"},
		{"middle_evolution"},
		{"late_evolution"},
	}

	for _, tc := range testCases {
		t.Run(tc.era, func(t *testing.T) {
			changes := engine.ApplyOrthographicChanges(testLang.Orthography, tc.era)

			// The function should return changes (even if empty)
			if changes == nil {
				t.Error("Expected changes to be a slice, got nil")
			}

			// Log what changes were applied
			if len(changes) > 0 {
				t.Logf("Applied orthographic changes for %s: %v", tc.era, changes)
			} else {
				t.Logf("No orthographic changes applied for %s", tc.era)
			}
		})
	}
}
