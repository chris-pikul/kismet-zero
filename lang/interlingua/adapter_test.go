package interlingua

import (
	"testing"
)

// MockRealizer is a mock implementation of Realizer for testing.
type MockRealizer struct {
	language string
}

// LanguageCode implements the Realizer interface.
func (mr *MockRealizer) LanguageCode() string {
	return mr.language
}

// Realize implements the Realizer interface.
func (mr *MockRealizer) Realize(doc Document) ([]string, Trace, error) {
	return []string{"mock", "realization", "in", mr.language}, Trace{}, nil
}

// MockAnalyzer is a mock implementation of Analyzer for testing.
type MockAnalyzer struct {
	language string
}

// LanguageCode implements the Analyzer interface.
func (ma *MockAnalyzer) LanguageCode() string {
	return ma.language
}

// Analyze implements the Analyzer interface.
func (ma *MockAnalyzer) Analyze(tokens []string) (Document, error) {
	return Document{
		SchemaVersion: "1.0",
		LanguageID:    ma.language,
		Entities:      []Entity{},
		Events:        []Event{},
	}, nil
}

func TestNewInterlinguaServices(t *testing.T) {
	services := NewInterlinguaServices()

	if services == nil {
		t.Error("NewInterlinguaServices should not return nil")
	}

	if services.Realizers == nil {
		t.Error("Realizers map should be initialized")
	}

	if services.Analyzers == nil {
		t.Error("Analyzers map should be initialized")
	}
}

func TestAddRealizer(t *testing.T) {
	services := NewInterlinguaServices()
	mockRealizer := &MockRealizer{language: "en"}

	services.AddRealizer("en", mockRealizer)

	if !services.HasRealizer("en") {
		t.Error("Realizer should be added and detectable")
	}

	retrieved, exists := services.GetRealizer("en")
	if !exists {
		t.Error("Realizer should exist after being added")
	}

	if retrieved != mockRealizer {
		t.Error("Retrieved realizer should match the added one")
	}
}

func TestAddAnalyzer(t *testing.T) {
	services := NewInterlinguaServices()
	mockAnalyzer := &MockAnalyzer{language: "en"}

	services.AddAnalyzer("en", mockAnalyzer)

	if !services.HasAnalyzer("en") {
		t.Error("Analyzer should be added and detectable")
	}

	retrieved, exists := services.GetAnalyzer("en")
	if !exists {
		t.Error("Analyzer should exist after being added")
	}

	if retrieved != mockAnalyzer {
		t.Error("Retrieved analyzer should match the added one")
	}
}

func TestGetRealizerNonExistent(t *testing.T) {
	services := NewInterlinguaServices()

	_, exists := services.GetRealizer("nonexistent")
	if exists {
		t.Error("Non-existent realizer should not be found")
	}
}

func TestGetAnalyzerNonExistent(t *testing.T) {
	services := NewInterlinguaServices()

	_, exists := services.GetAnalyzer("nonexistent")
	if exists {
		t.Error("Non-existent analyzer should not be found")
	}
}

func TestSupportedLanguages(t *testing.T) {
	services := NewInterlinguaServices()

	// Add realizers for multiple languages
	services.AddRealizer("en", &MockRealizer{language: "en"})
	services.AddRealizer("fr", &MockRealizer{language: "fr"})
	services.AddRealizer("de", &MockRealizer{language: "de"})

	languages := services.SupportedLanguages()

	if len(languages) != 3 {
		t.Errorf("Expected 3 supported languages, got %d", len(languages))
	}

	// Check that all expected languages are present
	expected := map[string]bool{"en": true, "fr": true, "de": true}
	for _, lang := range languages {
		if !expected[lang] {
			t.Errorf("Unexpected language: %s", lang)
		}
	}
}
