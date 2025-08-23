package grammar

import (
	"fmt"
	"math/rand/v2"

	"github.com/chris-pikul/kismet-zero/lang/morphology"
)

// AgreementRuleBase provides a base implementation for agreement rules.
type AgreementRuleBase struct {
	weight      float32
	description string
}

// Weight returns the weight of this rule.
func (ar *AgreementRuleBase) Weight() float32 {
	return ar.weight
}

// String returns the description of this rule.
func (ar *AgreementRuleBase) String() string {
	return ar.description
}

// CaseAgreementRule ensures proper case marking across sentence elements.
type CaseAgreementRule struct {
	AgreementRuleBase
	requiredCases []Case
}

// NewCaseAgreementRule creates a new case agreement rule.
func NewCaseAgreementRule(requiredCases []Case, weight float32) *CaseAgreementRule {
	return &CaseAgreementRule{
		AgreementRuleBase: AgreementRuleBase{
			weight:      weight,
			description: fmt.Sprintf("Case agreement rule for %v", requiredCases),
		},
		requiredCases: requiredCases,
	}
}

// Apply ensures proper case marking for nouns and pronouns.
func (car *CaseAgreementRule) Apply(words []morphology.Word, rng *rand.Rand) ([]morphology.Word, error) {
	if len(words) < 2 {
		return words, nil
	}

	// Find subject and object positions
	var subjectIndex int = -1
	var objectIndex int = -1

	for i, word := range words {
		switch word.Category {
		case morphology.WordCategoryNoun, morphology.WordCategoryPronoun:
			if subjectIndex == -1 {
				subjectIndex = i
			} else if objectIndex == -1 {
				objectIndex = i
			}
		}
	}

	// Apply case marking if we have both subject and object
	if subjectIndex != -1 && objectIndex != -1 {
		// Subject should be nominative, object should be accusative
		// This is a simplified case system - in practice, the morphology package
		// would handle the actual case inflections
		words[subjectIndex].Agreement.Case = CaseNominative.String()
		words[objectIndex].Agreement.Case = CaseAccusative.String()
	}

	return words, nil
}

// Validate checks if case agreement is maintained.
func (car *CaseAgreementRule) Validate(words []morphology.Word) bool {
	if len(words) < 2 {
		return true
	}

	// Check if required cases are present
	caseCount := make(map[string]int)
	for _, word := range words {
		if word.Category == morphology.WordCategoryNoun ||
			word.Category == morphology.WordCategoryPronoun {
			caseCount[word.Agreement.Case]++
		}
	}

	// Ensure we have at least the required cases
	for _, requiredCase := range car.requiredCases {
		if caseCount[requiredCase.String()] == 0 {
			return false
		}
	}

	return true
}

// NumberAgreementRule ensures number agreement across sentence elements.
type NumberAgreementRule struct {
	AgreementRuleBase
}

// NewNumberAgreementRule creates a new number agreement rule.
func NewNumberAgreementRule(weight float32) *NumberAgreementRule {
	return &NumberAgreementRule{
		AgreementRuleBase: AgreementRuleBase{
			weight:      weight,
			description: "Number agreement rule",
		},
	}
}

// Apply ensures number agreement between subject and verb.
func (nar *NumberAgreementRule) Apply(words []morphology.Word, rng *rand.Rand) ([]morphology.Word, error) {
	if len(words) < 2 {
		return words, nil
	}

	// Find subject and verb
	var subject *morphology.Word
	var verb *morphology.Word

	for i := range words {
		switch words[i].Category {
		case morphology.WordCategoryNoun, morphology.WordCategoryPronoun:
			if subject == nil {
				subject = &words[i]
			}
		case morphology.WordCategoryVerb:
			if verb == nil {
				verb = &words[i]
			}
		}
	}

	// Apply number agreement if both subject and verb are found
	if subject != nil && verb != nil {
		// Verb should agree with subject in number
		verb.Agreement.Number = subject.Agreement.Number
	}

	return words, nil
}

