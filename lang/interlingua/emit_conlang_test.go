package interlingua

import (
	"encoding/json"
	"testing"
)

func TestNewDocument(t *testing.T) {
	doc := NewDocument("en")

	if doc.SchemaVersion != SchemaVersion {
		t.Errorf("Expected schema version %s, got %s", SchemaVersion, doc.SchemaVersion)
	}
	if doc.LanguageID != "en" {
		t.Errorf("Expected language ID 'en', got %s", doc.LanguageID)
	}
	if doc.Entities == nil {
		t.Error("Expected entities slice to be initialized")
	}
	if doc.Events == nil {
		t.Error("Expected events slice to be initialized")
	}
	if doc.Trace.TokenToNode == nil {
		t.Error("Expected TokenToNode map to be initialized")
	}
	if doc.Trace.Notes == nil {
		t.Error("Expected Notes slice to be initialized")
	}
}

func TestAddEntity(t *testing.T) {
	doc := NewDocument("en")

	// Test with pre-assigned ID
	entity1 := Entity{ID: "custom-id", Concept: "boy"}
	entityID1 := AddEntity(&doc, entity1)
	if entityID1 != "custom-id" {
		t.Errorf("Expected entity ID 'custom-id', got %s", entityID1)
	}
	if len(doc.Entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(doc.Entities))
	}

	// Test with auto-generated ID
	entity2 := Entity{Concept: "book"}
	entityID2 := AddEntity(&doc, entity2)
	if entityID2 != "e2" {
		t.Errorf("Expected entity ID 'e2', got %s", entityID2)
	}
	if len(doc.Entities) != 2 {
		t.Errorf("Expected 2 entities, got %d", len(doc.Entities))
	}
}

func TestAddEvent(t *testing.T) {
	doc := NewDocument("en")

	// Test with pre-assigned ID
	event1 := Event{ID: "custom-event", Predicate: "give-01"}
	eventID1 := AddEvent(&doc, event1)
	if eventID1 != "custom-event" {
		t.Errorf("Expected event ID 'custom-event', got %s", eventID1)
	}
	if len(doc.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(doc.Events))
	}

	// Test with auto-generated ID
	event2 := Event{Predicate: "move-00"}
	eventID2 := AddEvent(&doc, event2)
	if eventID2 != "v2" {
		t.Errorf("Expected event ID 'v2', got %s", eventID2)
	}
	if len(doc.Events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(doc.Events))
	}
}

func TestSetRole(t *testing.T) {
	doc := NewDocument("en")

	// Add entities first
	agent := Entity{Concept: "boy"}
	patient := Entity{Concept: "book"}
	agentID := AddEntity(&doc, agent)
	patientID := AddEntity(&doc, patient)

	// Add event
	event := Event{Predicate: "give-01"}
	AddEvent(&doc, event)

	// Set roles
	err := SetRole(&doc.Events[0], RoleAgent, agentID)
	if err != nil {
		t.Errorf("Failed to set agent role: %v", err)
	}

	err = SetRole(&doc.Events[0], RolePatient, patientID)
	if err != nil {
		t.Errorf("Failed to set patient role: %v", err)
	}

	// Verify roles
	if doc.Events[0].Roles[RoleAgent] != agentID {
		t.Errorf("Expected agent role to map to %s, got %s", agentID, doc.Events[0].Roles[RoleAgent])
	}
	if doc.Events[0].Roles[RolePatient] != patientID {
		t.Errorf("Expected patient role to map to %s, got %s", patientID, doc.Events[0].Roles[RolePatient])
	}
}

func TestSetRoleInvalid(t *testing.T) {
	doc := NewDocument("en")

	// Add entity
	entity := Entity{Concept: "boy"}
	entityID := AddEntity(&doc, entity)

	// Add event with incompatible predicate
	event := Event{Predicate: "see-01"} // see-01 is not ditransitive
	AddEvent(&doc, event)

	// Try to set recipient role (should fail)
	err := SetRole(&doc.Events[0], RoleRecipient, entityID)
	if err == nil {
		t.Error("Expected error when setting recipient role for non-ditransitive predicate")
	}
}

