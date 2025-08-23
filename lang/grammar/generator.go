package grammar

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/chris-pikul/kismet-zero/lang/morphology"
)

// Observer interface for tracking sentence generation metadata.
type Observer interface {
	OnTemplateChosen(id string)
	OnToken(index int, source string) // "verb","subj","obj","iobj","adv-time"
}

// RealizationMeta contains metadata about sentence generation for interlingua integration.
type RealizationMeta struct {
	TemplateID    string         `json:"templateId"`
	TokenToSource map[int]string `json:"tokenToSource"` // e.g., "verb","subj","obj","iobj","adv-time"
}

// Options contains optional configuration for sentence generation.
type Options struct {
	Observer Observer
}

// SentenceGenerator creates sentences using the grammar rules and morphological building blocks.
type SentenceGenerator struct {
	grammar *Grammar
}

// NewSentenceGenerator creates a new sentence generator with the given grammar.
func NewSentenceGenerator(grammar *Grammar) *SentenceGenerator {
	return &SentenceGenerator{
		grammar: grammar,
	}
}

// GenerateSentence creates a sentence from a list of words, applying grammar rules.
func (sg *SentenceGenerator) GenerateSentence(
	words []morphology.Word,
	sentenceType SentenceType,
	rng *rand.Rand,
	opts *Options,
) (*Sentence, error) {
	if len(words) < 2 {
		return nil, fmt.Errorf("need at least 2 words to form a sentence")
	}

	// Validate that we have the required word categories for the sentence type
	if err := sg.validateWordCategories(words, sentenceType); err != nil {
		return nil, fmt.Errorf("word categories validation failed: %w", err)
	}

	// Apply syntax rules
	result, err := sg.applySyntaxRules(words, rng)
	if err != nil {
		return nil, fmt.Errorf("failed to apply syntax rules: %w", err)
	}

	// Notify observer about template selection and token sources
	if opts != nil && opts.Observer != nil {
		// Determine template ID based on word order and sentence type
		templateID := sg.determineTemplateID(result, sentenceType)
		opts.Observer.OnTemplateChosen(templateID)

		// Map tokens to their sources
		for i, word := range result {
			source := sg.mapWordToSource(word, i, result)
			opts.Observer.OnToken(i, source)
		}
	}

	// Apply agreement rules
	result, err = sg.applyAgreementRules(result, rng)
	if err != nil {
		return nil, fmt.Errorf("failed to apply agreement rules: %w", err)
	}

	// Generate sentence structure description
	structure := sg.generateStructureDescription(result, sentenceType)

	// Combine meanings for the overall sentence meaning
	meaning := sg.combineSentenceMeaning(result)

	// Calculate sentence weight
	weight := sg.calculateSentenceWeight(result)

	// Create the sentence
	sentence := &Sentence{
		ID:        sg.generateSentenceID(result, sentenceType),
		Type:      sentenceType,
		Words:     result,
		Structure: structure,
		Meaning:   meaning,
		Culture:   sg.grammar.GetCulture(),
		Weight:    weight,
	}

	return sentence, nil
}

// GeneratePhrase creates a grammatical phrase from a list of words.
func (sg *SentenceGenerator) GeneratePhrase(
	words []morphology.Word,
	phraseType string,
	rng *rand.Rand,
) (*Phrase, error) {
	if len(words) == 0 {
		return nil, fmt.Errorf("need at least one word to form a phrase")
	}

	// Find the head word (main word of the phrase)
	var head *morphology.Word
	var modifiers []morphology.Word

	for i, word := range words {
		if i == 0 || word.Category == morphology.WordCategoryNoun ||
			word.Category == morphology.WordCategoryVerb {
			if head == nil {
				head = &words[i]
			} else {
				modifiers = append(modifiers, word)
			}
		} else {
			modifiers = append(modifiers, word)
		}
	}

	// If no clear head was found, use the first word
	if head == nil {
		head = &words[0]
	}

	// Create agreement features based on the head word
	// Convert string agreement features to enum types
	var caseVal Case
	var numberVal Number
	var genderVal Gender

	// Parse case
	switch head.Agreement.Case {
	case "nominative":
		caseVal = CaseNominative
	case "accusative":
		caseVal = CaseAccusative
	case "genitive":
		caseVal = CaseGenitive
	case "dative":
		caseVal = CaseDative
	case "locative":
		caseVal = CaseLocative
	case "instrumental":
		caseVal = CaseInstrumental
	case "vocative":
		caseVal = CaseVocative
	case "ablative":
		caseVal = CaseAblative
	default:
		caseVal = CaseNominative
	}

	// Parse number
	switch head.Agreement.Number {
	case "singular":
		numberVal = NumberSingular
	case "dual":
		numberVal = NumberDual
	case "plural":
		numberVal = NumberPlural
	default:
		numberVal = NumberSingular
	}

	// Parse gender
	switch head.Agreement.Gender {
	case "masculine":
		genderVal = GenderMasculine
	case "feminine":
		genderVal = GenderFeminine
	case "neuter":
		genderVal = GenderNeuter
	case "animate":
		genderVal = GenderAnimate
	case "inanimate":
		genderVal = GenderInanimate
	default:
		genderVal = GenderNeuter
	}

	agreement := AgreementFeature{
		Case:   caseVal,
		Number: numberVal,
		Gender: genderVal,
	}

	phrase := &Phrase{
		Type:      phraseType,
		Words:     words,
		Head:      head,
		Modifiers: modifiers,
		Agreement: agreement,
	}

	return phrase, nil
}

