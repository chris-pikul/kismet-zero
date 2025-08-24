package lang

import (
	"fmt"
	"math/rand"
	"time"
)

// ContactType represents the type of language contact.
type ContactType byte

const (
	ContactTypeUnknown     ContactType = iota
	ContactTypeTrade                   // Trade relationships
	ContactTypeConquest                // Military conquest
	ContactTypeMigration               // Population migration
	ContactTypeCultural                // Cultural exchange
	ContactTypeReligious               // Religious influence
	ContactTypeEducational             // Educational exchange
)

var contactTypeEnum = []string{
	"unknown",
	"trade",
	"conquest",
	"migration",
	"cultural",
	"religious",
	"educational",
}

// String returns the string representation of the ContactType.
func (ct ContactType) String() string {
	if ct > ContactTypeEducational {
		return contactTypeEnum[0]
	}
	return contactTypeEnum[ct]
}

// BorrowingPattern represents a pattern of linguistic borrowing.
type BorrowingPattern struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"` // "lexical", "phonological", "grammatical", "orthographic"
	Description string      `json:"description"`
	Probability float32     `json:"probability"` // 0.0 to 1.0
	Intensity   float32     `json:"intensity"`   // 0.0 to 1.0
	ContactType ContactType `json:"contactType"`

	// Enhanced borrowing mechanisms
	PhonologicalRules []PhonologicalRule `json:"phonologicalRules,omitempty"`
	GrammaticalRules  []GrammaticalRule  `json:"grammaticalRules,omitempty"`
	LexicalRules      []LexicalRule      `json:"lexicalRules,omitempty"`
	OrthographicRules []OrthographicRule `json:"orthographicRules,omitempty"`
}

// PhonologicalRule represents a rule for phonological borrowing and adaptation.
type PhonologicalRule struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Type        string  `json:"type"` // "sound_addition", "sound_modification", "phonotactic_change"
	SourceSound string  `json:"sourceSound,omitempty"`
	TargetSound string  `json:"targetSound,omitempty"`
	Context     string  `json:"context,omitempty"` // "word_initial", "word_final", "intervocalic", etc.
	Probability float32 `json:"probability"`
}

// GrammaticalRule represents a rule for grammatical borrowing and adaptation.
type GrammaticalRule struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Type        string  `json:"type"`     // "morphology_addition", "syntax_change", "agreement_modification"
	Category    string  `json:"category"` // "noun", "verb", "adjective", "particle"
	Feature     string  `json:"feature"`  // "case", "tense", "aspect", "mood"
	Probability float32 `json:"probability"`
}

// LexicalRule represents a rule for lexical borrowing and adaptation.
type LexicalRule struct {
	ID             string  `json:"id"`
	Description    string  `json:"description"`
	Type           string  `json:"type"`           // "word_borrowing", "semantic_shift", "calque_formation"
	SemanticField  string  `json:"semanticField"`  // "trade", "religion", "technology", "culture"
	AdaptationType string  `json:"adaptationType"` // "phonological", "morphological", "semantic"
	Probability    float32 `json:"probability"`
}

// OrthographicRule represents a rule for orthographic borrowing and adaptation.
type OrthographicRule struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Type        string  `json:"type"`                 // "script_adoption", "spelling_reform", "diacritic_addition"
	ScriptType  string  `json:"scriptType,omitempty"` // "alphabet", "syllabary", "logographic"
	ReformType  string  `json:"reformType,omitempty"` // "simplification", "standardization", "foreign_adoption"
	Probability float32 `json:"probability"`
}

