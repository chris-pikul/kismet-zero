package grammar

import (
	"math/rand/v2"
	"testing"

	"github.com/chris-pikul/kismet-zero/lang/morphology"
	"github.com/chris-pikul/kismet-zero/lang/phoneme"
	"github.com/chris-pikul/kismet-zero/lang/phonology"
)

func TestWordOrderEnum(t *testing.T) {
	tests := []struct {
		name      string
		wordOrder WordOrder
		expected  string
	}{
		{"SVO", WordOrderSVO, "svo"},
		{"SOV", WordOrderSOV, "sov"},
		{"VSO", WordOrderVSO, "vso"},
		{"VOS", WordOrderVOS, "vos"},
		{"OVS", WordOrderOVS, "ovs"},
		{"OSV", WordOrderOSV, "osv"},
		{"Invalid", WordOrder(99), "svo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.wordOrder.String()
			if result != tt.expected {
				t.Errorf("WordOrder.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCaseEnum(t *testing.T) {
	tests := []struct {
		name     string
		caseVal  Case
		expected string
	}{
		{"Nominative", CaseNominative, "nominative"},
		{"Accusative", CaseAccusative, "accusative"},
		{"Genitive", CaseGenitive, "genitive"},
		{"Dative", CaseDative, "dative"},
		{"Locative", CaseLocative, "locative"},
		{"Instrumental", CaseInstrumental, "instrumental"},
		{"Vocative", CaseVocative, "vocative"},
		{"Ablative", CaseAblative, "ablative"},
		{"Invalid", Case(99), "nominative"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.caseVal.String()
			if result != tt.expected {
				t.Errorf("Case.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNumberEnum(t *testing.T) {
	tests := []struct {
		name     string
		number   Number
		expected string
	}{
		{"Singular", NumberSingular, "singular"},
		{"Dual", NumberDual, "dual"},
		{"Plural", NumberPlural, "plural"},
		{"Invalid", Number(99), "singular"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.number.String()
			if result != tt.expected {
				t.Errorf("Number.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGenderEnum(t *testing.T) {
	tests := []struct {
		name     string
		gender   Gender
		expected string
	}{
		{"Masculine", GenderMasculine, "masculine"},
		{"Feminine", GenderFeminine, "feminine"},
		{"Neuter", GenderNeuter, "neuter"},
		{"Animate", GenderAnimate, "animate"},
		{"Inanimate", GenderInanimate, "inanimate"},
		{"Invalid", Gender(99), "masculine"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.gender.String()
			if result != tt.expected {
				t.Errorf("Gender.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTenseEnum(t *testing.T) {
	tests := []struct {
		name     string
		tense    Tense
		expected string
	}{
		{"Present", TensePresent, "present"},
		{"Past", TensePast, "past"},
		{"Future", TenseFuture, "future"},
		{"Invalid", Tense(99), "present"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tense.String()
			if result != tt.expected {
				t.Errorf("Tense.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAspectEnum(t *testing.T) {
	tests := []struct {
		name     string
		aspect   Aspect
		expected string
	}{
		{"Perfective", AspectPerfective, "perfective"},
		{"Imperfective", AspectImperfective, "imperfective"},
		{"Habitual", AspectHabitual, "habitual"},
		{"Progressive", AspectProgressive, "progressive"},
		{"Invalid", Aspect(99), "perfective"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.aspect.String()
			if result != tt.expected {
				t.Errorf("Aspect.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMoodEnum(t *testing.T) {
	tests := []struct {
		name     string
		mood     Mood
		expected string
	}{
		{"Indicative", MoodIndicative, "indicative"},
		{"Subjunctive", MoodSubjunctive, "subjunctive"},
		{"Imperative", MoodImperative, "imperative"},
		{"Conditional", MoodConditional, "conditional"},
		{"Invalid", Mood(99), "indicative"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.mood.String()
			if result != tt.expected {
				t.Errorf("Mood.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSentenceTypeEnum(t *testing.T) {
	tests := []struct {
		name     string
		sentType SentenceType
		expected string
	}{
		{"Declarative", SentenceTypeDeclarative, "declarative"},
		{"Interrogative", SentenceTypeInterrogative, "interrogative"},
		{"Imperative", SentenceTypeImperative, "imperative"},
		{"Exclamatory", SentenceTypeExclamatory, "exclamatory"},
		{"Invalid", SentenceType(99), "declarative"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sentType.String()
			if result != tt.expected {
				t.Errorf("SentenceType.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewAgreementSystem(t *testing.T) {
	cases := []Case{CaseNominative, CaseAccusative}
	numbers := []Number{NumberSingular, NumberPlural}
	genders := []Gender{GenderMasculine, GenderFeminine}
	tense := []Tense{TensePresent, TensePast}
	aspects := []Aspect{AspectPerfective, AspectImperfective}
	moods := []Mood{MoodIndicative, MoodImperative}

	as := NewAgreementSystem(cases, numbers, genders, tense, aspects, moods)

	if len(as.cases) != 2 {
		t.Errorf("Expected 2 cases, got %d", len(as.cases))
	}
	if len(as.numbers) != 2 {
		t.Errorf("Expected 2 numbers, got %d", len(as.numbers))
	}
	if len(as.genders) != 2 {
		t.Errorf("Expected 2 genders, got %d", len(as.genders))
	}
	if len(as.tense) != 2 {
		t.Errorf("Expected 2 tenses, got %d", len(as.tense))
	}
	if len(as.aspects) != 2 {
		t.Errorf("Expected 2 aspects, got %d", len(as.aspects))
	}
	if len(as.moods) != 2 {
		t.Errorf("Expected 2 moods, got %d", len(as.moods))
	}
}

func TestNewGrammar(t *testing.T) {
	// Create a mock phoneme pool
	pool := phoneme.NewPool()

	// Create a mock phonology
	phon := phonology.NewPhonology(pool)

	// Create a mock word builder
	wordBuilder := &morphology.WordBuilder{}

	// Create agreement system
	agreement := NewAgreementSystem(
		[]Case{CaseNominative, CaseAccusative},
		[]Number{NumberSingular, NumberPlural},
		[]Gender{GenderMasculine, GenderFeminine},
		[]Tense{TensePresent, TensePast},
		[]Aspect{AspectPerfective, AspectImperfective},
		[]Mood{MoodIndicative, MoodImperative},
	)

	grammar := NewGrammar(WordOrderSVO, agreement, phon, wordBuilder, "test_culture")

	if grammar.GetWordOrder() != WordOrderSVO {
		t.Errorf("Expected WordOrderSVO, got %v", grammar.GetWordOrder())
	}
	if grammar.GetCulture() != "test_culture" {
		t.Errorf("Expected 'test_culture', got %s", grammar.GetCulture())
	}
}

func TestWordOrderRule(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 123))

	// Create test words
	subject := morphology.Word{
		ID:       "subject",
		Category: morphology.WordCategoryNoun,
		Meaning:  "cat",
		Agreement: morphology.AgreementFeatures{
			Case:   "nominative",
			Number: "singular",
		},
	}

	verb := morphology.Word{
		ID:       "verb",
		Category: morphology.WordCategoryVerb,
		Meaning:  "runs",
		Agreement: morphology.AgreementFeatures{
			Tense:  "present",
			Number: "singular",
		},
	}

	object := morphology.Word{
		ID:       "object",
		Category: morphology.WordCategoryNoun,
		Meaning:  "fast",
		Agreement: morphology.AgreementFeatures{
			Case:   "accusative",
			Number: "singular",
		},
	}

	words := []morphology.Word{subject, verb, object}

	// Test SVO word order
	svoRule := NewWordOrderRule(WordOrderSVO, 10.0)
	result, err := svoRule.Apply(words, rng)
	if err != nil {
		t.Fatalf("Failed to apply SVO rule: %v", err)
	}

	// Check that the order is maintained (SVO)
	if result[0].Category != morphology.WordCategoryNoun ||
		result[1].Category != morphology.WordCategoryVerb ||
		result[2].Category != morphology.WordCategoryNoun {
		t.Errorf("Expected SVO order, got: %s-%s-%s",
			result[0].Category.String(),
			result[1].Category.String(),
			result[2].Category.String())
	}

	// Test SOV word order
	sovRule := NewWordOrderRule(WordOrderSOV, 10.0)
	result, err = sovRule.Apply(words, rng)
	if err != nil {
		t.Fatalf("Failed to apply SOV rule: %v", err)
	}

	// Check that the order is SOV
	if result[0].Category != morphology.WordCategoryNoun ||
		result[1].Category != morphology.WordCategoryNoun ||
		result[2].Category != morphology.WordCategoryVerb {
		t.Errorf("Expected SOV order, got: %s-%s-%s",
			result[0].Category.String(),
			result[1].Category.String(),
			result[2].Category.String())
	}
}

func TestSentenceTemplateRule(t *testing.T) {
	// Create test words
	subject := morphology.Word{
		ID:       "subject",
		Category: morphology.WordCategoryNoun,
		Meaning:  "cat",
	}

	verb := morphology.Word{
		ID:       "verb",
		Category: morphology.WordCategoryVerb,
		Meaning:  "runs",
	}

	words := []morphology.Word{subject, verb}

	// Test SVO template rule
	template := TemplateSVO
	rule := NewSentenceTemplateRule(template, 10.0)

	if !rule.Validate(words) {
		t.Error("SVO template validation failed for valid words")
	}

	// Test with invalid words (missing verb)
	invalidWords := []morphology.Word{subject}
	if rule.Validate(invalidWords) {
		t.Error("SVO template validation passed for invalid words")
	}
}

func TestGetDefaultTemplates(t *testing.T) {
	templates := GetDefaultTemplates()

	if len(templates) == 0 {
		t.Error("Expected non-empty default templates")
	}

	// Check that we have templates for different sentence types
	hasDeclarative := false
	hasInterrogative := false
	hasImperative := false

	for _, template := range templates {
		switch template.Type {
		case SentenceTypeDeclarative:
			hasDeclarative = true
		case SentenceTypeInterrogative:
			hasInterrogative = true
		case SentenceTypeImperative:
			hasImperative = true
		}
	}

	if !hasDeclarative {
		t.Error("Missing declarative templates")
	}
	if !hasInterrogative {
		t.Error("Missing interrogative templates")
	}
	if !hasImperative {
		t.Error("Missing imperative templates")
	}
}

func TestGetDefaultAgreementRules(t *testing.T) {
	rules := GetDefaultAgreementRules()

	if len(rules) == 0 {
		t.Error("Expected non-empty default agreement rules")
	}

	// Check that we have the expected number of rules
	expectedRuleCount := 5
	if len(rules) != expectedRuleCount {
		t.Errorf("Expected %d agreement rules, got %d", expectedRuleCount, len(rules))
	}
}

func TestSentenceGenerator(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 123))

	// Create a mock phoneme pool
	pool := phoneme.NewPool()

	// Create a mock phonology
	phon := phonology.NewPhonology(pool)

	// Create a mock word builder
	wordBuilder := &morphology.WordBuilder{}

	// Create agreement system
	agreement := NewAgreementSystem(
		[]Case{CaseNominative, CaseAccusative},
		[]Number{NumberSingular, NumberPlural},
		[]Gender{GenderMasculine, GenderFeminine},
		[]Tense{TensePresent, TensePast},
		[]Aspect{AspectPerfective, AspectImperfective},
		[]Mood{MoodIndicative, MoodImperative},
	)

	// Create grammar
	grammar := NewGrammar(WordOrderSVO, agreement, phon, wordBuilder, "test_culture")

	// Create sentence generator
	generator := NewSentenceGenerator(grammar)

	// Create test words
	subject := morphology.Word{
		ID:       "subject",
		Category: morphology.WordCategoryNoun,
		Meaning:  "cat",
		Agreement: morphology.AgreementFeatures{
			Case:   "nominative",
			Number: "singular",
		},
	}

	verb := morphology.Word{
		ID:       "verb",
		Category: morphology.WordCategoryVerb,
		Meaning:  "runs",
		Agreement: morphology.AgreementFeatures{
			Tense:  "present",
			Number: "singular",
		},
	}

	words := []morphology.Word{subject, verb}

	// Generate a declarative sentence
	sentence, err := generator.GenerateSentence(words, SentenceTypeDeclarative, rng)
	if err != nil {
		t.Fatalf("Failed to generate sentence: %v", err)
	}

	if sentence.Type != SentenceTypeDeclarative {
		t.Errorf("Expected declarative sentence, got %s", sentence.Type.String())
	}

	if len(sentence.Words) != 2 {
		t.Errorf("Expected 2 words, got %d", len(sentence.Words))
	}

	if sentence.Culture != "test_culture" {
		t.Errorf("Expected culture 'test_culture', got %s", sentence.Culture)
	}

	if sentence.Structure == "" {
		t.Error("Expected non-empty sentence structure")
	}

	if sentence.Meaning == "" {
		t.Error("Expected non-empty sentence meaning")
	}
}

func TestGeneratePhrase(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 123))

	// Create a mock phoneme pool
	pool := phoneme.NewPool()

	// Create a mock phonology
	phon := phonology.NewPhonology(pool)

	// Create a mock word builder
	wordBuilder := &morphology.WordBuilder{}

	// Create agreement system
	agreement := NewAgreementSystem(
		[]Case{CaseNominative, CaseAccusative},
		[]Number{NumberSingular, NumberPlural},
		[]Gender{GenderMasculine, GenderFeminine},
		[]Tense{TensePresent, TensePast},
		[]Aspect{AspectPerfective, AspectImperfective},
		[]Mood{MoodIndicative, MoodImperative},
	)

	// Create grammar
	grammar := NewGrammar(WordOrderSVO, agreement, phon, wordBuilder, "test_culture")

	// Create sentence generator
	generator := NewSentenceGenerator(grammar)

	// Create test words for a noun phrase
	noun := morphology.Word{
		ID:       "noun",
		Category: morphology.WordCategoryNoun,
		Meaning:  "cat",
		Agreement: morphology.AgreementFeatures{
			Case:   "nominative",
			Number: "singular",
			Gender: "masculine",
		},
	}

	adjective := morphology.Word{
		ID:       "adjective",
		Category: morphology.WordCategoryAdjective,
		Meaning:  "black",
		Agreement: morphology.AgreementFeatures{
			Number: "singular",
			Gender: "masculine",
		},
	}

	words := []morphology.Word{noun, adjective}

	// Generate a noun phrase
	phrase, err := generator.GeneratePhrase(words, "noun", rng)
	if err != nil {
		t.Fatalf("Failed to generate phrase: %v", err)
	}

	if phrase.Type != "noun" {
		t.Errorf("Expected phrase type 'noun', got %s", phrase.Type)
	}

	if len(phrase.Words) != 2 {
		t.Errorf("Expected 2 words, got %d", len(phrase.Words))
	}

	if phrase.Head == nil {
		t.Error("Expected non-nil phrase head")
	}

	if phrase.Head.Category != morphology.WordCategoryNoun {
		t.Errorf("Expected noun as phrase head, got %s", phrase.Head.Category.String())
	}
}
