package lang

import (
	"testing"
)

// TestSocialStratificationModel tests the social stratification model system.
func TestSocialStratificationModel(t *testing.T) {
	model := NewSocialStratificationModel()

	t.Run("Class Effects", func(t *testing.T) {
		// Test noble class effects
		nobleEffect := model.ClassEffects["noble"]
		if nobleEffect.Description == "" {
			t.Error("Expected noble class effect to have description")
		}
		if nobleEffect.ComplexityImpact != 0.4 {
			t.Errorf("Expected noble complexity impact 0.4, got %f", nobleEffect.ComplexityImpact)
		}
		if nobleEffect.PrestigeLevel != 1.0 {
			t.Errorf("Expected noble prestige level 1.0, got %f", nobleEffect.PrestigeLevel)
		}
		if len(nobleEffect.PhonologicalChanges) == 0 {
			t.Error("Expected noble class to have phonological changes")
		}

		// Test lower class effects
		lowerEffect := model.ClassEffects["lower"]
		if lowerEffect.ComplexityImpact != -0.4 {
			t.Errorf("Expected lower complexity impact -0.4, got %f", lowerEffect.ComplexityImpact)
		}
		if lowerEffect.PrestigeLevel != 0.1 {
			t.Errorf("Expected lower prestige level 0.1, got %f", lowerEffect.PrestigeLevel)
		}

		t.Logf("Class effects: noble=%f (prestige: %f), lower=%f (prestige: %f)",
			nobleEffect.ComplexityImpact, nobleEffect.PrestigeLevel,
			lowerEffect.ComplexityImpact, lowerEffect.PrestigeLevel)
	})

	t.Run("Education Effects", func(t *testing.T) {
		// Test high education effects
		highEffect := model.EducationEffects["high"]
		if highEffect.Description == "" {
			t.Error("Expected high education effect to have description")
		}
		if highEffect.ComplexityImpact != 0.3 {
			t.Errorf("Expected high education complexity impact 0.3, got %f", highEffect.ComplexityImpact)
		}
		if highEffect.FormalityLevel != 0.8 {
			t.Errorf("Expected high education formality level 0.8, got %f", highEffect.FormalityLevel)
		}

		// Test no education effects
		noneEffect := model.EducationEffects["none"]
		if noneEffect.ComplexityImpact != -0.3 {
			t.Errorf("Expected no education complexity impact -0.3, got %f", noneEffect.ComplexityImpact)
		}
		if noneEffect.FormalityLevel != 0.0 {
			t.Errorf("Expected no education formality level 0.0, got %f", noneEffect.FormalityLevel)
		}

		t.Logf("Education effects: high=%f (formality: %f), none=%f (formality: %f)",
			highEffect.ComplexityImpact, highEffect.FormalityLevel,
			noneEffect.ComplexityImpact, noneEffect.FormalityLevel)
	})

	t.Run("Occupation Effects", func(t *testing.T) {
		// Test professional occupation effects
		professionalEffect := model.OccupationEffects["professional"]
		if professionalEffect.Description == "" {
			t.Error("Expected professional occupation effect to have description")
		}
		if professionalEffect.ComplexityImpact != 0.2 {
			t.Errorf("Expected professional complexity impact 0.2, got %f", professionalEffect.ComplexityImpact)
		}
		if professionalEffect.SpecializationLevel != 0.8 {
			t.Errorf("Expected professional specialization level 0.8, got %f", professionalEffect.SpecializationLevel)
		}

		// Test agricultural occupation effects
		agriculturalEffect := model.OccupationEffects["agricultural"]
		if agriculturalEffect.ComplexityImpact != 0.1 {
			t.Errorf("Expected agricultural complexity impact 0.1, got %f", agriculturalEffect.ComplexityImpact)
		}
		if agriculturalEffect.SpecializationLevel != 0.4 {
			t.Errorf("Expected agricultural specialization level 0.4, got %f", agriculturalEffect.SpecializationLevel)
		}

		t.Logf("Occupation effects: professional=%f (specialization: %f), agricultural=%f (specialization: %f)",
			professionalEffect.ComplexityImpact, professionalEffect.SpecializationLevel,
			agriculturalEffect.ComplexityImpact, agriculturalEffect.SpecializationLevel)
	})

	t.Run("Mobility Effects", func(t *testing.T) {
		// Test high mobility effects
		highEffect := model.MobilityEffects["high"]
		if highEffect.Description == "" {
			t.Error("Expected high mobility effect to have description")
		}
		if highEffect.ComplexityImpact != 0.1 {
			t.Errorf("Expected high mobility complexity impact 0.1, got %f", highEffect.ComplexityImpact)
		}
		if highEffect.AdaptationLevel != 0.8 {
			t.Errorf("Expected high mobility adaptation level 0.8, got %f", highEffect.AdaptationLevel)
		}

		// Test no mobility effects
		noneEffect := model.MobilityEffects["none"]
		if noneEffect.ComplexityImpact != 0.2 {
			t.Errorf("Expected no mobility complexity impact 0.2, got %f", noneEffect.ComplexityImpact)
		}
		if noneEffect.AdaptationLevel != 0.0 {
			t.Errorf("Expected no mobility adaptation level 0.0, got %f", noneEffect.AdaptationLevel)
		}

		t.Logf("Mobility effects: high=%f (adaptation: %f), none=%f (adaptation: %f)",
			highEffect.ComplexityImpact, highEffect.AdaptationLevel,
			noneEffect.ComplexityImpact, noneEffect.AdaptationLevel)
	})

	t.Run("Prestige Effects", func(t *testing.T) {
		// Test high prestige effects
		highEffect := model.PrestigeEffects["high"]
		if highEffect.Description == "" {
			t.Error("Expected high prestige effect to have description")
		}
		if highEffect.ComplexityImpact != 0.3 {
			t.Errorf("Expected high prestige complexity impact 0.3, got %f", highEffect.ComplexityImpact)
		}
		if highEffect.InfluenceLevel != 0.9 {
			t.Errorf("Expected high prestige influence level 0.9, got %f", highEffect.InfluenceLevel)
		}

		// Test low prestige effects
		lowEffect := model.PrestigeEffects["low"]
		if lowEffect.ComplexityImpact != -0.1 {
			t.Errorf("Expected low prestige complexity impact -0.1, got %f", lowEffect.ComplexityImpact)
		}
		if lowEffect.InfluenceLevel != 0.2 {
			t.Errorf("Expected low prestige influence level 0.2, got %f", lowEffect.InfluenceLevel)
		}

		t.Logf("Prestige effects: high=%f (influence: %f), low=%f (influence: %f)",
			highEffect.ComplexityImpact, highEffect.InfluenceLevel,
			lowEffect.ComplexityImpact, lowEffect.InfluenceLevel)
	})

	t.Run("Gender Effects", func(t *testing.T) {
		// Test male gender effects
		maleEffect := model.GenderEffects["male"]
		if maleEffect.Description == "" {
			t.Error("Expected male gender effect to have description")
		}
		if maleEffect.ComplexityImpact != 0.0 {
			t.Errorf("Expected male complexity impact 0.0, got %f", maleEffect.ComplexityImpact)
		}
		if maleEffect.DistinctionLevel != 0.3 {
			t.Errorf("Expected male distinction level 0.3, got %f", maleEffect.DistinctionLevel)
		}

		// Test neutral gender effects
		neutralEffect := model.GenderEffects["neutral"]
		if neutralEffect.ComplexityImpact != 0.0 {
			t.Errorf("Expected neutral complexity impact 0.0, got %f", neutralEffect.ComplexityImpact)
		}
		if neutralEffect.DistinctionLevel != 0.0 {
			t.Errorf("Expected neutral distinction level 0.0, got %f", neutralEffect.DistinctionLevel)
		}

		t.Logf("Gender effects: male=%f (distinction: %f), neutral=%f (distinction: %f)",
			maleEffect.ComplexityImpact, maleEffect.DistinctionLevel,
			neutralEffect.ComplexityImpact, neutralEffect.DistinctionLevel)
	})
}

