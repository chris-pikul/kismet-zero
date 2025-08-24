package lang

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/chris-pikul/kismet-zero/lang/grammar"
	"github.com/chris-pikul/kismet-zero/lang/interlingua"
	"github.com/chris-pikul/kismet-zero/lang/morphology"
	"github.com/chris-pikul/kismet-zero/lang/orthography"
	"github.com/chris-pikul/kismet-zero/lang/phoneme"
	"github.com/chris-pikul/kismet-zero/lang/phonology"
)

// LanguageConfig holds configuration parameters for language generation.
type LanguageConfig struct {
	// Phonological preferences
	PhonemeInventorySize int                `json:"phonemeInventorySize"` // Target size of phoneme inventory
	ConsonantVowelRatio  float64            `json:"consonantVowelRatio"`  // Ratio of consonants to vowels
	ComplexityTarget     LanguageComplexity `json:"complexityTarget"`     // Target complexity level

	// Cultural preferences
	WritingSystemStyle orthography.WritingStyle `json:"writingSystemStyle"` // Preferred writing system
	MorphologicalType  string                   `json:"morphologicalType"`  // agglutinative, fusional, isolating

	// Evolution parameters
	ChangeRate       float64 `json:"changeRate"`       // Rate of linguistic change over time
	ContactInfluence float64 `json:"contactInfluence"` // Susceptibility to external influence
	DialectFormation float64 `json:"dialectFormation"` // Probability of dialect formation
}

// DefaultLanguageConfig returns a sensible default configuration for language generation.
func DefaultLanguageConfig() LanguageConfig {
	return LanguageConfig{
		PhonemeInventorySize: 25,
		ConsonantVowelRatio:  2.5, // 2.5 consonants per vowel
		ComplexityTarget:     LanguageComplexityModerate,
		WritingSystemStyle:   orthography.WritingStyleAlphabetic,
		MorphologicalType:    "fusional",
		ChangeRate:           0.3,
		ContactInfluence:     0.4,
		DialectFormation:     0.2,
	}
}

// CulturalLanguageConfig returns a configuration tailored to a specific culture type.
func CulturalLanguageConfig(culture string) LanguageConfig {
	config := DefaultLanguageConfig()

	switch culture {
	case "elvish", "high_elf":
		config.PhonemeInventorySize = 35
		config.ConsonantVowelRatio = 2.0
		config.ComplexityTarget = LanguageComplexityComplex
		config.WritingSystemStyle = orthography.WritingStyleSyllabic
		config.MorphologicalType = "agglutinative"
		config.ChangeRate = 0.1 // Elves change slowly

	case "dwarvish", "mountain_dwarf":
		config.PhonemeInventorySize = 28
		config.ConsonantVowelRatio = 3.0
		config.ComplexityTarget = LanguageComplexityModerate
		config.WritingSystemStyle = orthography.WritingStyleAlphabetic
		config.MorphologicalType = "fusional"
		config.ChangeRate = 0.2

	case "orcish", "black_orc":
		config.PhonemeInventorySize = 20
		config.ConsonantVowelRatio = 1.5
		config.ComplexityTarget = LanguageComplexitySimple
		config.WritingSystemStyle = orthography.WritingStyleLogographic
		config.MorphologicalType = "isolating"
		config.ChangeRate = 0.5 // Orcs change quickly

	case "human", "northern_human":
		config.PhonemeInventorySize = 25
		config.ConsonantVowelRatio = 2.5
		config.ComplexityTarget = LanguageComplexityModerate
		config.WritingSystemStyle = orthography.WritingStyleAlphabetic
		config.MorphologicalType = "fusional"
		config.ChangeRate = 0.3

	case "ancient", "divine":
		config.PhonemeInventorySize = 40
		config.ConsonantVowelRatio = 1.8
		config.ComplexityTarget = LanguageComplexityVeryComplex
		config.WritingSystemStyle = orthography.WritingStyleLogographic
		config.MorphologicalType = "agglutinative"
		config.ChangeRate = 0.05 // Ancient languages change very slowly
	}

	return config
}

// CreateRandomLanguage generates a complete language with sensible defaults.
// This is the main entry point for quick language creation.
func CreateRandomLanguage(family, culture string, seed int64) (*Language, error) {
	config := CulturalLanguageConfig(culture)

	// Generate language ID
	langID := LanguageID{
		Family:   family,
		Branch:   culture,
		Language: generateLanguageName(family, culture, seed),
	}

	// Create base language
	lang := NewLanguage(langID, langID.Language, LanguageTypeConstructed, seed)
	lang.SetCulture(culture)
	lang.SetComplexity(config.ComplexityTarget)

	// Generate linguistic components
	if err := generatePhonology(lang, config); err != nil {
		return nil, fmt.Errorf("failed to generate phonology: %w", err)
	}

	if err := generateOrthography(lang, config); err != nil {
		return nil, fmt.Errorf("failed to generate orthography: %w", err)
	}

	if err := generateMorphology(lang, config); err != nil {
		return nil, fmt.Errorf("failed to generate morphology: %w", err)
	}

	if err := generateGrammar(lang, config); err != nil {
		return nil, fmt.Errorf("failed to generate grammar: %w", err)
	}

	// Set description
	lang.SetDescription(fmt.Sprintf("A %s language of the %s family, featuring %s morphology and %s writing system.",
		config.ComplexityTarget.String(), family, config.MorphologicalType, config.WritingSystemStyle.String()))

	return lang, nil
}

// generateLanguageName creates a culturally appropriate name for the language.
func generateLanguageName(family, culture string, seed int64) string {
	rand.Seed(seed)

	// Simple naming patterns based on culture
	switch culture {
	case "elvish", "high_elf":
		patterns := []string{"%sarin", "%sian", "%së", "%sëa", "%sëan"}
		pattern := patterns[rand.Intn(len(patterns))]
		return fmt.Sprintf(pattern, family)

	case "dwarvish", "mountain_dwarf":
		patterns := []string{"%suz", "%saz", "%siz", "%suzh", "%sazh"}
		pattern := patterns[rand.Intn(len(patterns))]
		return fmt.Sprintf(pattern, family)

	case "orcish", "black_orc":
		patterns := []string{"%s'%s", "%s-%s", "%s%s", "%s'%s'"}
		pattern := patterns[rand.Intn(len(patterns))]
		return fmt.Sprintf(pattern, family, "gul")

	case "human", "northern_human":
		patterns := []string{"%sish", "%sian", "%sic", "%sese", "%sian"}
		pattern := patterns[rand.Intn(len(patterns))]
		return fmt.Sprintf(pattern, family)

	case "ancient", "divine":
		patterns := []string{"Proto-%s", "Ancient %s", "%s Prime", "%s Root"}
		pattern := patterns[rand.Intn(len(patterns))]
		return fmt.Sprintf(pattern, family)

	default:
		return fmt.Sprintf("%s-%s", family, culture)
	}
}

// generatePhonology creates a phonology system for the language.
func generatePhonology(lang *Language, config LanguageConfig) error {
	// Start with common phonemes
	pool := phoneme.CommonPool

	// Add more phonemes based on inventory size target
	targetSize := config.PhonemeInventorySize
	currentSize := pool.Size()

	if currentSize < targetSize {
		// Add more consonants
		consonantTarget := int(float64(targetSize) * config.ConsonantVowelRatio / (1 + config.ConsonantVowelRatio))
		consonantCount := len(pool.Consonants)
		if consonantCount < consonantTarget {
			// Add more consonants from extended inventory
			extendedConsonants := getExtendedConsonants()
			for i := consonantCount; i < consonantTarget && i-consonantCount < len(extendedConsonants); i++ {
				pool.AddConsonant(extendedConsonants[i-consonantCount])
			}
		}

		// Add more vowels
		vowelTarget := targetSize - len(pool.Consonants)
		vowelCount := len(pool.Vowels)
		if vowelCount < vowelTarget {
			// Add more vowels from extended inventory
			extendedVowels := getExtendedVowels()
			for i := vowelCount; i < vowelTarget && i-vowelCount < len(extendedVowels); i++ {
				pool.AddVowel(extendedVowels[i-vowelCount])
			}
		}
	}

	// Create phonology with common syllable templates
	ph := phonology.NewPhonology(&pool)
	ph.AddTemplate(phonology.TemplateCV, 10)
	ph.AddTemplate(phonology.TemplateCVC, 8)
	ph.AddTemplate(phonology.TemplateV, 3)
	ph.AddTemplate(phonology.TemplateCCV, 5)

	lang.SetPhonology(ph)
	return nil
}

// getExtendedConsonants returns additional consonants beyond the common pool.
func getExtendedConsonants() []phoneme.Phoneme {
	return []phoneme.Phoneme{
		{Symbol: "θ", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{Manner: phoneme.ConsonantMannerFricative}, Weight: 0.6},
		{Symbol: "ð", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{Manner: phoneme.ConsonantMannerFricative}, Weight: 0.6},
		{Symbol: "ʃ", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{Manner: phoneme.ConsonantMannerSibilants}, Weight: 0.7},
		{Symbol: "ʒ", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{Manner: phoneme.ConsonantMannerSibilants}, Weight: 0.7},
		{Symbol: "x", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{Manner: phoneme.ConsonantMannerFricative}, Weight: 0.5},
		{Symbol: "ɣ", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{Manner: phoneme.ConsonantMannerFricative}, Weight: 0.5},
		{Symbol: "q", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{Manner: phoneme.ConsonantMannerPlosive}, Weight: 0.4},
		{Symbol: "ʔ", Type: phoneme.PhonemeTypeConsonant, Consonant: &phoneme.ConsonantSpec{Manner: phoneme.ConsonantMannerPlosive}, Weight: 0.3},
	}
}

// getExtendedVowels returns additional vowels beyond the common pool.
func getExtendedVowels() []phoneme.Phoneme {
	return []phoneme.Phoneme{
		{Symbol: "ɛ", Type: phoneme.PhonemeTypeVowel, Vowel: &phoneme.VowelSpec{Height: phoneme.VowelHeightMid, Backness: phoneme.VowelBacknessFront, Rounded: false}, Weight: 0.8},
		{Symbol: "ɔ", Type: phoneme.PhonemeTypeVowel, Vowel: &phoneme.VowelSpec{Height: phoneme.VowelHeightMid, Backness: phoneme.VowelBacknessBack, Rounded: true}, Weight: 0.8},
		{Symbol: "æ", Type: phoneme.PhonemeTypeVowel, Vowel: &phoneme.VowelSpec{Height: phoneme.VowelHeightLow, Backness: phoneme.VowelBacknessFront, Rounded: false}, Weight: 0.6},
		{Symbol: "ɑ", Type: phoneme.PhonemeTypeVowel, Vowel: &phoneme.VowelSpec{Height: phoneme.VowelHeightLow, Backness: phoneme.VowelBacknessBack, Rounded: false}, Weight: 0.6},
		{Symbol: "y", Type: phoneme.PhonemeTypeVowel, Vowel: &phoneme.VowelSpec{Height: phoneme.VowelHeightHigh, Backness: phoneme.VowelBacknessFront, Rounded: true}, Weight: 0.5},
		{Symbol: "ø", Type: phoneme.PhonemeTypeVowel, Vowel: &phoneme.VowelSpec{Height: phoneme.VowelHeightMid, Backness: phoneme.VowelBacknessFront, Rounded: true}, Weight: 0.5},
	}
}

// generateOrthography creates a writing system for the language.
func generateOrthography(lang *Language, config LanguageConfig) error {
	// Create a basic writing system based on the style preference
	ws := &orthography.WritingSystem{
		Style:   config.WritingSystemStyle,
		Name:    fmt.Sprintf("%s Script", lang.Name),
		Culture: lang.Culture,
	}

	// For now, create a simple mapping (this would be more sophisticated in practice)
	// The orthography package would handle the actual mapping logic
	lang.SetOrthography(ws)
	return nil
}

// generateMorphology creates a morphological system for the language.
func generateMorphology(lang *Language, config LanguageConfig) error {
	// Create a basic morpheme list
	// This is a simplified version - the morphology package would handle the actual generation
	morph := &morphology.MorphemeList{}

	// Add some basic morphemes based on the morphological type
	switch config.MorphologicalType {
	case "agglutinative":
		// Add affixes for agglutinative languages
		// This would be more sophisticated in practice
	case "fusional":
		// Add fusional morphemes
	case "isolating":
		// Add isolating morphemes
	}

	lang.SetMorphology(morph)
	return nil
}

// generateGrammar creates a grammatical system for the language.
func generateGrammar(lang *Language, config LanguageConfig) error {
	// Create a basic grammar using the grammar package

	// Create basic agreement system
	agreement := grammar.NewAgreementSystem(
		[]grammar.Case{grammar.CaseNominative, grammar.CaseAccusative},
		[]grammar.Number{grammar.NumberSingular, grammar.NumberPlural},
		[]grammar.Gender{grammar.GenderMasculine, grammar.GenderFeminine},
		[]grammar.Tense{grammar.TensePresent, grammar.TensePast},
		[]grammar.Aspect{grammar.AspectPerfective, grammar.AspectImperfective},
		[]grammar.Mood{grammar.MoodIndicative, grammar.MoodImperative},
	)

	// Create basic word builder
	wordBuilder := &morphology.WordBuilder{}

	// Create grammar with basic configuration
	gram := grammar.NewGrammar(
		grammar.WordOrderSVO,
		agreement,
		lang.Phonology,
		wordBuilder,
		lang.Culture,
	)

	lang.SetGrammar(gram)
	return nil
}

// LanguageInfo provides human-readable information about a language.
type LanguageInfo struct {
	BasicStats       BasicLanguageStats `json:"basicStats"`
	CulturalFeatures CulturalFeatures   `json:"culturalFeatures"`
	EvolutionHistory EvolutionHistory   `json:"evolutionHistory"`
	SampleVocabulary []string           `json:"sampleVocabulary"`
	GrammarPatterns  []string           `json:"grammarPatterns"`
}

// BasicLanguageStats provides basic statistical information about a language.
type BasicLanguageStats struct {
	PhonemeCount      int     `json:"phonemeCount"`
	ConsonantCount    int     `json:"consonantCount"`
	VowelCount        int     `json:"vowelCount"`
	ComplexityScore   float64 `json:"complexityScore"`
	MorphologicalType string  `json:"morphologicalType"`
	WritingSystem     string  `json:"writingSystem"`
}

// CulturalFeatures describes cultural aspects of the language.
type CulturalFeatures struct {
	Culture        string   `json:"culture"`
	LanguageFamily string   `json:"languageFamily"`
	TypicalDomains []string `json:"typicalDomains"`
	Influences     []string `json:"influences"`
	Prestige       string   `json:"prestige"`
}

// EvolutionHistory tracks the evolution of the language.
type EvolutionHistory struct {
	CreatedAt       time.Time `json:"createdAt"`
	LastEvolved     time.Time `json:"lastEvolved,omitempty"`
	ParentLanguage  string    `json:"parentLanguage,omitempty"`
	ChildLanguages  []string  `json:"childLanguages"`
	EvolutionEvents []string  `json:"evolutionEvents"`
}

// GetLanguageInfo extracts human-readable information about a language.
func GetLanguageInfo(language *Language) LanguageInfo {
	info := LanguageInfo{}

	// Basic stats
	if language.Phonology != nil {
		pool := language.Phonology.GetPool()
		if pool != nil {
			info.BasicStats.PhonemeCount = pool.Size()
			info.BasicStats.ConsonantCount = len(pool.Consonants)
			info.BasicStats.VowelCount = len(pool.Vowels)
		}
	}

	info.BasicStats.ComplexityScore = float64(language.Complexity)
	info.BasicStats.MorphologicalType = "unknown" // Would be extracted from morphology
	info.BasicStats.WritingSystem = "unknown"     // Would be extracted from orthography

	// Cultural features
	info.CulturalFeatures.Culture = language.Culture
	info.CulturalFeatures.LanguageFamily = language.ID.Family
	info.CulturalFeatures.TypicalDomains = []string{"general", "basic"} // Would be more sophisticated
	info.CulturalFeatures.Prestige = "medium"                           // Would be calculated

	// Evolution history
	info.EvolutionHistory.CreatedAt = language.CreatedAt
	info.EvolutionHistory.LastEvolved = language.EvolvedAt
	if language.ParentID != nil {
		info.EvolutionHistory.ParentLanguage = language.ParentID.String()
	}
	info.EvolutionHistory.ChildLanguages = make([]string, len(language.ChildIDs))
	for i, childID := range language.ChildIDs {
		info.EvolutionHistory.ChildLanguages[i] = childID.String()
	}

	// Sample vocabulary and grammar patterns would be generated here
	info.SampleVocabulary = []string{"sample", "words", "would", "go", "here"}
	info.GrammarPatterns = []string{"basic", "patterns", "would", "be", "listed"}

	return info
}

// ValidationResult provides information about language validation.
type ValidationResult struct {
	IsValid     bool     `json:"isValid"`
	Issues      []string `json:"issues"`
	Warnings    []string `json:"warnings"`
	Suggestions []string `json:"suggestions"`
}

// ValidateLanguage ensures language consistency and playability.
func ValidateLanguage(language *Language) ValidationResult {
	result := ValidationResult{IsValid: true}

	// Check if language is complete
	if !language.IsComplete() {
		result.IsValid = false
		missing := language.GetMissingComponents()
		for _, component := range missing {
			result.Issues = append(result.Issues, fmt.Sprintf("Missing %s component", component))
		}
	}

	// Check phonology consistency
	if language.Phonology != nil {
		pool := language.Phonology.GetPool()
		if pool != nil && pool.Size() < 10 {
			result.Warnings = append(result.Warnings, "Phoneme inventory is very small")
		}
		if pool != nil && pool.Size() > 50 {
			result.Warnings = append(result.Warnings, "Phoneme inventory is very large")
		}
	}

	// Check for reasonable complexity
	if language.Complexity == LanguageComplexityUnknown {
		result.Warnings = append(result.Warnings, "Language complexity is not set")
	}

	// Add suggestions for improvement
	if len(result.Issues) > 0 {
		result.Suggestions = append(result.Suggestions, "Complete all missing components before using the language")
	}

	return result
}

// CreateLanguageFamily generates a family of related languages from a proto-language.
func CreateLanguageFamily(protoLanguage *Language, numBranches int, evolutionTime string, seed int64) ([]*Language, error) {
	if protoLanguage == nil {
		return nil, fmt.Errorf("proto language cannot be nil")
	}

	if numBranches < 1 {
		return nil, fmt.Errorf("number of branches must be at least 1")
	}

	family := make([]*Language, 0, numBranches+1)
	family = append(family, protoLanguage)

	// Initialize child IDs slice if it doesn't exist
	if protoLanguage.ChildIDs == nil {
		protoLanguage.ChildIDs = make([]LanguageID, 0, numBranches)
	}

	// Create branch languages
	for i := 0; i < numBranches; i++ {
		branchSeed := seed + int64(i+1)
		branchID := LanguageID{
			Family:   protoLanguage.ID.Family,
			Branch:   fmt.Sprintf("branch_%d", i+1),
			Language: fmt.Sprintf("%s_branch_%d", protoLanguage.ID.Language, i+1),
		}

		// Clone the proto-language
		branchLang := protoLanguage.Clone(branchID, branchSeed)
		branchLang.SetDescription(fmt.Sprintf("A branch of %s, evolved over %s", protoLanguage.Name, evolutionTime))

		// Add to family
		family = append(family, branchLang)
		protoLanguage.ChildIDs = append(protoLanguage.ChildIDs, branchID)
	}

	return family, nil
}

// EvolveLanguage applies realistic language change over time.
func EvolveLanguage(language *Language, timePeriod string, culturalEvents []string, seed int64) (*Language, error) {
	if language == nil {
		return nil, fmt.Errorf("language cannot be nil")
	}

	// Parse time period into duration (for future use)
	_, err := parseTimePeriod(timePeriod)
	if err != nil {
		return nil, fmt.Errorf("invalid time period: %w", err)
	}

	// Create a copy for evolution
	evolvedLang := language.Clone(language.ID, seed)
	evolvedLang.EvolvedAt = time.Now()

	// Apply sound changes using the sound change engine
	soundChanges := applySoundChanges(evolvedLang, timePeriod, seed)

	// Apply morphological changes using the morphological evolution engine
	morphologicalChanges := applyMorphologicalChanges(evolvedLang, timePeriod, seed)

	// Apply orthographic changes using the orthographic evolution engine
	orthographicChanges := applyOrthographicChanges(evolvedLang, timePeriod, seed)

	// Store all linguistic changes in the language's change history
	StoreLinguisticChanges(evolvedLang, soundChanges, morphologicalChanges, orthographicChanges, timePeriod)

	// Start building the evolution description
	var evolutionParts []string

	// Add sound changes if any occurred
	if len(soundChanges) > 0 {
		evolutionParts = append(evolutionParts, fmt.Sprintf("sound changes: %v", soundChanges))
	}

	// Add morphological changes if any occurred
	if len(morphologicalChanges) > 0 {
		// Create a summary of morphological changes
		var morphChangeNames []string
		for _, change := range morphologicalChanges {
			morphChangeNames = append(morphChangeNames, change.Description)
		}
		evolutionParts = append(evolutionParts, fmt.Sprintf("morphological changes: %v", morphChangeNames))
	}

	// Add orthographic changes if any occurred
	if len(orthographicChanges) > 0 {
		// Create a summary of orthographic changes
		var orthoChangeNames []string
		for _, change := range orthographicChanges {
			orthoChangeNames = append(orthoChangeNames, change.Description)
		}
		evolutionParts = append(evolutionParts, fmt.Sprintf("orthographic changes: %v", orthoChangeNames))
	}

	// Apply basic evolution based on time period
	switch timePeriod {
	case "100_years", "short":
		// Minor phonological drift already handled by sound changes
		if evolvedLang.Phonology != nil && len(soundChanges) == 0 {
			evolutionParts = append(evolutionParts, "minor phonological drift")
		}

	case "500_years", "medium":
		// Moderate changes including some morphological simplification
		// Morphological changes are now handled by the morphological evolution engine
		if evolvedLang.Morphology != nil && len(morphologicalChanges) == 0 {
			evolutionParts = append(evolutionParts, "minor morphological drift")
		}

	case "1000_years", "long":
		// Significant changes including writing system evolution
		// Orthographic changes are now handled by the orthographic evolution engine
		if evolvedLang.Orthography != nil && len(orthographicChanges) == 0 {
			evolutionParts = append(evolutionParts, "minor orthographic drift")
		}

	case "2000_years", "very_long":
		// Major changes including potential dialect formation
		evolutionParts = append(evolutionParts, "significant linguistic divergence")

	default:
		// Custom time period
		evolutionParts = append(evolutionParts, "evolution")
	}

	// Build the final description
	if len(evolutionParts) > 0 {
		evolvedLang.SetDescription(fmt.Sprintf("%s (evolved over %s with %s)",
			language.Description, timePeriod, strings.Join(evolutionParts, ", ")))
	} else {
		evolvedLang.SetDescription(fmt.Sprintf("%s (evolved over %s)", language.Description, timePeriod))
	}

	// Apply cultural event influences
	if len(culturalEvents) > 0 {
		evolvedLang.SetDescription(fmt.Sprintf("%s, influenced by: %v", evolvedLang.Description, culturalEvents))

		// Apply cultural influence effects
		for _, event := range culturalEvents {
			switch event {
			case "trade":
				evolvedLang.SetDescription(fmt.Sprintf("%s with trade vocabulary", evolvedLang.Description))
			case "conquest":
				evolvedLang.SetDescription(fmt.Sprintf("%s with conquest influence", evolvedLang.Description))
			case "migration":
				evolvedLang.SetDescription(fmt.Sprintf("%s with migrant dialects", evolvedLang.Description))
			case "religious":
				evolvedLang.SetDescription(fmt.Sprintf("%s with religious terminology", evolvedLang.Description))
			case "magical":
				evolvedLang.SetDescription(fmt.Sprintf("%s with magical lexicon", evolvedLang.Description))
			}
		}
	}

	return evolvedLang, nil
}

// parseTimePeriod converts a time period string to a duration
func parseTimePeriod(period string) (time.Duration, error) {
	switch period {
	case "100_years", "short":
		return calculateYearsDuration(100), nil
	case "500_years", "medium":
		return calculateYearsDuration(500), nil
	case "1000_years", "long":
		return calculateYearsDuration(1000), nil
	case "2000_years", "very_long":
		return calculateYearsDuration(2000), nil
	default:
		// Try to parse as a custom duration
		if strings.HasSuffix(period, "_years") {
			yearsStr := strings.TrimSuffix(period, "_years")
			years, err := strconv.Atoi(yearsStr)
			if err != nil {
				return 0, fmt.Errorf("invalid year format: %s", period)
			}
			return calculateYearsDuration(years), nil
		}
		return 0, fmt.Errorf("unsupported time period: %s", period)
	}
}

// calculateYearsDuration calculates duration for a given number of years
// This avoids compile-time constant overflow issues
func calculateYearsDuration(years int) time.Duration {
	// Calculate in smaller steps to avoid overflow
	days := years * 365
	hours := days * 24
	return time.Duration(hours) * time.Hour
}

// parseDuration converts a string duration to time.Duration
func parseDuration(duration string) time.Duration {
	// Simple duration parsing for common formats
	switch duration {
	case "short", "brief":
		return 10 * time.Hour
	case "medium", "moderate":
		return 24 * time.Hour
	case "long", "extended":
		return 7 * 24 * time.Hour
	case "very_long", "extensive":
		return 30 * 24 * time.Hour
	default:
		// Try to parse as a number followed by a unit
		if d, err := time.ParseDuration(duration); err == nil {
			return d
		}
		// Default to medium duration
		return 24 * time.Hour
	}
}

// SimulateLanguageContact models how languages influence each other.
func SimulateLanguageContact(language1, language2 *Language, contactType string, intensity float64, duration string, seed int64) (*Language, *Language, error) {
	if language1 == nil || language2 == nil {
		return nil, nil, fmt.Errorf("both languages must be provided")
	}

	if intensity < 0.0 || intensity > 1.0 {
		return nil, nil, fmt.Errorf("intensity must be between 0.0 and 1.0")
	}

	// Convert string contact type to ContactType enum
	var enumContactType ContactType
	switch contactType {
	case "trade":
		enumContactType = ContactTypeTrade
	case "conquest":
		enumContactType = ContactTypeConquest
	case "migration":
		enumContactType = ContactTypeMigration
	case "cultural":
		enumContactType = ContactTypeCultural
	case "religious":
		enumContactType = ContactTypeReligious
	case "educational":
		enumContactType = ContactTypeEducational
	default:
		enumContactType = ContactTypeCultural // Default to cultural
	}

	// Apply cultural influence using the cultural influence engine
	result1, err := ApplyCulturalInfluence(language1, language2, enumContactType, float32(intensity), parseDuration(duration), seed)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to apply cultural influence: %w", err)
	}

	result2, err := ApplyCulturalInfluence(language2, language1, enumContactType, float32(intensity), parseDuration(duration), seed+1)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to apply cultural influence: %w", err)
	}

	// Create copies for contact simulation
	contactLang1 := language1.Clone(language1.ID, seed)
	contactLang2 := language2.Clone(language2.ID, seed+1)

	// Store cultural influence changes in the languages
	if result1 != nil {
		culturalChange1 := CreateCulturalInfluenceChange(result1, language2.Name)
		contactLang1.AddLinguisticChange(culturalChange1)
	}

	if result2 != nil {
		culturalChange2 := CreateCulturalInfluenceChange(result2, language1.Name)
		contactLang2.AddLinguisticChange(culturalChange2)
	}

	// Update descriptions to reflect the contact
	updateContactDescription(contactLang1, language2, contactType, duration, result1)
	updateContactDescription(contactLang2, language1, contactType, duration, result2)

	return contactLang1, contactLang2, nil
}

// updateContactDescription updates the language description to reflect contact effects.
func updateContactDescription(lang *Language, otherLang *Language, contactType string, duration string, result *CulturalInfluenceResult) {
	// Build a comprehensive description of the contact effects
	var effects []string

	if result != nil {
		if len(result.LexicalChanges) > 0 {
			effects = append(effects, fmt.Sprintf("lexical: %v", result.LexicalChanges))
		}
		if len(result.PhonologicalChanges) > 0 {
			effects = append(effects, fmt.Sprintf("phonological: %v", result.PhonologicalChanges))
		}
		if len(result.GrammaticalChanges) > 0 {
			effects = append(effects, fmt.Sprintf("grammatical: %v", result.GrammaticalChanges))
		}
		if len(result.OrthographicChanges) > 0 {
			effects = append(effects, fmt.Sprintf("orthographic: %v", result.OrthographicChanges))
		}
	}

	// Create the new description
	newDesc := fmt.Sprintf("%s (influenced by %s contact with %s over %s",
		lang.Description, contactType, otherLang.Name, duration)

	if len(effects) > 0 {
		newDesc += fmt.Sprintf(" - effects: %s", strings.Join(effects, ", "))
	}

	newDesc += ")"
	lang.SetDescription(newDesc)
}

// GenerateNames creates culturally appropriate names for people, places, etc.
func GenerateNames(language *Language, nameType string, count int, culturalContext string, seed int64) ([]string, error) {
	if language == nil {
		return nil, fmt.Errorf("language cannot be nil")
	}

	if count < 1 {
		return nil, fmt.Errorf("count must be at least 1")
	}

	names := make([]string, 0, count)
	rand.Seed(seed)

	// Generate names based on type and culture
	switch nameType {
	case "person":
		names = generatePersonNames(language, count, culturalContext, rand.New(rand.NewSource(seed)))

	case "place":
		names = generatePlaceNames(language, count, culturalContext, rand.New(rand.NewSource(seed)))

	case "deity":
		names = generateDeityNames(language, count, culturalContext, rand.New(rand.NewSource(seed)))

	case "artifact":
		names = generateArtifactNames(language, count, culturalContext, rand.New(rand.NewSource(seed)))

	default:
		names = generateGenericNames(language, count, culturalContext, rand.New(rand.NewSource(seed)))
	}

	return names, nil
}

// generatePersonNames creates culturally appropriate person names.
func generatePersonNames(language *Language, count int, culturalContext string, rng *rand.Rand) []string {
	names := make([]string, 0, count)

	// Simple name generation based on culture
	switch language.Culture {
	case "elvish", "high_elf":
		prefixes := []string{"El", "Gal", "Cele", "Ar", "Gil", "Fin", "Ecthel"}
		suffixes := []string{"ion", "adriel", "we", "wen", "mir", "dil", "las"}

		for i := 0; i < count; i++ {
			prefix := prefixes[rng.Intn(len(prefixes))]
			suffix := suffixes[rng.Intn(len(suffixes))]
			names = append(names, prefix+suffix)
		}

	case "dwarvish", "mountain_dwarf":
		prefixes := []string{"Thorin", "Balin", "Dwalin", "Bifur", "Bofur", "Bombur", "Dori"}
		suffixes := []string{"son", "sson", "dottir", "sson", "sson", "sson", "sson"}

		for i := 0; i < count; i++ {
			prefix := prefixes[rng.Intn(len(prefixes))]
			suffix := suffixes[rng.Intn(len(suffixes))]
			names = append(names, prefix+suffix)
		}

	case "orcish", "black_orc":
		prefixes := []string{"Grom", "Thrall", "Gul", "Mog", "Karg", "Rag", "Zug"}
		suffixes := []string{"'dan", "'gar", "'ash", "'osh", "'ash", "'osh", "'ash"}

		for i := 0; i < count; i++ {
			prefix := prefixes[rng.Intn(len(prefixes))]
			suffix := suffixes[rng.Intn(len(suffixes))]
			names = append(names, prefix+suffix)
		}

	default:
		// Generic fantasy names
		prefixes := []string{"A", "E", "I", "O", "U", "Y"}
		suffixes := []string{"ron", "lan", "dor", "mir", "wen", "las"}

		for i := 0; i < count; i++ {
			prefix := prefixes[rng.Intn(len(prefixes))]
			suffix := suffixes[rng.Intn(len(suffixes))]
			names = append(names, prefix+suffix)
		}
	}

	return names
}

// generatePlaceNames creates culturally appropriate place names.
func generatePlaceNames(language *Language, count int, culturalContext string, rng *rand.Rand) []string {
	names := make([]string, 0, count)

	// Simple place name generation
	prefixes := []string{"Mount", "River", "Forest", "City", "Valley", "Lake", "Hill"}
	suffixes := []string{"shire", "land", "dale", "burg", "ton", "ville", "ford"}

	for i := 0; i < count; i++ {
		prefix := prefixes[rng.Intn(len(prefixes))]
		suffix := suffixes[rng.Intn(len(suffixes))]
		names = append(names, prefix+suffix)
	}

	return names
}

// generateDeityNames creates culturally appropriate deity names.
func generateDeityNames(language *Language, count int, culturalContext string, rng *rand.Rand) []string {
	names := make([]string, 0, count)

	// Simple deity name generation
	prefixes := []string{"God", "Lord", "Lady", "Divine", "Sacred", "Holy", "Eternal"}
	suffixes := []string{"of Light", "of Darkness", "of Nature", "of War", "of Peace", "of Wisdom", "of Power"}

	for i := 0; i < count; i++ {
		prefix := prefixes[rng.Intn(len(prefixes))]
		suffix := suffixes[rng.Intn(len(suffixes))]
		names = append(names, prefix+" "+suffix)
	}

	return names
}

