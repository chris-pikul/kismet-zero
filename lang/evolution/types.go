package evolution

import (
	"time"

	"github.com/chris-pikul/kismet-zero/lang"
)

// ChangeType represents the category of linguistic change that occurred.
type ChangeType byte

const (
	ChangeTypeUnknown       ChangeType = iota
	ChangeTypeSoundShift               // Phonological changes (Grimm's Law, etc.)
	ChangeTypeMorphological            // Changes to word formation and grammar
	ChangeTypeLexical                  // New words, loanwords, semantic shifts
	ChangeTypeSyntactic                // Word order, sentence structure changes
	ChangeTypeOrthographic             // Writing system modifications
	ChangeTypeContact                  // Influence from other languages
)

var changeTypeEnum = []string{
	"unknown",
	"sound_shift",
	"morphological",
	"lexical",
	"syntactic",
	"orthographic",
	"contact",
}

// String returns the string representation of the ChangeType.
func (ct ChangeType) String() string {
	if ct > ChangeTypeContact {
		return changeTypeEnum[0]
	}
	return changeTypeEnum[ct]
}

// ChangeDirection represents whether a change is additive, subtractive, or modifying.
type ChangeDirection byte

const (
	ChangeDirectionUnknown     ChangeDirection = iota
	ChangeDirectionAdditive                    // New features added
	ChangeDirectionSubtractive                 // Features removed
	ChangeDirectionModifying                   // Existing features changed
	ChangeDirectionBlending                    // Features merged/combined
)

var changeDirectionEnum = []string{
	"unknown",
	"additive",
	"subtractive",
	"modifying",
	"blending",
}

// String returns the string representation of the ChangeDirection.
func (cd ChangeDirection) String() string {
	if cd > ChangeDirectionBlending {
		return changeDirectionEnum[0]
	}
	return changeDirectionEnum[cd]
}

// LinguisticChange represents a single change that occurred in a language's evolution.
type LinguisticChange struct {
	ID          string          `json:"id"`
	Type        ChangeType      `json:"type"`
	Direction   ChangeDirection `json:"direction"`
	Description string          `json:"description"`
	Details     string          `json:"details,omitempty"`

	// Affected components
	AffectedPhonemes  []string `json:"affectedPhonemes,omitempty"`
	AffectedMorphemes []string `json:"affectedMorphemes,omitempty"`
	AffectedRules     []string `json:"affectedRules,omitempty"`

	// Timing and context
	Timestamp time.Time `json:"timestamp"`
	Era       string    `json:"era,omitempty"`
	Trigger   string    `json:"trigger,omitempty"` // What caused this change

	// Cultural context
	CultureInfluence string  `json:"cultureInfluence,omitempty"`
	Intensity        float32 `json:"intensity"` // 0.0 to 1.0, strength of change
}

// EvolutionEvent represents a collection of related changes that occurred together.
type EvolutionEvent struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Timestamp   time.Time          `json:"timestamp"`
	Era         string             `json:"era,omitempty"`
	Changes     []LinguisticChange `json:"changes"`

	// Context and triggers
	TriggerType    string  `json:"triggerType"` // "natural", "contact", "cultural_shift"
	TriggerCulture string  `json:"triggerCulture,omitempty"`
	Intensity      float32 `json:"intensity"` // Overall impact of the event

	// Metadata
	Seed int64 `json:"seed"` // RNG seed for reproducibility
}

// LanguageNode represents a node in the language family tree.
type LanguageNode struct {
	Language  *lang.Language   `json:"language"`
	Parent    *LanguageNode    `json:"parent,omitempty"`
	Children  []*LanguageNode  `json:"children,omitempty"`
	Evolution []EvolutionEvent `json:"evolution"`

	// Relationship metadata
	DivergenceDate time.Time      `json:"divergenceDate,omitempty"`
	ContactHistory []ContactEvent `json:"contactHistory,omitempty"`
}

// ContactEvent represents an interaction between languages that may cause changes.
type ContactEvent struct {
	ID         string        `json:"id"`
	Timestamp  time.Time     `json:"timestamp"`
	SourceLang string        `json:"sourceLang"` // Language ID of influencing language
	TargetLang string        `json:"targetLang"` // Language ID of influenced language
	Type       string        `json:"type"`       // "trade", "conquest", "migration", "cultural"
	Intensity  float32       `json:"intensity"`  // 0.0 to 1.0, strength of influence
	Duration   time.Duration `json:"duration,omitempty"`

	// What was borrowed/influenced
	PhonologicalBorrowing bool `json:"phonologicalBorrowing"`
	LexicalBorrowing      bool `json:"lexicalBorrowing"`
	GrammaticalInfluence  bool `json:"grammaticalInfluence"`

	Description string `json:"description,omitempty"`
}

// EvolutionConfig holds configuration parameters for language evolution.
type EvolutionConfig struct {
	// Natural evolution parameters
	NaturalChangeRate     float32 `json:"naturalChangeRate"`     // Probability of natural changes per era
	SoundShiftProbability float32 `json:"soundShiftProbability"` // Likelihood of sound changes
	MorphologyChangeRate  float32 `json:"morphologyChangeRate"`  // Rate of morphological evolution

	// Contact evolution parameters
	ContactInfluenceRate float32 `json:"contactInfluenceRate"` // Base rate of contact influence
	BorrowingThreshold   float32 `json:"borrowingThreshold"`   // Minimum intensity for borrowing
	AdaptationStrength   float32 `json:"adaptationStrength"`   // How strongly borrowed features adapt

	// Cultural parameters
	CulturalInfluenceWeight float32       `json:"culturalInfluenceWeight"` // Weight of cultural factors
	EraDuration             time.Duration `json:"eraDuration"`             // Length of an evolution era

	// RNG and reproducibility
	Seed int64 `json:"seed"`
}

// DefaultEvolutionConfig returns a configuration with reasonable default values.
func DefaultEvolutionConfig(seed int64) EvolutionConfig {
	return EvolutionConfig{
		NaturalChangeRate:       0.3,
		SoundShiftProbability:   0.4,
		MorphologyChangeRate:    0.2,
		ContactInfluenceRate:    0.25,
		BorrowingThreshold:      0.3,
		AdaptationStrength:      0.6,
		CulturalInfluenceWeight: 0.4,
		EraDuration:             time.Hour * 24 * 365 * 100, // 100 years
		Seed:                    seed,
	}
}
