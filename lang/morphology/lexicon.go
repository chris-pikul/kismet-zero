package morphology

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"
)

// LexiconManager provides advanced lexicon management capabilities including
// cultural context, semantic domains, and cross-references.
type LexiconManager struct {
	lexicons map[string]*Lexicon
}

// NewLexiconManager creates a new lexicon manager.
func NewLexiconManager() *LexiconManager {
	return &LexiconManager{
		lexicons: make(map[string]*Lexicon),
	}
}

// CreateLexicon creates a new lexicon for the specified culture.
func (lm *LexiconManager) CreateLexicon(culture string) *Lexicon {
	lexicon := NewLexicon(culture)
	lm.lexicons[culture] = lexicon
	return lexicon
}

// GetLexicon retrieves a lexicon by culture name.
func (lm *LexiconManager) GetLexicon(culture string) *Lexicon {
	return lm.lexicons[culture]
}

// GetAllCultures returns all cultures that have lexicons.
func (lm *LexiconManager) GetAllCultures() []string {
	cultures := make([]string, 0, len(lm.lexicons))
	for culture := range lm.lexicons {
		cultures = append(cultures, culture)
	}
	sort.Strings(cultures)
	return cultures
}

// RemoveLexicon removes a lexicon and all its words.
func (lm *LexiconManager) RemoveLexicon(culture string) bool {
	if _, exists := lm.lexicons[culture]; exists {
		delete(lm.lexicons, culture)
		return true
	}
	return false
}

// GetTotalWordCount returns the total number of words across all lexicons.
func (lm *LexiconManager) GetTotalWordCount() int {
	total := 0
	for _, lexicon := range lm.lexicons {
		total += lexicon.GetWordCount()
	}
	return total
}

// GetWordCountByCulture returns the number of words in a specific culture's lexicon.
func (lm *LexiconManager) GetWordCountByCulture(culture string) int {
	if lexicon := lm.GetLexicon(culture); lexicon != nil {
		return lexicon.GetWordCount()
	}
	return 0
}

// SearchAcrossCultures searches for words across all lexicons.
func (lm *LexiconManager) SearchAcrossCultures(query string) map[string][]*Word {
	results := make(map[string][]*Word)
	query = strings.ToLower(query)

	for culture, lexicon := range lm.lexicons {
		words := lexicon.SearchWords(query)
		if len(words) > 0 {
			results[culture] = words
		}
	}

	return results
}

// GetRandomWordFromAnyCulture returns a random word from any lexicon.
func (lm *LexiconManager) GetRandomWordFromAnyCulture(rng *rand.Rand) (*Word, string) {
	cultures := lm.GetAllCultures()
	if len(cultures) == 0 {
		return nil, ""
	}

	// Select a random culture
	culture := cultures[rng.IntN(len(cultures))]
	lexicon := lm.GetLexicon(culture)
	if lexicon == nil {
		return nil, ""
	}

	word := lexicon.GetRandomWord(rng)
	return word, culture
}

// GetRandomWordByCategoryAcrossCultures returns a random word of a specific
// category from any lexicon.
func (lm *LexiconManager) GetRandomWordByCategoryAcrossCultures(
	category WordCategory,
	rng *rand.Rand,
) (*Word, string) {
	cultures := lm.GetAllCultures()
	if len(cultures) == 0 {
		return nil, ""
	}

	// Collect all words of the specified category
	var allWords []struct {
		word    *Word
		culture string
		lexicon *Lexicon
	}

	for _, culture := range cultures {
		lexicon := lm.GetLexicon(culture)
		if lexicon == nil {
			continue
		}

		words := lexicon.GetWordsByCategory(category)
		for _, word := range words {
			allWords = append(allWords, struct {
				word    *Word
				culture string
				lexicon *Lexicon
			}{word, culture, lexicon})
		}
	}

	if len(allWords) == 0 {
		return nil, ""
	}

	// Select a random word
	selection := allWords[rng.IntN(len(allWords))]
	return selection.word, selection.culture
}

