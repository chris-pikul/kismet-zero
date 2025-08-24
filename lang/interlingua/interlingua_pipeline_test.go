package interlingua

import (
	"testing"
	"time"

	"github.com/chris-pikul/kismet-zero/lang/interlingua"
)

// TestInterlinguaPipeline tests the complete interlingua pipeline system.
func TestInterlinguaPipeline(t *testing.T) {
	t.Run("Pipeline Creation", func(t *testing.T) {
		services := interlingua.NewInterlinguaServices()
		pipeline := NewInterlinguaPipeline(services)

		if pipeline == nil {
			t.Error("Expected pipeline to be created")
		}
		if pipeline.services != services {
			t.Error("Expected pipeline to have the provided services")
		}

		t.Logf("Successfully created interlingua pipeline")
	})
}

// TestTranslationResult tests the translation result structure.
func TestTranslationResult(t *testing.T) {
	t.Run("Result Structure", func(t *testing.T) {
		result := &TranslationResult{
			SourceText:     "Hello world",
			TargetText:     "Bonjour monde",
			SourceLanguage: "en",
			TargetLanguage: "fr",
			Confidence:     0.95,
			ProcessingTime: 100 * time.Millisecond,
		}

		if result.SourceText != "Hello world" {
			t.Errorf("Expected source text 'Hello world', got '%s'", result.SourceText)
		}
		if result.TargetText != "Bonjour monde" {
			t.Errorf("Expected target text 'Bonjour monde', got '%s'", result.TargetText)
		}
		if result.SourceLanguage != "en" {
			t.Errorf("Expected source language 'en', got '%s'", result.SourceLanguage)
		}
		if result.TargetLanguage != "fr" {
			t.Errorf("Expected target language 'fr', got '%s'", result.TargetText)
		}
		if result.Confidence != 0.95 {
			t.Errorf("Expected confidence 0.95, got %f", result.Confidence)
		}
		if result.ProcessingTime != 100*time.Millisecond {
			t.Errorf("Expected processing time 100ms, got %v", result.ProcessingTime)
		}

		t.Logf("Successfully tested translation result structure")
	})
}

// TestInterlinguaPipelineTranslation tests the core translation functionality.
func TestInterlinguaPipelineTranslation(t *testing.T) {
	t.Run("Basic Translation", func(t *testing.T) {
		// Create services with English support
		services := interlingua.NewInterlinguaServices()

		// Create English analyzer and realizer
		conceptResolver := interlingua.NewInMemoryConceptResolver()
		englishAnalyzer := interlingua.NewEnglishAnalyzer(conceptResolver)
		englishRealizer := interlingua.NewEnglishRealizer(42, conceptResolver)

		services.AddAnalyzer("en", englishAnalyzer)
		services.AddRealizer("en", englishRealizer)

		// Create pipeline
		pipeline := NewInterlinguaPipeline(services)

		// Test English to English translation (should work)
		result, err := pipeline.Translate("I walk", "en", "en")
		if err != nil {
			t.Fatalf("Failed to translate English to English: %v", err)
		}

		if result.SourceText != "I walk" {
			t.Errorf("Expected source text 'I walk', got '%s'", result.SourceText)
		}
		if result.SourceLanguage != "en" {
			t.Errorf("Expected source language 'en', got '%s'", result.SourceLanguage)
		}
		if result.TargetLanguage != "en" {
			t.Errorf("Expected target language 'en', got '%s'", result.TargetText)
		}
		if result.Confidence <= 0 {
			t.Errorf("Expected positive confidence, got %f", result.Confidence)
		}
		if result.ProcessingTime <= 0 {
			t.Errorf("Expected positive processing time, got %v", result.ProcessingTime)
		}

		t.Logf("Successfully translated 'I walk' to English with confidence %f", result.Confidence)
		t.Logf("Translation result: %s", result.TargetText)
	})

	t.Run("Translation with Missing Analyzer", func(t *testing.T) {
		services := interlingua.NewInterlinguaServices()
		pipeline := NewInterlinguaPipeline(services)

		// Try to translate from a language without an analyzer
		_, err := pipeline.Translate("Hello", "fr", "en")
		if err == nil {
			t.Error("Expected error when analyzer is missing")
		}
		if err.Error() != "no analyzer available for source language: fr" {
			t.Errorf("Expected specific error message, got: %v", err)
		}

		t.Logf("Successfully caught missing analyzer error")
	})

	t.Run("Translation with Missing Realizer", func(t *testing.T) {
		services := interlingua.NewInterlinguaServices()

		// Add analyzer but no realizer
		conceptResolver := interlingua.NewInMemoryConceptResolver()
		englishAnalyzer := interlingua.NewEnglishAnalyzer(conceptResolver)
		services.AddAnalyzer("en", englishAnalyzer)

		pipeline := NewInterlinguaPipeline(services)

		// Try to translate to a language without a realizer
		_, err := pipeline.Translate("Hello", "en", "fr")
		if err == nil {
			t.Error("Expected error when realizer is missing")
		}
		if err.Error() != "no realizer available for target language: fr" {
			t.Errorf("Expected specific error message, got: %v", err)
		}

		t.Logf("Successfully caught missing realizer error")
	})
}

