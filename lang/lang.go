package lang

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/chris-pikul/kismet-zero/lang/grammar"
	"github.com/chris-pikul/kismet-zero/lang/interlingua"
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

// LinguisticChange represents a unified change that has been applied to a language.
type LinguisticChange struct {
	ID          string               `json:"id"`
	Type        LinguisticChangeType `json:"type"`
	Description string               `json:"description"`
	Details     string               `json:"details,omitempty"`

	// Impact on complexity
	ComplexityChange float32 `json:"complexityChange"` // -1.0 to 1.0, negative = simplification

	// Context
	Timestamp time.Time `json:"timestamp"`
	Era       string    `json:"era,omitempty"`
	Trigger   string    `json:"trigger,omitempty"`
	Intensity float32   `json:"intensity"` // 0.0 to 1.0
}

// CulturalInfluenceResult represents the result of cultural influence on a language.
type CulturalInfluenceResult struct {
	LanguageID          string      `json:"languageId"`
	ContactType         ContactType `json:"contactType"`
	Duration            string      `json:"duration"`
	Intensity           float32     `json:"intensity"`
	ChangesApplied      []string    `json:"changesApplied"`
	ComplexityChange    float32     `json:"complexityChange"`
	Timestamp           time.Time   `json:"timestamp"`
	Description         string      `json:"description"`
	LexicalChanges      []string    `json:"lexicalChanges,omitempty"`
	PhonologicalChanges []string    `json:"phonologicalChanges,omitempty"`
	OrthographicChanges []string    `json:"orthographicChanges,omitempty"`
	GrammaticalChanges  []string    `json:"grammaticalChanges,omitempty"`
}

// LinguisticChangeType represents the type of linguistic change.
type LinguisticChangeType string

const (
	LinguisticChangeTypeSound         LinguisticChangeType = "sound"
	LinguisticChangeTypeMorphological LinguisticChangeType = "morphological"
	LinguisticChangeTypeOrthographic  LinguisticChangeType = "orthographic"
	LinguisticChangeTypeCultural      LinguisticChangeType = "cultural"
	LinguisticChangeTypeDialectal     LinguisticChangeType = "dialectal"
)

// String returns the string representation of the LinguisticChangeType.
func (lct LinguisticChangeType) String() string {
	return string(lct)
}

// SoundChange represents a phonological change that can occur in a language.
type SoundChange struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Probability float32 `json:"probability"` // 0.0 to 1.0
	Era         string  `json:"era"`         // When this change typically happens
}

// MorphologicalChangeType represents the specific type of morphological change.
type MorphologicalChangeType byte

const (
	MorphologicalChangeTypeUnknown        MorphologicalChangeType = iota
	MorphologicalChangeTypeSimplification                         // Reduction in complexity
	MorphologicalChangeTypeRegularization                         // Making irregular forms regular
	MorphologicalChangeTypeInnovation                             // New morphological features
	MorphologicalChangeTypeLoss                                   // Loss of morphological features
	MorphologicalChangeTypeAnalogy                                // Change by analogy with other forms
)

// String returns the string representation of the MorphologicalChangeType.
func (mct MorphologicalChangeType) String() string {
	if mct > MorphologicalChangeTypeAnalogy {
		return "unknown"
	}
	return []string{"unknown", "simplification", "regularization", "innovation", "loss", "analogy"}[mct]
}

// MorphologicalChange represents a change to the morphological system of a language.
type MorphologicalChange struct {
	ID          string                  `json:"id"`
	Type        MorphologicalChangeType `json:"type"`
	Description string                  `json:"description"`
	Details     string                  `json:"details,omitempty"`

	// Affected components
	AffectedMorphemes []string `json:"affectedMorphemes,omitempty"`
	AffectedRules     []string `json:"affectedRules,omitempty"`
	AffectedFeatures  []string `json:"affectedFeatures,omitempty"`

	// Impact on complexity
	ComplexityChange float32 `json:"complexityChange"` // -1.0 to 1.0, negative = simplification

	// Context
	Timestamp time.Time `json:"timestamp"`
	Era       string    `json:"era,omitempty"`
	Trigger   string    `json:"trigger,omitempty"`
	Intensity float32   `json:"intensity"` // 0.0 to 1.0
}

