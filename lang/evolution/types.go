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
	ChangeTypeDialectal                // Dialect-specific changes
)

var changeTypeEnum = []string{
	"unknown",
	"sound_shift",
	"morphological",
	"lexical",
	"syntactic",
	"orthographic",
	"contact",
	"dialectal",
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
	OrthographicBorrowing bool `json:"orthographicBorrowing"`

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

	// Orthographic evolution parameters
	OrthographyChangeRate float32 `json:"orthographyChangeRate"` // Rate of orthographic evolution
	ScriptReformRate      float32 `json:"scriptReformRate"`      // Rate of script reforms

	// Cultural parameters
	CulturalInfluenceWeight float32       `json:"culturalInfluenceWeight"` // Weight of cultural factors
	EraDuration             time.Duration `json:"eraDuration"`             // Length of an evolution era

	// Dialect formation parameters
	DialectFormationRate       float32 `json:"dialectFormationRate"`       // Probability of dialect formation per era
	GeographicIsolationWeight  float32 `json:"geographicIsolationWeight"`  // Weight of geographic factors in dialect formation
	SocialStratificationWeight float32 `json:"socialStratificationWeight"` // Weight of social factors in dialect formation
	UrbanRuralDivergenceRate   float32 `json:"urbanRuralDivergenceRate"`   // Rate of urban-rural dialect divergence

	// RNG and reproducibility
	Seed int64 `json:"seed"`
}

// DefaultEvolutionConfig returns a configuration with reasonable default values.
func DefaultEvolutionConfig(seed int64) EvolutionConfig {
	return EvolutionConfig{
		NaturalChangeRate:          0.3,
		SoundShiftProbability:      0.4,
		MorphologyChangeRate:       0.2,
		ContactInfluenceRate:       0.25,
		BorrowingThreshold:         0.3,
		AdaptationStrength:         0.6,
		OrthographyChangeRate:      0.15,
		ScriptReformRate:           0.1,
		CulturalInfluenceWeight:    0.4,
		EraDuration:                time.Hour * 24 * 365 * 100, // 100 years
		DialectFormationRate:       0.15,
		GeographicIsolationWeight:  0.6,
		SocialStratificationWeight: 0.4,
		UrbanRuralDivergenceRate:   0.25,
		Seed:                       seed,
	}
}

// DialectType represents the type of dialect formation.
type DialectType byte

const (
	DialectTypeUnknown    DialectType = iota
	DialectTypeGeographic             // Regional/geographic variation
	DialectTypeSocial                 // Social class variation
	DialectTypeUrban                  // Urban vs rural variation
	DialectTypeTemporal               // Historical/archaic variation
	DialectTypeContact                // Contact-induced variation
)

var dialectTypeEnum = []string{
	"unknown",
	"geographic",
	"social",
	"urban",
	"temporal",
	"contact",
}

// String returns the string representation of the DialectType.
func (dt DialectType) String() string {
	if dt > DialectTypeContact {
		return dialectTypeEnum[0]
	}
	return dialectTypeEnum[dt]
}

// GeographicRegion represents a geographical area where a dialect is spoken.
type GeographicRegion struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Climate      string  `json:"climate,omitempty"` // "tropical", "temperate", "arctic", etc.
	Terrain      string  `json:"terrain,omitempty"` // "mountain", "coastal", "plains", etc.
	Population   int     `json:"population,omitempty"`
	Urbanization float32 `json:"urbanization,omitempty"` // 0.0 to 1.0, rural to urban
}

// DialectFeatures represents the distinctive features of a dialect.
type DialectFeatures struct {
	ID                   string   `json:"id"`
	PhonologicalFeatures []string `json:"phonologicalFeatures,omitempty"` // Distinctive sound patterns
	LexicalFeatures      []string `json:"lexicalFeatures,omitempty"`      // Local vocabulary
	GrammaticalFeatures  []string `json:"grammaticalFeatures,omitempty"`  // Grammar variations
	PragmaticFeatures    []string `json:"pragmaticFeatures,omitempty"`    // Usage patterns
	IntelligibilityScore float32  `json:"intelligibilityScore"`           // 0.0 to 1.0, mutual intelligibility with parent
}

// Dialect represents a regional or social variant of a language.
type Dialect struct {
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	Type          DialectType      `json:"type"`
	ParentLang    string           `json:"parentLang"` // ID of parent language
	Region        GeographicRegion `json:"region"`
	Features      DialectFeatures  `json:"features"`
	FormationDate time.Time        `json:"formationDate"`
	Era           string           `json:"era"`

	// Evolution tracking
	EvolutionHistory []EvolutionEvent `json:"evolutionHistory,omitempty"`
	ContactHistory   []ContactEvent   `json:"contactHistory,omitempty"`

	// Metadata
	Description string `json:"description,omitempty"`
	Status      string `json:"status"` // "active", "archaic", "extinct"
	Seed        int64  `json:"seed"`
}

// DialectFormationEvent represents the creation of a new dialect.
type DialectFormationEvent struct {
	ID         string    `json:"id"`
	Timestamp  time.Time `json:"timestamp"`
	Era        string    `json:"era"`
	ParentLang string    `json:"parentLang"`
	NewDialect string    `json:"newDialect"`

	// Formation factors
	GeographicFactors []string `json:"geographicFactors,omitempty"`
	SocialFactors     []string `json:"socialFactors,omitempty"`
	ContactFactors    []string `json:"contactFactors,omitempty"`

	// Intensity and description
	Intensity   float32 `json:"intensity"`
	Description string  `json:"description"`
	Seed        int64   `json:"seed"`
}
