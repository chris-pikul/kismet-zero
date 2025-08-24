package lang

import (
	"fmt"
	"testing"
)

// TestMutualIntelligibilityEngine tests the mutual intelligibility engine system.
func TestMutualIntelligibilityEngine(t *testing.T) {
	engine := NewMutualIntelligibilityEngine()

	t.Run("Engine Configuration", func(t *testing.T) {
		if engine.PhonologicalWeight != 0.3 {
			t.Errorf("Expected phonological weight 0.3, got %f", engine.PhonologicalWeight)
		}
		if engine.LexicalWeight != 0.25 {
			t.Errorf("Expected lexical weight 0.25, got %f", engine.LexicalWeight)
		}
		if engine.GrammaticalWeight != 0.25 {
			t.Errorf("Expected grammatical weight 0.25, got %f", engine.GrammaticalWeight)
		}
		if engine.HistoricalWeight != 0.1 {
			t.Errorf("Expected historical weight 0.1, got %f", engine.HistoricalWeight)
		}
		if engine.GeographicWeight != 0.05 {
			t.Errorf("Expected geographic weight 0.05, got %f", engine.GeographicWeight)
		}
		if engine.ContactWeight != 0.05 {
			t.Errorf("Expected contact weight 0.05, got %f", engine.ContactWeight)
		}

		t.Logf("Mutual intelligibility engine configured with: phonological=%f, lexical=%f, grammatical=%f, historical=%f, geographic=%f, contact=%f",
			engine.PhonologicalWeight, engine.LexicalWeight, engine.GrammaticalWeight,
			engine.HistoricalWeight, engine.GeographicWeight, engine.ContactWeight)
	})
}

// TestIntelligibilityLevels tests the intelligibility level constants and determination.
func TestIntelligibilityLevels(t *testing.T) {
	t.Run("Intelligibility Level Constants", func(t *testing.T) {
		if IntelligibilityLevelNone != "none" {
			t.Errorf("Expected IntelligibilityLevelNone 'none', got '%s'", IntelligibilityLevelNone)
		}
		if IntelligibilityLevelMinimal != "minimal" {
			t.Errorf("Expected IntelligibilityLevelMinimal 'minimal', got '%s'", IntelligibilityLevelMinimal)
		}
		if IntelligibilityLevelLow != "low" {
			t.Errorf("Expected IntelligibilityLevelLow 'low', got '%s'", IntelligibilityLevelLow)
		}
		if IntelligibilityLevelModerate != "moderate" {
			t.Errorf("Expected IntelligibilityLevelModerate 'moderate', got '%s'", IntelligibilityLevelModerate)
		}
		if IntelligibilityLevelHigh != "high" {
			t.Errorf("Expected IntelligibilityLevelHigh 'high', got '%s'", IntelligibilityLevelHigh)
		}
		if IntelligibilityLevelVeryHigh != "very_high" {
			t.Errorf("Expected IntelligibilityLevelVeryHigh 'very_high', got '%s'", IntelligibilityLevelVeryHigh)
		}
		if IntelligibilityLevelNearNative != "near_native" {
			t.Errorf("Expected IntelligibilityLevelNearNative 'near_native', got '%s'", IntelligibilityLevelNearNative)
		}

		t.Logf("All intelligibility level constants properly defined")
	})

	t.Run("Intelligibility Level Determination", func(t *testing.T) {
		engine := NewMutualIntelligibilityEngine()

		testCases := []struct {
			score float32
			level IntelligibilityLevel
		}{
			{0.99, IntelligibilityLevelNearNative},
			{0.85, IntelligibilityLevelVeryHigh},
			{0.75, IntelligibilityLevelHigh},
			{0.50, IntelligibilityLevelModerate},
			{0.30, IntelligibilityLevelLow},
			{0.15, IntelligibilityLevelMinimal},
			{0.05, IntelligibilityLevelNone},
		}

		for _, tc := range testCases {
			result := engine.determineIntelligibilityLevel(tc.score)
			if result != tc.level {
				t.Errorf("Expected level %s for score %f, got %s", tc.level, tc.score, result)
			}
		}

		t.Logf("Intelligibility level determination working correctly")
	})
}

