package evolution

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/chris-pikul/kismet-zero/lang/orthography"
)

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

var orthographicChangeTypeEnum = []string{
	"unknown",
	"script_reform",
	"borrowing",
	"simplification",
	"standardization",
	"innovation",
	"adaptation",
}

// String returns the string representation of the OrthographicChangeType.
func (oct OrthographicChangeType) String() string {
	if oct > OrthographicChangeTypeAdaptation {
		return orthographicChangeTypeEnum[0]
	}
	return orthographicChangeTypeEnum[oct]
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

// OrthographicEvolutionEngine manages orthographic changes in a language.
type OrthographicEvolutionEngine struct {
	rng    *rand.Rand
	config EvolutionConfig
}

// NewOrthographicEvolutionEngine creates a new orthographic evolution engine.
func NewOrthographicEvolutionEngine(config EvolutionConfig) *OrthographicEvolutionEngine {
	rng := rand.New(rand.NewPCG(uint64(config.Seed), 0))

	return &OrthographicEvolutionEngine{
		rng:    rng,
		config: config,
	}
}

// ApplyOrthographicChanges applies orthographic evolution to a writing system.
// Returns a list of changes that were applied.
func (oee *OrthographicEvolutionEngine) ApplyOrthographicChanges(
	writingSystem *orthography.WritingSystem,
	era string,
) []OrthographicChange {
	changes := make([]OrthographicChange, 0)

	// Apply script reforms (less common but significant)
	if oee.rng.Float32() < oee.config.ScriptReformRate {
		reformChanges := oee.applyScriptReforms(writingSystem, era)
		changes = append(changes, reformChanges...)
	}

	// Apply simplification changes
	if oee.rng.Float32() < oee.config.OrthographyChangeRate*0.8 {
		simplificationChanges := oee.applySimplificationChanges(writingSystem, era)
		changes = append(changes, simplificationChanges...)
	}

	// Apply standardization changes
	if oee.rng.Float32() < oee.config.OrthographyChangeRate*0.6 {
		standardizationChanges := oee.applyStandardizationChanges(writingSystem, era)
		changes = append(changes, standardizationChanges...)
	}

	// Apply innovation changes (least common)
	if oee.rng.Float32() < oee.config.OrthographyChangeRate*0.4 {
		innovationChanges := oee.applyInnovationChanges(writingSystem, era)
		changes = append(changes, innovationChanges...)
	}

	return changes
}

// applyScriptReforms applies systematic script reforms to the writing system.
func (oee *OrthographicEvolutionEngine) applyScriptReforms(
	writingSystem *orthography.WritingSystem,
	era string,
) []OrthographicChange {
	var changes []OrthographicChange

	// Example: Spelling reform
	if oee.rng.Float32() < 0.6 {
		change := &OrthographicChange{
			ID:                fmt.Sprintf("spelling_reform_%d", time.Now().UnixNano()),
			Type:              OrthographicChangeTypeScriptReform,
			Description:       "Implemented spelling reform",
			Details:           "Standardized and simplified spelling conventions",
			AffectedGraphemes: []string{"spelling_system"},
			ComplexityChange:  -0.2,
			Timestamp:         time.Now(),
			Era:               era,
			Trigger:           "orthographic_reform",
			Intensity:         0.6,
			ScriptReformType:  "spelling",
		}
		changes = append(changes, *change)
	}

	// Example: Character simplification
	if oee.rng.Float32() < 0.4 {
		change := &OrthographicChange{
			ID:                fmt.Sprintf("character_simplification_%d", time.Now().UnixNano()),
			Type:              OrthographicChangeTypeScriptReform,
			Description:       "Simplified complex characters",
			Details:           "Reduced stroke count and complexity of characters",
			AffectedGraphemes: []string{"complex_characters"},
			ComplexityChange:  -0.3,
			Timestamp:         time.Now(),
			Era:               era,
			Trigger:           "orthographic_reform",
			Intensity:         0.5,
			ScriptReformType:  "character",
		}
		changes = append(changes, *change)
	}

	return changes
}

// applySimplificationChanges applies orthographic simplification over time.
func (oee *OrthographicEvolutionEngine) applySimplificationChanges(
	writingSystem *orthography.WritingSystem,
	era string,
) []OrthographicChange {
	var changes []OrthographicChange

	// Example: Remove redundant graphemes
	if oee.rng.Float32() < 0.5 {
		change := &OrthographicChange{
			ID:                fmt.Sprintf("remove_redundant_graphemes_%d", time.Now().UnixNano()),
			Type:              OrthographicChangeTypeSimplification,
			Description:       "Removed redundant graphemes",
			Details:           "Eliminated duplicate or unnecessary writing symbols",
			AffectedGraphemes: []string{"redundant_symbols"},
			ComplexityChange:  -0.15,
			Timestamp:         time.Now(),
			Era:               era,
			Trigger:           "natural_evolution",
			Intensity:         0.4,
		}
		changes = append(changes, *change)
	}

	// Example: Simplify complex mappings
	if oee.rng.Float32() < 0.4 {
		change := &OrthographicChange{
			ID:               fmt.Sprintf("simplify_mappings_%d", time.Now().UnixNano()),
			Type:             OrthographicChangeTypeSimplification,
			Description:      "Simplified phoneme-grapheme mappings",
			Details:          "Reduced irregular and complex spelling patterns",
			AffectedMappings: []string{"irregular_mappings"},
			ComplexityChange: -0.2,
			Timestamp:        time.Now(),
			Era:              era,
			Trigger:          "natural_evolution",
			Intensity:        0.3,
		}
		changes = append(changes, *change)
	}

	return changes
}

// applyStandardizationChanges applies orthographic standardization over time.
func (oee *OrthographicEvolutionEngine) applyStandardizationChanges(
	writingSystem *orthography.WritingSystem,
	era string,
) []OrthographicChange {
	var changes []OrthographicChange

	// Example: Standardize spelling
	if oee.rng.Float32() < 0.6 {
		change := &OrthographicChange{
			ID:               fmt.Sprintf("standardize_spelling_%d", time.Now().UnixNano()),
			Type:             OrthographicChangeTypeStandardization,
			Description:      "Standardized spelling conventions",
			Details:          "Established consistent spelling rules and patterns",
			AffectedRules:    []string{"spelling_rules"},
			ComplexityChange: -0.1,
			Timestamp:        time.Now(),
			Era:              era,
			Trigger:          "standardization",
			Intensity:        0.5,
		}
		changes = append(changes, *change)
	}

	// Example: Standardize character forms
	if oee.rng.Float32() < 0.4 {
		change := &OrthographicChange{
			ID:                fmt.Sprintf("standardize_characters_%d", time.Now().UnixNano()),
			Type:              OrthographicChangeTypeStandardization,
			Description:       "Standardized character forms",
			Details:           "Unified variant character forms and styles",
			AffectedGraphemes: []string{"character_variants"},
			ComplexityChange:  -0.15,
			Timestamp:         time.Now(),
			Era:               era,
			Trigger:           "standardization",
			Intensity:         0.4,
		}
		changes = append(changes, *change)
	}

	return changes
}

// applyInnovationChanges introduces new orthographic features.
func (oee *OrthographicEvolutionEngine) applyInnovationChanges(
	writingSystem *orthography.WritingSystem,
	era string,
) []OrthographicChange {
	var changes []OrthographicChange

	// Example: Add new punctuation
	if oee.rng.Float32() < 0.5 {
		change := &OrthographicChange{
			ID:                fmt.Sprintf("add_punctuation_%d", time.Now().UnixNano()),
			Type:              OrthographicChangeTypeInnovation,
			Description:       "Added new punctuation marks",
			Details:           "Introduced new punctuation for better text organization",
			AffectedGraphemes: []string{"punctuation"},
			ComplexityChange:  0.1,
			Timestamp:         time.Now(),
			Era:               era,
			Trigger:           "innovation",
			Intensity:         0.3,
		}
		changes = append(changes, *change)
	}

	// Example: Add diacritical marks
	if oee.rng.Float32() < 0.3 {
		change := &OrthographicChange{
			ID:                fmt.Sprintf("add_diacritics_%d", time.Now().UnixNano()),
			Type:              OrthographicChangeTypeInnovation,
			Description:       "Added diacritical marks",
			Details:           "Introduced accent marks and diacritics for phonetic precision",
			AffectedGraphemes: []string{"diacritics"},
			ComplexityChange:  0.2,
			Timestamp:         time.Now(),
			Era:               era,
			Trigger:           "innovation",
			Intensity:         0.4,
		}
		changes = append(changes, *change)
	}

	return changes
}

// SimulateOrthographicBorrowing simulates borrowing of writing system features from another language.
func (oee *OrthographicEvolutionEngine) SimulateOrthographicBorrowing(
	targetWritingSystem *orthography.WritingSystem,
	sourceWritingSystem *orthography.WritingSystem,
	borrowingType string,
	intensity float32,
) []OrthographicChange {
	var changes []OrthographicChange

	// Determine what gets borrowed based on type and intensity
	switch borrowingType {
	case "graphemes":
		if oee.rng.Float32() < intensity {
			change := &OrthographicChange{
				ID:                fmt.Sprintf("borrow_graphemes_%d", time.Now().UnixNano()),
				Type:              OrthographicChangeTypeBorrowing,
				Description:       fmt.Sprintf("Borrowed graphemes from %s", sourceWritingSystem.Name),
				Details:           "Adopted writing symbols from source writing system",
				AffectedGraphemes: []string{"borrowed_symbols"},
				ComplexityChange:  0.1,
				Timestamp:         time.Now(),
				Era:               "contact_evolution",
				Trigger:           "orthographic_borrowing",
				Intensity:         intensity,
				BorrowedFrom:      sourceWritingSystem.Name,
			}
			changes = append(changes, *change)
		}

	case "mappings":
		if oee.rng.Float32() < intensity*0.7 {
			change := &OrthographicChange{
				ID:               fmt.Sprintf("borrow_mappings_%d", time.Now().UnixNano()),
				Type:             OrthographicChangeTypeBorrowing,
				Description:      fmt.Sprintf("Borrowed mapping patterns from %s", sourceWritingSystem.Name),
				Details:          "Adopted phoneme-grapheme mapping conventions",
				AffectedMappings: []string{"borrowed_mappings"},
				ComplexityChange: 0.05,
				Timestamp:        time.Now(),
				Era:              "contact_evolution",
				Trigger:          "orthographic_borrowing",
				Intensity:        intensity * 0.7,
				BorrowedFrom:     sourceWritingSystem.Name,
			}
			changes = append(changes, *change)
		}

	case "style":
		if oee.rng.Float32() < intensity*0.5 {
			change := &OrthographicChange{
				ID:               fmt.Sprintf("borrow_style_%d", time.Now().UnixNano()),
				Type:             OrthographicChangeTypeBorrowing,
				Description:      fmt.Sprintf("Borrowed writing style elements from %s", sourceWritingSystem.Name),
				Details:          "Adopted writing system style and organization",
				AffectedRules:    []string{"writing_style"},
				ComplexityChange: 0.15,
				Timestamp:        time.Now(),
				Era:              "contact_evolution",
				Trigger:          "orthographic_borrowing",
				Intensity:        intensity * 0.5,
				BorrowedFrom:     sourceWritingSystem.Name,
			}
			changes = append(changes, *change)
		}
	}

	return changes
}

// AdaptToPhonologicalChanges adapts the writing system to phonological changes.
func (oee *OrthographicEvolutionEngine) AdaptToPhonologicalChanges(
	writingSystem *orthography.WritingSystem,
	phonologicalChanges []LinguisticChange,
	era string,
) []OrthographicChange {
	var changes []OrthographicChange

	for _, phonoChange := range phonologicalChanges {
		if phonoChange.Type == ChangeTypeSoundShift {
			// Create orthographic adaptation for sound changes
			change := &OrthographicChange{
				ID:               fmt.Sprintf("phonological_adaptation_%d", time.Now().UnixNano()),
				Type:             OrthographicChangeTypeAdaptation,
				Description:      "Adapted writing system to phonological changes",
				Details:          fmt.Sprintf("Updated orthography to reflect: %s", phonoChange.Description),
				AffectedMappings: []string{"phoneme_grapheme_mappings"},
				ComplexityChange: 0.05,
				Timestamp:        time.Now(),
				Era:              era,
				Trigger:          "phonological_adaptation",
				Intensity:        0.4,
			}
			changes = append(changes, *change)
		}
	}

	return changes
}

// CalculateOrthographicComplexity calculates the overall complexity of a writing system.
func (oee *OrthographicEvolutionEngine) CalculateOrthographicComplexity(
	writingSystem *orthography.WritingSystem,
) float32 {
	complexity := 0.0

	// Base complexity from grapheme count
	if writingSystem != nil {
		complexity += float64(len(writingSystem.Graphemes)) * 0.01

		// Add complexity from mapping irregularity
		irregularMappings := 0
		for _, mapping := range writingSystem.Mappings {
			if len(mapping.Graphemes) > 1 {
				irregularMappings++
			}
		}
		complexity += float64(irregularMappings) * 0.02

		// Add complexity from writing style
		switch writingSystem.Style {
		case orthography.WritingStyleLogographic:
			complexity += 0.3
		case orthography.WritingStyleSyllabic:
			complexity += 0.2
		case orthography.WritingStyleAlphabetic:
			complexity += 0.1
		}
	}

	// Normalize to 0.0-1.0 range
	if complexity > 1.0 {
		complexity = 1.0
	}

	return float32(complexity)
}

// GenerateOrthographicReform creates a comprehensive orthographic reform plan.
func (oee *OrthographicEvolutionEngine) GenerateOrthographicReform(
	writingSystem *orthography.WritingSystem,
	reformType string,
	era string,
) OrthographicChange {

	reformTypeEnum := OrthographicChangeTypeScriptReform

	var description, details string
	var complexityChange float32

	switch reformType {
	case "spelling":
		description = "Comprehensive spelling reform"
		details = "Standardized spelling rules and eliminated irregularities"
		complexityChange = -0.25
	case "character":
		description = "Character simplification reform"
		details = "Reduced stroke count and simplified complex characters"
		complexityChange = -0.3
	case "system":
		description = "System-wide orthographic reform"
		details = "Complete overhaul of writing system organization and rules"
		complexityChange = -0.2
	default:
		description = "General orthographic reform"
		details = "Various improvements to writing system consistency and usability"
		complexityChange = -0.2
	}

	return OrthographicChange{
		ID:                fmt.Sprintf("orthographic_reform_%s_%d", reformType, time.Now().UnixNano()),
		Type:              reformTypeEnum,
		Description:       description,
		Details:           details,
		AffectedGraphemes: []string{"writing_system"},
		ComplexityChange:  complexityChange,
		Timestamp:         time.Now(),
		Era:               era,
		Trigger:           "orthographic_reform",
		Intensity:         0.8,
		ScriptReformType:  reformType,
	}
}
