package evolution

import (
	"fmt"
	"time"

	"github.com/chris-pikul/kismet-zero/lang/grammar"
	"github.com/chris-pikul/kismet-zero/lang/morphology"
)

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

var morphologicalChangeTypeEnum = []string{
	"unknown",
	"simplification",
	"regularization",
	"innovation",
	"loss",
	"analogy",
}

// String returns the string representation of the MorphologicalChangeType.
func (mct MorphologicalChangeType) String() string {
	if mct > MorphologicalChangeTypeAnalogy {
		return morphologicalChangeTypeEnum[0]
	}
	return morphologicalChangeTypeEnum[mct]
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

// MorphologicalEvolutionEngine manages morphological changes in a language.
type MorphologicalEvolutionEngine struct {
	BaseEngine
}

// NewMorphologicalEvolutionEngine creates a new morphological evolution engine.
func NewMorphologicalEvolutionEngine(config EvolutionConfig) *MorphologicalEvolutionEngine {
	return &MorphologicalEvolutionEngine{
		BaseEngine: NewBaseEngine(config),
	}
}

// ApplyMorphologicalChanges applies morphological evolution to a language's morphology and grammar.
// Returns a list of changes that were applied.
func (mee *MorphologicalEvolutionEngine) ApplyMorphologicalChanges(
	morphemeList *morphology.MorphemeList,
	grammar *grammar.Grammar,
	era string,
) []MorphologicalChange {
	changes := make([]MorphologicalChange, 0)

	// Apply simplification changes
	if mee.GetRandomFloat32() < mee.GetConfig().MorphologyChangeRate {
		simplificationChanges := mee.applySimplificationChanges(morphemeList, grammar, era)
		changes = append(changes, simplificationChanges...)
	}

	// Apply regularization changes
	if mee.GetRandomFloat32() < mee.GetConfig().MorphologyChangeRate*0.7 {
		regularizationChanges := mee.applyRegularizationChanges(morphemeList, grammar, era)
		changes = append(changes, regularizationChanges...)
	}

	// Apply innovation changes (less common)
	if mee.GetRandomFloat32() < mee.GetConfig().MorphologyChangeRate*0.3 {
		innovationChanges := mee.applyInnovationChanges(morphemeList, grammar, era)
		changes = append(changes, innovationChanges...)
	}

	return changes
}

// applySimplificationChanges reduces morphological complexity over time.
func (mee *MorphologicalEvolutionEngine) applySimplificationChanges(
	morphemeList *morphology.MorphemeList,
	grammar *grammar.Grammar,
	era string,
) []MorphologicalChange {
	var changes []MorphologicalChange

	// Example: General morphological simplification
	if mee.GetRandomFloat32() < 0.4 {
		change := &MorphologicalChange{
			ID:               fmt.Sprintf("morphology_simplification_%d", time.Now().Unix()),
			Type:             MorphologicalChangeTypeSimplification,
			Description:      "Simplified morphological system",
			Details:          "Reduced morphological complexity through various changes",
			AffectedFeatures: []string{"morphology"},
			ComplexityChange: -0.3,
			Timestamp:        time.Now(),
			Era:              era,
			Trigger:          "natural_evolution",
			Intensity:        0.4,
		}
		changes = append(changes, *change)
	}

	// Example: Agreement system simplification
	if mee.GetRandomFloat32() < 0.3 {
		change := &MorphologicalChange{
			ID:               fmt.Sprintf("agreement_simplification_%d", time.Now().Unix()),
			Type:             MorphologicalChangeTypeSimplification,
			Description:      "Simplified agreement system",
			Details:          "Reduced number of agreement categories",
			AffectedFeatures: []string{"agreement_system"},
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

// applyRegularizationChanges makes irregular forms more regular over time.
func (mee *MorphologicalEvolutionEngine) applyRegularizationChanges(
	morphemeList *morphology.MorphemeList,
	grammar *grammar.Grammar,
	era string,
) []MorphologicalChange {
	var changes []MorphologicalChange

	// Example: Regularize irregular verb forms
	if mee.GetRandomFloat32() < 0.5 {
		change := &MorphologicalChange{
			ID:                fmt.Sprintf("verb_regularization_%d", time.Now().Unix()),
			Type:              MorphologicalChangeTypeRegularization,
			Description:       "Regularized irregular verb forms",
			Details:           "Applied regular conjugation patterns to irregular verbs",
			AffectedMorphemes: []string{"irregular_verbs"},
			ComplexityChange:  -0.1,
			Timestamp:         time.Now(),
			Era:               era,
			Trigger:           "natural_evolution",
			Intensity:         0.3,
		}
		changes = append(changes, *change)
	}

	return changes
}

// applyInnovationChanges introduces new morphological features.
func (mee *MorphologicalEvolutionEngine) applyInnovationChanges(
	morphemeList *morphology.MorphemeList,
	grammar *grammar.Grammar,
	era string,
) []MorphologicalChange {
	var changes []MorphologicalChange

	// Example: Introduce new aspect markers
	if mee.GetRandomFloat32() < 0.3 {
		change := &MorphologicalChange{
			ID:               fmt.Sprintf("aspect_innovation_%d", time.Now().Unix()),
			Type:             MorphologicalChangeTypeInnovation,
			Description:      "Introduced new aspect markers",
			Details:          "Developed progressive and perfective aspect forms",
			AffectedFeatures: []string{"aspect_system"},
			ComplexityChange: 0.2,
			Timestamp:        time.Now(),
			Era:              era,
			Trigger:          "natural_evolution",
			Intensity:        0.4,
		}
		changes = append(changes, *change)
	}

	return changes
}

// CalculateMorphologicalComplexity calculates the overall complexity of a language's morphology.
func (mee *MorphologicalEvolutionEngine) CalculateMorphologicalComplexity(
	morphemeList *morphology.MorphemeList,
	grammar *grammar.Grammar,
) float32 {
	complexity := 0.0

	// Base complexity from morpheme count
	if morphemeList != nil {
		complexity += float64(len(*morphemeList)) * 0.01
	}

	// Add complexity from grammar features
	if grammar != nil {
		// For now, add a base complexity for having grammar
		// In a full implementation, you'd check specific features
		complexity += 0.3
	}

	// Normalize to 0.0-1.0 range
	if complexity > 1.0 {
		complexity = 1.0
	}

	return float32(complexity)
}
