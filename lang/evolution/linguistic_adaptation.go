package evolution

import (
	"fmt"
	"time"

	"github.com/chris-pikul/kismet-zero/lang"
)

// LinguisticAdaptationEngine manages the application of linguistic adaptations to languages.
type LinguisticAdaptationEngine struct {
	BaseEngine
}

// NewLinguisticAdaptationEngine creates a new linguistic adaptation engine.
func NewLinguisticAdaptationEngine(config EvolutionConfig) *LinguisticAdaptationEngine {
	return &LinguisticAdaptationEngine{
		BaseEngine: NewBaseEngine(config),
	}
}

// ApplyPhonologicalAdaptation applies a phonological adaptation to a language.
func (lae *LinguisticAdaptationEngine) ApplyPhonologicalAdaptation(
	lang *lang.Language,
	adaptation *PhonologicalAdaptation,
) (*LinguisticChange, error) {
	if adaptation == nil {
		return nil, fmt.Errorf("adaptation cannot be nil")
	}

	// Create the linguistic change
	change := &LinguisticChange{
		ID:          adaptation.ID,
		Type:        ChangeTypeSoundShift,
		Direction:   ChangeDirectionModifying,
		Description: adaptation.Description,
		Details:     "Phonological adaptation applied",
		Timestamp:   time.Now(),
		Era:         "adaptation",
		Trigger:     "phonological_adaptation",
		Intensity:   adaptation.Intensity,
	}

	// Apply phoneme changes if specified
	if len(adaptation.AddedPhonemes) > 0 {
		// Add new phonemes to the phonology
		for _, phoneme := range adaptation.AddedPhonemes {
			change.AffectedPhonemes = append(change.AffectedPhonemes, phoneme.Symbol)
		}
	}

	if len(adaptation.ModifiedPhonemes) > 0 {
		// Mark modified phonemes
		for _, phoneme := range adaptation.ModifiedPhonemes {
			change.AffectedPhonemes = append(change.AffectedPhonemes, phoneme.Symbol)
		}
	}

	if len(adaptation.RemovedPhonemes) > 0 {
		// Mark removed phonemes
		change.AffectedPhonemes = append(change.AffectedPhonemes, adaptation.RemovedPhonemes...)
	}

	return change, nil
}

// ApplyGrammaticalAdaptation applies a grammatical adaptation to a language.
func (lae *LinguisticAdaptationEngine) ApplyGrammaticalAdaptation(
	lang *lang.Language,
	adaptation *GrammaticalAdaptation,
) (*LinguisticChange, error) {
	if adaptation == nil {
		return nil, fmt.Errorf("adaptation cannot be nil")
	}

	// Create the linguistic change
	change := &LinguisticChange{
		ID:          adaptation.ID,
		Type:        ChangeTypeMorphological,
		Direction:   ChangeDirectionAdditive,
		Description: adaptation.Description,
		Details:     "Grammatical adaptation applied",
		Timestamp:   time.Now(),
		Era:         "adaptation",
		Trigger:     "grammatical_adaptation",
		Intensity:   adaptation.Intensity,
	}

	// Apply morphological changes if specified
	if len(adaptation.NewCases) > 0 || len(adaptation.NewNumbers) > 0 || len(adaptation.NewGenders) > 0 {
		change.Direction = ChangeDirectionAdditive
		change.Details = fmt.Sprintf("%s: Added new grammatical features", change.Details)
	}

	if len(adaptation.NewTenses) > 0 || len(adaptation.NewAspects) > 0 || len(adaptation.NewMoods) > 0 {
		change.Direction = ChangeDirectionAdditive
		change.Details = fmt.Sprintf("%s: Added new verbal features", change.Details)
	}

	if len(adaptation.NewWordOrders) > 0 {
		change.Direction = ChangeDirectionModifying
		change.Details = fmt.Sprintf("%s: Modified word order", change.Details)
	}

	return change, nil
}

