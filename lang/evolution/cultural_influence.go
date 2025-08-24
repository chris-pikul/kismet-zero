package evolution

import (
	"fmt"
	"time"

	"github.com/chris-pikul/kismet-zero/lang"
)

// CulturalContactType represents the type of cultural contact between civilizations.
type CulturalContactType byte

const (
	CulturalContactTypeUnknown          CulturalContactType = iota
	CulturalContactTypeTrade                                // Commercial networks and trade routes
	CulturalContactTypeConquest                             // Military conquest and empire building
	CulturalContactTypeMigration                            // Population movements and diaspora
	CulturalContactTypeReligious                            // Religious influence and missionary work
	CulturalContactTypeDiplomatic                           // Formal diplomatic relations
	CulturalContactTypeAlliance                             // Military and cultural alliances
	CulturalContactTypeColonization                         // Settlement and colonization
	CulturalContactTypeCulturalExchange                     // Academic and artistic exchange
	CulturalContactTypeMagical                              // Magical influence and supernatural contact
	CulturalContactTypeAncient                              // Ancient civilization influence
)

var culturalContactTypeEnum = []string{
	"unknown",
	"trade",
	"conquest",
	"migration",
	"religious",
	"diplomatic",
	"alliance",
	"colonization",
	"cultural_exchange",
	"magical",
	"ancient",
}

// String returns the string representation of the CulturalContactType.
func (cct CulturalContactType) String() string {
	if cct > CulturalContactTypeAncient {
		return culturalContactTypeEnum[0]
	}
	return culturalContactTypeEnum[cct]
}

// CulturalCompatibility represents factors that influence cultural borrowing.
type CulturalCompatibility struct {
	LinguisticSimilarity    float32 `json:"linguisticSimilarity"`    // 0.0-1.0, how similar the languages are
	HistoricalRelationship  float32 `json:"historicalRelationship"`  // 0.0-1.0, previous contact history
	PowerBalance            float32 `json:"powerBalance"`            // -1.0 to 1.0, negative = source dominant, positive = target dominant
	GeographicProximity     float32 `json:"geographicProximity"`     // 0.0-1.0, physical distance and accessibility
	CulturalValues          float32 `json:"culturalValues"`          // 0.0-1.0, compatibility of cultural attitudes
	MagicalAffinity         float32 `json:"magicalAffinity"`         // 0.0-1.0, magical system compatibility
	ReligiousCompatibility  float32 `json:"religiousCompatibility"`  // 0.0-1.0, religious system compatibility
	EconomicInterdependence float32 `json:"economicInterdependence"` // 0.0-1.0, economic relationship strength
}

// CulturalIdentity represents a civilization's cultural characteristics.
type CulturalIdentity struct {
	CultureName           string  `json:"cultureName"`
	WritingSystemPrestige float32 `json:"writingSystemPrestige"` // 0.0-1.0, how prestigious their writing is
	InnovationTendency    float32 `json:"innovationTendency"`    // 0.0-1.0, how open to new ideas
	PreservationInstinct  float32 `json:"preservationInstinct"`  // 0.0-1.0, how protective of traditions
	MagicalTradition      float32 `json:"magicalTradition"`      // 0.0-1.0, strength of magical practices
	ReligiousInfluence    float32 `json:"religiousInfluence"`    // 0.0-1.0, religious system influence
	EconomicPower         float32 `json:"economicPower"`         // 0.0-1.0, economic strength and influence
	MilitaryPower         float32 `json:"militaryPower"`         // 0.0-1.0, military strength and influence
	CulturalConfidence    float32 `json:"culturalConfidence"`    // 0.0-1.0, confidence in cultural superiority
}

// BorrowingPattern represents how a culture approaches borrowing from others.
type BorrowingPattern struct {
	SelectiveAdoption     float32 `json:"selectiveAdoption"`     // 0.0-1.0, how selective in what to borrow
	AdaptationStrength    float32 `json:"adaptationStrength"`    // 0.0-1.0, how much borrowed elements are modified
	IntegrationDepth      float32 `json:"integrationDepth"`      // 0.0-1.0, how deeply borrowed elements are integrated
	ResistanceLevel       float32 `json:"resistanceLevel"`       // 0.0-1.0, resistance to cultural borrowing
	PrestigeSensitivity   float32 `json:"prestigeSensitivity"`   // 0.0-1.0, sensitivity to prestige-based borrowing
	HybridizationTendency float32 `json:"hybridizationTendency"` // 0.0-1.0, tendency to create hybrid systems
}