// CommonBorrowingPatterns provides realistic borrowing patterns for different contact types.
var CommonBorrowingPatterns = []BorrowingPattern{
	// Trade contact patterns
	{
		ID:          "trade_lexical",
		Type:        "lexical",
		Description: "Borrowing of trade-related vocabulary",
		Probability: 0.8,
		Intensity:   0.6,
		ContactType: ContactTypeTrade,
		LexicalRules: []LexicalRule{
			{
				ID:             "trade_commerce_terms",
				Description:    "Borrowing of commerce and trade terminology",
				Type:           "word_borrowing",
				SemanticField:  "trade",
				AdaptationType: "phonological",
				Probability:    0.9,
			},
			{
				ID:             "trade_currency_terms",
				Description:    "Adoption of foreign currency and measurement terms",
				Type:           "word_borrowing",
				SemanticField:  "trade",
				AdaptationType: "phonological",
				Probability:    0.8,
			},
		},
	},
	{
		ID:          "trade_phonological",
		Type:        "phonological",
		Description: "Adaptation of foreign sounds for trade terms",
		Probability: 0.4,
		Intensity:   0.3,
		ContactType: ContactTypeTrade,
		PhonologicalRules: []PhonologicalRule{
			{
				ID:          "trade_sound_addition",
				Description: "Addition of foreign sounds for trade vocabulary",
				Type:        "sound_addition",
				SourceSound: "foreign_fricative",
				TargetSound: "native_approximant",
				Context:     "word_initial",
				Probability: 0.6,
			},
		},
	},

	// Conquest contact patterns
	{
		ID:          "conquest_grammatical",
		Type:        "grammatical",
		Description: "Adoption of administrative grammar structures",
		Probability: 0.7,
		Intensity:   0.8,
		ContactType: ContactTypeConquest,
		GrammaticalRules: []GrammaticalRule{
			{
				ID:          "conquest_administrative_case",
				Description: "Adoption of administrative case marking",
				Type:        "morphology_addition",
				Category:    "noun",
				Feature:     "case",
				Probability: 0.8,
			},
			{
				ID:          "conquest_formal_register",
				Description: "Development of formal register markers",
				Type:        "syntax_change",
				Category:    "particle",
				Feature:     "register",
				Probability: 0.7,
			},
		},
	},
	{
		ID:          "conquest_orthographic",
		Type:        "orthographic",
		Description: "Adoption of writing system conventions",
		Probability: 0.6,
		Intensity:   0.7,
		ContactType: ContactTypeConquest,
		OrthographicRules: []OrthographicRule{
			{
				ID:          "conquest_script_adoption",
				Description: "Adoption of administrative writing conventions",
				Type:        "script_adoption",
				ScriptType:  "alphabet",
				ReformType:  "foreign_adoption",
				Probability: 0.7,
			},
		},
	},

	// Migration contact patterns
	{
		ID:          "migration_phonological",
		Type:        "phonological",
		Description: "Phonological accommodation to migrant speech",
		Probability: 0.6,
		Intensity:   0.5,
		ContactType: ContactTypeMigration,
		PhonologicalRules: []PhonologicalRule{
			{
				ID:          "migration_accommodation",
				Description: "Phonological accommodation to migrant speech patterns",
				Type:        "sound_modification",
				SourceSound: "migrant_vowel",
				TargetSound: "native_vowel",
				Context:     "intervocalic",
				Probability: 0.6,
			},
		},
	},
	{
		ID:          "migration_lexical",
		Type:        "lexical",
		Description: "Borrowing of cultural and domestic terms",
		Probability: 0.7,
		Intensity:   0.6,
		ContactType: ContactTypeMigration,
		LexicalRules: []LexicalRule{
			{
				ID:             "migration_domestic_terms",
				Description:    "Borrowing of domestic and cultural terminology",
				Type:           "word_borrowing",
				SemanticField:  "culture",
				AdaptationType: "morphological",
				Probability:    0.8,
			},
		},
	},

	// Cultural contact patterns
	{
		ID:          "cultural_lexical",
		Type:        "lexical",
		Description: "Borrowing of cultural and artistic terms",
		Probability: 0.8,
		Intensity:   0.5,
		ContactType: ContactTypeCultural,
		LexicalRules: []LexicalRule{
			{
				ID:             "cultural_artistic_terms",
				Description:    "Borrowing of artistic and cultural terminology",
				Type:           "word_borrowing",
				SemanticField:  "culture",
				AdaptationType: "semantic",
				Probability:    0.9,
			},
		},
	},
	{
		ID:          "cultural_grammatical",
		Type:        "grammatical",
		Description: "Adoption of cultural expression patterns",
		Probability: 0.5,
		Intensity:   0.4,
		ContactType: ContactTypeCultural,
		GrammaticalRules: []GrammaticalRule{
			{
				ID:          "cultural_expression_patterns",
				Description: "Adoption of cultural expression and politeness markers",
				Type:        "syntax_change",
				Category:    "particle",
				Feature:     "politeness",
				Probability: 0.6,
			},
		},
	},

	// Religious contact patterns
	{
		ID:          "religious_lexical",
		Type:        "lexical",
		Description: "Borrowing of religious and spiritual terms",
		Probability: 0.7,
		Intensity:   0.6,
		ContactType: ContactTypeReligious,
		LexicalRules: []LexicalRule{
			{
				ID:             "religious_spiritual_terms",
				Description:    "Borrowing of spiritual and religious terminology",
				Type:           "word_borrowing",
				SemanticField:  "religion",
				AdaptationType: "phonological",
				Probability:    0.8,
			},
		},
	},
	{
		ID:          "religious_grammatical",
		Type:        "grammatical",
		Description: "Adoption of religious speech patterns",
		Probability: 0.6,
		Intensity:   0.5,
		ContactType: ContactTypeReligious,
		GrammaticalRules: []GrammaticalRule{
			{
				ID:          "religious_honorifics",
				Description: "Adoption of religious honorific and respect markers",
				Type:        "morphology_addition",
				Category:    "noun",
				Feature:     "honorific",
				Probability: 0.7,
			},
		},
	},

	// Educational contact patterns
	{
		ID:          "educational_lexical",
		Type:        "lexical",
		Description: "Borrowing of academic and technical terms",
		Probability: 0.8,
		Intensity:   0.4,
		ContactType: ContactTypeEducational,
		LexicalRules: []LexicalRule{
			{
				ID:             "educational_academic_terms",
				Description:    "Borrowing of academic and technical terminology",
				Type:           "word_borrowing",
				SemanticField:  "technology",
				AdaptationType: "morphological",
				Probability:    0.9,
			},
		},
	},
	{
		ID:          "educational_orthographic",
		Type:        "orthographic",
		Description: "Adoption of standardized writing conventions",
		Probability: 0.7,
		Intensity:   0.5,
		ContactType: ContactTypeEducational,
		OrthographicRules: []OrthographicRule{
			{
				ID:          "educational_standardization",
				Description: "Adoption of standardized educational writing conventions",
				Type:        "spelling_reform",
				ScriptType:  "alphabet",
				ReformType:  "standardization",
				Probability: 0.8,
			},
		},
	},
}

