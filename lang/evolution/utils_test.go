package evolution

import (
	"strings"
	"testing"
	"time"
)

// TestIDGenerationFunctions tests all ID generation utility functions
func TestIDGenerationFunctions(t *testing.T) {
	// Test GenerateChangeID
	changeID := GenerateChangeID("test", "123")
	if !strings.Contains(changeID, "test_123") {
		t.Errorf("Expected change ID to contain 'test_123', got '%s'", changeID)
	}

	// Test GenerateContactID
	contactID := GenerateContactID("contact", "456")
	if !strings.Contains(contactID, "contact_456") {
		t.Errorf("Expected contact ID to contain 'contact_456', got '%s'", contactID)
	}

	// Test GenerateDialectID
	dialectID := GenerateDialectID("dialect", "789")
	if !strings.Contains(dialectID, "dialect_789") {
		t.Errorf("Expected dialect ID to contain 'dialect_789', got '%s'", dialectID)
	}

	// Test GenerateEvolutionID
	evolutionID := GenerateEvolutionID("evolution", "101")
	if !strings.Contains(evolutionID, "evolution_101") {
		t.Errorf("Expected evolution ID to contain 'evolution_101', got '%s'", evolutionID)
	}

	// Test GenerateSoundChangeID
	soundChangeID := GenerateSoundChangeID("sound", "202")
	if !strings.Contains(soundChangeID, "sound_202") {
		t.Errorf("Expected sound change ID to contain 'sound_202', got '%s'", soundChangeID)
	}

	// Test GenerateTimestampID
	timestampID := GenerateTimestampID("timestamp", "303")
	if !strings.Contains(timestampID, "timestamp_303") {
		t.Errorf("Expected timestamp ID to contain 'timestamp_303', got '%s'", timestampID)
	}

	// Test GenerateNanoID
	nanoID := GenerateNanoID("nano", "404")
	if !strings.Contains(nanoID, "nano_404") {
		t.Errorf("Expected nano ID to contain 'nano_404', got '%s'", nanoID)
	}
}

// TestEnumStringFunctions tests the enum string utility functions
func TestEnumStringFunctions(t *testing.T) {
	// Test EnumString with valid values
	testArray := []string{"unknown", "first", "second", "third"}

	result := EnumString(1, testArray, "unknown")
	if result != "first" {
		t.Errorf("Expected 'first', got '%s'", result)
	}

	result = EnumString(2, testArray, "unknown")
	if result != "second" {
		t.Errorf("Expected 'second', got '%s'", result)
	}

	// Test EnumString with out-of-bounds values
	result = EnumString(-1, testArray, "unknown")
	if result != "unknown" {
		t.Errorf("Expected 'unknown' for negative index, got '%s'", result)
	}

	result = EnumString(10, testArray, "unknown")
	if result != "unknown" {
		t.Errorf("Expected 'unknown' for out-of-bounds index, got '%s'", result)
	}

	// Test EnumStringByte with valid values
	result = EnumStringByte(1, testArray, "unknown")
	if result != "first" {
		t.Errorf("Expected 'first', got '%s'", result)
	}

	// Test EnumStringByte with out-of-bounds values
	result = EnumStringByte(255, testArray, "unknown")
	if result != "unknown" {
		t.Errorf("Expected 'unknown' for out-of-bounds byte index, got '%s'", result)
	}
}