// ExportLexicon exports a lexicon to a structured format for external use.
func (lm *LexiconManager) ExportLexicon(culture string) (map[string]interface{}, error) {
	lexicon := lm.GetLexicon(culture)
	if lexicon == nil {
		return nil, fmt.Errorf("lexicon for culture %s not found", culture)
	}

	export := make(map[string]interface{})
	export["culture"] = culture
	export["total_words"] = lexicon.GetWordCount()

	// Export words by category
	wordsByCategory := make(map[string][]map[string]interface{})
	for category := WordCategoryUnknown; category <= WordCategoryInterjection; category++ {
		words := lexicon.GetWordsByCategory(category)
		if len(words) == 0 {
			continue
		}

		var categoryWords []map[string]interface{}
		for _, word := range words {
			wordExport := map[string]interface{}{
				"id":        word.ID,
				"meaning":   word.Meaning,
				"written":   word.Written,
				"category":  word.Category.String(),
				"frequency": word.Frequency.String(),
				"weight":    word.Weight,
				"morphemes": len(word.Morphemes),
				"phonemes":  len(word.Phonemes),
			}
			categoryWords = append(categoryWords, wordExport)
		}

		wordsByCategory[category.String()] = categoryWords
	}

	export["words_by_category"] = wordsByCategory

	return export, nil
}

// ImportLexicon imports a lexicon from an external format.
func (lm *LexiconManager) ImportLexicon(export map[string]interface{}) error {
	culture, ok := export["culture"].(string)
	if !ok {
		return fmt.Errorf("invalid export format: missing culture")
	}

	lexicon := lm.CreateLexicon(culture)

	// Import words by category
	wordsByCategory, ok := export["words_by_category"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid export format: missing words_by_category")
	}

	for categoryStr, categoryWords := range wordsByCategory {
		// Parse category
		var category WordCategory
		switch strings.ToLower(categoryStr) {
		case "noun":
			category = WordCategoryNoun
		case "verb":
			category = WordCategoryVerb
		case "adjective":
			category = WordCategoryAdjective
		case "adverb":
			category = WordCategoryAdverb
		case "pronoun":
			category = WordCategoryPronoun
		case "preposition":
			category = WordCategoryPreposition
		case "conjunction":
			category = WordCategoryConjunction
		case "interjection":
			category = WordCategoryInterjection
		default:
			category = WordCategoryUnknown
		}

		// Import words in this category
		wordsList, ok := categoryWords.([]interface{})
		if !ok {
			continue
		}

		for _, wordData := range wordsList {
			wordMap, ok := wordData.(map[string]interface{})
			if !ok {
				continue
			}

			// Create a basic word structure
			// Note: This is a simplified import - in a real system, you'd want
			// to reconstruct the full morpheme and phoneme structures
			word := &Word{
				ID:        getString(wordMap, "id"),
				Meaning:   getString(wordMap, "meaning"),
				Written:   getString(wordMap, "written"),
				Category:  category,
				Culture:   culture,
				Frequency: parseWordFrequency(getString(wordMap, "frequency")),
				Weight:    getFloat32(wordMap, "weight"),
			}

			// Add to lexicon
			if err := lexicon.AddWord(word); err != nil {
				// Log error but continue with other words
				fmt.Printf("Warning: failed to import word %s: %v\n", word.ID, err)
			}
		}
	}

	return nil
}

// getString safely extracts a string value from a map.
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// getFloat32 safely extracts a float32 value from a map.
func getFloat32(m map[string]interface{}, key string) float32 {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case float32:
			return v
		case float64:
			return float32(v)
		case int:
			return float32(v)
		}
	}
	return 0.0
}

// parseWordFrequency parses a word frequency string.
func parseWordFrequency(freqStr string) WordFrequency {
	switch strings.ToLower(freqStr) {
	case "rare":
		return WordFrequencyRare
	case "uncommon":
		return WordFrequencyUncommon
	case "common":
		return WordFrequencyCommon
	case "very_common":
		return WordFrequencyVeryCommon
	default:
		return WordFrequencyUnknown
	}
}

// GetLexiconStatistics returns comprehensive statistics about all lexicons.
func (lm *LexiconManager) GetLexiconStatistics() map[string]interface{} {
	stats := make(map[string]interface{})

	stats["total_cultures"] = len(lm.lexicons)
	stats["total_words"] = lm.GetTotalWordCount()

	// Per-culture statistics
	cultureStats := make(map[string]interface{})
	for culture, lexicon := range lm.lexicons {
		cultureStats[culture] = map[string]interface{}{
			"total_words": lexicon.GetWordCount(),
			"categories":  lexicon.getCategoryBreakdown(),
		}
	}
	stats["cultures"] = cultureStats

	return stats
}

// getCategoryBreakdown returns a breakdown of words by category for a lexicon.
func (l *Lexicon) getCategoryBreakdown() map[string]int {
	breakdown := make(map[string]int)

	for category := WordCategoryUnknown; category <= WordCategoryInterjection; category++ {
		count := l.GetWordCountByCategory(category)
		if count > 0 {
			breakdown[category.String()] = count
		}
	}

	return breakdown
}