// CulturalInfluenceResult represents the result of cultural influence on a language.
type CulturalInfluenceResult struct {
	LanguageID          string             `json:"languageId"`
	ContactType         ContactType        `json:"contactType"`
	Intensity           float32            `json:"intensity"`
	Duration            string             `json:"duration"`
	BorrowingPatterns   []BorrowingPattern `json:"borrowingPatterns"`
	LexicalChanges      []string           `json:"lexicalChanges,omitempty"`
	PhonologicalChanges []string           `json:"phonologicalChanges,omitempty"`
	GrammaticalChanges  []string           `json:"grammaticalChanges,omitempty"`
	OrthographicChanges []string           `json:"orthographicChanges,omitempty"`
	ComplexityChange    float32            `json:"complexityChange"`
	Timestamp           time.Time          `json:"timestamp"`
}

// ApplyCulturalInfluence applies cultural influence between two languages, creating actual structural changes.
func ApplyCulturalInfluence(targetLang, sourceLang *Language, contactType ContactType, intensity float32, duration time.Duration, seed int64) (*CulturalInfluenceResult, error) {
	if targetLang == nil || sourceLang == nil {
		return nil, fmt.Errorf("both languages must be provided")
	}

	// Initialize contact intensity calculator
	calculator := NewContactIntensityCalculator()

	// Get or create contact history
	contactHistory := calculator.GetContactHistory(targetLang, sourceLang, contactType)

	// Calculate geographic proximity
	geographicProximity := calculator.CalculateGeographicProximity(targetLang, sourceLang)

	// Calculate effective intensity using the intensity model
	effectiveIntensity := calculator.CalculateEffectiveIntensity(
		contactType,
		intensity,
		duration,
		contactHistory,
		geographicProximity,
	)

	// Get applicable borrowing patterns for this contact type
	patterns := getApplicableBorrowingPatterns(contactType)
	if len(patterns) == 0 {
		// Return minimal result if no patterns available
		return &CulturalInfluenceResult{
			LanguageID:          targetLang.ID.String(),
			ContactType:         contactType,
			Intensity:           effectiveIntensity,
			Duration:            duration.String(),
			ComplexityChange:    0.0,
			PhonologicalChanges: []string{},
			GrammaticalChanges:  []string{},
			LexicalChanges:      []string{},
			OrthographicChanges: []string{},
			Timestamp:           time.Now(),
		}, nil
	}

	result := &CulturalInfluenceResult{
		LanguageID:          targetLang.ID.String(),
		ContactType:         contactType,
		Intensity:           effectiveIntensity,
		Duration:            duration.String(),
		ComplexityChange:    0.0,
		PhonologicalChanges: []string{},
		GrammaticalChanges:  []string{},
		LexicalChanges:      []string{},
		OrthographicChanges: []string{},
		Timestamp:           time.Now(),
	}

	// Determine if structural adaptation should occur based on effective intensity
	shouldAdapt := calculator.ShouldApplyAdaptation(contactType, effectiveIntensity)

	// Collect all changes for history tracking
	var allChanges []string

	// Apply borrowing patterns based on probability and intensity
	for _, pattern := range patterns {
		// Adjust pattern probability based on effective intensity
		adjustedProbability := pattern.Probability * effectiveIntensity

		if rand.Float32() <= adjustedProbability {
			// Create and apply linguistic adaptation based on the pattern
			adaptation := CreateAdaptationFromBorrowingPattern(pattern, sourceLang, effectiveIntensity)

			// Apply the adaptation based on its type
			switch adaptation.Type {
			case "phonological":
				if adaptation.PhonologicalChange != nil && shouldAdapt {
					err := ApplyPhonologicalAdaptation(targetLang, *adaptation.PhonologicalChange, sourceLang, seed)
					if err == nil {
						result.PhonologicalChanges = append(result.PhonologicalChanges, adaptation.Description)
						result.ComplexityChange += adaptation.ComplexityChange
						allChanges = append(allChanges, adaptation.Description)
					}
				}

			case "grammatical":
				if adaptation.GrammaticalChange != nil && shouldAdapt {
					err := ApplyGrammaticalAdaptation(targetLang, *adaptation.GrammaticalChange, sourceLang, seed)
					if err == nil {
						result.GrammaticalChanges = append(result.GrammaticalChanges, adaptation.Description)
						result.ComplexityChange += adaptation.ComplexityChange
						allChanges = append(allChanges, adaptation.Description)
					}
				}

			case "morphological":
				if adaptation.MorphologicalChange != nil && shouldAdapt {
					err := ApplyMorphologicalAdaptation(targetLang, *adaptation.MorphologicalChange, sourceLang, seed)
					if err == nil {
						result.GrammaticalChanges = append(result.GrammaticalChanges, adaptation.Description)
						result.ComplexityChange += adaptation.ComplexityChange
						allChanges = append(allChanges, adaptation.Description)
					}
				}

			case "lexical":
				lexicalChanges := applyLexicalBorrowing(targetLang, sourceLang, pattern, effectiveIntensity, seed)
				result.LexicalChanges = append(result.LexicalChanges, lexicalChanges...)
				result.ComplexityChange += 0.1 * float32(len(lexicalChanges)) // Lexical borrowing slightly increases complexity
				allChanges = append(allChanges, lexicalChanges...)

			case "orthographic":
				orthographicChanges := applyOrthographicBorrowing(targetLang, sourceLang, pattern, effectiveIntensity, seed)
				result.OrthographicChanges = append(result.OrthographicChanges, orthographicChanges...)
				result.ComplexityChange += 0.1 * float32(len(orthographicChanges)) // Orthographic borrowing minimal complexity impact
				allChanges = append(allChanges, orthographicChanges...)
			}
		}
	}

	// Update contact history with this contact event
	calculator.UpdateContactHistory(
		contactHistory,
		contactType,
		duration,
		effectiveIntensity,
		allChanges,
		result.ComplexityChange,
	)

	// Normalize complexity change
	if result.ComplexityChange > 1.0 {
		result.ComplexityChange = 1.0
	}

	return result, nil
}