// CulturalInfluenceEvent represents a cultural influence event with detailed metadata.
type CulturalInfluenceEvent struct {
	ID            string              `json:"id"`
	Timestamp     time.Time           `json:"timestamp"`
	SourceCulture string              `json:"sourceCulture"`
	TargetCulture string              `json:"targetCulture"`
	ContactType   CulturalContactType `json:"contactType"`
	Duration      time.Duration       `json:"duration"`
	Intensity     float32             `json:"intensity"`

	// Cultural context
	Compatibility  CulturalCompatibility `json:"compatibility"`
	SourceIdentity CulturalIdentity      `json:"sourceIdentity"`
	TargetIdentity CulturalIdentity      `json:"targetIdentity"`

	// What was influenced
	OrthographicInfluence bool `json:"orthographicInfluence"`
	LexicalInfluence      bool `json:"lexicalInfluence"`
	PhonologicalInfluence bool `json:"phonologicalInfluence"`
	GrammaticalInfluence  bool `json:"grammaticalInfluence"`

	// Influence details
	BorrowedFeatures  []string `json:"borrowedFeatures"`
	AdaptationNotes   string   `json:"adaptationNotes,omitempty"`
	ResistanceFactors []string `json:"resistanceFactors,omitempty"`
	SuccessFactors    []string `json:"successFactors,omitempty"`

	// Fantasy elements
	MagicalInfluence   bool `json:"magicalInfluence"`
	ReligiousInfluence bool `json:"religiousInfluence"`
	AncientInfluence   bool `json:"ancientInfluence"`

	Description string `json:"description,omitempty"`
}

// CulturalInfluenceEngine manages sophisticated cultural influence and borrowing.
type CulturalInfluenceEngine struct {
	BaseEngine
}

// NewCulturalInfluenceEngine creates a new cultural influence engine.
func NewCulturalInfluenceEngine(config EvolutionConfig) *CulturalInfluenceEngine {
	return &CulturalInfluenceEngine{
		BaseEngine: NewBaseEngine(config),
	}
}

