package evolution

import (
	"testing"
	"time"

	"github.com/chris-pikul/kismet-zero/lang"
)

// TestContactEvolutionEngineStruct tests the ContactEvolutionEngine struct and methods.
func TestContactEvolutionEngineStruct(t *testing.T) {
	config := EvolutionConfig{
		Seed:                    42,
		SoundShiftProbability:   0.3,
		MorphologyChangeRate:    0.2,
		OrthographyChangeRate:   0.1,
		ContactInfluenceRate:    0.4,
		CulturalInfluenceWeight: 0.3,
		DialectFormationRate:    0.2,
	}

	engine := NewContactEvolutionEngine(config)
	if engine == nil {
		t.Fatal("Expected ContactEvolutionEngine to be created")
	}

	// Test that the engine has the expected configuration
	if engine.config.Seed != 42 {
		t.Errorf("Expected seed 42, got %d", engine.config.Seed)
	}
}

// TestSimulatePhonologicalBorrowing tests the phonological borrowing simulation.
func TestSimulatePhonologicalBorrowing(t *testing.T) {
	config := EvolutionConfig{
		Seed:                    42,
		SoundShiftProbability:   0.3,
		MorphologyChangeRate:    0.2,
		OrthographyChangeRate:   0.1,
		ContactInfluenceRate:    0.4,
		CulturalInfluenceWeight: 0.3,
		DialectFormationRate:    0.2,
	}

	engine := NewContactEvolutionEngine(config)

	// Create test languages
	sourceLang, err := lang.CreateRandomLanguage("source", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create source language: %v", err)
	}

	targetLang, err := lang.CreateRandomLanguage("target", "ancient", 43)
	if err != nil {
		t.Fatalf("Failed to create target language: %v", err)
	}

	// Test with different intensity levels
	testCases := []struct {
		intensity float32
		name      string
	}{
		{0.1, "low_intensity"},
		{0.5, "medium_intensity"},
		{0.9, "high_intensity"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			changes := engine.simulatePhonologicalBorrowing(sourceLang, targetLang, tc.intensity)

			// The function should return changes (even if empty)
			if changes == nil {
				t.Logf("No phonological borrowing changes for intensity %f (this is expected behavior)", tc.intensity)
			} else {
				t.Logf("Applied phonological borrowing changes for intensity %f: %v", tc.intensity, changes)

				// Verify the changes have the expected properties
				for _, change := range changes {
					if change.Type != ChangeTypeSoundShift {
						t.Errorf("Expected change type ChangeTypeSoundShift, got %s", change.Type.String())
					}
					if change.Direction != ChangeDirectionAdditive {
						t.Errorf("Expected change direction ChangeDirectionAdditive, got %s", change.Direction.String())
					}
					if change.Era != "contact_evolution" {
						t.Errorf("Expected era 'contact_evolution', got %s", change.Era)
					}
					if change.Trigger != "phonological_contact" {
						t.Errorf("Expected trigger 'phonological_contact', got %s", change.Trigger)
					}
					if change.CultureInfluence != sourceLang.Culture {
						t.Errorf("Expected culture influence %s, got %s", sourceLang.Culture, change.CultureInfluence)
					}
				}
			}
		})
	}
}

// TestSimulateGrammaticalBorrowing tests the grammatical borrowing simulation.
func TestSimulateGrammaticalBorrowing(t *testing.T) {
	config := EvolutionConfig{
		Seed:                    42,
		SoundShiftProbability:   0.3,
		MorphologyChangeRate:    0.2,
		OrthographyChangeRate:   0.1,
		ContactInfluenceRate:    0.4,
		CulturalInfluenceWeight: 0.3,
		DialectFormationRate:    0.2,
	}

	engine := NewContactEvolutionEngine(config)

	// Create test languages
	sourceLang, err := lang.CreateRandomLanguage("source", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create source language: %v", err)
	}

	targetLang, err := lang.CreateRandomLanguage("target", "ancient", 43)
	if err != nil {
		t.Fatalf("Failed to create target language: %v", err)
	}

	// Test with different intensity levels
	testCases := []struct {
		intensity float32
		name      string
	}{
		{0.1, "low_intensity"},
		{0.5, "medium_intensity"},
		{0.9, "high_intensity"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			changes := engine.simulateGrammaticalBorrowing(sourceLang, targetLang, tc.intensity)

			// The function should return changes (even if empty)
			if changes == nil {
				t.Logf("No grammatical borrowing changes for intensity %f (this is expected behavior)", tc.intensity)
			} else {
				t.Logf("Applied grammatical borrowing changes for intensity %f: %v", tc.intensity, changes)

				// Verify the changes have the expected properties
				for _, change := range changes {
					if change.Type != ChangeTypeMorphological {
						t.Errorf("Expected change type ChangeTypeMorphological, got %s", change.Type.String())
					}
					if change.Direction != ChangeDirectionAdditive {
						t.Errorf("Expected change direction ChangeDirectionAdditive, got %s", change.Direction.String())
					}
					if change.Era != "contact_evolution" {
						t.Errorf("Expected era 'contact_evolution', got %s", change.Era)
					}
					if change.Trigger != "grammatical_contact" {
						t.Errorf("Expected trigger 'grammatical_contact', got %s", change.Trigger)
					}
					if change.CultureInfluence != sourceLang.Culture {
						t.Errorf("Expected culture influence %s, got %s", sourceLang.Culture, change.CultureInfluence)
					}
				}
			}
		})
	}
}