// Validate checks if number agreement is maintained.
func (nar *NumberAgreementRule) Validate(words []morphology.Word) bool {
	if len(words) < 2 {
		return true
	}

	// Find subject and verb
	var subject *morphology.Word
	var verb *morphology.Word

	for i := range words {
		switch words[i].Category {
		case morphology.WordCategoryNoun, morphology.WordCategoryPronoun:
			if subject == nil {
				subject = &words[i]
			}
		case morphology.WordCategoryVerb:
			if verb == nil {
				verb = &words[i]
			}
		}
	}

	// Check agreement if both are found
	if subject != nil && verb != nil {
		return subject.Agreement.Number == verb.Agreement.Number
	}

	return true
}

// GenderAgreementRule ensures gender agreement across sentence elements.
type GenderAgreementRule struct {
	AgreementRuleBase
}

// NewGenderAgreementRule creates a new gender agreement rule.
func NewGenderAgreementRule(weight float32) *GenderAgreementRule {
	return &GenderAgreementRule{
		AgreementRuleBase: AgreementRuleBase{
			weight:      weight,
			description: "Gender agreement rule",
		},
	}
}

// Apply ensures gender agreement between nouns and adjectives.
func (gar *GenderAgreementRule) Apply(words []morphology.Word, rng *rand.Rand) ([]morphology.Word, error) {
	if len(words) < 2 {
		return words, nil
	}

	// Find nouns and adjectives
	var nouns []*morphology.Word
	var adjectives []*morphology.Word

	for i := range words {
		switch words[i].Category {
		case morphology.WordCategoryNoun:
			nouns = append(nouns, &words[i])
		case morphology.WordCategoryAdjective:
			adjectives = append(adjectives, &words[i])
		}
	}

	// Apply gender agreement if we have both nouns and adjectives
	if len(nouns) > 0 && len(adjectives) > 0 {
		// Adjectives should agree with the nearest noun in gender
		// This is a simplified system - in practice, more complex rules would apply
		for _, adj := range adjectives {
			if len(nouns) > 0 {
				adj.Agreement.Gender = nouns[0].Agreement.Gender
			}
		}
	}

	return words, nil
}

// Validate checks if gender agreement is maintained.
func (gar *GenderAgreementRule) Validate(words []morphology.Word) bool {
	if len(words) < 2 {
		return true
	}

	// Find nouns and adjectives
	var nouns []*morphology.Word
	var adjectives []*morphology.Word

	for i := range words {
		switch words[i].Category {
		case morphology.WordCategoryNoun:
			nouns = append(nouns, &words[i])
		case morphology.WordCategoryAdjective:
			adjectives = append(adjectives, &words[i])
		}
	}

	// Check agreement if we have both
	if len(nouns) > 0 && len(adjectives) > 0 {
		// Adjectives should agree with nouns in gender
		// This is a simplified check
		for _, adj := range adjectives {
			hasAgreement := false
			for _, noun := range nouns {
				if adj.Agreement.Gender == noun.Agreement.Gender {
					hasAgreement = true
					break
				}
			}
			if !hasAgreement {
				return false
			}
		}
	}

	return true
}

// TenseAspectAgreementRule ensures tense and aspect agreement.
type TenseAspectAgreementRule struct {
	AgreementRuleBase
}

// NewTenseAspectAgreementRule creates a new tense-aspect agreement rule.
func NewTenseAspectAgreementRule(weight float32) *TenseAspectAgreementRule {
	return &TenseAspectAgreementRule{
		AgreementRuleBase: AgreementRuleBase{
			weight:      weight,
			description: "Tense-aspect agreement rule",
		},
	}
}

// Apply ensures consistent tense and aspect marking.
func (tar *TenseAspectAgreementRule) Apply(words []morphology.Word, rng *rand.Rand) ([]morphology.Word, error) {
	if len(words) < 2 {
		return words, nil
	}

	// Find verbs
	var verbs []*morphology.Word

	for i := range words {
		if words[i].Category == morphology.WordCategoryVerb {
			verbs = append(verbs, &words[i])
		}
	}

	// Apply consistent tense and aspect if we have multiple verbs
	if len(verbs) > 1 {
		// Use the first verb's tense and aspect as the base
		baseTense := verbs[0].Agreement.Tense
		baseAspect := verbs[0].Agreement.Aspect

		for _, verb := range verbs[1:] {
			verb.Agreement.Tense = baseTense
			verb.Agreement.Aspect = baseAspect
		}
	}

	return words, nil
}

