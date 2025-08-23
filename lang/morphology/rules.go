package morphology

import (
	"fmt"
	"math/rand/v2"

	"github.com/chris-pikul/kismet-zero/lang/phoneme"
	"github.com/chris-pikul/kismet-zero/lang/phonology"
)

// MorphologicalRule defines the contract for applying morphological rules
// to morphemes during word formation.
type MorphologicalRule interface {
	// Apply applies the rule to the input morphemes and returns the result.
	// The phonology is provided for validation and additional constraints.
	Apply(morphemes []Morpheme, phonology *phonology.Phonology) ([]Morpheme, error)

	// Validate checks if the input morphemes can have this rule applied.
	Validate(input []Morpheme) bool

	// Weight returns the probability weight for this rule.
	Weight() float32

	// Name returns a human-readable name for this rule.
	Name() string
}

// DerivationalRule represents a rule that changes the meaning or category
// of a word by adding or modifying morphemes.
type DerivationalRule struct {
	name       string
	inputType  MorphemeType
	outputType MorphemeType
	pattern    string
	weight     float32
}

// NewDerivationalRule creates a new derivational rule.
func NewDerivationalRule(
	name string,
	inputType, outputType MorphemeType,
	pattern string,
	weight float32,
) *DerivationalRule {
	return &DerivationalRule{
		name:       name,
		inputType:  inputType,
		outputType: outputType,
		pattern:    pattern,
		weight:     weight,
	}
}

// Apply implements MorphologicalRule.Apply for derivational rules.
func (dr *DerivationalRule) Apply(
	morphemes []Morpheme,
	phonology *phonology.Phonology,
) ([]Morpheme, error) {
	if !dr.Validate(morphemes) {
		return nil, fmt.Errorf("derivational rule %s cannot be applied to input morphemes", dr.name)
	}

	// For now, we'll implement a simple pattern where we add a suffix
	// In a more sophisticated system, this would handle various patterns
	if dr.outputType == MorphemeTypeSuffix {
		// Create a simple suffix morpheme
		suffix := &Morpheme{
			ID:        fmt.Sprintf("deriv_%s", dr.name),
			Type:      MorphemeTypeSuffix,
			Meaning:   dr.pattern,
			Phonemes:  []phoneme.Phoneme{}, // Would be generated based on pattern
			Syllables: []string{},
			Weight:    1.0,
			Frequency: MorphemeFrequencyUncommon,
		}

		// Add the suffix to the morphemes
		result := make([]Morpheme, len(morphemes)+1)
		copy(result, morphemes)
		result[len(morphemes)] = *suffix

		return result, nil
	}

	return morphemes, nil
}

// Validate implements MorphologicalRule.Validate for derivational rules.
func (dr *DerivationalRule) Validate(input []Morpheme) bool {
	if len(input) == 0 {
		return false
	}

	// Check if the first morpheme matches the expected input type
	return input[0].Type == dr.inputType
}

// Weight implements MorphologicalRule.Weight for derivational rules.
func (dr *DerivationalRule) Weight() float32 {
	return dr.weight
}

// Name implements MorphologicalRule.Name for derivational rules.
func (dr *DerivationalRule) Name() string {
	return dr.name
}

// InflectionalRule represents a rule that adds grammatical information
// without changing the core meaning of a word.
type InflectionalRule struct {
	name     string
	category WordCategory
	features []string
	pattern  string
	weight   float32
}

// NewInflectionalRule creates a new inflectional rule.
func NewInflectionalRule(
	name string,
	category WordCategory,
	features []string,
	pattern string,
	weight float32,
) *InflectionalRule {
	return &InflectionalRule{
		name:     name,
		category: category,
		features: features,
		pattern:  pattern,
		weight:   weight,
	}
}

// Apply implements MorphologicalRule.Apply for inflectional rules.
func (ir *InflectionalRule) Apply(
	morphemes []Morpheme,
	phonology *phonology.Phonology,
) ([]Morpheme, error) {
	if !ir.Validate(morphemes) {
		return nil, fmt.Errorf("inflectional rule %s cannot be applied to input morphemes", ir.name)
	}

	// Create an inflectional suffix
	suffix := &Morpheme{
		ID:        fmt.Sprintf("infl_%s", ir.name),
		Type:      MorphemeTypeSuffix,
		Meaning:   ir.pattern,
		Phonemes:  []phoneme.Phoneme{}, // Would be generated based on pattern
		Syllables: []string{},
		Weight:    1.0,
		Frequency: MorphemeFrequencyCommon,
	}

	// Add the suffix to the morphemes
	result := make([]Morpheme, len(morphemes)+1)
	copy(result, morphemes)
	result[len(morphemes)] = *suffix

	return result, nil
}