// TestInterlinguaPipelineConfidence tests the confidence calculation system.
func TestInterlinguaPipelineConfidence(t *testing.T) {
	t.Run("Confidence Calculation", func(t *testing.T) {
		services := interlingua.NewInterlinguaServices()

		conceptResolver := interlingua.NewInMemoryConceptResolver()
		englishAnalyzer := interlingua.NewEnglishAnalyzer(conceptResolver)
		englishRealizer := interlingua.NewEnglishRealizer(42, conceptResolver)

		services.AddAnalyzer("en", englishAnalyzer)
		services.AddRealizer("en", englishRealizer)

		pipeline := NewInterlinguaPipeline(services)

		// Test confidence calculation with different inputs
		testCases := []struct {
			input         string
			expected      string
			minConfidence float32
		}{
			{"I walk", "I walk", 0.8},
			{"you run", "you run", 0.8},
			{"he moves", "he moves", 0.8},
		}

		for _, tc := range testCases {
			result, err := pipeline.Translate(tc.input, "en", "en")
			if err != nil {
				t.Fatalf("Failed to translate '%s': %v", tc.input, err)
			}

			if result.Confidence < tc.minConfidence {
				t.Errorf("Expected confidence >= %f for '%s', got %f", tc.minConfidence, tc.input, result.Confidence)
			}

			t.Logf("'%s' -> '%s' with confidence %f", tc.input, result.TargetText, result.Confidence)
		}
	})
}

// TestInterlinguaPipelineBatch tests batch translation functionality.
func TestInterlinguaPipelineBatch(t *testing.T) {
	t.Run("Batch Translation", func(t *testing.T) {
		services := interlingua.NewInterlinguaServices()

		conceptResolver := interlingua.NewInMemoryConceptResolver()
		englishAnalyzer := interlingua.NewEnglishAnalyzer(conceptResolver)
		englishRealizer := interlingua.NewEnglishRealizer(42, conceptResolver)

		services.AddAnalyzer("en", englishAnalyzer)
		services.AddRealizer("en", englishRealizer)

		pipeline := NewInterlinguaPipeline(services)

		// Test batch translation
		batch := []struct {
			SourceText string
			SourceLang string
			TargetLang string
		}{
			{"I walk", "en", "en"},
			{"you run", "en", "en"},
			{"he moves", "en", "en"},
		}

		results, err := pipeline.BatchTranslate(batch)
		if err != nil {
			t.Fatalf("Failed to perform batch translation: %v", err)
		}

		if len(results) != len(batch) {
			t.Errorf("Expected %d results, got %d", len(batch), len(results))
		}

		for i, result := range results {
			if result.SourceText != batch[i].SourceText {
				t.Errorf("Result %d: expected source text '%s', got '%s'", i, batch[i].SourceText, result.SourceText)
			}
			if result.Confidence <= 0 {
				t.Errorf("Result %d: expected positive confidence, got %f", i, result.Confidence)
			}
		}

		t.Logf("Successfully performed batch translation of %d texts", len(batch))
	})
}