// TestMutualIntelligibilityCalculation tests the core mutual intelligibility calculation.
func TestMutualIntelligibilityCalculation(t *testing.T) {
	engine := NewMutualIntelligibilityEngine()

	t.Run("Same Language", func(t *testing.T) {
		// Create a language
		lang1, err := CreateRandomLanguage("test", "ancient", 42)
		if err != nil {
			t.Fatalf("Failed to create test language: %v", err)
		}

		// Calculate intelligibility with itself
		result, err := engine.CalculateMutualIntelligibility(lang1, lang1)
		if err != nil {
			t.Fatalf("Failed to calculate intelligibility: %v", err)
		}

		if result.OverallScore != 1.0 {
			t.Errorf("Expected overall score 1.0 for same language, got %f", result.OverallScore)
		}
		if result.IntelligibilityLevel != string(IntelligibilityLevelNearNative) {
			t.Errorf("Expected intelligibility level 'near_native', got '%s'", result.IntelligibilityLevel)
		}
		if result.Confidence != 1.0 {
			t.Errorf("Expected confidence 1.0 for same language, got %f", result.Confidence)
		}

		t.Logf("Same language intelligibility: score=%f, level=%s, confidence=%f",
			result.OverallScore, result.IntelligibilityLevel, result.Confidence)
	})

	t.Run("Different Languages", func(t *testing.T) {
		// Create two different languages
		lang1, err := CreateRandomLanguage("test1", "ancient", 42)
		if err != nil {
			t.Fatalf("Failed to create first test language: %v", err)
		}

		lang2, err := CreateRandomLanguage("test2", "ancient", 43)
		if err != nil {
			t.Fatalf("Failed to create second test language: %v", err)
		}

		// Calculate intelligibility between them
		result, err := engine.CalculateMutualIntelligibility(lang1, lang2)
		if err != nil {
			t.Fatalf("Failed to calculate intelligibility: %v", err)
		}

		if result.OverallScore < 0.0 || result.OverallScore > 1.0 {
			t.Errorf("Expected overall score between 0.0 and 1.0, got %f", result.OverallScore)
		}
		if result.SourceLanguage != lang1.Name {
			t.Errorf("Expected source language '%s', got '%s'", lang1.Name, result.SourceLanguage)
		}
		if result.TargetLanguage != lang2.Name {
			t.Errorf("Expected target language '%s', got '%s'", lang2.Name, result.TargetLanguage)
		}
		if result.Confidence < 0.0 || result.Confidence > 1.0 {
			t.Errorf("Expected confidence between 0.0 and 1.0, got %f", result.Confidence)
		}

		t.Logf("Different languages intelligibility: score=%f, level=%s, confidence=%f",
			result.OverallScore, result.IntelligibilityLevel, result.Confidence)
		t.Logf("Description: %s", result.Description)
		t.Logf("Factors: %v", result.Factors)
	})

	t.Run("Nil Language Error", func(t *testing.T) {
		lang1, err := CreateRandomLanguage("test", "ancient", 42)
		if err != nil {
			t.Fatalf("Failed to create test language: %v", err)
		}

		// Test with nil language
		_, err = engine.CalculateMutualIntelligibility(nil, lang1)
		if err == nil {
			t.Error("Expected error when first language is nil")
		}

		_, err = engine.CalculateMutualIntelligibility(lang1, nil)
		if err == nil {
			t.Error("Expected error when second language is nil")
		}

		t.Logf("Successfully caught nil language errors")
	})
}

// TestComponentScores tests the individual component score calculations.
func TestComponentScores(t *testing.T) {
	engine := NewMutualIntelligibilityEngine()

	t.Run("Component Score Ranges", func(t *testing.T) {
		lang1, err := CreateRandomLanguage("test1", "ancient", 42)
		if err != nil {
			t.Fatalf("Failed to create first test language: %v", err)
		}

		lang2, err := CreateRandomLanguage("test2", "ancient", 43)
		if err != nil {
			t.Fatalf("Failed to create second test language: %v", err)
		}

		result, err := engine.CalculateMutualIntelligibility(lang1, lang2)
		if err != nil {
			t.Fatalf("Failed to calculate intelligibility: %v", err)
		}

		// Check that all component scores are within valid range
		componentScores := []struct {
			name  string
			score float32
		}{
			{"Phonological", result.PhonologicalScore},
			{"Lexical", result.LexicalScore},
			{"Grammatical", result.GrammaticalScore},
			{"Historical", result.HistoricalScore},
			{"Geographic", result.GeographicScore},
			{"Contact", result.ContactScore},
		}

		for _, component := range componentScores {
			if component.score < 0.0 || component.score > 1.0 {
				t.Errorf("Expected %s score between 0.0 and 1.0, got %f", component.name, component.score)
			}
		}

		t.Logf("All component scores within valid range (0.0-1.0)")
		t.Logf("Phonological: %f, Lexical: %f, Grammatical: %f",
			result.PhonologicalScore, result.LexicalScore, result.GrammaticalScore)
		t.Logf("Historical: %f, Geographic: %f, Contact: %f",
			result.HistoricalScore, result.GeographicScore, result.ContactScore)
	})
}

