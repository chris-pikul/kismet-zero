package lang

import (
	"time"
)

// LinguisticChangeType represents the type of linguistic change.
type LinguisticChangeType string

const (
	LinguisticChangeTypeSound         LinguisticChangeType = "sound"
	LinguisticChangeTypeMorphological LinguisticChangeType = "morphological"
	LinguisticChangeTypeOrthographic  LinguisticChangeType = "orthographic"
	LinguisticChangeTypeCultural      LinguisticChangeType = "cultural"
	LinguisticChangeTypeDialectal     LinguisticChangeType = "dialectal"
)

var linguisticChangeTypeEnum = []string{
	"unknown",
	"sound",
	"morphological",
	"orthographic",
}

// String returns the string representation of the LinguisticChangeType.
func (lct LinguisticChangeType) String() string {
	return string(lct)
}

// LinguisticChange represents a unified change that has been applied to a language.
type LinguisticChange struct {
	ID          string               `json:"id"`
	Type        LinguisticChangeType `json:"type"`
	Description string               `json:"description"`
	Details     string               `json:"details,omitempty"`

	// Change-specific data
	SoundChange         *SoundChange         `json:"soundChange,omitempty"`
	MorphologicalChange *MorphologicalChange `json:"morphologicalChange,omitempty"`
	OrthographicChange  *OrthographicChange  `json:"orthographicChange,omitempty"`

	// Impact on complexity
	ComplexityChange float32 `json:"complexityChange"` // -1.0 to 1.0, negative = simplification

	// Context
	Timestamp time.Time `json:"timestamp"`
	Era       string    `json:"era,omitempty"`
	Trigger   string    `json:"trigger,omitempty"`
	Intensity float32   `json:"intensity"` // 0.0 to 1.0
}

// NewSoundChange creates a new linguistic change from a sound change.
func NewSoundChange(change SoundChange) LinguisticChange {
	return LinguisticChange{
		ID:               change.ID,
		Type:             LinguisticChangeTypeSound,
		Description:      change.Description,
		Details:          change.Description, // Use Description as Details since SoundChange doesn't have Details
		SoundChange:      &change,
		ComplexityChange: 0.0,        // Sound changes don't directly affect complexity
		Timestamp:        time.Now(), // Use current time since SoundChange doesn't have Timestamp
		Era:              change.Era,
		Trigger:          "natural_evolution", // Default trigger since SoundChange doesn't have Trigger
		Intensity:        0.5,                 // Default intensity since SoundChange doesn't have Intensity
	}
}

// NewMorphologicalChange creates a new linguistic change from a morphological change.
func NewMorphologicalChange(change MorphologicalChange) LinguisticChange {
	return LinguisticChange{
		ID:                  change.ID,
		Type:                LinguisticChangeTypeMorphological,
		Description:         change.Description,
		Details:             change.Details,
		MorphologicalChange: &change,
		ComplexityChange:    change.ComplexityChange,
		Timestamp:           change.Timestamp,
		Era:                 change.Era,
		Trigger:             change.Trigger,
		Intensity:           change.Intensity,
	}
}

// NewOrthographicChange creates a new linguistic change from an orthographic change.
func NewOrthographicChange(change OrthographicChange) LinguisticChange {
	return LinguisticChange{
		ID:                 change.ID,
		Type:               LinguisticChangeTypeOrthographic,
		Description:        change.Description,
		Details:            change.Details,
		OrthographicChange: &change,
		ComplexityChange:   change.ComplexityChange,
		Timestamp:          change.Timestamp,
		Era:                change.Era,
		Trigger:            change.Trigger,
		Intensity:          change.Intensity,
	}
}

// AddLinguisticChange adds a linguistic change to the language's change history.
func (l *Language) AddLinguisticChange(change LinguisticChange) {
	if l.LinguisticChanges == nil {
		l.LinguisticChanges = make([]LinguisticChange, 0)
	}
	l.LinguisticChanges = append(l.LinguisticChanges, change)
}

