package evolution

import (
	"fmt"
	"time"

	"github.com/chris-pikul/kismet-zero/lang"
)

// ContactType represents the type of linguistic contact between languages.
type ContactType byte

const (
	ContactTypeUnknown     ContactType = iota
	ContactTypeTrade                   // Commercial contact
	ContactTypeConquest                // Military conquest
	ContactTypeMigration               // Population movement
	ContactTypeCultural                // Cultural exchange
	ContactTypeReligious               // Religious influence
	ContactTypeEducational             // Educational/academic contact
)

var contactTypeEnum = []string{
	"unknown",
	"trade",
	"conquest",
	"migration",
	"cultural",
	"religious",
	"educational",
}

// String returns the string representation of the ContactType.
func (ct ContactType) String() string {
	if ct > ContactTypeEducational {
		return contactTypeEnum[0]
	}
	return contactTypeEnum[ct]
}

// BorrowingType represents what linguistic features are borrowed.
type BorrowingType byte

const (
	BorrowingTypeUnknown      BorrowingType = iota
	BorrowingTypeLexical                    // Words and vocabulary
	BorrowingTypePhonological               // Sound patterns
	BorrowingTypeGrammatical                // Grammar structures
	BorrowingTypeOrthographic               // Writing system features
)

var borrowingTypeEnum = []string{
	"unknown",
	"lexical",
	"phonological",
	"grammatical",
	"orthographic",
}

// String returns the string representation of the BorrowingType.
func (bt BorrowingType) String() string {
	if bt > BorrowingTypeOrthographic {
		return borrowingTypeEnum[0]
	}
	return borrowingTypeEnum[bt]
}

// ContactEvolutionEngine manages contact-induced language changes.
type ContactEvolutionEngine struct {
	BaseEngine
}

// NewContactEvolutionEngine creates a new contact evolution engine.
func NewContactEvolutionEngine(config EvolutionConfig) *ContactEvolutionEngine {
	return &ContactEvolutionEngine{
		BaseEngine: NewBaseEngine(config),
	}
}

// SimulateContactWithAdaptation simulates linguistic contact using the enhanced adaptation system.
// This method provides access to the sophisticated adaptation tracking and borrowing patterns.
func (cee *ContactEvolutionEngine) SimulateContactWithAdaptation(
	sourceLang *lang.Language,
	targetLang *lang.Language,
	contactType ContactType,
	intensity float32,
	duration time.Duration,
	adaptationEngine *LinguisticAdaptationEngine,
) ([]LinguisticChange, ContactEvent, error) {
	if adaptationEngine == nil {
		return nil, ContactEvent{}, fmt.Errorf("adaptation engine cannot be nil")
	}

	// Create the contact event
	contactEvent := ContactEvent{
		ID:          fmt.Sprintf("contact_adaptation_%s_%d", sourceLang.ID.String(), time.Now().Unix()),
		Timestamp:   time.Now(),
		SourceLang:  sourceLang.ID.String(),
		TargetLang:  targetLang.ID.String(),
		Type:        contactType.String(),
		Intensity:   intensity,
		Duration:    duration,
		Description: fmt.Sprintf("%s contact with adaptation tracking between %s and %s", contactType.String(), sourceLang.Name, targetLang.Name),
	}

	var changes []LinguisticChange

	// Determine what gets borrowed based on contact type and intensity
	if intensity >= cee.GetConfig().BorrowingThreshold {
		// Create borrowing patterns based on contact type and intensity
		borrowingPattern := cee.createBorrowingPattern(contactType, intensity)

		// Apply the borrowing pattern using the adaptation engine
		change, err := adaptationEngine.ApplyBorrowingPattern(targetLang, borrowingPattern, sourceLang.Name, contactType.String())
		if err != nil {
			return nil, ContactEvent{}, fmt.Errorf("failed to apply borrowing pattern: %w", err)
		}

		if change != nil {
			changes = append(changes, *change)

			// Update contact event based on the type of borrowing
			switch change.Type {
			case ChangeTypeLexical:
				contactEvent.LexicalBorrowing = true
			case ChangeTypeSoundShift:
				contactEvent.PhonologicalBorrowing = true
			case ChangeTypeMorphological:
				contactEvent.GrammaticalInfluence = true
			case ChangeTypeOrthographic:
				contactEvent.OrthographicBorrowing = true
			}
		}

		// Generate additional adaptations based on contact intensity
		if intensity > 0.7 {
			adaptation, err := adaptationEngine.GenerateRandomAdaptation(targetLang, "phonological")
			if err == nil && adaptation != nil {
				change, err := adaptationEngine.ApplyPhonologicalAdaptation(targetLang, adaptation.PhonologicalChange)
				if err == nil && change != nil {
					changes = append(changes, *change)
				}
			}
		}
	}

	return changes, contactEvent, nil
}

// createBorrowingPattern creates a borrowing pattern based on contact type and intensity.
func (cee *ContactEvolutionEngine) createBorrowingPattern(contactType ContactType, intensity float32) *BorrowingPattern {
	// Create a borrowing pattern with characteristics based on contact type
	pattern := &BorrowingPattern{
		SelectiveAdoption:     intensity * 0.8, // Higher intensity = more selective
		AdaptationStrength:    intensity * 0.9, // Higher intensity = stronger adaptation
		IntegrationDepth:      intensity * 0.7, // Higher intensity = deeper integration
		ResistanceLevel:       1.0 - intensity, // Higher intensity = lower resistance
		PrestigeSensitivity:   intensity * 0.6, // Higher intensity = more prestige sensitive
		HybridizationTendency: intensity * 0.5, // Higher intensity = more hybridization
	}

	return pattern
}

