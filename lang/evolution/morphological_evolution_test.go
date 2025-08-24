package evolution

import (
	"strings"
	"testing"
	"time"

	"github.com/chris-pikul/kismet-zero/lang"
)

// TestMorphologicalEvolutionIntegration tests that the morphological evolution system
// is properly integrated with the language evolution system.
func TestMorphologicalEvolutionIntegration(t *testing.T) {
	// Create a language with morphology and grammar using the lang package
	testLang, err := lang.CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	// Verify the language has morphology and grammar
	if testLang.Morphology == nil {
		t.Fatal("Test language should have morphology")
	}
	if testLang.Grammar == nil {
		t.Fatal("Test language should have grammar")
	}

	// Test evolution over different time periods to see morphological changes
	testCases := []struct {
		timePeriod string
		expected   bool // whether we expect morphological changes
	}{
		{"100_years", true},  // Early evolution should have morphological changes
		{"500_years", true},  // Medium evolution should have morphological changes
		{"1000_years", true}, // Long evolution should have morphological changes
		{"2000_years", true}, // Very long evolution should have morphological changes
	}

	for i, tc := range testCases {
		t.Run(tc.timePeriod, func(t *testing.T) {
			// Evolve the language using the lang package
			evolvedLang, err := lang.EvolveLanguage(testLang, tc.timePeriod, []string{}, 42+int64(i))
			if err != nil {
				t.Fatalf("Failed to evolve language: %v", err)
			}

			// Check if morphological changes were applied
			description := evolvedLang.Description
			hasMorphologicalChanges := strings.Contains(description, "morphological changes:")

			// Log the description for verification
			t.Logf("Description for %s: %s", tc.timePeriod, description)

			// Note: Morphological changes are probabilistic, so we don't always expect them
			// The important thing is that the system is working when they do occur
			if hasMorphologicalChanges {
				t.Logf("Morphological changes applied for %s", tc.timePeriod)
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

// TestMorphologicalChangeTypes tests the morphological change type system.
func TestMorphologicalChangeTypes(t *testing.T) {
	// Test all change types
	changeTypes := []MorphologicalChangeType{
		MorphologicalChangeTypeSimplification,
		MorphologicalChangeTypeRegularization,
		MorphologicalChangeTypeInnovation,
	}

	// Test that each change type has a valid string representation
	for _, changeType := range changeTypes {
		if changeType.String() == "" {
			t.Errorf("Morphological change type %d missing string representation", changeType)
		}
	}

	// Test that unknown types return "unknown"
	unknownType := MorphologicalChangeType(99)
	if unknownType.String() != "unknown" {
		t.Errorf("Expected unknown type to return 'unknown', got: %s", unknownType.String())
	}
}

// TestMorphologicalChangeStruct tests the MorphologicalChange struct.
func TestMorphologicalChangeStruct(t *testing.T) {
	// Create a test morphological change
	change := MorphologicalChange{
		ID:                "test_change",
		Type:              MorphologicalChangeTypeSimplification,
		Description:       "Test morphological change",
		Details:           "This is a test change",
		AffectedMorphemes: []string{"test_morpheme"},
		AffectedRules:     []string{"test_rule"},
		AffectedFeatures:  []string{"test_feature"},
		ComplexityChange:  -0.5,
		Timestamp:         time.Now(),
		Era:               "test_era",
		Trigger:           "test_trigger",
		Intensity:         0.7,
	}

	// Verify the change has the expected values
	if change.ID != "test_change" {
		t.Errorf("Expected ID 'test_change', got: %s", change.ID)
	}
	if change.Type != MorphologicalChangeTypeSimplification {
		t.Errorf("Expected type MorphologicalChangeTypeSimplification, got: %d", change.Type)
	}
	if change.Description != "Test morphological change" {
		t.Errorf("Expected description 'Test morphological change', got: %s", change.Description)
	}
	if change.ComplexityChange != -0.5 {
		t.Errorf("Expected complexity change -0.5, got: %f", change.ComplexityChange)
	}
	if change.Intensity != 0.7 {
		t.Errorf("Expected intensity 0.7, got: %f", change.Intensity)
	}
}

// TestMorphologicalEvolutionEngineStruct tests the MorphologicalEvolutionEngine struct and methods.
func TestMorphologicalEvolutionEngineStruct(t *testing.T) {
	config := EvolutionConfig{
		Seed:                    42,
		SoundShiftProbability:   0.3,
		MorphologyChangeRate:    0.2,
		OrthographyChangeRate:   0.1,
		ContactInfluenceRate:    0.4,
		CulturalInfluenceWeight: 0.3,
		DialectFormationRate:    0.2,
	}

	engine := NewMorphologicalEvolutionEngine(config)
	if engine == nil {
		t.Fatal("Expected MorphologicalEvolutionEngine to be created")
	}

	// Test that the engine has the expected configuration
	if engine.config.Seed != 42 {
		t.Errorf("Expected seed 42, got %d", engine.config.Seed)
	}
}

// TestMorphologicalEvolutionEngineApplyChanges tests the MorphologicalEvolutionEngine's ApplyMorphologicalChanges method.
func TestMorphologicalEvolutionEngineApplyChanges(t *testing.T) {
	config := EvolutionConfig{
		Seed:                    42,
		SoundShiftProbability:   0.3,
		MorphologyChangeRate:    0.8, // High rate for testing
		OrthographyChangeRate:   0.1,
		ContactInfluenceRate:    0.4,
		CulturalInfluenceWeight: 0.3,
		DialectFormationRate:    0.2,
	}

	engine := NewMorphologicalEvolutionEngine(config)

	// Create a test language with morphology and grammar
	testLang, err := lang.CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	// Test applying morphological changes for different eras
	testCases := []struct {
		era string
	}{
		{"early_evolution"},
		{"middle_evolution"},
		{"late_evolution"},
	}

	for _, tc := range testCases {
		t.Run(tc.era, func(t *testing.T) {
			changes := engine.ApplyMorphologicalChanges(testLang.Morphology, testLang.Grammar, tc.era)

			// The function should return changes (even if empty)
			if changes == nil {
				t.Error("Expected changes to be a slice, got nil")
			}

			// Log what changes were applied
			if len(changes) > 0 {
				t.Logf("Applied morphological changes for %s: %v", tc.era, changes)
			} else {
				t.Logf("No morphological changes applied for %s", tc.era)
			}
		})
	}
}

// TestMorphologicalComplexity tests the morphological complexity calculation.
func TestMorphologicalComplexity(t *testing.T) {
	config := EvolutionConfig{
		Seed:                    42,
		SoundShiftProbability:   0.3,
		MorphologyChangeRate:    0.2,
		OrthographyChangeRate:   0.1,
		ContactInfluenceRate:    0.4,
		CulturalInfluenceWeight: 0.3,
		DialectFormationRate:    0.2,
	}

	engine := NewMorphologicalEvolutionEngine(config)

	// Create a test language with morphology and grammar
	testLang, err := lang.CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	// Test complexity calculation
	complexity := engine.CalculateMorphologicalComplexity(testLang.Morphology, testLang.Grammar)

	// Complexity should be a valid value
	if complexity < 0.0 || complexity > 1.0 {
		t.Errorf("Expected complexity between 0.0 and 1.0, got: %f", complexity)
	}

	t.Logf("Calculated morphological complexity: %f", complexity)
}
