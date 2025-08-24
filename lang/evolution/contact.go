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
		baseInfluence = 0.4
	}

	// Adjust based on duration (longer contact = more influence)
	durationFactor := float32(duration.Hours()) / (24.0 * 365.0 * 100.0) // Normalize to 100 years
	if durationFactor > 1.0 {
		durationFactor = 1.0
	}

	// Adjust based on cultural compatibility
	culturalCompatibility := float32(0.5)                 // Default neutral value for now
	culturalFactor := 0.5 + (culturalCompatibility * 0.5) // 0.5 to 1.0 range

	finalInfluence := baseInfluence * durationFactor * culturalFactor

	// Ensure influence stays within reasonable bounds
	if finalInfluence > 1.0 {
		finalInfluence = 1.0
	}

	return finalInfluence
}