// validateWordCategories ensures the words have the required categories for the sentence type.
func (sg *SentenceGenerator) validateWordCategories(words []morphology.Word, sentenceType SentenceType) error {
	hasSubject := false
	hasVerb := false
	hasObject := false

	for _, word := range words {
		switch word.Category {
		case morphology.WordCategoryNoun, morphology.WordCategoryPronoun:
			if !hasSubject {
				hasSubject = true
			} else if !hasObject {
				hasObject = true
			}
		case morphology.WordCategoryVerb:
			hasVerb = true
		}
	}

	// Check requirements based on sentence type
	switch sentenceType {
	case SentenceTypeDeclarative:
		if !hasSubject || !hasVerb {
			return fmt.Errorf("declarative sentence requires subject and verb")
		}
	case SentenceTypeInterrogative:
		if !hasSubject || !hasVerb {
			return fmt.Errorf("interrogative sentence requires subject and verb")
		}
	case SentenceTypeImperative:
		if !hasVerb {
			return fmt.Errorf("imperative sentence requires verb")
		}
	case SentenceTypeExclamatory:
		if !hasSubject || !hasVerb {
			return fmt.Errorf("exclamatory sentence requires subject and verb")
		}
	}

	return nil
}

// applySyntaxRules applies all syntax rules to the words.
func (sg *SentenceGenerator) applySyntaxRules(words []morphology.Word, rng *rand.Rand) ([]morphology.Word, error) {
	result := words

	// Apply word order rule first
	wordOrderRule := NewWordOrderRule(sg.grammar.GetWordOrder(), 10.0)
	result, err := wordOrderRule.Apply(result, rng)
	if err != nil {
		return nil, fmt.Errorf("failed to apply word order rule: %w", err)
	}

	// Apply other syntax rules
	for _, rule := range sg.grammar.GetSyntaxRules() {
		result, err = rule.Apply(result, rng)
		if err != nil {
			return nil, fmt.Errorf("failed to apply syntax rule %s: %w", rule, err)
		}
	}

	return result, nil
}

// determineTemplateID determines the template ID based on word order and sentence type.
func (sg *SentenceGenerator) determineTemplateID(words []morphology.Word, sentenceType SentenceType) string {
	// Determine word order pattern
	var pattern []string
	for _, word := range words {
		switch word.Category {
		case morphology.WordCategoryNoun, morphology.WordCategoryPronoun:
			if len(pattern) == 0 {
				pattern = append(pattern, "S") // Subject
			} else {
				pattern = append(pattern, "O") // Object
			}
		case morphology.WordCategoryVerb:
			pattern = append(pattern, "V") // Verb
		case morphology.WordCategoryAdjective:
			pattern = append(pattern, "ADJ") // Adjective
		case morphology.WordCategoryAdverb:
			pattern = append(pattern, "ADV") // Adverb
		case morphology.WordCategoryInterjection:
			pattern = append(pattern, "Q") // Question/Interjection
		default:
			pattern = append(pattern, "X") // Unknown/Other
		}
	}

	patternStr := strings.Join(pattern, "")
	return fmt.Sprintf("%s.%s", sentenceType.String(), patternStr)
}

// mapWordToSource maps a word to its semantic source for interlingua integration.
func (sg *SentenceGenerator) mapWordToSource(word morphology.Word, index int, allWords []morphology.Word) string {
	switch word.Category {
	case morphology.WordCategoryVerb:
		return "verb"
	case morphology.WordCategoryNoun, morphology.WordCategoryPronoun:
		// First noun/pronoun is subject, second is object
		nounCount := 0
		for i := 0; i < index; i++ {
			if allWords[i].Category == morphology.WordCategoryNoun ||
				allWords[i].Category == morphology.WordCategoryPronoun {
				nounCount++
			}
		}
		if nounCount == 0 {
			return "subj"
		} else if nounCount == 1 {
			return "obj"
		} else {
			return "iobj" // indirect object
		}
	case morphology.WordCategoryAdjective:
		return "adj"
	case morphology.WordCategoryAdverb:
		return "adv"
	case morphology.WordCategoryInterjection:
		return "q"
	default:
		return "x"
	}
}