func TestAttachAdjunct(t *testing.T) {
	doc := NewDocument("en")

	// Add event
	event := Event{Predicate: "give-01"}
	AddEvent(&doc, event)

	// Add time adjunct
	timeEntity := Entity{Concept: "yesterday"}
	timeID := AddEntity(&doc, timeEntity)

	AttachAdjunct(&doc.Events[0], RoleTime, timeID)

	// Verify adjunct
	if doc.Events[0].Adjuncts[RoleTime] != timeID {
		t.Errorf("Expected time adjunct to map to %s, got %s", timeID, doc.Events[0].Adjuncts[RoleTime])
	}
}

func TestBuildFromTemplate(t *testing.T) {
	// Create entities
	agent := Entity{Concept: "boy", Feats: NominalFeatures{Number: "sg", Definiteness: "def"}}
	patient := Entity{Concept: "book", Feats: NominalFeatures{Number: "sg", Definiteness: "def"}}
	recipient := Entity{Concept: "girl", Feats: NominalFeatures{Number: "sg", Definiteness: "def"}}
	timeEntity := Entity{Concept: "yesterday"}

	// Create TAM
	tam := TAM{Tense: "past", Polarity: "neg"}

	// Build document
	args := map[Role]Entity{
		RoleAgent:     agent,
		RolePatient:   patient,
		RoleRecipient: recipient,
	}
	adjuncts := map[Role]Entity{
		RoleTime: timeEntity,
	}

	doc, err := BuildFromTemplate("elv-wood-sindarin", "decl.transitive.svo.neg.pst", "give-01", tam, args, adjuncts)
	if err != nil {
		t.Fatalf("Failed to build document: %v", err)
	}

	// Validate document
	if err := Validate(doc); err != nil {
		t.Fatalf("Document validation failed: %v", err)
	}

	// Check basic structure
	if doc.LanguageID != "elv-wood-sindarin" {
		t.Errorf("Expected language ID 'elv-wood-sindarin', got %s", doc.LanguageID)
	}
	if doc.Trace.ConstructionID != "decl.transitive.svo.neg.pst" {
		t.Errorf("Expected construction ID 'decl.transitive.svo.neg.pst', got %s", doc.Trace.ConstructionID)
	}

	// Check entities
	if len(doc.Entities) != 4 {
		t.Errorf("Expected 4 entities, got %d", len(doc.Entities))
	}

	// Check events
	if len(doc.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(doc.Events))
	}

	event := doc.Events[0]
	if event.Predicate != "give-01" {
		t.Errorf("Expected predicate 'give-01', got %s", event.Predicate)
	}
	if event.TAM.Tense != "past" {
		t.Errorf("Expected tense 'past', got %s", event.TAM.Tense)
	}
	if event.TAM.Polarity != "neg" {
		t.Errorf("Expected polarity 'neg', got %s", event.TAM.Polarity)
	}
	if event.Valency != 3 {
		t.Errorf("Expected valency 3, got %d", event.Valency)
	}

	// Check roles
	if len(event.Roles) != 3 {
		t.Errorf("Expected 3 roles, got %d", len(event.Roles))
	}

	// Check adjuncts
	if len(event.Adjuncts) != 1 {
		t.Errorf("Expected 1 adjunct, got %d", len(event.Adjuncts))
	}
}

func TestBuildSimpleTransitive(t *testing.T) {
	agent := Entity{Concept: "boy"}
	patient := Entity{Concept: "book"}
	tam := TAM{Tense: "pres", Polarity: "pos"}

	doc, err := BuildSimpleTransitive("en", "decl.transitive.svo.pos.pres", "give-01", tam, agent, patient)
	if err != nil {
		t.Fatalf("Failed to build simple transitive: %v", err)
	}

	if err := Validate(doc); err != nil {
		t.Fatalf("Document validation failed: %v", err)
	}

	if len(doc.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(doc.Events))
	}

	event := doc.Events[0]
	if len(event.Roles) != 2 {
		t.Errorf("Expected 2 roles, got %d", len(event.Roles))
	}
	if event.Valency != 2 {
		t.Errorf("Expected valency 2, got %d", event.Valency)
	}
}

func TestBuildDitransitive(t *testing.T) {
	agent := Entity{Concept: "boy"}
	patient := Entity{Concept: "book"}
	recipient := Entity{Concept: "girl"}
	tam := TAM{Tense: "past", Polarity: "pos"}

	doc, err := BuildDitransitive("en", "decl.ditransitive.svo.pos.pst", "give-01", tam, agent, patient, recipient)
	if err != nil {
		t.Fatalf("Failed to build ditransitive: %v", err)
	}

	if err := Validate(doc); err != nil {
		t.Fatalf("Document validation failed: %v", err)
	}

	if len(doc.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(doc.Events))
	}

	event := doc.Events[0]
	if len(event.Roles) != 3 {
		t.Errorf("Expected 3 roles, got %d", len(event.Roles))
	}
	if event.Valency != 3 {
		t.Errorf("Expected valency 3, got %d", event.Valency)
	}
}

