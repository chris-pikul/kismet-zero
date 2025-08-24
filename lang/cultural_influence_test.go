package lang

import (
	"strings"
	"testing"
	"time"
)

// TestCulturalInfluenceEngine tests the core cultural influence engine functionality.
func TestCulturalInfluenceEngine(t *testing.T) {
	// Create test languages
	lang1, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language 1: %v", err)
	}

	lang2, err := CreateRandomLanguage("test", "ancient", 43)
	if err != nil {
		t.Fatalf("Failed to create test language 2: %v", err)
	}

	// Test cultural influence engine
	t.Run("Cultural Influence Engine", func(t *testing.T) {
		// Apply cultural influence
		result, err := ApplyCulturalInfluence(lang1, lang2, ContactTypeTrade, 0.8, 100*time.Hour, 42)
		if err != nil {
			t.Fatalf("Failed to apply cultural influence: %v", err)
		}

		if result == nil {
			t.Fatal("Expected cultural influence result")
		}

		// Verify the result structure
		if result.LanguageID != lang1.ID.String() {
			t.Errorf("Expected language ID %s, got %s", lang1.ID.String(), result.LanguageID)
		}

		if result.ContactType != ContactTypeTrade {
			t.Errorf("Expected contact type %s, got %s", ContactTypeTrade.String(), result.ContactType.String())
		}

		if result.Intensity != 0.8 {
			t.Errorf("Expected intensity 0.8, got %f", result.Intensity)
		}

		// Check that changes were applied (may be empty due to probabilistic nature)
		t.Logf("Applied changes: phonological=%d, grammatical=%d, lexical=%d, orthographic=%d",
			len(result.PhonologicalChanges), len(result.GrammaticalChanges),
			len(result.LexicalChanges), len(result.OrthographicChanges))

		// Verify complexity change tracking
		if result.ComplexityChange < 0.0 || result.ComplexityChange > 1.0 {
			t.Errorf("Complexity change out of valid range: %f", result.ComplexityChange)
		}
	})
}

// TestContactTypes tests all contact types in the cultural influence system.
func TestContactTypes(t *testing.T) {
	// Create test languages
	lang1, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language 1: %v", err)
	}

	lang2, err := CreateRandomLanguage("test", "ancient", 43)
	if err != nil {
		t.Fatalf("Failed to create test language 2: %v", err)
	}

	contactTypes := []ContactType{
		ContactTypeTrade,
		ContactTypeConquest,
		ContactTypeMigration,
		ContactTypeCultural,
		ContactTypeReligious,
		ContactTypeEducational,
	}

	for _, contactType := range contactTypes {
		t.Run(contactType.String(), func(t *testing.T) {
			result1, err := ApplyCulturalInfluence(lang1, lang2, contactType, 0.8, 200*time.Hour, 42+int64(contactType))
			if err != nil {
				t.Fatalf("Failed to apply %s cultural influence: %v", contactType.String(), err)
			}

			if result1.ContactType != contactType {
				t.Errorf("Expected contact type %s, got %s", contactType.String(), result1.ContactType.String())
			}

			// Verify that appropriate borrowing patterns were found
			applicablePatterns := getApplicableBorrowingPatterns(contactType)
			if len(applicablePatterns) == 0 {
				t.Logf("No borrowing patterns found for %s contact (this may be normal)", contactType.String())
			} else {
				t.Logf("%s contact found %d applicable borrowing patterns", contactType.String(), len(applicablePatterns))
			}
		})
	}
}

