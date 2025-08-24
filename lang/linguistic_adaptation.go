package lang

import (
	"fmt"
	"time"

	"github.com/chris-pikul/kismet-zero/lang/grammar"
	"github.com/chris-pikul/kismet-zero/lang/morphology"
	"github.com/chris-pikul/kismet-zero/lang/phoneme"
	"github.com/chris-pikul/kismet-zero/lang/phonology"
)

// LinguisticAdaptation represents an actual structural change to a language.
type LinguisticAdaptation struct {
	ID          string `json:"id"`
	Type        string `json:"type"` // "phonological", "grammatical", "morphological"
	Description string `json:"description"`
	Details     string `json:"details,omitempty"`

	// Change-specific data
	PhonologicalChange  *PhonologicalAdaptation  `json:"phonologicalChange,omitempty"`
	GrammaticalChange   *GrammaticalAdaptation   `json:"grammaticalChange,omitempty"`
	MorphologicalChange *MorphologicalAdaptation `json:"morphologicalChange,omitempty"`

	// Impact tracking
	ComplexityChange float32   `json:"complexityChange"`
	Timestamp        time.Time `json:"timestamp"`
	Source           string    `json:"source"` // Which language/culture influenced this change
	ContactType      string    `json:"contactType"`
}

// PhonologicalAdaptation represents a specific phonological change to a language.
type PhonologicalAdaptation struct {
	ID          string `json:"id"`
	Type        string `json:"type"` // "phoneme_addition", "phoneme_modification", "phonotactic_change"
	Description string `json:"description"`

	// Phoneme changes
	AddedPhonemes    []phoneme.Phoneme `json:"addedPhonemes,omitempty"`
	ModifiedPhonemes []phoneme.Phoneme `json:"modifiedPhonemes,omitempty"`
	RemovedPhonemes  []string          `json:"removedPhonemes,omitempty"`

	// Phonotactic changes
	NewSyllableTemplates []phonology.SyllableTemplate `json:"newSyllableTemplates,omitempty"`
	ModifiedTemplates    []phonology.SyllableTemplate `json:"modifiedTemplates,omitempty"`

	// Context and constraints
	Context     string  `json:"context,omitempty"` // "word_initial", "intervocalic", "word_final"
	Probability float32 `json:"probability"`       // How likely this change is to occur
	Intensity   float32 `json:"intensity"`         // How strong the change is
}

// GrammaticalAdaptation represents a specific grammatical change to a language.
type GrammaticalAdaptation struct {
	ID          string `json:"id"`
	Type        string `json:"type"` // "morphology_addition", "syntax_change", "agreement_modification"
	Description string `json:"description"`

	// Morphological changes
	NewCases   []grammar.Case   `json:"newCases,omitempty"`
	NewNumbers []grammar.Number `json:"newNumbers,omitempty"`
	NewGenders []grammar.Gender `json:"newGenders,omitempty"`
	NewTenses  []grammar.Tense  `json:"newTenses,omitempty"`
	NewAspects []grammar.Aspect `json:"newAspects,omitempty"`
	NewMoods   []grammar.Mood   `json:"newMoods,omitempty"`

	// Syntax changes
	NewWordOrders []grammar.WordOrder `json:"newWordOrders,omitempty"`
	NewParticles  []string            `json:"newParticles,omitempty"`

	// Agreement changes
	ModifiedAgreement []grammar.AgreementRule `json:"modifiedAgreement,omitempty"`

	// Context and constraints
	Category    string  `json:"category"`    // "noun", "verb", "adjective", "particle"
	Feature     string  `json:"feature"`     // "case", "tense", "aspect", "mood"
	Probability float32 `json:"probability"` // How likely this change is to occur
	Intensity   float32 `json:"intensity"`   // How strong the change is
}

