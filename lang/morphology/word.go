package morphology

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
	"github.com/chris-pikul/kismet-zero/lang/phonology"
)

// WordBuilder constructs words from morphemes using morphological rules
// and ensures phonological well-formedness.
type WordBuilder struct {
	phonology *phonology.Phonology
	ruleSet   *RuleSet
}

// NewWordBuilder creates a new word builder with the given phonology and rule set.
func NewWordBuilder(phonology *phonology.Phonology, ruleSet *RuleSet) *WordBuilder {
	return &WordBuilder{
		phonology: phonology,
		ruleSet:   ruleSet,
	}
}

// BuildWord creates a word from a root morpheme, optionally applying
// derivational and inflectional rules.
func (wb *WordBuilder) BuildWord(
	root *Morpheme,
	category WordCategory,
	culture string,
	rng *rand.Rand,
) (*Word, error) {
	if root == nil {
		return nil, fmt.Errorf("root morpheme is required")
	}

	if root.Type != MorphemeTypeRoot {
		return nil, fmt.Errorf("root morpheme must be of type Root")
	}

	// Start with the root morpheme
	morphemes := []Morpheme{*root}

	// Apply derivational rules (0-2 times)
	numDerivations := rng.IntN(3)
	for i := 0; i < numDerivations; i++ {
		if result, err := wb.ruleSet.ApplyRandomRule(morphemes, wb.phonology, rng); err == nil {
			morphemes = result
		}
	}

	// Apply inflectional rules (0-1 times)
	if rng.IntN(2) == 0 {
		if result, err := wb.ruleSet.ApplyRandomRule(morphemes, wb.phonology, rng); err == nil {
			morphemes = result
		}
	}

	// Combine all phonemes from morphemes
	var allPhonemes []phoneme.Phoneme
	for _, m := range morphemes {
		allPhonemes = append(allPhonemes, m.Phonemes...)
	}

	// Generate written representation
	written := wb.generateWrittenForm(allPhonemes)

	// Calculate word weight based on morpheme weights
	weight := wb.calculateWordWeight(morphemes)

	// Determine frequency based on complexity
	frequency := wb.determineWordFrequency(len(morphemes))

	word := &Word{
		ID:        wb.generateWordID(root.Meaning, category, culture),
		Morphemes: morphemes,
		Meaning:   wb.combineMeanings(morphemes),
		Phonemes:  allPhonemes,
		Written:   written,
		Category:  category,
		Culture:   culture,
		Frequency: frequency,
		Weight:    weight,
	}

	return word, nil
}

// generateWrittenForm creates a written representation of the word from phonemes.
func (wb *WordBuilder) generateWrittenForm(phonemes []phoneme.Phoneme) string {
	var result string
	for _, p := range phonemes {
		result += p.Symbol
	}
	return result
}

// calculateWordWeight calculates the weight of a word based on its morphemes.
func (wb *WordBuilder) calculateWordWeight(morphemes []Morpheme) float32 {
	var totalWeight float32
	for _, m := range morphemes {
		totalWeight += m.Weight
	}

	// Normalize by number of morphemes
	if len(morphemes) > 0 {
		totalWeight /= float32(len(morphemes))
	}

	return totalWeight
}

// determineWordFrequency determines the frequency category based on word complexity.
func (wb *WordBuilder) determineWordFrequency(numMorphemes int) WordFrequency {
	switch numMorphemes {
	case 1:
		return WordFrequencyVeryCommon
	case 2:
		return WordFrequencyCommon
	case 3:
		return WordFrequencyUncommon
	default:
		return WordFrequencyRare
	}
}

// generateWordID creates a unique identifier for a word.
func (wb *WordBuilder) generateWordID(meaning string, category WordCategory, culture string) string {
	return fmt.Sprintf("%s_%s_%s_%d", culture, category.String(), meaning, len(meaning))
}

// combineMeanings combines the meanings of multiple morphemes into a single meaning.
func (wb *WordBuilder) combineMeanings(morphemes []Morpheme) string {
	if len(morphemes) == 0 {
		return ""
	}

	if len(morphemes) == 1 {
		return morphemes[0].Meaning
	}

	// Simple concatenation for now
	// In a more sophisticated system, this would handle semantic composition
	var meanings []string
	for _, m := range morphemes {
		meanings = append(meanings, m.Meaning)
	}

	return strings.Join(meanings, "+")
}

// Lexicon represents a collection of words for a specific culture or language.
type Lexicon struct {
	culture    string
	words      map[string]*Word
	byCategory map[WordCategory][]*Word
}

// NewLexicon creates a new lexicon for the specified culture.
func NewLexicon(culture string) *Lexicon {
	return &Lexicon{
		culture:    culture,
		words:      make(map[string]*Word),
		byCategory: make(map[WordCategory][]*Word),
	}
}

// AddWord adds a word to the lexicon.
func (l *Lexicon) AddWord(word *Word) error {
	if word == nil {
		return fmt.Errorf("cannot add nil word to lexicon")
	}

	if word.Culture != l.culture {
		return fmt.Errorf("word culture %s does not match lexicon culture %s", word.Culture, l.culture)
	}

	// Add to main word map
	l.words[word.ID] = word

	// Add to category index
	l.byCategory[word.Category] = append(l.byCategory[word.Category], word)

	return nil
}

// GetWord retrieves a word by its ID.
func (l *Lexicon) GetWord(id string) *Word {
	return l.words[id]
}

// GetWordsByCategory returns all words of a specific category.
func (l *Lexicon) GetWordsByCategory(category WordCategory) []*Word {
	words := l.byCategory[category]
	result := make([]*Word, len(words))
	copy(result, words)
	return result
}

// GetAllWords returns all words in the lexicon.
func (l *Lexicon) GetAllWords() []*Word {
	words := make([]*Word, 0, len(l.words))
	for _, word := range l.words {
		words = append(words, word)
	}
	return words
}

// GetWordCount returns the total number of words in the lexicon.
func (l *Lexicon) GetWordCount() int {
	return len(l.words)
}

// GetWordCountByCategory returns the number of words in a specific category.
func (l *Lexicon) GetWordCountByCategory(category WordCategory) int {
	return len(l.byCategory[category])
}

// GetCulture returns the culture associated with this lexicon.
func (l *Lexicon) GetCulture() string {
	return l.culture
}

// SearchWords searches for words by meaning substring.
func (l *Lexicon) SearchWords(query string) []*Word {
	var results []*Word
	query = strings.ToLower(query)

	for _, word := range l.words {
		if strings.Contains(strings.ToLower(word.Meaning), query) {
			results = append(results, word)
		}
	}

	return results
}

// GetRandomWord returns a random word from the lexicon.
func (l *Lexicon) GetRandomWord(rng *rand.Rand) *Word {
	words := l.GetAllWords()
	if len(words) == 0 {
		return nil
	}

	return words[rng.IntN(len(words))]
}

// GetRandomWordByCategory returns a random word of a specific category.
func (l *Lexicon) GetRandomWordByCategory(category WordCategory, rng *rand.Rand) *Word {
	words := l.GetWordsByCategory(category)
	if len(words) == 0 {
		return nil
	}

	return words[rng.IntN(len(words))]
}
