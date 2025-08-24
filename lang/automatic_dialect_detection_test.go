package lang

import (
	"fmt"
	"testing"
	"time"
)

// TestDialectDetectionEngine tests the dialect detection engine system.
func TestDialectDetectionEngine(t *testing.T) {
	engine := NewDialectDetectionEngine()

	t.Run("Engine Configuration", func(t *testing.T) {
		if engine.DivergenceThreshold != 0.3 {
			t.Errorf("Expected divergence threshold 0.3, got %f", engine.DivergenceThreshold)
		}
		if engine.SimilarityThreshold != 0.7 {
			t.Errorf("Expected similarity threshold 0.7, got %f", engine.SimilarityThreshold)
		}
		if engine.MinChangeCount != 5 {
			t.Errorf("Expected min change count 5, got %d", engine.MinChangeCount)
		}
		if engine.MaxDialectDistance != 0.5 {
			t.Errorf("Expected max dialect distance 0.5, got %f", engine.MaxDialectDistance)
		}

		t.Logf("Dialect detection engine configured with: divergence=%f, similarity=%f, min_changes=%d, max_distance=%f",
			engine.DivergenceThreshold, engine.SimilarityThreshold, engine.MinChangeCount, engine.MaxDialectDistance)
	})
}

// TestDialectFormationOpportunityDetection tests the detection of dialect formation opportunities.
func TestDialectFormationOpportunityDetection(t *testing.T) {
	engine := NewDialectDetectionEngine()

	// Create a base language
	baseLang, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create base language: %v", err)
	}

	t.Run("No Changes - No Opportunities", func(t *testing.T) {
		opportunities := engine.DetectDialectFormationOpportunities(baseLang)
		if len(opportunities) != 0 {
			t.Errorf("Expected no opportunities for language with no changes, got %d", len(opportunities))
		}
	})

	t.Run("Insufficient Changes - No Opportunities", func(t *testing.T) {
		// Add a few changes but not enough to meet threshold
		for i := 0; i < 3; i++ {
			change := LinguisticChange{
				ID:          fmt.Sprintf("test_change_%d", i),
				Type:        LinguisticChangeTypeSound,
				Description: fmt.Sprintf("Test change %d", i),
				Details:     fmt.Sprintf("Test change details %d", i),
				Timestamp:   time.Now(),
				Era:         "test",
				Trigger:     "test",
				Intensity:   0.5,
			}
			baseLang.AddLinguisticChange(change)
		}

		opportunities := engine.DetectDialectFormationOpportunities(baseLang)
		if len(opportunities) != 0 {
			t.Errorf("Expected no opportunities for language with insufficient changes, got %d", len(opportunities))
		}
	})

	t.Run("Sufficient Changes - Opportunities Detected", func(t *testing.T) {
		// Add more changes to meet threshold
		for i := 3; i < 8; i++ {
			change := LinguisticChange{
				ID:          fmt.Sprintf("test_change_%d", i),
				Type:        LinguisticChangeTypeMorphological,
				Description: fmt.Sprintf("Test change %d", i),
				Details:     fmt.Sprintf("Test change details %d", i),
				Timestamp:   time.Now(),
				Era:         "test",
				Trigger:     "test",
				Intensity:   0.6,
			}
			baseLang.AddLinguisticChange(change)
		}

		opportunities := engine.DetectDialectFormationOpportunities(baseLang)
		if len(opportunities) == 0 {
			t.Error("Expected opportunities for language with sufficient changes")
		}

		opportunity := opportunities[0]
		if opportunity.ChangeCount < 8 {
			t.Errorf("Expected at least 8 changes, got %d", opportunity.ChangeCount)
		}
		if opportunity.DivergenceScore == 0 {
			t.Error("Expected non-zero divergence score")
		}
		if opportunity.SuggestedDialect == "" {
			t.Error("Expected suggested dialect name")
		}
		if opportunity.Confidence == 0 {
			t.Error("Expected non-zero confidence")
		}
		if opportunity.Description == "" {
			t.Error("Expected opportunity description")
		}

		t.Logf("Detected opportunity: %s with %d%% divergence, %d changes, confidence %f",
			opportunity.SuggestedDialect, int(opportunity.DivergenceScore*100), opportunity.ChangeCount, opportunity.Confidence)
	})
}