// TestInterlinguaPipelineLanguagePairs tests the language pair discovery.
func TestInterlinguaPipelineLanguagePairs(t *testing.T) {
	t.Run("Supported Language Pairs", func(t *testing.T) {
		services := interlingua.NewInterlinguaServices()

		conceptResolver := interlingua.NewInMemoryConceptResolver()
		englishAnalyzer := interlingua.NewEnglishAnalyzer(conceptResolver)
		englishRealizer := interlingua.NewEnglishRealizer(42, conceptResolver)

		services.AddAnalyzer("en", englishAnalyzer)
		services.AddRealizer("en", englishRealizer)

		pipeline := NewInterlinguaPipeline(services)

		pairs := pipeline.GetSupportedLanguagePairs()
		if len(pairs) == 0 {
			t.Error("Expected at least one supported language pair")
		}

		// Should have English to English
		foundEnglishPair := false
		for _, pair := range pairs {
			if pair.Source == "en" && pair.Target == "en" {
				foundEnglishPair = true
				break
			}
		}

		if !foundEnglishPair {
			t.Error("Expected to find English to English language pair")
		}

		t.Logf("Found %d supported language pairs", len(pairs))
		for _, pair := range pairs {
			t.Logf("  %s -> %s", pair.Source, pair.Target)
		}
	})
}

// TestInterlinguaPipelineValidation tests the translation validation.
func TestInterlinguaPipelineValidation(t *testing.T) {
	t.Run("Translation Validation", func(t *testing.T) {
		services := interlingua.NewInterlinguaServices()

		conceptResolver := interlingua.NewInMemoryConceptResolver()
		englishAnalyzer := interlingua.NewEnglishAnalyzer(conceptResolver)
		englishRealizer := interlingua.NewEnglishRealizer(42, conceptResolver)

		services.AddAnalyzer("en", englishAnalyzer)
		services.AddRealizer("en", englishRealizer)

		pipeline := NewInterlinguaPipeline(services)

		// Test valid language pair
		err := pipeline.ValidateTranslation("en", "en")
		if err != nil {
			t.Errorf("Expected no error for valid language pair, got: %v", err)
		}

		// Test invalid source language
		err = pipeline.ValidateTranslation("fr", "en")
		if err == nil {
			t.Error("Expected error for invalid source language")
		}
		if err.Error() != "source language 'fr' not supported (no analyzer available)" {
			t.Errorf("Expected specific error message, got: %v", err)
		}

		// Test invalid target language
		err = pipeline.ValidateTranslation("en", "fr")
		if err == nil {
			t.Error("Expected error for invalid target language")
		}
		if err.Error() != "target language 'fr' not supported (no realizer available)" {
			t.Errorf("Expected specific error message, got: %v", err)
		}

		t.Logf("Successfully tested translation validation")
	})
}

// TestInterlinguaPipelineIntegration tests the integration with the language system.
func TestInterlinguaPipelineIntegration(t *testing.T) {
	t.Run("Language Integration", func(t *testing.T) {
		// Create a language with interlingua services
		language, err := CreateRandomLanguage("test", "ancient", 42)
		if err != nil {
			t.Fatalf("Failed to create test language: %v", err)
		}

		// Configure interlingua services
		err = ConfigureInterlinguaServices(language)
		if err != nil {
			t.Fatalf("Failed to configure interlingua services: %v", err)
		}

		// Test that the language has interlingua services
		services := language.Interlingua()
		if services == nil {
			t.Fatal("Expected language to have interlingua services")
		}

		// Test that English analyzer and realizer are available
		if !services.HasAnalyzer("en") {
			t.Error("Expected English analyzer to be available")
		}
		if !services.HasRealizer("en") {
			t.Error("Expected English realizer to be available")
		}

		// Test that generated language analyzer and realizer are available
		if !services.HasAnalyzer("generated") {
			t.Error("Expected generated language analyzer to be available")
		}
		if !services.HasRealizer("generated") {
			t.Error("Expected generated language realizer to be available")
		}

		t.Logf("Successfully integrated interlingua services with language: %s", language.Name)
	})
}