// generateArtifactNames creates culturally appropriate artifact names.
func generateArtifactNames(language *Language, count int, culturalContext string, rng *rand.Rand) []string {
	names := make([]string, 0, count)

	// Simple artifact name generation
	prefixes := []string{"Sword", "Staff", "Ring", "Crown", "Amulet", "Scroll", "Crystal"}
	suffixes := []string{"of Power", "of Destiny", "of Light", "of Darkness", "of Wisdom", "of Strength", "of Magic"}

	for i := 0; i < count; i++ {
		prefix := prefixes[rng.Intn(len(prefixes))]
		suffix := suffixes[rng.Intn(len(suffixes))]
		names = append(names, prefix+" "+suffix)
	}

	return names
}

// generateGenericNames creates generic names when type is not specified.
func generateGenericNames(language *Language, count int, culturalContext string, rng *rand.Rand) []string {
	names := make([]string, 0, count)

	// Simple generic name generation
	prefixes := []string{"A", "E", "I", "O", "U", "Y"}
	suffixes := []string{"ron", "lan", "dor", "mir", "wen", "las"}

	for i := 0; i < count; i++ {
		prefix := prefixes[rng.Intn(len(prefixes))]
		suffix := suffixes[rng.Intn(len(suffixes))]
		names = append(names, prefix+suffix)
	}

	return names
}

// GenerateText creates sample texts in the generated language.
func GenerateText(language *Language, textType string, length string, culturalContext string, seed int64) (string, error) {
	if language == nil {
		return "", fmt.Errorf("language cannot be nil")
	}

	rand.Seed(seed)

	// Generate text based on type and length
	switch textType {
	case "lore":
		return generateLoreText(language, length, culturalContext, rand.New(rand.NewSource(seed)))

	case "poetry":
		return generatePoetryText(language, length, culturalContext, rand.New(rand.NewSource(seed)))

	case "dialogue":
		return generateDialogueText(language, length, culturalContext, rand.New(rand.NewSource(seed)))

	case "inscription":
		return generateInscriptionText(language, length, culturalContext, rand.New(rand.NewSource(seed)))

	default:
		return generateGenericText(language, length, culturalContext, rand.New(rand.NewSource(seed)))
	}
}

// generateLoreText creates lore/legendary text.
func generateLoreText(language *Language, length string, culturalContext string, rng *rand.Rand) (string, error) {
	// Simple lore text generation
	texts := []string{
		"In the ancient times, when the world was young...",
		"Legends speak of a time when magic flowed freely...",
		"The old texts tell of heroes who walked these lands...",
		"Ancient wisdom passed down through generations...",
		"Stories of old speak of great battles and noble deeds...",
	}

	return texts[rng.Intn(len(texts))], nil
}

// generatePoetryText creates poetic text.
func generatePoetryText(language *Language, length string, culturalContext string, rng *rand.Rand) (string, error) {
	// Simple poetry generation
	poems := []string{
		"Moonlight dances on silver streams,\nAncient wisdom in elven dreams.",
		"Mountains tall and valleys deep,\nSecrets that the ages keep.",
		"Stars above and earth below,\nStories that the winds do know.",
	}

	return poems[rng.Intn(len(poems))], nil
}

// generateDialogueText creates dialogue text.
func generateDialogueText(language *Language, length string, culturalContext string, rng *rand.Rand) (string, error) {
	// Simple dialogue generation
	dialogues := []string{
		"\"Greetings, traveler. What brings you to these lands?\"",
		"\"The ancient ones speak of times long past.\"",
		"\"Magic flows through all things, if you know how to listen.\"",
	}

	return dialogues[rng.Intn(len(dialogues))], nil
}

// generateInscriptionText creates inscription text.
func generateInscriptionText(language *Language, length string, culturalContext string, rng *rand.Rand) (string, error) {
	// Simple inscription generation
	inscriptions := []string{
		"Here lies the resting place of ancient kings.",
		"Enter not, lest you disturb the eternal sleep.",
		"Knowledge and wisdom await those who seek.",
	}

	return inscriptions[rng.Intn(len(inscriptions))], nil
}

// generateGenericText creates generic text when type is not specified.
func generateGenericText(language *Language, length string, culturalContext string, rng *rand.Rand) (string, error) {
	// Simple generic text generation
	texts := []string{
		"This is a sample text in the constructed language.",
		"Words flow like water, carrying meaning and sound.",
		"The language speaks of culture and tradition.",
	}

	return texts[rng.Intn(len(texts))], nil
}

// CreateCulturalLanguageSet generates a set of languages for a fantasy culture.
func CreateCulturalLanguageSet(culture string, numLanguages int, diversity string, seed int64) ([]*Language, error) {
	if numLanguages < 1 {
		return nil, fmt.Errorf("number of languages must be at least 1")
	}

	languages := make([]*Language, 0, numLanguages)

	// Create main language
	mainLang, err := CreateRandomLanguage(culture, culture, seed)
	if err != nil {
		return nil, fmt.Errorf("failed to create main language: %w", err)
	}
	languages = append(languages, mainLang)

	// Create regional variants
	for i := 1; i < numLanguages; i++ {
		variantSeed := seed + int64(i)
		variantLang, err := CreateRandomLanguage(culture, culture, variantSeed)
		if err != nil {
			return nil, fmt.Errorf("failed to create variant language %d: %w", i, err)
		}

		// Ensure the variant has the correct culture name
		variantLang.SetCulture(culture)

		// Set appropriate complexity based on culture
		switch culture {
		case "elvish":
			variantLang.SetComplexity(LanguageComplexityComplex)
		case "dwarvish":
			variantLang.SetComplexity(LanguageComplexityModerate)
		case "human":
			variantLang.SetComplexity(LanguageComplexityModerate)
		case "orcish":
			variantLang.SetComplexity(LanguageComplexitySimple)
		case "ancient":
			variantLang.SetComplexity(LanguageComplexityVeryComplex)
		}

		// Make it a child of the main language
		variantLang.ParentID = &mainLang.ID
		mainLang.ChildIDs = append(mainLang.ChildIDs, variantLang.ID)

		languages = append(languages, variantLang)
	}

	return languages, nil
}

// ExportLanguage exports language data for external use.
func ExportLanguage(language *Language, format string) ([]byte, error) {
	if language == nil {
		return nil, fmt.Errorf("language cannot be nil")
	}

	// For now, just return a simple JSON representation
	// In a real implementation, this would use proper JSON marshaling
	exportData := fmt.Sprintf(`{
		"id": "%s",
		"name": "%s",
		"culture": "%s",
		"complexity": "%s",
		"description": "%s"
	}`, language.ID.String(), language.Name, language.Culture, language.Complexity.String(), language.Description)

	return []byte(exportData), nil
}

// ConfigureInterlinguaServices configures interlingua services for a language with English support.
func ConfigureInterlinguaServices(language *Language) error {
	if language == nil {
		return fmt.Errorf("language cannot be nil")
	}

	// Create new interlingua services
	services := interlingua.NewInterlinguaServices()

	// Create English analyzer and realizer
	conceptResolver := interlingua.NewInMemoryConceptResolver()
	englishAnalyzer := interlingua.NewEnglishAnalyzer(conceptResolver)
	englishRealizer := interlingua.NewEnglishRealizer(language.Seed, conceptResolver)

	// Add English services
	services.AddAnalyzer("en", englishAnalyzer)
	services.AddRealizer("en", englishRealizer)

	// For now, we'll create placeholder services for the generated language
	// In a full implementation, these would be actual analyzers/realizers for the generated language
	generatedAnalyzer := &GeneratedLanguageAnalyzer{
		language: language,
		resolver: conceptResolver,
	}
	generatedRealizer := &GeneratedLanguageRealizer{
		language: language,
		resolver: conceptResolver,
	}

	services.AddAnalyzer("generated", generatedAnalyzer)
	services.AddRealizer("generated", generatedRealizer)

	// Set the services on the language
	language.SetInterlinguaServices(services)

	return nil
}

// GeneratedLanguageAnalyzer is a placeholder analyzer for generated languages.
// In a full implementation, this would be a sophisticated analyzer that understands
// the generated language's grammar, morphology, and syntax.
type GeneratedLanguageAnalyzer struct {
	language *Language
	resolver interlingua.ConceptResolver
}

// LanguageCode returns the language code for the generated language.
func (a *GeneratedLanguageAnalyzer) LanguageCode() string {
	return "generated"
}

// Analyze converts generated language tokens to an Interlingua Document.
// This is a simplified implementation that creates basic semantic structures.
func (a *GeneratedLanguageAnalyzer) Analyze(tokens []string) (interlingua.Document, error) {
	if len(tokens) == 0 {
		return interlingua.Document{}, fmt.Errorf("no tokens to analyze")
	}

	// Create a basic document structure
	doc := interlingua.NewDocument("generated")
	doc.Trace.SetSourceLanguage("generated")

	// For now, create a simple entity and event structure
	// In a full implementation, this would parse the actual generated language grammar

	// Create a basic entity for the first token
	if len(tokens) > 0 {
		entity := interlingua.Entity{
			ID:        "e1",
			Concept:   interlingua.ConceptID(tokens[0]),
			LemmaHint: tokens[0],
			Feats: interlingua.NominalFeatures{
				Number: "sg",
				Case:   "nom",
			},
		}
		interlingua.AddEntity(&doc, entity)
	}

	// Create a basic event if we have multiple tokens
	if len(tokens) > 1 {
		event := interlingua.Event{
			ID:        "v1",
			Predicate: interlingua.ConceptID(tokens[1]),
			LemmaHint: tokens[1],
			TAM: interlingua.TAM{
				Tense:    "pres",
				Aspect:   "simple",
				Mood:     "ind",
				Polarity: "pos",
			},
			Roles: map[interlingua.Role]string{
				interlingua.RoleAgent: "e1",
			},
			Valency: 1,
		}
		interlingua.AddEvent(&doc, event)
	}

	return doc, nil
}

// GeneratedLanguageRealizer is a placeholder realizer for generated languages.
// In a full implementation, this would be a sophisticated realizer that generates
// text according to the generated language's grammar, morphology, and syntax.
type GeneratedLanguageRealizer struct {
	language *Language
	resolver interlingua.ConceptResolver
}

// LanguageCode returns the language code for the generated language.
func (r *GeneratedLanguageRealizer) LanguageCode() string {
	return "generated"
}

// Realize converts an Interlingua Document to generated language tokens.
// This is a simplified implementation that generates basic text structures.
func (r *GeneratedLanguageRealizer) Realize(doc interlingua.Document) ([]string, interlingua.Trace, error) {
	if err := interlingua.Validate(doc); err != nil {
		return nil, interlingua.Trace{}, fmt.Errorf("document validation failed: %w", err)
	}

	// Initialize trace
	trace := interlingua.Trace{
		SourceLanguageID: doc.Trace.SourceLanguageID,
		ConstructionID:   doc.Trace.ConstructionID,
		TokenToNode:      make(map[int]string),
		Notes:            make([]interlingua.Note, 0),
	}

	var tokens []string
	tokenIdx := 0

	// Build entity lookup map
	entityMap := make(map[string]interlingua.Entity)
	for _, entity := range doc.Entities {
		entityMap[entity.ID] = entity
	}

	// Process each event
	for _, event := range doc.Events {
		// For now, generate simple token sequences
		// In a full implementation, this would use the language's grammar and morphology

		// Add subject (agent)
		if agentID, hasAgent := event.Roles[interlingua.RoleAgent]; hasAgent {
			if agent, exists := entityMap[agentID]; exists {
				tokens = append(tokens, agent.LemmaHint)
				trace.TokenToNode[tokenIdx] = agent.ID
				tokenIdx++
			}
		}

		// Add verb
		if event.LemmaHint != "" {
			tokens = append(tokens, event.LemmaHint)
		} else {
			// Fallback to predicate
			predicate := string(event.Predicate)
			if idx := strings.Index(predicate, "-"); idx != -1 {
				predicate = predicate[:idx]
			}
			tokens = append(tokens, predicate)
		}
		trace.TokenToNode[tokenIdx] = event.ID
		tokenIdx++

		// Add other arguments
		for role, entityID := range event.Roles {
			if role == interlingua.RoleAgent {
				continue // Already handled
			}
			if entity, exists := entityMap[entityID]; exists {
				tokens = append(tokens, entity.LemmaHint)
				trace.TokenToNode[tokenIdx] = entity.ID
				tokenIdx++
			}
		}
	}

	return tokens, trace, nil
}

// TranslateToEnglish translates text from the generated language to English.
func TranslateToEnglish(language *Language, text string) (string, error) {
	if language == nil {
		return "", fmt.Errorf("language cannot be nil")
	}

	services := language.Interlingua()
	if !services.HasAnalyzer("en") {
		return "", fmt.Errorf("English analyzer not configured for this language")
	}

	// Create interlingua pipeline and perform translation
	pipeline := NewInterlinguaPipeline(services)

	// For now, we need to determine the source language code
	// This is a simplified approach - in a full implementation, we'd have proper language codes
	sourceLang := "generated" // Placeholder for generated language

	// Check if we have a realizer for the generated language
	if !services.HasRealizer(sourceLang) {
		// Fallback to placeholder translation
		return fmt.Sprintf("[English translation of: %s]", text), nil
	}

	result, err := pipeline.Translate(text, sourceLang, "en")
	if err != nil {
		// Fallback to placeholder translation
		return fmt.Sprintf("[English translation of: %s]", text), nil
	}

	return result.TargetText, nil
}

// TranslateFromEnglish translates English text to the generated language.
func TranslateFromEnglish(language *Language, englishText string) (string, error) {
	if language == nil {
		return "", fmt.Errorf("language cannot be nil")
	}

	services := language.Interlingua()
	if !services.HasRealizer("en") {
		return "", fmt.Errorf("English realizer not configured for this language")
	}

	// Create interlingua pipeline and perform translation
	pipeline := NewInterlinguaPipeline(services)

	// For now, we need to determine the target language code
	// This is a simplified approach - in a full implementation, we'd have proper language codes
	targetLang := "generated" // Placeholder for generated language

	// Check if we have an analyzer for the generated language
	if !services.HasAnalyzer(targetLang) {
		// Fallback to placeholder translation
		return fmt.Sprintf("[%s translation of: %s]", language.Name, englishText), nil
	}

	result, err := pipeline.Translate(englishText, "en", targetLang)
	if err != nil {
		// Fallback to placeholder translation
		return fmt.Sprintf("[%s translation of: %s]", language.Name, englishText), nil
	}

	return result.TargetText, nil
}

// CreateGeographicDialect creates a new geographic dialect based on regional factors.
func CreateGeographicDialect(
	parentLang *Language,
	dialectID LanguageID,
	dialectName string,
	region GeographicRegion,
	era string,
	seed int64,
) (*Language, error) {
	if parentLang == nil {
		return nil, fmt.Errorf("parent language cannot be nil")
	}

	// Clone the parent language and modify it to represent the dialect
	dialectLang := parentLang.Clone(dialectID, seed)
	dialectLang.Name = dialectName
	dialectLang.SetDescription(fmt.Sprintf("Geographic dialect of %s formed in %s region", parentLang.Name, region.Name))

	// Create geographic influence model and apply geographic effects
	geographicModel := NewGeographicInfluenceModel()
	geographicChanges := geographicModel.ApplyGeographicInfluence(dialectLang, region)

	// Store all geographic influence changes
	for _, change := range geographicChanges {
		dialectLang.AddLinguisticChange(change)
	}

	// Create a dialect formation change record
	dialectChange := LinguisticChange{
		ID:          fmt.Sprintf("dialect_formation_%s", dialectID.String()),
		Type:        LinguisticChangeTypeDialectal,
		Description: fmt.Sprintf("Formation of geographic dialect: %s", dialectName),
		Details:     fmt.Sprintf("Dialect formed in %s region with %s climate and %s terrain", region.Name, region.Climate, region.Terrain),
		Timestamp:   time.Now(),
		Era:         era,
		Trigger:     "dialect_formation",
		Intensity:   0.6, // Moderate intensity for geographic dialects
	}
	dialectLang.AddLinguisticChange(dialectChange)

	// Update parent-child relationships
	if dialectLang.ParentID == nil {
		dialectLang.ParentID = &parentLang.ID
	}
	if parentLang.ChildIDs == nil {
		parentLang.ChildIDs = make([]LanguageID, 0)
	}
	parentLang.ChildIDs = append(parentLang.ChildIDs, dialectLang.ID)

	return dialectLang, nil
}

// CreateSocialDialect creates a new social dialect based on social factors.
func CreateSocialDialect(
	parentLang *Language,
	dialectID LanguageID,
	dialectName string,
	socialClass string,
	urbanization float32,
	era string,
	seed int64,
) (*Language, error) {
	if parentLang == nil {
		return nil, fmt.Errorf("parent language cannot be nil")
	}

	// Clone the parent language and modify it to represent the dialect
	dialectLang := parentLang.Clone(dialectID, seed)
	dialectLang.Name = dialectName
	dialectLang.SetDescription(fmt.Sprintf("Social dialect of %s for %s class", parentLang.Name, socialClass))

	// Create social stratification model and apply social effects
	socialModel := NewSocialStratificationModel()

	// Determine education level based on social class
	educationLevel := getEducationLevelForClass(socialClass)

	// Determine occupation based on social class
	occupation := getOccupationForClass(socialClass)

	// Determine social mobility based on social class
	socialMobility := getSocialMobilityForClass(socialClass)

	// Determine prestige based on social class
	prestige := getPrestigeForClass(socialClass)

	// Determine gender (using neutral for now, can be enhanced later)
	gender := "neutral"

	// Apply social stratification effects
	socialChanges := socialModel.ApplySocialStratification(
		dialectLang,
		socialClass,
		educationLevel,
		occupation,
		socialMobility,
		prestige,
		gender,
	)

	// Store all social stratification changes
	for _, change := range socialChanges {
		dialectLang.AddLinguisticChange(change)
	}

	// Create a dialect formation change record
	dialectChange := LinguisticChange{
		ID:          fmt.Sprintf("dialect_formation_%s", dialectID.String()),
		Type:        LinguisticChangeTypeDialectal,
		Description: fmt.Sprintf("Formation of social dialect: %s", dialectName),
		Details:     fmt.Sprintf("Dialect formed for %s class with urbanization level %f", socialClass, urbanization),
		Timestamp:   time.Now(),
		Era:         era,
		Trigger:     "dialect_formation",
		Intensity:   0.5, // Moderate intensity for social dialects
	}
	dialectLang.AddLinguisticChange(dialectChange)

	// Update parent-child relationships
	if dialectLang.ParentID == nil {
		dialectLang.ParentID = &parentLang.ID
	}
	if parentLang.ChildIDs == nil {
		parentLang.ChildIDs = make([]LanguageID, 0)
	}
	parentLang.ChildIDs = append(parentLang.ChildIDs, dialectLang.ID)

	return dialectLang, nil
}

// getEducationLevelForClass determines the typical education level for a social class.
func getEducationLevelForClass(socialClass string) string {
	switch socialClass {
	case "noble", "upper":
		return "high"
	case "middle":
		return "medium"
	case "working":
		return "low"
	case "lower":
		return "none"
	default:
		return "medium"
	}
}

// getOccupationForClass determines the typical occupation for a social class.
func getOccupationForClass(socialClass string) string {
	switch socialClass {
	case "noble":
		return "professional"
	case "upper":
		return "professional"
	case "middle":
		return "skilled"
	case "working":
		return "unskilled"
	case "lower":
		return "unskilled"
	default:
		return "skilled"
	}
}

// getSocialMobilityForClass determines the typical social mobility for a social class.
func getSocialMobilityForClass(socialClass string) string {
	switch socialClass {
	case "noble":
		return "none"
	case "upper":
		return "low"
	case "middle":
		return "medium"
	case "working":
		return "medium"
	case "lower":
		return "high"
	default:
		return "medium"
	}
}

// getPrestigeForClass determines the typical prestige level for a social class.
func getPrestigeForClass(socialClass string) string {
	switch socialClass {
	case "noble":
		return "high"
	case "upper":
		return "high"
	case "middle":
		return "medium"
	case "working":
		return "low"
	case "lower":
		return "low"
	default:
		return "medium"
	}
}

// EvolveDialect evolves a dialect over time, potentially creating new sub-dialects.
func EvolveDialect(
	dialect *Language,
	era string,
	duration time.Duration,
	seed int64,
) (*Language, error) {
	if dialect == nil {
		return nil, fmt.Errorf("dialect cannot be nil")
	}

	// Create a copy for evolution
	evolvedLang := dialect.Clone(dialect.ID, seed)

	// Apply dialect-specific changes based on the era and duration
	changes := applyDialectalChanges(evolvedLang, era, duration, seed)

	// Store the evolution event in linguistic changes
	evolutionChange := LinguisticChange{
		ID:          fmt.Sprintf("dialect_evolution_%s", dialect.ID.String()),
		Type:        LinguisticChangeTypeDialectal,
		Description: fmt.Sprintf("Evolution of dialect: %s", dialect.Name),
		Details:     fmt.Sprintf("Dialect evolved over %v during %s era with %d changes", duration, era, len(changes)),
		Timestamp:   time.Now(),
		Era:         era,
		Trigger:     "dialect_evolution",
		Intensity:   0.4, // Moderate intensity for dialect evolution
	}
	evolvedLang.AddLinguisticChange(evolutionChange)

	// Update the description to reflect evolution
	evolvedLang.SetDescription(fmt.Sprintf("Evolved dialect of %s during %s era", dialect.Name, era))

	return evolvedLang, nil
}

// applyDialectalChanges applies dialect-specific changes over time.
func applyDialectalChanges(dialect *Language, era string, duration time.Duration, seed int64) []string {
	rand.Seed(seed)
	var changes []string

	// Apply changes based on dialect type and era
	if rand.Float32() < 0.3 { // 30% chance of phonological drift
		changes = append(changes, "phonological_drift")
	}

	if rand.Float32() < 0.2 { // 20% chance of lexical innovation
		changes = append(changes, "lexical_innovation")
	}

	if rand.Float32() < 0.15 { // 15% chance of grammatical simplification
		changes = append(changes, "grammatical_simplification")
	}

	// Era-specific changes
	switch era {
	case "early_evolution":
		if rand.Float32() < 0.4 {
			changes = append(changes, "conservative_retention")
		}
	case "middle_evolution":
		if rand.Float32() < 0.3 {
			changes = append(changes, "moderate_innovation")
		}
	case "late_evolution":
		if rand.Float32() < 0.25 {
			changes = append(changes, "rapid_change")
		}
	}

	return changes
}

// GeographicRegion represents a geographical area where a dialect is spoken.
type GeographicRegion struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Climate      string  `json:"climate,omitempty"` // "tropical", "temperate", "arctic", etc.
	Terrain      string  `json:"terrain,omitempty"` // "mountain", "coastal", "plains", etc.
	Population   int     `json:"population,omitempty"`
	Urbanization float32 `json:"urbanization,omitempty"` // 0.0 to 1.0, rural to urban
	Isolation    float32 `json:"isolation,omitempty"`    // 0.0 to 1.0, connected to isolated
	TradeRoutes  bool    `json:"tradeRoutes,omitempty"`  // Whether major trade routes pass through
	PortAccess   bool    `json:"portAccess,omitempty"`   // Whether has access to ports/waterways
}

// GeographicInfluenceModel represents the influence of geographic factors on linguistic development.
type GeographicInfluenceModel struct {
	ClimateEffects      map[string]ClimateEffect      `json:"climateEffects"`
	TerrainEffects      map[string]TerrainEffect      `json:"terrainEffects"`
	PopulationEffects   map[string]PopulationEffect   `json:"populationEffects"`
	UrbanizationEffects map[string]UrbanizationEffect `json:"urbanizationEffects"`
	IsolationEffects    map[string]IsolationEffect    `json:"isolationEffects"`
	TradeEffects        map[string]TradeEffect        `json:"tradeEffects"`
}

// ClimateEffect represents how climate affects linguistic development.
type ClimateEffect struct {
	PhonologicalChanges []string `json:"phonologicalChanges,omitempty"`
	LexicalChanges      []string `json:"lexicalChanges,omitempty"`
	GrammaticalChanges  []string `json:"grammaticalChanges,omitempty"`
	ComplexityImpact    float32  `json:"complexityImpact"` // -1.0 to 1.0
	Description         string   `json:"description"`
}

// TerrainEffect represents how terrain affects linguistic development.
type TerrainEffect struct {
	PhonologicalChanges []string `json:"phonologicalChanges,omitempty"`
	LexicalChanges      []string `json:"lexicalChanges,omitempty"`
	GrammaticalChanges  []string `json:"grammaticalChanges,omitempty"`
	ComplexityImpact    float32  `json:"complexityImpact"` // -1.0 to 1.0
	Description         string   `json:"description"`
}

// PopulationEffect represents how population density affects linguistic development.
type PopulationEffect struct {
	PhonologicalChanges []string `json:"phonologicalChanges,omitempty"`
	LexicalChanges      []string `json:"lexicalChanges,omitempty"`
	GrammaticalChanges  []string `json:"grammaticalChanges,omitempty"`
	ComplexityImpact    float32  `json:"complexityImpact"` // -1.0 to 1.0
	Description         string   `json:"description"`
}

// UrbanizationEffect represents how urbanization affects linguistic development.
type UrbanizationEffect struct {
	PhonologicalChanges []string `json:"phonologicalChanges,omitempty"`
	LexicalChanges      []string `json:"lexicalChanges,omitempty"`
	GrammaticalChanges  []string `json:"grammaticalChanges,omitempty"`
	ComplexityImpact    float32  `json:"complexityImpact"` // -1.0 to 1.0
	Description         string   `json:"description"`
}

// IsolationEffect represents how geographic isolation affects linguistic development.
type IsolationEffect struct {
	PhonologicalChanges []string `json:"phonologicalChanges,omitempty"`
	LexicalChanges      []string `json:"lexicalChanges,omitempty"`
	GrammaticalChanges  []string `json:"grammaticalChanges,omitempty"`
	ComplexityImpact    float32  `json:"complexityImpact"` // -1.0 to 1.0
	Description         string   `json:"description"`
}

// TradeEffect represents how trade routes affect linguistic development.
type TradeEffect struct {
	PhonologicalChanges []string `json:"phonologicalChanges,omitempty"`
	LexicalChanges      []string `json:"lexicalChanges,omitempty"`
	GrammaticalChanges  []string `json:"grammaticalChanges,omitempty"`
	ComplexityImpact    float32  `json:"complexityImpact"` // -1.0 to 1.0
	Description         string   `json:"description"`
}

// NewGeographicInfluenceModel creates a new geographic influence model with default effects.
func NewGeographicInfluenceModel() *GeographicInfluenceModel {
	model := &GeographicInfluenceModel{
		ClimateEffects:      make(map[string]ClimateEffect),
		TerrainEffects:      make(map[string]TerrainEffect),
		PopulationEffects:   make(map[string]PopulationEffect),
		UrbanizationEffects: make(map[string]UrbanizationEffect),
		IsolationEffects:    make(map[string]IsolationEffect),
		TradeEffects:        make(map[string]TradeEffect),
	}

	// Initialize climate effects
	model.initializeClimateEffects()
	model.initializeTerrainEffects()
	model.initializePopulationEffects()
	model.initializeUrbanizationEffects()
	model.initializeIsolationEffects()
	model.initializeTradeEffects()

	return model
}

// initializeClimateEffects sets up the default climate effects on linguistic development.
func (gim *GeographicInfluenceModel) initializeClimateEffects() {
	// Tropical climate effects
	gim.ClimateEffects["tropical"] = ClimateEffect{
		PhonologicalChanges: []string{"vowel_harmony_shift", "consonant_softening", "nasal_assimilation"},
		LexicalChanges:      []string{"agricultural_terms", "seasonal_vocabulary", "natural_phenomena"},
		GrammaticalChanges:  []string{"aspect_marking", "evidentiality", "classifier_system"},
		ComplexityImpact:    0.2,
		Description:         "Tropical climate promotes vowel harmony and agricultural vocabulary",
	}

	// Temperate climate effects
	gim.ClimateEffects["temperate"] = ClimateEffect{
		PhonologicalChanges: []string{"consonant_cluster_simplification", "vowel_reduction", "stress_patterns"},
		LexicalChanges:      []string{"seasonal_terms", "agricultural_vocabulary", "weather_expressions"},
		GrammaticalChanges:  []string{"tense_marking", "case_system", "agreement_patterns"},
		ComplexityImpact:    0.0,
		Description:         "Temperate climate allows for balanced linguistic complexity",
	}

	// Arctic climate effects
	gim.ClimateEffects["arctic"] = ClimateEffect{
		PhonologicalChanges: []string{"consonant_cluster_preservation", "vowel_quality_retention", "ejective_consonants"},
		LexicalChanges:      []string{"snow_ice_terms", "hunting_vocabulary", "survival_expressions"},
		GrammaticalChanges:  []string{"ergative_case", "incorporation", "polysynthesis"},
		ComplexityImpact:    0.3,
		Description:         "Arctic climate preserves complex phonological structures and promotes polysynthesis",
	}

	// Desert climate effects
	gim.ClimateEffects["desert"] = ClimateEffect{
		PhonologicalChanges: []string{"pharyngeal_consonants", "emphatic_consonants", "vowel_length"},
		LexicalChanges:      []string{"water_terms", "nomadic_vocabulary", "trade_expressions"},
		GrammaticalChanges:  []string{"root_pattern_system", "derivational_morphology", "gender_system"},
		ComplexityImpact:    0.1,
		Description:         "Desert climate promotes root-based morphology and trade vocabulary",
	}
}

// initializeTerrainEffects sets up the default terrain effects on linguistic development.
func (gim *GeographicInfluenceModel) initializeTerrainEffects() {
	// Mountain terrain effects
	gim.TerrainEffects["mountain"] = TerrainEffect{
		PhonologicalChanges: []string{"consonant_cluster_simplification", "vowel_harmony", "tone_system"},
		LexicalChanges:      []string{"elevation_terms", "mountain_vocabulary", "isolation_expressions"},
		GrammaticalChanges:  []string{"conservative_grammar", "case_system", "agreement_patterns"},
		ComplexityImpact:    0.1,
		Description:         "Mountain terrain preserves conservative grammar and promotes tone systems",
	}

	// Coastal terrain effects
	gim.TerrainEffects["coastal"] = TerrainEffect{
		PhonologicalChanges: []string{"liquid_consonants", "vowel_length", "stress_patterns"},
		LexicalChanges:      []string{"maritime_vocabulary", "fishing_terms", "trade_expressions"},
		GrammaticalChanges:  []string{"prepositional_system", "serial_verbs", "classifier_system"},
		ComplexityImpact:    -0.1,
		Description:         "Coastal terrain promotes maritime vocabulary and simplified grammar",
	}

	// Plains terrain effects
	gim.TerrainEffects["plains"] = TerrainEffect{
		PhonologicalChanges: []string{"vowel_reduction", "consonant_assimilation", "rhythm_patterns"},
		LexicalChanges:      []string{"agricultural_terms", "herding_vocabulary", "mobility_expressions"},
		GrammaticalChanges:  []string{"word_order_flexibility", "agreement_system", "derivational_morphology"},
		ComplexityImpact:    0.0,
		Description:         "Plains terrain allows for flexible word order and agricultural vocabulary",
	}

	// Forest terrain effects
	gim.TerrainEffects["forest"] = TerrainEffect{
		PhonologicalChanges: []string{"nasal_consonants", "vowel_quality", "tone_patterns"},
		LexicalChanges:      []string{"forest_vocabulary", "hunting_terms", "natural_expressions"},
		GrammaticalChanges:  []string{"classifier_system", "evidentiality", "aspect_marking"},
		ComplexityImpact:    0.2,
		Description:         "Forest terrain promotes classifier systems and hunting vocabulary",
	}
}

// initializePopulationEffects sets up the default population effects on linguistic development.
func (gim *GeographicInfluenceModel) initializePopulationEffects() {
	// High population density effects
	gim.PopulationEffects["high"] = PopulationEffect{
		PhonologicalChanges: []string{"fast_speech_patterns", "consonant_reduction", "vowel_merging"},
		LexicalChanges:      []string{"urban_slang", "specialized_terms", "borrowed_vocabulary"},
		GrammaticalChanges:  []string{"grammatical_simplification", "analytic_structures", "word_order_fixation"},
		ComplexityImpact:    -0.2,
		Description:         "High population density promotes grammatical simplification and urban vocabulary",
	}

	// Medium population density effects
	gim.PopulationEffects["medium"] = PopulationEffect{
		PhonologicalChanges: []string{"balanced_phonology", "moderate_assimilation", "stable_vowels"},
		LexicalChanges:      []string{"balanced_vocabulary", "traditional_terms", "innovative_expressions"},
		GrammaticalChanges:  []string{"balanced_grammar", "moderate_complexity", "stable_patterns"},
		ComplexityImpact:    0.0,
		Description:         "Medium population density maintains balanced linguistic complexity",
	}

	// Low population density effects
	gim.PopulationEffects["low"] = PopulationEffect{
		PhonologicalChanges: []string{"conservative_phonology", "phoneme_preservation", "vowel_quality"},
		LexicalChanges:      []string{"traditional_vocabulary", "archaic_terms", "conservative_expressions"},
		GrammaticalChanges:  []string{"complex_grammar", "irregular_patterns", "conservative_structures"},
		ComplexityImpact:    0.3,
		Description:         "Low population density preserves complex grammar and traditional vocabulary",
	}
}