// SimulateCulturalInfluence simulates sophisticated cultural influence between civilizations.
func (cie *CulturalInfluenceEngine) SimulateCulturalInfluence(
	sourceLang *lang.Language,
	targetLang *lang.Language,
	contactType CulturalContactType,
	duration time.Duration,
	sourceIdentity CulturalIdentity,
	targetIdentity CulturalIdentity,
) ([]LinguisticChange, CulturalInfluenceEvent) {

	// Calculate cultural compatibility
	compatibility := cie.calculateCulturalCompatibility(sourceLang, targetLang, sourceIdentity, targetIdentity)

	// Calculate base influence intensity
	baseIntensity := cie.calculateBaseInfluenceIntensity(contactType, duration, compatibility)

	// Apply cultural identity modifiers
	finalIntensity := cie.applyCulturalModifiers(baseIntensity, sourceIdentity, targetIdentity, compatibility)

	// Create the cultural influence event
	influenceEvent := CulturalInfluenceEvent{
		ID:             fmt.Sprintf("cultural_influence_%s_%d", contactType.String(), time.Now().Unix()),
		Timestamp:      time.Now(),
		SourceCulture:  sourceLang.Culture,
		TargetCulture:  targetLang.Culture,
		ContactType:    contactType,
		Duration:       duration,
		Intensity:      finalIntensity,
		Compatibility:  compatibility,
		SourceIdentity: sourceIdentity,
		TargetIdentity: targetIdentity,
	}

	var changes []LinguisticChange

	// Determine what gets influenced based on compatibility and intensity
	if finalIntensity >= cie.GetConfig().BorrowingThreshold {
		// Orthographic influence (writing system borrowing)
		if cie.shouldInfluenceOrthography(contactType, finalIntensity, compatibility) {
			orthoChanges := cie.simulateOrthographicInfluence(sourceLang, targetLang, finalIntensity, contactType, compatibility)
			changes = append(changes, orthoChanges...)
			influenceEvent.OrthographicInfluence = true
		}

		// Lexical influence (vocabulary borrowing)
		if cie.shouldInfluenceLexical(contactType, finalIntensity, compatibility) {
			lexicalChanges := cie.simulateLexicalInfluence(sourceLang, targetLang, finalIntensity, contactType, compatibility)
			changes = append(changes, lexicalChanges...)
			influenceEvent.LexicalInfluence = true
		}

		// Phonological influence (sound pattern borrowing)
		if cie.shouldInfluencePhonological(contactType, finalIntensity, compatibility) {
			phonologicalChanges := cie.simulatePhonologicalInfluence(sourceLang, targetLang, finalIntensity, compatibility)
			changes = append(changes, phonologicalChanges...)
			influenceEvent.PhonologicalInfluence = true
		}

		// Grammatical influence (grammar borrowing)
		if cie.shouldInfluenceGrammatical(contactType, finalIntensity, compatibility) {
			grammaticalChanges := cie.simulateGrammaticalInfluence(sourceLang, targetLang, finalIntensity, compatibility)
			changes = append(changes, grammaticalChanges...)
			influenceEvent.GrammaticalInfluence = true
		}

		// Fantasy-specific influences
		if sourceIdentity.MagicalTradition > 0.7 && targetIdentity.MagicalTradition > 0.3 {
			influenceEvent.MagicalInfluence = true
		}

		if sourceIdentity.ReligiousInfluence > 0.7 && targetIdentity.ReligiousInfluence > 0.3 {
			influenceEvent.ReligiousInfluence = true
		}

		if contactType == CulturalContactTypeAncient {
			influenceEvent.AncientInfluence = true
		}
	}

	// Update the influence event with details
	influenceEvent.BorrowedFeatures = cie.extractBorrowedFeatures(changes)
	influenceEvent.AdaptationNotes = cie.generateAdaptationNotes(changes, compatibility)
	influenceEvent.ResistanceFactors = cie.identifyResistanceFactors(targetIdentity, compatibility)
	influenceEvent.SuccessFactors = cie.identifySuccessFactors(sourceIdentity, compatibility)
	influenceEvent.Description = cie.generateInfluenceDescription(influenceEvent)

	return changes, influenceEvent
}

// calculateCulturalCompatibility calculates how compatible two cultures are for borrowing.
func (cie *CulturalInfluenceEngine) calculateCulturalCompatibility(
	sourceLang *lang.Language,
	targetLang *lang.Language,
	sourceIdentity CulturalIdentity,
	targetIdentity CulturalIdentity,
) CulturalCompatibility {

	// Linguistic similarity (simplified - in practice this would be more sophisticated)
	linguisticSimilarity := float32(0.5) // Base similarity

	// Historical relationship (random for now, could be enhanced with actual history)
	historicalRelationship := cie.GetRandomFloat32()

	// Power balance (based on cultural power indicators)
	sourcePower := (sourceIdentity.EconomicPower + sourceIdentity.MilitaryPower + sourceIdentity.CulturalConfidence) / 3.0
	targetPower := (targetIdentity.EconomicPower + targetIdentity.MilitaryPower + targetIdentity.CulturalConfidence) / 3.0
	powerBalance := sourcePower - targetPower // -1.0 to 1.0

	// Geographic proximity (random for now, could be enhanced with actual geography)
	geographicProximity := cie.GetRandomFloat32()

	// Cultural values compatibility
	culturalValues := 1.0 - abs(sourceIdentity.InnovationTendency-targetIdentity.InnovationTendency)

	// Magical affinity
	magicalAffinity := 1.0 - abs(sourceIdentity.MagicalTradition-targetIdentity.MagicalTradition)

	// Religious compatibility
	religiousCompatibility := 1.0 - abs(sourceIdentity.ReligiousInfluence-targetIdentity.ReligiousInfluence)

	// Economic interdependence (random for now)
	economicInterdependence := cie.GetRandomFloat32()

	return CulturalCompatibility{
		LinguisticSimilarity:    linguisticSimilarity,
		HistoricalRelationship:  historicalRelationship,
		PowerBalance:            powerBalance,
		GeographicProximity:     geographicProximity,
		CulturalValues:          culturalValues,
		MagicalAffinity:         magicalAffinity,
		ReligiousCompatibility:  religiousCompatibility,
		EconomicInterdependence: economicInterdependence,
	}
}

