package phonology

import (
	"github.com/chris-pikul/kismet-zero/lang/phoneme"
)

// MaxClusterSizeRule limits the maximum size of consonant clusters.
type MaxClusterSizeRule struct {
	MaxSize int
}

// NewMaxClusterSizeRule creates a new rule limiting consonant cluster size.
func NewMaxClusterSizeRule(maxSize int) *MaxClusterSizeRule {
	return &MaxClusterSizeRule{MaxSize: maxSize}
}

// Validate checks if consonant clusters don't exceed the maximum size.
func (r *MaxClusterSizeRule) Validate(seq []phoneme.Phoneme) bool {
	if len(seq) == 0 {
		return true
	}

	consonantCount := 0
	for _, p := range seq {
		if p.Type == phoneme.PhonemeTypeConsonant {
			consonantCount++
			if consonantCount > r.MaxSize {
				return false
			}
		} else {
			consonantCount = 0
		}
	}

	return true
}

// VelarNasalOnsetRule disallows velar+nasal onset clusters.
type VelarNasalOnsetRule struct{}

// NewVelarNasalOnsetRule creates a new rule disallowing velar+nasal onset clusters.
func NewVelarNasalOnsetRule() *VelarNasalOnsetRule {
	return &VelarNasalOnsetRule{}
}

// Validate checks if the sequence doesn't start with velar+nasal clusters.
func (r *VelarNasalOnsetRule) Validate(seq []phoneme.Phoneme) bool {
	if len(seq) < 2 {
		return true
	}

	// Check first two phonemes if they're both consonants
	if seq[0].Type == phoneme.PhonemeTypeConsonant && seq[1].Type == phoneme.PhonemeTypeConsonant {
		first := seq[0].Consonant
		second := seq[1].Consonant

		if first != nil && second != nil {
			// Check if first is velar and second is nasal
			if first.Place == phoneme.ConsonantPlaceVelar && second.Manner == phoneme.ConsonantMannerNasal {
				return false
			}
		}
	}

	return true
}

// VowelRequirementRule ensures that vowels are present in the nucleus.
type VowelRequirementRule struct{}

// NewVowelRequirementRule creates a new rule requiring vowels in the nucleus.
func NewVowelRequirementRule() *VowelRequirementRule {
	return &VowelRequirementRule{}
}

// Validate checks if the sequence contains at least one vowel.
func (r *VowelRequirementRule) Validate(seq []phoneme.Phoneme) bool {
	for _, p := range seq {
		if p.Type == phoneme.PhonemeTypeVowel {
			return true
		}
	}
	return false
}

// ConsonantHarmonyRule ensures consonant clusters follow certain harmony patterns.
type ConsonantHarmonyRule struct {
	AllowedCombinations map[phoneme.ConsonantPlace][]phoneme.ConsonantPlace
}

// NewConsonantHarmonyRule creates a new rule for consonant place harmony.
func NewConsonantHarmonyRule() *ConsonantHarmonyRule {
	// Define allowed combinations based on phonological naturalness
	allowed := map[phoneme.ConsonantPlace][]phoneme.ConsonantPlace{
		phoneme.ConsonantPlaceBilabial:     {phoneme.ConsonantPlaceBilabial, phoneme.ConsonantPlaceLabiodental},
		phoneme.ConsonantPlaceLabiodental:  {phoneme.ConsonantPlaceBilabial, phoneme.ConsonantPlaceLabiodental},
		phoneme.ConsonantPlaceDental:       {phoneme.ConsonantPlaceDental, phoneme.ConsonantPlaceAlveolar},
		phoneme.ConsonantPlaceAlveolar:     {phoneme.ConsonantPlaceAlveolar, phoneme.ConsonantPlacePostalveolar, phoneme.ConsonantPlaceDental},
		phoneme.ConsonantPlacePostalveolar: {phoneme.ConsonantPlaceAlveolar, phoneme.ConsonantPlacePostalveolar, phoneme.ConsonantPlaceRetroflex},
		phoneme.ConsonantPlaceRetroflex:    {phoneme.ConsonantPlacePostalveolar, phoneme.ConsonantPlaceRetroflex},
		phoneme.ConsonantPlacePalatal:      {phoneme.ConsonantPlacePalatal, phoneme.ConsonantPlaceVelar},
		phoneme.ConsonantPlaceVelar:        {phoneme.ConsonantPlaceVelar, phoneme.ConsonantPlaceUvular, phoneme.ConsonantPlacePalatal},
		phoneme.ConsonantPlaceUvular:       {phoneme.ConsonantPlaceVelar, phoneme.ConsonantPlaceUvular, phoneme.ConsonantPlacePharyngeal},
		phoneme.ConsonantPlacePharyngeal:   {phoneme.ConsonantPlacePharyngeal, phoneme.ConsonantPlaceGlottal, phoneme.ConsonantPlaceUvular},
		phoneme.ConsonantPlaceGlottal:      {phoneme.ConsonantPlaceGlottal, phoneme.ConsonantPlacePharyngeal},
	}

	return &ConsonantHarmonyRule{AllowedCombinations: allowed}
}