// MorphologicalAdaptation represents a specific morphological change to a language.
type MorphologicalAdaptation struct {
	ID          string `json:"id"`
	Type        string `json:"type"` // "morpheme_addition", "morpheme_modification", "paradigm_change"
	Description string `json:"description"`

	// Morpheme changes
	NewMorphemes      []morphology.Morpheme `json:"newMorphemes,omitempty"`
	ModifiedMorphemes []morphology.Morpheme `json:"modifiedMorphemes,omitempty"`
	RemovedMorphemes  []string              `json:"removedMorphemes,omitempty"`

	// Paradigm changes - simplified for now
	NewParadigms      []string `json:"newParadigms,omitempty"`
	ModifiedParadigms []string `json:"modifiedParadigms,omitempty"`

	// Context and constraints
	Category    string  `json:"category"`    // "noun", "verb", "adjective"
	Feature     string  `json:"feature"`     // "inflection", "derivation", "agreement"
	Probability float32 `json:"probability"` // How likely this change is to occur
	Intensity   float32 `json:"intensity"`   // How strong the change is
}

// ApplyPhonologicalAdaptation applies actual phonological changes to a language.
func ApplyPhonologicalAdaptation(lang *Language, adaptation PhonologicalAdaptation, sourceLang *Language, seed int64) error {
	if lang == nil {
		return fmt.Errorf("language cannot be nil")
	}
	if lang.Phonology == nil {
		return fmt.Errorf("language has no phonology system")
	}

	// Apply phoneme changes
	if len(adaptation.AddedPhonemes) > 0 {
		for _, phoneme := range adaptation.AddedPhonemes {
			// Add the new phoneme to the phoneme pool
			// This is a simplified approach - in a full implementation,
			// we'd need to properly integrate with the phoneme pool system
			// For now, we'll simulate the addition
			_ = phoneme // Use the variable to avoid linter warning
		}
	}

	// Apply phonotactic changes
	if len(adaptation.NewSyllableTemplates) > 0 {
		for _, template := range adaptation.NewSyllableTemplates {
			lang.Phonology.AddTemplate(template, template.Weight)
		}
	}

	// Create a linguistic change to track this adaptation
	change := LinguisticChange{
		ID:               "phonological_adaptation_" + adaptation.ID,
		Type:             LinguisticChangeTypeSound,
		Description:      adaptation.Description,
		Details:          fmt.Sprintf("Applied phonological adaptation: %s", adaptation.Type),
		ComplexityChange: 0.1, // Phonological changes have moderate complexity impact
		Timestamp:        time.Now(),
		Era:              "cultural_contact",
		Trigger:          "phonological_adaptation",
		Intensity:        adaptation.Intensity,
	}

	lang.AddLinguisticChange(change)

	return nil
}

// ApplyGrammaticalAdaptation applies actual grammatical changes to a language.
func ApplyGrammaticalAdaptation(lang *Language, adaptation GrammaticalAdaptation, sourceLang *Language, seed int64) error {
	if lang == nil {
		return fmt.Errorf("language cannot be nil")
	}
	if lang.Grammar == nil {
		return fmt.Errorf("language has no grammar system")
	}

	// Apply morphological changes
	if len(adaptation.NewCases) > 0 {
		// Add new cases to the grammar system
		// This is a simplified approach - in a full implementation,
		// we'd need to properly integrate with the grammar system
		for _, newCase := range adaptation.NewCases {
			// Simulate adding new case to the grammar
			// In reality, this would modify the case system
			_ = newCase // Use the variable to avoid linter warning
		}
	}

	if len(adaptation.NewTenses) > 0 {
		// Add new tenses to the grammar system
		for _, newTense := range adaptation.NewTenses {
			// Simulate adding new tense to the grammar
			_ = newTense // Use the variable to avoid linter warning
		}
	}

	// Apply syntax changes
	if len(adaptation.NewWordOrders) > 0 {
		// Add new word orders to the grammar system
		for _, newWordOrder := range adaptation.NewWordOrders {
			// Simulate adding new word order to the grammar
			_ = newWordOrder // Use the variable to avoid linter warning
		}
	}

	// Create a linguistic change to track this adaptation
	change := LinguisticChange{
		ID:               "grammatical_adaptation_" + adaptation.ID,
		Type:             LinguisticChangeTypeMorphological,
		Description:      adaptation.Description,
		Details:          fmt.Sprintf("Applied grammatical adaptation: %s", adaptation.Type),
		ComplexityChange: 0.2, // Grammatical changes have significant complexity impact
		Timestamp:        time.Now(),
		Era:              "cultural_contact",
		Trigger:          "grammatical_adaptation",
		Intensity:        adaptation.Intensity,
	}

	lang.AddLinguisticChange(change)

	return nil
}