// getApplicableBorrowingPatterns returns patterns that apply to the given contact type and intensity.
func getApplicableBorrowingPatterns(contactType ContactType) []BorrowingPattern {
	var patterns []BorrowingPattern

	for _, pattern := range CommonBorrowingPatterns {
		if pattern.ContactType == contactType {
			patterns = append(patterns, pattern)
		}
	}

	return patterns
}

// applyInfluenceToLanguage applies cultural influence to a single language.
func applyInfluenceToLanguage(targetLang, sourceLang *Language, patterns []BorrowingPattern, intensity float32, duration string, seed int64) *CulturalInfluenceResult {
	// Check if we have any patterns to apply
	if len(patterns) == 0 {
		// Return a minimal result with no changes
		return &CulturalInfluenceResult{
			LanguageID:        targetLang.ID.String(),
			ContactType:       ContactTypeUnknown,
			Intensity:         intensity,
			Duration:          duration,
			BorrowingPatterns: patterns,
			Timestamp:         time.Now(),
		}
	}

	result := &CulturalInfluenceResult{
		LanguageID:        targetLang.ID.String(),
		ContactType:       patterns[0].ContactType, // Use first pattern's contact type
		Intensity:         intensity,
		Duration:          duration,
		BorrowingPatterns: patterns,
		Timestamp:         time.Now(),
	}

	// Apply borrowing patterns based on probability
	for _, pattern := range patterns {
		// Use seed-based random generation
		rand.Seed(seed)
		if rand.Float32() < pattern.Probability {
			// Apply the specific type of borrowing
			switch pattern.Type {
			case "lexical":
				lexicalChanges := applyLexicalBorrowing(targetLang, sourceLang, pattern, intensity, seed)
				result.LexicalChanges = append(result.LexicalChanges, lexicalChanges...)
				result.ComplexityChange += 0.1 * float32(len(lexicalChanges)) // Lexical borrowing slightly increases complexity

			case "phonological":
				phonologicalChanges := applyPhonologicalBorrowing(targetLang, sourceLang, pattern, intensity, seed)
				result.PhonologicalChanges = append(result.PhonologicalChanges, phonologicalChanges...)
				result.ComplexityChange += 0.05 * float32(len(phonologicalChanges)) // Phonological borrowing minimal complexity impact

			case "grammatical":
				grammaticalChanges := applyGrammaticalBorrowing(targetLang, sourceLang, pattern, intensity, seed)
				result.GrammaticalChanges = append(result.GrammaticalChanges, grammaticalChanges...)
				result.ComplexityChange += 0.2 * float32(len(grammaticalChanges)) // Grammatical borrowing significant complexity impact

			case "orthographic":
				orthographicChanges := applyOrthographicBorrowing(targetLang, sourceLang, pattern, intensity, seed)
				result.OrthographicChanges = append(result.OrthographicChanges, orthographicChanges...)
				result.ComplexityChange += 0.15 * float32(len(orthographicChanges)) // Orthographic borrowing moderate complexity impact
			}
		}
	}

	// Normalize complexity change
	if result.ComplexityChange > 1.0 {
		result.ComplexityChange = 1.0
	}

	return result
}