// TestBatchIntelligibility tests batch intelligibility calculations.
func TestBatchIntelligibility(t *testing.T) {
	engine := NewMutualIntelligibilityEngine()

	t.Run("Batch Calculation", func(t *testing.T) {
		// Create multiple languages
		languages := make([]*Language, 0)
		for i := 0; i < 3; i++ {
			lang, err := CreateRandomLanguage(fmt.Sprintf("test%d", i), "ancient", int64(42+i))
			if err != nil {
				t.Fatalf("Failed to create test language %d: %v", i, err)
			}
			languages = append(languages, lang)
		}

		// Calculate batch intelligibility
		results, err := engine.BatchCalculateIntelligibility(languages)
		if err != nil {
			t.Fatalf("Failed to calculate batch intelligibility: %v", err)
		}

		// With 3 languages, we should have 3 choose 2 = 3 pairs
		expectedPairs := 3
		if len(results) != expectedPairs {
			t.Errorf("Expected %d language pairs, got %d", expectedPairs, len(results))
		}

		// Verify all results are valid
		for i, result := range results {
			if result.OverallScore < 0.0 || result.OverallScore > 1.0 {
				t.Errorf("Result %d: Expected overall score between 0.0 and 1.0, got %f", i, result.OverallScore)
			}
			if result.Confidence < 0.0 || result.Confidence > 1.0 {
				t.Errorf("Result %d: Expected confidence between 0.0 and 1.0, got %f", i, result.Confidence)
			}
		}

		t.Logf("Successfully calculated intelligibility for %d language pairs", len(results))
		for i, result := range results {
			t.Logf("Pair %d: %s ↔ %s = %.1f%% (%s)",
				i+1, result.SourceLanguage, result.TargetLanguage,
				result.OverallScore*100, result.IntelligibilityLevel)
		}
	})
}

// TestIntelligibilitySorting tests the sorting and filtering functions.
func TestIntelligibilitySorting(t *testing.T) {
	engine := NewMutualIntelligibilityEngine()

	t.Run("Most Intelligible Pairs", func(t *testing.T) {
		// Create test results
		results := []*IntelligibilityResult{
			{SourceLanguage: "A", TargetLanguage: "B", OverallScore: 0.3},
			{SourceLanguage: "C", TargetLanguage: "D", OverallScore: 0.8},
			{SourceLanguage: "E", TargetLanguage: "F", OverallScore: 0.5},
		}

		// Find most intelligible pairs
		mostIntelligible := engine.FindMostIntelligiblePairs(results, 2)
		if len(mostIntelligible) != 2 {
			t.Errorf("Expected 2 most intelligible pairs, got %d", len(mostIntelligible))
		}

		// Check that they're sorted by score (highest first)
		if mostIntelligible[0].OverallScore != 0.8 {
			t.Errorf("Expected highest score 0.8, got %f", mostIntelligible[0].OverallScore)
		}
		if mostIntelligible[1].OverallScore != 0.5 {
			t.Errorf("Expected second highest score 0.5, got %f", mostIntelligible[1].OverallScore)
		}

		t.Logf("Most intelligible pairs correctly sorted: %s↔%s (%.1f%%), %s↔%s (%.1f%%)",
			mostIntelligible[0].SourceLanguage, mostIntelligible[0].TargetLanguage, mostIntelligible[0].OverallScore*100,
			mostIntelligible[1].SourceLanguage, mostIntelligible[1].TargetLanguage, mostIntelligible[1].OverallScore*100)
	})

	t.Run("Least Intelligible Pairs", func(t *testing.T) {
		// Create test results
		results := []*IntelligibilityResult{
			{SourceLanguage: "A", TargetLanguage: "B", OverallScore: 0.3},
			{SourceLanguage: "C", TargetLanguage: "D", OverallScore: 0.8},
			{SourceLanguage: "E", TargetLanguage: "F", OverallScore: 0.5},
		}

		// Find least intelligible pairs
		leastIntelligible := engine.FindLeastIntelligiblePairs(results, 2)
		if len(leastIntelligible) != 2 {
			t.Errorf("Expected 2 least intelligible pairs, got %d", len(leastIntelligible))
		}

		// Check that they're sorted by score (lowest first)
		if leastIntelligible[0].OverallScore != 0.3 {
			t.Errorf("Expected lowest score 0.3, got %f", leastIntelligible[0].OverallScore)
		}
		if leastIntelligible[1].OverallScore != 0.5 {
			t.Errorf("Expected second lowest score 0.5, got %f", leastIntelligible[1].OverallScore)
		}

		t.Logf("Least intelligible pairs correctly sorted: %s↔%s (%.1f%%), %s↔%s (%.1f%%)",
			leastIntelligible[0].SourceLanguage, leastIntelligible[0].TargetLanguage, leastIntelligible[0].OverallScore*100,
			leastIntelligible[1].SourceLanguage, leastIntelligible[1].TargetLanguage, leastIntelligible[1].OverallScore*100)
	})
}