// TestSocialStratificationApplication tests the application of social stratification to languages.
func TestSocialStratificationApplication(t *testing.T) {
	model := NewSocialStratificationModel()

	// Create a test language
	baseLang, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	t.Run("Noble Class with High Education", func(t *testing.T) {
		changes := model.ApplySocialStratification(
			baseLang,
			"noble",        // social class
			"high",         // education level
			"professional", // occupation
			"none",         // social mobility
			"high",         // prestige
			"neutral",      // gender
		)

		// Should have changes for all 6 social factors
		expectedChanges := 6
		if len(changes) != expectedChanges {
			t.Errorf("Expected %d social changes, got %d", expectedChanges, len(changes))
		}

		// Verify social class change
		var classChange LinguisticChange
		for _, change := range changes {
			if change.Trigger == "social_class_influence" {
				classChange = change
				break
			}
		}
		if classChange.Trigger != "social_class_influence" {
			t.Error("Expected social class influence change to be applied")
		}
		if classChange.Description == "" {
			t.Error("Expected social class change to have description")
		}

		// Verify education change
		var educationChange LinguisticChange
		for _, change := range changes {
			if change.Trigger == "education_influence" {
				educationChange = change
				break
			}
		}
		if educationChange.Trigger != "education_influence" {
			t.Error("Expected education influence change to be applied")
		}

		t.Logf("Applied %d social changes for noble class with high education", len(changes))
	})

	t.Run("Working Class with Low Education", func(t *testing.T) {
		changes := model.ApplySocialStratification(
			baseLang,
			"working",   // social class
			"low",       // education level
			"unskilled", // occupation
			"medium",    // social mobility
			"low",       // prestige
			"neutral",   // gender
		)

		// Should have changes for all 6 social factors
		expectedChanges := 6
		if len(changes) != expectedChanges {
			t.Errorf("Expected %d social changes, got %d", expectedChanges, len(changes))
		}

		// Verify working class effect
		var classChange LinguisticChange
		for _, change := range changes {
			if change.Trigger == "social_class_influence" {
				classChange = change
				break
			}
		}
		if classChange.Trigger != "social_class_influence" {
			t.Error("Expected social class influence change to be applied")
		}

		t.Logf("Applied %d social changes for working class with low education", len(changes))
	})

	t.Run("Lower Class with No Education", func(t *testing.T) {
		changes := model.ApplySocialStratification(
			baseLang,
			"lower",     // social class
			"none",      // education level
			"unskilled", // occupation
			"high",      // social mobility
			"low",       // prestige
			"neutral",   // gender
		)

		// Should have changes for all 6 social factors
		expectedChanges := 6
		if len(changes) != expectedChanges {
			t.Errorf("Expected %d social changes, got %d", expectedChanges, len(changes))
		}

		// Verify lower class effect
		var classChange LinguisticChange
		for _, change := range changes {
			if change.Trigger == "social_class_influence" {
				classChange = change
				break
			}
		}
		if classChange.Trigger != "social_class_influence" {
			t.Error("Expected social class influence change to be applied")
		}

		t.Logf("Applied %d social changes for lower class with no education", len(changes))
	})
}

