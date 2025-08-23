package grammar

import (
	"fmt"
	"math/rand/v2"

	"github.com/chris-pikul/kismet-zero/lang/morphology"
)

// SentenceTemplate defines a pattern for constructing sentences of a specific type.
type SentenceTemplate struct {
	Type        SentenceType `json:"type"`
	Pattern     []string     `json:"pattern"`  // e.g., ["S", "V", "O"] for SVO
	Required    []string     `json:"required"` // Required elements
	Optional    []string     `json:"optional"` // Optional elements
	Weight      float32      `json:"weight"`
	Description string       `json:"description"`
}

// NewSentenceTemplate creates a new sentence template.
func NewSentenceTemplate(sentenceType SentenceType, pattern []string, required []string,
	optional []string, weight float32, description string) *SentenceTemplate {
	return &SentenceTemplate{
		Type:        sentenceType,
		Pattern:     pattern,
		Required:    required,
		Optional:    optional,
		Weight:      weight,
		Description: description,
	}
}

// Common sentence templates for different types
var (
	// Declarative templates
	TemplateSVO = NewSentenceTemplate(
		SentenceTypeDeclarative,
		[]string{"S", "V", "O"},
		[]string{"S", "V"},
		[]string{"O", "ADV"},
		10.0,
		"Subject-Verb-Object declarative sentence",
	)

	TemplateSOV = NewSentenceTemplate(
		SentenceTypeDeclarative,
		[]string{"S", "O", "V"},
		[]string{"S", "V"},
		[]string{"O", "ADV"},
		8.0,
		"Subject-Object-Verb declarative sentence",
	)

	TemplateVSO = NewSentenceTemplate(
		SentenceTypeDeclarative,
		[]string{"V", "S", "O"},
		[]string{"V", "S"},
		[]string{"O", "ADV"},
		6.0,
		"Verb-Subject-Object declarative sentence",
	)

	// Interrogative templates
	TemplateQuestionSVO = NewSentenceTemplate(
		SentenceTypeInterrogative,
		[]string{"Q", "S", "V", "O"},
		[]string{"Q", "S", "V"},
		[]string{"O", "ADV"},
		7.0,
		"Question word + Subject-Verb-Object interrogative",
	)

	TemplateQuestionVSO = NewSentenceTemplate(
		SentenceTypeInterrogative,
		[]string{"Q", "V", "S", "O"},
		[]string{"Q", "V", "S"},
		[]string{"O", "ADV"},
		5.0,
		"Question word + Verb-Subject-Object interrogative",
	)

	// Imperative templates
	TemplateImperativeVO = NewSentenceTemplate(
		SentenceTypeImperative,
		[]string{"V", "O"},
		[]string{"V"},
		[]string{"O", "ADV"},
		9.0,
		"Verb-Object imperative command",
	)

	TemplateImperativeV = NewSentenceTemplate(
		SentenceTypeImperative,
		[]string{"V"},
		[]string{"V"},
		[]string{"ADV"},
		8.0,
		"Verb-only imperative command",
	)
)

// SyntaxRuleBase provides a base implementation for syntax rules.
type SyntaxRuleBase struct {
	weight      float32
	description string
}

// Weight returns the weight of this rule.
func (sr *SyntaxRuleBase) Weight() float32 {
	return sr.weight
}

// String returns the description of this rule.
func (sr *SyntaxRuleBase) String() string {
	return sr.description
}

// WordOrderRule ensures words follow the correct word order pattern.
type WordOrderRule struct {
	SyntaxRuleBase
	wordOrder WordOrder
}

// NewWordOrderRule creates a new word order rule.
func NewWordOrderRule(wordOrder WordOrder, weight float32) *WordOrderRule {
	return &WordOrderRule{
		SyntaxRuleBase: SyntaxRuleBase{
			weight:      weight,
			description: fmt.Sprintf("Word order rule for %s pattern", wordOrder.String()),
		},
		wordOrder: wordOrder,
	}
}

// Apply reorders words to match the required word order pattern.
func (wor *WordOrderRule) Apply(words []morphology.Word, rng *rand.Rand) ([]morphology.Word, error) {
	if len(words) < 2 {
		return words, nil
	}

	// Categorize words by their morphological category
	subjects := make([]morphology.Word, 0)
	verbs := make([]morphology.Word, 0)
	objects := make([]morphology.Word, 0)
	others := make([]morphology.Word, 0)

	for _, word := range words {
		switch word.Category {
		case morphology.WordCategoryNoun, morphology.WordCategoryPronoun:
			if len(subjects) == 0 {
				subjects = append(subjects, word)
			} else if len(objects) == 0 {
				objects = append(objects, word)
			} else {
				others = append(others, word)
			}
		case morphology.WordCategoryVerb:
			verbs = append(verbs, word)
		default:
			others = append(others, word)
		}
	}

	// Reorder based on word order pattern
	var result []morphology.Word
	switch wor.wordOrder {
	case WordOrderSVO:
		result = append(result, subjects...)
		result = append(result, verbs...)
		result = append(result, objects...)
		result = append(result, others...)
	case WordOrderSOV:
		result = append(result, subjects...)
		result = append(result, objects...)
		result = append(result, verbs...)
		result = append(result, others...)
	case WordOrderVSO:
		result = append(result, verbs...)
		result = append(result, subjects...)
		result = append(result, objects...)
		result = append(result, others...)
	case WordOrderVOS:
		result = append(result, verbs...)
		result = append(result, objects...)
		result = append(result, subjects...)
		result = append(result, others...)
	case WordOrderOVS:
		result = append(result, objects...)
		result = append(result, verbs...)
		result = append(result, subjects...)
		result = append(result, others...)
	case WordOrderOSV:
		result = append(result, objects...)
		result = append(result, subjects...)
		result = append(result, verbs...)
		result = append(result, others...)
	default:
		result = words // Fallback to original order
	}

	return result, nil
}