// TestDialectSimilarityGrouping tests the grouping of similar dialects.
func TestDialectSimilarityGrouping(t *testing.T) {
	engine := NewDialectDetectionEngine()

	// Create a base language and several dialects
	baseLang, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create base language: %v", err)
	}

	// Create dialects with different characteristics
	dialect1 := baseLang.Clone(LanguageID{Family: "test", Branch: "test", Language: "test", Dialect: "dialect1"}, 43)
	dialect2 := baseLang.Clone(LanguageID{Family: "test", Branch: "test", Language: "test", Dialect: "dialect2"}, 44)
	dialect3 := baseLang.Clone(LanguageID{Family: "test", Branch: "test", Language: "test", Dialect: "dialect3"}, 45)

	// Set up parent-child relationships
	dialect1.ParentID = &baseLang.ID
	dialect2.ParentID = &baseLang.ID
	dialect3.ParentID = &baseLang.ID

	dialects := []*Language{dialect1, dialect2, dialect3}

	t.Run("Dialect Similarity Calculation", func(t *testing.T) {
		// Test similarity between dialects
		similarity := engine.calculateDialectSimilarity(dialect1, dialect2)
		if similarity < 0.5 {
			t.Errorf("Expected similarity >= 0.5 for dialects with shared parent, got %f", similarity)
		}

		// Test self-similarity
		selfSimilarity := engine.calculateDialectSimilarity(dialect1, dialect1)
		if selfSimilarity != 1.0 {
			t.Errorf("Expected self-similarity 1.0, got %f", selfSimilarity)
		}

		t.Logf("Dialect similarity: dialect1-dialect2=%f, dialect1-self=%f", similarity, selfSimilarity)
	})

	t.Run("Dialect Grouping", func(t *testing.T) {
		groups := engine.GroupSimilarDialects(dialects)
		if len(groups) == 0 {
			t.Error("Expected at least one dialect group")
		}

		// Should have one group with all dialects since they share a parent
		mainGroup := groups[0]
		if len(mainGroup.Dialects) < 3 {
			t.Errorf("Expected main group to contain at least 3 dialects, got %d", len(mainGroup.Dialects))
		}

		if mainGroup.SimilarityScore < 0.7 {
			t.Errorf("Expected similarity score >= 0.7, got %f", mainGroup.SimilarityScore)
		}

		if len(mainGroup.CommonFeatures) == 0 {
			t.Error("Expected common features to be identified")
		}

		t.Logf("Created dialect group: %s with %d dialects, similarity %f, features: %v",
			mainGroup.GroupName, len(mainGroup.Dialects), mainGroup.SimilarityScore, mainGroup.CommonFeatures)
	})
}

// TestDialectFormationSuggestion tests the suggestion of dialect formation opportunities.
func TestDialectFormationSuggestion(t *testing.T) {
	engine := NewDialectDetectionEngine()

	// Create multiple languages with different levels of change
	languages := make([]*Language, 0)

	// Language 1: No changes
	lang1, err := CreateRandomLanguage("test1", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create language 1: %v", err)
	}
	languages = append(languages, lang1)

	// Language 2: Some changes but not enough
	lang2, err := CreateRandomLanguage("test2", "ancient", 43)
	if err != nil {
		t.Fatalf("Failed to create language 2: %v", err)
	}
	for i := 0; i < 3; i++ {
		change := LinguisticChange{
			ID:          fmt.Sprintf("lang2_change_%d", i),
			Type:        LinguisticChangeTypeSound,
			Description: fmt.Sprintf("Language 2 change %d", i),
			Details:     fmt.Sprintf("Language 2 change details %d", i),
			Timestamp:   time.Now(),
			Era:         "test",
			Trigger:     "test",
			Intensity:   0.4,
		}
		lang2.AddLinguisticChange(change)
	}
	languages = append(languages, lang2)

	// Language 3: Sufficient changes for dialect formation
	lang3, err := CreateRandomLanguage("test3", "ancient", 44)
	if err != nil {
		t.Fatalf("Failed to create language 3: %v", err)
	}
	for i := 0; i < 8; i++ {
		change := LinguisticChange{
			ID:          fmt.Sprintf("lang3_change_%d", i),
			Type:        LinguisticChangeTypeMorphological,
			Description: fmt.Sprintf("Language 3 change %d", i),
			Details:     fmt.Sprintf("Language 3 change details %d", i),
			Timestamp:   time.Now(),
			Era:         "test",
			Trigger:     "test",
			Intensity:   0.6,
		}
		lang3.AddLinguisticChange(change)
	}
	languages = append(languages, lang3)

	t.Run("Dialect Formation Suggestions", func(t *testing.T) {
		suggestions := engine.SuggestDialectFormation(languages)
		if len(suggestions) == 0 {
			t.Error("Expected dialect formation suggestions")
		}

		// Should have suggestions for language 3
		var lang3Suggestion *DialectFormationOpportunity
		for i := range suggestions {
			if suggestions[i].LanguageID == lang3.ID {
				lang3Suggestion = &suggestions[i]
				break
			}
		}

		if lang3Suggestion == nil {
			t.Error("Expected suggestion for language 3")
		}

		if lang3Suggestion.ChangeCount < 8 {
			t.Errorf("Expected at least 8 changes in suggestion, got %d", lang3Suggestion.ChangeCount)
		}

		if lang3Suggestion.Confidence < 0.5 {
			t.Errorf("Expected confidence >= 0.5, got %f", lang3Suggestion.Confidence)
		}

		t.Logf("Generated %d dialect formation suggestions", len(suggestions))
		for i, suggestion := range suggestions {
			t.Logf("Suggestion %d: %s with %d%% divergence, confidence %f",
				i+1, suggestion.SuggestedDialect, int(suggestion.DivergenceScore*100), suggestion.Confidence)
		}
	})
}