// TestSocialStratificationIntegration tests the integration of social stratification with dialect creation.
func TestSocialStratificationIntegration(t *testing.T) {
	t.Run("Social Dialect with Full Stratification", func(t *testing.T) {
		// Create a base language
		baseLang, err := CreateRandomLanguage("test", "ancient", 42)
		if err != nil {
			t.Fatalf("Failed to create base language: %v", err)
		}

		// Create the social dialect
		dialectID := LanguageID{
			Family:   "test",
			Branch:   "social",
			Language: "dialect",
			Dialect:  "noble",
		}

		dialect, err := CreateSocialDialect(baseLang, dialectID, "Noble Dialect", "noble", 0.8, "middle_evolution", 43)
		if err != nil {
			t.Fatalf("Failed to create social dialect: %v", err)
		}

		// Verify dialect properties
		if dialect.Name != "Noble Dialect" {
			t.Errorf("Expected dialect name 'Noble Dialect', got '%s'", dialect.Name)
		}

		// Verify social stratification changes were applied
		socialChanges := dialect.GetLinguisticChanges(LinguisticChangeTypeDialectal)
		if len(socialChanges) < 7 { // 6 social + 1 formation
			t.Errorf("Expected at least 7 dialect changes (6 social + 1 formation), got %d", len(socialChanges))
		}

		// Verify specific social changes
		var classChange, educationChange, occupationChange LinguisticChange
		for _, change := range socialChanges {
			switch change.Trigger {
			case "social_class_influence":
				classChange = change
			case "education_influence":
				educationChange = change
			case "occupation_influence":
				occupationChange = change
			}
		}

		if classChange.Trigger != "social_class_influence" {
			t.Error("Expected social class influence change to be applied")
		}
		if educationChange.Trigger != "education_influence" {
			t.Error("Expected education influence change to be applied")
		}
		if occupationChange.Trigger != "occupation_influence" {
			t.Error("Expected occupation influence change to be applied")
		}

		// Verify complexity changes
		totalComplexityChange := float32(0.0)
		for _, change := range socialChanges {
			totalComplexityChange += change.ComplexityChange
		}

		t.Logf("Dialect created with %d social changes, total complexity impact: %f", len(socialChanges), totalComplexityChange)
		t.Logf("Class change: %s", classChange.Description)
		t.Logf("Education change: %s", educationChange.Description)
		t.Logf("Occupation change: %s", occupationChange.Description)
	})
}

