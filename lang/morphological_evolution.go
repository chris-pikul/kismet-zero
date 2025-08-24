package lang

import (
	"fmt"
	"math/rand/v2"
	"time"
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

// CommonMorphologicalChanges provides a library of common morphological changes.
var CommonMorphologicalChanges = []struct {
	Type        MorphologicalChangeType
	Description string
	Details     string
	Probability float32
	Era         string
	Complexity  float32
	Intensity   float32
}{
	{
		Type:        MorphologicalChangeTypeSimplification,
		Description: "Morphological Simplification",
		Details:     "Reduced morphological complexity through various changes",
		Probability: 0.6,
		Era:         "early_evolution",
		Complexity:  -0.3,
		Intensity:   0.4,
	},
	{
		Type:        MorphologicalChangeTypeSimplification,
		Description: "Agreement System Simplification",
		Details:     "Reduced number of agreement categories",
		Probability: 0.4,
		Era:         "middle_evolution",
		Complexity:  -0.2,
		Intensity:   0.3,
	},
	{
		Type:        MorphologicalChangeTypeRegularization,
		Description: "Verb Regularization",
		Details:     "Applied regular conjugation patterns to irregular verbs",
		Probability: 0.5,
		Era:         "middle_evolution",
		Complexity:  -0.1,
		Intensity:   0.3,
	},
	{
		Type:        MorphologicalChangeTypeInnovation,
		Description: "Aspect System Innovation",
		Details:     "Developed progressive and perfective aspect forms",
		Probability: 0.3,
		Era:         "late_evolution",
		Complexity:  0.2,
		Intensity:   0.4,
	},
	{
		Type:        MorphologicalChangeTypeLoss,
		Description: "Case System Reduction",
		Details:     "Lost some grammatical cases, simplified declension",
		Probability: 0.4,
		Era:         "late_evolution",
		Complexity:  -0.2,
		Intensity:   0.3,
	},
}

// applyMorphologicalChanges applies realistic morphological changes to a language.
func applyMorphologicalChanges(language *Language, timePeriod string, seed int64) []MorphologicalChange {
	if language.Morphology == nil && language.Grammar == nil {
		return nil
	}

	rng := rand.New(rand.NewPCG(uint64(seed), 0))
	var changes []MorphologicalChange

	// Determine which era we're in based on time period
	var era string
	switch timePeriod {
	case "100_years", "short":
		era = "early_evolution"
	case "500_years", "medium":
		era = "middle_evolution"
	case "1000_years", "long":
		era = "late_evolution"
	case "2000_years", "very_long":
		era = "late_evolution"
	default:
		era = "evolution"
	}

	// Apply applicable morphological changes
	for _, changeTemplate := range CommonMorphologicalChanges {
		if changeTemplate.Era == era && rng.Float32() < changeTemplate.Probability {
			change := MorphologicalChange{
				ID:               fmt.Sprintf("morphology_%s_%d", changeTemplate.Type.String(), time.Now().Unix()),
				Type:             changeTemplate.Type,
				Description:      changeTemplate.Description,
				Details:          changeTemplate.Details,
				AffectedFeatures: []string{"morphology"},
				ComplexityChange: changeTemplate.Complexity,
				Timestamp:        time.Now(),
				Era:              era,
				Trigger:          "natural_evolution",
				Intensity:        changeTemplate.Intensity,
			}
			changes = append(changes, change)
		}
	}

	return changes
}

// calculateMorphologicalComplexity calculates the overall complexity of a language's morphology.
func calculateMorphologicalComplexity(language *Language) float32 {
	complexity := 0.0

	// Base complexity from having morphology
	if language.Morphology != nil {
		complexity += 0.3
	}

	// Add complexity from grammar features
	if language.Grammar != nil {
		complexity += 0.3
	}

	// Normalize to 0.0-1.0 range
	if complexity > 1.0 {
		complexity = 1.0
	}

	return float32(complexity)
}
