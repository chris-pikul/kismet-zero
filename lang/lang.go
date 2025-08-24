package lang

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/chris-pikul/kismet-zero/lang/grammar"
	"github.com/chris-pikul/kismet-zero/lang/morphology"
	"github.com/chris-pikul/kismet-zero/lang/orthography"
	"github.com/chris-pikul/kismet-zero/lang/phonology"
)

// LanguageID represents a unique identifier for a language, inspired by BCP-47.
// Format: [family]-[branch]-[language]-[dialect]-[script]-[region]
// Examples: "elv-wood-sindarin", "hum-germanic-old-english-latin", "orc-black-speech"
type LanguageID struct {
	Family   string `json:"family"`   // Language family (elv, hum, orc, dwf, etc.)
	Branch   string `json:"branch"`   // Sub-family branch (wood, high, black, etc.)
	Language string `json:"language"` // Specific language name (sindarin, quenya, etc.)
	Dialect  string `json:"dialect"`  // Regional or social dialect variant
	Script   string `json:"script"`   // Writing system identifier
	Region   string `json:"region"`   // Geographic or cultural region
}

// String returns the canonical string representation of the language ID.
func (id LanguageID) String() string {
	parts := make([]string, 0, 6)

	if id.Family != "" {
		parts = append(parts, id.Family)
	}
	if id.Branch != "" {
		parts = append(parts, id.Branch)
	}
	if id.Language != "" {
		parts = append(parts, id.Language)
	}
	if id.Dialect != "" {
		parts = append(parts, id.Dialect)
	}
	if id.Script != "" {
		parts = append(parts, id.Script)
	}
	if id.Region != "" {
		parts = append(parts, id.Region)
	}

	return strings.Join(parts, "-")
}

// ParseLanguageID parses a string into a LanguageID struct.
func ParseLanguageID(idStr string) (LanguageID, error) {
	if strings.TrimSpace(idStr) == "" {
		return LanguageID{}, fmt.Errorf("language ID cannot be empty")
	}

	parts := strings.Split(idStr, "-")
	if len(parts) < 1 || len(parts) > 6 {
		return LanguageID{}, fmt.Errorf("invalid language ID format: %s (got %d parts)", idStr, len(parts))
	}

	id := LanguageID{}

	switch len(parts) {
	case 6:
		id.Region = parts[5]
		fallthrough
	case 5:
		id.Script = parts[4]
		fallthrough
	case 4:
		id.Dialect = parts[3]
		fallthrough
	case 3:
		id.Language = parts[2]
		fallthrough
	case 2:
		id.Branch = parts[1]
		fallthrough
	case 1:
		id.Family = parts[0]
	}

	return id, nil
}

// LanguageType represents the classification of a language based on its
// linguistic characteristics and cultural context.
type LanguageType byte

const (
	LanguageTypeUnknown     LanguageType = iota
	LanguageTypeNatural                  // Naturally evolved language
	LanguageTypeConstructed              // Artificially constructed language
	LanguageTypeDivine                   // Language of divine or supernatural origin
	LanguageTypeAncient                  // Ancient or extinct language
	LanguageTypeModern                   // Contemporary or actively used language
)

var languageTypeEnum = []string{
	"unknown",
	"natural",
	"constructed",
	"divine",
	"ancient",
	"modern",
}

// String returns the string representation of the LanguageType.
func (lt LanguageType) String() string {
	if lt > LanguageTypeModern {
		return languageTypeEnum[0]
	}
	return languageTypeEnum[lt]
}

// LanguageComplexity represents the overall complexity level of a language.
type LanguageComplexity byte

const (
	LanguageComplexityUnknown     LanguageComplexity = iota
	LanguageComplexitySimple                         // Simple grammar, small phoneme inventory
	LanguageComplexityModerate                       // Moderate complexity, balanced features
	LanguageComplexityComplex                        // Complex grammar, rich morphology
	LanguageComplexityVeryComplex                    // Highly complex, many grammatical categories
)

var languageComplexityEnum = []string{
	"unknown",
	"simple",
	"moderate",
	"complex",
	"very_complex",
}

// String returns the string representation of the LanguageComplexity.
func (lc LanguageComplexity) String() string {
	if lc > LanguageComplexityVeryComplex {
		return languageComplexityEnum[0]
	}
	return languageComplexityEnum[lc]
}

