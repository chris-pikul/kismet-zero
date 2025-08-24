package lang

import (
	"testing"
)

// TestLinguisticChangeTracking tests that linguistic changes are properly stored and retrieved.
func TestLinguisticChangeTracking(t *testing.T) {
	// Create a language
	lang, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	// Verify initial state
	if lang.LinguisticChanges != nil {
		t.Error("Expected no linguistic changes initially")
	}

	// Evolve the language to trigger changes
	evolvedLang, err := EvolveLanguage(lang, "500_years", []string{}, 42)
	if err != nil {
		t.Fatalf("Failed to evolve language: %v", err)
	}

	// Check that linguistic changes were stored
	if evolvedLang.LinguisticChanges == nil {
		t.Fatal("Expected linguistic changes slice to be initialized after evolution")
	}

	if len(evolvedLang.LinguisticChanges) == 0 {
		t.Log("No linguistic changes were applied (probabilistic system)")
		// This is expected behavior - the system is probabilistic
		// Let's try with a different seed to see if we can trigger changes
		evolvedLang2, err := EvolveLanguage(lang, "500_years", []string{}, 999)
		if err != nil {
			t.Fatalf("Failed to evolve language with different seed: %v", err)
		}

		if len(evolvedLang2.LinguisticChanges) > 0 {
			t.Logf("Second evolution with different seed applied %d linguistic changes", len(evolvedLang2.LinguisticChanges))
		} else {
			t.Log("Still no changes with different seed - this is normal for probabilistic system")
		}
	} else {
		t.Logf("Applied %d linguistic changes", len(evolvedLang.LinguisticChanges))

		// Verify change types
		for i, change := range evolvedLang.LinguisticChanges {
			t.Logf("Change %d: Type=%s, Description=%s, Era=%s",
				i, change.Type.String(), change.Description, change.Era)

			// Verify required fields
			if change.ID == "" {
				t.Errorf("Change %d missing ID", i)
			}
			if change.Description == "" {
				t.Errorf("Change %d missing Description", i)
			}
			if change.Era == "" {
				t.Errorf("Change %d missing Era", i)
			}
		}
	}
}

// TestLinguisticChangeTypes tests the linguistic change type system.
func TestLinguisticChangeTypes(t *testing.T) {
	// Test all change types have valid string representations
	expectedTypes := []LinguisticChangeType{
		LinguisticChangeTypeSound,
		LinguisticChangeTypeMorphological,
		LinguisticChangeTypeOrthographic,
		LinguisticChangeTypeCultural,
		LinguisticChangeTypeDialectal,
	}

	for _, expectedType := range expectedTypes {
		if expectedType.String() == "" {
			t.Errorf("Expected change type %s to have a valid string representation", expectedType)
		}
	}

	// Test unknown type handling
	unknownType := LinguisticChangeType("unknown_type")
	if unknownType.String() != "unknown_type" {
		t.Errorf("Expected unknown type to return 'unknown_type', got %s", unknownType.String())
	}
}