// TestSimulateOrthographicBorrowing tests the orthographic borrowing simulation.
func TestSimulateOrthographicBorrowing(t *testing.T) {
	config := EvolutionConfig{
		Seed:                    42,
		SoundShiftProbability:   0.3,
		MorphologyChangeRate:    0.2,
		OrthographyChangeRate:   0.1,
		ContactInfluenceRate:    0.4,
		CulturalInfluenceWeight: 0.3,
		DialectFormationRate:    0.2,
	}

	engine := NewContactEvolutionEngine(config)

	// Create test languages
	sourceLang, err := lang.CreateRandomLanguage("source", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create source language: %v", err)
	}

	targetLang, err := lang.CreateRandomLanguage("target", "ancient", 43)
	if err != nil {
		t.Fatalf("Failed to create target language: %v", err)
	}

	// Test with different contact types and intensity levels
	testCases := []struct {
		contactType ContactType
		intensity   float32
		name        string
	}{
		{ContactTypeTrade, 0.1, "trade_low_intensity"},
		{ContactTypeCultural, 0.5, "cultural_medium_intensity"},
		{ContactTypeConquest, 0.9, "conquest_high_intensity"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			changes := engine.simulateOrthographicBorrowing(sourceLang, targetLang, tc.intensity, tc.contactType)

			// The function should return changes (even if empty)
			if changes == nil {
				t.Logf("No orthographic borrowing changes for contact type %s and intensity %f (this is expected behavior)", tc.contactType.String(), tc.intensity)
			} else {
				t.Logf("Applied orthographic borrowing changes for contact type %s and intensity %f: %v", tc.contactType.String(), tc.intensity, changes)

				// Verify the changes have the expected properties
				for _, change := range changes {
					if change.Type != ChangeTypeOrthographic {
						t.Errorf("Expected change type ChangeTypeOrthographic, got %s", change.Type.String())
					}
					if change.Direction != ChangeDirectionAdditive {
						t.Errorf("Expected change direction ChangeDirectionAdditive, got %s", change.Direction.String())
					}
					if change.Era != "contact_evolution" {
						t.Errorf("Expected era 'contact_evolution', got %s", change.Era)
					}
					if change.Trigger != "orthographic_contact" {
						t.Errorf("Expected trigger 'orthographic_contact', got %s", change.Trigger)
					}
					if change.CultureInfluence != sourceLang.Culture {
						t.Errorf("Expected culture influence %s, got %s", sourceLang.Culture, change.CultureInfluence)
					}
				}
			}
		})
	}
}

// TestCalculateContactInfluence tests the contact influence calculation.
func TestCalculateContactInfluence(t *testing.T) {
	config := EvolutionConfig{
		Seed:                    42,
		SoundShiftProbability:   0.3,
		MorphologyChangeRate:    0.2,
		OrthographyChangeRate:   0.1,
		ContactInfluenceRate:    0.4,
		CulturalInfluenceWeight: 0.3,
		DialectFormationRate:    0.2,
	}

	engine := NewContactEvolutionEngine(config)

	// Create test languages
	sourceLang, err := lang.CreateRandomLanguage("source", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create source language: %v", err)
	}

	targetLang, err := lang.CreateRandomLanguage("target", "ancient", 43)
	if err != nil {
		t.Fatalf("Failed to create target language: %v", err)
	}

	// Test with different contact types and durations
	testCases := []struct {
		contactType ContactType
		duration    time.Duration
		name        string
	}{
		{ContactTypeTrade, time.Hour * 24 * 365, "trade_one_year"},
		{ContactTypeCultural, time.Hour * 24 * 365 * 5, "cultural_five_years"},
		{ContactTypeConquest, time.Hour * 24 * 365 * 10, "conquest_ten_years"},
		{ContactTypeMigration, time.Hour * 24 * 365 * 20, "migration_twenty_years"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			influence := engine.CalculateContactInfluence(sourceLang, targetLang, tc.contactType, tc.duration)

			// Influence should be a valid value between 0 and 1
			if influence < 0.0 || influence > 1.0 {
				t.Errorf("Expected influence between 0.0 and 1.0, got: %f", influence)
			}

			t.Logf("Contact influence for %s over %v: %f", tc.contactType.String(), tc.duration, influence)

			// Verify that conquest has higher influence than trade
			if tc.contactType == ContactTypeConquest && tc.duration == time.Hour*24*365 {
				tradeInfluence := engine.CalculateContactInfluence(sourceLang, targetLang, ContactTypeTrade, tc.duration)
				if influence <= tradeInfluence {
					t.Errorf("Expected conquest influence (%f) to be higher than trade influence (%f)", influence, tradeInfluence)
				}
			}
		})
	}
}

// TestContactTypeString tests the ContactType string representation.
func TestContactTypeString(t *testing.T) {
	testCases := []struct {
		contactType ContactType
		expected    string
	}{
		{ContactTypeUnknown, "unknown"},
		{ContactTypeTrade, "trade"},
		{ContactTypeCultural, "cultural"},
		{ContactTypeConquest, "conquest"},
		{ContactTypeMigration, "migration"},
		{ContactTypeReligious, "religious"},
		{ContactTypeEducational, "educational"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := tc.contactType.String()
			if result != tc.expected {
				t.Errorf("Expected ContactType %d to return '%s', got '%s'", tc.contactType, tc.expected, result)
			}
		})
	}

	// Test unknown contact type
	unknownType := ContactType(99)
	if unknownType.String() != "unknown" {
		t.Errorf("Expected unknown contact type to return 'unknown', got: %s", unknownType.String())
	}
}