// applyLexicalBorrowing applies lexical borrowing rules to the target language.
func applyLexicalBorrowing(targetLang, sourceLang *Language, pattern BorrowingPattern, intensity float32, seed int64) []string {
	var changes []string

	for _, rule := range pattern.LexicalRules {
		rand.Seed(seed + int64(len(changes)))
		if rand.Float32() < rule.Probability {
			// Apply the lexical rule
			change := applyLexicalRule(targetLang, sourceLang, rule, intensity)
			if change != "" {
				changes = append(changes, change)
			}
		}
	}

	return changes
}

// applyPhonologicalBorrowing applies phonological borrowing rules to the target language.
func applyPhonologicalBorrowing(targetLang, sourceLang *Language, pattern BorrowingPattern, intensity float32, seed int64) []string {
	var changes []string

	for _, rule := range pattern.PhonologicalRules {
		rand.Seed(seed + int64(len(changes)))
		if rand.Float32() < rule.Probability {
			// Apply the phonological rule
			change := applyPhonologicalRule(targetLang, sourceLang, rule, intensity)
			if change != "" {
				changes = append(changes, change)
			}
		}
	}

	return changes
}

// applyGrammaticalBorrowing applies grammatical borrowing rules to the target language.
func applyGrammaticalBorrowing(targetLang, sourceLang *Language, pattern BorrowingPattern, intensity float32, seed int64) []string {
	var changes []string

	for _, rule := range pattern.GrammaticalRules {
		rand.Seed(seed + int64(len(changes)))
		if rand.Float32() < rule.Probability {
			// Apply the grammatical rule
			change := applyGrammaticalRule(targetLang, sourceLang, rule, intensity)
			if change != "" {
				changes = append(changes, change)
			}
		}
	}

	return changes
}