// initializeUrbanizationEffects sets up the default urbanization effects on linguistic development.
func (gim *GeographicInfluenceModel) initializeUrbanizationEffects() {
	// High urbanization effects
	gim.UrbanizationEffects["high"] = UrbanizationEffect{
		PhonologicalChanges: []string{"fast_speech", "consonant_cluster_reduction", "vowel_centralization"},
		LexicalChanges:      []string{"urban_slang", "technical_terms", "borrowed_words"},
		GrammaticalChanges:  []string{"grammatical_simplification", "analytic_structures", "word_order_fixation"},
		ComplexityImpact:    -0.3,
		Description:         "High urbanization promotes grammatical simplification and urban vocabulary",
	}

	// Medium urbanization effects
	gim.UrbanizationEffects["medium"] = UrbanizationEffect{
		PhonologicalChanges: []string{"balanced_phonology", "moderate_assimilation", "stable_patterns"},
		LexicalChanges:      []string{"balanced_vocabulary", "traditional_terms", "innovative_words"},
		GrammaticalChanges:  []string{"balanced_grammar", "moderate_complexity", "stable_structures"},
		ComplexityImpact:    0.0,
		Description:         "Medium urbanization maintains balanced linguistic complexity",
	}

	// Low urbanization effects
	gim.UrbanizationEffects["low"] = UrbanizationEffect{
		PhonologicalChanges: []string{"conservative_phonology", "phoneme_preservation", "traditional_patterns"},
		LexicalChanges:      []string{"rural_vocabulary", "traditional_terms", "conservative_words"},
		GrammaticalChanges:  []string{"complex_grammar", "irregular_patterns", "conservative_structures"},
		ComplexityImpact:    0.2,
		Description:         "Low urbanization preserves complex grammar and traditional vocabulary",
	}
}

// initializeIsolationEffects sets up the default isolation effects on linguistic development.
func (gim *GeographicInfluenceModel) initializeIsolationEffects() {
	// High isolation effects
	gim.IsolationEffects["high"] = IsolationEffect{
		PhonologicalChanges: []string{"conservative_phonology", "phoneme_preservation", "archaic_patterns"},
		LexicalChanges:      []string{"archaic_vocabulary", "traditional_terms", "conservative_expressions"},
		GrammaticalChanges:  []string{"complex_grammar", "irregular_patterns", "conservative_structures"},
		ComplexityImpact:    0.4,
		Description:         "High isolation preserves complex grammar and archaic vocabulary",
	}

	// Medium isolation effects
	gim.IsolationEffects["medium"] = IsolationEffect{
		PhonologicalChanges: []string{"moderate_phonology", "balanced_assimilation", "stable_patterns"},
		LexicalChanges:      []string{"balanced_vocabulary", "traditional_terms", "innovative_expressions"},
		GrammaticalChanges:  []string{"balanced_grammar", "moderate_complexity", "stable_structures"},
		ComplexityImpact:    0.0,
		Description:         "Medium isolation maintains balanced linguistic complexity",
	}

	// Low isolation effects
	gim.IsolationEffects["low"] = IsolationEffect{
		PhonologicalChanges: []string{"innovative_phonology", "rapid_assimilation", "change_patterns"},
		LexicalChanges:      []string{"innovative_vocabulary", "borrowed_terms", "modern_expressions"},
		GrammaticalChanges:  []string{"simplified_grammar", "regular_patterns", "innovative_structures"},
		ComplexityImpact:    -0.2,
		Description:         "Low isolation promotes grammatical simplification and innovative vocabulary",
	}
}

// initializeTradeEffects sets up the default trade effects on linguistic development.
func (gim *GeographicInfluenceModel) initializeTradeEffects() {
	// Trade route effects
	gim.TradeEffects["trade_route"] = TradeEffect{
		PhonologicalChanges: []string{"phoneme_borrowing", "sound_adaptation", "accent_patterns"},
		LexicalChanges:      []string{"loanwords", "trade_terms", "cultural_vocabulary"},
		GrammaticalChanges:  []string{"grammatical_borrowing", "syntactic_influence", "structural_adaptation"},
		ComplexityImpact:    0.1,
		Description:         "Trade routes promote linguistic borrowing and cultural vocabulary",
	}

	// Port access effects
	gim.TradeEffects["port_access"] = TradeEffect{
		PhonologicalChanges: []string{"accent_adaptation", "phoneme_borrowing", "rhythm_patterns"},
		LexicalChanges:      []string{"maritime_terms", "trade_vocabulary", "cultural_words"},
		GrammaticalChanges:  []string{"grammatical_influence", "syntactic_borrowing", "structural_adaptation"},
		ComplexityImpact:    0.0,
		Description:         "Port access promotes maritime vocabulary and grammatical influence",
	}

	// No trade effects
	gim.TradeEffects["no_trade"] = TradeEffect{
		PhonologicalChanges: []string{"conservative_phonology", "phoneme_preservation", "traditional_patterns"},
		LexicalChanges:      []string{"traditional_vocabulary", "local_terms", "conservative_expressions"},
		GrammaticalChanges:  []string{"conservative_grammar", "traditional_patterns", "stable_structures"},
		ComplexityImpact:    0.1,
		Description:         "No trade preserves traditional vocabulary and conservative grammar",
	}
}

// ApplyGeographicInfluence applies geographic factors to a language, creating dialect-specific features.
func (gim *GeographicInfluenceModel) ApplyGeographicInfluence(language *Language, region GeographicRegion) []LinguisticChange {
	var changes []LinguisticChange

	// Apply climate effects
	if region.Climate != "" {
		if effect, exists := gim.ClimateEffects[region.Climate]; exists {
			climateChange := LinguisticChange{
				ID:               fmt.Sprintf("geographic_climate_%s_%s", language.ID.String(), region.Climate),
				Type:             LinguisticChangeTypeDialectal,
				Description:      fmt.Sprintf("Climate influence: %s", effect.Description),
				Details:          fmt.Sprintf("Applied %s climate effects to dialect", region.Climate),
				ComplexityChange: effect.ComplexityImpact,
				Timestamp:        time.Now(),
				Era:              "geographic_formation",
				Trigger:          "climate_influence",
				Intensity:        0.6,
			}
			changes = append(changes, climateChange)
		}
	}

	// Apply terrain effects
	if region.Terrain != "" {
		if effect, exists := gim.TerrainEffects[region.Terrain]; exists {
			terrainChange := LinguisticChange{
				ID:               fmt.Sprintf("geographic_terrain_%s_%s", language.ID.String(), region.Terrain),
				Type:             LinguisticChangeTypeDialectal,
				Description:      fmt.Sprintf("Terrain influence: %s", effect.Description),
				Details:          fmt.Sprintf("Applied %s terrain effects to dialect", region.Terrain),
				ComplexityChange: effect.ComplexityImpact,
				Timestamp:        time.Now(),
				Era:              "geographic_formation",
				Trigger:          "terrain_influence",
				Intensity:        0.5,
			}
			changes = append(changes, terrainChange)
		}
	}

	// Apply population effects
	populationLevel := gim.getPopulationLevel(region.Population)
	if effect, exists := gim.PopulationEffects[populationLevel]; exists {
		populationChange := LinguisticChange{
			ID:               fmt.Sprintf("geographic_population_%s_%s", language.ID.String(), populationLevel),
			Type:             LinguisticChangeTypeDialectal,
			Description:      fmt.Sprintf("Population influence: %s", effect.Description),
			Details:          fmt.Sprintf("Applied %s population effects to dialect", populationLevel),
			ComplexityChange: effect.ComplexityImpact,
			Timestamp:        time.Now(),
			Era:              "geographic_formation",
			Trigger:          "population_influence",
			Intensity:        0.4,
		}
		changes = append(changes, populationChange)
	}

	// Apply urbanization effects
	urbanizationLevel := gim.getUrbanizationLevel(region.Urbanization)
	if effect, exists := gim.UrbanizationEffects[urbanizationLevel]; exists {
		urbanizationChange := LinguisticChange{
			ID:               fmt.Sprintf("geographic_urbanization_%s_%s", language.ID.String(), urbanizationLevel),
			Type:             LinguisticChangeTypeDialectal,
			Description:      fmt.Sprintf("Urbanization influence: %s", effect.Description),
			Details:          fmt.Sprintf("Applied %s urbanization effects to dialect", urbanizationLevel),
			ComplexityChange: effect.ComplexityImpact,
			Timestamp:        time.Now(),
			Era:              "geographic_formation",
			Trigger:          "urbanization_influence",
			Intensity:        0.5,
		}
		changes = append(changes, urbanizationChange)
	}

	// Apply isolation effects
	isolationLevel := gim.getIsolationLevel(region.Isolation)
	if effect, exists := gim.IsolationEffects[isolationLevel]; exists {
		isolationChange := LinguisticChange{
			ID:               fmt.Sprintf("geographic_isolation_%s_%s", language.ID.String(), isolationLevel),
			Type:             LinguisticChangeTypeDialectal,
			Description:      fmt.Sprintf("Isolation influence: %s", effect.Description),
			Details:          fmt.Sprintf("Applied %s isolation effects to dialect", isolationLevel),
			ComplexityChange: effect.ComplexityImpact,
			Timestamp:        time.Now(),
			Era:              "geographic_formation",
			Trigger:          "isolation_influence",
			Intensity:        0.7,
		}
		changes = append(changes, isolationChange)
	}

	// Apply trade effects
	tradeLevel := gim.getTradeLevel(region.TradeRoutes, region.PortAccess)
	if effect, exists := gim.TradeEffects[tradeLevel]; exists {
		tradeChange := LinguisticChange{
			ID:               fmt.Sprintf("geographic_trade_%s_%s", language.ID.String(), tradeLevel),
			Type:             LinguisticChangeTypeDialectal,
			Description:      fmt.Sprintf("Trade influence: %s", effect.Description),
			Details:          fmt.Sprintf("Applied %s trade effects to dialect", tradeLevel),
			ComplexityChange: effect.ComplexityImpact,
			Timestamp:        time.Now(),
			Era:              "geographic_formation",
			Trigger:          "trade_influence",
			Intensity:        0.4,
		}
		changes = append(changes, tradeChange)
	}

	return changes
}

// getPopulationLevel categorizes population into high, medium, or low.
func (gim *GeographicInfluenceModel) getPopulationLevel(population int) string {
	if population > 10000 {
		return "high"
	} else if population > 1000 {
		return "medium"
	} else {
		return "low"
	}
}

// getUrbanizationLevel categorizes urbanization into high, medium, or low.
func (gim *GeographicInfluenceModel) getUrbanizationLevel(urbanization float32) string {
	if urbanization > 0.7 {
		return "high"
	} else if urbanization > 0.3 {
		return "medium"
	} else {
		return "low"
	}
}

// getIsolationLevel categorizes isolation into high, medium, or low.
func (gim *GeographicInfluenceModel) getIsolationLevel(isolation float32) string {
	if isolation > 0.7 {
		return "high"
	} else if isolation > 0.3 {
		return "medium"
	} else {
		return "low"
	}
}

// getTradeLevel categorizes trade access based on trade routes and port access.
func (gim *GeographicInfluenceModel) getTradeLevel(tradeRoutes, portAccess bool) string {
	if tradeRoutes {
		return "trade_route"
	} else if portAccess {
		return "port_access"
	} else {
		return "no_trade"
	}
}

// SocialStratificationModel represents the influence of social factors on linguistic development.
type SocialStratificationModel struct {
	ClassEffects      map[string]ClassEffect      `json:"classEffects"`
	EducationEffects  map[string]EducationEffect  `json:"educationEffects"`
	OccupationEffects map[string]OccupationEffect `json:"occupationEffects"`
	MobilityEffects   map[string]MobilityEffect   `json:"mobilityEffects"`
	PrestigeEffects   map[string]PrestigeEffect   `json:"prestigeEffects"`
	GenderEffects     map[string]GenderEffect     `json:"genderEffects"`
}

// ClassEffect represents how social class affects linguistic development.
type ClassEffect struct {
	PhonologicalChanges []string `json:"phonologicalChanges,omitempty"`
	LexicalChanges      []string `json:"lexicalChanges,omitempty"`
	GrammaticalChanges  []string `json:"grammaticalChanges,omitempty"`
	ComplexityImpact    float32  `json:"complexityImpact"` // -1.0 to 1.0
	Description         string   `json:"description"`
	PrestigeLevel       float32  `json:"prestigeLevel"` // 0.0 to 1.0
}

// EducationEffect represents how education level affects linguistic development.
type EducationEffect struct {
	PhonologicalChanges []string `json:"phonologicalChanges,omitempty"`
	LexicalChanges      []string `json:"lexicalChanges,omitempty"`
	GrammaticalChanges  []string `json:"grammaticalChanges,omitempty"`
	ComplexityImpact    float32  `json:"complexityImpact"` // -1.0 to 1.0
	Description         string   `json:"description"`
	FormalityLevel      float32  `json:"formalityLevel"` // 0.0 to 1.0
}

// OccupationEffect represents how occupation affects linguistic development.
type OccupationEffect struct {
	PhonologicalChanges []string `json:"phonologicalChanges,omitempty"`
	LexicalChanges      []string `json:"lexicalChanges,omitempty"`
	GrammaticalChanges  []string `json:"grammaticalChanges,omitempty"`
	ComplexityImpact    float32  `json:"complexityImpact"` // -1.0 to 1.0
	Description         string   `json:"description"`
	SpecializationLevel float32  `json:"specializationLevel"` // 0.0 to 1.0
}

// MobilityEffect represents how social mobility affects linguistic development.
type MobilityEffect struct {
	PhonologicalChanges []string `json:"phonologicalChanges,omitempty"`
	LexicalChanges      []string `json:"lexicalChanges,omitempty"`
	GrammaticalChanges  []string `json:"grammaticalChanges,omitempty"`
	ComplexityImpact    float32  `json:"complexityImpact"` // -1.0 to 1.0
	Description         string   `json:"description"`
	AdaptationLevel     float32  `json:"adaptationLevel"` // 0.0 to 1.0
}

// PrestigeEffect represents how social prestige affects linguistic development.
type PrestigeEffect struct {
	PhonologicalChanges []string `json:"phonologicalChanges,omitempty"`
	LexicalChanges      []string `json:"lexicalChanges,omitempty"`
	GrammaticalChanges  []string `json:"grammaticalChanges,omitempty"`
	ComplexityImpact    float32  `json:"complexityImpact"` // -1.0 to 1.0
	Description         string   `json:"description"`
	InfluenceLevel      float32  `json:"influenceLevel"` // 0.0 to 1.0
}

// GenderEffect represents how gender affects linguistic development.
type GenderEffect struct {
	PhonologicalChanges []string `json:"phonologicalChanges,omitempty"`
	LexicalChanges      []string `json:"lexicalChanges,omitempty"`
	GrammaticalChanges  []string `json:"grammaticalChanges,omitempty"`
	ComplexityImpact    float32  `json:"complexityImpact"` // -1.0 to 1.0
	Description         string   `json:"description"`
	DistinctionLevel    float32  `json:"distinctionLevel"` // 0.0 to 1.0
}

// NewSocialStratificationModel creates a new social stratification model with default effects.
func NewSocialStratificationModel() *SocialStratificationModel {
	model := &SocialStratificationModel{
		ClassEffects:      make(map[string]ClassEffect),
		EducationEffects:  make(map[string]EducationEffect),
		OccupationEffects: make(map[string]OccupationEffect),
		MobilityEffects:   make(map[string]MobilityEffect),
		PrestigeEffects:   make(map[string]PrestigeEffect),
		GenderEffects:     make(map[string]GenderEffect),
	}

	// Initialize all effect types
	model.initializeClassEffects()
	model.initializeEducationEffects()
	model.initializeOccupationEffects()
	model.initializeMobilityEffects()
	model.initializePrestigeEffects()
	model.initializeGenderEffects()

	return model
}

// initializeClassEffects sets up the default social class effects on linguistic development.
func (ssm *SocialStratificationModel) initializeClassEffects() {
	// Upper class effects
	ssm.ClassEffects["upper"] = ClassEffect{
		PhonologicalChanges: []string{"prestige_pronunciation", "conservative_phonology", "standard_speech"},
		LexicalChanges:      []string{"formal_vocabulary", "prestige_terms", "educated_expressions"},
		GrammaticalChanges:  []string{"complex_grammar", "formal_register", "conservative_structures"},
		ComplexityImpact:    0.3,
		Description:         "Upper class promotes conservative grammar and formal vocabulary",
		PrestigeLevel:       0.9,
	}

	// Middle class effects
	ssm.ClassEffects["middle"] = ClassEffect{
		PhonologicalChanges: []string{"standard_pronunciation", "moderate_phonology", "balanced_speech"},
		LexicalChanges:      []string{"standard_vocabulary", "professional_terms", "balanced_expressions"},
		GrammaticalChanges:  []string{"standard_grammar", "moderate_formality", "balanced_structures"},
		ComplexityImpact:    0.0,
		Description:         "Middle class maintains standard grammar and balanced vocabulary",
		PrestigeLevel:       0.6,
	}

	// Working class effects
	ssm.ClassEffects["working"] = ClassEffect{
		PhonologicalChanges: []string{"colloquial_pronunciation", "simplified_phonology", "informal_speech"},
		LexicalChanges:      []string{"colloquial_vocabulary", "occupational_terms", "informal_expressions"},
		GrammaticalChanges:  []string{"simplified_grammar", "informal_register", "practical_structures"},
		ComplexityImpact:    -0.2,
		Description:         "Working class promotes simplified grammar and practical vocabulary",
		PrestigeLevel:       0.3,
	}

	// Lower class effects
	ssm.ClassEffects["lower"] = ClassEffect{
		PhonologicalChanges: []string{"vernacular_pronunciation", "simplified_phonology", "casual_speech"},
		LexicalChanges:      []string{"vernacular_vocabulary", "basic_terms", "casual_expressions"},
		GrammaticalChanges:  []string{"basic_grammar", "casual_register", "simple_structures"},
		ComplexityImpact:    -0.4,
		Description:         "Lower class uses basic grammar and vernacular vocabulary",
		PrestigeLevel:       0.1,
	}

	// Noble class effects
	ssm.ClassEffects["noble"] = ClassEffect{
		PhonologicalChanges: []string{"aristocratic_pronunciation", "archaic_phonology", "refined_speech"},
		LexicalChanges:      []string{"aristocratic_vocabulary", "noble_terms", "refined_expressions"},
		GrammaticalChanges:  []string{"archaic_grammar", "formal_register", "traditional_structures"},
		ComplexityImpact:    0.4,
		Description:         "Noble class preserves archaic grammar and aristocratic vocabulary",
		PrestigeLevel:       1.0,
	}
}

// initializeEducationEffects sets up the default education effects on linguistic development.
func (ssm *SocialStratificationModel) initializeEducationEffects() {
	// High education effects
	ssm.EducationEffects["high"] = EducationEffect{
		PhonologicalChanges: []string{"precise_pronunciation", "standard_phonology", "educated_speech"},
		LexicalChanges:      []string{"academic_vocabulary", "technical_terms", "educated_expressions"},
		GrammaticalChanges:  []string{"complex_grammar", "formal_register", "educated_structures"},
		ComplexityImpact:    0.3,
		Description:         "High education promotes complex grammar and academic vocabulary",
		FormalityLevel:      0.8,
	}

	// Medium education effects
	ssm.EducationEffects["medium"] = EducationEffect{
		PhonologicalChanges: []string{"standard_pronunciation", "moderate_phonology", "standard_speech"},
		LexicalChanges:      []string{"standard_vocabulary", "general_terms", "standard_expressions"},
		GrammaticalChanges:  []string{"standard_grammar", "moderate_formality", "standard_structures"},
		ComplexityImpact:    0.0,
		Description:         "Medium education maintains standard grammar and vocabulary",
		FormalityLevel:      0.5,
	}

	// Low education effects
	ssm.EducationEffects["low"] = EducationEffect{
		PhonologicalChanges: []string{"basic_pronunciation", "simplified_phonology", "basic_speech"},
		LexicalChanges:      []string{"basic_vocabulary", "simple_terms", "basic_expressions"},
		GrammaticalChanges:  []string{"basic_grammar", "informal_register", "simple_structures"},
		ComplexityImpact:    -0.2,
		Description:         "Low education uses basic grammar and simple vocabulary",
		FormalityLevel:      0.2,
	}

	// No education effects
	ssm.EducationEffects["none"] = EducationEffect{
		PhonologicalChanges: []string{"vernacular_pronunciation", "simplified_phonology", "vernacular_speech"},
		LexicalChanges:      []string{"vernacular_vocabulary", "basic_terms", "vernacular_expressions"},
		GrammaticalChanges:  []string{"vernacular_grammar", "casual_register", "vernacular_structures"},
		ComplexityImpact:    -0.3,
		Description:         "No education results in vernacular grammar and basic vocabulary",
		FormalityLevel:      0.0,
	}
}

// initializeOccupationEffects sets up the default occupation effects on linguistic development.
func (ssm *SocialStratificationModel) initializeOccupationEffects() {
	// Professional occupation effects
	ssm.OccupationEffects["professional"] = OccupationEffect{
		PhonologicalChanges: []string{"precise_pronunciation", "standard_phonology", "professional_speech"},
		LexicalChanges:      []string{"professional_vocabulary", "technical_terms", "specialized_expressions"},
		GrammaticalChanges:  []string{"formal_grammar", "professional_register", "technical_structures"},
		ComplexityImpact:    0.2,
		Description:         "Professional occupations promote formal grammar and technical vocabulary",
		SpecializationLevel: 0.8,
	}

	// Skilled occupation effects
	ssm.OccupationEffects["skilled"] = OccupationEffect{
		PhonologicalChanges: []string{"standard_pronunciation", "moderate_phonology", "skilled_speech"},
		LexicalChanges:      []string{"skilled_vocabulary", "trade_terms", "practical_expressions"},
		GrammaticalChanges:  []string{"standard_grammar", "moderate_formality", "practical_structures"},
		ComplexityImpact:    0.0,
		Description:         "Skilled occupations maintain standard grammar and trade vocabulary",
		SpecializationLevel: 0.6,
	}

	// Unskilled occupation effects
	ssm.OccupationEffects["unskilled"] = OccupationEffect{
		PhonologicalChanges: []string{"basic_pronunciation", "simplified_phonology", "basic_speech"},
		LexicalChanges:      []string{"basic_vocabulary", "simple_terms", "basic_expressions"},
		GrammaticalChanges:  []string{"basic_grammar", "informal_register", "simple_structures"},
		ComplexityImpact:    -0.1,
		Description:         "Unskilled occupations use basic grammar and simple vocabulary",
		SpecializationLevel: 0.2,
	}

	// Agricultural occupation effects
	ssm.OccupationEffects["agricultural"] = OccupationEffect{
		PhonologicalChanges: []string{"rural_pronunciation", "conservative_phonology", "rural_speech"},
		LexicalChanges:      []string{"agricultural_vocabulary", "rural_terms", "traditional_expressions"},
		GrammaticalChanges:  []string{"traditional_grammar", "rural_register", "conservative_structures"},
		ComplexityImpact:    0.1,
		Description:         "Agricultural occupations preserve traditional grammar and rural vocabulary",
		SpecializationLevel: 0.4,
	}
}

// initializeMobilityEffects sets up the default social mobility effects on linguistic development.
func (ssm *SocialStratificationModel) initializeMobilityEffects() {
	// High mobility effects
	ssm.MobilityEffects["high"] = MobilityEffect{
		PhonologicalChanges: []string{"adaptive_pronunciation", "flexible_phonology", "adaptive_speech"},
		LexicalChanges:      []string{"adaptive_vocabulary", "borrowed_terms", "flexible_expressions"},
		GrammaticalChanges:  []string{"adaptive_grammar", "flexible_register", "innovative_structures"},
		ComplexityImpact:    0.1,
		Description:         "High social mobility promotes adaptive grammar and flexible vocabulary",
		AdaptationLevel:     0.8,
	}

	// Medium mobility effects
	ssm.MobilityEffects["medium"] = MobilityEffect{
		PhonologicalChanges: []string{"moderate_pronunciation", "balanced_phonology", "moderate_speech"},
		LexicalChanges:      []string{"balanced_vocabulary", "mixed_terms", "balanced_expressions"},
		GrammaticalChanges:  []string{"balanced_grammar", "moderate_register", "balanced_structures"},
		ComplexityImpact:    0.0,
		Description:         "Medium social mobility maintains balanced grammar and vocabulary",
		AdaptationLevel:     0.5,
	}

	// Low mobility effects
	ssm.MobilityEffects["low"] = MobilityEffect{
		PhonologicalChanges: []string{"conservative_pronunciation", "stable_phonology", "conservative_speech"},
		LexicalChanges:      []string{"conservative_vocabulary", "traditional_terms", "stable_expressions"},
		GrammaticalChanges:  []string{"conservative_grammar", "stable_register", "traditional_structures"},
		ComplexityImpact:    0.1,
		Description:         "Low social mobility preserves conservative grammar and traditional vocabulary",
		AdaptationLevel:     0.2,
	}

	// No mobility effects
	ssm.MobilityEffects["none"] = MobilityEffect{
		PhonologicalChanges: []string{"fixed_pronunciation", "archaic_phonology", "fixed_speech"},
		LexicalChanges:      []string{"archaic_vocabulary", "fixed_terms", "archaic_expressions"},
		GrammaticalChanges:  []string{"archaic_grammar", "fixed_register", "archaic_structures"},
		ComplexityImpact:    0.2,
		Description:         "No social mobility preserves archaic grammar and fixed vocabulary",
		AdaptationLevel:     0.0,
	}
}

// initializePrestigeEffects sets up the default prestige effects on linguistic development.
func (ssm *SocialStratificationModel) initializePrestigeEffects() {
	// High prestige effects
	ssm.PrestigeEffects["high"] = PrestigeEffect{
		PhonologicalChanges: []string{"prestige_pronunciation", "refined_phonology", "prestige_speech"},
		LexicalChanges:      []string{"prestige_vocabulary", "refined_terms", "elegant_expressions"},
		GrammaticalChanges:  []string{"refined_grammar", "elegant_register", "sophisticated_structures"},
		ComplexityImpact:    0.3,
		Description:         "High prestige promotes refined grammar and elegant vocabulary",
		InfluenceLevel:      0.9,
	}

	// Medium prestige effects
	ssm.PrestigeEffects["medium"] = PrestigeEffect{
		PhonologicalChanges: []string{"standard_pronunciation", "moderate_phonology", "standard_speech"},
		LexicalChanges:      []string{"standard_vocabulary", "moderate_terms", "standard_expressions"},
		GrammaticalChanges:  []string{"standard_grammar", "moderate_register", "standard_structures"},
		ComplexityImpact:    0.0,
		Description:         "Medium prestige maintains standard grammar and vocabulary",
		InfluenceLevel:      0.5,
	}

	// Low prestige effects
	ssm.PrestigeEffects["low"] = PrestigeEffect{
		PhonologicalChanges: []string{"basic_pronunciation", "simple_phonology", "basic_speech"},
		LexicalChanges:      []string{"basic_vocabulary", "simple_terms", "basic_expressions"},
		GrammaticalChanges:  []string{"basic_grammar", "simple_register", "basic_structures"},
		ComplexityImpact:    -0.1,
		Description:         "Low prestige uses basic grammar and simple vocabulary",
		InfluenceLevel:      0.2,
	}
}

// initializeGenderEffects sets up the default gender effects on linguistic development.
func (ssm *SocialStratificationModel) initializeGenderEffects() {
	// Male gender effects
	ssm.GenderEffects["male"] = GenderEffect{
		PhonologicalChanges: []string{"masculine_pronunciation", "standard_phonology", "masculine_speech"},
		LexicalChanges:      []string{"masculine_vocabulary", "male_terms", "masculine_expressions"},
		GrammaticalChanges:  []string{"standard_grammar", "masculine_register", "standard_structures"},
		ComplexityImpact:    0.0,
		Description:         "Male gender maintains standard grammar with masculine vocabulary",
		DistinctionLevel:    0.3,
	}

	// Female gender effects
	ssm.GenderEffects["female"] = GenderEffect{
		PhonologicalChanges: []string{"feminine_pronunciation", "standard_phonology", "feminine_speech"},
		LexicalChanges:      []string{"feminine_vocabulary", "female_terms", "feminine_expressions"},
		GrammaticalChanges:  []string{"standard_grammar", "feminine_register", "standard_structures"},
		ComplexityImpact:    0.0,
		Description:         "Female gender maintains standard grammar with feminine vocabulary",
		DistinctionLevel:    0.3,
	}

	// Neutral gender effects
	ssm.GenderEffects["neutral"] = GenderEffect{
		PhonologicalChanges: []string{"neutral_pronunciation", "standard_phonology", "neutral_speech"},
		LexicalChanges:      []string{"neutral_vocabulary", "neutral_terms", "neutral_expressions"},
		GrammaticalChanges:  []string{"standard_grammar", "neutral_register", "standard_structures"},
		ComplexityImpact:    0.0,
		Description:         "Neutral gender maintains standard grammar with neutral vocabulary",
		DistinctionLevel:    0.0,
	}
}

// ApplySocialStratification applies social factors to a language, creating dialect-specific features.
func (ssm *SocialStratificationModel) ApplySocialStratification(
	language *Language,
	socialClass string,
	educationLevel string,
	occupation string,
	socialMobility string,
	prestige string,
	gender string,
) []LinguisticChange {
	var changes []LinguisticChange

	// Apply social class effects
	if effect, exists := ssm.ClassEffects[socialClass]; exists {
		classChange := LinguisticChange{
			ID:               fmt.Sprintf("social_class_%s_%s", language.ID.String(), socialClass),
			Type:             LinguisticChangeTypeDialectal,
			Description:      fmt.Sprintf("Social class influence: %s", effect.Description),
			Details:          fmt.Sprintf("Applied %s class effects to dialect", socialClass),
			ComplexityChange: effect.ComplexityImpact,
			Timestamp:        time.Now(),
			Era:              "social_formation",
			Trigger:          "social_class_influence",
			Intensity:        0.7,
		}
		changes = append(changes, classChange)
	}

	// Apply education effects
	if effect, exists := ssm.EducationEffects[educationLevel]; exists {
		educationChange := LinguisticChange{
			ID:               fmt.Sprintf("social_education_%s_%s", language.ID.String(), educationLevel),
			Type:             LinguisticChangeTypeDialectal,
			Description:      fmt.Sprintf("Education influence: %s", effect.Description),
			Details:          fmt.Sprintf("Applied %s education effects to dialect", educationLevel),
			ComplexityChange: effect.ComplexityImpact,
			Timestamp:        time.Now(),
			Era:              "social_formation",
			Trigger:          "education_influence",
			Intensity:        0.6,
		}
		changes = append(changes, educationChange)
	}

	// Apply occupation effects
	if effect, exists := ssm.OccupationEffects[occupation]; exists {
		occupationChange := LinguisticChange{
			ID:               fmt.Sprintf("social_occupation_%s_%s", language.ID.String(), occupation),
			Type:             LinguisticChangeTypeDialectal,
			Description:      fmt.Sprintf("Occupation influence: %s", effect.Description),
			Details:          fmt.Sprintf("Applied %s occupation effects to dialect", occupation),
			ComplexityChange: effect.ComplexityImpact,
			Timestamp:        time.Now(),
			Era:              "social_formation",
			Trigger:          "occupation_influence",
			Intensity:        0.5,
		}
		changes = append(changes, occupationChange)
	}

	// Apply social mobility effects
	if effect, exists := ssm.MobilityEffects[socialMobility]; exists {
		mobilityChange := LinguisticChange{
			ID:               fmt.Sprintf("social_mobility_%s_%s", language.ID.String(), socialMobility),
			Type:             LinguisticChangeTypeDialectal,
			Description:      fmt.Sprintf("Social mobility influence: %s", effect.Description),
			Details:          fmt.Sprintf("Applied %s social mobility effects to dialect", socialMobility),
			ComplexityChange: effect.ComplexityImpact,
			Timestamp:        time.Now(),
			Era:              "social_formation",
			Trigger:          "social_mobility_influence",
			Intensity:        0.4,
		}
		changes = append(changes, mobilityChange)
	}

	// Apply prestige effects
	if effect, exists := ssm.PrestigeEffects[prestige]; exists {
		prestigeChange := LinguisticChange{
			ID:               fmt.Sprintf("social_prestige_%s_%s", language.ID.String(), prestige),
			Type:             LinguisticChangeTypeDialectal,
			Description:      fmt.Sprintf("Prestige influence: %s", effect.Description),
			Details:          fmt.Sprintf("Applied %s prestige effects to dialect", prestige),
			ComplexityChange: effect.ComplexityImpact,
			Timestamp:        time.Now(),
			Era:              "social_formation",
			Trigger:          "prestige_influence",
			Intensity:        0.6,
		}
		changes = append(changes, prestigeChange)
	}

	// Apply gender effects
	if effect, exists := ssm.GenderEffects[gender]; exists {
		genderChange := LinguisticChange{
			ID:               fmt.Sprintf("social_gender_%s_%s", language.ID.String(), gender),
			Type:             LinguisticChangeTypeDialectal,
			Description:      fmt.Sprintf("Gender influence: %s", effect.Description),
			Details:          fmt.Sprintf("Applied %s gender effects to dialect", gender),
			ComplexityChange: effect.ComplexityImpact,
			Timestamp:        time.Now(),
			Era:              "social_formation",
			Trigger:          "gender_influence",
			Intensity:        0.3,
		}
		changes = append(changes, genderChange)
	}

	return changes
}

// DialectDetectionEngine represents an engine for automatically detecting dialect formation opportunities.
type DialectDetectionEngine struct {
	DivergenceThreshold float32 `json:"divergenceThreshold"` // Threshold for detecting significant divergence
	SimilarityThreshold float32 `json:"similarityThreshold"` // Threshold for grouping similar dialects
	MinChangeCount      int     `json:"minChangeCount"`      // Minimum changes required for dialect formation
	MaxDialectDistance  float32 `json:"maxDialectDistance"`  // Maximum distance for dialect grouping
}

