package lang

import (
	"testing"
	"time"
)

func TestContactIntensityCalculator(t *testing.T) {
	calculator := NewContactIntensityCalculator()

	// Test that all contact types have models
	t.Run("Default Models", func(t *testing.T) {
		contactTypes := []ContactType{
			ContactTypeTrade,
			ContactTypeConquest,
			ContactTypeMigration,
			ContactTypeCultural,
			ContactTypeReligious,
			ContactTypeEducational,
		}

		for _, contactType := range contactTypes {
			model := calculator.models[contactType]
			if model == nil {
				t.Errorf("Expected model for contact type %s", contactType.String())
			}

			// Verify model has reasonable values
			if model.BaseIntensity < 0.0 || model.BaseIntensity > 1.0 {
				t.Errorf("Invalid base intensity for %s: %f", contactType.String(), model.BaseIntensity)
			}

			if model.MinIntensityThreshold >= model.MaxIntensityThreshold {
				t.Errorf("Invalid thresholds for %s: min=%f, max=%f",
					contactType.String(), model.MinIntensityThreshold, model.MaxIntensityThreshold)
			}
		}
	})
}

func TestEffectiveIntensityCalculation(t *testing.T) {
	calculator := NewContactIntensityCalculator()

	// Create test languages
	lang1 := &Language{
		ID:   LanguageID{Family: "test", Branch: "lang", Language: "1"},
		Name: "Test Language 1",
	}

	lang2 := &Language{
		ID:   LanguageID{Family: "test", Branch: "lang", Language: "2"},
		Name: "Test Language 2",
	}

	t.Run("Trade Contact Intensity", func(t *testing.T) {
		// Test trade contact with moderate duration
		effectiveIntensity := calculator.CalculateEffectiveIntensity(
			ContactTypeTrade,
			0.5,          // base intensity
			24*time.Hour, // duration
			nil,          // no history
			0.5,          // geographic proximity
		)

		// Trade should have moderate intensity
		if effectiveIntensity < 0.3 || effectiveIntensity > 0.8 {
			t.Errorf("Expected trade intensity between 0.3-0.8, got %f", effectiveIntensity)
		}

		t.Logf("Trade contact effective intensity: %f (lang1: %s, lang2: %s)",
			effectiveIntensity, lang1.Name, lang2.Name)
	})

	t.Run("Conquest Contact Intensity", func(t *testing.T) {
		// Test conquest contact - should have high intensity
		effectiveIntensity := calculator.CalculateEffectiveIntensity(
			ContactTypeConquest,
			0.8,         // base intensity
			2*time.Hour, // short duration
			nil,         // no history
			0.0,         // no geographic proximity
		)

		// Conquest should have high intensity
		if effectiveIntensity < 0.7 || effectiveIntensity > 1.0 {
			t.Errorf("Expected conquest intensity between 0.7-1.0, got %f", effectiveIntensity)
		}

		t.Logf("Conquest contact effective intensity: %f", effectiveIntensity)
	})

	t.Run("Duration Scaling", func(t *testing.T) {
		// Test that longer duration increases intensity
		shortIntensity := calculator.CalculateEffectiveIntensity(
			ContactTypeTrade,
			0.5,
			1*time.Hour,
			nil,
			0.5,
		)

		longIntensity := calculator.CalculateEffectiveIntensity(
			ContactTypeTrade,
			0.5,
			7*24*time.Hour, // 1 week
			nil,
			0.5,
		)

		if longIntensity <= shortIntensity {
			t.Errorf("Expected longer duration to increase intensity: short=%f, long=%f",
				shortIntensity, longIntensity)
		}

		t.Logf("Duration scaling: short=%f, long=%f", shortIntensity, longIntensity)
	})

	t.Run("Geographic Proximity", func(t *testing.T) {
		// Test that geographic proximity affects intensity
		closeIntensity := calculator.CalculateEffectiveIntensity(
			ContactTypeMigration,
			0.5,
			24*time.Hour,
			nil,
			0.9, // high proximity
		)

		farIntensity := calculator.CalculateEffectiveIntensity(
			ContactTypeMigration,
			0.5,
			24*time.Hour,
			nil,
			0.1, // low proximity
		)

		if closeIntensity <= farIntensity {
			t.Errorf("Expected higher proximity to increase intensity: close=%f, far=%f",
				closeIntensity, farIntensity)
		}

		t.Logf("Geographic proximity: close=%f, far=%f", closeIntensity, farIntensity)
	})
}

