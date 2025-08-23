package morphology

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
	"github.com/chris-pikul/kismet-zero/lang/phonology"
)

// NameGenerator creates culturally appropriate names using morphological
// building blocks and cultural patterns.
type NameGenerator struct {
	phonology *phonology.Phonology
	lexicon   *Lexicon
}

// NewNameGenerator creates a new name generator with the given phonology and lexicon.
func NewNameGenerator(phonology *phonology.Phonology, lexicon *Lexicon) *NameGenerator {
	return &NameGenerator{
		phonology: phonology,
		lexicon:   lexicon,
	}
}

// GeneratePersonalName creates a personal name appropriate for the given culture.
func (ng *NameGenerator) GeneratePersonalName(
	culture string,
	gender string,
	formality string,
	rng *rand.Rand,
) (*Name, error) {
	if ng.phonology == nil {
		return nil, fmt.Errorf("phonology is required for name generation")
	}

	// Generate 1-3 syllables for a personal name
	numSyllables := rng.IntN(3) + 1
	var allPhonemes []phoneme.Phoneme
	var syllables []string

	for i := 0; i < numSyllables; i++ {
		syllablePhonemes := ng.generateSyllable(rng)
		if len(syllablePhonemes) == 0 {
			continue
		}

		syllableStr := ng.phonemesToString(syllablePhonemes)
		syllables = append(syllables, syllableStr)
		allPhonemes = append(allPhonemes, syllablePhonemes...)
	}

	if len(allPhonemes) == 0 {
		return nil, fmt.Errorf("failed to generate phonemes for personal name")
	}

	// Create a meaningful name based on cultural patterns
	meaning := ng.generatePersonalNameMeaning(culture, gender, formality)
	value := ng.phonemesToString(allPhonemes)

	name := &Name{
		ID:        ng.generateNameID(value, NameTypePersonal, culture),
		Type:      NameTypePersonal,
		Value:     value,
		Meaning:   meaning,
		Culture:   culture,
		Gender:    gender,
		Formality: formality,
	}

	return name, nil
}

// GeneratePlaceName creates a place name appropriate for the given culture.
func (ng *NameGenerator) GeneratePlaceName(
	culture string,
	placeType string,
	rng *rand.Rand,
) (*Name, error) {
	if ng.phonology == nil {
		return nil, fmt.Errorf("phonology is required for name generation")
	}

	// Place names are typically longer than personal names (2-4 syllables)
	numSyllables := rng.IntN(3) + 2
	var allPhonemes []phoneme.Phoneme
	var syllables []string

	for i := 0; i < numSyllables; i++ {
		syllablePhonemes := ng.generateSyllable(rng)
		if len(syllablePhonemes) == 0 {
			continue
		}

		syllableStr := ng.phonemesToString(syllablePhonemes)
		syllables = append(syllables, syllableStr)
		allPhonemes = append(allPhonemes, syllablePhonemes...)
	}

	if len(allPhonemes) == 0 {
		return nil, fmt.Errorf("failed to generate phonemes for place name")
	}

	// Create a meaningful place name
	meaning := ng.generatePlaceNameMeaning(culture, placeType)
	value := ng.phonemesToString(allPhonemes)

	name := &Name{
		ID:      ng.generateNameID(value, NameTypePlace, culture),
		Type:    NameTypePlace,
		Value:   value,
		Meaning: meaning,
		Culture: culture,
	}

	return name, nil
}

// GenerateTitle creates a title or honorific appropriate for the given culture.
func (ng *NameGenerator) GenerateTitle(
	culture string,
	rank string,
	rng *rand.Rand,
) (*Name, error) {
	if ng.phonology == nil {
		return nil, fmt.Errorf("phonology is required for name generation")
	}

	// Titles are typically short (1-2 syllables)
	numSyllables := rng.IntN(2) + 1
	var allPhonemes []phoneme.Phoneme
	var syllables []string

	for i := 0; i < numSyllables; i++ {
		syllablePhonemes := ng.generateSyllable(rng)
		if len(syllablePhonemes) == 0 {
			continue
		}

		syllableStr := ng.phonemesToString(syllablePhonemes)
		syllables = append(syllables, syllableStr)
		allPhonemes = append(allPhonemes, syllablePhonemes...)
	}

	if len(allPhonemes) == 0 {
		return nil, fmt.Errorf("failed to generate phonemes for title")
	}

	// Create a meaningful title
	meaning := ng.generateTitleMeaning(culture, rank)
	value := ng.phonemesToString(allPhonemes)

	name := &Name{
		ID:      ng.generateNameID(value, NameTypeTitle, culture),
		Type:    NameTypeTitle,
		Value:   value,
		Meaning: meaning,
		Culture: culture,
	}

	return name, nil
}