// TestAutomaticDialectCreation tests the automatic creation of dialects.
func TestAutomaticDialectCreation(t *testing.T) {
	engine := NewDialectDetectionEngine()

	// Create a language with sufficient changes
	baseLang, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create base language: %v", err)
	}

	// Add changes to trigger dialect formation
	for i := 0; i < 10; i++ {
		change := LinguisticChange{
			ID:          fmt.Sprintf("auto_change_%d", i),
			Type:        LinguisticChangeTypeSound,
			Description: fmt.Sprintf("Auto change %d", i),
			Details:     fmt.Sprintf("Auto change details %d", i),
			Timestamp:   time.Now(),
			Era:         "test",
			Trigger:     "test",
			Intensity:   0.7,
		}
		baseLang.AddLinguisticChange(change)
	}

	t.Run("Automatic Dialect Creation", func(t *testing.T) {
		// Detect opportunities
		opportunities := engine.DetectDialectFormationOpportunities(baseLang)
		if len(opportunities) == 0 {
			t.Fatal("Expected dialect formation opportunities")
		}

		opportunity := opportunities[0]

		// Automatically create dialect
		dialect, err := engine.AutoCreateDialect(baseLang, opportunity, "test_era", 43)
		if err != nil {
			t.Fatalf("Failed to automatically create dialect: %v", err)
		}

		// Verify dialect properties
		if dialect.Name != opportunity.SuggestedDialect {
			t.Errorf("Expected dialect name '%s', got '%s'", opportunity.SuggestedDialect, dialect.Name)
		}

		if dialect.ParentID == nil || *dialect.ParentID != baseLang.ID {
			t.Error("Expected dialect to have correct parent ID")
		}

		// Verify automatic formation change was recorded
		dialectChanges := dialect.GetLinguisticChanges(LinguisticChangeTypeDialectal)
		if len(dialectChanges) == 0 {
			t.Error("Expected dialect formation change to be recorded")
		}

		var formationChange LinguisticChange
		for _, change := range dialectChanges {
			if change.Trigger == "automatic_detection" {
				formationChange = change
				break
			}
		}

		if formationChange.Trigger != "automatic_detection" {
			t.Error("Expected automatic detection change to be recorded")
		}

		// Verify parent language has dialect as child
		if len(baseLang.ChildIDs) == 0 {
			t.Error("Expected parent language to have child dialects")
		}

		foundChild := false
		for _, childID := range baseLang.ChildIDs {
			if childID == dialect.ID {
				foundChild = true
				break
			}
		}

		if !foundChild {
			t.Error("Expected parent language to have created dialect as child")
		}

		t.Logf("Successfully created automatic dialect: %s", dialect.Name)
		t.Logf("Dialect has %d linguistic changes", len(dialect.LinguisticChanges))
		t.Logf("Parent language has %d child dialects", len(baseLang.ChildIDs))
	})
}

// TestDialectDetectionIntegration tests the integration of dialect detection with the overall system.
func TestDialectDetectionIntegration(t *testing.T) {
	t.Run("Complete Dialect Detection Workflow", func(t *testing.T) {
		engine := NewDialectDetectionEngine()

		// Create a language family
		protoLang, err := CreateRandomLanguage("proto", "ancient", 42)
		if err != nil {
			t.Fatalf("Failed to create proto language: %v", err)
		}

		// Evolve the language to create divergence
		evolvedLang, err := EvolveLanguage(protoLang, "1000_years", []string{"cultural_shift", "migration"}, 43)
		if err != nil {
			t.Fatalf("Failed to evolve language: %v", err)
		}

		// Check for dialect formation opportunities
		opportunities := engine.DetectDialectFormationOpportunities(evolvedLang)
		t.Logf("Detected %d dialect formation opportunities", len(opportunities))

		// If opportunities exist, create dialects automatically
		if len(opportunities) > 0 {
			opportunity := opportunities[0]
			dialect, err := engine.AutoCreateDialect(evolvedLang, opportunity, "auto_era", 44)
			if err != nil {
				t.Fatalf("Failed to create automatic dialect: %v", err)
			}

			t.Logf("Automatically created dialect: %s", dialect.Name)
			t.Logf("Dialect divergence: %d%%", int(opportunity.DivergenceScore*100))
			t.Logf("Dialect confidence: %f", opportunity.Confidence)

			// Group similar dialects
			allLanguages := []*Language{protoLang, evolvedLang, dialect}
			groups := engine.GroupSimilarDialects(allLanguages)
			t.Logf("Created %d dialect groups", len(groups))

			for i, group := range groups {
				t.Logf("Group %d: %s with %d dialects, similarity %f",
					i+1, group.GroupName, len(group.Dialects), group.SimilarityScore)
			}
		}

		t.Logf("Successfully completed dialect detection workflow")
	})
}
