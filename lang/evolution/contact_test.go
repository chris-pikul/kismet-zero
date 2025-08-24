package evolution

import (
	"testing"
	"time"

	"github.com/chris-pikul/kismet-zero/lang"
)

// TestContactEvolutionEngine tests the contact evolution engine creation and basic functionality.
func TestContactEvolutionEngine(t *testing.T) {
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
		t.Fatal("Expected engine to be created, got nil")
	}

	// Test that the engine has the expected configuration
	if engine.GetConfig().Seed != 42 {
		t.Errorf("Expected seed 42, got %d", engine.GetConfig().Seed)
	}

	// Test that the engine has RNG initialized
	if engine.GetRNG() == nil {
		t.Error("Expected RNG to be initialized")
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