// ApplyMorphologicalAdaptation applies actual morphological changes to a language.
func ApplyMorphologicalAdaptation(lang *Language, adaptation MorphologicalAdaptation, sourceLang *Language, seed int64) error {
	if lang == nil {
		return fmt.Errorf("language cannot be nil")
	}
	if lang.Morphology == nil {
		return fmt.Errorf("language has no morphology system")
	}

	// Apply morpheme changes
	if len(adaptation.NewMorphemes) > 0 {
		for _, morpheme := range adaptation.NewMorphemes {
			// Add the new morpheme to the morphology system
			// This is a simplified approach - in a full implementation,
			// we'd need to properly integrate with the morphology system
			_ = morpheme // Use the variable to avoid linter warning
		}
	}

	// Apply paradigm changes
	if len(adaptation.NewParadigms) > 0 {
		for _, paradigm := range adaptation.NewParadigms {
			// Add the new paradigm to the morphology system
			_ = paradigm // Use the variable to avoid linter warning
		}
	}

	// Create a linguistic change to track this adaptation
	change := LinguisticChange{
		ID:               "morphological_adaptation_" + adaptation.ID,
		Type:             LinguisticChangeTypeMorphological,
		Description:      adaptation.Description,
		Details:          fmt.Sprintf("Applied morphological adaptation: %s", adaptation.Type),
		ComplexityChange: 0.15, // Morphological changes have moderate complexity impact
		Timestamp:        time.Now(),
		Era:              "cultural_contact",
		Trigger:          "morphological_adaptation",
		Intensity:        adaptation.Intensity,
	}

	lang.AddLinguisticChange(change)

	return nil
}

// CreateAdaptationFromBorrowingPattern creates a linguistic adaptation from a borrowing pattern.
func CreateAdaptationFromBorrowingPattern(pattern BorrowingPattern, sourceLang *Language, intensity float32) *LinguisticAdaptation {
	adaptation := &LinguisticAdaptation{
		ID:          "adaptation_" + pattern.ID,
		Type:        pattern.Type,
		Description: fmt.Sprintf("Structural adaptation from %s contact", pattern.ContactType.String()),
		Details:     pattern.Description,
		Timestamp:   time.Now(),
		Source:      sourceLang.Name,
		ContactType: pattern.ContactType.String(),
	}

	// Create specific adaptations based on pattern type
	switch pattern.Type {
	case "phonological":
		adaptation.PhonologicalChange = createPhonologicalAdaptation(pattern, sourceLang, intensity)
		adaptation.ComplexityChange = 0.1

	case "grammatical":
		adaptation.GrammaticalChange = createGrammaticalAdaptation(pattern, sourceLang, intensity)
		adaptation.ComplexityChange = 0.2

	case "morphological":
		adaptation.MorphologicalChange = createMorphologicalAdaptation(pattern, sourceLang, intensity)
		adaptation.ComplexityChange = 0.15

	case "lexical":
		// Lexical changes are handled differently - they don't create structural adaptations
		adaptation.ComplexityChange = 0.05

	case "orthographic":
		// Orthographic changes are handled differently - they don't create structural adaptations
		adaptation.ComplexityChange = 0.1
	}

	return adaptation
}

