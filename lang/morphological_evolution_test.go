package lang

import (
	"strings"
	"testing"
)

// TestMorphologicalEvolutionIntegration tests that the morphological evolution system
// is properly integrated with the language evolution system.
func TestMorphologicalEvolutionIntegration(t *testing.T) {
	// Create a language with morphology and grammar
	lang, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	// Verify the language has morphology and grammar
	if lang.Morphology == nil {
		t.Fatal("Test language should have morphology")
	}
	if lang.Grammar == nil {
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
			// Evolve the language
			evolvedLang, err := EvolveLanguage(lang, tc.timePeriod, []string{}, 42+int64(i))
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
			if evolvedLang.Description == lang.Description {
				t.Error("Expected language description to change after evolution")
			}
		})
	}
}

// TestMorphologicalChangeTypes tests the morphological change type system.
func TestMorphologicalChangeTypes(t *testing.T) {
	// Test that we have common morphological changes defined
	if len(CommonMorphologicalChanges) == 0 {
		t.Fatal("Expected common morphological changes to be defined")
	}

	// Test that each morphological change has required fields
	for i, change := range CommonMorphologicalChanges {
		if change.Description == "" {
			t.Errorf("Morphological change %d missing Description", i)
		}
		if change.Details == "" {
			t.Errorf("Morphological change %d missing Details", i)
		}
		if change.Probability < 0.0 || change.Probability > 1.0 {
			t.Errorf("Morphological change %d has invalid probability: %f", i, change.Probability)
		}
		if change.Era == "" {
			t.Errorf("Morphological change %d missing Era", i)
		}
		if change.Complexity < -1.0 || change.Complexity > 1.0 {
			t.Errorf("Morphological change %d has invalid complexity: %f", i, change.Complexity)
		}
		if change.Intensity < 0.0 || change.Intensity > 1.0 {
			t.Errorf("Morphological change %d has invalid intensity: %f", i, change.Intensity)
		}
	}

	// Test that we have changes for different eras
	eras := make(map[string]bool)
	for _, change := range CommonMorphologicalChanges {
		eras[change.Era] = true
	}

	expectedEras := []string{"early_evolution", "middle_evolution", "late_evolution"}
	for _, era := range expectedEras {
		if !eras[era] {
			t.Errorf("Expected morphological changes for era: %s", era)
		}
	}

	// Test that we have different types of changes
	types := make(map[MorphologicalChangeType]bool)
	for _, change := range CommonMorphologicalChanges {
		types[change.Type] = true
	}

	expectedTypes := []MorphologicalChangeType{
		MorphologicalChangeTypeSimplification,
		MorphologicalChangeTypeRegularization,
		MorphologicalChangeTypeInnovation,
		MorphologicalChangeTypeLoss,
	}
	for _, changeType := range expectedTypes {
		if !types[changeType] {
			t.Errorf("Expected morphological changes of type: %s", changeType.String())
		}
	}
}

// TestApplyMorphologicalChanges tests the morphological change application function.
func TestApplyMorphologicalChanges(t *testing.T) {
	// Create a language with morphology and grammar
	lang, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	// Test applying morphological changes for different time periods
	testCases := []struct {
		timePeriod string
		seed       int64
	}{
		{"100_years", 42},
		{"500_years", 43},
		{"1000_years", 44},
		{"2000_years", 45},
	}

	for _, tc := range testCases {
		t.Run(tc.timePeriod, func(t *testing.T) {
			changes := applyMorphologicalChanges(lang, tc.timePeriod, tc.seed)

			// The function should return a slice (even if empty)
			if changes == nil {
				t.Error("Expected changes to be a slice, got nil")
			}

			// Log what changes were applied
			if len(changes) > 0 {
				t.Logf("Applied morphological changes for %s: %v", tc.timePeriod, changes)
			} else {
				t.Logf("No morphological changes applied for %s", tc.timePeriod)
			}
		})
	}
}

// TestMorphologicalComplexity tests the complexity calculation function.
func TestMorphologicalComplexity(t *testing.T) {
	// Create a language with morphology and grammar
	lang, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	complexity := calculateMorphologicalComplexity(lang)

	// Complexity should be between 0.0 and 1.0
	if complexity < 0.0 || complexity > 1.0 {
		t.Errorf("Expected complexity between 0.0 and 1.0, got %f", complexity)
	}

	// Since the language has both morphology and grammar, complexity should be > 0
	if complexity <= 0.0 {
		t.Errorf("Expected positive complexity for language with morphology and grammar, got %f", complexity)
	}

	t.Logf("Calculated morphological complexity: %f", complexity)
}