// GetLinguisticChanges returns all linguistic changes for a specific type.
func (l *Language) GetLinguisticChanges(changeType LinguisticChangeType) []LinguisticChange {
	if l.LinguisticChanges == nil {
		return nil
	}

	var changes []LinguisticChange
	for _, change := range l.LinguisticChanges {
		if change.Type == changeType {
			changes = append(changes, change)
		}
	}
	return changes
}

// GetLinguisticChangesByEra returns all linguistic changes for a specific era.
func (l *Language) GetLinguisticChangesByEra(era string) []LinguisticChange {
	if l.LinguisticChanges == nil {
		return nil
	}

	var changes []LinguisticChange
	for _, change := range l.LinguisticChanges {
		if change.Era == era {
			changes = append(changes, change)
		}
	}
	return changes
}

// GetTotalComplexityChange calculates the total complexity change from all linguistic changes.
func (l *Language) GetTotalComplexityChange() float32 {
	if l.LinguisticChanges == nil {
		return 0.0
	}

	totalChange := 0.0
	for _, change := range l.LinguisticChanges {
		totalChange += float64(change.ComplexityChange)
	}

	// Normalize to -1.0 to 1.0 range
	if totalChange < -1.0 {
		totalChange = -1.0
	} else if totalChange > 1.0 {
		totalChange = 1.0
	}

	return float32(totalChange)
}

// GetEvolutionSummary returns a summary of all linguistic changes applied to the language.
func (l *Language) GetEvolutionSummary() map[string]interface{} {
	if l.LinguisticChanges == nil {
		return map[string]interface{}{
			"totalChanges":     0,
			"eras":             []string{},
			"types":            []string{},
			"complexityChange": 0.0,
		}
	}

	eras := make(map[string]bool)
	types := make(map[string]bool)
	totalComplexityChange := l.GetTotalComplexityChange()

	for _, change := range l.LinguisticChanges {
		eras[change.Era] = true
		types[change.Type.String()] = true
	}

	// Convert maps to slices
	eraList := make([]string, 0, len(eras))
	for era := range eras {
		eraList = append(eraList, era)
	}

	typeList := make([]string, 0, len(types))
	for changeType := range types {
		typeList = append(typeList, changeType)
	}

	return map[string]interface{}{
		"totalChanges":     len(l.LinguisticChanges),
		"eras":             eraList,
		"types":            typeList,
		"complexityChange": totalComplexityChange,
		"lastChange":       l.LinguisticChanges[len(l.LinguisticChanges)-1].Timestamp,
	}
}

// StoreLinguisticChanges stores all linguistic changes in the language's change history.
func StoreLinguisticChanges(language *Language, soundChanges []string, morphologicalChanges []MorphologicalChange, orthographicChanges []OrthographicChange, timePeriod string) {
	// Initialize the linguistic changes slice if it doesn't exist
	if language.LinguisticChanges == nil {
		language.LinguisticChanges = make([]LinguisticChange, 0)
	}

	// Store sound changes
	for _, soundChangeName := range soundChanges {
		// Find the corresponding sound change template
		for _, template := range CommonSoundChanges {
			if template.Name == soundChangeName {
				// Create a sound change with current timestamp
				soundChange := SoundChange{
					ID:          template.ID,
					Name:        template.Name,
					Description: template.Description,
					Probability: template.Probability,
					Era:         template.Era,
				}

				// Convert to linguistic change and store
				linguisticChange := NewSoundChange(soundChange)
				language.AddLinguisticChange(linguisticChange)
				break
			}
		}
	}

	// Store morphological changes
	for _, morphChange := range morphologicalChanges {
		linguisticChange := NewMorphologicalChange(morphChange)
		language.AddLinguisticChange(linguisticChange)
	}

	// Store orthographic changes
	for _, orthoChange := range orthographicChanges {
		linguisticChange := NewOrthographicChange(orthoChange)
		language.AddLinguisticChange(linguisticChange)
	}
}