// Validate implements MorphologicalRule.Validate for inflectional rules.
func (ir *InflectionalRule) Validate(input []Morpheme) bool {
	if len(input) == 0 {
		return false
	}

	// For inflectional rules, we typically expect root morphemes
	return input[0].Type == MorphemeTypeRoot
}

// Weight implements MorphologicalRule.Weight for inflectional rules.
func (ir *InflectionalRule) Weight() float32 {
	return ir.weight
}

// Name implements MorphologicalRule.Name for inflectional rules.
func (ir *InflectionalRule) Name() string {
	return ir.name
}

// RuleSet manages a collection of morphological rules and provides
// methods for applying them to morphemes.
type RuleSet struct {
	derivational []MorphologicalRule
	inflectional []MorphologicalRule
}

// NewRuleSet creates a new rule set.
func NewRuleSet() *RuleSet {
	return &RuleSet{
		derivational: make([]MorphologicalRule, 0),
		inflectional: make([]MorphologicalRule, 0),
	}
}

// AddDerivationalRule adds a derivational rule to the rule set.
func (rs *RuleSet) AddDerivationalRule(rule MorphologicalRule) {
	rs.derivational = append(rs.derivational, rule)
}

// AddInflectionalRule adds an inflectional rule to the rule set.
func (rs *RuleSet) AddInflectionalRule(rule MorphologicalRule) {
	rs.inflectional = append(rs.inflectional, rule)
}

// GetApplicableRules returns all rules that can be applied to the given morphemes.
func (rs *RuleSet) GetApplicableRules(morphemes []Morpheme) []MorphologicalRule {
	var applicable []MorphologicalRule

	// Check derivational rules
	for _, rule := range rs.derivational {
		if rule.Validate(morphemes) {
			applicable = append(applicable, rule)
		}
	}

	// Check inflectional rules
	for _, rule := range rs.inflectional {
		if rule.Validate(morphemes) {
			applicable = append(applicable, rule)
		}
	}

	return applicable
}

// ApplyRandomRule randomly selects and applies an applicable rule to the morphemes.
func (rs *RuleSet) ApplyRandomRule(
	morphemes []Morpheme,
	phonology *phonology.Phonology,
	rng *rand.Rand,
) ([]Morpheme, error) {
	applicable := rs.GetApplicableRules(morphemes)
	if len(applicable) == 0 {
		return morphemes, nil // No rules to apply
	}

	// Select a rule using weighted random selection
	rule := rs.selectRule(applicable, rng)
	if rule == nil {
		return morphemes, nil
	}

	return rule.Apply(morphemes, phonology)
}

// selectRule selects a rule using weighted random selection.
func (rs *RuleSet) selectRule(rules []MorphologicalRule, rng *rand.Rand) MorphologicalRule {
	if len(rules) == 0 {
		return nil
	}

	// Calculate total weight
	var totalWeight float32
	for _, rule := range rules {
		totalWeight += rule.Weight()
	}

	if totalWeight == 0 {
		return rules[0] // Fallback to first rule
	}

	// Weighted selection
	target := rng.Float32() * totalWeight
	var cumulative float32
	for _, rule := range rules {
		cumulative += rule.Weight()
		if cumulative >= target {
			return rule
		}
	}

	return rules[len(rules)-1] // Fallback to last rule
}

// GetRuleCount returns the total number of rules in the rule set.
func (rs *RuleSet) GetRuleCount() int {
	return len(rs.derivational) + len(rs.inflectional)
}

// GetDerivationalRuleCount returns the number of derivational rules.
func (rs *RuleSet) GetDerivationalRuleCount() int {
	return len(rs.derivational)
}

// GetInflectionalRuleCount returns the number of inflectional rules.
func (rs *RuleSet) GetInflectionalRuleCount() int {
	return len(rs.inflectional)
}
