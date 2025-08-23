package morphology

import (
	"fmt"
	"math/rand/v2"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
	"github.com/chris-pikul/kismet-zero/lang/phonology"
)

// MorphemeBuilder constructs morphemes from phonological elements and
// applies morphological constraints.
type MorphemeBuilder struct {
	phonology *phonology.Phonology
}

// NewMorphemeBuilder creates a new morpheme builder with the given phonology.
func NewMorphemeBuilder(phonology *phonology.Phonology) *MorphemeBuilder {
	return &MorphemeBuilder{
		phonology: phonology,
	}
}

// BuildRootMorpheme creates a root morpheme with the given meaning and
// phonological constraints.
func (mb *MorphemeBuilder) BuildRootMorpheme(
	meaning string,
	culture string,
	rng *rand.Rand,
) (*Morpheme, error) {
	if mb.phonology == nil {
		return nil, fmt.Errorf("phonology is required for morpheme building")
	}

	// Generate 1-3 syllables for a root morpheme
	numSyllables := rng.IntN(3) + 1
	var allPhonemes []phoneme.Phoneme
	var syllables []string

	for i := 0; i < numSyllables; i++ {
		// Generate a syllable using the phonology
		syllablePhonemes := mb.generateSyllable(rng)
		if len(syllablePhonemes) == 0 {
			continue
		}

		// Convert to string representation
		syllableStr := mb.phonemesToString(syllablePhonemes)
		syllables = append(syllables, syllableStr)
		allPhonemes = append(allPhonemes, syllablePhonemes...)
	}

	if len(allPhonemes) == 0 {
		return nil, fmt.Errorf("failed to generate phonemes for morpheme")
	}

	// Calculate weight based on frequency and complexity
	weight := mb.calculateMorphemeWeight(allPhonemes, MorphemeFrequencyCommon)

	morpheme := &Morpheme{
		ID:        mb.generateMorphemeID(meaning, culture),
		Type:      MorphemeTypeRoot,
		Meaning:   meaning,
		Phonemes:  allPhonemes,
		Syllables: syllables,
		Weight:    weight,
		Culture:   culture,
		Frequency: MorphemeFrequencyCommon,
	}

	return morpheme, nil
}

// BuildAffixMorpheme creates an affix morpheme (prefix, suffix, infix) with
// the given type and meaning.
func (mb *MorphemeBuilder) BuildAffixMorpheme(
	morphemeType MorphemeType,
	meaning string,
	culture string,
	rng *rand.Rand,
) (*Morpheme, error) {
	if morphemeType == MorphemeTypeRoot {
		return nil, fmt.Errorf("use BuildRootMorpheme for root morphemes")
	}

	if mb.phonology == nil {
		return nil, fmt.Errorf("phonology is required for morpheme building")
	}

	// Affixes are typically shorter than roots (1-2 syllables)
	numSyllables := rng.IntN(2) + 1
	var allPhonemes []phoneme.Phoneme
	var syllables []string

	for i := 0; i < numSyllables; i++ {
		syllablePhonemes := mb.generateSyllable(rng)
		if len(syllablePhonemes) == 0 {
			continue
		}

		syllableStr := mb.phonemesToString(syllablePhonemes)
		syllables = append(syllables, syllableStr)
		allPhonemes = append(allPhonemes, syllablePhonemes...)
	}

	if len(allPhonemes) == 0 {
		return nil, fmt.Errorf("failed to generate phonemes for affix")
	}

	// Affixes typically have lower weight than roots
	weight := mb.calculateMorphemeWeight(allPhonemes, MorphemeFrequencyUncommon)

	morpheme := &Morpheme{
		ID:        mb.generateMorphemeID(meaning, culture),
		Type:      morphemeType,
		Meaning:   meaning,
		Phonemes:  allPhonemes,
		Syllables: syllables,
		Weight:    weight,
		Culture:   culture,
		Frequency: MorphemeFrequencyUncommon,
	}

	return morpheme, nil
}