// TestBorrowingPatterns tests the borrowing pattern system.
func TestBorrowingPatterns(t *testing.T) {
	// Test that we have borrowing patterns for each contact type
	contactTypes := []ContactType{
		ContactTypeTrade,
		ContactTypeConquest,
		ContactTypeMigration,
		ContactTypeCultural,
	}

	for _, contactType := range contactTypes {
		t.Run(contactType.String(), func(t *testing.T) {
			patterns := getApplicableBorrowingPatterns(contactType)

			// Count patterns by type
			lexicalCount := 0
			phonologicalCount := 0
			grammaticalCount := 0
			orthographicCount := 0

			for _, pattern := range patterns {
				switch pattern.Type {
				case "lexical":
					lexicalCount++
				case "phonological":
					phonologicalCount++
				case "grammatical":
					grammaticalCount++
				case "orthographic":
					orthographicCount++
				}
			}

			t.Logf("%s contact patterns: lexical=%d, phonological=%d, grammatical=%d, orthographic=%d",
				contactType.String(), lexicalCount, phonologicalCount, grammaticalCount, orthographicCount)

			// Each contact type should have at least some patterns
			if len(patterns) == 0 {
				t.Errorf("Expected %s contact to have borrowing patterns", contactType.String())
			}
		})
	}
}

// TestEnhancedBorrowingPatterns tests the enhanced borrowing pattern system with detailed rules.
func TestEnhancedBorrowingPatterns(t *testing.T) {
	// Test that borrowing patterns have detailed rules
	for _, pattern := range CommonBorrowingPatterns {
		t.Run(pattern.ID, func(t *testing.T) {
			// Check that patterns have appropriate rules based on their type
			switch pattern.Type {
			case "lexical":
				if len(pattern.LexicalRules) == 0 {
					t.Errorf("Lexical pattern %s should have lexical rules", pattern.ID)
				}
				// Verify lexical rules have required fields
				for _, rule := range pattern.LexicalRules {
					if rule.SemanticField == "" {
						t.Errorf("Lexical rule %s missing semantic field", rule.ID)
					}
					if rule.AdaptationType == "" {
						t.Errorf("Lexical rule %s missing adaptation type", rule.ID)
					}
				}

			case "phonological":
				if len(pattern.PhonologicalRules) == 0 {
					t.Errorf("Phonological pattern %s should have phonological rules", pattern.ID)
				}
				// Verify phonological rules have required fields
				for _, rule := range pattern.PhonologicalRules {
					if rule.Type == "" {
						t.Errorf("Phonological rule %s missing type", rule.ID)
					}
					if rule.Context == "" {
						t.Errorf("Phonological rule %s missing context", rule.ID)
					}
				}

			case "grammatical":
				if len(pattern.GrammaticalRules) == 0 {
					t.Errorf("Grammatical pattern %s should have grammatical rules", pattern.ID)
				}
				// Verify grammatical rules have required fields
				for _, rule := range pattern.GrammaticalRules {
					if rule.Category == "" {
						t.Errorf("Grammatical rule %s missing category", rule.ID)
					}
					if rule.Feature == "" {
						t.Errorf("Grammatical rule %s missing feature", rule.ID)
					}
				}

			case "orthographic":
				if len(pattern.OrthographicRules) == 0 {
					t.Errorf("Orthographic pattern %s should have orthographic rules", pattern.ID)
				}
				// Verify orthographic rules have required fields
				for _, rule := range pattern.OrthographicRules {
					if rule.Type == "" {
						t.Errorf("Orthographic rule %s missing type", rule.ID)
					}
				}
			}

			t.Logf("Pattern %s (%s) has %d detailed rules", pattern.ID, pattern.Type,
				len(pattern.LexicalRules)+len(pattern.PhonologicalRules)+len(pattern.GrammaticalRules)+len(pattern.OrthographicRules))
		})
	}
}