// GenerateDeityName creates a deity name appropriate for the given culture.
func (ng *NameGenerator) GenerateDeityName(
	culture string,
	domain string,
	rng *rand.Rand,
) (*Name, error) {
	if ng.phonology == nil {
		return nil, fmt.Errorf("phonology is required for name generation")
	}

	// Deity names are typically longer and more complex (2-4 syllables)
	numSyllables := rng.IntN(3) + 2
	var allPhonemes []phoneme.Phoneme
	var syllables []string

	for i := 0; i < numSyllables; i++ {
		syllablePhonemes := ng.generateSyllable(rng)
		if len(syllablePhonemes) == 0 {
			continue
		}

		syllableStr := ng.phonemesToString(syllablePhonemes)
		syllables = append(syllables, syllableStr)
		allPhonemes = append(allPhonemes, syllablePhonemes...)
	}

	if len(allPhonemes) == 0 {
		return nil, fmt.Errorf("failed to generate phonemes for deity name")
	}

	// Create a meaningful deity name
	meaning := ng.generateDeityNameMeaning(culture, domain)
	value := ng.phonemesToString(allPhonemes)

	name := &Name{
		ID:      ng.generateNameID(value, NameTypeDeity, culture),
		Type:    NameTypeDeity,
		Value:   value,
		Meaning: meaning,
		Culture: culture,
	}

	return name, nil
}

// generateSyllable creates a single syllable using the phonology's rules.
func (ng *NameGenerator) generateSyllable(rng *rand.Rand) []phoneme.Phoneme {
	pool := ng.phonology.GetPool()
	if pool == nil {
		return nil
	}

	// Simple CV or CVC syllable for names
	var result []phoneme.Phoneme

	// Add onset consonant (optional)
	if rng.IntN(2) == 0 {
		if consonant := pool.Consonants.WeightedChoice(rng); consonant != nil {
			result = append(result, *consonant)
		}
	}

	// Add nucleus vowel (required)
	if vowel := pool.Vowels.WeightedChoice(rng); vowel != nil {
		result = append(result, *vowel)
	}

	// Add coda consonant (optional)
	if rng.IntN(2) == 0 {
		if consonant := pool.Consonants.WeightedChoice(rng); consonant != nil {
			result = append(result, *consonant)
		}
	}

	return result
}

// phonemesToString converts a slice of phonemes to a string representation.
func (ng *NameGenerator) phonemesToString(phonemes []phoneme.Phoneme) string {
	var result string
	for _, p := range phonemes {
		result += p.Symbol
	}
	return result
}

// generatePersonalNameMeaning creates a meaningful personal name based on cultural patterns.
func (ng *NameGenerator) generatePersonalNameMeaning(culture, gender, formality string) string {
	var meaning string

	// Base meaning based on culture
	switch strings.ToLower(culture) {
	case "nordic", "norse":
		meaning = "brave warrior"
	case "elvish", "elven":
		meaning = "wise guardian"
	case "dwarven", "dwarf":
		meaning = "strong craftsman"
	case "orcish", "orc":
		meaning = "fierce fighter"
	default:
		meaning = "noble person"
	}

	// Add gender-specific elements
	if gender != "" {
		switch strings.ToLower(gender) {
		case "male":
			meaning += " (male)"
		case "female":
			meaning += " (female)"
		}
	}

	// Add formality elements
	if formality != "" {
		switch strings.ToLower(formality) {
		case "formal":
			meaning += " (honored)"
		case "informal":
			meaning += " (friendly)"
		}
	}

	return meaning
}

// generatePlaceNameMeaning creates a meaningful place name based on cultural patterns.
func (ng *NameGenerator) generatePlaceNameMeaning(culture, placeType string) string {
	var meaning string

	// Base meaning based on place type
	switch strings.ToLower(placeType) {
	case "city", "town":
		meaning = "settlement"
	case "mountain", "peak":
		meaning = "high place"
	case "river", "stream":
		meaning = "flowing water"
	case "forest", "wood":
		meaning = "tree land"
	case "castle", "fort":
		meaning = "stronghold"
	default:
		meaning = "landmark"
	}

	// Add cultural elements
	if culture != "" {
		meaning = fmt.Sprintf("%s of %s", meaning, culture)
	}

	return meaning
}

// generateTitleMeaning creates a meaningful title based on cultural patterns.
func (ng *NameGenerator) generateTitleMeaning(culture, rank string) string {
	var meaning string

	// Base meaning based on rank
	switch strings.ToLower(rank) {
	case "king", "queen":
		meaning = "ruler"
	case "lord", "lady":
		meaning = "noble"
	case "chief", "elder":
		meaning = "leader"
	case "priest", "priestess":
		meaning = "spiritual guide"
	default:
		meaning = "honored one"
	}

	// Add cultural elements
	if culture != "" {
		meaning = fmt.Sprintf("%s of %s", meaning, culture)
	}

	return meaning
}

// generateDeityNameMeaning creates a meaningful deity name based on cultural patterns.
func (ng *NameGenerator) generateDeityNameMeaning(culture, domain string) string {
	var meaning string

	// Base meaning based on domain
	switch strings.ToLower(domain) {
	case "war", "battle":
		meaning = "god of war"
	case "wisdom", "knowledge":
		meaning = "god of wisdom"
	case "nature", "earth":
		meaning = "god of nature"
	case "sun", "light":
		meaning = "god of light"
	case "moon", "night":
		meaning = "god of night"
	default:
		meaning = "divine being"
	}

	// Add cultural elements
	if culture != "" {
		meaning = fmt.Sprintf("%s of %s", meaning, culture)
	}

	return meaning
}

// generateNameID creates a unique identifier for a name.
func (ng *NameGenerator) generateNameID(value string, nameType NameType, culture string) string {
	return fmt.Sprintf("%s_%s_%s_%d", culture, nameType.String(), value, len(value))
}