// Validate checks if tense-aspect agreement is maintained.
func (tar *TenseAspectAgreementRule) Validate(words []morphology.Word) bool {
	if len(words) < 2 {
		return true
	}

	// Find verbs
	var verbs []*morphology.Word

	for i := range words {
		if words[i].Category == morphology.WordCategoryVerb {
			verbs = append(verbs, &words[i])
		}
	}

	// Check agreement if we have multiple verbs
	if len(verbs) > 1 {
		baseTense := verbs[0].Agreement.Tense
		baseAspect := verbs[0].Agreement.Aspect

		for _, verb := range verbs[1:] {
			if verb.Agreement.Tense != baseTense || verb.Agreement.Aspect != baseAspect {
				return false
			}
		}
	}

	return true
}

// PersonAgreementRule ensures person agreement between subject and verb.
type PersonAgreementRule struct {
	AgreementRuleBase
}

// NewPersonAgreementRule creates a new person agreement rule.
func NewPersonAgreementRule(weight float32) *PersonAgreementRule {
	return &PersonAgreementRule{
		AgreementRuleBase: AgreementRuleBase{
			weight:      weight,
			description: "Person agreement rule",
		},
	}
}

// Apply ensures person agreement between subject and verb.
func (par *PersonAgreementRule) Apply(words []morphology.Word, rng *rand.Rand) ([]morphology.Word, error) {
	if len(words) < 2 {
		return words, nil
	}

	// Find subject and verb
	var subject *morphology.Word
	var verb *morphology.Word

	for i := range words {
		switch words[i].Category {
		case morphology.WordCategoryNoun, morphology.WordCategoryPronoun:
			if subject == nil {
				subject = &words[i]
			}
		case morphology.WordCategoryVerb:
			if verb == nil {
				verb = &words[i]
			}
		}
	}

	// Apply person agreement if both subject and verb are found
	if subject != nil && verb != nil {
		// Verb should agree with subject in person
		// This is a simplified system - in practice, the morphology package
		// would handle the actual person inflections
		verb.Agreement.Person = subject.Agreement.Person
	}

	return words, nil
}

// Validate checks if person agreement is maintained.
func (par *PersonAgreementRule) Validate(words []morphology.Word) bool {
	if len(words) < 2 {
		return true
	}

	// Find subject and verb
	var subject *morphology.Word
	var verb *morphology.Word

	for i := range words {
		switch words[i].Category {
		case morphology.WordCategoryNoun, morphology.WordCategoryPronoun:
			if subject == nil {
				subject = &words[i]
			}
		case morphology.WordCategoryVerb:
			if verb == nil {
				verb = &words[i]
			}
		}
	}

	// Check agreement if both are found
	if subject != nil && verb != nil {
		return subject.Agreement.Person == verb.Agreement.Person
	}

	return true
}

// GetDefaultAgreementRules returns a set of common agreement rules.
func GetDefaultAgreementRules() []AgreementRule {
	return []AgreementRule{
		NewCaseAgreementRule([]Case{CaseNominative, CaseAccusative}, 10.0),
		NewNumberAgreementRule(9.0),
		NewGenderAgreementRule(8.0),
		NewTenseAspectAgreementRule(9.0),
		NewPersonAgreementRule(8.0),
	}
}

// ApplyAgreementRules applies all agreement rules to the given words.
func ApplyAgreementRules(words []morphology.Word, rules []AgreementRule, rng *rand.Rand) ([]morphology.Word, error) {
	result := words

	for _, rule := range rules {
		var err error
		result, err = rule.Apply(result, rng)
		if err != nil {
			return nil, fmt.Errorf("failed to apply agreement rule %s: %w", rule, err)
		}
	}

	return result, nil
}

// ValidateAgreementRules validates that all agreement rules are satisfied.
func ValidateAgreementRules(words []morphology.Word, rules []AgreementRule) bool {
	for _, rule := range rules {
		if !rule.Validate(words) {
			return false
		}
	}
	return true
}
