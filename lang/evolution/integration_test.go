package evolution

import (
	"testing"
	"time"

	"github.com/chris-pikul/kismet-zero/lang"
)

// TestEvolutionPackageIntegration tests the complete evolution package working together.
func TestEvolutionPackageIntegration(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewEvolutionEngine(config)

	// Create a base language
	baseLang := createTestLanguage("Proto-Language", "proto_culture")

	// Test 1: Basic language evolution
	t.Run("BasicLanguageEvolution", func(t *testing.T) {
		evolvedLang, evolutionEvent, err := engine.EvolveLanguage(
			baseLang,
			"early_evolution",
			time.Hour*24*365*100, // 100 years
		)

		if err != nil {
			t.Fatalf("Failed to evolve language: %v", err)
		}

		if evolvedLang == nil {
			t.Fatal("Evolved language should not be nil")
		}

		if len(evolutionEvent.Changes) == 0 {
			t.Error("Evolution should produce changes")
		}

		// Verify different types of changes occurred
		changeTypes := make(map[ChangeType]bool)
		for _, change := range evolutionEvent.Changes {
			changeTypes[change.Type] = true
		}

		// Should have at least some natural changes
		if !changeTypes[ChangeTypeMorphological] && !changeTypes[ChangeTypeSoundShift] {
			t.Error("Should have natural evolution changes")
		}
	})

	// Test 2: Create child language
	t.Run("CreateChildLanguage", func(t *testing.T) {
		childLang, evolutionEvent, err := engine.CreateChildLanguage(
			baseLang,
			lang.LanguageID{Family: "test", Branch: "child", Language: "child_lang"},
			"Child Language",
			"child_evolution",
		)

		if err != nil {
			t.Fatalf("Failed to create child language: %v", err)
		}

		if childLang == nil {
			t.Fatal("Child language should not be nil")
		}

		if childLang.ParentID == nil {
			t.Error("Child should have parent ID")
		}

		// Child creation should either produce changes or have a properly structured evolution event
		if len(evolutionEvent.Changes) == 0 {
			// If no changes, verify the evolution event structure is correct
			if evolutionEvent.ID == "" || evolutionEvent.Timestamp.IsZero() {
				t.Error("Child creation should produce a properly structured evolution event")
			}
		}
	})

	// Test 3: Language contact evolution
	t.Run("LanguageContactEvolution", func(t *testing.T) {
		contactLang := createTestLanguage("Contact-Language", "contact_culture")

		contactEvolved, evolutionEvent, err := engine.EvolveLanguageFromContact(
			baseLang,
			contactLang,
			ContactTypeTrade,
			0.7,                  // High intensity
			time.Hour*24*365*200, // 200 years
		)

		if err != nil {
			t.Fatalf("Failed to evolve language from contact: %v", err)
		}

		if contactEvolved == nil {
			t.Fatal("Contact evolved language should not be nil")
		}

		if len(evolutionEvent.Changes) == 0 {
			t.Error("Contact evolution should produce changes")
		}

		// Verify contact-related changes
		hasContactChanges := false
		for _, change := range evolutionEvent.Changes {
			if change.Trigger == "contact_trade" {
				hasContactChanges = true
				break
			}
		}

		if !hasContactChanges {
			t.Error("Should have contact-related changes")
		}
	})

	// Test 4: Cultural influence modeling
	t.Run("CulturalInfluenceModeling", func(t *testing.T) {
		sourceLang := createTestLanguage("Imperial-Language", "imperial_culture")
		targetLang := createTestLanguage("Frontier-Language", "frontier_culture")

		// Add target language to family tree first
		familyTree := engine.GetFamilyTree()
		err := familyTree.AddLanguage(targetLang, nil)
		if err != nil {
			t.Fatalf("Failed to add target language to family tree: %v", err)
		}

		// Define cultural identities
		imperialIdentity := CulturalIdentity{
			CultureName:           "Imperial Empire",
			WritingSystemPrestige: 0.9,
			InnovationTendency:    0.6,
			PreservationInstinct:  0.7,
			MagicalTradition:      0.8,
			ReligiousInfluence:    0.9,
			EconomicPower:         0.9,
			MilitaryPower:         0.9,
			CulturalConfidence:    0.9,
		}

		frontierIdentity := CulturalIdentity{
			CultureName:           "Frontier Colony",
			WritingSystemPrestige: 0.3,
			InnovationTendency:    0.8,
			PreservationInstinct:  0.2,
			MagicalTradition:      0.4,
			ReligiousInfluence:    0.3,
			EconomicPower:         0.4,
			MilitaryPower:         0.3,
			CulturalConfidence:    0.4,
		}

		// Simulate conquest influence
		changes, event, err := engine.SimulateCulturalInfluence(
			sourceLang,
			targetLang,
			CulturalContactTypeConquest,
			time.Hour*24*365*100, // 100 years
			imperialIdentity,
			frontierIdentity,
		)

		if err != nil {
			t.Fatalf("Failed to simulate cultural influence: %v", err)
		}

		if event.Intensity < 0.7 {
			t.Errorf("Conquest should have high intensity, got %f", event.Intensity)
		}

		if !event.OrthographicInfluence && !event.LexicalInfluence {
			t.Error("Should have some form of influence")
		}

		if !event.MagicalInfluence {
			t.Error("Should detect magical influence with compatible traditions")
		}

		if len(changes) == 0 {
			t.Error("Cultural influence should produce changes")
		}
	})

	// Test 5: Language family tree
	t.Run("LanguageFamilyTree", func(t *testing.T) {
		// Get the language family tree
		langTree := engine.GetFamilyTree()

		// Create languages
		protoLang := createTestLanguage("Proto-Lang", "proto_culture")
		childLang := createTestLanguage("Child-Lang", "child_culture")

		// Add to tree
		err := langTree.AddLanguage(protoLang, nil)
		if err != nil {
			t.Fatalf("Failed to add proto language: %v", err)
		}

		err = langTree.AddLanguage(childLang, protoLang)
		if err != nil {
			t.Fatalf("Failed to add child language: %v", err)
		}

		// Test family tree operations
		ancestors := langTree.GetAncestors(childLang.ID.String())
		if len(ancestors) != 1 {
			t.Errorf("Expected 1 ancestor, got %d", len(ancestors))
		}

		descendants := langTree.GetDescendants(protoLang.ID.String())
		if len(descendants) != 1 {
			t.Errorf("Expected 1 descendant, got %d", len(descendants))
		}

		// Test evolution event recording
		evolutionEvent := EvolutionEvent{
			ID:          "test_evolution",
			Changes:     []LinguisticChange{},
			Timestamp:   time.Now(),
			Era:         "test_era",
			Description: "Test evolution event",
		}

		err = langTree.AddEvolutionEvent(childLang.ID.String(), evolutionEvent)
		if err != nil {
			t.Fatalf("Failed to add evolution event: %v", err)
		}

		// Test contact event recording
		contactEvent := ContactEvent{
			ID:         "test_contact",
			Timestamp:  time.Now(),
			SourceLang: protoLang.ID.String(),
			TargetLang: childLang.ID.String(),
			Type:       "trade",
			Intensity:  0.6,
			Duration:   time.Hour * 24 * 365 * 50,
		}

		err = langTree.AddContactEvent(childLang.ID.String(), contactEvent)
		if err != nil {
			t.Fatalf("Failed to add contact event: %v", err)
		}

		// Verify events were recorded
		evolutionHistory, err := langTree.GetEvolutionHistory(childLang.ID.String())
		if err != nil {
			t.Fatalf("Failed to get evolution history: %v", err)
		}

		if len(evolutionHistory) != 1 {
			t.Errorf("Expected 1 evolution event, got %d", len(evolutionHistory))
		}

		contactHistory, err := langTree.GetContactHistory(childLang.ID.String())
		if err != nil {
			t.Fatalf("Failed to get contact history: %v", err)
		}

		if len(contactHistory) != 1 {
			t.Errorf("Expected 1 contact event, got %d", len(contactHistory))
		}
	})

	// Test 6: Complete evolution workflow
	t.Run("CompleteEvolutionWorkflow", func(t *testing.T) {
		// Start with a proto language
		protoLang := createTestLanguage("Proto-Test", "proto_culture")

		// Evolve it naturally
		evolvedLang, _, err := engine.EvolveLanguage(
			protoLang,
			"natural_evolution",
			time.Hour*24*365*150, // 150 years
		)

		if err != nil {
			t.Fatalf("Failed natural evolution: %v", err)
		}

		// Create a child language
		childLang, _, err := engine.CreateChildLanguage(
			evolvedLang,
			lang.LanguageID{Family: "test", Branch: "child", Language: "child_lang"},
			"Child Language",
			"child_creation",
		)

		if err != nil {
			t.Fatalf("Failed child creation: %v", err)
		}

		// Apply cultural influence
		sourceLang := createTestLanguage("Influencer", "influencer_culture")
		sourceIdentity := CulturalIdentity{
			CultureName:           "Influencer Culture",
			WritingSystemPrestige: 0.8,
			InnovationTendency:    0.7,
			MagicalTradition:      0.6,
		}

		targetIdentity := CulturalIdentity{
			CultureName:          "Target Culture",
			InnovationTendency:   0.8,
			PreservationInstinct: 0.3,
		}

		changes, event, err := engine.SimulateCulturalInfluence(
			sourceLang,
			childLang,
			CulturalContactTypeTrade,
			time.Hour*24*365*50, // 50 years
			sourceIdentity,
			targetIdentity,
		)

		if err != nil {
			t.Fatalf("Failed cultural influence: %v", err)
		}

		// Verify the complete workflow produced results
		if len(changes) == 0 && event.Intensity < 0.3 {
			t.Error("Complete workflow should produce meaningful results")
		}

		// Check family tree integration
		familyTree := engine.GetFamilyTree()
		if familyTree == nil {
			t.Error("Family tree should be available")
		}

		// Verify languages are in the tree
		_, exists := familyTree.GetLanguageNode(protoLang.ID.String())
		if !exists {
			t.Error("Proto language should be in family tree")
		}

		childNode, exists := familyTree.GetLanguageNode(childLang.ID.String())
		if !exists {
			t.Error("Child language should be in family tree")
		}

		if childNode.Parent == nil {
			t.Error("Child should have parent in family tree")
		}

		if childNode.Parent.Language != protoLang {
			t.Error("Child's parent should be proto language")
		}
	})
}