// NewDialectDetectionEngine creates a new dialect detection engine with default settings.
func NewDialectDetectionEngine() *DialectDetectionEngine {
	return &DialectDetectionEngine{
		DivergenceThreshold: 0.3, // 30% divergence threshold
		SimilarityThreshold: 0.7, // 70% similarity threshold
		MinChangeCount:      5,   // Minimum 5 changes for dialect formation
		MaxDialectDistance:  0.5, // Maximum 50% distance for grouping
	}
}

// DialectFormationOpportunity represents a detected opportunity for dialect formation.
type DialectFormationOpportunity struct {
	LanguageID       LanguageID `json:"languageId"`
	DivergenceScore  float32    `json:"divergenceScore"`
	ChangeCount      int        `json:"changeCount"`
	ChangeTypes      []string   `json:"changeTypes"`
	SuggestedDialect string     `json:"suggestedDialect"`
	Confidence       float32    `json:"confidence"`
	Description      string     `json:"description"`
}

// DialectSimilarityGroup represents a group of similar dialects.
type DialectSimilarityGroup struct {
	GroupID         string       `json:"groupId"`
	GroupName       string       `json:"groupName"`
	Dialects        []LanguageID `json:"dialects"`
	SimilarityScore float32      `json:"similarityScore"`
	CommonFeatures  []string     `json:"commonFeatures"`
	Description     string       `json:"description"`
}

// DetectDialectFormationOpportunities analyzes a language and detects opportunities for dialect formation.
func (dde *DialectDetectionEngine) DetectDialectFormationOpportunities(language *Language) []DialectFormationOpportunity {
	var opportunities []DialectFormationOpportunity

	// Calculate divergence score based on linguistic changes
	divergenceScore := dde.calculateDivergenceScore(language)
	changeCount := len(language.LinguisticChanges)

	// Check if divergence meets threshold
	if divergenceScore >= dde.DivergenceThreshold && changeCount >= dde.MinChangeCount {
		// Analyze change types
		changeTypes := dde.analyzeChangeTypes(language.LinguisticChanges)

		// Generate dialect suggestion
		suggestedDialect := dde.generateDialectSuggestion(language, changeTypes)

		// Calculate confidence
		confidence := dde.calculateConfidence(divergenceScore, changeCount, changeTypes)

		// Create opportunity
		opportunity := DialectFormationOpportunity{
			LanguageID:       language.ID,
			DivergenceScore:  divergenceScore,
			ChangeCount:      changeCount,
			ChangeTypes:      changeTypes,
			SuggestedDialect: suggestedDialect,
			Confidence:       confidence,
			Description:      dde.generateOpportunityDescription(language, divergenceScore, changeTypes),
		}

		opportunities = append(opportunities, opportunity)
	}

	return opportunities
}

// calculateDivergenceScore calculates how much a language has diverged from its parent.
func (dde *DialectDetectionEngine) calculateDivergenceScore(language *Language) float32 {
	if language == nil || len(language.LinguisticChanges) == 0 {
		return 0.0
	}

	var totalDivergence float32
	var totalIntensity float32

	for _, change := range language.LinguisticChanges {
		// Weight changes by intensity and type
		weight := dde.getChangeWeight(change)
		totalDivergence += change.Intensity * weight
		totalIntensity += change.Intensity
	}

	if totalIntensity == 0 {
		return 0.0
	}

	// Normalize by total intensity and scale to 0-1 range
	return (totalDivergence / totalIntensity) * 0.8 // Scale factor for realistic divergence
}

// getChangeWeight returns the weight for a specific type of linguistic change.
func (dde *DialectDetectionEngine) getChangeWeight(change LinguisticChange) float32 {
	switch change.Type {
	case LinguisticChangeTypeSound:
		return 1.2 // Sound changes are highly significant
	case LinguisticChangeTypeMorphological:
		return 1.0 // Morphological changes are significant
	case LinguisticChangeTypeOrthographic:
		return 0.8 // Orthographic changes are moderately significant
	case LinguisticChangeTypeDialectal:
		return 1.1 // Dialectal changes are very significant
	default:
		return 1.0 // Default weight
	}
}

// analyzeChangeTypes analyzes the types of linguistic changes present.
func (dde *DialectDetectionEngine) analyzeChangeTypes(changes []LinguisticChange) []string {
	typeCounts := make(map[string]int)

	for _, change := range changes {
		typeCounts[string(change.Type)]++
	}

	var changeTypes []string
	for changeType, count := range typeCounts {
		if count > 0 {
			changeTypes = append(changeTypes, fmt.Sprintf("%s:%d", changeType, count))
		}
	}

	return changeTypes
}

// generateDialectSuggestion generates a suggested name for the new dialect.
func (dde *DialectDetectionEngine) generateDialectSuggestion(language *Language, changeTypes []string) string {
	// Base dialect name from language name
	baseName := language.Name

	// Analyze dominant change type
	dominantType := dde.getDominantChangeType(changeTypes)

	// Generate descriptive suffix
	suffix := dde.getDialectSuffix(dominantType)

	return fmt.Sprintf("%s (%s)", baseName, suffix)
}

// getDominantChangeType determines the most common type of linguistic change.
func (dde *DialectDetectionEngine) getDominantChangeType(changeTypes []string) string {
	if len(changeTypes) == 0 {
		return "general"
	}

	// Parse change types to find the most frequent
	typeCounts := make(map[string]int)
	for _, changeType := range changeTypes {
		parts := strings.Split(changeType, ":")
		if len(parts) == 2 {
			typeName := parts[0]
			if count, err := strconv.Atoi(parts[1]); err == nil {
				typeCounts[typeName] = count
			}
		}
	}

	// Find the most frequent type
	var dominantType string
	maxCount := 0
	for changeType, count := range typeCounts {
		if count > maxCount {
			maxCount = count
			dominantType = changeType
		}
	}

	return dominantType
}

// getDialectSuffix returns a descriptive suffix based on the dominant change type.
func (dde *DialectDetectionEngine) getDialectSuffix(changeType string) string {
	switch changeType {
	case "sound":
		return "phonetic"
	case "morphological":
		return "morphological"
	case "orthographic":
		return "orthographic"
	case "dialectal":
		return "regional"
	default:
		return "evolved"
	}
}

// calculateConfidence calculates the confidence level for dialect formation.
func (dde *DialectDetectionEngine) calculateConfidence(divergenceScore float32, changeCount int, changeTypes []string) float32 {
	// Base confidence from divergence score
	confidence := divergenceScore

	// Boost confidence based on change count
	if changeCount >= 10 {
		confidence += 0.2
	} else if changeCount >= 7 {
		confidence += 0.1
	}

	// Boost confidence based on change type diversity
	if len(changeTypes) >= 3 {
		confidence += 0.1
	}

	// Cap confidence at 1.0
	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

// generateOpportunityDescription generates a human-readable description of the opportunity.
func (dde *DialectDetectionEngine) generateOpportunityDescription(language *Language, divergenceScore float32, changeTypes []string) string {
	divergencePercent := int(divergenceScore * 100)

	return fmt.Sprintf(
		"Language '%s' has diverged %d%% from its parent with %d linguistic changes (%s). "+
			"Consider forming a new dialect to reflect these developments.",
		language.Name,
		divergencePercent,
		len(language.LinguisticChanges),
		strings.Join(changeTypes, ", "),
	)
}

// GroupSimilarDialects groups dialects based on linguistic similarity.
func (dde *DialectDetectionEngine) GroupSimilarDialects(dialects []*Language) []DialectSimilarityGroup {
	var groups []DialectSimilarityGroup

	// Create similarity matrix
	similarityMatrix := dde.createSimilarityMatrix(dialects)

	// Group dialects based on similarity
	groupedDialects := make(map[int]bool)
	groupID := 0

	for i, dialect1 := range dialects {
		if groupedDialects[i] {
			continue
		}

		// Start new group
		group := DialectSimilarityGroup{
			GroupID:         fmt.Sprintf("group_%d", groupID),
			GroupName:       fmt.Sprintf("Dialect Group %d", groupID+1),
			Dialects:        []LanguageID{dialect1.ID},
			SimilarityScore: 1.0,
			CommonFeatures:  []string{},
			Description:     fmt.Sprintf("Group of similar dialects starting with %s", dialect1.Name),
		}

		groupedDialects[i] = true

		// Find similar dialects
		for j, dialect2 := range dialects {
			if i == j || groupedDialects[j] {
				continue
			}

			similarity := similarityMatrix[i][j]
			if similarity >= dde.SimilarityThreshold {
				group.Dialects = append(group.Dialects, dialect2.ID)
				group.SimilarityScore = (group.SimilarityScore + similarity) / 2
				groupedDialects[j] = true
			}
		}

		// Analyze common features
		group.CommonFeatures = dde.analyzeCommonFeatures(group.Dialects, dialects)

		groups = append(groups, group)
		groupID++
	}

	return groups
}

// createSimilarityMatrix creates a matrix of similarity scores between dialects.
func (dde *DialectDetectionEngine) createSimilarityMatrix(dialects []*Language) [][]float32 {
	n := len(dialects)
	matrix := make([][]float32, n)

	for i := range matrix {
		matrix[i] = make([]float32, n)
		matrix[i][i] = 1.0 // Self-similarity is 1.0
	}

	// Calculate similarity between different dialects
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			similarity := dde.calculateDialectSimilarity(dialects[i], dialects[j])
			matrix[i][j] = similarity
			matrix[j][i] = similarity
		}
	}

	return matrix
}

// calculateDialectSimilarity calculates the similarity between two dialects.
func (dde *DialectDetectionEngine) calculateDialectSimilarity(dialect1, dialect2 *Language) float32 {
	if dialect1 == nil || dialect2 == nil {
		return 0.0
	}

	// Check if they are the same dialect (self-similarity)
	if dialect1.ID == dialect2.ID {
		return 1.0
	}

	// Check if they share a common parent
	if dialect1.ParentID != nil && dialect2.ParentID != nil && *dialect1.ParentID == *dialect2.ParentID {
		// Siblings get high similarity
		return 0.8
	}

	// Calculate similarity based on shared linguistic features
	sharedFeatures := 0
	totalFeatures := 0

	// Compare phonology systems (check if both have phonology)
	if dialect1.Phonology != nil && dialect2.Phonology != nil {
		totalFeatures++
		// For now, assume dialects with phonology systems are similar
		// In a more sophisticated system, we could compare actual phoneme inventories
		sharedFeatures++
	}

	// Compare grammar systems (check if both have grammar)
	if dialect1.Grammar != nil && dialect2.Grammar != nil {
		totalFeatures++
		// For now, assume dialects with grammar systems are similar
		// In a more sophisticated system, we could compare word order, agreement, etc.
		sharedFeatures++
	}

	// Compare morphology systems (check if both have morphology)
	if dialect1.Morphology != nil && dialect2.Morphology != nil {
		totalFeatures++
		// For now, assume dialects with morphology systems are similar
		// In a more sophisticated system, we could compare morpheme inventories
		sharedFeatures++
	}

	// Compare orthography systems (check if both have orthography)
	if dialect1.Orthography != nil && dialect2.Orthography != nil {
		totalFeatures++
		// For now, assume dialects with orthography systems are similar
		// In a more sophisticated system, we could compare writing systems
		sharedFeatures++
	}

	if totalFeatures == 0 {
		return 0.5 // Default similarity for dialects with no comparable features
	}

	return float32(sharedFeatures) / float32(totalFeatures)
}

// analyzeCommonFeatures analyzes features common to a group of dialects.
func (dde *DialectDetectionEngine) analyzeCommonFeatures(dialectIDs []LanguageID, allDialects []*Language) []string {
	if len(dialectIDs) < 2 {
		return []string{"single_dialect"}
	}

	var commonFeatures []string

	// Find dialects in the group
	var groupDialects []*Language
	for _, dialectID := range dialectIDs {
		for _, dialect := range allDialects {
			if dialect.ID == dialectID {
				groupDialects = append(groupDialects, dialect)
				break
			}
		}
	}

	if len(groupDialects) < 2 {
		return []string{"insufficient_dialects"}
	}

	// Check for common parent
	firstParent := groupDialects[0].ParentID
	if firstParent != nil {
		allShareParent := true
		for _, dialect := range groupDialects[1:] {
			if dialect.ParentID == nil || *dialect.ParentID != *firstParent {
				allShareParent = false
				break
			}
		}
		if allShareParent {
			commonFeatures = append(commonFeatures, "shared_parent")
		}
	}

	// Check for common linguistic systems
	if dde.checkCommonLinguisticSystem(groupDialects, "phonology") {
		commonFeatures = append(commonFeatures, "shared_phonology")
	}
	if dde.checkCommonLinguisticSystem(groupDialects, "grammar") {
		commonFeatures = append(commonFeatures, "shared_grammar")
	}
	if dde.checkCommonLinguisticSystem(groupDialects, "morphology") {
		commonFeatures = append(commonFeatures, "shared_morphology")
	}
	if dde.checkCommonLinguisticSystem(groupDialects, "orthography") {
		commonFeatures = append(commonFeatures, "shared_orthography")
	}

	if len(commonFeatures) == 0 {
		commonFeatures = append(commonFeatures, "minimal_shared_features")
	}

	return commonFeatures
}

// checkCommonLinguisticSystem checks if dialects share a common linguistic system.
func (dde *DialectDetectionEngine) checkCommonLinguisticSystem(dialects []*Language, systemType string) bool {
	if len(dialects) < 2 {
		return false
	}

	// Check if all dialects have the specified linguistic system
	allHaveSystem := true
	for _, dialect := range dialects {
		var hasSystem bool
		switch systemType {
		case "phonology":
			hasSystem = dialect.Phonology != nil
		case "grammar":
			hasSystem = dialect.Grammar != nil
		case "morphology":
			hasSystem = dialect.Morphology != nil
		case "orthography":
			hasSystem = dialect.Orthography != nil
		}

		if !hasSystem {
			allHaveSystem = false
			break
		}
	}

	return allHaveSystem
}

// SuggestDialectFormation analyzes a language family and suggests where dialect formation would be beneficial.
func (dde *DialectDetectionEngine) SuggestDialectFormation(languages []*Language) []DialectFormationOpportunity {
	var suggestions []DialectFormationOpportunity

	for _, language := range languages {
		opportunities := dde.DetectDialectFormationOpportunities(language)
		suggestions = append(suggestions, opportunities...)
	}

	// Sort by confidence (highest first)
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Confidence > suggestions[j].Confidence
	})

	return suggestions
}

// AutoCreateDialect automatically creates a dialect based on detected opportunities.
func (dde *DialectDetectionEngine) AutoCreateDialect(
	language *Language,
	opportunity DialectFormationOpportunity,
	era string,
	seed int64,
) (*Language, error) {
	if language == nil {
		return nil, fmt.Errorf("language cannot be nil")
	}

	// Generate dialect ID
	dialectID := LanguageID{
		Family:   language.ID.Family,
		Branch:   language.ID.Branch,
		Language: language.ID.Language,
		Dialect:  opportunity.SuggestedDialect,
	}

	// Create dialect name
	dialectName := opportunity.SuggestedDialect

	// Determine if this should be a geographic or social dialect
	// For now, create a general dialect
	dialectLang := language.Clone(dialectID, seed)
	dialectLang.Name = dialectName
	dialectLang.SetDescription(fmt.Sprintf("Automatically detected dialect: %s", opportunity.Description))

	// Create automatic dialect formation change
	formationChange := LinguisticChange{
		ID:          fmt.Sprintf("auto_dialect_formation_%s", dialectID.String()),
		Type:        LinguisticChangeTypeDialectal,
		Description: fmt.Sprintf("Automatic dialect formation: %s", dialectName),
		Details:     fmt.Sprintf("Dialect automatically formed based on %d%% divergence with %d changes", int(opportunity.DivergenceScore*100), opportunity.ChangeCount),
		Timestamp:   time.Now(),
		Era:         era,
		Trigger:     "automatic_detection",
		Intensity:   opportunity.Confidence,
	}
	dialectLang.AddLinguisticChange(formationChange)

	// Update parent-child relationships
	if dialectLang.ParentID == nil {
		dialectLang.ParentID = &language.ID
	}
	if language.ChildIDs == nil {
		language.ChildIDs = make([]LanguageID, 0)
	}
	language.ChildIDs = append(language.ChildIDs, dialectLang.ID)

	return dialectLang, nil
}

// InterlinguaPipeline represents a complete translation pipeline between any two languages.
type InterlinguaPipeline struct {
	services *interlingua.InterlinguaServices
}

// NewInterlinguaPipeline creates a new interlingua pipeline with the given services.
func NewInterlinguaPipeline(services *interlingua.InterlinguaServices) *InterlinguaPipeline {
	return &InterlinguaPipeline{
		services: services,
	}
}

// TranslationResult represents the result of a translation operation.
type TranslationResult struct {
	SourceText     string               `json:"sourceText"`
	TargetText     string               `json:"targetText"`
	SourceLanguage string               `json:"sourceLanguage"`
	TargetLanguage string               `json:"targetLanguage"`
	InterlinguaDoc interlingua.Document `json:"interlinguaDoc"`
	Confidence     float32              `json:"confidence"`
	Notes          []interlingua.Note   `json:"notes"`
	ProcessingTime time.Duration        `json:"processingTime"`
}

// Translate translates text from source language to target language using the interlingua pipeline.
func (ip *InterlinguaPipeline) Translate(
	sourceText string,
	sourceLang string,
	targetLang string,
) (*TranslationResult, error) {
	startTime := time.Now()

	// Validate that we have the necessary services
	if !ip.services.HasAnalyzer(sourceLang) {
		return nil, fmt.Errorf("no analyzer available for source language: %s", sourceLang)
	}
	if !ip.services.HasRealizer(targetLang) {
		return nil, fmt.Errorf("no realizer available for target language: %s", targetLang)
	}

	// Step 1: Analyze source text to interlingua
	analyzer := ip.services.Analyzers[sourceLang]
	interlinguaDoc, err := analyzer.Analyze(strings.Fields(sourceText))
	if err != nil {
		return nil, fmt.Errorf("failed to analyze source text: %w", err)
	}

	// Step 2: Realize interlingua to target language
	realizer := ip.services.Realizers[targetLang]
	targetTokens, trace, err := realizer.Realize(interlinguaDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to realize target text: %w", err)
	}

	// Step 3: Calculate confidence based on analysis quality
	confidence := ip.calculateTranslationConfidence(interlinguaDoc, trace)

	// Step 4: Build result
	result := &TranslationResult{
		SourceText:     sourceText,
		TargetText:     strings.Join(targetTokens, " "),
		SourceLanguage: sourceLang,
		TargetLanguage: targetLang,
		InterlinguaDoc: interlinguaDoc,
		Confidence:     confidence,
		Notes:          trace.Notes,
		ProcessingTime: time.Since(startTime),
	}

	return result, nil
}

// calculateTranslationConfidence calculates the confidence level of a translation.
func (ip *InterlinguaPipeline) calculateTranslationConfidence(doc interlingua.Document, trace interlingua.Trace) float32 {
	// Base confidence starts at 1.0
	confidence := float32(1.0)

	// Reduce confidence for each warning or error note
	for _, note := range trace.Notes {
		switch note.Severity {
		case "warn":
			confidence -= 0.1
		case "error":
			confidence -= 0.3
		case "loss":
			confidence -= 0.2
		}
	}

	// Boost confidence for well-formed documents
	if len(doc.Entities) > 0 && len(doc.Events) > 0 {
		confidence += 0.1
	}

	// Ensure confidence stays in valid range
	if confidence < 0.0 {
		confidence = 0.0
	}
	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

// TranslateToEnglish translates text from any language to English using the interlingua pipeline.
func (ip *InterlinguaPipeline) TranslateToEnglish(sourceText string, sourceLang string) (*TranslationResult, error) {
	return ip.Translate(sourceText, sourceLang, "en")
}

// TranslateFromEnglish translates English text to any target language using the interlingua pipeline.
func (ip *InterlinguaPipeline) TranslateFromEnglish(englishText string, targetLang string) (*TranslationResult, error) {
	return ip.Translate(englishText, "en", targetLang)
}

// BatchTranslate translates multiple texts in a single operation for efficiency.
func (ip *InterlinguaPipeline) BatchTranslate(
	translations []struct {
		SourceText string
		SourceLang string
		TargetLang string
	},
) ([]*TranslationResult, error) {
	var results []*TranslationResult

	for _, translation := range translations {
		result, err := ip.Translate(translation.SourceText, translation.SourceLang, translation.TargetLang)
		if err != nil {
			return nil, fmt.Errorf("batch translation failed at item %d: %w", len(results), err)
		}
		results = append(results, result)
	}

	return results, nil
}

// GetSupportedLanguagePairs returns all possible language pairs that can be translated.
func (ip *InterlinguaPipeline) GetSupportedLanguagePairs() []struct {
	Source string
	Target string
} {
	var pairs []struct {
		Source string
		Target string
	}

	sourceLangs := make([]string, 0)
	for langCode := range ip.services.Analyzers {
		sourceLangs = append(sourceLangs, langCode)
	}

	targetLangs := make([]string, 0)
	for langCode := range ip.services.Realizers {
		targetLangs = append(targetLangs, langCode)
	}

	for _, source := range sourceLangs {
		for _, target := range targetLangs {
			pairs = append(pairs, struct {
				Source string
				Target string
			}{Source: source, Target: target})
		}
	}

	return pairs
}

// ValidateTranslation validates that a translation can be performed between two languages.
func (ip *InterlinguaPipeline) ValidateTranslation(sourceLang string, targetLang string) error {
	if !ip.services.HasAnalyzer(sourceLang) {
		return fmt.Errorf("source language '%s' not supported (no analyzer available)", sourceLang)
	}
	if !ip.services.HasRealizer(targetLang) {
		return fmt.Errorf("target language '%s' not supported (no realizer available)", targetLang)
	}
	return nil
}

// MutualIntelligibilityEngine calculates how well speakers of different languages can understand each other.
type MutualIntelligibilityEngine struct {
	// Configuration parameters
	PhonologicalWeight float32 `json:"phonologicalWeight"` // Weight for phonological similarity (0.0-1.0)
	LexicalWeight      float32 `json:"lexicalWeight"`      // Weight for lexical similarity (0.0-1.0)
	GrammaticalWeight  float32 `json:"grammaticalWeight"`  // Weight for grammatical similarity (0.0-1.0)
	HistoricalWeight   float32 `json:"historicalWeight"`   // Weight for historical relationship (0.0-1.0)
	GeographicWeight   float32 `json:"geographicWeight"`   // Weight for geographic proximity (0.0-1.0)
	ContactWeight      float32 `json:"contactWeight"`      // Weight for contact history (0.0-1.0)
	MinIntelligibility float32 `json:"minIntelligibility"` // Minimum intelligibility threshold
	MaxIntelligibility float32 `json:"maxIntelligibility"` // Maximum intelligibility threshold
}

// NewMutualIntelligibilityEngine creates a new mutual intelligibility engine with default settings.
func NewMutualIntelligibilityEngine() *MutualIntelligibilityEngine {
	return &MutualIntelligibilityEngine{
		PhonologicalWeight: 0.3,  // 30% weight for phonological similarity
		LexicalWeight:      0.25, // 25% weight for lexical similarity
		GrammaticalWeight:  0.25, // 25% weight for grammatical similarity
		HistoricalWeight:   0.1,  // 10% weight for historical relationship
		GeographicWeight:   0.05, // 5% weight for geographic proximity
		ContactWeight:      0.05, // 5% weight for contact history
		MinIntelligibility: 0.0,  // 0% minimum intelligibility
		MaxIntelligibility: 1.0,  // 100% maximum intelligibility
	}
}

// IntelligibilityResult represents the result of a mutual intelligibility calculation.
type IntelligibilityResult struct {
	SourceLanguage       string   `json:"sourceLanguage"`
	TargetLanguage       string   `json:"targetLanguage"`
	OverallScore         float32  `json:"overallScore"`
	PhonologicalScore    float32  `json:"phonologicalScore"`
	LexicalScore         float32  `json:"lexicalScore"`
	GrammaticalScore     float32  `json:"grammaticalScore"`
	HistoricalScore      float32  `json:"historicalScore"`
	GeographicScore      float32  `json:"geographicScore"`
	ContactScore         float32  `json:"contactScore"`
	IntelligibilityLevel string   `json:"intelligibilityLevel"`
	Description          string   `json:"description"`
	Factors              []string `json:"factors"`
	Confidence           float32  `json:"confidence"`
}

// IntelligibilityLevel represents the level of mutual understanding between languages.
type IntelligibilityLevel string

const (
	IntelligibilityLevelNone       IntelligibilityLevel = "none"        // 0-10%: No mutual understanding
	IntelligibilityLevelMinimal    IntelligibilityLevel = "minimal"     // 10-25%: Minimal understanding
	IntelligibilityLevelLow        IntelligibilityLevel = "low"         // 25-40%: Low understanding
	IntelligibilityLevelModerate   IntelligibilityLevel = "moderate"    // 40-60%: Moderate understanding
	IntelligibilityLevelHigh       IntelligibilityLevel = "high"        // 60-80%: High understanding
	IntelligibilityLevelVeryHigh   IntelligibilityLevel = "very_high"   // 80-95%: Very high understanding
	IntelligibilityLevelNearNative IntelligibilityLevel = "near_native" // 95-100%: Near-native understanding
)

// CalculateMutualIntelligibility calculates the mutual intelligibility between two languages.
func (mie *MutualIntelligibilityEngine) CalculateMutualIntelligibility(
	lang1, lang2 *Language,
) (*IntelligibilityResult, error) {
	if lang1 == nil || lang2 == nil {
		return nil, fmt.Errorf("both languages must be provided")
	}

	if lang1.ID == lang2.ID {
		return &IntelligibilityResult{
			SourceLanguage:       lang1.Name,
			TargetLanguage:       lang2.Name,
			OverallScore:         1.0,
			PhonologicalScore:    1.0,
			LexicalScore:         1.0,
			GrammaticalScore:     1.0,
			HistoricalScore:      1.0,
			GeographicScore:      1.0,
			ContactScore:         1.0,
			IntelligibilityLevel: string(IntelligibilityLevelNearNative),
			Description:          "Same language - perfect mutual intelligibility",
			Factors:              []string{"identical_language"},
			Confidence:           1.0,
		}, nil
	}

	// Calculate individual component scores
	phonologicalScore := mie.calculatePhonologicalSimilarity(lang1, lang2)
	lexicalScore := mie.calculateLexicalSimilarity(lang1, lang2)
	grammaticalScore := mie.calculateGrammaticalSimilarity(lang1, lang2)
	historicalScore := mie.calculateHistoricalSimilarity(lang1, lang2)
	geographicScore := mie.calculateGeographicSimilarity(lang1, lang2)
	contactScore := mie.calculateContactSimilarity(lang1, lang2)

	// Calculate weighted overall score
	overallScore := (phonologicalScore * mie.PhonologicalWeight) +
		(lexicalScore * mie.LexicalWeight) +
		(grammaticalScore * mie.GrammaticalWeight) +
		(historicalScore * mie.HistoricalWeight) +
		(geographicScore * mie.GeographicWeight) +
		(contactScore * mie.ContactWeight)

	// Ensure score is within bounds
	if overallScore < mie.MinIntelligibility {
		overallScore = mie.MinIntelligibility
	}
	if overallScore > mie.MaxIntelligibility {
		overallScore = mie.MaxIntelligibility
	}

	// Determine intelligibility level
	intelligibilityLevel := mie.determineIntelligibilityLevel(overallScore)

	// Generate description and factors
	description, factors := mie.generateIntelligibilityDescription(
		lang1, lang2, overallScore, intelligibilityLevel,
		phonologicalScore, lexicalScore, grammaticalScore,
		historicalScore, geographicScore, contactScore,
	)

	// Calculate confidence based on data availability
	confidence := mie.calculateConfidence(lang1, lang2)

	return &IntelligibilityResult{
		SourceLanguage:       lang1.Name,
		TargetLanguage:       lang2.Name,
		OverallScore:         overallScore,
		PhonologicalScore:    phonologicalScore,
		LexicalScore:         lexicalScore,
		GrammaticalScore:     grammaticalScore,
		HistoricalScore:      historicalScore,
		GeographicScore:      geographicScore,
		ContactScore:         contactScore,
		IntelligibilityLevel: string(intelligibilityLevel),
		Description:          description,
		Factors:              factors,
		Confidence:           confidence,
	}, nil
}

// calculatePhonologicalSimilarity calculates the phonological similarity between two languages.
func (mie *MutualIntelligibilityEngine) calculatePhonologicalSimilarity(lang1, lang2 *Language) float32 {
	if lang1.Phonology == nil || lang2.Phonology == nil {
		return 0.5 // Default similarity if phonology data is missing
	}

	// For now, use a simplified approach
	// In a full implementation, this would compare actual phoneme inventories and phonotactic rules

	// Check if they share the same phonology system
	if lang1.Phonology == lang2.Phonology {
		return 1.0
	}

	// Check if they have similar phoneme pools
	// This is a placeholder - in reality, we'd compare actual phoneme inventories
	return 0.7 // Assume moderate similarity for different phonology systems
}

// calculateLexicalSimilarity calculates the lexical similarity between two languages.
func (mie *MutualIntelligibilityEngine) calculateLexicalSimilarity(lang1, lang2 *Language) float32 {
	// Check for shared vocabulary through cultural influence
	sharedVocabulary := 0
	totalVocabulary := 0

	// Count shared concepts in linguistic changes
	for _, change := range lang1.LinguisticChanges {
		if change.Type == LinguisticChangeTypeCultural {
			totalVocabulary++
			// Check if this change is also present in lang2
			for _, change2 := range lang2.LinguisticChanges {
				if change2.Type == LinguisticChangeTypeCultural && change.Trigger == change2.Trigger {
					sharedVocabulary++
					break
				}
			}
		}
	}

	if totalVocabulary == 0 {
		return 0.5 // Default similarity if no vocabulary data
	}

	return float32(sharedVocabulary) / float32(totalVocabulary)
}

// calculateGrammaticalSimilarity calculates the grammatical similarity between two languages.
func (mie *MutualIntelligibilityEngine) calculateGrammaticalSimilarity(lang1, lang2 *Language) float32 {
	if lang1.Grammar == nil || lang2.Grammar == nil {
		return 0.5 // Default similarity if grammar data is missing
	}

	// Check if they share the same grammar system
	if lang1.Grammar == lang2.Grammar {
		return 1.0
	}

	// For now, use a simplified approach
	// In a full implementation, this would compare actual grammatical features
	return 0.6 // Assume moderate similarity for different grammar systems
}

// calculateHistoricalSimilarity calculates the historical relationship between two languages.
func (mie *MutualIntelligibilityEngine) calculateHistoricalSimilarity(lang1, lang2 *Language) float32 {
	// Check if they share a common ancestor
	if lang1.ParentID != nil && lang2.ParentID != nil {
		if *lang1.ParentID == *lang2.ParentID {
			return 0.9 // Siblings have high historical similarity
		}
	}

	// Check if one is a descendant of the other
	if lang1.ParentID != nil && *lang1.ParentID == lang2.ID {
		return 0.8 // Direct parent-child relationship
	}
	if lang2.ParentID != nil && *lang2.ParentID == lang1.ID {
		return 0.8 // Direct parent-child relationship
	}

	// Check for shared ancestors further back
	// This is a simplified approach - in reality, we'd traverse the full family tree
	return 0.3 // Assume low similarity for unrelated languages
}

// calculateGeographicSimilarity calculates the geographic proximity between two languages.
func (mie *MutualIntelligibilityEngine) calculateGeographicSimilarity(lang1, lang2 *Language) float32 {
	// For now, use a simplified approach
	// In a full implementation, this would use actual geographic coordinates

	// Check if they have similar geographic regions in their linguistic changes
	sharedGeographicFeatures := 0
	totalGeographicFeatures := 0

	for _, change := range lang1.LinguisticChanges {
		if change.Type == LinguisticChangeTypeDialectal && strings.Contains(change.Trigger, "geographic") {
			totalGeographicFeatures++
			// Check if lang2 has similar geographic features
			for _, change2 := range lang2.LinguisticChanges {
				if change2.Type == LinguisticChangeTypeDialectal && strings.Contains(change2.Trigger, "geographic") {
					sharedGeographicFeatures++
					break
				}
			}
		}
	}

	if totalGeographicFeatures == 0 {
		return 0.5 // Default similarity if no geographic data
	}

	return float32(sharedGeographicFeatures) / float32(totalGeographicFeatures)
}

// calculateContactSimilarity calculates the contact history between two languages.
func (mie *MutualIntelligibilityEngine) calculateContactSimilarity(lang1, lang2 *Language) float32 {
	// Check for shared contact events in linguistic changes
	sharedContacts := 0
	totalContacts := 0

	for _, change := range lang1.LinguisticChanges {
		if change.Type == LinguisticChangeTypeCultural && strings.Contains(change.Trigger, "contact") {
			totalContacts++
			// Check if lang2 has similar contact events
			for _, change2 := range lang2.LinguisticChanges {
				if change2.Type == LinguisticChangeTypeCultural && strings.Contains(change2.Trigger, "contact") {
					sharedContacts++
					break
				}
			}
		}
	}

	if totalContacts == 0 {
		return 0.5 // Default similarity if no contact data
	}

	return float32(sharedContacts) / float32(totalContacts)
}

// determineIntelligibilityLevel determines the intelligibility level based on the overall score.
func (mie *MutualIntelligibilityEngine) determineIntelligibilityLevel(score float32) IntelligibilityLevel {
	switch {
	case score >= 0.95:
		return IntelligibilityLevelNearNative
	case score >= 0.80:
		return IntelligibilityLevelVeryHigh
	case score >= 0.60:
		return IntelligibilityLevelHigh
	case score >= 0.40:
		return IntelligibilityLevelModerate
	case score >= 0.25:
		return IntelligibilityLevelLow
	case score >= 0.10:
		return IntelligibilityLevelMinimal
	default:
		return IntelligibilityLevelNone
	}
}

// generateIntelligibilityDescription generates a human-readable description of the intelligibility result.
func (mie *MutualIntelligibilityEngine) generateIntelligibilityDescription(
	lang1, lang2 *Language,
	overallScore float32,
	intelligibilityLevel IntelligibilityLevel,
	phonologicalScore, lexicalScore, grammaticalScore, historicalScore, geographicScore, contactScore float32,
) (string, []string) {
	var factors []string
	var description strings.Builder

	// Base description
	description.WriteString(fmt.Sprintf("Speakers of %s and %s have %s mutual intelligibility (%.1f%%). ",
		lang1.Name, lang2.Name, intelligibilityLevel, overallScore*100))

	// Add specific factors
	if phonologicalScore > 0.7 {
		factors = append(factors, "similar_phonology")
		description.WriteString("Phonological systems are similar. ")
	}
	if lexicalScore > 0.7 {
		factors = append(factors, "shared_vocabulary")
		description.WriteString("Vocabulary shows significant overlap. ")
	}
	if grammaticalScore > 0.7 {
		factors = append(factors, "similar_grammar")
		description.WriteString("Grammatical structures are comparable. ")
	}
	if historicalScore > 0.7 {
		factors = append(factors, "historical_relationship")
		description.WriteString("Languages share a common ancestor. ")
	}
	if geographicScore > 0.7 {
		factors = append(factors, "geographic_proximity")
		description.WriteString("Languages developed in similar geographic regions. ")
	}
	if contactScore > 0.7 {
		factors = append(factors, "contact_history")
		description.WriteString("Languages have a history of contact and borrowing. ")
	}

	// Add intelligibility level description
	switch intelligibilityLevel {
	case IntelligibilityLevelNearNative:
		description.WriteString("Speakers can understand each other with near-native fluency.")
	case IntelligibilityLevelVeryHigh:
		description.WriteString("Speakers can communicate with very high understanding.")
	case IntelligibilityLevelHigh:
		description.WriteString("Speakers can communicate with high understanding.")
	case IntelligibilityLevelModerate:
		description.WriteString("Speakers can communicate with moderate understanding.")
	case IntelligibilityLevelLow:
		description.WriteString("Speakers have difficulty understanding each other.")
	case IntelligibilityLevelMinimal:
		description.WriteString("Speakers have minimal understanding of each other.")
	case IntelligibilityLevelNone:
		description.WriteString("Speakers cannot understand each other.")
	}

	return description.String(), factors
}

// calculateConfidence calculates the confidence level of the intelligibility calculation.
func (mie *MutualIntelligibilityEngine) calculateConfidence(lang1, lang2 *Language) float32 {
	confidence := float32(1.0)

	// Reduce confidence if linguistic components are missing
	if lang1.Phonology == nil || lang2.Phonology == nil {
		confidence -= 0.1
	}
	if lang1.Grammar == nil || lang2.Grammar == nil {
		confidence -= 0.1
	}
	if len(lang1.LinguisticChanges) == 0 || len(lang2.LinguisticChanges) == 0 {
		confidence -= 0.2
	}

	// Ensure confidence stays in valid range
	if confidence < 0.0 {
		confidence = 0.0
	}
	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

// BatchCalculateIntelligibility calculates mutual intelligibility for multiple language pairs.
func (mie *MutualIntelligibilityEngine) BatchCalculateIntelligibility(
	languages []*Language,
) ([]*IntelligibilityResult, error) {
	var results []*IntelligibilityResult

	for i := 0; i < len(languages); i++ {
		for j := i + 1; j < len(languages); j++ {
			result, err := mie.CalculateMutualIntelligibility(languages[i], languages[j])
			if err != nil {
				return nil, fmt.Errorf("failed to calculate intelligibility for %s and %s: %w",
					languages[i].Name, languages[j].Name, err)
			}
			results = append(results, result)
		}
	}

	return results, nil
}

// FindMostIntelligiblePairs finds the most intelligible language pairs.
func (mie *MutualIntelligibilityEngine) FindMostIntelligiblePairs(
	results []*IntelligibilityResult,
	limit int,
) []*IntelligibilityResult {
	if limit <= 0 {
		limit = len(results)
	}

	// Sort by overall score (highest first)
	sorted := make([]*IntelligibilityResult, len(results))
	copy(sorted, results)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].OverallScore > sorted[j].OverallScore
	})

	if len(sorted) <= limit {
		return sorted
	}
	return sorted[:limit]
}

