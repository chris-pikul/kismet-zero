package evolution

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/chris-pikul/kismet-zero/lang"
)

// EvolutionEngine is the main orchestrator for language evolution.
// It coordinates sound changes, morphological changes, orthographic changes, contact influence, and cultural influence.
type EvolutionEngine struct {
	soundChangeEngine       *SoundChangeEngine
	morphologyEngine        *MorphologicalEvolutionEngine
	orthographyEngine       *OrthographicEvolutionEngine
	contactEngine           *ContactEvolutionEngine
	culturalInfluenceEngine *CulturalInfluenceEngine
	familyTree              *LanguageFamilyTree
	config                  EvolutionConfig
	rng                     *rand.Rand
}

// NewEvolutionEngine creates a new evolution engine with the given configuration.
func NewEvolutionEngine(config EvolutionConfig) *EvolutionEngine {
	rng := rand.New(rand.NewPCG(uint64(config.Seed), 0))

	return &EvolutionEngine{
		soundChangeEngine:       NewSoundChangeEngine(config),
		morphologyEngine:        NewMorphologicalEvolutionEngine(config),
		orthographyEngine:       NewOrthographicEvolutionEngine(config),
		contactEngine:           NewContactEvolutionEngine(config),
		culturalInfluenceEngine: NewCulturalInfluenceEngine(config),
		familyTree:              NewLanguageFamilyTree(config),
		config:                  config,
		rng:                     rng,
	}
}

// EvolveLanguage evolves a language over a specified time period.
// Returns the evolved language and a summary of all changes.
func (ee *EvolutionEngine) EvolveLanguage(
	language *lang.Language,
	era string,
	duration time.Duration,
) (*lang.Language, EvolutionEvent, error) {

	if language == nil {
		return nil, EvolutionEvent{}, fmt.Errorf("cannot evolve nil language")
	}

	// Create a copy of the language for evolution
	evolvedLang := ee.cloneLanguage(language)

	// Track all changes
	allChanges := make([]LinguisticChange, 0)

	// Apply natural evolution
	naturalChanges := ee.applyNaturalEvolution(evolvedLang, era)
	allChanges = append(allChanges, naturalChanges...)

	// Apply sound changes
	soundChanges, modifiedPhonology := ee.soundChangeEngine.ApplySoundChanges(evolvedLang.Phonology, era)
	allChanges = append(allChanges, soundChanges...)
	if modifiedPhonology != nil {
		evolvedLang.SetPhonology(modifiedPhonology)
	}

	// Apply morphological changes
	morphologyChanges := ee.morphologyEngine.ApplyMorphologicalChanges(
		evolvedLang.Morphology,
		evolvedLang.Grammar,
		era,
	)
	// Convert morphological changes to linguistic changes
	for _, morphChange := range morphologyChanges {
		linguisticChange := LinguisticChange{
			ID:               morphChange.ID,
			Type:             ChangeTypeMorphological,
			Direction:        ChangeDirectionModifying,
			Description:      morphChange.Description,
			Details:          morphChange.Details,
			Timestamp:        morphChange.Timestamp,
			Era:              morphChange.Era,
			Trigger:          morphChange.Trigger,
			CultureInfluence: "",
			Intensity:        morphChange.Intensity,
		}
		allChanges = append(allChanges, linguisticChange)
	}

	// Apply orthographic changes
	orthographyChanges := ee.orthographyEngine.ApplyOrthographicChanges(
		evolvedLang.Orthography,
		era,
	)

	// Convert orthographic changes to linguistic changes
	for _, orthoChange := range orthographyChanges {
		linguisticChange := LinguisticChange{
			ID:               orthoChange.ID,
			Type:             ChangeTypeOrthographic,
			Direction:        ChangeDirectionModifying,
			Description:      orthoChange.Description,
			Details:          orthoChange.Details,
			Timestamp:        orthoChange.Timestamp,
			Era:              orthoChange.Era,
			Trigger:          orthoChange.Trigger,
			CultureInfluence: "",
			Intensity:        orthoChange.Intensity,
		}
		allChanges = append(allChanges, linguisticChange)
	}

	// Update language metadata
	evolvedLang.EvolvedAt = time.Now()
	evolvedLang.Seed = ee.rng.Int64() // New seed for future evolution

	// Create evolution event
	evolutionEvent := EvolutionEvent{
		ID:          fmt.Sprintf("evolution_%s_%d", language.ID.String(), time.Now().Unix()),
		Name:        fmt.Sprintf("Evolution of %s", language.Name),
		Description: fmt.Sprintf("Language evolved over %v during %s era", duration, era),
		Timestamp:   time.Now(),
		Era:         era,
		Changes:     allChanges,
		TriggerType: "natural_evolution",
		Intensity:   ee.calculateOverallIntensity(allChanges),
		Seed:        ee.rng.Int64(),
	}

	// Add to family tree if not already present
	if node, exists := ee.familyTree.GetLanguageNode(language.ID.String()); !exists || node == nil {
		err := ee.familyTree.AddLanguage(language, nil)
		if err != nil {
			return nil, EvolutionEvent{}, fmt.Errorf("failed to add language to family tree: %w", err)
		}
	}

	// Record the evolution event
	err := ee.familyTree.AddEvolutionEvent(language.ID.String(), evolutionEvent)
	if err != nil {
		return nil, EvolutionEvent{}, fmt.Errorf("failed to record evolution event: %w", err)
	}

	return evolvedLang, evolutionEvent, nil
}