// simulateLexicalBorrowing simulates the borrowing of words from one language to another.
func (cee *ContactEvolutionEngine) simulateLexicalBorrowing(
	sourceLang *lang.Language,
	targetLang *lang.Language,
	intensity float32,
	contactType ContactType,
) []LinguisticChange {

	var changes []LinguisticChange

	// Number of words to borrow based on intensity
	numWords := int(intensity * 20) // 0-20 words based on intensity

	for i := 0; i < numWords; i++ {
		if cee.GetRandomFloat32() < 0.7 { // 70% chance per word
			change := &LinguisticChange{
				ID:               fmt.Sprintf("lexical_borrowing_%d", time.Now().UnixNano()),
				Type:             ChangeTypeLexical,
				Direction:        ChangeDirectionAdditive,
				Description:      fmt.Sprintf("Borrowed word from %s", sourceLang.Name),
				Details:          "New vocabulary item adapted from source language",
				Timestamp:        time.Now(),
				Era:              "contact_evolution",
				Trigger:          fmt.Sprintf("contact_%s", contactType.String()),
				CultureInfluence: sourceLang.Culture,
				Intensity:        intensity * 0.8,
			}
			changes = append(changes, *change)
		}
	}

	return changes
}

// simulatePhonologicalBorrowing simulates the adoption of sound patterns from another language.
func (cee *ContactEvolutionEngine) simulatePhonologicalBorrowing(
	sourceLang *lang.Language,
	targetLang *lang.Language,
	intensity float32,
) []LinguisticChange {

	var changes []LinguisticChange

	// Phonological borrowing is more selective
	if cee.GetRandomFloat32() < intensity {
		change := &LinguisticChange{
			ID:               fmt.Sprintf("phonological_borrowing_%d", time.Now().UnixNano()),
			Type:             ChangeTypeSoundShift,
			Direction:        ChangeDirectionAdditive,
			Description:      fmt.Sprintf("Adopted phonological features from %s", sourceLang.Name),
			Details:          "New sound patterns or phoneme inventory changes",
			Timestamp:        time.Now(),
			Era:              "contact_evolution",
			Trigger:          "phonological_contact",
			CultureInfluence: sourceLang.Culture,
			Intensity:        intensity * 0.6,
		}
		changes = append(changes, *change)
	}

	return changes
}

// simulateGrammaticalBorrowing simulates the adoption of grammatical structures from another language.
func (cee *ContactEvolutionEngine) simulateGrammaticalBorrowing(
	sourceLang *lang.Language,
	targetLang *lang.Language,
	intensity float32,
) []LinguisticChange {

	var changes []LinguisticChange

	// Grammatical borrowing is the rarest and most significant
	if cee.GetRandomFloat32() < intensity*0.5 {
		change := &LinguisticChange{
			ID:               fmt.Sprintf("grammatical_borrowing_%d", time.Now().UnixNano()),
			Type:             ChangeTypeMorphological,
			Direction:        ChangeDirectionAdditive,
			Description:      fmt.Sprintf("Adopted grammatical features from %s", sourceLang.Name),
			Details:          "New grammatical structures or patterns",
			Timestamp:        time.Now(),
			Era:              "contact_evolution",
			Trigger:          "grammatical_contact",
			CultureInfluence: sourceLang.Culture,
			Intensity:        intensity * 0.4,
		}
		changes = append(changes, *change)
	}

	return changes
}

// CalculateContactInfluence calculates how much influence one language has on another.
func (cee *ContactEvolutionEngine) CalculateContactInfluence(
	sourceLang *lang.Language,
	targetLang *lang.Language,
	contactType ContactType,
	duration time.Duration,
) float32 {

	baseInfluence := float32(0.3)

	// Adjust based on contact type
	switch contactType {
	case ContactTypeConquest:
		baseInfluence = 0.8
	case ContactTypeMigration:
		baseInfluence = 0.6
	case ContactTypeTrade:
		baseInfluence = 0.4
	case ContactTypeCultural:
		baseInfluence = 0.3
	case ContactTypeReligious:
		baseInfluence = 0.5
	case ContactTypeEducational:
		baseInfluence = 0.2
	}

	// Adjust based on duration (longer contact = more influence)
	durationFactor := float32(duration.Hours() / (24 * 365 * 10)) // 10 years as baseline
	if durationFactor > 1.0 {
		durationFactor = 1.0
	}

	// Adjust based on cultural similarity (placeholder for future implementation)
	culturalSimilarity := float32(0.5) // Default neutral value

	return baseInfluence * (0.7 + 0.3*durationFactor) * (0.8 + 0.2*culturalSimilarity)
}

// simulateOrthographicBorrowing simulates the borrowing of writing system features from one language to another.
func (cee *ContactEvolutionEngine) simulateOrthographicBorrowing(
	sourceLang *lang.Language,
	targetLang *lang.Language,
	intensity float32,
	contactType ContactType,
) []LinguisticChange {

	var changes []LinguisticChange

	// Orthographic borrowing is very selective and depends on writing system compatibility
	if cee.GetRandomFloat32() < intensity*0.3 {
		change := &LinguisticChange{
			ID:               fmt.Sprintf("orthographic_borrowing_%d", time.Now().UnixNano()),
			Type:             ChangeTypeOrthographic,
			Direction:        ChangeDirectionAdditive,
			Description:      fmt.Sprintf("Adopted orthographic features from %s", sourceLang.Name),
			Details:          "New writing system features or conventions",
			Timestamp:        time.Now(),
			Era:              "contact_evolution",
			Trigger:          "orthographic_contact",
			CultureInfluence: sourceLang.Culture,
			Intensity:        intensity * 0.4,
		}
		changes = append(changes, *change)
	}

	return changes
}