// applyAgreementRules applies all agreement rules to the words.
func (sg *SentenceGenerator) applyAgreementRules(words []morphology.Word, rng *rand.Rand) ([]morphology.Word, error) {
	if sg.grammar.agreement == nil {
		return words, nil
	}

	return ApplyAgreementRules(words, sg.grammar.agreement.GetRules(), rng)
}

// generateStructureDescription creates a human-readable description of the sentence structure.
func (sg *SentenceGenerator) generateStructureDescription(words []morphology.Word, sentenceType SentenceType) string {
	var parts []string

	for _, word := range words {
		switch word.Category {
		case morphology.WordCategoryNoun, morphology.WordCategoryPronoun:
			if len(parts) == 0 {
				parts = append(parts, "S") // Subject
			} else {
				parts = append(parts, "O") // Object
			}
		case morphology.WordCategoryVerb:
			parts = append(parts, "V") // Verb
		case morphology.WordCategoryAdjective:
			parts = append(parts, "ADJ") // Adjective
		case morphology.WordCategoryAdverb:
			parts = append(parts, "ADV") // Adverb
		case morphology.WordCategoryInterjection:
			parts = append(parts, "Q") // Question/Interjection
		default:
			parts = append(parts, "X") // Unknown/Other
		}
	}

	structure := strings.Join(parts, "-")
	return fmt.Sprintf("%s: %s", sentenceType.String(), structure)
}

// combineSentenceMeaning creates a combined meaning from all words in the sentence.
func (sg *SentenceGenerator) combineSentenceMeaning(words []morphology.Word) string {
	if len(words) == 0 {
		return ""
	}

	if len(words) == 1 {
		return words[0].Meaning
	}

	var meanings []string
	for _, word := range words {
		meanings = append(meanings, word.Meaning)
	}

	return strings.Join(meanings, " ")
}

// calculateSentenceWeight calculates the overall weight of the sentence.
func (sg *SentenceGenerator) calculateSentenceWeight(words []morphology.Word) float32 {
	if len(words) == 0 {
		return 0.0
	}

	var totalWeight float32
	for _, word := range words {
		totalWeight += word.Weight
	}

	// Normalize by number of words
	return totalWeight / float32(len(words))
}

// generateSentenceID creates a unique identifier for the sentence.
func (sg *SentenceGenerator) generateSentenceID(words []morphology.Word, sentenceType SentenceType) string {
	wordCount := len(words)
	firstWord := ""
	if len(words) > 0 {
		firstWord = words[0].Meaning
	}

	return fmt.Sprintf("%s_%s_%s_%d",
		sg.grammar.GetCulture(),
		sentenceType.String(),
		firstWord,
		wordCount)
}

// GenerateRandomSentence creates a sentence with random word selection from a lexicon.
func (sg *SentenceGenerator) GenerateRandomSentence(
	lexicon *morphology.Lexicon,
	sentenceType SentenceType,
	rng *rand.Rand,
) (*Sentence, error) {
	// Get words of different categories
	subjects := lexicon.GetWordsByCategory(morphology.WordCategoryNoun)
	verbs := lexicon.GetWordsByCategory(morphology.WordCategoryVerb)
	objects := lexicon.GetWordsByCategory(morphology.WordCategoryNoun)

	if len(subjects) == 0 || len(verbs) == 0 {
		return nil, fmt.Errorf("insufficient words in lexicon for sentence generation")
	}

	// Select random words
	var words []morphology.Word
	words = append(words, *subjects[rng.IntN(len(subjects))])
	words = append(words, *verbs[rng.IntN(len(verbs))])

	// Add object if available
	if len(objects) > 0 {
		words = append(words, *objects[rng.IntN(len(objects))])
	}

	return sg.GenerateSentence(words, sentenceType, rng, nil)
}

// GenerateSimpleSentence creates a basic sentence with minimal complexity.
func (sg *SentenceGenerator) GenerateSimpleSentence(
	subject *morphology.Word,
	verb *morphology.Word,
	object *morphology.Word,
	sentenceType SentenceType,
	rng *rand.Rand,
) (*Sentence, error) {
	var words []morphology.Word

	if subject != nil {
		words = append(words, *subject)
	}
	if verb != nil {
		words = append(words, *verb)
	}
	if object != nil {
		words = append(words, *object)
	}

	return sg.GenerateSentence(words, sentenceType, rng, nil)
}
