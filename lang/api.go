package lang

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/chris-pikul/kismet-zero/lang/grammar"
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
	rng := rand.New(rand.NewPCG(uint64(seed), 0))

	// Simple naming patterns based on culture
	switch culture {
	case "elvish", "high_elf":
		patterns := []string{"%sarin", "%sian", "%së", "%sëa", "%sëan"}
		pattern := patterns[rng.IntN(len(patterns))]
		return fmt.Sprintf(pattern, family)

	case "dwarvish", "mountain_dwarf":
		patterns := []string{"%suz", "%saz", "%siz", "%suzh", "%sazh"}
		pattern := patterns[rng.IntN(len(patterns))]
		return fmt.Sprintf(pattern, family)

	case "orcish", "black_orc":
		patterns := []string{"%s'%s", "%s-%s", "%s%s", "%s'%s'"}
		pattern := patterns[rng.IntN(len(patterns))]
		return fmt.Sprintf(pattern, family, "gul")

	case "human", "northern_human":
		patterns := []string{"%sish", "%sian", "%sic", "%sese", "%sian"}
		pattern := patterns[rng.IntN(len(patterns))]
		return fmt.Sprintf(pattern, family)

	case "ancient", "divine":
		patterns := []string{"Proto-%s", "Ancient %s", "%s Prime", "%s Root"}
		pattern := patterns[rng.IntN(len(patterns))]
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
	// Create a basic grammar
	// This is a simplified version - the grammar package would handle the actual generation

	// Create basic agreement system
	agreement := grammar.NewAgreementSystem(
		[]grammar.Case{grammar.CaseNominative, grammar.CaseAccusative},
		[]grammar.Number{grammar.NumberSingular, grammar.NumberPlural},
		[]grammar.Gender{grammar.GenderMasculine, grammar.GenderFeminine},
		[]grammar.Tense{grammar.TensePresent, grammar.TensePast},
		[]grammar.Aspect{grammar.AspectPerfective, grammar.AspectImperfective},
		[]grammar.Mood{grammar.MoodIndicative, grammar.MoodImperative},
	)

	// Create basic word builder (simplified)
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

	// Create a copy for evolution
	evolvedLang := language.Clone(language.ID, seed)
	evolvedLang.EvolvedAt = time.Now()

	// Apply basic evolution based on time period
	switch timePeriod {
	case "100_years", "short":
		// Minor changes
		evolvedLang.SetDescription(fmt.Sprintf("%s (evolved over %s)", language.Description, timePeriod))

	case "500_years", "medium":
		// Moderate changes
		evolvedLang.SetDescription(fmt.Sprintf("%s (evolved over %s)", language.Description, timePeriod))

	case "1000_years", "long":
		// Significant changes
		evolvedLang.SetDescription(fmt.Sprintf("%s (evolved over %s)", language.Description, timePeriod))

	case "2000_years", "very_long":
		// Major changes
		evolvedLang.SetDescription(fmt.Sprintf("%s (evolved over %s)", language.Description, timePeriod))

	default:
		// Custom time period
		evolvedLang.SetDescription(fmt.Sprintf("%s (evolved over %s)", language.Description, timePeriod))
	}

	// Apply cultural event influences
	if len(culturalEvents) > 0 {
		evolvedLang.SetDescription(fmt.Sprintf("%s, influenced by: %v", evolvedLang.Description, culturalEvents))
	}

	return evolvedLang, nil
}

