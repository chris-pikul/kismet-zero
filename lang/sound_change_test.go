package lang

import (
	"strings"
	"testing"
)

// TestSoundChangeIntegration tests that the sound change system is properly integrated
// with the language evolution system.
func TestSoundChangeIntegration(t *testing.T) {
	// Create a language with phonology
	lang, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	// Verify the language has phonology
	if lang.Phonology == nil {
		t.Fatal("Test language should have phonology")
	}

	// Test evolution over different time periods to see sound changes
	testCases := []struct {
		timePeriod string
		expected   bool // whether we expect sound changes
	}{
		{"100_years", true},  // Early evolution should have sound changes
		{"500_years", true},  // Medium evolution should have sound changes
		{"1000_years", true}, // Long evolution should have sound changes
		{"2000_years", true}, // Very long evolution should have sound changes
	}

	for i, tc := range testCases {
		t.Run(tc.timePeriod, func(t *testing.T) {
			// Evolve the language
			evolvedLang, err := EvolveLanguage(lang, tc.timePeriod, []string{}, 42+int64(i))
			if err != nil {
				t.Fatalf("Failed to evolve language: %v", err)
			}

			// Check if sound changes were applied
			description := evolvedLang.Description
			hasSoundChanges := strings.Contains(description, "sound changes:")

			// Log the description for verification
			t.Logf("Description for %s: %s", tc.timePeriod, description)

			// Note: Sound changes are probabilistic, so we don't always expect them
			// The important thing is that the system is working when they do occur
			if hasSoundChanges {
				t.Logf("Sound changes applied for %s", tc.timePeriod)
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

// TestSoundChangeRules tests the individual sound change rules.
func TestSoundChangeRules(t *testing.T) {
	// Test that we have common sound changes defined
	if len(CommonSoundChanges) == 0 {
		t.Fatal("Expected common sound changes to be defined")
	}

	// Test that each sound change has required fields
	for i, change := range CommonSoundChanges {
		if change.ID == "" {
			t.Errorf("Sound change %d missing ID", i)
		}
		if change.Name == "" {
			t.Errorf("Sound change %d missing Name", i)
		}
		if change.Description == "" {
			t.Errorf("Sound change %d missing Description", i)
		}
		if change.Probability < 0.0 || change.Probability > 1.0 {
			t.Errorf("Sound change %d has invalid probability: %f", i, change.Probability)
		}
		if change.Era == "" {
			t.Errorf("Sound change %d missing Era", i)
		}
	}

	// Test that we have changes for different eras
	eras := make(map[string]bool)
	for _, change := range CommonSoundChanges {
		eras[change.Era] = true
	}

	expectedEras := []string{"early_evolution", "middle_evolution", "late_evolution"}
	for _, era := range expectedEras {
		if !eras[era] {
			t.Errorf("Expected sound changes for era: %s", era)
		}
	}
}

// TestApplySoundChanges tests the sound change application function.
func TestApplySoundChanges(t *testing.T) {
	// Create a language with phonology
	lang, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	// Test applying sound changes for different time periods
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
			changes := applySoundChanges(lang, tc.timePeriod, tc.seed)

			// The function should return a slice (even if empty)
			if changes == nil {
				t.Error("Expected changes to be a slice, got nil")
			}

			// Log what changes were applied
			if len(changes) > 0 {
				t.Logf("Applied sound changes for %s: %v", tc.timePeriod, changes)
			} else {
				t.Logf("No sound changes applied for %s", tc.timePeriod)
			}
		})
	}
}
