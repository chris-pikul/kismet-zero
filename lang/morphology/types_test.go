package morphology

import (
	"encoding/json"
	"testing"
)

func TestWordConceptIDs(t *testing.T) {
	// Test Word with ConceptIDs
	word := Word{
		ID:         "test-word",
		Meaning:    "test meaning",
		Category:   WordCategoryNoun,
		ConceptIDs: []ConceptID{"test-concept", "another-concept"},
	}

	// Test JSON marshaling
	data, err := json.Marshal(word)
	if err != nil {
		t.Fatalf("Failed to marshal Word with ConceptIDs: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Word
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Word with ConceptIDs: %v", err)
	}

	// Verify ConceptIDs are preserved
	if len(unmarshaled.ConceptIDs) != 2 {
		t.Errorf("Expected 2 ConceptIDs, got %d", len(unmarshaled.ConceptIDs))
	}
	if unmarshaled.ConceptIDs[0] != "test-concept" {
		t.Errorf("Expected first ConceptID 'test-concept', got %s", unmarshaled.ConceptIDs[0])
	}
	if unmarshaled.ConceptIDs[1] != "another-concept" {
		t.Errorf("Expected second ConceptID 'another-concept', got %s", unmarshaled.ConceptIDs[1])
	}
}

func TestWordWithoutConceptIDs(t *testing.T) {
	// Test Word without ConceptIDs (backward compatibility)
	word := Word{
		ID:       "test-word",
		Meaning:  "test meaning",
		Category: WordCategoryNoun,
		// ConceptIDs field omitted
	}

	// Test JSON marshaling
	data, err := json.Marshal(word)
	if err != nil {
		t.Fatalf("Failed to marshal Word without ConceptIDs: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Word
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Word without ConceptIDs: %v", err)
	}

	// Verify ConceptIDs is nil (zero value when field is omitted)
	if unmarshaled.ConceptIDs != nil {
		t.Error("Expected ConceptIDs to be nil when field is omitted, got non-nil")
	}
}

func TestMorphemeFeatures(t *testing.T) {
	// Test Morpheme with Features
	morpheme := Morpheme{
		ID:      "test-morpheme",
		Type:    MorphemeTypeSuffix,
		Meaning: "past tense",
		Features: map[FeatureKey]FeatureVal{
			FeatureTense:  "past",
			FeatureAspect: "perfective",
		},
	}

	// Test JSON marshaling
	data, err := json.Marshal(morpheme)
	if err != nil {
		t.Fatalf("Failed to marshal Morpheme with Features: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Morpheme
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Morpheme with Features: %v", err)
	}

	// Verify Features are preserved
	if len(unmarshaled.Features) != 2 {
		t.Errorf("Expected 2 Features, got %d", len(unmarshaled.Features))
	}
	if unmarshaled.Features[FeatureTense] != "past" {
		t.Errorf("Expected Tense 'past', got %s", unmarshaled.Features[FeatureTense])
	}
	if unmarshaled.Features[FeatureAspect] != "perfective" {
		t.Errorf("Expected Aspect 'perfective', got %s", unmarshaled.Features[FeatureAspect])
	}
}

func TestMorphemeWithoutFeatures(t *testing.T) {
	// Test Morpheme without Features (backward compatibility)
	morpheme := Morpheme{
		ID:      "test-morpheme",
		Type:    MorphemeTypeSuffix,
		Meaning: "past tense",
		// Features field omitted
	}

	// Test JSON marshaling
	data, err := json.Marshal(morpheme)
	if err != nil {
		t.Fatalf("Failed to marshal Morpheme without Features: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Morpheme
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Morpheme without Features: %v", err)
	}

	// Verify Features is nil (not empty map)
	if unmarshaled.Features != nil {
		t.Error("Expected Features to be nil, got non-nil")
	}
}

func TestFeatureConstants(t *testing.T) {
	// Test that feature constants are properly defined
	expectedFeatures := map[FeatureKey]string{
		FeatureTense:      "tense",
		FeatureAspect:     "aspect",
		FeatureMood:       "mood",
		FeaturePolarity:   "polarity",
		FeatureCase:       "case",
		FeatureVoice:      "voice",
		FeatureEvidential: "evidential",
		FeaturePerson:     "person",
		FeatureNumber:     "number",
		FeatureGender:     "gender",
	}

	for key, expected := range expectedFeatures {
		if string(key) != expected {
			t.Errorf("Expected FeatureKey %s, got %s", expected, string(key))
		}
	}
}