// TestHelperFunctions tests the helper functions for determining social factors.
func TestHelperFunctions(t *testing.T) {
	t.Run("Education Level for Class", func(t *testing.T) {
		if getEducationLevelForClass("noble") != "high" {
			t.Error("Expected noble class to have high education level")
		}
		if getEducationLevelForClass("upper") != "high" {
			t.Error("Expected upper class to have high education level")
		}
		if getEducationLevelForClass("middle") != "medium" {
			t.Error("Expected middle class to have medium education level")
		}
		if getEducationLevelForClass("working") != "low" {
			t.Error("Expected working class to have low education level")
		}
		if getEducationLevelForClass("lower") != "none" {
			t.Error("Expected lower class to have no education level")
		}
		if getEducationLevelForClass("unknown") != "medium" {
			t.Error("Expected unknown class to have medium education level")
		}
	})

	t.Run("Occupation for Class", func(t *testing.T) {
		if getOccupationForClass("noble") != "professional" {
			t.Error("Expected noble class to have professional occupation")
		}
		if getOccupationForClass("upper") != "professional" {
			t.Error("Expected upper class to have professional occupation")
		}
		if getOccupationForClass("middle") != "skilled" {
			t.Error("Expected middle class to have skilled occupation")
		}
		if getOccupationForClass("working") != "unskilled" {
			t.Error("Expected working class to have unskilled occupation")
		}
		if getOccupationForClass("lower") != "unskilled" {
			t.Error("Expected lower class to have unskilled occupation")
		}
	})

	t.Run("Social Mobility for Class", func(t *testing.T) {
		if getSocialMobilityForClass("noble") != "none" {
			t.Error("Expected noble class to have no social mobility")
		}
		if getSocialMobilityForClass("upper") != "low" {
			t.Error("Expected upper class to have low social mobility")
		}
		if getSocialMobilityForClass("middle") != "medium" {
			t.Error("Expected middle class to have medium social mobility")
		}
		if getSocialMobilityForClass("working") != "medium" {
			t.Error("Expected working class to have medium social mobility")
		}
		if getSocialMobilityForClass("lower") != "high" {
			t.Error("Expected lower class to have high social mobility")
		}
	})

	t.Run("Prestige for Class", func(t *testing.T) {
		if getPrestigeForClass("noble") != "high" {
			t.Error("Expected noble class to have high prestige")
		}
		if getPrestigeForClass("upper") != "high" {
			t.Error("Expected upper class to have high prestige")
		}
		if getPrestigeForClass("middle") != "medium" {
			t.Error("Expected middle class to have medium prestige")
		}
		if getPrestigeForClass("working") != "low" {
			t.Error("Expected working class to have low prestige")
		}
		if getPrestigeForClass("lower") != "low" {
			t.Error("Expected lower class to have low prestige")
		}
	})

	t.Logf("Successfully tested all helper functions for social stratification")
}