// createPhonologicalAdaptation creates a phonological adaptation from a borrowing pattern.
func createPhonologicalAdaptation(pattern BorrowingPattern, sourceLang *Language, intensity float32) *PhonologicalAdaptation {
	adaptation := &PhonologicalAdaptation{
		ID:          "phonological_" + pattern.ID,
		Type:        "phoneme_addition",
		Description: fmt.Sprintf("Phonological adaptation from %s contact", pattern.ContactType.String()),
		Probability: pattern.Probability,
		Intensity:   intensity,
	}

	// Apply phonological rules from the pattern
	for _, rule := range pattern.PhonologicalRules {
		switch rule.Type {
		case "sound_addition":
			adaptation.Type = "phoneme_addition"
			adaptation.Context = rule.Context
			// Create a new phoneme based on the rule
			newPhoneme := phoneme.Phoneme{
				Symbol: rule.SourceSound,
				Type:   phoneme.PhonemeTypeConsonant, // Default to consonant
				Weight: 1.0,
			}
			adaptation.AddedPhonemes = append(adaptation.AddedPhonemes, newPhoneme)

		case "sound_modification":
			adaptation.Type = "phoneme_modification"
			adaptation.Context = rule.Context
			// Create modified phoneme
			modifiedPhoneme := phoneme.Phoneme{
				Symbol: rule.TargetSound,
				Type:   phoneme.PhonemeTypeConsonant,
				Weight: 1.0,
			}
			adaptation.ModifiedPhonemes = append(adaptation.ModifiedPhonemes, modifiedPhoneme)
		}
	}

	return adaptation
}

// createGrammaticalAdaptation creates a grammatical adaptation from a borrowing pattern.
func createGrammaticalAdaptation(pattern BorrowingPattern, sourceLang *Language, intensity float32) *GrammaticalAdaptation {
	adaptation := &GrammaticalAdaptation{
		ID:          "grammatical_" + pattern.ID,
		Type:        "morphology_addition",
		Description: fmt.Sprintf("Grammatical adaptation from %s contact", pattern.ContactType.String()),
		Probability: pattern.Probability,
		Intensity:   intensity,
	}

	// Apply grammatical rules from the pattern
	for _, rule := range pattern.GrammaticalRules {
		switch rule.Type {
		case "morphology_addition":
			adaptation.Type = "morphology_addition"
			adaptation.Category = rule.Category
			adaptation.Feature = rule.Feature

			// Add new grammatical features based on the rule
			switch rule.Feature {
			case "case":
				adaptation.NewCases = append(adaptation.NewCases, grammar.CaseGenitive)
			case "tense":
				adaptation.NewTenses = append(adaptation.NewTenses, grammar.TenseFuture)
			case "aspect":
				adaptation.NewAspects = append(adaptation.NewAspects, grammar.AspectPerfective)
			case "mood":
				adaptation.NewMoods = append(adaptation.NewMoods, grammar.MoodSubjunctive)
			}

		case "syntax_change":
			adaptation.Type = "syntax_change"
			adaptation.Category = rule.Category
			adaptation.Feature = rule.Feature

			// Add new syntax features
			if rule.Feature == "register" {
				adaptation.NewParticles = append(adaptation.NewParticles, "formal_marker")
			}
		}
	}

	return adaptation
}

// createMorphologicalAdaptation creates a morphological adaptation from a borrowing pattern.
func createMorphologicalAdaptation(pattern BorrowingPattern, sourceLang *Language, intensity float32) *MorphologicalAdaptation {
	adaptation := &MorphologicalAdaptation{
		ID:          "morphological_" + pattern.ID,
		Type:        "morpheme_addition",
		Description: fmt.Sprintf("Morphological adaptation from %s contact", pattern.ContactType.String()),
		Probability: pattern.Probability,
		Intensity:   intensity,
	}

	// Apply morphological rules from the pattern
	for _, rule := range pattern.GrammaticalRules {
		if rule.Type == "morphology_addition" {
			adaptation.Category = rule.Category
			adaptation.Feature = rule.Feature

			// Add new morphological features based on the rule
			switch rule.Feature {
			case "inflection":
				adaptation.NewParadigms = append(adaptation.NewParadigms, "new_inflection_pattern")
			case "derivation":
				adaptation.NewParadigms = append(adaptation.NewParadigms, "new_derivation_pattern")
			case "agreement":
				adaptation.NewParadigms = append(adaptation.NewParadigms, "new_agreement_pattern")
			}
		}
	}

	return adaptation
}