// calculateBaseInfluenceIntensity calculates the base intensity of cultural influence.
func (cie *CulturalInfluenceEngine) calculateBaseInfluenceIntensity(
	contactType CulturalContactType,
	duration time.Duration,
	compatibility CulturalCompatibility,
) float32 {

	// Base intensity by contact type
	var baseIntensity float32
	switch contactType {
	case CulturalContactTypeConquest:
		baseIntensity = 0.9
	case CulturalContactTypeColonization:
		baseIntensity = 0.8
	case CulturalContactTypeReligious:
		baseIntensity = 0.7
	case CulturalContactTypeTrade:
		baseIntensity = 0.6
	case CulturalContactTypeMigration:
		baseIntensity = 0.6
	case CulturalContactTypeAlliance:
		baseIntensity = 0.5
	case CulturalContactTypeDiplomatic:
		baseIntensity = 0.4
	case CulturalContactTypeCulturalExchange:
		baseIntensity = 0.4
	case CulturalContactTypeMagical:
		baseIntensity = 0.7
	case CulturalContactTypeAncient:
		baseIntensity = 0.8
	default:
		baseIntensity = 0.5
	}

	// Duration factor (longer contact = more influence)
	durationFactor := float32(duration.Hours()) / (24.0 * 365.0 * 100.0) // Normalize to 100 years
	if durationFactor > 1.0 {
		durationFactor = 1.0
	}

	// Compatibility factor
	compatibilityFactor := (compatibility.LinguisticSimilarity + compatibility.CulturalValues + compatibility.GeographicProximity) / 3.0

	// Calculate final intensity
	finalIntensity := baseIntensity * (0.6 + 0.4*compatibilityFactor) * (0.7 + 0.3*durationFactor)

	if finalIntensity > 1.0 {
		finalIntensity = 1.0
	}

	return finalIntensity
}

// applyCulturalModifiers applies cultural identity modifiers to influence intensity.
func (cie *CulturalInfluenceEngine) applyCulturalModifiers(
	baseIntensity float32,
	sourceIdentity CulturalIdentity,
	targetIdentity CulturalIdentity,
	compatibility CulturalCompatibility,
) float32 {

	modifiedIntensity := baseIntensity

	// Source writing system prestige
	if sourceIdentity.WritingSystemPrestige > 0.7 {
		modifiedIntensity *= 1.2 // Prestigious writing systems have more influence
	}

	// Target innovation tendency
	if targetIdentity.InnovationTendency > 0.7 {
		modifiedIntensity *= 1.1 // Innovative cultures are more open to borrowing
	}

	// Target preservation instinct
	if targetIdentity.PreservationInstinct > 0.7 {
		modifiedIntensity *= 0.8 // Traditional cultures resist borrowing
	}

	// Power balance effects
	if compatibility.PowerBalance < -0.5 {
		modifiedIntensity *= 1.3 // Dominant source culture has more influence
	} else if compatibility.PowerBalance > 0.5 {
		modifiedIntensity *= 0.7 // Dominant target culture resists influence
	}

	// Magical tradition compatibility
	if compatibility.MagicalAffinity > 0.7 {
		modifiedIntensity *= 1.2 // Magical compatibility increases influence
	}

	// Religious compatibility
	if compatibility.ReligiousCompatibility > 0.7 {
		modifiedIntensity *= 1.1 // Religious compatibility increases influence
	}

	if modifiedIntensity > 1.0 {
		modifiedIntensity = 1.0
	}

	return modifiedIntensity
}

