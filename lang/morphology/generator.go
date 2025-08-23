package morphology

import (
	"fmt"
	"math/rand/v2"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
	"github.com/chris-pikul/kismet-zero/lang/phonology"
)

// MorphologyGenerator provides a high-level API for generating morphemes,
// words, and names using the morphology package components.
type MorphologyGenerator struct {
	phonology       *phonology.Phonology
	morphemeBuilder *MorphemeBuilder
	wordBuilder     *WordBuilder
	nameGenerator   *NameGenerator
	ruleSet         *RuleSet
	lexicon         *Lexicon
}

// NewMorphologyGenerator creates a new morphology generator with the given phonology.
func NewMorphologyGenerator(phonology *phonology.Phonology) *MorphologyGenerator {
	// Create rule set with some default rules
	ruleSet := NewRuleSet()

	// Add some common derivational rules
	ruleSet.AddDerivationalRule(NewDerivationalRule(
		"agent_noun",
		MorphemeTypeRoot,
		MorphemeTypeSuffix,
		"one who does",
		1.0,
	))

	ruleSet.AddDerivationalRule(NewDerivationalRule(
		"abstract_noun",
		MorphemeTypeRoot,
		MorphemeTypeSuffix,
		"quality of",
		0.8,
	))

	// Add some common inflectional rules
	ruleSet.AddInflectionalRule(NewInflectionalRule(
		"plural",
		WordCategoryNoun,
		[]string{"plural"},
		"many",
		1.0,
	))

	ruleSet.AddInflectionalRule(NewInflectionalRule(
		"past_tense",
		WordCategoryVerb,
		[]string{"past"},
		"completed",
		1.0,
	))

	// Create lexicon for the culture
	lexicon := NewLexicon("default")

	// Create components
	morphemeBuilder := NewMorphemeBuilder(phonology)
	wordBuilder := NewWordBuilder(phonology, ruleSet)
	nameGenerator := NewNameGenerator(phonology, lexicon)

	return &MorphologyGenerator{
		phonology:       phonology,
		morphemeBuilder: morphemeBuilder,
		wordBuilder:     wordBuilder,
		nameGenerator:   nameGenerator,
		ruleSet:         ruleSet,
		lexicon:         lexicon,
	}
}

// SetCulture sets the culture for the generator and updates the lexicon.
func (mg *MorphologyGenerator) SetCulture(culture string) {
	mg.lexicon = NewLexicon(culture)
	mg.nameGenerator = NewNameGenerator(mg.phonology, mg.lexicon)
}

// GetCulture returns the current culture of the generator.
func (mg *MorphologyGenerator) GetCulture() string {
	return mg.lexicon.GetCulture()
}

// GenerateRootMorpheme creates a root morpheme with the given meaning.
func (mg *MorphologyGenerator) GenerateRootMorpheme(
	meaning string,
	rng *rand.Rand,
) (*Morpheme, error) {
	culture := mg.GetCulture()
	return mg.morphemeBuilder.BuildRootMorpheme(meaning, culture, rng)
}

// GenerateAffixMorpheme creates an affix morpheme of the specified type.
func (mg *MorphologyGenerator) GenerateAffixMorpheme(
	morphemeType MorphemeType,
	meaning string,
	rng *rand.Rand,
) (*Morpheme, error) {
	culture := mg.GetCulture()
	return mg.morphemeBuilder.BuildAffixMorpheme(morphemeType, meaning, culture, rng)
}

// GenerateWord creates a word from a root morpheme with the specified category.
func (mg *MorphologyGenerator) GenerateWord(
	root *Morpheme,
	category WordCategory,
	rng *rand.Rand,
) (*Word, error) {
	culture := mg.GetCulture()
	word, err := mg.wordBuilder.BuildWord(root, category, culture, rng)
	if err != nil {
		return nil, err
	}

	// Add the word to the lexicon
	if err := mg.lexicon.AddWord(word); err != nil {
		return nil, fmt.Errorf("failed to add word to lexicon: %w", err)
	}

	return word, nil
}

// GeneratePersonalName creates a personal name for the current culture.
func (mg *MorphologyGenerator) GeneratePersonalName(
	gender string,
	formality string,
	rng *rand.Rand,
) (*Name, error) {
	culture := mg.GetCulture()
	return mg.nameGenerator.GeneratePersonalName(culture, gender, formality, rng)
}

// GeneratePlaceName creates a place name for the current culture.
func (mg *MorphologyGenerator) GeneratePlaceName(
	placeType string,
	rng *rand.Rand,
) (*Name, error) {
	culture := mg.GetCulture()
	return mg.nameGenerator.GeneratePlaceName(culture, placeType, rng)
}

// GenerateTitle creates a title for the current culture.
func (mg *MorphologyGenerator) GenerateTitle(
	rank string,
	rng *rand.Rand,
) (*Name, error) {
	culture := mg.GetCulture()
	return mg.nameGenerator.GenerateTitle(culture, rank, rng)
}