// OrthographicChangeType represents the specific type of orthographic change.
type OrthographicChangeType byte

const (
	OrthographicChangeTypeUnknown         OrthographicChangeType = iota
	OrthographicChangeTypeScriptReform                           // Systematic script reforms
	OrthographicChangeTypeBorrowing                              // Borrowing from other writing systems
	OrthographicChangeTypeSimplification                         // Simplifying complex characters
	OrthographicChangeTypeStandardization                        // Standardizing spelling/characters
	OrthographicChangeTypeInnovation                             // New characters or writing features
	OrthographicChangeTypeAdaptation                             // Adapting to phonological changes
)

// String returns the string representation of the OrthographicChangeType.
func (oct OrthographicChangeType) String() string {
	if oct > OrthographicChangeTypeAdaptation {
		return "unknown"
	}
	return []string{"unknown", "script_reform", "borrowing", "simplification", "standardization", "innovation", "adaptation"}[oct]
}

// OrthographicChange represents a change to the writing system of a language.
type OrthographicChange struct {
	ID          string                 `json:"id"`
	Type        OrthographicChangeType `json:"type"`
	Description string                 `json:"description"`
	Details     string                 `json:"details,omitempty"`

	// Affected components
	AffectedGraphemes []string `json:"affectedGraphemes,omitempty"`
	AffectedMappings  []string `json:"affectedMappings,omitempty"`
	AffectedRules     []string `json:"affectedRules,omitempty"`

	// Impact on complexity
	ComplexityChange float32 `json:"complexityChange"` // -1.0 to 1.0, negative = simplification

	// Context
	Timestamp time.Time `json:"timestamp"`
	Era       string    `json:"era,omitempty"`
	Trigger   string    `json:"trigger,omitempty"`
	Intensity float32   `json:"intensity"` // 0.0 to 1.0

	// Orthography-specific fields
	ScriptReformType string `json:"scriptReformType,omitempty"` // "spelling", "character", "system"
	BorrowedFrom     string `json:"borrowedFrom,omitempty"`     // Source writing system
}

// BorrowingPattern represents a pattern of linguistic borrowing.
type BorrowingPattern struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"` // "lexical", "phonological", "grammatical", "orthographic"
	Description string      `json:"description"`
	Probability float32     `json:"probability"` // 0.0 to 1.0
	Intensity   float32     `json:"intensity"`   // 0.0 to 1.0
	ContactType ContactType `json:"contactType"`

	// Enhanced borrowing mechanisms
	PhonologicalRules []PhonologicalRule `json:"phonologicalRules,omitempty"`
	GrammaticalRules  []GrammaticalRule  `json:"grammaticalRules,omitempty"`
	LexicalRules      []LexicalRule      `json:"lexicalRules,omitempty"`
	OrthographicRules []OrthographicRule `json:"orthographicRules,omitempty"`
}

// PhonologicalRule represents a rule for phonological borrowing and adaptation.
type PhonologicalRule struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Type        string  `json:"type"` // "sound_addition", "sound_modification", "phonotactic_change"
	SourceSound string  `json:"sourceSound,omitempty"`
	TargetSound string  `json:"targetSound,omitempty"`
	Context     string  `json:"context,omitempty"` // "word_initial", "word_final", "intervocalic", etc.
	Probability float32 `json:"probability"`
}

// GrammaticalRule represents a rule for grammatical borrowing and adaptation.
type GrammaticalRule struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Type        string  `json:"type"`     // "morphology_addition", "syntax_change", "agreement_modification"
	Category    string  `json:"category"` // "noun", "verb", "adjective", "particle"
	Feature     string  `json:"feature"`  // "case", "tense", "aspect", "mood"
	Probability float32 `json:"probability"`
}

// LexicalRule represents a rule for lexical borrowing and adaptation.
type LexicalRule struct {
	ID             string  `json:"id"`
	Description    string  `json:"description"`
	Type           string  `json:"type"`           // "word_borrowing", "semantic_shift", "calque_formation"
	SemanticField  string  `json:"semanticField"`  // "trade", "religion", "technology", "culture"
	AdaptationType string  `json:"adaptationType"` // "phonological", "morphological", "semantic"
	Probability    float32 `json:"probability"`
}