// Helper functions for determining what gets influenced
func (cie *CulturalInfluenceEngine) shouldInfluenceOrthography(contactType CulturalContactType, intensity float32, compatibility CulturalCompatibility) bool {
	// Orthographic influence is most likely with high intensity and good compatibility
	return intensity > 0.7 && compatibility.LinguisticSimilarity > 0.6
}

func (cie *CulturalInfluenceEngine) shouldInfluenceLexical(contactType CulturalContactType, intensity float32, compatibility CulturalCompatibility) bool {
	// Lexical borrowing is common with moderate intensity
	return intensity > 0.5
}

func (cie *CulturalInfluenceEngine) shouldInfluencePhonological(contactType CulturalContactType, intensity float32, compatibility CulturalCompatibility) bool {
	// Phonological borrowing requires high intensity and good linguistic compatibility
	return intensity > 0.7 && compatibility.LinguisticSimilarity > 0.7
}

func (cie *CulturalInfluenceEngine) shouldInfluenceGrammatical(contactType CulturalContactType, intensity float32, compatibility CulturalCompatibility) bool {
	// Grammatical borrowing requires very high intensity and excellent compatibility
	return intensity > 0.8 && compatibility.LinguisticSimilarity > 0.8
}

// Helper function for absolute value
func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

// simulateOrthographicInfluence implements sophisticated orthographic borrowing.
func (cie *CulturalInfluenceEngine) simulateOrthographicInfluence(
	sourceLang *lang.Language,
	targetLang *lang.Language,
	intensity float32,
	contactType CulturalContactType,
	compatibility CulturalCompatibility,
) []LinguisticChange {

	var changes []LinguisticChange

	// Determine borrowing probability based on intensity and compatibility
	borrowingProb := intensity * compatibility.LinguisticSimilarity * 0.8

	if cie.GetRandomFloat32() < borrowingProb {
		// Create orthographic change
		change := LinguisticChange{
			ID:               fmt.Sprintf("orthographic_influence_%d", time.Now().UnixNano()),
			Type:             ChangeTypeOrthographic,
			Direction:        ChangeDirectionAdditive,
			Description:      fmt.Sprintf("Adopted orthographic features from %s", sourceLang.Culture),
			Details:          "Writing system features borrowed through cultural influence",
			Timestamp:        time.Now(),
			Era:              "cultural_influence",
			Trigger:          fmt.Sprintf("cultural_%s", contactType.String()),
			CultureInfluence: sourceLang.Culture,
			Intensity:        intensity * 0.7,
		}
		changes = append(changes, change)
	}

	// Additional changes based on contact type
	switch contactType {
	case CulturalContactTypeConquest, CulturalContactTypeColonization:
		if cie.GetRandomFloat32() < intensity*0.6 {
			// Script reform under pressure
			change := LinguisticChange{
				ID:               fmt.Sprintf("script_reform_%d", time.Now().UnixNano()),
				Type:             ChangeTypeOrthographic,
				Direction:        ChangeDirectionModifying,
				Description:      "Script reform under cultural pressure",
				Details:          "Writing system adapted to new cultural context",
				Timestamp:        time.Now(),
				Era:              "cultural_influence",
				Trigger:          "cultural_pressure",
				CultureInfluence: sourceLang.Culture,
				Intensity:        intensity * 0.8,
			}
			changes = append(changes, change)
		}
	}

	return changes
}