// TestIntelligibilityIntegration tests the integration with the language system.
func TestIntelligibilityIntegration(t *testing.T) {
	engine := NewMutualIntelligibilityEngine()

	t.Run("Language Family Intelligibility", func(t *testing.T) {
		// Create a language family
		protoLang, err := CreateRandomLanguage("proto", "ancient", 42)
		if err != nil {
			t.Fatalf("Failed to create proto language: %v", err)
		}

		// Create dialects
		dialect1, err := CreateGeographicDialect(
			protoLang,
			LanguageID{Family: "proto", Branch: "proto", Language: "proto", Dialect: "dialect1"},
			"Coastal Dialect",
			GeographicRegion{ID: "coastal", Name: "Coastal", Climate: "temperate", Terrain: "coastal"},
			"ancient",
			43,
		)
		if err != nil {
			t.Fatalf("Failed to create first dialect: %v", err)
		}

		dialect2, err := CreateGeographicDialect(
			protoLang,
			LanguageID{Family: "proto", Branch: "proto", Language: "proto", Dialect: "dialect2"},
			"Inland Dialect",
			GeographicRegion{ID: "inland", Name: "Inland", Climate: "temperate", Terrain: "plains"},
			"ancient",
			44,
		)
		if err != nil {
			t.Fatalf("Failed to create second dialect: %v", err)
		}

		// Calculate intelligibility between proto and dialects
		protoDialect1, err := engine.CalculateMutualIntelligibility(protoLang, dialect1)
		if err != nil {
			t.Fatalf("Failed to calculate proto-dialect1 intelligibility: %v", err)
		}

		protoDialect2, err := engine.CalculateMutualIntelligibility(protoLang, dialect2)
		if err != nil {
			t.Fatalf("Failed to calculate proto-dialect2 intelligibility: %v", err)
		}

		dialect1Dialect2, err := engine.CalculateMutualIntelligibility(dialect1, dialect2)
		if err != nil {
			t.Fatalf("Failed to calculate dialect1-dialect2 intelligibility: %v", err)
		}

		// Dialects should have high intelligibility with proto language
		if protoDialect1.OverallScore < 0.6 {
			t.Errorf("Expected high intelligibility between proto and dialect1, got %f", protoDialect1.OverallScore)
		}
		if protoDialect2.OverallScore < 0.6 {
			t.Errorf("Expected high intelligibility between proto and dialect2, got %f", protoDialect2.OverallScore)
		}

		// Dialects should have moderate to high intelligibility with each other
		if dialect1Dialect2.OverallScore < 0.4 {
			t.Errorf("Expected moderate intelligibility between dialects, got %f", dialect1Dialect2.OverallScore)
		}

		t.Logf("Language family intelligibility analysis:")
		t.Logf("Proto ↔ Coastal Dialect: %.1f%% (%s)", protoDialect1.OverallScore*100, protoDialect1.IntelligibilityLevel)
		t.Logf("Proto ↔ Inland Dialect: %.1f%% (%s)", protoDialect2.OverallScore*100, protoDialect2.IntelligibilityLevel)
		t.Logf("Coastal ↔ Inland Dialect: %.1f%% (%s)", dialect1Dialect2.OverallScore*100, dialect1Dialect2.IntelligibilityLevel)
	})
}
