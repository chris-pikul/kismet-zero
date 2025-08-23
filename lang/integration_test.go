package lang

import (
	"testing"

	"github.com/chris-pikul/kismet-zero/lang/grammar"
	"github.com/chris-pikul/kismet-zero/lang/morphology"
)

func TestLanguageInterlinguaIntegration(t *testing.T) {
	// Create a new language
	langID := LanguageID{Family: "test", Language: "lang"}
	lang := NewLanguage(langID, "Test Language", LanguageTypeConstructed, 42)

	// Access interlingua services
	services := lang.Interlingua()
	if services == nil {
		t.Fatal("Expected interlingua services to be initialized")
	}

	// Verify services are empty initially
	if len(services.SupportedLanguages()) != 0 {
		t.Errorf("Expected 0 supported languages initially, got %d", len(services.SupportedLanguages()))
	}

	// Test that we can set services externally
	externalServices := NewInterlinguaServices()
	lang.SetInterlinguaServices(externalServices)

	// Verify the services were set
	if lang.interlinguaServices != externalServices {
		t.Error("Expected interlingua services to be set externally")
	}
}

func TestMorphologyGrammarIntegration(t *testing.T) {
	// This test demonstrates how the new fields would work together
	// in a real-world scenario where words have concept IDs and
	// morphemes have feature annotations

	// Create a word with concept IDs
	word := &morphology.Word{
		ID:         "walk-word",
		Meaning:    "to walk",
		Category:   morphology.WordCategoryVerb,
		ConceptIDs: []morphology.ConceptID{"walk-01", "move-01"},
	}

	// Create a morpheme with tense features
	morpheme := &morphology.Morpheme{
		ID:      "past-suffix",
		Type:    morphology.MorphemeTypeSuffix,
		Meaning: "past tense",
		Features: map[morphology.FeatureKey]morphology.FeatureVal{
			morphology.FeatureTense: "past",
		},
	}

	// Verify the new fields work
	if len(word.ConceptIDs) != 2 {
		t.Errorf("Expected 2 concept IDs, got %d", len(word.ConceptIDs))
	}

	if word.ConceptIDs[0] != "walk-01" {
		t.Errorf("Expected first concept ID 'walk-01', got %s", word.ConceptIDs[0])
	}

	if morpheme.Features[morphology.FeatureTense] != "past" {
		t.Errorf("Expected tense feature 'past', got %s", morpheme.Features[morphology.FeatureTense])
	}
}

func TestGrammarTemplateIntegration(t *testing.T) {
	// This test demonstrates how the new template fields work

	// Create a template with semantic roles
	template := &grammar.SentenceTemplate{
		Type:        grammar.SentenceTypeDeclarative,
		Pattern:     []string{"S", "V", "O"},
		Required:    []string{"S", "V"},
		Optional:    []string{"O", "ADV"},
		Weight:      10.0,
		Description: "Subject-Verb-Object declarative sentence",
		ID:          "decl.transitive.svo",
		Roles: map[string]grammar.Role{
			"S": "agent",
			"O": "patient",
		},
	}

	// Verify the new fields work
	if template.ID != "decl.transitive.svo" {
		t.Errorf("Expected template ID 'decl.transitive.svo', got %s", template.ID)
	}

	if len(template.Roles) != 2 {
		t.Errorf("Expected 2 roles, got %d", len(template.Roles))
	}

	if template.Roles["S"] != "agent" {
		t.Errorf("Expected subject role 'agent', got %s", template.Roles["S"])
	}

	if template.Roles["O"] != "patient" {
		t.Errorf("Expected object role 'patient', got %s", template.Roles["O"])
	}
}