func (cie *CulturalInfluenceEngine) simulateLexicalInfluence(
	sourceLang *lang.Language,
	targetLang *lang.Language,
	intensity float32,
	contactType CulturalContactType,
	compatibility CulturalCompatibility,
) []LinguisticChange {

	var changes []LinguisticChange

	// Calculate number of words to borrow based on intensity and contact type
	var baseWordCount int
	switch contactType {
	case CulturalContactTypeConquest, CulturalContactTypeColonization:
		baseWordCount = int(intensity * 25) // 0-25 words
	case CulturalContactTypeTrade, CulturalContactTypeReligious:
		baseWordCount = int(intensity * 15) // 0-15 words
	case CulturalContactTypeMagical, CulturalContactTypeAncient:
		baseWordCount = int(intensity * 20) // 0-20 words
	default:
		baseWordCount = int(intensity * 10) // 0-10 words
	}

	// Apply cultural compatibility modifier
	compatibilityModifier := (compatibility.CulturalValues + compatibility.GeographicProximity) / 2.0
	adjustedWordCount := int(float32(baseWordCount) * compatibilityModifier)

	// Generate lexical changes
	for i := 0; i < adjustedWordCount; i++ {
		if cie.GetRandomFloat32() < 0.8 { // 80% chance per word

			// Determine word category based on contact type
			var wordCategory string
			switch contactType {
			case CulturalContactTypeTrade:
				wordCategory = "commercial_terms"
			case CulturalContactTypeReligious:
				wordCategory = "religious_vocabulary"
			case CulturalContactTypeMagical:
				wordCategory = "magical_terms"
			case CulturalContactTypeAncient:
				wordCategory = "ancient_knowledge"
			case CulturalContactTypeConquest:
				wordCategory = "military_terms"
			default:
				wordCategory = "general_vocabulary"
			}

			change := LinguisticChange{
				ID:               fmt.Sprintf("lexical_influence_%d_%d", time.Now().UnixNano(), i),
				Type:             ChangeTypeLexical,
				Direction:        ChangeDirectionAdditive,
				Description:      fmt.Sprintf("Borrowed %s from %s", wordCategory, sourceLang.Culture),
				Details:          fmt.Sprintf("Vocabulary item adapted from %s cultural context", sourceLang.Culture),
				Timestamp:        time.Now(),
				Era:              "cultural_influence",
				Trigger:          fmt.Sprintf("cultural_%s", contactType.String()),
				CultureInfluence: sourceLang.Culture,
				Intensity:        intensity * 0.6,
			}
			changes = append(changes, change)
		}
	}

	return changes
}

func (cie *CulturalInfluenceEngine) simulatePhonologicalInfluence(
	sourceLang *lang.Language,
	targetLang *lang.Language,
	intensity float32,
	compatibility CulturalCompatibility,
) []LinguisticChange {

	var changes []LinguisticChange

	// Phonological borrowing requires high linguistic compatibility
	if compatibility.LinguisticSimilarity < 0.6 {
		return changes // Too different to borrow sounds
	}

	// Calculate borrowing probability
	borrowingProb := intensity * compatibility.LinguisticSimilarity * 0.5

	if cie.GetRandomFloat32() < borrowingProb {
		// Create phonological change
		change := LinguisticChange{
			ID:               fmt.Sprintf("phonological_influence_%d", time.Now().UnixNano()),
			Type:             ChangeTypeSoundShift,
			Direction:        ChangeDirectionAdditive,
			Description:      fmt.Sprintf("Adopted phonological patterns from %s", sourceLang.Culture),
			Details:          "Sound patterns and phoneme inventory influenced by cultural contact",
			Timestamp:        time.Now(),
			Era:              "cultural_influence",
			Trigger:          "phonological_cultural_contact",
			CultureInfluence: sourceLang.Culture,
			Intensity:        intensity * 0.5,
		}
		changes = append(changes, change)
	}

	// Additional phonological changes for high-intensity contact
	if intensity > 0.8 && compatibility.LinguisticSimilarity > 0.7 {
		if cie.GetRandomFloat32() < 0.4 {
			// Phoneme inventory expansion
			change := LinguisticChange{
				ID:               fmt.Sprintf("phoneme_expansion_%d", time.Now().UnixNano()),
				Type:             ChangeTypeSoundShift,
				Direction:        ChangeDirectionAdditive,
				Description:      "Phoneme inventory expanded through cultural contact",
				Details:          "New sounds adopted to accommodate borrowed vocabulary",
				Timestamp:        time.Now(),
				Era:              "cultural_influence",
				Trigger:          "phoneme_expansion",
				CultureInfluence: sourceLang.Culture,
				Intensity:        intensity * 0.6,
			}
			changes = append(changes, change)
		}
	}

	return changes
}