// ApplyMorphologicalAdaptation applies a morphological adaptation to a language.
func (lae *LinguisticAdaptationEngine) ApplyMorphologicalAdaptation(
	lang *lang.Language,
	adaptation *MorphologicalAdaptation,
) (*LinguisticChange, error) {
	if adaptation == nil {
		return nil, fmt.Errorf("adaptation cannot be nil")
	}

	// Create the linguistic change
	change := &LinguisticChange{
		ID:          adaptation.ID,
		Type:        ChangeTypeMorphological,
		Direction:   ChangeDirectionModifying,
		Description: adaptation.Description,
		Details:     "Morphological adaptation applied",
		Timestamp:   time.Now(),
		Era:         "adaptation",
		Trigger:     "morphological_adaptation",
		Intensity:   adaptation.Intensity,
	}

	// Apply morpheme changes if specified
	if len(adaptation.NewMorphemes) > 0 {
		change.Direction = ChangeDirectionAdditive
		change.Details = fmt.Sprintf("%s: Added new morphemes", change.Details)
	}

	if len(adaptation.ModifiedMorphemes) > 0 {
		change.Direction = ChangeDirectionModifying
		change.Details = fmt.Sprintf("%s: Modified existing morphemes", change.Details)
	}

	if len(adaptation.RemovedMorphemes) > 0 {
		change.Direction = ChangeDirectionSubtractive
		change.Details = fmt.Sprintf("%s: Removed morphemes", change.Details)
	}

	return change, nil
}

// ApplyBorrowingPattern applies a borrowing pattern to a language.
func (lae *LinguisticAdaptationEngine) ApplyBorrowingPattern(
	lang *lang.Language,
	pattern *BorrowingPattern,
	source string,
	contactType string,
) (*LinguisticChange, error) {
	if pattern == nil {
		return nil, fmt.Errorf("pattern cannot be nil")
	}

	// Create the linguistic change based on the borrowing pattern characteristics
	change := &LinguisticChange{
		ID:               fmt.Sprintf("borrowing_%s_%d", source, time.Now().Unix()),
		Type:             ChangeTypeContact,
		Direction:        ChangeDirectionAdditive,
		Description:      fmt.Sprintf("Applied borrowing pattern from %s", source),
		Details:          "Borrowing pattern applied based on cultural compatibility",
		Timestamp:        time.Now(),
		Era:              "contact_evolution",
		Trigger:          fmt.Sprintf("contact_%s", contactType),
		CultureInfluence: source,
		Intensity:        pattern.AdaptationStrength, // Use adaptation strength as intensity
	}

	return change, nil
}

// GenerateRandomAdaptation generates a random linguistic adaptation based on the current language state.
func (lae *LinguisticAdaptationEngine) GenerateRandomAdaptation(
	lang *lang.Language,
	adaptationType string,
) (*LinguisticAdaptation, error) {
	adaptation := &LinguisticAdaptation{
		ID:          fmt.Sprintf("random_adaptation_%d", time.Now().Unix()),
		Type:        adaptationType,
		Description: fmt.Sprintf("Random %s adaptation", adaptationType),
		Details:     "Generated automatically based on language state",
		Timestamp:   time.Now(),
		Source:      "internal_generation",
		ContactType: "none",
	}

	switch adaptationType {
	case "phonological":
		adaptation.PhonologicalChange = &PhonologicalAdaptation{
			ID:          fmt.Sprintf("phonological_%s", adaptation.ID),
			Type:        "phoneme_addition",
			Description: "Random phonological change",
			Probability: 0.3,
			Intensity:   0.5,
		}
	case "grammatical":
		adaptation.GrammaticalChange = &GrammaticalAdaptation{
			ID:          fmt.Sprintf("grammatical_%s", adaptation.ID),
			Type:        "morphology_addition",
			Description: "Random grammatical change",
			Probability: 0.2,
			Intensity:   0.4,
		}
	case "morphological":
		adaptation.MorphologicalChange = &MorphologicalAdaptation{
			ID:          fmt.Sprintf("morphological_%s", adaptation.ID),
			Type:        "morpheme_addition",
			Description: "Random morphological change",
			Probability: 0.25,
			Intensity:   0.45,
		}
	}

	return adaptation, nil
}
