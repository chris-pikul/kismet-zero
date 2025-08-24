package lang

import (
	"fmt"
	"math/rand/v2"
	"time"
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

// CommonOrthographicChanges provides a library of common orthographic changes.
var CommonOrthographicChanges = []struct {
	Type        OrthographicChangeType
	Description string
	Details     string
	Probability float32
	Era         string
	Complexity  float32
	Intensity   float32
	ScriptType  string
}{
	{
		Type:        OrthographicChangeTypeScriptReform,
		Description: "Spelling Reform",
		Details:     "Standardized and simplified spelling conventions",
		Probability: 0.4,
		Era:         "middle_evolution",
		Complexity:  -0.2,
		Intensity:   0.6,
		ScriptType:  "spelling",
	},
	{
		Type:        OrthographicChangeTypeScriptReform,
		Description: "Character Simplification",
		Details:     "Reduced stroke count and complexity of characters",
		Probability: 0.3,
		Era:         "late_evolution",
		Complexity:  -0.3,
		Intensity:   0.5,
		ScriptType:  "character",
	},
	{
		Type:        OrthographicChangeTypeSimplification,
		Description: "Remove Redundant Graphemes",
		Details:     "Eliminated duplicate or unnecessary writing symbols",
		Probability: 0.5,
		Era:         "early_evolution",
		Complexity:  -0.15,
		Intensity:   0.4,
		ScriptType:  "",
	},
	{
		Type:        OrthographicChangeTypeSimplification,
		Description: "Simplify Mappings",
		Details:     "Reduced irregular and complex spelling patterns",
		Probability: 0.4,
		Era:         "middle_evolution",
		Complexity:  -0.2,
		Intensity:   0.4,
		ScriptType:  "",
	},
	{
		Type:        OrthographicChangeTypeStandardization,
		Description: "Standardize Orthography",
		Details:     "Established consistent writing conventions",
		Probability: 0.6,
		Era:         "early_evolution",
		Complexity:  -0.1,
		Intensity:   0.5,
		ScriptType:  "",
	},
	{
		Type:        OrthographicChangeTypeInnovation,
		Description: "New Writing Features",
		Details:     "Introduced punctuation and formatting marks",
		Probability: 0.3,
		Era:         "late_evolution",
		Complexity:  0.2,
		Intensity:   0.4,
		ScriptType:  "",
	},
	{
		Type:        OrthographicChangeTypeAdaptation,
		Description: "Phonological Adaptation",
		Details:     "Updated writing system to reflect sound changes",
		Probability: 0.5,
		Era:         "middle_evolution",
		Complexity:  0.1,
		Intensity:   0.4,
		ScriptType:  "",
	},
}

// applyOrthographicChanges applies realistic orthographic changes to a language.
func applyOrthographicChanges(language *Language, timePeriod string, seed int64) []OrthographicChange {
	if language.Orthography == nil {
		return nil
	}

	rng := rand.New(rand.NewPCG(uint64(seed), 0))
	var changes []OrthographicChange

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

	// Apply applicable orthographic changes
	for _, changeTemplate := range CommonOrthographicChanges {
		if changeTemplate.Era == era && rng.Float32() < changeTemplate.Probability {
			change := OrthographicChange{
				ID:                fmt.Sprintf("orthography_%s_%d", changeTemplate.Type.String(), time.Now().Unix()),
				Type:              changeTemplate.Type,
				Description:       changeTemplate.Description,
				Details:           changeTemplate.Details,
				AffectedGraphemes: []string{"writing_system"},
				ComplexityChange:  changeTemplate.Complexity,
				Timestamp:         time.Now(),
				Era:               era,
				Trigger:           "natural_evolution",
				Intensity:         changeTemplate.Intensity,
				ScriptReformType:  changeTemplate.ScriptType,
			}
			changes = append(changes, change)
		}
	}

	return changes
}

// calculateOrthographicComplexity calculates the overall complexity of a language's writing system.
func calculateOrthographicComplexity(language *Language) float32 {
	complexity := 0.0

	// Base complexity from having orthography
	if language.Orthography != nil {
		complexity += 0.3
	}

	// Add complexity from writing system type
	if language.Orthography != nil {
		// This is a simplified calculation - in a full implementation,
		// you'd analyze the actual writing system characteristics
		complexity += 0.2
	}

	// Normalize to 0.0-1.0 range
	if complexity > 1.0 {
		complexity = 1.0
	}

	return float32(complexity)
}
