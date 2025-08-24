package evolution

import (
	"fmt"
	"math/rand/v2"
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
	rng    *rand.Rand
	config EvolutionConfig
}

// NewContactEvolutionEngine creates a new contact evolution engine.
func NewContactEvolutionEngine(config EvolutionConfig) *ContactEvolutionEngine {
	rng := rand.New(rand.NewPCG(uint64(config.Seed), 0))

	return &ContactEvolutionEngine{
		rng:    rng,
		config: config,
	}
}

// SimulateContact simulates linguistic contact between two languages.
// Returns changes to the target language and a contact event record.
func (cee *ContactEvolutionEngine) SimulateContact(
	sourceLang *lang.Language,
	targetLang *lang.Language,
	contactType ContactType,
	intensity float32,
	duration time.Duration,
) ([]LinguisticChange, ContactEvent) {

	contactEvent := ContactEvent{
		ID:          fmt.Sprintf("contact_%s_%d", sourceLang.ID.String(), time.Now().Unix()),
		Timestamp:   time.Now(),
		SourceLang:  sourceLang.ID.String(),
		TargetLang:  targetLang.ID.String(),
		Type:        contactType.String(),
		Intensity:   intensity,
		Duration:    duration,
		Description: fmt.Sprintf("%s contact between %s and %s", contactType.String(), sourceLang.Name, targetLang.Name),
	}

	var changes []LinguisticChange

	// Determine what gets borrowed based on contact type and intensity
	if intensity >= cee.config.BorrowingThreshold {
		// Lexical borrowing (most common)
		if cee.rng.Float32() < 0.8 {
			lexicalChanges := cee.simulateLexicalBorrowing(sourceLang, targetLang, intensity, contactType)
			changes = append(changes, lexicalChanges...)
			contactEvent.LexicalBorrowing = true
		}

		// Phonological borrowing (less common)
		if cee.rng.Float32() < 0.4 && intensity > 0.6 {
			phonologicalChanges := cee.simulatePhonologicalBorrowing(sourceLang, targetLang, intensity)
			changes = append(changes, phonologicalChanges...)
			contactEvent.PhonologicalBorrowing = true
		}

		// Grammatical borrowing (least common)
		if cee.rng.Float32() < 0.2 && intensity > 0.8 {
			grammaticalChanges := cee.simulateGrammaticalBorrowing(sourceLang, targetLang, intensity)
			changes = append(changes, grammaticalChanges...)
			contactEvent.GrammaticalInfluence = true
		}
	}

	return changes, contactEvent
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
		if cee.rng.Float32() < 0.7 { // 70% chance per word
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
	if cee.rng.Float32() < intensity {
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
	if cee.rng.Float32() < intensity*0.5 {
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