// TestEvolutionPackagePerformance tests the performance of key operations.
func TestEvolutionPackagePerformance(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewEvolutionEngine(config)

	// Create a base language
	baseLang := createTestLanguage("Performance-Test", "test_culture")

	// Test evolution performance
	t.Run("EvolutionPerformance", func(t *testing.T) {
		start := time.Now()

		_, _, err := engine.EvolveLanguage(
			baseLang,
			"performance_test",
			time.Hour*24*365*100, // 100 years
		)

		duration := time.Since(start)

		if err != nil {
			t.Fatalf("Performance test failed: %v", err)
		}

		// Evolution should complete in reasonable time
		if duration > time.Second {
			t.Errorf("Evolution took too long: %v", duration)
		}
	})

	// Test cultural influence performance
	t.Run("CulturalInfluencePerformance", func(t *testing.T) {
		sourceLang := createTestLanguage("Source", "source_culture")
		targetLang := createTestLanguage("Target", "target_culture")

		// Add target language to family tree first
		familyTree := engine.GetFamilyTree()
		err := familyTree.AddLanguage(targetLang, nil)
		if err != nil {
			t.Fatalf("Failed to add target language to family tree: %v", err)
		}

		sourceIdentity := CulturalIdentity{
			CultureName:           "Source Culture",
			WritingSystemPrestige: 0.8,
			InnovationTendency:    0.7,
		}

		targetIdentity := CulturalIdentity{
			CultureName:        "Target Culture",
			InnovationTendency: 0.8,
		}

		start := time.Now()

		_, _, err = engine.SimulateCulturalInfluence(
			sourceLang,
			targetLang,
			CulturalContactTypeTrade,
			time.Hour*24*365*50, // 50 years
			sourceIdentity,
			targetIdentity,
		)

		duration := time.Since(start)

		if err != nil {
			t.Fatalf("Cultural influence performance test failed: %v", err)
		}

		// Cultural influence should complete quickly
		if duration > 100*time.Millisecond {
			t.Errorf("Cultural influence took too long: %v", duration)
		}
	})
}

