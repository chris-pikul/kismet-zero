package lang

import (
	"testing"

	"github.com/chris-pikul/kismet-zero/lang/interlingua"
)

func TestNewInterlinguaServices(t *testing.T) {
	services := NewInterlinguaServices()

	if services == nil {
		t.Fatal("Expected services to be created")
	}

	if services.Realizers == nil {
		t.Error("Expected Realizers map to be initialized")
	}

	if services.Analyzers == nil {
		t.Error("Expected Analyzers map to be initialized")
	}
}

func TestAddAndGetRealizer(t *testing.T) {
	services := NewInterlinguaServices()

	// Mock realizer (we don't need the actual interlingua package for this test)
	mockRealizer := &mockRealizer{langCode: "en"}

	services.AddRealizer("en", mockRealizer)

	// Test GetRealizer
	realizer, exists := services.GetRealizer("en")
	if !exists {
		t.Error("Expected realizer to exist")
	}

	if realizer == nil {
		t.Error("Expected realizer to be returned")
	}

	// Test HasRealizer
	if !services.HasRealizer("en") {
		t.Error("Expected HasRealizer to return true")
	}

	if services.HasRealizer("fr") {
		t.Error("Expected HasRealizer to return false for non-existent language")
	}
}

func TestAddAndGetAnalyzer(t *testing.T) {
	services := NewInterlinguaServices()

	// Mock analyzer
	mockAnalyzer := &mockAnalyzer{langCode: "en"}

	services.AddAnalyzer("en", mockAnalyzer)

	// Test GetAnalyzer
	analyzer, exists := services.GetAnalyzer("en")
	if !exists {
		t.Error("Expected analyzer to exist")
	}

	if analyzer == nil {
		t.Error("Expected analyzer to be returned")
	}

	// Test HasAnalyzer
	if !services.HasAnalyzer("en") {
		t.Error("Expected HasAnalyzer to return true")
	}

	if services.HasAnalyzer("fr") {
		t.Error("Expected HasAnalyzer to return false for non-existent language")
	}
}

func TestSupportedLanguages(t *testing.T) {
	services := NewInterlinguaServices()

	// Add realizers for multiple languages
	mockRealizer1 := &mockRealizer{langCode: "en"}
	mockRealizer2 := &mockRealizer{langCode: "fr"}

	services.AddRealizer("en", mockRealizer1)
	services.AddRealizer("fr", mockRealizer2)

	languages := services.SupportedLanguages()

	if len(languages) != 2 {
		t.Errorf("Expected 2 supported languages, got %d", len(languages))
	}

	// Check that both languages are present
	hasEN := false
	hasFR := false
	for _, lang := range languages {
		if lang == "en" {
			hasEN = true
		}
		if lang == "fr" {
			hasFR = true
		}
	}

	if !hasEN {
		t.Error("Expected English to be in supported languages")
	}

	if !hasFR {
		t.Error("Expected French to be in supported languages")
	}
}

// Mock implementations for testing
type mockRealizer struct {
	langCode string
}

func (m *mockRealizer) LanguageCode() string {
	return m.langCode
}

func (m *mockRealizer) Realize(doc interlingua.Document) ([]string, interlingua.Trace, error) {
	return []string{"mock", "output"}, interlingua.Trace{}, nil
}

type mockAnalyzer struct {
	langCode string
}

func (m *mockAnalyzer) LanguageCode() string {
	return m.langCode
}

func (m *mockAnalyzer) Analyze(tokens []string) (interlingua.Document, error) {
	return interlingua.Document{}, nil
}
