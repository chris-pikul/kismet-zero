package grammar

import (
	"math/rand/v2"

	"github.com/chris-pikul/kismet-zero/lang/morphology"
	"github.com/chris-pikul/kismet-zero/lang/phonology"
)

// WordOrder defines the basic sentence structure pattern for a language.
type WordOrder byte

const (
	WordOrderSVO WordOrder = iota // Subject-Verb-Object (English, Chinese)
	WordOrderSOV                  // Subject-Object-Verb (Japanese, Turkish)
	WordOrderVSO                  // Verb-Subject-Object (Arabic, Welsh)
	WordOrderVOS                  // Verb-Object-Subject (Malagasy)
	WordOrderOVS                  // Object-Verb-Subject (Hixkaryana)
	WordOrderOSV                  // Object-Subject-Verb (Apurinã)
)

var wordOrderEnum = []string{
	"svo",
	"sov",
	"vso",
	"vos",
	"ovs",
	"osv",
}

// String returns the string representation of the WordOrder.
func (wo WordOrder) String() string {
	if wo > WordOrderOSV {
		return wordOrderEnum[0]
	}
	return wordOrderEnum[wo]
}

// Case represents grammatical case marking for nouns and pronouns.
type Case byte

const (
	CaseNominative   Case = iota // Subject of the sentence
	CaseAccusative               // Direct object
	CaseGenitive                 // Possession, origin
	CaseDative                   // Indirect object, recipient
	CaseLocative                 // Location, place
	CaseInstrumental             // Means, instrument
	CaseVocative                 // Direct address
	CaseAblative                 // Source, separation
)

var caseEnum = []string{
	"nominative",
	"accusative",
	"genitive",
	"dative",
	"locative",
	"instrumental",
	"vocative",
	"ablative",
}

// String returns the string representation of the Case.
func (c Case) String() string {
	if c > CaseAblative {
		return caseEnum[0]
	}
	return caseEnum[c]
}

// Number represents grammatical number agreement.
type Number byte

const (
	NumberSingular Number = iota // One item
	NumberDual                   // Two items
	NumberPlural                 // Three or more items
)

var numberEnum = []string{
	"singular",
	"dual",
	"plural",
}

// String returns the string representation of the Number.
func (n Number) String() string {
	if n > NumberPlural {
		return numberEnum[0]
	}
	return numberEnum[n]
}

// Gender represents grammatical gender agreement.
type Gender byte

const (
	GenderMasculine Gender = iota // Masculine gender
	GenderFeminine                // Feminine gender
	GenderNeuter                  // Neuter gender
	GenderAnimate                 // Animate beings
	GenderInanimate               // Inanimate objects
)

var genderEnum = []string{
	"masculine",
	"feminine",
	"neuter",
	"animate",
	"inanimate",
}

// String returns the string representation of the Gender.
func (g Gender) String() string {
	if g > GenderInanimate {
		return genderEnum[0]
	}
	return genderEnum[g]
}

// Tense represents grammatical tense marking.
type Tense byte

const (
	TensePresent Tense = iota // Current time
	TensePast                 // Past time
	TenseFuture               // Future time
)

var tenseEnum = []string{
	"present",
	"past",
	"future",
}

// String returns the string representation of the Tense.
func (t Tense) String() string {
	if t > TenseFuture {
		return tenseEnum[0]
	}
	return tenseEnum[t]
}

// Aspect represents grammatical aspect marking.
type Aspect byte

const (
	AspectPerfective   Aspect = iota // Completed action
	AspectImperfective               // Ongoing action
	AspectHabitual                   // Habitual action
	AspectProgressive                // Progressive action
)

var aspectEnum = []string{
	"perfective",
	"imperfective",
	"habitual",
	"progressive",
}

// String returns the string representation of the Aspect.
func (a Aspect) String() string {
	if a > AspectProgressive {
		return aspectEnum[0]
	}
	return aspectEnum[a]
}

// Mood represents grammatical mood marking.
type Mood byte

const (
	MoodIndicative  Mood = iota // Statement of fact
	MoodSubjunctive             // Subjective, hypothetical
	MoodImperative              // Command, request
	MoodConditional             // Conditional statement
)

var moodEnum = []string{
	"indicative",
	"subjunctive",
	"imperative",
	"conditional",
}

// String returns the string representation of the Mood.
func (m Mood) String() string {
	if m > MoodConditional {
		return moodEnum[0]
	}
	return moodEnum[m]
}

// SentenceType represents the type of sentence being generated.
type SentenceType byte

const (
	SentenceTypeDeclarative   SentenceType = iota // Statement
	SentenceTypeInterrogative                     // Question
	SentenceTypeImperative                        // Command
	SentenceTypeExclamatory                       // Exclamation
)

var sentenceTypeEnum = []string{
	"declarative",
	"interrogative",
	"imperative",
	"exclamatory",
}