// applyOrthographicBorrowing applies orthographic borrowing rules to the target language.
func applyOrthographicBorrowing(targetLang, sourceLang *Language, pattern BorrowingPattern, intensity float32, seed int64) []string {
	var changes []string

	for _, rule := range pattern.OrthographicRules {
		rand.Seed(seed + int64(len(changes)))
		if rand.Float32() < rule.Probability {
			// Apply the orthographic rule
			change := applyOrthographicRule(targetLang, sourceLang, rule, intensity)
			if change != "" {
				changes = append(changes, change)
			}
		}
	}

	return changes
}

// applyLexicalRule applies a specific lexical rule to the target language.
func applyLexicalRule(targetLang, sourceLang *Language, rule LexicalRule, intensity float32) string {
	// Simulate lexical borrowing based on the rule type
	switch rule.Type {
	case "word_borrowing":
		return fmt.Sprintf("Borrowed %s terms from %s", rule.SemanticField, sourceLang.Name)
	case "semantic_shift":
		return fmt.Sprintf("Semantic shift in %s vocabulary", rule.SemanticField)
	case "calque_formation":
		return fmt.Sprintf("Formed calques for %s concepts", rule.SemanticField)
	default:
		return fmt.Sprintf("Applied lexical rule: %s", rule.Description)
	}
}

// applyPhonologicalRule applies a specific phonological rule to the target language.
func applyPhonologicalRule(targetLang, sourceLang *Language, rule PhonologicalRule, intensity float32) string {
	// Simulate phonological borrowing based on the rule type
	switch rule.Type {
	case "sound_addition":
		return fmt.Sprintf("Added foreign sound %s in %s context", rule.SourceSound, rule.Context)
	case "sound_modification":
		return fmt.Sprintf("Modified %s to %s in %s context", rule.SourceSound, rule.TargetSound, rule.Context)
	case "phonotactic_change":
		return fmt.Sprintf("Changed phonotactic rules for %s context", rule.Context)
	default:
		return fmt.Sprintf("Applied phonological rule: %s", rule.Description)
	}
}

// applyGrammaticalRule applies a specific grammatical rule to the target language.
func applyGrammaticalRule(targetLang, sourceLang *Language, rule GrammaticalRule, intensity float32) string {
	// Simulate grammatical borrowing based on the rule type
	switch rule.Type {
	case "morphology_addition":
		return fmt.Sprintf("Added %s %s marking", rule.Category, rule.Feature)
	case "syntax_change":
		return fmt.Sprintf("Modified %s syntax for %s", rule.Category, rule.Feature)
	case "agreement_modification":
		return fmt.Sprintf("Modified agreement patterns for %s", rule.Feature)
	default:
		return fmt.Sprintf("Applied grammatical rule: %s", rule.Description)
	}
}

// applyOrthographicRule applies a specific orthographic rule to the target language.
func applyOrthographicRule(targetLang, sourceLang *Language, rule OrthographicRule, intensity float32) string {
	// Simulate orthographic borrowing based on the rule type
	switch rule.Type {
	case "script_adoption":
		return fmt.Sprintf("Adopted %s writing conventions", rule.ScriptType)
	case "spelling_reform":
		return fmt.Sprintf("Implemented %s spelling reform", rule.ReformType)
	case "diacritic_addition":
		return fmt.Sprintf("Added diacritics for %s sounds", rule.ScriptType)
	default:
		return fmt.Sprintf("Applied orthographic rule: %s", rule.Description)
	}
}

// CreateCulturalInfluenceChange creates a linguistic change from cultural influence.
func CreateCulturalInfluenceChange(result *CulturalInfluenceResult, sourceLanguage string) LinguisticChange {
	// Create a description of the cultural influence
	description := "Cultural influence from " + sourceLanguage
	if len(result.BorrowingPatterns) > 0 {
		description += " (" + result.BorrowingPatterns[0].ContactType.String() + " contact)"
	}

	return LinguisticChange{
		ID:               "cultural_influence_" + result.LanguageID,
		Type:             LinguisticChangeTypeMorphological, // Cultural influence is typically morphological
		Description:      description,
		Details:          "Applied cultural influence patterns",
		ComplexityChange: result.ComplexityChange,
		Timestamp:        result.Timestamp,
		Era:              "cultural_contact",
		Trigger:          "cultural_influence",
		Intensity:        result.Intensity,
	}
}