func TestAddTimeAdjunct(t *testing.T) {
	doc := NewDocument("en")

	// Add event
	event := Event{Predicate: "give-01"}
	AddEvent(&doc, event)

	// Add time adjunct
	timeEntity := Entity{Concept: "yesterday"}
	err := AddTimeAdjunct(&doc, 0, timeEntity)
	if err != nil {
		t.Errorf("Failed to add time adjunct: %v", err)
	}

	// Verify adjunct was added
	if len(doc.Events[0].Adjuncts) != 1 {
		t.Errorf("Expected 1 adjunct, got %d", len(doc.Events[0].Adjuncts))
	}

	// Verify token mapping
	if len(doc.Trace.TokenToNode) != 1 {
		t.Errorf("Expected 1 token mapping, got %d", len(doc.Trace.TokenToNode))
	}
}

func TestAddLocationAdjunct(t *testing.T) {
	doc := NewDocument("en")

	// Add event
	event := Event{Predicate: "give-01"}
	AddEvent(&doc, event)

	// Add location adjunct
	locationEntity := Entity{Concept: "house"}
	err := AddLocationAdjunct(&doc, 0, locationEntity)
	if err != nil {
		t.Errorf("Failed to add location adjunct: %v", err)
	}

	// Verify adjunct was added
	if len(doc.Events[0].Adjuncts) != 1 {
		t.Errorf("Expected 1 adjunct, got %d", len(doc.Events[0].Adjuncts))
	}

	// Verify token mapping
	if len(doc.Trace.TokenToNode) != 1 {
		t.Errorf("Expected 1 token mapping, got %d", len(doc.Trace.TokenToNode))
	}
}

func TestMapTokensToNodes(t *testing.T) {
	doc := NewDocument("en")

	tokens := []string{"neg", "give", "boy", "book", "girl", "yesterday"}
	nodeIDs := []string{"v1", "e1", "e2", "e3", "t1"}

	MapTokensToNodes(&doc, tokens, nodeIDs)

	// Should have 5 mappings (shorter of the two slices)
	if len(doc.Trace.TokenToNode) != 5 {
		t.Errorf("Expected 5 token mappings, got %d", len(doc.Trace.TokenToNode))
	}

	// Check specific mappings
	if doc.Trace.TokenToNode[0] != "v1" {
		t.Errorf("Expected token 0 to map to 'v1', got %s", doc.Trace.TokenToNode[0])
	}
	if doc.Trace.TokenToNode[4] != "t1" {
		t.Errorf("Expected token 4 to map to 't1', got %s", doc.Trace.TokenToNode[4])
	}

	// Should have a warning note about mismatch
	warnings := doc.Trace.GetNotesByCode("token.mapping.mismatch")
	if len(warnings) != 1 {
		t.Errorf("Expected 1 warning note, got %d", len(warnings))
	}
}

func TestSetSourceLanguage(t *testing.T) {
	doc := NewDocument("en")

	SetSourceLanguage(&doc, "elv-wood-sindarin")

	if doc.Trace.SourceLanguageID != "elv-wood-sindarin" {
		t.Errorf("Expected source language 'elv-wood-sindarin', got %s", doc.Trace.SourceLanguageID)
	}
}

func TestAddConstructionNote(t *testing.T) {
	doc := NewDocument("en")

	AddConstructionNote(&doc, "Selected SVO transitive construction")

	notes := doc.Trace.GetNotesByCode("construction.info")
	if len(notes) != 1 {
		t.Errorf("Expected 1 construction note, got %d", len(notes))
	}

	if notes[0].Message != "Selected SVO transitive construction" {
		t.Errorf("Expected note message 'Selected SVO transitive construction', got %s", notes[0].Message)
	}
}