// String returns the string representation of the SentenceType.
func (st SentenceType) String() string {
	if st > SentenceTypeExclamatory {
		return sentenceTypeEnum[0]
	}
	return sentenceTypeEnum[st]
}

// AgreementFeature represents a grammatical feature that requires agreement.
type AgreementFeature struct {
	Case   Case   `json:"case"`
	Number Number `json:"number"`
	Gender Gender `json:"gender"`
}

// VerbFeatures represents the grammatical features of a verb.
type VerbFeatures struct {
	Tense  Tense  `json:"tense"`
	Aspect Aspect `json:"aspect"`
	Mood   Mood   `json:"mood"`
	Number Number `json:"number"`
	Person byte   `json:"person"` // 1st, 2nd, 3rd person
}

// AgreementRule defines how grammatical features should agree across sentence elements.
type AgreementRule interface {
	// Apply applies agreement rules to the given words.
	Apply(words []morphology.Word, rng *rand.Rand) ([]morphology.Word, error)

	// Validate checks if the given words follow this agreement rule.
	Validate(words []morphology.Word) bool

	// Weight returns the weight of this rule for selection.
	Weight() float32
}

// SyntaxRule defines how to construct valid sentence patterns.
type SyntaxRule interface {
	// Apply applies the syntax rule to the given words.
	Apply(words []morphology.Word, rng *rand.Rand) ([]morphology.Word, error)

	// Validate checks if the given words follow this syntax rule.
	Validate(words []morphology.Word) bool

	// Weight returns the weight of this rule for selection.
	Weight() float32
}

// Sentence represents a complete grammatical sentence.
type Sentence struct {
	ID        string            `json:"id"`
	Type      SentenceType      `json:"type"`
	Words     []morphology.Word `json:"words"`
	Structure string            `json:"structure"` // String representation of structure
	Meaning   string            `json:"meaning"`
	Culture   string            `json:"culture"`
	Weight    float32           `json:"weight"`
}

// Phrase represents a grammatical phrase within a sentence.
type Phrase struct {
	Type      string            `json:"type"` // "noun", "verb", "prepositional"
	Words     []morphology.Word `json:"words"`
	Head      *morphology.Word  `json:"head"`      // Main word of the phrase
	Modifiers []morphology.Word `json:"modifiers"` // Adjectives, adverbs, etc.
	Agreement AgreementFeature  `json:"agreement"`
}

// Grammar represents the syntactic rules and patterns for a language.
type Grammar struct {
	wordOrder   WordOrder
	agreement   *AgreementSystem
	syntaxRules []SyntaxRule
	culture     string
	phonology   *phonology.Phonology
	wordBuilder *morphology.WordBuilder
}

// AgreementSystem handles grammatical agreement rules for a language.
type AgreementSystem struct {
	cases   []Case
	numbers []Number
	genders []Gender
	tense   []Tense
	aspects []Aspect
	moods   []Mood
	rules   []AgreementRule
}

// NewAgreementSystem creates a new agreement system with the specified features.
func NewAgreementSystem(cases []Case, numbers []Number, genders []Gender,
	tense []Tense, aspects []Aspect, moods []Mood) *AgreementSystem {
	return &AgreementSystem{
		cases:   cases,
		numbers: numbers,
		genders: genders,
		tense:   tense,
		aspects: aspects,
		moods:   moods,
		rules:   make([]AgreementRule, 0),
	}
}

// AddRule adds an agreement rule to the system.
func (as *AgreementSystem) AddRule(rule AgreementRule) {
	as.rules = append(as.rules, rule)
}

// GetRules returns all agreement rules.
func (as *AgreementSystem) GetRules() []AgreementRule {
	return as.rules
}

// NewGrammar creates a new grammar with the given configuration.
func NewGrammar(wordOrder WordOrder, agreement *AgreementSystem,
	phonology *phonology.Phonology, wordBuilder *morphology.WordBuilder,
	culture string) *Grammar {
	return &Grammar{
		wordOrder:   wordOrder,
		agreement:   agreement,
		syntaxRules: make([]SyntaxRule, 0),
		culture:     culture,
		phonology:   phonology,
		wordBuilder: wordBuilder,
	}
}

// AddSyntaxRule adds a syntax rule to the grammar.
func (g *Grammar) AddSyntaxRule(rule SyntaxRule) {
	g.syntaxRules = append(g.syntaxRules, rule)
}

// GetSyntaxRules returns all syntax rules.
func (g *Grammar) GetSyntaxRules() []SyntaxRule {
	return g.syntaxRules
}

// GetWordOrder returns the word order of this grammar.
func (g *Grammar) GetWordOrder() WordOrder {
	return g.wordOrder
}

// GetCulture returns the culture identifier for this grammar.
func (g *Grammar) GetCulture() string {
	return g.culture
}