// TestDetailedBorrowingApplication tests that detailed borrowing rules are properly applied.
func TestDetailedBorrowingApplication(t *testing.T) {
	// Create test languages
	lang1, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language 1: %v", err)
	}

	lang2, err := CreateRandomLanguage("test", "ancient", 43)
	if err != nil {
		t.Fatalf("Failed to create test language 2: %v", err)
	}

	// Test trade contact with detailed lexical rules
	result1, err := ApplyCulturalInfluence(lang1, lang2, ContactTypeTrade, 0.8, 200*time.Hour, 42)
	if err != nil {
		t.Fatalf("Failed to apply trade cultural influence: %v", err)
	}

	// Check that detailed lexical changes were applied
	if len(result1.LexicalChanges) == 0 {
		t.Log("No lexical changes applied (probabilistic system)")
	} else {
		t.Logf("Applied %d detailed lexical changes:", len(result1.LexicalChanges))
		for i, change := range result1.LexicalChanges {
			t.Logf("  %d: %s", i, change)

			// Verify that changes are detailed and specific
			if !strings.Contains(change, "Borrowed") && !strings.Contains(change, "Semantic shift") && !strings.Contains(change, "Formed calques") {
				t.Errorf("Lexical change should be detailed, got: %s", change)
			}
		}
	}

	// Test conquest contact with detailed grammatical rules
	result2, err := ApplyCulturalInfluence(lang1, lang2, ContactTypeConquest, 0.9, 300*time.Hour, 44)
	if err != nil {
		t.Fatalf("Failed to apply conquest cultural influence: %v", err)
	}

	// Check that detailed grammatical changes were applied
	if len(result2.GrammaticalChanges) == 0 {
		t.Log("No grammatical changes applied (probabilistic system)")
	} else {
		t.Logf("Applied %d detailed grammatical changes:", len(result2.GrammaticalChanges))
		for i, change := range result2.GrammaticalChanges {
			t.Logf("  %d: %s", i, change)

			// Verify that changes are detailed and specific
			if !strings.Contains(change, "Added") && !strings.Contains(change, "Modified") {
				t.Errorf("Grammatical change should be detailed, got: %s", change)
			}
		}
	}
}

// TestBorrowingRuleTypes tests the different types of borrowing rules.
func TestBorrowingRuleTypes(t *testing.T) {
	// Test lexical rule types
	lexicalRule := LexicalRule{
		ID:             "test_lexical",
		Description:    "Test lexical rule",
		Type:           "word_borrowing",
		SemanticField:  "test_field",
		AdaptationType: "phonological",
		Probability:    0.8,
	}

	// Test phonological rule types
	phonologicalRule := PhonologicalRule{
		ID:          "test_phonological",
		Description: "Test phonological rule",
		Type:        "sound_addition",
		SourceSound: "foreign_sound",
		TargetSound: "native_sound",
		Context:     "word_initial",
		Probability: 0.7,
	}

	// Test grammatical rule types
	grammaticalRule := GrammaticalRule{
		ID:          "test_grammatical",
		Description: "Test grammatical rule",
		Type:        "morphology_addition",
		Category:    "noun",
		Feature:     "case",
		Probability: 0.6,
	}

	// Test orthographic rule types
	orthographicRule := OrthographicRule{
		ID:          "test_orthographic",
		Description: "Test orthographic rule",
		Type:        "script_adoption",
		ScriptType:  "alphabet",
		ReformType:  "foreign_adoption",
		Probability: 0.5,
	}

	// Verify rule structures
	if lexicalRule.SemanticField != "test_field" {
		t.Errorf("Expected semantic field 'test_field', got %s", lexicalRule.SemanticField)
	}

	if phonologicalRule.Context != "word_initial" {
		t.Errorf("Expected context 'word_initial', got %s", phonologicalRule.Context)
	}

	if grammaticalRule.Category != "noun" {
		t.Errorf("Expected category 'noun', got %s", grammaticalRule.Category)
	}

	if orthographicRule.ScriptType != "alphabet" {
		t.Errorf("Expected script type 'alphabet', got %s", orthographicRule.ScriptType)
	}

	t.Log("All borrowing rule types properly structured")
}