// GenerateDeityName creates a deity name for the current culture.
func (mg *MorphologyGenerator) GenerateDeityName(
	domain string,
	rng *rand.Rand,
) (*Name, error) {
	culture := mg.GetCulture()
	return mg.nameGenerator.GenerateDeityName(culture, domain, rng)
}

// GenerateLexicon generates a basic lexicon for the current culture with
// common words in various categories.
func (mg *MorphologyGenerator) GenerateLexicon(
	numWords int,
	rng *rand.Rand,
) error {
	if numWords <= 0 {
		return fmt.Errorf("number of words must be positive")
	}

	// Common meanings for different word categories
	meanings := map[WordCategory][]string{
		WordCategoryNoun: {
			"person", "water", "fire", "earth", "air", "sun", "moon", "star",
			"tree", "stone", "metal", "food", "house", "road", "mountain",
		},
		WordCategoryVerb: {
			"walk", "run", "eat", "drink", "sleep", "speak", "hear", "see",
			"think", "feel", "love", "hate", "build", "destroy", "create",
		},
		WordCategoryAdjective: {
			"big", "small", "hot", "cold", "fast", "slow", "strong", "weak",
			"bright", "dark", "good", "bad", "new", "old", "beautiful",
		},
	}

	// Generate words for each category
	for category, categoryMeanings := range meanings {
		wordsPerCategory := numWords / len(meanings)
		if wordsPerCategory == 0 {
			wordsPerCategory = 1
		}

		for i := 0; i < wordsPerCategory && i < len(categoryMeanings); i++ {
			meaning := categoryMeanings[i]

			// Generate root morpheme
			root, err := mg.GenerateRootMorpheme(meaning, rng)
			if err != nil {
				continue // Skip this word if generation fails
			}

			// Generate word
			word, err := mg.GenerateWord(root, category, rng)
			if err != nil {
				continue // Skip this word if generation fails
			}

			// Word is automatically added to lexicon by GenerateWord
			_ = word // Use the word to avoid unused variable warning
		}
	}

	return nil
}

// GetLexicon returns the current lexicon.
func (mg *MorphologyGenerator) GetLexicon() *Lexicon {
	return mg.lexicon
}

// GetRuleSet returns the current rule set.
func (mg *MorphologyGenerator) GetRuleSet() *RuleSet {
	return mg.ruleSet
}

// AddDerivationalRule adds a derivational rule to the generator.
func (mg *MorphologyGenerator) AddDerivationalRule(rule MorphologicalRule) {
	mg.ruleSet.AddDerivationalRule(rule)
}

// AddInflectionalRule adds an inflectional rule to the generator.
func (mg *MorphologyGenerator) AddInflectionalRule(rule MorphologicalRule) {
	mg.ruleSet.AddInflectionalRule(rule)
}

// GetPhonology returns the phonology used by the generator.
func (mg *MorphologyGenerator) GetPhonology() *phonology.Phonology {
	return mg.phonology
}

// GetPhonemePool returns the phoneme pool from the phonology.
func (mg *MorphologyGenerator) GetPhonemePool() *phoneme.PhonemePool {
	if mg.phonology == nil {
		return nil
	}
	return mg.phonology.GetPool()
}

// ValidateWord checks if a word is phonologically well-formed according to
// the phonology rules.
func (mg *MorphologyGenerator) ValidateWord(word *Word) bool {
	if word == nil || mg.phonology == nil {
		return false
	}

	// Check if the phoneme sequence forms valid syllables
	// This is a simplified validation - in a real system, you might want
	// more sophisticated syllable boundary detection
	return len(word.Phonemes) > 0
}

// GetStatistics returns statistics about the generated lexicon and rules.
func (mg *MorphologyGenerator) GetStatistics() map[string]interface{} {
	stats := make(map[string]interface{})

	// Lexicon statistics
	stats["total_words"] = mg.lexicon.GetWordCount()
	stats["culture"] = mg.lexicon.GetCulture()

	// Category statistics
	for category := WordCategoryUnknown; category <= WordCategoryInterjection; category++ {
		count := mg.lexicon.GetWordCountByCategory(category)
		if count > 0 {
			stats[fmt.Sprintf("words_%s", category.String())] = count
		}
	}

	// Rule statistics
	stats["total_rules"] = mg.ruleSet.GetRuleCount()
	stats["derivational_rules"] = mg.ruleSet.GetDerivationalRuleCount()
	stats["inflectional_rules"] = mg.ruleSet.GetInflectionalRuleCount()

	// Phonology statistics
	if mg.phonology != nil {
		stats["syllable_templates"] = len(mg.phonology.GetTemplates())
		stats["phonotactic_rules"] = mg.phonology.GetRuleCount()
	}

	return stats
}