// Validate checks if consonant clusters follow harmony patterns.
func (r *ConsonantHarmonyRule) Validate(seq []phoneme.Phoneme) bool {
	if len(seq) < 2 {
		return true
	}

	for i := 0; i < len(seq)-1; i++ {
		if seq[i].Type == phoneme.PhonemeTypeConsonant && seq[i+1].Type == phoneme.PhonemeTypeConsonant {
			first := seq[i].Consonant
			second := seq[i+1].Consonant

			if first != nil && second != nil {
				allowedPlaces, exists := r.AllowedCombinations[first.Place]
				if !exists {
					// If no specific rule exists, allow the combination
					continue
				}

				// Check if second place is allowed
				allowed := false
				for _, place := range allowedPlaces {
					if place == second.Place {
						allowed = true
						break
					}
				}

				if !allowed {
					return false
				}
			}
		}
	}

	return true
}

// SyllableLengthRule limits the total length of syllables.
type SyllableLengthRule struct {
	MinLength int
	MaxLength int
}

// NewSyllableLengthRule creates a new rule limiting syllable length.
func NewSyllableLengthRule(minLength, maxLength int) *SyllableLengthRule {
	return &SyllableLengthRule{
		MinLength: minLength,
		MaxLength: maxLength,
	}
}

// Validate checks if the syllable length is within the allowed range.
func (r *SyllableLengthRule) Validate(seq []phoneme.Phoneme) bool {
	length := len(seq)
	return length >= r.MinLength && length <= r.MaxLength
}

// VoicingAssimilationRule ensures consonant clusters follow voicing patterns.
type VoicingAssimilationRule struct {
	RequireMatching bool // If true, adjacent consonants must have the same voicing
}

// NewVoicingAssimilationRule creates a new rule for voicing assimilation.
func NewVoicingAssimilationRule(requireMatching bool) *VoicingAssimilationRule {
	return &VoicingAssimilationRule{RequireMatching: requireMatching}
}

// Validate checks if consonant clusters follow voicing patterns.
func (r *VoicingAssimilationRule) Validate(seq []phoneme.Phoneme) bool {
	if len(seq) < 2 || !r.RequireMatching {
		return true
	}

	for i := 0; i < len(seq)-1; i++ {
		if seq[i].Type == phoneme.PhonemeTypeConsonant && seq[i+1].Type == phoneme.PhonemeTypeConsonant {
			first := seq[i].Consonant
			second := seq[i+1].Consonant

			if first != nil && second != nil {
				// Check if voicing matches
				if first.Voicing != second.Voicing &&
					first.Voicing != phoneme.ConsonantVoicingUnknown &&
					second.Voicing != phoneme.ConsonantVoicingUnknown {
					return false
				}
			}
		}
	}

	return true
}

// MannerRestrictionsRule disallows certain manner combinations.
type MannerRestrictionsRule struct {
	ForbiddenCombinations map[phoneme.ConsonantManner][]phoneme.ConsonantManner
}

// NewMannerRestrictionsRule creates a new rule for manner-based restrictions.
func NewMannerRestrictionsRule() *MannerRestrictionsRule {
	// Define forbidden combinations based on phonological constraints
	forbidden := map[phoneme.ConsonantManner][]phoneme.ConsonantManner{
		// Fricatives don't typically cluster with other fricatives in most languages
		phoneme.ConsonantMannerFricative: {phoneme.ConsonantMannerFricative, phoneme.ConsonantMannerSibilants},
		phoneme.ConsonantMannerSibilants: {phoneme.ConsonantMannerFricative, phoneme.ConsonantMannerSibilants},
		// Nasals don't typically cluster with other nasals
		phoneme.ConsonantMannerNasal: {phoneme.ConsonantMannerNasal},
	}

	return &MannerRestrictionsRule{ForbiddenCombinations: forbidden}
}

// Validate checks if consonant clusters avoid forbidden manner combinations.
func (r *MannerRestrictionsRule) Validate(seq []phoneme.Phoneme) bool {
	if len(seq) < 2 {
		return true
	}

	for i := 0; i < len(seq)-1; i++ {
		if seq[i].Type == phoneme.PhonemeTypeConsonant && seq[i+1].Type == phoneme.PhonemeTypeConsonant {
			first := seq[i].Consonant
			second := seq[i+1].Consonant

			if first != nil && second != nil {
				forbiddenManners, exists := r.ForbiddenCombinations[first.Manner]
				if exists {
					for _, forbidden := range forbiddenManners {
						if second.Manner == forbidden {
							return false
						}
					}
				}
			}
		}
	}

	return true
}

// DefaultRules returns a set of commonly used phonotactic rules.
func DefaultRules() []PhonotacticRule {
	return []PhonotacticRule{
		NewMaxClusterSizeRule(3),
		NewVelarNasalOnsetRule(),
		NewVowelRequirementRule(),
		NewSyllableLengthRule(1, 6),
		NewVoicingAssimilationRule(false), // Allow mixed voicing by default
		NewMannerRestrictionsRule(),
	}
}

// StrictDefaultRules returns a more restrictive set of phonotactic rules.
func StrictDefaultRules() []PhonotacticRule {
	return []PhonotacticRule{
		NewMaxClusterSizeRule(2), // Stricter cluster size
		NewVelarNasalOnsetRule(),
		NewVowelRequirementRule(),
		NewSyllableLengthRule(1, 4),      // Shorter syllables
		NewConsonantHarmonyRule(),        // Add harmony requirements
		NewVoicingAssimilationRule(true), // Require voicing matching
		NewMannerRestrictionsRule(),
	}
}