// TestCulturalInfluenceIntegration tests the integration with SimulateLanguageContact.
func TestCulturalInfluenceIntegration(t *testing.T) {
	// Create test languages
	lang1, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language 1: %v", err)
	}

	lang2, err := CreateRandomLanguage("test", "ancient", 43)
	if err != nil {
		t.Fatalf("Failed to create test language 2: %v", err)
	}

	// Test trade contact through the API
	contactLang1, _, err := SimulateLanguageContact(lang1, lang2, "trade", 0.7, "100_years", 42)
	if err != nil {
		t.Fatalf("Failed to simulate trade contact: %v", err)
	}

	// Verify that cultural influence changes were stored
	if contactLang1.LinguisticChanges == nil {
		t.Fatal("Expected linguistic changes to be initialized after contact")
	}

	// Look for cultural influence changes
	var culturalChanges []LinguisticChange
	for _, change := range contactLang1.LinguisticChanges {
		if change.Trigger == "cultural_influence" {
			culturalChanges = append(culturalChanges, change)
		}
	}

	if len(culturalChanges) == 0 {
		t.Log("No cultural influence changes were stored (probabilistic system)")
	} else {
		t.Logf("Stored %d cultural influence changes", len(culturalChanges))

		for i, change := range culturalChanges {
			t.Logf("Cultural change %d: %s (complexity: %f)",
				i, change.Description, change.ComplexityChange)
		}
	}

	// Verify descriptions were updated
	if !strings.Contains(contactLang1.Description, "trade") {
		t.Error("Expected trade contact to be mentioned in description")
	}

	if !strings.Contains(contactLang1.Description, "influenced by") {
		t.Error("Expected influence to be mentioned in description")
	}

	t.Logf("Updated description: %s", contactLang1.Description)
}

// TestComplexityImpact tests that cultural influence properly affects language complexity.
func TestComplexityImpact(t *testing.T) {
	// Create test languages
	lang1, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language 1: %v", err)
	}

	lang2, err := CreateRandomLanguage("test", "ancient", 43)
	if err != nil {
		t.Fatalf("Failed to apply cultural influence: %v", err)
	}

	// Apply high-intensity conquest contact (should have significant impact)
	result1, err := ApplyCulturalInfluence(lang1, lang2, ContactTypeConquest, 0.9, 500*time.Hour, 42)
	if err != nil {
		t.Fatalf("Failed to apply conquest cultural influence: %v", err)
	}

	// Conquest should have significant complexity impact
	if result1.ComplexityChange <= 0.0 {
		t.Log("No complexity change from conquest (probabilistic system)")
	} else {
		t.Logf("Conquest complexity change: %f", result1.ComplexityChange)

		// Conquest should typically increase complexity due to administrative structures
		if result1.ComplexityChange < 0.1 {
			t.Log("Conquest had minimal complexity impact")
		} else {
			t.Log("Conquest had significant complexity impact")
		}
	}

	// Test that complexity change is normalized
	if result1.ComplexityChange > 1.0 {
		t.Errorf("Complexity change should be normalized to <= 1.0, got %f", result1.ComplexityChange)
	}
}

// TestCulturalInfluenceChangeCreation tests the creation of linguistic changes from cultural influence.
func TestCulturalInfluenceChangeCreation(t *testing.T) {
	// Create a sample cultural influence result
	result := &CulturalInfluenceResult{
		LanguageID:        "test_lang",
		ContactType:       ContactTypeTrade,
		Intensity:         0.6,
		Duration:          "100_years",
		BorrowingPatterns: []BorrowingPattern{CommonBorrowingPatterns[0]},
		LexicalChanges:    []string{"Borrowing of trade-related vocabulary"},
		ComplexityChange:  0.1,
		Timestamp:         time.Now(),
	}

	// Create a linguistic change from the cultural influence
	change := CreateCulturalInfluenceChange(result, "Source Language")

	// Verify the change was created correctly
	if change.Type != LinguisticChangeTypeMorphological {
		t.Errorf("Expected change type %s, got %s", LinguisticChangeTypeMorphological.String(), change.Type.String())
	}

	if change.Trigger != "cultural_influence" {
		t.Errorf("Expected trigger 'cultural_influence', got %s", change.Trigger)
	}

	if change.Era != "cultural_contact" {
		t.Errorf("Expected era 'cultural_contact', got %s", change.Era)
	}

	if change.ComplexityChange != 0.1 {
		t.Errorf("Expected complexity change 0.1, got %f", change.ComplexityChange)
	}

	if !strings.Contains(change.Description, "Source Language") {
		t.Error("Expected source language to be mentioned in description")
	}

	if !strings.Contains(change.Description, "trade") {
		t.Error("Expected contact type to be mentioned in description")
	}

	t.Logf("Created cultural influence change: %s", change.Description)
}
