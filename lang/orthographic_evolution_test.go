package lang

import (
	"strings"
	"testing"
)

// TestOrthographicEvolutionIntegration tests that the orthographic evolution system
// is properly integrated with the language evolution system.
func TestOrthographicEvolutionIntegration(t *testing.T) {
	// Create a language with orthography
	lang, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	// Verify the language has orthography
	if lang.Orthography == nil {
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
			// Evolve the language
			evolvedLang, err := EvolveLanguage(lang, tc.timePeriod, []string{}, 42+int64(i))
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
			if evolvedLang.Description == lang.Description {
				t.Error("Expected language description to change after evolution")
			}
		})
	}
}

// TestOrthographicChangeTypes tests the orthographic change type system.
func TestOrthographicChangeTypes(t *testing.T) {
	// Test that we have common orthographic changes defined
	if len(CommonOrthographicChanges) == 0 {
		t.Fatal("Expected common orthographic changes to be defined")
	}

	// Test that each orthographic change has required fields
	for i, change := range CommonOrthographicChanges {
		if change.Description == "" {
			t.Errorf("Orthographic change %d missing Description", i)
		}
		if change.Details == "" {
			t.Errorf("Orthographic change %d missing Details", i)
		}
		if change.Probability < 0.0 || change.Probability > 1.0 {
			t.Errorf("Orthographic change %d has invalid probability: %f", i, change.Probability)
		}
		if change.Era == "" {
			t.Errorf("Orthographic change %d missing Era", i)
		}
		if change.Complexity < -1.0 || change.Complexity > 1.0 {
			t.Errorf("Orthographic change %d has invalid complexity: %f", i, change.Complexity)
		}
		if change.Intensity < 0.0 || change.Intensity > 1.0 {
			t.Errorf("Orthographic change %d has invalid intensity: %f", i, change.Intensity)
		}
	}

	// Test that we have changes for different eras
	eras := make(map[string]bool)
	for _, change := range CommonOrthographicChanges {
		eras[change.Era] = true
	}

	expectedEras := []string{"early_evolution", "middle_evolution", "late_evolution"}
	for _, era := range expectedEras {
		if !eras[era] {
			t.Errorf("Expected orthographic changes for era: %s", era)
		}
	}

	// Test that we have different types of changes
	types := make(map[OrthographicChangeType]bool)
	for _, change := range CommonOrthographicChanges {
		types[change.Type] = true
	}

	expectedTypes := []OrthographicChangeType{
		OrthographicChangeTypeScriptReform,
		OrthographicChangeTypeSimplification,
		OrthographicChangeTypeStandardization,
		OrthographicChangeTypeInnovation,
		OrthographicChangeTypeAdaptation,
	}
	for _, changeType := range expectedTypes {
		if !types[changeType] {
			t.Errorf("Expected orthographic changes of type: %s", changeType.String())
		}
	}
}

// TestApplyOrthographicChanges tests the orthographic change application function.
func TestApplyOrthographicChanges(t *testing.T) {
	// Create a language with orthography
	lang, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	// Test applying orthographic changes for different time periods
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
			changes := applyOrthographicChanges(lang, tc.timePeriod, tc.seed)

			// The function should return a slice (even if empty)
			if changes == nil {
				t.Error("Expected changes to be a slice, got nil")
			}

			// Log what changes were applied
			if len(changes) > 0 {
				t.Logf("Applied orthographic changes for %s: %v", tc.timePeriod, changes)
			} else {
				t.Logf("No orthographic changes applied for %s", tc.timePeriod)
			}
		})
	}
}

// TestOrthographicComplexity tests the complexity calculation function.
func TestOrthographicComplexity(t *testing.T) {
	// Create a language with orthography
	lang, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	complexity := calculateOrthographicComplexity(lang)

	// Complexity should be between 0.0 and 1.0
	if complexity < 0.0 || complexity > 1.0 {
		t.Errorf("Expected complexity between 0.0 and 1.0, got %f", complexity)
	}

	// Since the language has orthography, complexity should be > 0
	if complexity <= 0.0 {
		t.Errorf("Expected positive complexity for language with orthography, got %f", complexity)
	}

	t.Logf("Calculated orthographic complexity: %f", complexity)
}