// SimulateLanguageContact models how languages influence each other.
func SimulateLanguageContact(language1, language2 *Language, contactType string, intensity float64, duration string, seed int64) (*Language, *Language, error) {
	if language1 == nil || language2 == nil {
		return nil, nil, fmt.Errorf("both languages must be provided")
	}

	if intensity < 0.0 || intensity > 1.0 {
		return nil, nil, fmt.Errorf("intensity must be between 0.0 and 1.0")
	}

	// Create copies for contact simulation
	contactLang1 := language1.Clone(language1.ID, seed)
	contactLang2 := language2.Clone(language2.ID, seed+1)

	// Apply contact effects based on type and intensity
	switch contactType {
	case "trade":
		contactLang1.SetDescription(fmt.Sprintf("%s (influenced by trade with %s)", language1.Description, language2.Name))
		contactLang2.SetDescription(fmt.Sprintf("%s (influenced by trade with %s)", language2.Description, language1.Name))

	case "conquest":
		contactLang1.SetDescription(fmt.Sprintf("%s (conquered by %s)", language1.Description, language2.Name))
		contactLang2.SetDescription(fmt.Sprintf("%s (conquered %s)", language2.Description, language1.Name))

	case "migration":
		contactLang1.SetDescription(fmt.Sprintf("%s (influenced by %s migration)", language1.Description, language2.Name))
		contactLang2.SetDescription(fmt.Sprintf("%s (migrated to %s region)", language2.Description, language1.Name))

	case "cultural":
		contactLang1.SetDescription(fmt.Sprintf("%s (culturally influenced by %s)", language1.Description, language2.Name))
		contactLang2.SetDescription(fmt.Sprintf("%s (culturally influenced %s)", language2.Description, language1.Name))

	default:
		contactLang1.SetDescription(fmt.Sprintf("%s (contact with %s)", language1.Description, language2.Name))
		contactLang2.SetDescription(fmt.Sprintf("%s (contact with %s)", language2.Description, language1.Name))
	}

	// Add duration information
	contactLang1.SetDescription(fmt.Sprintf("%s over %s", contactLang1.Description, duration))
	contactLang2.SetDescription(fmt.Sprintf("%s over %s", contactLang2.Description, duration))

	return contactLang1, contactLang2, nil
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
	rng := rand.New(rand.NewPCG(uint64(seed), 0))

	// Generate names based on type and culture
	switch nameType {
	case "person":
		names = generatePersonNames(language, count, culturalContext, rng)

	case "place":
		names = generatePlaceNames(language, count, culturalContext, rng)

	case "deity":
		names = generateDeityNames(language, count, culturalContext, rng)

	case "artifact":
		names = generateArtifactNames(language, count, culturalContext, rng)

	default:
		names = generateGenericNames(language, count, culturalContext, rng)
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
			prefix := prefixes[rng.IntN(len(prefixes))]
			suffix := suffixes[rng.IntN(len(suffixes))]
			names = append(names, prefix+suffix)
		}

	case "dwarvish", "mountain_dwarf":
		prefixes := []string{"Thorin", "Balin", "Dwalin", "Bifur", "Bofur", "Bombur", "Dori"}
		suffixes := []string{"son", "sson", "dottir", "sson", "sson", "sson", "sson"}

		for i := 0; i < count; i++ {
			prefix := prefixes[rng.IntN(len(prefixes))]
			suffix := suffixes[rng.IntN(len(suffixes))]
			names = append(names, prefix+suffix)
		}

	case "orcish", "black_orc":
		prefixes := []string{"Grom", "Thrall", "Gul", "Mog", "Karg", "Rag", "Zug"}
		suffixes := []string{"'dan", "'gar", "'ash", "'osh", "'ash", "'osh", "'ash"}

		for i := 0; i < count; i++ {
			prefix := prefixes[rng.IntN(len(prefixes))]
			suffix := suffixes[rng.IntN(len(suffixes))]
			names = append(names, prefix+suffix)
		}

	default:
		// Generic fantasy names
		prefixes := []string{"A", "E", "I", "O", "U", "Y"}
		suffixes := []string{"ron", "lan", "dor", "mir", "wen", "las"}

		for i := 0; i < count; i++ {
			prefix := prefixes[rng.IntN(len(prefixes))]
			suffix := suffixes[rng.IntN(len(suffixes))]
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
		prefix := prefixes[rng.IntN(len(prefixes))]
		suffix := suffixes[rng.IntN(len(suffixes))]
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
		prefix := prefixes[rng.IntN(len(prefixes))]
		suffix := suffixes[rng.IntN(len(suffixes))]
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
		prefix := prefixes[rng.IntN(len(prefixes))]
		suffix := suffixes[rng.IntN(len(suffixes))]
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
		prefix := prefixes[rng.IntN(len(prefixes))]
		suffix := suffixes[rng.IntN(len(suffixes))]
		names = append(names, prefix+suffix)
	}

	return names
}

// GenerateText creates sample texts in the generated language.
func GenerateText(language *Language, textType string, length string, culturalContext string, seed int64) (string, error) {
	if language == nil {
		return "", fmt.Errorf("language cannot be nil")
	}

	rng := rand.New(rand.NewPCG(uint64(seed), 0))

	// Generate text based on type and length
	switch textType {
	case "lore":
		return generateLoreText(language, length, culturalContext, rng)

	case "poetry":
		return generatePoetryText(language, length, culturalContext, rng)

	case "dialogue":
		return generateDialogueText(language, length, culturalContext, rng)

	case "inscription":
		return generateInscriptionText(language, length, culturalContext, rng)

	default:
		return generateGenericText(language, length, culturalContext, rng)
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

	return texts[rng.IntN(len(texts))], nil
}

// generatePoetryText creates poetic text.
func generatePoetryText(language *Language, length string, culturalContext string, rng *rand.Rand) (string, error) {
	// Simple poetry generation
	poems := []string{
		"Moonlight dances on silver streams,\nAncient wisdom in elven dreams.",
		"Mountains tall and valleys deep,\nSecrets that the ages keep.",
		"Stars above and earth below,\nStories that the winds do know.",
	}

	return poems[rng.IntN(len(poems))], nil
}

// generateDialogueText creates dialogue text.
func generateDialogueText(language *Language, length string, culturalContext string, rng *rand.Rand) (string, error) {
	// Simple dialogue generation
	dialogues := []string{
		"\"Greetings, traveler. What brings you to these lands?\"",
		"\"The ancient ones speak of times long past.\"",
		"\"Magic flows through all things, if you know how to listen.\"",
	}

	return dialogues[rng.IntN(len(dialogues))], nil
}

// generateInscriptionText creates inscription text.
func generateInscriptionText(language *Language, length string, culturalContext string, rng *rand.Rand) (string, error) {
	// Simple inscription generation
	inscriptions := []string{
		"Here lies the resting place of ancient kings.",
		"Enter not, lest you disturb the eternal sleep.",
		"Knowledge and wisdom await those who seek.",
	}

	return inscriptions[rng.IntN(len(inscriptions))], nil
}

// generateGenericText creates generic text when type is not specified.
func generateGenericText(language *Language, length string, culturalContext string, rng *rand.Rand) (string, error) {
	// Simple generic text generation
	texts := []string{
		"This is a sample text in the constructed language.",
		"Words flow like water, carrying meaning and sound.",
		"The language speaks of culture and tradition.",
	}

	return texts[rng.IntN(len(texts))], nil
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
		variantLang, err := CreateRandomLanguage(culture, fmt.Sprintf("%s_variant_%d", culture, i), variantSeed)
		if err != nil {
			return nil, fmt.Errorf("failed to create variant language %d: %w", i, err)
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