// FindLeastIntelligiblePairs finds the least intelligible language pairs.
func (mie *MutualIntelligibilityEngine) FindLeastIntelligiblePairs(
	results []*IntelligibilityResult,
	limit int,
) []*IntelligibilityResult {
	if limit <= 0 {
		limit = len(results)
	}

	// Sort by overall score (lowest first)
	sorted := make([]*IntelligibilityResult, len(results))
	copy(sorted, results)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].OverallScore < sorted[j].OverallScore
	})

	if len(sorted) <= limit {
		return sorted
	}
	return sorted[:limit]
}

// HistoricalReconstructionEngine reconstructs ancestral languages from their descendants.
type HistoricalReconstructionEngine struct {
	// Configuration parameters
	ReconstructionDepth     int     `json:"reconstructionDepth"`     // How many generations back to reconstruct
	ConfidenceThreshold     float32 `json:"confidenceThreshold"`     // Minimum confidence for reconstruction
	FeaturePreservationRate float32 `json:"featurePreservationRate"` // Rate at which features are preserved (0.0-1.0)
	ChangeReversalRate      float32 `json:"changeReversalRate"`      // Rate at which changes are reversed (0.0-1.0)
	MaxReconstructionSteps  int     `json:"maxReconstructionSteps"`  // Maximum reconstruction steps per feature
}

// NewHistoricalReconstructionEngine creates a new historical reconstruction engine with default settings.
func NewHistoricalReconstructionEngine() *HistoricalReconstructionEngine {
	return &HistoricalReconstructionEngine{
		ReconstructionDepth:     3,   // Reconstruct up to 3 generations back
		ConfidenceThreshold:     0.6, // 60% minimum confidence
		FeaturePreservationRate: 0.8, // 80% feature preservation rate
		ChangeReversalRate:      0.7, // 70% change reversal rate
		MaxReconstructionSteps:  10,  // Maximum 10 reconstruction steps per feature
	}
}

// ReconstructionResult represents the result of a historical language reconstruction.
type ReconstructionResult struct {
	AncestralLanguage     *Language   `json:"ancestralLanguage"`
	ReconstructionPath    []*Language `json:"reconstructionPath"` // Path from descendant to ancestor
	Confidence            float32     `json:"confidence"`
	ReconstructionSteps   int         `json:"reconstructionSteps"`
	PreservedFeatures     []string    `json:"preservedFeatures"`
	ReconstructedFeatures []string    `json:"reconstructedFeatures"`
	LostFeatures          []string    `json:"lostFeatures"`
	Description           string      `json:"description"`
	Notes                 []string    `json:"notes"`
}

// ReconstructAncestralLanguage reconstructs an ancestral language from a descendant.
func (hre *HistoricalReconstructionEngine) ReconstructAncestralLanguage(
	descendant *Language,
	targetGeneration int,
	seed int64,
) (*ReconstructionResult, error) {
	if descendant == nil {
		return nil, fmt.Errorf("descendant language cannot be nil")
	}
	if targetGeneration < 0 {
		return nil, fmt.Errorf("target generation cannot be negative")
	}
	if targetGeneration > hre.ReconstructionDepth {
		return nil, fmt.Errorf("target generation %d exceeds reconstruction depth %d", targetGeneration, hre.ReconstructionDepth)
	}

	// Start reconstruction from the descendant
	currentLang := descendant
	reconstructionPath := []*Language{currentLang}
	preservedFeatures := make([]string, 0)
	reconstructedFeatures := make([]string, 0)
	lostFeatures := make([]string, 0)
	totalSteps := 0

	// Reconstruct each generation
	for generation := 0; generation < targetGeneration; generation++ {
		// Find the parent language
		if currentLang.ParentID == nil {
			// No parent available, create a reconstructed ancestor
			reconstructedAncestor, err := hre.reconstructFromDescendant(currentLang, generation+1, seed)
			if err != nil {
				return nil, fmt.Errorf("failed to reconstruct ancestor at generation %d: %w", generation+1, err)
			}
			currentLang = reconstructedAncestor
		} else {
			// Parent exists, use it
			parentLang, err := hre.findLanguageByID(*currentLang.ParentID)
			if err != nil {
				// Parent not found, create reconstructed ancestor
				reconstructedAncestor, err := hre.reconstructFromDescendant(currentLang, generation+1, seed)
				if err != nil {
					return nil, fmt.Errorf("failed to reconstruct ancestor at generation %d: %w", generation+1, err)
				}
				currentLang = reconstructedAncestor
			} else {
				currentLang = parentLang
			}
		}

		reconstructionPath = append(reconstructionPath, currentLang)
		totalSteps++

		// Analyze what features were preserved, reconstructed, or lost
		hre.analyzeFeatureChanges(reconstructionPath, &preservedFeatures, &reconstructedFeatures, &lostFeatures)
	}

	// Calculate overall confidence
	confidence := hre.calculateReconstructionConfidence(reconstructionPath, preservedFeatures, reconstructedFeatures, lostFeatures)

	// Generate description
	description := hre.generateReconstructionDescription(descendant, currentLang, targetGeneration, confidence)

	// Create result
	result := &ReconstructionResult{
		AncestralLanguage:     currentLang,
		ReconstructionPath:    reconstructionPath,
		Confidence:            confidence,
		ReconstructionSteps:   totalSteps,
		PreservedFeatures:     preservedFeatures,
		ReconstructedFeatures: reconstructedFeatures,
		LostFeatures:          lostFeatures,
		Description:           description,
		Notes:                 hre.generateReconstructionNotes(reconstructionPath),
	}

	return result, nil
}

// reconstructFromDescendant reconstructs an ancestral language from a descendant's features.
func (hre *HistoricalReconstructionEngine) reconstructFromDescendant(
	descendant *Language,
	generation int,
	seed int64,
) (*Language, error) {
	// Create a new language ID for the reconstructed ancestor
	ancestorID := LanguageID{
		Family:   descendant.ID.Family,
		Branch:   descendant.ID.Branch,
		Language: descendant.ID.Language,
		Dialect:  fmt.Sprintf("reconstructed_gen_%d", generation),
	}

	// Clone the descendant and reverse changes to reconstruct the ancestor
	ancestor := descendant.Clone(ancestorID, seed)
	ancestor.Name = fmt.Sprintf("%s (Reconstructed Ancestor Gen %d)", descendant.Name, generation)
	ancestor.SetDescription(fmt.Sprintf("Reconstructed ancestral language from generation %d", generation))

	// Reverse linguistic changes to reconstruct ancestral features
	hre.reverseLinguisticChanges(ancestor, generation)

	return ancestor, nil
}

// reverseLinguisticChanges reverses linguistic changes to reconstruct ancestral features.
func (hre *HistoricalReconstructionEngine) reverseLinguisticChanges(
	language *Language,
	generation int,
) {
	// Reverse changes in reverse chronological order
	reversedChanges := make([]LinguisticChange, 0)

	for i := len(language.LinguisticChanges) - 1; i >= 0; i-- {
		change := language.LinguisticChanges[i]

		// Determine if this change should be reversed based on generation and confidence
		if hre.shouldReverseChange(change, generation) {
			reversedChange := hre.reverseChange(change)
			reversedChanges = append(reversedChanges, reversedChange)
		}
	}

	// Apply reversed changes
	for _, change := range reversedChanges {
		language.AddLinguisticChange(change)
	}
}

// shouldReverseChange determines if a linguistic change should be reversed.
func (hre *HistoricalReconstructionEngine) shouldReverseChange(
	change LinguisticChange,
	generation int,
) bool {
	// Higher generations (more distant ancestors) have higher reversal rates
	generationMultiplier := float32(generation+1) / float32(hre.ReconstructionDepth+1)
	effectiveReversalRate := hre.ChangeReversalRate * generationMultiplier

	// Use the change's intensity to determine reversal probability
	reversalProbability := effectiveReversalRate * change.Intensity

	// Simulate probabilistic reversal
	return rand.Float32() < reversalProbability
}

// reverseChange creates a reversed version of a linguistic change.
func (hre *HistoricalReconstructionEngine) reverseChange(change LinguisticChange) LinguisticChange {
	reversedChange := LinguisticChange{
		ID:          fmt.Sprintf("reversed_%s", change.ID),
		Type:        change.Type,
		Description: fmt.Sprintf("Reversed: %s", change.Description),
		Details:     fmt.Sprintf("Reconstruction reversal of: %s", change.Details),
		Timestamp:   time.Now(),
		Era:         "reconstruction",
		Trigger:     "historical_reconstruction",
		Intensity:   change.Intensity * hre.FeaturePreservationRate,
	}

	return reversedChange
}

// findLanguageByID finds a language by its ID.
// This is a simplified implementation - in a full system, this would search through a language database.
func (hre *HistoricalReconstructionEngine) findLanguageByID(id LanguageID) (*Language, error) {
	// For now, return an error indicating the language wasn't found
	// In a full implementation, this would search through available languages
	return nil, fmt.Errorf("language with ID %s not found in current context", id.String())
}

// analyzeFeatureChanges analyzes what features were preserved, reconstructed, or lost during reconstruction.
func (hre *HistoricalReconstructionEngine) analyzeFeatureChanges(
	reconstructionPath []*Language,
	preservedFeatures, reconstructedFeatures, lostFeatures *[]string,
) {
	if len(reconstructionPath) < 2 {
		return
	}

	descendant := reconstructionPath[0]
	ancestor := reconstructionPath[len(reconstructionPath)-1]

	// Analyze phonological features
	if descendant.Phonology != nil && ancestor.Phonology != nil {
		if descendant.Phonology == ancestor.Phonology {
			*preservedFeatures = append(*preservedFeatures, "phonology")
		} else {
			*reconstructedFeatures = append(*reconstructedFeatures, "phonology")
		}
	} else if descendant.Phonology != nil {
		*lostFeatures = append(*lostFeatures, "phonology")
	}

	// Analyze grammatical features
	if descendant.Grammar != nil && ancestor.Grammar != nil {
		if descendant.Grammar == ancestor.Grammar {
			*preservedFeatures = append(*preservedFeatures, "grammar")
		} else {
			*reconstructedFeatures = append(*reconstructedFeatures, "grammar")
		}
	} else if descendant.Grammar != nil {
		*lostFeatures = append(*lostFeatures, "grammar")
	}

	// Analyze morphological features
	if descendant.Morphology != nil && ancestor.Morphology != nil {
		if descendant.Morphology == ancestor.Morphology {
			*preservedFeatures = append(*preservedFeatures, "morphology")
		} else {
			*reconstructedFeatures = append(*reconstructedFeatures, "morphology")
		}
	} else if descendant.Morphology != nil {
		*lostFeatures = append(*lostFeatures, "morphology")
	}

	// Analyze orthographic features
	if descendant.Orthography != nil && ancestor.Orthography != nil {
		if descendant.Orthography == ancestor.Orthography {
			*preservedFeatures = append(*preservedFeatures, "orthography")
		} else {
			*reconstructedFeatures = append(*reconstructedFeatures, "orthography")
		}
	} else if descendant.Orthography != nil {
		*lostFeatures = append(*lostFeatures, "orthography")
	}
}

// calculateReconstructionConfidence calculates the confidence level of the reconstruction.
func (hre *HistoricalReconstructionEngine) calculateReconstructionConfidence(
	reconstructionPath []*Language,
	preservedFeatures, reconstructedFeatures, lostFeatures []string,
) float32 {
	confidence := float32(1.0)

	// Reduce confidence for each reconstruction step
	stepPenalty := float32(len(reconstructionPath)-1) * 0.1
	confidence -= stepPenalty

	// Boost confidence for preserved features
	preservedBonus := float32(len(preservedFeatures)) * 0.05
	confidence += preservedBonus

	// Reduce confidence for lost features
	lostPenalty := float32(len(lostFeatures)) * 0.1
	confidence -= lostPenalty

	// Ensure confidence stays in valid range
	if confidence < 0.0 {
		confidence = 0.0
	}
	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

// generateReconstructionDescription generates a human-readable description of the reconstruction.
func (hre *HistoricalReconstructionEngine) generateReconstructionDescription(
	descendant, ancestor *Language,
	generation int,
	confidence float32,
) string {
	confidencePercent := int(confidence * 100)

	description := fmt.Sprintf(
		"Reconstructed ancestral language '%s' from descendant '%s' (%d generations back). "+
			"Reconstruction confidence: %d%%. ",
		ancestor.Name, descendant.Name, generation, confidencePercent,
	)

	if confidence >= 0.8 {
		description += "High confidence reconstruction with strong evidence."
	} else if confidence >= 0.6 {
		description += "Moderate confidence reconstruction with reasonable evidence."
	} else {
		description += "Low confidence reconstruction with limited evidence."
	}

	return description
}

// generateReconstructionNotes generates notes about the reconstruction process.
func (hre *HistoricalReconstructionEngine) generateReconstructionNotes(reconstructionPath []*Language) []string {
	var notes []string

	if len(reconstructionPath) > 1 {
		notes = append(notes, fmt.Sprintf("Reconstruction path: %d steps", len(reconstructionPath)-1))
	}

	// Add notes about reconstruction quality
	if len(reconstructionPath) > 2 {
		notes = append(notes, "Multiple generation reconstruction may reduce accuracy")
	}

	// Add notes about available data
	notes = append(notes, "Reconstruction based on linguistic change analysis")
	notes = append(notes, "Confidence varies based on feature preservation and change reversal")

	return notes
}

// BatchReconstructAncestors reconstructs ancestral languages for multiple descendants.
func (hre *HistoricalReconstructionEngine) BatchReconstructAncestors(
	descendants []*Language,
	targetGeneration int,
	seed int64,
) ([]*ReconstructionResult, error) {
	var results []*ReconstructionResult

	for _, descendant := range descendants {
		result, err := hre.ReconstructAncestralLanguage(descendant, targetGeneration, seed)
		if err != nil {
			return nil, fmt.Errorf("failed to reconstruct ancestor for %s: %w", descendant.Name, err)
		}
		results = append(results, result)
	}

	return results, nil
}

// FindCommonAncestor finds the common ancestor of multiple languages.
func (hre *HistoricalReconstructionEngine) FindCommonAncestor(
	languages []*Language,
	maxGenerations int,
	seed int64,
) (*ReconstructionResult, error) {
	if len(languages) < 2 {
		return nil, fmt.Errorf("at least 2 languages required to find common ancestor")
	}

	// Start with the first language and find its ancestors
	baseLanguage := languages[0]
	baseAncestors := make([]*Language, 0)

	for generation := 1; generation <= maxGenerations; generation++ {
		result, err := hre.ReconstructAncestralLanguage(baseLanguage, generation, seed)
		if err != nil {
			break
		}
		baseAncestors = append(baseAncestors, result.AncestralLanguage)
	}

	// For each ancestor, check if it's a common ancestor of all languages
	for i, ancestor := range baseAncestors {
		generation := i + 1
		isCommonAncestor := true

		for _, language := range languages[1:] {
			// Check if this language can trace back to the ancestor
			if !hre.canTraceToAncestor(language, ancestor, generation) {
				isCommonAncestor = false
				break
			}
		}

		if isCommonAncestor {
			// Found common ancestor
			return &ReconstructionResult{
				AncestralLanguage:     ancestor,
				ReconstructionPath:    []*Language{baseLanguage, ancestor},
				Confidence:            0.8, // High confidence for common ancestor
				ReconstructionSteps:   generation,
				PreservedFeatures:     []string{"common_ancestry"},
				ReconstructedFeatures: []string{"shared_features"},
				LostFeatures:          []string{},
				Description:           fmt.Sprintf("Common ancestor of %d languages at generation %d", len(languages), generation),
				Notes:                 []string{"Common ancestor identified through shared ancestry"},
			}, nil
		}
	}

	return nil, fmt.Errorf("no common ancestor found within %d generations", maxGenerations)
}

// canTraceToAncestor checks if a language can trace back to a specific ancestor.
func (hre *HistoricalReconstructionEngine) canTraceToAncestor(
	language, ancestor *Language,
	generations int,
) bool {
	// Simplified implementation - check if they share the same family and branch
	// In a full implementation, this would traverse the actual family tree
	return language.ID.Family == ancestor.ID.Family && language.ID.Branch == ancestor.ID.Branch
}

// SetRandomSeed sets the random seed for deterministic reconstruction.
func (hre *HistoricalReconstructionEngine) SetRandomSeed(seed int64) {
	rand.Seed(seed)
}

// WritingSystemGenealogy tracks the evolution and relationships between writing systems.
type WritingSystemGenealogy struct {
	// Core genealogy information
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Family         string  `json:"family"`
	ScriptType     string  `json:"scriptType"`
	OriginDate     string  `json:"originDate"`
	ExtinctionDate *string `json:"extinctionDate,omitempty"`

	// Genealogical relationships
	ParentID     *string  `json:"parentID,omitempty"`
	ChildIDs     []string `json:"childIDs"`
	InfluencedBy []string `json:"influencedBy"`
	Influenced   []string `json:"influenced"`

	// Writing system characteristics
	Direction     string `json:"direction"`     // left-to-right, right-to-left, top-to-bottom, etc.
	AlphabetSize  int    `json:"alphabetSize"`  // Number of basic characters
	HasVowels     bool   `json:"hasVowels"`     // Whether vowels are explicitly written
	HasDiacritics bool   `json:"hasDiacritics"` // Whether diacritical marks are used
	IsLogographic bool   `json:"isLogographic"` // Whether it's primarily logographic
	IsSyllabic    bool   `json:"isSyllabic"`    // Whether it's primarily syllabic
	IsAlphabetic  bool   `json:"isAlphabetic"`  // Whether it's primarily alphabetic

	// Evolution tracking
	EvolutionSteps []WritingSystemChange `json:"evolutionSteps"`
	Complexity     float32               `json:"complexity"`
	Elegance       float32               `json:"elegance"`

	// Geographic and cultural context
	GeographicOrigin string   `json:"geographicOrigin"`
	CulturalOrigin   string   `json:"culturalOrigin"`
	UsageRegions     []string `json:"usageRegions"`

	// Metadata
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Notes     []string  `json:"notes"`
}

// WritingSystemChange represents a change in a writing system over time.
type WritingSystemChange struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // addition, removal, modification, borrowing
	Description string    `json:"description"`
	Details     string    `json:"details"`
	Timestamp   time.Time `json:"timestamp"`
	Era         string    `json:"era"`
	Trigger     string    `json:"trigger"`    // cultural_contact, technological_advance, etc.
	Intensity   float32   `json:"intensity"`  // 0.0-1.0
	Confidence  float32   `json:"confidence"` // 0.0-1.0
}

// NewWritingSystemGenealogy creates a new writing system genealogy entry.
func NewWritingSystemGenealogy(
	id, name, description, family, scriptType, originDate string,
) *WritingSystemGenealogy {
	now := time.Now()
	return &WritingSystemGenealogy{
		ID:             id,
		Name:           name,
		Description:    description,
		Family:         family,
		ScriptType:     scriptType,
		OriginDate:     originDate,
		ChildIDs:       make([]string, 0),
		InfluencedBy:   make([]string, 0),
		Influenced:     make([]string, 0),
		Direction:      "left-to-right",
		AlphabetSize:   0,
		HasVowels:      false,
		HasDiacritics:  false,
		IsLogographic:  false,
		IsSyllabic:     false,
		IsAlphabetic:   false,
		EvolutionSteps: make([]WritingSystemChange, 0),
		Complexity:     0.5,
		Elegance:       0.5,
		UsageRegions:   make([]string, 0),
		CreatedAt:      now,
		UpdatedAt:      now,
		Notes:          make([]string, 0),
	}
}

// AddChild adds a child writing system to this genealogy.
func (wsg *WritingSystemGenealogy) AddChild(childID string) {
	if !wsg.hasChild(childID) {
		wsg.ChildIDs = append(wsg.ChildIDs, childID)
		wsg.UpdatedAt = time.Now()
	}
}

// RemoveChild removes a child writing system from this genealogy.
func (wsg *WritingSystemGenealogy) RemoveChild(childID string) {
	for i, id := range wsg.ChildIDs {
		if id == childID {
			wsg.ChildIDs = append(wsg.ChildIDs[:i], wsg.ChildIDs[i+1:]...)
			wsg.UpdatedAt = time.Now()
			break
		}
	}
}

// hasChild checks if a writing system has a specific child.
func (wsg *WritingSystemGenealogy) hasChild(childID string) bool {
	for _, id := range wsg.ChildIDs {
		if id == childID {
			return true
		}
	}
	return false
}

// AddInfluence adds an influence relationship to this genealogy.
func (wsg *WritingSystemGenealogy) AddInfluence(influencedID string) {
	if !wsg.hasInfluence(influencedID) {
		wsg.Influenced = append(wsg.Influenced, influencedID)
		wsg.UpdatedAt = time.Now()
	}
}

// AddInfluencedBy adds an "influenced by" relationship to this genealogy.
func (wsg *WritingSystemGenealogy) AddInfluencedBy(influencerID string) {
	if !wsg.hasInfluencedBy(influencerID) {
		wsg.InfluencedBy = append(wsg.InfluencedBy, influencerID)
		wsg.UpdatedAt = time.Now()
	}
}

// hasInfluence checks if a writing system has influenced a specific system.
func (wsg *WritingSystemGenealogy) hasInfluence(influencedID string) bool {
	for _, id := range wsg.Influenced {
		if id == influencedID {
			return true
		}
	}
	return false
}

// hasInfluencedBy checks if a writing system has been influenced by a specific system.
func (wsg *WritingSystemGenealogy) hasInfluencedBy(influencerID string) bool {
	for _, id := range wsg.InfluencedBy {
		if id == influencerID {
			return true
		}
	}
	return false
}

// AddEvolutionStep adds an evolution step to the writing system.
func (wsg *WritingSystemGenealogy) AddEvolutionStep(change WritingSystemChange) {
	wsg.EvolutionSteps = append(wsg.EvolutionSteps, change)
	wsg.UpdatedAt = time.Now()

	// Update complexity and elegance based on the change
	wsg.updateMetrics(change)
}

// updateMetrics updates the complexity and elegance metrics based on a change.
func (wsg *WritingSystemGenealogy) updateMetrics(change WritingSystemChange) {
	switch change.Type {
	case "addition":
		wsg.Complexity = min(1.0, wsg.Complexity+change.Intensity*0.1)
		wsg.Elegance = min(1.0, wsg.Elegance+change.Intensity*0.05)
	case "removal":
		wsg.Complexity = max(0.0, wsg.Complexity-change.Intensity*0.1)
		wsg.Elegance = max(0.0, wsg.Elegance-change.Intensity*0.05)
	case "modification":
		wsg.Complexity = min(1.0, wsg.Complexity+change.Intensity*0.05)
		wsg.Elegance = min(1.0, wsg.Elegance+change.Intensity*0.1)
	case "borrowing":
		wsg.Complexity = min(1.0, wsg.Complexity+change.Intensity*0.08)
		wsg.Elegance = min(1.0, wsg.Elegance+change.Intensity*0.08)
	}
}

// min returns the minimum of two float32 values.
func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two float32 values.
func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

// SetParent sets the parent writing system for this genealogy.
func (wsg *WritingSystemGenealogy) SetParent(parentID string) {
	wsg.ParentID = &parentID
	wsg.UpdatedAt = time.Now()
}

// GetAncestors returns all ancestor writing systems in the genealogy tree.
func (wsg *WritingSystemGenealogy) GetAncestors(genealogies map[string]*WritingSystemGenealogy) []*WritingSystemGenealogy {
	var ancestors []*WritingSystemGenealogy
	current := wsg

	for current.ParentID != nil {
		parent, exists := genealogies[*current.ParentID]
		if !exists {
			break
		}
		ancestors = append(ancestors, parent)
		current = parent
	}

	return ancestors
}

// GetDescendants returns all descendant writing systems in the genealogy tree.
func (wsg *WritingSystemGenealogy) GetDescendants(genealogies map[string]*WritingSystemGenealogy) []*WritingSystemGenealogy {
	var descendants []*WritingSystemGenealogy
	visited := make(map[string]bool)

	var collectDescendants func(ws *WritingSystemGenealogy)
	collectDescendants = func(ws *WritingSystemGenealogy) {
		if visited[ws.ID] {
			return
		}
		visited[ws.ID] = true

		for _, childID := range ws.ChildIDs {
			if child, exists := genealogies[childID]; exists {
				descendants = append(descendants, child)
				collectDescendants(child)
			}
		}
	}

	collectDescendants(wsg)
	return descendants
}

// GetSiblings returns all sibling writing systems (those with the same parent).
func (wsg *WritingSystemGenealogy) GetSiblings(genealogies map[string]*WritingSystemGenealogy) []*WritingSystemGenealogy {
	if wsg.ParentID == nil {
		return nil
	}

	var siblings []*WritingSystemGenealogy
	parent, exists := genealogies[*wsg.ParentID]
	if !exists {
		return nil
	}

	for _, childID := range parent.ChildIDs {
		if childID != wsg.ID {
			if child, exists := genealogies[childID]; exists {
				siblings = append(siblings, child)
			}
		}
	}

	return siblings
}

// CalculateSimilarity calculates the similarity between two writing systems.
func (wsg *WritingSystemGenealogy) CalculateSimilarity(other *WritingSystemGenealogy) float32 {
	if wsg.ID == other.ID {
		return 1.0
	}

	// Calculate feature similarity
	featureSimilarity := wsg.calculateFeatureSimilarity(other)

	// Calculate genealogical similarity
	genealogicalSimilarity := wsg.calculateGenealogicalSimilarity(other)

	// Weighted combination
	return (featureSimilarity * 0.6) + (genealogicalSimilarity * 0.4)
}

// calculateFeatureSimilarity calculates similarity based on writing system features.
func (wsg *WritingSystemGenealogy) calculateFeatureSimilarity(other *WritingSystemGenealogy) float32 {
	similarity := float32(0.0)

	// Script type similarity
	if wsg.ScriptType == other.ScriptType {
		similarity += 0.3
	}

	// Direction similarity
	if wsg.Direction == other.Direction {
		similarity += 0.2
	}

	// Vowel representation similarity
	if wsg.HasVowels == other.HasVowels {
		similarity += 0.15
	}

	// Diacritic similarity
	if wsg.HasDiacritics == other.HasDiacritics {
		similarity += 0.15
	}

	// Writing system type similarity
	if wsg.IsLogographic == other.IsLogographic {
		similarity += 0.1
	}
	if wsg.IsSyllabic == other.IsSyllabic {
		similarity += 0.1
	}
	if wsg.IsAlphabetic == other.IsAlphabetic {
		similarity += 0.1
	}

	return similarity
}

// calculateGenealogicalSimilarity calculates similarity based on genealogical relationships.
func (wsg *WritingSystemGenealogy) calculateGenealogicalSimilarity(other *WritingSystemGenealogy) float32 {
	// Check if they're the same family
	if wsg.Family == other.Family {
		return 0.8
	}

	// Check if they share a parent
	if wsg.ParentID != nil && other.ParentID != nil && *wsg.ParentID == *other.ParentID {
		return 0.9
	}

	// Check if one is the parent of the other
	if wsg.ParentID != nil && *wsg.ParentID == other.ID {
		return 0.7
	}
	if other.ParentID != nil && *other.ParentID == wsg.ID {
		return 0.7
	}

	// Check for influence relationships
	if wsg.hasInfluence(other.ID) || wsg.hasInfluencedBy(other.ID) {
		return 0.6
	}

	return 0.3 // Default similarity for unrelated systems
}

// WritingSystemGenealogyManager manages a collection of writing system genealogies.
type WritingSystemGenealogyManager struct {
	Genealogies map[string]*WritingSystemGenealogy `json:"genealogies"`
	NextID      int                                `json:"nextID"`
}

// NewWritingSystemGenealogyManager creates a new writing system genealogy manager.
func NewWritingSystemGenealogyManager() *WritingSystemGenealogyManager {
	return &WritingSystemGenealogyManager{
		Genealogies: make(map[string]*WritingSystemGenealogy),
		NextID:      1,
	}
}

// CreateWritingSystem creates a new writing system genealogy entry.
func (wsgm *WritingSystemGenealogyManager) CreateWritingSystem(
	name, description, family, scriptType, originDate string,
) *WritingSystemGenealogy {
	id := fmt.Sprintf("ws_%d", wsgm.NextID)
	wsgm.NextID++

	genealogy := NewWritingSystemGenealogy(id, name, description, family, scriptType, originDate)
	wsgm.Genealogies[id] = genealogy

	return genealogy
}

// GetWritingSystem retrieves a writing system by ID.
func (wsgm *WritingSystemGenealogyManager) GetWritingSystem(id string) (*WritingSystemGenealogy, bool) {
	genealogy, exists := wsgm.Genealogies[id]
	return genealogy, exists
}

