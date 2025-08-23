package phonology

import (
	"math/rand/v2"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
)

// SyllableTemplate defines the structure of a syllable with onset, nucleus, and coda slots.
// Each slot can be optional or required, and the template has a weight for selection.
type SyllableTemplate struct {
	Onset   SyllableSlot `json:"onset"`
	Nucleus SyllableSlot `json:"nucleus"`
	Coda    SyllableSlot `json:"coda"`
	Weight  int          `json:"weight"`
}

// SyllableSlot defines a position in a syllable that can hold phonemes.
type SyllableSlot struct {
	Required bool `json:"required"`
	MaxSize  int  `json:"maxSize"` // Maximum number of phonemes in this slot
}

// NewSyllableTemplate creates a new syllable template with the specified configuration.
func NewSyllableTemplate(onset, nucleus, coda SyllableSlot, weight int) SyllableTemplate {
	return SyllableTemplate{
		Onset:   onset,
		Nucleus: nucleus,
		Coda:    coda,
		Weight:  weight,
	}
}

// Common syllable templates for convenience
var (
	// CV - Consonant-Vowel (most common)
	TemplateCV = NewSyllableTemplate(
		SyllableSlot{Required: true, MaxSize: 1},
		SyllableSlot{Required: true, MaxSize: 1},
		SyllableSlot{Required: false, MaxSize: 0},
		10,
	)

	// CVC - Consonant-Vowel-Consonant
	TemplateCVC = NewSyllableTemplate(
		SyllableSlot{Required: true, MaxSize: 1},
		SyllableSlot{Required: true, MaxSize: 1},
		SyllableSlot{Required: true, MaxSize: 1},
		8,
	)

	// V - Vowel only
	TemplateV = NewSyllableTemplate(
		SyllableSlot{Required: false, MaxSize: 0},
		SyllableSlot{Required: true, MaxSize: 1},
		SyllableSlot{Required: false, MaxSize: 0},
		3,
	)

	// CCV - Consonant cluster + Vowel
	TemplateCCV = NewSyllableTemplate(
		SyllableSlot{Required: true, MaxSize: 2},
		SyllableSlot{Required: true, MaxSize: 1},
		SyllableSlot{Required: false, MaxSize: 0},
		5,
	)
)

// PhonotacticRule defines the contract for checking whether a syllable is valid.
type PhonotacticRule interface {
	// Validate checks if the given phoneme sequence follows this rule.
	// Returns true if the sequence is valid according to this rule.
	Validate(seq []phoneme.Phoneme) bool
}

// Phonology holds a phoneme inventory, syllable templates, and phonotactic rules.
// It provides methods for generating and validating syllables.
type Phonology struct {
	pool      *phoneme.PhonemePool
	templates []SyllableTemplate
	rules     []PhonotacticRule
}

// NewPhonology creates a new phonology with the given phoneme pool.
func NewPhonology(pool *phoneme.PhonemePool) *Phonology {
	return &Phonology{
		pool:      pool,
		templates: make([]SyllableTemplate, 0),
		rules:     make([]PhonotacticRule, 0),
	}
}

// AddTemplate adds a syllable template to the phonology with the specified weight.
func (p *Phonology) AddTemplate(template SyllableTemplate, weight int) {
	template.Weight = weight
	p.templates = append(p.templates, template)
}

// PickTemplate selects a syllable template using weighted random selection.
func (p *Phonology) PickTemplate(rng *rand.Rand) SyllableTemplate {
	if len(p.templates) == 0 {
		// Return a default CV template if none are available
		return TemplateCV
	}

	var total int
	for _, t := range p.templates {
		total += t.Weight
	}

	if total == 0 {
		return p.templates[rng.IntN(len(p.templates))]
	}

	target := rng.IntN(total)
	var cumulative int
	for _, t := range p.templates {
		cumulative += t.Weight
		if cumulative > target {
			return t
		}
	}

	return p.templates[len(p.templates)-1]
}

// AddRule adds a phonotactic rule to the phonology.
func (p *Phonology) AddRule(rule PhonotacticRule) {
	p.rules = append(p.rules, rule)
}

// ApplyRules validates a phoneme sequence against all phonotactic rules.
// Returns true if all rules pass, false if any rule fails.
func (p *Phonology) ApplyRules(seq []phoneme.Phoneme) bool {
	for _, rule := range p.rules {
		if !rule.Validate(seq) {
			return false
		}
	}
	return true
}

// GenerateSyllable generates a valid syllable using the phonology's templates and rules.
// It selects a template, fills the slots with appropriate phonemes, and validates the result.
func (p *Phonology) GenerateSyllable(rng *rand.Rand) string {
	maxAttempts := 100 // Prevent infinite loops

	for attempt := 0; attempt < maxAttempts; attempt++ {
		template := p.PickTemplate(rng)
		syllable := p.fillTemplate(template, rng)

		if p.ValidateSyllable(syllable) {
			return p.phonemesToString(syllable)
		}
	}

	// Fallback: return a simple CV syllable
	return "a"
}

// ValidateSyllable checks if a phoneme sequence forms a valid syllable.
func (p *Phonology) ValidateSyllable(seq []phoneme.Phoneme) bool {
	if len(seq) == 0 {
		return false
	}

	// Check if we have at least one vowel (nucleus requirement)
	hasVowel := false
	for _, p := range seq {
		if p.Type == phoneme.PhonemeTypeVowel {
			hasVowel = true
			break
		}
	}

	if !hasVowel {
		return false
	}

	// Apply all phonotactic rules
	return p.ApplyRules(seq)
}

// fillTemplate fills a syllable template with phonemes from the pool.
func (p *Phonology) fillTemplate(template SyllableTemplate, rng *rand.Rand) []phoneme.Phoneme {
	var result []phoneme.Phoneme

	// Fill onset
	if template.Onset.Required && template.Onset.MaxSize > 0 {
		size := rng.IntN(template.Onset.MaxSize) + 1
		for i := 0; i < size; i++ {
			if phoneme := p.pool.Consonants.WeightedChoice(rng); phoneme != nil {
				result = append(result, *phoneme)
			}
		}
	}

	// Fill nucleus
	if template.Nucleus.Required && template.Nucleus.MaxSize > 0 {
		size := rng.IntN(template.Nucleus.MaxSize) + 1
		for i := 0; i < size; i++ {
			if phoneme := p.pool.Vowels.WeightedChoice(rng); phoneme != nil {
				result = append(result, *phoneme)
			}
		}
	}

	// Fill coda
	if template.Coda.Required && template.Coda.MaxSize > 0 {
		size := rng.IntN(template.Coda.MaxSize) + 1
		for i := 0; i < size; i++ {
			if phoneme := p.pool.Consonants.WeightedChoice(rng); phoneme != nil {
				result = append(result, *phoneme)
			}
		}
	}

	return result
}

// phonemesToString converts a slice of phonemes to a string representation.
func (p *Phonology) phonemesToString(phonemes []phoneme.Phoneme) string {
	var result string
	for _, p := range phonemes {
		result += p.Symbol
	}
	return result
}

// GetPool returns the phoneme pool used by this phonology.
func (p *Phonology) GetPool() *phoneme.PhonemePool {
	return p.pool
}

// GetTemplates returns a copy of the syllable templates.
func (p *Phonology) GetTemplates() []SyllableTemplate {
	result := make([]SyllableTemplate, len(p.templates))
	copy(result, p.templates)
	return result
}

// GetRules returns the number of phonotactic rules.
func (p *Phonology) GetRuleCount() int {
	return len(p.rules)
}