// TestChangeCreationFactories tests all change creation factory functions
func TestChangeCreationFactories(t *testing.T) {
	// Test NewLinguisticChange
	change := NewLinguisticChange(
		ChangeTypeSoundShift,
		ChangeDirectionModifying,
		"Test change",
		"test_era",
		"test_trigger",
		0.5,
	)

	if change.Type != ChangeTypeSoundShift {
		t.Errorf("Expected ChangeTypeSoundShift, got %v", change.Type)
	}
	if change.Direction != ChangeDirectionModifying {
		t.Errorf("Expected ChangeDirectionModifying, got %v", change.Direction)
	}
	if change.Description != "Test change" {
		t.Errorf("Expected 'Test change', got '%s'", change.Description)
	}
	if change.Era != "test_era" {
		t.Errorf("Expected 'test_era', got '%s'", change.Era)
	}
	if change.Trigger != "test_trigger" {
		t.Errorf("Expected 'test_trigger', got '%s'", change.Trigger)
	}
	if change.Intensity != 0.5 {
		t.Errorf("Expected 0.5, got %f", change.Intensity)
	}

	// Test NewSoundChange
	soundChange := NewSoundChange("Sound test", "sound_era", 0.7)
	if soundChange.Type != ChangeTypeSoundShift {
		t.Errorf("Expected ChangeTypeSoundShift, got %v", soundChange.Type)
	}
	if soundChange.Direction != ChangeDirectionModifying {
		t.Errorf("Expected ChangeDirectionModifying, got %v", soundChange.Direction)
	}
	if soundChange.Description != "Sound test" {
		t.Errorf("Expected 'Sound test', got '%s'", soundChange.Description)
	}

	// Test NewMorphologicalChange
	morphChange := NewMorphologicalChange("Morph test", "morph_era", 0.6)
	if morphChange.Type != ChangeTypeMorphological {
		t.Errorf("Expected ChangeTypeMorphological, got %v", morphChange.Type)
	}
	if morphChange.Direction != ChangeDirectionModifying {
		t.Errorf("Expected ChangeDirectionModifying, got %v", morphChange.Direction)
	}
	if morphChange.Description != "Morph test" {
		t.Errorf("Expected 'Morph test', got '%s'", morphChange.Description)
	}

	// Test NewContactChange
	contactChange := NewContactChange(ChangeTypeLexical, "Contact test", "contact_era", 0.8, "test_culture")
	if contactChange.Type != ChangeTypeLexical {
		t.Errorf("Expected ChangeTypeLexical, got %v", contactChange.Type)
	}
	if contactChange.Direction != ChangeDirectionAdditive {
		t.Errorf("Expected ChangeDirectionAdditive, got %v", contactChange.Direction)
	}
	if contactChange.CultureInfluence != "test_culture" {
		t.Errorf("Expected 'test_culture', got '%s'", contactChange.CultureInfluence)
	}

	// Test NewAdaptationChange
	adaptChange := NewAdaptationChange(ChangeTypeMorphological, "Adapt test", "adapt_era", 0.9)
	if adaptChange.Type != ChangeTypeMorphological {
		t.Errorf("Expected ChangeTypeMorphological, got %v", adaptChange.Type)
	}
	if adaptChange.Direction != ChangeDirectionModifying {
		t.Errorf("Expected ChangeDirectionModifying, got %v", adaptChange.Direction)
	}
	if adaptChange.Description != "Adapt test" {
		t.Errorf("Expected 'Adapt test', got '%s'", adaptChange.Description)
	}
}

// TestIDGenerationConsistency tests that ID generation is consistent
func TestIDGenerationConsistency(t *testing.T) {
	// Test that different prefixes produce different results
	id1 := GenerateChangeID("test", "123")
	id3 := GenerateChangeID("different", "123")
	if strings.Contains(id1, id3) || strings.Contains(id3, id1) {
		t.Error("Expected different prefixes to produce different IDs")
	}

	// Test that nanosecond IDs are always different
	nanoID1 := GenerateNanoID("test", "123")
	time.Sleep(1 * time.Millisecond)
	nanoID2 := GenerateNanoID("test", "123")
	if nanoID1 == nanoID2 {
		t.Error("Expected nanosecond IDs to be different")
	}

	// Test that same inputs with different timestamps produce different results
	// Note: Unix() returns seconds, so we need to wait at least 1 second
	time.Sleep(1100 * time.Millisecond) // Wait 1.1 seconds to ensure different Unix timestamps
	id2 := GenerateChangeID("test", "123")
	if id1 == id2 {
		t.Error("Expected different IDs for different timestamps")
	}
}

// TestEnumStringEdgeCases tests edge cases for enum string functions
func TestEnumStringEdgeCases(t *testing.T) {
	// Test with empty array
	result := EnumString(0, []string{}, "default")
	if result != "default" {
		t.Errorf("Expected 'default' for empty array, got '%s'", result)
	}

	// Test with nil array
	result = EnumString(0, nil, "default")
	if result != "default" {
		t.Errorf("Expected 'default' for nil array, got '%s'", result)
	}

	// Test with single element array
	singleArray := []string{"only"}
	result = EnumString(0, singleArray, "default")
	if result != "only" {
		t.Errorf("Expected 'only', got '%s'", result)
	}

	result = EnumString(1, singleArray, "default")
	if result != "default" {
		t.Errorf("Expected 'default' for out-of-bounds index, got '%s'", result)
	}
}