// OrthographicRule represents a rule for orthographic borrowing and adaptation.
type OrthographicRule struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Type        string  `json:"type"`                 // "script_adoption", "spelling_reform", "diacritic_addition"
	ScriptType  string  `json:"scriptType,omitempty"` // "alphabet", "syllabary", "logographic"
	ReformType  string  `json:"reformType,omitempty"` // "simplification", "standardization", "foreign_adoption"
	Probability float32 `json:"probability"`
}

// ContactType represents the type of language contact.
type ContactType byte

const (
	ContactTypeUnknown     ContactType = iota
	ContactTypeTrade                   // Trade relationships
	ContactTypeConquest                // Military conquest
	ContactTypeMigration               // Population migration
	ContactTypeCultural                // Cultural exchange
	ContactTypeReligious               // Religious influence
	ContactTypeEducational             // Educational exchange
)

// String returns the string representation of the ContactType.
func (ct ContactType) String() string {
	if ct > ContactTypeEducational {
		return "unknown"
	}
	return []string{"unknown", "trade", "conquest", "migration", "cultural", "religious", "educational"}[ct]
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
	interlinguaServices *interlingua.InterlinguaServices `json:"-"`
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
func (l *Language) Interlingua() *interlingua.InterlinguaServices {
	if l.interlinguaServices == nil {
		l.interlinguaServices = interlingua.NewInterlinguaServices()
	}
	return l.interlinguaServices
}

// SetInterlinguaServices sets the interlingua services for this language.
func (l *Language) SetInterlinguaServices(services *interlingua.InterlinguaServices) {
	l.interlinguaServices = services
}

// applySoundChanges applies realistic sound changes to a language's phonology.
func applySoundChanges(language *Language, timePeriod string, seed int64) []string {
	// This is a placeholder implementation
	// In a full implementation, this would apply actual sound changes
	return []string{"placeholder_sound_change"}
}

// applyMorphologicalChanges applies realistic morphological changes to a language.
func applyMorphologicalChanges(language *Language, timePeriod string, seed int64) []MorphologicalChange {
	// This is a placeholder implementation
	// In a full implementation, this would apply actual morphological changes
	return []MorphologicalChange{}
}

// applyOrthographicChanges applies realistic orthographic changes to a language.
func applyOrthographicChanges(language *Language, timePeriod string, seed int64) []OrthographicChange {
	// This is a placeholder implementation
	// In a full implementation, this would apply actual orthographic changes
	return []OrthographicChange{}
}

// StoreLinguisticChanges stores linguistic changes in the language.
func StoreLinguisticChanges(language *Language, soundChanges []string, morphologicalChanges []MorphologicalChange, orthographicChanges []OrthographicChange, timePeriod string) {
	// Store sound changes
	for _, change := range soundChanges {
		language.AddLinguisticChange(LinguisticChange{
			ID:          fmt.Sprintf("sound_%d", time.Now().Unix()),
			Type:        LinguisticChangeTypeSound,
			Description: change,
			Timestamp:   time.Now(),
			Era:         timePeriod,
			Trigger:     "natural_evolution",
			Intensity:   0.5,
		})
	}

	// Store morphological changes
	for _, change := range morphologicalChanges {
		language.AddLinguisticChange(LinguisticChange{
			ID:               change.ID,
			Type:             LinguisticChangeTypeMorphological,
			Description:      change.Description,
			Details:          change.Details,
			ComplexityChange: change.ComplexityChange,
			Timestamp:        change.Timestamp,
			Era:              change.Era,
			Trigger:          change.Trigger,
			Intensity:        change.Intensity,
		})
	}

	// Store orthographic changes
	for _, change := range orthographicChanges {
		language.AddLinguisticChange(LinguisticChange{
			ID:               change.ID,
			Type:             LinguisticChangeTypeOrthographic,
			Description:      change.Description,
			Details:          change.Details,
			ComplexityChange: change.ComplexityChange,
			Timestamp:        change.Timestamp,
			Era:              change.Era,
			Trigger:          change.Trigger,
			Intensity:        change.Intensity,
		})
	}
}

// AddLinguisticChange adds a linguistic change to the language's change history.
func (l *Language) AddLinguisticChange(change LinguisticChange) {
	if l.LinguisticChanges == nil {
		l.LinguisticChanges = make([]LinguisticChange, 0)
	}
	l.LinguisticChanges = append(l.LinguisticChanges, change)
}

// GetLinguisticChanges returns all linguistic changes for a language.
func (l *Language) GetLinguisticChanges() []LinguisticChange {
	if l.LinguisticChanges == nil {
		return make([]LinguisticChange, 0)
	}
	return l.LinguisticChanges
}

// GetLinguisticChangesByType returns linguistic changes of a specific type.
func (l *Language) GetLinguisticChangesByType(changeType LinguisticChangeType) []LinguisticChange {
	var filteredChanges []LinguisticChange

	for _, change := range l.LinguisticChanges {
		if change.Type == changeType {
			filteredChanges = append(filteredChanges, change)
		}
	}

	return filteredChanges
}

// GetLinguisticChangesByEra returns linguistic changes from a specific era.
func (l *Language) GetLinguisticChangesByEra(era string) []LinguisticChange {
	var filteredChanges []LinguisticChange

	for _, change := range l.LinguisticChanges {
		if change.Era == era {
			filteredChanges = append(filteredChanges, change)
		}
	}

	return filteredChanges
}

// ApplyCulturalInfluence applies cultural influence to a language.
func ApplyCulturalInfluence(language1 *Language, language2 *Language, contactType ContactType, intensity float32, duration time.Duration, seed int64) (*CulturalInfluenceResult, error) {
	// This is a placeholder implementation
	// In a full implementation, this would apply actual cultural influence
	return &CulturalInfluenceResult{
		LanguageID:          language1.ID.String(),
		ContactType:         contactType,
		Duration:            duration.String(),
		Intensity:           intensity,
		ChangesApplied:      []string{"placeholder_cultural_change"},
		ComplexityChange:    0.1,
		Timestamp:           time.Now(),
		Description:         "Placeholder cultural influence",
		LexicalChanges:      []string{},
		PhonologicalChanges: []string{},
		OrthographicChanges: []string{},
		GrammaticalChanges:  []string{},
	}, nil
}

// CreateCulturalInfluenceChange creates a cultural influence change from an existing result.
func CreateCulturalInfluenceChange(result *CulturalInfluenceResult, description string) LinguisticChange {
	return LinguisticChange{
		ID:               fmt.Sprintf("cultural_%d", time.Now().Unix()),
		Type:             LinguisticChangeTypeCultural,
		Description:      description,
		Details:          result.Description,
		ComplexityChange: result.ComplexityChange,
		Timestamp:        result.Timestamp,
		Era:              "contact",
		Trigger:          "cultural_influence",
		Intensity:        result.Intensity,
	}
}

// GetTotalComplexityChange calculates the total complexity change from all linguistic changes.
func (l *Language) GetTotalComplexityChange() float32 {
	totalChange := float32(0.0)
	for _, change := range l.LinguisticChanges {
		totalChange += change.ComplexityChange
	}
	return totalChange
}

// GetEvolutionSummary returns a summary of the language's evolution.
func (l *Language) GetEvolutionSummary() map[string]interface{} {
	summary := make(map[string]interface{})

	// Count total changes
	summary["totalChanges"] = len(l.LinguisticChanges)

	// Count changes by type
	typeCounts := make(map[LinguisticChangeType]int)
	for _, change := range l.LinguisticChanges {
		typeCounts[change.Type]++
	}
	summary["changesByType"] = typeCounts

	// Count changes by era
	eraCounts := make(map[string]int)
	for _, change := range l.LinguisticChanges {
		eraCounts[change.Era]++
	}
	summary["changesByEra"] = eraCounts

	// Total complexity change
	summary["totalComplexityChange"] = l.GetTotalComplexityChange()

	return summary
}
