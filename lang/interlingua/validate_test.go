package interlingua

import (
	"testing"
)

func TestValidateValidDocument(t *testing.T) {
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
	}

	err := Validate(doc)
	if err != nil {
		t.Errorf("Expected valid document, got error: %v", err)
	}
}

func TestValidateMissingSchemaVersion(t *testing.T) {
	doc := Document{
		LanguageID: "en",
		Entities:   []Entity{},
		Events:     []Event{},
	}

	err := Validate(doc)
	if err == nil {
		t.Error("Expected error for missing schema version")
	}

	validationErr, ok := err.(ValidationError)
	if !ok {
		t.Fatalf("Expected ValidationError, got %T", err)
	}

	if validationErr.Code != "missing.schema_version" {
		t.Errorf("Expected code 'missing.schema_version', got %s", validationErr.Code)
	}
}

func TestValidateUnsupportedSchemaVersion(t *testing.T) {
	doc := Document{
		SchemaVersion: "2.0",
		LanguageID:    "en",
		Entities:      []Entity{},
		Events:        []Event{},
	}

	err := Validate(doc)
	if err == nil {
		t.Error("Expected error for unsupported schema version")
	}

	validationErr, ok := err.(ValidationError)
	if !ok {
		t.Fatalf("Expected ValidationError, got %T", err)
	}

	if validationErr.Code != "unsupported.schema_version" {
		t.Errorf("Expected code 'unsupported.schema_version', got %s", validationErr.Code)
	}
}

func TestValidateMissingEntityID(t *testing.T) {
	doc := Document{
		SchemaVersion: SchemaVersion,
		LanguageID:    "en",
		Entities: []Entity{
			{
				ID:      "", // Missing ID
				Concept: "boy",
			},
		},
		Events: []Event{},
	}

	err := Validate(doc)
	if err == nil {
		t.Error("Expected error for missing entity ID")
	}

	validationErr, ok := err.(ValidationError)
	if !ok {
		t.Fatalf("Expected ValidationError, got %T", err)
	}

	if validationErr.Code != "missing.entity_id" {
		t.Errorf("Expected code 'missing.entity_id', got %s", validationErr.Code)
	}
}

func TestValidateMissingEntityConcept(t *testing.T) {
	doc := Document{
		SchemaVersion: SchemaVersion,
		LanguageID:    "en",
		Entities: []Entity{
			{
				ID:      "e1",
				Concept: "", // Missing concept
			},
		},
		Events: []Event{},
	}

	err := Validate(doc)
	if err == nil {
		t.Error("Expected error for missing entity concept")
	}

	validationErr, ok := err.(ValidationError)
	if !ok {
		t.Fatalf("Expected ValidationError, got %T", err)
	}

	if validationErr.Code != "missing.entity_concept" {
		t.Errorf("Expected code 'missing.entity_concept', got %s", validationErr.Code)
	}
}

func TestValidateMissingEventID(t *testing.T) {
	doc := Document{
		SchemaVersion: SchemaVersion,
		LanguageID:    "en",
		Entities:      []Entity{},
		Events: []Event{
			{
				ID:        "", // Missing ID
				Predicate: "give-01",
			},
		},
	}

	err := Validate(doc)
	if err == nil {
		t.Error("Expected error for missing event ID")
	}

	validationErr, ok := err.(ValidationError)
	if !ok {
		t.Fatalf("Expected ValidationError, got %T", err)
	}

	if validationErr.Code != "missing.event_id" {
		t.Errorf("Expected code 'missing.event_id', got %s", validationErr.Code)
	}
}

func TestValidateMissingEventPredicate(t *testing.T) {
	doc := Document{
		SchemaVersion: SchemaVersion,
		LanguageID:    "en",
		Entities:      []Entity{},
		Events: []Event{
			{
				ID:        "v1",
				Predicate: "", // Missing predicate
			},
		},
	}

	err := Validate(doc)
	if err == nil {
		t.Error("Expected error for missing event predicate")
	}

	validationErr, ok := err.(ValidationError)
	if !ok {
		t.Fatalf("Expected ValidationError, got %T", err)
	}

	if validationErr.Code != "missing.event_predicate" {
		t.Errorf("Expected code 'missing.event_predicate', got %s", validationErr.Code)
	}
}

func TestValidateInvalidRoleReference(t *testing.T) {
	doc := Document{
		SchemaVersion: SchemaVersion,
		LanguageID:    "en",
		Entities: []Entity{
			{
				ID:      "e1",
				Concept: "boy",
			},
		},
		Events: []Event{
			{
				ID:        "v1",
				Predicate: "give-01",
				Roles: map[Role]string{
					RoleAgent:   "e1",
					RolePatient: "e2", // References non-existent entity
				},
			},
		},
	}

	err := Validate(doc)
	if err == nil {
		t.Error("Expected error for invalid role reference")
	}

	validationErr, ok := err.(ValidationError)
	if !ok {
		t.Fatalf("Expected ValidationError, got %T", err)
	}

	if validationErr.Code != "invalid.role_reference" {
		t.Errorf("Expected code 'invalid.role_reference', got %s", validationErr.Code)
	}
}

func TestValidateInconsistentValency(t *testing.T) {
	doc := Document{
		SchemaVersion: SchemaVersion,
		LanguageID:    "en",
		Entities: []Entity{
			{
				ID:      "e1",
				Concept: "boy",
			},
			{
				ID:      "e2",
				Concept: "book",
			},
		},
		Events: []Event{
			{
				ID:        "v1",
				Predicate: "give-01",
				Roles: map[Role]string{
					RoleAgent:   "e1",
					RolePatient: "e2",
				},
				Valency: 3, // Declares 3 but only has 2 roles
			},
		},
	}

	err := Validate(doc)
	if err == nil {
		t.Error("Expected error for inconsistent valency")
	}

	validationErr, ok := err.(ValidationError)
	if !ok {
		t.Fatalf("Expected ValidationError, got %T", err)
	}

	if validationErr.Code != "inconsistent.valency" {
		t.Errorf("Expected code 'inconsistent.valency', got %s", validationErr.Code)
	}
}

func TestValidateRoleCompatibility(t *testing.T) {
	// Test valid role combinations
	err := ValidateRoleCompatibility("give-01", RoleRecipient)
	if err != nil {
		t.Errorf("Expected valid role combination, got error: %v", err)
	}

	err = ValidateRoleCompatibility("move-00", RoleSource)
	if err != nil {
		t.Errorf("Expected valid role combination, got error: %v", err)
	}

	// Test invalid role combinations
	err = ValidateRoleCompatibility("see-01", RoleRecipient)
	if err == nil {
		t.Error("Expected error for incompatible role")
	}

	validationErr, ok := err.(ValidationError)
	if !ok {
		t.Fatalf("Expected ValidationError, got %T", err)
	}

	if validationErr.Code != "incompatible.role" {
		t.Errorf("Expected code 'incompatible.role', got %s", validationErr.Code)
	}
}

func TestValidationErrorString(t *testing.T) {
	err := ValidationError{
		Code:    "test.code",
		Message: "test message",
	}

	expected := "validation error [test.code]: test message"
	if err.Error() != expected {
		t.Errorf("Expected error string '%s', got '%s'", expected, err.Error())
	}
}