func (cie *CulturalInfluenceEngine) simulateGrammaticalInfluence(
	sourceLang *lang.Language,
	targetLang *lang.Language,
	intensity float32,
	compatibility CulturalCompatibility,
) []LinguisticChange {

	var changes []LinguisticChange

	// Grammatical borrowing requires very high linguistic compatibility and intensity
	if compatibility.LinguisticSimilarity < 0.7 || intensity < 0.7 {
		return changes // Too different or too weak to borrow grammar
	}

	// Calculate borrowing probability (much lower than other types)
	borrowingProb := intensity * compatibility.LinguisticSimilarity * 0.3

	if cie.GetRandomFloat32() < borrowingProb {
		// Create grammatical change
		change := LinguisticChange{
			ID:               fmt.Sprintf("grammatical_influence_%d", time.Now().UnixNano()),
			Type:             ChangeTypeMorphological,
			Direction:        ChangeDirectionAdditive,
			Description:      fmt.Sprintf("Adopted grammatical patterns from %s", sourceLang.Culture),
			Details:          "Syntactic structures and morphological patterns influenced by cultural contact",
			Timestamp:        time.Now(),
			Era:              "cultural_influence",
			Trigger:          "grammatical_cultural_contact",
			CultureInfluence: sourceLang.Culture,
			Intensity:        intensity * 0.4,
		}
		changes = append(changes, change)
	}

	// Additional grammatical changes for very high-intensity contact
	if intensity > 0.9 && compatibility.LinguisticSimilarity > 0.8 {
		if cie.GetRandomFloat32() < 0.3 {
			// Syntax restructuring
			change := LinguisticChange{
				ID:               fmt.Sprintf("syntax_restructuring_%d", time.Now().UnixNano()),
				Type:             ChangeTypeMorphological,
				Direction:        ChangeDirectionModifying,
				Description:      "Syntax restructured under cultural influence",
				Details:          "Word order and sentence structure adapted to new cultural patterns",
				Timestamp:        time.Now(),
				Era:              "cultural_influence",
				Trigger:          "syntax_restructuring",
				CultureInfluence: sourceLang.Culture,
				Intensity:        intensity * 0.5,
			}
			changes = append(changes, change)
		}
	}

	return changes
}

// Helper functions for the influence event
func (cie *CulturalInfluenceEngine) extractBorrowedFeatures(changes []LinguisticChange) []string {
	features := make([]string, 0)
	for _, change := range changes {
		features = append(features, change.Type.String())
	}
	return features
}

func (cie *CulturalInfluenceEngine) generateAdaptationNotes(changes []LinguisticChange, compatibility CulturalCompatibility) string {
	if len(changes) == 0 {
		return "No borrowing occurred"
	}
	return fmt.Sprintf("Adapted %d features with compatibility factor %.2f", len(changes), compatibility.CulturalValues)
}

func (cie *CulturalInfluenceEngine) identifyResistanceFactors(targetIdentity CulturalIdentity, compatibility CulturalCompatibility) []string {
	factors := make([]string, 0)

	if targetIdentity.PreservationInstinct > 0.7 {
		factors = append(factors, "high_traditional_preservation")
	}

	if compatibility.PowerBalance > 0.5 {
		factors = append(factors, "target_culture_dominance")
	}

	if targetIdentity.CulturalConfidence > 0.8 {
		factors = append(factors, "cultural_superiority_complex")
	}

	return factors
}

func (cie *CulturalInfluenceEngine) identifySuccessFactors(sourceIdentity CulturalIdentity, compatibility CulturalCompatibility) []string {
	factors := make([]string, 0)

	if sourceIdentity.WritingSystemPrestige > 0.7 {
		factors = append(factors, "prestigious_writing_system")
	}

	if compatibility.LinguisticSimilarity > 0.7 {
		factors = append(factors, "high_linguistic_compatibility")
	}

	if compatibility.GeographicProximity > 0.7 {
		factors = append(factors, "geographic_accessibility")
	}

	return factors
}

func (cie *CulturalInfluenceEngine) generateInfluenceDescription(event CulturalInfluenceEvent) string {
	return fmt.Sprintf("%s influence from %s to %s with intensity %.2f",
		event.ContactType.String(), event.SourceCulture, event.TargetCulture, event.Intensity)
}