// generateSyllable creates a single syllable using the phonology's rules.
func (mb *MorphemeBuilder) generateSyllable(rng *rand.Rand) []phoneme.Phoneme {
	// Use the phonology to generate a valid syllable
	pool := mb.phonology.GetPool()
	if pool == nil {
		return nil
	}

	// Simple CV or CVC syllable for affixes
	templates := mb.phonology.GetTemplates()
	if len(templates) == 0 {
		// Fallback: create a simple CV syllable
		var result []phoneme.Phoneme
		if consonant := pool.Consonants.WeightedChoice(rng); consonant != nil {
			result = append(result, *consonant)
		}
		if vowel := pool.Vowels.WeightedChoice(rng); vowel != nil {
			result = append(result, *vowel)
		}
		return result
	}

	// Pick a template and fill it
	template := mb.pickTemplate(templates, rng)
	return mb.fillTemplate(template, pool, rng)
}

// pickTemplate selects a syllable template using weighted random selection.
func (mb *MorphemeBuilder) pickTemplate(templates []phonology.SyllableTemplate, rng *rand.Rand) phonology.SyllableTemplate {
	if len(templates) == 0 {
		// Return a simple CV template as fallback
		return phonology.SyllableTemplate{
			Onset:   phonology.SyllableSlot{Required: true, MaxSize: 1},
			Nucleus: phonology.SyllableSlot{Required: true, MaxSize: 1},
			Coda:    phonology.SyllableSlot{Required: false, MaxSize: 0},
			Weight:  1,
		}
	}

	// Calculate total weight
	var totalWeight int
	for _, t := range templates {
		totalWeight += t.Weight
	}

	if totalWeight == 0 {
		return templates[0] // Fallback to first template
	}

	// Weighted selection
	target := rng.IntN(totalWeight)
	var cumulative int
	for _, t := range templates {
		cumulative += t.Weight
		if cumulative > target {
			return t
		}
	}

	return templates[len(templates)-1] // Fallback to last template
}

// fillTemplate fills a syllable template with phonemes from the pool.
func (mb *MorphemeBuilder) fillTemplate(
	template phonology.SyllableTemplate,
	pool *phoneme.PhonemePool,
	rng *rand.Rand,
) []phoneme.Phoneme {
	var result []phoneme.Phoneme

	// Fill onset
	if template.Onset.Required && template.Onset.MaxSize > 0 {
		size := rng.IntN(template.Onset.MaxSize) + 1
		for i := 0; i < size; i++ {
			if phoneme := pool.Consonants.WeightedChoice(rng); phoneme != nil {
				result = append(result, *phoneme)
			}
		}
	}

	// Fill nucleus
	if template.Nucleus.Required && template.Nucleus.MaxSize > 0 {
		size := rng.IntN(template.Nucleus.MaxSize) + 1
		for i := 0; i < size; i++ {
			if phoneme := pool.Vowels.WeightedChoice(rng); phoneme != nil {
				result = append(result, *phoneme)
			}
		}
	}

	// Fill coda
	if template.Coda.Required && template.Coda.MaxSize > 0 {
		size := rng.IntN(template.Coda.MaxSize) + 1
		for i := 0; i < size; i++ {
			if phoneme := pool.Consonants.WeightedChoice(rng); phoneme != nil {
				result = append(result, *phoneme)
			}
		}
	}

	return result
}

// phonemesToString converts a slice of phonemes to a string representation.
func (mb *MorphemeBuilder) phonemesToString(phonemes []phoneme.Phoneme) string {
	var result string
	for _, p := range phonemes {
		result += p.Symbol
	}
	return result
}

// calculateMorphemeWeight calculates the weight of a morpheme based on its
// phonological complexity and frequency.
func (mb *MorphemeBuilder) calculateMorphemeWeight(
	phonemes []phoneme.Phoneme,
	frequency MorphemeFrequency,
) float32 {
	baseWeight := float32(1.0)

	// Adjust weight based on frequency
	switch frequency {
	case MorphemeFrequencyRare:
		baseWeight *= 0.5
	case MorphemeFrequencyUncommon:
		baseWeight *= 0.8
	case MorphemeFrequencyCommon:
		baseWeight *= 1.0
	case MorphemeFrequencyVeryCommon:
		baseWeight *= 1.5
	}

	// Adjust weight based on complexity (more phonemes = higher weight)
	complexityFactor := float32(len(phonemes)) / 5.0 // Normalize to reasonable range
	baseWeight *= (1.0 + complexityFactor)

	return baseWeight
}

// generateMorphemeID creates a unique identifier for a morpheme.
func (mb *MorphemeBuilder) generateMorphemeID(meaning, culture string) string {
	// Simple ID generation - in a real system, this might be more sophisticated
	return fmt.Sprintf("%s_%s_%d", culture, meaning, len(meaning))
}