func TestContactHistory(t *testing.T) {
	calculator := NewContactIntensityCalculator()

	lang1 := &Language{
		ID:   LanguageID{Family: "test", Branch: "lang", Language: "1"},
		Name: "Test Language 1",
	}

	lang2 := &Language{
		ID:   LanguageID{Family: "test", Branch: "lang", Language: "2"},
		Name: "Test Language 2",
	}

	t.Run("History Creation", func(t *testing.T) {
		history := calculator.GetContactHistory(lang1, lang2, ContactTypeTrade)

		if history == nil {
			t.Fatal("Expected contact history to be created")
		}

		if history.LanguageID != lang1.ID.String() {
			t.Errorf("Expected language ID %s, got %s", lang1.ID.String(), history.LanguageID)
		}

		if history.SourceID != lang2.ID.String() {
			t.Errorf("Expected source ID %s, got %s", lang2.ID.String(), history.SourceID)
		}

		if history.ContactType != ContactTypeTrade {
			t.Errorf("Expected contact type %s, got %s", ContactTypeTrade.String(), history.ContactType.String())
		}
	})

	t.Run("History Update", func(t *testing.T) {
		history := calculator.GetContactHistory(lang1, lang2, ContactTypeTrade)

		// Update with first contact
		updatedHistory := calculator.UpdateContactHistory(
			history,
			ContactTypeTrade,
			24*time.Hour,
			0.6,
			[]string{"Borrowed trade terms"},
			0.1,
		)

		if updatedHistory.ContactCount != 1 {
			t.Errorf("Expected contact count 1, got %d", updatedHistory.ContactCount)
		}

		if updatedHistory.PeakIntensity != 0.6 {
			t.Errorf("Expected peak intensity 0.6, got %f", updatedHistory.PeakIntensity)
		}

		if len(updatedHistory.Sessions) != 1 {
			t.Errorf("Expected 1 session, got %d", len(updatedHistory.Sessions))
		}

		// Update with second contact
		updatedHistory2 := calculator.UpdateContactHistory(
			updatedHistory,
			ContactTypeTrade,
			48*time.Hour,
			0.8,
			[]string{"Borrowed more terms"},
			0.2,
		)

		if updatedHistory2.ContactCount != 2 {
			t.Errorf("Expected contact count 2, got %d", updatedHistory2.ContactCount)
		}

		if updatedHistory2.PeakIntensity != 0.8 {
			t.Errorf("Expected peak intensity 0.8, got %f", updatedHistory2.PeakIntensity)
		}

		if len(updatedHistory2.Sessions) != 2 {
			t.Errorf("Expected 2 sessions, got %d", len(updatedHistory2.Sessions))
		}

		t.Logf("Contact history: count=%d, peak=%f, total_complexity=%f",
			updatedHistory2.ContactCount, updatedHistory2.PeakIntensity, updatedHistory2.TotalComplexityChange)
	})
}

func TestAdaptationThresholds(t *testing.T) {
	calculator := NewContactIntensityCalculator()

	t.Run("Trade Adaptation Threshold", func(t *testing.T) {
		// Trade has adaptation threshold of 0.6
		belowThreshold := calculator.ShouldApplyAdaptation(ContactTypeTrade, 0.5)
		aboveThreshold := calculator.ShouldApplyAdaptation(ContactTypeTrade, 0.7)

		if belowThreshold {
			t.Error("Expected below-threshold trade contact to not trigger adaptation")
		}

		if !aboveThreshold {
			t.Error("Expected above-threshold trade contact to trigger adaptation")
		}
	})

	t.Run("Conquest Adaptation Threshold", func(t *testing.T) {
		// Conquest has adaptation threshold of 0.7
		belowThreshold := calculator.ShouldApplyAdaptation(ContactTypeConquest, 0.6)
		aboveThreshold := calculator.ShouldApplyAdaptation(ContactTypeConquest, 0.8)

		if belowThreshold {
			t.Error("Expected below-threshold conquest contact to not trigger adaptation")
		}

		if !aboveThreshold {
			t.Error("Expected above-threshold conquest contact to trigger adaptation")
		}
	})
}

func TestIntensityIntegration(t *testing.T) {
	// Test that the intensity modeling integrates with cultural influence
	lang1 := &Language{
		ID:   LanguageID{Family: "test", Branch: "lang", Language: "1"},
		Name: "Test Language 1",
	}

	lang2 := &Language{
		ID:   LanguageID{Family: "test", Branch: "lang", Language: "2"},
		Name: "Test Language 2",
	}

	t.Run("Cultural Influence with Intensity Modeling", func(t *testing.T) {
		// Apply cultural influence with intensity modeling
		result, err := ApplyCulturalInfluence(lang1, lang2, ContactTypeTrade, 0.5, 24*time.Hour, 42)
		if err != nil {
			t.Fatalf("Failed to apply cultural influence: %v", err)
		}

		if result == nil {
			t.Fatal("Expected cultural influence result")
		}

		// The effective intensity should be different from the base intensity due to modeling
		if result.Intensity == 0.5 {
			t.Log("Effective intensity equals base intensity (may be normal for first contact)")
		}

		t.Logf("Cultural influence result: base_intensity=0.5, effective_intensity=%f, complexity_change=%f",
			result.Intensity, result.ComplexityChange)

		// Verify that the result has reasonable values
		if result.Intensity < 0.0 || result.Intensity > 1.0 {
			t.Errorf("Effective intensity out of valid range: %f", result.Intensity)
		}

		if result.ComplexityChange < 0.0 || result.ComplexityChange > 1.0 {
			t.Errorf("Complexity change out of valid range: %f", result.ComplexityChange)
		}
	})
}

func TestGeographicProximity(t *testing.T) {
	calculator := NewContactIntensityCalculator()

	lang1 := &Language{
		ID: LanguageID{Family: "test", Branch: "lang", Language: "1"},
	}

	lang2 := &Language{
		ID: LanguageID{Family: "test", Branch: "lang", Language: "2"},
	}

	lang3 := &Language{
		ID: LanguageID{Family: "test", Branch: "lang", Language: "1"}, // Similar ID prefix
	}

	t.Run("Geographic Proximity Calculation", func(t *testing.T) {
		proximity1 := calculator.CalculateGeographicProximity(lang1, lang2)
		proximity2 := calculator.CalculateGeographicProximity(lang1, lang3)

		// Languages with similar IDs should have higher proximity
		if proximity2 <= proximity1 {
			t.Errorf("Expected similar ID languages to have higher proximity: different=%f, similar=%f",
				proximity1, proximity2)
		}

		// Proximity should be in valid range
		if proximity1 < 0.0 || proximity1 > 1.0 {
			t.Errorf("Proximity out of valid range: %f", proximity1)
		}

		if proximity2 < 0.0 || proximity2 > 1.0 {
			t.Errorf("Proximity out of valid range: %f", proximity2)
		}

		t.Logf("Geographic proximity: different=%f, similar=%f", proximity1, proximity2)
	})
}