// EvolveLanguageFromContact evolves a language through contact with another language.
// This simulates how languages influence each other through cultural interaction.
func (ee *EvolutionEngine) EvolveLanguageFromContact(
	targetLang *lang.Language,
	sourceLang *lang.Language,
	contactType ContactType,
	intensity float32,
	duration time.Duration,
) (*lang.Language, EvolutionEvent, error) {

	if targetLang == nil || sourceLang == nil {
		return nil, EvolutionEvent{}, fmt.Errorf("both languages must be provided")
	}

	// Create a copy of the target language
	evolvedLang := ee.cloneLanguage(targetLang)

	// Simulate contact and get changes
	contactChanges, contactEvent := ee.contactEngine.SimulateContact(
		sourceLang,
		evolvedLang,
		contactType,
		intensity,
		duration,
	)

	// Apply the contact changes
	// Note: In a full implementation, these changes would actually modify the language
	// For now, we just record them

	// Update language metadata
	evolvedLang.EvolvedAt = time.Now()
	evolvedLang.Seed = ee.rng.Int64()

	// Create evolution event
	evolutionEvent := EvolutionEvent{
		ID:             fmt.Sprintf("contact_evolution_%s_%d", targetLang.ID.String(), time.Now().Unix()),
		Name:           fmt.Sprintf("Contact Evolution of %s", targetLang.Name),
		Description:    fmt.Sprintf("Language evolved through %s contact with %s", contactType.String(), sourceLang.Name),
		Timestamp:      time.Now(),
		Era:            "contact_evolution",
		Changes:        contactChanges,
		TriggerType:    "contact",
		TriggerCulture: sourceLang.Culture,
		Intensity:      intensity,
		Seed:           ee.rng.Int64(),
	}

	// Add both languages to family tree if not present
	if node, exists := ee.familyTree.GetLanguageNode(targetLang.ID.String()); !exists || node == nil {
		err := ee.familyTree.AddLanguage(targetLang, nil)
		if err != nil {
			return nil, EvolutionEvent{}, fmt.Errorf("failed to add target language to family tree: %w", err)
		}
	}

	if node, exists := ee.familyTree.GetLanguageNode(sourceLang.ID.String()); !exists || node == nil {
		err := ee.familyTree.AddLanguage(sourceLang, nil)
		if err != nil {
			return nil, EvolutionEvent{}, fmt.Errorf("failed to add source language to family tree: %w", err)
		}
	}

	// Record events
	err := ee.familyTree.AddEvolutionEvent(targetLang.ID.String(), evolutionEvent)
	if err != nil {
		return nil, EvolutionEvent{}, fmt.Errorf("failed to record evolution event: %w", err)
	}

	err = ee.familyTree.AddContactEvent(targetLang.ID.String(), contactEvent)
	if err != nil {
		return nil, EvolutionEvent{}, fmt.Errorf("failed to record contact event: %w", err)
	}

	return evolvedLang, evolutionEvent, nil
}

// CreateChildLanguage creates a new language that has evolved from a parent language.
// This simulates language divergence and the creation of new dialects/languages.
func (ee *EvolutionEngine) CreateChildLanguage(
	parentLang *lang.Language,
	childID lang.LanguageID,
	childName string,
	era string,
) (*lang.Language, EvolutionEvent, error) {

	if parentLang == nil {
		return nil, EvolutionEvent{}, fmt.Errorf("parent language cannot be nil")
	}

	// Create the child language
	childLang := lang.NewLanguage(childID, childName, parentLang.Type, ee.rng.Int64())

	// Copy linguistic components from parent
	childLang.Culture = parentLang.Culture
	childLang.Description = fmt.Sprintf("Evolved from %s", parentLang.Name)
	childLang.ParentID = &parentLang.ID

	// Set parent's evolved components
	if parentLang.Phonology != nil {
		childLang.SetPhonology(parentLang.Phonology)
	}
	if parentLang.Orthography != nil {
		childLang.SetOrthography(parentLang.Orthography)
	}
	if parentLang.Morphology != nil {
		childLang.Morphology = parentLang.Morphology
	}
	if parentLang.Grammar != nil {
		childLang.Grammar = parentLang.Grammar
	}

	// Apply some initial evolution to differentiate from parent
	evolvedChild, evolutionEvent, err := ee.EvolveLanguage(childLang, era, time.Hour*24*365*50) // 50 years
	if err != nil {
		return nil, EvolutionEvent{}, fmt.Errorf("failed to evolve child language: %w", err)
	}

	// Use the evolved version
	childLang = evolvedChild

	// Add to family tree with parent relationship
	err = ee.familyTree.AddLanguage(childLang, parentLang)
	if err != nil {
		return nil, EvolutionEvent{}, fmt.Errorf("failed to add child language to family tree: %w", err)
	}

	// Update parent's child list
	parentLang.ChildIDs = append(parentLang.ChildIDs, childID)

	return childLang, evolutionEvent, nil
}