// CreateDerivedWritingSystem creates a new writing system derived from an existing one.
func (wsgm *WritingSystemGenealogyManager) CreateDerivedWritingSystem(
	parentID, name, description, scriptType, originDate string,
) (*WritingSystemGenealogy, error) {
	parent, exists := wsgm.Genealogies[parentID]
	if !exists {
		return nil, fmt.Errorf("parent writing system %s not found", parentID)
	}

	// Create the derived writing system
	derived := wsgm.CreateWritingSystem(name, description, parent.Family, scriptType, originDate)

	// Set up the parent-child relationship
	derived.SetParent(parentID)
	parent.AddChild(derived.ID)

	// Inherit some characteristics from the parent
	derived.Direction = parent.Direction
	derived.HasVowels = parent.HasVowels
	derived.HasDiacritics = parent.HasDiacritics
	derived.IsLogographic = parent.IsLogographic
	derived.IsSyllabic = parent.IsSyllabic
	derived.IsAlphabetic = parent.IsAlphabetic

	// Add evolution step
	evolutionStep := WritingSystemChange{
		ID:          fmt.Sprintf("evolution_%s", derived.ID),
		Type:        "derivation",
		Description: fmt.Sprintf("Derived from %s", parent.Name),
		Details:     fmt.Sprintf("Created new writing system based on %s", parent.Name),
		Timestamp:   time.Now(),
		Era:         "creation",
		Trigger:     "language_evolution",
		Intensity:   0.8,
		Confidence:  0.9,
	}
	derived.AddEvolutionStep(evolutionStep)

	return derived, nil
}

// FindCommonAncestor finds the common ancestor of two writing systems.
func (wsgm *WritingSystemGenealogyManager) FindCommonAncestor(
	ws1ID, ws2ID string,
) (*WritingSystemGenealogy, error) {
	ws1, exists := wsgm.Genealogies[ws1ID]
	if !exists {
		return nil, fmt.Errorf("writing system %s not found", ws1ID)
	}

	ws2, exists := wsgm.Genealogies[ws2ID]
	if !exists {
		return nil, fmt.Errorf("writing system %s not found", ws2ID)
	}

	// Get ancestors of both writing systems
	ancestors1 := ws1.GetAncestors(wsgm.Genealogies)
	ancestors2 := ws2.GetAncestors(wsgm.Genealogies)

	// Find the first common ancestor
	for _, ancestor1 := range ancestors1 {
		for _, ancestor2 := range ancestors2 {
			if ancestor1.ID == ancestor2.ID {
				return ancestor1, nil
			}
		}
	}

	// Check if one is the ancestor of the other
	for _, ancestor := range ancestors1 {
		if ancestor.ID == ws2.ID {
			return ws2, nil
		}
	}

	for _, ancestor := range ancestors2 {
		if ancestor.ID == ws1.ID {
			return ws1, nil
		}
	}

	return nil, fmt.Errorf("no common ancestor found")
}

// GetWritingSystemTree returns a tree representation of writing system relationships.
func (wsgm *WritingSystemGenealogyManager) GetWritingSystemTree(rootID string) (map[string]interface{}, error) {
	root, exists := wsgm.Genealogies[rootID]
	if !exists {
		return nil, fmt.Errorf("writing system %s not found", rootID)
	}

	return wsgm.buildTree(root), nil
}

// buildTree recursively builds a tree representation of writing system relationships.
func (wsgm *WritingSystemGenealogyManager) buildTree(ws *WritingSystemGenealogy) map[string]interface{} {
	tree := map[string]interface{}{
		"id":          ws.ID,
		"name":        ws.Name,
		"description": ws.Description,
		"scriptType":  ws.ScriptType,
		"complexity":  ws.Complexity,
		"elegance":    ws.Elegance,
		"children":    make([]map[string]interface{}, 0),
	}

	// Add children
	for _, childID := range ws.ChildIDs {
		if child, exists := wsgm.Genealogies[childID]; exists {
			childTree := wsgm.buildTree(child)
			tree["children"] = append(tree["children"].([]map[string]interface{}), childTree)
		}
	}

	return tree
}

// CalculateGenealogicalDistance calculates the genealogical distance between two writing systems.
func (wsgm *WritingSystemGenealogyManager) CalculateGenealogicalDistance(
	ws1ID, ws2ID string,
) (int, error) {
	commonAncestor, err := wsgm.FindCommonAncestor(ws1ID, ws2ID)
	if err != nil {
		return -1, err
	}

	// Calculate distance from each writing system to the common ancestor
	distance1 := wsgm.calculateDistanceToAncestor(ws1ID, commonAncestor.ID)
	distance2 := wsgm.calculateDistanceToAncestor(ws2ID, commonAncestor.ID)

	return distance1 + distance2, nil
}

// calculateDistanceToAncestor calculates the distance from a writing system to an ancestor.
func (wsgm *WritingSystemGenealogyManager) calculateDistanceToAncestor(
	wsID, ancestorID string,
) int {
	distance := 0
	current, exists := wsgm.Genealogies[wsID]
	if !exists {
		return -1
	}

	// If the current writing system is the ancestor, distance is 0
	if current.ID == ancestorID {
		return 0
	}

	// Count steps up to the ancestor
	for current.ParentID != nil {
		if *current.ParentID == ancestorID {
			return distance + 1
		}
		distance++
		current, exists = wsgm.Genealogies[*current.ParentID]
		if !exists {
			return -1
		}
	}

	return -1 // Ancestor not found
}

// SemanticField represents a conceptual category of related words and concepts.
type SemanticField struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"` // basic, cultural, technical, abstract, etc.
	Domain      string `json:"domain"`   // family, nature, technology, religion, etc.
	Era         string `json:"era"`      // ancient, medieval, modern, etc.

	// Core concepts and vocabulary
	CoreConcepts []string `json:"coreConcepts"` // Fundamental concepts in this field
	Vocabulary   []string `json:"vocabulary"`   // Words associated with this field
	Synonyms     []string `json:"synonyms"`     // Alternative terms for concepts

	// Evolution tracking
	EvolutionHistory []SemanticFieldChange `json:"evolutionHistory"`
	Complexity       float32               `json:"complexity"` // 0.0-1.0
	Richness         float32               `json:"richness"`   // 0.0-1.0
	Stability        float32               `json:"stability"`  // 0.0-1.0

	// Relationships
	RelatedFields []string `json:"relatedFields"`         // IDs of semantically related fields
	ParentField   *string  `json:"parentField,omitempty"` // ID of broader parent field
	ChildFields   []string `json:"childFields"`           // IDs of more specific sub-fields

	// Cultural and historical context
	CulturalOrigin   string   `json:"culturalOrigin"`
	GeographicOrigin string   `json:"geographicOrigin"`
	UsageRegions     []string `json:"usageRegions"`

	// Metadata
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Notes     []string  `json:"notes"`
}

// SemanticFieldChange represents a change in a semantic field over time.
type SemanticFieldChange struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // expansion, contraction, shift, borrowing, specialization
	Description string    `json:"description"`
	Details     string    `json:"details"`
	Timestamp   time.Time `json:"timestamp"`
	Era         string    `json:"era"`
	Trigger     string    `json:"trigger"`    // cultural_change, technological_advance, social_shift, etc.
	Intensity   float32   `json:"intensity"`  // 0.0-1.0
	Confidence  float32   `json:"confidence"` // 0.0-1.0

	// Specific changes
	AddedConcepts    []string `json:"addedConcepts,omitempty"`
	RemovedConcepts  []string `json:"removedConcepts,omitempty"`
	ModifiedConcepts []string `json:"modifiedConcepts,omitempty"`
	BorrowedTerms    []string `json:"borrowedTerms,omitempty"`
}

// NewSemanticField creates a new semantic field.
func NewSemanticField(
	id, name, description, category, domain, era string,
) *SemanticField {
	now := time.Now()
	return &SemanticField{
		ID:               id,
		Name:             name,
		Description:      description,
		Category:         category,
		Domain:           domain,
		Era:              era,
		CoreConcepts:     make([]string, 0),
		Vocabulary:       make([]string, 0),
		Synonyms:         make([]string, 0),
		EvolutionHistory: make([]SemanticFieldChange, 0),
		Complexity:       0.5,
		Richness:         0.5,
		Stability:        0.8,
		RelatedFields:    make([]string, 0),
		ChildFields:      make([]string, 0),
		UsageRegions:     make([]string, 0),
		CreatedAt:        now,
		UpdatedAt:        now,
		Notes:            make([]string, 0),
	}
}

// AddCoreConcept adds a core concept to the semantic field.
func (sf *SemanticField) AddCoreConcept(concept string) {
	if !sf.hasConcept(concept) {
		sf.CoreConcepts = append(sf.CoreConcepts, concept)
		sf.UpdatedAt = time.Now()
		sf.updateMetrics()
	}
}

// RemoveCoreConcept removes a core concept from the semantic field.
func (sf *SemanticField) RemoveCoreConcept(concept string) {
	for i, c := range sf.CoreConcepts {
		if c == concept {
			sf.CoreConcepts = append(sf.CoreConcepts[:i], sf.CoreConcepts[i+1:]...)
			sf.UpdatedAt = time.Now()
			sf.updateMetrics()
			break
		}
	}
}

// AddVocabulary adds vocabulary to the semantic field.
func (sf *SemanticField) AddVocabulary(terms ...string) {
	for _, term := range terms {
		if !sf.hasVocabulary(term) {
			sf.Vocabulary = append(sf.Vocabulary, term)
		}
	}
	sf.UpdatedAt = time.Now()
	sf.updateMetrics()
}

// RemoveVocabulary removes vocabulary from the semantic field.
func (sf *SemanticField) RemoveVocabulary(term string) {
	for i, v := range sf.Vocabulary {
		if v == term {
			sf.Vocabulary = append(sf.Vocabulary[:i], sf.Vocabulary[i+1:]...)
			sf.UpdatedAt = time.Now()
			sf.updateMetrics()
			break
		}
	}
}

// AddSynonym adds a synonym for existing concepts.
func (sf *SemanticField) AddSynonym(synonym string) {
	if !sf.hasSynonym(synonym) {
		sf.Synonyms = append(sf.Synonyms, synonym)
		sf.UpdatedAt = time.Now()
	}
}

// hasConcept checks if a concept exists in the semantic field.
func (sf *SemanticField) hasConcept(concept string) bool {
	for _, c := range sf.CoreConcepts {
		if c == concept {
			return true
		}
	}
	return false
}

// hasVocabulary checks if vocabulary exists in the semantic field.
func (sf *SemanticField) hasVocabulary(term string) bool {
	for _, v := range sf.Vocabulary {
		if v == term {
			return true
		}
	}
	return false
}

// hasSynonym checks if a synonym exists in the semantic field.
func (sf *SemanticField) hasSynonym(synonym string) bool {
	for _, s := range sf.Synonyms {
		if s == synonym {
			return true
		}
	}
	return false
}

// AddEvolutionStep adds an evolution step to the semantic field.
func (sf *SemanticField) AddEvolutionStep(change SemanticFieldChange) {
	sf.EvolutionHistory = append(sf.EvolutionHistory, change)
	sf.UpdatedAt = time.Now()

	// Update metrics based on the change
	sf.updateMetrics()
}

// updateMetrics updates the complexity, richness, and stability metrics.
func (sf *SemanticField) updateMetrics() {
	// Complexity: based on number of concepts and relationships
	conceptComplexity := float32(len(sf.CoreConcepts)) * 0.1
	vocabularyComplexity := float32(len(sf.Vocabulary)) * 0.05
	relationshipComplexity := float32(len(sf.RelatedFields)) * 0.05
	sf.Complexity = min(1.0, 0.3+conceptComplexity+vocabularyComplexity+relationshipComplexity)

	// Richness: based on vocabulary diversity and synonyms
	vocabularyRichness := float32(len(sf.Vocabulary)) * 0.08
	synonymRichness := float32(len(sf.Synonyms)) * 0.1
	sf.Richness = min(1.0, 0.2+vocabularyRichness+synonymRichness)

	// Stability: decreases with more recent changes
	if len(sf.EvolutionHistory) > 0 {
		recentChanges := 0
		for _, change := range sf.EvolutionHistory {
			if time.Since(change.Timestamp) < 24*time.Hour {
				recentChanges++
			}
		}
		stabilityPenalty := float32(recentChanges) * 0.1
		sf.Stability = max(0.0, 0.8-stabilityPenalty)
	}
}

// AddRelatedField adds a related semantic field.
func (sf *SemanticField) AddRelatedField(fieldID string) {
	if !sf.hasRelatedField(fieldID) {
		sf.RelatedFields = append(sf.RelatedFields, fieldID)
		sf.UpdatedAt = time.Now()
		sf.updateMetrics()
	}
}

// RemoveRelatedField removes a related semantic field.
func (sf *SemanticField) RemoveRelatedField(fieldID string) {
	for i, id := range sf.RelatedFields {
		if id == fieldID {
			sf.RelatedFields = append(sf.RelatedFields[:i], sf.RelatedFields[i+1:]...)
			sf.UpdatedAt = time.Now()
			sf.updateMetrics()
			break
		}
	}
}

// hasRelatedField checks if a field is related to this semantic field.
func (sf *SemanticField) hasRelatedField(fieldID string) bool {
	for _, id := range sf.RelatedFields {
		if id == fieldID {
			return true
		}
	}
	return false
}

// SetParentField sets the parent semantic field.
func (sf *SemanticField) SetParentField(parentID string) {
	sf.ParentField = &parentID
	sf.UpdatedAt = time.Now()
}

// AddChildField adds a child semantic field.
func (sf *SemanticField) AddChildField(childID string) {
	if !sf.hasChildField(childID) {
		sf.ChildFields = append(sf.ChildFields, childID)
		sf.UpdatedAt = time.Now()
	}
}

// RemoveChildField removes a child semantic field.
func (sf *SemanticField) RemoveChildField(childID string) {
	for i, id := range sf.ChildFields {
		if id == childID {
			sf.ChildFields = append(sf.ChildFields[:i], sf.ChildFields[i+1:]...)
			sf.UpdatedAt = time.Now()
			break
		}
	}
}

// hasChildField checks if a field is a child of this semantic field.
func (sf *SemanticField) hasChildField(childID string) bool {
	for _, id := range sf.ChildFields {
		if id == childID {
			return true
		}
	}
	return false
}

// GetAncestors returns all ancestor semantic fields in the hierarchy.
func (sf *SemanticField) GetAncestors(fields map[string]*SemanticField) []*SemanticField {
	var ancestors []*SemanticField
	current := sf

	for current.ParentField != nil {
		parent, exists := fields[*current.ParentField]
		if !exists {
			break
		}
		ancestors = append(ancestors, parent)
		current = parent
	}

	return ancestors
}

// GetDescendants returns all descendant semantic fields in the hierarchy.
func (sf *SemanticField) GetDescendants(fields map[string]*SemanticField) []*SemanticField {
	var descendants []*SemanticField
	visited := make(map[string]bool)

	var collectDescendants func(field *SemanticField)
	collectDescendants = func(field *SemanticField) {
		if visited[field.ID] {
			return
		}
		visited[field.ID] = true

		for _, childID := range field.ChildFields {
			if child, exists := fields[childID]; exists {
				descendants = append(descendants, child)
				collectDescendants(child)
			}
		}
	}

	collectDescendants(sf)
	return descendants
}

// CalculateSimilarity calculates the similarity between two semantic fields.
func (sf *SemanticField) CalculateSimilarity(other *SemanticField) float32 {
	if sf.ID == other.ID {
		return 1.0
	}

	// Calculate concept similarity
	conceptSimilarity := sf.calculateConceptSimilarity(other)

	// Calculate domain similarity
	domainSimilarity := sf.calculateDomainSimilarity(other)

	// Calculate structural similarity
	structuralSimilarity := sf.calculateStructuralSimilarity(other)

	// Weighted combination
	return (conceptSimilarity * 0.4) + (domainSimilarity * 0.4) + (structuralSimilarity * 0.2)
}

// calculateConceptSimilarity calculates similarity based on shared concepts and vocabulary.
func (sf *SemanticField) calculateConceptSimilarity(other *SemanticField) float32 {
	// Count shared core concepts
	sharedConcepts := 0
	for _, concept := range sf.CoreConcepts {
		if other.hasConcept(concept) {
			sharedConcepts++
		}
	}

	// Count shared vocabulary
	sharedVocabulary := 0
	for _, term := range sf.Vocabulary {
		if other.hasVocabulary(term) {
			sharedVocabulary++
		}
	}

	// Calculate similarity scores
	conceptScore := float32(sharedConcepts) / float32(maxInt(len(sf.CoreConcepts), len(other.CoreConcepts)))
	vocabularyScore := float32(sharedVocabulary) / float32(maxInt(len(sf.Vocabulary), len(other.Vocabulary)))

	return (conceptScore + vocabularyScore) / 2.0
}

// calculateDomainSimilarity calculates similarity based on domain and category.
func (sf *SemanticField) calculateDomainSimilarity(other *SemanticField) float32 {
	similarity := float32(0.0)

	// Domain similarity
	if sf.Domain == other.Domain {
		similarity += 0.6
	}

	// Category similarity
	if sf.Category == other.Category {
		similarity += 0.4
	}

	return similarity
}

// calculateStructuralSimilarity calculates similarity based on structural characteristics.
func (sf *SemanticField) calculateStructuralSimilarity(other *SemanticField) float32 {
	similarity := float32(0.0)

	// Era similarity
	if sf.Era == other.Era {
		similarity += 0.3
	}

	// Complexity similarity (closer complexity = higher similarity)
	complexityDiff := abs(sf.Complexity - other.Complexity)
	similarity += (1.0 - complexityDiff) * 0.3

	// Richness similarity
	richnessDiff := abs(sf.Richness - other.Richness)
	similarity += (1.0 - richnessDiff) * 0.2

	// Stability similarity
	stabilityDiff := abs(sf.Stability - other.Stability)
	similarity += (1.0 - stabilityDiff) * 0.2

	return similarity
}

// abs returns the absolute value of a float32.
func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

// maxInt returns the maximum of two integers.
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// SemanticFieldEvolutionEngine manages the evolution of semantic fields over time.
type SemanticFieldEvolutionEngine struct {
	// Evolution parameters
	ExpansionRate      float32 `json:"expansionRate"`      // Rate of field expansion
	ContractionRate    float32 `json:"contractionRate"`    // Rate of field contraction
	ShiftRate          float32 `json:"shiftRate"`          // Rate of semantic shifts
	BorrowingRate      float32 `json:"borrowingRate"`      // Rate of borrowing from other fields
	SpecializationRate float32 `json:"specializationRate"` // Rate of field specialization

	// Change triggers
	CulturalChangeWeight      float32 `json:"culturalChangeWeight"`      // Weight of cultural changes
	TechnologicalWeight       float32 `json:"technologicalWeight"`       // Weight of technological advances
	SocialShiftWeight         float32 `json:"socialShiftWeight"`         // Weight of social shifts
	GeographicExpansionWeight float32 `json:"geographicExpansionWeight"` // Weight of geographic expansion

	// Evolution history
	EvolutionEvents []EvolutionEvent `json:"evolutionEvents"`
}

// EvolutionEvent represents a significant evolution event affecting multiple semantic fields.
type EvolutionEvent struct {
	ID             string    `json:"id"`
	Type           string    `json:"type"` // cultural_shift, technological_breakthrough, social_change, etc.
	Description    string    `json:"description"`
	Timestamp      time.Time `json:"timestamp"`
	Era            string    `json:"era"`
	Intensity      float32   `json:"intensity"`      // 0.0-1.0
	AffectedFields []string  `json:"affectedFields"` // IDs of semantic fields affected
	Details        string    `json:"details"`
}

// NewSemanticFieldEvolutionEngine creates a new semantic field evolution engine.
func NewSemanticFieldEvolutionEngine() *SemanticFieldEvolutionEngine {
	return &SemanticFieldEvolutionEngine{
		ExpansionRate:             0.3,
		ContractionRate:           0.1,
		ShiftRate:                 0.2,
		BorrowingRate:             0.25,
		SpecializationRate:        0.15,
		CulturalChangeWeight:      0.3,
		TechnologicalWeight:       0.25,
		SocialShiftWeight:         0.25,
		GeographicExpansionWeight: 0.2,
		EvolutionEvents:           make([]EvolutionEvent, 0),
	}
}

// EvolveSemanticField evolves a semantic field based on various factors.
func (sfe *SemanticFieldEvolutionEngine) EvolveSemanticField(
	field *SemanticField,
	fields map[string]*SemanticField,
	era string,
	triggers []string,
) error {
	// Determine evolution type based on triggers and rates
	evolutionType := sfe.determineEvolutionType(triggers)

	// Create evolution change
	change := SemanticFieldChange{
		ID:          fmt.Sprintf("evolution_%s_%d", field.ID, len(field.EvolutionHistory)+1),
		Type:        evolutionType,
		Description: sfe.generateEvolutionDescription(evolutionType, triggers),
		Details:     sfe.generateEvolutionDetails(evolutionType, field, triggers),
		Timestamp:   time.Now(),
		Era:         era,
		Trigger:     strings.Join(triggers, "_"),
		Intensity:   sfe.calculateIntensity(triggers),
		Confidence:  0.8,
	}

	// Apply evolution changes
	switch evolutionType {
	case "expansion":
		sfe.applyExpansion(field, change, fields)
	case "contraction":
		sfe.applyContraction(field, change)
	case "shift":
		sfe.applySemanticShift(field, change)
	case "borrowing":
		sfe.applyBorrowing(field, change, fields)
	case "specialization":
		sfe.applySpecialization(field, change)
	}

	// Add evolution step
	field.AddEvolutionStep(change)

	// Record evolution event
	sfe.recordEvolutionEvent(evolutionType, field.ID, triggers, change.Details)

	return nil
}

// determineEvolutionType determines the type of evolution based on triggers and rates.
func (sfe *SemanticFieldEvolutionEngine) determineEvolutionType(triggers []string) string {
	// Simple weighted random selection based on rates
	totalRate := sfe.ExpansionRate + sfe.ContractionRate + sfe.ShiftRate + sfe.BorrowingRate + sfe.SpecializationRate
	random := rand.Float32() * totalRate

	current := float32(0.0)

	current += sfe.ExpansionRate
	if random <= current {
		return "expansion"
	}

	current += sfe.ContractionRate
	if random <= current {
		return "contraction"
	}

	current += sfe.ShiftRate
	if random <= current {
		return "shift"
	}

	current += sfe.BorrowingRate
	if random <= current {
		return "borrowing"
	}

	return "specialization"
}

// generateEvolutionDescription generates a description for the evolution change.
func (sfe *SemanticFieldEvolutionEngine) generateEvolutionDescription(
	evolutionType string,
	triggers []string,
) string {
	switch evolutionType {
	case "expansion":
		return "Semantic field expanded due to external influences"
	case "contraction":
		return "Semantic field contracted due to changing circumstances"
	case "shift":
		return "Semantic field experienced conceptual shifts"
	case "borrowing":
		return "Semantic field borrowed concepts from related fields"
	case "specialization":
		return "Semantic field became more specialized"
	default:
		return "Semantic field evolved in response to external factors"
	}
}

// generateEvolutionDetails generates detailed information about the evolution change.
func (sf *SemanticFieldEvolutionEngine) generateEvolutionDetails(
	evolutionType string,
	field *SemanticField,
	triggers []string,
) string {
	switch evolutionType {
	case "expansion":
		return fmt.Sprintf("Field '%s' expanded with %d new concepts and %d vocabulary terms",
			field.Name, len(field.CoreConcepts), len(field.Vocabulary))
	case "contraction":
		return fmt.Sprintf("Field '%s' contracted, reducing complexity and vocabulary", field.Name)
	case "shift":
		return fmt.Sprintf("Field '%s' experienced semantic shifts in %d concepts",
			field.Name, len(field.CoreConcepts))
	case "borrowing":
		return fmt.Sprintf("Field '%s' borrowed concepts from %d related fields",
			field.Name, len(field.RelatedFields))
	case "specialization":
		return fmt.Sprintf("Field '%s' became more specialized with focused concepts", field.Name)
	default:
		return fmt.Sprintf("Field '%s' evolved in response to %s", field.Name, strings.Join(triggers, ", "))
	}
}

// calculateIntensity calculates the intensity of the evolution change.
func (sfe *SemanticFieldEvolutionEngine) calculateIntensity(triggers []string) float32 {
	intensity := float32(0.0)

	for _, trigger := range triggers {
		switch trigger {
		case "cultural_change":
			intensity += sfe.CulturalChangeWeight
		case "technological_advance":
			intensity += sfe.TechnologicalWeight
		case "social_shift":
			intensity += sfe.SocialShiftWeight
		case "geographic_expansion":
			intensity += sfe.GeographicExpansionWeight
		default:
			intensity += 0.1
		}
	}

	return min(1.0, intensity)
}

// applyExpansion applies expansion evolution to a semantic field.
func (sfe *SemanticFieldEvolutionEngine) applyExpansion(
	field *SemanticField,
	change SemanticFieldChange,
	fields map[string]*SemanticField,
) {
	// Add new concepts
	newConcepts := []string{"expanded_concept_1", "expanded_concept_2", "expanded_concept_3"}
	for _, concept := range newConcepts {
		field.AddCoreConcept(concept)
		change.AddedConcepts = append(change.AddedConcepts, concept)
	}

	// Add new vocabulary
	newVocabulary := []string{"expanded_term_1", "expanded_term_2", "expanded_term_3", "expanded_term_4"}
	for _, term := range newVocabulary {
		field.AddVocabulary(term)
		change.AddedConcepts = append(change.AddedConcepts, term)
	}

	// Add related fields if available
	for _, relatedField := range fields {
		if relatedField.ID != field.ID && !field.hasRelatedField(relatedField.ID) {
			if rand.Float32() < 0.3 { // 30% chance
				field.AddRelatedField(relatedField.ID)
			}
		}
	}
}

// applyContraction applies contraction evolution to a semantic field.
func (sfe *SemanticFieldEvolutionEngine) applyContraction(
	field *SemanticField,
	change SemanticFieldChange,
) {
	// Remove some concepts if available
	if len(field.CoreConcepts) > 2 {
		removeCount := minInt(2, len(field.CoreConcepts)/3)
		for i := 0; i < removeCount; i++ {
			if len(field.CoreConcepts) > 0 {
				removed := field.CoreConcepts[len(field.CoreConcepts)-1]
				field.RemoveCoreConcept(removed)
				change.RemovedConcepts = append(change.RemovedConcepts, removed)
			}
		}
	}

	// Remove some vocabulary if available
	if len(field.Vocabulary) > 3 {
		removeCount := minInt(3, len(field.Vocabulary)/4)
		for i := 0; i < removeCount; i++ {
			if len(field.Vocabulary) > 0 {
				removed := field.Vocabulary[len(field.Vocabulary)-1]
				field.RemoveVocabulary(removed)
				change.RemovedConcepts = append(change.RemovedConcepts, removed)
			}
		}
	}
}

// applySemanticShift applies semantic shift evolution to a semantic field.
func (sfe *SemanticFieldEvolutionEngine) applySemanticShift(
	field *SemanticField,
	change SemanticFieldChange,
) {
	// Modify existing concepts
	if len(field.CoreConcepts) > 0 {
		shiftedConcepts := []string{}
		for _, concept := range field.CoreConcepts {
			shifted := fmt.Sprintf("shifted_%s", concept)
			shiftedConcepts = append(shiftedConcepts, shifted)
		}
		change.ModifiedConcepts = shiftedConcepts
	}

	// Add new shifted concepts
	newShiftedConcepts := []string{"shifted_concept_1", "shifted_concept_2"}
	for _, concept := range newShiftedConcepts {
		field.AddCoreConcept(concept)
		change.AddedConcepts = append(change.AddedConcepts, concept)
	}
}

// applyBorrowing applies borrowing evolution to a semantic field.
func (sfe *SemanticFieldEvolutionEngine) applyBorrowing(
	field *SemanticField,
	change SemanticFieldChange,
	fields map[string]*SemanticField,
) {
	// Borrow concepts from related fields
	borrowedTerms := []string{}
	for _, relatedID := range field.RelatedFields {
		if relatedField, exists := fields[relatedID]; exists {
			if len(relatedField.CoreConcepts) > 0 {
				borrowed := fmt.Sprintf("borrowed_%s", relatedField.CoreConcepts[0])
				field.AddCoreConcept(borrowed)
				borrowedTerms = append(borrowedTerms, borrowed)
			}
		}
	}
	change.BorrowedTerms = borrowedTerms
}

// applySpecialization applies specialization evolution to a semantic field.
func (sfe *SemanticFieldEvolutionEngine) applySpecialization(
	field *SemanticField,
	change SemanticFieldChange,
) {
	// Add specialized concepts
	specializedConcepts := []string{"specialized_concept_1", "specialized_concept_2"}
	for _, concept := range specializedConcepts {
		field.AddCoreConcept(concept)
		change.AddedConcepts = append(change.AddedConcepts, concept)
	}

	// Add specialized vocabulary
	specializedVocabulary := []string{"specialized_term_1", "specialized_term_2"}
	for _, term := range specializedVocabulary {
		field.AddVocabulary(term)
		change.AddedConcepts = append(change.AddedConcepts, term)
	}
}

// recordEvolutionEvent records a significant evolution event.
func (sfe *SemanticFieldEvolutionEngine) recordEvolutionEvent(
	eventType, fieldID string,
	triggers []string,
	details string,
) {
	event := EvolutionEvent{
		ID:             fmt.Sprintf("event_%s_%d", eventType, len(sfe.EvolutionEvents)+1),
		Type:           eventType,
		Description:    fmt.Sprintf("%s evolution in semantic field", eventType),
		Timestamp:      time.Now(),
		Era:            "current",
		Intensity:      0.7,
		AffectedFields: []string{fieldID},
		Details:        details,
	}

	sfe.EvolutionEvents = append(sfe.EvolutionEvents, event)
}

// minInt returns the minimum of two integers.
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// LanguageRegister represents a specific style or register of a language.
type LanguageRegister struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`    // formal, informal, technical, literary, etc.
	Level       string `json:"level"`   // high, medium, low
	Context     string `json:"context"` // academic, business, casual, religious, etc.

	// Linguistic characteristics
	PhonologicalFeatures  []string `json:"phonologicalFeatures"`  // Sound patterns specific to this register
	MorphologicalFeatures []string `json:"morphologicalFeatures"` // Word formation patterns
	SyntacticFeatures     []string `json:"syntacticFeatures"`     // Sentence structure patterns
	LexicalFeatures       []string `json:"lexicalFeatures"`       // Vocabulary choices

	// Social and cultural context
	SocialClass       string `json:"socialClass"`       // upper, middle, lower, etc.
	AgeGroup          string `json:"ageGroup"`          // young, adult, elderly, etc.
	Gender            string `json:"gender"`            // masculine, feminine, neutral, etc.
	EducationLevel    string `json:"educationLevel"`    // primary, secondary, university, etc.
	OccupationalField string `json:"occupationalField"` // academic, medical, legal, etc.

	// Usage patterns
	FormalityLevel  float32 `json:"formalityLevel"`  // 0.0-1.0 (informal to formal)
	ComplexityLevel float32 `json:"complexityLevel"` // 0.0-1.0 (simple to complex)
	PrestigeLevel   float32 `json:"prestigeLevel"`   // 0.0-1.0 (low to high prestige)

	// Evolution tracking
	EvolutionHistory []RegisterChange `json:"evolutionHistory"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt"`
	Notes            []string         `json:"notes"`
}

// RegisterChange represents a change in a language register over time.
type RegisterChange struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // phonological, morphological, syntactic, lexical, social
	Description string    `json:"description"`
	Details     string    `json:"details"`
	Timestamp   time.Time `json:"timestamp"`
	Era         string    `json:"era"`
	Trigger     string    `json:"trigger"`    // social_change, technological_advance, cultural_shift, etc.
	Intensity   float32   `json:"intensity"`  // 0.0-1.0
	Confidence  float32   `json:"confidence"` // 0.0-1.0

	// Specific changes
	AddedFeatures    []string `json:"addedFeatures,omitempty"`
	RemovedFeatures  []string `json:"removedFeatures,omitempty"`
	ModifiedFeatures []string `json:"modifiedFeatures,omitempty"`
}