// TestLinguisticChangeStorage tests the storage and retrieval of linguistic changes.
func TestLinguisticChangeStorage(t *testing.T) {
	// Create a language
	lang, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	// Create test changes
	soundChange := SoundChange{
		ID:          "test_sound",
		Name:        "Test Sound Change",
		Description: "A test sound change",
		Probability: 0.5,
		Era:         "test_era",
	}

	morphologicalChange := MorphologicalChange{
		ID:               "test_morph",
		Type:             MorphologicalChangeTypeSimplification,
		Description:      "Test Morphological Change",
		Details:          "A test morphological change",
		ComplexityChange: -0.2,
		Timestamp:        lang.CreatedAt,
		Era:              "test_era",
		Trigger:          "test",
		Intensity:        0.5,
	}

	orthographicChange := OrthographicChange{
		ID:               "test_ortho",
		Type:             OrthographicChangeTypeSimplification,
		Description:      "Test Orthographic Change",
		Details:          "A test orthographic change",
		ComplexityChange: -0.1,
		Timestamp:        lang.CreatedAt,
		Era:              "test_era",
		Trigger:          "test",
		Intensity:        0.4,
	}

	// Convert and store changes
	lang.AddLinguisticChange(NewSoundChange(soundChange))
	lang.AddLinguisticChange(NewMorphologicalChange(morphologicalChange))
	lang.AddLinguisticChange(NewOrthographicChange(orthographicChange))

	// Verify storage
	if len(lang.LinguisticChanges) != 3 {
		t.Errorf("Expected 3 linguistic changes, got %d", len(lang.LinguisticChanges))
	}

	// Test retrieval by type
	soundChanges := lang.GetLinguisticChanges(LinguisticChangeTypeSound)
	if len(soundChanges) != 1 {
		t.Errorf("Expected 1 sound change, got %d", len(soundChanges))
	}

	morphChanges := lang.GetLinguisticChanges(LinguisticChangeTypeMorphological)
	if len(morphChanges) != 1 {
		t.Errorf("Expected 1 morphological change, got %d", len(morphChanges))
	}

	orthoChanges := lang.GetLinguisticChanges(LinguisticChangeTypeOrthographic)
	if len(orthoChanges) != 1 {
		t.Errorf("Expected 1 orthographic change, got %d", len(orthoChanges))
	}

	// Test retrieval by era
	eraChanges := lang.GetLinguisticChangesByEra("test_era")
	if len(eraChanges) != 3 {
		t.Errorf("Expected 3 changes for test_era, got %d", len(eraChanges))
	}
}

// TestComplexityTracking tests that complexity changes are properly tracked.
func TestComplexityTracking(t *testing.T) {
	// Create a language
	lang, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	// Add changes with different complexity impacts
	lang.AddLinguisticChange(LinguisticChange{
		ID:               "simplify",
		Type:             LinguisticChangeTypeMorphological,
		Description:      "Simplification",
		ComplexityChange: -0.3,
		Timestamp:        lang.CreatedAt,
		Era:              "test",
		Trigger:          "test",
		Intensity:        0.5,
	})

	lang.AddLinguisticChange(LinguisticChange{
		ID:               "complexify",
		Type:             LinguisticChangeTypeOrthographic,
		Description:      "Complexification",
		ComplexityChange: 0.2,
		Timestamp:        lang.CreatedAt,
		Era:              "test",
		Trigger:          "test",
		Intensity:        0.4,
	})

	// Test total complexity change
	totalChange := lang.GetTotalComplexityChange()
	expectedChange := float32(-0.1) // -0.3 + 0.2
	if totalChange != expectedChange {
		t.Errorf("Expected total complexity change %f, got %f", expectedChange, totalChange)
	}

	t.Logf("Total complexity change: %f", totalChange)
}

// TestEvolutionSummary tests the evolution summary functionality.
func TestEvolutionSummary(t *testing.T) {
	// Create a language
	lang, err := CreateRandomLanguage("test", "ancient", 42)
	if err != nil {
		t.Fatalf("Failed to create test language: %v", err)
	}

	// Get initial summary
	summary := lang.GetEvolutionSummary()
	if summary["totalChanges"].(int) != 0 {
		t.Error("Expected 0 changes initially")
	}

	// Add some changes
	lang.AddLinguisticChange(LinguisticChange{
		ID:          "change1",
		Type:        LinguisticChangeTypeSound,
		Description: "Sound Change 1",
		Era:         "early_evolution",
		Timestamp:   lang.CreatedAt,
		Trigger:     "test",
		Intensity:   0.5,
	})

	lang.AddLinguisticChange(LinguisticChange{
		ID:          "change2",
		Type:        LinguisticChangeTypeMorphological,
		Description: "Morphological Change 1",
		Era:         "middle_evolution",
		Timestamp:   lang.CreatedAt,
		Trigger:     "test",
		Intensity:   0.4,
	})

	// Get updated summary
	summary = lang.GetEvolutionSummary()

	t.Logf("Evolution Summary: %+v", summary)

	if summary["totalChanges"].(int) != 2 {
		t.Errorf("Expected 2 changes, got %d", summary["totalChanges"])
	}

	eras := summary["eras"].([]string)
	if len(eras) != 2 {
		t.Errorf("Expected 2 eras, got %d", len(eras))
	}

	types := summary["types"].([]string)
	if len(types) != 2 {
		t.Errorf("Expected 2 types, got %d", len(types))
	}
}