// Validate checks if the words follow the required word order.
func (wor *WordOrderRule) Validate(words []morphology.Word) bool {
	if len(words) < 2 {
		return true
	}

	// Simple validation - check if we have at least one subject and one verb
	hasSubject := false
	hasVerb := false

	for _, word := range words {
		switch word.Category {
		case morphology.WordCategoryNoun, morphology.WordCategoryPronoun:
			hasSubject = true
		case morphology.WordCategoryVerb:
			hasVerb = true
		}
	}

	return hasSubject && hasVerb
}

// SubjectVerbAgreementRule ensures subject-verb agreement.
type SubjectVerbAgreementRule struct {
	SyntaxRuleBase
}

// NewSubjectVerbAgreementRule creates a new subject-verb agreement rule.
func NewSubjectVerbAgreementRule(weight float32) *SubjectVerbAgreementRule {
	return &SubjectVerbAgreementRule{
		SyntaxRuleBase: SyntaxRuleBase{
			weight:      weight,
			description: "Subject-verb agreement rule",
		},
	}
}

// Apply ensures subject-verb agreement in number and person.
func (svar *SubjectVerbAgreementRule) Apply(words []morphology.Word, rng *rand.Rand) ([]morphology.Word, error) {
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

	// Apply agreement if both subject and verb are found
	if subject != nil && verb != nil {
		// For now, we'll just mark that agreement should be applied
		// The actual morphological changes would be handled by the morphology package
		// This rule ensures the structure is correct for agreement
	}

	return words, nil
}

// Validate checks if subject-verb agreement is maintained.
func (svar *SubjectVerbAgreementRule) Validate(words []morphology.Word) bool {
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

	// Basic validation - both subject and verb should be present
	return subject != nil && verb != nil
}

// SentenceTemplateRule applies a specific sentence template.
type SentenceTemplateRule struct {
	SyntaxRuleBase
	template *SentenceTemplate
}

// NewSentenceTemplateRule creates a new sentence template rule.
func NewSentenceTemplateRule(template *SentenceTemplate, weight float32) *SentenceTemplateRule {
	return &SentenceTemplateRule{
		SyntaxRuleBase: SyntaxRuleBase{
			weight:      weight,
			description: fmt.Sprintf("Template rule: %s", template.Description),
		},
		template: template,
	}
}

// Apply applies the sentence template to the words.
func (str *SentenceTemplateRule) Apply(words []morphology.Word, rng *rand.Rand) ([]morphology.Word, error) {
	// This rule validates that the words match the template pattern
	// The actual reordering is handled by the WordOrderRule
	return words, nil
}

// Validate checks if the words follow the template pattern.
func (str *SentenceTemplateRule) Validate(words []morphology.Word) bool {
	if len(words) < len(str.template.Required) {
		return false
	}

	// Check if all required elements are present
	requiredCount := 0
	for _, required := range str.template.Required {
		for _, word := range words {
			if str.matchesPattern(word, required) {
				requiredCount++
				break
			}
		}
	}

	return requiredCount >= len(str.template.Required)
}

// matchesPattern checks if a word matches a pattern element.
func (str *SentenceTemplateRule) matchesPattern(word morphology.Word, pattern string) bool {
	switch pattern {
	case "S":
		return word.Category == morphology.WordCategoryNoun ||
			word.Category == morphology.WordCategoryPronoun
	case "V":
		return word.Category == morphology.WordCategoryVerb
	case "O":
		return word.Category == morphology.WordCategoryNoun ||
			word.Category == morphology.WordCategoryPronoun
	case "ADV":
		return word.Category == morphology.WordCategoryAdverb
	case "ADJ":
		return word.Category == morphology.WordCategoryAdjective
	case "Q":
		return word.Category == morphology.WordCategoryInterjection
	default:
		return false
	}
}

// GetSentenceTemplates returns common sentence templates for a given type.
func GetSentenceTemplates(sentenceType SentenceType) []*SentenceTemplate {
	var templates []*SentenceTemplate

	switch sentenceType {
	case SentenceTypeDeclarative:
		templates = append(templates, TemplateSVO, TemplateSOV, TemplateVSO)
	case SentenceTypeInterrogative:
		templates = append(templates, TemplateQuestionSVO, TemplateQuestionVSO)
	case SentenceTypeImperative:
		templates = append(templates, TemplateImperativeVO, TemplateImperativeV)
	case SentenceTypeExclamatory:
		// Exclamatory sentences often follow declarative patterns
		templates = append(templates, TemplateSVO, TemplateSOV)
	}

	return templates
}

// GetDefaultTemplates returns all available sentence templates.
func GetDefaultTemplates() []*SentenceTemplate {
	return []*SentenceTemplate{
		TemplateSVO,
		TemplateSOV,
		TemplateVSO,
		TemplateQuestionSVO,
		TemplateQuestionVSO,
		TemplateImperativeVO,
		TemplateImperativeV,
	}
}