func TestAddFeatureNote(t *testing.T) {
	doc := NewDocument("en")

	AddFeatureNote(&doc, "warn", "defaulted", "Number defaulted to singular")

	notes := doc.Trace.GetNotesByCode("feature.warn.defaulted")
	if len(notes) != 1 {
		t.Errorf("Expected 1 feature note, got %d", len(notes))
	}

	if notes[0].Severity != "warn" {
		t.Errorf("Expected severity 'warn', got %s", notes[0].Severity)
	}
	if notes[0].Message != "Number defaulted to singular" {
		t.Errorf("Expected note message 'Number defaulted to singular', got %s", notes[0].Message)
	}
}

func TestBuildFromTemplateREADMEExample(t *testing.T) {
	// This test recreates the README example exactly
	// Surface (conlang, glossed): `neg-give.pst boy-nom book-acc girl-dat yesterday`

	// Create entities as specified in README
	agent := Entity{
		Concept: "boy",
		Feats:   NominalFeatures{Number: "sg", Definiteness: "def"},
	}
	patient := Entity{
		Concept: "book",
		Feats:   NominalFeatures{Number: "sg", Definiteness: "def"},
	}
	recipient := Entity{
		Concept: "girl",
		Feats:   NominalFeatures{Number: "sg", Definiteness: "def"},
	}
	timeEntity := Entity{
		Concept: "yesterday",
	}

	// Create TAM as specified
	tam := TAM{
		Tense:    "past",
		Polarity: "neg",
	}

	// Build document using template
	args := map[Role]Entity{
		RoleAgent:     agent,
		RolePatient:   patient,
		RoleRecipient: recipient,
	}
	adjuncts := map[Role]Entity{
		RoleTime: timeEntity,
	}

	doc, err := BuildFromTemplate("elv-wood-sindarin", "decl.transitive.svo.neg.pst", "give-01", tam, args, adjuncts)
	if err != nil {
		t.Fatalf("Failed to build README example: %v", err)
	}

	// Validate the document
	if err := Validate(doc); err != nil {
		t.Fatalf("README example validation failed: %v", err)
	}

	// Verify it matches README specification
	if doc.LanguageID != "elv-wood-sindarin" {
		t.Errorf("Expected language ID 'elv-wood-sindarin', got %s", doc.LanguageID)
	}
	if doc.Trace.ConstructionID != "decl.transitive.svo.neg.pst" {
		t.Errorf("Expected construction ID 'decl.transitive.svo.neg.pst', got %s", doc.Trace.ConstructionID)
	}

	// Verify entities
	if len(doc.Entities) != 4 {
		t.Errorf("Expected 4 entities, got %d", len(doc.Entities))
	}

	// Verify event
	if len(doc.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(doc.Events))
	}

	event := doc.Events[0]
	if event.Predicate != "give-01" {
		t.Errorf("Expected predicate 'give-01', got %s", event.Predicate)
	}
	if event.TAM.Tense != "past" {
		t.Errorf("Expected tense 'past', got %s", event.TAM.Tense)
	}
	if event.TAM.Polarity != "neg" {
		t.Errorf("Expected polarity 'neg', got %s", event.TAM.Polarity)
	}
	if event.Valency != 3 {
		t.Errorf("Expected valency 3, got %d", event.Valency)
	}

	// Verify roles
	if len(event.Roles) != 3 {
		t.Errorf("Expected 3 roles, got %d", len(event.Roles))
	}

	// Verify adjuncts
	if len(event.Adjuncts) != 1 {
		t.Errorf("Expected 1 adjunct, got %d", len(event.Adjuncts))
	}

	// Test JSON marshaling (as specified in README)
	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("Failed to marshal README example to JSON: %v", err)
	}

	// Verify JSON contains expected content
	jsonStr := string(data)
	if !contains(jsonStr, "give-01") {
		t.Error("JSON should contain predicate 'give-01'")
	}
	if !contains(jsonStr, "past") {
		t.Error("JSON should contain tense 'past'")
	}
	if !contains(jsonStr, "neg") {
		t.Error("JSON should contain polarity 'neg'")
	}
	if !contains(jsonStr, "boy") {
		t.Error("JSON should contain entity 'boy'")
	}
	if !contains(jsonStr, "book") {
		t.Error("JSON should contain entity 'book'")
	}
	if !contains(jsonStr, "girl") {
		t.Error("JSON should contain entity 'girl'")
	}
	if !contains(jsonStr, "yesterday") {
		t.Error("JSON should contain time adjunct 'yesterday'")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || contains(s[1:], substr)))
}