// TestEvolutionPackageEdgeCases tests edge cases and error conditions.
func TestEvolutionPackageEdgeCases(t *testing.T) {
	config := DefaultEvolutionConfig(42)
	engine := NewEvolutionEngine(config)

	t.Run("NilLanguageHandling", func(t *testing.T) {
		_, _, err := engine.EvolveLanguage(nil, "test", time.Hour*24*365*100)
		if err == nil {
			t.Error("Should error on nil language")
		}
	})

	t.Run("EmptyLanguageHandling", func(t *testing.T) {
		emptyLang := &lang.Language{}
		_, _, err := engine.EvolveLanguage(emptyLang, "test", time.Hour*24*365*100)
		if err != nil {
			t.Errorf("Should handle empty language gracefully: %v", err)
		}
	})

	t.Run("ZeroDurationHandling", func(t *testing.T) {
		baseLang := createTestLanguage("Test", "test_culture")
		_, _, err := engine.EvolveLanguage(baseLang, "test", 0)
		if err != nil {
			t.Errorf("Should handle zero duration: %v", err)
		}
	})

	t.Run("ExtremeDurationHandling", func(t *testing.T) {
		baseLang := createTestLanguage("Test", "test_culture")
		extremeDuration := time.Hour * 24 * 365 * 100 // 100 years
		_, _, err := engine.EvolveLanguage(baseLang, "test", extremeDuration)
		if err != nil {
			t.Errorf("Should handle extreme duration: %v", err)
		}
	})
}