// applyNaturalEvolution applies natural, internal changes to a language.
func (ee *EvolutionEngine) applyNaturalEvolution(language *lang.Language, era string) []LinguisticChange {
	var changes []LinguisticChange

	// Apply natural changes based on configuration
	if ee.rng.Float32() < ee.config.NaturalChangeRate {
		change := &LinguisticChange{
			ID:          fmt.Sprintf("natural_change_%d", time.Now().UnixNano()),
			Type:        ChangeTypeMorphological,
			Direction:   ChangeDirectionModifying,
			Description: "Natural internal evolution",
			Details:     "Gradual internal changes over time",
			Timestamp:   time.Now(),
			Era:         era,
			Trigger:     "natural_evolution",
			Intensity:   0.3,
		}
		changes = append(changes, *change)
	}

	return changes
}

// cloneLanguage creates a deep copy of a language for evolution.
func (ee *EvolutionEngine) cloneLanguage(language *lang.Language) *lang.Language {
	// This is a simplified clone - in practice, you'd need to deep copy
	// all linguistic components
	clone := *language
	clone.EvolvedAt = time.Time{} // Reset evolution timestamp
	return &clone
}

// calculateOverallIntensity calculates the overall intensity of all changes.
func (ee *EvolutionEngine) calculateOverallIntensity(changes []LinguisticChange) float32 {
	if len(changes) == 0 {
		return 0.0
	}

	totalIntensity := float32(0.0)
	for _, change := range changes {
		totalIntensity += change.Intensity
	}

	return totalIntensity / float32(len(changes))
}

// GetFamilyTree returns the language family tree.
func (ee *EvolutionEngine) GetFamilyTree() *LanguageFamilyTree {
	return ee.familyTree
}

// GetEvolutionHistory returns the evolution history of a language.
func (ee *EvolutionEngine) GetEvolutionHistory(languageID string) ([]EvolutionEvent, error) {
	return ee.familyTree.GetEvolutionHistory(languageID)
}

// GetContactHistory returns the contact history of a language.
func (ee *EvolutionEngine) GetContactHistory(languageID string) ([]ContactEvent, error) {
	return ee.familyTree.GetContactHistory(languageID)
}

// SimulateCulturalInfluence simulates cultural influence between two languages.
// This method provides a high-level interface to the Cultural Influence Engine.
func (ee *EvolutionEngine) SimulateCulturalInfluence(
	sourceLang *lang.Language,
	targetLang *lang.Language,
	contactType CulturalContactType,
	duration time.Duration,
	sourceIdentity CulturalIdentity,
	targetIdentity CulturalIdentity,
) ([]LinguisticChange, CulturalInfluenceEvent, error) {

	if sourceLang == nil || targetLang == nil {
		return nil, CulturalInfluenceEvent{}, fmt.Errorf("both source and target languages must be provided")
	}

	// Use the cultural influence engine to simulate the influence
	changes, event := ee.culturalInfluenceEngine.SimulateCulturalInfluence(
		sourceLang,
		targetLang,
		contactType,
		duration,
		sourceIdentity,
		targetIdentity,
	)

	// Record the cultural influence event in the family tree
	if len(changes) > 0 {
		// Create a contact event for the family tree
		contactEvent := ContactEvent{
			ID:                    event.ID,
			Timestamp:             event.Timestamp,
			SourceLang:            event.SourceCulture,
			TargetLang:            event.TargetCulture,
			Type:                  event.ContactType.String(),
			Intensity:             event.Intensity,
			Duration:              event.Duration,
			Description:           event.Description,
			LexicalBorrowing:      event.LexicalInfluence,
			PhonologicalBorrowing: event.PhonologicalInfluence,
			GrammaticalInfluence:  event.GrammaticalInfluence,
			OrthographicBorrowing: event.OrthographicInfluence,
		}

		// Add to the family tree
		err := ee.familyTree.AddContactEvent(targetLang.ID.String(), contactEvent)
		if err != nil {
			return changes, event, fmt.Errorf("failed to record cultural influence in family tree: %w", err)
		}
	}

	return changes, event, nil
}

// GetCulturalInfluenceEngine returns the cultural influence engine for direct access.
func (ee *EvolutionEngine) GetCulturalInfluenceEngine() *CulturalInfluenceEngine {
	return ee.culturalInfluenceEngine
}
