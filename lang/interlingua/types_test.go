package interlingua

import (
	"encoding/json"
	"testing"
)

func TestDocumentJSONMarshaling(t *testing.T) {
	doc := Document{
		SchemaVersion: SchemaVersion,
		LanguageID:    "en",
		Entities: []Entity{
			{
				ID:      "e1",
				Concept: "boy",
				Feats: NominalFeatures{
					Number:       "sg",
					Definiteness: "def",
				},
			},
			{
				ID:      "e2",
				Concept: "book",
				Feats: NominalFeatures{
					Number:       "sg",
					Definiteness: "def",
				},
			},
		},
		Events: []Event{
			{
				ID:        "v1",
				Predicate: "give-01",
				TAM: TAM{
					Tense:    "past",
					Polarity: "pos",
				},
				Roles: map[Role]string{
					RoleAgent:   "e1",
					RolePatient: "e2",
				},
				Valency: 2,
			},
		},
		Trace: Trace{
			SourceLanguageID: "en",
			ConstructionID:   "decl.transitive.svo.pos.pst",
			TokenToNode: map[int]string{
				0: "v1",
				1: "e1",
				2: "e2",
			},
		},
	}

	// Test marshaling
	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("Failed to marshal document: %v", err)
	}

	// Test unmarshaling
	var unmarshaled Document
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal document: %v", err)
	}

	// Verify key fields
	if unmarshaled.SchemaVersion != SchemaVersion {
		t.Errorf("Expected schema version %s, got %s", SchemaVersion, unmarshaled.SchemaVersion)
	}
	if unmarshaled.LanguageID != "en" {
		t.Errorf("Expected language ID 'en', got %s", unmarshaled.LanguageID)
	}
	if len(unmarshaled.Entities) != 2 {
		t.Errorf("Expected 2 entities, got %d", len(unmarshaled.Entities))
	}
	if len(unmarshaled.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(unmarshaled.Events))
	}
	if unmarshaled.Events[0].Predicate != "give-01" {
		t.Errorf("Expected predicate 'give-01', got %s", unmarshaled.Events[0].Predicate)
	}
}

func TestTAMJSONMarshaling(t *testing.T) {
	tam := TAM{
		Tense:      "past",
		Aspect:     "perf",
		Mood:       "ind",
		Polarity:   "pos",
		Evidential: "dir",
		Voice:      "act",
		Deixis:     "proximal",
	}

	data, err := json.Marshal(tam)
	if err != nil {
		t.Fatalf("Failed to marshal TAM: %v", err)
	}

	var unmarshaled TAM
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal TAM: %v", err)
	}

	if unmarshaled.Tense != "past" {
		t.Errorf("Expected tense 'past', got %s", unmarshaled.Tense)
	}
	if unmarshaled.Aspect != "perf" {
		t.Errorf("Expected aspect 'perf', got %s", unmarshaled.Aspect)
	}
}

func TestNominalFeaturesJSONMarshaling(t *testing.T) {
	feats := NominalFeatures{
		Person:       3,
		Number:       "sg",
		Gender:       "masc",
		Animacy:      "anim",
		Definiteness: "def",
		Case:         "nom",
		Classifier:   "human",
		Deixis:       "proximal",
	}

	data, err := json.Marshal(feats)
	if err != nil {
		t.Fatalf("Failed to marshal NominalFeatures: %v", err)
	}

	var unmarshaled NominalFeatures
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal NominalFeatures: %v", err)
	}

	if unmarshaled.Person != 3 {
		t.Errorf("Expected person 3, got %d", unmarshaled.Person)
	}
	if unmarshaled.Number != "sg" {
		t.Errorf("Expected number 'sg', got %s", unmarshaled.Number)
	}
}

func TestRoleConstants(t *testing.T) {
	expectedRoles := map[Role]string{
		RoleAgent:       "agent",
		RolePatient:     "patient",
		RoleRecipient:   "recipient",
		RoleExperiencer: "experiencer",
		RoleStimulus:    "stimulus",
		RoleInstrument:  "instrument",
		RoleLocation:    "location",
		RoleTime:        "time",
		RoleManner:      "manner",
		RoleCause:       "cause",
		RoleBeneficiary: "beneficiary",
		RoleSource:      "source",
		RoleGoal:        "goal",
	}

	for role, expected := range expectedRoles {
		if string(role) != expected {
			t.Errorf("Expected role %s, got %s", expected, string(role))
		}
	}
}

func TestEntityCreation(t *testing.T) {
	entity := Entity{
		ID:      "test-entity",
		Concept: "test-concept",
		Feats: NominalFeatures{
			Number: "sg",
		},
	}

	if entity.ID != "test-entity" {
		t.Errorf("Expected ID 'test-entity', got %s", entity.ID)
	}
	if entity.Concept != "test-concept" {
		t.Errorf("Expected concept 'test-concept', got %s", entity.Concept)
	}
	if entity.Feats.Number != "sg" {
		t.Errorf("Expected number 'sg', got %s", entity.Feats.Number)
	}
}

func TestEventCreation(t *testing.T) {
	event := Event{
		ID:        "test-event",
		Predicate: "test-predicate",
		TAM: TAM{
			Tense: "pres",
		},
		Roles: map[Role]string{
			RoleAgent: "e1",
		},
		Valency: 1,
	}

	if event.ID != "test-event" {
		t.Errorf("Expected ID 'test-event', got %s", event.ID)
	}
	if event.Predicate != "test-predicate" {
		t.Errorf("Expected predicate 'test-predicate', got %s", event.Predicate)
	}
	if event.TAM.Tense != "pres" {
		t.Errorf("Expected tense 'pres', got %s", event.TAM.Tense)
	}
	if len(event.Roles) != 1 {
		t.Errorf("Expected 1 role, got %d", len(event.Roles))
	}
	if event.Valency != 1 {
		t.Errorf("Expected valency 1, got %d", event.Valency)
	}
}