// CreoleLanguage represents a pidgin that has become a native language.
type CreoleLanguage struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Era         string `json:"era"`
	Type        string `json:"type"`    // plantation, trade, urban, etc.
	Context     string `json:"context"` // agricultural, commercial, urban, etc.

	// Origin information
	PidginID        string   `json:"pidginID"`        // ID of the pidgin that became a creole
	SourceLanguages []string `json:"sourceLanguages"` // IDs of original source languages
	PrimarySource   string   `json:"primarySource"`   // ID of the dominant source language
	SecondarySource string   `json:"secondarySource"` // ID of the secondary source language

	// Linguistic characteristics
	PhonologicalFeatures  []string `json:"phonologicalFeatures"`  // Developed sound patterns
	MorphologicalFeatures []string `json:"morphologicalFeatures"` // Developed word formation
	SyntacticFeatures     []string `json:"syntacticFeatures"`     // Developed sentence structure
	LexicalFeatures       []string `json:"lexicalFeatures"`       // Developed vocabulary

	// Development characteristics
	ComplexityLevel    float32 `json:"complexityLevel"`    // 0.0-1.0 (simple to complex)
	StabilityLevel     float32 `json:"stabilityLevel"`     // 0.0-1.0 (unstable to stable)
	FunctionalityLevel float32 `json:"functionalityLevel"` // 0.0-1.0 (limited to comprehensive)
	NativeSpeakerCount int     `json:"nativeSpeakerCount"` // Number of native speakers

	// Usage patterns
	TotalSpeakerCount int      `json:"totalSpeakerCount"` // Total number of speakers
	UsageRegions      []string `json:"usageRegions"`      // Regions where creole is used
	UsageDomains      []string `json:"usageDomains"`      // Domains where creole is used
	GenerationCount   int      `json:"generationCount"`   // Number of generations using the creole

	// Evolution tracking
	EvolutionHistory []CreoleChange `json:"evolutionHistory"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	Notes            []string       `json:"notes"`
}

// CreoleChange represents a change in a creole language over time.
type CreoleChange struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // phonological, morphological, syntactic, lexical, functional
	Description string    `json:"description"`
	Details     string    `json:"details"`
	Timestamp   time.Time `json:"timestamp"`
	Era         string    `json:"era"`
	Trigger     string    `json:"trigger"`    // generational_change, social_mobility, etc.
	Intensity   float32   `json:"intensity"`  // 0.0-1.0
	Confidence  float32   `json:"confidence"` // 0.0-1.0

	// Specific changes
	AddedFeatures    []string `json:"addedFeatures,omitempty"`
	RemovedFeatures  []string `json:"removedFeatures,omitempty"`
	ModifiedFeatures []string `json:"modifiedFeatures,omitempty"`
}

// NewPidginLanguage creates a new pidgin language.
func NewPidginLanguage(
	id, name, description, era, pidginType, context string,
	sourceLanguages []string,
	primarySource, secondarySource string,
) *PidginLanguage {
	now := time.Now()
	return &PidginLanguage{
		ID:                    id,
		Name:                  name,
		Description:           description,
		Era:                   era,
		Type:                  pidginType,
		Context:               context,
		SourceLanguages:       sourceLanguages,
		PrimarySource:         primarySource,
		SecondarySource:       secondarySource,
		PhonologicalFeatures:  make([]string, 0),
		MorphologicalFeatures: make([]string, 0),
		SyntacticFeatures:     make([]string, 0),
		LexicalFeatures:       make([]string, 0),
		ComplexityLevel:       0.3,
		StabilityLevel:        0.4,
		FunctionalityLevel:    0.5,
		SpeakerCount:          0,
		UsageRegions:          make([]string, 0),
		UsageDomains:          make([]string, 0),
		GenerationCount:       1,
		EvolutionHistory:      make([]PidginChange, 0),
		CreatedAt:             now,
		UpdatedAt:             now,
		Notes:                 make([]string, 0),
	}
}

// NewCreoleLanguage creates a new creole language from a pidgin.
func NewCreoleLanguage(
	id, name, description, era, creoleType, context string,
	pidginID string,
	sourceLanguages []string,
	primarySource, secondarySource string,
) *CreoleLanguage {
	now := time.Now()
	return &CreoleLanguage{
		ID:                    id,
		Name:                  name,
		Description:           description,
		Era:                   era,
		Type:                  creoleType,
		Context:               context,
		PidginID:              pidginID,
		SourceLanguages:       sourceLanguages,
		PrimarySource:         primarySource,
		SecondarySource:       secondarySource,
		PhonologicalFeatures:  make([]string, 0),
		MorphologicalFeatures: make([]string, 0),
		SyntacticFeatures:     make([]string, 0),
		LexicalFeatures:       make([]string, 0),
		ComplexityLevel:       0.6,
		StabilityLevel:        0.7,
		FunctionalityLevel:    0.8,
		NativeSpeakerCount:    0,
		TotalSpeakerCount:     0,
		UsageRegions:          make([]string, 0),
		UsageDomains:          make([]string, 0),
		GenerationCount:       2,
		EvolutionHistory:      make([]CreoleChange, 0),
		CreatedAt:             now,
		UpdatedAt:             now,
		Notes:                 make([]string, 0),
	}
}

// AddPhonologicalFeature adds a phonological feature to the pidgin.
func (pl *PidginLanguage) AddPhonologicalFeature(feature string) {
	if !pl.hasPhonologicalFeature(feature) {
		pl.PhonologicalFeatures = append(pl.PhonologicalFeatures, feature)
		pl.UpdatedAt = time.Now()
		pl.updateMetrics()
	}
}

// hasPhonologicalFeature checks if a phonological feature exists in the pidgin.
func (pl *PidginLanguage) hasPhonologicalFeature(feature string) bool {
	for _, f := range pl.PhonologicalFeatures {
		if f == feature {
			return true
		}
	}
	return false
}

// AddMorphologicalFeature adds a morphological feature to the pidgin.
func (pl *PidginLanguage) AddMorphologicalFeature(feature string) {
	if !pl.hasMorphologicalFeature(feature) {
		pl.MorphologicalFeatures = append(pl.MorphologicalFeatures, feature)
		pl.UpdatedAt = time.Now()
		pl.updateMetrics()
	}
}

// hasMorphologicalFeature checks if a morphological feature exists in the pidgin.
func (pl *PidginLanguage) hasMorphologicalFeature(feature string) bool {
	for _, f := range pl.MorphologicalFeatures {
		if f == feature {
			return true
		}
	}
	return false
}

// AddSyntacticFeature adds a syntactic feature to the pidgin.
func (pl *PidginLanguage) AddSyntacticFeature(feature string) {
	if !pl.hasSyntacticFeature(feature) {
		pl.SyntacticFeatures = append(pl.SyntacticFeatures, feature)
		pl.UpdatedAt = time.Now()
		pl.updateMetrics()
	}
}

// hasSyntacticFeature checks if a syntactic feature exists in the pidgin.
func (pl *PidginLanguage) hasSyntacticFeature(feature string) bool {
	for _, f := range pl.SyntacticFeatures {
		if f == feature {
			return true
		}
	}
	return false
}

// AddLexicalFeature adds a lexical feature to the pidgin.
func (pl *PidginLanguage) AddLexicalFeature(feature string) {
	if !pl.hasLexicalFeature(feature) {
		pl.LexicalFeatures = append(pl.LexicalFeatures, feature)
		pl.UpdatedAt = time.Now()
		pl.updateMetrics()
	}
}

// hasLexicalFeature checks if a lexical feature exists in the pidgin.
func (pl *PidginLanguage) hasLexicalFeature(feature string) bool {
	for _, f := range pl.LexicalFeatures {
		if f == feature {
			return true
		}
	}
	return false
}

// SetSpeakerCount sets the number of speakers of the pidgin.
func (pl *PidginLanguage) SetSpeakerCount(count int) {
	pl.SpeakerCount = count
	pl.UpdatedAt = time.Now()
	pl.updateMetrics()
}

// AddUsageRegion adds a usage region to the pidgin.
func (pl *PidginLanguage) AddUsageRegion(region string) {
	if !pl.hasUsageRegion(region) {
		pl.UsageRegions = append(pl.UsageRegions, region)
		pl.UpdatedAt = time.Now()
	}
}

// hasUsageRegion checks if a region exists in the usage regions list.
func (pl *PidginLanguage) hasUsageRegion(region string) bool {
	for _, r := range pl.UsageRegions {
		if r == region {
			return true
		}
	}
	return false
}

// AddUsageDomain adds a usage domain to the pidgin.
func (pl *PidginLanguage) AddUsageDomain(domain string) {
	if !pl.hasUsageDomain(domain) {
		pl.UsageDomains = append(pl.UsageDomains, domain)
		pl.UpdatedAt = time.Now()
	}
}

// hasUsageDomain checks if a domain exists in the usage domains list.
func (pl *PidginLanguage) hasUsageDomain(domain string) bool {
	for _, d := range pl.UsageDomains {
		if d == domain {
			return true
		}
	}
	return false
}

// updateMetrics updates the complexity, stability, and functionality metrics.
func (pl *PidginLanguage) updateMetrics() {
	// Complexity level: based on features and speaker count
	complexityScore := float32(0.0)
	complexityScore += float32(len(pl.PhonologicalFeatures)) * 0.08
	complexityScore += float32(len(pl.LexicalFeatures)) * 0.06

	// Speaker count bonus (more speakers = more complexity)
	if pl.SpeakerCount > 0 {
		speakerBonus := min(0.2, float32(pl.SpeakerCount)*0.001)
		complexityScore += speakerBonus
	}

	pl.ComplexityLevel = min(1.0, 0.2+complexityScore)

	// Stability level: based on generation count and speaker count
	stabilityScore := float32(0.0)
	stabilityScore += float32(pl.GenerationCount) * 0.15
	if pl.SpeakerCount > 100 {
		stabilityScore += 0.2
	}

	pl.StabilityLevel = min(1.0, 0.3+stabilityScore)

	// Functionality level: based on features and usage domains
	functionalityScore := float32(0.0)
	functionalityScore += float32(len(pl.UsageDomains)) * 0.1
	functionalityScore += float32(len(pl.LexicalFeatures)) * 0.05
	functionalityScore += float32(len(pl.SyntacticFeatures)) * 0.08

	pl.FunctionalityLevel = min(1.0, 0.4+functionalityScore)
}

// AddEvolutionStep adds an evolution step to the pidgin.
func (pl *PidginLanguage) AddEvolutionStep(change PidginChange) {
	pl.EvolutionHistory = append(pl.EvolutionHistory, change)
	pl.UpdatedAt = time.Now()
}

// AddPhonologicalFeature adds a phonological feature to the creole.
func (cl *CreoleLanguage) AddPhonologicalFeature(feature string) {
	if !cl.hasPhonologicalFeature(feature) {
		cl.PhonologicalFeatures = append(cl.PhonologicalFeatures, feature)
		cl.UpdatedAt = time.Now()
		cl.updateMetrics()
	}
}

// hasPhonologicalFeature checks if a phonological feature exists in the creole.
func (cl *CreoleLanguage) hasPhonologicalFeature(feature string) bool {
	for _, f := range cl.PhonologicalFeatures {
		if f == feature {
			return true
		}
	}
	return false
}

// AddMorphologicalFeature adds a morphological feature to the creole.
func (cl *CreoleLanguage) AddMorphologicalFeature(feature string) {
	if !cl.hasMorphologicalFeature(feature) {
		cl.MorphologicalFeatures = append(cl.MorphologicalFeatures, feature)
		cl.UpdatedAt = time.Now()
		cl.updateMetrics()
	}
}

// hasMorphologicalFeature checks if a morphological feature exists in the creole.
func (cl *CreoleLanguage) hasMorphologicalFeature(feature string) bool {
	for _, f := range cl.MorphologicalFeatures {
		if f == feature {
			return true
		}
	}
	return false
}

// AddSyntacticFeature adds a syntactic feature to the creole.
func (cl *CreoleLanguage) AddSyntacticFeature(feature string) {
	if !cl.hasSyntacticFeature(feature) {
		cl.SyntacticFeatures = append(cl.SyntacticFeatures, feature)
		cl.UpdatedAt = time.Now()
		cl.updateMetrics()
	}
}

// hasSyntacticFeature checks if a syntactic feature exists in the creole.
func (cl *CreoleLanguage) hasSyntacticFeature(feature string) bool {
	for _, f := range cl.SyntacticFeatures {
		if f == feature {
			return true
		}
	}
	return false
}

// AddLexicalFeature adds a lexical feature to the creole.
func (cl *CreoleLanguage) AddLexicalFeature(feature string) {
	if !cl.hasLexicalFeature(feature) {
		cl.LexicalFeatures = append(cl.LexicalFeatures, feature)
		cl.UpdatedAt = time.Now()
		cl.updateMetrics()
	}
}

// hasLexicalFeature checks if a lexical feature exists in the creole.
func (cl *CreoleLanguage) hasLexicalFeature(feature string) bool {
	for _, f := range cl.LexicalFeatures {
		if f == feature {
			return true
		}
	}
	return false
}

// SetNativeSpeakerCount sets the number of native speakers of the creole.
func (cl *CreoleLanguage) SetNativeSpeakerCount(count int) {
	cl.NativeSpeakerCount = count
	cl.UpdatedAt = time.Now()
	cl.updateMetrics()
}

// SetTotalSpeakerCount sets the total number of speakers of the creole.
func (cl *CreoleLanguage) SetTotalSpeakerCount(count int) {
	cl.TotalSpeakerCount = count
	cl.UpdatedAt = time.Now()
	cl.updateMetrics()
}

// AddUsageRegion adds a usage region to the creole.
func (cl *CreoleLanguage) AddUsageRegion(region string) {
	if !cl.hasUsageRegion(region) {
		cl.UsageRegions = append(cl.UsageRegions, region)
		cl.UpdatedAt = time.Now()
	}
}

// hasUsageRegion checks if a region exists in the usage regions list.
func (cl *CreoleLanguage) hasUsageRegion(region string) bool {
	for _, r := range cl.UsageRegions {
		if r == region {
			return true
		}
	}
	return false
}

// AddUsageDomain adds a usage domain to the creole.
func (cl *CreoleLanguage) AddUsageDomain(domain string) {
	if !cl.hasUsageDomain(domain) {
		cl.UsageDomains = append(cl.UsageDomains, domain)
		cl.UpdatedAt = time.Now()
	}
}

// hasUsageDomain checks if a domain exists in the usage domains list.
func (cl *CreoleLanguage) hasUsageDomain(domain string) bool {
	for _, d := range cl.UsageDomains {
		if d == domain {
			return true
		}
	}
	return false
}

// updateMetrics updates the complexity, stability, and functionality metrics.
func (cl *CreoleLanguage) updateMetrics() {
	// Complexity level: based on features and speaker count
	complexityScore := float32(0.0)
	complexityScore += float32(len(cl.PhonologicalFeatures)) * 0.08
	complexityScore += float32(len(cl.MorphologicalFeatures)) * 0.1
	complexityScore += float32(len(cl.SyntacticFeatures)) * 0.12
	complexityScore += float32(len(cl.LexicalFeatures)) * 0.06

	// Native speaker bonus (native speakers = more complexity)
	if cl.NativeSpeakerCount > 0 {
		speakerBonus := min(0.2, float32(cl.NativeSpeakerCount)*0.002)
		complexityScore += speakerBonus
	}

	cl.ComplexityLevel = min(1.0, 0.5+complexityScore)

	// Stability level: based on generation count and native speakers
	stabilityScore := float32(0.0)
	stabilityScore += float32(cl.GenerationCount) * 0.15
	if cl.NativeSpeakerCount > 50 {
		stabilityScore += 0.2
	}

	cl.StabilityLevel = min(1.0, 0.6+stabilityScore)

	// Functionality level: based on features and usage domains
	functionalityScore := float32(0.0)
	functionalityScore += float32(len(cl.UsageDomains)) * 0.1
	functionalityScore += float32(len(cl.LexicalFeatures)) * 0.05
	functionalityScore += float32(len(cl.SyntacticFeatures)) * 0.08

	cl.FunctionalityLevel = min(1.0, 0.7+functionalityScore)
}

// AddEvolutionStep adds an evolution step to the creole.
func (cl *CreoleLanguage) AddEvolutionStep(change CreoleChange) {
	cl.EvolutionHistory = append(cl.EvolutionHistory, change)
	cl.UpdatedAt = time.Now()
}

// NewLanguageRegister creates a new language register.
func NewLanguageRegister(
	id, name, description, registerType, level, context string,
) *LanguageRegister {
	now := time.Now()
	return &LanguageRegister{
		ID:                    id,
		Name:                  name,
		Description:           description,
		Type:                  registerType,
		Level:                 level,
		Context:               context,
		PhonologicalFeatures:  make([]string, 0),
		MorphologicalFeatures: make([]string, 0),
		SyntacticFeatures:     make([]string, 0),
		LexicalFeatures:       make([]string, 0),
		SocialClass:           "",
		AgeGroup:              "",
		Gender:                "",
		EducationLevel:        "",
		OccupationalField:     "",
		FormalityLevel:        0.5,
		ComplexityLevel:       0.5,
		PrestigeLevel:         0.5,
		EvolutionHistory:      make([]RegisterChange, 0),
		CreatedAt:             now,
		UpdatedAt:             now,
		Notes:                 make([]string, 0),
	}
}

// AddPhonologicalFeature adds a phonological feature to the register.
func (lr *LanguageRegister) AddPhonologicalFeature(feature string) {
	if !lr.hasPhonologicalFeature(feature) {
		lr.PhonologicalFeatures = append(lr.PhonologicalFeatures, feature)
		lr.UpdatedAt = time.Now()
		lr.updateMetrics()
	}
}

// RemovePhonologicalFeature removes a phonological feature from the register.
func (lr *LanguageRegister) RemovePhonologicalFeature(feature string) {
	for i, f := range lr.PhonologicalFeatures {
		if f == feature {
			lr.PhonologicalFeatures = append(lr.PhonologicalFeatures[:i], lr.PhonologicalFeatures[i+1:]...)
			lr.UpdatedAt = time.Now()
			lr.updateMetrics()
			break
		}
	}
}

// AddMorphologicalFeature adds a morphological feature to the register.
func (lr *LanguageRegister) AddMorphologicalFeature(feature string) {
	if !lr.hasMorphologicalFeature(feature) {
		lr.MorphologicalFeatures = append(lr.MorphologicalFeatures, feature)
		lr.UpdatedAt = time.Now()
		lr.updateMetrics()
	}
}

// RemoveMorphologicalFeature removes a morphological feature from the register.
func (lr *LanguageRegister) RemoveMorphologicalFeature(feature string) {
	for i, f := range lr.MorphologicalFeatures {
		if f == feature {
			lr.MorphologicalFeatures = append(lr.MorphologicalFeatures[:i], lr.MorphologicalFeatures[i+1:]...)
			lr.UpdatedAt = time.Now()
			lr.updateMetrics()
			break
		}
	}
}

// AddSyntacticFeature adds a syntactic feature to the register.
func (lr *LanguageRegister) AddSyntacticFeature(feature string) {
	if !lr.hasSyntacticFeature(feature) {
		lr.SyntacticFeatures = append(lr.SyntacticFeatures, feature)
		lr.UpdatedAt = time.Now()
		lr.updateMetrics()
	}
}

// RemoveSyntacticFeature removes a syntactic feature from the register.
func (lr *LanguageRegister) RemoveSyntacticFeature(feature string) {
	for i, f := range lr.SyntacticFeatures {
		if f == feature {
			lr.SyntacticFeatures = append(lr.SyntacticFeatures[:i], lr.SyntacticFeatures[i+1:]...)
			lr.UpdatedAt = time.Now()
			lr.updateMetrics()
			break
		}
	}
}

// AddLexicalFeature adds a lexical feature to the register.
func (lr *LanguageRegister) AddLexicalFeature(feature string) {
	if !lr.hasLexicalFeature(feature) {
		lr.LexicalFeatures = append(lr.LexicalFeatures, feature)
		lr.UpdatedAt = time.Now()
		lr.updateMetrics()
	}
}

// RemoveLexicalFeature removes a lexical feature from the register.
func (lr *LanguageRegister) RemoveLexicalFeature(feature string) {
	for i, f := range lr.LexicalFeatures {
		if f == feature {
			lr.LexicalFeatures = append(lr.LexicalFeatures[:i], lr.LexicalFeatures[i+1:]...)
			lr.UpdatedAt = time.Now()
			lr.updateMetrics()
			break
		}
	}
}

// hasPhonologicalFeature checks if a phonological feature exists in the register.
func (lr *LanguageRegister) hasPhonologicalFeature(feature string) bool {
	for _, f := range lr.PhonologicalFeatures {
		if f == feature {
			return true
		}
	}
	return false
}

// hasMorphologicalFeature checks if a morphological feature exists in the register.
func (lr *LanguageRegister) hasMorphologicalFeature(feature string) bool {
	for _, f := range lr.MorphologicalFeatures {
		if f == feature {
			return true
		}
	}
	return false
}

// hasSyntacticFeature checks if a syntactic feature exists in the register.
func (lr *LanguageRegister) hasSyntacticFeature(feature string) bool {
	for _, f := range lr.SyntacticFeatures {
		if f == feature {
			return true
		}
	}
	return false
}

// hasLexicalFeature checks if a lexical feature exists in the register.
func (lr *LanguageRegister) hasLexicalFeature(feature string) bool {
	for _, f := range lr.LexicalFeatures {
		if f == feature {
			return true
		}
	}
	return false
}

// SetSocialContext sets the social context information for the register.
func (lr *LanguageRegister) SetSocialContext(
	socialClass, ageGroup, gender, educationLevel, occupationalField string,
) {
	lr.SocialClass = socialClass
	lr.AgeGroup = ageGroup
	lr.Gender = gender
	lr.EducationLevel = educationLevel
	lr.OccupationalField = occupationalField
	lr.UpdatedAt = time.Now()
	lr.updateMetrics()
}

// updateMetrics updates the formality, complexity, and prestige metrics.
func (lr *LanguageRegister) updateMetrics() {
	// Formality level: based on type, level, and context
	formalityScore := float32(0.0)
	switch lr.Type {
	case "formal":
		formalityScore += 0.4
	case "informal":
		formalityScore += 0.1
	case "technical":
		formalityScore += 0.3
	case "literary":
		formalityScore += 0.35
	}

	switch lr.Level {
	case "high":
		formalityScore += 0.3
	case "medium":
		formalityScore += 0.2
	case "low":
		formalityScore += 0.1
	}

	switch lr.Context {
	case "academic":
		formalityScore += 0.2
	case "business":
		formalityScore += 0.15
	case "casual":
		formalityScore += 0.05
	case "religious":
		formalityScore += 0.25
	}

	lr.FormalityLevel = min(1.0, formalityScore)

	// Complexity level: based on features and education
	complexityScore := float32(0.0)
	complexityScore += float32(len(lr.PhonologicalFeatures)) * 0.05
	complexityScore += float32(len(lr.MorphologicalFeatures)) * 0.08
	complexityScore += float32(len(lr.SyntacticFeatures)) * 0.1
	complexityScore += float32(len(lr.LexicalFeatures)) * 0.06

	switch lr.EducationLevel {
	case "university":
		complexityScore += 0.3
	case "secondary":
		complexityScore += 0.2
	case "primary":
		complexityScore += 0.1
	}

	lr.ComplexityLevel = min(1.0, 0.2+complexityScore)

	// Prestige level: based on social class, education, and context
	prestigeScore := float32(0.0)
	switch lr.SocialClass {
	case "upper":
		prestigeScore += 0.4
	case "middle":
		prestigeScore += 0.25
	case "lower":
		prestigeScore += 0.1
	}

	switch lr.EducationLevel {
	case "university":
		prestigeScore += 0.3
	case "secondary":
		prestigeScore += 0.2
	case "primary":
		prestigeScore += 0.1
	}

	switch lr.Context {
	case "academic":
		prestigeScore += 0.2
	case "business":
		prestigeScore += 0.15
	case "legal":
		prestigeScore += 0.25
	case "casual":
		prestigeScore += 0.05
	}

	lr.PrestigeLevel = min(1.0, prestigeScore)
}

// AddEvolutionStep adds an evolution step to the register.
func (lr *LanguageRegister) AddEvolutionStep(change RegisterChange) {
	lr.EvolutionHistory = append(lr.EvolutionHistory, change)
	lr.UpdatedAt = time.Now()
}

// CalculateSimilarity calculates the similarity between two language registers.
func (lr *LanguageRegister) CalculateSimilarity(other *LanguageRegister) float32 {
	if lr.ID == other.ID {
		return 1.0
	}

	// Calculate type similarity
	typeSimilarity := lr.calculateTypeSimilarity(other)

	// Calculate feature similarity
	featureSimilarity := lr.calculateFeatureSimilarity(other)

	// Calculate social similarity
	socialSimilarity := lr.calculateSocialSimilarity(other)

	// Calculate metric similarity
	metricSimilarity := lr.calculateMetricSimilarity(other)

	// Weighted combination
	return (typeSimilarity * 0.3) + (featureSimilarity * 0.3) + (socialSimilarity * 0.2) + (metricSimilarity * 0.2)
}

// calculateTypeSimilarity calculates similarity based on register type and context.
func (lr *LanguageRegister) calculateTypeSimilarity(other *LanguageRegister) float32 {
	similarity := float32(0.0)

	// Type similarity
	if lr.Type == other.Type {
		similarity += 0.4
	}

	// Level similarity
	if lr.Level == other.Level {
		similarity += 0.3
	}

	// Context similarity
	if lr.Context == other.Context {
		similarity += 0.3
	}

	return similarity
}

// calculateFeatureSimilarity calculates similarity based on linguistic features.
func (lr *LanguageRegister) calculateFeatureSimilarity(other *LanguageRegister) float32 {
	similarity := float32(0.0)

	// Count shared features
	sharedPhonological := 0
	for _, feature := range lr.PhonologicalFeatures {
		if other.hasPhonologicalFeature(feature) {
			sharedPhonological++
		}
	}

	sharedMorphological := 0
	for _, feature := range lr.MorphologicalFeatures {
		if other.hasMorphologicalFeature(feature) {
			sharedMorphological++
		}
	}

	sharedSyntactic := 0
	for _, feature := range lr.SyntacticFeatures {
		if other.hasSyntacticFeature(feature) {
			sharedSyntactic++
		}
	}

	sharedLexical := 0
	for _, feature := range lr.LexicalFeatures {
		if other.hasLexicalFeature(feature) {
			sharedLexical++
		}
	}

	// Calculate similarity scores
	totalFeatures := len(lr.PhonologicalFeatures) + len(lr.MorphologicalFeatures) + len(lr.SyntacticFeatures) + len(lr.LexicalFeatures)
	otherTotalFeatures := len(other.PhonologicalFeatures) + len(other.MorphologicalFeatures) + len(other.SyntacticFeatures) + len(other.LexicalFeatures)

	if totalFeatures > 0 && otherTotalFeatures > 0 {
		sharedTotal := sharedPhonological + sharedMorphological + sharedSyntactic + sharedLexical
		maxTotal := maxInt(totalFeatures, otherTotalFeatures)
		similarity = float32(sharedTotal) / float32(maxTotal)
	}

	return similarity
}

// calculateSocialSimilarity calculates similarity based on social characteristics.
func (lr *LanguageRegister) calculateSocialSimilarity(other *LanguageRegister) float32 {
	similarity := float32(0.0)

	// Social class similarity
	if lr.SocialClass == other.SocialClass {
		similarity += 0.3
	}

	// Age group similarity
	if lr.AgeGroup == other.AgeGroup {
		similarity += 0.2
	}

	// Education level similarity
	if lr.EducationLevel == other.EducationLevel {
		similarity += 0.3
	}

	// Occupational field similarity
	if lr.OccupationalField == other.OccupationalField {
		similarity += 0.2
	}

	return similarity
}

// calculateMetricSimilarity calculates similarity based on metric values.
func (lr *LanguageRegister) calculateMetricSimilarity(other *LanguageRegister) float32 {
	similarity := float32(0.0)

	// Formality similarity (closer formality = higher similarity)
	formalityDiff := abs(lr.FormalityLevel - other.FormalityLevel)
	similarity += (1.0 - formalityDiff) * 0.4

	// Complexity similarity
	complexityDiff := abs(lr.ComplexityLevel - other.ComplexityLevel)
	similarity += (1.0 - complexityDiff) * 0.3

	// Prestige similarity
	prestigeDiff := abs(lr.PrestigeLevel - other.PrestigeLevel)
	similarity += (1.0 - prestigeDiff) * 0.3

	return similarity
}

// RegisterEvolutionEngine manages the evolution of language registers over time.
type RegisterEvolutionEngine struct {
	// Evolution parameters
	FeatureAdditionRate     float32 `json:"featureAdditionRate"`     // Rate of adding new features
	FeatureRemovalRate      float32 `json:"featureRemovalRate"`      // Rate of removing features
	FeatureModificationRate float32 `json:"featureModificationRate"` // Rate of modifying features
	SocialChangeRate        float32 `json:"socialChangeRate"`        // Rate of social context changes

	// Change triggers
	SocialMobilityWeight       float32 `json:"socialMobilityWeight"`       // Weight of social mobility
	EducationalReformWeight    float32 `json:"educationalReformWeight"`    // Weight of educational reforms
	TechnologicalAdvanceWeight float32 `json:"technologicalAdvanceWeight"` // Weight of technological advances
	CulturalShiftWeight        float32 `json:"culturalShiftWeight"`        // Weight of cultural shifts

	// Evolution history
	EvolutionEvents []RegisterEvolutionEvent `json:"evolutionEvents"`
}

// RegisterEvolutionEvent represents a significant evolution event affecting language registers.
type RegisterEvolutionEvent struct {
	ID                string    `json:"id"`
	Type              string    `json:"type"` // social_mobility, educational_reform, technological_advance, etc.
	Description       string    `json:"description"`
	Timestamp         time.Time `json:"timestamp"`
	Era               string    `json:"era"`
	Intensity         float32   `json:"intensity"`         // 0.0-1.0
	AffectedRegisters []string  `json:"affectedRegisters"` // IDs of registers affected
	Details           string    `json:"details"`
}

// NewRegisterEvolutionEngine creates a new register evolution engine.
func NewRegisterEvolutionEngine() *RegisterEvolutionEngine {
	return &RegisterEvolutionEngine{
		FeatureAdditionRate:        0.3,
		FeatureRemovalRate:         0.1,
		FeatureModificationRate:    0.2,
		SocialChangeRate:           0.25,
		SocialMobilityWeight:       0.3,
		EducationalReformWeight:    0.25,
		TechnologicalAdvanceWeight: 0.25,
		CulturalShiftWeight:        0.2,
		EvolutionEvents:            make([]RegisterEvolutionEvent, 0),
	}
}

// EvolveRegister evolves a language register based on various factors.
func (ree *RegisterEvolutionEngine) EvolveRegister(
	register *LanguageRegister,
	era string,
	triggers []string,
) error {
	// Determine evolution type based on triggers and rates
	evolutionType := ree.determineEvolutionType(triggers)

	// Create evolution change
	change := RegisterChange{
		ID:          fmt.Sprintf("evolution_%s_%d", register.ID, len(register.EvolutionHistory)+1),
		Type:        evolutionType,
		Description: ree.generateEvolutionDescription(evolutionType, triggers),
		Details:     ree.generateEvolutionDetails(evolutionType, register, triggers),
		Timestamp:   time.Now(),
		Era:         era,
		Trigger:     strings.Join(triggers, "_"),
		Intensity:   ree.calculateIntensity(triggers),
		Confidence:  0.8,
	}

	// Apply evolution changes
	switch evolutionType {
	case "phonological":
		ree.applyPhonologicalEvolution(register, change)
	case "morphological":
		ree.applyMorphologicalEvolution(register, change)
	case "syntactic":
		ree.applySyntacticEvolution(register, change)
	case "lexical":
		ree.applyLexicalEvolution(register, change)
	case "social":
		ree.applySocialEvolution(register, change)
	}

	// Add evolution step
	register.AddEvolutionStep(change)

	// Record evolution event
	ree.recordEvolutionEvent(evolutionType, register.ID, triggers, change.Details)

	return nil
}

// determineEvolutionType determines the type of evolution based on triggers and rates.
func (ree *RegisterEvolutionEngine) determineEvolutionType(triggers []string) string {
	// Simple weighted random selection based on rates
	totalRate := ree.FeatureAdditionRate + ree.FeatureRemovalRate + ree.FeatureModificationRate + ree.SocialChangeRate
	random := rand.Float32() * totalRate

	current := float32(0.0)

	current += ree.FeatureAdditionRate
	if random <= current {
		return "phonological"
	}

	current += ree.FeatureRemovalRate
	if random <= current {
		return "morphological"
	}

	current += ree.FeatureModificationRate
	if random <= current {
		return "syntactic"
	}

	return "lexical"
}

// generateEvolutionDescription generates a description for the evolution change.
func (ree *RegisterEvolutionEngine) generateEvolutionDescription(
	evolutionType string,
	triggers []string,
) string {
	switch evolutionType {
	case "phonological":
		return "Register experienced phonological feature changes"
	case "morphological":
		return "Register experienced morphological feature changes"
	case "syntactic":
		return "Register experienced syntactic feature changes"
	case "lexical":
		return "Register experienced lexical feature changes"
	case "social":
		return "Register experienced social context changes"
	default:
		return "Register evolved in response to external factors"
	}
}

// generateEvolutionDetails generates detailed information about the evolution change.
func (ree *RegisterEvolutionEngine) generateEvolutionDetails(
	evolutionType string,
	register *LanguageRegister,
	triggers []string,
) string {
	switch evolutionType {
	case "phonological":
		return fmt.Sprintf("Register '%s' added %d new phonological features",
			register.Name, len(register.PhonologicalFeatures))
	case "morphological":
		return fmt.Sprintf("Register '%s' modified %d morphological features",
			register.Name, len(register.MorphologicalFeatures))
	case "syntactic":
		return fmt.Sprintf("Register '%s' updated %d syntactic features",
			register.Name, len(register.SyntacticFeatures))
	case "lexical":
		return fmt.Sprintf("Register '%s' expanded with %d new lexical features",
			register.Name, len(register.LexicalFeatures))
	case "social":
		return fmt.Sprintf("Register '%s' experienced social context changes", register.Name)
	default:
		return fmt.Sprintf("Register '%s' evolved in response to %s", register.Name, strings.Join(triggers, ", "))
	}
}

// calculateIntensity calculates the intensity of the evolution change.
func (ree *RegisterEvolutionEngine) calculateIntensity(triggers []string) float32 {
	intensity := float32(0.0)

	for _, trigger := range triggers {
		switch trigger {
		case "social_mobility":
			intensity += ree.SocialMobilityWeight
		case "educational_reform":
			intensity += ree.EducationalReformWeight
		case "technological_advance":
			intensity += ree.TechnologicalAdvanceWeight
		case "cultural_shift":
			intensity += ree.CulturalShiftWeight
		default:
			intensity += 0.1
		}
	}

	return min(1.0, intensity)
}

// applyPhonologicalEvolution applies phonological evolution to a register.
func (ree *RegisterEvolutionEngine) applyPhonologicalEvolution(
	register *LanguageRegister,
	change RegisterChange,
) {
	// Add new phonological features
	newFeatures := []string{"new_phonological_feature_1", "new_phonological_feature_2"}
	for _, feature := range newFeatures {
		register.AddPhonologicalFeature(feature)
		change.AddedFeatures = append(change.AddedFeatures, feature)
	}
}

// applyMorphologicalEvolution applies morphological evolution to a register.
func (ree *RegisterEvolutionEngine) applyMorphologicalEvolution(
	register *LanguageRegister,
	change RegisterChange,
) {
	// Modify existing morphological features if available
	if len(register.MorphologicalFeatures) > 0 {
		modifiedFeatures := []string{}
		for _, feature := range register.MorphologicalFeatures {
			modified := fmt.Sprintf("modified_%s", feature)
			modifiedFeatures = append(modifiedFeatures, modified)
		}
		change.ModifiedFeatures = modifiedFeatures
	}

	// Add new morphological features
	newFeatures := []string{"new_morphological_feature_1"}
	for _, feature := range newFeatures {
		register.AddMorphologicalFeature(feature)
		change.AddedFeatures = append(change.AddedFeatures, feature)
	}
}

// applySyntacticEvolution applies syntactic evolution to a register.
func (ree *RegisterEvolutionEngine) applySyntacticEvolution(
	register *LanguageRegister,
	change RegisterChange,
) {
	// Add new syntactic features
	newFeatures := []string{"new_syntactic_feature_1", "new_syntactic_feature_2"}
	for _, feature := range newFeatures {
		register.AddSyntacticFeature(feature)
		change.AddedFeatures = append(change.AddedFeatures, feature)
	}
}

// applyLexicalEvolution applies lexical evolution to a register.
func (ree *RegisterEvolutionEngine) applyLexicalEvolution(
	register *LanguageRegister,
	change RegisterChange,
) {
	// Add new lexical features
	newFeatures := []string{"new_lexical_feature_1", "new_lexical_feature_2", "new_lexical_feature_3"}
	for _, feature := range newFeatures {
		register.AddLexicalFeature(feature)
		change.AddedFeatures = append(change.AddedFeatures, feature)
	}
}

// applySocialEvolution applies social evolution to a register.
func (ree *RegisterEvolutionEngine) applySocialEvolution(
	register *LanguageRegister,
	change RegisterChange,
) {
	// Update social context
	register.SetSocialContext("middle", "adult", "neutral", "secondary", "technical")

	// Add note about social change
	change.Details = "Updated social context to reflect changing demographics"
}

// recordEvolutionEvent records a significant evolution event.
func (ree *RegisterEvolutionEngine) recordEvolutionEvent(
	eventType, registerID string,
	triggers []string,
	details string,
) {
	event := RegisterEvolutionEvent{
		ID:                fmt.Sprintf("event_%s_%d", eventType, len(ree.EvolutionEvents)+1),
		Type:              eventType,
		Description:       fmt.Sprintf("%s evolution in language register", eventType),
		Timestamp:         time.Now(),
		Era:               "current",
		Intensity:         0.7,
		AffectedRegisters: []string{registerID},
		Details:           details,
	}

	ree.EvolutionEvents = append(ree.EvolutionEvents, event)
}

// LanguageDeath represents the death or extinction of a language.
type LanguageDeath struct {
	ID          string    `json:"id"`
	LanguageID  string    `json:"languageID"`
	DeathDate   time.Time `json:"deathDate"`
	Era         string    `json:"era"`
	Type        string    `json:"type"`  // natural, forced, gradual, sudden
	Cause       string    `json:"cause"` // conquest, assimilation, migration, disease, etc.
	Description string    `json:"description"`
	Details     string    `json:"details"`

	// Death characteristics
	LastSpeakers      int      `json:"lastSpeakers"`      // Number of speakers at death
	LastRegions       []string `json:"lastRegions"`       // Regions where language was last spoken
	LastDocuments     []string `json:"lastDocuments"`     // Last known documents in the language
	PreservationLevel float32  `json:"preservationLevel"` // How well the language was preserved (0.0-1.0)

	// Impact assessment
	InfluenceOnSurvivors []string `json:"influenceOnSurvivors"` // Languages influenced by the dead language
	LostKnowledge        []string `json:"lostKnowledge"`        // Knowledge lost with the language
	CulturalImpact       string   `json:"culturalImpact"`       // Description of cultural impact

	// Revival potential
	RevivalPotential float32  `json:"revivalPotential"` // Likelihood of revival (0.0-1.0)
	RevivalFactors   []string `json:"revivalFactors"`   // Factors that could aid revival
	RevivalBarriers  []string `json:"revivalBarriers"`  // Barriers to revival

	// Metadata
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Notes     []string  `json:"notes"`
}

// LanguageRevival represents the revival or revitalization of a dead language.
type LanguageRevival struct {
	ID          string    `json:"id"`
	LanguageID  string    `json:"languageID"`
	RevivalDate time.Time `json:"revivalDate"`
	Era         string    `json:"era"`
	Type        string    `json:"type"`    // academic, cultural, religious, practical
	Trigger     string    `json:"trigger"` // cultural_renaissance, academic_interest, etc.
	Description string    `json:"description"`
	Details     string    `json:"details"`

	// Revival characteristics
	RevivalMethod string   `json:"revivalMethod"` // reconstruction, documentation, teaching
	RevivalScope  string   `json:"revivalScope"`  // limited, moderate, comprehensive
	NewSpeakers   int      `json:"newSpeakers"`   // Number of new speakers
	NewRegions    []string `json:"newRegions"`    // Regions where language is revived
	Modernized    bool     `json:"modernized"`    // Whether the language was modernized

	// Revival process
	ReconstructionSteps []string `json:"reconstructionSteps"` // Steps taken to reconstruct the language
	DocumentationUsed   []string `json:"documentationUsed"`   // Sources used for reconstruction
	ModernAdaptations   []string `json:"modernAdaptations"`   // Modern adaptations made

	// Success metrics
	SuccessLevel         float32  `json:"successLevel"`         // How successful the revival was (0.0-1.0)
	CommunitySize        int      `json:"communitySize"`        // Size of the revived language community
	InstitutionalSupport []string `json:"institutionalSupport"` // Institutions supporting the revival

	// Metadata
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Notes     []string  `json:"notes"`
}

// NewLanguageDeath creates a new language death record.
func NewLanguageDeath(
	id, languageID, deathType, cause, description, details string,
	lastSpeakers int,
	lastRegions []string,
) *LanguageDeath {
	now := time.Now()
	return &LanguageDeath{
		ID:                   id,
		LanguageID:           languageID,
		DeathDate:            now,
		Era:                  "current",
		Type:                 deathType,
		Cause:                cause,
		Description:          description,
		Details:              details,
		LastSpeakers:         lastSpeakers,
		LastRegions:          lastRegions,
		LastDocuments:        make([]string, 0),
		PreservationLevel:    0.5,
		InfluenceOnSurvivors: make([]string, 0),
		LostKnowledge:        make([]string, 0),
		CulturalImpact:       "Language extinction resulted in loss of cultural knowledge",
		RevivalPotential:     0.3,
		RevivalFactors:       make([]string, 0),
		RevivalBarriers:      make([]string, 0),
		CreatedAt:            now,
		UpdatedAt:            now,
		Notes:                make([]string, 0),
	}
}

// NewLanguageRevival creates a new language revival record.
func NewLanguageRevival(
	id, languageID, revivalType, trigger, description, details string,
	revivalMethod, revivalScope string,
) *LanguageRevival {
	now := time.Now()
	return &LanguageRevival{
		ID:                   id,
		LanguageID:           languageID,
		RevivalDate:          now,
		Era:                  "current",
		Type:                 revivalType,
		Trigger:              trigger,
		Description:          description,
		Details:              details,
		RevivalMethod:        revivalMethod,
		RevivalScope:         revivalScope,
		NewSpeakers:          0,
		NewRegions:           make([]string, 0),
		Modernized:           false,
		ReconstructionSteps:  make([]string, 0),
		DocumentationUsed:    make([]string, 0),
		ModernAdaptations:    make([]string, 0),
		SuccessLevel:         0.5,
		CommunitySize:        0,
		InstitutionalSupport: make([]string, 0),
		CreatedAt:            now,
		UpdatedAt:            now,
		Notes:                make([]string, 0),
	}
}

// AddLastDocument adds a last known document to the death record.
func (ld *LanguageDeath) AddLastDocument(document string) {
	if !ld.hasLastDocument(document) {
		ld.LastDocuments = append(ld.LastDocuments, document)
		ld.UpdatedAt = time.Now()
		ld.updatePreservationLevel()
	}
}

// hasLastDocument checks if a document exists in the last documents list.
func (ld *LanguageDeath) hasLastDocument(document string) bool {
	for _, doc := range ld.LastDocuments {
		if doc == document {
			return true
		}
	}
	return false
}

// AddInfluenceOnSurvivor adds a language influenced by the dead language.
func (ld *LanguageDeath) AddInfluenceOnSurvivor(languageID string) {
	if !ld.hasInfluenceOnSurvivor(languageID) {
		ld.InfluenceOnSurvivors = append(ld.InfluenceOnSurvivors, languageID)
		ld.UpdatedAt = time.Now()
	}
}

// hasInfluenceOnSurvivor checks if a language is in the influence list.
func (ld *LanguageDeath) hasInfluenceOnSurvivor(languageID string) bool {
	for _, id := range ld.InfluenceOnSurvivors {
		if id == languageID {
			return true
		}
	}
	return false
}

// AddLostKnowledge adds knowledge lost with the language.
func (ld *LanguageDeath) AddLostKnowledge(knowledge string) {
	if !ld.hasLostKnowledge(knowledge) {
		ld.LostKnowledge = append(ld.LostKnowledge, knowledge)
		ld.UpdatedAt = time.Now()
	}
}

// hasLostKnowledge checks if knowledge exists in the lost knowledge list.
func (ld *LanguageDeath) hasLostKnowledge(knowledge string) bool {
	for _, k := range ld.LostKnowledge {
		if k == knowledge {
			return true
		}
	}
	return false
}

// AddRevivalFactor adds a factor that could aid revival.
func (ld *LanguageDeath) AddRevivalFactor(factor string) {
	if !ld.hasRevivalFactor(factor) {
		ld.RevivalFactors = append(ld.RevivalFactors, factor)
		ld.UpdatedAt = time.Now()
		ld.updateRevivalPotential()
	}
}

// hasRevivalFactor checks if a factor exists in the revival factors list.
func (ld *LanguageDeath) hasRevivalFactor(factor string) bool {
	for _, f := range ld.RevivalFactors {
		if f == factor {
			return true
		}
	}
	return false
}

// AddRevivalBarrier adds a barrier to revival.
func (ld *LanguageDeath) AddRevivalBarrier(barrier string) {
	if !ld.hasRevivalBarrier(barrier) {
		ld.RevivalBarriers = append(ld.RevivalBarriers, barrier)
		ld.UpdatedAt = time.Now()
		ld.updateRevivalPotential()
	}
}

// hasRevivalBarrier checks if a barrier exists in the revival barriers list.
func (ld *LanguageDeath) hasRevivalBarrier(barrier string) bool {
	for _, b := range ld.RevivalBarriers {
		if b == barrier {
			return true
		}
	}
	return false
}

// updatePreservationLevel updates the preservation level based on available documentation.
func (ld *LanguageDeath) updatePreservationLevel() {
	// Base preservation level
	baseLevel := float32(0.3)

	// Document preservation bonus
	documentBonus := float32(len(ld.LastDocuments)) * 0.1
	documentBonus = min(0.4, documentBonus)

	// Speaker preservation bonus
	speakerBonus := float32(0.0)
	if ld.LastSpeakers > 0 {
		speakerBonus = min(0.2, float32(ld.LastSpeakers)*0.01)
	}

	// Region preservation bonus
	regionBonus := float32(len(ld.LastRegions)) * 0.05
	regionBonus = min(0.1, regionBonus)

	ld.PreservationLevel = min(1.0, baseLevel+documentBonus+speakerBonus+regionBonus)
}

// updateRevivalPotential updates the revival potential based on factors and barriers.
func (ld *LanguageDeath) updateRevivalPotential() {
	// Base revival potential
	basePotential := float32(0.2)

	// Factor bonus
	factorBonus := float32(len(ld.RevivalFactors)) * 0.1
	factorBonus = min(0.4, factorBonus)

	// Barrier penalty
	barrierPenalty := float32(len(ld.RevivalBarriers)) * 0.08
	barrierPenalty = min(0.3, barrierPenalty)

	// Preservation bonus
	preservationBonus := ld.PreservationLevel * 0.2

	ld.RevivalPotential = max(0.0, min(1.0, basePotential+factorBonus-barrierPenalty+preservationBonus))
}

// AddReconstructionStep adds a reconstruction step to the revival.
func (lr *LanguageRevival) AddReconstructionStep(step string) {
	if !lr.hasReconstructionStep(step) {
		lr.ReconstructionSteps = append(lr.ReconstructionSteps, step)
		lr.UpdatedAt = time.Now()
		lr.updateSuccessLevel()
	}
}

// hasReconstructionStep checks if a step exists in the reconstruction steps list.
func (lr *LanguageRevival) hasReconstructionStep(step string) bool {
	for _, s := range lr.ReconstructionSteps {
		if s == step {
			return true
		}
	}
	return false
}

// AddDocumentationUsed adds documentation used for reconstruction.
func (lr *LanguageRevival) AddDocumentationUsed(document string) {
	if !lr.hasDocumentationUsed(document) {
		lr.DocumentationUsed = append(lr.DocumentationUsed, document)
		lr.UpdatedAt = time.Now()
	}
}

// hasDocumentationUsed checks if documentation exists in the used documents list.
func (lr *LanguageRevival) hasDocumentationUsed(document string) bool {
	for _, doc := range lr.DocumentationUsed {
		if doc == document {
			return true
		}
	}
	return false
}

// AddModernAdaptation adds a modern adaptation made to the language.
func (lr *LanguageRevival) AddModernAdaptation(adaptation string) {
	if !lr.hasModernAdaptation(adaptation) {
		lr.ModernAdaptations = append(lr.ModernAdaptations, adaptation)
		lr.UpdatedAt = time.Now()
		lr.Modernized = true
		lr.updateSuccessLevel()
	}
}

// hasModernAdaptation checks if an adaptation exists in the modern adaptations list.
func (lr *LanguageRevival) hasModernAdaptation(adaptation string) bool {
	for _, a := range lr.ModernAdaptations {
		if a == adaptation {
			return true
		}
	}
	return false
}

// SetCommunitySize sets the size of the revived language community.
func (lr *LanguageRevival) SetCommunitySize(size int) {
	lr.CommunitySize = size
	lr.UpdatedAt = time.Now()
	lr.updateSuccessLevel()
}

// AddInstitutionalSupport adds institutional support for the revival.
func (lr *LanguageRevival) AddInstitutionalSupport(institution string) {
	if !lr.hasInstitutionalSupport(institution) {
		lr.InstitutionalSupport = append(lr.InstitutionalSupport, institution)
		lr.UpdatedAt = time.Now()
		lr.updateSuccessLevel()
	}
}

// hasInstitutionalSupport checks if an institution exists in the support list.
func (lr *LanguageRevival) hasInstitutionalSupport(institution string) bool {
	for _, inst := range lr.InstitutionalSupport {
		if inst == institution {
			return true
		}
	}
	return false
}

// updateSuccessLevel updates the success level based on various factors.
func (lr *LanguageRevival) updateSuccessLevel() {
	// Base success level
	baseLevel := float32(0.3)

	// Reconstruction steps bonus
	stepBonus := float32(len(lr.ReconstructionSteps)) * 0.08
	stepBonus = min(0.2, stepBonus)

	// Documentation bonus
	docBonus := float32(len(lr.DocumentationUsed)) * 0.05
	docBonus = min(0.15, docBonus)

	// Community size bonus
	communityBonus := float32(0.0)
	if lr.CommunitySize > 0 {
		communityBonus = min(0.2, float32(lr.CommunitySize)*0.001)
	}

	// Institutional support bonus
	institutionBonus := float32(len(lr.InstitutionalSupport)) * 0.1
	institutionBonus = min(0.15, institutionBonus)

	lr.SuccessLevel = min(1.0, baseLevel+stepBonus+docBonus+communityBonus+institutionBonus)
}

// LanguageDeathAndRevivalEngine manages the death and revival of languages.
type LanguageDeathAndRevivalEngine struct {
	// Death parameters
	NaturalDeathRate float32 `json:"naturalDeathRate"` // Rate of natural language death
	ForcedDeathRate  float32 `json:"forcedDeathRate"`  // Rate of forced language death
	GradualDeathRate float32 `json:"gradualDeathRate"` // Rate of gradual language death
	SuddenDeathRate  float32 `json:"suddenDeathRate"`  // Rate of sudden language death

	// Revival parameters
	AcademicRevivalRate  float32 `json:"academicRevivalRate"`  // Rate of academic revivals
	CulturalRevivalRate  float32 `json:"culturalRevivalRate"`  // Rate of cultural revivals
	ReligiousRevivalRate float32 `json:"religiousRevivalRate"` // Rate of religious revivals
	PracticalRevivalRate float32 `json:"practicalRevivalRate"` // Rate of practical revivals

	// Trigger weights
	ConquestWeight            float32 `json:"conquestWeight"`            // Weight of conquest as death cause
	AssimilationWeight        float32 `json:"assimilationWeight"`        // Weight of assimilation as death cause
	MigrationWeight           float32 `json:"migrationWeight"`           // Weight of migration as death cause
	DiseaseWeight             float32 `json:"diseaseWeight"`             // Weight of disease as death cause
	CulturalRenaissanceWeight float32 `json:"culturalRenaissanceWeight"` // Weight of cultural renaissance as revival trigger
	AcademicInterestWeight    float32 `json:"academicInterestWeight"`    // Weight of academic interest as revival trigger

	// History tracking
	DeathEvents   []LanguageDeathEvent   `json:"deathEvents"`
	RevivalEvents []LanguageRevivalEvent `json:"revivalEvents"`
}

// LanguageDeathEvent represents a significant death event affecting languages.
type LanguageDeathEvent struct {
	ID                string    `json:"id"`
	Type              string    `json:"type"` // conquest, assimilation, migration, disease, etc.
	Description       string    `json:"description"`
	Timestamp         time.Time `json:"timestamp"`
	Era               string    `json:"era"`
	Intensity         float32   `json:"intensity"`         // 0.0-1.0
	AffectedLanguages []string  `json:"affectedLanguages"` // IDs of languages affected
	Details           string    `json:"details"`
}

// LanguageRevivalEvent represents a significant revival event affecting languages.
type LanguageRevivalEvent struct {
	ID                string    `json:"id"`
	Type              string    `json:"type"` // academic, cultural, religious, practical
	Description       string    `json:"description"`
	Timestamp         time.Time `json:"timestamp"`
	Era               string    `json:"era"`
	Intensity         float32   `json:"intensity"`         // 0.0-1.0
	AffectedLanguages []string  `json:"affectedLanguages"` // IDs of languages affected
	Details           string    `json:"details"`
}

// NewLanguageDeathAndRevivalEngine creates a new language death and revival engine.
func NewLanguageDeathAndRevivalEngine() *LanguageDeathAndRevivalEngine {
	return &LanguageDeathAndRevivalEngine{
		NaturalDeathRate:          0.05,
		ForcedDeathRate:           0.03,
		GradualDeathRate:          0.04,
		SuddenDeathRate:           0.02,
		AcademicRevivalRate:       0.02,
		CulturalRevivalRate:       0.03,
		ReligiousRevivalRate:      0.01,
		PracticalRevivalRate:      0.02,
		ConquestWeight:            0.4,
		AssimilationWeight:        0.3,
		MigrationWeight:           0.2,
		DiseaseWeight:             0.1,
		CulturalRenaissanceWeight: 0.4,
		AcademicInterestWeight:    0.3,
		DeathEvents:               make([]LanguageDeathEvent, 0),
		RevivalEvents:             make([]LanguageRevivalEvent, 0),
	}
}

// SimulateLanguageDeath simulates the death of a language based on various factors.
func (ldre *LanguageDeathAndRevivalEngine) SimulateLanguageDeath(
	language *Language,
	deathType, cause, description, details string,
	lastSpeakers int,
	lastRegions []string,
) (*LanguageDeath, error) {
	// Create death record
	death := NewLanguageDeath(
		fmt.Sprintf("death_%s_%d", language.ID.String(), len(ldre.DeathEvents)+1),
		language.ID.String(),
		deathType,
		cause,
		description,
		details,
		lastSpeakers,
		lastRegions,
	)

	// Note: Language status and extinction date would be updated in a full implementation
	// For now, we just create the death record

	// Add death event
	ldre.recordDeathEvent(deathType, language.ID.String(), cause, details)

	return death, nil
}

// SimulateLanguageRevival simulates the revival of a dead language.
func (ldre *LanguageDeathAndRevivalEngine) SimulateLanguageRevival(
	language *Language,
	revivalType, trigger, description, details string,
	revivalMethod, revivalScope string,
) (*LanguageRevival, error) {
	// Note: In a full implementation, we would check if the language is actually dead
	// For now, we just create the revival record

	// Create revival record
	revival := NewLanguageRevival(
		fmt.Sprintf("revival_%s_%d", language.ID.String(), len(ldre.RevivalEvents)+1),
		language.ID.String(),
		revivalType,
		trigger,
		description,
		details,
		revivalMethod,
		revivalScope,
	)

	// Note: Language status and revival date would be updated in a full implementation
	// For now, we just create the revival record

	// Add revival event
	ldre.recordRevivalEvent(revivalType, language.ID.String(), trigger, details)

	return revival, nil
}

// AssessRevivalPotential assesses the potential for reviving a dead language.
func (ldre *LanguageDeathAndRevivalEngine) AssessRevivalPotential(
	language *Language,
	death *LanguageDeath,
) float32 {
	// Base potential
	potential := death.RevivalPotential

	// Language complexity factor (simpler languages are easier to revive)
	// Note: Using a default complexity value since Language struct doesn't have this field
	complexityFactor := float32(0.5) // Default complexity
	potential += (1.0 - complexityFactor) * 0.2

	// Documentation factor
	documentationFactor := float32(len(death.LastDocuments)) * 0.05
	potential += min(0.2, documentationFactor)

	// Cultural significance factor
	// Note: Using a default value since Language struct doesn't have this field
	culturalSignificance := float32(0.5) // Default cultural significance
	if culturalSignificance > 0.7 {
		potential += 0.1
	}

	// Geographic distribution factor
	distributionFactor := float32(len(death.LastRegions)) * 0.03
	potential += min(0.1, distributionFactor)

	return min(1.0, potential)
}

// recordDeathEvent records a significant death event.
func (ldre *LanguageDeathAndRevivalEngine) recordDeathEvent(
	eventType, languageID string,
	cause, details string,
) {
	event := LanguageDeathEvent{
		ID:                fmt.Sprintf("death_event_%s_%d", eventType, len(ldre.DeathEvents)+1),
		Type:              eventType,
		Description:       fmt.Sprintf("%s death of language", eventType),
		Timestamp:         time.Now(),
		Era:               "current",
		Intensity:         0.8,
		AffectedLanguages: []string{languageID},
		Details:           details,
	}

	ldre.DeathEvents = append(ldre.DeathEvents, event)
}

// recordRevivalEvent records a significant revival event.
func (ldre *LanguageDeathAndRevivalEngine) recordRevivalEvent(
	eventType, languageID string,
	trigger, details string,
) {
	event := LanguageRevivalEvent{
		ID:                fmt.Sprintf("revival_event_%s_%d", eventType, len(ldre.RevivalEvents)+1),
		Type:              eventType,
		Description:       fmt.Sprintf("%s revival of language", eventType),
		Timestamp:         time.Now(),
		Era:               "current",
		Intensity:         0.7,
		AffectedLanguages: []string{languageID},
		Details:           details,
	}

	ldre.RevivalEvents = append(ldre.RevivalEvents, event)
}

// PidginAndCreoleFormationEngine manages the formation and evolution of pidgins and creoles.
type PidginAndCreoleFormationEngine struct {
	// Formation parameters
	PidginFormationRate float32 `json:"pidginFormationRate"` // Rate of pidgin formation
	CreoleFormationRate float32 `json:"creoleFormationRate"` // Rate of creole formation from pidgins

	// Contact type weights
	TradeContactWeight    float32 `json:"tradeContactWeight"`    // Weight of trade contact
	WorkContactWeight     float32 `json:"workContactWeight"`     // Weight of work contact
	SocialContactWeight   float32 `json:"socialContactWeight"`   // Weight of social contact
	MilitaryContactWeight float32 `json:"militaryContactWeight"` // Weight of military contact

	// Evolution parameters
	PidginEvolutionRate float32 `json:"pidginEvolutionRate"` // Rate of pidgin evolution
	CreoleEvolutionRate float32 `json:"creoleEvolutionRate"` // Rate of creole evolution

	// History tracking
	FormationEvents []PidginCreoleFormationEvent `json:"formationEvents"`
}

// PidginCreoleFormationEvent represents a significant formation event.
type PidginCreoleFormationEvent struct {
	ID                string    `json:"id"`
	Type              string    `json:"type"` // pidgin_formation, creole_formation
	Description       string    `json:"description"`
	Timestamp         time.Time `json:"timestamp"`
	Era               string    `json:"era"`
	Intensity         float32   `json:"intensity"`         // 0.0-1.0
	AffectedLanguages []string  `json:"affectedLanguages"` // IDs of languages affected
	Details           string    `json:"details"`
}

// NewPidginAndCreoleFormationEngine creates a new pidgin and creole formation engine.
func NewPidginAndCreoleFormationEngine() *PidginAndCreoleFormationEngine {
	return &PidginAndCreoleFormationEngine{
		PidginFormationRate:   0.3,
		CreoleFormationRate:   0.1,
		TradeContactWeight:    0.4,
		WorkContactWeight:     0.3,
		SocialContactWeight:   0.2,
		MilitaryContactWeight: 0.1,
		PidginEvolutionRate:   0.2,
		CreoleEvolutionRate:   0.15,
		FormationEvents:       make([]PidginCreoleFormationEvent, 0),
	}
}

// SimulatePidginFormation simulates the formation of a pidgin language.
func (pcfe *PidginAndCreoleFormationEngine) SimulatePidginFormation(
	sourceLanguages []string,
	primarySource, secondarySource string,
	pidginType, context string,
	era string,
) (*PidginLanguage, error) {
	if len(sourceLanguages) < 2 {
		return nil, fmt.Errorf("need at least 2 source languages for pidgin formation")
	}

	// Create pidgin
	pidgin := NewPidginLanguage(
		fmt.Sprintf("pidgin_%s_%s_%d", primarySource, secondarySource, len(pcfe.FormationEvents)+1),
		fmt.Sprintf("%s-%s Pidgin", primarySource, secondarySource),
		fmt.Sprintf("Pidgin formed between %s and %s speakers in %s context", primarySource, secondarySource, context),
		era,
		pidginType,
		context,
		sourceLanguages,
		primarySource,
		secondarySource,
	)

	// Add basic features based on source languages
	pidgin.AddPhonologicalFeature("simplified_phonology")
	pidgin.AddMorphologicalFeature("reduced_morphology")
	pidgin.AddSyntacticFeature("basic_syntax")
	pidgin.AddLexicalFeature("core_vocabulary")

	// Add usage information
	pidgin.AddUsageDomain(pidginType)
	pidgin.AddUsageRegion("contact_zone")

	// Record formation event
	pcfe.recordFormationEvent("pidgin_formation", pidgin.ID, sourceLanguages,
		fmt.Sprintf("Pidgin formed between %s and %s in %s context", primarySource, secondarySource, context))

	return pidgin, nil
}

// SimulateCreoleFormation simulates the formation of a creole from a pidgin.
func (pcfe *PidginAndCreoleFormationEngine) SimulateCreoleFormation(
	pidgin *PidginLanguage,
	creoleType, context string,
	era string,
) (*CreoleLanguage, error) {
	// Create creole
	creole := NewCreoleLanguage(
		fmt.Sprintf("creole_%s_%d", pidgin.ID, len(pcfe.FormationEvents)+1),
		fmt.Sprintf("%s Creole", pidgin.Name),
		fmt.Sprintf("Creole developed from %s in %s context", pidgin.Name, context),
		era,
		creoleType,
		context,
		pidgin.ID,
		pidgin.SourceLanguages,
		pidgin.PrimarySource,
		pidgin.SecondarySource,
	)

	// Inherit features from pidgin
	for _, feature := range pidgin.PhonologicalFeatures {
		creole.AddPhonologicalFeature(feature)
	}
	for _, feature := range pidgin.MorphologicalFeatures {
		creole.AddMorphologicalFeature(feature)
	}
	for _, feature := range pidgin.SyntacticFeatures {
		creole.AddSyntacticFeature(feature)
	}
	for _, feature := range pidgin.LexicalFeatures {
		creole.AddLexicalFeature(feature)
	}

	// Add creole-specific features
	creole.AddPhonologicalFeature("developed_phonology")
	creole.AddMorphologicalFeature("developed_morphology")
	creole.AddSyntacticFeature("developed_syntax")
	creole.AddLexicalFeature("expanded_vocabulary")

	// Add usage information
	creole.AddUsageDomain(creoleType)
	creole.AddUsageRegion("creole_community")

	// Record formation event
	pcfe.recordFormationEvent("creole_formation", creole.ID, []string{pidgin.ID},
		fmt.Sprintf("Creole developed from %s in %s context", pidgin.Name, context))

	return creole, nil
}

// AssessCreolizationPotential assesses the potential for a pidgin to become a creole.
func (pcfe *PidginAndCreoleFormationEngine) AssessCreolizationPotential(pidgin *PidginLanguage) float32 {
	// Base potential
	potential := float32(0.3)

	// Stability factor (more stable pidgins are more likely to creolize)
	stabilityFactor := pidgin.StabilityLevel * 0.3
	potential += stabilityFactor

	// Speaker count factor (more speakers = higher potential)
	if pidgin.SpeakerCount > 100 {
		potential += 0.2
	} else if pidgin.SpeakerCount > 50 {
		potential += 0.1
	}

	// Generation factor (more generations = higher potential)
	generationFactor := float32(pidgin.GenerationCount) * 0.1
	potential += min(0.2, generationFactor)

	// Functionality factor (more functional pidgins are more likely to creolize)
	functionalityFactor := pidgin.FunctionalityLevel * 0.2
	potential += functionalityFactor

	return min(1.0, potential)
}

// recordFormationEvent records a significant formation event.
func (pcfe *PidginAndCreoleFormationEngine) recordFormationEvent(
	eventType, languageID string,
	affectedLanguages []string,
	details string,
) {
	event := PidginCreoleFormationEvent{
		ID:                fmt.Sprintf("formation_event_%s_%d", eventType, len(pcfe.FormationEvents)+1),
		Type:              eventType,
		Description:       fmt.Sprintf("%s occurred", eventType),
		Timestamp:         time.Now(),
		Era:               "current",
		Intensity:         0.8,
		AffectedLanguages: affectedLanguages,
		Details:           details,
	}

	pcfe.FormationEvents = append(pcfe.FormationEvents, event)
}

// PidginLanguage represents a simplified contact language that develops between speakers of different languages.
type PidginLanguage struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Era         string `json:"era"`
	Type        string `json:"type"`    // trade, work, social, military, etc.
	Context     string `json:"context"` // marketplace, plantation, port, etc.

	// Source languages
	SourceLanguages []string `json:"sourceLanguages"` // IDs of languages that contributed to the pidgin
	PrimarySource   string   `json:"primarySource"`   // ID of the dominant source language
	SecondarySource string   `json:"secondarySource"` // ID of the secondary source language

	// Linguistic characteristics
	PhonologicalFeatures  []string `json:"phonologicalFeatures"`  // Simplified sound patterns
	MorphologicalFeatures []string `json:"morphologicalFeatures"` // Simplified word formation
	SyntacticFeatures     []string `json:"syntacticFeatures"`     // Simplified sentence structure
	LexicalFeatures       []string `json:"lexicalFeatures"`       // Simplified vocabulary

	// Development characteristics
	ComplexityLevel    float32 `json:"complexityLevel"`    // 0.0-1.0 (simplified to complex)
	StabilityLevel     float32 `json:"stabilityLevel"`     // 0.0-1.0 (unstable to stable)
	FunctionalityLevel float32 `json:"functionalityLevel"` // 0.0-1.0 (limited to comprehensive)

	// Usage patterns
	SpeakerCount    int      `json:"speakerCount"`    // Number of speakers
	UsageRegions    []string `json:"usageRegions"`    // Regions where pidgin is used
	UsageDomains    []string `json:"usageDomains"`    // Domains where pidgin is used (trade, work, etc.)
	GenerationCount int      `json:"generationCount"` // Number of generations using the pidgin

	// Evolution tracking
	EvolutionHistory []PidginChange `json:"evolutionHistory"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	Notes            []string       `json:"notes"`
}

// PidginChange represents a change in a pidgin language over time.
type PidginChange struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // phonological, morphological, syntactic, lexical, functional
	Description string    `json:"description"`
	Details     string    `json:"details"`
	Timestamp   time.Time `json:"timestamp"`
	Era         string    `json:"era"`
	Trigger     string    `json:"trigger"`    // increased_contact, generational_change, etc.
	Intensity   float32   `json:"intensity"`  // 0.0-1.0
	Confidence  float32   `json:"confidence"` // 0.0-1.0

	// Specific changes
	AddedFeatures    []string `json:"addedFeatures,omitempty"`
	RemovedFeatures  []string `json:"removedFeatures,omitempty"`
	ModifiedFeatures []string `json:"modifiedFeatures,omitempty"`
}