// Language represents a complete linguistic system with all its components.
// It serves as the central orchestrator for all language-related operations.
type Language struct {
	ID          LanguageID         `json:"id"`
	Name        string             `json:"name"`
	Type        LanguageType       `json:"type"`
	Complexity  LanguageComplexity `json:"complexity"`
	Culture     string             `json:"culture,omitempty"`
	Description string             `json:"description,omitempty"`

	// Core linguistic components
	Phonology   *phonology.Phonology       `json:"phonology"`
	Orthography *orthography.WritingSystem `json:"orthography"`
	Morphology  *morphology.MorphemeList   `json:"morphology"`
	Grammar     *grammar.Grammar           `json:"grammar"`

	// Metadata and evolution tracking
	CreatedAt time.Time    `json:"createdAt"`
	EvolvedAt time.Time    `json:"evolvedAt,omitempty"`
	ParentID  *LanguageID  `json:"parentId,omitempty"`
	ChildIDs  []LanguageID `json:"childIds,omitempty"`

	// Linguistic change tracking
	LinguisticChanges []LinguisticChange `json:"linguisticChanges,omitempty"`

	// Configuration and generation parameters
	Seed int64      `json:"seed"`
	RNG  *rand.Rand `json:"-"`

	// Optional interlingua services for semantic analysis and generation
	interlinguaServices *InterlinguaServices `json:"-"`
}

// NewLanguage creates a new language with the specified configuration.
// All linguistic components are initialized with default values that can
// be customized after creation.
func NewLanguage(id LanguageID, name string, langType LanguageType, seed int64) *Language {
	rng := rand.New(rand.NewPCG(uint64(seed), 0))

	return &Language{
		ID:         id,
		Name:       name,
		Type:       langType,
		Complexity: LanguageComplexityModerate,
		CreatedAt:  time.Now(),
		Seed:       seed,
		RNG:        rng,
	}
}

// SetPhonology assigns a phonology to the language.
func (l *Language) SetPhonology(ph *phonology.Phonology) {
	l.Phonology = ph
}

// SetOrthography assigns a writing system to the language.
func (l *Language) SetOrthography(orth *orthography.WritingSystem) {
	l.Orthography = orth
}

// SetMorphology assigns a morpheme list to the language.
func (l *Language) SetMorphology(morph *morphology.MorphemeList) {
	l.Morphology = morph
}

// SetGrammar assigns grammar rules to the language.
func (l *Language) SetGrammar(gram *grammar.Grammar) {
	l.Grammar = gram
}

// SetComplexity updates the language complexity level.
func (l *Language) SetComplexity(complexity LanguageComplexity) {
	l.Complexity = complexity
}

// SetCulture links the language to a specific culture.
func (l *Language) SetCulture(culture string) {
	l.Culture = culture
}

// SetDescription provides a description of the language.
func (l *Language) SetDescription(description string) {
	l.Description = description
}

// IsComplete checks if the language has all required components.
func (l *Language) IsComplete() bool {
	return l.Phonology != nil && l.Orthography != nil &&
		l.Morphology != nil && l.Grammar != nil
}

// GetMissingComponents returns a list of missing linguistic components.
func (l *Language) GetMissingComponents() []string {
	missing := make([]string, 0)

	if l.Phonology == nil {
		missing = append(missing, "phonology")
	}
	if l.Orthography == nil {
		missing = append(missing, "orthography")
	}
	if l.Morphology == nil {
		missing = append(missing, "morphology")
	}
	if l.Grammar == nil {
		missing = append(missing, "grammar")
	}

	return missing
}

// Clone creates a deep copy of the language with a new ID and seed.
func (l *Language) Clone(newID LanguageID, newSeed int64) *Language {
	clone := NewLanguage(newID, l.Name+" (clone)", l.Type, newSeed)
	clone.Complexity = l.Complexity
	clone.Culture = l.Culture
	clone.Description = l.Description
	clone.ParentID = &l.ID
	clone.CreatedAt = time.Now()

	// Copy linguistic components
	if l.Phonology != nil {
		clone.Phonology = l.Phonology
	}
	if l.Orthography != nil {
		clone.Orthography = l.Orthography
	}
	if l.Morphology != nil {
		clone.Morphology = l.Morphology
	}
	if l.Grammar != nil {
		clone.Grammar = l.Grammar
	}

	// Copy linguistic changes
	if l.LinguisticChanges != nil {
		clone.LinguisticChanges = make([]LinguisticChange, len(l.LinguisticChanges))
		copy(clone.LinguisticChanges, l.LinguisticChanges)
	}

	return clone
}

// String returns a string representation of the language.
func (l *Language) String() string {
	return fmt.Sprintf("Language(%s: %s)", l.ID.String(), l.Name)
}

// Interlingua returns the interlingua services for this language, initializing if needed.
func (l *Language) Interlingua() *InterlinguaServices {
	if l.interlinguaServices == nil {
		l.interlinguaServices = NewInterlinguaServices()
	}
	return l.interlinguaServices
}

// SetInterlinguaServices sets the interlingua services for this language.
func (l *Language) SetInterlinguaServices(services *InterlinguaServices) {
	l.interlinguaServices = services
}
